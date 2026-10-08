/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scope

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
)

// AzureClients contains all the Azure clients used by the scopes.
type AzureClients struct {
	TokenCredential            azcore.TokenCredential
	cloudConfig                cloud.Configuration
	ResourceManagerEndpoint    string
	ResourceManagerVMDNSSuffix string
	activeDirectoryEndpoint    string
	tokenAudience              string

	authType infrav1.IdentityType

	cloudEnvironment string
	tenantID         string
	clientID         string
	clientSecret     string
	subscriptionID   string
}

// CloudConfiguration returns the Azure cloud.Configuration for SDK v2 clients.
func (c *AzureClients) CloudConfiguration() cloud.Configuration {
	return c.cloudConfig
}

// CloudEnvironment returns the Azure environment the controller runs in.
func (c *AzureClients) CloudEnvironment() string {
	return c.cloudEnvironment
}

// TenantID returns the Azure tenant id the controller runs in.
func (c *AzureClients) TenantID() string {
	return c.tenantID
}

// ClientID returns the Azure client id from the controller environment.
func (c *AzureClients) ClientID() string {
	return c.clientID
}

// ClientSecret returns the Azure client secret from the controller environment.
func (c *AzureClients) ClientSecret() string {
	return c.clientSecret
}

// SubscriptionID returns the Azure subscription id of the cluster,
// either specified or from the environment.
func (c *AzureClients) SubscriptionID() string {
	return c.subscriptionID
}

// Token returns the Azure token credential of the cluster used for SDKv2 services.
func (c *AzureClients) Token() azcore.TokenCredential {
	return c.TokenCredential
}

// HashKey returns a base64 url encoded sha256 hash for the Auth scope (Azure TenantID + CloudEnv + SubscriptionID +
// ClientID).
func (c *AzureClients) HashKey() string {
	hasher := sha256.New()
	_, _ = hasher.Write([]byte(c.TenantID() + c.CloudEnvironment() + c.SubscriptionID() + c.ClientID() + string(c.authType)))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (c *AzureClients) setCredentialsWithProvider(ctx context.Context, subscriptionID, environmentName string, credentialsProvider CredentialsProvider) error {
	if credentialsProvider == nil {
		return fmt.Errorf("credentials provider cannot have an empty value")
	}

	err := c.getSettingsFromEnvironment(environmentName)
	if err != nil {
		return err
	}

	if subscriptionID == "" {
		subscriptionID = c.SubscriptionID()
		if subscriptionID == "" {
			return fmt.Errorf("error creating azure services. subscriptionID is not set in cluster or AZURE_SUBSCRIPTION_ID env var")
		}
	}

	c.subscriptionID = strings.TrimSuffix(subscriptionID, "\n")
	c.tenantID = strings.TrimSuffix(credentialsProvider.GetTenantID(), "\n")
	c.clientID = strings.TrimSuffix(credentialsProvider.GetClientID(), "\n")

	clientSecret, err := credentialsProvider.GetClientSecret(ctx)
	if err != nil {
		return err
	}
	c.clientSecret = strings.TrimSuffix(clientSecret, "\n")

	c.authType = credentialsProvider.Type()

	tokenCredential, err := credentialsProvider.GetTokenCredential(ctx, c.ResourceManagerEndpoint, c.activeDirectoryEndpoint, c.tokenAudience)
	if err != nil {
		return err
	}
	c.TokenCredential = tokenCredential
	return err
}

func (c *AzureClients) getSettingsFromEnvironment(environmentName string) error {
	if environmentName == "" {
		environmentName = azure.PublicCloudName
	}
	setValue(&c.subscriptionID, "AZURE_SUBSCRIPTION_ID")

	// These strings were well-known by go-autorest which we don't use anymore.
	// This translates those strings into the corresponding SDKv2 configuration.
	var cloudConfig cloud.Configuration
	switch environmentName {
	case azure.ChinaCloudName:
		cloudConfig = cloud.AzureChina
		c.ResourceManagerVMDNSSuffix = "cloudapp.chinacloudapi.cn"
	case azure.GermanCloudName:
		// Not built in to SDKv2.
		// https://github.com/Azure/go-autorest/blob/33e12ab7683c1c236a863ccfbfdd78c626f7fe28/autorest/azure/environments.go#L243
		cloudConfig = cloud.Configuration{
			ActiveDirectoryAuthorityHost: "https://login.microsoftonline.de/",
			Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
				cloud.ResourceManager: {
					Audience: "https://management.microsoftazure.de/",
					Endpoint: "https://management.microsoftazure.de/",
				},
			},
		}
		c.ResourceManagerVMDNSSuffix = "cloudapp.microsoftazure.de"
	case azure.PublicCloudName:
		cloudConfig = cloud.AzurePublic
		c.ResourceManagerVMDNSSuffix = "cloudapp.azure.com"
	case azure.USGovernmentCloudName:
		cloudConfig = cloud.AzureGovernment
		c.ResourceManagerVMDNSSuffix = "cloudapp.usgovcloudapi.net"
	case azure.USSecCloudName:
		env, err := loadCloudEnvironmentFromFile()
		if err != nil {
			return err
		}
		cloudConfig = cloud.Configuration{
			ActiveDirectoryAuthorityHost: env.ActiveDirectoryEndpoint,
			Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
				cloud.ResourceManager: {
					Audience: env.TokenAudience,
					Endpoint: env.ResourceManagerEndpoint,
				},
			},
		}
		if env.ResourceManagerVMDNSSuffix != "" {
			c.ResourceManagerVMDNSSuffix = env.ResourceManagerVMDNSSuffix
		}
	default:
		return fmt.Errorf("invalid cloud environment name %q", c.CloudEnvironment())
	}

	c.cloudEnvironment = environmentName
	c.cloudConfig = cloudConfig
	c.activeDirectoryEndpoint = cloudConfig.ActiveDirectoryAuthorityHost
	c.ResourceManagerEndpoint = cloudConfig.Services[cloud.ResourceManager].Endpoint
	c.tokenAudience = cloudConfig.Services[cloud.ResourceManager].Audience

	return nil
}

