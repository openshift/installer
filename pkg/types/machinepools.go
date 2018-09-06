package types

// MachinePool is a pool of machines to be installed.
type MachinePool struct {
	// Name is the name of the machine pool.
	Name string `json:"name"`

	// Replicas is the count of machines for this machine pool.
	// Default is 1.
	Replicas *int64 `json:"replicas"`

	// PlatformConfig is configuration for machine pool specific to the platfrom.
	PlatformConfig MachinePoolPlatformConfig `json:"platformConfig"`
}

// MachinePoolPlatformConfig is the platform-specific configuration for a machine
// pool. Only one of the platforms should be set.
type MachinePoolPlatformConfig struct {
	// AWS is the configuration used when installing on AWS.
	AWS *AWSMachinePoolPlatformConfig `json:"aws,omitempty"`
	// Libvirt is the configuration used when installing on libvirt.
	Libvirt *LibvirtMachinePoolPlatformConfig `json:"libvirt,omitempty"`
}

// AWSMachinePoolPlatformConfig stores the configuration for a machine pool
// installed on AWS.
type AWSMachinePoolPlatformConfig struct {
	// InstanceType defines the ec2 instance type.
	// eg. m4-large
	InstanceType string `json:"type"`

	// IAMRoleName defines the IAM role associated
	// with the ec2 instance.
	IAMRoleName string `json:"iamRoleName"`

	// EC2RootVolume defines the storage for ec2 instance.
	EC2RootVolume `json:"rootVolume"`
}

// EC2RootVolume defines the storage for an ec2 instance.
type EC2RootVolume struct {
	// IOPS defines the iops for the instance.
	IOPS int `json:"iops"`
	// Size defines the size of the instance.
	Size int `json:"size"`
	// Type defines the type of the instance.
	Type string `json:"type"`
}

// LibvirtMachinePoolPlatformConfig stores the configuration for a machine pool
// installed on libvirt.
type LibvirtMachinePoolPlatformConfig struct {
	// QCOWImagePath
	QCOWImagePath string `json:"qcowImagePath"`
}
