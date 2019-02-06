package defaults

import (
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	libvirtdefaults "github.com/openshift/installer/pkg/types/libvirt/defaults"
	nonedefaults "github.com/openshift/installer/pkg/types/none/defaults"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

var (
	defaultMachineCIDR      = ipnet.MustParseCIDR("10.0.0.0/16")
	defaultServiceCIDR      = ipnet.MustParseCIDR("172.30.0.0/16")
	defaultClusterCIDR      = ipnet.MustParseCIDR("10.128.0.0/14")
	defaultHostSubnetLength = 9 // equivalent to a /23 per node
	defaultNetworkPlugin    = "OpenShiftSDN"
)

// SetInstallConfigDefaults sets the defaults for the install config.
func SetInstallConfigDefaults(c *types.InstallConfig) {
	if c.Networking == nil {
		c.Networking = &types.Networking{}
	}
	if c.Networking.MachineCIDR == nil {
		c.Networking.MachineCIDR = defaultMachineCIDR
		if c.Platform.Libvirt != nil {
			c.Networking.MachineCIDR = libvirtdefaults.DefaultMachineCIDR
		}
	}
	if c.Networking.Type == "" {
		c.Networking.Type = defaultNetworkPlugin
	}
	if c.Networking.ServiceCIDR == nil {
		c.Networking.ServiceCIDR = defaultServiceCIDR
	}
	if len(c.Networking.ClusterNetworks) == 0 {
		c.Networking.ClusterNetworks = []types.ClusterNetworkEntry{
			{
				CIDR:             *defaultClusterCIDR,
				HostSubnetLength: int32(defaultHostSubnetLength),
			},
		}
	}
	defaultReplicaCount := int64(3)
	if c.Platform.Libvirt != nil {
		defaultReplicaCount = 1
	}
	if c.ControlPlane == nil {
		c.ControlPlane = &types.MachinePool{
			Replicas: &defaultReplicaCount,
		}
	}
	c.ControlPlane.Name = "control-plane"
	if len(c.Compute) == 0 {
		c.Compute = []types.MachinePool{
			{
				Name:     "worker",
				Replicas: &defaultReplicaCount,
			},
		}
	}
	for i, p := range c.Compute {
		if p.Replicas == nil {
			c.Compute[i].Replicas = &defaultReplicaCount
		}
	}
	switch {
	case c.Platform.AWS != nil:
		awsdefaults.SetPlatformDefaults(c.Platform.AWS)
	case c.Platform.Libvirt != nil:
		libvirtdefaults.SetPlatformDefaults(c.Platform.Libvirt)
	case c.Platform.OpenStack != nil:
		openstackdefaults.SetPlatformDefaults(c.Platform.OpenStack)
	case c.Platform.None != nil:
		nonedefaults.SetPlatformDefaults(c.Platform.None)
	}
}
