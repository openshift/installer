package vmware

// Config defines the VMware configuraiton for a cluster.
type Config struct {
	ControllerDomain string    `yaml:"controllerDomain,omitempty"`
	Etcd             component `yaml:"etcd,omitempty"`
	Folder           string    `yaml:"folder,omitempty"`
	IngressDomain    string    `yaml:"ingressDomain,omitempty"`
	Master           component `yaml:"master,omitempty"`
	NodeDNS          string    `yaml:"nodeDNS,omitempty"`
	Server           string    `yaml:"server,omitempty"`
	SSH              ssh       `yaml:"ssh,omitempty"`
	SSLSelfSigned    string    `yaml:"sslSelfsigned,omitempty"`
	Type             string    `yaml:"type,omitempty"`
	VM               vm        `yaml:"vm,omitempty"`
	Worker           component `yaml:"worker,omitempty"`
}
