package azure

// Metadata contains Azure metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	CloudName CloudEnvironment `json:"cloudName"`
	Region    string           `json:"region"`
}
