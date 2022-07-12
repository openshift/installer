package openstack

// FailureDomain holds the information for a failure domain
type FailureDomain struct {
	// Name defines the name of the OpenStackPlatformFailureDomainSpec
	Name string `json:"name"`

	// ComputeZone is the compute zone on which the nodes belonging to the
	// failure domain must be provisioned.
	// If not specified, the nodes are provisioned in the OpenStack Nova default availabity zone.
	// +optional
	ComputeZone string `json:"computeZone,omitempty"`

	// StorageZone is the storage zone from where volumes should be provisioned
	// for the nodes belonging to the failure domain.
	// If not specified, volumes are provisioned from the default storage availabity zone.
	// +optional
	StorageZone string `json:"storageZone,omitempty"`

	// Subnet is the UUID of the OpenStack subnet nodes will be provisioned on
	Subnet string `json:"subnet"`
}
