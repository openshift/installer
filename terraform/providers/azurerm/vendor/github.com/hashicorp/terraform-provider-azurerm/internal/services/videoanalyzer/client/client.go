package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/videoanalyzer/2021-05-01-preview/edgemodules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/videoanalyzer/2021-05-01-preview/videoanalyzers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	EdgeModuleClient     *edgemodules.EdgeModulesClient
	VideoAnalyzersClient *videoanalyzers.VideoAnalyzersClient
}

func NewClient(o *common.ClientOptions) *Client {
	edgeModulesClient := edgemodules.NewEdgeModulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&edgeModulesClient.Client, o.ResourceManagerAuthorizer)

	videoAnalyzersClient := videoanalyzers.NewVideoAnalyzersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&videoAnalyzersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		EdgeModuleClient:     &edgeModulesClient,
		VideoAnalyzersClient: &videoAnalyzersClient,
	}
}
