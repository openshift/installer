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
func SetPlatformDefaults(p *openstack.Platform, installConfig *types.InstallConfig) {
	if p.Cloud == "" {
		p.Cloud = os.Getenv("OS_CLOUD")
		if p.Cloud == "" {
			p.Cloud = DefaultCloudName
		}
	}
        if p.AciNetExt.Mtu == "" {
                p.AciNetExt.Mtu = "1500"
        }
}

// APIVIP returns the internal virtual IP address (VIP) put in front
// of the Kubernetes API server for use by components inside the
// cluster. The DNS static pods running on the nodes resolve the
// api-int record to APIVIP.
func APIVIP(p *openstack.Platform) (net.IP, error) {
       neutronCIDR := &p.AciNetExt.NeutronCIDR.IPNet
       return cidr.Host(neutronCIDR, 5)
 }

// DNSVIP returns the internal virtual IP address (VIP) put in front
// of the DNS static pods running on the nodes. Unlike the DNS
// operator these services provide name resolution for the nodes
// themselves.
func DNSVIP(p *openstack.Platform) (net.IP, error) {
       neutronCIDR := &p.AciNetExt.NeutronCIDR.IPNet
       return cidr.Host(neutronCIDR, 6)
 }
// IngressVIP returns the internal virtual IP address (VIP) put in
// front of the OpenShift router pods. This provides the internal
// accessibility to the internal pods running on the worker nodes,
// e.g. `console`. The DNS static pods running on the nodes resolve
// the wildcard apps record to IngressVIP.
func IngressVIP(p *openstack.Platform) (net.IP, error) {
       neutronCIDR := &p.AciNetExt.NeutronCIDR.IPNet
       return cidr.Host(neutronCIDR, 7)
 }
