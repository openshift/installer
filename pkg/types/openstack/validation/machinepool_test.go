package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/openstack"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name     string
		pool     *openstack.MachinePool
		expected string
	}{
		{
			name: "empty",
			pool: &openstack.MachinePool{},
		},
		{
			name: "invalid size",
			pool: &openstack.MachinePool{
				RootVolume: &openstack.RootVolume{
					Size: -10,
					Type: "default",
				},
			},
			expected: `^test-path\.rootVolume\.size: Invalid value: -10: Volume size must be greater than zero to use root volumes$`,
		},
		{
			name: "missing size",
			pool: &openstack.MachinePool{
				RootVolume: &openstack.RootVolume{
					Type: "default",
				},
			},
			expected: `^test-path\.rootVolume\.size: Invalid value: 0: Volume size must be greater than zero to use root volumes$`,
		},
		{
			name: "missing type",
			pool: &openstack.MachinePool{
				RootVolume: &openstack.RootVolume{
					Size: 10,
				},
			},
			expected: `^test-path\.rootVolume\.type: Invalid value: "": Volume type must be specified to use root volumes$`,
		},
		{
			name: "valid root volume",
			pool: &openstack.MachinePool{
				RootVolume: &openstack.RootVolume{
					Size: 10,
					Type: "default",
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}
