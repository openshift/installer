package alibabacloud

// Metadata contains Alibaba Cloud metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	Region          string `json:"region"`
	ResourceGroupID string `json:"resourceGroupID"`
	ClusterDomain   string `json:"clusterDomain"`
}
