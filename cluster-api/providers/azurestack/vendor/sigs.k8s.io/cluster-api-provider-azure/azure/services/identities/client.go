/*
Copyright 2022 The Kubernetes Authors.

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

package identities

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// Client wraps go-sdk.
type Client interface {
	Get(ctx context.Context, resourceGroupName, name string) (armmsi.Identity, error)
	GetClientID(ctx context.Context, providerID string) (string, error)
}

// AzureClient contains the Azure go-sdk Client.
type AzureClient struct {
	userAssignedIdentities *armmsi.UserAssignedIdentitiesClient
}

// NewClient creates a new MSI client from an authorizer.
func NewClient(auth azure.Authorizer) (Client, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create identities client options")
	}
	factory, err := armmsi.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armmsi client factory")
	}
	return &AzureClient{factory.NewUserAssignedIdentitiesClient()}, nil
}

// NewClientBySub creates a new MSI client with a given subscriptionID.
func NewClientBySub(auth azure.Authorizer, subscriptionID string) (Client, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create identities client options")
	}
	factory, err := armmsi.NewClientFactory(subscriptionID, auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armmsi client factory")
	}
	return &AzureClient{factory.NewUserAssignedIdentitiesClient()}, nil
}

// Get returns a managed service identity.
func (ac *AzureClient) Get(ctx context.Context, resourceGroupName, name string) (armmsi.Identity, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "identities.AzureClient.Get")
	defer done()

	resp, err := ac.userAssignedIdentities.Get(ctx, resourceGroupName, name, nil)
	if err != nil {
		return armmsi.Identity{}, err
	}
	return resp.Identity, nil
}

// GetClientID returns the client ID of a managed service identity, given its full URL identifier.
func (ac *AzureClient) GetClientID(ctx context.Context, providerID string) (string, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "identities.AzureClient.GetClientID")
	defer done()

	parsed, err := azureutil.ParseResourceID(providerID)
	if err != nil {
		return "", err
	}
	ident, err := ac.Get(ctx, parsed.ResourceGroupName, parsed.Name)
	if err != nil {
		return "", err
	}
	return ptr.Deref(ident.Properties.ClientID, ""), nil
}
