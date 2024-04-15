/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	asocontainerservicev1preview "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20230202preview"
	asocontainerservicev1 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001"
	"github.com/pkg/errors"
	azprovider "sigs.k8s.io/cloud-provider-azure/pkg/provider"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/scope"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/agentpools"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/scalesets"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

type (
	// azureManagedMachinePoolService contains the services required by the cluster controller.
	azureManagedMachinePoolService struct {
		scope         agentpools.AgentPoolScope
		agentPoolsSvc azure.Reconciler
		scaleSetsSvc  NodeLister
	}

	// AgentPoolVMSSNotFoundError represents a reconcile error when the VMSS for an agent pool can't be found.
	AgentPoolVMSSNotFoundError struct {
		NodeResourceGroup string
		PoolName          string
	}

	// NodeLister is a service interface for returning generic lists.
	NodeLister interface {
		ListInstances(context.Context, string, string) ([]armcompute.VirtualMachineScaleSetVM, error)
		List(context.Context, string) ([]armcompute.VirtualMachineScaleSet, error)
	}
)

// NewAgentPoolVMSSNotFoundError creates a new AgentPoolVMSSNotFoundError.
func NewAgentPoolVMSSNotFoundError(nodeResourceGroup, poolName string) *AgentPoolVMSSNotFoundError {
	return &AgentPoolVMSSNotFoundError{
		NodeResourceGroup: nodeResourceGroup,
		PoolName:          poolName,
	}
}

func (a *AgentPoolVMSSNotFoundError) Error() string {
	msgFmt := "failed to find vm scale set in resource group %s matching pool named %s"
	return fmt.Sprintf(msgFmt, a.NodeResourceGroup, a.PoolName)
}

// Is returns true if the target error is an `AgentPoolVMSSNotFoundError`.
func (a *AgentPoolVMSSNotFoundError) Is(target error) bool {
	var err *AgentPoolVMSSNotFoundError
	ok := errors.As(target, &err)
	return ok
}

// newAzureManagedMachinePoolService populates all the services based on input scope.
func newAzureManagedMachinePoolService(scope *scope.ManagedMachinePoolScope, apiCallTimeout time.Duration) (*azureManagedMachinePoolService, error) {
	scaleSetAuthorizer, err := scaleSetAuthorizer(scope)
	if err != nil {
		return nil, err
	}
	scaleSetsClient, err := scalesets.NewClient(scaleSetAuthorizer, apiCallTimeout)
	if err != nil {
		return nil, err
	}
	return &azureManagedMachinePoolService{
		scope:         scope,
		agentPoolsSvc: agentpools.New(scope),
		scaleSetsSvc:  scaleSetsClient,
	}, nil
}

// scaleSetAuthorizer takes a scope and determines if a regional authorizer is needed for scale sets
// see https://github.com/kubernetes-sigs/cluster-api-provider-azure/pull/1850 for context on region based authorizer.
func scaleSetAuthorizer(scope *scope.ManagedMachinePoolScope) (azure.Authorizer, error) {
	if scope.ControlPlane.Spec.AzureEnvironment == azure.PublicCloudName {
		return azure.WithRegionalBaseURI(scope, scope.Location()) // public cloud supports regional end points
	}

	return scope, nil
}

// Reconcile reconciles all the services in a predetermined order.
func (s *azureManagedMachinePoolService) Reconcile(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.azureManagedMachinePoolService.Reconcile")
	defer done()

	s.scope.SetSubnetName()

	log.Info("reconciling managed machine pool")
	agentPool, err := s.scope.AgentPoolSpec().Parameters(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to get agent pool parameters")
	}
	var agentPoolName string
	if s.scope.IsPreviewEnabled() {
		agentPoolTyped := agentPool.(*asocontainerservicev1preview.ManagedClustersAgentPool)
		agentPoolName = agentPoolTyped.AzureName()
	} else {
		agentPoolTyped := agentPool.(*asocontainerservicev1.ManagedClustersAgentPool)
		agentPoolName = agentPoolTyped.AzureName()
	}

	if err := s.agentPoolsSvc.Reconcile(ctx); err != nil {
		return errors.Wrapf(err, "failed to reconcile machine pool %s", agentPoolName)
	}

	nodeResourceGroup := s.scope.NodeResourceGroup()
	vmss, err := s.scaleSetsSvc.List(ctx, nodeResourceGroup)
	if err != nil {
		return errors.Wrapf(err, "failed to list vmss in resource group %s", nodeResourceGroup)
	}

	var match *armcompute.VirtualMachineScaleSet
	for _, ss := range vmss {
		ss := ss
		if ss.Tags["poolName"] != nil && *ss.Tags["poolName"] == agentPoolName {
			match = &ss
			break
		}

		if ss.Tags["aks-managed-poolName"] != nil && *ss.Tags["aks-managed-poolName"] == agentPoolName {
			match = &ss
			break
		}
	}

	if match == nil {
		return azure.WithTransientError(NewAgentPoolVMSSNotFoundError(nodeResourceGroup, agentPoolName), 20*time.Second)
	}

	instances, err := s.scaleSetsSvc.ListInstances(ctx, nodeResourceGroup, *match.Name)
	if err != nil {
		return errors.Wrapf(err, "failed to reconcile machine pool %s", agentPoolName)
	}

	var providerIDs = make([]string, len(instances))
	for i := 0; i < len(instances); i++ {
		// Transform the VMSS instance resource representation to conform to the cloud-provider-azure representation
		providerID, err := azprovider.ConvertResourceGroupNameToLower(azureutil.ProviderIDPrefix + *instances[i].ID)
		if err != nil {
			return errors.Wrapf(err, "failed to parse instance ID %s", *instances[i].ID)
		}
		providerIDs[i] = providerID
	}

	s.scope.SetAgentPoolProviderIDList(providerIDs)
	s.scope.SetAgentPoolReplicas(int32(len(providerIDs)))
	s.scope.SetAgentPoolReady(true)

	log.Info("reconciled managed machine pool successfully")
	return nil
}

// Pause pauses all components making up the machine pool.
func (s *azureManagedMachinePoolService) Pause(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.azureManagedMachinePoolService.Pause")
	defer done()

	pauser, ok := s.agentPoolsSvc.(azure.Pauser)
	if !ok {
		return nil
	}
	if err := pauser.Pause(ctx); err != nil {
		return errors.Wrapf(err, "failed to pause machine pool %s", s.scope.Name())
	}

	return nil
}

// Delete reconciles all the services in a predetermined order.
func (s *azureManagedMachinePoolService) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.azureManagedMachinePoolService.Delete")
	defer done()

	if err := s.agentPoolsSvc.Delete(ctx); err != nil {
		return errors.Wrapf(err, "failed to delete machine pool %s", s.scope.Name())
	}

	return nil
}
