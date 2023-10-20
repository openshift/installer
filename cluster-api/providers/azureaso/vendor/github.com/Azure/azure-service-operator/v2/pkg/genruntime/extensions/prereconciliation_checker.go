/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package extensions

import (
	"context"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

// PreReconciliationChecker is implemented by resources that want to do extra checks before proceeding with
// a full ARM reconcile.
type PreReconciliationChecker interface {
	// PreReconcileCheck does a pre-reconcile check to see if the resource is in a state that can be reconciled.
	// ARM resources should implement this to avoid reconciliation attempts that cannot possibly succeed.
	// Returns ProceedWithReconcile if the reconciliation should go ahead.
	// Returns BlockReconcile and a human-readable reason if the reconciliation should be skipped.
	// ctx is the current operation context.
	// obj is the resource about to be reconciled. The resource's State will be freshly updated.
	// kubeClient allows access to the cluster for any required queries.
	// armClient allows access to ARM for any required queries.
	// log is the logger for the current operation.
	// next is the next (nested) implementation to call.
	PreReconcileCheck(
		ctx context.Context,
		obj genruntime.MetaObject,
		owner genruntime.MetaObject,
		resourceResolver *resolver.Resolver,
		armClient *genericarmclient.GenericClient,
		log logr.Logger,
		next PreReconcileCheckFunc,
	) (PreReconcileCheckResult, error)
}

type PreReconcileCheckFunc func(
	ctx context.Context,
	obj genruntime.MetaObject,
	owner genruntime.MetaObject,
	resourceResolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (PreReconcileCheckResult, error)

type PreReconcileCheckResult struct {
	action   preReconcileCheckResultType
	severity conditions.ConditionSeverity
	reason   conditions.Reason
	message  string
}

// ProceedWithReconcile indicates that a resource is ready for reconciliation by returning a PreReconcileCheckResult
// with action `Proceed`.
func ProceedWithReconcile() PreReconcileCheckResult {
	return PreReconcileCheckResult{
		action: preReconcileCheckResultTypeProceed,
	}
}

// BlockReconcile indicates reconciliation of a resource is currently blocked by returning a PreReconcileCheckResult
// with action `Block`.
// reason is an explanatory reason to show to the user via a warning condition on the resource.
func BlockReconcile(reason string) PreReconcileCheckResult {
	return PreReconcileCheckResult{
		action:   preReconcileCheckResultTypeBlock,
		severity: conditions.ConditionSeverityWarning,
		reason:   conditions.ReasonReconcileBlocked,
		message:  reason,
	}
}

// PostponeReconcile indicates reconciliation of a resource is not currently required by returning a
// PreReconcileCheckResult with action `Postpone`.
func PostponeReconcile() PreReconcileCheckResult {
	return PreReconcileCheckResult{
		action:   preReconcileCheckResultTypePostpone,
		severity: conditions.ConditionSeverityInfo,
		reason:   conditions.ReasonReconcilePostponed,
	}
}

func (r PreReconcileCheckResult) PostponeReconciliation() bool {
	return r.action == preReconcileCheckResultTypePostpone
}

func (r PreReconcileCheckResult) BlockReconciliation() bool {
	return r.action == preReconcileCheckResultTypeBlock
}

func (r PreReconcileCheckResult) Message() string {
	return r.message
}

// CreateConditionError returns an error that can be used to set a condition on the resource.
func (r PreReconcileCheckResult) CreateConditionError() error {
	return conditions.NewReadyConditionImpactingError(
		errors.New(r.message),
		r.severity,
		r.reason)
}

// PreReconcileCheckResultType is the type of result returned by PreReconcileCheck.
type preReconcileCheckResultType string

const (
	preReconcileCheckResultTypeBlock    preReconcileCheckResultType = "Block"
	preReconcileCheckResultTypeProceed  preReconcileCheckResultType = "Proceed"
	preReconcileCheckResultTypePostpone preReconcileCheckResultType = "Postpone"
)

// CreatePreReconciliationChecker creates a checker that can be used to check if a resource is ready for reconciliation.
// If the resource in question has not implemented the PreReconciliationChecker interface, the provided default checker
// is returned directly.
// We also return a bool indicating whether the resource extension implements the PreReconciliationChecker interface.
// host is a resource extension that may implement the PreReconciliationChecker interface.
func CreatePreReconciliationChecker(
	host genruntime.ResourceExtension,
) (PreReconcileCheckFunc, bool) {
	impl, ok := host.(PreReconciliationChecker)
	if !ok {
		return alwaysReconcile, false
	}

	return func(
		ctx context.Context,
		obj genruntime.MetaObject,
		owner genruntime.MetaObject,
		resourceResolver *resolver.Resolver,
		armClient *genericarmclient.GenericClient,
		log logr.Logger,
	) (PreReconcileCheckResult, error) {
		log.V(Status).Info("Extension pre-reconcile check running")

		result, err := impl.PreReconcileCheck(ctx, obj, owner, resourceResolver, armClient, log, alwaysReconcile)
		if err != nil {
			log.V(Status).Info(
				"Extension pre-reconcile check failed",
				"Error", err.Error())

			// We choose to skip here so that things are definitely broken and the user will notice
			// If we defaulted to always reconciling, the user might not notice that something is wrong
			return BlockReconcile("Extension PreReconcileCheck failed"), err
		}

		log.V(Status).Info(
			"Extension pre-reconcile check succeeded",
			"Result", result)

		return result, nil
	}, true
}

// alwaysReconcile is a PreReconciliationChecker that always indicates a reconciliation is required.
// We have this here so we can set up a chain, even if it's only one link long.
// When we start doing proper comparisons between Spec and Status, we'll have an actual chain of checkers.
func alwaysReconcile(
	_ context.Context,
	_ genruntime.MetaObject,
	_ genruntime.MetaObject,
	_ *resolver.Resolver,
	_ *genericarmclient.GenericClient,
	_ logr.Logger,
) (PreReconcileCheckResult, error) {
	return ProceedWithReconcile(), nil
}
