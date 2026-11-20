package azure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"

	azuretypes "github.com/openshift/installer/pkg/types/azure"
)

const (
	usSecEndpoints = `{
	"portal": "https://portal.azure.microsoft.scloud",
	"authentication": {
		"loginEndpoint": "https://login.microsoftonline.microsoft.scloud",
		"audiences": ["https://management.core.microsoft.scloud/", "https://management.azure.microsoft.scloud/"],
		"tenant": "common",
		"identityProvider": "AAD"
	},
	"name": "AzureUSSecCloud",
	"resourceManager": "https://management.azure.microsoft.scloud",
	"suffixes": {
		"storage": "core.microsoft.scloud",
		"keyVaultDns": "vault.cloudapi.microsoft.scloud"
	}
}`

	multipleEndpoints = `[
	{"name": "AzureCloud", "resourceManager": "https://management.azure.com"},
	{"name": "AzureChinaCloud", "resourceManager": "https://management.chinacloudapi.cn"},
	{"name": "AzureUSGovernment", "resourceManager": "https://management.usgovcloudapi.net"}
]`
)

// newEndpointsServer starts a TLS server that serves body with the given
// status code on /metadata/endpoints, and points httpClient at it for the
// duration of the test.
func newEndpointsServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metadata/endpoints" || r.URL.Query().Get("api-version") == "" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		w.Write([]byte(body)) //nolint:all
	}))
	t.Cleanup(server.Close)

	orig := httpClient
	httpClient = server.Client()
	t.Cleanup(func() { httpClient = orig })

	return server
}

func TestEnvironmentFromURL(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		body      string
		cloudName string
		expectEnv bool
		envName   string
		err       string
	}{
		{
			name:      "single object response",
			status:    http.StatusOK,
			body:      usSecEndpoints,
			cloudName: "AzureUSSecCloud",
			expectEnv: true,
			envName:   "AzureUSSecCloud",
		},
		{
			name:      "array response with matching cloud name",
			status:    http.StatusOK,
			body:      multipleEndpoints,
			cloudName: "AzureChinaCloud",
			expectEnv: true,
			envName:   "AzureChinaCloud",
		},
		{
			name:      "array response with empty cloud name returns first",
			status:    http.StatusOK,
			body:      multipleEndpoints,
			expectEnv: true,
			envName:   "AzureCloud",
		},
		{
			name:      "array response with unmatched cloud name returns first and an error",
			status:    http.StatusOK,
			body:      multipleEndpoints,
			cloudName: "AzureUSNatCloud",
			expectEnv: true,
			envName:   "AzureCloud",
			err:       "cloudName not found",
		},
		{
			name:   "empty body",
			status: http.StatusOK,
			body:   "  ",
			err:    "no endpoints returned",
		},
		{
			name:   "empty array",
			status: http.StatusOK,
			body:   "[]",
			err:    "no endpoints returned",
		},
		{
			name:   "unexpected body",
			status: http.StatusOK,
			body:   `"AzureUSSecCloud"`,
			err:    "error unmarshalling endpoint response: unexpected type",
		},
		{
			name:   "non-200 response",
			status: http.StatusInternalServerError,
			body:   "",
			err:    "error getting endpoints: error getting response: 500 Internal Server Error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := newEndpointsServer(t, tc.status, tc.body)

			env, err := EnvironmentFromURL(server.URL, "", tc.cloudName)
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.err)
			}
			if !tc.expectEnv {
				assert.Nil(t, env)
				return
			}
			if assert.NotNil(t, env) {
				assert.Equal(t, tc.envName, env.Name)
			}
		})
	}
}

func TestEnvironmentFromURLFields(t *testing.T) {
	server := newEndpointsServer(t, http.StatusOK, usSecEndpoints)

	env, err := EnvironmentFromURL(server.URL, "", "AzureUSSecCloud")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "https://login.microsoftonline.microsoft.scloud", env.ActiveDirectoryEndpoint)
	assert.Equal(t, "https://management.azure.microsoft.scloud", env.ResourceManagerEndpoint)
	assert.Equal(t, "https://management.core.microsoft.scloud/", env.TokenAudience)
	assert.Equal(t, "https://management.core.microsoft.scloud/", env.ServiceManagementEndpoint)
	assert.Equal(t, "core.microsoft.scloud", env.StorageEndpointSuffix)
	assert.Equal(t, "https://vault.cloudapi.microsoft.scloud/", env.KeyVaultEndpoint)
	assert.Empty(t, env.ManagedHSMEndpoint)
}

