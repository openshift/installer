package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
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
			expectedErrMsg: `^test-path\.coresPerSocket: Invalid value: 8: cores per socket must be less than the number of CPUs \(which is by default \d+\)$`,
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
			expectedErrMsg: `^test-path.cpus: Invalid value: 7: numCPUs specified should be a multiple of cores per socket \(which is by default \d+\)$`,
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
			expectedErrMsg: `^test-path.cpus: Invalid value: 7: numCPUs specified should be a multiple of cores per socket \(which is by default \d+\)$`,
		},
		{
			name: "multi-zone invalid zone name",
			platform: func() *vsphere.Platform {
				platform := validPlatform()
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
			platform: validPlatform(),
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
			platform: validPlatform(),
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
			name:     "multi-zone duplicate zones defined for control plane pool",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Name: types.MachinePoolControlPlaneRoleName,
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						Zones: []string{
							"test-east-1a",
							"test-east-1a",
						},
					},
				},
			},
			expectedErrMsg: `test-path.zones\[1]: Duplicate value: "test-east-1a"`,
		},
		{
			name:     "multi-zone no zones defined for compute pool",
			platform: validPlatform(),
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
			name:     "multi-zone duplicate zones defined for compute pool",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Name: types.MachinePoolComputeRoleName,
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						Zones: []string{
							"test-east-1a",
							"test-east-1a",
						},
					},
				},
			},
			expectedErrMsg: `test-path.zones\[1]: Duplicate value: "test-east-1a"`,
		},
		{
			name:     "multi-zone undefined zone",
			platform: validPlatform(),
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
		{
			name:     "data disk valid config",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						DataDisks: []vsphere.DataDisk{
							{
								Name:             "Disk1",
								SizeGiB:          10,
								ProvisioningMode: vsphere.ProvisioningModeThin,
							},
						},
					},
				},
			},
		},
		{
			name:     "data disk name not set",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						DataDisks: []vsphere.DataDisk{
							{
								SizeGiB:          10,
								ProvisioningMode: vsphere.ProvisioningModeThin,
							},
						},
					},
				},
			},
			expectedErrMsg: "test-path.disks\\[0].name: Required value: data disk name must be set",
		},
		{
			name:     "data disk name invalid characters",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						DataDisks: []vsphere.DataDisk{
							{
								Name:             "bad disk name",
								SizeGiB:          10,
								ProvisioningMode: vsphere.ProvisioningModeThin,
							},
						},
					},
				},
			},
			expectedErrMsg: "test-path.disks\\[0].name: Invalid value: \"bad disk name\": data disk name must consist only of alphanumeric characters, hyphens and underscores, and must start and end with an alphanumeric character.",
		},
		{
			name:     "data disk size too large",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						DataDisks: []vsphere.DataDisk{
							{
								Name:             "Disk1",
								SizeGiB:          20000,
								ProvisioningMode: vsphere.ProvisioningModeThin,
							},
						},
					},
				},
			},
			expectedErrMsg: "test-path.disks\\[0].sizeGiB: Invalid value: 20000: data disk size \\(GiB\\) must not exceed 16384",
		},
		{
			name:     "data disk size not set",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						DataDisks: []vsphere.DataDisk{
							{
								Name:             "Disk1",
								ProvisioningMode: vsphere.ProvisioningModeThin,
							},
						},
					},
				},
			},
			expectedErrMsg: "test-path.disks\\[0].sizeGiB: Required value: data disk size must be set",
		},
		{
			name:     "data disk provisioning mode not set",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						DataDisks: []vsphere.DataDisk{
							{
								Name:    "Disk1",
								SizeGiB: 10,
							},
						},
					},
				},
			},
		},
		{
			name:     "data disk invalid provisioning mode",
			platform: validPlatform(),
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					VSphere: &vsphere.MachinePool{
						DataDisks: []vsphere.DataDisk{
							{
								Name:             "Disk1",
								SizeGiB:          10,
								ProvisioningMode: "Fake",
							},
						},
					},
				},
			},
			expectedErrMsg: "test-path.disks\\[0]: Unsupported value: \"Fake\": supported values: \"EagerlyZeroed\", \"Thick\", \"Thin\"",
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
