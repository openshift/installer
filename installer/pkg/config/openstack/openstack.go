package openstack

const (
	// DefaultCloud is the default OpenStack cloud to use.
	DefaultCloud = "default"
)

// Openstack converts OpenStack related config.
type Openstack struct {
	Cloud   string `json:"tectonic_openstack_cloud,omitempty" yaml:"cloud,omitempty"`
	Etcd    `json:",inline" yaml:"etcd,omitempty"`
	KeyPair string `json:"tectonic_openstack_key_pair,omitempty" yaml:"keyPair,omitempty"`
	Master  `json:",inline" yaml:"master,omitempty"`
	SSHKey  string `json:"tectonic_openstack_ssh_key,omitempty" yaml:"sshKey,omitempty"`
	Worker  `json:",inline" yaml:"worker,omitempty"`
}

// Etcd converts etcd related config.
type Etcd struct {
	ImageName  string   `json:"tectonic_openstack_etcd_image_name,omitempty" yaml:"imageName,omitempty"`
	FlavorName string   `json:"tectonic_openstack_etcd_flavor_name,omitempty" yaml:"flavorName,omitempty"`
	ExtraSGIDs []string `json:"tectonic_openstack_etcd_extra_sg_ids,omitempty" yaml:"extraSGIDs,omitempty"`
}

// Master converts master related config.
type Master struct {
	ImageName  string   `json:"tectonic_openstack_master_image_name,omitempty" yaml:"imageName,omitempty"`
	FlavorName string   `json:"tectonic_openstack_master_flavor_name,omitempty" yaml:"flavorName,omitempty"`
	ExtraSGIDs []string `json:"tectonic_openstack_master,omitempty" yaml:"extraSGIDs,omitempty"`
}

// Worker converts worker related config.
type Worker struct {
	ImageName  string   `json:"tectonic_openstack_worker_image_name,omitempty" yaml:"imageName,omitempty"`
	FlavorName string   `json:"tectonic_openstack_worker_flavor_name,omitempty" yaml:"flavorName,omitempty"`
	ExtraSGIDs []string `json:"tectonic_openstack_worker,omitempty" yaml:"extraSGIDs,omitempty"`
}
