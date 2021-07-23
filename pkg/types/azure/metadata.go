package azure

// Metadata contains Azure metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	ARMEndpoint                 string           `json:"armEndpoint"`
	CloudName                   CloudEnvironment `json:"cloudName"`
	Region                      string           `json:"region"`
	ResourceGroupName           string           `json:"resourceGroupName"`
	ClusterName                 string           `json:"clusterName"`
	BaseDomainResourceGroupName string           `json:"baseDomainResourceGroupName"`
}
