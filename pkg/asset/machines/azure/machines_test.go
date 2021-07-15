package azure

import (
	"fmt"
	"testing"

	"github.com/openshift/installer/pkg/types/azure"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	azureprovider "sigs.k8s.io/cluster-api-provider-azure/pkg/apis/azureprovider/v1beta1"
)

func TestProvider(t *testing.T) {
	clusterID := "00000000-0000-0000-0000-000000000000"
	rg := fmt.Sprintf("%s-rg", clusterID)
	zone := "1"

	expectedMachineProviderSpec := func() *azureprovider.AzureMachineProviderSpec {
		return &azureprovider.AzureMachineProviderSpec{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "azureproviderconfig.openshift.io/v1beta1",
				Kind:       "AzureMachineProviderSpec",
			},
			UserDataSecret:    &corev1.SecretReference{Name: ""},
			CredentialsSecret: &corev1.SecretReference{Name: cloudsSecret, Namespace: cloudsSecretNamespace},
			Location:          "fake-region",
			VMSize:            "",
			Image: azureprovider.Image{
				ResourceID: fmt.Sprintf("/resourceGroups/%s/providers/Microsoft.Compute/images/%s", rg, clusterID),
			},
			OSDisk: azureprovider.OSDisk{
				OSType:     "Linux",
				DiskSizeGB: 123,
				ManagedDisk: azureprovider.ManagedDiskParameters{
					StorageAccountType: "Premium_LRS",
				},
			},
			Zone:                 &zone,
			Subnet:               fmt.Sprintf("%s-master-subnet", clusterID),
			ManagedIdentity:      fmt.Sprintf("%s-identity", clusterID),
			Vnet:                 fmt.Sprintf("%s-vnet", clusterID),
			ResourceGroup:        rg,
			NetworkResourceGroup: rg,
			PublicLoadBalancer:   clusterID,
		}
	}

	tests := []struct {
		name   string
		mocks  func(platform *azure.Platform, mpool *azure.MachinePool)
		expect func(spec *azureprovider.AzureMachineProviderSpec)
	}{
		{
			name: "no customisations",
		},
		{
			name: "DiskType specified",
			mocks: func(platform *azure.Platform, mpool *azure.MachinePool) {
				mpool.OSDisk.DiskType = "Standard_LRS"
			},
			expect: func(spec *azureprovider.AzureMachineProviderSpec) {
				spec.OSDisk.ManagedDisk.StorageAccountType = "Standard_LRS"
			},
		},
		{
			name: "OutboundType set to UserDefinedRouting",
			mocks: func(platform *azure.Platform, mpool *azure.MachinePool) {
				platform.OutboundType = azure.UserDefinedRoutingOutboundType
			},
			expect: func(spec *azureprovider.AzureMachineProviderSpec) {
				spec.PublicLoadBalancer = ""
			},
		},
		{
			name: "Image from ResourceID",
			mocks: func(platform *azure.Platform, mpool *azure.MachinePool) {
				platform.Image = &azure.Image{
					ResourceID: "/resourceGroups/fake-rg/providers/Microsoft.Compute/images/fake-image",
				}
			},
			expect: func(spec *azureprovider.AzureMachineProviderSpec) {
				spec.Image = azureprovider.Image{
					ResourceID: "/resourceGroups/fake-rg/providers/Microsoft.Compute/images/fake-image",
				}
			},
		},
		{
			name: "Image from marketplace",
			mocks: func(platform *azure.Platform, mpool *azure.MachinePool) {
				platform.Image = &azure.Image{
					Publisher: "fake-publisher",
					Offer:     "fake-offer",
					SKU:       "fake-sku",
					Version:   "fake-version",
				}
			},
			expect: func(spec *azureprovider.AzureMachineProviderSpec) {
				spec.Image = azureprovider.Image{
					Publisher: "fake-publisher",
					Offer:     "fake-offer",
					SKU:       "fake-sku",
					Version:   "fake-version",
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			platform := &azure.Platform{
				Region: "fake-region",
			}
			mpool := &azure.MachinePool{
				Zones: []string{"0", "1", "2"},
				OSDisk: azure.OSDisk{
					DiskSizeGB: 123,
				},
			}
			azIdx := 1

			if tt.mocks != nil {
				tt.mocks(platform, mpool)
			}

			spec, err := provider(platform, mpool, "", "", clusterID, "master", &azIdx)
			if err != nil {
				t.Error(err)
			}

			expectedSpec := expectedMachineProviderSpec()
			if tt.expect != nil {
				tt.expect(expectedSpec)
			}

			assert.Equal(t, expectedSpec, spec)
		})
	}
}

func TestGetNetworkInfo(t *testing.T) {
	clusterID := "00000000-0000-0000-0000-000000000000"

	for _, role := range []string{"master", "worker"} {
		tests := []struct {
			name                       string
			platform                   *azure.Platform
			expectNetworkResourceGroup string
			expectVirtualNetwork       string
			expectSubnet               string
		}{
			{
				name:                       "no VirtualNetwork customisation",
				platform:                   &azure.Platform{},
				expectNetworkResourceGroup: fmt.Sprintf("%s-rg", clusterID),
				expectVirtualNetwork:       fmt.Sprintf("%s-vnet", clusterID),
				expectSubnet:               fmt.Sprintf("%s-%s-subnet", clusterID, role),
			},
			{
				name: "VirtualNetwork customisation",
				platform: &azure.Platform{
					NetworkResourceGroupName: "fake-vnet-rg",
					VirtualNetwork:           "fake-vnet",
					ControlPlaneSubnet:       "fake-master-subnet",
					ComputeSubnet:            "fake-worker-subnet",
				},
				expectNetworkResourceGroup: "fake-vnet-rg",
				expectVirtualNetwork:       "fake-vnet",
				expectSubnet:               fmt.Sprintf("fake-%s-subnet", role),
			},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s: %s", role, tt.name), func(t *testing.T) {
				networkResourceGroup, virtualNetwork, subnet, err := getNetworkInfo(tt.platform, clusterID, role)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, tt.expectNetworkResourceGroup, networkResourceGroup)
				assert.Equal(t, tt.expectVirtualNetwork, virtualNetwork)
				assert.Equal(t, tt.expectSubnet, subnet)
			})
		}
	}
}
