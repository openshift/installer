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
	return &vsphere.Platform{
		Workspace: vsphere.Workspace{
			Folder: testClusterName,
		},
		SCSIControllerType: "pvscsi",
	}
}

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name     string
		platform *vsphere.Platform
		expected *vsphere.Platform
	}{
		{
			name:     "empty",
			platform: &vsphere.Platform{},
			expected: defaultPlatform(),
		},
		{
			name: "single vCenter",
			platform: &vsphere.Platform{
				VirtualCenters: []vsphere.VirtualCenter{{Name: "test-server"}},
			},
			expected: func() *vsphere.Platform {
				p := defaultPlatform()
				p.VirtualCenters = []vsphere.VirtualCenter{{Name: "test-server"}}
				p.Workspace.Server = "test-server"
				return p
			}(),
		},
		{
			name: "multiple vCenters",
			platform: &vsphere.Platform{
				VirtualCenters: []vsphere.VirtualCenter{
					{Name: "test-server1"},
					{Name: "test-server2"},
				},
			},
			expected: func() *vsphere.Platform {
				p := defaultPlatform()
				p.VirtualCenters = []vsphere.VirtualCenter{
					{Name: "test-server1"},
					{Name: "test-server2"},
				}
				return p
			}(),
		},
		{
			name: "vCenter set",
			platform: &vsphere.Platform{
				VirtualCenters: []vsphere.VirtualCenter{{Name: "test-server"}},
				Workspace: vsphere.Workspace{
					Server: "other-server",
				},
			},
			expected: func() *vsphere.Platform {
				p := defaultPlatform()
				p.VirtualCenters = []vsphere.VirtualCenter{{Name: "test-server"}}
				p.Workspace.Server = "other-server"
				return p
			}(),
		},
		{
			name: "single datacenter",
			platform: &vsphere.Platform{
				VirtualCenters: []vsphere.VirtualCenter{
					{
						Name:        "test-server",
						Datacenters: []string{"test-datacenter"},
					},
				},
				Workspace: vsphere.Workspace{
					Server: "test-server",
				},
			},
			expected: func() *vsphere.Platform {
				p := defaultPlatform()
				p.VirtualCenters = []vsphere.VirtualCenter{
					{
						Name:        "test-server",
						Datacenters: []string{"test-datacenter"},
					},
				}
				p.Workspace.Server = "test-server"
				p.Workspace.Datacenter = "test-datacenter"
				return p
			}(),
		},
		{
			name: "multiple datacenter",
			platform: &vsphere.Platform{
				VirtualCenters: []vsphere.VirtualCenter{
					{
						Name:        "test-server",
						Datacenters: []string{"test-datacenter1", "test-datacenter2"},
					},
				},
				Workspace: vsphere.Workspace{
					Server: "test-server",
				},
			},
			expected: func() *vsphere.Platform {
				p := defaultPlatform()
				p.VirtualCenters = []vsphere.VirtualCenter{
					{
						Name:        "test-server",
						Datacenters: []string{"test-datacenter1", "test-datacenter2"},
					},
				}
				p.Workspace.Server = "test-server"
				return p
			}(),
		},
		{
			name: "datacenter set",
			platform: &vsphere.Platform{
				VirtualCenters: []vsphere.VirtualCenter{
					{
						Name:        "test-server",
						Datacenters: []string{"test-datacenter"},
					},
				},
				Workspace: vsphere.Workspace{
					Server:     "test-server",
					Datacenter: "other-datacenter",
				},
			},
			expected: func() *vsphere.Platform {
				p := defaultPlatform()
				p.VirtualCenters = []vsphere.VirtualCenter{
					{
						Name:        "test-server",
						Datacenters: []string{"test-datacenter"},
					},
				}
				p.Workspace.Server = "test-server"
				p.Workspace.Datacenter = "other-datacenter"
				return p
			}(),
		},
		{
			name: "folder set",
			platform: &vsphere.Platform{
				Workspace: vsphere.Workspace{
					Folder: "test-folder",
				},
			},
			expected: func() *vsphere.Platform {
				p := defaultPlatform()
				p.Workspace.Folder = "test-folder"
				return p
			}(),
		},
		{
			name: "SCSI controller type set",
			platform: &vsphere.Platform{
				SCSIControllerType: "test-controller-type",
			},
			expected: func() *vsphere.Platform {
				p := defaultPlatform()
				p.SCSIControllerType = "test-controller-type"
				return p
			}(),
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
		})
	}
}
