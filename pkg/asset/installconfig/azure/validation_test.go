package azure

import (
	"fmt"
	"net"
	"testing"

	aznetwork "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	azres "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	azsubs "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-06-01/subscriptions"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/golang/mock/gomock"
	"github.com/openshift/installer/pkg/asset/installconfig/azure/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type editFunctions []func(ic *types.InstallConfig)

var (
	validVirtualNetwork            = "valid-virtual-network"
	validNetworkResourceGroup      = "valid-network-resource-group"
	validRegion                    = "centralus"
	validRegionsList               = []string{"centralus", "westus", "australiacentral2"}
	resourcesCapableRegionsList    = []string{"centralus", "westus"}
	validComputeSubnet             = "valid-compute-subnet"
	validControlPlaneSubnet        = "valid-controlplane-subnet"
	validCIDR                      = "10.0.0.0/16"
	validComputeSubnetCIDR         = "10.0.0.0/24"
	validControlPlaneSubnetCIDR    = "10.0.32.0/24"
	validResourceGroupNamespace    = "Microsoft.Resources"
	validResourceGroupResourceType = "resourceGroups"
	validResourceSkuRegions        = "southeastasia"

	invalidateMachineCIDR = func(ic *types.InstallConfig) {
		_, newCidr, _ := net.ParseCIDR("192.168.111.0/24")
		ic.MachineNetwork = []types.MachineNetworkEntry{
			{CIDR: ipnet.IPNet{IPNet: *newCidr}},
		}
	}
	invalidResourceSkuRegion = "centralus"

	invalidateVirtualNetwork     = func(ic *types.InstallConfig) { ic.Azure.VirtualNetwork = "invalid-virtual-network" }
	invalidateComputeSubnet      = func(ic *types.InstallConfig) { ic.Azure.ComputeSubnet = "invalid-compute-subnet" }
	invalidateControlPlaneSubnet = func(ic *types.InstallConfig) { ic.Azure.ControlPlaneSubnet = "invalid-controlplane-subnet" }
	invalidateRegion             = func(ic *types.InstallConfig) { ic.Azure.Region = "neverland" }
	invalidateRegionCapabilities = func(ic *types.InstallConfig) { ic.Azure.Region = "australiacentral2" }
	invalidateRegionLetterCase   = func(ic *types.InstallConfig) { ic.Azure.Region = "Central US" }
	removeVirtualNetwork         = func(ic *types.InstallConfig) { ic.Azure.VirtualNetwork = "" }
	removeSubnets                = func(ic *types.InstallConfig) { ic.Azure.ComputeSubnet, ic.Azure.ControlPlaneSubnet = "", "" }

	virtualNetworkAPIResult = &aznetwork.VirtualNetwork{
		Name: &validVirtualNetwork,
	}
	computeSubnetAPIResult = &aznetwork.Subnet{
		Name: &validComputeSubnet,
		SubnetPropertiesFormat: &aznetwork.SubnetPropertiesFormat{
			AddressPrefix: &validComputeSubnetCIDR,
		},
	}
	controlPlaneSubnetAPIResult = &aznetwork.Subnet{
		Name: &validControlPlaneSubnet,
		SubnetPropertiesFormat: &aznetwork.SubnetPropertiesFormat{
			AddressPrefix: &validControlPlaneSubnetCIDR,
		},
	}
	locationsAPIResult = func() *[]azsubs.Location {
		r := []azsubs.Location{}
		for i := 0; i < len(validRegionsList); i++ {
			r = append(r, azsubs.Location{
				Name:        &validRegionsList[i],
				DisplayName: &validRegionsList[i],
			})
		}
		return &r
	}()
	resourcesProviderAPIResult = &azres.Provider{
		Namespace: &validResourceGroupNamespace,
		ResourceTypes: &[]azres.ProviderResourceType{
			{
				ResourceType: &validResourceGroupResourceType,
				Locations:    &resourcesCapableRegionsList,
			},
		},
	}
)

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Platform: types.Platform{
			Azure: &azure.Platform{
				Region:                   validRegion,
				NetworkResourceGroupName: validNetworkResourceGroup,
				VirtualNetwork:           validVirtualNetwork,
				ComputeSubnet:            validComputeSubnet,
				ControlPlaneSubnet:       validControlPlaneSubnet,
				DefaultMachinePlatform:   &azure.MachinePool{},
			},
		},
	}
}

