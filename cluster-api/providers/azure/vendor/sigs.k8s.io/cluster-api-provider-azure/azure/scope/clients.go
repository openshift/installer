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
	default:
		return fmt.Errorf("invalid cloud environment name %q", c.CloudEnvironment())
	}

	c.cloudEnvironment = environmentName
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
