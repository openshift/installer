package metal

// Client converts client related config.
type Client struct {
	Cert string `json:"tectonic_metal_matchbox_client_cert,omitempty" yaml:"cert,omitempty"`
	Key  string `json:"tectonic_metal_matchbox_client_key,omitempty" yaml:"key,omitempty"`
}

// Controller converts controller related config.
type Controller struct {
	Domain  string `json:"tectonic_metal_controller_domain,omitempty" yaml:"domain,omitempty"`
	Domains string `json:"tectonic_metal_controller_domains,omitempty" yaml:"domains,omitempty"`
	MACs    string `json:"tectonic_metal_controller_macs,omitempty" yaml:"macs,omitempty"`
	Names   string `json:"tectonic_metal_controller_names,omitempty" yaml:"names,omitempty"`
}

// Matchbox converts matchbox related config.
type Matchbox struct {
	CA          string `json:"tectonic_metal_matchbox_ca,omitempty" yaml:"ca,omitempty"`
	Client      `json:",inline" yaml:"client,omitempty"`
	HTTPURL     string `json:"tectonic_metal_matchbox_http_url,omitempty" yaml:"httpURL,omitempty"`
	RPCEndpoint string `json:"tectonic_metal_matchbox_rpc_endpoint,omitempty" yaml:"rpcEndpoint,omitempty"`
}

// Metal converts metal related config.
type Metal struct {
	CalicoMTU        string `json:"tectonic_metal_calico_mtu,omitempty" yaml:"calicoMTU,omitempty"`
	Controller       `json:",inline" yaml:"controller,omitempty"`
	IngressDomain    string `json:"tectonic_metal_ingress_domain,omitempty" yaml:"ingressDomain,omitempty"`
	Matchbox         `json:",inline" yaml:"matchbox,omitempty"`
	SSHAuthorizedKey string `json:"tectonic_ssh_authorized_key,omitempty" yaml:"sshAuthorizedKey,omitempty"` // TODO:(spangenberg): Prefix with metal.
	Worker           `json:",inline" yaml:"worker,omitempty"`
}

// Worker converts worker related config.
type Worker struct {
	Domains string `json:"tectonic_metal_worker_domains,omitempty" yaml:"domains,omitempty"`
	MACs    string `json:"tectonic_metal_worker_macs,omitempty" yaml:"macs,omitempty"`
	Names   string `json:"tectonic_metal_worker_names,omitempty" yaml:"names,omitempty"`
}
