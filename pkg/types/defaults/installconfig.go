package defaults

import (
	netopv1 "github.com/openshift/cluster-network-operator/pkg/apis/networkoperator/v1"

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
	defaultClusterCIDR      = "10.128.0.0/14"
	defaultHostSubnetLength = 9 // equivalent to a /23 per node
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
		c.Networking.Type = netopv1.NetworkTypeOpenshiftSDN
	}
	if c.Networking.ServiceCIDR == nil {
		c.Networking.ServiceCIDR = defaultServiceCIDR
	}
	if len(c.Networking.ClusterNetworks) == 0 && c.Networking.PodCIDR == nil {
		c.Networking.ClusterNetworks = []netopv1.ClusterNetwork{
			{
				CIDR:             defaultClusterCIDR,
				HostSubnetLength: uint32(defaultHostSubnetLength),
			},
		}
	}
	if len(c.Machines) == 0 {
		numberOfMasters := int64(3)
		numberOfWorkers := int64(3)
		if c.Platform.Libvirt != nil {
			numberOfMasters = 1
			numberOfWorkers = 1
		}
		c.Machines = []types.MachinePool{
			{
				Name:     "master",
				Replicas: func(x int64) *int64 { return &x }(numberOfMasters),
			},
			{
				Name:     "worker",
				Replicas: func(x int64) *int64 { return &x }(numberOfWorkers),
			},
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
