package client

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2022-02-01/signalr"
	webpubsub_v2021_10_01 "github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2021-10-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SignalRClient   *signalr.SignalRClient
	WebPubSubClient *webpubsub_v2021_10_01.Client
}

func NewClient(o *common.ClientOptions) *Client {
	signalRClient := signalr.NewSignalRClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&signalRClient.Client, o.ResourceManagerAuthorizer)

	webPubSubClient := webpubsub_v2021_10_01.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	return &Client{
		SignalRClient:   &signalRClient,
		WebPubSubClient: &webPubSubClient,
	}
}
