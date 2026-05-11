package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/ipnet"
)

func TestBuildNoProxySet(t *testing.T) {
	cases := []struct {
		name        string
		networking  *Networking
		userNoProxy string
		expected    []string
	}{
		{
			name:        "empty networking and no user entries",
			networking:  &Networking{},
			userNoProxy: "",
			expected:    []string{"127.0.0.1", "localhost", ".svc", ".cluster.local"},
		},
		{
			name: "service network entries are included",
			networking: &Networking{
				ServiceNetwork: []ipnet.IPNet{
					*ipnet.MustParseCIDR("172.30.0.0/16"),
				},
			},
			userNoProxy: "",
			expected:    []string{"127.0.0.1", "localhost", ".svc", ".cluster.local", "172.30.0.0/16"},
		},
		{
			name: "machine network entries are included",
			networking: &Networking{
				MachineNetwork: []MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
				},
			},
			userNoProxy: "",
			expected:    []string{"127.0.0.1", "localhost", ".svc", ".cluster.local", "10.0.0.0/16"},
		},
		{
			name: "cluster network entries are included",
			networking: &Networking{
				ClusterNetwork: []ClusterNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("10.128.0.0/14")},
				},
			},
			userNoProxy: "",
			expected:    []string{"127.0.0.1", "localhost", ".svc", ".cluster.local", "10.128.0.0/14"},
		},
		{
			name:        "single user no-proxy entry is included",
			networking:  &Networking{},
			userNoProxy: "example.com",
			expected:    []string{"127.0.0.1", "localhost", ".svc", ".cluster.local", "example.com"},
		},
		{
			name:        "multiple comma-separated user entries are included",
			networking:  &Networking{},
			userNoProxy: "example.com,internal.corp,192.168.1.0/24",
			expected:    []string{"127.0.0.1", "localhost", ".svc", ".cluster.local", "example.com", "internal.corp", "192.168.1.0/24"},
		},
		{
			name:        "user entries with surrounding whitespace are trimmed",
			networking:  &Networking{},
			userNoProxy: " example.com , internal.corp ",
			expected:    []string{"127.0.0.1", "localhost", ".svc", ".cluster.local", "example.com", "internal.corp"},
		},
		{
			name:        "empty segments in user no-proxy are ignored",
			networking:  &Networking{},
			userNoProxy: "example.com,,internal.corp,",
			expected:    []string{"127.0.0.1", "localhost", ".svc", ".cluster.local", "example.com", "internal.corp"},
		},
		{
			name: "all network types and user entries combined",
			networking: &Networking{
				ServiceNetwork: []ipnet.IPNet{
					*ipnet.MustParseCIDR("172.30.0.0/16"),
				},
				MachineNetwork: []MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
				},
				ClusterNetwork: []ClusterNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("10.128.0.0/14")},
				},
			},
			userNoProxy: "example.com",
			expected: []string{
				"127.0.0.1", "localhost", ".svc", ".cluster.local",
				"172.30.0.0/16", "10.0.0.0/16", "10.128.0.0/14", "example.com",
			},
		},
		{
			name: "duplicate user entry matching a built-in entry is deduplicated",
			networking: &Networking{
				ServiceNetwork: []ipnet.IPNet{
					*ipnet.MustParseCIDR("172.30.0.0/16"),
				},
			},
			userNoProxy: "172.30.0.0/16,localhost",
			expected:    []string{"127.0.0.1", "localhost", ".svc", ".cluster.local", "172.30.0.0/16"},
		},
		{
			name: "multiple entries in each network type",
			networking: &Networking{
				ServiceNetwork: []ipnet.IPNet{
					*ipnet.MustParseCIDR("172.30.0.0/16"),
					*ipnet.MustParseCIDR("fd02::/112"),
				},
				MachineNetwork: []MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
					{CIDR: *ipnet.MustParseCIDR("fd00::/48")},
				},
				ClusterNetwork: []ClusterNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("10.128.0.0/14")},
					{CIDR: *ipnet.MustParseCIDR("fd01::/48")},
				},
			},
			userNoProxy: "",
			expected: []string{
				"127.0.0.1", "localhost", ".svc", ".cluster.local",
				"172.30.0.0/16", "fd02::/112", "10.0.0.0/16", "fd00::/48", "10.128.0.0/14", "fd01::/48",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := BuildNoProxySet(tc.networking, tc.userNoProxy)
			assert.ElementsMatch(t, tc.expected, result.List())
		})
	}
}
