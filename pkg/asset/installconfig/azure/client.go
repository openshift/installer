package azure

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplaceordering/armmarketplaceordering"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	azstorage "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"k8s.io/apimachinery/pkg/util/sets"
)

//go:generate mockgen -source=./client.go -destination=mock/azureclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	GetVirtualNetwork(ctx context.Context, resourceGroupName, virtualNetwork string) (*armnetwork.VirtualNetwork, error)
	CheckIPAddressAvailability(ctx context.Context, resourceGroupName, virtualNetwork, ipAddr string) (*armnetwork.IPAddressAvailabilityResult, error)
	GetComputeSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subnet string) (*armnetwork.Subnet, error)
	GetControlPlaneSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subnet string) (*armnetwork.Subnet, error)
	ListLocations(ctx context.Context) ([]*armsubscriptions.Location, error)
	GetResourcesProvider(ctx context.Context, resourceProviderNamespace string) (*armresources.Provider, error)
	GetVirtualMachineSku(ctx context.Context, name, region string) (*armcompute.ResourceSKU, error)
	GetVirtualMachineFamily(ctx context.Context, name, region string) (string, error)
	GetDiskSkus(ctx context.Context, region string) ([]*armcompute.ResourceSKU, error)
	GetGroup(ctx context.Context, groupName string) (*armresources.ResourceGroup, error)
	ListResourceIDsByGroup(ctx context.Context, groupName string) ([]string, error)
	GetStorageEndpointSuffix(ctx context.Context) (string, error)
	GetDiskEncryptionSet(ctx context.Context, subscriptionID, groupName string, diskEncryptionSetName string) (*armcompute.DiskEncryptionSet, error)
	GetMarketplaceImage(ctx context.Context, region, publisher, offer, sku, version string) (*armcompute.VirtualMachineImage, error)
	AreMarketplaceImageTermsAccepted(ctx context.Context, publisher, offer, sku string) (bool, error)
	GetVMCapabilities(ctx context.Context, instanceType, region string) (map[string]string, error)
	GetAvailabilityZones(ctx context.Context, region string, instanceType string) ([]string, error)
	GetLocationInfo(ctx context.Context, region string, instanceType string) (*armcompute.ResourceSKULocationInfo, error)
	CheckIfExistsStorageAccount(ctx context.Context, resourceGroup, storageAccountName, region string) error
	GetRegionAvailabilityZones(ctx context.Context, region string) ([]string, error)
	CheckSubnetNatgateway(ctx context.Context, resourceGroup, virtualNetwork, subnet string) (bool, error)
	GetUserAssignedIdentity(ctx context.Context, subscriptionID, resourceGroup, name string) error
}

// Deprecated: Use ClientConfig.ClientOptions(ServiceNetwork) which automatically applies the correct
// API version based on the cloud environment. This constant is retained for backward compatibility.
// APIVersion describes the version to use for Azure API calls that support both azure and azurestack.
const APIVersion = "2019-11-01"

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
func (c *Client) GetVirtualNetwork(ctx context.Context, resourceGroupName, virtualNetwork string) (*armnetwork.VirtualNetwork, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	vnetClient, err := c.getVirtualNetworksClient()
	if err != nil {
		return nil, err
	}

	resp, err := vnetClient.Get(ctx, resourceGroupName, virtualNetwork, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual network %s: %w", virtualNetwork, err)
	}

	return &resp.VirtualNetwork, nil
}

// CheckIPAddressAvailability checks availability of an IP address in an Azure virtual network.
func (c *Client) CheckIPAddressAvailability(ctx context.Context, resourceGroupName, virtualNetwork, ipAddr string) (*armnetwork.IPAddressAvailabilityResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	vnetClient, err := c.getVirtualNetworksClient()
	if err != nil {
		return nil, err
	}

	resp, err := vnetClient.CheckIPAddressAvailability(ctx, resourceGroupName, virtualNetwork, ipAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get azure ip availability: %w", err)
	}

	return &resp.IPAddressAvailabilityResult, nil
}

// getSubnet gets an Azure subnet by name
func (c *Client) getSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subNetwork string) (*armnetwork.Subnet, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	subnetsClient, err := c.getSubnetsClient()
	if err != nil {
		return nil, err
	}

	resp, err := subnetsClient.Get(ctx, resourceGroupName, virtualNetwork, subNetwork, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get subnet %s: %w", subNetwork, err)
	}

	return &resp.Subnet, nil
}

// GetComputeSubnet gets the Azure compute subnet
func (c *Client) GetComputeSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subNetwork string) (*armnetwork.Subnet, error) {
	return c.getSubnet(ctx, resourceGroupName, virtualNetwork, subNetwork)
}

