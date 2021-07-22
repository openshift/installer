package azurestackprivate

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/storage/mgmt/storage"
	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2016-04-01/dns"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	mainStorage "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// ArmClient contains the handles to all the specific Azure Resource Manager
// resource classes' respective clients.
type ArmClient struct {
	clientID                 string
	tenantID                 string
	subscriptionID           string
	terraformVersion         string
	usingServicePrincipal    bool
	environment              azure.Environment
	skipProviderRegistration bool

	StopContext context.Context

	// Authentication
	servicePrincipalsClient graphrbac.ServicePrincipalsClient

	// DNS
	dnsClient   dns.RecordSetsClient
	zonesClient dns.ZonesClient

	// Networking
	ifaceClient                  network.InterfacesClient
	localNetConnClient           network.LocalNetworkGatewaysClient
	secRuleClient                network.SecurityRulesClient
	vnetGatewayClient            network.VirtualNetworkGatewaysClient
	vnetGatewayConnectionsClient network.VirtualNetworkGatewayConnectionsClient

	// Resources
	providersClient resources.ProvidersClient
	resourcesClient resources.Client

	resourceGroupsClient resources.GroupsClient
	deploymentsClient    resources.DeploymentsClient

	// Compute
	availSetClient       compute.AvailabilitySetsClient
	diskClient           compute.DisksClient
	vmExtensionClient    compute.VirtualMachineExtensionsClient
	vmClient             compute.VirtualMachinesClient
	vmImageClient        compute.VirtualMachineImagesClient
	vmScaleSetClient     compute.VirtualMachineScaleSetsClient
	storageServiceClient storage.AccountsClient

	// Network
	vnetClient         network.VirtualNetworksClient
	secGroupClient     network.SecurityGroupsClient
	publicIPClient     network.PublicIPAddressesClient
	subnetClient       network.SubnetsClient
	loadBalancerClient network.LoadBalancersClient
	routesClient       network.RoutesClient
	routeTablesClient  network.RouteTablesClient
}

func (c *ArmClient) configureClient(client *autorest.Client, auth autorest.Authorizer) {
	setUserAgent(client, c.terraformVersion)
	client.Authorizer = auth
	client.Sender = sender.BuildSender("AzureStack")
	client.SkipResourceProviderRegistration = c.skipProviderRegistration
	client.PollingDuration = 60 * time.Minute
}

func setUserAgent(client *autorest.Client, tfVersion string) {
	tfUserAgent := httpclient.TerraformUserAgent(tfVersion)

	// if the user agent already has a value append the Terraform user agent string
	if curUserAgent := client.UserAgent; curUserAgent != "" {
		client.UserAgent = fmt.Sprintf("%s %s", curUserAgent, tfUserAgent)
	} else {
		client.UserAgent = tfUserAgent
	}

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		client.UserAgent = fmt.Sprintf("%s %s", client.UserAgent, azureAgent)
	}
}

// getArmClient is a helper method which returns a fully instantiated
// *ArmClient based on the Config's current settings.
func getArmClient(authCfg *authentication.Config, tfVersion string, skipProviderRegistration bool) (*ArmClient, error) {
	env, err := authentication.LoadEnvironmentFromUrl(authCfg.CustomResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}

	// client declarations:
	client := ArmClient{
		clientID:                 authCfg.ClientID,
		tenantID:                 authCfg.TenantID,
		subscriptionID:           authCfg.SubscriptionID,
		terraformVersion:         tfVersion,
		environment:              *env,
		usingServicePrincipal:    authCfg.AuthenticatedAsAServicePrincipal,
		skipProviderRegistration: skipProviderRegistration,
	}

	oauth, err := authCfg.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, err
	}

	sender := sender.BuildSender("AzureStack")

	// Resource Manager endpoints
	endpoint := env.ResourceManagerEndpoint

	// Instead of the same endpoint use token audience to get the correct token.
	auth, err := authCfg.GetAuthorizationToken(sender, oauth, env.TokenAudience)
	if err != nil {
		return nil, err
	}

	// Graph Endpoints
	graphEndpoint := env.GraphEndpoint
	graphAuth, err := authCfg.GetAuthorizationToken(sender, oauth, graphEndpoint)
	if err != nil {
		return nil, err
	}

	client.registerAuthentication(graphEndpoint, client.tenantID, graphAuth, sender)
	client.registerComputeClients(endpoint, client.subscriptionID, auth)
	client.registerDNSClients(endpoint, client.subscriptionID, auth)
	client.registerNetworkingClients(endpoint, client.subscriptionID, auth)
	client.registerResourcesClients(endpoint, client.subscriptionID, auth)
	client.registerStorageClients(endpoint, client.subscriptionID, auth)

	return &client, nil
}

