package azureprivatedns

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/terraform-providers/terraform-provider-azurerm/version"
)

// ArmClient contains the handles to all the specific Azure Resource Manager
// resource classes' respective clients.
type ArmClient struct {
	clientID              string
	tenantID              string
	subscriptionID        string
	usingServicePrincipal bool
	environment           azure.Environment

	StopContext context.Context

	recordSetsClient          *privatedns.RecordSetsClient
	privateZonesClient        *privatedns.PrivateZonesClient
	virtualNetworkLinksClient *privatedns.VirtualNetworkLinksClient
}

func (c *ArmClient) configureClient(client *autorest.Client, auth autorest.Authorizer) {
	setUserAgent(client)
	client.Authorizer = auth
	client.Sender = sender.BuildSender("terraform-provider-azurerm-client-sender")
	client.SkipResourceProviderRegistration = true
	client.PollingDuration = 60 * time.Minute
}

func setUserAgent(client *autorest.Client) {
	// TODO: This is the SDK version not the CLI version, once we are on 0.12, should revisit
	tfUserAgent := httpclient.TerraformUserAgent("azure-private-dns")

	pv := version.ProviderVersion
	providerUserAgent := fmt.Sprintf("%s terraform-provider-azurerm/%s", tfUserAgent, pv)
	client.UserAgent = strings.TrimSpace(fmt.Sprintf("%s %s", client.UserAgent, providerUserAgent))

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		client.UserAgent = fmt.Sprintf("%s %s", client.UserAgent, azureAgent)
	}

	log.Printf("[DEBUG] AzureRM Client User Agent: %s\n", client.UserAgent)
}

// getArmClient is a helper method which returns a fully instantiated
// *ArmClient based on the Config's current settings.
func getArmClient(c *authentication.Config) (*ArmClient, error) {
	env, err := authentication.DetermineEnvironment(c.Environment)
	if err != nil {
		return nil, err
	}

	// client declarations:
	client := ArmClient{
		clientID:              c.ClientID,
		tenantID:              c.TenantID,
		subscriptionID:        c.SubscriptionID,
		environment:           *env,
		usingServicePrincipal: c.AuthenticatedAsAServicePrincipal,
	}

	adalOAuthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, c.TenantID)
	if err != nil {
		return nil, err
	}

	newOAuthConfig := &authentication.OAuthConfig{
		OAuth: adalOAuthConfig,
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if newOAuthConfig.OAuth == nil {
		return nil, fmt.Errorf("unable to configure OAuthConfig for tenant %s", c.TenantID)
	}

	sender := sender.BuildSender("terraform-provider-azurerm-sender")

	// Resource Manager endpoints
	endpoint := env.ResourceManagerEndpoint
	auth, err := c.GetAuthorizationToken(sender, newOAuthConfig, env.TokenAudience)
	if err != nil {
		return nil, err
	}

	client.registerPrivateDNSClients(endpoint, c.SubscriptionID, auth)

	return &client, nil
}

func (c *ArmClient) registerPrivateDNSClients(endpoint, subscriptionID string, auth autorest.Authorizer) {
	rs := privatedns.NewRecordSetsClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&rs.Client, auth)
	c.recordSetsClient = &rs

	zo := privatedns.NewPrivateZonesClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&zo.Client, auth)
	c.privateZonesClient = &zo

	vnl := privatedns.NewVirtualNetworkLinksClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&vnl.Client, auth)
	c.virtualNetworkLinksClient = &vnl
}
