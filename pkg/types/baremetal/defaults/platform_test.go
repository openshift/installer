package defaults

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

const testClusterName = "test-cluster"

func TestSetPlatformDefaults(t *testing.T) {
	// Stub the call to net.LookupHost
	lookupHost = func (host string) (addrs []string, err error) {
		if host == "api.test-cluster.test" {
			ips := []string{"192.168.111.2",}
			return  ips, nil
		} else if host == "test.apps.test-cluster.test" {
			ips := []string{"192.168.111.3",}
			return  ips, nil
		} else {
			return nil, errors.New("Unknown Host " + host)
		}
	}
	cases := []struct {
		name     string
		platform *baremetal.Platform
		expected *baremetal.Platform
	}{
		{
			name:     "default_empty",
			platform: &baremetal.Platform{},
			expected: &baremetal.Platform{
				LibvirtURI: "qemu:///system",
				ClusterProvisioningIP: "172.22.0.3",
				BootstrapProvisioningIP: "172.22.0.2",
				ExternalBridge: "baremetal",
				ProvisioningBridge: "provisioning",
				APIVIP: "192.168.111.2",
				IngressVIP: "192.168.111.3",
				ProvisioningInterface: "ens3",
				ProvisioningNetworkCIDR: "172.22.0.0/24",
				ProvisioningDHCPStart: "172.22.0.20",
				ProvisioningDHCPEnd: "172.22.0.200",
				CachedImageURL: "http://192.168.111.1/images",
			},
		},
		{
			name:     "alternate_cidr",
			platform: &baremetal.Platform{
				ProvisioningNetworkCIDR: "172.23.0.0/24",
			},
			expected: &baremetal.Platform{
				LibvirtURI: "qemu:///system",
				ClusterProvisioningIP: "172.23.0.3",
				BootstrapProvisioningIP: "172.23.0.2",
				ExternalBridge: "baremetal",
				ProvisioningBridge: "provisioning",
				APIVIP: "192.168.111.2",
				IngressVIP: "192.168.111.3",
				ProvisioningInterface: "ens3",
				ProvisioningNetworkCIDR: "172.23.0.0/24",
				ProvisioningDHCPStart: "172.23.0.20",
				ProvisioningDHCPEnd: "172.23.0.200",
				CachedImageURL: "http://192.168.111.1/images",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ic := &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: testClusterName,
				},
				BaseDomain: "test",
			}
			SetPlatformDefaults(tc.platform, ic)
			assert.Equal(t, tc.expected, tc.platform, "unexpected platform")
		})
	}
}
