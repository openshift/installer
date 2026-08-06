package gcp

// Metadata contains GCP metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	// NetworkProjectID is used for shared VPC setups
	// Deprecated. Use the Secret referenced by ClusterMetadata.MetadataJSONSecretRef instead. We
	// may stop populating this section in the future.
	// +optional
	NetworkProjectID *string `json:"networkProjectID,omitempty"`
}
