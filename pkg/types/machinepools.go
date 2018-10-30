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

// Name returns a string representation of the platform (e.g. "aws" if
// AWS is non-nil).  It returns an empty string if no platform is
// configured.
func (p *MachinePoolPlatform) Name() string {
	if p == nil {
		return ""
	}
	if p.AWS != nil {
		return PlatformNameAWS
	}
	if p.Libvirt != nil {
		return PlatformNameLibvirt
	}
	if p.OpenStack != nil {
		return PlatformNameOpenstack
	}
	return ""
}

// AWSMachinePoolPlatform stores the configuration for a machine pool
// installed on AWS.
type AWSMachinePoolPlatform struct {
	// Zones is list of availability zones that can be used.
	Zones []string `json:"zones,omitempty"`

	// AMIID defines the AMI that should be used.
	AMIID string `json:"amiID,omitempty"`

	// InstanceType defines the ec2 instance type.
	// eg. m4-large
	InstanceType string `json:"type"`

	// IAMRoleName defines the IAM role associated
	// with the ec2 instance.
	IAMRoleName string `json:"iamRoleName"`

	// EC2RootVolume defines the storage for ec2 instance.
	EC2RootVolume `json:"rootVolume"`
}

// Set sets the values from `required` to `a`.
func (a *AWSMachinePoolPlatform) Set(required *AWSMachinePoolPlatform) {
	if required == nil || a == nil {
		return
	}

	if len(required.Zones) > 0 {
		a.Zones = required.Zones
	}
	if required.AMIID != "" {
		a.AMIID = required.AMIID
	}
	if required.InstanceType != "" {
		a.InstanceType = required.InstanceType
	}
	if required.IAMRoleName != "" {
		a.IAMRoleName = required.IAMRoleName
	}

	if required.EC2RootVolume.IOPS != 0 {
		a.EC2RootVolume.IOPS = required.EC2RootVolume.IOPS
	}
	if required.EC2RootVolume.Size != 0 {
		a.EC2RootVolume.Size = required.EC2RootVolume.Size
	}
	if required.EC2RootVolume.Type != "" {
		a.EC2RootVolume.Type = required.EC2RootVolume.Type
	}
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

// Set sets the values from `required` to `a`.
func (o *OpenStackMachinePoolPlatform) Set(required *OpenStackMachinePoolPlatform) {
	if required == nil || o == nil {
		return
	}

	if required.FlavorName != "" {
		o.FlavorName = required.FlavorName
	}

	if required.OpenStackRootVolume.IOPS != 0 {
		o.OpenStackRootVolume.IOPS = required.OpenStackRootVolume.IOPS
	}
	if required.OpenStackRootVolume.Size != 0 {
		o.OpenStackRootVolume.Size = required.OpenStackRootVolume.Size
	}
	if required.OpenStackRootVolume.Type != "" {
		o.OpenStackRootVolume.Type = required.OpenStackRootVolume.Type
	}
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

// Set sets the values from `required` to `a`.
func (l *LibvirtMachinePoolPlatform) Set(required *LibvirtMachinePoolPlatform) {
	if required == nil || l == nil {
		return
	}

	if required.ImagePool != "" {
		l.ImagePool = required.ImagePool
	}
	if required.ImageVolume != "" {
		l.ImageVolume = required.ImageVolume
	}
	if required.Image != "" {
		l.Image = required.Image
	}
}
