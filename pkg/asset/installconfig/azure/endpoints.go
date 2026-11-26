package azure

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/azure"

	azuretypes "github.com/openshift/installer/pkg/types/azure"
)

// hybridEnvironmentName is the default environment name for Azure Stack Hub.
const hybridEnvironmentName = "HybridEnvironment"

// httpClient is the HTTP client used for metadata requests.
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

	// Legacy api-version=1.0 fields used by Azure Stack Hub.
	GalleryEndpoint string `json:"galleryEndpoint"`
	GraphEndpoint   string `json:"graphEndpoint"`
	PortalEndpoint  string `json:"portalEndpoint"`
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

// toEnvironment converts a metadata endpoints Response into an azure.Environment.
func (r Response) toEnvironment(armEndpoint, cloudName string) *azure.Environment {
	env := azure.Environment{
		Name:                    r.Name,
		ManagementPortalURL:     cmp.Or(r.Portal, r.PortalEndpoint),
		ActiveDirectoryEndpoint: r.Authentication.LoginEndpoint,

		ResourceManagerEndpoint: cmp.Or(r.ResourceManager, armEndpoint),
		GalleryEndpoint:         r.GalleryEndpoint,
		GraphEndpoint:           cmp.Or(r.Graph, r.GraphEndpoint),
		MicrosoftGraphEndpoint:  r.MicrosoftGraphResourceID,
		BatchManagementEndpoint: r.Batch,

		StorageEndpointSuffix:       cmp.Or(r.Suffixes.Storage, storageSuffixFromARM(armEndpoint)),
		SQLDatabaseDNSSuffix:        r.Suffixes.SQLServerHostname,
		KeyVaultDNSSuffix:           r.Suffixes.KeyVaultDNS,
		ManagedHSMDNSSuffix:         r.Suffixes.MHSMDNS,
		MariaDBDNSSuffix:            r.Suffixes.MariaDBServerEndpoint,
		MySQLDatabaseDNSSuffix:      r.Suffixes.MySQLServerEndpoint,
		PostgresqlDatabaseDNSSuffix: r.Suffixes.PostgreSQLServerEndpoint,
		ContainerRegistryDNSSuffix:  r.Suffixes.ACRLoginServer,
		SynapseEndpointSuffix:       r.Suffixes.SynapseAnalytics,
		DatalakeSuffix:              r.Suffixes.AzureDataLakeStoreFileSystem,

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

	if env.Name == "" && cloudName == azuretypes.StackCloud.Name() {
		env.Name = hybridEnvironmentName
	}

	if env.KeyVaultDNSSuffix == "" && env.StorageEndpointSuffix != "" {
		env.KeyVaultDNSSuffix = "vault." + env.StorageEndpointSuffix
	}

	if len(r.Authentication.Audiences) > 0 {
		env.ServiceManagementEndpoint = r.Authentication.Audiences[0]
		env.TokenAudience = r.Authentication.Audiences[0]
	}

	if env.KeyVaultDNSSuffix != "" {
		env.KeyVaultEndpoint = fmt.Sprintf("https://%s/", env.KeyVaultDNSSuffix)
		env.ResourceIdentifiers.KeyVault = fmt.Sprintf("https://%s", env.KeyVaultDNSSuffix)
	}
	if env.ManagedHSMDNSSuffix != "" {
		env.ManagedHSMEndpoint = fmt.Sprintf("https://%s/", env.ManagedHSMDNSSuffix)
		env.ResourceIdentifiers.ManagedHSM = fmt.Sprintf("https://%s", env.ManagedHSMDNSSuffix)
	}

	return &env
}

// storageSuffixFromARM returns the armEndpoint host without its first label.
func storageSuffixFromARM(armEndpoint string) string {
	u, err := url.Parse(armEndpoint)
	if err != nil {
		return ""
	}
	host := u.Hostname()
	if net.ParseIP(host) != nil {
		return ""
	}
	_, suffix, found := strings.Cut(host, ".")
	if !found {
		return ""
	}
	return suffix
}

// EnvironmentFromURL builds an azure.Environment from the metadata endpoints at armEndpoint.
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
			return r.toEnvironment(armEndpoint, cloudName), nil
		}
	}

	return responses[0].toEnvironment(armEndpoint, cloudName), fmt.Errorf("cloudName not found")
}
