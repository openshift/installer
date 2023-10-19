/*
Copyright 2023 The Kubernetes Authors.

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
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/resourceskus"
	"sigs.k8s.io/cluster-api-provider-azure/util/generators"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// ScaleSetSpec defines the specification for a Scale Set.
type ScaleSetSpec struct {
	Name                         string
	ResourceGroup                string
	Size                         string
	Capacity                     int64
	SSHKeyData                   string
	OSDisk                       infrav1.OSDisk
	DataDisks                    []infrav1.DataDisk
	SubnetName                   string
	VNetName                     string
	VNetResourceGroup            string
	PublicLBName                 string
	PublicLBAddressPoolName      string
	AcceleratedNetworking        *bool
	TerminateNotificationTimeout *int
	Identity                     infrav1.VMIdentity
	UserAssignedIdentities       []infrav1.UserAssignedIdentity
	SecurityProfile              *infrav1.SecurityProfile
	SpotVMOptions                *infrav1.SpotVMOptions
	AdditionalCapabilities       *infrav1.AdditionalCapabilities
	DiagnosticsProfile           *infrav1.Diagnostics
	FailureDomains               []string
	VMExtensions                 []infrav1.VMExtension
	NetworkInterfaces            []infrav1.NetworkInterface
	IPv6Enabled                  bool
	OrchestrationMode            infrav1.OrchestrationModeType
	Location                     string
	SubscriptionID               string
	SKU                          resourceskus.SKU
	VMSSExtensionSpecs           []azure.ResourceSpecGetter
	VMImage                      *infrav1.Image
	BootstrapData                string
	VMSSInstances                []armcompute.VirtualMachineScaleSetVM
	MaxSurge                     int
	ClusterName                  string
	ShouldPatchCustomData        bool
	HasReplicasExternallyManaged bool
	AdditionalTags               infrav1.Tags
}

// ResourceName returns the name of the Scale Set.
func (s *ScaleSetSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group for this Scale Set.
func (s *ScaleSetSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for Scale Sets.
func (s *ScaleSetSpec) OwnerResourceName() string {
	return ""
}

func (s *ScaleSetSpec) existingParameters(ctx context.Context, existing interface{}) (parameters interface{}, err error) {
	existingVMSS, ok := existing.(armcompute.VirtualMachineScaleSet)
	if !ok {
		return nil, errors.Errorf("%T is not an armcompute.VirtualMachineScaleSet", existing)
	}

	existingInfraVMSS := converters.SDKToVMSS(existingVMSS, s.VMSSInstances)

	params, err := s.Parameters(ctx, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate scale set update parameters for %s", s.Name)
	}

	vmss, ok := params.(armcompute.VirtualMachineScaleSet)
	if !ok {
		return nil, errors.Errorf("%T is not an armcompute.VirtualMachineScaleSet", existing)
	}

	vmss.Properties.VirtualMachineProfile.NetworkProfile = nil
	vmss.ID = existingVMSS.ID

	hasModelChanges := hasModelModifyingDifferences(&existingInfraVMSS, vmss)
	isFlex := s.OrchestrationMode == infrav1.FlexibleOrchestrationMode
	updated := true
	if !isFlex {
		updated = existingInfraVMSS.HasEnoughLatestModelOrNotMixedModel()
	}
	if s.MaxSurge > 0 && (hasModelChanges || !updated) && !s.HasReplicasExternallyManaged {
		// surge capacity with the intention of lowering during instance reconciliation
		surge := s.Capacity + int64(s.MaxSurge)
		vmss.SKU.Capacity = ptr.To[int64](surge)
	}

	// If there are no model changes and no increase in the replica count, do not update the VMSS.
	// Decreases in replica count is handled by deleting AzureMachinePoolMachine instances in the MachinePoolScope
	if *vmss.SKU.Capacity <= existingInfraVMSS.Capacity && !hasModelChanges && !s.ShouldPatchCustomData {
		// up to date, nothing to do
		return nil, nil
	}

	return vmss, nil
}

// Parameters returns the parameters for the Scale Set.
func (s *ScaleSetSpec) Parameters(ctx context.Context, existing interface{}) (parameters interface{}, err error) {
	if existing != nil {
		return s.existingParameters(ctx, existing)
	}

	if s.AcceleratedNetworking == nil {
		// set accelerated networking to the capability of the VMSize
		accelNet := s.SKU.HasCapability(resourceskus.AcceleratedNetworking)
		s.AcceleratedNetworking = &accelNet
	}

	extensions, err := s.generateExtensions(ctx)
	if err != nil {
		return armcompute.VirtualMachineScaleSet{}, err
	}

	storageProfile, err := s.generateStorageProfile(ctx)
	if err != nil {
		return armcompute.VirtualMachineScaleSet{}, err
	}

	securityProfile, err := s.getSecurityProfile()
	if err != nil {
		return armcompute.VirtualMachineScaleSet{}, err
	}

	priority, evictionPolicy, billingProfile, err := converters.GetSpotVMOptions(s.SpotVMOptions, s.OSDisk.DiffDiskSettings)
	if err != nil {
		return armcompute.VirtualMachineScaleSet{}, errors.Wrapf(err, "failed to get Spot VM options")
	}

	diagnosticsProfile := converters.GetDiagnosticsProfile(s.DiagnosticsProfile)

	osProfile, err := s.generateOSProfile(ctx)
	if err != nil {
		return armcompute.VirtualMachineScaleSet{}, err
	}

	orchestrationMode := converters.GetOrchestrationMode(s.OrchestrationMode)

	vmss := armcompute.VirtualMachineScaleSet{
		Location: ptr.To(s.Location),
		SKU: &armcompute.SKU{
			Name:     ptr.To(s.Size),
			Tier:     ptr.To("Standard"),
			Capacity: ptr.To[int64](s.Capacity),
		},
		Zones: azure.PtrSlice(&s.FailureDomains),
		Plan:  s.generateImagePlan(ctx),
		Properties: &armcompute.VirtualMachineScaleSetProperties{
			OrchestrationMode:    ptr.To(orchestrationMode),
			SinglePlacementGroup: ptr.To(false),
			VirtualMachineProfile: &armcompute.VirtualMachineScaleSetVMProfile{
				OSProfile:          osProfile,
				StorageProfile:     storageProfile,
				SecurityProfile:    securityProfile,
				DiagnosticsProfile: diagnosticsProfile,
				NetworkProfile: &armcompute.VirtualMachineScaleSetNetworkProfile{
					NetworkInterfaceConfigurations: azure.PtrSlice(s.getVirtualMachineScaleSetNetworkConfiguration()),
				},
				Priority:       priority,
				EvictionPolicy: evictionPolicy,
				BillingProfile: billingProfile,
				ExtensionProfile: &armcompute.VirtualMachineScaleSetExtensionProfile{
					Extensions: azure.PtrSlice(&extensions),
				},
			},
		},
	}

	// Set properties specific to VMSS orchestration mode
	// See https://learn.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-orchestration-modes for more details
	switch orchestrationMode {
	case armcompute.OrchestrationModeUniform: // Uniform VMSS
		vmss.Properties.Overprovision = ptr.To(false)
		vmss.Properties.UpgradePolicy = &armcompute.UpgradePolicy{Mode: ptr.To(armcompute.UpgradeModeManual)}
	case armcompute.OrchestrationModeFlexible: // VMSS Flex, VMs are treated as individual virtual machines
		vmss.Properties.VirtualMachineProfile.NetworkProfile.NetworkAPIVersion =
			ptr.To(armcompute.NetworkAPIVersionTwoThousandTwenty1101)
		vmss.Properties.PlatformFaultDomainCount = ptr.To[int32](1)
		if len(s.FailureDomains) > 1 {
			vmss.Properties.PlatformFaultDomainCount = ptr.To[int32](int32(len(s.FailureDomains)))
		}
	}

	// Assign Identity to VMSS
	if s.Identity == infrav1.VMIdentitySystemAssigned {
		vmss.Identity = &armcompute.VirtualMachineScaleSetIdentity{
			Type: ptr.To(armcompute.ResourceIdentityTypeSystemAssigned),
		}
	} else if s.Identity == infrav1.VMIdentityUserAssigned {
		userIdentitiesMap, err := converters.UserAssignedIdentitiesToVMSSSDK(s.UserAssignedIdentities)
		if err != nil {
			return vmss, errors.Wrapf(err, "failed to assign identity %q", s.Name)
		}
		vmss.Identity = &armcompute.VirtualMachineScaleSetIdentity{
			Type:                   ptr.To(armcompute.ResourceIdentityTypeUserAssigned),
			UserAssignedIdentities: userIdentitiesMap,
		}
	}

	// Provisionally detect whether there is any Data Disk defined which uses UltraSSDs.
	// If that's the case, enable the UltraSSD capability.
	for _, dataDisk := range s.DataDisks {
		if dataDisk.ManagedDisk != nil && dataDisk.ManagedDisk.StorageAccountType == string(armcompute.StorageAccountTypesUltraSSDLRS) {
			vmss.Properties.AdditionalCapabilities = &armcompute.AdditionalCapabilities{
				UltraSSDEnabled: ptr.To(true),
			}
		}
	}

	// Set Additional Capabilities if any is present on the spec.
	if s.AdditionalCapabilities != nil {
		// Set UltraSSDEnabled if a specific value is set on the spec for it.
		if s.AdditionalCapabilities.UltraSSDEnabled != nil {
			vmss.Properties.AdditionalCapabilities.UltraSSDEnabled = s.AdditionalCapabilities.UltraSSDEnabled
		}
	}

	if s.TerminateNotificationTimeout != nil {
		vmss.Properties.VirtualMachineProfile.ScheduledEventsProfile = &armcompute.ScheduledEventsProfile{
			TerminateNotificationProfile: &armcompute.TerminateNotificationProfile{
				NotBeforeTimeout: ptr.To(fmt.Sprintf("PT%dM", *s.TerminateNotificationTimeout)),
				Enable:           ptr.To(true),
			},
		}
	}

	tags := infrav1.Build(infrav1.BuildParams{
		ClusterName: s.ClusterName,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        ptr.To(s.Name),
		Role:        ptr.To(infrav1.Node),
		Additional:  s.AdditionalTags,
	})

	vmss.Tags = converters.TagsToMap(tags)
	return vmss, nil
}

func hasModelModifyingDifferences(infraVMSS *azure.VMSS, vmss armcompute.VirtualMachineScaleSet) bool {
	other := converters.SDKToVMSS(vmss, []armcompute.VirtualMachineScaleSetVM{})
	return infraVMSS.HasModelChanges(other)
}

func (s *ScaleSetSpec) generateExtensions(ctx context.Context) ([]armcompute.VirtualMachineScaleSetExtension, error) {
	extensions := make([]armcompute.VirtualMachineScaleSetExtension, len(s.VMSSExtensionSpecs))
	for i, extensionSpec := range s.VMSSExtensionSpecs {
		extensionSpec := extensionSpec
		parameters, err := extensionSpec.Parameters(ctx, nil)
		if err != nil {
			return nil, err
		}
		vmssextension, ok := parameters.(armcompute.VirtualMachineScaleSetExtension)
		if !ok {
			return nil, errors.Errorf("%T is not an armcompute.VirtualMachineScaleSetExtension", parameters)
		}
		extensions[i] = vmssextension
	}

	return extensions, nil
}

func (s *ScaleSetSpec) getVirtualMachineScaleSetNetworkConfiguration() *[]armcompute.VirtualMachineScaleSetNetworkConfiguration {
	var backendAddressPools []armcompute.SubResource
	if s.PublicLBName != "" {
		if s.PublicLBAddressPoolName != "" {
			backendAddressPools = append(backendAddressPools,
				armcompute.SubResource{
					ID: ptr.To(azure.AddressPoolID(s.SubscriptionID, s.ResourceGroup, s.PublicLBName, s.PublicLBAddressPoolName)),
				})
		}
	}
	nicConfigs := []armcompute.VirtualMachineScaleSetNetworkConfiguration{}
	for i, n := range s.NetworkInterfaces {
		nicConfig := armcompute.VirtualMachineScaleSetNetworkConfiguration{}
		nicConfig.Properties = &armcompute.VirtualMachineScaleSetNetworkConfigurationProperties{}
		nicConfig.Name = ptr.To(s.Name + "-nic-" + strconv.Itoa(i))
		nicConfig.Properties.EnableIPForwarding = ptr.To(true)
		if n.AcceleratedNetworking != nil {
			nicConfig.Properties.EnableAcceleratedNetworking = n.AcceleratedNetworking
		} else {
			// If AcceleratedNetworking is not specified, use the value from the VMSS spec.
			// It will be set to true if the VMSS SKU supports it.
			nicConfig.Properties.EnableAcceleratedNetworking = s.AcceleratedNetworking
		}

		// Create IPConfigs
		ipconfigs := []armcompute.VirtualMachineScaleSetIPConfiguration{}
		for j := 0; j < n.PrivateIPConfigs; j++ {
			ipconfig := armcompute.VirtualMachineScaleSetIPConfiguration{
				Name: ptr.To(fmt.Sprintf("ipConfig" + strconv.Itoa(j))),
				Properties: &armcompute.VirtualMachineScaleSetIPConfigurationProperties{
					PrivateIPAddressVersion: ptr.To(armcompute.IPVersionIPv4),
					Subnet: &armcompute.APIEntityReference{
						ID: ptr.To(azure.SubnetID(s.SubscriptionID, s.VNetResourceGroup, s.VNetName, n.SubnetName)),
					},
				},
			}

			if j == 0 {
				// Always use the first IPConfig as the Primary
				ipconfig.Properties.Primary = ptr.To(true)
			}
			ipconfigs = append(ipconfigs, ipconfig)
		}
		if s.IPv6Enabled {
			ipv6Config := armcompute.VirtualMachineScaleSetIPConfiguration{
				Name: ptr.To("ipConfigv6"),
				Properties: &armcompute.VirtualMachineScaleSetIPConfigurationProperties{
					PrivateIPAddressVersion: ptr.To(armcompute.IPVersionIPv6),
					Primary:                 ptr.To(false),
					Subnet: &armcompute.APIEntityReference{
						ID: ptr.To(azure.SubnetID(s.SubscriptionID, s.VNetResourceGroup, s.VNetName, n.SubnetName)),
					},
				},
			}
			ipconfigs = append(ipconfigs, ipv6Config)
		}
		if i == 0 {
			ipconfigs[0].Properties.LoadBalancerBackendAddressPools = azure.PtrSlice(&backendAddressPools)
			nicConfig.Properties.Primary = ptr.To(true)
		}
		nicConfig.Properties.IPConfigurations = azure.PtrSlice(&ipconfigs)
		nicConfigs = append(nicConfigs, nicConfig)
	}
	return &nicConfigs
}

// generateStorageProfile generates a pointer to an armcompute.VirtualMachineScaleSetStorageProfile which can utilized for VM creation.
func (s *ScaleSetSpec) generateStorageProfile(ctx context.Context) (*armcompute.VirtualMachineScaleSetStorageProfile, error) {
	_, _, done := tele.StartSpanWithLogger(ctx, "scalesets.ScaleSetSpec.generateStorageProfile")
	defer done()

	storageProfile := &armcompute.VirtualMachineScaleSetStorageProfile{
		OSDisk: &armcompute.VirtualMachineScaleSetOSDisk{
			OSType:       ptr.To(armcompute.OperatingSystemTypes(s.OSDisk.OSType)),
			CreateOption: ptr.To(armcompute.DiskCreateOptionTypesFromImage),
			DiskSizeGB:   s.OSDisk.DiskSizeGB,
		},
	}

	// enable ephemeral OS
	if s.OSDisk.DiffDiskSettings != nil {
		if !s.SKU.HasCapability(resourceskus.EphemeralOSDisk) {
			return nil, fmt.Errorf("vm size %s does not support ephemeral os. select a different vm size or disable ephemeral os", s.Size)
		}

		storageProfile.OSDisk.DiffDiskSettings = &armcompute.DiffDiskSettings{
			Option: ptr.To(armcompute.DiffDiskOptions(s.OSDisk.DiffDiskSettings.Option)),
		}
	}

	if s.OSDisk.ManagedDisk != nil {
		storageProfile.OSDisk.ManagedDisk = &armcompute.VirtualMachineScaleSetManagedDiskParameters{}
		if s.OSDisk.ManagedDisk.StorageAccountType != "" {
			storageProfile.OSDisk.ManagedDisk.StorageAccountType = ptr.To(armcompute.StorageAccountTypes(s.OSDisk.ManagedDisk.StorageAccountType))
		}
		if s.OSDisk.ManagedDisk.DiskEncryptionSet != nil {
			storageProfile.OSDisk.ManagedDisk.DiskEncryptionSet = &armcompute.DiskEncryptionSetParameters{ID: ptr.To(s.OSDisk.ManagedDisk.DiskEncryptionSet.ID)}
		}
	}

	if s.OSDisk.CachingType != "" {
		storageProfile.OSDisk.Caching = ptr.To(armcompute.CachingTypes(s.OSDisk.CachingType))
	}

	dataDisks := make([]armcompute.VirtualMachineScaleSetDataDisk, len(s.DataDisks))
	for i, disk := range s.DataDisks {
		dataDisks[i] = armcompute.VirtualMachineScaleSetDataDisk{
			CreateOption: ptr.To(armcompute.DiskCreateOptionTypesEmpty),
			DiskSizeGB:   ptr.To[int32](disk.DiskSizeGB),
			Lun:          disk.Lun,
			Name:         ptr.To(azure.GenerateDataDiskName(s.Name, disk.NameSuffix)),
		}

		if disk.ManagedDisk != nil {
			dataDisks[i].ManagedDisk = &armcompute.VirtualMachineScaleSetManagedDiskParameters{
				StorageAccountType: ptr.To(armcompute.StorageAccountTypes(disk.ManagedDisk.StorageAccountType)),
			}

			if disk.ManagedDisk.DiskEncryptionSet != nil {
				dataDisks[i].ManagedDisk.DiskEncryptionSet = &armcompute.DiskEncryptionSetParameters{ID: ptr.To(disk.ManagedDisk.DiskEncryptionSet.ID)}
			}
		}
	}
	storageProfile.DataDisks = azure.PtrSlice(&dataDisks)

	if s.VMImage == nil {
		return nil, errors.Errorf("vm image is nil")
	}
	imageRef, err := converters.ImageToSDK(s.VMImage)
	if err != nil {
		return nil, err
	}

	storageProfile.ImageReference = imageRef

	return storageProfile, nil
}

func (s *ScaleSetSpec) generateOSProfile(_ context.Context) (*armcompute.VirtualMachineScaleSetOSProfile, error) {
	sshKey, err := base64.StdEncoding.DecodeString(s.SSHKeyData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode ssh public key")
	}

	osProfile := &armcompute.VirtualMachineScaleSetOSProfile{
		ComputerNamePrefix: ptr.To(s.Name),
		AdminUsername:      ptr.To(azure.DefaultUserName),
		CustomData:         ptr.To(s.BootstrapData),
	}

	switch s.OSDisk.OSType {
	case string(armcompute.OperatingSystemTypesWindows):
		// Cloudbase-init is used to generate a password.
		// https://cloudbase-init.readthedocs.io/en/latest/plugins.html#setting-password-main
		//
		// We generate a random password here in case of failure
		// but the password on the VM will NOT be the same as created here.
		// Access is provided via SSH public key that is set during deployment
		// Azure also provides a way to reset user passwords in the case of need.
		osProfile.AdminPassword = ptr.To(generators.SudoRandomPassword(123))
		osProfile.WindowsConfiguration = &armcompute.WindowsConfiguration{
			EnableAutomaticUpdates: ptr.To(false),
		}
	default:
		osProfile.LinuxConfiguration = &armcompute.LinuxConfiguration{
			DisablePasswordAuthentication: ptr.To(true),
			SSH: &armcompute.SSHConfiguration{
				PublicKeys: []*armcompute.SSHPublicKey{
					{
						Path:    ptr.To(fmt.Sprintf("/home/%s/.ssh/authorized_keys", azure.DefaultUserName)),
						KeyData: ptr.To(string(sshKey)),
					},
				},
			},
		}
	}

	return osProfile, nil
}

func (s *ScaleSetSpec) generateImagePlan(ctx context.Context) *armcompute.Plan {
	_, log, done := tele.StartSpanWithLogger(ctx, "scalesets.ScaleSetSpec.generateImagePlan")
	defer done()

	if s.VMImage == nil {
		log.V(2).Info("no vm image found, disabling plan")
		return nil
	}

	if s.VMImage.SharedGallery != nil && s.VMImage.SharedGallery.Publisher != nil && s.VMImage.SharedGallery.SKU != nil && s.VMImage.SharedGallery.Offer != nil {
		return &armcompute.Plan{
			Publisher: s.VMImage.SharedGallery.Publisher,
			Name:      s.VMImage.SharedGallery.SKU,
			Product:   s.VMImage.SharedGallery.Offer,
		}
	}

	if s.VMImage.Marketplace == nil || !s.VMImage.Marketplace.ThirdPartyImage {
		return nil
	}

	if s.VMImage.Marketplace.Publisher == "" || s.VMImage.Marketplace.SKU == "" || s.VMImage.Marketplace.Offer == "" {
		return nil
	}

	return &armcompute.Plan{
		Publisher: ptr.To(s.VMImage.Marketplace.Publisher),
		Name:      ptr.To(s.VMImage.Marketplace.SKU),
		Product:   ptr.To(s.VMImage.Marketplace.Offer),
	}
}

func (s *ScaleSetSpec) getSecurityProfile() (*armcompute.SecurityProfile, error) {
	if s.SecurityProfile == nil {
		return nil, nil
	}

	if !s.SKU.HasCapability(resourceskus.EncryptionAtHost) {
		return nil, azure.WithTerminalError(errors.Errorf("encryption at host is not supported for VM type %s", s.Size))
	}

	return &armcompute.SecurityProfile{
		EncryptionAtHost: ptr.To(*s.SecurityProfile.EncryptionAtHost),
	}, nil
}
