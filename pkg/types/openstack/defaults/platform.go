package defaults

import (
	"net"

	"github.com/apparentlymart/go-cidr/cidr"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *openstack.Platform) {
}

// APIVIP returns the internal virtual IP address (VIP) put in front
// of the Kubernetes API server for use by components inside the
// cluster. The DNS static pods running on the nodes resolve the
// api-int record to APIVIP.
func APIVIP(networking *types.Networking) (net.IP, error) {
	return cidr.Host(&networking.MachineCIDR.IPNet, 5)
}

// DNSVIP returns the internal virtual IP address (VIP) put in front
// of the DNS static pods running on the nodes. Unlike the DNS
// operator these services provide name resolution for the nodes
// themselves.
func DNSVIP(networking *types.Networking) (net.IP, error) {
	return cidr.Host(&networking.MachineCIDR.IPNet, 6)
}

// IngressVIP returns the internal virtual IP address (VIP) put in
// front of the OpenShift router pods. This provides the internal
// accessibility to the internal pods running on the worker nodes,
// e.g. `console`. The DNS static pods running on the nodes resolve
// the wildcard apps record to IngressVIP.
func IngressVIP(networking *types.Networking) (net.IP, error) {
	return cidr.Host(&networking.MachineCIDR.IPNet, 7)
}
