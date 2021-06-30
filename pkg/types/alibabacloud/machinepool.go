package alibabacloud

// MachinePool stores the configuration for a machine pool installed
// on Alibaba Cloud.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	// eg. ["cn-hangzhou-i", "cn-hangzhou-h", "cn-hangzhou-j"]
	// +optional
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the ECS instance type.
	// eg. ecs.xn4.small
	//
	// +optional
	InstanceType string `json:"instanceType"`

	// SystemDisk defines the system disk for ECS instances in the machine pool.
	//
	// +optional
	SystemDisk `json:"systemDisk"`
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

	if required.SystemDisk.Size != 0 {
		a.SystemDisk.Size = required.SystemDisk.Size
	}
	if required.SystemDisk.Category != "" {
		a.SystemDisk.Category = required.SystemDisk.Category
	}
}

// DefaultDiskCategory holds the default Alibaba Cloud disk type used by the ECS.
const DefaultDiskCategory string = "cloud_essd"

// SystemDisk defines the storage for an ecs instance.
type SystemDisk struct {
	// Size defines the size of the disk in gibibytes (GiB).
	//
	// +kubebuilder:validation:Minimum=20
	Size int `json:"size"`

	// Category defines the category of the disk.
	Category string `json:"category"`
}
