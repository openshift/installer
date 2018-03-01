package metal

// Config defines the Metal configuraiton for a cluster.
type Config struct {
	CalicoMTU        string    `yaml:"calicoMTU,omitempty"`
	Controller       component `yaml:"controller,omitempty"`
	IngressDomain    string    `yaml:"ingressDomain,omitempty"`
	Matchbox         matchbox  `yaml:"matchbox,omitempty"`
	SSHAuthorizedKey string    `yaml:"sshAuthorizedKey,omitempty"`
	Worker           component `yaml:"worker,omitempty"`
}
