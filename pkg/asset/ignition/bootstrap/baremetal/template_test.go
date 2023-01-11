package baremetal

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/baremetal"
)

func TestTemplatingIPv4(t *testing.T) {
	bareMetalConfig := baremetal.Platform{
		ProvisioningNetworkCIDR: ipnet.MustParseCIDR("172.22.0.0/24"),
		BootstrapProvisioningIP: "172.22.0.2",
		ProvisioningNetwork:     baremetal.ManagedProvisioningNetwork,
		ProvisioningDHCPRange:   "172.22.0.10,172.22.0.100",
		Hosts: []*baremetal.Host{
			{
				Role:           "master",
				BootMACAddress: "c0:ff:ee:ca:fe:00",
			},
			{
				Role:           "master",
				BootMACAddress: "c0:ff:ee:ca:fe:01",
			},
			{
				Role:           "master",
				BootMACAddress: "c0:ff:ee:ca:fe:02",
			},
			{
				Role:           "worker",
				BootMACAddress: "c0:ff:ee:ca:fe:03",
			},
		},
	}

	result := GetTemplateData(&bareMetalConfig, nil, "bootstrap-ironic-user", "passw0rd")

	assert.Equal(t, result.ProvisioningDHCPRange, "172.22.0.10,172.22.0.100,24")
	assert.Equal(t, result.ProvisioningCIDR, 24)
	assert.Equal(t, result.ProvisioningIPv6, false)
	assert.Equal(t, result.ProvisioningIP, "172.22.0.2")
	assert.Equal(t, result.ProvisioningDHCPAllowList, "c0:ff:ee:ca:fe:00 c0:ff:ee:ca:fe:01 c0:ff:ee:ca:fe:02")
	assert.Equal(t, result.IronicUsername, "bootstrap-ironic-user")
	assert.Equal(t, result.IronicPassword, "passw0rd")
}

func TestTemplatingManagedIPv6(t *testing.T) {
	bareMetalConfig := baremetal.Platform{
		ProvisioningNetworkCIDR: ipnet.MustParseCIDR("fd2e:6f44:5dd8:b856::0/80"),
		ProvisioningDHCPRange:   "fd2e:6f44:5dd8:b856::1,fd2e:6f44:5dd8::ff",
		BootstrapProvisioningIP: "fd2e:6f44:5dd8:b856::2",
		ProvisioningNetwork:     baremetal.ManagedProvisioningNetwork,
	}

	result := GetTemplateData(&bareMetalConfig, nil, "bootstrap-ironic-user", "passw0rd")

	assert.Equal(t, result.ProvisioningDHCPRange, "fd2e:6f44:5dd8:b856::1,fd2e:6f44:5dd8::ff,80")
	assert.Equal(t, result.ProvisioningCIDR, 80)
	assert.Equal(t, result.ProvisioningIPv6, true)
	assert.Equal(t, result.ProvisioningIP, "fd2e:6f44:5dd8:b856::2")
	assert.Equal(t, result.IronicUsername, "bootstrap-ironic-user")
	assert.Equal(t, result.IronicPassword, "passw0rd")
}

func TestTemplatingUnmanagedIPv6(t *testing.T) {
	bareMetalConfig := baremetal.Platform{
		ProvisioningNetworkCIDR: ipnet.MustParseCIDR("fd2e:6f44:5dd8:b856::0/64"),
		BootstrapProvisioningIP: "fd2e:6f44:5dd8:b856::2",
		ProvisioningNetwork:     baremetal.UnmanagedProvisioningNetwork,
	}

	result := GetTemplateData(&bareMetalConfig, nil, "bootstrap-ironic-user", "passw0rd")

	assert.Equal(t, result.ProvisioningDHCPRange, "")
	assert.Equal(t, result.ProvisioningCIDR, 64)
	assert.Equal(t, result.ProvisioningIPv6, true)
	assert.Equal(t, result.ProvisioningIP, "fd2e:6f44:5dd8:b856::2")
	assert.Equal(t, result.ProvisioningDHCPAllowList, "")
	assert.Equal(t, result.IronicUsername, "bootstrap-ironic-user")
	assert.Equal(t, result.IronicPassword, "passw0rd")
}
