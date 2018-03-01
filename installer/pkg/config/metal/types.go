package metal

type client struct {
	Cert string `yaml:"cert,omitempty"`
	Key  string `yaml:"key,omitempty"`
}

type component struct {
	Domain  string `yaml:"domain,omitempty"`
	Domains string `yaml:"domains,omitempty"`
	MACs    string `yaml:"macs,omitempty"`
	Names   string `yaml:"names,omitempty"`
}

type matchbox struct {
	CA          string `yaml:"ca,omitempty"`
	Client      client `yaml:"client,omitempty"`
	HTTPURL     string `yaml:"httpURL,omitempty"`
	RPCEndpoint string `yaml:"rpcEndpoint,omitempty"`
}
