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

package scalesets

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	azprovider "sigs.k8s.io/cloud-provider-azure/pkg/provider"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/resourceskus"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/slice"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "scalesets"

type (
	// ScaleSetScope defines the scope interface for a scale sets service.
	ScaleSetScope interface {
		azure.ClusterDescriber
		azure.AsyncStatusUpdater
		ScaleSetSpec(context.Context) azure.ResourceSpecGetter
		VMSSExtensionSpecs() []azure.ResourceSpecGetter
		SetAnnotation(string, string)
		SetProviderID(string)
		SetVMSSState(*azure.VMSS)
		ReconcileReplicas(context.Context, *azure.VMSS) error
	}

	// Service provides operations on Azure resources.
	Service struct {
		Scope ScaleSetScope
		Client
		resourceSKUCache *resourceskus.Cache
		async.Reconciler
	}
)

// New creates a new service.
func New(scope ScaleSetScope, skuCache *resourceskus.Cache) (*Service, error) {
	client, err := NewClient(scope, scope.DefaultedAzureCallTimeout())
	if err != nil {
		return nil, err
	}
	return &Service{
		Reconciler: async.New[armcompute.VirtualMachineScaleSetsClientCreateOrUpdateResponse,
			armcompute.VirtualMachineScaleSetsClientDeleteResponse](scope, client, client),
		Client:           client,
		Scope:            scope,
		resourceSKUCache: skuCache,
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently gets, creates, and updates a scale set.
func (s *Service) Reconcile(ctx context.Context) (retErr error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scalesets.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	if err := s.validateSpec(ctx); err != nil {
		// do as much early validation as possible to limit calls to Azure
		return err
	}

	spec := s.Scope.ScaleSetSpec(ctx)
	scaleSetSpec, ok := spec.(*ScaleSetSpec)
	if !ok {
		return errors.Errorf("%T is not of type ScaleSetSpec", spec)
	}

	result, err := s.Client.Get(ctx, spec)
	if err == nil {
		// We can only get the existing instances if the VMSS already exists
		scaleSetSpec.VMSSInstances, err = s.Client.ListInstances(ctx, spec.ResourceGroupName(), spec.ResourceName())
		if err != nil {
			err = errors.Wrapf(err, "failed to get existing VMSS instances")
			s.Scope.UpdatePutStatus(infrav1.BootstrapSucceededCondition, serviceName, err)
			return err
		}
		if result != nil {
			if err := s.updateScopeState(ctx, result, scaleSetSpec); err != nil {
				return err
			}
		}
	} else if !azure.ResourceNotFound(err) {
		return errors.Wrapf(err, "failed to get existing VMSS")
	}

	result, err = s.CreateOrUpdateResource(ctx, scaleSetSpec, serviceName)
	s.Scope.UpdatePutStatus(infrav1.BootstrapSucceededCondition, serviceName, err)

	if err == nil && result != nil {
		if err := s.updateScopeState(ctx, result, scaleSetSpec); err != nil {
			return err
		}
	}

	return err
}

// updateScopeState updates the scope's VMSS state and provider ID
//
// Code later in the reconciler uses scope's VMSS state for determining scale status and whether to create/delete
// AzureMachinePoolMachines.
// N.B.: before calling this function, make sure scaleSetSpec.VMSSInstances is updated to the latest state.
func (s *Service) updateScopeState(ctx context.Context, result interface{}, scaleSetSpec *ScaleSetSpec) error {
	vmss, ok := result.(armcompute.VirtualMachineScaleSet)
	if !ok {
		return errors.Errorf("%T is not an armcompute.VirtualMachineScaleSet", result)
	}

	fetchedVMSS := converters.SDKToVMSS(vmss, scaleSetSpec.VMSSInstances)
	if err := s.Scope.ReconcileReplicas(ctx, &fetchedVMSS); err != nil {
		return errors.Wrap(err, "unable to reconcile VMSS replicas")
	}

	// Transform the VMSS resource representation to conform to the cloud-provider-azure representation
	providerID, err := azprovider.ConvertResourceGroupNameToLower(azureutil.ProviderIDPrefix + fetchedVMSS.ID)
	if err != nil {
		return errors.Wrapf(err, "failed to parse VMSS ID %s", fetchedVMSS.ID)
	}
	s.Scope.SetProviderID(providerID)
	s.Scope.SetVMSSState(&fetchedVMSS)
	return nil
}

// Delete deletes a scale set asynchronously. Delete sends a DELETE request to Azure and if accepted without error,
// the VMSS will be considered deleted. The actual delete in Azure may take longer, but should eventually complete.
func (s *Service) Delete(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "scalesets.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	scaleSetSpec := s.Scope.ScaleSetSpec(ctx)

	defer func() {
		fetchedVMSS, err := s.getVirtualMachineScaleSet(ctx, scaleSetSpec)
		if err != nil && !azure.ResourceNotFound(err) {
			log.Error(err, "failed to get vmss in deferred update")
		}

		if fetchedVMSS != nil {
			s.Scope.SetVMSSState(fetchedVMSS)
		}
	}()

	err := s.DeleteResource(ctx, scaleSetSpec, serviceName)

	s.Scope.UpdateDeleteStatus(infrav1.BootstrapSucceededCondition, serviceName, err)

	return err
}

func (s *Service) validateSpec(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scalesets.Service.validateSpec")
	defer done()

	spec := s.Scope.ScaleSetSpec(ctx)
	scaleSetSpec, ok := spec.(*ScaleSetSpec)
	if !ok {
		return errors.Errorf("%T is not a ScaleSetSpec", spec)
	}

	sku, err := s.resourceSKUCache.Get(ctx, scaleSetSpec.Size, resourceskus.VirtualMachines)
	if err != nil {
		return errors.Wrapf(err, "failed to get SKU %s in compute api", scaleSetSpec.Size)
	}

	// Checking if the requested VM size has at least 2 vCPUS
	vCPUCapability, err := sku.HasCapabilityWithCapacity(resourceskus.VCPUs, resourceskus.MinimumVCPUS)
	if err != nil {
		return azure.WithTerminalError(errors.Wrap(err, "failed to validate the vCPU capability"))
	}

	if !vCPUCapability {
		return azure.WithTerminalError(errors.New("vm size should be bigger or equal to at least 2 vCPUs"))
	}

	// Checking if the requested VM size has at least 2 Gi of memory
	MemoryCapability, err := sku.HasCapabilityWithCapacity(resourceskus.MemoryGB, resourceskus.MinimumMemory)
	if err != nil {
		return azure.WithTerminalError(errors.Wrap(err, "failed to validate the memory capability"))
	}

	if !MemoryCapability {
		return azure.WithTerminalError(errors.New("vm memory should be bigger or equal to at least 2Gi"))
	}

	// enable ephemeral OS
	if scaleSetSpec.OSDisk.DiffDiskSettings != nil && !sku.HasCapability(resourceskus.EphemeralOSDisk) {
		return azure.WithTerminalError(fmt.Errorf("vm size %s does not support ephemeral os. select a different vm size or disable ephemeral os", scaleSetSpec.Size))
	}

	if scaleSetSpec.SecurityProfile != nil && !sku.HasCapability(resourceskus.EncryptionAtHost) {
		return azure.WithTerminalError(errors.Errorf("encryption at host is not supported for VM type %s", scaleSetSpec.Size))
	}

	// Fetch location and zone to check for their support of ultra disks.
	zones, err := s.resourceSKUCache.GetZones(ctx, scaleSetSpec.Location)
	if err != nil {
		return azure.WithTerminalError(errors.Wrapf(err, "failed to get the zones for location %s", scaleSetSpec.Location))
	}

	for _, zone := range zones {
		hasLocationCapability := sku.HasLocationCapability(resourceskus.UltraSSDAvailable, scaleSetSpec.Location, zone)
		err := fmt.Errorf("vm size %s does not support ultra disks in location %s. select a different vm size or disable ultra disks", scaleSetSpec.Size, scaleSetSpec.Location)

		// Check support for ultra disks as data disks.
		for _, disks := range scaleSetSpec.DataDisks {
			if disks.ManagedDisk != nil &&
				disks.ManagedDisk.StorageAccountType == string(armcompute.StorageAccountTypesUltraSSDLRS) &&
				!hasLocationCapability {
				return azure.WithTerminalError(err)
			}
		}
		// Check support for ultra disks as persistent volumes.
		if scaleSetSpec.AdditionalCapabilities != nil && scaleSetSpec.AdditionalCapabilities.UltraSSDEnabled != nil {
			if *scaleSetSpec.AdditionalCapabilities.UltraSSDEnabled &&
				!hasLocationCapability {
				return azure.WithTerminalError(err)
			}
		}
	}

	// Validate DiagnosticProfile spec
	if scaleSetSpec.DiagnosticsProfile != nil && scaleSetSpec.DiagnosticsProfile.Boot != nil {
		if scaleSetSpec.DiagnosticsProfile.Boot.StorageAccountType == infrav1.UserManagedDiagnosticsStorage {
			if scaleSetSpec.DiagnosticsProfile.Boot.UserManaged == nil {
				return azure.WithTerminalError(fmt.Errorf("userManaged must be specified when storageAccountType is '%s'", infrav1.UserManagedDiagnosticsStorage))
			} else if scaleSetSpec.DiagnosticsProfile.Boot.UserManaged.StorageAccountURI == "" {
				return azure.WithTerminalError(fmt.Errorf("storageAccountURI cannot be empty when storageAccountType is '%s'", infrav1.UserManagedDiagnosticsStorage))
			}
		}

		possibleStorageAccountTypeValues := []string{
			string(infrav1.DisabledDiagnosticsStorage),
			string(infrav1.ManagedDiagnosticsStorage),
			string(infrav1.UserManagedDiagnosticsStorage),
		}

		if !slice.Contains(possibleStorageAccountTypeValues, string(scaleSetSpec.DiagnosticsProfile.Boot.StorageAccountType)) {
			return azure.WithTerminalError(fmt.Errorf("invalid storageAccountType: %s. Allowed values are %v",
				scaleSetSpec.DiagnosticsProfile.Boot.StorageAccountType, possibleStorageAccountTypeValues))
		}
	}

	// Checking if selected availability zones are available selected VM type in location
	azsInLocation, err := s.resourceSKUCache.GetZonesWithVMSize(ctx, scaleSetSpec.Size, scaleSetSpec.Location)
	if err != nil {
		return errors.Wrapf(err, "failed to get zones for VM type %s in location %s", scaleSetSpec.Size, scaleSetSpec.Location)
	}

	for _, az := range scaleSetSpec.FailureDomains {
		if !slice.Contains(azsInLocation, az) {
			return azure.WithTerminalError(errors.Errorf("availability zone %s is not available for VM type %s in location %s", az, scaleSetSpec.Size, scaleSetSpec.Location))
		}
	}

	return nil
}

// getVirtualMachineScaleSet provides information about a Virtual Machine Scale Set and its instances.
func (s *Service) getVirtualMachineScaleSet(ctx context.Context, spec azure.ResourceSpecGetter) (*azure.VMSS, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scalesets.Service.getVirtualMachineScaleSet")
	defer done()

	vmssResult, err := s.Client.Get(ctx, spec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get existing VMSS")
	}
	vmss, ok := vmssResult.(armcompute.VirtualMachineScaleSet)
	if !ok {
		return nil, errors.Errorf("%T is not an armcompute.VirtualMachineScaleSet", vmssResult)
	}

	vmssInstances, err := s.Client.ListInstances(ctx, spec.ResourceGroupName(), spec.ResourceName())
	if err != nil {
		return nil, errors.Wrap(err, "failed to list instances")
	}

	result := converters.SDKToVMSS(vmss, vmssInstances)

	return &result, nil
}

// IsManaged returns always returns true as CAPZ does not support BYO scale set.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	return true, nil
}
