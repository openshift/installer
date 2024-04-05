package azure

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/coreos/stream-metadata-go/arch"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	aztypes "github.com/openshift/installer/pkg/types/azure"
)

// Provider implements Azure CAPI installation.
type Provider struct {
	ResourceGroupName    string
	StorageAccountName   string
	StorageURL           string
	StorageAccount       *armstorage.Account
	StorageClientFactory *armstorage.ClientFactory
	StorageAccountKeys   []armstorage.AccountKey
	Tags                 map[string]*string
	lbBackendAddressPool *armnetwork.BackendAddressPool
}

var _ clusterapi.PreProvider = (*Provider)(nil)
var _ clusterapi.InfraReadyProvider = (*Provider)(nil)
var _ clusterapi.PostProvider = (*Provider)(nil)
var _ clusterapi.IgnitionProvider = (*Provider)(nil)

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
	p.Tags = tags

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

	userTags := platform.UserTags
	tags := make(map[string]*string, len(userTags)+1)
	tags[fmt.Sprintf("kubernetes.io_cluster.%s", in.InfraID)] = ptr.To("owned")
	for k, v := range userTags {
		tags[k] = ptr.To(v)
	}

	tokenCredential := session.TokenCreds
	storageURL := fmt.Sprintf("https://%s.blob.core.windows.net", storageAccountName)
	blobURL := fmt.Sprintf("%s/%s/%s", storageURL, containerName, blobName)

	// Create user assigned identity
	userAssignedIdentityName := fmt.Sprintf("%s-identity", in.InfraID)
	armmsiClientFactory, err := armmsi.NewClientFactory(
		subscriptionID,
		tokenCredential,
		&arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Cloud: cloudConfiguration,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create armmsi client: %w", err)
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
		return fmt.Errorf("failed to create user assigned identity %s: %w", userAssignedIdentityName, err)
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

	networkClientFactory, err := armnetwork.NewClientFactory(subscriptionID, session.TokenCreds,
		&arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Cloud: cloudConfiguration,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error creating network client factory: %w", err)
	}

	lbClient := networkClientFactory.NewLoadBalancersClient()

	lbInput := &lbInput{
		infraID:        in.InfraID,
		region:         in.InstallConfig.Config.Azure.Region,
		resourceGroup:  resourceGroupName,
		subscriptionID: session.Credentials.SubscriptionID,
		lbClient:       lbClient,
		pipClient:      networkClientFactory.NewPublicIPAddressesClient(),
		tags:           p.Tags,
	}

	intLoadBalancer, err := updateInternalLoadBalancer(ctx, lbInput)
	if err != nil {
		return fmt.Errorf("failed to update internal load balancer: %w", err)
	}
	logrus.Debugf("updated internal load balancer: %s", *intLoadBalancer.ID)

	var lbBap *armnetwork.BackendAddressPool
	var extLBFQDN string
	if in.InstallConfig.Config.Publish == types.ExternalPublishingStrategy {
		publicIP, err := createPublicIP(ctx, lbInput)
		if err != nil {
			return fmt.Errorf("failed to create public ip: %w", err)
		}
		logrus.Debugf("created public ip: %s", *publicIP.ID)

		loadBalancer, err := createExternalLoadBalancer(ctx, publicIP, lbInput)
		if err != nil {
			return fmt.Errorf("failed to create load balancer: %w", err)
		}
		logrus.Debugf("created load balancer: %s", *loadBalancer.ID)
		lbBap = loadBalancer.Properties.BackendAddressPools[0]
		extLBFQDN = *publicIP.Properties.DNSSettings.Fqdn
	}

	// Save context for other hooks
	p.ResourceGroupName = resourceGroupName
	p.StorageAccountName = storageAccountName
	p.StorageURL = storageURL
	p.StorageAccount = storageAccount
	p.StorageClientFactory = storageClientFactory
	p.StorageAccountKeys = storageAccountKeys
	p.lbBackendAddressPool = lbBap

	if err := createDNSEntries(ctx, in, extLBFQDN, resourceGroupName); err != nil {
		return fmt.Errorf("error creating DNS records: %w", err)
	}

	return nil
}

