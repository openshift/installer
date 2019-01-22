package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/google"
)

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *google.Platform
		valid    bool
	}{
		{
			name: "minimal",
			platform: &google.Platform{
				Region: "us-central1",
			},
			valid: true,
		},
		{
			name: "invalid region",
			platform: &google.Platform{
				Region: "bad-region",
			},
			valid: false,
		},
		{
			name: "valid machine pool",
			platform: &google.Platform{
				Region:                 "us-central1",
				DefaultMachinePlatform: &google.MachinePool{},
			},
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
