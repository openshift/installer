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
	stackAPIVersion = "2020-06-01"
)

// Provider implements Azure CAPI installation.
type Provider struct {
	ResourceGroupName     string
	StorageAccountName    string
	StorageURL            string
	StorageAccount        *armstorage.Account
	StorageClientFactory  *armstorage.ClientFactory
	StorageAccountKeys    []armstorage.AccountKey
	NetworkClientFactory  *armnetwork.ClientFactory
	lbBackendAddressPools []*armnetwork.BackendAddressPool
	CloudConfiguration    cloud.Configuration
	TokenCredential       azcore.TokenCredential
	Tags                  map[string]*string
	clientOptions         *arm.ClientOptions
}

var _ clusterapi.PreProvider = (*Provider)(nil)
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
	p.ResourceGroupName = platform.ClusterResourceGroupName(in.InfraID)

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
	if platform.CloudName == aztypes.StackCloud {
		opts.APIVersion = stackAPIVersion
	}
	p.clientOptions = opts

	// // Create resource group
	// resourcesClientFactory, err := armresources.NewClientFactory(
	// 	subscriptionID,
	// 	tokenCredential,
	// 	&arm.ClientOptions{
	// 		ClientOptions: policy.ClientOptions{
	// 			Cloud: cloudConfiguration,
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to get azure resource groups factory: %w", err)
	// }
	// resourceGroupsClient := resourcesClientFactory.NewResourceGroupsClient()
	// _, err = resourceGroupsClient.CreateOrUpdate(
	// 	ctx,
	// 	resourceGroupName,
	// 	armresources.ResourceGroup{
	// 		Location:  ptr.To(platform.Region),
	// 		ManagedBy: nil,
	// 		Tags:      tags,
	// 	},
	// 	nil,
	// )
	// if err != nil {
	// 	return fmt.Errorf("error creating resource group %s: %w", resourceGroupName, err)
	// }
	// resourceGroup, err := resourceGroupsClient.Get(ctx, resourceGroupName, nil)
	// if err != nil {
	// 	return fmt.Errorf("error getting resource group %s: %w", resourceGroupName, err)
	// }

	// logrus.Debugf("ResourceGroup.ID=%s", *resourceGroup.ID)
	// p.ResourceGroupName = resourceGroupName

	// // Create user assigned identity
	// userAssignedIdentityName := fmt.Sprintf("%s-identity", in.InfraID)
	// armmsiClientFactory, err := armmsi.NewClientFactory(
	// 	subscriptionID,
	// 	tokenCredential,
	// 	&arm.ClientOptions{
	// 		ClientOptions: policy.ClientOptions{
	// 			Cloud: cloudConfiguration,
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to create armmsi client: %w", err)
	// }
	// _, err = armmsiClientFactory.NewUserAssignedIdentitiesClient().CreateOrUpdate(
	// 	ctx,
	// 	resourceGroupName,
	// 	userAssignedIdentityName,
	// 	armmsi.Identity{
	// 		Location: ptr.To(platform.Region),
	// 		Tags:     tags,
	// 	},
	// 	nil,
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to create user assigned identity %s: %w", userAssignedIdentityName, err)
	// }
	// userAssignedIdentity, err := armmsiClientFactory.NewUserAssignedIdentitiesClient().Get(
	// 	ctx,
	// 	resourceGroupName,
	// 	userAssignedIdentityName,
	// 	nil,
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to get user assigned identity %s: %w", userAssignedIdentityName, err)
	// }
	// principalID := *userAssignedIdentity.Properties.PrincipalID

	// logrus.Debugf("UserAssignedIdentity.ID=%s", *userAssignedIdentity.ID)
	// logrus.Debugf("PrinciapalID=%s", principalID)

	// clientFactory, err := armauthorization.NewClientFactory(
	// 	subscriptionID,
	// 	tokenCredential,
	// 	&arm.ClientOptions{
	// 		ClientOptions: policy.ClientOptions{
	// 			Cloud: cloudConfiguration,
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to create armauthorization client: %w", err)
	// }

	// roleDefinitionsClient := clientFactory.NewRoleDefinitionsClient()

	// var contributor *armauthorization.RoleDefinition
	// roleDefinitionsPager := roleDefinitionsClient.NewListPager(*resourceGroup.ID, nil)
	// for roleDefinitionsPager.More() {
	// 	roleDefinitionsList, err := roleDefinitionsPager.NextPage(ctx)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to find any role definitions: %w", err)
	// 	}
	// 	for _, roleDefinition := range roleDefinitionsList.Value {
	// 		if *roleDefinition.Properties.RoleName == "Contributor" {
	// 			contributor = roleDefinition
	// 			break
	// 		}
	// 	}
	// }
	// if contributor == nil {
	// 	return fmt.Errorf("failed to find contributor definition")
	// }

	// roleAssignmentsClient := clientFactory.NewRoleAssignmentsClient()
	// scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionID, resourceGroupName)
	// roleAssignmentUUID := uuid.New().String()

	// // XXX: Azure doesn't like creating an identity and immediately
	// // creating a role assignment for the identity. There can be
	// // replication delays. So, retry every 10 seconds for a minute until
	// // the role assignment gets created.
	// //
	// // See https://aka.ms/docs-principaltype
	// for i := 0; i < retryCount; i++ {
	// 	_, err = roleAssignmentsClient.Create(ctx, scope, roleAssignmentUUID,
	// 		armauthorization.RoleAssignmentCreateParameters{
	// 			Properties: &armauthorization.RoleAssignmentProperties{
	// 				PrincipalID:      ptr.To(principalID),
	// 				RoleDefinitionID: contributor.ID,
	// 			},
	// 		},
	// 		nil,
	// 	)
	// 	if err == nil {
	// 		break
	// 	}
	// 	time.Sleep(retryTime)
	// }
	// if err != nil {
	// 	return fmt.Errorf("failed to create role assignment: %w", err)
	// }

	// Creating a dummy nsg for existing vnets installation to appease the ingress operator.
	if in.InstallConfig.Config.Azure.VirtualNetwork != "" {
		networkClientFactory, err := armnetwork.NewClientFactory(
			subscriptionID,
			tokenCredential,
			p.clientOptions,
		)
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
				Tags:     tags,
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

	var architecture armcompute.Architecture
	if installConfig.ControlPlane.Architecture == types.ArchitectureARM64 {
		architecture = armcompute.ArchitectureArm64
	} else {
		architecture = armcompute.ArchitectureX64
	}

	resourceGroupName := p.ResourceGroupName
	storageAccountName := getStorageAccountName(in.InfraID)
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

	userTags := platform.UserTags
	tags := make(map[string]*string, len(userTags)+1)
	tags[fmt.Sprintf("kubernetes.io_cluster.%s", in.InfraID)] = ptr.To("owned")
	for k, v := range userTags {
		tags[k] = ptr.To(v)
	}

	tokenCredential := session.TokenCreds
	storageURL := fmt.Sprintf("https://%s.blob.%s", storageAccountName, session.Environment.StorageEndpointSuffix)
	blobURL := fmt.Sprintf("%s/%s/%s", storageURL, containerName, blobName)

	// Create storage account
	createStorageAccountOutput, err := CreateStorageAccount(ctx, &CreateStorageAccountInput{
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

	storageAccount := createStorageAccountOutput.StorageAccount
	storageClientFactory := createStorageAccountOutput.StorageClientFactory
	storageAccountKeys := createStorageAccountOutput.StorageAccountKeys

	logrus.Debugf("StorageAccount.ID=%s", *storageAccount.ID)

	// Create blob storage container
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

	// Upload the image to the container
	if _, ok := os.LookupEnv("OPENSHIFT_INSTALL_SKIP_IMAGE_UPLOAD"); !ok {
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

	networkClientFactory, err := armnetwork.NewClientFactory(subscriptionID, session.TokenCreds, p.clientOptions)
	if err != nil {
		return fmt.Errorf("error creating network client factory: %w", err)
	}

	lbClient := networkClientFactory.NewLoadBalancersClient()
	lbInput := &lbInput{
		loadBalancerName:       fmt.Sprintf("%s-internal", in.InfraID),
		infraID:                in.InfraID,
		region:                 platform.Region,
		resourceGroup:          resourceGroupName,
		subscriptionID:         session.Credentials.SubscriptionID,
		frontendIPConfigName:   "public-lb-ip-v4",
		backendAddressPoolName: fmt.Sprintf("%s-internal", in.InfraID),
		idPrefix: fmt.Sprintf("subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers",
			session.Credentials.SubscriptionID,
			resourceGroupName,
		),
		lbClient: lbClient,
		tags:     p.Tags,
	}

	intLoadBalancer, err := updateInternalLoadBalancer(ctx, lbInput)
	if err != nil {
		return fmt.Errorf("failed to update internal load balancer: %w", err)
	}
	logrus.Debugf("updated internal load balancer: %s", *intLoadBalancer.ID)

	var lbBaps []*armnetwork.BackendAddressPool
	var extLBFQDN string
	if in.InstallConfig.Config.PublicAPI() {
		publicIP, err := createPublicIP(ctx, &pipInput{
			name:          fmt.Sprintf("%s-pip-v4", in.InfraID),
			infraID:       in.InfraID,
			region:        in.InstallConfig.Config.Azure.Region,
			resourceGroup: resourceGroupName,
			pipClient:     networkClientFactory.NewPublicIPAddressesClient(),
			tags:          p.Tags,
		})
		if err != nil {
			return fmt.Errorf("failed to create public ip: %w", err)
		}
		logrus.Debugf("created public ip: %s", *publicIP.ID)

		lbInput.loadBalancerName = in.InfraID
		lbInput.backendAddressPoolName = in.InfraID

		var loadBalancer *armnetwork.LoadBalancer
		if platform.OutboundType == aztypes.UserDefinedRoutingOutboundType {
			loadBalancer, err = createAPILoadBalancer(ctx, publicIP, lbInput)
			if err != nil {
				return fmt.Errorf("failed to create API load balancer: %w", err)
			}
		} else {
			loadBalancer, err = updateOutboundLoadBalancerToAPILoadBalancer(ctx, publicIP, lbInput)
			if err != nil {
				return fmt.Errorf("failed to update external load balancer: %w", err)
			}
		}

		logrus.Debugf("updated external load balancer: %s", *loadBalancer.ID)
		lbBaps = loadBalancer.Properties.BackendAddressPools
		extLBFQDN = *publicIP.Properties.DNSSettings.Fqdn
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

	if err := createDNSEntries(ctx, in, extLBFQDN, resourceGroupName, p.clientOptions); err != nil {
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
		vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, ssn.TokenCreds, p.clientOptions)
		if err != nil {
			return fmt.Errorf("error creating vm client: %w", err)
		}

		vmIDs, err := getControlPlaneIDs(in.Client, in.InstallConfig.Config.ControlPlane.Replicas, in.InfraID)
		if err != nil {
			return fmt.Errorf("failed to get control plane VM IDs: %w", err)
		}

		vmInput := &vmInput{
			infraID:             in.InfraID,
			resourceGroup:       p.ResourceGroupName,
			vmClient:            vmClient,
			nicClient:           p.NetworkClientFactory.NewInterfacesClient(),
			ids:                 vmIDs,
			backendAddressPools: p.lbBackendAddressPools,
		}

		if err = associateVMToBackendPool(ctx, *vmInput); err != nil {
			return fmt.Errorf("failed to associate control plane VMs with external load balancer: %w", err)
		}

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

		loadBalancerName := in.InfraID
		frontendIPConfigName := "public-lb-ip-v4"
		frontendIPConfigID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/frontendIPConfigurations/%s",
			subscriptionID,
			p.ResourceGroupName,
			loadBalancerName,
			frontendIPConfigName,
		)

		// Create an inbound nat rule that forwards port 22 on the
		// public load balancer to the bootstrap host. This takes 2
		// stages to accomplish. First, the nat rule needs to be added
		// to the frontend IP configuration on the public load
		// balancer. Second, the nat rule needs to be addded to the
		// bootstrap interface with the association to the rule on the
		// public load balancer.
		inboundNatRule, err := addInboundNatRuleToLoadBalancer(ctx, &inboundNatRuleInput{
			resourceGroupName:    p.ResourceGroupName,
			loadBalancerName:     loadBalancerName,
			frontendIPConfigID:   frontendIPConfigID,
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
			frontendIPConfigID:   frontendIPConfigID,
			inboundNatRuleID:     *inboundNatRule.ID,
			inboundNatRuleName:   sshRuleName,
			inboundNatRulePort:   22,
			networkClientFactory: p.NetworkClientFactory,
		})
		if err != nil {
			return fmt.Errorf("failed to associate inbound nat rule to interface: %w", err)
		}
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

// Storage account names can't be more than 24 characters.
func getStorageAccountName(infraID string) string {
	storageAccountNameMax := 24

	storageAccountName := strings.ReplaceAll(infraID, "-", "")
	if len(storageAccountName) > storageAccountNameMax-2 {
		storageAccountName = storageAccountName[:storageAccountNameMax-2]
	}
	storageAccountName = fmt.Sprintf("%ssa", storageAccountName)

	return storageAccountName
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

	blobIgnitionContainer := createBlobContainerOutput.BlobContainer
	logrus.Debugf("BlobIgnitionContainer.ID=%s", *blobIgnitionContainer.ID)

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
