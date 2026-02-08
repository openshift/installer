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

	if err := s.validateVMCapabilities(sku, scaleSetSpec); err != nil {
		return err
	}

	if err := s.validateUltraDiskSupport(ctx, sku, scaleSetSpec); err != nil {
		return err
	}

	if err := validateDiagnosticsProfile(scaleSetSpec); err != nil {
		return err
	}

	return s.validateAvailabilityZones(ctx, scaleSetSpec)
}

// validateVMCapabilities validates the VM size capabilities (vCPU, memory, ephemeral OS, encryption).
func (s *Service) validateVMCapabilities(sku resourceskus.SKU, spec *ScaleSetSpec) error {
	// Validate minimum vCPU requirement.
	vCPUCapability, err := sku.HasCapabilityWithCapacity(resourceskus.VCPUs, resourceskus.MinimumVCPUS)
	if err != nil {
		return azure.WithTerminalError(errors.Wrap(err, "failed to validate the vCPU capability"))
	}
	if !vCPUCapability {
		return azure.WithTerminalError(errors.New("vm size should be bigger or equal to at least 2 vCPUs"))
	}

	// Validate minimum memory requirement.
	memoryCapability, err := sku.HasCapabilityWithCapacity(resourceskus.MemoryGB, resourceskus.MinimumMemory)
	if err != nil {
		return azure.WithTerminalError(errors.Wrap(err, "failed to validate the memory capability"))
	}
	if !memoryCapability {
		return azure.WithTerminalError(errors.New("vm memory should be bigger or equal to at least 2Gi"))
	}

	// Validate ephemeral OS disk support.
	if spec.OSDisk.DiffDiskSettings != nil && !sku.HasCapability(resourceskus.EphemeralOSDisk) {
		return azure.WithTerminalError(fmt.Errorf("vm size %s does not support ephemeral os. select a different vm size or disable ephemeral os", spec.Size))
	}

	// Validate encryption at host support.
	if spec.SecurityProfile != nil && !sku.HasCapability(resourceskus.EncryptionAtHost) {
		return azure.WithTerminalError(errors.Errorf("encryption at host is not supported for VM type %s", spec.Size))
	}

	return nil
}

// validateUltraDiskSupport validates ultra disk support for the specified location and zones.
func (s *Service) validateUltraDiskSupport(ctx context.Context, sku resourceskus.SKU, spec *ScaleSetSpec) error {
	if !s.requiresUltraDiskValidation(spec) {
		return nil
	}

	zones, err := s.resourceSKUCache.GetZones(ctx, spec.Location)
	if err != nil {
		return azure.WithTerminalError(errors.Wrapf(err, "failed to get the zones for location %s", spec.Location))
	}

	for _, zone := range zones {
		hasLocationCapability := sku.HasLocationCapability(resourceskus.UltraSSDAvailable, spec.Location, zone)
		if err := s.validateUltraDiskInZone(spec, hasLocationCapability); err != nil {
			return err
		}
	}

	return nil
}

// requiresUltraDiskValidation checks if ultra disk validation is needed.
func (s *Service) requiresUltraDiskValidation(spec *ScaleSetSpec) bool {
	// Check if any data disk uses ultra SSD.
	for _, disk := range spec.DataDisks {
		if disk.ManagedDisk != nil && disk.ManagedDisk.StorageAccountType == string(armcompute.StorageAccountTypesUltraSSDLRS) {
			return true
		}
	}

	// Check if ultra SSD is enabled via additional capabilities.
	if spec.AdditionalCapabilities != nil && spec.AdditionalCapabilities.UltraSSDEnabled != nil {
		return *spec.AdditionalCapabilities.UltraSSDEnabled
	}

	return false
}

// validateUltraDiskInZone validates ultra disk support for a specific zone.
func (s *Service) validateUltraDiskInZone(spec *ScaleSetSpec, hasLocationCapability bool) error {
	if hasLocationCapability {
		return nil
	}

	ultraDiskErr := fmt.Errorf("vm size %s does not support ultra disks in location %s. select a different vm size or disable ultra disks", spec.Size, spec.Location)

	// Check data disks.
	for _, disk := range spec.DataDisks {
		if disk.ManagedDisk != nil && disk.ManagedDisk.StorageAccountType == string(armcompute.StorageAccountTypesUltraSSDLRS) {
			return azure.WithTerminalError(ultraDiskErr)
		}
	}

	// Check additional capabilities.
	if spec.AdditionalCapabilities != nil && spec.AdditionalCapabilities.UltraSSDEnabled != nil {
		if *spec.AdditionalCapabilities.UltraSSDEnabled {
			return azure.WithTerminalError(ultraDiskErr)
		}
	}

	return nil
}

// validateDiagnosticsProfile validates the diagnostics profile configuration.
func validateDiagnosticsProfile(spec *ScaleSetSpec) error {
	if spec.DiagnosticsProfile == nil || spec.DiagnosticsProfile.Boot == nil {
		return nil
	}

	boot := spec.DiagnosticsProfile.Boot

	// Validate user-managed storage account configuration.
	if boot.StorageAccountType == infrav1.UserManagedDiagnosticsStorage {
		if boot.UserManaged == nil {
			return azure.WithTerminalError(fmt.Errorf("userManaged must be specified when storageAccountType is '%s'", infrav1.UserManagedDiagnosticsStorage))
		}
		if boot.UserManaged.StorageAccountURI == "" {
			return azure.WithTerminalError(fmt.Errorf("storageAccountURI cannot be empty when storageAccountType is '%s'", infrav1.UserManagedDiagnosticsStorage))
		}
	}

	// Validate storage account type is valid.
	validStorageTypes := []string{
		string(infrav1.DisabledDiagnosticsStorage),
		string(infrav1.ManagedDiagnosticsStorage),
		string(infrav1.UserManagedDiagnosticsStorage),
	}
	if !slice.Contains(validStorageTypes, string(boot.StorageAccountType)) {
		return azure.WithTerminalError(fmt.Errorf("invalid storageAccountType: %s. Allowed values are %v", boot.StorageAccountType, validStorageTypes))
	}

	return nil
}

// validateAvailabilityZones validates that the requested availability zones are available for the VM size.
func (s *Service) validateAvailabilityZones(ctx context.Context, spec *ScaleSetSpec) error {
	if len(spec.FailureDomains) == 0 {
		return nil
	}

	azsInLocation, err := s.resourceSKUCache.GetZonesWithVMSize(ctx, spec.Size, spec.Location)
	if err != nil {
		return errors.Wrapf(err, "failed to get zones for VM type %s in location %s", spec.Size, spec.Location)
	}

	for _, az := range spec.FailureDomains {
		if !slice.Contains(azsInLocation, az) {
			return azure.WithTerminalError(errors.Errorf("availability zone %s is not available for VM type %s in location %s", az, spec.Size, spec.Location))
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
func (s *Service) IsManaged(_ context.Context) (bool, error) {
	return true, nil
}
