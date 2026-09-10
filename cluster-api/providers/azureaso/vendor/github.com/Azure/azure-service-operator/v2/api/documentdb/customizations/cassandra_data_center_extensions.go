/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package customizations

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	hub "github.com/Azure/azure-service-operator/v2/api/documentdb/v20251015/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.PreReconciliationOwnerChecker = &CassandraDataCenterExtension{}

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
func (ext *CassandraDataCenterExtension) PreReconcileOwnerCheck(
	ctx context.Context,
	owner genruntime.MetaObject,
	resourceResolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
	next extensions.PreReconcileOwnerCheckFunc,
) (extensions.PreReconcileCheckResult, error) {
	// Ensure we have the expected owner type.
	if server, ok := owner.(*hub.CassandraCluster); ok {
		serverState := server.Status.Properties.ProvisioningState
		if serverState != nil && cassandraClusterProvisioningStateBlocksReconciliation(*serverState) {
			return extensions.BlockReconcile(
				fmt.Sprintf(
					"Owning CassandraCluster is in provisioning state %q",
					*serverState)), nil
		}
	}

	return next(ctx, owner, resourceResolver, armClient, log)
}

func cassandraClusterProvisioningStateBlocksReconciliation(state string) bool {
	// If the cluster is in a provisioning state that isn't "succeeded", then we should block reconciliation of the
	// data center. A transient state (e.g. `Creating`) will block because things aren't ready; a `Failed` will block
	// because the cluster is in a bad state.
	return state != "Succeeded"
}
