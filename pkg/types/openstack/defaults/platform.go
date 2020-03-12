package defaults

import (
	"net"
	"os"

	"github.com/apparentlymart/go-cidr/cidr"
        "github.com/openshift/installer/pkg/ipnet"
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
        // Panic if input neutron CIDR and machine CIDR don't have equal masks.
        // If no neutron CIDR provided, set it to 192.168.0.0 with the machine CIDR mask
        machineMask := installConfig.Networking.DeprecatedMachineCIDR.Mask
        if p.NeutronCIDR.String() != "" {
                neutronMask := p.NeutronCIDR.Mask
                if machineMask.String() != neutronMask.String() {
                        panic("Machine CIDR and Neutron CIDR have different subnet masks")
                        
                }
        } else {
                neutronIP := net.ParseIP("192.168.0.0")
                p.NeutronCIDR = &ipnet.IPNet{
                                        IPNet: net.IPNet{
                                                IP:   neutronIP,
                                                Mask: machineMask,
                                        },
                                }
        }
        installConfig.Networking.NeutronCIDR = p.NeutronCIDR

        ipnet.MustParseCIDR(p.InstallerHostSubnet)
}

// APIVIP returns the internal virtual IP address (VIP) put in front
// of the Kubernetes API server for use by components inside the
// cluster. The DNS static pods running on the nodes resolve the
// api-int record to APIVIP.
func APIVIP(networking *types.Networking) (net.IP, error) {
       neutronCIDR := &networking.NeutronCIDR.IPNet
       return cidr.Host(neutronCIDR, 5)
 }

// DNSVIP returns the internal virtual IP address (VIP) put in front
// of the DNS static pods running on the nodes. Unlike the DNS
// operator these services provide name resolution for the nodes
// themselves.
func DNSVIP(networking *types.Networking) (net.IP, error) {
       neutronCIDR := &networking.NeutronCIDR.IPNet
       return cidr.Host(neutronCIDR, 6)
 }
// IngressVIP returns the internal virtual IP address (VIP) put in
// front of the OpenShift router pods. This provides the internal
// accessibility to the internal pods running on the worker nodes,
// e.g. `console`. The DNS static pods running on the nodes resolve
// the wildcard apps record to IngressVIP.
func IngressVIP(networking *types.Networking) (net.IP, error) {
       neutronCIDR := &networking.NeutronCIDR.IPNet
       return cidr.Host(neutronCIDR, 7)
 }
