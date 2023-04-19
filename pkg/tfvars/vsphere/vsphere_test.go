package vsphere

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

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

func createTFVarsSources(createHosts bool, ipTypes int64, cpc []*v1beta1.VSphereMachineProviderSpec) TFVarsSources {
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
		ControlPlaneConfigs: cpc,
	}

	if createHosts {
		tvs.InstallConfig.Config.VSphere.Hosts = createValidHosts(ipTypes)
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
		hosts[0].NetworkDevice.Gateway4 = ipv4Gateway
		hosts[0].NetworkDevice.Nameservers = append(hosts[0].NetworkDevice.Nameservers, "8.8.8.8")
		hosts[0].NetworkDevice.IPAddrs = append(hosts[0].NetworkDevice.IPAddrs, fmt.Sprintf(ipv4Template, 0))
	} else if ipTypes&ipv6 != 0 {
		hosts[0].NetworkDevice.Gateway6 = ipv6Gateway
		hosts[0].NetworkDevice.Nameservers = append(hosts[0].NetworkDevice.Nameservers, "2001::100")
	}

	if ipTypes&ipv6 != 0 {
		hosts[0].NetworkDevice.IPAddrs = append(hosts[0].NetworkDevice.IPAddrs, fmt.Sprintf(ipv6Template, 0))
	}

	return hosts
}

func createControlPlaneConfigs(ipTypes int64) []*v1beta1.VSphereMachineProviderSpec {
	var machines []*v1beta1.VSphereMachineProviderSpec

	for i := 1; i <= 3; i++ {
		machine := &v1beta1.VSphereMachineProviderSpec{
			Network: v1beta1.NetworkSpec{
				Devices: []v1beta1.NetworkDeviceSpec{
					{
						// Create default first device
					},
				},
			},
		}

		if ipTypes&ipv4 != 0 {
			machine.Network.Devices[0].Gateway4 = ipv4Gateway
			machine.Network.Devices[0].Nameservers = append(machine.Network.Devices[0].Nameservers, "8.8.8.8")
		} else if ipTypes&ipv6 != 0 {
			machine.Network.Devices[0].Gateway6 = ipv6Gateway
			machine.Network.Devices[0].Nameservers = append(machine.Network.Devices[0].Nameservers, "2001::100")
		}

		if ipTypes&ipv4 != 0 {
			machine.Network.Devices[0].IPAddrs = append(machine.Network.Devices[0].IPAddrs, fmt.Sprintf(ipv4Template, i))
		}
		if ipTypes&ipv6 != 0 {
			machine.Network.Devices[0].IPAddrs = append(machine.Network.Devices[0].IPAddrs, fmt.Sprintf(ipv6Template, i))
		}
		machines = append(machines, machine)
	}

	return machines
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
			source:               createTFVarsSources(false, 0, nil),
			expectedBootKargs:    "",
			expectedControlKargs: []string(nil),
		},
		{
			name:              "Hosts - Single IPV4",
			config:            createVsphereConfig(),
			source:            createTFVarsSources(true, ipv4, createControlPlaneConfigs(ipv4)),
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
			source:            createTFVarsSources(true, ipv6, createControlPlaneConfigs(ipv6)),
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
			source:            createTFVarsSources(true, ipv4|ipv6, createControlPlaneConfigs(ipv4|ipv6)),
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
