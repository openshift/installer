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

package virtualmachines

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/resourceskus"
	"sigs.k8s.io/cluster-api-provider-azure/util/generators"
)

// VMSpec defines the specification for a Virtual Machine.
type VMSpec struct {
	Name                       string
	ResourceGroup              string
	Location                   string
	ExtendedLocation           *infrav1.ExtendedLocationSpec
	ClusterName                string
	Role                       string
	NICIDs                     []string
	SSHKeyData                 string
	Size                       string
	AvailabilitySetID          string
	Zone                       string
	Identity                   infrav1.VMIdentity
	OSDisk                     infrav1.OSDisk
	DataDisks                  []infrav1.DataDisk
	UserAssignedIdentities     []infrav1.UserAssignedIdentity
	SpotVMOptions              *infrav1.SpotVMOptions
	SecurityProfile            *infrav1.SecurityProfile
	AdditionalTags             infrav1.Tags
	AdditionalCapabilities     *infrav1.AdditionalCapabilities
	DiagnosticsProfile         *infrav1.Diagnostics
	DisableExtensionOperations bool
	CapacityReservationGroupID string
	SKU                        resourceskus.SKU
	Image                      *infrav1.Image
	BootstrapData              string
	ProviderID                 string
}

// ResourceName returns the name of the virtual machine.
func (s *VMSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the virtual machine.
func (s *VMSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for virtual machines.
func (s *VMSpec) OwnerResourceName() string {
	return ""
}

// Parameters returns the parameters for the virtual machine.
func (s *VMSpec) Parameters(ctx context.Context, existing interface{}) (params interface{}, err error) {
	if existing != nil {
		if _, ok := existing.(armcompute.VirtualMachine); !ok {
			return nil, errors.Errorf("%T is not an armcompute.VirtualMachine", existing)
		}
		// vm already exists
		return nil, nil
	}

	// VM got deleted outside of capz, do not recreate it as Machines are immutable.
	if s.ProviderID != "" {
		return nil, azure.VMDeletedError{ProviderID: s.ProviderID}
	}

	storageProfile, err := s.generateStorageProfile()
	if err != nil {
		return nil, err
	}

	securityProfile, err := s.generateSecurityProfile(storageProfile)
	if err != nil {
		return nil, err
	}

	osProfile, err := s.generateOSProfile()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate OS Profile")
	}

	priority, evictionPolicy, billingProfile, err := converters.GetSpotVMOptions(s.SpotVMOptions, s.OSDisk.DiffDiskSettings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Spot VM options")
	}

	identity, err := converters.VMIdentityToVMSDK(s.Identity, s.UserAssignedIdentities)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate VM identity")
	}

	return armcompute.VirtualMachine{
		Plan:             converters.ImageToPlan(s.Image),
		Location:         ptr.To(s.Location),
		ExtendedLocation: converters.ExtendedLocationToComputeSDK(s.ExtendedLocation),
		Tags: converters.TagsToMap(infrav1.Build(infrav1.BuildParams{
			ClusterName: s.ClusterName,
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Name:        ptr.To(s.Name),
			Role:        ptr.To(s.Role),
			Additional:  s.AdditionalTags,
		})),
		Properties: &armcompute.VirtualMachineProperties{
			AdditionalCapabilities: s.generateAdditionalCapabilities(),
			AvailabilitySet:        s.getAvailabilitySet(),
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: ptr.To(armcompute.VirtualMachineSizeTypes(s.Size)),
			},
			StorageProfile:  storageProfile,
			SecurityProfile: securityProfile,
			OSProfile:       osProfile,
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: s.generateNICRefs(),
			},
			Priority:            priority,
			EvictionPolicy:      evictionPolicy,
			BillingProfile:      billingProfile,
			DiagnosticsProfile:  converters.GetDiagnosticsProfile(s.DiagnosticsProfile),
			CapacityReservation: s.getCapacityReservationProfile(),
		},
		Identity: identity,
		Zones:    s.getZones(),
	}, nil
}

