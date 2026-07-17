package manifests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
)

func TestBuildNoProxySet(t *testing.T) {
	cases := []struct {
		name             string
		config           *types.InstallConfig
		expected         []string
		expectedWildcard bool
	}{
		{
			name: "empty networking and no user entries",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{},
			},
			expected: []string{".cluster.local", ".svc", "127.0.0.1", "api-int.test.example.com", "localhost"},
		},
		{
			name: "service network entries are included",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{
					ServiceNetwork: []ipnet.IPNet{
						*ipnet.MustParseCIDR("172.30.0.0/16"),
					},
				},
			},
			expected: []string{".cluster.local", ".svc", "127.0.0.1", "172.30.0.0/16", "api-int.test.example.com", "localhost"},
		},
		{
			name: "machine network entries are included",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
					},
				},
			},
			expected: []string{".cluster.local", ".svc", "10.0.0.0/16", "127.0.0.1", "api-int.test.example.com", "localhost"},
		},
		{
			name: "cluster network entries are included",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{
					ClusterNetwork: []types.ClusterNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("10.128.0.0/14")},
					},
				},
			},
			expected: []string{".cluster.local", ".svc", "10.128.0.0/14", "127.0.0.1", "api-int.test.example.com", "localhost"},
		},
		{
			name: "single user no-proxy entry is included",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{},
				Proxy: &types.Proxy{
					NoProxy: "example.com",
				},
			},
			expected: []string{".cluster.local", ".svc", "127.0.0.1", "api-int.test.example.com", "example.com", "localhost"},
		},
		{
			name: "multiple comma-separated user entries are included",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{},
				Proxy: &types.Proxy{
					NoProxy: "example.com,internal.corp,192.168.1.0/24",
				},
			},
			expected: []string{".cluster.local", ".svc", "127.0.0.1", "192.168.1.0/24", "api-int.test.example.com", "example.com", "internal.corp", "localhost"},
		},
		{
			name: "user entries with surrounding whitespace are trimmed",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{},
				Proxy: &types.Proxy{
					NoProxy: " example.com , internal.corp ",
				},
			},
			expected: []string{".cluster.local", ".svc", "127.0.0.1", "api-int.test.example.com", "example.com", "internal.corp", "localhost"},
		},
		{
			name: "empty segments in user no-proxy are ignored",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{},
				Proxy: &types.Proxy{
					NoProxy: "example.com,,internal.corp,",
				},
			},
			expected: []string{".cluster.local", ".svc", "127.0.0.1", "api-int.test.example.com", "example.com", "internal.corp", "localhost"},
		},
		{
			name: "all network types and user entries combined",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{
					ServiceNetwork: []ipnet.IPNet{
						*ipnet.MustParseCIDR("172.30.0.0/16"),
					},
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
					},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("10.128.0.0/14")},
					},
				},
				Proxy: &types.Proxy{
					NoProxy: "example.com",
				},
			},
			expected: []string{
				".cluster.local", ".svc", "10.0.0.0/16", "10.128.0.0/14",
				"127.0.0.1", "172.30.0.0/16", "api-int.test.example.com", "example.com", "localhost",
			},
		},
		{
			name: "duplicate user entry matching a built-in entry is deduplicated",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{
					ServiceNetwork: []ipnet.IPNet{
						*ipnet.MustParseCIDR("172.30.0.0/16"),
					},
				},
				Proxy: &types.Proxy{
					NoProxy: "172.30.0.0/16,localhost",
				},
			},
			expected: []string{".cluster.local", ".svc", "127.0.0.1", "172.30.0.0/16", "api-int.test.example.com", "localhost"},
		},
		{
			name: "multiple entries in each network type",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{
					ServiceNetwork: []ipnet.IPNet{
						*ipnet.MustParseCIDR("172.30.0.0/16"),
						*ipnet.MustParseCIDR("fd02::/112"),
					},
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
						{CIDR: *ipnet.MustParseCIDR("fd00::/48")},
					},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("10.128.0.0/14")},
						{CIDR: *ipnet.MustParseCIDR("fd01::/48")},
					},
				},
			},
			expected: []string{
				".cluster.local", ".svc", "10.0.0.0/16", "10.128.0.0/14",
				"127.0.0.1", "172.30.0.0/16", "api-int.test.example.com",
				"fd00::/48", "fd01::/48", "fd02::/112", "localhost",
			},
		},
		{
			name: "wildcard no-proxy returns nil set",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{
					ServiceNetwork: []ipnet.IPNet{
						*ipnet.MustParseCIDR("172.30.0.0/16"),
					},
				},
				Proxy: &types.Proxy{
					NoProxy: "*",
				},
			},
			expectedWildcard: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, wildcard := BuildNoProxySet(tc.config)
			assert.Equal(t, tc.expectedWildcard, wildcard)
			if wildcard {
				assert.Nil(t, result)
			} else {
				assert.ElementsMatch(t, tc.expected, sets.List(result))
			}
		})
	}
}
