package openstack

// Config defines the OpenStack configuraiton for a cluster.
type Config struct {
	DisableFloatingIP string `yaml:"DisableFloatingIP,omitempty"`
	DNSNameservers    string `yaml:"DNSNameservers,omitempty"`
	EtcdFlavor        flavor `yaml:"EtcdFlavor,omitempty"`
	ExternalGatewayID string `yaml:"ExternalGatewayID,omitempty"`
	FloatingIPPool    string `yaml:"FloatingIPPool,omitempty"`
	Image             flavor `yaml:"Image,omitempty"`
	LBProvider        string `yaml:"LBProvider,omitempty"`
	MasterFlavor      flavor `yaml:"MasterFlavor,omitempty"`
	SubnetCIDR        string `yaml:"SubnetCIDR,omitempty"`
	WorkerFlavor      flavor `yaml:"WorkerFlavor,omitempty"`
}
