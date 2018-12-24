package openstack

import (
	"github.com/openshift/installer/pkg/ipnet"
)

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
	// Region specifies the OpenStack region where the cluster will be created.
	Region string `json:"region"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on OpenStack for machine pools which do not define their own
	// platform configuration.
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// NetworkCIDRBlock
	NetworkCIDRBlock ipnet.IPNet `json:"NetworkCIDRBlock"`

	// BaseImage
	// Name of image to use from OpenStack cloud
	BaseImage string `json:"baseImage"`

	// Cloud
	// Name of OpenStack cloud to use from clouds.yaml
	Cloud string `json:"cloud"`

	// ExternalNetwork
	// The OpenStack external network to be used for installation.
	ExternalNetwork string `json:"externalNetwork"`

	// FlavorName
	// The OpenStack compute flavor to use for servers.
	FlavorName string `json:"computeFlavor"`
}
