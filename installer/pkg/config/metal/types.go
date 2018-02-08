package metal

type client struct {
	Cert string `yaml:"Cert,omitempty"`
	Key  string `yaml:"Key,omitempty"`
}

type component struct {
	Domain  string `yaml:"Domain,omitempty"`
	Domains string `yaml:"Domains,omitempty"`
	MACs    string `yaml:"MACs,omitempty"`
	Names   string `yaml:"Names,omitempty"`
}

type matchbox struct {
	CA          string `yaml:"CA,omitempty"`
	Client      client `yaml:"Client,omitempty"`
	HTTPURL     string `yaml:"HTTPURL,omitempty"`
	RPCEndpoint string `yaml:"RPCEndpoint,omitempty"`
}
