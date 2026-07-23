package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/gcp"
)

func TestValidateOSImageForSovereignCloud(t *testing.T) {
	validOSImage := &gcp.OSImage{Name: "my-image", Project: "my-project"}

	cases := []struct {
		name          string
		platform      *gcp.Platform
		pool          *gcp.MachinePool
		expectedError string
	}{
		{
			name:     "non-sovereign cloud skips validation",
			platform: &gcp.Platform{ProjectID: "my-project"},
			pool:     &gcp.MachinePool{},
		},
		{
			name:     "sovereign cloud with os image on pool",
			platform: &gcp.Platform{ProjectID: "eu0:my-project"},
			pool:     &gcp.MachinePool{OSImage: validOSImage},
		},
		{
			name:          "sovereign cloud missing os image on pool",
			platform:      &gcp.Platform{ProjectID: "eu0:my-project"},
			pool:          &gcp.MachinePool{},
			expectedError: `test-path.osImage: Required value: must specify an OS image for sovereign cloud environments (domain-scoped project ID)`,
		},
		{
			name:          "sovereign cloud nil pool",
			platform:      &gcp.Platform{ProjectID: "eu0:my-project"},
			pool:          nil,
			expectedError: `test-path.osImage: Required value: must specify an OS image for sovereign cloud environments (domain-scoped project ID)`,
		},
		{
			name:          "sovereign cloud os image missing name",
			platform:      &gcp.Platform{ProjectID: "eu0:my-project"},
			pool:          &gcp.MachinePool{OSImage: &gcp.OSImage{Project: "my-project"}},
			expectedError: `test-path.osImage.name: Required value: must specify an OS image name for sovereign cloud environments`,
		},
		{
			name:          "sovereign cloud os image missing project",
			platform:      &gcp.Platform{ProjectID: "eu0:my-project"},
			pool:          &gcp.MachinePool{OSImage: &gcp.OSImage{Name: "my-image"}},
			expectedError: `test-path.osImage.project: Required value: must specify an OS image project for sovereign cloud environments`,
		},
		{
			name:     "org-scoped project ID is not sovereign",
			platform: &gcp.Platform{ProjectID: "myorg:my-project"},
			pool:     &gcp.MachinePool{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateOSImageForSovereignCloud(tc.platform, tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

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
