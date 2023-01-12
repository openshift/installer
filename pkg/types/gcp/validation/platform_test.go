package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name            string
		platform        *gcp.Platform
		credentialsMode types.CredentialsMode
		valid           bool
	}{
		{
			name: "minimal",
			platform: &gcp.Platform{
				Region: "us-east1",
			},
			valid: true,
		},
		{
			name: "invalid region",
			platform: &gcp.Platform{
				Region: "",
			},
			valid: false,
		},
		{
			name: "valid machine pool",
			platform: &gcp.Platform{
				Region:                 "us-east1",
				DefaultMachinePlatform: &gcp.MachinePool{},
			},
			valid: true,
		},
		{
			name: "valid subnets & network",
			platform: &gcp.Platform{
				Region:             "us-east1",
				Network:            "valid-vpc",
				ComputeSubnet:      "valid-compute-subnet",
				ControlPlaneSubnet: "valid-cp-subnet",
			},
			valid: true,
		},
		{
			name: "missing subnets",
			platform: &gcp.Platform{
				Region:  "us-east1",
				Network: "valid-vpc",
			},
			valid: false,
		},
		{
			name: "subnets missing network",
			platform: &gcp.Platform{
				Region:        "us-east1",
				ComputeSubnet: "valid-compute-subnet",
			},
			valid: false,
		},
		{
			name: "unsupported GCP disk type",
			platform: &gcp.Platform{
				Region: "us-east1",
				DefaultMachinePlatform: &gcp.MachinePool{
					OSDisk: gcp.OSDisk{
						DiskType: "pd-standard",
					},
				},
			},
			valid: false,
		},

		{
			name: "supported GCP disk type",
			platform: &gcp.Platform{
				Region: "us-east1",
				DefaultMachinePlatform: &gcp.MachinePool{
					OSDisk: gcp.OSDisk{
						DiskType: "pd-ssd",
					},
				},
			},
			valid: true,
		},
		{
			name: "GCP valid network project data",
			platform: &gcp.Platform{
				Region:             "us-east1",
				NetworkProjectID:   "valid-network-project",
				ProjectID:          "valid-project",
				Network:            "valid-vpc",
				ComputeSubnet:      "valid-compute-subnet",
				ControlPlaneSubnet: "valid-cp-subnet",
			},
			credentialsMode: types.PassthroughCredentialsMode,
			valid:           true,
		},
		{
			name: "GCP invalid network project missing network",
			platform: &gcp.Platform{
				Region:             "us-east1",
				NetworkProjectID:   "valid-network-project",
				ProjectID:          "valid-project",
				ComputeSubnet:      "valid-compute-subnet",
				ControlPlaneSubnet: "valid-cp-subnet",
			},
			credentialsMode: types.PassthroughCredentialsMode,
			valid:           false,
		},
		{
			name: "GCP invalid network project missing compute subnet",
			platform: &gcp.Platform{
				Region:             "us-east1",
				NetworkProjectID:   "valid-network-project",
				ProjectID:          "valid-project",
				Network:            "valid-vpc",
				ControlPlaneSubnet: "valid-cp-subnet",
			},
			credentialsMode: types.PassthroughCredentialsMode,
			valid:           false,
		},
		{
			name: "GCP invalid network project missing control plane subnet",
			platform: &gcp.Platform{
				Region:           "us-east1",
				NetworkProjectID: "valid-network-project",
				ProjectID:        "valid-project",
				Network:          "valid-vpc",
				ComputeSubnet:    "valid-compute-subnet",
			},
			credentialsMode: types.PassthroughCredentialsMode,
			valid:           false,
		},
		{
			name: "GCP invalid network project bad credentials mode",
			platform: &gcp.Platform{
				Region:             "us-east1",
				NetworkProjectID:   "valid-network-project",
				ProjectID:          "valid-project",
				Network:            "valid-vpc",
				ComputeSubnet:      "valid-compute-subnet",
				ControlPlaneSubnet: "valid-cp-subnet",
			},
			credentialsMode: types.MintCredentialsMode,
			valid:           false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			credentialsMode := tc.credentialsMode
			if credentialsMode == "" {
				credentialsMode = types.MintCredentialsMode
			}

			// the only item currently used is the credentialsMode
			ic := types.InstallConfig{
				CredentialsMode: credentialsMode,
			}

			err := ValidatePlatform(tc.platform, field.NewPath("test-path"), &ic).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
