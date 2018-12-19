package aws

import (
	"github.com/openshift/installer/pkg/ipnet"
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
