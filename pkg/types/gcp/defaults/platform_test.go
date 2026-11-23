package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/network"
)

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name     string
		platform *gcp.Platform
		expected *gcp.Platform
	}{
		{
			name:     "nil platform is a no-op",
			platform: nil,
			expected: nil,
		},
		{
			name: "empty ipFamily defaults to IPv4",
			platform: &gcp.Platform{
				Region:    "us-east1",
				ProjectID: "valid-project",
			},
			expected: &gcp.Platform{
				Region:    "us-east1",
				ProjectID: "valid-project",
				IPFamily:  network.IPv4,
			},
		},
		{
			name: "ipFamily is not overridden when set to IPv4",
			platform: &gcp.Platform{
				Region:   "us-east1",
				IPFamily: network.IPv4,
			},
			expected: &gcp.Platform{
				Region:   "us-east1",
				IPFamily: network.IPv4,
			},
		},
		{
			name: "ipFamily is not overridden when set to DualStackIPv4Primary",
			platform: &gcp.Platform{
				Region:   "us-east1",
				IPFamily: network.DualStackIPv4Primary,
			},
			expected: &gcp.Platform{
				Region:   "us-east1",
				IPFamily: network.DualStackIPv4Primary,
			},
		},
		{
			name: "ipFamily is not overridden when set to DualStackIPv6Primary",
			platform: &gcp.Platform{
				Region:   "us-east1",
				IPFamily: network.DualStackIPv6Primary,
			},
			expected: &gcp.Platform{
				Region:   "us-east1",
				IPFamily: network.DualStackIPv6Primary,
			},
		},
		{
			name: "kms key defaults are still applied alongside ipFamily",
			platform: &gcp.Platform{
				Region:    "us-east1",
				ProjectID: "valid-project",
				DefaultMachinePlatform: &gcp.MachinePool{
					OSDisk: gcp.OSDisk{
						EncryptionKey: &gcp.EncryptionKeyReference{
							KMSKey: &gcp.KMSKeyReference{Name: "key"},
						},
					},
				},
			},
			expected: &gcp.Platform{
				Region:    "us-east1",
				ProjectID: "valid-project",
				IPFamily:  network.IPv4,
				DefaultMachinePlatform: &gcp.MachinePool{
					OSDisk: gcp.OSDisk{
						EncryptionKey: &gcp.EncryptionKeyReference{
							KMSKey: &gcp.KMSKeyReference{
								Name:      "key",
								ProjectID: "valid-project",
								Location:  "us-east1",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetPlatformDefaults(tc.platform)
			assert.Equal(t, tc.expected, tc.platform)
		})
	}
}
