package aws

type component struct {
	CustomSubnets string     `yaml:"CustomSubnets,omitempty"`
	EC2Type       string     `yaml:"EC2Type,omitempty"`
	ExtraSGIDs    string     `yaml:"ExtraSGIDs,omitempty"`
	IAMRoleName   string     `yaml:"IAMRoleName,omitempty"`
	LoadBalancers string     `yaml:"LoadBalancers,omitempty"`
	RootVolume    rootVolume `yaml:"RootVolume,omitempty"`
}

type external struct {
	MasterSubnetIDs string `yaml:"MasterSubnetIDs,omitempty"`
	PrivateZone     string `yaml:"PrivateZone,omitempty"`
	VPCID           string `yaml:"VPCIC,omitempty"`
	WorkerSubnetIDs string `yaml:"WorkerSubnetIDs,omitempty"`
}

type rootVolume struct {
	IOPS int    `yaml:"RootVolumeIOPS,omitempty"`
	Size int    `yaml:"RootVolumeSize,omitempty"`
	Type string `yaml:"RootVolumeType,omitempty"`
}
