package defaults

import (
	operv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	"github.com/openshift/installer/pkg/types/azure"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
	baremetaldefaults "github.com/openshift/installer/pkg/types/baremetal/defaults"
	gcpdefaults "github.com/openshift/installer/pkg/types/gcp/defaults"
	ibmclouddefaults "github.com/openshift/installer/pkg/types/ibmcloud/defaults"
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
	defaultNetworkType    = string(operv1.NetworkTypeOpenShiftSDN)
	defaultOKDNetworkType = string(operv1.NetworkTypeOVNKubernetes)
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
		if c.IsOKD() {
			c.Networking.NetworkType = defaultOKDNetworkType
		} else {
			c.Networking.NetworkType = defaultNetworkType
		}
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

	if c.CredentialsMode == "" {
		if c.Platform.Azure != nil && c.Platform.Azure.CloudName == azure.StackCloud {
			c.CredentialsMode = types.ManualCredentialsMode
		}
	}

	switch {
	case c.Platform.AWS != nil:
		awsdefaults.SetPlatformDefaults(c.Platform.AWS)
	case c.Platform.Azure != nil:
		azuredefaults.SetPlatformDefaults(c.Platform.Azure)
	case c.Platform.GCP != nil:
		gcpdefaults.SetPlatformDefaults(c.Platform.GCP)
	case c.Platform.IBMCloud != nil:
		ibmclouddefaults.SetPlatformDefaults(c.Platform.IBMCloud)
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
		ovirtdefaults.SetControlPlaneDefaults(c.Platform.Ovirt, c.ControlPlane)
		for i := range c.Compute {
			ovirtdefaults.SetComputeDefaults(c.Platform.Ovirt, &c.Compute[i])
		}
	case c.Platform.None != nil:
		nonedefaults.SetPlatformDefaults(c.Platform.None)
	}
}
