/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package extensions

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/go-logr/logr"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// PreReconciliationChecker is implemented by resources that want to do extra checks before proceeding with
// a full ARM reconcile. This extension is invoked before sending any requests to Azure, giving resources
// the ability to block reconciliation until prerequisites are met.
// Implement this extension when:
// - Parent/owner resources must reach certain states before child creation
// - External dependencies must be satisfied before reconciling
// - Reconciliation would fail due to known prerequisites not being met
type PreReconciliationChecker interface {
	// PreReconcileCheck does a pre-reconcile check to see if the resource is in a state that can be reconciled.
	// ARM resources should implement this to avoid reconciliation attempts that cannot possibly succeed.
	// Returns ProceedWithReconcile if the reconciliation should go ahead.
	// Returns BlockReconcile and a human-readable reason if the reconciliation should be skipped.
	// ctx is the current operation context.
	// obj is the resource about to be reconciled. The resource's status will be freshly updated.
	// resourceResolver helps resolve resource references.
	// armClient allows access to ARM for any required queries.
	// log is the logger for the current operation.
	// next is the next (nested) implementation to call.
	PreReconcileCheck(
		ctx context.Context,
		obj genruntime.MetaObject,
		resourceResolver *resolver.Resolver,
		armClient *genericarmclient.GenericClient,
		log logr.Logger,
		next PreReconcileCheckFunc,
	) (PreReconcileCheckResult, error)
}

type PreReconcileCheckFunc func(
	ctx context.Context,
	obj genruntime.MetaObject,
	resourceResolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (PreReconcileCheckResult, error)

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
		return preReconciliationCheckerAlways, false
	}

	return func(
		ctx context.Context,
		obj genruntime.MetaObject,
		resourceResolver *resolver.Resolver,
		armClient *genericarmclient.GenericClient,
		log logr.Logger,
	) (PreReconcileCheckResult, error) {
		log.V(Status).Info("Extension pre-reconcile check running")

		result, err := impl.PreReconcileCheck(ctx, obj, resourceResolver, armClient, log, preReconciliationCheckerAlways)
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

// preReconciliationCheckerAlways is a PreReconciliationChecker that always indicates a reconciliation is required.
// We have this here so we can set up a chain, even if it's only one link long.
// When we start doing proper comparisons between Spec and Status, we'll have an actual chain of checkers.
func preReconciliationCheckerAlways(
	_ context.Context,
	_ genruntime.MetaObject,
	_ *resolver.Resolver,
	_ *genericarmclient.GenericClient,
	_ logr.Logger,
) (PreReconcileCheckResult, error) {
	return ProceedWithReconcile(), nil
}
