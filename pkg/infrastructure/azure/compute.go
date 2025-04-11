package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/sirupsen/logrus"
)

// CreateImageGalleryInput contains the input parameters for creating a image
// gallery.
type CreateImageGalleryInput struct {
	SubscriptionID     string
	ResourceGroupName  string
	GalleryName        string
	Region             string
	Tags               map[string]*string
	TokenCredential    azcore.TokenCredential
	CloudConfiguration cloud.Configuration
}

// CreateImageGalleryOutput contains the return values after creating a image
// gallery.
type CreateImageGalleryOutput struct {
	ComputeClientFactory *armcompute.ClientFactory
	Gallery              *armcompute.Gallery
}

// CreateImageGallery creates a image gallery.
func CreateImageGallery(ctx context.Context, in *CreateImageGalleryInput) (*CreateImageGalleryOutput, error) {
	logrus.Debugf("Creating image gallery: %s", in.GalleryName)
	computeClientFactory, err := armcompute.NewClientFactory(
		in.SubscriptionID,
		in.TokenCredential,
		&arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Cloud: in.CloudConfiguration,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get compute client factory: %w", err)
	}

	galleriesClient := computeClientFactory.NewGalleriesClient()

	galleriesPoller, err := galleriesClient.BeginCreateOrUpdate(
		ctx,
		in.ResourceGroupName,
		in.GalleryName,
		armcompute.Gallery{
			Location:   to.Ptr(in.Region),
			Properties: &armcompute.GalleryProperties{
				// Fill this in
			},
			Tags: in.Tags,
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gallery %s: %w", in.GalleryName, err)
	}

	galleriesPollDone, err := galleriesPoller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to finish creating gallery %s: %w", in.GalleryName, err)
	}

	logrus.Debugf("Image gallery %s successfully created", in.GalleryName)
	return &CreateImageGalleryOutput{
		ComputeClientFactory: computeClientFactory,
		Gallery:              to.Ptr(galleriesPollDone.Gallery),
	}, nil
}

// CreateGalleryImageInput contains the input parameters for creating a gallery
// image.
type CreateGalleryImageInput struct {
	ResourceGroupName    string
	GalleryName          string
	GalleryImageName     string
	Region               string
	Publisher            string
	Offer                string
	SKU                  string
	Tags                 map[string]*string
	TokenCredential      azcore.TokenCredential
	CloudConfiguration   cloud.Configuration
	Architecture         armcompute.Architecture
	OSType               armcompute.OperatingSystemTypes
	OSState              armcompute.OperatingSystemStateTypes
	HyperVGeneration     armcompute.HyperVGeneration
	ComputeClientFactory *armcompute.ClientFactory
	SecurityType         string
}

// CreateGalleryImageOutput contains the return values after creating a gallery
// image.
type CreateGalleryImageOutput struct {
	GalleryImage *armcompute.GalleryImage
}

// CreateGalleryImage creates a gallery image.
func CreateGalleryImage(ctx context.Context, in *CreateGalleryImageInput) (*CreateGalleryImageOutput, error) {
	logrus.Debugf("Creating gallery image: %s", in.GalleryImageName)

	// Set `Features` within the image properties of Gen V2 images
	// so that they can be used when encryptionAtHost is enabled.
	var features []*armcompute.GalleryImageFeature
	if in.SecurityType != "" {
		feature := armcompute.GalleryImageFeature{
			Name:  to.Ptr("SecurityType"),
			Value: to.Ptr(in.SecurityType),
		}
		features = append(features, to.Ptr(feature))
	}
	if in.HyperVGeneration == armcompute.HyperVGenerationV2 {
		feature := armcompute.GalleryImageFeature{
			Name:  to.Ptr("DiskControllerTypes"),
			Value: to.Ptr("SCSI, NVMe"),
		}
		features = append(features, to.Ptr(feature))
	}

	galleryImagesClient := in.ComputeClientFactory.NewGalleryImagesClient()
	galleryImagesPoller, err := galleryImagesClient.BeginCreateOrUpdate(
		ctx,
		in.ResourceGroupName,
		in.GalleryName,
		in.GalleryImageName,
		armcompute.GalleryImage{
			Location: to.Ptr(in.Region),
			Properties: &armcompute.GalleryImageProperties{
				Architecture:     to.Ptr(in.Architecture),
				OSType:           to.Ptr(in.OSType),
				OSState:          to.Ptr(in.OSState),
				HyperVGeneration: to.Ptr(in.HyperVGeneration),
				Identifier: &armcompute.GalleryImageIdentifier{
					Publisher: to.Ptr(in.Publisher),
					Offer:     to.Ptr(in.Offer),
					SKU:       to.Ptr(in.SKU),
				},
				Features: features,
			},
			Tags: in.Tags,
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gallery image %s: %w", in.GalleryImageName, err)
	}

	galleryImagesPollDone, err := galleryImagesPoller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to finish creating gallery image %s: %w", in.GalleryImageName, err)
	}
	galleryImage := galleryImagesPollDone.GalleryImage
	logrus.Infof("GalleryImage.ID=%s", *galleryImage.ID)

	logrus.Debugf("Gallery image %s successfully created", in.GalleryImageName)
	return &CreateGalleryImageOutput{
		GalleryImage: to.Ptr(galleryImage),
	}, nil
}

