package azure

// Metadata contains Azure metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	// ResourceGroupName is the name of the resource group in which the cluster resources were created.
	// Deprecated. Use the Secret referenced by ClusterMetadata.MetadataJSONSecretRef instead. We
	// may stop populating this section in the future.
	ResourceGroupName *string `json:"resourceGroupName"`
}
