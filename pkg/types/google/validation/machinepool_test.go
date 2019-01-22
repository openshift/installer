package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/google"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name  string
		pool  *google.MachinePool
		valid bool
	}{
		{
			name:  "empty",
			pool:  &google.MachinePool{},
			valid: true,
		},
		{
			name: "valid size",
			pool: &google.MachinePool{
				RootVolume: google.RootVolume{
					Size: 10,
				},
			},
			valid: true,
		},
		{
			name: "invalid size",
			pool: &google.MachinePool{
				RootVolume: google.RootVolume{
					Size: -10,
				},
			},
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
