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
func SetPlatformDefaults(p *openstack.Platform) {
	if p.Cloud == "" {
		p.Cloud = os.Getenv("OS_CLOUD")
		if p.Cloud == "" {
			p.Cloud = DefaultCloudName
		}
	}
}

// APIVIP returns the internal virtual IP address (VIP) put in front
// of the Kubernetes API server for use by components inside the
// cluster. The DNS static pods running on the nodes resolve the
// api-int record to APIVIP.
func APIVIP(networking *types.Networking) (net.IP, error) {
       return cidr.Host(&net.IPNet{IP: net.IP{192,168,0,0},
       Mask: net.IPMask{255,255,0,0}},5)
 }

// DNSVIP returns the internal virtual IP address (VIP) put in front
// of the DNS static pods running on the nodes. Unlike the DNS
// operator these services provide name resolution for the nodes
// themselves.
func DNSVIP(networking *types.Networking) (net.IP, error) {
       return cidr.Host(&net.IPNet{IP: net.IP{192,168,0,0},
               Mask: net.IPMask{255,255,0,0}},6)
 }
// IngressVIP returns the internal virtual IP address (VIP) put in
// front of the OpenShift router pods. This provides the internal
// accessibility to the internal pods running on the worker nodes,
// e.g. `console`. The DNS static pods running on the nodes resolve
// the wildcard apps record to IngressVIP.
func IngressVIP(networking *types.Networking) (net.IP, error) {
       return cidr.Host(&net.IPNet{IP: net.IP{192,168,0,0},
              Mask: net.IPMask{255,255,0,0}},7)
 }
