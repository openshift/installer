package ovirt

// Metadata contains ovirt metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	ClusterID      string `json:"cluster_id"`
	RemoveTemplate bool   `json:"remove_template"`
}
