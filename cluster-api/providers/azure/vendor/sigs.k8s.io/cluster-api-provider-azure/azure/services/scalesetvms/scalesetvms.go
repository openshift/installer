/*
Copyright 2021 The Kubernetes Authors.

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

package scalesetvms

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/virtualmachines"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "scalesetvms"

type (
	// ScaleSetVMScope defines the scope interface for a scale sets service.
	ScaleSetVMScope interface {
		azure.ClusterDescriber
		azure.AsyncStatusUpdater
		ScaleSetVMSpec() azure.ResourceSpecGetter
		SetVMSSVM(vmssvm *azure.VMSSVM)
		SetVMSSVMState(state infrav1.ProvisioningState)
	}

	// Service provides operations on Azure resources.
	Service struct {
		Scope ScaleSetVMScope
		async.Reconciler
		VMReconciler async.Reconciler
	}
)

// NewService creates a new service.
func NewService(scope ScaleSetVMScope) (*Service, error) {
	client, err := newClient(scope)
	if err != nil {
		return nil, err
	}
	vmClient, err := virtualmachines.NewClient(scope)
	if err != nil {
		return nil, err
	}
	return &Service{
		Reconciler: async.New[armcompute.VirtualMachineScaleSetVMsClientUpdateResponse,
			armcompute.VirtualMachineScaleSetVMsClientDeleteResponse](scope, client, client),
		VMReconciler: async.New[armcompute.VirtualMachinesClientCreateOrUpdateResponse,
			armcompute.VirtualMachinesClientDeleteResponse](scope, vmClient, vmClient),
		Scope: scope,
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently gets, creates, and updates a scale set.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "scalesetvms.Service.Reconcile")
	defer done()

	spec := s.Scope.ScaleSetVMSpec()
	scaleSetVMSpec, ok := spec.(*ScaleSetVMSpec)
	if !ok {
		return errors.Errorf("%T is not of type ScaleSetVMSpec", spec)
	}

	reconciler := s.Reconciler
	var getter azure.ResourceSpecGetter = scaleSetVMSpec
	var result interface{}
	var err error
	// Fetch the latest instance or VM data. AzureMachinePoolReconciler handles model mutations.
	if scaleSetVMSpec.IsFlex {
		log.V(4).Info("VMSS is flex", "vmssName", scaleSetVMSpec.Name, "providerID", scaleSetVMSpec.ProviderID, "resourceID", scaleSetVMSpec.ResourceID)
		getter, err = scaleSetVMSpecToVMSpec(*scaleSetVMSpec)
		if err != nil {
			return errors.Wrap(err, "failed to convert scaleSetVMSpec to vmSpec")
		}
		reconciler = s.VMReconciler
	} else {
		log.V(4).Info("VMSS is uniform", "vmssName", scaleSetVMSpec.Name, "providerID", scaleSetVMSpec.ProviderID, "instanceID", scaleSetVMSpec.InstanceID)
	}

	// We only want to get the resource if it exists and handle the not found error.
	// We're using CreateOrUpdateResource() to do so but it doesn't actually create or update anything since getter.Parameters() always returns nil.
	result, err = reconciler.CreateOrUpdateResource(ctx, getter, serviceName)
	if err != nil {
		return err
	} else if result == nil {
		return azure.WithTransientError(fmt.Errorf("instance does not exist yet"), time.Second*30)
	}

	if scaleSetVMSpec.IsFlex {
		vm, ok := result.(armcompute.VirtualMachine)
		if !ok {
			return errors.Errorf("%T is not of type armcompute.VirtualMachine", result)
		}
		s.Scope.SetVMSSVM(converters.SDKVMToVMSSVM(vm, infrav1.FlexibleOrchestrationMode))
	} else {
		instance, ok := result.(armcompute.VirtualMachineScaleSetVM)
		if !ok {
			return errors.Errorf("%T is not of type armcompute.VirtualMachineScaleSetVM", result)
		}
		s.Scope.SetVMSSVM(converters.SDKToVMSSVM(instance))
	}

	return nil
}

// Delete deletes a scaleset instance asynchronously returning a future which encapsulates the long-running operation.
func (s *Service) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scalesetvms.Service.Delete")

	defer done()

	spec := s.Scope.ScaleSetVMSpec()
	scaleSetVMSpec, ok := spec.(*ScaleSetVMSpec)
	if !ok {
		return errors.Errorf("%T is not of type ScaleSetVMSpec", spec)
	}

	reconciler := s.Reconciler
	var getter azure.ResourceSpecGetter = scaleSetVMSpec
	var err error
	if scaleSetVMSpec.IsFlex {
		getter, err = scaleSetVMSpecToVMSpec(*scaleSetVMSpec)
		if err != nil {
			return errors.Wrap(err, "failed to convert scaleSetVMSpec to vmSpec")
		}
		reconciler = s.VMReconciler
	}

	err = reconciler.DeleteResource(ctx, getter, serviceName)
	if err != nil {
		s.Scope.SetVMSSVMState(infrav1.Deleting)
	} else {
		s.Scope.SetVMSSVMState(infrav1.Deleted)
	}

	return err
}

func scaleSetVMSpecToVMSpec(scaleSetVMSpec ScaleSetVMSpec) (*VMSSFlexGetter, error) {
	parsed, err := azureutil.ParseResourceID(scaleSetVMSpec.ResourceID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse resource id %q", scaleSetVMSpec.ResourceID))
	}
	resourceGroup, resourceName := parsed.ResourceGroupName, parsed.Name

	return &VMSSFlexGetter{
		Name:          resourceName,
		ResourceGroup: resourceGroup,
	}, nil
}
