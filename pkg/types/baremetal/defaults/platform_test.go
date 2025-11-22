package defaults

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

const testClusterName = "test-cluster"

var host = &baremetal.Host{
	Name: "test-host",
}

func TestSetPlatformDefaults(t *testing.T) {
	// Stub the call to net.LookupHost
	lookupHost = func(host string) (addrs []string, err error) {
		if host == "api.test-cluster.test" {
			ips := []string{"192.168.111.2"}
			return ips, nil
		} else if host == "test.apps.test-cluster.test" {
			ips := []string{"192.168.111.3"}
			return ips, nil
		} else {
			return nil, errors.New("Unknown Host " + host)
		}
	}

	machineNetwork := ipnet.MustParseCIDR("192.168.111.0/24")

	cases := []struct {
		name         string
		platform     *baremetal.Platform
		expected     *baremetal.Platform
		expectedHost *baremetal.Host
	}{
		{
			name:     "default_empty",
			platform: &baremetal.Platform{},
			expected: &baremetal.Platform{
				LibvirtURI:              "qemu:///system",
				ClusterProvisioningIP:   "172.22.0.3",
				BootstrapProvisioningIP: "172.22.0.2",
				ExternalBridge:          "baremetal",
				ProvisioningBridge:      "provisioning",
				ProvisioningNetwork:     baremetal.ManagedProvisioningNetwork,
				APIVIPs:                 []string{"192.168.111.2"},
				IngressVIPs:             []string{"192.168.111.3"},
				ProvisioningNetworkCIDR: ipnet.MustParseCIDR("172.22.0.0/24"),
				ProvisioningDHCPRange:   "172.22.0.10,172.22.0.254",
				DNSRecordsType:          configv1.DNSRecordsTypeInternal,
			},
		},
		{
			name: "alternate_cidr",
			platform: &baremetal.Platform{
				ProvisioningNetworkCIDR: ipnet.MustParseCIDR("172.23.0.0/24"),
			},
			expected: &baremetal.Platform{
				LibvirtURI:              "qemu:///system",
				ClusterProvisioningIP:   "172.23.0.3",
				BootstrapProvisioningIP: "172.23.0.2",
				ExternalBridge:          "baremetal",
				ProvisioningBridge:      "provisioning",
				ProvisioningNetwork:     baremetal.ManagedProvisioningNetwork,
				APIVIPs:                 []string{"192.168.111.2"},
				IngressVIPs:             []string{"192.168.111.3"},
				ProvisioningNetworkCIDR: ipnet.MustParseCIDR("172.23.0.0/24"),
				ProvisioningDHCPRange:   "172.23.0.10,172.23.0.254",
				DNSRecordsType:          configv1.DNSRecordsTypeInternal,
			},
		},
		{
			name: "alternate_cidr_ipv6",
			platform: &baremetal.Platform{
				ProvisioningNetworkCIDR: ipnet.MustParseCIDR("fd2e:6f44:5dd8:b856::/64"),
			},
			expected: &baremetal.Platform{
				LibvirtURI:              "qemu:///system",
				ClusterProvisioningIP:   "fd2e:6f44:5dd8:b856::3",
				BootstrapProvisioningIP: "fd2e:6f44:5dd8:b856::2",
				ExternalBridge:          "baremetal",
				ProvisioningBridge:      "provisioning",
				ProvisioningNetwork:     baremetal.ManagedProvisioningNetwork,
				APIVIPs:                 []string{"192.168.111.2"},
				IngressVIPs:             []string{"192.168.111.3"},
				ProvisioningNetworkCIDR: ipnet.MustParseCIDR("fd2e:6f44:5dd8:b856::/64"),
				ProvisioningDHCPRange:   "fd2e:6f44:5dd8:b856::a,fd2e:6f44:5dd8:b856:ffff:ffff:ffff:fffe",
				DNSRecordsType:          configv1.DNSRecordsTypeInternal,
			},
		},
		{
			name: "alternate_cidr_dhcp_disabled",
			platform: &baremetal.Platform{
				ProvisioningNetworkCIDR: ipnet.MustParseCIDR("172.23.0.0/24"),
				ProvisioningNetwork:     baremetal.UnmanagedProvisioningNetwork,
			},
			expected: &baremetal.Platform{
				LibvirtURI:              "qemu:///system",
				ClusterProvisioningIP:   "172.23.0.3",
				BootstrapProvisioningIP: "172.23.0.2",
				ExternalBridge:          "baremetal",
				ProvisioningBridge:      "provisioning",
				APIVIPs:                 []string{"192.168.111.2"},
				IngressVIPs:             []string{"192.168.111.3"},
				ProvisioningNetworkCIDR: ipnet.MustParseCIDR("172.23.0.0/24"),
				ProvisioningNetwork:     baremetal.UnmanagedProvisioningNetwork,
				DNSRecordsType:          configv1.DNSRecordsTypeInternal,
			},
		},
		{
			name: "disabled_provisioning_network",
			platform: &baremetal.Platform{
				ProvisioningNetwork:     baremetal.DisabledProvisioningNetwork,
				BootstrapProvisioningIP: "192.168.111.7",
				ClusterProvisioningIP:   "192.168.111.8",
			},
			expected: &baremetal.Platform{
				BootstrapProvisioningIP: "192.168.111.7",
				ClusterProvisioningIP:   "192.168.111.8",
				LibvirtURI:              "qemu:///system",
				ExternalBridge:          "baremetal",
				ProvisioningBridge:      "",
				ProvisioningNetwork:     baremetal.DisabledProvisioningNetwork,
				ProvisioningNetworkCIDR: machineNetwork,
				APIVIPs:                 []string{"192.168.111.2"},
				IngressVIPs:             []string{"192.168.111.3"},
				DNSRecordsType:          configv1.DNSRecordsTypeInternal,
			},
		},
		{
			name: "disabled_provisioning_network_no_bootstrap_ip",
			platform: &baremetal.Platform{
				ProvisioningNetwork:   baremetal.DisabledProvisioningNetwork,
				ClusterProvisioningIP: "192.168.111.8",
			},
			expected: &baremetal.Platform{
				BootstrapProvisioningIP: "",
				ClusterProvisioningIP:   "192.168.111.8",
				LibvirtURI:              "qemu:///system",
				ExternalBridge:          "baremetal",
				ProvisioningBridge:      "",
				ProvisioningNetwork:     baremetal.DisabledProvisioningNetwork,
				ProvisioningNetworkCIDR: machineNetwork,
				APIVIPs:                 []string{"192.168.111.2"},
				IngressVIPs:             []string{"192.168.111.3"},
				DNSRecordsType:          configv1.DNSRecordsTypeInternal,
			},
		},
		{
			name: "defaults_for_hosts",
			platform: &baremetal.Platform{
				Hosts: []*baremetal.Host{
					host,
				},
			},
			expectedHost: &baremetal.Host{
				Name:            "test-host",
				BootMode:        "UEFI",
				HardwareProfile: "default",
			},
		},
		{
			name: "deprecated VIPs fields set",
			platform: &baremetal.Platform{
				DeprecatedAPIVIP:     "10.0.0.1",
				DeprecatedIngressVIP: "10.0.0.2",
			},
			expected: &baremetal.Platform{
				LibvirtURI:              "qemu:///system",
				ClusterProvisioningIP:   "172.22.0.3",
				BootstrapProvisioningIP: "172.22.0.2",
				ExternalBridge:          "baremetal",
				ProvisioningBridge:      "provisioning",
				ProvisioningNetwork:     baremetal.ManagedProvisioningNetwork,
				ProvisioningNetworkCIDR: ipnet.MustParseCIDR("172.22.0.0/24"),
				ProvisioningDHCPRange:   "172.22.0.10,172.22.0.254",
				DeprecatedAPIVIP:        "10.0.0.1",
				DeprecatedIngressVIP:    "10.0.0.2",
				APIVIPs:                 nil,
				IngressVIPs:             nil,
				DNSRecordsType:          configv1.DNSRecordsTypeInternal,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			ic := &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: testClusterName,
				},
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{
							CIDR: ipnet.IPNet{IPNet: machineNetwork.IPNet},
						},
					},
				},
				BaseDomain: "test",
			}
			SetPlatformDefaults(tc.platform, ic)
			if tc.expected != nil {
				//Check Generated Provisioning MAC address
				provisioningmac, err := net.ParseMAC(tc.platform.ProvisioningMACAddress)
				if err != nil {
					assert.Fail(t, "Invalid Provisioning MAC %s", string(provisioningmac), err)
				} else {
					tc.expected.ProvisioningMACAddress = tc.platform.ProvisioningMACAddress
				}
				//Check Generated External MAC address
				externalmac, err := net.ParseMAC(tc.platform.ExternalMACAddress)
				if err != nil {
					assert.Fail(t, "Invalid External MAC %s", string(externalmac), err)
				} else {
					tc.expected.ExternalMACAddress = tc.platform.ExternalMACAddress
				}

				assert.Equal(t, tc.expected, tc.platform, "unexpected platform")
			}

			if tc.expectedHost != nil {
				assert.Contains(t, tc.platform.Hosts, tc.expectedHost, "expected host not found")
			}
		})
	}
}

