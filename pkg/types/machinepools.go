package types

type MachinePools struct {
	// Name is the name of the machine pool.
	Name string

	// This is count of machines for this machine pool.
	// Default is 1.
	Replicas *int64 `json:"replicas"`

	// PlatformConfig is configuration for machine pool specfic to the platfrom.
	PlatformConfig MachinePoolPlatformConfig `json:"platformConfig"`
}

type MachinePoolPlatformConfig struct {
	AWS     *AWSMachinePoolPlatformConfig     `json:"aws,omitempty"`
	Libvirt *LibvirtMachinePoolPlatformConfig `json:"libvirt,omitempty"`
}

type AWSMachinePoolPlatformConfig struct {
	// InstanceType defines the ec2 instance type.
	// eg. m4-large
	InstanceType string `json:"type"`

	// IAMRoleName defines the IAM role associated
	// with the ec2 instance.
	IAMRoleName string `json:"iamRoleName"`

	// RootVolume defines the storage for ec2 instance.
	EC2RootVolume `json:"rootVolume"`
}

type EC2RootVolume struct {
	IOPS int    `json:"iops"`
	Size int    `json:"size"`
	Type string `json:"type"`
}

type LibvirtMachinePoolPlatformConfig struct {
	// QCOWImagePath
	QCOWImagePath string `json:"qcowImagePath"`
}
