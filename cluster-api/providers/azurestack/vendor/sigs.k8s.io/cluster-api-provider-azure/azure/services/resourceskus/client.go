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

package resourceskus

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// Client wraps go-sdk.
type Client interface {
	List(context.Context, string) ([]armcompute.ResourceSKU, error)
}

// AzureClient contains the Azure go-sdk Client.
type AzureClient struct {
	skus *armcompute.ResourceSKUsClient
}

var _ Client = &AzureClient{}

// NewClient creates a new Resource SKUs client from an authorizer.
func NewClient(auth azure.Authorizer) (*AzureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create resourceskus client options")
	}
	factory, err := armcompute.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armcompute client factory")
	}
	return &AzureClient{factory.NewResourceSKUsClient()}, nil
}

// List returns all Resource SKUs available to the subscription.
func (ac *AzureClient) List(ctx context.Context, filter string) ([]armcompute.ResourceSKU, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "resourceskus.AzureClient.List")
	defer done()

	var skus []armcompute.ResourceSKU
	opts := armcompute.ResourceSKUsClientListOptions{Filter: &filter}
	pager := ac.skus.NewListPager(&opts)
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return skus, errors.Wrap(err, "could not iterate resource skus")
		}
		for _, sku := range resp.Value {
			skus = append(skus, *sku)
		}
	}

	return skus, nil
}
