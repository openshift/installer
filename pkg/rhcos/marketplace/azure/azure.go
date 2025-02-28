package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/coreos/stream-metadata-go/stream/rhcos"
	"github.com/sirupsen/logrus"

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

	x86   = "x86_64"
	arm64 = "aarch64"
)

type MarketplaceStream struct {
	client *armcompute.VirtualMachineImagesClient
}

func NewStreamClient() (*MarketplaceStream, error) {
	cl, err := getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create azure marketplace stream client: %w", err)
	}
	return &MarketplaceStream{cl}, nil
}

func (az *MarketplaceStream) NoPurchasePlan(ctx context.Context, arch, release string) (*rhcos.AzureMarketplaceImages, error) {
	imgs := &rhcos.AzureMarketplaceImages{}
	gen1SKU, gen2SKU, version := parse(release, arch)
	if gen1SKU != "" {
		img, err := az.getImage(ctx, pubARO, offerARO, gen1SKU, version)
		if err != nil {
			return nil, err
		}
		imgs.Gen1 = img
	}
	if gen2SKU != "" {
		img, err := az.getImage(ctx, pubARO, offerARO, gen2SKU, version)
		if err != nil {
			return nil, err
		}
		imgs.Gen2 = img
	}
	return imgs, nil
}

func (az *MarketplaceStream) getImage(ctx context.Context, pub, offer, sku, version string) (*rhcos.AzureMarketplaceImage, error) {
	foundVersion := version
	resp, err := az.client.List(ctx, region, pub, offer, sku, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}
	for _, v := range resp.VirtualMachineImageResourceArray {
		if *v.Name == version {
			logrus.Infof("Found exact version match for %s:%s:%s:%s \n", pub, offer, sku, version)
			foundVersion = version
			break
		} else {
			logrus.Infof("Found image version %s", *v.Name)
			if *v.Name > foundVersion {
				logrus.Infof("Image version %s is the latest version found", *v.Name)
				foundVersion = *v.Name
			}
		}
	}

	//TODO(padillon): handle thirdParty field. Currently defaulting to false is fine with unpaid images.
	return &rhcos.AzureMarketplaceImage{
		Publisher: pub,
		Offer:     offer,
		SKU:       sku,
		Version:   foundVersion,
	}, nil
}

// parse takes the release from coreos stream and
// uses conventions to generate the SKU (gen1 & gen2) and version.
// For instance, with a coreos release of "418.94.202410090804-0"
// gen1SKU: "aro_418"
// gen2SKU: "418-v2"
// version: "418.94.20241009" (removes timestamp & build number)
func parse(release, arch string) (string, string, string) {
	xyVersion := strings.Split(release, ".")[0]
	var gen1SKU, gen2SKU string
	switch arch {
	case x86:
		gen1SKU = fmt.Sprintf("aro_%s", xyVersion)
		gen2SKU = fmt.Sprintf("%s-v2", xyVersion)
	case arm64:
		gen1SKU = ""
		gen2SKU = fmt.Sprintf("%s-arm", xyVersion)
	}
	return gen1SKU, gen2SKU, release[:len(release)-6]
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
