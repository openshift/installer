package terraformgenerator

import (
	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

// GCP defines all variables for this platform.
type GCP struct {
	ConfigVersion            string `json:"tectonic_gcp_config_version,omitempty"`
	EtcdDiskSize             string `json:"tectonic_gcp_etcd_disk_size,omitempty"`
	EtcdDiskType             string `json:"tectonic_gcp_etcd_disktype,omitempty"`
	EtcdGCEType              string `json:"tectonic_gcp_etcd_gce_type,omitempty"`
	ExtGoogleManagedZoneName string `json:"tectonic_gcp_ext_google_managedzone_name,omitempty"`
	MasterDiskSize           string `json:"tectonic_gcp_master_disk_size,omitempty"`
	MasterDiskType           string `json:"tectonic_gcp_master_disktype,omitempty"`
	MasterGCEType            string `json:"tectonic_gcp_master_gce_type,omitempty"`
	Region                   string `json:"tectonic_gcp_region,omitempty"`
	SSHKey                   string `json:"tectonic_gcp_ssh_key,omitempty"`
	WorkerDiskSize           string `json:"tectonic_gcp_worker_disk_size,omitempty"`
	WorkerDiskType           string `json:"tectonic_gcp_worker_disktype,omitempty"`
	WorkerGCEType            string `json:"tectonic_gcp_worker_gce_type,omitempty"`
}

// NewGCP returns the config for GCP.
func NewGCP(cluster config.Cluster) GCP {
	return GCP{
		ConfigVersion:            cluster.GCP.ConfigVersion,
		EtcdDiskSize:             cluster.GCP.Etcd.DiskSize,
		EtcdDiskType:             cluster.GCP.Etcd.DiskType,
		EtcdGCEType:              cluster.GCP.Etcd.GCEType,
		ExtGoogleManagedZoneName: cluster.GCP.ExtGoogleManagedZoneName,
		MasterDiskSize:           cluster.GCP.Master.DiskSize,
		MasterDiskType:           cluster.GCP.Master.DiskType,
		MasterGCEType:            cluster.GCP.Master.GCEType,
		Region:                   cluster.GCP.Region,
		SSHKey:                   cluster.GCP.SSHKey,
		WorkerDiskSize:           cluster.GCP.Worker.DiskSize,
		WorkerDiskType:           cluster.GCP.Worker.DiskType,
		WorkerGCEType:            cluster.GCP.Worker.GCEType,
	}
}
