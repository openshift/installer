package vmware

// Config defines the VMware configuraiton for a cluster.
type Config struct {
	ControllerDomain string    `yaml:"ControllerDomain,omitempty"`
	Etcd             component `yaml:"Etcd,omitempty"`
	Folder           string    `yaml:"Folder,omitempty"`
	IngressDomain    string    `yaml:"IngressDomain,omitempty"`
	Master           component `yaml:"Master,omitempty"`
	NodeDNS          string    `yaml:"NodeDNS,omitempty"`
	Server           string    `yaml:"Server,omitempty"`
	SSH              ssh       `yaml:"SSH,omitempty"`
	SSLSelfSigned    string    `yaml:"SSLSelfsigned,omitempty"`
	Type             string    `yaml:"Type,omitempty"`
	VM               vm        `yaml:"VM,omitempty"`
	Worker           component `yaml:"Worker,omitempty"`
}
