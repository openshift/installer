package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/kubevirt"
)

func validPlatform() *kubevirt.Platform {
	return &kubevirt.Platform{
		Namespace:                  "test-namespace",
		StorageClass:               "",
		NetworkName:                "test network",
		APIVIP:                     "192.168.123.15",
		IngressVIP:                 "192.168.123.16",
		PersistentVolumeAccessMode: "ReadWriteMany",
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *kubevirt.Platform
		valid    bool
	}{
		{
			name:     "valid",
			platform: validPlatform(),
			valid:    true,
		},
		{
			name: "empty namespace",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.Namespace = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "empty network name",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.NetworkName = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "empty API VIP",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.APIVIP = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid API VIP",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.APIVIP = "invalid API VIP"
				return p
			}(),
			valid: false,
		},
		{
			name: "API VIP not in CIDR",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.APIVIP = "10.0.0.1"
				return p
			}(),
			valid: false,
		},
		{
			name: "empty ingress VIP",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.IngressVIP = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid ingress VIP",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.IngressVIP = "invalid ingress VIP"
				return p
			}(),
			valid: false,
		},
		{
			name: "ingress VIP not in CIDR",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.IngressVIP = "1.1.1.1"
				return p
			}(),
			valid: false,
		},
		{
			name: "valid - empty access mode",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.PersistentVolumeAccessMode = ""
				return p
			}(),
			valid: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ic := &types.InstallConfig{
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("192.168.123.0/24")},
					},
				},
			}
			err := ValidatePlatform(tc.platform, field.NewPath("test-path"), ic).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
