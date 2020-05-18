package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/gcp"
)

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *gcp.Platform
		valid    bool
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
				Region: "bad-region",
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
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
