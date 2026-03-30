package defaults

import (
	"path"

	configv1 "github.com/openshift/api/config/v1"
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
	for i := range p.FailureDomains {
		if p.FailureDomains[i].Topology.ResourcePool == "" && p.FailureDomains[i].Topology.ComputeCluster != "" {
			p.FailureDomains[i].Topology.ResourcePool = path.Join(p.FailureDomains[i].Topology.ComputeCluster, "Resources")
		}
	}

	if p.DNSRecordsType == "" {
		p.DNSRecordsType = configv1.DNSRecordsTypeInternal
	}
}
