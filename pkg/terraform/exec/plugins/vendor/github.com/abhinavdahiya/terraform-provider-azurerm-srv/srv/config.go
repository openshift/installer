package srv

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	az "github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
)

// armClient contains the handles to all the specific Azure Resource Manager
// resource classes' respective clients.
type armClient struct {
	clientID              string
	tenantID              string
	subscriptionID        string
	usingServicePrincipal bool
	environment           az.Environment

	StopContext context.Context

	dnsClient   dns.RecordSetsClient
	zonesClient dns.ZonesClient
}

// getArmClient is a helper method which returns a fully instantiated
// *armClient based on the Config's current settings.
func getArmClient(c *authentication.Config) (*armClient, error) {
	env, err := authentication.DetermineEnvironment(c.Environment)
	if err != nil {
		return nil, err
	}

	// client declarations:
	client := armClient{
		clientID:              c.ClientID,
		tenantID:              c.TenantID,
		subscriptionID:        c.SubscriptionID,
		environment:           *env,
		usingServicePrincipal: c.AuthenticatedAsAServicePrincipal,
	}

	oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, c.TenantID)
	if err != nil {
		return nil, err
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauthConfig == nil {
		return nil, fmt.Errorf("Unable to configure OAuthConfig for tenant %s", c.TenantID)
	}

	// Resource Manager endpoints
	endpoint := env.ResourceManagerEndpoint
	auth, err := c.GetAuthorizationToken(oauthConfig, env.TokenAudience)
	if err != nil {
		return nil, err
	}

	client.registerDNSClients(endpoint, c.SubscriptionID, auth)
	return &client, nil
}

func (c *armClient) registerDNSClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	dn := dns.NewRecordSetsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&dn.Client, auth)
	c.dnsClient = dn

	zo := dns.NewZonesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&zo.Client, auth)
	c.zonesClient = zo
}

func (c *armClient) configureClient(client *autorest.Client, auth autorest.Authorizer) {
	client.Authorizer = auth
	client.Sender = buildSender()
	client.PollingDuration = 60 * time.Minute
}

func buildSender() autorest.Sender {
	return autorest.DecorateSender(&http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}, withRequestLogging())
}

func withRequestLogging() autorest.SendDecorator {
	return func(s autorest.Sender) autorest.Sender {
		return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
			// strip the authorization header prior to printing
			authHeaderName := "Authorization"
			auth := r.Header.Get(authHeaderName)
			if auth != "" {
				r.Header.Del(authHeaderName)
			}

			// dump request to wire format
			if dump, err := httputil.DumpRequestOut(r, true); err == nil {
				log.Printf("[DEBUG] AzureRM Request: \n%s\n", dump)
			} else {
				// fallback to basic message
				log.Printf("[DEBUG] AzureRM Request: %s to %s\n", r.Method, r.URL)
			}

			// add the auth header back
			if auth != "" {
				r.Header.Add(authHeaderName, auth)
			}

			resp, err := s.Do(r)
			if resp != nil {
				// dump response to wire format
				if dump, err2 := httputil.DumpResponse(resp, true); err2 == nil {
					log.Printf("[DEBUG] AzureRM Response for %s: \n%s\n", r.URL, dump)
				} else {
					// fallback to basic message
					log.Printf("[DEBUG] AzureRM Response: %s for %s\n", resp.Status, r.URL)
				}
			} else {
				log.Printf("[DEBUG] Request to %s completed with no response", r.URL)
			}
			return resp, err
		})
	}
}
