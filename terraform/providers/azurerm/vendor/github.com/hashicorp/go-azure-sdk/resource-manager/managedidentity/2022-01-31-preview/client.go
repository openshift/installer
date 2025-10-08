package v2022_01_31_preview

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2022-01-31-preview/managedidentities"
)

type Client struct {
	ManagedIdentities *managedidentities.ManagedIdentitiesClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	managedIdentitiesClient := managedidentities.NewManagedIdentitiesClientWithBaseURI(endpoint)
	configureAuthFunc(&managedIdentitiesClient.Client)

	return Client{
		ManagedIdentities: &managedIdentitiesClient,
	}
}
