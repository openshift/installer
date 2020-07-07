package client

import (
<<<<<<< HEAD
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
=======
	"github.com/Azure/azure-sdk-for-go/services/policyinsights/mgmt/2019-10-01/policyinsights"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-09-01/policy"
>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AssignmentsClient    *policy.AssignmentsClient
	DefinitionsClient    *policy.DefinitionsClient
	SetDefinitionsClient *policy.SetDefinitionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	assignmentsClient := policy.NewAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&assignmentsClient.Client, o.ResourceManagerAuthorizer)

	definitionsClient := policy.NewDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&definitionsClient.Client, o.ResourceManagerAuthorizer)

	setDefinitionsClient := policy.NewSetDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&setDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AssignmentsClient:    &assignmentsClient,
		DefinitionsClient:    &definitionsClient,
		SetDefinitionsClient: &setDefinitionsClient,
	}
}
