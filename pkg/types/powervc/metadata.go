package powervc

// Metadata contains PowerVC and OpenStack metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	Cloud string `json:"cloud"`
	// Most OpenStack resources are tagged with these tags as identifier.
	Identifier map[string]string `json:"identifier"`
}
