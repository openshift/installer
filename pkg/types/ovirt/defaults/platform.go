package defaults

import (
	"github.com/openshift/installer/pkg/types/ovirt"
)

// DefaultNetworkName is the default network name to use in a cluster.
const DefaultNetworkName = "ovirtmgmt"

// DefaultControlPlaneAffinityGroupName is the default affinity group name for the control plane VMs.
const DefaultControlPlaneAffinityGroupName = "controlplane"

// DefaultComputeAffinityGroupName is the default affinity group name for the compute VMs.
const DefaultComputeAffinityGroupName = "compute"

func defaultControlPlaneAffinityGroup() ovirt.AffinityGroup {
	return ovirt.AffinityGroup{
		Name:        DefaultControlPlaneAffinityGroupName,
		Priority:    5,
		Description: "AffinityGroup for spreading each control plane machine to a different host",
		Enforcing:   true,
	}
}

func defaultComputeAffinityGroup() ovirt.AffinityGroup {
	return ovirt.AffinityGroup{
		Name:        DefaultComputeAffinityGroupName,
		Priority:    3,
		Description: "AffinityGroup for spreading each compute machine to a different host",
		Enforcing:   true,
	}
}

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *ovirt.Platform) {
	if p.NetworkName == "" {
		p.NetworkName = DefaultNetworkName
	}
	if p.AffinityGroups == nil {
		// No affinity group field, using the default settings
		p.AffinityGroups = []ovirt.AffinityGroup{
			defaultComputeAffinityGroup(),
			defaultControlPlaneAffinityGroup()}
	}
}
