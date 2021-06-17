package alibabacloud

// MachinePool stores the configuration for a machine pool installed
// on Alibaba Cloud.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	// eg. ["cn-hangzhou-i", "cn-hangzhou-h", "cn-hangzhou-j"]
	// +optional
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the ECS instance type.
	// eg. ecs.g6.large
	//
	// +optional
	InstanceType string `json:"instanceType"`

	// SystemDiskCategory defines the category of the system disk.
	//
	// +kubebuilder:validation:Enum=cloud_efficiency;cloud_essd
	// +optional
	SystemDiskCategory string `json:"systemDiskCategory"`

	// SystemDiskSize defines the size of the system disk in gibibytes (GiB).
	//
	// +kubebuilder:validation:Minimum=120
	// +optional
	SystemDiskSize int `json:"systemDiskSize"`

	// ImageID is the Image ID that should be used to create ECS instance.
	// If set, the ImageID should belong to the same region as the cluster.
	//
	// +optional
	ImageID string `json:"imageID"`
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

	if required.SystemDiskSize != 0 {
		a.SystemDiskSize = required.SystemDiskSize
	}
	if required.SystemDiskCategory != "" {
		a.SystemDiskCategory = required.SystemDiskCategory
	}
	if required.ImageID != "" {
		a.ImageID = required.ImageID
	}
}

// DefaultDiskCategory holds the default Alibaba Cloud disk type used by the ECS.
const DefaultDiskCategory string = "cloud_essd"

// SystemDisk defines the storage for an ecs instance.
type SystemDisk struct {
	// Size defines the size of the disk in gibibytes (GiB).
	//
	// +kubebuilder:validation:Minimum=120
	Size int `json:"size"`

	// Category defines the category of the disk.
	Category string `json:"category"`
}
