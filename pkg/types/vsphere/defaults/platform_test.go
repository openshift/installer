package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

const testClusterName = "test-cluster"

func defaultPlatform() *vsphere.Platform {
	return &vsphere.Platform{}
}

func validPlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenters: []vsphere.VCenter{
			{
				Server:   "test-vcenter",
				Port:     443,
				Username: "test-username",
				Password: "test-password",
				Datacenters: []string{
					"test-datacenter",
				},
			},
		},
		FailureDomains: []vsphere.FailureDomain{
			{
				Name:   "test-east-1a",
				Region: "test-east",
				Zone:   "test-east-1a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter",
					ComputeCluster: "/test-datacenter/host/test-cluster",
					Datastore:      "/test-datacenter/datastore/test-datastore",
					Networks:       []string{"test-portgroup"},
					ResourcePool:   "/test-datacenter/host/test-cluster/Resources",
					Folder:         "/test-datacenter/vm/test-folder",
				},
			},
			{
				Name:   "test-east-2a",
				Region: "test-east",
				Zone:   "test-east-2a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter",
					ComputeCluster: "/test-datacenter/host/test-cluster",
					Datastore:      "/test-datacenter/datastore/test-datastore",
					Networks:       []string{"test-portgroup"},
					ResourcePool:   "/test-datacenter/host/test-cluster/Resources",
					Folder:         "/test-datacenter/vm/test-folder",
				},
			},
		},
	}
}

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name       string
		platform   *vsphere.Platform
		expected   *vsphere.Platform
		expectedRP string
	}{
		{
			name:     "empty",
			platform: &vsphere.Platform{},
			expected: defaultPlatform(),
		},
		{
			name:       "completed",
			platform:   validPlatform(),
			expected:   validPlatform(),
			expectedRP: "/test-datacenter/host/test-cluster/Resources",
		},
		{
			name: "incomplete",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				for _, fd := range p.FailureDomains {
					fd.Topology.ResourcePool = ""
				}
				return p
			}(),
			expected:   validPlatform(),
			expectedRP: "/test-datacenter/host/test-cluster/Resources",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ic := &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: testClusterName,
				},
			}
			SetPlatformDefaults(tc.platform, ic)
			assert.Equal(t, tc.expected, tc.platform, "unexpected platform")

			for _, fd := range tc.platform.FailureDomains {
				assert.Equal(t, tc.expectedRP, fd.Topology.ResourcePool, "invalid resource pool")
			}
		})
	}
}
