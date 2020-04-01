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
		{
			name: "valid additional network ids",
			pool: &openstack.MachinePool{
				AdditionalNetworkIDs: []string{
					"51e5fe10-5325-4a32-bce8-7ebe9708c453",
					"3ade1375-acfd-4eda-90be-3530af4f25ec",
					"460e993b-e932-43c6-a7a2-e51ca58f4eae",
				},
			},
		},
		{
			name: "invalid additional network ids",
			pool: &openstack.MachinePool{
				AdditionalNetworkIDs: []string{
					"51e5fe10-5325-4a32-bce8-7ebe9708c453",
					"INVALID",
					"",
				},
			},
			expected: `^\[test-path.additionalNetworkIDs\[1\]: Invalid value: \"INVALID\": valid UUID v4 must be specified, test-path.additionalNetworkIDs\[2\]: Invalid value: \"\": valid UUID v4 must be specified\]$`,
		},
		{
			name: "wrong additional network ids version",
			pool: &openstack.MachinePool{
				AdditionalNetworkIDs: []string{
					"25b91ff0-75c3-11ea-9aff-4fc68ed06d45", // VERSION_1
					"bd15ec47-a3ec-329b-812b-9c617ca86881", // VERSION_3
					"39499c61-11eb-5f02-a519-f2e38575cedd", // VERSION_5
				},
			},
			expected: `^\[test-path.additionalNetworkIDs\[0\]: Invalid value: \"25b91ff0-75c3-11ea-9aff-4fc68ed06d45\": valid UUID v4 must be specified, test-path.additionalNetworkIDs\[1\]: Invalid value: \"bd15ec47-a3ec-329b-812b-9c617ca86881\": valid UUID v4 must be specified, test-path.additionalNetworkIDs\[2\]: Invalid value: \"39499c61-11eb-5f02-a519-f2e38575cedd\": valid UUID v4 must be specified\]$`,
		},
		{
			name: "valid additional security group ids",
			pool: &openstack.MachinePool{
				AdditionalSecurityGroupIDs: []string{
					"51e5fe10-5325-4a32-bce8-7ebe9708c453",
					"3ade1375-acfd-4eda-90be-3530af4f25ec",
					"460e993b-e932-43c6-a7a2-e51ca58f4eae",
				},
			},
		},
		{
			name: "invalid additional security group ids",
			pool: &openstack.MachinePool{
				AdditionalSecurityGroupIDs: []string{
					"51e5fe10-5325-4a32-bce8-7ebe9708c453",
					"INVALID",
					"",
				},
			},
			expected: `^\[test-path.additionalSecurityGroupIDs\[1\]: Invalid value: \"INVALID\": valid UUID v4 must be specified, test-path.additionalSecurityGroupIDs\[2\]: Invalid value: \"\": valid UUID v4 must be specified\]$`,
		},
		{
			name: "wrong additional security group ids version",
			pool: &openstack.MachinePool{
				AdditionalSecurityGroupIDs: []string{
					"25b91ff0-75c3-11ea-9aff-4fc68ed06d45", // VERSION_1
					"bd15ec47-a3ec-329b-812b-9c617ca86881", // VERSION_3
					"39499c61-11eb-5f02-a519-f2e38575cedd", // VERSION_5
				},
			},
			expected: `^\[test-path.additionalSecurityGroupIDs\[0\]: Invalid value: \"25b91ff0-75c3-11ea-9aff-4fc68ed06d45\": valid UUID v4 must be specified, test-path.additionalSecurityGroupIDs\[1\]: Invalid value: \"bd15ec47-a3ec-329b-812b-9c617ca86881\": valid UUID v4 must be specified, test-path.additionalSecurityGroupIDs\[2\]: Invalid value: \"39499c61-11eb-5f02-a519-f2e38575cedd\": valid UUID v4 must be specified\]$`,
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
