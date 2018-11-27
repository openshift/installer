package aws

import (
	"github.com/openshift/installer/pkg/ipnet"
)

var (
	// ValidRegions is a map of the known AWS regions. The key of the map is
	// the short name of the region. The value of the map is the long name of
	// the region.
	ValidRegions = map[string]string{
		"ap-northeast-1": "Tokyo",
		"ap-northeast-2": "Seoul",
		"ap-northeast-3": "Osaka-Local",
		"ap-south-1":     "Mumbai",
		"ap-southeast-1": "Singapore",
		"ap-southeast-2": "Sydney",
		"ca-central-1":   "Central",
		"cn-north-1":     "Beijing",
		"cn-northwest-1": "Ningxia",
		"eu-central-1":   "Frankfurt",
		"eu-west-1":      "Ireland",
		"eu-west-2":      "London",
		"eu-west-3":      "Paris",
		"sa-east-1":      "São Paulo",
		"us-east-1":      "N. Virginia",
		"us-east-2":      "Ohio",
		"us-west-1":      "N. California",
		"us-west-2":      "Oregon",
	}
)

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// Region specifies the AWS region where the cluster will be created.
	Region string `json:"region"`

	// UserTags specifies additional tags for AWS resources created for the cluster.
	UserTags map[string]string `json:"userTags,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on AWS for machine pools which do not define their own
	// platform configuration.
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// VPCCIDRBlock
	// +optional
	VPCCIDRBlock *ipnet.IPNet `json:"vpcCIDRBlock,omitempty"`
}
