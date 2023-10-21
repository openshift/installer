package client

import (
	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"                // nolint: staticcheck
	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2021-05-01/mysqlflexibleservers" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/serverfailover"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient               *mysql.ConfigurationsClient
	DatabasesClient                    *mysql.DatabasesClient
	FirewallRulesClient                *mysql.FirewallRulesClient
	FlexibleDatabasesClient            *mysqlflexibleservers.DatabasesClient
	FlexibleServerConfigurationsClient *mysqlflexibleservers.ConfigurationsClient
	FlexibleServerClient               *servers.ServersClient
	FlexibleServerFailoverClient       *serverfailover.ServerFailoverClient
	FlexibleServerFirewallRulesClient  *mysqlflexibleservers.FirewallRulesClient
	ServersClient                      *mysql.ServersClient
	ServerKeysClient                   *mysql.ServerKeysClient
	ServerSecurityAlertPoliciesClient  *mysql.ServerSecurityAlertPoliciesClient
	VirtualNetworkRulesClient          *mysql.VirtualNetworkRulesClient
	ServerAdministratorsClient         *mysql.ServerAdministratorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ConfigurationsClient := mysql.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	DatabasesClient := mysql.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabasesClient.Client, o.ResourceManagerAuthorizer)

	FirewallRulesClient := mysql.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&FirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	flexibleDatabasesClient := mysqlflexibleservers.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&flexibleDatabasesClient.Client, o.ResourceManagerAuthorizer)

	flexibleServerClient := servers.NewServersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&flexibleServerClient.Client, o.ResourceManagerAuthorizer)

	flexibleServerFailoverClient := serverfailover.NewServerFailoverClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&flexibleServerFailoverClient.Client, o.ResourceManagerAuthorizer)

	flexibleServerFirewallRulesClient := mysqlflexibleservers.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&flexibleServerFirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	flexibleServerConfigurationsClient := mysqlflexibleservers.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&flexibleServerConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	ServersClient := mysql.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServersClient.Client, o.ResourceManagerAuthorizer)

	ServerKeysClient := mysql.NewServerKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServerKeysClient.Client, o.ResourceManagerAuthorizer)

	serverSecurityAlertPoliciesClient := mysql.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	VirtualNetworkRulesClient := mysql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VirtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	serverAdministratorsClient := mysql.NewServerAdministratorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverAdministratorsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient:               &ConfigurationsClient,
		DatabasesClient:                    &DatabasesClient,
		FirewallRulesClient:                &FirewallRulesClient,
		FlexibleDatabasesClient:            &flexibleDatabasesClient,
		FlexibleServerClient:               &flexibleServerClient,
		FlexibleServerFailoverClient:       &flexibleServerFailoverClient,
		FlexibleServerFirewallRulesClient:  &flexibleServerFirewallRulesClient,
		FlexibleServerConfigurationsClient: &flexibleServerConfigurationsClient,
		ServersClient:                      &ServersClient,
		ServerKeysClient:                   &ServerKeysClient,
		ServerSecurityAlertPoliciesClient:  &serverSecurityAlertPoliciesClient,
		VirtualNetworkRulesClient:          &VirtualNetworkRulesClient,
		ServerAdministratorsClient:         &serverAdministratorsClient,
	}
}
