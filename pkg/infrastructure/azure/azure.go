package azure

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/coreos/stream-metadata-go/arch"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	azconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	aztypes "github.com/openshift/installer/pkg/types/azure"
)

const (
	retryTime        = 10 * time.Second
	retryCount       = 6
	confidentialVMST = "ConfidentialVMSupported"
	trustedLaunchST  = "TrustedLaunchsupported"

	// stackAPIVersion is the Azure Stack compatible API version.
	stackAPIVersion = "2019-06-01"

	// stackComputeAPIVersion is the Azure Stack compatible API version for compute resources (VMs & images).
	stackComputeAPIVersion = "2020-06-01"

	// stackDNSAPIVersion is the Azure Stack compatible API version for DNS resources.
	stackDNSAPIVersion = "2018-05-01"

/*
publicFrontendIPv4ConfigName   = "public-lb-ipv4"
publicFrontendIPv6ConfigName   = "public-lb-ipv6"
internalFrontendIPv4ConfigName = "internal-lb-ipv4"
internalFrontendIPv6ConfigName = "internal-lb-ipv6"
*/
)

// Provider implements Azure CAPI installation.
type Provider struct {
	ResourceGroupName      string
	StorageAccountName     string
	StorageURL             string
	FrontendIPConfigName   string
	BackendAddressPoolName string
	DualStack              bool
	StackType              aztypes.StackType
	StorageAccount         *armstorage.Account
	StorageClientFactory   *armstorage.ClientFactory
	StorageAccountKeys     []armstorage.AccountKey
	NetworkClientFactory   *armnetwork.ClientFactory
	lbBackendAddressPools  []*armnetwork.BackendAddressPool
	CloudConfiguration     cloud.Configuration
	TokenCredential        azcore.TokenCredential
	Tags                   map[string]*string
	clientOptions          *arm.ClientOptions
	computeClientOptions   *arm.ClientOptions
}

var _ clusterapi.InfraReadyProvider = (*Provider)(nil)
var _ clusterapi.PostProvider = (*Provider)(nil)
var _ clusterapi.IgnitionProvider = (*Provider)(nil)
var _ clusterapi.PostDestroyer = (*Provider)(nil)
var _ clusterapi.Timeouts = (*Provider)(nil)

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return aztypes.Name
}

// NetworkTimeout uses the default timeout of 15 minutes to satisfy the Timeouts interface.
// Azure only needs special handling for machine provisioning timeouts.
func (p Provider) NetworkTimeout() time.Duration {
	return 15 * time.Minute
}

// ProvisionTimeout bumps the machine provisioning timeout due to
// https://issues.redhat.com/browse/OCPBUGS-43625.
func (p Provider) ProvisionTimeout() time.Duration {
	return 20 * time.Minute
}

// PublicGatherEndpoint indicates that machine ready checks should NOT wait for an ExternalIP
// in the status and should use the API load balancer when gathering bootstrap log bundles.
func (*Provider) PublicGatherEndpoint() clusterapi.GatherEndpoint { return clusterapi.APILoadBalancer }

