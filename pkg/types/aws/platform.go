package aws

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// AMIID is the AMI that should be used to boot machines for the cluster.
	// If set, the AMI should belong to the same region as the cluster.
	AMIID string `json:"amiID,omitempty"`

	// Region specifies the AWS region where the cluster will be created.
	Region string `json:"region"`

	// PublicZoneID specifies the ID for the public DNS Zone
	PublicZoneID string `json:"publicZoneID"`

	// UserTags specifies additional tags for AWS resources created for the cluster.
	// +optional
	UserTags map[string]string `json:"userTags,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on AWS for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
}

//SetBaseDomain parses the baseDomainID and sets the related fields on azure.Platform
func (p *Platform) SetBaseDomain(baseDomainID string) error {
	p.PublicZoneID = baseDomainID
	return nil
}