func TestEnvironmentFromURLResourceManagerFallback(t *testing.T) {
	server := newEndpointsServer(t, http.StatusOK, `{"name": "AzureUSNatCloud"}`)

	env, err := EnvironmentFromURL(server.URL, "", "AzureUSNatCloud")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, server.URL, env.ResourceManagerEndpoint)
}

func TestEnvironmentFromURLInvalidEndpoint(t *testing.T) {
	cases := []struct {
		name     string
		endpoint string
		err      string
	}{
		{
			name:     "empty endpoint",
			endpoint: "",
			err:      "armEndpoint is empty",
		},
		{
			name:     "http endpoint",
			endpoint: "http://management.azure.microsoft.scloud",
			err:      `error getting endpoints: invalid ARM endpoint "http://management.azure.microsoft.scloud": scheme must be https`,
		},
		{
			name:     "missing host",
			endpoint: "https://",
			err:      `error getting endpoints: invalid ARM endpoint "https://": host cannot be empty`,
		},
		//nolint:all
		{
			name:     "user information",
			endpoint: "https://user:pass@management.azure.microsoft.scloud",
			err:      `error getting endpoints: invalid ARM endpoint "https://user:pass@management.azure.microsoft.scloud": user information is not allowed`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env, err := EnvironmentFromURL(tc.endpoint, "", "")
			assert.EqualError(t, err, tc.err)
			assert.Nil(t, env)
		})
	}
}