// InfraReady is called once the installer infrastructure is ready.
//
//nolint:gocyclo //TODO(padillon): forthcoming marketplace image support should help reduce complexity here.
func (p *Provider) InfraReady(ctx context.Context, in clusterapi.InfraReadyInput) error {
	session, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	installConfig := in.InstallConfig.Config
	platform := installConfig.Platform.Azure
	subscriptionID := session.Credentials.SubscriptionID
	cloudConfiguration := session.CloudConfig
	tokenCredential := session.TokenCreds
	p.ResourceGroupName = platform.ClusterResourceGroupName(in.InfraID)

	p.DualStack = false
	if installConfig.Networking.IsDualStack() {
		logrus.Debugf("XXX: Setting configuration to dual stack")
		p.DualStack = true
		//p.FrontendIPConfigName = frontendIPv4ConfigName
		//p.BackendAddressPoolName = fmt.Sprintf("%s-ipv4-api-internal", in.InfraID)
	} else if installConfig.Networking.IsIPv4() {
		logrus.Debugf("XXX: Setting configuration to ipv4")
		//p.FrontendIPConfigName = frontendIPv4ConfigName
		//p.BackendAddressPoolName = fmt.Sprintf("%s-ipv4-api-internal", in.InfraID)
	} else if installConfig.Networking.IsIPv6() {
		logrus.Debugf("XXX: Setting configuration to ipv6")
		//p.FrontendIPConfigName = frontendIPv6ConfigName
		//p.BackendAddressPoolName = fmt.Sprintf("%s-ipv6-api-internal", in.InfraID)
	}

	/*
		p.StackType = aztypes.StackTypeIPv4
		if installConfig.Networking.IsIPv6() {
			p.FrontendIPConfigName = frontendIPv6ConfigName
			p.StackType = aztypes.StackTypeIPv6
		} else {
			p.FrontendIPConfigName = frontendIPv4ConfigName
		}
	*/

	userTags := platform.UserTags
	tags := make(map[string]*string, len(userTags)+1)
	tags[fmt.Sprintf("kubernetes.io_cluster.%s", in.InfraID)] = ptr.To("owned")
	for k, v := range userTags {
		tags[k] = ptr.To(v)
	}
	p.Tags = tags

	opts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: cloudConfiguration,
		},
	}
	computeClientOpts := opts
	if platform.CloudName == aztypes.StackCloud {
		opts.APIVersion = stackAPIVersion
		computeClientOpts = &arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Cloud:      cloudConfiguration,
				APIVersion: stackComputeAPIVersion,
			},
		}
	}
	p.clientOptions = opts
	p.computeClientOptions = computeClientOpts

	if err = handleIdentity(ctx, identityInput{
		installConfig:     installConfig,
		region:            platform.Region,
		resourceGroupName: p.ResourceGroupName,
		subscriptionID:    subscriptionID,
		tokenCredential:   tokenCredential,
		infraID:           in.InfraID,
		clientOpts:        p.clientOptions,
		tags:              p.Tags,
	}); err != nil {
		errMsg := "error creating user-assigned identity: please ensure your user credentials " +
			"have the User Access Admin Role or if you are not utilizing an Azure Container Registry " +
			"you can set installconfig.platform.azure.defaultMachinePlatform.identity.type: None to skip " +
			"the creation of the identity: creation failed with: %w"
		return fmt.Errorf(errMsg, err)
	}

	// Creating a dummy nsg for existing vnets installation to appease the ingress operator.
	if in.InstallConfig.Config.Azure.VirtualNetwork != "" {
		networkClientFactory, err := armnetwork.NewClientFactory(subscriptionID, tokenCredential, p.clientOptions)
		if err != nil {
			return fmt.Errorf("failed to create azure network factory: %w", err)
		}
		securityGroupName := in.InstallConfig.Config.Platform.Azure.NetworkSecurityGroupName(in.InfraID)
		securityGroupsClient := networkClientFactory.NewSecurityGroupsClient()
		pollerResp, err := securityGroupsClient.BeginCreateOrUpdate(
			ctx,
			p.ResourceGroupName,
			securityGroupName,
			armnetwork.SecurityGroup{
				Location: to.Ptr(platform.Region),
				Tags:     p.Tags,
			},
			nil)
		if err != nil {
			return fmt.Errorf("failed to create network security group: %w", err)
		}
		nsg, err := pollerResp.PollUntilDone(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to create network security group: %w", err)
		}
		logrus.Infof("nsg=%s", *nsg.ID)
	}

	var architecture armcompute.Architecture
	if installConfig.ControlPlane.Architecture == types.ArchitectureARM64 {
		architecture = armcompute.ArchitectureArm64
	} else {
		architecture = armcompute.ArchitectureX64
	}

	resourceGroupName := p.ResourceGroupName
	storageAccountName := aztypes.GetStorageAccountName(in.InfraID)
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
		return fmt.Errorf("image length is not aligned on a 512 byte boundary")
	}

	storageURL := fmt.Sprintf("https://%s.blob.%s", storageAccountName, session.Environment.StorageEndpointSuffix)
	blobURL := fmt.Sprintf("%s/%s/%s", storageURL, containerName, blobName)

	var storageAccount *armstorage.Account
	var storageClientFactory *armstorage.ClientFactory
	var storageAccountKeys []armstorage.AccountKey

	var createStorageAccountOutput *CreateStorageAccountOutput
	if platform.CloudName != aztypes.StackCloud {
		// Create storage account
		createStorageAccountOutput, err = CreateStorageAccount(ctx, &CreateStorageAccountInput{
			SubscriptionID:     subscriptionID,
			ResourceGroupName:  resourceGroupName,
			StorageAccountName: storageAccountName,
			CloudName:          platform.CloudName,
			Region:             platform.Region,
			AuthType:           session.AuthType,
			Tags:               tags,
			CustomerManagedKey: platform.CustomerManagedKey,
			TokenCredential:    tokenCredential,
			ClientOpts:         p.clientOptions,
		})
		if err != nil {
			return err
		}
		storageAccount = createStorageAccountOutput.StorageAccount
		storageClientFactory = createStorageAccountOutput.StorageClientFactory
		storageAccountKeys = createStorageAccountOutput.StorageAccountKeys

		logrus.Debugf("StorageAccount.ID=%s", *storageAccount.ID)
	}

	// Upload the image to the container
	_, skipImageUpload := os.LookupEnv("OPENSHIFT_INSTALL_SKIP_IMAGE_UPLOAD")
	if !(skipImageUpload || platform.CloudName == aztypes.StackCloud) {
		// Create vhd blob storage container
		publicAccess := armstorage.PublicAccessNone
		createBlobContainerOutput, err := CreateBlobContainer(ctx, &CreateBlobContainerInput{
			SubscriptionID:       subscriptionID,
			ResourceGroupName:    resourceGroupName,
			StorageAccountName:   storageAccountName,
			ContainerName:        containerName,
			PublicAccess:         to.Ptr(publicAccess),
			StorageClientFactory: storageClientFactory,
		})
		if err != nil {
			return err
		}

		blobContainer := createBlobContainerOutput.BlobContainer
		logrus.Debugf("BlobContainer.ID=%s", *blobContainer.ID)

		_, err = CreatePageBlob(ctx, &CreatePageBlobInput{
			StorageURL:         storageURL,
			BlobURL:            blobURL,
			ImageURL:           imageURL,
			ImageLength:        imageLength,
			StorageAccountName: storageAccountName,
			StorageAccountKeys: storageAccountKeys,
			ClientOpts:         p.clientOptions,
		})
		if err != nil {
			return err
		}

		// Create image gallery
		createImageGalleryOutput, err := CreateImageGallery(ctx, &CreateImageGalleryInput{
			SubscriptionID:    subscriptionID,
			ResourceGroupName: resourceGroupName,
			GalleryName:       galleryName,
			Region:            platform.Region,
			Tags:              tags,
			TokenCredential:   tokenCredential,
			ClientOpts:        p.clientOptions,
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
			Publisher:            "RedHat",
			Offer:                "rhcos",
			SKU:                  "basic",
			Tags:                 tags,
			TokenCredential:      tokenCredential,
			ClientOpts:           p.clientOptions,
			Architecture:         architecture,
			OSType:               armcompute.OperatingSystemTypesLinux,
			OSState:              armcompute.OperatingSystemStateTypesGeneralized,
			HyperVGeneration:     armcompute.HyperVGenerationV1,
			ComputeClientFactory: computeClientFactory,
			SecurityType:         "",
		})
		if err != nil {
			return err
		}
		// If Control Plane Security Type is provided, then pass that along
		// during Gen V2 Gallery Image creation. It will be added as a
		// supported feature of the image.
		securityType, err := getMachinePoolSecurityType(in)
		if err != nil {
			return err
		}

		_, err = CreateGalleryImage(ctx, &CreateGalleryImageInput{
			ResourceGroupName:    resourceGroupName,
			GalleryName:          galleryName,
			GalleryImageName:     galleryGen2ImageName,
			Region:               platform.Region,
			Publisher:            "RedHat-gen2",
			Offer:                "rhcos-gen2",
			SKU:                  "gen2",
			Tags:                 tags,
			TokenCredential:      tokenCredential,
			ClientOpts:           p.clientOptions,
			Architecture:         architecture,
			OSType:               armcompute.OperatingSystemTypesLinux,
			OSState:              armcompute.OperatingSystemStateTypesGeneralized,
			HyperVGeneration:     armcompute.HyperVGenerationV2,
			ComputeClientFactory: computeClientFactory,
			SecurityType:         securityType,
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
	}

	if installConfig.Azure.CloudName == aztypes.StackCloud {
		client, err := armcompute.NewImagesClient(subscriptionID, tokenCredential, p.computeClientOptions)
		if err != nil {
			return fmt.Errorf("error creating stack managed images client: %w", err)
		}
		createManagedImageInput := CreateManagedImageInput{
			VHDBlobURL:        platform.ClusterOSImage,
			ResourceGroupName: resourceGroupName,
			Region:            platform.Region,
			InfraID:           in.InfraID,
			Tags:              tags,
			Client:            client,
		}
		if err = CreateManagedImage(ctx, &createManagedImageInput); err != nil {
			return fmt.Errorf("error creating stack managed image: %w", err)
		}
	}

	networkClientFactory, err := armnetwork.NewClientFactory(subscriptionID, session.TokenCreds, p.clientOptions)
	if err != nil {
		return fmt.Errorf("error creating network client factory: %w", err)
	}
	lbClient := networkClientFactory.NewLoadBalancersClient()

	idPrefix := fmt.Sprintf("subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers",
		session.Credentials.SubscriptionID,
		resourceGroupName,
	)

	lbInput := &lbInput{
		loadBalancerName:     fmt.Sprintf("%s-internal", in.InfraID),
		infraID:              in.InfraID,
		region:               platform.Region,
		resourceGroupName:    resourceGroupName,
		subscriptionID:       session.Credentials.SubscriptionID,
		idPrefix:             idPrefix,
		networkClientFactory: networkClientFactory,
		lbClient:             lbClient,
		tags:                 p.Tags,
	}

	/*
		// XXX 6443
		apiProbe := apiProbe()
		_, err = addProbeToLoadBalancer(ctx, apiProbe, lbInput)
		if err != nil {
			return fmt.Errorf("failed to add api probe to internal load balancer: %w", err)
		}

		// XXX 22623
		mcsProbe := mcsProbe()
		_, err = addProbeToLoadBalancer(ctx, mcsProbe, lbInput)
		if err != nil {
			return fmt.Errorf("failed to add mcs probe to internal load balancer: %w", err)
		}
	*/

	// Get the virtual network
	virtualNetwork, err := getVirtualNetwork(ctx, &vnetInput{
		resourceGroupName:    resourceGroupName,
		virtualNetworkName:   installConfig.Azure.VirtualNetworkName(in.InfraID),
		networkClientFactory: networkClientFactory,
	})
	if err != nil {
		return fmt.Errorf("failed to get virtual network: %w", err)
	}

	// Get the control plane subnet
	controlPlaneSubnet, err := getSubnet(ctx, &subnetInput{
		resourceGroupName:    resourceGroupName,
		virtualNetworkName:   *virtualNetwork.Name,
		subnetName:           installConfig.Azure.ControlPlaneSubnetName(in.InfraID),
		networkClientFactory: networkClientFactory,
	})
	if err != nil {
		return fmt.Errorf("failed to get control plane subnet: %w", err)
	}

	////////////////////////////////////////////////////////////////////////
	//
	// Public load balancer
	//
	//     For all of the following, the IPv4 configurations already exist.
	//     So we just need to create the IPv6 configurations. The exception
	//     being the probes and associated rules.
	//
	// Frontend IP Configuration:
	//
	//   Name                        IP address (backend pool)
	//   ${infraID}-frontEnd         ${infraID}-controlplane-outbound
	//   ${infraID}-frontEnd-ipv6    ${infraID}-controlplane-outbound-ipv6
	//   public-lb-ipv4              ${infraID}-pip-ipv4
	//   public-lb-ipv6              ${infraID}-pip-ipv6
	//
	// Backend pools:
	//
	//   Backend pool
	//   ${infraID}
	//   ${infraID}-ipv6
	//   ${infraID}-outbound-lb-outboundBackendPool
	//   ${infraID}-outbound-lb-outboundBackendPool-ipv6
	//
	// setup / create health probes
	// setup / create load balancing rules
	// setup / create outbound rules
	//
	////////////////////////////////////////////////////////////////////////

	if installConfig.Networking.IsDualStack() {

		// Setup the internal load balancer for dual stack networking
		internalLoadBalancerName := fmt.Sprintf("%s-internal", in.InfraID)
		lbInput.loadBalancerName = internalLoadBalancerName

		logrus.Debugf("XXX: Adding IPv6 frontend IP configuration to internal load balancer")

		// Create IPv6 frontend IP configuration
		frontendIPConfiguration := newFrontendIPConfigurationIPv6(fmt.Sprintf("%s-internal-frontEnd-ipv6", in.InfraID), controlPlaneSubnet)
		lbInput.frontendIPConfiguration = frontendIPConfiguration
		_, err = addFrontendIPConfigurationToLoadBalancer(ctx, lbInput)
		if err != nil {
			return fmt.Errorf("failed to add ipv6 frontend IP configuration to internal load balancer: %w", err)
		}

		logrus.Debugf("XXX: Adding IPv4 backend address pool to internal load balancer")

		// Create IPv4 backend address pool
		/*
			ipv4BackendAddressPool := newBackendAddressPool(fmt.Sprintf("%s-internal-backendPool-ipv4", in.InfraID), virtualNetwork)
			lbInput.backendAddressPool = ipv4BackendAddressPool
			_, err := addBackendAddressPoolToLoadBalancer(ctx, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add ipv4 backend address pool to internal load balancer: %w", err)
			}
		*/

		logrus.Debugf("XXX: Adding IPv6 backend address pool to internal load balancer")

		// Create IPv6 backend address pool
		ipv6BackendAddressPool := newBackendAddressPool(fmt.Sprintf("%s-internal-ipv6", in.InfraID), virtualNetwork)
		lbInput.backendAddressPool = ipv6BackendAddressPool
		_, err = addBackendAddressPoolToLoadBalancer(ctx, lbInput)
		if err != nil {
			return fmt.Errorf("failed to add ipv6 backend address pool to internal load balancer: %w", err)
		}

		logrus.Debugf("XXX: Adding health probes to internal load balancer")

		// Add health probes to the internal load balancer
		/*
			apiProbe := apiProbe()
			_, err = addProbeToLoadBalancer(ctx, apiProbe, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add api probe to internal load balancer: %w", err)
			}
		*/

		mcsProbe := mcsProbe()
		_, err = addProbeToLoadBalancer(ctx, mcsProbe, lbInput)
		if err != nil {
			return fmt.Errorf("failed to add mcs probe to internal load balancer: %w", err)
		}

		logrus.Debugf("XXX: Adding IPv4 rules to internal load balancer")

		// Add IPv4 rules to the internal load balancer
		/*
			apiRuleIPv4 := apiRule(&lbRuleInput{
				idPrefix:               idPrefix,
				loadBalancerName:       internalLoadBalancerName,
				probeName:              *apiProbe.Name,
				ruleName:               "api-internal-ipv4",
				frontendIPConfigName:   fmt.Sprintf("%s-internal-frontEnd-ipv4", in.InfraID),
				backendAddressPoolName: fmt.Sprintf("%s-internal-backendPool-ipv4", in.InfraID),
			})
			_, err = addLoadBalancingRuleToLoadBalancer(ctx, apiRuleIPv4, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add ipv4 api rule to internal load balancer: %w", err)
			}
		*/

		mcsRuleIPv4 := mcsRule(&lbRuleInput{
			idPrefix:               idPrefix,
			loadBalancerName:       internalLoadBalancerName,
			probeName:              *mcsProbe.Name,
			ruleName:               "sint-ipv4",
			frontendIPConfigName:   fmt.Sprintf("%s-internal-frontEnd", in.InfraID),
			backendAddressPoolName: fmt.Sprintf("%s-internal", in.InfraID),
		})
		_, err = addLoadBalancingRuleToLoadBalancer(ctx, mcsRuleIPv4, lbInput)
		if err != nil {
			return fmt.Errorf("failed to add ipv4 mcs rule to internal load balancer: %w", err)
		}

		logrus.Debugf("XXX: Adding IPv6 rules to internal load balancer")

		// Add IPv6 rules to the internal load balancer
		apiRuleIPv6 := apiRule(&lbRuleInput{
			idPrefix:         idPrefix,
			loadBalancerName: internalLoadBalancerName,
			//probeName:              *apiProbe.Name,
			probeName:              "HTTPSProbe",
			ruleName:               "api-internal-ipv6",
			frontendIPConfigName:   fmt.Sprintf("%s-internal-frontEnd-ipv6", in.InfraID),
			backendAddressPoolName: fmt.Sprintf("%s-internal-ipv6", in.InfraID),
		})
		_, err = addLoadBalancingRuleToLoadBalancer(ctx, apiRuleIPv6, lbInput)
		if err != nil {
			return fmt.Errorf("failed to add ipv6 api rule to internal load balancer: %w", err)
		}

		mcsRuleIPv6 := mcsRule(&lbRuleInput{
			idPrefix:               idPrefix,
			loadBalancerName:       internalLoadBalancerName,
			probeName:              *mcsProbe.Name,
			ruleName:               "sint-ipv6",
			frontendIPConfigName:   fmt.Sprintf("%s-internal-frontEnd-ipv6", in.InfraID),
			backendAddressPoolName: fmt.Sprintf("%s-internal-ipv6", in.InfraID),
		})
		_, err = addLoadBalancingRuleToLoadBalancer(ctx, mcsRuleIPv6, lbInput)
		if err != nil {
			return fmt.Errorf("failed to add ipv6 mcs rule to internal load balancer: %w", err)
		}

		/*
			lbInput.stackType = aztypes.StackTypeIPv4
			lbInput.frontendIPConfigName = aztypes.InternalFrontendIPv4ConfigName
			lbInput.backendAddressPoolName = azure.InternalBackendAddressPoolIPv4Name
			intLoadBalancer, err = updateInternalLoadBalancer(ctx, lbInput)
			if err != nil {
				return fmt.Errorf("failed to update internal load balancer: %w", err)
			}
			logrus.Debugf("updated internal load balancer: %s", *intLoadBalancer.ID)

			lbInput.stackType = aztypes.StackTypeIPv6
			lbInput.frontendIPConfigName = aztypes.InternalFrontendIPv6ConfigName
			lbInput.backendAddressPoolName = azure.InternalBackendAddressPoolIPv6Name
			intLoadBalancer, err := updateInternalLoadBalancer(ctx, lbInput)
			if err != nil {
				return fmt.Errorf("failed to update internal load balancer: %w", err)
			}
			logrus.Debugf("updated internal load balancer: %s", *intLoadBalancer.ID)
		*/

	} else if installConfig.Networking.IsIPv4() {
		/*
			logrus.Debugf("XXX: updating internal API load balancer for ipv4 configuration")
			lbInput.stackType = aztypes.StackTypeIPv4
			lbInput.frontendIPConfigName = aztypes.InternalFrontendIPv4ConfigName
			lbInput.backendAddressPoolName = azure.InternalBackendAddressPoolIPv4Name
			intLoadBalancer, err = updateInternalLoadBalancer(ctx, lbInput)
			if err != nil {
				return fmt.Errorf("failed to update internal load balancer: %w", err)
			}
			logrus.Debugf("updated internal load balancer: %s", *intLoadBalancer.ID)
		*/

	} else if installConfig.Networking.IsIPv6() {
		/*
			logrus.Debugf("XXX: updating internal API load balancer for ipv6 configuration")

			lbInput.stackType = aztypes.StackTypeIPv6
			lbInput.frontendIPConfigName = aztypes.InternalFrontendIPv6ConfigName
			lbInput.backendAddressPoolName = azure.InternalBackendAddressPoolIPv6Name
			intLoadBalancer, err := updateInternalLoadBalancer(ctx, lbInput)
			if err != nil {
				return fmt.Errorf("failed to update internal load balancer: %w", err)
			}
			logrus.Debugf("updated internal load balancer: %s", *intLoadBalancer.ID)
		*/
	}

	var lbBaps []*armnetwork.BackendAddressPool
	var extLBFQDNIPv4, extLBFQDNIPv6 string
	var publicIPv4, publicIPv6, frontendIPv6 *armnetwork.PublicIPAddress = nil, nil, nil
	if in.InstallConfig.Config.PublicAPI() {

		publicIPv4, err = createPublicIP(ctx, &pipInput{
			name:              fmt.Sprintf("%s-pip-v4", in.InfraID),
			infraID:           in.InfraID,
			region:            in.InstallConfig.Config.Azure.Region,
			resourceGroupName: resourceGroupName,
			stackType:         aztypes.StackTypeIPv4,
			pipClient:         networkClientFactory.NewPublicIPAddressesClient(),
			tags:              p.Tags,
		})
		if err != nil {
			return fmt.Errorf("failed to create public ipv4 address: %w", err)
		}
		logrus.Debugf("created public ipv4 address: %s", *publicIPv4.ID)

		if p.DualStack {
			publicIPv6, err = createPublicIP(ctx, &pipInput{
				name:              fmt.Sprintf("%s-pip-v6", in.InfraID),
				infraID:           in.InfraID,
				region:            in.InstallConfig.Config.Azure.Region,
				resourceGroupName: resourceGroupName,
				stackType:         aztypes.StackTypeIPv6,
				pipClient:         networkClientFactory.NewPublicIPAddressesClient(),
				tags:              p.Tags,
			})
			if err != nil {
				return fmt.Errorf("failed to create public ipv6 address: %w", err)
			}
			logrus.Debugf("created public ipv6 address: %s", *publicIPv6.ID)

			//52.241.251.184 (pip-jhipv6-yqobuv-crgfk-controlplane-outbound)

			frontendIPv6, err = createPublicIP(ctx, &pipInput{
				name:              fmt.Sprintf("%s-controlplane-outbound-ipv6", in.InfraID),
				infraID:           in.InfraID,
				region:            in.InstallConfig.Config.Azure.Region,
				resourceGroupName: resourceGroupName,
				stackType:         aztypes.StackTypeIPv6,
				pipClient:         networkClientFactory.NewPublicIPAddressesClient(),
				tags:              p.Tags,
			})
			if err != nil {
				return fmt.Errorf("failed to create frontend ipv6 address: %w", err)
			}
			logrus.Debugf("created frontend ipv6 address: %s", *frontendIPv6.ID)
		}

		publicLoadBalancerName := in.InfraID
		lbInput.loadBalancerName = publicLoadBalancerName
		lbInput.backendAddressPoolName = publicLoadBalancerName

		var loadBalancer *armnetwork.LoadBalancer
		if platform.OutboundType == aztypes.UserDefinedRoutingOutboundType {
			loadBalancer, err = createAPILoadBalancer(ctx, publicIPv4, publicIPv6, lbInput)
			if err != nil {
				return fmt.Errorf("failed to create API load balancer: %w", err)
			}
			lbBaps = append(lbBaps, loadBalancer.Properties.BackendAddressPools...)
		} else {

			// Setup the public load balancer for dual stack networking
			logrus.Debugf("XXX: Adding IPv4 frontend IP configuration to public load balancer")

			// Create IPv4 frontend IP configuration
			ipv4PublicFrontendIPConfiguration := newFrontendIPConfigurationIPv6("public-lb-ipv4", controlPlaneSubnet)
			ipv4PublicFrontendIPConfiguration.Properties.PublicIPAddress = publicIPv4
			ipv4PublicFrontendIPConfiguration.Properties.Subnet = nil
			//ipv4PublicFrontendIPConfiguration.PrivateIPAllocationMethod = to.Ptr(armnetwork.IPAllocationMethodDynamic)
			lbInput.frontendIPConfiguration = ipv4PublicFrontendIPConfiguration
			_, err = addFrontendIPConfigurationToLoadBalancer(ctx, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add ipv4 public frontend IP configuration to public load balancer: %w", err)
			}

			logrus.Debugf("XXX: Adding IPv6 frontend IP configuration to public load balancer")

			ipv6PublicFrontendIPConfiguration := newFrontendIPConfigurationIPv6("public-lb-ipv6", controlPlaneSubnet)
			ipv6PublicFrontendIPConfiguration.Properties.PublicIPAddress = publicIPv6
			ipv6PublicFrontendIPConfiguration.Properties.Subnet = nil
			lbInput.frontendIPConfiguration = ipv6PublicFrontendIPConfiguration
			//ipv6PublicFrontendIPConfiguration.PrivateIPAllocationMethod = to.Ptr(armnetwork.IPAllocationMethodDynamic)
			_, err = addFrontendIPConfigurationToLoadBalancer(ctx, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add ipv6 public frontend IP configuration to public load balancer: %w", err)
			}

			// Create IPv6 frontend IP configurations
			ipv6FrontendIPConfiguration := newFrontendIPConfigurationIPv6(fmt.Sprintf("%s-frontEnd-ipv6", in.InfraID), controlPlaneSubnet)
			ipv6FrontendIPConfiguration.Properties.PublicIPAddress = frontendIPv6
			ipv6FrontendIPConfiguration.Properties.Subnet = nil
			lbInput.frontendIPConfiguration = ipv6FrontendIPConfiguration
			_, err = addFrontendIPConfigurationToLoadBalancer(ctx, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add ipv6 frontend IP configuration to public load balancer: %w", err)
			}

			// Create IPv4 backend address pool
			ipv4BackendAddressPool := newBackendAddressPool(fmt.Sprintf("%s", in.InfraID), virtualNetwork)
			lbInput.backendAddressPool = ipv4BackendAddressPool
			_, err = addBackendAddressPoolToLoadBalancer(ctx, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add ipv4 backend address pool to public load balancer: %w", err)
			}

			logrus.Debugf("XXX: Adding IPv6 backend address pools to public load balancer")

			// Create IPv6 backend address pools

			// ${infraid}-frontEnd-ipv6
			ipv6BackendAddressPool := newBackendAddressPool(fmt.Sprintf("%s-ipv6", in.InfraID), virtualNetwork)
			lbInput.backendAddressPool = ipv6BackendAddressPool
			_, err = addBackendAddressPoolToLoadBalancer(ctx, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add ipv6 backend address pool to public load balancer: %w", err)
			}

			// ${infraid}-frontEnd-ipv6
			outboundBackendAddressPool := newBackendAddressPool(fmt.Sprintf("%s-outbound-lb-outboundBackendPool-ipv6", in.InfraID), virtualNetwork)
			lbInput.backendAddressPool = outboundBackendAddressPool
			_, err = addBackendAddressPoolToLoadBalancer(ctx, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add ipv6 outbound backend address pool to public load balancer: %w", err)
			}
			/*
				natbackendAddressPool := newBackendAddressPool("OutboundNATAllProtocols-ipv6", virtualNetwork)
				lbInput.backendAddressPool = natbackendAddressPool
				_, err = addBackendAddressPoolToLoadBalancer(ctx, lbInput)
				if err != nil {
					return fmt.Errorf("failed to add ipv6 nat backend address pool to public load balancer: %w", err)
				}
			*/

			logrus.Debugf("XXX: Adding health probes to public load balancer")

			// Add health probes to the public load balancer
			apiProbe := apiProbe()
			_, err = addProbeToLoadBalancer(ctx, apiProbe, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add api probe to public load balancer: %w", err)
			}

			/*
				mcsProbe := mcsProbe()
				_, err = addProbeToLoadBalancer(ctx, mcsProbe, lbInput)
				if err != nil {
					return fmt.Errorf("failed to add mcs probe to public load balancer: %w", err)
				}
			*/

			logrus.Debugf("XXX: Adding IPv4 rules to public load balancer")

			// Add IPv4 rules to the public load balancer
			apiRuleIPv4 := apiRule(&lbRuleInput{
				idPrefix:             idPrefix,
				loadBalancerName:     publicLoadBalancerName,
				probeName:            *apiProbe.Name,
				ruleName:             "api-ipv4",
				frontendIPConfigName: "public-lb-ipv4",
				//frontendIPConfigName:   fmt.Sprintf("%s-frontEnd", in.InfraID),
				backendAddressPoolName: fmt.Sprintf("%s", in.InfraID),
			})
			_, err = addLoadBalancingRuleToLoadBalancer(ctx, apiRuleIPv4, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add ipv4 api rule to public load balancer: %w", err)
			}

			/*
				outboundApiRuleIPv4 := apiRule(&lbRuleInput{
					idPrefix:               idPrefix,
					loadBalancerName:       publicLoadBalancerName,
					probeName:              *apiProbe.Name,
					ruleName:               "api-ipv4",
					frontendIPConfigName:   "public-lb-ipv4",
					backendAddressPoolName: fmt.Sprintf("%s-outbound-lb-outboundBackendPool", in.InfraID),
				})
				_, err = addLoadBalancingRuleToLoadBalancer(ctx, outboundApiRuleIPv4, lbInput)
				if err != nil {
					return fmt.Errorf("failed to add outbound ipv4 api rule to public load balancer: %w", err)
				}
			*/

			logrus.Debugf("XXX: Adding IPv6 rules to internal load balancer")

			// Add IPv6 rules to the internal load balancer
			apiRuleIPv6 := apiRule(&lbRuleInput{
				idPrefix:         idPrefix,
				loadBalancerName: publicLoadBalancerName,
				probeName:        *apiProbe.Name,
				ruleName:         "api-ipv6",
				//frontendIPConfigName:   fmt.Sprintf("%s-frontEnd-ipv6", in.InfraID),
				frontendIPConfigName:   "public-lb-ipv6",
				backendAddressPoolName: fmt.Sprintf("%s-ipv6", in.InfraID),
			})
			_, err = addLoadBalancingRuleToLoadBalancer(ctx, apiRuleIPv6, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add ipv6 api rule to public load balancer: %w", err)
			}

			/*
				outboundApiRuleIPv6 := apiRule(&lbRuleInput{
					idPrefix:               idPrefix,
					loadBalancerName:       publicLoadBalancerName,
					probeName:              *apiProbe.Name,
					ruleName:               "api-ipv6",
					frontendIPConfigName:   "public-lb-ipv6",
					backendAddressPoolName: fmt.Sprintf("%s-outbound-lb-outboundBackendPool-ipv6", in.InfraID),
				})
				_, err = addLoadBalancingRuleToLoadBalancer(ctx, outboundApiRuleIPv6, lbInput)
				if err != nil {
					return fmt.Errorf("failed to add outbound ipv6 api rule to public load balancer: %w", err)
				}
			*/

			outboundNatRuleIPv6 := outboundRule(&lbRuleInput{
				idPrefix:               idPrefix,
				loadBalancerName:       publicLoadBalancerName,
				ruleName:               "OutboundNATAllProtocols-ipv6",
				frontendIPConfigName:   fmt.Sprintf("%s-frontEnd-ipv6", in.InfraID),
				backendAddressPoolName: fmt.Sprintf("%s-outbound-lb-outboundBackendPool-ipv6", in.InfraID),
			})
			_, err = addOutboundRuleToLoadBalancer(ctx, outboundNatRuleIPv6, lbInput)
			if err != nil {
				return fmt.Errorf("failed to add nat ipv6 api rule to public load balancer: %w", err)
			}

			lbInput.loadBalancerName = publicLoadBalancerName
			loadBalancer, err = getLoadBalancer(ctx, lbInput)
			if err != nil {
				return fmt.Errorf("failed to get public load balancer: %w", err)
			}

			lbBaps = append(lbBaps, loadBalancer.Properties.BackendAddressPools...)

			// XXX: handle dual & single stack here TODO
			/*
				logrus.Debugf("XXX: updating API load balancer with IPv4 address")
				lbInput.frontendIPConfigName = aztypes.PublicFrontendIPv4ConfigName
				lbInput.frontendIPConfigName = fmt.Sprintf("%s-frontEnd", in.InfraID)
				lbInput.backendAddressPoolName = azure.InternalBackendAddressPoolIPv4Name
				lbInput.stackType = aztypes.StackTypeIPv4
				lbInput.pip = publicIPv4

				loadBalancer, err = updateOutboundLoadBalancerToAPILoadBalancer(ctx, lbInput)
				if err != nil {
					return fmt.Errorf("failed to update external load balancer: %w", err)
				}

				logrus.Debugf("XXX: updating API load balancer with IPv6 address")
				lbInput.frontendIPConfigName = aztypes.PublicFrontendIPv6ConfigName
				lbInput.backendAddressPoolName = azure.InternalBackendAddressPoolIPv6Name
				lbInput.outboundAddressPoolName = aztypes.OutboundBackendAddressPoolIPv6Name
				lbInput.stackType = aztypes.StackTypeIPv6
				lbInput.pip = publicIPv6

				loadBalancer, err = updateOutboundLoadBalancerToAPILoadBalancer(ctx, lbInput)
				if err != nil {
					return fmt.Errorf("failed to update external load balancer: %w", err)
				}

				lbBaps = append(lbBaps, loadBalancer.Properties.BackendAddressPools...)
			*/
		}

		logrus.Debugf("updated external load balancer: %s", *loadBalancer.ID)
		//lbBaps = loadBalancer.Properties.BackendAddressPools
		// XXX: this *should* be the same for IPv4 and IPv6
		extLBFQDNIPv4 = *publicIPv4.Properties.DNSSettings.Fqdn
		extLBFQDNIPv6 = *publicIPv6.Properties.DNSSettings.Fqdn
		logrus.Debugf("XXX: IPv4 FQDN=%s IPv6 FQDN=%s", *publicIPv4.Properties.DNSSettings.Fqdn, *publicIPv6.Properties.DNSSettings.Fqdn)
	}

	// Save context for other hooks
	p.ResourceGroupName = resourceGroupName
	p.StorageAccountName = storageAccountName
	p.StorageURL = storageURL
	p.StorageAccount = storageAccount
	p.StorageAccountKeys = storageAccountKeys
	p.StorageClientFactory = storageClientFactory
	p.NetworkClientFactory = networkClientFactory
	p.lbBackendAddressPools = lbBaps

	logrus.Debugf("XXX: publicIPv4=%s publicIPv6=%s", *publicIPv4.Properties.IPAddress, *publicIPv6.Properties.IPAddress)
	if err := createDNSEntries(ctx, &createDNSEntriesInput{
		infra:                in,
		extLBFQDNIPv4:        extLBFQDNIPv4,
		extLBFQDNIPv6:        extLBFQDNIPv6,
		publicIPv4:           *publicIPv4.Properties.IPAddress,
		publicIPv6:           *publicIPv6.Properties.IPAddress,
		resourceGroupName:    resourceGroupName,
		networkClientFactory: networkClientFactory,
		opts:                 p.clientOptions,
	}); err != nil {
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

	if in.InstallConfig.Config.PublicAPI() {
		vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, ssn.TokenCreds, p.computeClientOptions)
		if err != nil {
			return fmt.Errorf("error creating vm client: %w", err)
		}

		vmIDs, err := getControlPlaneIDs(in.Client, in.InstallConfig.Config.ControlPlane.Replicas, in.InfraID)
		if err != nil {
			return fmt.Errorf("failed to get control plane VM IDs: %w", err)
		}

		vmInput := vmInput{
			infraID:             fmt.Sprintf("%s-internal", in.InfraID),
			resourceGroupName:   p.ResourceGroupName,
			vmClient:            vmClient,
			nicClient:           p.NetworkClientFactory.NewInterfacesClient(),
			ids:                 vmIDs,
			backendAddressPools: p.lbBackendAddressPools,
		}

		publicBackendAddressPools := p.lbBackendAddressPools

		internalLB, _ := getLoadBalancer(ctx, &lbInput{
			loadBalancerName:     fmt.Sprintf("%s-internal", in.InfraID),
			resourceGroupName:    p.ResourceGroupName,
			networkClientFactory: p.NetworkClientFactory,
		})
		vmInput.backendAddressPools = internalLB.Properties.BackendAddressPools

		logrus.Debugf("XXX: associating to VM internal backend pools")
		if err = associateVMToBackendPool(ctx, vmInput); err != nil {
			return fmt.Errorf("failed to associate control plane VMs with internal load balancer: %w", err)
		}

		vmInput.backendAddressPools = publicBackendAddressPools
		logrus.Debugf("XXX: associating to VM public backend pools")
		if err = associateVMToBackendPool(ctx, vmInput); err != nil {
			return fmt.Errorf("failed to associate control plane VMs with public load balancer: %w", err)
		}

		loadBalancerName := in.InfraID
		installConfig := in.InstallConfig.Config
		lbFrontendClient := p.NetworkClientFactory.NewLoadBalancerFrontendIPConfigurationsClient()

		/*
			sshRuleName := fmt.Sprintf("%s_ssh_in", in.InfraID)
			if err = addSecurityGroupRule(ctx, &securityGroupInput{
				resourceGroupName:    p.ResourceGroupName,
				securityGroupName:    fmt.Sprintf("%s-nsg", in.InfraID),
				securityRuleName:     sshRuleName,
				securityRulePort:     "22",
				securityRulePriority: 220,
				networkClientFactory: p.NetworkClientFactory,
			}); err != nil {
				return fmt.Errorf("failed to add ipv4 security rule: %w", err)
			}
		*/

		// XXX: Get back to this XXX, we can indeed have 2 SSH rules for
		// each protocol
		if installConfig.Networking.IsDualStack() {
			logrus.Debugf("XXX: Adding Nat rule to dual stack configuration")

			interfacesClient := p.NetworkClientFactory.NewInterfacesClient()
			interfaceResp, err := interfacesClient.Get(ctx,
				p.ResourceGroupName,
				fmt.Sprintf("%s-bootstrap-nic", in.InfraID),
				nil,
			)
			if err != nil {
				return fmt.Errorf("failed to get bootstrap interface: %w", err)
			}
			bootstrapInterface := interfaceResp.Interface

			var ipv4Config, ipv6Config *armnetwork.InterfaceIPConfiguration
			for _, ipConfig := range bootstrapInterface.Properties.IPConfigurations {
				if *ipConfig.Properties.PrivateIPAddressVersion == armnetwork.IPVersionIPv4 {
					ipv4Config = ipConfig
				} else if *ipConfig.Properties.PrivateIPAddressVersion == armnetwork.IPVersionIPv6 {
					ipv6Config = ipConfig
				}
			}

			if ipv4Config == nil || ipv6Config == nil {
				return fmt.Errorf("failed to get IP configurations for bootstrap interface")
			}

			logrus.Debugf("XXX: Adding Nat rule to ipv4 load balancer")
			ipv4Frontend, err := lbFrontendClient.Get(ctx, p.ResourceGroupName, loadBalancerName, "public-lb-ipv4", nil)
			if err != nil {
				return fmt.Errorf("failed to get ipv4 frontend: %w", err)
			}

			natRuleName := fmt.Sprintf("%s_ssh_in_ipv4", in.InfraID)
			inboundNatRule, err := addInboundNatRuleToLoadBalancer(ctx, &inboundNatRuleInput{
				resourceGroupName:    p.ResourceGroupName,
				loadBalancerName:     loadBalancerName,
				frontendIPConfigID:   *ipv4Frontend.ID,
				backendIPConfigID:    *ipv4Config.ID,
				inboundNatRuleName:   natRuleName,
				inboundNatRulePort:   22,
				networkClientFactory: p.NetworkClientFactory,
			})
			if err != nil {
				return fmt.Errorf("failed to create inbound nat rule: %w", err)
			}
			ipv4Config.Properties.LoadBalancerInboundNatRules = append(ipv4Config.Properties.LoadBalancerInboundNatRules, inboundNatRule)

			/*
				_, err = associateInboundNatRuleToInterface(ctx, &inboundNatRuleInput{
					resourceGroupName:    p.ResourceGroupName,
					loadBalancerName:     loadBalancerName,
					bootstrapNicName:     fmt.Sprintf("%s-bootstrap-nic", in.InfraID),
					frontendIPConfigName: *ipv4Frontend.Name,
					inboundNatRuleID:     *inboundNatRule.ID,
					inboundNatRuleName:   natRuleName,
					inboundNatRulePort:   22,
					networkClientFactory: p.NetworkClientFactory,
				})
				if err != nil {
					return fmt.Errorf("failed to associate inbound nat rule to interface: %w", err)
				}
			*/

			logrus.Debugf("XXX: Adding Nat rule to ipv6 load balancer")
			ipv6Frontend, err := lbFrontendClient.Get(ctx, p.ResourceGroupName, loadBalancerName, "public-lb-ipv6", nil)
			if err != nil {
				return fmt.Errorf("failed to get ipv6 frontend: %w", err)
			}

			natRuleName = fmt.Sprintf("%s_ssh_in_ipv6", in.InfraID)
			inboundNatRule, err = addInboundNatRuleToLoadBalancer(ctx, &inboundNatRuleInput{
				resourceGroupName:    p.ResourceGroupName,
				loadBalancerName:     loadBalancerName,
				frontendIPConfigID:   *ipv6Frontend.ID,
				backendIPConfigID:    *ipv6Config.ID,
				inboundNatRuleName:   natRuleName,
				inboundNatRulePort:   22,
				networkClientFactory: p.NetworkClientFactory,
			})
			if err != nil {
				return fmt.Errorf("failed to create inbound nat rule: %w", err)
			}
			ipv6Config.Properties.LoadBalancerInboundNatRules = append(ipv6Config.Properties.LoadBalancerInboundNatRules, inboundNatRule)

			/*
				_, err = associateInboundNatRuleToInterface(ctx, &inboundNatRuleInput{
					resourceGroupName:    p.ResourceGroupName,
					loadBalancerName:     loadBalancerName,
					bootstrapNicName:     fmt.Sprintf("%s-bootstrap-nic", in.InfraID),
					frontendIPConfigName: *ipv6Frontend.Name,
					inboundNatRuleID:     *inboundNatRule.ID,
					inboundNatRuleName:   natRuleName,
					inboundNatRulePort:   22,
					networkClientFactory: p.NetworkClientFactory,
				})
				if err != nil {
					return fmt.Errorf("failed to associate inbound nat rule to interface: %w", err)
				}
			*/
			for i, ipConfig := range bootstrapInterface.Properties.IPConfigurations {
				if *ipConfig.Properties.PrivateIPAddressVersion == armnetwork.IPVersionIPv4 {
					bootstrapInterface.Properties.IPConfigurations[i] = ipv4Config

				} else if *ipConfig.Properties.PrivateIPAddressVersion == armnetwork.IPVersionIPv6 {
					bootstrapInterface.Properties.IPConfigurations[i] = ipv6Config
				}
			}

			interfacesPollerResp, err := interfacesClient.BeginCreateOrUpdate(ctx,
				p.ResourceGroupName,
				fmt.Sprintf("%s-bootstrap-nic", in.InfraID),
				bootstrapInterface,
				nil,
			)
			if err != nil {
				return fmt.Errorf("failed to add inbound nat rule to interface: %w", err)
			}

			_, err = interfacesPollerResp.PollUntilDone(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to add inbound nat rule to interface: %w", err)
			}

			securityRuleName := fmt.Sprintf("%s_ssh_in", in.InfraID)
			if err = addSecurityGroupRule(ctx, &securityGroupInput{
				resourceGroupName:    p.ResourceGroupName,
				securityGroupName:    fmt.Sprintf("%s-nsg", in.InfraID),
				securityRuleName:     securityRuleName,
				securityRulePort:     "22",
				securityRulePriority: 220,
				networkClientFactory: p.NetworkClientFactory,
			}); err != nil {
				return fmt.Errorf("failed to add ssh security rule: %w", err)
			}

		} else if installConfig.Networking.IsIPv4() {
			logrus.Debugf("XXX: Adding Nat rule to ipv4 configuration")

			sshRuleName := fmt.Sprintf("%s_ssh_in_ipv4", in.InfraID)
			if err = addSecurityGroupRule(ctx, &securityGroupInput{
				resourceGroupName:    p.ResourceGroupName,
				securityGroupName:    fmt.Sprintf("%s-nsg", in.InfraID),
				securityRuleName:     sshRuleName,
				securityRulePort:     "22",
				securityRulePriority: 224,
				networkClientFactory: p.NetworkClientFactory,
			}); err != nil {
				return fmt.Errorf("failed to add ipv4 security rule: %w", err)
			}

			publicFrontendIPConfigName := aztypes.PublicFrontendIPv4ConfigName
			publicFrontendIPConfigName = fmt.Sprintf("%s-frontEnd", in.InfraID)
			publicFrontendIPConfigID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/frontendIPConfigurations/%s",
				subscriptionID,
				p.ResourceGroupName,
				loadBalancerName,
				publicFrontendIPConfigName,
			)
			inboundNatRule, err := addInboundNatRuleToLoadBalancer(ctx, &inboundNatRuleInput{
				resourceGroupName:    p.ResourceGroupName,
				loadBalancerName:     loadBalancerName,
				frontendIPConfigID:   publicFrontendIPConfigID,
				inboundNatRuleName:   sshRuleName,
				inboundNatRulePort:   22,
				networkClientFactory: p.NetworkClientFactory,
			})
			if err != nil {
				return fmt.Errorf("failed to create inbound nat rule: %w", err)
			}
			_, err = associateInboundNatRuleToInterface(ctx, &inboundNatRuleInput{
				resourceGroupName:    p.ResourceGroupName,
				loadBalancerName:     loadBalancerName,
				bootstrapNicName:     fmt.Sprintf("%s-bootstrap-nic", in.InfraID),
				frontendIPConfigName: publicFrontendIPConfigName,
				inboundNatRuleID:     *inboundNatRule.ID,
				inboundNatRuleName:   sshRuleName,
				inboundNatRulePort:   22,
				networkClientFactory: p.NetworkClientFactory,
			})
			if err != nil {
				return fmt.Errorf("failed to associate inbound nat rule to interface: %w", err)
			}

		} else if installConfig.Networking.IsIPv6() {
			logrus.Debugf("XXX: Adding Nat rule to ipv6 configuration")

			sshRuleName := fmt.Sprintf("%s_ssh_in_ipv6", in.InfraID)
			if err = addSecurityGroupRule(ctx, &securityGroupInput{
				resourceGroupName:    p.ResourceGroupName,
				securityGroupName:    fmt.Sprintf("%s-nsg", in.InfraID),
				securityRuleName:     sshRuleName,
				securityRulePort:     "22",
				securityRulePriority: 226,
				networkClientFactory: p.NetworkClientFactory,
			}); err != nil {
				return fmt.Errorf("failed to add ipv4 security rule: %w", err)
			}

			//publicFrontendIPConfigName := fmt.Sprintf("%s-frontEnd-ipv6", in.InfraID)
			publicFrontendIPConfigName := "public-lb-ipv6"
			publicFrontendIPConfigID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/frontendIPConfigurations/%s",
				subscriptionID,
				p.ResourceGroupName,
				loadBalancerName,
				publicFrontendIPConfigName,
			)
			inboundNatRule, err := addInboundNatRuleToLoadBalancer(ctx, &inboundNatRuleInput{
				resourceGroupName:    p.ResourceGroupName,
				loadBalancerName:     loadBalancerName,
				frontendIPConfigID:   publicFrontendIPConfigID,
				inboundNatRuleName:   sshRuleName,
				inboundNatRulePort:   22,
				networkClientFactory: p.NetworkClientFactory,
			})
			if err != nil {
				return fmt.Errorf("failed to create inbound nat rule: %w", err)
			}
			_, err = associateInboundNatRuleToInterface(ctx, &inboundNatRuleInput{
				resourceGroupName:    p.ResourceGroupName,
				loadBalancerName:     loadBalancerName,
				bootstrapNicName:     fmt.Sprintf("%s-bootstrap-nic", in.InfraID),
				frontendIPConfigName: publicFrontendIPConfigName,
				inboundNatRuleID:     *inboundNatRule.ID,
				inboundNatRuleName:   sshRuleName,
				inboundNatRulePort:   22,
				networkClientFactory: p.NetworkClientFactory,
			})
			if err != nil {
				return fmt.Errorf("failed to associate inbound nat rule to interface: %w", err)
			}
		}

		/*
			sshRuleName := fmt.Sprintf("%s_ssh_in", in.InfraID)
			if err = addSecurityGroupRule(ctx, &securityGroupInput{
				resourceGroupName:    p.ResourceGroupName,
				securityGroupName:    fmt.Sprintf("%s-nsg", in.InfraID),
				securityRuleName:     sshRuleName,
				securityRulePort:     "22",
				securityRulePriority: 220,
				networkClientFactory: p.NetworkClientFactory,
			}); err != nil {
				return fmt.Errorf("failed to add security rule: %w", err)
			}
		*/

		/*
			publicFrontendIPConfigName := publicFrontendIPv4ConfigName
			publicFrontendIPConfigID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/frontendIPConfigurations/%s",
				subscriptionID,
				p.ResourceGroupName,
				loadBalancerName,
				publicFrontendIPConfigName,
			)
		*/

		// Create an inbound nat rule that forwards port 22 on the
		// public load balancer to the bootstrap host. This takes 2
		// stages to accomplish. First, the nat rule needs to be added
		// to the frontend IP configuration on the public load
		// balancer. Second, the nat rule needs to be addded to the
		// bootstrap interface with the association to the rule on the
		// public load balancer.
		/*
			inboundNatRule, err := addInboundNatRuleToLoadBalancer(ctx, &inboundNatRuleInput{
				resourceGroupName:    p.ResourceGroupName,
				loadBalancerName:     loadBalancerName,
				frontendIPConfigID:   publicFrontendIPConfigID,
				inboundNatRuleName:   sshRuleName,
				inboundNatRulePort:   22,
				networkClientFactory: p.NetworkClientFactory,
			})
			if err != nil {
				return fmt.Errorf("failed to create inbound nat rule: %w", err)
			}
			_, err = associateInboundNatRuleToInterface(ctx, &inboundNatRuleInput{
				resourceGroupName:    p.ResourceGroupName,
				loadBalancerName:     loadBalancerName,
				bootstrapNicName:     fmt.Sprintf("%s-bootstrap-nic", in.InfraID),
				frontendIPConfigName: publicFrontendIPConfigName,
				inboundNatRuleID:     *inboundNatRule.ID,
				inboundNatRuleName:   sshRuleName,
				inboundNatRulePort:   22,
				networkClientFactory: p.NetworkClientFactory,
			})
			if err != nil {
				return fmt.Errorf("failed to associate inbound nat rule to interface: %w", err)
			}
		*/
	}

	return nil
}

