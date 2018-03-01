package aws

type component struct {
	CustomSubnets string     `yaml:"customSubnets,omitempty"`
	EC2Type       string     `yaml:"ec2Type,omitempty"`
	ExtraSGIDs    string     `yaml:"extraSGIDs,omitempty"`
	IAMRoleName   string     `yaml:"iamRoleName,omitempty"`
	LoadBalancers string     `yaml:"loadBalancers,omitempty"`
	RootVolume    rootVolume `yaml:"rootVolume,omitempty"`
}

type external struct {
	MasterSubnetIDs string `yaml:"masterSubnetIDs,omitempty"`
	PrivateZone     string `yaml:"privateZone,omitempty"`
	VPCID           string `yaml:"vpcID,omitempty"`
	WorkerSubnetIDs string `yaml:"workerSubnetIDs,omitempty"`
}

type rootVolume struct {
	IOPS int    `yaml:"iops,omitempty"`
	Size int    `yaml:"size,omitempty"`
	Type string `yaml:"type,omitempty"`
}