// GetControlPlaneSubnet gets the Azure control plane subnet
func (c *Client) GetControlPlaneSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subNetwork string) (*armnetwork.Subnet, error) {
	return c.getSubnet(ctx, resourceGroupName, virtualNetwork, subNetwork)
}

// getVnetsClient sets up a new client to retrieve vnets
func (c *Client) getVirtualNetworksClient() (*armnetwork.VirtualNetworksClient, error) {
	return armnetwork.NewVirtualNetworksClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.ClientOptions(ServiceNetwork))
}

// GetStorageEndpointSuffix retrieves the StorageEndpointSuffix from the
// session environment
func (c *Client) GetStorageEndpointSuffix(ctx context.Context) (string, error) {
	return c.ssn.Environment.StorageEndpointSuffix, nil
}

// getSubnetsClient sets up a new client to retrieve a subnet
func (c *Client) getSubnetsClient() (*armnetwork.SubnetsClient, error) {
	return armnetwork.NewSubnetsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.ClientOptions(ServiceNetwork))
}

// ListLocations lists the Azure regions dir the given subscription
func (c *Client) ListLocations(ctx context.Context) ([]*armsubscriptions.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	subsClient, err := c.getSubscriptionsClient()
	if err != nil {
		return nil, err
	}

	var locations []*armsubscriptions.Location
	pager := subsClient.NewListLocationsPager(c.ssn.Credentials.SubscriptionID, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list locations: %w", err)
		}
		locations = append(locations, page.Value...)
	}

	return locations, nil
}

// getSubscriptionsClient sets up a new client to retrieve subscription data
func (c *Client) getSubscriptionsClient() (*armsubscriptions.Client, error) {
	return armsubscriptions.NewClient(c.ssn.TokenCreds, c.ssn.ClientConfig.DefaultClientOptions())
}

// GetResourcesProvider gets the Azure resource provider
func (c *Client) GetResourcesProvider(ctx context.Context, resourceProviderNamespace string) (*armresources.Provider, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	providersClient, err := c.getProvidersClient()
	if err != nil {
		return nil, err
	}

	resp, err := providersClient.Get(ctx, resourceProviderNamespace, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource provider %s: %w", resourceProviderNamespace, err)
	}

	return &resp.Provider, nil
}

// getProvidersClient sets up a new client to retrieve providers data
func (c *Client) getProvidersClient() (*armresources.ProvidersClient, error) {
	return armresources.NewProvidersClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.DefaultClientOptions())
}

// GetDiskSkus returns all the disk SKU pages for a given region.
func (c *Client) GetDiskSkus(ctx context.Context, region string) ([]*armcompute.ResourceSKU, error) {
	client, err := armcompute.NewResourceSKUsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.ClientOptions(ServiceCompute))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource SKUs client: %w", err)
	}
	// See https://issues.redhat.com/browse/OCPBUGS-29469 before changing this timeout
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	var skus []*armcompute.ResourceSKU
	filter := fmt.Sprintf("location eq '%s'", region)
	pager := client.NewListPager(&armcompute.ResourceSKUsClientListOptions{
		Filter:                   &filter,
		IncludeExtendedLocations: to.Ptr("false"),
	})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error fetching SKU pages: %w", err)
		}
		for _, sku := range page.Value {
			if sku.Locations != nil {
				for _, diskRegion := range sku.Locations {
					if diskRegion != nil && strings.EqualFold(*diskRegion, region) {
						skus = append(skus, sku)
						break
					}
				}
			}
		}
	}

	if len(skus) != 0 {
		return skus, nil
	}

	return nil, fmt.Errorf("no disks for specified subscription in region %s", region)
}

// GetGroup returns resource group for the groupName.
func (c *Client) GetGroup(ctx context.Context, groupName string) (*armresources.ResourceGroup, error) {
	client, err := armresources.NewResourceGroupsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.DefaultClientOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to create resource groups client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, groupName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource group: %w", err)
	}
	return &resp.ResourceGroup, nil
}

// ListResourceIDsByGroup returns a list of resource IDs for resource group groupName.
func (c *Client) ListResourceIDsByGroup(ctx context.Context, groupName string) ([]string, error) {
	client, err := armresources.NewClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.DefaultClientOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to create resources client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var res []string
	pager := client.NewListByResourceGroupPager(groupName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error fetching resource pages: %w", err)
		}
		for _, resource := range page.Value {
			if resource.ID != nil {
				res = append(res, *resource.ID)
			}
		}
	}
	return res, nil
}

