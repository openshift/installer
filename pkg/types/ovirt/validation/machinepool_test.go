package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/ovirt"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name  string
		pool  *ovirt.MachinePool
		valid bool
	}{
		{
			name:  "empty",
			pool:  &ovirt.MachinePool{},
			valid: true,
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
