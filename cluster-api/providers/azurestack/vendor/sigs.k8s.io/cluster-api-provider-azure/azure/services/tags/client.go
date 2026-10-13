/*
Copyright 2020 The Kubernetes Authors.

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

package tags

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// client wraps go-sdk.
type client interface {
	GetAtScope(context.Context, string) (armresources.TagsResource, error)
	UpdateAtScope(context.Context, string, armresources.TagsPatchResource) (armresources.TagsResource, error)
}

// AzureClient contains the Azure go-sdk client.
type AzureClient struct {
	tags *armresources.TagsClient
}

var _ client = (*AzureClient)(nil)

// NewClient creates a tags client from an authorizer.
func NewClient(auth azure.Authorizer) (*AzureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tags client options")
	}
	factory, err := armresources.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armresources client factory")
	}
	return &AzureClient{factory.NewTagsClient()}, nil
}

// GetAtScope sends the get at scope request.
func (ac *AzureClient) GetAtScope(ctx context.Context, scope string) (armresources.TagsResource, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "tags.AzureClient.GetAtScope")
	defer done()

	resp, err := ac.tags.GetAtScope(ctx, scope, nil)
	if err != nil {
		return armresources.TagsResource{}, err
	}

	return resp.TagsResource, nil
}

// UpdateAtScope this operation allows replacing, merging or selectively deleting tags on the specified resource or
// subscription.
func (ac *AzureClient) UpdateAtScope(ctx context.Context, scope string, parameters armresources.TagsPatchResource) (armresources.TagsResource, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "tags.AzureClient.UpdateAtScope")
	defer done()

	resp, err := ac.tags.UpdateAtScope(ctx, scope, parameters, nil)
	if err != nil {
		return armresources.TagsResource{}, err
	}

	return resp.TagsResource, nil
}
