package azure

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	azres "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/resources/mgmt/resources"
	azsubs "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/resources/mgmt/subscriptions"
	aznetwork "github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
	azenc "github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	azmarketplace "github.com/Azure/azure-sdk-for-go/profiles/latest/marketplaceordering/mgmt/marketplaceordering"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azstorage "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/go-autorest/autorest/to"
	"k8s.io/apimachinery/pkg/util/sets"
)

//go:generate mockgen -source=./client.go -destination=mock/azureclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	GetVirtualNetwork(ctx context.Context, resourceGroupName, virtualNetwork string) (*aznetwork.VirtualNetwork, error)
	CheckIPAddressAvailability(ctx context.Context, resourceGroupName, virtualNetwork, ipAddr string) (*aznetwork.IPAddressAvailabilityResult, error)
	GetComputeSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subnet string) (*aznetwork.Subnet, error)
	GetControlPlaneSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subnet string) (*aznetwork.Subnet, error)
	ListLocations(ctx context.Context) (*[]azsubs.Location, error)
	GetResourcesProvider(ctx context.Context, resourceProviderNamespace string) (*azres.Provider, error)
	GetVirtualMachineSku(ctx context.Context, name, region string) (*azenc.ResourceSku, error)
	GetVirtualMachineFamily(ctx context.Context, name, region string) (string, error)
	GetDiskSkus(ctx context.Context, region string) ([]azenc.ResourceSku, error)
	GetGroup(ctx context.Context, groupName string) (*azres.Group, error)
	ListResourceIDsByGroup(ctx context.Context, groupName string) ([]string, error)
	GetStorageEndpointSuffix(ctx context.Context) (string, error)
	GetDiskEncryptionSet(ctx context.Context, subscriptionID, groupName string, diskEncryptionSetName string) (*azenc.DiskEncryptionSet, error)
	GetHyperVGenerationVersion(ctx context.Context, instanceType string, region string, imageHyperVGen string) (string, error)
	GetMarketplaceImage(ctx context.Context, region, publisher, offer, sku, version string) (azenc.VirtualMachineImage, error)
	AreMarketplaceImageTermsAccepted(ctx context.Context, publisher, offer, sku string) (bool, error)
	GetVMCapabilities(ctx context.Context, instanceType, region string) (map[string]string, error)
	GetAvailabilityZones(ctx context.Context, region string, instanceType string) ([]string, error)
	GetLocationInfo(ctx context.Context, region string, instanceType string) (*azenc.ResourceSkuLocationInfo, error)
	CheckIfExistsStorageAccount(ctx context.Context, resourceGroup, storageAccountName, region string) error
}

// Client makes calls to the Azure API.
type Client struct {
	ssn *Session
}

// NewClient initializes a client with a session.
func NewClient(ssn *Session) *Client {
	client := &Client{
		ssn: ssn,
	}
	return client
}

// GetVirtualNetwork gets an Azure virtual network by name
func (c *Client) GetVirtualNetwork(ctx context.Context, resourceGroupName, virtualNetwork string) (*aznetwork.VirtualNetwork, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	vnetClient, err := c.getVirtualNetworksClient(ctx)
	if err != nil {
		return nil, err
	}

	vnet, err := vnetClient.Get(ctx, resourceGroupName, virtualNetwork, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual network %s: %w", virtualNetwork, err)
	}

	return &vnet, nil
}

// CheckIPAddressAvailability checks availability of an IP address in an Azure virtual network.
func (c *Client) CheckIPAddressAvailability(ctx context.Context, resourceGroupName, virtualNetwork, ipAddr string) (*aznetwork.IPAddressAvailabilityResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	vnetClient, err := c.getVirtualNetworksClient(ctx)
	if err != nil {
		return nil, err
	}

	availability, err := vnetClient.CheckIPAddressAvailability(ctx, resourceGroupName, virtualNetwork, ipAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get azure ip availability: %w", err)
	}

	return &availability, nil
}

// getSubnet gets an Azure subnet by name
func (c *Client) getSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subNetwork string) (*aznetwork.Subnet, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	subnetsClient, err := c.getSubnetsClient(ctx)
	if err != nil {
		return nil, err
	}

	subnet, err := subnetsClient.Get(ctx, resourceGroupName, virtualNetwork, subNetwork, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get subnet %s: %w", subNetwork, err)
	}

	return &subnet, nil
}

