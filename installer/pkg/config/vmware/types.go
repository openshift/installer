package vmware

type component struct {
	Clusters     string `yaml:"clusters,omitempty"`
	Datacenters  string `yaml:"datacenters,omitempty"`
	Datastores   string `yaml:"datastores,omitempty"`
	Gateways     string `yaml:"gateways,omitempty"`
	Hostnames    string `yaml:"hostnames,omitempty"`
	IP           string `yaml:"ip,omitempty"`
	Memory       string `yaml:"memory,omitempty"`
	Networks     string `yaml:"networks,omitempty"`
	ResourcePool string `yaml:"resourcePool,omitempty"`
	VCPU         string `yaml:"vCPU,omitempty"`
}

type ssh struct {
	AuthorizedKey  string `yaml:"authorizedKey,omitempty"`
	PrivateKeyPath string `yaml:"privateKeyPath,omitempty"`
}

type vm struct {
	Template       string `yaml:"template,omitempty"`
	TemplateFolder string `yaml:"templateFolder,omitempty"`
}
