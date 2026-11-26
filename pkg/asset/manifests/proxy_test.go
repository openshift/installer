package manifests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/none"
)

func TestBuildNoProxySet(t *testing.T) {
	cases := []struct {
		name     string
		config   *types.InstallConfig
		expected []string
	}{
		{
			name: "empty networking and no user entries",
			config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				BaseDomain: "example.com",
				Networking: &types.Networking{},
				Proxy:      &types.Proxy{},
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
				Proxy: &types.Proxy{},
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
				Proxy: &types.Proxy{},
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
				Proxy: &types.Proxy{},
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
				Proxy: &types.Proxy{},
			},
			expected: []string{
				".cluster.local", ".svc", "10.0.0.0/16", "10.128.0.0/14",
				"127.0.0.1", "172.30.0.0/16", "api-int.test.example.com",
				"fd00::/48", "fd01::/48", "fd02::/112", "localhost",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := BuildNoProxySet(tc.config)
			assert.ElementsMatch(t, tc.expected, sets.List(result))
		})
	}
}

func TestCreateNoProxyIBMCloudIncludesIMDS(t *testing.T) {
	cases := []struct {
		name             string
		platform         types.Platform
		wantIMDS         bool
		wantIBMHost      bool
		wantAWSEC2Suffix string
	}{
		{
			name:        "ibmcloud includes VPC IMDS address and hostname",
			platform:    types.Platform{IBMCloud: &ibmcloud.Platform{Region: "us-south"}},
			wantIMDS:    true,
			wantIBMHost: true,
		},
		{
			name:             "aws still includes IMDS",
			platform:         types.Platform{AWS: &aws.Platform{Region: "us-east-1"}},
			wantIMDS:         true,
			wantAWSEC2Suffix: ".ec2.internal",
		},
		{
			name:     "none does not include IMDS",
			platform: types.Platform{None: &none.Platform{}},
			wantIMDS: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			noProxy, err := createNoProxy(proxyTestInstallConfig(tc.platform))
			assert.NoError(t, err)
			entries := strings.Split(noProxy, ",")
			if tc.wantIMDS {
				assert.Contains(t, entries, "169.254.169.254")
			} else {
				assert.NotContains(t, entries, "169.254.169.254")
			}
			if tc.wantIBMHost {
				assert.Contains(t, entries, "api.metadata.cloud.ibm.com")
			} else {
				assert.NotContains(t, entries, "api.metadata.cloud.ibm.com")
			}
			if tc.wantAWSEC2Suffix != "" {
				assert.Contains(t, entries, tc.wantAWSEC2Suffix)
			}
			assert.Contains(t, entries, "127.0.0.1")
			assert.Contains(t, entries, "api-int.test.example.com")
		})
	}
}

func proxyTestInstallConfig(platform types.Platform) *installconfig.InstallConfig {
	return installconfig.MakeAsset(&types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		BaseDomain: "example.com",
		Networking: &types.Networking{
			ClusterNetwork: []types.ClusterNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.128.0.0/14"), HostPrefix: 23},
			},
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
			},
			ServiceNetwork: []ipnet.IPNet{
				*ipnet.MustParseCIDR("172.30.0.0/16"),
			},
		},
		Platform: platform,
		Proxy: &types.Proxy{
			HTTPProxy:  "http://proxy.example.com:3128",
			HTTPSProxy: "http://proxy.example.com:3128",
		},
	})
}
