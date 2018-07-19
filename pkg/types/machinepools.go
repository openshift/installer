package types

// MachinePool is a pool of machines to be installed.
type MachinePool struct {
	// Name is the name of the machine pool.
	Name string `json:"name"`

	// Replicas is the count of machines for this machine pool.
	// Default is 1.
	Replicas *int64 `json:"replicas"`

	// Platform is configuration for machine pool specific to the platfrom.
	Platform MachinePoolPlatform `json:"platform"`
}

// MachinePoolPlatform is the platform-specific configuration for a machine
// pool. Only one of the platforms should be set.
type MachinePoolPlatform struct {
	// AWS is the configuration used when installing on AWS.
	AWS *AWSMachinePoolPlatform `json:"aws,omitempty"`

	// Libvirt is the configuration used when installing on libvirt.
	Libvirt *LibvirtMachinePoolPlatform `json:"libvirt,omitempty"`

	// OpenStack is the configuration used when installing on OpenStack.
	OpenStack *OpenStackMachinePoolPlatform `json:"openstack,omitempty"`
}

// AWSMachinePoolPlatform stores the configuration for a machine pool
// installed on AWS.
type AWSMachinePoolPlatform struct {
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

// OpenStackMachinePoolPlatform stores the configuration for a machine pool
// installed on OpenStack.
type OpenStackMachinePoolPlatform struct {
	// FlavorName defines the OpenStack Nova flavor.
	// eg. m1.large
	FlavorName string `json:"type"`

	// OpenStackRootVolume defines the storage for Nova instance.
	OpenStackRootVolume `json:"rootVolume"`
}

// OpenStackRootVolume defines the storage for a Nova instance.
type OpenStackRootVolume struct {
	// IOPS defines the iops for the instance.
	IOPS int `json:"iops"`
	// Size defines the size of the instance.
	Size int `json:"size"`
	// Type defines the type of the instance.
	Type string `json:"type"`
}

// LibvirtMachinePoolPlatform stores the configuration for a machine pool
// installed on libvirt.
type LibvirtMachinePoolPlatform struct {
	// ImagePool is the name of the libvirt storage pool to which the storage
	// volume containing the OS image belongs.
	ImagePool string `json:"imagePool,omitempty"`
	// ImageVolume is the name of the libvirt storage volume containing the OS
	// image.
	ImageVolume string `json:"imageVolume,omitempty"`

	// Image is the URL to the OS image.
	// E.g. "http://aos-ostree.rhev-ci-vms.eng.rdu2.redhat.com/rhcos/images/cloud/latest/rhcos-qemu.qcow2.gz"
	Image string `json:"image"`
}