// setValue adds the specified environment variable value to the Values map if it exists.
func setValue(value *string, key string) {
	if v := os.Getenv(key); v != "" {
		*value = v
	}
}

// cloudEnvironmentFile represents the JSON structure of the Azure environment
// file used for air-gapped clouds (USSec, Stack). This format matches the
// Environment struct in cloud-provider-azure/pkg/azclient/cloud.go.
type cloudEnvironmentFile struct {
	Name                       string `json:"name"`
	ResourceManagerEndpoint    string `json:"resourceManagerEndpoint"`
	ActiveDirectoryEndpoint    string `json:"activeDirectoryEndpoint"`
	TokenAudience              string `json:"tokenAudience"`
	ResourceManagerVMDNSSuffix string `json:"resourceManagerVMDNSSuffix,omitempty"`
	StorageEndpointSuffix      string `json:"storageEndpointSuffix,omitempty"`
	KeyVaultDNSSuffix          string `json:"keyVaultDNSSuffix,omitempty"`
	ContainerRegistryDNSSuffix string `json:"containerRegistryDNSSuffix,omitempty"`
}

// loadCloudEnvironmentFromFile reads the Azure environment JSON file specified
// by the AZURE_ENVIRONMENT_FILEPATH environment variable. This is the same
// mechanism used by cloud-provider-azure for Azure Stack and USSec clouds.
// The file is produced by an external process (the installer or an
// administrator) that discovers endpoints from the ARM metadata service
// within the air-gapped network.
func loadCloudEnvironmentFromFile() (*cloudEnvironmentFile, error) {
	envFilePath, ok := os.LookupEnv("AZURE_ENVIRONMENT_FILEPATH")
	if !ok || envFilePath == "" {
		return nil, fmt.Errorf(
			"AZURE_ENVIRONMENT_FILEPATH environment variable must be set when using cloud environment %q", azure.USSecCloudName)
	}
	content, err := os.ReadFile(envFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cloud environment file %q: %w",
			envFilePath, err)
	}
	var env cloudEnvironmentFile
	if err := json.Unmarshal(content, &env); err != nil {
		return nil, fmt.Errorf("failed to parse cloud environment file %q: %w",
			envFilePath, err)
	}
	if env.ResourceManagerEndpoint == "" {
		return nil, fmt.Errorf("resourceManagerEndpoint is required in cloud environment file %q", envFilePath)
	}
	if env.ActiveDirectoryEndpoint == "" {
		return nil, fmt.Errorf("activeDirectoryEndpoint is required in cloud environment file %q", envFilePath)
	}
	if env.TokenAudience == "" {
		return nil, fmt.Errorf("tokenAudience is required in cloud environment file %q", envFilePath)
	}
	return &env, nil
}
