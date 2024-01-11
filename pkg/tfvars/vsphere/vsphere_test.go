package vsphere

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1alpha1"

	"github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

const ipv4Gateway = "192.168.101.200"
const ipv6Gateway = "2001::200"
const ipv4Template = "192.168.101.24%d/24"
const ipv6Template = "2001::24%d/64"

const (
	ipv4 = 0x01
	ipv6 = 0x02
)

func createTFVarsSources(createHosts bool, ipTypes int64) TFVarsSources {
	tvs := TFVarsSources{
		InstallConfig: &installconfig.InstallConfig{
			AssetBase: installconfig.AssetBase{
				Config: &types.InstallConfig{
					Platform: types.Platform{
						VSphere: &vsphere.Platform{},
					},
				},
			},
		},
	}

	if createHosts {
		tvs.InstallConfig.Config.VSphere.Hosts = createValidHosts(ipTypes)
		machines, cpc, ips := createControlPlaneConfigs(ipTypes)
		tvs.ControlPlaneMachines = machines
		tvs.ControlPlaneConfigs = cpc
		tvs.IPAddresses = ips
	}

	return tvs
}

func createValidHosts(ipTypes int64) []*vsphere.Host {
	hosts := []*vsphere.Host{
		{
			Role:          "bootstrap",
			NetworkDevice: &vsphere.NetworkDeviceSpec{},
		},
		// NOTE: All control plane and compute hosts info is missing since not needed for test.
	}

	if ipTypes&ipv4 != 0 {
		hosts[0].NetworkDevice.Gateway = ipv4Gateway
		hosts[0].NetworkDevice.Nameservers = append(hosts[0].NetworkDevice.Nameservers, "8.8.8.8")
		hosts[0].NetworkDevice.IPAddrs = append(hosts[0].NetworkDevice.IPAddrs, fmt.Sprintf(ipv4Template, 0))
	} else if ipTypes&ipv6 != 0 {
		hosts[0].NetworkDevice.Gateway = ipv6Gateway
		hosts[0].NetworkDevice.Nameservers = append(hosts[0].NetworkDevice.Nameservers, "2001::100")
	}

	if ipTypes&ipv6 != 0 {
		hosts[0].NetworkDevice.IPAddrs = append(hosts[0].NetworkDevice.IPAddrs, fmt.Sprintf(ipv6Template, 0))
	}

	return hosts
}

func createControlPlaneConfigs(ipTypes int64) ([]v1beta1.Machine, []*v1beta1.VSphereMachineProviderSpec, []ipamv1.IPAddress) {
	var machines []v1beta1.Machine
	var specs []*v1beta1.VSphereMachineProviderSpec
	var ips []ipamv1.IPAddress

	for i := 1; i <= 3; i++ {
		machine := v1beta1.Machine{}
		machine.Name = fmt.Sprintf("master-%d", i-1)
		spec := &v1beta1.VSphereMachineProviderSpec{
			Network: v1beta1.NetworkSpec{
				Devices: []v1beta1.NetworkDeviceSpec{
					{
						// Create default first device
					},
				},
			},
		}

		var gateway string
		if ipTypes&ipv4 != 0 {
			gateway = ipv4Gateway
			spec.Network.Devices[0].Nameservers = append(spec.Network.Devices[0].Nameservers, "8.8.8.8")
		} else if ipTypes&ipv6 != 0 {
			gateway = ipv6Gateway
			spec.Network.Devices[0].Nameservers = append(spec.Network.Devices[0].Nameservers, "2001::100")
		}

		idx := 0
		if ipTypes&ipv4 != 0 {
			ipAddress := ipamv1.IPAddress{}
			ipAddress.Name = fmt.Sprintf("%s-claim-%d-%d", machine.Name, 0, idx)
			pool := v1beta1.AddressesFromPool{
				Group:    "installer.openshift.io",
				Resource: "IPPool",
				Name:     fmt.Sprintf("default-%d", idx),
			}
			spec.Network.Devices[0].AddressesFromPools = append(spec.Network.Devices[0].AddressesFromPools, pool)
			ip := fmt.Sprintf(ipv4Template, i)
			separatorIndex := strings.Index(ip, "/")
			if separatorIndex > 0 {
				ipAddress.Spec.Address = ip[:separatorIndex]
				prefix, err := strconv.Atoi(ip[strings.Index(ip, "/")+1:])
				if err == nil {
					ipAddress.Spec.Prefix = prefix
				}
			}
			ipAddress.Spec.Gateway = gateway
			ips = append(ips, ipAddress)
			idx++
		}
		if ipTypes&ipv6 != 0 {
			ipAddress := ipamv1.IPAddress{}
			ipAddress.Name = fmt.Sprintf("%s-claim-%d-%d", machine.Name, 0, idx)
			pool := v1beta1.AddressesFromPool{
				Group:    "installer.openshift.io",
				Resource: "IPPool",
				Name:     fmt.Sprintf("default-%d", idx),
			}
			spec.Network.Devices[0].AddressesFromPools = append(spec.Network.Devices[0].AddressesFromPools, pool)
			ip := fmt.Sprintf(ipv6Template, i)
			separatorIndex := strings.Index(ip, "/")
			if separatorIndex > 0 {
				ipAddress.Spec.Address = ip[:separatorIndex]
				prefix, err := strconv.Atoi(ip[strings.Index(ip, "/")+1:])
				if err == nil {
					ipAddress.Spec.Prefix = prefix
				}
			}
			ipAddress.Spec.Gateway = gateway
			ips = append(ips, ipAddress)
		}
		specs = append(specs, spec)
		machines = append(machines, machine)
	}

	return machines, specs, ips
}

