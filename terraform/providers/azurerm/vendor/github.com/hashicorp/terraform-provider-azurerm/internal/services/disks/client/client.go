package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/diskpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/iscsitargets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DiskPoolsClient            *diskpools.DiskPoolsClient
	DisksPoolIscsiTargetClient *iscsitargets.IscsiTargetsClient
}

func NewClient(o *common.ClientOptions) *Client {
	diskPoolsClient := diskpools.NewDiskPoolsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&diskPoolsClient.Client, o.ResourceManagerAuthorizer)

	iscsiTargetClient := iscsitargets.NewIscsiTargetsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&iscsiTargetClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DiskPoolsClient:            &diskPoolsClient,
		DisksPoolIscsiTargetClient: &iscsiTargetClient,
	}
}
