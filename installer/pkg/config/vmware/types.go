package vmware

type component struct {
	Clusters     string `yaml:"Clusters,omitempty"`
	Datacenters  string `yaml:"Datacenters,omitempty"`
	Datastores   string `yaml:"Datastores,omitempty"`
	Gateways     string `yaml:"Gateways,omitempty"`
	Hostnames    string `yaml:"Hostnames,omitempty"`
	IP           string `yaml:"IP,omitempty"`
	Memory       string `yaml:"Memory,omitempty"`
	Networks     string `yaml:"Networks,omitempty"`
	ResourcePool string `yaml:"ResourcePool,omitempty"`
	VCPU         string `yaml:"VCPU,omitempty"`
}

type ssh struct {
	AuthorizedKey  string `yaml:"AuthorizedKey,omitempty"`
	PrivateKeyPath string `yaml:"PrivateKeyPath,omitempty"`
}

type vm struct {
	Template       string `yaml:"Template,omitempty"`
	TemplateFolder string `yaml:"TemplateFolder,omitempty"`
}
