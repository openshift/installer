package types

// MachinePool is a pool of machines to be installed.
type MachinePool struct {
	// Name is the name of the machine pool.
	Name string `json:"name" yaml:"name"`

	// Replicas is the count of machines for this machine pool.
	// Default is 1.
	Replicas *int64 `json:"replicas,omitempty" yaml:"replicas,omitempty"`

	// Platform is configuration for machine pool specific to the platfrom.
	Platform MachinePoolPlatform `json:"platform" yaml:"platform"`
}

// MachinePoolPlatform is the platform-specific configuration for a machine
// pool. Only one of the platforms should be set.
type MachinePoolPlatform struct {
	// AWS is the configuration used when installing on AWS.
	AWS *AWSMachinePoolPlatform `json:"aws,omitempty" yaml:"aws,omitempty"`

	// Libvirt is the configuration used when installing on libvirt.
	Libvirt *LibvirtMachinePoolPlatform `json:"libvirt,omitempty" yaml:"libvirt,omitempty"`
}

// AWSMachinePoolPlatform stores the configuration for a machine pool
// installed on AWS.
type AWSMachinePoolPlatform struct {
	// InstanceType defines the ec2 instance type.
	// eg. m4-large
	InstanceType string `json:"type,omitempty" yaml:"type,omitempty"`

	// IAMRoleName defines the IAM role associated
	// with the ec2 instance.
	IAMRoleName string `json:"iamRoleName,omitempty" yaml:"iamRoleName,omitempty"`

	// EC2RootVolume defines the storage for ec2 instance.
	EC2RootVolume `json:"rootVolume" yaml:"rootVolume"`
}

// EC2RootVolume defines the storage for an ec2 instance.
type EC2RootVolume struct {
	// IOPS defines the iops for the instance.
	IOPS int `json:"iops,omitempty" yaml:"iops,omitempty"`
	// Size defines the size of the instance.
	Size int `json:"size,omitempty" yaml:"size,omitempty"`
	// Type defines the type of the instance.
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

// LibvirtMachinePoolPlatform stores the configuration for a machine pool
// installed on libvirt.
type LibvirtMachinePoolPlatform struct {
	// QCOWImagePath
	QCOWImagePath string `json:"qcowImagePath,omitempty" yaml:"qcowImagePath,omitempty"`
}
