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
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	containerservice "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001/storage"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.PreReconciliationChecker = &ManagedClustersAgentPoolExtension{}

// If an agent pool has a provisioningState not in this set, it will reject any attempt to PUT a new state out of
// hand; so there's no point in even trying. This is true even if the PUT we're doing will have no effect on the state
// of the agent pool.
// These are all listed lowercase, so we can do a case-insensitive match.
var nonBlockingManagedClustersAgentPoolProvisioningStates = set.Make(
	"succeeded",
	"failed",
	"canceled",
)

func (ext *ManagedClustersAgentPoolExtension) PreReconcileCheck(
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
	agentPool, ok := obj.(*containerservice.ManagedClustersAgentPool)
	if !ok {
		return extensions.PreReconcileCheckResult{},
			errors.Errorf("cannot run on unknown resource type %T, expected *containerservice.ManagedCluster", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = agentPool

	// Check to see if the owning cluster is in a state that will block us from reconciling
	if owner != nil {
		if managedCluster, ok := owner.(*containerservice.ManagedCluster); ok {
			state := managedCluster.Status.ProvisioningState
			if state != nil && clusterProvisioningStateBlocksReconciliation(state) {
				return extensions.BlockReconcile(
						fmt.Sprintf("Managed cluster %q is in provisioning state %q", owner.GetName(), *state)),
					nil
			}

		}
	}

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

	return extensions.ProceedWithReconcile(), nil
}

func agentPoolProvisioningStateBlocksReconciliation(provisioningState *string) bool {
	if provisioningState == nil {
		return false
	}

	return !nonBlockingManagedClustersAgentPoolProvisioningStates.Contains(strings.ToLower(*provisioningState))
}
