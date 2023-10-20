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

package token

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// AzureClient to get azure active directory token.
type AzureClient struct {
	aadToken *azidentity.ClientSecretCredential
}

// NewClient creates a new azure active directory token client from an authorizer.
func NewClient(auth azure.Authorizer) (*AzureClient, error) {
	aadToken, err := newAzureActiveDirectoryTokenClient(auth.TenantID(),
		auth.ClientID(),
		auth.ClientSecret(),
		auth.CloudEnvironment())
	if err != nil {
		return nil, err
	}
	return &AzureClient{
		aadToken: aadToken,
	}, nil
}

// newAzureActiveDirectoryTokenClient creates a new aad token client from an authorizer.
func newAzureActiveDirectoryTokenClient(tenantID, clientID, clientSecret, envName string) (*azidentity.ClientSecretCredential, error) {
	cliOpts, err := azure.ARMClientOptions(envName)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting client options")
	}
	clientOptions := &azidentity.ClientSecretCredentialOptions{
		ClientOptions: cliOpts.ClientOptions,
	}
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting az client secret credentials")
	}
	return cred, nil
}

// GetAzureActiveDirectoryToken gets the token for authentication with azure active directory.
func (ac *AzureClient) GetAzureActiveDirectoryToken(ctx context.Context, resourceID string) (string, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "aadToken.GetToken")
	defer done()

	spnAccessToken, err := ac.aadToken.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{resourceID + "/.default"}})
	if err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}
	return spnAccessToken.Token, nil
}