func TestAzureInstallConfigValidation(t *testing.T) {
	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		{
			name:     "Valid virtual network & subnets",
			edits:    editFunctions{},
			errorMsg: "",
		},
		{
			name:     "Valid install config without virtual network & subnets",
			edits:    editFunctions{removeVirtualNetwork, removeSubnets},
			errorMsg: "",
		},
		{
			name:     "Invalid subnet range",
			edits:    editFunctions{invalidateMachineCIDR},
			errorMsg: "subnet .+ address prefix is outside of the specified machine networks",
		},
		{
			name:     "Invalid virtual network",
			edits:    editFunctions{invalidateVirtualNetwork},
			errorMsg: "invalid virtual network",
		},
		{
			name:     "Invalid compute subnet",
			edits:    editFunctions{invalidateComputeSubnet},
			errorMsg: "failed to retrieve compute subnet",
		},
		{
			name:     "Invalid control plane subnet",
			edits:    editFunctions{invalidateControlPlaneSubnet},
			errorMsg: "failed to retrieve control plane subnet",
		},
		{
			name:     "Invalid both subnets",
			edits:    editFunctions{invalidateControlPlaneSubnet, invalidateComputeSubnet},
			errorMsg: "failed to retrieve compute subnet",
		},
		{
			name:     "Invalid region",
			edits:    editFunctions{invalidateRegion},
			errorMsg: "region \"neverland\" is not valid or not available for this account$",
		},
		{
			name:     "Invalid region uncapable",
			edits:    editFunctions{invalidateRegionCapabilities},
			errorMsg: "region \"australiacentral2\" does not support resource creation$",
		},
		{
			name:     "Invalid region letter case",
			edits:    editFunctions{invalidateRegionLetterCase},
			errorMsg: "region \"Central US\" is not valid or not available for this account, did you mean \"centralus\"\\?$",
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	azureClient := mock.NewMockAPI(mockCtrl)

	// VirtualNetwork
	azureClient.EXPECT().GetVirtualNetwork(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork).Return(virtualNetworkAPIResult, nil).AnyTimes()
	azureClient.EXPECT().GetVirtualNetwork(gomock.Any(), gomock.Not(validNetworkResourceGroup), gomock.Not(validVirtualNetwork)).Return(&aznetwork.VirtualNetwork{}, fmt.Errorf("invalid network resource group")).AnyTimes()
	azureClient.EXPECT().GetVirtualNetwork(gomock.Any(), validNetworkResourceGroup, gomock.Not(validVirtualNetwork)).Return(&aznetwork.VirtualNetwork{}, fmt.Errorf("invalid virtual network")).AnyTimes()

	// ComputeSubnet
	azureClient.EXPECT().GetComputeSubnet(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork, validComputeSubnet).Return(computeSubnetAPIResult, nil).AnyTimes()
	azureClient.EXPECT().GetComputeSubnet(gomock.Any(), gomock.Not(validNetworkResourceGroup), validVirtualNetwork, validComputeSubnet).Return(&aznetwork.Subnet{}, fmt.Errorf("invalid network resource group")).AnyTimes()
	azureClient.EXPECT().GetComputeSubnet(gomock.Any(), validNetworkResourceGroup, gomock.Not(validVirtualNetwork), validComputeSubnet).Return(&aznetwork.Subnet{}, fmt.Errorf("invalid virtual network")).AnyTimes()
	azureClient.EXPECT().GetComputeSubnet(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork, gomock.Not(validComputeSubnet)).Return(&aznetwork.Subnet{}, fmt.Errorf("invalid compute subnet")).AnyTimes()

	// ControlPlaneSubnet
	azureClient.EXPECT().GetControlPlaneSubnet(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork, validControlPlaneSubnet).Return(controlPlaneSubnetAPIResult, nil).AnyTimes()
	azureClient.EXPECT().GetControlPlaneSubnet(gomock.Any(), gomock.Not(validNetworkResourceGroup), validVirtualNetwork, validControlPlaneSubnet).Return(&aznetwork.Subnet{}, fmt.Errorf("invalid network resource group")).AnyTimes()
	azureClient.EXPECT().GetControlPlaneSubnet(gomock.Any(), validNetworkResourceGroup, gomock.Not(validVirtualNetwork), validControlPlaneSubnet).Return(&aznetwork.Subnet{}, fmt.Errorf("invalid virtual network")).AnyTimes()
	azureClient.EXPECT().GetControlPlaneSubnet(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork, gomock.Not(validControlPlaneSubnet)).Return(&aznetwork.Subnet{}, fmt.Errorf("invalid control plane subnet")).AnyTimes()

	// Location
	azureClient.EXPECT().ListLocations(gomock.Any()).Return(locationsAPIResult, nil).AnyTimes()

	// ResourceProvider
	azureClient.EXPECT().GetResourcesProvider(gomock.Any(), validResourceGroupNamespace).Return(resourcesProviderAPIResult, nil).AnyTimes()

	//Resource SKUs
	azureClient.EXPECT().GetDiskSkus(gomock.Any(), validResourceSkuRegions).Return(nil, fmt.Errorf("invalid disk type")).AnyTimes()
	azureClient.EXPECT().GetDiskSkus(gomock.Any(), invalidResourceSkuRegion).Return(nil, fmt.Errorf("invalid region")).AnyTimes()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			aggregatedErrors := Validate(azureClient, editedInstallConfig)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}

var validGroupResult = &azres.Group{
	ID:       to.StringPtr("valid-resource-group"),
	Location: to.StringPtr("centralus"),
}

var invalidGroupOutsideRegionResult = &azres.Group{
	ID:       to.StringPtr("invalid-resource-group-useast2"),
	Location: to.StringPtr("useast2"),
}

var validGroupWithTagsResult = &azres.Group{
	ID:       to.StringPtr("valid-resource-group-tags"),
	Location: to.StringPtr("centralus"),
	Tags: map[string]*string{
		"key": to.StringPtr("value"),
	},
}

var validGroupWithConflictinsTagsResult = &azres.Group{
	ID:       to.StringPtr("valid-resource-group-conf-tags"),
	Location: to.StringPtr("centralus"),
	Tags: map[string]*string{
		"kubernetes.io_cluster.test-cluster-12345": to.StringPtr("owned"),
	},
}

func Test_validateResourceGroup(t *testing.T) {
	cases := []struct {
		groupName string
		err       string
	}{{
		groupName: "non-existent-group",
		err:       `^\Qplatform.azure.resourceGroupName: Internal error: failed to get resource group: resource group /resourceGroups/non-existent-group was not found\E$`,
	}, {
		groupName: "valid-resource-group",
	}, {
		groupName: "invalid-resource-group-useast2",
		err:       `^\Qplatform.azure.resourceGroupName: Invalid value: "invalid-resource-group-useast2": expected to in region centralus, but found it to be in useast2\E$`,
	}, {
		groupName: "valid-resource-group-tags",
	}, {
		groupName: "valid-resource-group-conf-tags",
		err:       `^\Qplatform.azure.resourceGroupName: Invalid value: "valid-resource-group-conf-tags": resource group has conflicting tags kubernetes.io_cluster.test-cluster-12345\E$`,
	}, {
		groupName: "valid-resource-group-with-resources",
		err:       `^\Qplatform.azure.resourceGroupName: Invalid value: "valid-resource-group-with-resources": resource group must be empty but it has 3 resources like id1, id2 ...\E$`,
	}}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	azureClient := mock.NewMockAPI(mockCtrl)
	azureClient.EXPECT().GetGroup(gomock.Any(), "non-existent-group").Return(nil, fmt.Errorf("resource group /resourceGroups/non-existent-group was not found")).AnyTimes()
	azureClient.EXPECT().GetGroup(gomock.Any(), "valid-resource-group").Return(validGroupResult, nil).AnyTimes()
	azureClient.EXPECT().GetGroup(gomock.Any(), "invalid-resource-group-useast2").Return(invalidGroupOutsideRegionResult, nil).AnyTimes()
	azureClient.EXPECT().GetGroup(gomock.Any(), "valid-resource-group-with-resources").Return(validGroupResult, nil).AnyTimes()
	azureClient.EXPECT().GetGroup(gomock.Any(), "valid-resource-group-tags").Return(validGroupWithTagsResult, nil).AnyTimes()
	azureClient.EXPECT().GetGroup(gomock.Any(), "valid-resource-group-conf-tags").Return(validGroupWithConflictinsTagsResult, nil).AnyTimes()
	azureClient.EXPECT().ListResourceIDsByGroup(gomock.Any(), gomock.Not("valid-resource-group-with-resources")).Return(nil, nil).AnyTimes()
	azureClient.EXPECT().ListResourceIDsByGroup(gomock.Any(), "valid-resource-group-with-resources").Return([]string{"id1", "id2", "id3"}, nil).AnyTimes()

	for _, test := range cases {
		t.Run("", func(t *testing.T) {
			err := validateResourceGroup(azureClient, field.NewPath("platform").Child("azure"), &azure.Platform{ResourceGroupName: test.groupName, Region: "centralus"})
			if test.err != "" {
				assert.Regexp(t, test.err, err.ToAggregate())
			} else {
				assert.NoError(t, err.ToAggregate())
			}
		})
	}
}
