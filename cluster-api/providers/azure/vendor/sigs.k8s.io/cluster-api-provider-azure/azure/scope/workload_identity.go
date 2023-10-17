/*
Copyright 2023 The Kubernetes Authors.

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
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/pkg/errors"
)

/*

For workload identity to work we need the following.

|-----------------------------------------------------------------------------------|
|AZURE_AUTHORITY_HOST       | The Azure Active Directory (AAD) endpoint.            |
|AZURE_CLIENT_ID            | The client ID of the Azure AD                         |
|                           | application or user-assigned managed identity.        |
|AZURE_TENANT_ID            | The tenant ID of the Azure subscription.              |
|AZURE_FEDERATED_TOKEN_FILE | The path of the projected service account token file. |
|-----------------------------------------------------------------------------------|

With the current implementation, AZURE_CLIENT_ID and AZURE_TENANT_ID are read via AzureClusterIdentity
object and fallback to reading from env variables if not found on AzureClusterIdentity.

AZURE_FEDERATED_TOKEN_FILE is the path of the projected service account token which is by default
"/var/run/secrets/azure/tokens/azure-identity-token".
The path can be overridden by setting "AZURE_FEDERATED_TOKEN_FILE" env variable.

*/

const (
	// azureFederatedTokenFileEnvKey is the env key for AZURE_FEDERATED_TOKEN_FILE.
	azureFederatedTokenFileEnvKey = "AZURE_FEDERATED_TOKEN_FILE"
	// azureClientIDEnvKey is the env key for AZURE_CLIENT_ID.
	azureClientIDEnvKey = "AZURE_CLIENT_ID"
	// azureTenantIDEnvKey is the env key for AZURE_TENANT_ID.
	azureTenantIDEnvKey = "AZURE_TENANT_ID"
	// azureTokenFilePath is the path of the projected token.
	azureTokenFilePath = "/var/run/secrets/azure/tokens/azure-identity-token" // #nosec G101
	// azureFederatedTokenFileRefreshTime is the time interval after which it should be read again.
	azureFederatedTokenFileRefreshTime = 5 * time.Minute
)

type workloadIdentityCredential struct {
	assertion string
	file      string
	cred      *azidentity.ClientAssertionCredential
	lastRead  time.Time
}

// WorkloadIdentityCredentialOptions contains the configurable options for azwi.
type WorkloadIdentityCredentialOptions struct {
	azcore.ClientOptions
	ClientID      string
	TenantID      string
	TokenFilePath string
}

// NewWorkloadIdentityCredentialOptions returns an empty instance of WorkloadIdentityCredentialOptions.
func NewWorkloadIdentityCredentialOptions() *WorkloadIdentityCredentialOptions {
	return &WorkloadIdentityCredentialOptions{}
}

// WithClientID sets client ID to WorkloadIdentityCredentialOptions.
func (w *WorkloadIdentityCredentialOptions) WithClientID(clientID string) *WorkloadIdentityCredentialOptions {
	w.ClientID = strings.TrimSpace(clientID)
	return w
}

// WithTenantID sets tenant ID to WorkloadIdentityCredentialOptions.
func (w *WorkloadIdentityCredentialOptions) WithTenantID(tenantID string) *WorkloadIdentityCredentialOptions {
	w.TenantID = strings.TrimSpace(tenantID)
	return w
}

// getProjectedTokenPath return projected token file path from the env variable.
func getProjectedTokenPath() string {
	tokenPath := strings.TrimSpace(os.Getenv(azureFederatedTokenFileEnvKey))
	if tokenPath == "" {
		return azureTokenFilePath
	}
	return tokenPath
}

// WithDefaults sets token file path. It also sets the client tenant ID from injected env in
// case empty values are passed.
func (w *WorkloadIdentityCredentialOptions) WithDefaults() (*WorkloadIdentityCredentialOptions, error) {
	w.TokenFilePath = getProjectedTokenPath()

	// Fallback to using client ID from env variable if not set.
	if w.ClientID == "" {
		w.ClientID = strings.TrimSpace(os.Getenv(azureClientIDEnvKey))
		if w.ClientID == "" {
			return nil, errors.New("empty client ID")
		}
	}

	// Fallback to using tenant ID from env variable.
	if w.TenantID == "" {
		w.TenantID = strings.TrimSpace(os.Getenv(azureTenantIDEnvKey))
		if w.TenantID == "" {
			return nil, errors.New("empty tenant ID")
		}
	}
	return w, nil
}

// NewWorkloadIdentityCredential returns a workload identity credential.
func NewWorkloadIdentityCredential(options *WorkloadIdentityCredentialOptions) (azcore.TokenCredential, error) {
	w := &workloadIdentityCredential{file: options.TokenFilePath}
	cred, err := azidentity.NewClientAssertionCredential(options.TenantID, options.ClientID, w.getAssertion, &azidentity.ClientAssertionCredentialOptions{ClientOptions: options.ClientOptions})
	if err != nil {
		return nil, err
	}
	w.cred = cred
	return w, nil
}

// GetToken returns the token for azwi.
func (w *workloadIdentityCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return w.cred.GetToken(ctx, opts)
}

func (w *workloadIdentityCredential) getAssertion(context.Context) (string, error) {
	if now := time.Now(); w.lastRead.Add(azureFederatedTokenFileRefreshTime).Before(now) {
		content, err := os.ReadFile(w.file)
		if err != nil {
			return "", err
		}
		w.assertion = string(content)
		w.lastRead = now
	}
	return w.assertion, nil
}
