package types

// ClusterMetadata contains information
// regarding the cluster that was created by installer.
type ClusterMetadata struct {
	ClusterName             string `json:"clusterName"`
	ClusterPlatformMetadata `json:",inline"`
}

// ClusterPlatformMetadata contains metadata for platfrom.
type ClusterPlatformMetadata struct {
	AWS       *ClusterAWSPlatformMetadata       `json:"aws,omitempty"`
	OpenStack *ClusterOpenStackPlatformMetadata `json:"openstack,omitempty"`
	Libvirt   *ClusterLibvirtPlatformMetadata   `json:"libvirt,omitempty"`
}

// Platform returns a string representation of the platform
// (e.g. "aws" if AWS is non-nil).  It returns an empty string if no
// platform is configured.
func (cpm *ClusterPlatformMetadata) Platform() string {
	if cpm == nil {
		return ""
	}
	if cpm.AWS != nil {
		return "aws"
	}
	if cpm.Libvirt != nil {
		return "libvirt"
	}
	if cpm.OpenStack != nil {
		return "openstack"
	}
	return ""
}

// ClusterAWSPlatformMetadata contains AWS metadata.
type ClusterAWSPlatformMetadata struct {
	Region string `json:"region"`

	// Identifier holds a slice of filter maps.  The maps hold the
	// key/value pairs for the tags we will be matching against.  A
	// resource matches the map if all of the key/value pairs are in its
	// tags.  A resource matches Identifier if it matches any of the maps.
	Identifier []map[string]string `json:"identifier"`
}

// ClusterOpenStackPlatformMetadata contains OpenStack metadata.
type ClusterOpenStackPlatformMetadata struct {
	Region string `json:"region"`
	// Most OpenStack resources are tagged with these tags as identifier.
	Identifier map[string]string `json:"identifier"`
}

// ClusterLibvirtPlatformMetadata contains libvirt metadata.
type ClusterLibvirtPlatformMetadata struct {
	URI string `json:"uri"`
}