// GetVirtualMachineSku retrieves the resource SKU of a specified virtual machine SKU in the specified region.
func (c *Client) GetVirtualMachineSku(ctx context.Context, name, region string) (*armcompute.ResourceSKU, error) {
	client, err := armcompute.NewResourceSKUsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.ClientOptions(ServiceCompute))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource SKUs client: %w", err)
	}

	// See https://issues.redhat.com/browse/OCPBUGS-29469 before changing this timeout
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	filter := fmt.Sprintf("location eq '%s'", region)
	pager := client.NewListPager(&armcompute.ResourceSKUsClientListOptions{
		Filter:                   &filter,
		IncludeExtendedLocations: to.Ptr("false"),
	})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error fetching SKU pages: %w", err)
		}
		for _, sku := range page.Value {
			// Filter out resources that are not virtualMachines
			if sku.ResourceType == nil || !strings.EqualFold("virtualMachines", *sku.ResourceType) {
				continue
			}
			// Filter out resources that do not match the provided name
			if sku.Name == nil || !strings.EqualFold(name, *sku.Name) {
				continue
			}
			// Return the resource from the provided region
			if sku.Locations != nil {
				for _, location := range sku.Locations {
					if location != nil && strings.EqualFold(*location, region) {
						return sku, nil
					}
				}
			}
		}
	}

	return nil, nil
}

// GetDiskEncryptionSet retrieves the specified disk encryption set.
func (c *Client) GetDiskEncryptionSet(ctx context.Context, subscriptionID, groupName, diskEncryptionSetName string) (*armcompute.DiskEncryptionSet, error) {
	if !strings.EqualFold(c.ssn.Credentials.SubscriptionID, subscriptionID) {
		return nil, fmt.Errorf("different subscription from resource group subscription. Azure does not support cross subscription encryption sets")
	}
	client, err := armcompute.NewDiskEncryptionSetsClient(subscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.ClientOptions(ServiceCompute))
	if err != nil {
		return nil, fmt.Errorf("failed to create disk encryption sets client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, groupName, diskEncryptionSetName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk encryption set: %w", err)
	}
	return &resp.DiskEncryptionSet, nil
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

	return *typeMeta.Family, nil
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
	if typeMeta.Capabilities != nil {
		for _, capability := range typeMeta.Capabilities {
			if capability.Name != nil && capability.Value != nil {
				capabilities[*capability.Name] = *capability.Value
			}
		}
	}

	return capabilities, nil
}

// GetMarketplaceImage get the specified marketplace VM image.
func (c *Client) GetMarketplaceImage(ctx context.Context, region, publisher, offer, sku, version string) (*armcompute.VirtualMachineImage, error) {
	client, err := armcompute.NewVirtualMachineImagesClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.ClientOptions(ServiceCompute))
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual machine images client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, region, publisher, offer, sku, version, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get marketplace image: %w", err)
	}
	return &resp.VirtualMachineImage, nil
}

// AreMarketplaceImageTermsAccepted tests whether the terms have been accepted for the specified marketplace VM image.
func (c *Client) AreMarketplaceImageTermsAccepted(ctx context.Context, publisher, offer, sku string) (bool, error) {
	client, err := armmarketplaceordering.NewMarketplaceAgreementsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.DefaultClientOptions())
	if err != nil {
		return false, fmt.Errorf("failed to create marketplace agreements client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, armmarketplaceordering.OfferTypeVirtualmachine, publisher, offer, sku, nil)
	if err != nil {
		return false, err
	}

	if resp.Properties == nil {
		return false, errors.New("no agreement properties for image")
	}

	return resp.Properties.Accepted != nil && *resp.Properties.Accepted, nil
}

// GetAvailabilityZones retrieves a list of availability zones for the given region, and instance type.
func (c *Client) GetAvailabilityZones(ctx context.Context, region string, instanceType string) ([]string, error) {
	locationInfo, err := c.GetLocationInfo(ctx, region, instanceType)
	if err != nil {
		return nil, err
	}
	if locationInfo != nil && locationInfo.Zones != nil {
		zones := make([]string, 0, len(locationInfo.Zones))
		for _, z := range locationInfo.Zones {
			if z != nil {
				zones = append(zones, *z)
			}
		}
		return zones, nil
	}

	return nil, fmt.Errorf("error retrieving availability zones for %s in %s", instanceType, region)
}

// GetLocationInfo retrieves the location info associated with the instance type in region
func (c *Client) GetLocationInfo(ctx context.Context, region string, instanceType string) (*armcompute.ResourceSKULocationInfo, error) {
	client, err := armcompute.NewResourceSKUsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.ClientOptions(ServiceCompute))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource SKUs client: %w", err)
	}

	// Only supported filter atm is `location`
	filter := fmt.Sprintf("location eq '%s'", region)
	pager := client.NewListPager(&armcompute.ResourceSKUsClientListOptions{
		Filter:                   &filter,
		IncludeExtendedLocations: to.Ptr("false"),
	})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, resSku := range page.Value {
			if resSku.ResourceType == nil || !strings.EqualFold(*resSku.ResourceType, "virtualMachines") {
				continue
			}
			if resSku.Name != nil && strings.EqualFold(*resSku.Name, instanceType) {
				if len(resSku.LocationInfo) > 0 {
					return resSku.LocationInfo[0], nil
				}
			}
		}
	}

	return nil, fmt.Errorf("location information not found for %s in %s", instanceType, region)
}

