package powervs

// Metadata contains Power VS metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	Region string `json:"region"`
	Zone   string `json:"zone"`
}
