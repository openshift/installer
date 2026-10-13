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
	azureautorest "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// AzureClients contains all the Azure clients used by the scopes.
type AzureClients struct {
	auth.EnvironmentSettings

	TokenCredential            azcore.TokenCredential
	ResourceManagerEndpoint    string
	ResourceManagerVMDNSSuffix string

	authType infrav1.IdentityType
}

// CloudEnvironment returns the Azure environment the controller runs in.
func (c *AzureClients) CloudEnvironment() string {
	return c.Environment.Name
}

// TenantID returns the Azure tenant id the controller runs in.
func (c *AzureClients) TenantID() string {
	return c.Values["AZURE_TENANT_ID"]
}

// ClientID returns the Azure client id from the controller environment.
func (c *AzureClients) ClientID() string {
	return c.Values["AZURE_CLIENT_ID"]
}

// ClientSecret returns the Azure client secret from the controller environment.
func (c *AzureClients) ClientSecret() string {
	return c.Values["AZURE_CLIENT_SECRET"]
}

// SubscriptionID returns the Azure subscription id of the cluster,
// either specified or from the environment.
func (c *AzureClients) SubscriptionID() string {
	return c.Values["AZURE_SUBSCRIPTION_ID"]
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

func (c *AzureClients) setCredentialsWithProvider(ctx context.Context, subscriptionID, environmentName, armEndpoint string, credentialsProvider CredentialsProvider) error {
	if credentialsProvider == nil {
		return fmt.Errorf("credentials provider cannot have an empty value")
	}

	settings, err := c.getSettingsFromEnvironment(environmentName, armEndpoint)
	if err != nil {
		return err
	}

	if subscriptionID == "" {
		subscriptionID = settings.GetSubscriptionID()
		if subscriptionID == "" {
			return fmt.Errorf("error creating azure services. subscriptionID is not set in cluster or AZURE_SUBSCRIPTION_ID env var")
		}
	}

	c.EnvironmentSettings = settings
	c.ResourceManagerEndpoint = settings.Environment.ResourceManagerEndpoint
	c.ResourceManagerVMDNSSuffix = settings.Environment.ResourceManagerVMDNSSuffix
	c.Values["AZURE_SUBSCRIPTION_ID"] = strings.TrimSuffix(subscriptionID, "\n")
	c.Values["AZURE_TENANT_ID"] = strings.TrimSuffix(credentialsProvider.GetTenantID(), "\n")
	c.Values["AZURE_CLIENT_ID"] = strings.TrimSuffix(credentialsProvider.GetClientID(), "\n")

	clientSecret, err := credentialsProvider.GetClientSecret(ctx)
	if err != nil {
		return err
	}
	c.Values["AZURE_CLIENT_SECRET"] = strings.TrimSuffix(clientSecret, "\n")

	c.authType = credentialsProvider.Type()

	tokenCredential, err := credentialsProvider.GetTokenCredential(ctx, c.ResourceManagerEndpoint, c.Environment.ActiveDirectoryEndpoint, c.Environment.TokenAudience)
	if err != nil {
		return err
	}
	c.TokenCredential = tokenCredential
	return err
}

func (c *AzureClients) getSettingsFromEnvironment(environmentName, armEndpoint string) (s auth.EnvironmentSettings, err error) {
	s = auth.EnvironmentSettings{
		Values: map[string]string{},
	}
	s.Values["AZURE_ENVIRONMENT"] = environmentName
	setValue(s, "AZURE_SUBSCRIPTION_ID")
	setValue(s, "AZURE_TENANT_ID")
	setValue(s, "AZURE_AUXILIARY_TENANT_IDS")
	setValue(s, "AZURE_CLIENT_ID")
	setValue(s, "AZURE_CLIENT_SECRET")
	setValue(s, "AZURE_CERTIFICATE_PATH")
	setValue(s, "AZURE_CERTIFICATE_PASSWORD")
	setValue(s, "AZURE_USERNAME")
	setValue(s, "AZURE_PASSWORD")
	setValue(s, "AZURE_AD_RESOURCE")
	if v := s.Values["AZURE_ENVIRONMENT"]; v == "" {
		s.Environment = azureautorest.PublicCloud
	} else if len(armEndpoint) > 0 {
		s.Environment, err = azureautorest.EnvironmentFromURL(armEndpoint)
	} else {
		s.Environment, err = azureautorest.EnvironmentFromName(v)
	}
	if s.Values["AZURE_AD_RESOURCE"] == "" {
		s.Values["AZURE_AD_RESOURCE"] = s.Environment.ResourceManagerEndpoint
	}
	return
}

// setValue adds the specified environment variable value to the Values map if it exists.
func setValue(settings auth.EnvironmentSettings, key string) {
	if v := os.Getenv(key); v != "" {
		settings.Values[key] = v
	}
}
