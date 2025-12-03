package validation

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/baremetal/defaults"
)

/*
This test uses a Fluent Builder and Object Mother patterns to increase the readability of the test code,
so that only the relevant values could be easily exposed for each case, thus allowing the reader to
immediately catch the important pieces.

Every builder exposes one or more factory methods to create canned objects that could be further customized
using the fluent interface by chaining the exposed functions accordingly.
*/

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		config   *types.InstallConfig
		platform *baremetal.Platform
		expected string
	}{
		{
			name:     "valid",
			platform: platform().build(),
		},
		{
			name: "valid_ipv6_provisioning",
			platform: platform().
				ProvisioningNetworkCIDR("fd2e:6f44:5dd8:b856::/64").
				ClusterProvisioningIP("fd2e:6f44:5dd8:b856::3").
				BootstrapProvisioningIP("fd2e:6f44:5dd8:b856::2").build(),
		},
		{
			name: "invalid_hosts",
			platform: platform().
				Hosts().build(),
			expected: "bare metal hosts are missing",
		},
		{
			name: "no_hosts_machineapi_disabled",
			config: installConfig().
				Capabilities(types.Capabilities{BaselineCapabilitySet: configv1.ClusterVersionCapabilitySetNone,
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{configv1.ClusterVersionCapabilityIngress,
						configv1.ClusterVersionCapabilityBaremetal}}).
				BareMetalPlatform(
					platform().Hosts()).build(),
		},
		{
			name: "toofew_masters_norole",
			config: installConfig().
				BareMetalPlatform(
					platform().Hosts(
						host1().Role("worker"),
						host2())).
				ControlPlane(
					machinePool().Replicas(3)).
				Compute(
					machinePool().Replicas(1)).build(),
			expected: "baremetal.hosts: Required value: not enough hosts found \\(1\\) to support all the configured ControlPlane replicas \\(3\\)",
		},
		{
			name: "toofew_masters",
			config: installConfig().
				BareMetalPlatform(
					platform().Hosts(
						host1().Role("master"),
						host2().Role("worker"))).
				ControlPlane(
					machinePool().Replicas(3)).
				Compute(
					machinePool().Replicas(1)).build(),
			expected: "baremetal.hosts: Required value: not enough hosts found \\(1\\) to support all the configured ControlPlane replicas \\(3\\)",
		},
		{
			name: "toofew_workers",
			config: installConfig().
				BareMetalPlatform(
					platform().Hosts(
						host1().Role("master"),
						host2().Role("worker"))).
				ControlPlane(
					machinePool().Replicas(1)).
				Compute(
					machinePool().Replicas(3)).build(),
			expected: "baremetal.hosts: Required value: not enough hosts found \\(1\\) to support all the configured Compute replicas \\(3\\)",
		},
		{
			name: "enough_hosts",
			config: installConfig().
				BareMetalPlatform(
					platform().Hosts(
						host1().Role("master"),
						host2().Role("worker"))).
				ControlPlane(
					machinePool().Replicas(1)).
				Compute(
					machinePool().Replicas(1)).build(),
		},
		{
			name: "enough_hosts_norole",
			config: installConfig().
				BareMetalPlatform(
					platform().Hosts(
						host1(),
						host2(),
						host3())).
				ControlPlane(
					machinePool().Replicas(1)).
				Compute(
					machinePool().Replicas(2)).build(),
		},
		{
			name: "enough_hosts_mixed",
			config: installConfig().
				BareMetalPlatform(
					platform().Hosts(
						host1().Role("master"),
						host2().Role("worker"),
						host3(),
						host4())).
				ControlPlane(
					machinePool().Replicas(2)).
				Compute(
					machinePool().Replicas(2)).build(),
		},
		{
			name: "not_enough_hosts_norole",
			config: installConfig().
				BareMetalPlatform(
					platform().Hosts(
						host1(),
						host2(),
						host3())).
				ControlPlane(
					machinePool().Replicas(2)).
				Compute(
					machinePool().Replicas(2)).build(),
			expected: "baremetal.hosts: Required value: not enough hosts found \\(1\\) to support all the configured Compute replicas \\(2\\)",
		},
		{
			name: "more_than_enough_hosts",
			config: installConfig().
				BareMetalPlatform(
					platform().Hosts(
						host1().Role("master"),
						host2().Role("master"),
						host3().Role("worker"),
						host4().Role("worker"),
						host5())).
				ControlPlane(
					machinePool().Replicas(1)).
				Compute(
					machinePool().Replicas(1)).build(),
		},
		{
			name: "norole_for_workers",
			config: installConfig().
				BareMetalPlatform(
					platform().Hosts(
						host1().Role("master"),
						host2().Role("master"),
						host3().Role("master"),
						host4(),
						host5())).
				ControlPlane(
					machinePool().Replicas(3)).
				Compute(
					machinePool().Replicas(2)).build(),
		},
		{
			name: "missing_name",
			config: installConfig().
				BareMetalPlatform(
					platform().Hosts(
						host1().Name(""))).
				ControlPlane(machinePool().Replicas(1)).build(),
			expected: "baremetal.hosts\\[0\\].name: Required value: missing Name",
		},
		{
			name: "allowed_feature_loadbalancer_openshift_managed_default",
			config: installConfig().
				BareMetalPlatform(
					platform().LoadBalancerType("OpenShiftManagedDefault")).
				FeatureSet(configv1.TechPreviewNoUpgrade).build(),
		},
		{
			name: "allowed_feature_loadbalancer_user_managed",
			config: installConfig().
				BareMetalPlatform(
					platform().LoadBalancerType("UserManaged")).
				FeatureSet(configv1.TechPreviewNoUpgrade).build(),
		},
		{
			name: "allowed_feature_loadbalancer_invalid",
			config: installConfig().
				BareMetalPlatform(
					platform().LoadBalancerType("FooBar")).
				FeatureSet(configv1.TechPreviewNoUpgrade).build(),
			expected: "baremetal.loadBalancer.type: Invalid value: \"FooBar\": invalid load balancer type",
		},
		{
			name: "missing_mac",
			platform: platform().
				Hosts(host1().BootMACAddress("")).build(),
			expected: "baremetal.hosts\\[0\\].bootMACAddress: Required value: missing BootMACAddress",
		},
		{
			name: "duplicate_host_name",
			platform: platform().
				Hosts(
					host1().Name("host1"),
					host2().Name("host1")).build(),
			expected: "baremetal.hosts\\[1\\].name: Duplicate value: \"host1\"",
		},
		{
			name: "valid_host_name",
			platform: platform().
				Hosts(host1().Name("host1")).build(),
			expected: "",
		},
		{
			name: "valid_host_name_fqdn",
			platform: platform().
				Hosts(host1().Name("test.example.com")).build(),
			expected: "",
		},
		{
			name: "invalid_host_name_char",
			platform: platform().
				Hosts(host1().Name("test,example.com")).build(),
			expected: "baremetal.hosts\\[0\\].name: Invalid value: \"test,example.com\"",
		},
		{
			name: "invalid_host_name_uppercase",
			platform: platform().
				Hosts(host1().Name("Host1")).build(),
			expected: "baremetal.hosts\\[0\\].name: Invalid value: \"Host1\"",
		},
		{
			name: "invalid_host_name_length",
			platform: platform().
				Hosts(host1().Name(strings.Repeat("a", 300))).build(),
			expected: "baremetal.hosts\\[0\\].name: Invalid value: \"aaaaaaaaa",
		},
		{
			name: "duplicate_host_mac",
			platform: platform().
				Hosts(
					host1().BootMACAddress("CA:FE:CA:FE:CA:FE"),
					host2().BootMACAddress("CA:FE:CA:FE:CA:FE")).build(),
			expected: "baremetal.hosts\\[1\\].bootMACAddress: Duplicate value: \"CA:FE:CA:FE:CA:FE\"",
		},
		{
			name: "invalid_boot_mode",
			platform: platform().
				Hosts(host1().BootMode("not-a-valid-value")).build(),
			expected: "baremetal.hosts\\[0\\].bootMode: Unsupported value: \"not-a-valid-value\": supported values: \"UEFI\", \"UEFISecureBoot\", \"legacy\"",
		},
		{
			name: "uefi_boot_mode",
			platform: platform().
				Hosts(host1().BootMode("UEFI")).build(),
			expected: "",
		},
		{
			name: "uefi_secure_boot_mode",
			platform: platform().
				Hosts(host1().BMCAddress("redfish://example.com/redfish/v1").BootMode("UEFISecureBoot")).build(),
			expected: "",
		},
		{
			name: "unsupported_uefi_secure_boot_mode",
			platform: platform().
				Hosts(host1().BootMode("UEFISecureBoot")).build(),
			expected: "baremetal.hosts\\[0\\].bootMode: Invalid value: \"UEFISecureBoot\": driver ipmi does not support UEFI secure boot",
		},
		{
			name: "legacy_boot_mode",
			platform: platform().
				Hosts(host1().BootMode("legacy")).build(),
			expected: "",
		},
		{
			name:     "provisioningNetwork_disabled_valid",
			platform: platform().ProvisioningNetwork(baremetal.DisabledProvisioningNetwork).build(),
		},
		{
			name:     "provisioningNetwork_unmanaged_valid",
			platform: platform().ProvisioningNetwork(baremetal.UnmanagedProvisioningNetwork).build(),
		},
		{
			name:     "provisioningNetwork_invalid",
			platform: platform().ProvisioningNetwork("Invalid").build(),
			expected: `Unsupported value: "Invalid": supported values: "Disabled", "Managed", "Unmanaged"`,
		},
		{
			name:     "networkConfig_invalid",
			platform: platform().Hosts(host1().NetworkConfig("Not a valid yaml content")).build(),
			expected: ".*Invalid value.*Not a valid yaml: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type map\\[string\\]interface \\{\\}",
		},
		{
			name: "networkConfig_valid_yml",
			platform: platform().Hosts(host1().NetworkConfig(`
interfaces:
- name: eth1
  type: ethernet
  state: up
- name: linux-br0
  type: linux-bridge
  state: up
  bridge:
    options:
      group-forward-mask: 0
      mac-ageing-time: 300
      multicast-snooping: true
      stp:
        enabled: true
        forward-delay: 15
        hello-time: 2
        max-age: 20
        priority: 32768
      port:
        - name: eth1
          stp-hairpin-mode: false
          stp-path-cost: 100
          stp-priority: 32`)).build(),
			expected: "",
		},
		{
			name: "networkConfig_invalid_mtu_string",
			platform: platform().Hosts(
				host1().NetworkConfig(`
interfaces:
- name: eth0
  type: ethernet
  state: up
  mtu: "1500"
`)).build(),
			expected: `baremetal.hosts\[0\].networkConfig.interfaces\[0\].mtu: Invalid value: "1500": mtu must be an integer \(not quoted string\)`,
		},
		{
			name: "networkConfig_valid_mtu_integer",
			platform: platform().Hosts(
				host1().NetworkConfig(`
interfaces:
- name: eth0
  type: ethernet
  state: up
  mtu: 1500
`)).build(),
			expected: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Build default wrapping installConfig
			if tc.config == nil {
				tc.config = installConfig().build()
				tc.config.BareMetal = tc.platform
			}

			err := ValidatePlatform(tc.config.BareMetal, false, network(), field.NewPath("baremetal"), tc.config).ToAggregate()

			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}

func TestValidateHostRootDeviceHints(t *testing.T) {
	cases := []struct {
		name            string
		rootDeviceHints *baremetal.RootDeviceHints
		expectedSuccess bool
	}{
		{
			name:            "nil hints",
			expectedSuccess: true,
		},
		{
			name:            "no hints",
			rootDeviceHints: &baremetal.RootDeviceHints{},
			expectedSuccess: true,
		},
		{
			name: "non /dev path",
			rootDeviceHints: &baremetal.RootDeviceHints{
				DeviceName: "sda",
			},
		},
		{
			name: "/dev path",
			rootDeviceHints: &baremetal.RootDeviceHints{
				DeviceName: "/dev/sda",
			},
			expectedSuccess: true,
		},
		{
			name: "by-path path",
			rootDeviceHints: &baremetal.RootDeviceHints{
				DeviceName: "/dev/disk/by-path/pci-0000:01:00.0-scsi-0:2:0:0",
			},
			expectedSuccess: true,
		},
		{
			name: "by-id path",
			rootDeviceHints: &baremetal.RootDeviceHints{
				DeviceName: "/dev/disk/by-id/wwn-0x600508e000000000ce506dc50ab0ad05",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateHostRootDeviceHints(tc.rootDeviceHints, field.NewPath("rootDeviceHints"))

			if tc.expectedSuccess {
				assert.Empty(t, errs)
			} else {
				assert.NotEmpty(t, errs)
			}
		})
	}
}

func TestValidateProvisioning(t *testing.T) {
	//Used for url validations
	imagesServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		if strings.Contains(r.RequestURI, "notexistent") {
			http.NotFound(w, r)
		}
	}))
	defer imagesServer.Close()

	origInterfaceValidator := interfaceValidator
	t.Cleanup(func() {
		interfaceValidator = origInterfaceValidator
	})
	interfaceValidator = func(libvirtURI string) (func(string) error, error) {
		return func(interfaceName string) error {
			interfaceNames := []string{"br0", "br1"}
			for _, foundInterface := range interfaceNames {
				if foundInterface == interfaceName {
					return nil
				}
			}

			return fmt.Errorf("could not find interface %q, valid interfaces are %s", interfaceName, strings.Join(interfaceNames, ", "))
		}, nil
	}

	cases := []struct {
		name     string
		config   *types.InstallConfig
		platform *baremetal.Platform
		expected string
	}{
		{
			name: "duplicate_bmc_address",
			platform: platform().
				Hosts(
					host1().BMCAddress("ipmi://192.168.111.1"),
					host2().BMCAddress("ipmi://192.168.111.1")).build(),
			expected: "baremetal.hosts\\[1\\].bmc.address: Duplicate value: \"ipmi://192.168.111.1\"",
		},
		{
			name: "bmc_address_required",
			platform: platform().
				Hosts(host1().BMCAddress("")).build(),
			expected: "baremetal.hosts\\[0\\].bmc.address: Required value: missing Address",
		},
		{
			name: "bmc_username_required",
			platform: platform().
				Hosts(host1().BMCUsername("")).build(),
			expected: "baremetal.hosts\\[0\\].bmc.username: Required value: missing Username",
		},
		{
			name: "bmc_password_required",
			platform: platform().
				Hosts(host1().BMCPassword("")).build(),
			expected: "baremetal.hosts\\[0\\].bmc.password: Required value: missing Password",
		},
		{
			name: "valid_with_os_image_overrides",
			platform: platform().
				BootstrapOSImage(imagesServer.URL + "/images/qemu.x86_64.qcow2.gz?sha256=3b5a882c2af3e19d515b961855d144f293cab30190c2bdedd661af31a1fc4e2f").
				ClusterOSImage(imagesServer.URL + "/images/metal.x86_64.qcow2.gz?sha256=340dfa4d92450f2eee852ed1e2d02e3138cc68d824827ef9cf0a40a7ea2f93da").build(),
		},
		{
			name: "invalid_bootstraposimage",
			platform: platform().
				BootstrapOSImage("192.168.111.1/images/qemu.x86_64.qcow2.gz?sha256=3b5a882c2af3e19d515b961855d144f293cab30190c2bdedd661af31a1fc4e2f").build(),
			expected: "baremetal.bootstrapOSImage: Invalid value:.*: the URI provided:.*is invalid",
		},
		{
			name: "invalid_clusterosimage",
			platform: platform().
				ClusterOSImage("http//192.168.111.1/images/metal.x86_64.qcow2.gz?sha256=340dfa4d92450f2eee852ed1e2d02e3138cc68d824827ef9cf0a40a7ea2f93da").build(),
			expected: "baremetal.clusterOSImage: Invalid value:.*: the URI provided:.*is invalid",
		},
		{
			name: "invalid_bootstraposimage_checksum",
			platform: platform().
				BootstrapOSImage("http://192.168.111.1/images/qemu.x86_64.qcow2.gz?md5sum=3b5a882c2af3e19d515b961855d144f293cab30190c2bdedd661af31a1fc4e2f").build(),
			expected: "baremetal.bootstrapOSImage: Invalid value:.*: the sha256 parameter in the.*URI is missing",
		},
		{
			name: "invalid_clusterosimage_checksum",
			platform: platform().
				ClusterOSImage("http://192.168.111.1/images/metal.x86_64.qcow2.gz?sha256=3ee852ed1e2d02e3138cc68d824827ef9cf0a40a7ea2f93da").build(),
			expected: "baremetal.clusterOSImage: Invalid value:.*: the sha256 parameter in the.*URI is invalid",
		},
		{
			name: "invalid_bootstraposimage_uri_scheme",
			platform: platform().
				BootstrapOSImage("xttp://192.168.111.1/images/qemu.x86_64.qcow2.gz?sha256=3b5a882c2af3e19d515b961855d144f293cab30190c2bdedd661af31a1fc4e2f").build(),
			expected: "baremetal.bootstrapOSImage: Invalid value:.*: the URI provided.*must begin with http/https",
		},
		{
			name: "invalid_clusterosimage_uri_scheme",
			platform: platform().
				ClusterOSImage("xttp://192.168.111.1/images/qemu.x86_64.qcow2.gz?sha256=3b5a882c2af3e19d515b961855d144f293cab30190c2bdedd661af31a1fc4e2f").build(),
			expected: "baremetal.clusterOSImage: Invalid value:.*: the URI provided.*must begin with http/https",
		},
		{
			name: "notfound_bootstraposimage",
			platform: platform().
				BootstrapOSImage(imagesServer.URL + "/images/notexistent.x86_64.qcow2.gz?sha256=3b5a882c2af3e19d515b961855d144f293cab30190c2bdedd661af31a1fc4e2f").build(),
			expected: "baremetal.bootstrapOSImage: Not found:.*",
		},
		{
			name: "notfound_clusterosimageimage",
			platform: platform().
				ClusterOSImage(imagesServer.URL + "/images/notexistent.x86_64.qcow2.gz?sha256=3b5a882c2af3e19d515b961855d144f293cab30190c2bdedd661af31a1fc4e2f").build(),
			expected: "baremetal.clusterOSImage: Not found:.*",
		},
		{
			name: "invalid_extbridge",
			platform: platform().
				ExternalBridge("noexist").build(),
			expected: "Invalid value: \"noexist\": could not find interface \"noexist\", valid interfaces are br0, br1",
		},
		{
			name:     "valid_extbridge_mac",
			platform: platform().ExternalMACAddress("CA:FE:CA:FE:CA:FE").build(),
		},
		{
			name: "invalid_provbridge",
			platform: platform().
				ProvisioningBridge("noexist").build(),
			expected: "Invalid value: \"noexist\": could not find interface \"noexist\", valid interfaces are br0, br1",
		},
		{
			name:     "valid_provbridge_mac",
			platform: platform().ProvisioningMACAddress("CA:FE:CA:FE:CA:FE").build(),
		},
		{
			name: "invalid_duplicate_bridge_macs",
			platform: platform().
				ProvisioningMACAddress("CA:FE:CA:FE:CA:FE").
				ExternalMACAddress("CA:FE:CA:FE:CA:FE").
				build(),
			expected: "Duplicate value: \"provisioning and external MAC addresses may not be identical\"",
		},
		{
			name: "valid_both_macs_specified",
			platform: platform().
				ProvisioningMACAddress("CA:FE:CA:FE:CA:FD").
				ExternalMACAddress("CA:FE:CA:FE:CA:FE").
				build(),
		},
		{
			name: "invalid_multicast_mac",
			platform: platform().
				ExternalMACAddress("7D:CE:E3:29:35:6F").
				build(),
			expected: "expected unicast mac address, found multicast",
		},
		{
			name: "invalid_bootstrapprovip_wrongCIDR",
			platform: platform().
				BootstrapProvisioningIP("192.168.128.1").build(),
			expected: "Invalid value: \"192.168.128.1\": \"192.168.128.1\" is not in the provisioning network",
		},
		{
			name: "invalid_bootstrapprovip_machineCIDR",
			platform: platform().
				BootstrapProvisioningIP("192.168.111.5").build(),
			expected: "Invalid value: \"192.168.111.5\": \"192.168.111.5\" is not in the provisioning network",
		},

		{
			name: "invalid_clusterprovip_machineCIDR",
			platform: platform().
				ClusterProvisioningIP("192.168.111.5").build(),
			expected: "Invalid value: \"192.168.111.5\": \"192.168.111.5\" is not in the provisioning network",
		},
		{
			name: "invalid_clusterprovip_wrongCIDR",
			platform: platform().
				ClusterProvisioningIP("192.168.128.1").build(),
			expected: "Invalid value: \"192.168.128.1\": \"192.168.128.1\" is not in the provisioning network",
		},
		{
			name: "invalid_provisioning_network_overlapping_CIDR",
			platform: platform().
				ProvisioningNetworkCIDR("192.168.110.0/23").build(),
			expected: "Invalid value: \"192.168.110.0/23\": cannot overlap with machine network: 192.168.111.0/24 overlaps with 192.168.110.0/23",
		},
		{
			name: "valid_provisioningDHCPRange",
			platform: platform().
				ProvisioningDHCPRange("172.22.0.10,172.22.0.50").build(),
		},
		{
			name: "invalid_provisioningDHCPRange_missing_pair",
			platform: platform().
				ProvisioningDHCPRange("172.22.0.10,").build(),
			expected: "provisioningDHCPRange: Invalid value: \"172.22.0.10,\": : \"\" is not a valid IP",
		},
		{
			name: "invalid_provisioningDHCPRange_not_a_range",
			platform: platform().
				ProvisioningDHCPRange("172.22.0.19").build(),
			expected: "Invalid value: \"172.22.0.19\": provisioning DHCP range should be in format: start_ip,end_ip",
		},
		{
			name: "invalid_provisioningDHCPRange_wrong_CIDR",
			platform: platform().
				ProvisioningDHCPRange("192.168.128.1,172.22.0.100").build(),
			expected: "Invalid value: \"192.168.128.1,172.22.0.100\": \"192.168.128.1\" is not in the provisioning network",
		},
		{
			name: "invalid_clusterprovip_overlapDHCPRange",
			platform: platform().
				ClusterProvisioningIP("172.22.0.10").build(),
			expected: "Invalid value: \"172.22.0.10\": \"172.22.0.10\" overlaps with the allocated DHCP range",
		},
		{
			name: "invalid_bootstrapprovip_overlapDHCPRange",
			platform: platform().
				BootstrapProvisioningIP("172.22.0.20").build(),
			expected: "Invalid value: \"172.22.0.20\": \"172.22.0.20\" overlaps with the allocated DHCP range",
		},
		{
			name: "invalid_libvirturi",
			platform: platform().
				LibvirtURI("bad").build(),
			expected: "invalid URI \"bad\"",
		},
		{
			name:     "valid_provisioning_network_ipv4",
			platform: platform().ProvisioningNetworkCIDR("172.22.0.0/24").build(),
			expected: "",
		},
		{
			name:     "invalid_provisioning_network_need_network_address",
			platform: platform().ProvisioningNetworkCIDR("172.22.0.2/24").build(),
			expected: "provisioningNetworkCIDR has host bits set, expected 172.22.0.0/24",
		},
		{
			name: "valid_provisioning_network_ipv6",
			platform: platform().
				ProvisioningNetworkCIDR("fd00:0111:0::/64").
				ClusterProvisioningIP("fd00:0111::3").
				BootstrapProvisioningIP("fd00:0111::2").build(),
			expected: "",
		},
		{
			name: "valid_provisioning_network_ipv6_long",
			platform: platform().
				ProvisioningNetworkCIDR("fd00:0111:0000:0000:0000:0000:0000:0000/64").
				ClusterProvisioningIP("fd00:0111:0000:0000:0000:0000:0000:0003").
				BootstrapProvisioningIP("fd00:0111:0000:0000:0000:0000:0000:0002").build(),
			expected: "",
		},
		{
			name: "valid_provisioning_network_ipv6_mixed",
			platform: platform().
				ProvisioningNetworkCIDR("fd00:0111::/64").
				ClusterProvisioningIP("fd00:0111:0000:0000:0000:0000:0000:0003").
				BootstrapProvisioningIP("fd00:0111:0000:0000:0000:0000:0000:0002").build(),
			expected: "",
		},
		{
			name: "invalid_provisioning_network_need_network_address_ipv6",
			platform: platform().
				ProvisioningNetworkCIDR("fd00:0111:0::1/64").
				ClusterProvisioningIP("fd00:0111::3").
				BootstrapProvisioningIP("fd00:0111::2").build(),
			expected: "provisioningNetworkCIDR has host bits set, expected fd00:111::/64",
		},
		{
			name: "ipv6_CIDR_too_large",
			platform: platform().
				ProvisioningNetworkCIDR("fd2e:6f44:5dd8:b856::/32").
				ClusterProvisioningIP("fd2e:6f44:5dd8:b856::3").
				BootstrapProvisioningIP("fd2e:6f44:5dd8:b856::2").build(),
			expected: "provisioningNetworkCIDR mask must be greater than or equal to 64 for IPv6 networks",
		},

		// Disabled provisioning network
		{
			name:   "valid_provisioningDisabled_noProvisioningInterface",
			config: installConfig().Network(networking().Network("192.168.111.0/24")).build(),
			platform: platform().
				ProvisioningNetwork(baremetal.DisabledProvisioningNetwork).
				ClusterProvisioningIP("192.168.111.2").
				BootstrapProvisioningIP("192.168.111.3").
				ProvisioningNetworkInterface("").build(),
		},
		{
			name:   "valid_provisioningDisabled_IPs_in_machineCIDR",
			config: installConfig().Network(networking().Network("192.168.111.0/24")).build(),
			platform: platform().
				ProvisioningNetwork(baremetal.DisabledProvisioningNetwork).
				ClusterProvisioningIP("192.168.111.2").
				BootstrapProvisioningIP("192.168.111.3").build(),
		},
		{
			name:   "valid_provisioningDisabled_no_provisioning_ips",
			config: installConfig().Network(networking().Network("192.168.111.0/24")).build(),
			platform: platform().
				ProvisioningNetwork(baremetal.DisabledProvisioningNetwork).
				ClusterProvisioningIP("").
				BootstrapProvisioningIP("").build(),
		},
		{
			name:   "invalid_provisioningDisabled_IPs_not_in_machineCIDR",
			config: installConfig().Network(networking().Network("192.168.111.0/24")).build(),
			platform: platform().
				ProvisioningNetwork(baremetal.DisabledProvisioningNetwork).
				BootstrapProvisioningIP("192.168.111.3").
				ClusterProvisioningIP("192.168.0.2").build(),
			expected: "Invalid value: \"192.168.0.2\": provisioning network is disabled, IP expected to be in one of the machine networks: 192.168.111.0/24",
		},
		{
			name:   "not_supported_bmc_driver_provisioning_network_disabled",
			config: installConfig().Network(networking().Network("192.168.111.0/24")).build(),
			platform: platform().
				ProvisioningNetwork(baremetal.DisabledProvisioningNetwork).
				ClusterProvisioningIP("192.168.111.2").
				BootstrapProvisioningIP("192.168.111.3").
				Hosts(host1().BMCAddress("ipmi://192.168.111.1")).
				build(),
			expected: "baremetal.hosts\\[0\\].bmc: Invalid value: \"ipmi://192.168.111.1\": driver ipmi requires provisioning network",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Build default wrapping installConfig
			if tc.config == nil {
				tc.config = installConfig().build()
			}
			tc.config.BareMetal = tc.platform

			defaults.SetPlatformDefaults(tc.config.BareMetal, tc.config)

			err := ValidateProvisioning(tc.config.BareMetal, network(), field.NewPath("baremetal")).ToAggregate()

			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}

type hostBuilder struct {
	baremetal.Host
}

func host1() *hostBuilder {
	return &hostBuilder{
		baremetal.Host{
			Name:           "host1",
			BootMACAddress: "CA:FE:CA:FE:00:00",
			BMC: baremetal.BMC{
				Username: "root",
				Password: "password",
				Address:  "ipmi://192.168.111.1",
			},
		},
	}
}

func host2() *hostBuilder {
	return &hostBuilder{
		baremetal.Host{
			Name:           "host2",
			BootMACAddress: "CA:FE:CA:FE:00:01",
			BMC: baremetal.BMC{
				Username: "root",
				Password: "password",
				Address:  "ipmi://192.168.111.2",
			},
		},
	}
}

func host3() *hostBuilder {
	return &hostBuilder{
		baremetal.Host{
			Name:           "host3",
			BootMACAddress: "CA:FE:CA:FE:00:02",
			BMC: baremetal.BMC{
				Username: "root",
				Password: "password",
				Address:  "ipmi://192.168.111.3",
			},
		},
	}
}

func host4() *hostBuilder {
	return &hostBuilder{
		baremetal.Host{
			Name:           "host4",
			BootMACAddress: "CA:FE:CA:FE:00:03",
			BMC: baremetal.BMC{
				Username: "root",
				Password: "password",
				Address:  "ipmi://192.168.111.4",
			},
		},
	}
}

func host5() *hostBuilder {
	return &hostBuilder{
		baremetal.Host{
			Name:           "host5",
			BootMACAddress: "CA:FE:CA:FE:00:04",
			BMC: baremetal.BMC{
				Username: "root",
				Password: "password",
				Address:  "ipmi://192.168.111.5",
			},
		},
	}
}

func (hb *hostBuilder) build() *baremetal.Host {
	return &hb.Host
}

func (hb *hostBuilder) Name(value string) *hostBuilder {
	hb.Host.Name = value
	return hb
}

func (hb *hostBuilder) BootMACAddress(value string) *hostBuilder {
	hb.Host.BootMACAddress = value
	return hb
}

func (hb *hostBuilder) BootMode(value string) *hostBuilder {
	hb.Host.BootMode = baremetal.BootMode(value)
	return hb
}

func (hb *hostBuilder) BMCAddress(value string) *hostBuilder {
	hb.Host.BMC.Address = value
	return hb
}

func (hb *hostBuilder) BMCUsername(value string) *hostBuilder {
	hb.Host.BMC.Username = value
	return hb
}

func (hb *hostBuilder) BMCPassword(value string) *hostBuilder {
	hb.Host.BMC.Password = value
	return hb
}

func (hb *hostBuilder) Role(value string) *hostBuilder {
	hb.Host.Role = value
	return hb
}

func (hb *hostBuilder) NetworkConfig(value string) *hostBuilder {
	yaml.Unmarshal([]byte(value), &hb.Host.NetworkConfig)
	return hb
}

type platformBuilder struct {
	baremetal.Platform
}

func platform() *platformBuilder {
	return &platformBuilder{
		baremetal.Platform{
			APIVIPs:                      []string{"192.168.111.2"},
			IngressVIPs:                  []string{"192.168.111.4"},
			Hosts:                        []*baremetal.Host{},
			LibvirtURI:                   "qemu://system",
			ProvisioningNetworkCIDR:      ipnet.MustParseCIDR("172.22.0.0/24"),
			ProvisioningNetwork:          baremetal.ManagedProvisioningNetwork,
			ClusterProvisioningIP:        "172.22.0.3",
			BootstrapProvisioningIP:      "172.22.0.2",
			ExternalBridge:               "br0",
			ProvisioningBridge:           "br1",
			ProvisioningNetworkInterface: "ens3",
		}}
}

func (pb *platformBuilder) build() *baremetal.Platform {
	return &pb.Platform
}

func (pb *platformBuilder) ProvisioningNetworkCIDR(value string) *platformBuilder {
	pb.Platform.ProvisioningNetworkCIDR = ipnet.MustParseCIDR(value)
	return pb
}

func (pb *platformBuilder) ProvisioningNetwork(value baremetal.ProvisioningNetwork) *platformBuilder {
	pb.Platform.ProvisioningNetwork = value
	return pb
}

func (pb *platformBuilder) ClusterProvisioningIP(value string) *platformBuilder {
	pb.Platform.ClusterProvisioningIP = value
	return pb
}

func (pb *platformBuilder) BootstrapProvisioningIP(value string) *platformBuilder {
	pb.Platform.BootstrapProvisioningIP = value
	return pb
}

func (pb *platformBuilder) BootstrapOSImage(value string) *platformBuilder {
	pb.Platform.DeprecatedBootstrapOSImage = value
	return pb
}

func (pb *platformBuilder) ClusterOSImage(value string) *platformBuilder {
	pb.Platform.DeprecatedClusterOSImage = value
	return pb
}

func (pb *platformBuilder) ProvisioningDHCPRange(value string) *platformBuilder {
	pb.Platform.ProvisioningDHCPRange = value
	return pb
}

func (pb *platformBuilder) Hosts(builders ...*hostBuilder) *platformBuilder {
	pb.Platform.Hosts = nil
	for _, builder := range builders {
		pb.Platform.Hosts = append(pb.Platform.Hosts, builder.build())
	}
	return pb
}

func (pb *platformBuilder) LoadBalancerType(value string) *platformBuilder {
	pb.Platform.LoadBalancer = &configv1.BareMetalPlatformLoadBalancer{
		Type: configv1.PlatformLoadBalancerType(value),
	}
	return pb
}

func (pb *platformBuilder) LibvirtURI(value string) *platformBuilder {
	pb.Platform.LibvirtURI = value
	return pb
}

func (pb *platformBuilder) ExternalBridge(value string) *platformBuilder {
	pb.Platform.ExternalBridge = value
	return pb
}

func (pb *platformBuilder) ExternalMACAddress(value string) *platformBuilder {
	pb.Platform.ExternalMACAddress = value
	return pb
}

func (pb *platformBuilder) ProvisioningBridge(value string) *platformBuilder {
	pb.Platform.ProvisioningBridge = value
	return pb
}

func (pb *platformBuilder) ProvisioningMACAddress(value string) *platformBuilder {
	pb.Platform.ProvisioningMACAddress = value
	return pb
}

func (pb *platformBuilder) ProvisioningNetworkInterface(value string) *platformBuilder {
	pb.Platform.ProvisioningNetworkInterface = value
	return pb
}

func network() *types.Networking {
	return &types.Networking{MachineNetwork: []types.MachineNetworkEntry{{CIDR: *ipnet.MustParseCIDR("192.168.111.0/24")}}}
}

type installConfigBuilder struct {
	types.InstallConfig
}

func installConfig() *installConfigBuilder {
	return &installConfigBuilder{
		InstallConfig: types.InstallConfig{},
	}
}

func (icb *installConfigBuilder) build() *types.InstallConfig {
	return &icb.InstallConfig
}

func (icb *installConfigBuilder) BareMetalPlatform(builder *platformBuilder) *installConfigBuilder {
	icb.InstallConfig.Platform = types.Platform{
		BareMetal: builder.build(),
	}
	return icb
}

func (icb *installConfigBuilder) Capabilities(capabilities types.Capabilities) *installConfigBuilder {
	icb.InstallConfig.Capabilities = &capabilities
	return icb
}

func (icb *installConfigBuilder) FeatureSet(value configv1.FeatureSet) *installConfigBuilder {
	icb.InstallConfig.FeatureSet = value
	return icb
}

func (icb *installConfigBuilder) ControlPlane(builder *machinePoolBuilder) *installConfigBuilder {
	icb.InstallConfig.ControlPlane = builder.build()

	return icb
}

func (icb *installConfigBuilder) Network(builder *networkingBuilder) *installConfigBuilder {
	icb.InstallConfig.Networking = builder.build()

	return icb
}

func (icb *installConfigBuilder) Compute(builders ...*machinePoolBuilder) *installConfigBuilder {
	icb.InstallConfig.Compute = nil
	for _, builder := range builders {
		icb.InstallConfig.Compute = append(icb.InstallConfig.Compute, *builder.build())
	}
	return icb
}

type machinePoolBuilder struct {
	types.MachinePool
}

func machinePool() *machinePoolBuilder {
	return &machinePoolBuilder{
		MachinePool: types.MachinePool{},
	}
}

func (mpb *machinePoolBuilder) build() *types.MachinePool {
	return &mpb.MachinePool
}

func (mpb *machinePoolBuilder) Replicas(count int64) *machinePoolBuilder {
	mpb.MachinePool.Replicas = &count
	return mpb
}

type networkingBuilder struct {
	types.Networking
}

func networking() *networkingBuilder {
	return &networkingBuilder{
		Networking: types.Networking{},
	}
}

func (nb *networkingBuilder) Network(cidr string) *networkingBuilder {
	network := ipnet.MustParseCIDR(cidr)

	nb.MachineNetwork = append(nb.MachineNetwork, types.MachineNetworkEntry{
		CIDR: *network,
	})

	return nb
}

func (nb *networkingBuilder) build() *types.Networking {
	return &nb.Networking
}
