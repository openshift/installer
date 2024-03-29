package azure

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/coreos/stream-metadata-go/arch"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/rhcos"
	aztypes "github.com/openshift/installer/pkg/types/azure"
)

// Provider implements Azure CAPI installation.
type Provider struct {
	clusterapi.InfraProvider
	ResourceGroupName    string
	StorageAccountName   string
	StorageURL           string
	StorageAccount       *armstorage.Account
	StorageClientFactory *armstorage.ClientFactory
	StorageAccountKeys   []armstorage.AccountKey
	Tags                 map[string]*string
}

var _ clusterapi.PreProvider = (*Provider)(nil)
var _ clusterapi.InfraReadyProvider = (*Provider)(nil)

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return aztypes.Name
}

// PreProvision is called before provisioning using CAPI controllers has begun.
func (p *Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
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

	userTags := platform.UserTags
	tags := make(map[string]*string, len(userTags)+1)
	tags[fmt.Sprintf("kubernetes.io_cluster.%s", in.InfraID)] = ptr.To("owned")
	for k, v := range userTags {
		tags[k] = ptr.To(v)
	}

	// Create resource group
	resourcesClientFactory, err := armresources.NewClientFactory(
		subscriptionID,
		tokenCredential,
		&arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Cloud: cloudConfiguration,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to get azure resource groups factory: %w", err)
	}
	resourceGroupsClient := resourcesClientFactory.NewResourceGroupsClient()
	resourceGroup, err := resourceGroupsClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		armresources.ResourceGroup{
			Location:  ptr.To(platform.Region),
			ManagedBy: nil,
			Tags:      tags,
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("error creating resource group %s: %w", resourceGroupName, err)
	}
	logrus.Debugf("ResourceGroup.ID=%s", *resourceGroup.ID)

	p.ResourceGroupName = resourceGroupName
	p.Tags = tags

	return nil
}

// InfraReady is called once the installer infrastructure is ready.
func (p *Provider) InfraReady(ctx context.Context, in clusterapi.InfraReadyInput) error {
	session, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	installConfig := in.InstallConfig.Config
	platform := installConfig.Platform.Azure
	subscriptionID := session.Credentials.SubscriptionID
	cloudConfiguration := session.CloudConfig
	resourceGroupName := p.ResourceGroupName
	tags := p.Tags

	storageAccountName := fmt.Sprintf("cluster%s", randomString(5))
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

	tokenCredential := session.TokenCreds
	storageURL := fmt.Sprintf("https://%s.blob.core.windows.net", storageAccountName)
	blobURL := fmt.Sprintf("%s/%s/%s", storageURL, containerName, blobName)

	// Create user assigned identity
	userAssignedIdentityName := fmt.Sprintf("%s-identity", in.InfraID)
	armmsiClientFactory, err := armmsi.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		log.Fatalf("failed to create armmsi client: %v", err)
	}
	userAssignedIdentity, err := armmsiClientFactory.NewUserAssignedIdentitiesClient().CreateOrUpdate(
		ctx,
		resourceGroupName,
		userAssignedIdentityName,
		armmsi.Identity{
			Location: ptr.To(platform.Region),
			Tags:     tags,
		},
		nil,
	)
	if err != nil {
		log.Fatalf("failed to create user assigned identity %s: %v", userAssignedIdentityName, err)
	}
	principalID := *userAssignedIdentity.Properties.PrincipalID

	logrus.Debugf("UserAssignedIdentity.ID=%s", *userAssignedIdentity.ID)
	logrus.Debugf("PrinciapalID=%s", principalID)

	// Create storage account
	createStorageAccountOutput, err := CreateStorageAccount(ctx, &CreateStorageAccountInput{
		SubscriptionID:     subscriptionID,
		ResourceGroupName:  resourceGroupName,
		StorageAccountName: storageAccountName,
		CloudName:          platform.CloudName,
		Region:             platform.Region,
		Tags:               tags,
		TokenCredential:    tokenCredential,
		CloudConfiguration: cloudConfiguration,
	})
	if err != nil {
		return err
	}

	storageAccount := createStorageAccountOutput.StorageAccount
	storageClientFactory := createStorageAccountOutput.StorageClientFactory
	storageAccountKeys := createStorageAccountOutput.StorageAccountKeys

	logrus.Debugf("StorageAccount.ID=%s", *storageAccount.ID)

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
		return err
	}

	// Create image gallery
	createImageGalleryOutput, err := CreateImageGallery(ctx, &CreateImageGalleryInput{
		SubscriptionID:     subscriptionID,
		ResourceGroupName:  resourceGroupName,
		GalleryName:        galleryName,
		Region:             platform.Region,
		Tags:               tags,
		TokenCredential:    tokenCredential,
		CloudConfiguration: cloudConfiguration,
	})
	if err != nil {
		return err
	}

	computeClientFactory := createImageGalleryOutput.ComputeClientFactory

	// Create gallery images
	_, err = CreateGalleryImage(ctx, &CreateGalleryImageInput{
		ResourceGroupName:    resourceGroupName,
		GalleryName:          galleryName,
		GalleryImageName:     galleryImageName,
		Region:               platform.Region,
		Tags:                 tags,
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
		return err
	}

	_, err = CreateGalleryImage(ctx, &CreateGalleryImageInput{
		ResourceGroupName:    resourceGroupName,
		GalleryName:          galleryName,
		GalleryImageName:     galleryGen2ImageName,
		Region:               platform.Region,
		Tags:                 tags,
		TokenCredential:      tokenCredential,
		CloudConfiguration:   cloudConfiguration,
		OSType:               armcompute.OperatingSystemTypesLinux,
		OSState:              armcompute.OperatingSystemStateTypesGeneralized,
		HyperVGeneration:     armcompute.HyperVGenerationV1,
		Publisher:            "RedHat-gen2",
		Offer:                "rhcos-gen2",
		SKU:                  "gen2",
		ComputeClientFactory: computeClientFactory,
	})
	if err != nil {
		return err
	}

	// Create gallery image versions
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
		return err
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
		return err
	}

	// Save context for other hooks
	p.ResourceGroupName = resourceGroupName
	p.StorageAccountName = storageAccountName
	p.StorageURL = storageURL
	p.StorageAccount = storageAccount
	p.StorageClientFactory = storageClientFactory
	p.StorageAccountKeys = storageAccountKeys

	return createDNSEntries(ctx, in)
}

func randomString(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source) // nolint:gosec
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"

	s := make([]byte, length)
	for i := range s {
		s[i] = chars[rng.Intn(len(chars))]
	}

	return string(s)
}
