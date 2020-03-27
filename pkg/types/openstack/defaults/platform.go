package defaults

import (
	"net"
	"os"

	"github.com/apparentlymart/go-cidr/cidr"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

const (
	// DefaultCloudName is the default name of the cloud in clouds.yaml file.
	DefaultCloudName = "openstack"
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
	if p.APIVIP == "" {
		vip, _ := cidr.Host(&n.MachineNetwork[0].CIDR.IPNet, 5)
		p.APIVIP = vip.String()
	}

	// IngressVIP returns the internal virtual IP address (VIP) put in
	// front of the OpenShift router pods. This provides the internal
	// accessibility to the internal pods running on the worker nodes,
	// e.g. `console`. The DNS static pods running on the nodes resolve
	// the wildcard apps record to IngressVIP.
	if p.IngressVIP == "" {
		vip, _ := cidr.Host(&n.MachineNetwork[0].CIDR.IPNet, 7)
		p.IngressVIP = vip.String()
	}
}

// DNSVIP returns the internal virtual IP address (VIP) put in front
// of the DNS static pods running on the nodes. Unlike the DNS
// operator these services provide name resolution for the nodes
// themselves.
func DNSVIP(networking *types.Networking) (net.IP, error) {
	return cidr.Host(&networking.MachineNetwork[0].CIDR.IPNet, 6)
}
