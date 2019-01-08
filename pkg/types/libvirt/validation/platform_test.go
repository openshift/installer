package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/libvirt"
)

func validPlatform() *libvirt.Platform {
	return &libvirt.Platform{
		URI: "qemu+tcp://192.168.122.1/system",
		Network: libvirt.Network{
			IfName: "tt0",
		},
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *libvirt.Platform
		valid    bool
	}{
		{
			name:     "minimal",
			platform: validPlatform(),
			valid:    true,
		},
		{
			name: "invalid uri",
			platform: func() *libvirt.Platform {
				p := validPlatform()
				p.URI = "bad-uri"
				return p
			}(),
			valid: false,
		},
		{
			name: "missing interface name",
			platform: func() *libvirt.Platform {
				p := validPlatform()
				p.Network.IfName = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "valid machine pool",
			platform: func() *libvirt.Platform {
				p := validPlatform()
				p.DefaultMachinePlatform = &libvirt.MachinePool{}
				return p
			}(),
			valid: true,
		},
		{
			name: "invalid machine pool",
			platform: func() *libvirt.Platform {
				p := validPlatform()
				p.DefaultMachinePlatform = &libvirt.MachinePool{
					Image: "bad-image",
				}
				return p
			}(),
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
