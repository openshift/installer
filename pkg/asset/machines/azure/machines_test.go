package azure

import (
	"encoding/json"
	"testing"

	azureenv "github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	machineapi "github.com/openshift/api/machine/v1beta1"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/types"
	aztypes "github.com/openshift/installer/pkg/types/azure"
)

func TestGenerateSecurityProfile(t *testing.T) {
	enabled := "Enabled"
	disabled := "Disabled"
	encryptionAtHost := true

	testCases := []struct {
		name     string
		mpool    *aztypes.MachinePool
		expected *machineapi.SecurityProfile
	}{
		{
			name:     "no security features",
			mpool:    &aztypes.MachinePool{},
			expected: nil,
		},
		{
			name: "encryption at host only",
			mpool: &aztypes.MachinePool{
				EncryptionAtHost: true,
			},
			expected: &machineapi.SecurityProfile{
				EncryptionAtHost: &encryptionAtHost,
			},
		},
		{
			name: "trusted launch",
			mpool: &aztypes.MachinePool{
				Settings: &aztypes.SecuritySettings{
					SecurityType: aztypes.SecurityTypesTrustedLaunch,
					TrustedLaunch: &aztypes.TrustedLaunch{
						UEFISettings: &aztypes.UEFISettings{
							SecureBoot:                       &enabled,
							VirtualizedTrustedPlatformModule: &disabled,
						},
					},
				},
			},
			expected: &machineapi.SecurityProfile{
				Settings: machineapi.SecuritySettings{
					SecurityType: machineapi.SecurityTypesTrustedLaunch,
					TrustedLaunch: &machineapi.TrustedLaunch{
						UEFISettings: machineapi.UEFISettings{
							SecureBoot:                       machineapi.SecureBootPolicyEnabled,
							VirtualizedTrustedPlatformModule: machineapi.VirtualizedTrustedPlatformModulePolicyDisabled,
						},
					},
				},
			},
		},
		{
			name: "confidential VM",
			mpool: &aztypes.MachinePool{
				Settings: &aztypes.SecuritySettings{
					SecurityType: aztypes.SecurityTypesConfidentialVM,
					ConfidentialVM: &aztypes.ConfidentialVM{
						UEFISettings: &aztypes.UEFISettings{
							SecureBoot:                       &disabled,
							VirtualizedTrustedPlatformModule: &enabled,
						},
					},
				},
			},
			expected: &machineapi.SecurityProfile{
				Settings: machineapi.SecuritySettings{
					SecurityType: machineapi.SecurityTypesConfidentialVM,
					ConfidentialVM: &machineapi.ConfidentialVM{
						UEFISettings: machineapi.UEFISettings{
							SecureBoot:                       machineapi.SecureBootPolicyDisabled,
							VirtualizedTrustedPlatformModule: machineapi.VirtualizedTrustedPlatformModulePolicyEnabled,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, generateSecurityProfile(tc.mpool))
		})
	}
}

func TestEmptySecurityProfileIsOmittedFromProvider(t *testing.T) {
	azIdx := 0
	spec, err := provider(
		&aztypes.Platform{Region: "eastus"},
		&aztypes.MachinePool{
			InstanceType: "Standard_D2s_v3",
			Zones:        []string{"1"},
			Identity:     &aztypes.VMIdentity{},
		},
		"",
		"user-data",
		"test-cluster",
		"master",
		&azIdx,
		map[string]string{"HyperVGenerations": "V1,V2"},
		&icazure.Session{
			Credentials: icazure.Credentials{SubscriptionID: "test-subscription"},
			Environment: azureenv.PublicCloud,
		},
		"test-network-resource-group",
		"test-vnet",
		"test-subnet",
	)
	assert.NoError(t, err)
	assert.Nil(t, spec.SecurityProfile)

	serialized, err := json.Marshal(spec)
	assert.NoError(t, err)

	var serializedSpec map[string]json.RawMessage
	err = json.Unmarshal(serialized, &serializedSpec)
	assert.NoError(t, err)
	_, securityProfilePresent := serializedSpec["securityProfile"]
	assert.False(t, securityProfilePresent)
}

func TestGenerateMachinesSecurityProfile(t *testing.T) {
	enabled := "Enabled"
	disabled := "Disabled"

	testCases := []struct {
		name     string
		settings *aztypes.MachinePool
		expected *capz.SecurityProfile
	}{
		{
			name:     "no security features",
			settings: &aztypes.MachinePool{},
			expected: nil,
		},
		{
			name: "encryption at host only",
			settings: &aztypes.MachinePool{
				EncryptionAtHost: true,
			},
			expected: &capz.SecurityProfile{
				EncryptionAtHost: ptr.To(true),
			},
		},
		{
			name: "trusted launch",
			settings: &aztypes.MachinePool{
				Settings: &aztypes.SecuritySettings{
					SecurityType: aztypes.SecurityTypesTrustedLaunch,
					TrustedLaunch: &aztypes.TrustedLaunch{
						UEFISettings: &aztypes.UEFISettings{
							SecureBoot:                       &enabled,
							VirtualizedTrustedPlatformModule: &disabled,
						},
					},
				},
			},
			expected: &capz.SecurityProfile{
				SecurityType: capz.SecurityTypesTrustedLaunch,
				UefiSettings: &capz.UefiSettings{
					SecureBootEnabled: ptr.To(true),
					VTpmEnabled:       ptr.To(false),
				},
			},
		},
		{
			name: "confidential VM",
			settings: &aztypes.MachinePool{
				Settings: &aztypes.SecuritySettings{
					SecurityType: aztypes.SecurityTypesConfidentialVM,
					ConfidentialVM: &aztypes.ConfidentialVM{
						UEFISettings: &aztypes.UEFISettings{
							SecureBoot:                       &disabled,
							VirtualizedTrustedPlatformModule: &enabled,
						},
					},
				},
			},
			expected: &capz.SecurityProfile{
				SecurityType: capz.SecurityTypesConfidentialVM,
				UefiSettings: &capz.UefiSettings{
					SecureBootEnabled: ptr.To(false),
					VTpmEnabled:       ptr.To(true),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			machinePool := &types.MachinePool{
				Name:     "master",
				Replicas: ptr.To[int64](1),
				Platform: types.MachinePoolPlatform{
					Azure: tc.settings,
				},
			}
			machinePool.Platform.Azure.Identity = &aztypes.VMIdentity{}

			files, err := GenerateMachines(
				"test-cluster",
				"test-resource-group",
				"test-subscription",
				&icazure.Session{Environment: azureenv.PublicCloud},
				&MachineInput{
					Environment: aztypes.PublicCloud,
					Platform:    &aztypes.Platform{},
					Pool:        machinePool,
				},
			)
			assert.NoError(t, err)

			azureMachineCount := 0
			for _, file := range files {
				azureMachine, ok := file.Object.(*capz.AzureMachine)
				if !ok {
					continue
				}
				azureMachineCount++
				assert.Equal(t, tc.expected, azureMachine.Spec.SecurityProfile)

				if tc.expected == nil {
					serialized, err := json.Marshal(azureMachine)
					assert.NoError(t, err)

					var serializedMachine struct {
						Spec map[string]json.RawMessage `json:"spec"`
					}
					err = json.Unmarshal(serialized, &serializedMachine)
					assert.NoError(t, err)
					_, securityProfilePresent := serializedMachine.Spec["securityProfile"]
					assert.False(t, securityProfilePresent)
				}
			}
			assert.Equal(t, 2, azureMachineCount)
		})
	}
}

func TestCapzImage(t *testing.T) {
	testCases := []struct {
		name           string
		osImage        aztypes.OSImage
		azEnv          aztypes.CloudEnvironment
		confidentialVM bool
		gen            string
		rg             string
		sub            string
		infraID        string
		rhcosImg       string
		expectedImage  *capz.Image
	}{
		{
			name: "marketplace image from osImage",
			osImage: aztypes.OSImage{
				Publisher: "RedHat",
				Offer:     "RHEL",
				SKU:       "8-lvm-gen2",
				Version:   "latest",
				Plan:      aztypes.ImageWithPurchasePlan,
			},
			azEnv:          aztypes.PublicCloud,
			confidentialVM: false,
			gen:            "V2",
			rg:             "test-rg",
			sub:            "test-sub",
			infraID:        "test-cluster",
			rhcosImg:       "",
			expectedImage: &capz.Image{
				Marketplace: &capz.AzureMarketplaceImage{
					ImagePlan: capz.ImagePlan{
						Publisher: "RedHat",
						Offer:     "RHEL",
						SKU:       "8-lvm-gen2",
					},
					Version:         "latest",
					ThirdPartyImage: true,
				},
			},
		},
		{
			name: "marketplace image no purchase plan from osImage",
			osImage: aztypes.OSImage{
				Publisher: "RedHat",
				Offer:     "RHEL",
				SKU:       "8-lvm-gen2",
				Version:   "latest",
				Plan:      aztypes.ImageNoPurchasePlan,
			},
			azEnv:          aztypes.PublicCloud,
			confidentialVM: false,
			gen:            "V2",
			rg:             "test-rg",
			sub:            "test-sub",
			infraID:        "test-cluster",
			rhcosImg:       "",
			expectedImage: &capz.Image{
				Marketplace: &capz.AzureMarketplaceImage{
					ImagePlan: capz.ImagePlan{
						Publisher: "RedHat",
						Offer:     "RHEL",
						SKU:       "8-lvm-gen2",
					},
					Version:         "latest",
					ThirdPartyImage: false,
				},
			},
		},
		{
			name:           "azure stack cloud managed image",
			osImage:        aztypes.OSImage{},
			azEnv:          aztypes.StackCloud,
			confidentialVM: false,
			gen:            "V1",
			rg:             "test-rg",
			sub:            "test-sub-123",
			infraID:        "test-cluster-abc",
			rhcosImg:       "",
			expectedImage: &capz.Image{
				ID: strPtr("/subscriptions/test-sub-123/resourceGroups/test-rg/providers/Microsoft.Compute/images/test-cluster-abc"),
			},
		},
		{
			name:           "marketplace URN format in rhcosImg",
			osImage:        aztypes.OSImage{},
			azEnv:          aztypes.PublicCloud,
			confidentialVM: false,
			gen:            "V2",
			rg:             "test-rg",
			sub:            "test-sub",
			infraID:        "test-cluster",
			rhcosImg:       "RedHat:RHEL:8-lvm:latest",
			expectedImage: &capz.Image{
				Marketplace: &capz.AzureMarketplaceImage{
					ImagePlan: capz.ImagePlan{
						Publisher: "RedHat",
						Offer:     "RHEL",
						SKU:       "8-lvm",
					},
					Version:         "latest",
					ThirdPartyImage: false,
				},
			},
		},
		{
			name:           "explicit subscription resource path for managed image",
			osImage:        aztypes.OSImage{},
			azEnv:          aztypes.PublicCloud,
			confidentialVM: false,
			gen:            "V2",
			rg:             "test-rg",
			sub:            "test-sub",
			infraID:        "test-cluster",
			rhcosImg:       "/subscriptions/custom-sub/resourceGroups/custom-rg/providers/Microsoft.Compute/images/custom-image",
			expectedImage: &capz.Image{
				ID: strPtr("/subscriptions/custom-sub/resourceGroups/custom-rg/providers/Microsoft.Compute/images/custom-image"),
			},
		},
		{
			name:           "explicit subscription resource path for gallery image",
			osImage:        aztypes.OSImage{},
			azEnv:          aztypes.PublicCloud,
			confidentialVM: false,
			gen:            "V2",
			rg:             "test-rg",
			sub:            "test-sub",
			infraID:        "test-cluster",
			rhcosImg:       "/subscriptions/custom-sub/resourceGroups/custom-rg/providers/Microsoft.Compute/galleries/gallery_name/images/image_def",
			expectedImage: &capz.Image{
				ID: strPtr("/subscriptions/custom-sub/resourceGroups/custom-rg/providers/Microsoft.Compute/galleries/gallery_name/images/image_def"),
			},
		},
		{
			name:           "empty rhcosImg for non-confidential VM",
			osImage:        aztypes.OSImage{},
			azEnv:          aztypes.PublicCloud,
			confidentialVM: false,
			gen:            "V2",
			rg:             "test-rg",
			sub:            "test-sub",
			infraID:        "test-cluster",
			rhcosImg:       "",
			expectedImage:  &capz.Image{},
		},
		{
			name:           "installer-created gallery image for confidential VM gen V2",
			osImage:        aztypes.OSImage{},
			azEnv:          aztypes.PublicCloud,
			confidentialVM: true,
			gen:            "V2",
			rg:             "test-rg",
			sub:            "test-sub-456",
			infraID:        "test-cluster-xyz",
			rhcosImg:       "",
			expectedImage: &capz.Image{
				ID: strPtr("/subscriptions/test-sub-456/resourceGroups/test-rg/providers/Microsoft.Compute/galleries/gallery_test_cluster_xyz/images/test-cluster-xyz-gen2"),
			},
		},
		{
			name:           "installer-created gallery image for confidential VM gen V1",
			osImage:        aztypes.OSImage{},
			azEnv:          aztypes.PublicCloud,
			confidentialVM: true,
			gen:            "V1",
			rg:             "test-rg",
			sub:            "test-sub-456",
			infraID:        "test-cluster-xyz",
			rhcosImg:       "",
			expectedImage: &capz.Image{
				ID: strPtr("/subscriptions/test-sub-456/resourceGroups/test-rg/providers/Microsoft.Compute/galleries/gallery_test_cluster_xyz/images/test-cluster-xyz"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := capzImage(tc.osImage, tc.azEnv, tc.confidentialVM, tc.gen, tc.rg, tc.sub, tc.infraID, tc.rhcosImg)
			assert.Equal(t, tc.expectedImage, result)
		})
	}
}

func TestTrimSubscriptionPrefix(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid full subscription path for managed image",
			input:    "/subscriptions/12345678-1234-1234-1234-123456789abc/resourceGroups/my-rg/providers/Microsoft.Compute/images/my-image",
			expected: "/resourceGroups/my-rg/providers/Microsoft.Compute/images/my-image",
		},
		{
			name:     "valid full subscription path for gallery image",
			input:    "/subscriptions/12345678-1234-1234-1234-123456789abc/resourceGroups/my-rg/providers/Microsoft.Compute/galleries/my_gallery/images/my_image_def",
			expected: "/resourceGroups/my-rg/providers/Microsoft.Compute/galleries/my_gallery/images/my_image_def",
		},
		{
			name:     "already trimmed path",
			input:    "/resourceGroups/my-rg/providers/Microsoft.Compute/images/my-image",
			expected: "/resourceGroups/my-rg/providers/Microsoft.Compute/images/my-image",
		},
		{
			name:     "invalid path without resourceGroups",
			input:    "/subscriptions/12345678-1234-1234-1234-123456789abc/providers/Microsoft.Compute/images/my-image",
			expected: "/subscriptions/12345678-1234-1234-1234-123456789abc/providers/Microsoft.Compute/images/my-image",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "marketplace URN not an ID",
			input:    "RedHat:RHEL:8-lvm:latest",
			expected: "RedHat:RHEL:8-lvm:latest",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := trimSubscriptionPrefix(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMapiImage(t *testing.T) {
	testCases := []struct {
		name           string
		osImage        aztypes.OSImage
		azEnv          aztypes.CloudEnvironment
		confidentialVM bool
		gen            string
		rg             string
		sub            string
		infraID        string
		rhcosImg       string
		expected       interface{}
	}{
		{
			name: "marketplace image with purchase plan",
			osImage: aztypes.OSImage{
				Publisher: "RedHat",
				Offer:     "RHEL",
				SKU:       "8-lvm-gen2",
				Version:   "latest",
				Plan:      aztypes.ImageWithPurchasePlan,
			},
			azEnv:          aztypes.PublicCloud,
			confidentialVM: false,
			gen:            "V2",
			rg:             "test-rg",
			sub:            "test-sub",
			infraID:        "test-cluster",
			rhcosImg:       "",
			expected: struct {
				publisher string
				offer     string
				sku       string
				version   string
				imageType string
			}{
				publisher: "RedHat",
				offer:     "RHEL",
				sku:       "8-lvm-gen2",
				version:   "latest",
				imageType: "MarketplaceWithPlan",
			},
		},
		{
			name: "marketplace image no purchase plan",
			osImage: aztypes.OSImage{
				Publisher: "RedHat",
				Offer:     "RHEL",
				SKU:       "8-lvm-gen2",
				Version:   "latest",
				Plan:      aztypes.ImageNoPurchasePlan,
			},
			azEnv:          aztypes.PublicCloud,
			confidentialVM: false,
			gen:            "V2",
			rg:             "test-rg",
			sub:            "test-sub",
			infraID:        "test-cluster",
			rhcosImg:       "",
			expected: struct {
				publisher string
				offer     string
				sku       string
				version   string
				imageType string
			}{
				publisher: "RedHat",
				offer:     "RHEL",
				sku:       "8-lvm-gen2",
				version:   "latest",
				imageType: "MarketplaceNoPlan",
			},
		},
		{
			name:           "explicit subscription resource path",
			osImage:        aztypes.OSImage{},
			azEnv:          aztypes.PublicCloud,
			confidentialVM: false,
			gen:            "V2",
			rg:             "test-rg",
			sub:            "test-sub",
			infraID:        "test-cluster",
			rhcosImg:       "/subscriptions/custom-sub/resourceGroups/custom-rg/providers/Microsoft.Compute/images/custom-image",
			expected: struct {
				resourceID string
			}{
				resourceID: "/resourceGroups/custom-rg/providers/Microsoft.Compute/images/custom-image",
			},
		},
		{
			name:           "gallery image path",
			osImage:        aztypes.OSImage{},
			azEnv:          aztypes.PublicCloud,
			confidentialVM: false,
			gen:            "V2",
			rg:             "test-rg",
			sub:            "test-sub",
			infraID:        "test-cluster",
			rhcosImg:       "/subscriptions/custom-sub/resourceGroups/custom-rg/providers/Microsoft.Compute/galleries/gallery_name/images/image_def",
			expected: struct {
				resourceID string
			}{
				resourceID: "/resourceGroups/custom-rg/providers/Microsoft.Compute/galleries/gallery_name/images/image_def",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := mapiImage(tc.osImage, tc.azEnv, tc.confidentialVM, tc.gen, tc.rg, tc.sub, tc.infraID, tc.rhcosImg)

			switch expected := tc.expected.(type) {
			case struct {
				publisher string
				offer     string
				sku       string
				version   string
				imageType string
			}:
				assert.Equal(t, expected.publisher, result.Publisher)
				assert.Equal(t, expected.offer, result.Offer)
				assert.Equal(t, expected.sku, result.SKU)
				assert.Equal(t, expected.version, result.Version)
				assert.Equal(t, expected.imageType, string(result.Type))
			case struct{ resourceID string }:
				assert.Equal(t, expected.resourceID, result.ResourceID)
			}
		})
	}
}

func TestProviderDataDiskEncryptionSet(t *testing.T) {
	validDESID := "/subscriptions/sub-123/resourceGroups/rg-456/providers/Microsoft.Compute/diskEncryptionSets/des-789"
	lun := int32(1)

	basePlatform := &aztypes.Platform{
		Region: "eastus",
	}
	baseSession := &icazure.Session{
		Credentials: icazure.Credentials{
			SubscriptionID: "test-sub",
		},
		Environment: azureenv.Environment{
			StorageEndpointSuffix: "core.windows.net",
		},
	}
	baseMpool := func() *aztypes.MachinePool {
		return &aztypes.MachinePool{
			InstanceType: "Standard_D2s_v3",
			OSDisk: aztypes.OSDisk{
				DiskSizeGB: 128,
				DiskType:   "Premium_LRS",
			},
			Zones:    []string{"1"},
			Identity: &aztypes.VMIdentity{},
		}
	}
	capabilities := map[string]string{
		"HyperVGenerations": "V1,V2",
	}
	azIdx := 0

	testCases := []struct {
		name          string
		dataDisks     []capz.DataDisk
		expectError   bool
		errorContains string
		validate      func(t *testing.T, spec *machineapi.AzureMachineProviderSpec)
	}{
		{
			name: "data disk with valid DiskEncryptionSet",
			dataDisks: []capz.DataDisk{
				{
					NameSuffix:  "disk1",
					DiskSizeGB:  256,
					Lun:         &lun,
					CachingType: "ReadOnly",
					ManagedDisk: &capz.ManagedDiskParameters{
						StorageAccountType: "Premium_LRS",
						DiskEncryptionSet: &capz.DiskEncryptionSetParameters{
							ID: validDESID,
						},
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, spec *machineapi.AzureMachineProviderSpec) {
				t.Helper()
				assert.Len(t, spec.DataDisks, 1)
				assert.NotNil(t, spec.DataDisks[0].ManagedDisk.DiskEncryptionSet)
				assert.Equal(t, validDESID, spec.DataDisks[0].ManagedDisk.DiskEncryptionSet.ID)
			},
		},
		{
			name: "data disk with DiskEncryptionSet but empty ID",
			dataDisks: []capz.DataDisk{
				{
					NameSuffix:  "disk-empty-id",
					DiskSizeGB:  256,
					CachingType: "ReadOnly",
					ManagedDisk: &capz.ManagedDiskParameters{
						StorageAccountType: "Premium_LRS",
						DiskEncryptionSet: &capz.DiskEncryptionSetParameters{
							ID: "",
						},
					},
				},
			},
			expectError:   true,
			errorContains: "data disk disk-empty-id has invalid disk encryption set: empty ID",
		},
		{
			name: "data disk without DiskEncryptionSet",
			dataDisks: []capz.DataDisk{
				{
					NameSuffix:  "disk-no-des",
					DiskSizeGB:  256,
					CachingType: "ReadOnly",
					ManagedDisk: &capz.ManagedDiskParameters{
						StorageAccountType: "Premium_LRS",
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, spec *machineapi.AzureMachineProviderSpec) {
				t.Helper()
				assert.Len(t, spec.DataDisks, 1)
				assert.Nil(t, spec.DataDisks[0].ManagedDisk.DiskEncryptionSet)
			},
		},
		{
			name: "data disk with DiskEncryptionSet but nil SecurityProfile - no panic",
			dataDisks: []capz.DataDisk{
				{
					NameSuffix:  "disk-nil-secprofile",
					DiskSizeGB:  256,
					CachingType: "ReadOnly",
					ManagedDisk: &capz.ManagedDiskParameters{
						StorageAccountType: "Premium_LRS",
						DiskEncryptionSet: &capz.DiskEncryptionSetParameters{
							ID: validDESID,
						},
						SecurityProfile: nil,
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, spec *machineapi.AzureMachineProviderSpec) {
				t.Helper()
				assert.Len(t, spec.DataDisks, 1)
				assert.NotNil(t, spec.DataDisks[0].ManagedDisk.DiskEncryptionSet)
				assert.Equal(t, validDESID, spec.DataDisks[0].ManagedDisk.DiskEncryptionSet.ID)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mpool := baseMpool()
			mpool.DataDisks = tc.dataDisks

			spec, err := provider(
				basePlatform,
				mpool,
				"",             // osImage
				"user-data",    // userDataSecret
				"test-cluster", // clusterID
				"worker",       // role
				&azIdx,         // azIdx
				capabilities,
				baseSession,
				"test-nrg",    // networkResourceGroup
				"test-vnet",   // virtualNetwork
				"test-subnet", // subnet
			)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
				assert.Nil(t, spec)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, spec)
				if tc.validate != nil {
					tc.validate(t, spec)
				}
			}
		})
	}
}

// strPtr returns a pointer to the given string.
func strPtr(s string) *string {
	return &s
}
