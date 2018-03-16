package gcp

// Etcd converts etcd related config.
type Etcd struct {
	DiskSize string `json:"tectonic_gcp_etcd_disk_size,omitempty" yaml:"diskSize,omitempty"`
	DiskType string `json:"tectonic_gcp_etcd_disktype,omitempty" yaml:"diskType,omitempty"`
	GCEType  string `json:"tectonic_gcp_etcd_gce_type,omitempty" yaml:"gceType,omitempty"`
}

// GCP converts GCP related config.
type GCP struct {
	ConfigVersion            string `json:"tectonic_gcp_config_version,omitempty" yaml:"configVersion,omitempty"`
	Etcd                     `json:",inline" yaml:"etcd,omitempty"`
	ExtGoogleManagedZoneName string `json:"tectonic_gcp_ext_google_managedzone_name,omitempty" yaml:"extGoogleManagedZoneName,omitempty"`
	Master                   `json:",inline" yaml:"master,omitempty"`
	Region                   string `json:"tectonic_gcp_region,omitempty" yaml:"region,omitempty"`
	SSHKey                   string `json:"tectonic_gcp_ssh_key,omitempty" yaml:"sshKey,omitempty"`
	Worker                   `json:",inline" yaml:"worker,omitempty"`
}

// Master converts master related config.
type Master struct {
	DiskSize string `json:"tectonic_gcp_master_disk_size,omitempty" yaml:"diskSize,omitempty"`
	DiskType string `json:"tectonic_gcp_master_disktype,omitempty" yaml:"diskType,omitempty"`
	GCEType  string `json:"tectonic_gcp_master_gce_type,omitempty" yaml:"gceType,omitempty"`
}

// Worker converts worker related config.
type Worker struct {
	DiskSize string `json:"tectonic_gcp_worker_disk_size,omitempty" yaml:"diskSize,omitempty"`
	DiskType string `json:"tectonic_gcp_worker_disktype,omitempty" yaml:"diskType,omitempty"`
	GCEType  string `json:"tectonic_gcp_worker_gce_type,omitempty" yaml:"gceType,omitempty"`
}
