package rhcos

import "fmt"

// Extensions is data specific to Red Hat Enterprise Linux CoreOS
type Extensions struct {
	AzureDisk   *AzureDisk   `json:"azure-disk,omitempty"`
	Marketplace *Marketplace `json:"marketplace,omitempty"`
}

// AzureDisk represents an Azure disk image that can be imported
// into an image gallery or otherwise replicated, and then used
// as a boot source for virtual machines.
type AzureDisk struct {
	// Release is the source release version
	Release string `json:"release"`
	// URL to an image already stored in Azure infrastructure
	// that can be copied into an image gallery.  Avoid creating VMs directly
	// from this URL as that may lead to performance limitations.
	URL string `json:"url,omitempty"`
}

// Marketplace contains marketplace images for all clouds.
type Marketplace struct {
	Azure *AzureMarketplace `json:"azure,omitempty"`
}

// AzureMarketplaceImage defines the attributes for an Azure
// marketplace image.
type AzureMarketplaceImage struct {
	Publisher       string `json:"publisher"`
	Offer           string `json:"offer"`
	SKU             string `json:"sku"`
	Version         string `json:"version"`
	ThirdPartyImage bool   `json:"thirdParty"`
}

// URN returns the image URN for the marketplace image.
func (i *AzureMarketplaceImage) URN() string {
	return fmt.Sprintf("%s:%s:%s:%s", i.Publisher, i.Offer, i.SKU, i.Version)
}

// AzureMarketplaceImages contains both the HyperV- Gen1 & Gen2
// images for a purchase plan.
type AzureMarketplaceImages struct {
	Gen1 *AzureMarketplaceImage `json:"hyperVGen1,omitempty"`
	Gen2 *AzureMarketplaceImage `json:"hyperVGen2,omitempty"`
}

// AzureMarketplace lists images, both paid and
// unpaid, available in the Azure marketplace.
type AzureMarketplace struct {
	// NoPurchasePlan is the standard, unpaid RHCOS image.
	NoPurchasePlan *AzureMarketplaceImages `json:"no-purchase-plan,omitempty"`
}
