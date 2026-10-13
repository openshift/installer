package baremetal

// Metadata contains baremetal metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	LibvirtURI              string `json:"libvirtURI"`
	BootstrapProvisioningIP string `json:"bootstrapProvisioningIP"`
	ClusterProvisioningIP   string `json:"provisioningHostIP"`
}
