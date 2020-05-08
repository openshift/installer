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

	// OSDisk defines the storage for instance.
	//
	// +optional
	OSDisk `json:"osDisk"`
}

// OSDisk defines the disk for machines on Azure.
type OSDisk struct {
	// DiskSizeGB defines the size of disk in GB.
	//
	// +kubebuilder:validation:Minimum=0
	DiskSizeGB int32 `json:"diskSizeGB"`
	// DiskType defines the type of disk.
	// The valid values are Standard_LRS, Premium_LRS, StandardSSD_LRS
	// For control plane nodes, the valid values are Premium_LRS and StandardSSD_LRS.
	// Default is Premium_LRS
	// +optional
	// +kubebuilder:validation:Enum=Standard_LRS;Premium_LRS;StandardSSD_LRS
	DiskType string `json:"diskType"`
}

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

	if required.OSDisk.DiskSizeGB != 0 {
		a.OSDisk.DiskSizeGB = required.OSDisk.DiskSizeGB
	}

	if required.OSDisk.DiskType != "" {
		a.OSDisk.DiskType = required.OSDisk.DiskType
	}
}