func (c *ArmClient) registerAuthentication(graphEndpoint, tenantID string, graphAuth autorest.Authorizer, sender autorest.Sender) {
	servicePrincipalsClient := graphrbac.NewServicePrincipalsClientWithBaseURI(graphEndpoint, tenantID)
	setUserAgent(&servicePrincipalsClient.Client, c.terraformVersion)
	servicePrincipalsClient.Authorizer = graphAuth
	servicePrincipalsClient.Sender = sender
	servicePrincipalsClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.servicePrincipalsClient = servicePrincipalsClient
}

func (c *ArmClient) registerComputeClients(endpoint, subscriptionID string, auth autorest.Authorizer) {
	availabilitySetsClient := compute.NewAvailabilitySetsClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&availabilitySetsClient.Client, auth)
	c.availSetClient = availabilitySetsClient

	diskClient := compute.NewDisksClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&diskClient.Client, auth)
	c.diskClient = diskClient

	extensionsClient := compute.NewVirtualMachineExtensionsClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&extensionsClient.Client, auth)
	c.vmExtensionClient = extensionsClient

	scaleSetsClient := compute.NewVirtualMachineScaleSetsClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&scaleSetsClient.Client, auth)
	c.vmScaleSetClient = scaleSetsClient

	virtualMachinesClient := compute.NewVirtualMachinesClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&virtualMachinesClient.Client, auth)
	c.vmClient = virtualMachinesClient

	virtualMachineImagesClient := compute.NewVirtualMachineImagesClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&virtualMachineImagesClient.Client, auth)
	c.vmImageClient = virtualMachineImagesClient
}

func (c *ArmClient) registerDNSClients(endpoint, subscriptionID string, auth autorest.Authorizer) {
	dn := dns.NewRecordSetsClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&dn.Client, auth)
	c.dnsClient = dn

	zo := dns.NewZonesClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&zo.Client, auth)
	c.zonesClient = zo
}

func (c *ArmClient) registerNetworkingClients(endpoint, subscriptionID string, auth autorest.Authorizer) {
	interfacesClient := network.NewInterfacesClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&interfacesClient.Client, auth)
	c.ifaceClient = interfacesClient

	gatewaysClient := network.NewVirtualNetworkGatewaysClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&gatewaysClient.Client, auth)
	c.vnetGatewayClient = gatewaysClient

	gatewayConnectionsClient := network.NewVirtualNetworkGatewayConnectionsClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&gatewayConnectionsClient.Client, auth)
	c.vnetGatewayConnectionsClient = gatewayConnectionsClient

	localNetworkGatewaysClient := network.NewLocalNetworkGatewaysClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&localNetworkGatewaysClient.Client, auth)
	c.localNetConnClient = localNetworkGatewaysClient

	loadBalancersClient := network.NewLoadBalancersClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&loadBalancersClient.Client, auth)
	c.loadBalancerClient = loadBalancersClient

	networksClient := network.NewVirtualNetworksClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&networksClient.Client, auth)
	c.vnetClient = networksClient

	publicIPAddressesClient := network.NewPublicIPAddressesClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&publicIPAddressesClient.Client, auth)
	c.publicIPClient = publicIPAddressesClient

	securityGroupsClient := network.NewSecurityGroupsClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&securityGroupsClient.Client, auth)
	c.secGroupClient = securityGroupsClient

	securityRulesClient := network.NewSecurityRulesClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&securityRulesClient.Client, auth)
	c.secRuleClient = securityRulesClient

	subnetsClient := network.NewSubnetsClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&subnetsClient.Client, auth)
	c.subnetClient = subnetsClient

	routeTablesClient := network.NewRouteTablesClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&routeTablesClient.Client, auth)
	c.routeTablesClient = routeTablesClient

	routesClient := network.NewRoutesClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&routesClient.Client, auth)
	c.routesClient = routesClient
}

