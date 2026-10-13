/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	kusto "github.com/Azure/azure-service-operator/v2/api/kusto/v1api20230815/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.PreReconciliationChecker = &ClusterExtension{}

// PreReconcileCheck is called before the reconciliation of the resource to see if the cluster
// is in a state that will allow reconciliation to proceed.
// We can't try to create/update a Database unless the cluster is in a state that allows it.
func (ext *DatabaseExtension) PreReconcileCheck(
	_ context.Context,
	obj genruntime.MetaObject,
	owner genruntime.MetaObject,
	_ *resolver.Resolver,
	_ *genericarmclient.GenericClient,
	_ logr.Logger,
	_ extensions.PreReconcileCheckFunc,
) (extensions.PreReconcileCheckResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	cluster, ok := owner.(*kusto.Cluster)
	if !ok {
		return extensions.PreReconcileCheckResult{},
			eris.Errorf("cannot run on unknown resource type %T, expected owner to be *kusto.Cluster", owner)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = cluster

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

	return extensions.ProceedWithReconcile(), nil
}
