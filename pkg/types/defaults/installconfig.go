package defaults

import (
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
	baremetaldefaults "github.com/openshift/installer/pkg/types/baremetal/defaults"
	gcpdefaults "github.com/openshift/installer/pkg/types/gcp/defaults"
	libvirtdefaults "github.com/openshift/installer/pkg/types/libvirt/defaults"
	nonedefaults "github.com/openshift/installer/pkg/types/none/defaults"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
	ovirtdefaults "github.com/openshift/installer/pkg/types/ovirt/defaults"
	vspheredefaults "github.com/openshift/installer/pkg/types/vsphere/defaults"
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
	if len(c.Networking.MachineNetwork) == 0 {
		c.Networking.MachineNetwork = []types.MachineNetworkEntry{
			{CIDR: *defaultMachineCIDR},
		}
		if c.Platform.Libvirt != nil {
			c.Networking.MachineNetwork = []types.MachineNetworkEntry{
				{CIDR: *libvirtdefaults.DefaultMachineCIDR},
			}
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

	if c.Publish == "" {
		c.Publish = types.ExternalPublishingStrategy
	}

	if c.ControlPlane == nil {
		c.ControlPlane = &types.MachinePool{}
	}
	c.ControlPlane.Name = "master"
	SetMachinePoolDefaults(c.ControlPlane, c.Platform.Name())
	if len(c.Compute) == 0 {
		c.Compute = []types.MachinePool{{Name: "worker"}}
	}
	for i := range c.Compute {
		SetMachinePoolDefaults(&c.Compute[i], c.Platform.Name())
	}
	switch {
	case c.Platform.AWS != nil:
		awsdefaults.SetPlatformDefaults(c.Platform.AWS)
	case c.Platform.Azure != nil:
		azuredefaults.SetPlatformDefaults(c.Platform.Azure)
	case c.Platform.GCP != nil:
		gcpdefaults.SetPlatformDefaults(c.Platform.GCP)
	case c.Platform.Libvirt != nil:
		libvirtdefaults.SetPlatformDefaults(c.Platform.Libvirt)
	case c.Platform.OpenStack != nil:
		openstackdefaults.SetPlatformDefaults(c.Platform.OpenStack, c.Networking)
	case c.Platform.VSphere != nil:
		vspheredefaults.SetPlatformDefaults(c.Platform.VSphere, c)
	case c.Platform.BareMetal != nil:
		baremetaldefaults.SetPlatformDefaults(c.Platform.BareMetal, c)
	case c.Platform.Ovirt != nil:
		ovirtdefaults.SetPlatformDefaults(c.Platform.Ovirt)
	case c.Platform.None != nil:
		nonedefaults.SetPlatformDefaults(c.Platform.None)
	}
}
