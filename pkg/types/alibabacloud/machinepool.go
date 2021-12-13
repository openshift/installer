package alibabacloud

// DiskCategory is the category of the ECS disk. Supported disk category:
// cloud_essd(ESSD disk), cloud_efficiency(ultra disk).
//
// +kubebuilder:validation:Enum="";cloud_efficiency;cloud_essd
type DiskCategory string

const (
	// CloudEfficiencyDiskCategory is the disk category for ultra disk.
	CloudEfficiencyDiskCategory DiskCategory = "cloud_efficiency"

	// CloudESSDDiskCategory is the disk category for ESSD.
	CloudESSDDiskCategory DiskCategory = "cloud_essd"
)

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
const DefaultDiskCategory DiskCategory = CloudESSDDiskCategory

// DefaultSystemDiskSize holds the default Alibaba Cloud system disk size used by the ECS.
const DefaultSystemDiskSize int = 120

// DefaultMasterInstanceType holds the default Alibaba Cloud instacne type used by the master.
const DefaultMasterInstanceType string = "ecs.g6.xlarge"

// DefaultWorkerInstanceType holds the default Alibaba Cloud instacne type used by the worker.
const DefaultWorkerInstanceType string = "ecs.g6.large"

// DefaultMasterMachinePoolPlatform returns the default machine pool configuration used by the master
func DefaultMasterMachinePoolPlatform() MachinePool {
	return MachinePool{
		InstanceType:       DefaultMasterInstanceType,
		SystemDiskCategory: DefaultDiskCategory,
		SystemDiskSize:     DefaultSystemDiskSize,
	}
}

// DefaultWorkerMachinePoolPlatform returns the default machine pool configuration used by the worker
func DefaultWorkerMachinePoolPlatform() MachinePool {
	return MachinePool{
		InstanceType:       DefaultWorkerInstanceType,
		SystemDiskCategory: DefaultDiskCategory,
		SystemDiskSize:     DefaultSystemDiskSize,
	}
}
