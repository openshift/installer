package metal

// Config defines the Metal configuraiton for a cluster.
type Config struct {
	CalicoMTU     string    `yaml:"CalicoMTU,omitempty"`
	Controller    component `yaml:"Controller,omitempty"`
	IngressDomain string    `yaml:"IngressDomain,omitempty"`
	Matchbox      matchbox  `yaml:"Matchbox,omitempty"`
	Worker        component `yaml:"Worker,omitempty"`
}
