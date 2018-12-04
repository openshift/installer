package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/openstack"
)

func validPlatform() *openstack.Platform {
	return &openstack.Platform{
		Region:           "test-region",
		NetworkCIDRBlock: *ipnet.MustParseCIDR("10.0.0.0/16"),
		BaseImage:        "https://example.com/rhcos-qemu.qcow2",
		Cloud:            "test-cloud",
		ExternalNetwork:  "test-ext-net",
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *openstack.Platform
		valid    bool
	}{
		{
			name:     "minimal",
			platform: validPlatform(),
			valid:    true,
		},
		{
			name: "missing region",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.Region = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid base image",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.BaseImage = "bad-image"
				return p
			}(),
			valid: false,
		},
		{
			name: "missing cloud",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.Cloud = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "missing external network",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ExternalNetwork = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "valid default machine pool",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.DefaultMachinePlatform = &openstack.MachinePool{}
				return p
			}(),
			valid: true,
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
