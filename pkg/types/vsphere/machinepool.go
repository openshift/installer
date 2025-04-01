package vsphere

// MachinePool stores the configuration for a machine pool installed
// on vSphere.
type MachinePool struct {
	// NumCPUs is the total number of virtual processor cores to assign a vm.
	//
	// +optional
	NumCPUs int32 `json:"cpus"`

	// NumCoresPerSocket is the number of cores per socket in a vm. The number
	// of vCPUs on the vm will be NumCPUs/NumCoresPerSocket.
	//
	// +optional
	NumCoresPerSocket int32 `json:"coresPerSocket"`

	// Memory is the size of a VM's memory in MB.
	//
	// +optional
	MemoryMiB int64 `json:"memoryMB"`

	// OSDisk defines the storage for instance.
	//
	// +optional
	OSDisk `json:"osDisk"`

	// DataDisks are additional disks to add to the VM that are not part of the VM's OVA template.
	//
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=29
	DataDisks []DataDisk `json:"dataDisks"`

	// Zones defines available zones
	// Zones is available in TechPreview.
	//
	// +omitempty
	Zones []string `json:"zones,omitempty"`
}

// OSDisk defines the disk for a virtual machine.
type OSDisk struct {
	// DiskSizeGB defines the size of disk in GB.
	//
	// +optional
	DiskSizeGB int32 `json:"diskSizeGB"`
}

// DataDisk defines a data disk to add to the VM that is not part of the VM OVA template.
type DataDisk struct {
	// name is used to identify the disk definition. name is required needs to be unique so that it can be used to
	// clearly identify purpose of the disk.
	// +kubebuilder:example=images_1
	// +kubebuilder:validation:MaxLength=80
	// +kubebuilder:validation:Pattern="^[a-zA-Z0-9]([-_a-zA-Z0-9]*[a-zA-Z0-9])?$"
	// +required
	Name string `json:"name"`
	// sizeGiB is the size of the disk in GiB.
	// The maximum supported size is 16384 GiB.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=16384
	// +required
	SizeGiB int32 `json:"sizeGiB"`
	// provisioningMode is an optional field that specifies the provisioning type to be used by this vSphere data disk.
	// Allowed values are "Thin", "Thick", "EagerlyZeroed", and omitted.
	// When set to Thin, the disk will be made using thin provisioning allocating the bare minimum space.
	// When set to Thick, the full disk size will be allocated when disk is created.
	// When set to EagerlyZeroed, the disk will be created using eager zero provisioning. An eager zeroed thick disk has all space allocated and wiped clean of any previous contents on the physical media at creation time. Such disks may take longer time during creation compared to other disk formats.
	// When omitted, no setting will be applied to the data disk and the provisioning mode for the disk will be determined by the default storage policy configured for the datastore in vSphere.
	// +optional
	ProvisioningMode ProvisioningMode `json:"provisioningMode,omitempty"`
}

// ProvisioningMode represents the various provisioning types available to a VMs disk.
// +kubebuilder:validation:Enum=Thin;Thick;EagerlyZeroed
type ProvisioningMode string

const (
	// ProvisioningModeThin creates the disk using thin provisioning. This means a sparse (allocate on demand)
	// format with additional space optimizations.
	ProvisioningModeThin ProvisioningMode = "Thin"

	// ProvisioningModeThick creates the disk with all space allocated.
	ProvisioningModeThick ProvisioningMode = "Thick"

	// ProvisioningModeEagerlyZeroed creates the disk using eager zero provisioning. An eager zeroed thick disk
	// has all space allocated and wiped clean of any previous contents on the physical media at
	// creation time. Such disks may take longer time during creation compared to other disk formats.
	ProvisioningModeEagerlyZeroed ProvisioningMode = "EagerlyZeroed"
)

// Set sets the values from `required` to `p`.
func (p *MachinePool) Set(required *MachinePool) {
	if required == nil || p == nil {
		return
	}

	if required.NumCPUs != 0 {
		p.NumCPUs = required.NumCPUs
	}

	if required.NumCoresPerSocket != 0 {
		p.NumCoresPerSocket = required.NumCoresPerSocket
	}

	if required.MemoryMiB != 0 {
		p.MemoryMiB = required.MemoryMiB
	}

	if required.OSDisk.DiskSizeGB != 0 {
		p.OSDisk.DiskSizeGB = required.OSDisk.DiskSizeGB
	}

	if len(required.Zones) > 0 {
		p.Zones = required.Zones
	}

	if len(required.DataDisks) > 0 {
		p.DataDisks = required.DataDisks
	}
}