func (c *ArmClient) registerResourcesClients(endpoint, subscriptionID string, auth autorest.Authorizer) {
	resourcesClient := resources.NewClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&resourcesClient.Client, auth)
	c.resourcesClient = resourcesClient

	deploymentsClient := resources.NewDeploymentsClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&deploymentsClient.Client, auth)
	c.deploymentsClient = deploymentsClient

	resourceGroupsClient := resources.NewGroupsClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&resourceGroupsClient.Client, auth)
	c.resourceGroupsClient = resourceGroupsClient

	providersClient := resources.NewProvidersClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&providersClient.Client, auth)
	c.providersClient = providersClient
}

func (c *ArmClient) registerStorageClients(endpoint, subscriptionID string, auth autorest.Authorizer) {
	accountsClient := storage.NewAccountsClientWithBaseURI(endpoint, subscriptionID)
	c.configureClient(&accountsClient.Client, auth)
	c.storageServiceClient = accountsClient
}

var (
	storageKeyCacheMu sync.RWMutex
	storageKeyCache   = make(map[string]string)
)

func (c *ArmClient) getKeyForStorageAccount(ctx context.Context, resourceGroupName, storageAccountName string) (string, bool, error) {
	cacheIndex := resourceGroupName + "/" + storageAccountName
	storageKeyCacheMu.RLock()
	key, ok := storageKeyCache[cacheIndex]
	storageKeyCacheMu.RUnlock()

	if ok {
		return key, true, nil
	}

	storageKeyCacheMu.Lock()
	defer storageKeyCacheMu.Unlock()
	key, ok = storageKeyCache[cacheIndex]
	if !ok {
		accountKeys, err := c.storageServiceClient.ListKeys(ctx, resourceGroupName, storageAccountName)
		if utils.ResponseWasNotFound(accountKeys.Response) {
			return "", false, nil
		}
		if err != nil {
			// We assume this is a transient error rather than a 404 (which is caught above),  so assume the
			// account still exists.
			return "", true, fmt.Errorf("error retrieving keys for storage account %q: %s", storageAccountName, err)
		}

		if accountKeys.Keys == nil {
			return "", false, fmt.Errorf("nil key returned for storage account %q", storageAccountName)
		}

		keys := *accountKeys.Keys
		if len(keys) <= 0 {
			return "", false, fmt.Errorf("no keys returned for storage account %q", storageAccountName)
		}

		keyPtr := keys[0].Value
		if keyPtr == nil {
			return "", false, fmt.Errorf("the first key returned is nil for storage account %q", storageAccountName)
		}

		key = *keyPtr
		storageKeyCache[cacheIndex] = key
	}

	return key, true, nil
}

func (c *ArmClient) getBlobStorageClientForStorageAccount(ctx context.Context, resourceGroupName, storageAccountName string) (*mainStorage.BlobStorageClient, bool, error) {
	key, accountExists, err := c.getKeyForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, accountExists, err
	}
	if !accountExists {
		return nil, false, nil
	}

	storageClient, err := mainStorage.NewClient(storageAccountName, key, c.environment.StorageEndpointSuffix,
		"2016-05-31", true)
	if err != nil {
		return nil, true, fmt.Errorf("error creating storage client for storage account %q: %s", storageAccountName, err)
	}

	blobClient := storageClient.GetBlobService()
	return &blobClient, true, nil
}