// GetComputeSubnet gets the Azure compute subnet
func (c *Client) GetComputeSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subNetwork string) (*aznetwork.Subnet, error) {
	return c.getSubnet(ctx, resourceGroupName, virtualNetwork, subNetwork)
}

// GetControlPlaneSubnet gets the Azure control plane subnet
func (c *Client) GetControlPlaneSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subNetwork string) (*aznetwork.Subnet, error) {
	return c.getSubnet(ctx, resourceGroupName, virtualNetwork, subNetwork)
}

// getVnetsClient sets up a new client to retrieve vnets
func (c *Client) getVirtualNetworksClient(ctx context.Context) (*aznetwork.VirtualNetworksClient, error) {
	vnetsClient := aznetwork.NewVirtualNetworksClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	vnetsClient.Authorizer = c.ssn.Authorizer
	return &vnetsClient, nil
}

// GetStorageEndpointSuffix retrieves the StorageEndpointSuffix from the
// session environment
func (c *Client) GetStorageEndpointSuffix(ctx context.Context) (string, error) {
	return c.ssn.Environment.StorageEndpointSuffix, nil
}

// getSubnetsClient sets up a new client to retrieve a subnet
func (c *Client) getSubnetsClient(ctx context.Context) (*aznetwork.SubnetsClient, error) {
	subnetClient := aznetwork.NewSubnetsClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	subnetClient.Authorizer = c.ssn.Authorizer
	return &subnetClient, nil
}

// ListLocations lists the Azure regions dir the given subscription
func (c *Client) ListLocations(ctx context.Context) (*[]azsubs.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	subsClient, err := c.getSubscriptionsClient(ctx)
	if err != nil {
		return nil, err
	}

	locations, err := subsClient.ListLocations(ctx, c.ssn.Credentials.SubscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list locations: %w", err)
	}

	return locations.Value, nil
}

// getSubscriptionsClient sets up a new client to retrieve subscription data
func (c *Client) getSubscriptionsClient(ctx context.Context) (azsubs.Client, error) {
	client := azsubs.NewClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint)
	client.Authorizer = c.ssn.Authorizer
	return client, nil
}

// GetResourcesProvider gets the Azure resource provider
func (c *Client) GetResourcesProvider(ctx context.Context, resourceProviderNamespace string) (*azres.Provider, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	providersClient, err := c.getProvidersClient(ctx)
	if err != nil {
		return nil, err
	}

	provider, err := providersClient.Get(ctx, resourceProviderNamespace, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get resource provider %s: %w", resourceProviderNamespace, err)
	}

	return &provider, nil
}

// getProvidersClient sets up a new client to retrieve providers data
func (c *Client) getProvidersClient(ctx context.Context) (azres.ProvidersClient, error) {
	client := azres.NewProvidersClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	return client, nil
}

// GetDiskSkus returns all the disk SKU pages for a given region.
func (c *Client) GetDiskSkus(ctx context.Context, region string) ([]azenc.ResourceSku, error) {
	client := azenc.NewResourceSkusClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	// See https://issues.redhat.com/browse/OCPBUGS-29469 before changing this timeout
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	var sku []azenc.ResourceSku
	filter := fmt.Sprintf("location eq '%s'", region)
	// This has to be initialized outside the `for` because we need access to
	// `err`. If initialized in the loop and the API call fails right away,
	// `page.NotDone()` will return `false` and we'll never check for the error
	skuPage, err := client.List(ctx, filter, "false")
	if err != nil {
		return nil, fmt.Errorf("failed to list SKUs: %w", err)
	}
	for ; skuPage.NotDone(); err = skuPage.NextWithContext(ctx) {
		if err != nil {
			return nil, fmt.Errorf("error fetching SKU pages: %w", err)
		}
		for _, page := range skuPage.Values() {
			for _, diskRegion := range to.StringSlice(page.Locations) {
				if strings.EqualFold(diskRegion, region) {
					sku = append(sku, page)
				}
			}
		}
	}

	if len(sku) != 0 {
		return sku, nil
	}

	return nil, fmt.Errorf("no disks for specified subscription in region %s", region)
}

