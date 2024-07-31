package validation

import (
	"testing"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"

	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/types/nutanix"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name           string
		role           string
		pool           *nutanix.MachinePool
		expectedErrMsg string
	}{
		{
			name:           "empty",
			pool:           &nutanix.MachinePool{},
			expectedErrMsg: "",
		}, {
			name: "negative disk size",
			pool: &nutanix.MachinePool{
				OSDisk: nutanix.OSDisk{
					DiskSizeGiB: -1,
				},
			},
			expectedErrMsg: `test-path.diskSizeGiB: Invalid value: -1: storage disk size must be positive`,
		}, {
			name: "negative CPUs",
			pool: &nutanix.MachinePool{
				NumCPUs: -1,
			},
			expectedErrMsg: `test-path.cpus: Invalid value: -1: number of CPUs must be positive`,
		}, {
			name: "negative cores",
			pool: &nutanix.MachinePool{
				NumCoresPerSocket: -1,
			},
			expectedErrMsg: `test-path.coresPerSocket: Invalid value: -1: cores per socket must be positive`,
		}, {
			name: "negative memory",
			pool: &nutanix.MachinePool{
				MemoryMiB: -1,
			},
			expectedErrMsg: `test-path.memoryMiB: Invalid value: -1: memory size must be positive`,
		}, {
			name: "less CPUs than cores per socket",
			pool: &nutanix.MachinePool{
				NumCPUs:           1,
				NumCoresPerSocket: 8,
			},
			expectedErrMsg: `test-path.coresPerSocket: Invalid value: 8: cores per socket must be less than number of CPUs`,
		}, {
			name: "gpus not supported for master nodes",
			role: "master",
			pool: &nutanix.MachinePool{
				GPUs: []machinev1.NutanixGPU{
					{Type: machinev1.NutanixGPUIdentifierName, Name: ptr.To("gpu-1")},
				},
			},
			expectedErrMsg: `'gpus' are not supported for 'master' nodes`,
		}, {
			name: "dataDisks not supported for master nodes",
			role: "master",
			pool: &nutanix.MachinePool{
				DataDisks: []nutanix.DataDisk{{
					DiskSize: resource.MustParse("1Gi"),
				}},
			},
			expectedErrMsg: `'dataDisks' are not supported for 'master' nodes`,
		}, {
			name: "dataDisk size less than 1GB",
			role: "worker",
			pool: &nutanix.MachinePool{
				DataDisks: []nutanix.DataDisk{{
					DiskSize: resource.MustParse("0.5Gi"),
				}},
			},
			expectedErrMsg: `The minimum diskSize is 1Gi bytes.`,
		}, {
			name: "negative dataDisk deviceIndex",
			role: "worker",
			pool: &nutanix.MachinePool{
				DataDisks: []nutanix.DataDisk{{
					DiskSize: resource.MustParse("1Gi"),
					DeviceProperties: &machinev1.NutanixVMDiskDeviceProperties{
						DeviceType:  machinev1.NutanixDiskDeviceTypeDisk,
						AdapterType: machinev1.NutanixDiskAdapterTypeSCSI,
						DeviceIndex: int32(-1),
					},
				}},
			},
			expectedErrMsg: `invalid device index, the valid values are non-negative integers.`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gs := gomega.NewWithT(t)

			err := ValidateMachinePool(tc.pool, field.NewPath("test-path"), tc.role).ToAggregate()
			if tc.expectedErrMsg == "" {
				assert.NoError(t, err)
			} else {
				gs.Expect(err.Error()).To(gomega.ContainSubstring(tc.expectedErrMsg))
			}
		})
	}
}
