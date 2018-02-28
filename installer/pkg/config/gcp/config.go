package gcp

// Config defines the GCP configuraiton for a cluster.
type Config struct {
	ConfigVersion            string    `yaml:"ConfigVersion,omitempty"`
	Etcd                     component `yaml:"Etcd,omitempty"`
	ExtGoogleManagedZoneName string    `yaml:"ExtGoogleManagedZoneName,omitempty"`
	Master                   component `yaml:"Master,omitempty"`
	Region                   string    `yaml:"Region,omitempty"`
	SSHKey                   string    `yaml:"SSHKey,omitempty"`
	Worker                   component `yaml:"Worker,omitempty"`
}
