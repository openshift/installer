package alibabacloud

// Metadata contains Alibaba Cloud metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	Region string `json:"region"`
	// The system checks whether the resource group contains resources after deleted,
	// This process takes 7 days. The resource group remains in the 'Deleting' state
	// and couldn't create resource group with the same name during this period.
	// Before deploying the cluster, the user must manually create a resource group.
	// The parameter ResourceGroupID is required.
	ResourceGroupID string `json:"resourceGroupID"`
	ClusterDomain   string `json:"clusterDomain"`
}