// PostProvision provisions an external Load Balancer (when appropriate), and adds configuration
// for the MCS to the CAPI-provisioned internal LB.
func (p *Provider) PostProvision(ctx context.Context, in clusterapi.PostProvisionInput) error {
	ssn, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return fmt.Errorf("error retrieving Azure session: %w", err)
	}
	subscriptionID := ssn.Credentials.SubscriptionID
	cloudConfiguration := ssn.CloudConfig

	if in.InstallConfig.Config.Publish == types.ExternalPublishingStrategy {
		vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, ssn.TokenCreds,
			&arm.ClientOptions{
				ClientOptions: policy.ClientOptions{
					Cloud: cloudConfiguration,
				},
			},
		)
		if err != nil {
			return fmt.Errorf("error creating vm client: %w", err)
		}
		nicClient, err := armnetwork.NewInterfacesClient(ssn.Credentials.SubscriptionID, ssn.TokenCreds,
			&arm.ClientOptions{
				ClientOptions: policy.ClientOptions{
					Cloud: cloudConfiguration,
				},
			},
		)
		if err != nil {
			return fmt.Errorf("error creating nic client: %w", err)
		}

		vmIDs, err := getControlPlaneIDs(in.Client, in.InstallConfig.Config.ControlPlane.Replicas, in.InfraID)
		if err != nil {
			return fmt.Errorf("failed to get control plane VM IDs: %w", err)
		}

		vmInput := &vmInput{
			infraID:       in.InfraID,
			resourceGroup: p.ResourceGroupName,
			vmClient:      vmClient,
			nicClient:     nicClient,
			ids:           vmIDs,
			bap:           p.lbBackendAddressPool,
		}

		if err = associateVMToBackendPool(ctx, *vmInput); err != nil {
			return fmt.Errorf("failed to associate control plane VMs with external load balancer: %w", err)
		}
	}

	return nil
}

func getControlPlaneIDs(cl client.Client, replicas *int64, infraID string) ([]string, error) {
	res := []string{}
	total := int64(1)
	if replicas != nil {
		total = *replicas
	}
	for i := int64(0); i < total; i++ {
		machineName := fmt.Sprintf("%s-master-%d", infraID, i)
		key := client.ObjectKey{
			Name:      machineName,
			Namespace: capiutils.Namespace,
		}
		azureMachine := &capz.AzureMachine{}
		if err := cl.Get(context.Background(), key, azureMachine); err != nil {
			return nil, fmt.Errorf("failed to get AzureMachine: %w", err)
		}
		if vmID := azureMachine.Spec.ProviderID; vmID != nil && len(*vmID) != 0 {
			res = append(res, *azureMachine.Spec.ProviderID)
		} else {
			return nil, fmt.Errorf("%s .Spec.ProviderID is empty", machineName)
		}
	}

	bootstrapName := capiutils.GenerateBoostrapMachineName(infraID)
	key := client.ObjectKey{
		Name:      bootstrapName,
		Namespace: capiutils.Namespace,
	}
	azureMachine := &capz.AzureMachine{}
	if err := cl.Get(context.Background(), key, azureMachine); err != nil {
		return nil, fmt.Errorf("failed to get AzureMachine: %w", err)
	}
	if vmID := azureMachine.Spec.ProviderID; vmID != nil && len(*vmID) != 0 {
		res = append(res, *azureMachine.Spec.ProviderID)
	} else {
		return nil, fmt.Errorf("%s .Spec.ProviderID is empty", bootstrapName)
	}
	return res, nil
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

// Ignition provisions the Azure container that holds the bootstrap ignition
// file.
func (p Provider) Ignition(ctx context.Context, in clusterapi.IgnitionInput) ([]byte, error) {
	session, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	bootstrapIgnData := in.BootstrapIgnData
	subscriptionID := session.Credentials.SubscriptionID
	cloudConfiguration := session.CloudConfig

	ignitionContainerName := "ignition"
	blobName := "bootstrap.ign"
	blobURL := fmt.Sprintf("%s/%s/%s", p.StorageURL, ignitionContainerName, blobName)

	// Create ignition blob storage container
	createBlobContainerOutput, err := CreateBlobContainer(ctx, &CreateBlobContainerInput{
		ContainerName:        ignitionContainerName,
		SubscriptionID:       subscriptionID,
		ResourceGroupName:    p.ResourceGroupName,
		StorageAccountName:   p.StorageAccountName,
		StorageClientFactory: p.StorageClientFactory,
	})
	if err != nil {
		return nil, err
	}

	blobIgnitionContainer := createBlobContainerOutput.BlobContainer
	logrus.Debugf("BlobIgnitionContainer.ID=%s", *blobIgnitionContainer.ID)

	sasURL, err := CreateBlockBlob(ctx, &CreateBlockBlobInput{
		StorageURL:         p.StorageURL,
		BlobURL:            blobURL,
		StorageAccountName: p.StorageAccountName,
		StorageAccountKeys: p.StorageAccountKeys,
		CloudConfiguration: cloudConfiguration,
		BootstrapIgnData:   bootstrapIgnData,
	})
	if err != nil {
		return nil, err
	}
	ignShim, err := bootstrap.GenerateIgnitionShimWithCertBundleAndProxy(sasURL, in.InstallConfig.Config.AdditionalTrustBundle, in.InstallConfig.Config.Proxy)
	if err != nil {
		return nil, fmt.Errorf("failed to create ignition shim: %w", err)
	}

	return ignShim, nil
}
