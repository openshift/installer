package azure

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/Azure/go-autorest/autorest/azure"
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

// fetchEndpoints retrieves the metadata endpoints document for armEndpoint.
func fetchEndpoints(armEndpoint, apiVersion string) ([]byte, error) {
	base, err := url.Parse(armEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error parsing resource URL: %w", err)
	}
	if base.Scheme != "https" {
		return nil, fmt.Errorf("invalid ARM endpoint %q: scheme must be https", armEndpoint)
	}
	if base.Hostname() == "" {
		return nil, fmt.Errorf("invalid ARM endpoint %q: host cannot be empty", armEndpoint)
	}
	if base.User != nil {
		return nil, fmt.Errorf("invalid ARM endpoint %q: user information is not allowed", armEndpoint)
	}

	endpointsURL := &url.URL{
		Scheme:   "https",
		Host:     base.Host,
		Path:     path.Join(base.Path, "/metadata/endpoints"),
		RawQuery: url.Values{"api-version": []string{apiVersion}}.Encode(),
	}

	req, err := http.NewRequest(http.MethodGet, endpointsURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req) //nolint:gosec // G704: armEndpoint is validated above
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
// /metadata/endpoints document are left at their zero value, except
// ResourceManagerEndpoint, which falls back to armEndpoint (the URL the
// document was fetched from).
func (r Response) toEnvironment(armEndpoint string) *azure.Environment {
	env := azure.Environment{
		// Identity and portals.
		Name:                    r.Name,
		ManagementPortalURL:     r.Portal,
		ActiveDirectoryEndpoint: r.Authentication.LoginEndpoint,

		// Core service endpoints.
		ResourceManagerEndpoint: cmp.Or(r.ResourceManager, armEndpoint),
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

	// Audiences[0] is the classic service-management audience for the cloud.
	if len(r.Authentication.Audiences) > 0 {
		env.ServiceManagementEndpoint = r.Authentication.Audiences[0]
		env.TokenAudience = r.Authentication.Audiences[0]
	}

	// The metadata doc only provides DNS suffixes for Key Vault / HSM,
	// so synthesize the endpoints from them.
	if r.Suffixes.KeyVaultDNS != "" {
		env.KeyVaultEndpoint = fmt.Sprintf("https://%s/", r.Suffixes.KeyVaultDNS)
	}
	if r.Suffixes.MHSMDNS != "" {
		env.ManagedHSMEndpoint = fmt.Sprintf("https://%s/", r.Suffixes.MHSMDNS)
	}

	return &env
}

// EnvironmentFromURL builds an azure.Environment from the metadata endpoints
// document published by armEndpoint. This allows the installer to use a
// non-global management URL for endpoints.
//
// Multiple results are returned for some versions of the API. Microsoft does not
// specify a return order, but in testing every single API version available, the
// results (when multiple results are returned), are always the same:
// ["AzureCloud","AzureChinaCloud","AzureUSGovernment"]
// It is possible this could change in the future. This function is currently only used
// for AzureStackCloud and AzureUSSecCloud so it is assumed to be safe.
// In the case that multiple results are returned and the cloud name is not matched,
// an error is returnerd with the first environment result. It's up to the caller what
// to do with it.
func EnvironmentFromURL(armEndpoint, apiVersion, cloudName string) (environment *azure.Environment, err error) {
	if armEndpoint == "" {
		return nil, fmt.Errorf("armEndpoint is empty")
	}

	if apiVersion == "" {
		apiVersion = "2022-09-01"
	}

	body, err := fetchEndpoints(armEndpoint, apiVersion)
	if err != nil {
		return nil, fmt.Errorf("error getting endpoints: %w", err)
	}

	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("no endpoints returned")
	}

	var responses []Response
	switch trimmed[0] {
	case '{':
		// A single cloud is returned as a bare object rather than an array.
		var single Response
		if err := json.Unmarshal(trimmed, &single); err != nil {
			return nil, fmt.Errorf("error unmarshalling endpoint response: %w", err)
		}
		responses = []Response{single}
	case '[':
		if err := json.Unmarshal(trimmed, &responses); err != nil {
			return nil, fmt.Errorf("error unmarshalling endpoint response: %w", err)
		}
	default:
		err = errors.New("unexpected type")
	}

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling endpoint response: %w", err)
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("no endpoints returned")
	}

	for _, r := range responses {
		if cloudName == "" || cloudName == r.Name {
			return r.toEnvironment(armEndpoint), nil
		}
	}

	return responses[0].toEnvironment(armEndpoint), fmt.Errorf("cloudName not found")
}