// CreateGalleryImageVersionInput contains the input parameters for creating a
// gallery image version.
type CreateGalleryImageVersionInput struct {
	ResourceGroupName       string
	GalleryName             string
	GalleryImageName        string
	GalleryImageVersionName string
	Region                  string
	StorageAccountID        string
	BlobURL                 string
	Tags                    map[string]*string
	RegionalReplicaCount    int32
	ComputeClientFactory    *armcompute.ClientFactory
}

// CreateGalleryImageVersionOutput contains the return values after create a
// gallery image version.
type CreateGalleryImageVersionOutput struct {
	GalleryImageVersion *armcompute.GalleryImageVersion
}

// CreateGalleryImageVersion creates a gallery image version.
func CreateGalleryImageVersion(ctx context.Context, in *CreateGalleryImageVersionInput) (*CreateGalleryImageVersionOutput, error) {
	logrus.Debugf("Creating gallery image version: %s", in.GalleryImageVersionName)
	galleryImageVersionsClient := in.ComputeClientFactory.NewGalleryImageVersionsClient()

	galleryImageVersionProperties := armcompute.GalleryImageVersionProperties{
		StorageProfile: &armcompute.GalleryImageVersionStorageProfile{
			OSDiskImage: &armcompute.GalleryOSDiskImage{
				Source: &armcompute.GalleryDiskImageSource{
					StorageAccountID: to.Ptr(in.StorageAccountID),
					URI:              to.Ptr(in.BlobURL),
				},
			},
		},
		PublishingProfile: &armcompute.GalleryImageVersionPublishingProfile{
			TargetRegions: []*armcompute.TargetRegion{
				{
					Name:                 to.Ptr(in.Region),
					RegionalReplicaCount: to.Ptr(in.RegionalReplicaCount),
				},
			},
		},
	}

	galleryImageVersionPoller, err := galleryImageVersionsClient.BeginCreateOrUpdate(
		ctx,
		in.ResourceGroupName,
		in.GalleryName,
		in.GalleryImageName,
		in.GalleryImageVersionName,
		armcompute.GalleryImageVersion{
			Location:   to.Ptr(in.Region),
			Properties: &galleryImageVersionProperties,
		},
		&armcompute.GalleryImageVersionsClientBeginCreateOrUpdateOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gallery image version %s: %w", in.GalleryImageVersionName, err)
	}

	galleryImageVersionPollDone, err := galleryImageVersionPoller.PollUntilDone(ctx,
		&runtime.PollUntilDoneOptions{
			Frequency: time.Minute * 5,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to finish creating gallery image version %s: %w", in.GalleryImageVersionName, err)
	}

	logrus.Debugf("Gallery image version %s successfully created", in.GalleryImageVersionName)
	return &CreateGalleryImageVersionOutput{
		GalleryImageVersion: to.Ptr(galleryImageVersionPollDone.GalleryImageVersion),
	}, nil
}
