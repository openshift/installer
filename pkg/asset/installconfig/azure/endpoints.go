package azure

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/azure"

	aztypes "github.com/openshift/installer/pkg/types/azure"
)

// httpClient is the shared HTTP client used for all requests. It has a
// timeout so requests cannot hang indefinitely.
var httpClient = &http.Client{Timeout: 30 * time.Second}

// Response represents the Azure cloud metadata endpoints document.
type Response struct {
	Portal                                string         `json:"portal"`
	Authentication                        Authentication `json:"authentication"`
	Media                                 string         `json:"media"`
	GraphAudience                         string         `json:"graphAudience"`
	Graph                                 string         `json:"graph"`
	Name                                  string         `json:"name"`
	Suffixes                              Suffixes       `json:"suffixes"`
	Batch                                 string         `json:"batch"`
	ResourceManager                       string         `json:"resourceManager"`
	VMImageAliasDoc                       string         `json:"vmImageAliasDoc"`
	ActiveDirectoryDataLake               string         `json:"activeDirectoryDataLake"`
	SQLManagement                         string         `json:"sqlManagement"`
	MicrosoftGraphResourceID              string         `json:"microsoftGraphResourceId"`
	AppInsightsResourceID                 string         `json:"appInsightsResourceId"`
	AppServiceResourceID                  string         `json:"appServiceResourceId"`
	AppInsightsTelemetryChannelResourceID string         `json:"appInsightsTelemetryChannelResourceId"`
	AttestationResourceID                 string         `json:"attestationResourceId"`
	SynapseAnalyticsResourceID            string         `json:"synapseAnalyticsResourceId"`
	LogAnalyticsResourceID                string         `json:"logAnalyticsResourceId"`
	OssrDbmsResourceID                    string         `json:"ossrDbmsResourceId"`
}

// Authentication holds the authentication configuration for the cloud.
type Authentication struct {
	LoginEndpoint    string   `json:"loginEndpoint"`
	Audiences        []string `json:"audiences"`
	Tenant           string   `json:"tenant"`
	IdentityProvider string   `json:"identityProvider"`
}

// Suffixes holds the DNS suffixes for various Azure services.
type Suffixes struct {
	AzureDataLakeStoreFileSystem        string `json:"azureDataLakeStoreFileSystem"`
	ACRLoginServer                      string `json:"acrLoginServer"`
	SQLServerHostname                   string `json:"sqlServerHostname"`
	AzureDataLakeAnalyticsCatalogAndJob string `json:"azureDataLakeAnalyticsCatalogAndJob"`
	KeyVaultDNS                         string `json:"keyVaultDns"`
	Storage                             string `json:"storage"`
	AzureFrontDoorEndpointSuffix        string `json:"azureFrontDoorEndpointSuffix"`
	StorageSyncEndpointSuffix           string `json:"storageSyncEndpointSuffix"`
	MHSMDNS                             string `json:"mhsmDns"`
	MySQLServerEndpoint                 string `json:"mysqlServerEndpoint"`
	PostgreSQLServerEndpoint            string `json:"postgresqlServerEndpoint"`
	MariaDBServerEndpoint               string `json:"mariadbServerEndpoint"`
	SynapseAnalytics                    string `json:"synapseAnalytics"`
	AttestationEndpoint                 string `json:"attestationEndpoint"`
}

