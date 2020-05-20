package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
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
				CloudName:                   azure.PublicCloud,
			},
			valid: false,
		},
		{
			name: "invalid baseDomainResourceGroupName",
			platform: &azure.Platform{
				Region:                      "eastus",
				BaseDomainResourceGroupName: "",
				CloudName:                   azure.PublicCloud,
			},
			valid: false,
		},
		{
			name: "minimal",
			platform: &azure.Platform{
				Region:                      "eastus",
				BaseDomainResourceGroupName: "group",
				CloudName:                   azure.PublicCloud,
			},
			valid: true,
		},
		{
			name: "valid machine pool",
			platform: &azure.Platform{
				Region:                      "eastus",
				BaseDomainResourceGroupName: "group",
				DefaultMachinePlatform:      &azure.MachinePool{},
				CloudName:                   azure.PublicCloud,
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
				CloudName:                   azure.PublicCloud,
			},
			valid: true,
		},
		{
			name: "missing subnets",
			platform: &azure.Platform{
				Region:                   "eastus",
				NetworkResourceGroupName: "networkresourcegroup",
				VirtualNetwork:           "virtualnetwork",
				CloudName:                azure.PublicCloud,
			},
			valid: false,
		},
		{
			name: "subnets missing virtual network",
			platform: &azure.Platform{
				Region:                   "eastus",
				NetworkResourceGroupName: "networkresourcegroup",
				ComputeSubnet:            "computesubnet",
				CloudName:                azure.PublicCloud,
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
				CloudName:          azure.PublicCloud,
			},
			valid: false,
		},
		{
			name: "missing cloud name",
			platform: &azure.Platform{
				Region:                      "eastus",
				BaseDomainResourceGroupName: "group",
			},
			valid: false,
		},
		{
			name: "invalid cloud name",
			platform: &azure.Platform{
				Region:                      "eastus",
				BaseDomainResourceGroupName: "group",
				CloudName:                   azure.CloudEnvironment("AzureOtherCloud"),
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