func createVsphereConfig() *config {
	return &config{}
}

func TestProcessGuestNetworkConfiguration(t *testing.T) {
	cases := []struct {
		name                 string
		config               *config
		source               TFVarsSources
		expectedBootKargs    string
		expectedControlKargs []string
		expectedError        string
	}{
		{
			name:                 "No Hosts",
			config:               createVsphereConfig(),
			source:               createTFVarsSources(false, 0),
			expectedBootKargs:    "",
			expectedControlKargs: []string(nil),
		},
		{
			name:              "Hosts - Single IPV4",
			config:            createVsphereConfig(),
			source:            createTFVarsSources(true, ipv4),
			expectedBootKargs: "ip=192.168.101.240::192.168.101.200:255.255.255.0:::none nameserver=8.8.8.8",
			expectedControlKargs: []string{
				"ip=192.168.101.241::192.168.101.200:255.255.255.0:::none nameserver=8.8.8.8",
				"ip=192.168.101.242::192.168.101.200:255.255.255.0:::none nameserver=8.8.8.8",
				"ip=192.168.101.243::192.168.101.200:255.255.255.0:::none nameserver=8.8.8.8",
			},
		},
		{
			name:              "Hosts - Single IPV6",
			config:            createVsphereConfig(),
			source:            createTFVarsSources(true, ipv6),
			expectedBootKargs: "ip=[2001::240]::[2001::200]:64:::none nameserver=[2001::100]",
			expectedControlKargs: []string{
				"ip=[2001::241]::[2001::200]:64:::none nameserver=[2001::100]",
				"ip=[2001::242]::[2001::200]:64:::none nameserver=[2001::100]",
				"ip=[2001::243]::[2001::200]:64:::none nameserver=[2001::100]",
			},
		},
		{
			name:              "Hosts - Dual Stack",
			config:            createVsphereConfig(),
			source:            createTFVarsSources(true, ipv4|ipv6),
			expectedBootKargs: "ip=192.168.101.240::192.168.101.200:255.255.255.0:::none ip=[2001::240]:::64:::none nameserver=8.8.8.8",
			expectedControlKargs: []string{
				"ip=192.168.101.241::192.168.101.200:255.255.255.0:::none ip=[2001::241]:::64:::none nameserver=8.8.8.8",
				"ip=192.168.101.242::192.168.101.200:255.255.255.0:::none ip=[2001::242]:::64:::none nameserver=8.8.8.8",
				"ip=192.168.101.243::192.168.101.200:255.255.255.0:::none ip=[2001::243]:::64:::none nameserver=8.8.8.8",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := processGuestNetworkConfiguration(tc.config, tc.source)
			if tc.expectedError == "" {
				assert.NoError(t, err)

				// Verify values
				assert.Equal(t, tc.expectedBootKargs, tc.config.BootStrapNetworkKargs)
				assert.Equal(t, tc.expectedControlKargs, tc.config.ControlPlaneNetworkKargs)
			} else {
				assert.Regexp(t, regexp.MustCompile(tc.expectedError), err)
			}
		})
	}
}