// generateStorageProfile generates a pointer to an armcompute.StorageProfile which can utilized for VM creation.
func (s *VMSpec) generateStorageProfile() (*armcompute.StorageProfile, error) {
	osDisk := &armcompute.OSDisk{
		Name:         ptr.To(azure.GenerateOSDiskName(s.Name)),
		OSType:       ptr.To(armcompute.OperatingSystemTypes(s.OSDisk.OSType)),
		CreateOption: ptr.To(armcompute.DiskCreateOptionTypesFromImage),
		DiskSizeGB:   s.OSDisk.DiskSizeGB,
	}
	if s.OSDisk.CachingType != "" {
		osDisk.Caching = ptr.To(armcompute.CachingTypes(s.OSDisk.CachingType))
	}
	storageProfile := &armcompute.StorageProfile{
		OSDisk: osDisk,
	}

	// Checking if the requested VM size has at least 2 vCPUS
	vCPUCapability, err := s.SKU.HasCapabilityWithCapacity(resourceskus.VCPUs, resourceskus.MinimumVCPUS)
	if err != nil {
		return nil, azure.WithTerminalError(errors.Wrap(err, "failed to validate the vCPU capability"))
	}
	if !vCPUCapability {
		return nil, azure.WithTerminalError(errors.New("VM size should be bigger or equal to at least 2 vCPUs"))
	}

	// Checking if the requested VM size has at least 2 Gi of memory
	MemoryCapability, err := s.SKU.HasCapabilityWithCapacity(resourceskus.MemoryGB, resourceskus.MinimumMemory)
	if err != nil {
		return nil, azure.WithTerminalError(errors.Wrap(err, "failed to validate the memory capability"))
	}

	if !MemoryCapability {
		return nil, azure.WithTerminalError(errors.New("VM memory should be bigger or equal to at least 2Gi"))
	}
	// enable ephemeral OS
	if s.OSDisk.DiffDiskSettings != nil {
		if !s.SKU.HasCapability(resourceskus.EphemeralOSDisk) {
			return nil, azure.WithTerminalError(fmt.Errorf("VM size %s does not support ephemeral os. Select a different VM size or disable ephemeral os", s.Size))
		}

		storageProfile.OSDisk.DiffDiskSettings = &armcompute.DiffDiskSettings{
			Option: ptr.To(armcompute.DiffDiskOptions(s.OSDisk.DiffDiskSettings.Option)),
		}
	}

	if s.OSDisk.ManagedDisk != nil {
		storageProfile.OSDisk.ManagedDisk = &armcompute.ManagedDiskParameters{}
		if s.OSDisk.ManagedDisk.StorageAccountType != "" {
			storageProfile.OSDisk.ManagedDisk.StorageAccountType = ptr.To(armcompute.StorageAccountTypes(s.OSDisk.ManagedDisk.StorageAccountType))
		}
		if s.OSDisk.ManagedDisk.DiskEncryptionSet != nil {
			storageProfile.OSDisk.ManagedDisk.DiskEncryptionSet = &armcompute.DiskEncryptionSetParameters{ID: ptr.To(s.OSDisk.ManagedDisk.DiskEncryptionSet.ID)}
		}
		if s.OSDisk.ManagedDisk.SecurityProfile != nil {
			if _, exists := s.SKU.GetCapability(resourceskus.ConfidentialComputingType); !exists {
				return nil, azure.WithTerminalError(fmt.Errorf("VM size %s does not support confidential computing. Select a different VM size or remove the security profile of the OS disk", s.Size))
			}

			storageProfile.OSDisk.ManagedDisk.SecurityProfile = &armcompute.VMDiskSecurityProfile{}

			if s.OSDisk.ManagedDisk.SecurityProfile.DiskEncryptionSet != nil {
				storageProfile.OSDisk.ManagedDisk.SecurityProfile.DiskEncryptionSet = &armcompute.DiskEncryptionSetParameters{ID: ptr.To(s.OSDisk.ManagedDisk.SecurityProfile.DiskEncryptionSet.ID)}
			}
			if s.OSDisk.ManagedDisk.SecurityProfile.SecurityEncryptionType != "" {
				storageProfile.OSDisk.ManagedDisk.SecurityProfile.SecurityEncryptionType = ptr.To(armcompute.SecurityEncryptionTypes(string(s.OSDisk.ManagedDisk.SecurityProfile.SecurityEncryptionType)))
			}
		}
	}

	dataDisks := make([]*armcompute.DataDisk, len(s.DataDisks))
	for i, disk := range s.DataDisks {
		dataDisks[i] = &armcompute.DataDisk{
			CreateOption: ptr.To(armcompute.DiskCreateOptionTypesEmpty),
			DiskSizeGB:   ptr.To[int32](disk.DiskSizeGB),
			Lun:          disk.Lun,
			Name:         ptr.To(azure.GenerateDataDiskName(s.Name, disk.NameSuffix)),
		}
		if disk.CachingType != "" {
			dataDisks[i].Caching = ptr.To(armcompute.CachingTypes(disk.CachingType))
		}

		if disk.ManagedDisk != nil {
			dataDisks[i].ManagedDisk = &armcompute.ManagedDiskParameters{
				StorageAccountType: ptr.To(armcompute.StorageAccountTypes(disk.ManagedDisk.StorageAccountType)),
			}

			if disk.ManagedDisk.DiskEncryptionSet != nil {
				dataDisks[i].ManagedDisk.DiskEncryptionSet = &armcompute.DiskEncryptionSetParameters{ID: ptr.To(disk.ManagedDisk.DiskEncryptionSet.ID)}
			}

			// check the support for ultra disks based on location and vm size
			if disk.ManagedDisk.StorageAccountType == string(armcompute.StorageAccountTypesUltraSSDLRS) && !s.SKU.HasLocationCapability(resourceskus.UltraSSDAvailable, s.Location, s.Zone) {
				return nil, azure.WithTerminalError(fmt.Errorf("VM size %s does not support ultra disks in location %s. Select a different VM size or disable ultra disks", s.Size, s.Location))
			}
		}
	}
	storageProfile.DataDisks = dataDisks

	imageRef, err := converters.ImageToSDK(s.Image)
	if err != nil {
		return nil, err
	}

	storageProfile.ImageReference = imageRef

	return storageProfile, nil
}

