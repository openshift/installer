package aws

// Config defines the AWS configuraiton for a cluster.
type Config struct {
	AssetsS3BucketName        string    `yaml:"AssetsS3BucketName,omitempty"`
	AutoScalingGroupExtraTags string    `yaml:"AutoScalingGroupExtraTags,omitempty"`
	EC2AMIOverride            string    `yaml:"EC2AMIOverride,omitempty"`
	Etcd                      component `yaml:"Etcd,omitempty"`
	External                  external  `yaml:"External,omitempty"`
	ExtraTags                 string    `yaml:"ExtraTags,omitempty"`
	InstallerRole             string    `yaml:"InstallerRole,omitempty"`
	Master                    component `yaml:"Master,omitempty"`
	PrivateEndpoints          bool      `yaml:"PrivateEndpoints,omitempty"`
	Profile                   string    `yaml:"Profile,omitempty"`
	PublicEndpoints           bool      `yaml:"PublicEndpoints,omitempty"`
	Region                    string    `yaml:"Region,omitempty"`
	SSHKey                    string    `yaml:"SSHKey,omitempty"`
	VPCCIDRBlock              string    `yaml:"VPCCIDRBlock,omitempty"`
	Worker                    component `yaml:"Worker,omitempty"`
}
