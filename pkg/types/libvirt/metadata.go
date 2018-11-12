package libvirt

// Metadata contains libvirt metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	URI string `json:"uri"`
}
