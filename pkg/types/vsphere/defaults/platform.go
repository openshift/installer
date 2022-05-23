package defaults

import (
	"fmt"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *vsphere.Platform, installConfig *types.InstallConfig) {

	// We need to deploy templates (OVA) via DeploymentZones
	// since we could have compute (workers) in those zones
	// but _not_ control plane nodes. If the placementConstraints
	// are not defined we must use the default for the datacenter
	// and cluster.
	for i, v := range p.DeploymentZones {
		var failureDomain vsphere.FailureDomainSpec
		for _, f := range p.FailureDomains {
			if v.FailureDomain == f.Name {
				failureDomain = f
			}
		}

		if v.PlacementConstraint.ResourcePool == "" {
			p.DeploymentZones[i].PlacementConstraint.ResourcePool = fmt.Sprintf("%s/%s", failureDomain.Topology.ComputeCluster, "/Resources")
		}
	}
}
