package ovirt

//config is the oVirt's cloud provider config
type config struct {
	// StorageDomainID the id of the storage domain for the OS disks of VMs.
	StorageDomainID string `json:"storageDomainId"`
	// ClusterID the id of the cluster of the VMs.
	ClusterID string `json:"clusterId"`
	// NetworkName the name of the network used for the VMs network interfaces.
	NetworkName string `json:"networkName"`
}
