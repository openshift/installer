package aws

// Config defines the AWS configuraiton for a cluster.
type Config struct {
	AssetsS3BucketName        string    `yaml:"assetsS3BucketName,omitempty"`
	AutoScalingGroupExtraTags string    `yaml:"autoScalingGroupExtraTags,omitempty"`
	EC2AMIOverride            string    `yaml:"ec2AMIOverride,omitempty"`
	Etcd                      component `yaml:"etcd,omitempty"`
	External                  external  `yaml:"external,omitempty"`
	ExtraTags                 string    `yaml:"extraTags,omitempty"`
	InstallerRole             string    `yaml:"installerRole,omitempty"`
	Master                    component `yaml:"master,omitempty"`
	PrivateEndpoints          bool      `yaml:"privateEndpoints,omitempty"`
	Profile                   string    `yaml:"profile,omitempty"`
	PublicEndpoints           bool      `yaml:"publicEndpoints,omitempty"`
	Region                    string    `yaml:"region,omitempty"`
	SSHKey                    string    `yaml:"sshKey,omitempty"`
	VPCCIDRBlock              string    `yaml:"vpcCIDRBlock,omitempty"`
	Worker                    component `yaml:"worker,omitempty"`
}
