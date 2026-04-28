package validation

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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
			name: "valid custom CPU, memory, and disk configuration",
			pool: &nutanix.MachinePool{
				NumCPUs:           8,
				NumCoresPerSocket: 2,
				MemoryMiB:         16384,
				OSDisk: nutanix.OSDisk{
					DiskSizeGiB: 120,
				},
			},
			expectedErrMsg: "",
		}, {
			name: "valid custom CPU and memory for master",
			role: "master",
			pool: &nutanix.MachinePool{
				NumCPUs:           4,
				NumCoresPerSocket: 2,
				MemoryMiB:         32768,
				OSDisk: nutanix.OSDisk{
					DiskSizeGiB: 200,
				},
			},
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
			name: "multiple gpus not supported for master nodes",
			role: "master",
			pool: &nutanix.MachinePool{
				GPUs: []machinev1.NutanixGPU{
					{Type: machinev1.NutanixGPUIdentifierName, Name: ptr.To("gpu-1")},
					{Type: machinev1.NutanixGPUIdentifierDeviceID, DeviceID: ptr.To(int32(42))},
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
		}, {
			name: "valid multiple data disks on worker",
			role: "worker",
			pool: &nutanix.MachinePool{
				DataDisks: []nutanix.DataDisk{
					{
						DiskSize: resource.MustParse("100Gi"),
						DeviceProperties: &machinev1.NutanixVMDiskDeviceProperties{
							DeviceType:  machinev1.NutanixDiskDeviceTypeDisk,
							AdapterType: machinev1.NutanixDiskAdapterTypeSCSI,
							DeviceIndex: 0,
						},
					},
					{
						DiskSize: resource.MustParse("200Gi"),
						DeviceProperties: &machinev1.NutanixVMDiskDeviceProperties{
							DeviceType:  machinev1.NutanixDiskDeviceTypeDisk,
							AdapterType: machinev1.NutanixDiskAdapterTypePCI,
							DeviceIndex: 1,
						},
					},
					{
						DiskSize: resource.MustParse("50Gi"),
					},
				},
			},
			expectedErrMsg: "",
		}, {
			name: "valid dataDisk with CDRom device type and IDE adapter",
			role: "worker",
			pool: &nutanix.MachinePool{
				DataDisks: []nutanix.DataDisk{{
					DiskSize: resource.MustParse("1Gi"),
					DeviceProperties: &machinev1.NutanixVMDiskDeviceProperties{
						DeviceType:  machinev1.NutanixDiskDeviceTypeCDROM,
						AdapterType: machinev1.NutanixDiskAdapterTypeIDE,
						DeviceIndex: 0,
					},
				}},
			},
			expectedErrMsg: "",
		}, {
			name: "invalid dataDisk adapter type for Disk device",
			role: "worker",
			pool: &nutanix.MachinePool{
				DataDisks: []nutanix.DataDisk{{
					DiskSize: resource.MustParse("1Gi"),
					DeviceProperties: &machinev1.NutanixVMDiskDeviceProperties{
						DeviceType:  machinev1.NutanixDiskDeviceTypeDisk,
						AdapterType: "InvalidAdapter",
						DeviceIndex: 0,
					},
				}},
			},
			expectedErrMsg: `invalid adapter type for the "Disk" device type`,
		}, {
			name: "invalid dataDisk adapter type for CDRom device",
			role: "worker",
			pool: &nutanix.MachinePool{
				DataDisks: []nutanix.DataDisk{{
					DiskSize: resource.MustParse("1Gi"),
					DeviceProperties: &machinev1.NutanixVMDiskDeviceProperties{
						DeviceType:  machinev1.NutanixDiskDeviceTypeCDROM,
						AdapterType: machinev1.NutanixDiskAdapterTypeSCSI,
						DeviceIndex: 0,
					},
				}},
			},
			expectedErrMsg: `invalid adapter type for the "CDRom" device type`,
		}, {
			name: "invalid dataDisk device type",
			role: "worker",
			pool: &nutanix.MachinePool{
				DataDisks: []nutanix.DataDisk{{
					DiskSize: resource.MustParse("1Gi"),
					DeviceProperties: &machinev1.NutanixVMDiskDeviceProperties{
						DeviceType:  "InvalidType",
						AdapterType: machinev1.NutanixDiskAdapterTypeSCSI,
						DeviceIndex: 0,
					},
				}},
			},
			expectedErrMsg: `invalid device type, the valid types are`,
		}, {
			name: "invalid dataDisk storage config disk mode",
			role: "worker",
			pool: &nutanix.MachinePool{
				DataDisks: []nutanix.DataDisk{{
					DiskSize: resource.MustParse("1Gi"),
					StorageConfig: &nutanix.StorageConfig{
						DiskMode: "InvalidMode",
					},
				}},
			},
			expectedErrMsg: `invalid disk mode, the valid values`,
		}, {
			name: "valid bootType Legacy",
			pool: &nutanix.MachinePool{
				BootType: machinev1.NutanixLegacyBoot,
			},
			expectedErrMsg: "",
		}, {
			name: "valid bootType UEFI",
			pool: &nutanix.MachinePool{
				BootType: machinev1.NutanixUEFIBoot,
			},
			expectedErrMsg: "",
		}, {
			name: "valid bootType SecureBoot",
			pool: &nutanix.MachinePool{
				BootType: machinev1.NutanixSecureBoot,
			},
			expectedErrMsg: "",
		}, {
			name: "invalid bootType",
			pool: &nutanix.MachinePool{
				BootType: "InvalidBoot",
			},
			expectedErrMsg: `valid bootType: "", "Legacy", "UEFI", "SecureBoot".`,
		}, {
			name: "project with invalid identifier type",
			pool: &nutanix.MachinePool{
				Project: &machinev1.NutanixResourceIdentifier{
					Type: "invalid",
				},
			},
			expectedErrMsg: `invalid project identifier type`,
		}, {
			name: "project with name identifier missing name",
			pool: &nutanix.MachinePool{
				Project: &machinev1.NutanixResourceIdentifier{
					Type: machinev1.NutanixIdentifierName,
				},
			},
			expectedErrMsg: `missing project name`,
		}, {
			name: "project with uuid identifier missing uuid",
			pool: &nutanix.MachinePool{
				Project: &machinev1.NutanixResourceIdentifier{
					Type: machinev1.NutanixIdentifierUUID,
				},
			},
			expectedErrMsg: `missing project uuid`,
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

// testServerPort extracts the port from an httptest.Server as an int32.
func testServerPort(t *testing.T, ts *httptest.Server) int32 {
	t.Helper()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("failed to parse test server URL: %v", err)
	}
	port, err := strconv.ParseInt(u.Port(), 10, 32)
	if err != nil {
		t.Fatalf("failed to parse test server port: %v", err)
	}
	return int32(port)
}

// newFakePlatform returns a *nutanix.Platform configured to use a local httptest server.
// The http:// scheme is required: the Nutanix client defaults to HTTPS for
// schemeless URLs (WithBaseURL prepends "https://"), which would fail against
// the plain-HTTP test server. This address would not pass ValidatePlatform,
// but ValidateConfig (called here) does not re-validate the endpoint format.
func newFakePlatform(t *testing.T, peUUID string, port int32) *nutanix.Platform {
	t.Helper()
	return &nutanix.Platform{
		PrismCentral: nutanix.PrismCentral{
			Endpoint: nutanix.PrismEndpoint{
				Address: "http://127.0.0.1",
				Port:    port,
			},
			Username: "test-user",
			Password: "test-pass",
		},
		PrismElements: []nutanix.PrismElement{{
			UUID:     peUUID,
			Endpoint: nutanix.PrismEndpoint{Address: "test-pe", Port: 8081},
		}},
		SubnetUUIDs: []string{"test-subnet-uuid"},
	}
}

// fakeNutanixGPUServer returns an httptest.Server that responds to POST /api/nutanix/v3/hosts/list
// with a single host containing one UNUSED GPU (Name="Tesla-T4", DeviceID=1234).
func fakeNutanixGPUServer(t *testing.T, peUUID string) *httptest.Server {
	t.Helper()

	gpuResponse, err := json.Marshal(map[string]interface{}{
		"metadata": map[string]interface{}{"total_matches": 1, "length": 1, "offset": 0, "kind": "host"},
		"entities": []map[string]interface{}{{
			"status": map[string]interface{}{
				"cluster_reference": map[string]interface{}{"uuid": peUUID},
				"resources": map[string]interface{}{
					"gpu_list": []map[string]interface{}{{
						"device_id": 1234,
						"name":      "Tesla-T4",
						"status":    "UNUSED",
						"vendor":    "NVIDIA",
						"mode":      "PASSTHROUGH_GRAPHICS",
					}},
				},
			},
		}},
	})
	if err != nil {
		t.Fatalf("failed to encode GPU response: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/nutanix/v3/hosts/list", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(gpuResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return httptest.NewServer(mux)
}

// TestValidateMachinePoolGPUFieldValidation uses a fake Nutanix API server so
// that validateGPUsConfig proceeds past the inventory check and reaches per-field
// validation (invalid type, missing name, missing deviceID, not-found, etc.).
func TestValidateMachinePoolGPUFieldValidation(t *testing.T) {
	const peUUID = "test-pe-uuid"
	ts := fakeNutanixGPUServer(t, peUUID)
	defer ts.Close()

	port := testServerPort(t, ts)
	platform := newFakePlatform(t, peUUID, port)

	cases := []struct {
		name           string
		pool           *nutanix.MachinePool
		expectedErrMsg string
	}{
		{
			name: "gpu with invalid identifier type",
			pool: &nutanix.MachinePool{
				GPUs: []machinev1.NutanixGPU{
					{Type: "InvalidType", DeviceID: ptr.To(int32(7864))},
				},
			},
			expectedErrMsg: `invalid gpu identifier type`,
		},
		{
			name: "gpu with nil deviceID for DeviceID type",
			pool: &nutanix.MachinePool{
				GPUs: []machinev1.NutanixGPU{
					{Type: machinev1.NutanixGPUIdentifierDeviceID},
				},
			},
			expectedErrMsg: `missing gpu deviceID`,
		},
		{
			name: "gpu with nil name for Name type",
			pool: &nutanix.MachinePool{
				GPUs: []machinev1.NutanixGPU{
					{Type: machinev1.NutanixGPUIdentifierName},
				},
			},
			expectedErrMsg: `missing gpu name`,
		},
		{
			name: "gpu with empty name string",
			pool: &nutanix.MachinePool{
				GPUs: []machinev1.NutanixGPU{
					{Type: machinev1.NutanixGPUIdentifierName, Name: ptr.To("")},
				},
			},
			expectedErrMsg: `missing gpu name`,
		},
		{
			name: "gpu name not found in inventory",
			pool: &nutanix.MachinePool{
				GPUs: []machinev1.NutanixGPU{
					{Type: machinev1.NutanixGPUIdentifierName, Name: ptr.To("nonexistent-gpu")},
				},
			},
			expectedErrMsg: `no available GPU found that matches required GPU inputs`,
		},
		{
			name: "gpu deviceID not found in inventory",
			pool: &nutanix.MachinePool{
				GPUs: []machinev1.NutanixGPU{
					{Type: machinev1.NutanixGPUIdentifierDeviceID, DeviceID: ptr.To(int32(9999))},
				},
			},
			expectedErrMsg: `no available GPU found that matches required GPU inputs`,
		},
		{
			name: "gpu with mixed invalid configs",
			pool: &nutanix.MachinePool{
				GPUs: []machinev1.NutanixGPU{
					{Type: "InvalidType", DeviceID: ptr.To(int32(7864))},
					{Type: machinev1.NutanixGPUIdentifierName, Name: ptr.To("nonexistent-gpu")},
					{Type: machinev1.NutanixGPUIdentifierDeviceID, DeviceID: ptr.To(int32(9999))},
				},
			},
			expectedErrMsg: `invalid gpu identifier type`,
		},
		{
			name: "gpu name found in inventory succeeds",
			pool: &nutanix.MachinePool{
				GPUs: []machinev1.NutanixGPU{
					{Type: machinev1.NutanixGPUIdentifierName, Name: ptr.To("Tesla-T4")},
				},
			},
			expectedErrMsg: "",
		},
		{
			name: "gpu deviceID found in inventory succeeds",
			pool: &nutanix.MachinePool{
				GPUs: []machinev1.NutanixGPU{
					{Type: machinev1.NutanixGPUIdentifierDeviceID, DeviceID: ptr.To(int32(1234))},
				},
			},
			expectedErrMsg: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.pool.ValidateConfig(platform, "worker")
			if tc.expectedErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			}
		})
	}
}

// TestValidateMachinePoolCategoryValidation uses a fake Nutanix API server so
// that category validation reaches the GetCategoryValue call and receives a
// deterministic 404 response, avoiding any real network calls.
func TestValidateMachinePoolCategoryValidation(t *testing.T) {
	const peUUID = "test-pe-uuid"

	categoryNotFound, err := json.Marshal(map[string]interface{}{
		"state":        "ERROR",
		"message_list": []map[string]interface{}{{"message": "category not found"}},
	})
	if err != nil {
		t.Fatalf("failed to encode category response: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/nutanix/v3/categories/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write(categoryNotFound); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	port := testServerPort(t, ts)
	platform := newFakePlatform(t, peUUID, port)

	cases := []struct {
		name           string
		pool           *nutanix.MachinePool
		expectedErrMsg string
	}{
		{
			name: "category not found returns deterministic error",
			pool: &nutanix.MachinePool{
				Categories: []machinev1.NutanixCategory{
					{Key: "test-key", Value: "test-value"},
				},
			},
			expectedErrMsg: `Failed to find the category with key`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.pool.ValidateConfig(platform, "worker")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErrMsg)
		})
	}
}
