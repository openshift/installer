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

// PreReconciliationOwnerChecker is implemented by resources that want to check their owner's state before
// proceeding with a full ARM reconcile. This is a specialized variant of PreReconciliationChecker that only
// checks the owner, avoiding any GET operations on the resource itself.
// Implement this extension when:
// - The owner's state can block all access to the resource, including GET operations
// - You need to avoid attempting operations that will fail due to owner state
// - The resource cannot be accessed when its owner is in certain states (e.g., powered off, updating)
type PreReconciliationOwnerChecker interface {
	// PreReconcileOwnerCheck does a pre-reconcile check to see if the owner of a resource is in a state that permits
	// the resource to be reconciled. For a limited number of resources, the state of their owner can block all access
	// to the resource, including GETs. One example is a Kusto Cluster, where you can't even GET the database if the
	// cluster is powered off.
	// Prefer to implement PreReconciliationChecker unless you specifically need to avoid GETs on the resource itself.
	// ARM resources should implement this to avoid reconciliation attempts that cannot possibly succeed.
	// Returns ProceedWithReconcile if the reconciliation should go ahead.
	// Returns BlockReconcile and a human-readable reason if the reconciliation should be skipped.
	// ctx is the current operation context.
	// owner is the owner of the resource about to be reconciled. The owner's status will be freshly updated. May be nil
	// if the resource has no owner, or if it has been referenced via ARMID directly.
	// resourceResolver helps resolve resource references.
	// armClient allows access to ARM for any required queries.
	// log is the logger for the current operation.
	// next is the next (nested) implementation to call.
	PreReconcileOwnerCheck(
		ctx context.Context,
		owner genruntime.MetaObject,
		resourceResolver *resolver.Resolver,
		armClient *genericarmclient.GenericClient,
		log logr.Logger,
		next PreReconcileOwnerCheckFunc,
	) (PreReconcileCheckResult, error)
}

type PreReconcileOwnerCheckFunc func(
	ctx context.Context,
	owner genruntime.MetaObject,
	resourceResolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (PreReconcileCheckResult, error)

// CreatePreReconciliationChecker creates a checker that can be used to check if a resource is ready for reconciliation.
// If the resource in question has not implemented the PreReconciliationChecker interface, the provided default checker
// is returned directly.
// We also return a bool indicating whether the resource extension implements the PreReconciliationChecker interface.
// host is a resource extension that may implement the PreReconciliationChecker interface.
func CreatePreReconciliationOwnerChecker(
	host genruntime.ResourceExtension,
) (PreReconcileOwnerCheckFunc, bool) {
	impl, ok := host.(PreReconciliationOwnerChecker)
	if !ok {
		return preReconciliationOwnerCheckerAlways, false
	}

	return func(
		ctx context.Context,
		owner genruntime.MetaObject,
		resourceResolver *resolver.Resolver,
		armClient *genericarmclient.GenericClient,
		log logr.Logger,
	) (PreReconcileCheckResult, error) {
		log.V(Status).Info("Extension pre-reconcile check running")

		result, err := impl.PreReconcileOwnerCheck(ctx, owner, resourceResolver, armClient, log, preReconciliationOwnerCheckerAlways)
		if err != nil {
			log.V(Status).Info(
				"Extension pre-reconcile owner check failed",
				"Error", err.Error())

			// We choose to skip here so that things are definitely broken and the user will notice
			// If we defaulted to always reconciling, the user might not notice that something is wrong
			return BlockReconcile("Extension PreReconcileOwnerCheck failed"), err
		}

		log.V(Status).Info(
			"Extension pre-reconcile owner check succeeded",
			"Result", result)

		return result, nil
	}, true
}

// preReconciliationOwnerCheckerAlways is a PreReconciliationOwnerChecker that always indicates a reconciliation is required.
// We have this here so we can set up a chain, even if it's only one link long.
// When we start doing proper comparisons between Spec and Status, we'll have an actual chain of checkers.
func preReconciliationOwnerCheckerAlways(
	_ context.Context,
	_ genruntime.MetaObject,
	_ *resolver.Resolver,
	_ *genericarmclient.GenericClient,
	_ logr.Logger,
) (PreReconcileCheckResult, error) {
	return ProceedWithReconcile(), nil
}
