package gcp

// MachinePool stores the configuration for a machine pool installed on GCP.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	//
	// +optional
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the GCP instance type.
	// eg. n1-standard-4
	//
	// +optional
	InstanceType string `json:"type"`

	// OSDisk defines the storage for instance.
	//
	// +optional
	OSDisk `json:"osDisk"`
}

// OSDisk defines the disk for machines on GCP.
type OSDisk struct {
	// DiskType defines the type of disk.
	// The valid values are pd-standard and pd-ssd
	// For control plane nodes, the valid value is pd-ssd.
	// +optional
	// +kubebuilder:validation:Enum=pd-ssd;pd-standard
	DiskType string `json:"DiskType"`

	// DiskSizeGB defines the size of disk in GB.
	//
	// +kubebuilder:validation:Minimum=16
	// +kubebuilder:validation:Maximum=65536
	DiskSizeGB int64 `json:"DiskSizeGB"`
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

	if required.OSDisk.DiskSizeGB > 0 {
		a.OSDisk.DiskSizeGB = required.OSDisk.DiskSizeGB
	}

	if required.OSDisk.DiskType != "" {
		a.OSDisk.DiskType = required.OSDisk.DiskType
	}
}
