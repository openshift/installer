package aws

// Metadata contains AWS metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	// HostedZoneRole is the role to assume when performing operations
	// on a hosted zone owned by another account.
	// Deprecated. Use the Secret referenced by ClusterMetadata.MetadataJSONSecretRef instead. We
	// may stop populating this section in the future.
	HostedZoneRole *string `json:"hostedZoneRole,omitempty"`
}
