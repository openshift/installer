package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	aztypes "github.com/openshift/installer/pkg/types/azure"
)

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

// strPtr returns a pointer to the given string.
func strPtr(s string) *string {
	return &s
}
