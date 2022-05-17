package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/vsphere"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name           string
		pool           *vsphere.MachinePool
		platform       *vsphere.Platform
		expectedErrMsg string
	}{
		{
			name:           "empty",
			pool:           &vsphere.MachinePool{},
			platform:       validPlatform(),
			expectedErrMsg: "",
		}, {
			name:     "negative disk size",
			platform: validPlatform(),
			pool: &vsphere.MachinePool{
				OSDisk: vsphere.OSDisk{
					DiskSizeGB: -1,
				},
			},
			expectedErrMsg: `^test-path\.diskSizeGB: Invalid value: -1: storage disk size must be positive$`,
		}, {
			name:     "negative CPUs",
			platform: validPlatform(),
			pool: &vsphere.MachinePool{
				NumCPUs: -1,
			},
			expectedErrMsg: `^test-path\.cpus: Invalid value: -1: number of CPUs must be positive$`,
		}, {
			name:     "negative cores",
			platform: validPlatform(),
			pool: &vsphere.MachinePool{
				NumCoresPerSocket: -1,
			},
			expectedErrMsg: `^test-path\.coresPerSocket: Invalid value: -1: cores per socket must be positive$`,
		}, {
			name:     "negative memory",
			platform: validPlatform(),
			pool: &vsphere.MachinePool{
				MemoryMiB: -1,
			},
			expectedErrMsg: `^test-path\.memoryMB: Invalid value: -1: memory size must be positive$`,
		}, {
			name:     "less CPUs than cores per socket",
			platform: validPlatform(),
			pool: &vsphere.MachinePool{
				NumCPUs:           1,
				NumCoresPerSocket: 8,
			},
			expectedErrMsg: `^test-path\.coresPerSocket: Invalid value: 8: cores per socket must be less than number of CPUs$`,
		},
		{
			name:     "numCPUs not a multiple of cores per socket",
			platform: validPlatform(),
			pool: &vsphere.MachinePool{
				NumCPUs:           7,
				NumCoresPerSocket: 4,
			},
			expectedErrMsg: `^test-path.cpus: Invalid value: 7: numCPUs specified should be a multiple of cores per socket$`,
		},
		{
			name:     "numCPUs not a multiple of default cores per socket",
			platform: validPlatform(),
			pool: &vsphere.MachinePool{
				NumCPUs: 7,
			},
			expectedErrMsg: `^test-path.cpus: Invalid value: 7: numCPUs specified should be a multiple of cores per socket which is by default 4$`,
		},
		{
			name: "multi-zone invalid zone name",
			platform: func() *vsphere.Platform {
				platform := validMultiVCenterPlatform()
				platform.DeploymentZones[0].Name = "Zone%^@112233"
				return platform
			}(),
			pool: &vsphere.MachinePool{
				Zones: []string{
					"Zone%^@112233",
				},
			},
			expectedErrMsg: `^test-path.zones: Invalid value: \[\]string{"Zone%\^@112233"}: cluster name must begin with a lower-case letter$`,
		},
		{
			name:     "multi-zone valid",
			platform: validMultiVCenterPlatform(),
			pool: &vsphere.MachinePool{
				Zones: []string{
					"test-dz-east-1a",
				},
			},
		},
		{
			name:     "multi-zone missing zones",
			platform: validMultiVCenterPlatform(),
			pool: &vsphere.MachinePool{
				Zones: []string{},
			},
			expectedErrMsg: `^test-path.zones: Required value: zones must be defined if deploymentZones are defined$`,
		},
		{
			name:     "multi-zone undefined zone",
			platform: validMultiVCenterPlatform(),
			pool: &vsphere.MachinePool{
				Zones: []string{
					"unknown-zone",
				},
			},
			expectedErrMsg: `^test-path.zones: Invalid value: "unknown-zone": zone not defined in deploymentZones$`,
		},
		{
			name: "multi-zone missing deploymentZones",
			platform: func() *vsphere.Platform {
				platform := validMultiVCenterPlatform()
				platform.DeploymentZones = make([]vsphere.DeploymentZoneSpec, 0)
				return platform
			}(),
			pool: &vsphere.MachinePool{
				Zones: []string{
					"test-dz-east-1a",
				},
			},
			expectedErrMsg: `^test-path.zones: Required value: deploymentZones must be defined if zones are defined$`,
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
		})
	}
}
