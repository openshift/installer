package azure

// SecurityTypes represents the SecurityType of the virtual machine.
type SecurityTypes string

const (
	// SecurityTypesConfidentialVM defines the SecurityType of the virtual machine as a Confidential VM.
	SecurityTypesConfidentialVM SecurityTypes = "ConfidentialVM"
	// SecurityTypesTrustedLaunch defines the SecurityType of the virtual machine as a Trusted Launch VM.
	SecurityTypesTrustedLaunch SecurityTypes = "TrustedLaunch"
)

// MachinePool stores the configuration for a machine pool installed
// on Azure.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	// eg. ["1", "2", "3"]
	//
	// +optional
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the azure instance type.
	// eg. Standard_DS_V2
	//
	// +optional
	InstanceType string `json:"type"`

	// EncryptionAtHost enables encryption at the VM host.
	//
	// +optional
	EncryptionAtHost bool `json:"encryptionAtHost,omitempty"`

	// OSDisk defines the storage for instance.
	//
	// +optional
	OSDisk `json:"osDisk"`

	// ultraSSDCapability defines if the instance should use Ultra SSD disks.
	//
	// +optional
	// +kubebuilder:validation:Enum=Enabled;Disabled
	UltraSSDCapability string `json:"ultraSSDCapability,omitempty"`

	// VMNetworkingType specifies whether to enable accelerated networking.
	// Accelerated networking enables single root I/O virtualization (SR-IOV) to a VM, greatly improving its
	// networking performance.
	// eg. values: "Accelerated", "Basic"
	//
	// +kubebuilder:validation:Enum="Accelerated"; "Basic"
	// +optional
	VMNetworkingType string `json:"vmNetworkingType,omitempty"`

	// OSImage defines the image to use for the OS.
	// +optional
	OSImage OSImage `json:"osImage,omitempty"`

	// Settings specify the security type and the UEFI settings of the virtual machine. This field can
	// be set for Confidential VMs and Trusted Launch for VMs.
	// +optional
	Settings *SecuritySettings `json:"settings,omitempty"`
}

// SecuritySettings define the security type and the UEFI settings of the virtual machine.
type SecuritySettings struct {
	// SecurityType specifies the SecurityType of the virtual machine. It has to be set to any specified value to
	// enable secure boot and vTPM. The default behavior is: secure boot and vTPM will not be enabled unless this property is set.
	// +kubebuilder:validation:Enum=ConfidentialVM;TrustedLaunch
	// +kubebuilder:validation:Required
	SecurityType SecurityTypes `json:"securityType,omitempty"`

	// ConfidentialVM specifies the security configuration of the virtual machine.
	// For more information regarding Confidential VMs, please refer to:
	// https://learn.microsoft.com/azure/confidential-computing/confidential-vm-overview
	// +optional
	ConfidentialVM *ConfidentialVM `json:"confidentialVM,omitempty"`

	// TrustedLaunch specifies the security configuration of the virtual machine.
	// For more information regarding TrustedLaunch for VMs, please refer to:
	// https://learn.microsoft.com/azure/virtual-machines/trusted-launch
	// +optional
	TrustedLaunch *TrustedLaunch `json:"trustedLaunch,omitempty"`
}

// ConfidentialVM defines the UEFI settings for the virtual machine.
type ConfidentialVM struct {
	// UEFISettings specifies the security settings like secure boot and vTPM used while creating the virtual machine.
	// +kubebuilder:validation:Required
	UEFISettings *UEFISettings `json:"uefiSettings,omitempty"`
}

// TrustedLaunch defines the UEFI settings for the virtual machine.
type TrustedLaunch struct {
	// UEFISettings specifies the security settings like secure boot and vTPM used while creating the virtual machine.
	// +kubebuilder:validation:Required
	UEFISettings *UEFISettings `json:"uefiSettings,omitempty"`
}

