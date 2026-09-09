/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	kusto "github.com/Azure/azure-service-operator/v2/api/kusto/v1api20240413/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.PreReconciliationOwnerChecker = &DatabaseExtension{}

// Ensure we're dealing with the hub version of Cluster
// If this no longer compiles (due to a newer version being imported), change the import above to that version.
var _ conversion.Hub = &kusto.Cluster{}

// PreReconcileOwnerCheck is called before the reconciliation of the resource to see if the cluster
// is in a state that will allow reconciliation to proceed.
// We can't try to create/update a Database unless the cluster is in a state that allows it.
// We use PreReconcileOwnerCheck instead of PreReconcileCheck because you can't even GET a database if the cluster is powered down.
func (ext *DatabaseExtension) PreReconcileOwnerCheck(
	ctx context.Context,
	owner genruntime.MetaObject,
	resourceResolver *resolver.Resolver,
	armClient_ *genericarmclient.GenericClient,
	log logr.Logger,
	next extensions.PreReconcileOwnerCheckFunc,
) (extensions.PreReconcileCheckResult, error) {
	// Check to see if the owning cluster is in a state that will block us from reconciling
	if cluster, ok := owner.(*kusto.Cluster); ok {
		// If our owning *cluster* is in a state that will reject any PUT, then we should skip
		// reconciliation of the database as there's no point in even trying.
		// One way this can happen is when we reconcile the cluster, putting it into an `Updating`
		// state for a period. While it's in that state, we can't even try to reconcile the database as
		// the operation will fail with a `Conflict` error.
		// Checking the state of our owning cluster allows us to "play nice with others" and not use up
		// request quota attempting to make changes when we already know those attempts will fail.
		state := cluster.Status.ProvisioningState
		if state != nil && clusterProvisioningStateBlocksReconciliation(state) {
			return extensions.BlockReconcile(
					fmt.Sprintf("Owning cluster is in provisioning state %q", *state)),
				nil
		}

	}

	return next(ctx, owner, resourceResolver, armClient_, log)
}
