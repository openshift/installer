package azure

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
	// The valid values are Enabled, Disabled.
	//
	// +optional
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
}

// OSImage is the image to use for the OS of a machine.
type OSImage struct {
	// Publisher is the publisher of the image.
	Publisher string `json:"publisher"`
	// Offer is the offer of the image.
	Offer string `json:"offer"`
	// SKU is the SKU of the image.
	SKU string `json:"sku"`
	// Version is the version of the image.
	Version string `json:"version"`
}
