package alibabacloud

// DiskCategory is the category of the ECS disk. Supported disk category:
// cloud_essd(ESSD disk), cloud_efficiency(ultra disk).
//
// +kubebuilder:validation:Enum="";cloud_efficiency;cloud_essd
type DiskCategory string

// MachinePool stores the configuration for a machine pool installed
// on Alibaba Cloud.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	// eg. ["cn-hangzhou-i", "cn-hangzhou-h", "cn-hangzhou-j"]
	//
	// +optional
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the ECS instance type.
	// eg. ecs.g6.large
	//
	// +optional
	InstanceType string `json:"instanceType,omitempty"`

	// SystemDiskCategory defines the category of the system disk.
	//
	// +optional
	SystemDiskCategory DiskCategory `json:"systemDiskCategory,omitempty"`

	// SystemDiskSize defines the size of the system disk in gibibytes (GiB).
	//
	// +kubebuilder:validation:Type=integer
	// +kubebuilder:validation:Minimum=120
	// +optional
	SystemDiskSize int `json:"systemDiskSize,omitempty"`

	// ImageID is the Image ID that should be used to create ECS instance.
	// If set, the ImageID should belong to the same region as the cluster.
	//
	// +optional
	ImageID string `json:"imageID,omitempty"`
}
