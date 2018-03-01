package gcp

// Config defines the GCP configuraiton for a cluster.
type Config struct {
	ConfigVersion            string    `yaml:"configVersion,omitempty"`
	Etcd                     component `yaml:"etcd,omitempty"`
	ExtGoogleManagedZoneName string    `yaml:"extGoogleManagedZoneName,omitempty"`
	Master                   component `yaml:"master,omitempty"`
	Region                   string    `yaml:"region,omitempty"`
	SSHKey                   string    `yaml:"sshKey,omitempty"`
	Worker                   component `yaml:"worker,omitempty"`
}
