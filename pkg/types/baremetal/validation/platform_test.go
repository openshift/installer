package validation

import (
	"net"
	"testing"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidatePlatform(t *testing.T) {
	iface, _ := net.Interfaces()
	network := &types.Networking{MachineCIDR: ipnet.MustParseCIDR("192.168.111.0/24")}
	cases := []struct {
		name     string
		platform *baremetal.Platform
		network  *types.Networking
		expected string
	}{
		{
			name: "valid",
			platform: &baremetal.Platform{
				APIVIP:                  "192.168.111.2",
				DNSVIP:                  "192.168.111.3",
				IngressVIP:              "192.168.111.4",
				Hosts:                   []*baremetal.Host{},
				LibvirtURI:              "qemu://system",
				ClusterProvisioningIP:   "172.22.0.3",
				BootstrapProvisioningIP: "172.22.0.2",
				ExternalBridge:          iface[0].Name,
				ProvisioningBridge:      iface[0].Name,
			},
			network: network,
		},
		{
			name: "invalid_apivip",
			platform: &baremetal.Platform{
				APIVIP:                  "192.168.222.2",
				DNSVIP:                  "192.168.111.3",
				IngressVIP:              "192.168.111.4",
				Hosts:                   []*baremetal.Host{},
				LibvirtURI:              "qemu://system",
				ClusterProvisioningIP:   "172.22.0.3",
				BootstrapProvisioningIP: "172.22.0.2",
				ExternalBridge:          iface[0].Name,
				ProvisioningBridge:      iface[0].Name,
			},
			network:  network,
			expected: "Invalid value: \"192.168.222.2\": the virtual IP is expected to be in 192.168.111.0/24 subnet",
		},
		{
			name: "invalid_dnsvip",
			platform: &baremetal.Platform{
				APIVIP:                  "192.168.111.2",
				DNSVIP:                  "192.168.222.3",
				IngressVIP:              "192.168.111.4",
				Hosts:                   []*baremetal.Host{},
				LibvirtURI:              "qemu://system",
				ClusterProvisioningIP:   "172.22.0.3",
				BootstrapProvisioningIP: "172.22.0.2",
				ExternalBridge:          iface[0].Name,
				ProvisioningBridge:      iface[0].Name,
			},
			network:  network,
			expected: "Invalid value: \"192.168.222.3\": the virtual IP is expected to be in 192.168.111.0/24 subnet",
		},
		{
			name: "invalid_ingressvip",
			platform: &baremetal.Platform{
				APIVIP:                  "192.168.111.2",
				DNSVIP:                  "192.168.111.3",
				IngressVIP:              "192.168.222.4",
				Hosts:                   []*baremetal.Host{},
				LibvirtURI:              "qemu://system",
				ClusterProvisioningIP:   "172.22.0.3",
				BootstrapProvisioningIP: "172.22.0.2",
				ExternalBridge:          iface[0].Name,
				ProvisioningBridge:      iface[0].Name,
			},
			network:  network,
			expected: "Invalid value: \"192.168.222.4\": the virtual IP is expected to be in 192.168.111.0/24 subnet",
		},
		{
			name: "invalid_hosts",
			platform: &baremetal.Platform{
				APIVIP:                  "192.168.111.2",
				DNSVIP:                  "192.168.111.3",
				IngressVIP:              "192.168.111.4",
				Hosts:                   nil,
				LibvirtURI:              "qemu://system",
				ClusterProvisioningIP:   "172.22.0.3",
				BootstrapProvisioningIP: "172.22.0.2",
				ExternalBridge:          iface[0].Name,
				ProvisioningBridge:      iface[0].Name,
			},
			network:  network,
			expected: "bare metal hosts are missing",
		},
		{
			name: "invalid_libvirturi",
			platform: &baremetal.Platform{
				APIVIP:                  "192.168.111.2",
				DNSVIP:                  "192.168.111.3",
				IngressVIP:              "192.168.111.4",
				Hosts:                   []*baremetal.Host{},
				LibvirtURI:              "",
				ClusterProvisioningIP:   "172.22.0.3",
				BootstrapProvisioningIP: "172.22.0.2",
				ExternalBridge:          iface[0].Name,
				ProvisioningBridge:      iface[0].Name,
			},
			network:  network,
			expected: "invalid URI \"\"",
		},
		{
			name: "invalid_extbridge",
			platform: &baremetal.Platform{
				APIVIP:                  "192.168.111.2",
				DNSVIP:                  "192.168.111.3",
				IngressVIP:              "192.168.111.4",
				Hosts:                   []*baremetal.Host{},
				LibvirtURI:              "qemu://system",
				ClusterProvisioningIP:   "172.22.0.3",
				BootstrapProvisioningIP: "172.22.0.2",
				ExternalBridge:          "noexist",
				ProvisioningBridge:      iface[0].Name,
			},
			network:  network,
			expected: "Invalid value: \"noexist\": noexist is not a valid network interface",
		},
		{
			name: "invalid_provbridge",
			platform: &baremetal.Platform{
				APIVIP:                  "192.168.111.2",
				DNSVIP:                  "192.168.111.3",
				IngressVIP:              "192.168.111.4",
				Hosts:                   []*baremetal.Host{},
				LibvirtURI:              "qemu://system",
				ClusterProvisioningIP:   "172.22.0.3",
				BootstrapProvisioningIP: "172.22.0.2",
				ExternalBridge:          iface[0].Name,
				ProvisioningBridge:      "noexist",
			},
			network:  network,
			expected: "Invalid value: \"noexist\": noexist is not a valid network interface",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, tc.network, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}
