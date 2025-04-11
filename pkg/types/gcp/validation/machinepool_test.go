package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/gcp"
)

func TestValidateMachinePool(t *testing.T) {
	platform := &gcp.Platform{Region: "us-east1"}
	cases := []struct {
		name     string
		pool     *gcp.MachinePool
		expected string
	}{
		{
			name: "empty",
			pool: &gcp.MachinePool{},
		},
		{
			name: "valid zone",
			pool: &gcp.MachinePool{
				Zones: []string{"us-east1-b", "us-east1-c"},
			},
		},
		{
			name: "invalid zone",
			pool: &gcp.MachinePool{
				Zones: []string{"us-east1-b", "us-central1-f"},
			},
			expected: `^test-path\.zones\[1]: Invalid value: "us-central1-f": Zone not in configured region \(us-east1\)$`,
		},
		{
			name: "valid disk type",
			pool: &gcp.MachinePool{
				OSDisk: gcp.OSDisk{
					DiskType: "pd-standard",
				},
			},
		},
		{
			name: "invalid disk type",
			pool: &gcp.MachinePool{
				OSDisk: gcp.OSDisk{
					DiskType: "pd-",
				},
			},
			expected: `^test-path\.diskType: Unsupported value: "pd-": supported values: "hyperdisk-balanced", "pd-balanced", "pd-ssd", "pd-standard"$`,
		},
		{
			name: "valid disk size",
			pool: &gcp.MachinePool{
				OSDisk: gcp.OSDisk{
					DiskSizeGB: 100,
				},
			},
		},
		{
			name: "invalid disk size",
			pool: &gcp.MachinePool{
				OSDisk: gcp.OSDisk{
					DiskSizeGB: -120,
				},
			},
			expected: `^test-path\.diskSizeGB: Invalid value: -120: must be at least 16GB in size$`,
		},
		{
			name: "insufficient disk size",
			pool: &gcp.MachinePool{
				OSDisk: gcp.OSDisk{
					DiskSizeGB: 11,
				},
			},
			expected: `^test-path\.diskSizeGB: Invalid value: 11: must be at least 16GB in size$`,
		},
		{
			name: "exceeded disk size",
			pool: &gcp.MachinePool{
				OSDisk: gcp.OSDisk{
					DiskSizeGB: 66000,
				},
			},
			expected: `^test-path\.diskSizeGB: Invalid value: 66000: exceeding maximum GCP disk size limit, must be below 65536$`,
		},
		{
			name: "enable confidential compute with correct on host maintenance",
			pool: &gcp.MachinePool{
				InstanceType:        "c2d-standard-4",
				ConfidentialCompute: string(gcp.EnabledFeature),
				OnHostMaintenance:   string(gcp.OnHostMaintenanceTerminate),
			},
		},
		{
			name: "enable confidential compute with incorrect on host maintenance",
			pool: &gcp.MachinePool{
				ConfidentialCompute: string(gcp.EnabledFeature),
				OnHostMaintenance:   string(gcp.OnHostMaintenanceMigrate),
			},
			expected: `test-path.OnHostMaintenance: Invalid value: "Migrate": OnHostMaintenace must be set to Terminate when ConfidentialCompute is Enabled`,
		},
		{
			name: "AMDEncryptedVirtualization confidential compute with incorrect on host maintenance",
			pool: &gcp.MachinePool{
				ConfidentialCompute: string(gcp.ConfidentialComputePolicySEV),
				OnHostMaintenance:   string(gcp.OnHostMaintenanceMigrate),
			},
			expected: `test-path.OnHostMaintenance: Invalid value: "Migrate": OnHostMaintenace must be set to Terminate when ConfidentialCompute is AMDEncryptedVirtualization`,
		},
		{
			name: "AMDEncryptedVirtualizationNestedPaging confidential compute with incorrect on host maintenance",
			pool: &gcp.MachinePool{
				ConfidentialCompute: string(gcp.ConfidentialComputePolicySEVSNP),
				OnHostMaintenance:   string(gcp.OnHostMaintenanceMigrate),
			},
			expected: `test-path.OnHostMaintenance: Invalid value: "Migrate": OnHostMaintenace must be set to Terminate when ConfidentialCompute is AMDEncryptedVirtualizationNestedPaging`,
		},
		{
			name: "IntelTrustedDomainExtensions confidential compute with incorrect on host maintenance",
			pool: &gcp.MachinePool{
				ConfidentialCompute: string(gcp.ConfidentialComputePolicyTDX),
				OnHostMaintenance:   string(gcp.OnHostMaintenanceMigrate),
			},
			expected: `test-path.OnHostMaintenance: Invalid value: "Migrate": OnHostMaintenace must be set to Terminate when ConfidentialCompute is IntelTrustedDomainExtensions`,
		},
		{
			name: "AMDEncryptedVirtualization confidential compute with unsupported machine type",
			pool: &gcp.MachinePool{
				InstanceType:        "c3-standard-4",
				ConfidentialCompute: string(gcp.ConfidentialComputePolicySEV),
				OnHostMaintenance:   string(gcp.OnHostMaintenanceTerminate),
			},
			expected: `test-path.type: Invalid value: "c3-standard-4": Machine type do not support AMDEncryptedVirtualization. Machine types supporting AMDEncryptedVirtualization: c2d, n2d, c3d`,
		},
		{
			name: "AMDEncryptedVirtualization confidential compute with supported machine type",
			pool: &gcp.MachinePool{
				InstanceType:        "c3d-standard-4",
				ConfidentialCompute: string(gcp.ConfidentialComputePolicySEV),
				OnHostMaintenance:   string(gcp.OnHostMaintenanceTerminate),
			},
		},
		{
			name: "AMDEncryptedVirtualizationNestedPaging confidential compute with unsupported machine type",
			pool: &gcp.MachinePool{
				InstanceType:        "c2d-standard-4",
				ConfidentialCompute: string(gcp.ConfidentialComputePolicySEVSNP),
				OnHostMaintenance:   string(gcp.OnHostMaintenanceTerminate),
			},
			expected: `test-path.type: Invalid value: "c2d-standard-4": Machine type do not support AMDEncryptedVirtualizationNestedPaging. Machine types supporting AMDEncryptedVirtualizationNestedPaging: n2d`,
		},
		{
			name: "AMDEncryptedVirtualizationNestedPaging confidential compute with supported machine type",
			pool: &gcp.MachinePool{
				InstanceType:        "n2d-standard-4",
				ConfidentialCompute: string(gcp.ConfidentialComputePolicySEVSNP),
				OnHostMaintenance:   string(gcp.OnHostMaintenanceTerminate),
			},
		},
		{
			name: "IntelTrustedDomainExtensions confidential compute with unsupported machine type",
			pool: &gcp.MachinePool{
				InstanceType:        "n2d-standard-4",
				ConfidentialCompute: string(gcp.ConfidentialComputePolicyTDX),
				OnHostMaintenance:   string(gcp.OnHostMaintenanceTerminate),
			},
			expected: `test-path.type: Invalid value: "n2d-standard-4": Machine type do not support IntelTrustedDomainExtensions. Machine types supporting IntelTrustedDomainExtensions: c3`,
		},
		{
			name: "IntelTrustedDomainExtensions confidential compute with supported machine type",
			pool: &gcp.MachinePool{
				InstanceType:        "c3-standard-4",
				ConfidentialCompute: string(gcp.ConfidentialComputePolicyTDX),
				OnHostMaintenance:   string(gcp.OnHostMaintenanceTerminate),
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(platform, tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}