// PostDestroy removes SSH access from the network security rules and removes
// SSH port forwarding off the public load balancer when the bootstrap machine
// is destroyed.
func (p *Provider) PostDestroy(ctx context.Context, in clusterapi.PostDestroyerInput) error {
	session, err := azconfig.GetSession(in.Metadata.Azure.CloudName, in.Metadata.Azure.ARMEndpoint)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Construct client options here, rather than relying on p.clientOptions,
	// as PostDestroy can be called as part of destroy bootstrap, in which case
	// p.clientOption would not be populated.
	opts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: session.CloudConfig,
		},
	}
	if in.Metadata.Azure.CloudName == aztypes.StackCloud {
		opts.APIVersion = stackAPIVersion
	}

	networkClientFactory, err := armnetwork.NewClientFactory(
		session.Credentials.SubscriptionID,
		session.TokenCreds,
		opts,
	)
	if err != nil {
		return fmt.Errorf("error creating network client factory: %w", err)
	}

	resourceGroupName := fmt.Sprintf("%s-rg", in.Metadata.InfraID)
	if in.Metadata.Azure.ResourceGroupName != "" {
		resourceGroupName = in.Metadata.Azure.ResourceGroupName
	}
	securityGroupName := fmt.Sprintf("%s-nsg", in.Metadata.InfraID)
	sshRuleName := fmt.Sprintf("%s_ssh_in", in.Metadata.InfraID)

	// See if a security group rule exists with the name ${InfraID}_ssh_in.
	// If it does, this is a private cluster. If it does not, this is a
	// public cluster and we need to delete the SSH forward rule and
	// security group rule.
	_, err = networkClientFactory.NewSecurityRulesClient().Get(ctx,
		resourceGroupName,
		securityGroupName,
		sshRuleName,
		nil,
	)
	if err == nil {
		err = deleteSecurityGroupRule(ctx, &securityGroupInput{
			resourceGroupName:    resourceGroupName,
			securityGroupName:    securityGroupName,
			securityRuleName:     sshRuleName,
			securityRulePort:     "22",
			networkClientFactory: networkClientFactory,
		})
		if err != nil {
			return fmt.Errorf("failed to delete security rule: %w", err)
		}

		err = deleteInboundNatRule(ctx, &inboundNatRuleInput{
			resourceGroupName:    resourceGroupName,
			loadBalancerName:     in.Metadata.InfraID,
			inboundNatRuleName:   sshRuleName,
			networkClientFactory: networkClientFactory,
		})
		if err != nil {
			return fmt.Errorf("failed to delete inbound nat rule: %w", err)
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
func (p Provider) Ignition(ctx context.Context, in clusterapi.IgnitionInput) ([]*corev1.Secret, error) {
	session, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	bootstrapIgnData := in.BootstrapIgnData
	subscriptionID := session.Credentials.SubscriptionID

	ignitionContainerName := "ignition"
	blobName := "bootstrap.ign"
	blobURL := fmt.Sprintf("%s/%s/%s", p.StorageURL, ignitionContainerName, blobName)
	publicAccess := armstorage.PublicAccessNone
	// Create ignition blob storage container
	var blobIgnitionContainer *armstorage.BlobContainer
	if in.InstallConfig.Azure.CloudName != aztypes.StackCloud {
		createBlobContainerOutput, err := CreateBlobContainer(ctx, &CreateBlobContainerInput{
			ContainerName:        ignitionContainerName,
			SubscriptionID:       subscriptionID,
			ResourceGroupName:    p.ResourceGroupName,
			StorageAccountName:   p.StorageAccountName,
			PublicAccess:         to.Ptr(publicAccess),
			StorageClientFactory: p.StorageClientFactory,
		})
		if err != nil {
			return nil, err
		}
		blobIgnitionContainer = createBlobContainerOutput.BlobContainer
		logrus.Debugf("BlobIgnitionContainer.ID=%s", *blobIgnitionContainer.ID)
	}

	sasURL := ""

	if in.InstallConfig.Config.Azure.CustomerManagedKey == nil {
		logrus.Debugf("Creating a Block Blob for ignition shim")
		sasURL, err = CreateBlockBlob(ctx, &CreateBlockBlobInput{
			StorageURL:         p.StorageURL,
			BlobURL:            blobURL,
			StorageAccountName: p.StorageAccountName,
			StorageAccountKeys: p.StorageAccountKeys,
			ClientOpts:         p.clientOptions,
			BootstrapIgnData:   bootstrapIgnData,
			CloudEnvironment:   in.InstallConfig.Azure.CloudName,
			ContainerName:      ignitionContainerName,
			BlobName:           blobName,
			StorageSuffix:      session.Environment.StorageEndpointSuffix,
			ARMEndpoint:        in.InstallConfig.Azure.ARMEndpoint,
			Session:            session,
			Region:             in.InstallConfig.Config.Azure.Region,
			Tags:               p.Tags,
			ResourceGroupName:  p.ResourceGroupName,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create BlockBlob for ignition shim: %w", err)
		}
	} else {
		logrus.Debugf("Creating a Page Blob for ignition shim because Customer Managed Key is provided")
		lengthBootstrapFile := int64(len(bootstrapIgnData))
		if lengthBootstrapFile%512 != 0 {
			lengthBootstrapFile = (((lengthBootstrapFile / 512) + 1) * 512)
		}

		sasURL, err = CreatePageBlob(ctx, &CreatePageBlobInput{
			StorageURL:         p.StorageURL,
			BlobURL:            blobURL,
			ImageURL:           "",
			StorageAccountName: p.StorageAccountName,
			BootstrapIgnData:   bootstrapIgnData,
			ImageLength:        lengthBootstrapFile,
			StorageAccountKeys: p.StorageAccountKeys,
			ClientOpts:         p.clientOptions,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create PageBlob for ignition shim: %w", err)
		}
	}
	ignShim, err := bootstrap.GenerateIgnitionShimWithCertBundleAndProxy(sasURL, in.InstallConfig.Config.AdditionalTrustBundle, in.InstallConfig.Config.Proxy)
	if err != nil {
		return nil, fmt.Errorf("failed to create ignition shim: %w", err)
	}

	ignSecrets := []*corev1.Secret{
		clusterapi.IgnitionSecret(ignShim, in.InfraID, "bootstrap"),
		clusterapi.IgnitionSecret(in.MasterIgnData, in.InfraID, "master"),
	}

	return ignSecrets, nil
}

func getMachinePoolSecurityType(in clusterapi.InfraReadyInput) (string, error) {
	var securityType aztypes.SecurityTypes
	if in.InstallConfig.Config.ControlPlane != nil && in.InstallConfig.Config.ControlPlane.Platform.Azure != nil {
		pool := in.InstallConfig.Config.ControlPlane.Platform.Azure
		if pool.Settings != nil {
			securityType = pool.Settings.SecurityType
		}
	}
	if securityType == "" && in.InstallConfig.Config.Compute != nil {
		for _, compute := range in.InstallConfig.Config.Compute {
			if compute.Platform.Azure != nil {
				pool := compute.Platform.Azure
				if pool.Settings != nil {
					securityType = pool.Settings.SecurityType
					break
				}
			}
		}
	}
	if securityType == "" && in.InstallConfig.Config.Platform.Azure.DefaultMachinePlatform != nil {
		pool := in.InstallConfig.Config.Platform.Azure.DefaultMachinePlatform
		if pool.Settings != nil {
			securityType = pool.Settings.SecurityType
		}
	}
	switch securityType {
	case aztypes.SecurityTypesTrustedLaunch:
		return trustedLaunchST, nil
	case aztypes.SecurityTypesConfidentialVM:
		return confidentialVMST, nil
	}
	return "", nil
}
