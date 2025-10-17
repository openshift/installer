package rhcos

import "fmt"

// Extensions is data specific to Red Hat Enterprise Linux CoreOS
type Extensions struct {
	AwsWinLi    *AwsWinLi    `json:"aws-winli,omitempty"`
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

// AwsWinLi represents prebuilt AWS Windows License Included Images.
type AwsWinLi = ReplicatedImage

// ReplicatedImage represents an image in all regions of an AWS-like cloud
// This struct was copied from the release package to avoid an import cycle,
// and is used to describe all AWS WinLI Images in all regions.
type ReplicatedImage struct {
	Regions map[string]SingleImage `json:"regions,omitempty"`
}

// SingleImage represents a globally-accessible image or an image in a
// single region of an AWS-like cloud
// This struct was copied from the release package to avoid an import cycle,
// and is used to describe individual AWS WinLI Images.
type SingleImage struct {
	Release string `json:"release"`
	Image   string `json:"image"`
}

// Marketplace contains marketplace images for all clouds.
type Marketplace struct {
	Azure *AzureMarketplace `json:"azure,omitempty"`
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

	// OCP is the paid marketplace image for OpenShift Container Platform.
	OCP *AzureMarketplaceImages `json:"ocp,omitempty"`

	// OPP is the paid marketplace image for OpenShift Platform Plus.
	OPP *AzureMarketplaceImages `json:"opp,omitempty"`

	// OKE is the paid marketplace image for OpenShift Kubernetes Engine.
	OKE *AzureMarketplaceImages `json:"oke,omitempty"`

	// OCPEMEA is the paid marketplace image for OpenShift Container Platform in EMEA regions.
	OCPEMEA *AzureMarketplaceImages `json:"ocp-emea,omitempty"`

	// OPPEMEA is the paid marketplace image for OpenShift Platform Plus in EMEA regions.
	OPPEMEA *AzureMarketplaceImages `json:"opp-emea,omitempty"`

	// OKEEMEA is the paid marketplace image for OpenShift Kubernetes Engine in EMEA regions.
	OKEEMEA *AzureMarketplaceImages `json:"oke-emea,omitempty"`
}

// AzureMarketplaceImage defines the attributes for an Azure
// marketplace image.
type AzureMarketplaceImage struct {
	Publisher string `json:"publisher"`
	Offer     string `json:"offer"`
	SKU       string `json:"sku"`
	Version   string `json:"version"`
}

// URN returns the image URN for the marketplace image.
func (i *AzureMarketplaceImage) URN() string {
	return fmt.Sprintf("%s:%s:%s:%s", i.Publisher, i.Offer, i.SKU, i.Version)
}