// GetGroup returns resource group for the groupName.
func (c *Client) GetGroup(ctx context.Context, groupName string) (*azres.Group, error) {
	client := azres.NewGroupsClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res, err := client.Get(ctx, groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource group: %w", err)
	}
	return &res, nil
}

// ListResourceIDsByGroup returns a list of resource IDs for resource group groupName.
func (c *Client) ListResourceIDsByGroup(ctx context.Context, groupName string) ([]string, error) {
	client := azres.NewClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var res []string
	resPage, err := client.ListByResourceGroup(ctx, groupName, "", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}
	for ; resPage.NotDone(); err = resPage.NextWithContext(ctx) {
		if err != nil {
			return nil, fmt.Errorf("error fetching resource pages: %w", err)
		}
		for _, page := range resPage.Values() {
			res = append(res, to.String(page.ID))
		}
	}
	return res, nil
}

// GetVirtualMachineSku retrieves the resource SKU of a specified virtual machine SKU in the specified region.
func (c *Client) GetVirtualMachineSku(ctx context.Context, name, region string) (*azenc.ResourceSku, error) {
	client := azenc.NewResourceSkusClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer

	// See https://issues.redhat.com/browse/OCPBUGS-29469 before chaging this timeout
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	filter := fmt.Sprintf("location eq '%s'", region)
	// This has to be initialized outside the `for` because we need access to
	// `err`. If initialized in the loop and the API call fails right away,
	// `page.NotDone()` will return `false` and we'll never check for the error
	page, err := client.List(ctx, filter, "false")
	if err != nil {
		return nil, fmt.Errorf("failed to list SKUs: %w", err)
	}
	for ; page.NotDone(); err = page.NextWithContext(ctx) {
		if err != nil {
			return nil, fmt.Errorf("error fetching SKU pages: %w", err)
		}
		for _, sku := range page.Values() {
			// Filter out resources that are not virtualMachines
			if !strings.EqualFold("virtualMachines", *sku.ResourceType) {
				continue
			}
			// Filter out resources that do not match the provided name
			if !strings.EqualFold(name, *sku.Name) {
				continue
			}
			// Return the resource from the provided region
			for _, location := range to.StringSlice(sku.Locations) {
				if strings.EqualFold(location, region) {
					return &sku, nil
				}
			}
		}
	}

	return nil, nil
}

// GetDiskEncryptionSet retrieves the specified disk encryption set.
func (c *Client) GetDiskEncryptionSet(ctx context.Context, subscriptionID, groupName, diskEncryptionSetName string) (*azenc.DiskEncryptionSet, error) {
	client := azenc.NewDiskEncryptionSetsClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = c.ssn.Authorizer
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	diskEncryptionSet, err := client.Get(ctx, groupName, diskEncryptionSetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk encryption set: %w", err)
	}

	return &diskEncryptionSet, nil
}

// GetVirtualMachineFamily retrieves the VM family of an instance type.
func (c *Client) GetVirtualMachineFamily(ctx context.Context, name, region string) (string, error) {
	typeMeta, err := c.GetVirtualMachineSku(ctx, name, region)
	if err != nil {
		return "", fmt.Errorf("error connecting to Azure client: %w", err)
	}
	if typeMeta == nil {
		return "", fmt.Errorf("not found in region %s", region)
	}
	if typeMeta.Family == nil {
		return "", fmt.Errorf("error getting resource family")
	}

	return to.String(typeMeta.Family), nil
}

// GetVMCapabilities retrieves the capabilities of an instant type in a specific region. Returns these values
// in a map with the capability name as the key and the corresponding value.
func (c *Client) GetVMCapabilities(ctx context.Context, instanceType, region string) (map[string]string, error) {
	typeMeta, err := c.GetVirtualMachineSku(ctx, instanceType, region)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Azure client: %w", err)
	}
	if typeMeta == nil {
		return nil, fmt.Errorf("not found in region %s", region)
	}
	capabilities := make(map[string]string)
	for _, capability := range *typeMeta.Capabilities {
		capabilities[to.String(capability.Name)] = to.String(capability.Value)
	}

	return capabilities, nil
}

