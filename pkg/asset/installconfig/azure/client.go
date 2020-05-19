package azure

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	azsku "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	aznetwork "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	azres "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	azsubs "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-06-01/subscriptions"
	"github.com/Azure/go-autorest/autorest/to"
)

//go:generate mockgen -source=./client.go -destination=mock/azureclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	GetVirtualNetwork(ctx context.Context, resourceGroupName, virtualNetwork string) (*aznetwork.VirtualNetwork, error)
	GetComputeSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subnet string) (*aznetwork.Subnet, error)
	GetControlPlaneSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subnet string) (*aznetwork.Subnet, error)
	ListLocations(ctx context.Context) (*[]azsubs.Location, error)
	GetResourcesProvider(ctx context.Context, resourceProviderNamespace string) (*azres.Provider, error)
	GetDiskSkus(ctx context.Context, region string) ([]azsku.ResourceSku, error)
}

// Client makes calls to the Azure API.
type Client struct {
	ssn *Session
}

// NewClient initializes a client with a session.
func NewClient(ctx context.Context) (*Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	ssn, err := GetSession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	client := &Client{
		ssn: ssn,
	}
	return client, nil
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
		return nil, errors.Wrapf(err, "failed to get virtual network %s", virtualNetwork)
	}

	return &vnet, nil
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
		return nil, errors.Wrapf(err, "failed to get subnet %s", subNetwork)
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
	vnetsClient := aznetwork.NewVirtualNetworksClient(c.ssn.Credentials.SubscriptionID)
	vnetsClient.Authorizer = c.ssn.Authorizer
	return &vnetsClient, nil
}

// getSubnetsClient sets up a new client to retrieve a subnet
func (c *Client) getSubnetsClient(ctx context.Context) (*aznetwork.SubnetsClient, error) {
	subnetClient := aznetwork.NewSubnetsClient(c.ssn.Credentials.SubscriptionID)
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
		return nil, errors.Wrapf(err, "failed to list locations")
	}

	return locations.Value, nil
}

// getSubscriptionsClient sets up a new client to retrieve subscription data
func (c *Client) getSubscriptionsClient(ctx context.Context) (azsubs.Client, error) {
	client := azsubs.NewClient()
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
		return nil, errors.Wrapf(err, "failed to get resource provider %s", resourceProviderNamespace)
	}

	return &provider, nil
}

// getProvidersClient sets up a new client to retrieve providers data
func (c *Client) getProvidersClient(ctx context.Context) (azres.ProvidersClient, error) {
	client := azres.NewProvidersClient(c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	return client, nil
}

// GetDiskSkus returns all the disk SKU pages for a given region.
func (c *Client) GetDiskSkus(ctx context.Context, region string) ([]azsku.ResourceSku, error) {
	client := azsku.NewResourceSkusClient(c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var sku []azsku.ResourceSku

	for skuPage, err := client.List(ctx); skuPage.NotDone(); err = skuPage.NextWithContext(ctx) {
		if err != nil {
			return nil, errors.Wrap(err, "error fetching SKU pages")
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

	return nil, errors.Errorf("no disks for specified subscription in region %s", region)
}
