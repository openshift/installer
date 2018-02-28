package azure

// Config defines the Azure configuraiton for a cluster.
type Config struct {
	CloudEnvironment string    `yaml:"CloudEnvironment,omitempty"`
	Etcd             component `yaml:"Etcd,omitempty"`
	External         external  `yaml:"External,omitempty"`
	ExtraTags        string    `yaml:"ExtraTags,omitempty"`
	Master           component `yaml:"Master,omitempty"`
	PrivateCluster   string    `yaml:"PrivateCluster,omitempty"`
	SSH              ssh       `yaml:"SSH,omitempty"`
	VNetCIDRBlock    string    `yaml:"VNetCIDRBlock,omitempty"`
	Worker           component `yaml:"Worker,omitempty"`
}
