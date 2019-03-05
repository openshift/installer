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
	defaultMachineCIDR    = ipnet.MustParseCIDR("10.0.0.0/16")
	defaultServiceNetwork = ipnet.MustParseCIDR("172.30.0.0/16")
	defaultClusterNetwork = ipnet.MustParseCIDR("10.128.0.0/14")
	defaultHostPrefix     = 23
	defaultNetworkType    = "OpenShiftSDN"
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
	if c.Networking.NetworkType == "" {
		c.Networking.NetworkType = defaultNetworkType
	}
	if len(c.Networking.ServiceNetwork) == 0 {
		c.Networking.ServiceNetwork = []ipnet.IPNet{*defaultServiceNetwork}
	}
	if len(c.Networking.ClusterNetwork) == 0 {
		c.Networking.ClusterNetwork = []types.ClusterNetworkEntry{
			{
				CIDR:       *defaultClusterNetwork,
				HostPrefix: int32(defaultHostPrefix),
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
	c.ControlPlane.Name = "master"
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
