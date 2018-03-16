package openstack

// EtcdFlavor converts etcd flavor related config.
type EtcdFlavor struct {
	ID   string `json:"tectonic_openstack_etcd_flavor_id,omitempty" yaml:"id,omitempty"`
	Name string `json:"tectonic_openstack_etcd_flavor_name,omitempty" yaml:"name,omitempty"`
}

// Image converts image related config.
type Image struct {
	ID   string `json:"tectonic_openstack_image_id,omitempty" yaml:"id,omitempty"`
	Name string `json:"tectonic_openstack_image_name,omitempty" yaml:"name,omitempty"`
}

// MasterFlavor converts master flavor related config.
type MasterFlavor struct {
	ID   string `json:"tectonic_openstack_master_flavor_id,omitempty" yaml:"id,omitempty"`
	Name string `json:"tectonic_openstack_master_flavor_name,omitempty" yaml:"name,omitempty"`
}

// OpenStack converts OpenStack related config.
type OpenStack struct {
	DisableFloatingIP string `json:"tectonic_openstack_disable_floatingip,omitempty" yaml:"disableFloatingIP,omitempty"`
	DNSNameservers    string `json:"tectonic_openstack_dns_nameservers,omitempty" yaml:"dnsNameservers,omitempty"`
	EtcdFlavor        `json:",inline" yaml:"etcdFlavor,omitempty"`
	ExternalGatewayID string `json:"tectonic_openstack_external_gateway_id,omitempty" yaml:"externalGatewayID,omitempty"`
	FloatingIPPool    string `json:"tectonic_openstack_floatingip_pool,omitempty" yaml:"floatingIPPool,omitempty"`
	Image             `json:",inline" yaml:"image,omitempty"`
	LBProvider        string `json:"tectonic_openstack_lb_provider,omitempty" yaml:"lbProvider,omitempty"`
	MasterFlavor      `json:",inline" yaml:"masterFlavor,omitempty"`
	SubnetCIDR        string `json:"tectonic_openstack_subnet_cidr,omitempty" yaml:"subnetCIDR,omitempty"`
	WorkerFlavor      `json:",inline" yaml:"workerFlavor,omitempty"`
}

// WorkerFlavor converts worker flavor related config.
type WorkerFlavor struct {
	ID   string `json:"tectonic_openstack_worker_flavor_id,omitempty" yaml:"id,omitempty"`
	Name string `json:"tectonic_openstack_worker_flavor_name,omitempty" yaml:"name,omitempty"`
}
