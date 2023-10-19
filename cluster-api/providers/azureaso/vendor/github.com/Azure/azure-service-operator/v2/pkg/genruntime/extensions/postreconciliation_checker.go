/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package extensions

import (
	"context"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
)

// PostReconciliationChecker is implemented by resources that want to do extra status checks after
// a full ARM reconcile.
type PostReconciliationChecker interface {
	// PostReconcileCheck does a post-reconcile check to see if the resource is in a state to set 'Ready' condition.
	// ARM resources should implement this if they need to defer the Ready condition until later.
	// Returns PostReconcileCheckResultSuccess if the reconciliation is successful.
	// Returns PostReconcileCheckResultFailure and a human-readable reason if the reconciliation should put a condition on resource.
	// ctx is the current operation context.
	// obj is the resource about to be reconciled. The resource's State will be freshly updated.
	// owner is the parent resource of obj. This can be nil in some cases like `ResourceGroups` and `Alias`.
	// kubeClient allows access to the cluster for any required queries.
	// armClient allows access to ARM for any required queries.
	// log is the logger for the current operation.
	PostReconcileCheck(
		ctx context.Context,
		obj genruntime.MetaObject,
		owner genruntime.MetaObject,
		resourceResolver *resolver.Resolver,
		armClient *genericarmclient.GenericClient,
		log logr.Logger,
		next PostReconcileCheckFunc,
	) (PostReconcileCheckResult, error)
}

type PostReconcileCheckFunc func(
	ctx context.Context,
	obj genruntime.MetaObject,
	owner genruntime.MetaObject,
	resourceResolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (PostReconcileCheckResult, error)

type PostReconcileCheckResult struct {
	action   postReconcileCheckResultType
	severity conditions.ConditionSeverity
	reason   conditions.Reason
	message  string
}

// PostReconcileCheckResultSuccess indicates that a resource is ready after reconciliation by returning a PostReconcileCheckResult
// with action `Success`.
func PostReconcileCheckResultSuccess() PostReconcileCheckResult {
	return PostReconcileCheckResult{
		action: postReconcileCheckResultTypeSuccess,
	}
}

// PostReconcileCheckResultFailure indicates post reconciliation check of a resource is currently failed by returning a PostReconcileCheckResult
// with action `Failure`.
// reason is an explanatory reason to show to the user via a warning condition on the resource.
func PostReconcileCheckResultFailure(reason string) PostReconcileCheckResult {
	return PostReconcileCheckResult{
		action:   postReconcileCheckResultTypeFailure,
		severity: conditions.ConditionSeverityWarning,
		reason:   conditions.ReasonPostReconcileFailure,
		message:  reason,
	}
}

func (r PostReconcileCheckResult) ReconciliationFailed() bool {
	return r.action == postReconcileCheckResultTypeFailure
}

func (r PostReconcileCheckResult) ReconciliationSucceeded() bool {
	return r.action == postReconcileCheckResultTypeSuccess
}

func (r PostReconcileCheckResult) Message() string {
	return r.message
}

// CreateConditionError returns an error that can be used to set a condition on the resource.
func (r PostReconcileCheckResult) CreateConditionError() error {
	return conditions.NewReadyConditionImpactingError(
		errors.New(r.message),
		r.severity,
		r.reason)
}

// postReconcileCheckResultType is the type of result returned by PreReconcileCheck.
type postReconcileCheckResultType string

const (
	postReconcileCheckResultTypeFailure postReconcileCheckResultType = "Failure"
	postReconcileCheckResultTypeSuccess postReconcileCheckResultType = "Success"
)

// CreatePostReconciliationChecker creates a checker that can be used to check if we want to customise the condition on the resource after reconciliation.
// If the resource in question has not implemented the PostReconciliationChecker interface, the provided default checker
// is returned directly.
// We also return a bool indicating whether the resource extension implements the PostReconciliationChecker interface.
// host is a resource extension that may implement the PostReconciliationChecker interface.
func CreatePostReconciliationChecker(
	host genruntime.ResourceExtension,
) (PostReconcileCheckFunc, bool) {
	impl, ok := host.(PostReconciliationChecker)
	if !ok {
		return nil, false
	}

	return func(
		ctx context.Context,
		obj genruntime.MetaObject,
		owner genruntime.MetaObject,
		resourceResolver *resolver.Resolver,
		armClient *genericarmclient.GenericClient,
		log logr.Logger,
	) (PostReconcileCheckResult, error) {
		log.V(Status).Info("Extension post-reconcile check running")

		result, err := impl.PostReconcileCheck(ctx, obj, owner, resourceResolver, armClient, log, alwaysSucceed)
		if err != nil {
			log.V(Status).Info(
				"Extension post-reconcile check failed",
				"Error", err.Error())

			// We choose to skip here so that things are definitely broken and the user will notice
			// If we defaulted to always reconciling, the user might not notice that something is wrong
			return PostReconcileCheckResultFailure("Extension PostReconcileCheck failed"), err
		}

		log.V(Status).Info("Extension post-reconcile check succeeded", "Result", result)

		return result, nil
	}, true
}

// alwaysSucceed is a PostReconciliationChecker that always indicates a reconciliation is successful.
// We have this here, so we can set up a chain, even if it's only one link long.
func alwaysSucceed(
	_ context.Context,
	_ genruntime.MetaObject,
	_ genruntime.MetaObject,
	_ *resolver.Resolver,
	_ *genericarmclient.GenericClient,
	_ logr.Logger,
) (PostReconcileCheckResult, error) {
	return PostReconcileCheckResultSuccess(), nil
}
