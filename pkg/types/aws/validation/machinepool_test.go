package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/aws"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name  string
		pool  *aws.MachinePool
		valid bool
	}{
		{
			name:  "empty",
			pool:  &aws.MachinePool{},
			valid: true,
		},
		{
			name: "valid iops",
			pool: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					IOPS: 10,
				},
			},
			valid: true,
		},
		{
			name: "invalid iops",
			pool: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					IOPS: -10,
				},
			},
			valid: false,
		},
		{
			name: "valid size",
			pool: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					Size: 10,
				},
			},
			valid: true,
		},
		{
			name: "invalid size",
			pool: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
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
