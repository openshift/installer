/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package interval

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/internal/util/randextensions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
)

// Calculator calculates an interval
type Calculator interface {
	NextInterval(req ctrl.Request, result ctrl.Result, err error) (ctrl.Result, error)
}

type CalculatorParameters struct {
	Rand                 *rand.Rand
	ErrorBaseDelay       time.Duration
	ErrorMaxFastDelay    time.Duration
	ErrorMaxSlowDelay    time.Duration
	SyncPeriod           *time.Duration
	RequeueDelayOverride time.Duration
}

// NewCalculator creates a new Calculator.
func NewCalculator(params CalculatorParameters) Calculator {
	return &calculator{
		failures:             make(map[ctrl.Request]int),
		errorBaseDelay:       params.ErrorBaseDelay,
		errorMaxSlowDelay:    params.ErrorMaxSlowDelay,
		errorMaxFastDelay:    params.ErrorMaxFastDelay,
		syncPeriod:           params.SyncPeriod,
		requeueDelayOverride: params.RequeueDelayOverride,
		rand:                 params.Rand,
	}
}

type calculator struct {
	failuresLock sync.Mutex
	failures     map[ctrl.Request]int

	// Used only if there is an error
	errorBaseDelay    time.Duration
	errorMaxSlowDelay time.Duration
	errorMaxFastDelay time.Duration

	// Used only if there is not an error
	syncPeriod           *time.Duration
	requeueDelayOverride time.Duration
	rand                 *rand.Rand
}

var _ Calculator = &calculator{}

// NextInterval calculates the next interval for a given request, result, and error.
// Remember: There is also a controller-runtime RateLimiter that also can determine intervals. This implementation
// takes ownership of specific scenarios while leaving the rest to the standard RateLimiter.
// The scenarios that this implementation targets are:
// 1. Errors of type ReadyConditionImpactingError. These have a RetryClassification on them which is honored here.
// 2. If syncPeriod is set, success results that would normally be terminal are instead configured to try again in syncPeriod.
// 3. If requeueDelayOverride is set, all happy-path requests have requeueDelayOverride set.
// The scenarios that this handler doesn't target:
//  1. Any error other than ReadyConditionImpacting error.
//  2. Happy-path requests when requeueDelayOverride is not set. These are scenarios where the operator is working
//     as expected and we're just doing something like polling an async operation.
func (i *calculator) NextInterval(req ctrl.Request, result ctrl.Result, err error) (ctrl.Result, error) {
	i.failuresLock.Lock()
	defer i.failuresLock.Unlock()

	if err != nil {
		return i.failureResult(req, err)
	}

	// Happy path
	if (result == ctrl.Result{}) {
		// If result is a success, ensure that we requeue for monitoring state in Azure
		result = i.makeSuccessResult()
	}

	delete(i.failures, req) // On reconcile without an error, forget any previous failures

	hasRequeueDelayOverride := i.requeueDelayOverride != time.Duration(0)
	isRequeueing := result.Requeue || result.RequeueAfter > time.Duration(0)
	if hasRequeueDelayOverride && isRequeueing {
		result.RequeueAfter = i.requeueDelayOverride
		result.Requeue = true
	}
	return result, nil
}

func (i *calculator) failureResult(req ctrl.Request, err error) (ctrl.Result, error) {
	exp := i.failures[req]
	i.failures[req] = i.failures[req] + 1

	readyErr, ok := conditions.AsReadyConditionImpactingError(err)
	if !ok {
		// NotFound is a superfluous error as per https://github.com/kubernetes-sigs/controller-runtime/issues/377
		// The correct handling is just to ignore it and we will get an event shortly with the updated version to patch
		// We do NOT ignore conflict here because it's hard to tell if it's coming from an attempt to update a non-existing resource
		// (see https://github.com/kubernetes/kubernetes/issues/89985), or if it's from an attempt to update a resource which
		// was updated by a user. If we ignore the user-update case, we MIGHT get another event since they changed the resource,
		// but since we don't trigger updates on all changes (some annotations are ignored) we also MIGHT NOT get a fresh event
		// and get stuck. The solution is to let the GET at the top of the controller check for the not-found case and requeue
		// on everything else.
		err = kubeclient.IgnoreNotFound(err)
		if err == nil {
			// Since we're ignoring this error and counting it as a success, stop tracking the req
			delete(i.failures, req)
		}
		return ctrl.Result{}, err
	}

	// Now we have a readyErr
	if readyErr.Severity == conditions.ConditionSeverityError {
		// Severity error is fatal, return fast and block requeue
		// Since this is fatal, stop tracking the req
		delete(i.failures, req)
		return ctrl.Result{}, nil
	} else if readyErr.Severity == conditions.ConditionSeverityWarning {
		switch readyErr.RetryClassification {
		case conditions.RetrySlow:
			delay := i.calculateExponentialDelay(i.errorBaseDelay, exp, i.errorMaxSlowDelay)
			return ctrl.Result{RequeueAfter: delay}, nil
		case conditions.RetryFast:
			delay := i.calculateExponentialDelay(i.errorBaseDelay, exp, i.errorMaxFastDelay)
			return ctrl.Result{RequeueAfter: delay}, nil
		case conditions.RetryNone:
			// This shouldn't happen, return an error
			return ctrl.Result{}, errors.New("didn't expect RetryNone classification for error")
		default:
			// This shouldn't happen, return an error
			return ctrl.Result{}, errors.Errorf("unknown RetryClassification %q", readyErr.RetryClassification)
		}
	}

	// This shouldn't happen, return an error
	return ctrl.Result{}, errors.Errorf("Error with severity %q is unexpected", readyErr.Severity)
}

func (i *calculator) makeSuccessResult() ctrl.Result {
	result := ctrl.Result{}
	// This has a RequeueAfter because we want to force a re-sync at some point in the future in order to catch
	// potential drift from the state in Azure. Note that we cannot use mgr.Options.SyncPeriod for this because we filter
	// our events by predicate.GenerationChangedPredicate and the generation will not have changed.
	if i.syncPeriod != nil {
		result.RequeueAfter = randextensions.Jitter(i.rand, *i.syncPeriod, 0.25)
	}

	return result
}

func (i *calculator) calculateExponentialDelay(base time.Duration, exp int, max time.Duration) time.Duration {
	// The backoff is capped such that 'calculated' value never overflows.
	backoff := float64(base.Nanoseconds()) * math.Pow(2, float64(exp))
	if backoff > math.MaxInt64 {
		return max
	}

	calculated := time.Duration(backoff)
	if calculated > max {
		return max
	}

	return calculated
}
