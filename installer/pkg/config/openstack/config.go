package openstack

// Config defines the OpenStack configuraiton for a cluster.
type Config struct {
	DisableFloatingIP string `yaml:"disableFloatingIP,omitempty"`
	DNSNameservers    string `yaml:"dnsNameservers,omitempty"`
	EtcdFlavor        flavor `yaml:"etcdFlavor,omitempty"`
	ExternalGatewayID string `yaml:"externalGatewayID,omitempty"`
	FloatingIPPool    string `yaml:"floatingIPPool,omitempty"`
	Image             flavor `yaml:"image,omitempty"`
	LBProvider        string `yaml:"lbProvider,omitempty"`
	MasterFlavor      flavor `yaml:"masterFlavor,omitempty"`
	SubnetCIDR        string `yaml:"subnetCIDR,omitempty"`
	WorkerFlavor      flavor `yaml:"workerFlavor,omitempty"`
}
