package azure

// Metadata contains Azure metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	ARMEndpoint                 string           `json:"armEndpoint"`
	CloudName                   CloudEnvironment `json:"cloudName"`
	Region                      string           `json:"region"`
	ResourceGroupName           string           `json:"resourceGroupName"`
	BaseDomainResourceGroupName string           `json:"baseDomainResourceGroupName"`
}

// Keys used to save Metadata information as tags.
const (
	TagMetadataRegion       = "openshift_region"
	TagMetadataBaseDomainRG = "openshift_basedomainRG"
	TagMetadataNetworkRG    = "openshift_networkRG"
)
