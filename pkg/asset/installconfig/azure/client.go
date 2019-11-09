package azure

import (
	"context"
	"time"

	"github.com/pkg/errors"

	aznetwork "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
)

//go:generate mockgen -source=./client.go -destination=mock/azureclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	GetVirtualNetwork(ctx context.Context, resourceGroupName, virtualNetwork string) (*aznetwork.VirtualNetwork, error)
	GetComputeSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subnet string) (*aznetwork.Subnet, error)
	GetControlPlaneSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subnet string) (*aznetwork.Subnet, error)
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
