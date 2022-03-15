package validation

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/powervs"
)

var (
	validRegion = powervs.RegionShortNames()[rand.Intn(len(powervs.RegionShortNames()))]
)

func validMinimalPlatform() *powervs.Platform {
	return &powervs.Platform{
		Region: validRegion,
	}
}

func validMachinePool() *powervs.MachinePool {
	return &powervs.MachinePool{}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *powervs.Platform
		valid    bool
	}{
		{
			name:     "Config: Minimal platform config",
			platform: validMinimalPlatform(),
			valid:    true,
		},
		{
			name: "Region: Invalid region",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.Region = "invalid"
				return p
			}(),
			valid: false,
		},
		{
			name: "Region: Missing region",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.Region = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "Machine Pool: Valid machine pool",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.DefaultMachinePlatform = validMachinePool()
				return p
			}(),
			valid: true,
		},
		{
			name: "ServiceInstanceID: Valid Power VS ServiceInstanceID",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceInstanceID = "05d5dbfd-2a62-4d01-b37b-71211be442f6"
				return p
			}(),
			valid: true,
		},
		{
			name: "ServiceInstanceID: Invalid Power VS ServiceInstanceID",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceInstanceID = "abc123"
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