// fetchEndpoints retrieves the metadata endpoints from the specified endPoint URL.
func fetchEndpoints(armEndpoint, apiVersion, accessToken string) ([]byte, error) {
	endpointsURL, err := url.Parse(fmt.Sprintf("%s/metadata/endpoints", strings.TrimSuffix(armEndpoint, "/")))
	if err != nil {
		return nil, fmt.Errorf("error parsing resource URL: %w", err)
	}

	params := endpointsURL.Query()
	params.Add("api-version", apiVersion)
	endpointsURL.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, endpointsURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	if accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req) //nolint:gosec // the ARM endpoint is a user-provided install-config value, not attacker-controlled input
	if err != nil {
		return nil, fmt.Errorf("error performing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return body, nil
}

// toEnvironment converts a metadata endpoints Response into an
// azure.Environment. Fields that have no counterpart in the
// /metadata/endpoints document are left at their zero value.
func (r Response) toEnvironment() azure.Environment {
	env := azure.Environment{
		// Identity and portals.
		Name:                    r.Name,
		ManagementPortalURL:     r.Portal,
		ActiveDirectoryEndpoint: r.Authentication.LoginEndpoint,

		// Core service endpoints.
		ResourceManagerEndpoint: r.ResourceManager,
		GraphEndpoint:           r.Graph,
		MicrosoftGraphEndpoint:  r.MicrosoftGraphResourceID,
		BatchManagementEndpoint: r.Batch,

		// DNS / endpoint suffixes.
		StorageEndpointSuffix:       r.Suffixes.Storage,
		SQLDatabaseDNSSuffix:        r.Suffixes.SQLServerHostname,
		KeyVaultDNSSuffix:           r.Suffixes.KeyVaultDNS,
		ManagedHSMDNSSuffix:         r.Suffixes.MHSMDNS,
		MariaDBDNSSuffix:            r.Suffixes.MariaDBServerEndpoint,
		MySQLDatabaseDNSSuffix:      r.Suffixes.MySQLServerEndpoint,
		PostgresqlDatabaseDNSSuffix: r.Suffixes.PostgreSQLServerEndpoint,
		ContainerRegistryDNSSuffix:  r.Suffixes.ACRLoginServer,
		SynapseEndpointSuffix:       r.Suffixes.SynapseAnalytics,
		DatalakeSuffix:              r.Suffixes.AzureDataLakeStoreFileSystem,

		// Token audience matches the ARM endpoint (as in PublicCloud).
		TokenAudience: r.ResourceManager,

		// Per-service resource identifiers.
		ResourceIdentifiers: azure.ResourceIdentifier{
			Graph:               r.GraphAudience,
			MicrosoftGraph:      r.MicrosoftGraphResourceID,
			Batch:               r.Batch,
			Datalake:            r.ActiveDirectoryDataLake,
			OperationalInsights: r.LogAnalyticsResourceID,
			OSSRDBMS:            r.OssrDbmsResourceID,
			Synapse:             r.SynapseAnalyticsResourceID,
		},
	}

	// Audiences[0] is the classic service-management audience
	// (https://management.core.windows.net/); guard the slice.
	if len(r.Authentication.Audiences) > 0 {
		env.ServiceManagementEndpoint = r.Authentication.Audiences[0]
	}

	// The metadata doc only provides DNS suffixes for Key Vault / HSM,
	// so synthesize the endpoints from them.
	if r.Suffixes.KeyVaultDNS != "" {
		env.KeyVaultEndpoint = fmt.Sprintf("https://%s/", strings.TrimPrefix(r.Suffixes.KeyVaultDNS, "."))
	}
	if r.Suffixes.MHSMDNS != "" {
		env.ManagedHSMEndpoint = fmt.Sprintf("https://%s/", strings.TrimPrefix(r.Suffixes.MHSMDNS, "."))
	}

	return env
}

// EnvironmentFromURL gets endpoints from a mangement URL.
// This allows the installer to use a non-global management URL for endpoints.
func EnvironmentFromURL(cloudName aztypes.CloudEnvironment, armEndpoint string, credentials *Credentials) (environment azure.Environment, err error) {
	body, err := fetchEndpoints(armEndpoint, "invalid", "")
	if err != nil {
		return environment, fmt.Errorf("error getting endpoints: %w", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return environment, fmt.Errorf("error unmarshalling endpoint response: %w", err)
	}

	env := response.toEnvironment()
	if env.ResourceManagerEndpoint == "" {
		env.ResourceManagerEndpoint = armEndpoint
	}
	// For sovereign Azure clouds (Public, Government, USSec) the ARM token
	// audience is the Resource Manager endpoint, as in the built-in SDK
	// environments. authentication.audiences[0] is the classic
	// service-management audience (https://management.core.windows.net/), so
	// it must not be used as the ARM token audience. Fall back to the resolved
	// Resource Manager endpoint when the metadata document omits
	// resourceManager.
	if env.TokenAudience == "" {
		env.TokenAudience = env.ResourceManagerEndpoint
	}
	env.Name = string(cloudName)

	return env, nil
}
