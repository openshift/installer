package aws

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// AMIID is the AMI that should be used to boot machines for the cluster.
	// If set, the AMI should belong to the same region as the cluster.
	AMIID string `json:"amiID,omitempty"`

	// Region specifies the AWS region where the cluster will be created.
	Region string `json:"region"`

	// UserTags specifies additional tags for AWS resources created for the cluster.
	// +optional
	UserTags map[string]string `json:"userTags,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on AWS for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// CustomRegionOverride list contains custom endpoints which will overridde default
	// service endpoint of AWS Services.
	// +optional
	CustomRegionOverride []CustomEndpoint `json:"customRegionOverride,omitempty"`
}

// FetchCustomEndpointsMap converts List of custom endpoints to a map
func (p *Platform) FetchCustomEndpointsMap() map[string]string {
	endpointMap := make(map[string]string)
	if len(p.CustomRegionOverride) == 0 {
		return nil
	}
	for _, endpoint := range p.CustomRegionOverride {
		endpointMap[endpoint.Service] = endpoint.URL
	}
	return endpointMap
}
