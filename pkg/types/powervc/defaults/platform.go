package defaults

import (
	"os"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/powervc"
)

const (
	// DefaultCloudName is the default name of the cloud in clouds.yaml file.
	DefaultCloudName = "powervc"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(pPvc *powervc.Platform, pO *openstack.Platform, n *types.Networking) {
	if pPvc.Cloud == "" {
		pPvc.Cloud = os.Getenv("OS_CLOUD")
		if pPvc.Cloud == "" {
			pPvc.Cloud = DefaultCloudName
		}
	}

	// Set the LB to User Managed
	pPvc.LoadBalancer = &configv1.OpenStackPlatformLoadBalancer{
		Type: configv1.LoadBalancerTypeUserManaged,
	}
	// For the LB, we also need to set the OpenStack default
	pO.LoadBalancer = &configv1.OpenStackPlatformLoadBalancer{
		Type: configv1.LoadBalancerTypeUserManaged,
	}
}
