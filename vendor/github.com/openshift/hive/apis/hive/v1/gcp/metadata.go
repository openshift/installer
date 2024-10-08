package gcp

// Metadata contains GCP metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	// NetworkProjectID is used for shared VPC setups
	// +optional
	NetworkProjectID *string `json:"networkProjectID,omitempty"`
}
