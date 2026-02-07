package defaults

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/network"
)

const (
	defaultInstanceSizeHighAvailabilityTopology = "xlarge"
	defaultInstanceSizeSingleReplicaTopology    = "2xlarge"
)

var (
	defaultMachineTypes = map[types.Architecture]map[string][]string{
		types.ArchitectureAMD64: {
			// Example region default machine class override for AMD64:
			// "ap-east-1":      {"m6i.xlarge", "m5.xlarge"},
		},
		types.ArchitectureARM64: {
			// Example region default machine class override for ARM64:
			// "us-east-1":      {"m6g.xlarge", "m6gd.xlarge"},
		},
	}
	// Skip undesired zones in the discovery.
	// - us-east-1e is a well-known limited zone, with limited offering of
	// 	 instance types supported by installer.
	skippedZones = []string{"us-east-1e"}
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *aws.Platform) {
	if p.IPFamily == "" {
		p.IPFamily = network.IPv4
		logrus.Infof("ipFamily is not specified in install-config; defaulting to %q", network.IPv4)
	}

	if p.LBType == "" {
		if p.IPFamily.DualStackEnabled() {
			p.LBType = configv1.NLB
		}
	}
}

// InstanceTypes returns a list of instance types, in decreasing priority order, which we should use for a given
// region. Default is m6i.xlarge, m5.xlarge, lastly c5d.2xlarge unless a region override
// is defined in defaultMachineTypes.
// c5d.2xlarge is in the most locations of availability for Local Zone offerings.
// https://aws.amazon.com/about-aws/global-infrastructure/localzones/features
// https://aws.amazon.com/ec2/pricing/on-demand/
func InstanceTypes(region string, arch types.Architecture, topology configv1.TopologyMode) []string {
	if classesForArch, ok := defaultMachineTypes[arch]; ok {
		if classes, ok := classesForArch[region]; ok {
			return classes
		}
	}

	instanceSize := defaultInstanceSizeHighAvailabilityTopology
	// If the control plane is single node, we need to use a larger
	// instance type for that node, as the minimum requirement for
	// single-node control-plane nodes is 8 cores, and xlarge only has
	// 4. Unfortunately 2xlarge has twice as much RAM as we need, but
	// we default to it because AWS doesn't offer an 8-core 16GiB
	// instance type
	if topology == configv1.SingleReplicaTopologyMode {
		instanceSize = defaultInstanceSizeSingleReplicaTopology
	}

	switch arch {
	case types.ArchitectureARM64:
		return []string{
			fmt.Sprintf("m6g.%s", instanceSize),
		}
	default:
		return []string{
			fmt.Sprintf("m6i.%s", instanceSize),
			fmt.Sprintf("m5.%s", instanceSize),
			// For Local Zone compatibility
			fmt.Sprintf("r5.%s", instanceSize),
			"c5.2xlarge",
			"m5.2xlarge",
			"c5d.2xlarge",
			"r5.2xlarge",
		}
	}
}

// SupportedZones returns the list of supported zones.
func SupportedZones(zones []string) []string {
	return sets.List(sets.New(zones...).Difference(sets.New(skippedZones...)))
}
