package azure

import (
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/azure"
)

func TestResolveCloudEndpointsWellKnown(t *testing.T) {
	cases := []struct {
		cloudName   azure.CloudEnvironment
		cloudConfig cloud.Configuration
	}{
		{
			cloudName:   azure.PublicCloud,
			cloudConfig: cloud.AzurePublic,
		},
		{
			cloudName:   azure.USGovernmentCloud,
			cloudConfig: cloud.AzureGovernment,
		},
		{
			cloudName:   azure.ChinaCloud,
			cloudConfig: cloud.AzureChina,
		},
	}

	for _, tc := range cases {
		t.Run(string(tc.cloudName), func(t *testing.T) {
			endpoints, err := ResolveCloudEndpoints(tc.cloudName, "")
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tc.cloudConfig, endpoints.CloudConfig)
		})
	}
}

func TestResolveCloudEndpointsDiscovery(t *testing.T) {
	cloudNames := []azure.CloudEnvironment{
		azure.StackCloud,
		azure.USSecCloud,
		azure.USNatCloud,
	}

	for _, cloudName := range cloudNames {
		t.Run(string(cloudName), func(t *testing.T) {
			// Discovery clouds require an ARM endpoint; well-known clouds ignore it.
			_, err := ResolveCloudEndpoints(cloudName, "")
			assert.EqualError(t, err, "ARM endpoint URL must not be empty for endpoint discovery")

			body := `{
	"name": "` + string(cloudName) + `",
	"authentication": {
		"loginEndpoint": "https://login.example.test/",
		"audiences": ["https://management.core.example.test/"]
	}
}`
			server := newEndpointsServer(t, http.StatusOK, body)

			endpoints, err := ResolveCloudEndpoints(cloudName, server.URL)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, string(cloudName), endpoints.Environment.Name)
			assert.Equal(t, cloud.Configuration{
				ActiveDirectoryAuthorityHost: "https://login.example.test/",
				Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
					cloud.ResourceManager: {
						Audience: "https://management.core.example.test/",
						Endpoint: server.URL,
					},
				},
			}, endpoints.CloudConfig)
		})
	}
}

func TestResolveCloudEndpointsDiscoveryNameOverride(t *testing.T) {
	// When the metadata document does not contain the requested cloud, the
	// first environment is used and its name is overridden.
	server := newEndpointsServer(t, http.StatusOK, `[{"name": "AzureCloud"}, {"name": "AzureUSGovernment"}]`)

	endpoints, err := ResolveCloudEndpoints(azure.USNatCloud, server.URL)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, string(azure.USNatCloud), endpoints.Environment.Name)
}
