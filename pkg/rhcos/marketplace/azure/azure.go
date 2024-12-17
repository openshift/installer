package azure

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/coreos/stream-metadata-go/stream/rhcos"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"

	azuresession "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/types/azure"
)

const (
	// region is an arbitrarily chosen region. Marketplace
	// images are published globally, we just need to verify
	// the image exists, so we can use any region.

	region = "centralus"

	// image attributes for the NoPurchasePlan image,
	// published by ARO.
	pubARO   = "azureopenshift"
	offerARO = "aro4"

	// images attributes for paid marketplace images
	pubRH    = "redhat"
	offerOCP = "rh-ocp-worker"
	offerOPP = "rh-opp-worker"
	offerOKE = "rh-oke-worker"

	// supported architectures
	x86   = "x86_64"
	arm64 = "aarch64"
)

type MarketplaceStream struct {
	client *armcompute.VirtualMachineImagesClient
}

type imgsQuery struct {
	gen1, gen2 *imgQuery
}

type imgQuery struct {
	publisher, offer, sku, xyVersion string
}

func NewStreamClient() (*MarketplaceStream, error) {
	cl, err := getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create azure marketplace stream client: %w", err)
	}
	return &MarketplaceStream{cl}, nil
}

func (az *MarketplaceStream) NoPurchasePlan(ctx context.Context, arch, release string) (*rhcos.AzureMarketplaceImages, error) {
	logrus.Info("Retrieving NoPurchase Plan Images for release: ", release)
	gen1SKU, gen2SKU := parseAROSKUs(release, arch)
	q := imgsQuery{
		gen1: &imgQuery{
			publisher: pubARO,
			offer:     offerARO,
			sku:       gen1SKU,
			xyVersion: release,
		},
		gen2: &imgQuery{
			publisher: pubARO,
			offer:     offerARO,
			sku:       gen2SKU,
			xyVersion: release,
		},
	}
	return az.getImages(ctx, q, arch)
}

func (az *MarketplaceStream) OCP(ctx context.Context, arch, release string) (*rhcos.AzureMarketplaceImages, error) {
	if arch != x86 {
		return nil, nil
	}
	return az.getImages(ctx, paidImageQuery(release, offerOCP), arch)
}

func (az *MarketplaceStream) OPP(ctx context.Context, arch, release string) (*rhcos.AzureMarketplaceImages, error) {
	if arch != x86 {
		return nil, nil
	}
	return az.getImages(ctx, paidImageQuery(release, offerOPP), arch)
}

func (az *MarketplaceStream) OKE(ctx context.Context, arch, release string) (*rhcos.AzureMarketplaceImages, error) {
	if arch != x86 {
		return nil, nil
	}
	return az.getImages(ctx, paidImageQuery(release, offerOKE), arch)
}

func paidImageQuery(release, offer string) imgsQuery {
	return imgsQuery{
		gen1: &imgQuery{
			publisher: pubRH,
			offer:     offer,
			sku:       fmt.Sprintf("%s-gen1", offer),
			xyVersion: release,
		},
		gen2: &imgQuery{
			publisher: pubRH,
			offer:     offer,
			sku:       offer,
			xyVersion: release,
		},
	}
}

func (az *MarketplaceStream) getImages(ctx context.Context, query imgsQuery, arch string) (*rhcos.AzureMarketplaceImages, error) {
	imgs := &rhcos.AzureMarketplaceImages{}
	if gen1 := query.gen1; gen1 != nil && gen1.sku != "" {
		logrus.Infof("Searching for image with publisher: %s, offer %s, sku %s architecture %s in release %s", gen1.publisher, gen1.offer, gen1.sku, arch, gen1.xyVersion)
		img, err := az.getImage(ctx, gen1.publisher, gen1.offer, gen1.sku, gen1.xyVersion, arch)
		if err != nil {
			return nil, err
		}
		imgs.Gen1 = img
	}
	if gen2 := query.gen2; gen2 != nil && gen2.sku != "" {
		logrus.Infof("Searching for image with publisher: %s, offer %s, sku %s architecture %s in release %s", gen2.publisher, gen2.offer, gen2.sku, arch, gen2.xyVersion)
		img, err := az.getImage(ctx, gen2.publisher, gen2.offer, gen2.sku, gen2.xyVersion, arch)
		if err != nil {
			return nil, err
		}
		imgs.Gen2 = img
	}
	if imgs.Gen1 == nil && imgs.Gen2 == nil {
		return nil, nil
	}
	return imgs, nil
}