func TestBaremetalHostsSortByRole(t *testing.T) {
	cases := []struct {
		name          string
		hosts         []*baremetal.Host
		expectedHosts []string
	}{
		{
			name: "default",
			hosts: []*baremetal.Host{
				{Name: "master-0", Role: "master"},
				{Name: "master-1", Role: "master"},
				{Name: "master-2", Role: "master"},
				{Name: "worker-0", Role: "worker"},
				{Name: "worker-1", Role: "worker"},
			},
			expectedHosts: []string{
				"master-0", "master-1", "master-2", "worker-0", "worker-1",
			},
		},
		{
			name: "norole",
			hosts: []*baremetal.Host{
				{Name: "master-0"},
				{Name: "master-1"},
				{Name: "master-2"},
				{Name: "worker-0"},
				{Name: "worker-1"},
			},
			expectedHosts: []string{
				"master-0", "master-1", "master-2", "worker-0", "worker-1",
			},
		},
		{
			name: "mixed",
			hosts: []*baremetal.Host{
				{Name: "worker-0", Role: "worker"},
				{Name: "master-0", Role: "master"},
				{Name: "worker-1", Role: "worker"},
				{Name: "master-1", Role: "master"},
				{Name: "master-2", Role: "master"},
			},
			expectedHosts: []string{
				"master-0", "master-1", "master-2", "worker-0", "worker-1",
			},
		},
		{
			name: "mixed-norole",
			hosts: []*baremetal.Host{
				{Name: "worker-0", Role: "worker"},
				{Name: "master-0", Role: "master"},
				{Name: "worker-1", Role: ""},
				{Name: "master-1", Role: "master"},
				{Name: "master-2", Role: "master"},
			},
			expectedHosts: []string{
				"master-0", "master-1", "master-2", "worker-0", "worker-1",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			ic := &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: testClusterName,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						Hosts: tc.hosts,
					},
				},
				BaseDomain: "test",
			}
			SetPlatformDefaults(ic.Platform.BareMetal, ic)

			for i, h := range ic.Platform.BareMetal.Hosts {
				assert.Equal(t, h.Name, tc.expectedHosts[i])
			}
		})
	}
}
