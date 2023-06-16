package defaults

import (
	"fmt"
	"os"

	"github.com/apparentlymart/go-cidr/cidr"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

const (
	// DefaultCloudName is the default name of the cloud in clouds.yaml file.
	DefaultCloudName = "openstack"
	// DualStackVIPsPortTag is the identifier of VIPs Port with dual-stack addresses.
	DualStackVIPsPortTag = "-dual-stack-vips-port"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *openstack.Platform, n *types.Networking) {
	if p.Cloud == "" {
		p.Cloud = os.Getenv("OS_CLOUD")
		if p.Cloud == "" {
			p.Cloud = DefaultCloudName
		}
	}
	// APIVIP returns the internal virtual IP address (VIP) put in front
	// of the Kubernetes API server for use by components inside the
	// cluster. The DNS static pods running on the nodes resolve the
	// api-int record to APIVIP.
	if len(p.APIVIPs) == 0 && p.DeprecatedAPIVIP == "" {
		vip, err := cidr.Host(&n.MachineNetwork[0].CIDR.IPNet, 5)
		if err != nil {
			// This will fail validation and abort the install
			p.APIVIPs = []string{fmt.Sprintf("could not derive API VIP from machine networks: %s", err.Error())}
		} else {
			p.APIVIPs = []string{vip.String()}
		}
	}

	// IngressVIP returns the internal virtual IP address (VIP) put in
	// front of the OpenShift router pods. This provides the internal
	// accessibility to the internal pods running on the worker nodes,
	// e.g. `console`. The DNS static pods running on the nodes resolve
	// the wildcard apps record to IngressVIP.
	if len(p.IngressVIPs) == 0 && p.DeprecatedIngressVIP == "" {
		vip, err := cidr.Host(&n.MachineNetwork[0].CIDR.IPNet, 7)
		if err != nil {
			// This will fail validation and abort the install
			p.IngressVIPs = []string{fmt.Sprintf("could not derive Ingress VIP from machine networks: %s", err.Error())}
		} else {
			p.IngressVIPs = []string{vip.String()}
		}
	}
}