func (s *VMSpec) generateOSProfile() (*armcompute.OSProfile, error) {
	sshKey, err := base64.StdEncoding.DecodeString(s.SSHKeyData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode ssh public key")
	}

	osProfile := &armcompute.OSProfile{
		ComputerName:             ptr.To(s.Name),
		AdminUsername:            ptr.To(azure.DefaultUserName),
		CustomData:               ptr.To(s.BootstrapData),
		AllowExtensionOperations: ptr.To(!s.DisableExtensionOperations),
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

func (s *VMSpec) generateSecurityProfile(storageProfile *armcompute.StorageProfile) (*armcompute.SecurityProfile, error) {
	if s.SecurityProfile == nil {
		return nil, nil
	}

	securityProfile := &armcompute.SecurityProfile{}

	if storageProfile.OSDisk.ManagedDisk != nil &&
		storageProfile.OSDisk.ManagedDisk.SecurityProfile != nil &&
		ptr.Deref(storageProfile.OSDisk.ManagedDisk.SecurityProfile.SecurityEncryptionType, "") != "" {
		if s.SecurityProfile.EncryptionAtHost != nil && *s.SecurityProfile.EncryptionAtHost &&
			*storageProfile.OSDisk.ManagedDisk.SecurityProfile.SecurityEncryptionType == armcompute.SecurityEncryptionTypesDiskWithVMGuestState {
			return nil, azure.WithTerminalError(errors.Errorf("encryption at host is not supported when securityEncryptionType is set to %s", armcompute.SecurityEncryptionTypesDiskWithVMGuestState))
		}

		if s.SecurityProfile.SecurityType != infrav1.SecurityTypesConfidentialVM {
			return nil, azure.WithTerminalError(errors.Errorf("securityType should be set to %s when securityEncryptionType is set", infrav1.SecurityTypesConfidentialVM))
		}

		if s.SecurityProfile.UefiSettings == nil {
			return nil, azure.WithTerminalError(errors.New("vTpmEnabled should be true when securityEncryptionType is set"))
		}

		if ptr.Deref(storageProfile.OSDisk.ManagedDisk.SecurityProfile.SecurityEncryptionType, "") == armcompute.SecurityEncryptionTypesDiskWithVMGuestState &&
			!*s.SecurityProfile.UefiSettings.SecureBootEnabled {
			return nil, azure.WithTerminalError(errors.Errorf("secureBootEnabled should be true when securityEncryptionType is set to %s", armcompute.SecurityEncryptionTypesDiskWithVMGuestState))
		}

		if s.SecurityProfile.UefiSettings.VTpmEnabled != nil && !*s.SecurityProfile.UefiSettings.VTpmEnabled {
			return nil, azure.WithTerminalError(errors.New("vTpmEnabled should be true when securityEncryptionType is set"))
		}

		securityProfile.SecurityType = ptr.To(armcompute.SecurityTypesConfidentialVM)

		securityProfile.UefiSettings = &armcompute.UefiSettings{
			SecureBootEnabled: s.SecurityProfile.UefiSettings.SecureBootEnabled,
			VTpmEnabled:       s.SecurityProfile.UefiSettings.VTpmEnabled,
		}

		return securityProfile, nil
	}

	if s.SecurityProfile.EncryptionAtHost != nil {
		if !s.SKU.HasCapability(resourceskus.EncryptionAtHost) && *s.SecurityProfile.EncryptionAtHost {
			return nil, azure.WithTerminalError(errors.Errorf("encryption at host is not supported for VM type %s", s.Size))
		}

		securityProfile.EncryptionAtHost = s.SecurityProfile.EncryptionAtHost
	}

	hasTrustedLaunchDisabled := s.SKU.HasCapability(resourceskus.TrustedLaunchDisabled)

	if s.SecurityProfile.UefiSettings != nil {
		securityProfile.UefiSettings = &armcompute.UefiSettings{}

		if s.SecurityProfile.UefiSettings.SecureBootEnabled != nil && *s.SecurityProfile.UefiSettings.SecureBootEnabled {
			if hasTrustedLaunchDisabled {
				return nil, azure.WithTerminalError(errors.Errorf("secure boot is not supported for VM type %s", s.Size))
			}

			if s.SecurityProfile.SecurityType != infrav1.SecurityTypesTrustedLaunch {
				return nil, azure.WithTerminalError(errors.Errorf("securityType should be set to %s when secureBootEnabled is true", infrav1.SecurityTypesTrustedLaunch))
			}

			securityProfile.SecurityType = ptr.To(armcompute.SecurityTypesTrustedLaunch)
			securityProfile.UefiSettings.SecureBootEnabled = ptr.To(true)
		}

		if s.SecurityProfile.UefiSettings.VTpmEnabled != nil && *s.SecurityProfile.UefiSettings.VTpmEnabled {
			if hasTrustedLaunchDisabled {
				return nil, azure.WithTerminalError(errors.Errorf("vTPM is not supported for VM type %s", s.Size))
			}

			if s.SecurityProfile.SecurityType != infrav1.SecurityTypesTrustedLaunch {
				return nil, azure.WithTerminalError(errors.Errorf("securityType should be set to %s when vTpmEnabled is true", infrav1.SecurityTypesTrustedLaunch))
			}

			securityProfile.SecurityType = ptr.To(armcompute.SecurityTypesTrustedLaunch)
			securityProfile.UefiSettings.VTpmEnabled = ptr.To(true)
		}
	}

	return securityProfile, nil
}

func (s *VMSpec) generateNICRefs() []*armcompute.NetworkInterfaceReference {
	nicRefs := make([]*armcompute.NetworkInterfaceReference, len(s.NICIDs))
	for i, id := range s.NICIDs {
		primary := i == 0
		nicRefs[i] = &armcompute.NetworkInterfaceReference{
			ID: ptr.To(id),
			Properties: &armcompute.NetworkInterfaceReferenceProperties{
				Primary: ptr.To(primary),
			},
		}
	}
	return nicRefs
}

func (s *VMSpec) generateAdditionalCapabilities() *armcompute.AdditionalCapabilities {
	var capabilities *armcompute.AdditionalCapabilities

	// Provisionally detect whether there is any Data Disk defined which uses UltraSSDs.
	// If that's the case, enable the UltraSSD capability.
	for _, dataDisk := range s.DataDisks {
		if dataDisk.ManagedDisk != nil && dataDisk.ManagedDisk.StorageAccountType == string(armcompute.StorageAccountTypesUltraSSDLRS) {
			capabilities = &armcompute.AdditionalCapabilities{
				UltraSSDEnabled: ptr.To(true),
			}
			break
		}
	}

	// Set Additional Capabilities if any is present on the spec.
	if s.AdditionalCapabilities != nil {
		if capabilities == nil {
			capabilities = &armcompute.AdditionalCapabilities{}
		}
		// Set UltraSSDEnabled if a specific value is set on the spec for it.
		if s.AdditionalCapabilities.UltraSSDEnabled != nil {
			capabilities.UltraSSDEnabled = s.AdditionalCapabilities.UltraSSDEnabled
		}
	}

	return capabilities
}

func (s *VMSpec) getAvailabilitySet() *armcompute.SubResource {
	var as *armcompute.SubResource
	if s.AvailabilitySetID != "" {
		as = &armcompute.SubResource{ID: &s.AvailabilitySetID}
	}
	return as
}

func (s *VMSpec) getZones() []*string {
	var zones []*string
	if s.Zone != "" {
		zones = []*string{ptr.To(s.Zone)}
	}
	return zones
}

func (s *VMSpec) getCapacityReservationProfile() *armcompute.CapacityReservationProfile {
	var crf *armcompute.CapacityReservationProfile
	if s.CapacityReservationGroupID != "" {
		crf = &armcompute.CapacityReservationProfile{
			CapacityReservationGroup: &armcompute.SubResource{ID: &s.CapacityReservationGroupID},
		}
	}
	return crf
}