func TestToEnvironment(t *testing.T) {
	const stackARM = "https://management.local.azurestack.external"

	cases := []struct {
		name      string
		response  Response
		arm       string
		cloudName string
		check     func(t *testing.T, env *azure.Environment)
	}{
		{
			name: "stack cloud without name defaults to HybridEnvironment",
			response: Response{
				GalleryEndpoint: "https://portal.local.azurestack.external:30015/",
				GraphEndpoint:   "https://graph.windows.net/",
				PortalEndpoint:  "https://portal.local.azurestack.external/",
				Authentication: Authentication{
					LoginEndpoint: "https://login.microsoftonline.com/",
					Audiences:     []string{"https://management.adfs.azurestack.local/abc"},
				},
			},
			arm:       stackARM,
			cloudName: string(azuretypes.StackCloud),
			check: func(t *testing.T, env *azure.Environment) {
				t.Helper()
				assert.Equal(t, "HybridEnvironment", env.Name)
				assert.Equal(t, "https://portal.local.azurestack.external:30015/", env.GalleryEndpoint)
				assert.Equal(t, "https://graph.windows.net/", env.GraphEndpoint)
				assert.Equal(t, "https://portal.local.azurestack.external/", env.ManagementPortalURL)
				assert.Equal(t, "https://login.microsoftonline.com/", env.ActiveDirectoryEndpoint)
				assert.Equal(t, "https://management.adfs.azurestack.local/abc", env.TokenAudience)
				assert.Equal(t, stackARM, env.ResourceManagerEndpoint)
				assert.Equal(t, "local.azurestack.external", env.StorageEndpointSuffix)
				assert.Equal(t, "vault.local.azurestack.external", env.KeyVaultDNSSuffix)
				assert.Equal(t, "https://vault.local.azurestack.external/", env.KeyVaultEndpoint)
				assert.Equal(t, "https://vault.local.azurestack.external", env.ResourceIdentifiers.KeyVault)
			},
		},
		{
			name:      "stack cloud keeps name from metadata",
			response:  Response{Name: "MyStack"},
			arm:       stackARM,
			cloudName: string(azuretypes.StackCloud),
			check: func(t *testing.T, env *azure.Environment) {
				t.Helper()
				assert.Equal(t, "MyStack", env.Name)
			},
		},
		{
			name:      "non-stack cloud without name stays empty",
			response:  Response{},
			arm:       "https://management.azure.microsoft.scloud",
			cloudName: string(azuretypes.USSecCloud),
			check: func(t *testing.T, env *azure.Environment) {
				t.Helper()
				assert.Empty(t, env.Name)
				assert.Equal(t, "azure.microsoft.scloud", env.StorageEndpointSuffix)
				assert.Equal(t, "vault.azure.microsoft.scloud", env.KeyVaultDNSSuffix)
				assert.Equal(t, "https://vault.azure.microsoft.scloud/", env.KeyVaultEndpoint)
			},
		},
		{
			name: "metadata values take precedence over fallbacks",
			response: Response{
				Name:            "AzureUSSecCloud",
				Portal:          "https://portal.azure.microsoft.scloud",
				PortalEndpoint:  "https://legacy-portal.example",
				Graph:           "https://graph.microsoft.scloud/",
				GraphEndpoint:   "https://legacy-graph.example",
				ResourceManager: "https://management.azure.microsoft.scloud",
				Suffixes: Suffixes{
					Storage:     "core.microsoft.scloud",
					KeyVaultDNS: "vault.cloudapi.microsoft.scloud",
					MHSMDNS:     "managedhsm.microsoft.scloud",
				},
			},
			arm:       "https://other.example.com",
			cloudName: string(azuretypes.USSecCloud),
			check: func(t *testing.T, env *azure.Environment) {
				t.Helper()
				assert.Equal(t, "https://portal.azure.microsoft.scloud", env.ManagementPortalURL)
				assert.Equal(t, "https://graph.microsoft.scloud/", env.GraphEndpoint)
				assert.Equal(t, "https://management.azure.microsoft.scloud", env.ResourceManagerEndpoint)
				assert.Equal(t, "core.microsoft.scloud", env.StorageEndpointSuffix)
				assert.Equal(t, "vault.cloudapi.microsoft.scloud", env.KeyVaultDNSSuffix)
				assert.Equal(t, "https://vault.cloudapi.microsoft.scloud/", env.KeyVaultEndpoint)
				assert.Equal(t, "https://vault.cloudapi.microsoft.scloud", env.ResourceIdentifiers.KeyVault)
				assert.Equal(t, "https://managedhsm.microsoft.scloud/", env.ManagedHSMEndpoint)
				assert.Equal(t, "https://managedhsm.microsoft.scloud", env.ResourceIdentifiers.ManagedHSM)
			},
		},
		{
			name:     "no audiences leaves token audience empty",
			response: Response{},
			arm:      stackARM,
			check: func(t *testing.T, env *azure.Environment) {
				t.Helper()
				assert.Empty(t, env.TokenAudience)
				assert.Empty(t, env.ServiceManagementEndpoint)
				assert.Empty(t, env.ManagedHSMEndpoint)
				assert.Empty(t, env.ResourceIdentifiers.ManagedHSM)
			},
		},
		{
			name:     "underivable storage suffix leaves key vault empty",
			response: Response{},
			arm:      "https://127.0.0.1:8443",
			check: func(t *testing.T, env *azure.Environment) {
				t.Helper()
				assert.Empty(t, env.StorageEndpointSuffix)
				assert.Empty(t, env.KeyVaultDNSSuffix)
				assert.Empty(t, env.KeyVaultEndpoint)
				assert.Empty(t, env.ResourceIdentifiers.KeyVault)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, tc.response.toEnvironment(tc.arm, tc.cloudName))
		})
	}
}

func TestStorageSuffixFromARM(t *testing.T) {
	cases := []struct {
		arm    string
		expect string
	}{
		{arm: "https://management.local.azurestack.external", expect: "local.azurestack.external"},
		{arm: "https://management.local.azurestack.external/", expect: "local.azurestack.external"},
		{arm: "https://management.local.azurestack.external:443/some/path", expect: "local.azurestack.external"},
		{arm: "https://localhost", expect: ""},
		{arm: "https://127.0.0.1:8443", expect: ""},
		{arm: "https://[::1]:8443", expect: ""},
		{arm: "://bad", expect: ""},
		{arm: "", expect: ""},
	}

	for _, tc := range cases {
		t.Run(tc.arm, func(t *testing.T) {
			assert.Equal(t, tc.expect, storageSuffixFromARM(tc.arm))
		})
	}
}
