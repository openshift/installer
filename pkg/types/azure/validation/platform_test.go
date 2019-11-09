package validation

import (
	"testing"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *azure.Platform
		valid    bool
	}{
		{
			name: "invalid region",
			platform: &azure.Platform{
				Region:                      "",
				BaseDomainResourceGroupName: "group",
			},
			valid: false,
		},
		{
			name: "invalid baseDomainResourceGroupName",
			platform: &azure.Platform{
				Region:                      "eastus",
				BaseDomainResourceGroupName: "",
			},
			valid: false,
		},
		{
			name: "minimal",
			platform: &azure.Platform{
				Region:                      "eastus",
				BaseDomainResourceGroupName: "group",
			},
			valid: true,
		},
		{
			name: "valid machine pool",
			platform: &azure.Platform{
				Region:                      "eastus",
				BaseDomainResourceGroupName: "group",
				DefaultMachinePlatform:      &azure.MachinePool{},
			},
			valid: true,
		},
		{
			name: "valid subnets & virtual network",
			platform: &azure.Platform{
				Region:                      "eastus",
				BaseDomainResourceGroupName: "group",
				NetworkResourceGroupName:    "networkresourcegroup",
				VirtualNetwork:              "virtualnetwork",
				ComputeSubnet:               "computesubnet",
				ControlPlaneSubnet:          "controlplanesubnet",
			},
			valid: true,
		},
		{
			name: "missing subnets",
			platform: &azure.Platform{
				Region:                   "eastus",
				NetworkResourceGroupName: "networkresourcegroup",
				VirtualNetwork:           "virtualnetwork",
			},
			valid: false,
		},
		{
			name: "subnets missing virtual network",
			platform: &azure.Platform{
				Region:                   "eastus",
				NetworkResourceGroupName: "networkresourcegroup",
				ComputeSubnet:            "computesubnet",
			},
			valid: false,
		},
		{
			name: "missing network resource group",
			platform: &azure.Platform{
				Region:             "eastus",
				VirtualNetwork:     "virtualnetwork",
				ComputeSubnet:      "computesubnet",
				ControlPlaneSubnet: "controlplanesubnet",
			},
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, types.ExternalPublishingStrategy, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
