package azure

// Config defines the Azure configuraiton for a cluster.
type Config struct {
	CloudEnvironment string    `yaml:"cloudEnvironment,omitempty"`
	Etcd             component `yaml:"etcd,omitempty"`
	External         external  `yaml:"external,omitempty"`
	ExtraTags        string    `yaml:"extraTags,omitempty"`
	Master           component `yaml:"master,omitempty"`
	PrivateCluster   string    `yaml:"privateCluster,omitempty"`
	SSH              ssh       `yaml:"ssh,omitempty"`
	VNetCIDRBlock    string    `yaml:"vNetCIDRBlock,omitempty"`
	Worker           component `yaml:"worker,omitempty"`
}
