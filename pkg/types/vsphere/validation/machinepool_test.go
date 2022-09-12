package validation

import (
	"testing"

	"github.com/openshift/installer/pkg/types"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/vsphere"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name           string
		pool           *types.MachinePool
		platform       *vsphere.Platform
		expectedErrMsg string
		expectedZones  *[]string
	}{
		{
			name: "empty",
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{},
				},
			},
			platform:       validPlatform(),
			expectedErrMsg: "",
		}, {
			name:     "negative disk size",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						OSDisk: vsphere.OSDisk{
							DiskSizeGB: -1,
						},
					},
				},
			},
			expectedErrMsg: `^test-path\.diskSizeGB: Invalid value: -1: storage disk size must be positive$`,
		}, {
			name:     "negative CPUs",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						NumCPUs: -1,
					},
				},
			},
			expectedErrMsg: `^test-path\.cpus: Invalid value: -1: number of CPUs must be positive$`,
		}, {
			name:     "negative cores",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						NumCoresPerSocket: -1,
					},
				},
			},
			expectedErrMsg: `^test-path\.coresPerSocket: Invalid value: -1: cores per socket must be positive$`,
		}, {
			name:     "negative memory",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						MemoryMiB: -1,
					},
				},
			},
			expectedErrMsg: `^test-path\.memoryMB: Invalid value: -1: memory size must be positive$`,
		}, {
			name:     "less CPUs than cores per socket",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						NumCPUs:           1,
						NumCoresPerSocket: 8,
					},
				},
			},
			expectedErrMsg: `^test-path\.coresPerSocket: Invalid value: 8: cores per socket must be less than number of CPUs$`,
		},
		{
			name:     "numCPUs not a multiple of cores per socket",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						NumCPUs:           7,
						NumCoresPerSocket: 4,
					},
				},
			},
			expectedErrMsg: `^test-path.cpus: Invalid value: 7: numCPUs specified should be a multiple of cores per socket$`,
		},
		{
			name:     "numCPUs not a multiple of default cores per socket",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						NumCPUs: 7,
					},
				},
			},
			expectedErrMsg: `^test-path.cpus: Invalid value: 7: numCPUs specified should be a multiple of cores per socket which is by default 4$`,
		},
		{
			name: "multi-zone invalid zone name",
			platform: func() *vsphere.Platform {
				platform := validMultiVCenterPlatform()
				platform.FailureDomains[0].Name = "Zone%^@112233"
				return platform
			}(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						Zones: []string{
							"Zone%^@112233",
						},
					},
				},
			},
			expectedErrMsg: `^test-path.zones: Invalid value: \[\]string{"Zone%\^@112233"}: cluster name must begin with a lower-case letter$`,
		},
		{
			name:     "multi-zone valid",
			platform: validMultiVCenterPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						Zones: []string{
							"test-east-1a",
						},
					},
				},
			},
		},
		{
			name:     "multi-zone no zones defined for control plane pool",
			platform: validMultiVCenterPlatform(),
			pool: &types.MachinePool{
				Name: types.MachinePoolControlPlaneRoleName,
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{},
				},
			},
			expectedZones:  &[]string{"test-east-1a", "test-east-2a"},
			expectedErrMsg: "",
		},
		{
			name:     "multi-zone no zones defined for compute pool",
			platform: validMultiVCenterPlatform(),
			pool: &types.MachinePool{
				Name: types.MachinePoolComputeRoleName,
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{},
				},
			},
			expectedZones:  &[]string{"test-east-1a", "test-east-2a"},
			expectedErrMsg: "",
		},
		{
			name:     "multi-zone undefined zone",
			platform: validMultiVCenterPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						Zones: []string{
							"unknown-zone",
						},
					},
				},
			},
			expectedErrMsg: `^test-path.zones: Invalid value: "unknown-zone": zone not defined in failureDomains$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.platform, tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.expectedErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedErrMsg, err)
			}
			if tc.expectedZones != nil {
				zones := tc.pool.Platform.VSphere.Zones
				for _, expectedZone := range *tc.expectedZones {
					found := false
					for _, zone := range zones {
						if zone == expectedZone {
							found = true
							break
						}
					}
					if found == false {
						t.Errorf("expected zone not found %s", expectedZone)
					}
				}
				for _, zone := range zones {
					found := false
					for _, expectedZone := range *tc.expectedZones {
						if zone == expectedZone {
							found = true
							break
						}
					}
					if found == false {
						t.Errorf("unexpected zone %s", zone)
					}

				}
			}
		})
	}
}
