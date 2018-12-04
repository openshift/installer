package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/libvirt"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name  string
		pool  *libvirt.MachinePool
		valid bool
	}{
		{
			name:  "empty",
			pool:  &libvirt.MachinePool{},
			valid: true,
		},
		{
			name: "valid image",
			pool: &libvirt.MachinePool{
				Image: "https://example.com/rhcos-qemu.qcow2",
			},
			valid: true,
		},
		{
			name: "invalid image",
			pool: &libvirt.MachinePool{
				Image: "bad-image",
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
