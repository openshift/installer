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
	"github.com/Azure/go-autorest/autorest"
	azureautorest "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// AzureClients contains all the Azure clients used by the scopes.
type AzureClients struct {
	auth.EnvironmentSettings

	Authorizer                 autorest.Authorizer
	TokenCredential            azcore.TokenCredential
	ResourceManagerEndpoint    string
	ResourceManagerVMDNSSuffix string
}

// CloudEnvironment returns the Azure environment the controller runs in.
func (c *AzureClients) CloudEnvironment() string {
	return c.Environment.Name
}

// TenantID returns the Azure tenant id the controller runs in.
func (c *AzureClients) TenantID() string {
	return c.Values[auth.TenantID]
}

// ClientID returns the Azure client id from the controller environment.
func (c *AzureClients) ClientID() string {
	return c.Values[auth.ClientID]
}

// ClientSecret returns the Azure client secret from the controller environment.
func (c *AzureClients) ClientSecret() string {
	return c.Values[auth.ClientSecret]
}

// SubscriptionID returns the Azure subscription id of the cluster,
// either specified or from the environment.
func (c *AzureClients) SubscriptionID() string {
	return c.Values[auth.SubscriptionID]
}

// Token returns the Azure token credential of the cluster used for SDKv2 services.
func (c *AzureClients) Token() azcore.TokenCredential {
	return c.TokenCredential
}

// HashKey returns a base64 url encoded sha256 hash for the Auth scope (Azure TenantID + CloudEnv + SubscriptionID +
// ClientID).
func (c *AzureClients) HashKey() string {
	hasher := sha256.New()
	_, _ = hasher.Write([]byte(c.TenantID() + c.CloudEnvironment() + c.SubscriptionID() + c.ClientID()))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (c *AzureClients) setCredentialsWithProvider(ctx context.Context, subscriptionID, environmentName string, credentialsProvider CredentialsProvider) error {
	if credentialsProvider == nil {
		return fmt.Errorf("credentials provider cannot have an empty value")
	}

	settings, err := c.getSettingsFromEnvironment(environmentName)
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
	c.Values[auth.SubscriptionID] = strings.TrimSuffix(subscriptionID, "\n")
	c.Values[auth.TenantID] = strings.TrimSuffix(credentialsProvider.GetTenantID(), "\n")
	c.Values[auth.ClientID] = strings.TrimSuffix(credentialsProvider.GetClientID(), "\n")

	clientSecret, err := credentialsProvider.GetClientSecret(ctx)
	if err != nil {
		return err
	}
	c.Values[auth.ClientSecret] = strings.TrimSuffix(clientSecret, "\n")

	tokenCredential, err := credentialsProvider.GetTokenCredential(ctx, c.ResourceManagerEndpoint, c.Environment.ActiveDirectoryEndpoint, c.Environment.TokenAudience)
	if err != nil {
		return err
	}
	c.TokenCredential = tokenCredential
	c.Authorizer, err = credentialsProvider.GetAuthorizer(ctx, tokenCredential, c.Environment.TokenAudience)
	return err
}

func (c *AzureClients) getSettingsFromEnvironment(environmentName string) (s auth.EnvironmentSettings, err error) {
	s = auth.EnvironmentSettings{
		Values: map[string]string{},
	}
	s.Values[auth.EnvironmentName] = environmentName
	setValue(s, auth.SubscriptionID)
	setValue(s, auth.TenantID)
	setValue(s, auth.AuxiliaryTenantIDs)
	setValue(s, auth.ClientID)
	setValue(s, auth.ClientSecret)
	setValue(s, auth.CertificatePath)
	setValue(s, auth.CertificatePassword)
	setValue(s, auth.Username)
	setValue(s, auth.Password)
	setValue(s, auth.Resource)
	if v := s.Values[auth.EnvironmentName]; v == "" {
		s.Environment = azureautorest.PublicCloud
	} else {
		s.Environment, err = azureautorest.EnvironmentFromName(v)
	}
	if s.Values[auth.Resource] == "" {
		s.Values[auth.Resource] = s.Environment.ResourceManagerEndpoint
	}
	return
}

// setValue adds the specified environment variable value to the Values map if it exists.
func setValue(settings auth.EnvironmentSettings, key string) {
	if v := os.Getenv(key); v != "" {
		settings.Values[key] = v
	}
}