// UEFISettings specifies the security settings like secure boot and vTPM used while creating the
// virtual machine.
type UEFISettings struct {
	// SecureBoot specifies whether secure boot should be enabled on the virtual machine.
	// Secure Boot verifies the digital signature of all boot components and halts the boot process if
	// signature verification fails.
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is disabled.
	// +kubebuilder:validation:Enum=Enabled;Disabled
	// +optional
	SecureBoot *string `json:"secureBoot,omitempty"`

	// VirtualizedTrustedPlatformModule specifies whether vTPM should be enabled on the virtual machine.
	// When enabled the virtualized trusted platform module measurements are used to create a known good boot integrity policy baseline.
	// The integrity policy baseline is used for comparison with measurements from subsequent VM boots to determine if anything has changed.
	// This is required to be set to enabled if the SecurityEncryptionType is defined.
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is disabled.
	// +kubebuilder:validation:Enum=Enabled;Disabled
	// +optional
	VirtualizedTrustedPlatformModule *string `json:"virtualizedTrustedPlatformModule,omitempty"`
}

// VMNetworkingCapability defines the states for accelerated networking feature
type VMNetworkingCapability string

const (
	// AcceleratedNetworkingEnabled is string representation of the VMNetworkingType / AcceleratedNetworking Capability
	// provided by the Azure API
	AcceleratedNetworkingEnabled = "AcceleratedNetworkingEnabled"

	// VMNetworkingTypeBasic enum attribute that is the default setting which means AcceleratedNetworking is disabled.
	VMNetworkingTypeBasic VMNetworkingCapability = "Basic"

	// VMnetworkingTypeAccelerated enum attribute that enables AcceleratedNetworking on a VM NIC.
	VMnetworkingTypeAccelerated VMNetworkingCapability = "Accelerated"
)

// Set sets the values from `required` to `a`.
func (a *MachinePool) Set(required *MachinePool) {
	if required == nil || a == nil {
		return
	}

	if len(required.Zones) > 0 {
		a.Zones = required.Zones
	}

	if required.InstanceType != "" {
		a.InstanceType = required.InstanceType
	}

	if required.EncryptionAtHost {
		a.EncryptionAtHost = required.EncryptionAtHost
	}

	if required.OSDisk.DiskSizeGB != 0 {
		a.OSDisk.DiskSizeGB = required.OSDisk.DiskSizeGB
	}

	if required.OSDisk.DiskType != "" {
		a.OSDisk.DiskType = required.OSDisk.DiskType
	}

	if required.DiskEncryptionSet != nil {
		a.DiskEncryptionSet = required.DiskEncryptionSet
	}

	if required.UltraSSDCapability != "" {
		a.UltraSSDCapability = required.UltraSSDCapability
	}

	if required.VMNetworkingType != "" {
		a.VMNetworkingType = required.VMNetworkingType
	}

	var emptyOSImage OSImage
	if required.OSImage != emptyOSImage {
		a.OSImage = required.OSImage
	}

	if required.OSDisk.SecurityProfile != nil {
		a.OSDisk.SecurityProfile = required.OSDisk.SecurityProfile
	}

	if required.Settings != nil {
		a.Settings = required.Settings
	}
}

// ImagePurchasePlan defines the purchase plan of a Marketplace image.
// +kubebuilder:validation:Enum=WithPurchasePlan;NoPurchasePlan
type ImagePurchasePlan string

const (
	// ImageWithPurchasePlan enum attribute which is the default setting.
	ImageWithPurchasePlan ImagePurchasePlan = "WithPurchasePlan"
	// ImageNoPurchasePlan  enum attribute which speficies the image does not need a purchase plan.
	ImageNoPurchasePlan ImagePurchasePlan = "NoPurchasePlan"
)

// OSImage is the image to use for the OS of a machine.
type OSImage struct {
	// Plan is the purchase plan of the image.
	// If omitted, it defaults to "WithPurchasePlan".
	// +optional
	Plan ImagePurchasePlan `json:"plan"`
	// Publisher is the publisher of the image.
	Publisher string `json:"publisher"`
	// Offer is the offer of the image.
	Offer string `json:"offer"`
	// SKU is the SKU of the image.
	SKU string `json:"sku"`
	// Version is the version of the image.
	Version string `json:"version"`
}
