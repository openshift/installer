package types

// ClusterMetadata contains information
// regarding the cluster that was created by installer.
type ClusterMetadata struct {
	ClusterName             string `json:"clusterName"`
	ClusterPlatformMetadata `json:",inline"`
}

// ClusterPlatformMetadata contains metadata for platfrom.
type ClusterPlatformMetadata struct {
	AWS     *ClusterAWSPlatformMetadata     `json:"aws,omitempty"`
	Libvirt *ClusterLibvirtPlatformMetadata `json:"libvirt,omitempty"`
}

// ClusterAWSPlatformMetadata contains AWS metadata.
type ClusterAWSPlatformMetadata struct {
	Region string `json:"region"`
	// Most AWS resources are tagged with these tags as identifier.
	Identifier map[string]string `json:"identifier"`
}

// ClusterLibvirtPlatformMetadata contains libvirt metadata.
type ClusterLibvirtPlatformMetadata struct {
	URI string `json:"uri"`
}
