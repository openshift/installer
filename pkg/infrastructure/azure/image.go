package azure

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/coreos/stream-metadata-go/arch"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/sirupsen/logrus"
)

func (p *Provider) createImages(ctx context.Context, in clusterapi.PreProvisionInput) error {
	session, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	installConfig := in.InstallConfig.Config
	platform := installConfig.Platform.Azure
	subscriptionID := session.Credentials.SubscriptionID
	cloudConfiguration := session.CloudConfig
	tokenCredential := session.TokenCreds
	resourceGroupName := platform.ClusterResourceGroupName(in.InfraID)

	// Create storage account
	storageAccountName := fmt.Sprintf("%ssa", strings.ReplaceAll(in.InfraID, "-", ""))
	createStorageAccountOutput, err := CreateStorageAccount(ctx, &CreateStorageAccountInput{
		SubscriptionID:     subscriptionID,
		ResourceGroupName:  resourceGroupName,
		StorageAccountName: storageAccountName,
		CloudName:          platform.CloudName,
		Region:             platform.Region,
		Tags:               p.Tags,
		TokenCredential:    tokenCredential,
		CloudConfiguration: cloudConfiguration,
	})
	if err != nil {
		return err
	}

	storageAccount := createStorageAccountOutput.StorageAccount
	storageClientFactory := createStorageAccountOutput.StorageClientFactory
	storageAccountKeys := createStorageAccountOutput.StorageAccountKeys

	containerName := "vhd"
	blobName := fmt.Sprintf("rhcos%s.vhd", randomString(5))

	stream, err := rhcos.FetchCoreOSBuild(ctx)
	if err != nil {
		return fmt.Errorf("failed to get rhcos stream: %w", err)
	}
	archName := arch.RpmArch(string(installConfig.ControlPlane.Architecture))
	streamArch, err := stream.GetArchitecture(archName)
	if err != nil {
		return fmt.Errorf("failed to get rhcos architecture: %w", err)
	}

	azureDisk := streamArch.RHELCoreOSExtensions.AzureDisk
	imageURL := azureDisk.URL

	rawImageVersion := strings.ReplaceAll(azureDisk.Release, "-", "_")
	imageVersion := rawImageVersion[:len(rawImageVersion)-6]

	galleryName := fmt.Sprintf("gallery_%s", strings.ReplaceAll(in.InfraID, "-", "_"))
	galleryImageName := in.InfraID
	galleryImageVersionName := imageVersion
	galleryGen2ImageName := fmt.Sprintf("%s-gen2", in.InfraID)
	galleryGen2ImageVersionName := imageVersion

	headResponse, err := http.Head(imageURL) // nolint:gosec
	if err != nil {
		return fmt.Errorf("failed HEAD request for image URL %s: %w", imageURL, err)
	}

	imageLength := headResponse.ContentLength
	if imageLength%512 != 0 {
		return fmt.Errorf("image length is not alisnged on a 512 byte boundary")
	}

	storageURL := fmt.Sprintf("https://%s.blob.core.windows.net", storageAccountName)
	blobURL := fmt.Sprintf("%s/%s/%s", storageURL, containerName, blobName)

	// Create blob storage container
	createBlobContainerOutput, err := CreateBlobContainer(ctx, &CreateBlobContainerInput{
		SubscriptionID:       subscriptionID,
		ResourceGroupName:    resourceGroupName,
		StorageAccountName:   storageAccountName,
		ContainerName:        containerName,
		StorageClientFactory: storageClientFactory,
	})
	if err != nil {
		return err
	}

	blobContainer := createBlobContainerOutput.BlobContainer
	logrus.Debugf("BlobContainer.ID=%s", *blobContainer.ID)

	// Upload the image to the container
	if _, ok := os.LookupEnv("OPENSHIFT_INSTALL_SKIP_IMAGE_UPLOAD"); !ok {
		logrus.Debug("Starting image upload creation")
		p.imageWaitGroup.Add(1)
		go func() {
			defer p.imageWaitGroup.Done()
			_, err = CreatePageBlob(ctx, &CreatePageBlobInput{
				StorageURL:         storageURL,
				BlobURL:            blobURL,
				ImageURL:           imageURL,
				ImageLength:        imageLength,
				StorageAccountName: storageAccountName,
				StorageAccountKeys: storageAccountKeys,
				CloudConfiguration: cloudConfiguration,
			})
			if err != nil {
				p.imageErrors <- fmt.Errorf("failed to create page blob: %w", err)
				return
			}

			// Create image gallery
			createImageGalleryOutput, err := CreateImageGallery(ctx, &CreateImageGalleryInput{
				SubscriptionID:     subscriptionID,
				ResourceGroupName:  resourceGroupName,
				GalleryName:        galleryName,
				Region:             platform.Region,
				Tags:               p.Tags,
				TokenCredential:    tokenCredential,
				CloudConfiguration: cloudConfiguration,
			})
			if err != nil {
				p.imageErrors <- fmt.Errorf("failed to create image gallery: %w", err)
				return
			}

			computeClientFactory := createImageGalleryOutput.ComputeClientFactory

			// Create gen1 gallery image & version
			p.imageWaitGroup.Add(1)
			go func() {
				defer p.imageWaitGroup.Done()
				logrus.Debug("Starting gen1 image creation")
				_, err = CreateGalleryImage(ctx, &CreateGalleryImageInput{
					ResourceGroupName:    resourceGroupName,
					GalleryName:          galleryName,
					GalleryImageName:     galleryImageName,
					Region:               platform.Region,
					Tags:                 p.Tags,
					TokenCredential:      tokenCredential,
					CloudConfiguration:   cloudConfiguration,
					OSType:               armcompute.OperatingSystemTypesLinux,
					OSState:              armcompute.OperatingSystemStateTypesGeneralized,
					HyperVGeneration:     armcompute.HyperVGenerationV1,
					Publisher:            "RedHat",
					Offer:                "rhcos",
					SKU:                  "basic",
					ComputeClientFactory: computeClientFactory,
				})
				if err != nil {
					p.imageErrors <- fmt.Errorf("failed to create gen1 gallery image: %w", err)
					return
				}

				_, err = CreateGalleryImageVersion(ctx, &CreateGalleryImageVersionInput{
					ResourceGroupName:       resourceGroupName,
					StorageAccountID:        *storageAccount.ID,
					GalleryName:             galleryName,
					GalleryImageName:        galleryImageName,
					GalleryImageVersionName: galleryImageVersionName,
					Region:                  platform.Region,
					BlobURL:                 blobURL,
					RegionalReplicaCount:    int32(1),
					ComputeClientFactory:    computeClientFactory,
				})
				if err != nil {
					p.imageErrors <- fmt.Errorf("failed to create gen1 gallery image version: %w", err)
					return
				}
				logrus.Debug("Successfully created gen1 image")
			}()

			p.imageWaitGroup.Add(1)
			go func() {
				defer p.imageWaitGroup.Done()
				logrus.Debug("Starting gen2 image creation")
				_, err = CreateGalleryImage(ctx, &CreateGalleryImageInput{
					ResourceGroupName:    resourceGroupName,
					GalleryName:          galleryName,
					GalleryImageName:     galleryGen2ImageName,
					Region:               platform.Region,
					Tags:                 p.Tags,
					TokenCredential:      tokenCredential,
					CloudConfiguration:   cloudConfiguration,
					OSType:               armcompute.OperatingSystemTypesLinux,
					OSState:              armcompute.OperatingSystemStateTypesGeneralized,
					HyperVGeneration:     armcompute.HyperVGenerationV2,
					Publisher:            "RedHat-gen2",
					Offer:                "rhcos-gen2",
					SKU:                  "gen2",
					ComputeClientFactory: computeClientFactory,
				})
				if err != nil {
					p.imageErrors <- fmt.Errorf("failed to create gen2 gallery image: %w", err)
					return
				}

				_, err = CreateGalleryImageVersion(ctx, &CreateGalleryImageVersionInput{
					ResourceGroupName:       resourceGroupName,
					StorageAccountID:        *storageAccount.ID,
					GalleryName:             galleryName,
					GalleryImageName:        galleryGen2ImageName,
					GalleryImageVersionName: galleryGen2ImageVersionName,
					Region:                  platform.Region,
					BlobURL:                 blobURL,
					RegionalReplicaCount:    int32(1),
					ComputeClientFactory:    computeClientFactory,
				})
				if err != nil {
					p.imageErrors <- fmt.Errorf("failed to create gen1 gallery image: %w", err)
					return
				}
				logrus.Debug("Successfully created gen1 image")
			}()
		}()
	}

	logrus.Debugf("StorageAccount.ID=%s", *storageAccount.ID)

	// Save context for other hooks
	p.ResourceGroupName = resourceGroupName
	p.StorageAccountName = storageAccountName
	p.StorageURL = storageURL
	p.StorageAccount = storageAccount
	p.StorageAccountKeys = storageAccountKeys
	p.StorageClientFactory = storageClientFactory
	return nil
}