// getImage finds the latest version matching the x.y version of the release.
func (az *MarketplaceStream) getImage(ctx context.Context, pub, offer, sku, xyVersion, arch string) (*rhcos.AzureMarketplaceImage, error) {
	resp, err := az.client.List(ctx, region, pub, offer, sku, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	if len(resp.VirtualMachineImageResourceArray) == 0 {
		logrus.Infof("Found no images for publisher: %s, offer: %s, sku %s for the architecture %s in the release %s", pub, offer, sku, arch, xyVersion)
		return nil, nil
	}

	spew.Dump(resp.VirtualMachineImageResourceArray)

	var foundVersion string
	var greatestSemver string
	for _, v := range resp.VirtualMachineImageResourceArray {
		v := *v.Name
		semVer := convertToSemver(v)

		// Ensure that the image is not from a later Y stream,
		// e.g. if we are populating a 4.19 stream, we don't want 4.20 images,
		// but 4.18 would be ok if 4.19 is not available yet.
		if checkIfNewer(semVer, xyVersion) {
			logrus.Infof("Skipping version %s as it is released after %s", v, xyVersion)
			continue
		}

		if semver.Compare(greatestSemver, semVer) < 0 {
			greatestSemver = semVer
			foundVersion = v
		}
	}

	logrus.Infof("Found version: %s", foundVersion)
	// Now that we've found the version, check the architecture and the plan.
	img, err := az.client.Get(ctx, region, pub, offer, sku, foundVersion, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get the image for the found version, urn: %s:%s:%s:%s in region %s: %w", pub, offer, sku, foundVersion, region, err)
	}

	// This way of checking architecture is works,
	// but may be unnecessary. We would only need to do something
	// like this if the URN for different architectures can be the same;
	// otherwise we know before generating the query which architecture we are looking for.
	azureArch := map[string]armcompute.ArchitectureTypes{
		x86:   armcompute.ArchitectureTypesX64,
		arm64: armcompute.ArchitectureTypesArm64,
	}

	if *img.Properties.Architecture != azureArch[arch] {
		return nil, nil
	}

	logrus.Infof("Found matching image %s:%s:%s:%s", pub, offer, sku, foundVersion)
	return &rhcos.AzureMarketplaceImage{
		Publisher:    pub,
		Offer:        offer,
		SKU:          sku,
		Version:      foundVersion,
		PurchasePlan: img.Properties.Plan != nil,
	}, nil
}

// parseARO takes the release from coreos stream and
// uses conventions to generate the SKU (gen1 & gen2) and version.
// For instance, with a coreos release of "4.19"
// gen1SKU: "aro_418"
// gen2SKU: "418-v2"
// version: "418.94.20241009" (removes timestamp & build number)
func parseAROSKUs(release, arch string) (string, string) {
	xyVersion := strings.ReplaceAll(release, ".", "")
	var gen1SKU, gen2SKU string
	switch arch {
	case x86:
		gen1SKU = fmt.Sprintf("aro_%s", xyVersion)
		gen2SKU = fmt.Sprintf("%s-v2", xyVersion)
	case arm64:
		gen1SKU = ""
		gen2SKU = fmt.Sprintf("%s-arm", xyVersion)
	}
	return gen1SKU, gen2SKU
}

func getClient() (*armcompute.VirtualMachineImagesClient, error) {
	// TODO(padillon): check if we need to expand to other clouds (govcloud) or specify endpoint
	ssn, err := azuresession.GetSession(azure.PublicCloud, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	client, err := armcompute.NewVirtualMachineImagesClient(ssn.Credentials.SubscriptionID, ssn.TokenCreds, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	return client, nil
}

// ARO & the paid marketplace images format their version strings differently
// this function takes either one, and converts it to a go semantic version string.
// ARO combines xy into what should be the x, and includes the RHEL version in y; e.g. 418.94.20250122
// Paid marketplace images use a correct semantic version (well, they use a timestamp for z, but good enough): 4.17.2024101109
func convertToSemver(ver string) string {
	// Using RHEL versioning
	if major := strings.Split(ver, ".")[0]; major == "9" || major == "10" {
		return fmt.Sprintf("v%s", ver)
	}

	if segments := strings.Split(ver, "."); len(segments[0]) == 1 {
		semV := fmt.Sprintf("v%s", ver)
		return semV
	} else if len(segments[0]) == 3 {
		combinedXY := segments[0]
		semV := fmt.Sprintf("v%s", strings.Join([]string{combinedXY[:1], combinedXY[1:], segments[2]}, "."))
		return semV
	}
	return ""
}

func checkIfNewer(candidate, release string) bool {
	img, err := strconv.Atoi(strings.Split(semver.MajorMinor(candidate), ".")[1])
	if err != nil {
		logrus.Infof("Error converting minor version to int with version %s", candidate)
		return true
	}
	rel, err := strconv.Atoi(strings.Split(release, ".")[1])
	if err != nil {
		logrus.Infof("Error converting minor version to int with version %s", release)
		return true
	}
	return img > rel

}