// GetHyperVGenerationVersion gets the HyperVGeneration version for the given instance type and marketplace image version, if specified. Defaults to V2 if either V1 or V2
// available.
func (c *Client) GetHyperVGenerationVersion(ctx context.Context, instanceType string, region string, imageHyperVGen string) (version string, err error) {
	capabilities, err := c.GetVMCapabilities(ctx, instanceType, region)
	if err != nil {
		return "", err
	}

	return GetHyperVGenerationVersion(capabilities, imageHyperVGen)
}

// GetMarketplaceImage get the specified marketplace VM image.
func (c *Client) GetMarketplaceImage(ctx context.Context, region, publisher, offer, sku, version string) (azenc.VirtualMachineImage, error) {
	client := azenc.NewVirtualMachineImagesClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	image, err := client.Get(ctx, region, publisher, offer, sku, version)
	if err != nil {
		return image, fmt.Errorf("could not get marketplace image: %w", err)
	}
	return image, nil
}

// AreMarketplaceImageTermsAccepted tests whether the terms have been accepted for the specified marketplace VM image.
func (c *Client) AreMarketplaceImageTermsAccepted(ctx context.Context, publisher, offer, sku string) (bool, error) {
	client := azmarketplace.NewMarketplaceAgreementsClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	terms, err := client.Get(ctx, publisher, offer, sku)
	if err != nil {
		return false, err
	}

	if terms.AgreementProperties == nil {
		return false, errors.New("no agreement properties for image")
	}

	return terms.AgreementProperties.Accepted != nil && *terms.AgreementProperties.Accepted, nil
}

// GetAvailabilityZones retrieves a list of availability zones for the given region, and instance type.
func (c *Client) GetAvailabilityZones(ctx context.Context, region string, instanceType string) ([]string, error) {
	locationInfo, err := c.GetLocationInfo(ctx, region, instanceType)
	if err != nil {
		return nil, err
	}
	if locationInfo != nil {
		return to.StringSlice(locationInfo.Zones), nil
	}

	return nil, fmt.Errorf("error retrieving availability zones for %s in %s", instanceType, region)
}

// GetLocationInfo retrieves the location info associated with the instance type in region
func (c *Client) GetLocationInfo(ctx context.Context, region string, instanceType string) (*azenc.ResourceSkuLocationInfo, error) {
	client := azenc.NewResourceSkusClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer

	// Only supported filter atm is `location`
	filter := fmt.Sprintf("location eq '%s'", region)
	res, err := client.List(ctx, filter, "false")
	if err != nil {
		return nil, fmt.Errorf("failed to list SKUs: %w", err)
	}
	for ; res.NotDone(); err = res.NextWithContext(ctx) {
		if err != nil {
			return nil, err
		}

		for _, resSku := range res.Values() {
			if !strings.EqualFold(to.String(resSku.ResourceType), "virtualMachines") {
				continue
			}
			if strings.EqualFold(to.String(resSku.Name), instanceType) {
				for _, locationInfo := range *resSku.LocationInfo {
					return &locationInfo, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("location information not found for %s in %s", instanceType, region)
}

// CheckIfExistsStorageAccount checks if the storage account provided already exists for diagnostics
// purposes.
func (c *Client) CheckIfExistsStorageAccount(ctx context.Context, resourceGroup, storageAccountName, region string) error {
	accountClientOptions := arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			// NOTE: the api version must support AzureStack
			APIVersion: "2019-04-01",
			Cloud:      c.ssn.CloudConfig,
		},
	}
	storageClient, err := azstorage.NewAccountsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, &accountClientOptions)
	if err != nil {
		return err
	}
	resp, err := storageClient.GetProperties(ctx, resourceGroup, storageAccountName, nil)
	if err != nil {
		return err
	}
	if resp.Account.Name != nil && region != *resp.Account.Location {
		return fmt.Errorf("%s is an invalid location for storage account. must be in the same region as the cluster", *resp.Account.Location)
	}
	validSKUs := sets.NewString(string(azstorage.SKUNameStandardGRS), string(azstorage.SKUNameStandardRAGRS), string(azstorage.SKUNameStandardLRS))
	if resp.Account.SKU != nil && resp.Account.SKU.Name != nil && !validSKUs.Has(string(*resp.Account.SKU.Name)) {
		stringSKUs := validSKUs.List()
		return fmt.Errorf("%s is not supported, supported values are %s,%s,%s", string(*resp.Account.SKU.Name), stringSKUs[0], stringSKUs[1], stringSKUs[2])
	}
	return err
}
