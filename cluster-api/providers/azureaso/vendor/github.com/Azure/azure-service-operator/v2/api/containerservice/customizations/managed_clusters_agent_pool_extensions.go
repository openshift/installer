/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	containerservice "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20250801/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var (
	_ extensions.PreReconciliationChecker      = &ManagedClustersAgentPoolExtension{}
	_ extensions.PreReconciliationOwnerChecker = &ManagedClustersAgentPoolExtension{}
)

// If an agent pool has a provisioningState not in this set, it will reject any attempt to PUT a new state out of
// hand; so there's no point in even trying. This is true even if the PUT we're doing will have no effect on the state
// of the agent pool.
// These are all listed lowercase, so we can do a case-insensitive match.
var nonBlockingManagedClustersAgentPoolProvisioningStates = set.Make(
	"succeeded",
	"failed",
	"canceled",
)

func (ext *ManagedClustersAgentPoolExtension) PreReconcileOwnerCheck(
	ctx context.Context,
	owner genruntime.MetaObject,
	resourceResolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
	next extensions.PreReconcileOwnerCheckFunc,
) (extensions.PreReconcileCheckResult, error) {
	// Check to see if the owning cluster is in a state that will block us from reconciling
	if managedCluster, ok := owner.(*containerservice.ManagedCluster); ok {
		state := managedCluster.Status.ProvisioningState
		if state != nil && clusterProvisioningStateBlocksReconciliation(state) {
			return extensions.BlockReconcile(
					fmt.Sprintf("Managed cluster %q is in provisioning state %q", owner.GetName(), *state)),
				nil
		}
	}

	return next(ctx, owner, resourceResolver, armClient, log)
}

func (ext *ManagedClustersAgentPoolExtension) PreReconcileCheck(
	ctx context.Context,
	obj genruntime.MetaObject,
	resourceResolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
	next extensions.PreReconcileCheckFunc,
) (extensions.PreReconcileCheckResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	agentPool, ok := obj.(*containerservice.ManagedClustersAgentPool)
	if !ok {
		return extensions.PreReconcileCheckResult{},
			eris.Errorf("cannot run on unknown resource type %T, expected *containerservice.ManagedCluster", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = agentPool

	// If the agent pool is in a state that will reject any PUT, then we should skip reconciliation
	// as there's no point in even trying.
	// This allows us to "play nice with others" and not use up request quota attempting to make changes when we
	// already know those attempts will fail.
	state := agentPool.Status.ProvisioningState
	if state != nil && agentPoolProvisioningStateBlocksReconciliation(state) {
		return extensions.BlockReconcile(
				fmt.Sprintf("Managed cluster agent pool is in provisioning state %q", *state)),
			nil
	}

	return next(ctx, obj, resourceResolver, armClient, log)
}

func agentPoolProvisioningStateBlocksReconciliation(provisioningState *string) bool {
	if provisioningState == nil {
		return false
	}

	return !nonBlockingManagedClustersAgentPoolProvisioningStates.Contains(strings.ToLower(*provisioningState))
}