// CheckIfExistsStorageAccount checks if the storage account provided already exists for diagnostics
// purposes.
func (c *Client) CheckIfExistsStorageAccount(ctx context.Context, resourceGroup, storageAccountName, region string) error {
	storageClient, err := azstorage.NewAccountsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.ClientOptions(ServiceStorage))
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

// GetRegionAvailabilityZones checks if a given region has availabililty zones for the nat gateways to use.
func (c *Client) GetRegionAvailabilityZones(ctx context.Context, region string) ([]string, error) {
	providersClient, err := armresources.NewProvidersClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.ClientOptions(ServiceNetwork))
	if err != nil {
		return nil, fmt.Errorf("failed to create providers client: %w", err)
	}

	provider, err := providersClient.Get(ctx, "Microsoft.Network", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get Microsoft.Network provider: %w", err)
	}

	if provider.ResourceTypes == nil {
		return nil, fmt.Errorf("no resource types found in Microsoft.Network provider")
	}

	// Find natGateways resource type
	for _, rt := range provider.ResourceTypes {
		if rt.ResourceType == nil || *rt.ResourceType != "natGateways" {
			continue
		}
		if rt.ZoneMappings != nil {
			for _, zm := range rt.ZoneMappings {
				if zones := getZoneMappings(zm, region); zones != nil {
					return zones, nil
				}
			}
		}
		if rt.Locations != nil {
			for _, loc := range rt.Locations {
				if loc != nil && strings.EqualFold(*loc, region) {
					return nil, nil // NAT gateway available but no zones
				}
			}
		}
		return nil, fmt.Errorf("NAT gateway not available in region %s", region)
	}

	return nil, fmt.Errorf("natGateways resource type not found in Microsoft.Network provider")
}

func getZoneMappings(zm *armresources.ZoneMapping, region string) []string {
	if zm.Location == nil || len(zm.Zones) == 0 {
		return nil
	}
	if !strings.EqualFold(strings.ReplaceAll(strings.ToLower(*zm.Location), " ", ""), region) {
		return nil
	}
	zones := []string{}
	for _, zone := range zm.Zones {
		if zone != nil {
			zones = append(zones, *zone)
		}
	}
	if len(zones) == 0 {
		return nil
	}
	return zones
}

// CheckSubnetNatgateway checks if there is an existing NAT gateway in a subnet.
func (c *Client) CheckSubnetNatgateway(ctx context.Context, resourceGroup, virtualNetwork, subnet string) (bool, error) {
	clientFactory, err := armnetwork.NewClientFactory(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, c.ssn.ClientConfig.ClientOptions(ServiceNetwork))
	if err != nil {
		return false, fmt.Errorf("failed to create client factory: %w", err)
	}

	res, err := clientFactory.NewSubnetsClient().Get(
		ctx,
		resourceGroup,
		virtualNetwork,
		subnet,
		&armnetwork.SubnetsClientGetOptions{Expand: nil},
	)
	if err != nil {
		return false, fmt.Errorf("failed to get subnet %s: %w", subnet, err)
	}

	if res.Subnet.Properties != nil {
		return res.Subnet.Properties.NatGateway != nil, nil
	}
	return false, fmt.Errorf("unable to get subnet nat gateway")
}

// GetUserAssignedIdentity checks if a user-assigned identity exists in the specified resource group.
func (c *Client) GetUserAssignedIdentity(ctx context.Context, subscriptionID, resourceGroup, name string) error {
	// Use the subscription ID from the function parameter if provided, otherwise use session default
	subID := subscriptionID
	if subID == "" {
		subID = c.ssn.Credentials.SubscriptionID
	}

	// Don't override APIVersion for managed identities - let SDK use the default
	// API version which supports user-assigned identities.
	client, err := armmsi.NewUserAssignedIdentitiesClient(subID, c.ssn.TokenCreds, c.ssn.ClientConfig.DefaultClientOptions())
	if err != nil {
		return fmt.Errorf("failed to create user-assigned identities client: %w", err)
	}

	_, err = client.Get(ctx, resourceGroup, name, nil)
	return err
}
