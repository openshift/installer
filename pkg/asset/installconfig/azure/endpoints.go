package azure

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/azure"
	//"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	//"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// httpClient is the shared HTTP client used for all requests. It has a
// timeout so requests cannot hang indefinitely.
var httpClient = &http.Client{Timeout: 30 * time.Second}

// ErrorResponse is the top-level wrapper for the API error payload.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail holds the actual error code and message.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

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

// ServicePrincipal holds the Azure credentials read from the service
// principal JSON file.
type ServicePrincipal struct {
	SubscriptionId string `json:"subscriptionId"`
	ClientId       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`
	TenantId       string `json:"tenantId"`
}

// accessTokenResponse holds the OAuth2 token response returned by the Azure
// login endpoint.
type accessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	Resource     string `json:"resource"`
	TokenType    string `json:"token_type"`
	NotBefore    string `json:"not_before"`
	ExpiresIn    string `json:"expires_in"`
	ExtExpiresIn string `json:"ext_expires_in"`
	ExpiresOn    string `json:"expires_on"`
}

// loadServicePrincipal reads and parses the service principal file at path.
func loadServicePrincipal(path string) (ServicePrincipal, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return ServicePrincipal{}, fmt.Errorf("error opening %s: %w", path, err)
	}

	var sp ServicePrincipal
	if err := json.Unmarshal(content, &sp); err != nil {
		return ServicePrincipal{}, fmt.Errorf("error during unmarshaling: %w", err)
	}

	return sp, nil
}

// requestAccessToken exchanges the service principal credentials for an OAuth2
// access token and returns the token together with the resource URL the token
// is scoped to.
func requestAccessToken(armEndpoint string, credentials *Credentials) (accessToken string, err error) {
	authURL := fmt.Sprintf("%s/%s/oauth2/token", credentials.TenantID)

	postData := url.Values{}
	postData.Set("grant_type", "client_credentials")
	postData.Set("client_id", credentials.ClientID)

	// XXX: TODO
	switch {
	case credentials.ClientCertificatePath != "":
		//authType = ClientCertificateAuth
		postData.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	case credentials.ClientSecret != "":
		//authType = ClientSecretAuth
		postData.Set("client_secret", credentials.ClientSecret)
	default:
		//authType = ManagedIdentityAuth
	}

	postData.Set("resource", armEndpoint)

	req, err := http.NewRequest(http.MethodPost, authURL, strings.NewReader(postData.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting response: %s", resp.Status)
	}

	var data accessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("error while decoding response: %w", err)
	}

	return data.AccessToken, nil
}

func getAPIVersions(armEndpoint, subscriptionId, accessToken string) ([]string, error) {
	endpointsURL, err := url.Parse(fmt.Sprintf("%s/subscriptions/%s", armEndpoint, subscriptionId))
	if err != nil {
		return nil, fmt.Errorf("error parsing resource URL: %w", err)
	}

	params := endpointsURL.Query()
	params.Add("api-version", "Invalid")
	endpointsURL.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, endpointsURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	var apiVersions []string
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response: %w", err)
		}

		var errResponse ErrorResponse
		if err := json.Unmarshal(body, &errResponse); err != nil {
			return nil, fmt.Errorf("error unmarshalling error: %w", err)
		}

		re := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}(?:-[a-zA-Z]+)*")
		matches := re.FindAllString(errResponse.Error.Message, -1)
		for _, match := range matches {
			apiVersions = append(apiVersions, match)
		}

		count := len(apiVersions)
		if count > 0 {
			sortAPIVersions(apiVersions)
			return apiVersions, nil
		}

	}

	return nil, fmt.Errorf("error getting API versions")
}

// splitAPIVersion splits an api-version into its YYYY-MM-DD date prefix and any
// trailing suffix (e.g. "2021-01-01-preview" -> "2021-01-01", "-preview").
func splitAPIVersion(v string) (date, suffix string) {
	if len(v) >= 10 {
		return v[:10], v[10:]
	}
	return v, ""
}

// sortAPIVersions sorts api-version strings in place, newest first. Each entry
// has the form YYYY-MM-DD with an optional -suffix (e.g. -preview). For the
// same date, a stable version (no suffix) is ordered ahead of a suffixed one.
func sortAPIVersions(versions []string) {
	sort.Slice(versions, func(i, j int) bool {
		di, si := splitAPIVersion(versions[i])
		dj, sj := splitAPIVersion(versions[j])
		if di != dj {
			return di > dj // newer date first
		}
		if si == "" || sj == "" {
			return si == "" // stable (no suffix) before suffixed
		}
		return si < sj // stable ordering for two suffixed versions
	})
}

// latestStableAPIVersion returns the newest stable (non-preview) api-version
// from versions — the entry with the latest YYYY-MM-DD date and no suffix.
// It reports false if versions contains no stable entry.
func latestStableAPIVersion(versions []string) (string, bool) {
	var latest string
	for _, v := range versions {
		if _, suffix := splitAPIVersion(v); suffix != "" {
			continue // skip preview/beta/etc.
		}
		if latest == "" || v > latest {
			latest = v
		}
	}
	return latest, latest != ""
}

// fetchEndpoints retrieves the metadata endpoints document for ,
// authenticated with the given access token.
func fetchEndpoints(armEndpoint, apiVersion, accessToken string) ([]byte, error) {
	endpointsURL, err := url.Parse(fmt.Sprintf("%s/metadata/endpoints", armEndpoint))
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

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status: ", resp.Status)
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
	r.Name = "AzureUSSecCloud"
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
		env.KeyVaultEndpoint = fmt.Sprintf("https://%s/", r.Suffixes.KeyVaultDNS)
	}
	if r.Suffixes.MHSMDNS != "" {
		env.ManagedHSMEndpoint = fmt.Sprintf("https://%s/", r.Suffixes.MHSMDNS)
	}

	return env
}

// EnvironmentFromURLAuthenticated gets endpoints from a mangement URL while authenticated.
// This allows the installer to use a non-global management URL for endpoints.
func EnvironmentFromURLAuthenticated(armEndpoint string, credentials *Credentials) (environment azure.Environment, err error) {
	accessToken, err := requestAccessToken(armEndpoint, credentials)
	if err != nil {
		return environment, fmt.Errorf("error getting access token: %w", err)
	}

	apiVersions, err := getAPIVersions(armEndpoint, credentials.SubscriptionID, accessToken)
	if err != nil {
		return environment, fmt.Errorf("error getting API versions: %w", err)
	}

	apiVersion, ok := latestStableAPIVersion(apiVersions)
	if !ok {
		return environment, fmt.Errorf("no stable api-version available")
	}

	body, err := fetchEndpoints(armEndpoint, apiVersion, accessToken)
	if err != nil {
		return environment, fmt.Errorf("error getting endpoints: %w", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return environment, fmt.Errorf("error unmarshalling endpoint response: %w", err)
	}

	return response.toEnvironment(), nil
}

// EnvironmentFromURL gets endpoints from a mangement URL.
// This allows the installer to use a non-global management URL for endpoints.
func EnvironmentFromURL(armEndpoint string, credentials *Credentials) (environment azure.Environment, err error) {
	// XXX: This works... should probably use an actual valid API version though
	body, err := fetchEndpoints(armEndpoint, "invalid", "")
	if err != nil {
		return environment, fmt.Errorf("error getting endpoints: %w", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return environment, fmt.Errorf("error unmarshalling endpoint response: %w", err)
	}

	return response.toEnvironment(), nil
}
