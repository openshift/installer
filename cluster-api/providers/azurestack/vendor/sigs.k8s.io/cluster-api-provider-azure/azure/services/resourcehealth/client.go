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

package resourcehealth

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcehealth/armresourcehealth"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// client wraps go-sdk.
type client interface {
	GetByResource(context.Context, string) (armresourcehealth.AvailabilityStatus, error)
}

// azureClient contains the Azure go-sdk Client.
type azureClient struct {
	availabilityStatuses *armresourcehealth.AvailabilityStatusesClient
}

// newClient creates a new resource health client from an authorizer.
func newClient(auth azure.Authorizer) (*azureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create resourcehealth client options")
	}
	factory, err := armresourcehealth.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armresourcehealth client factory")
	}
	return &azureClient{factory.NewAvailabilityStatusesClient()}, nil
}

// GetByResource gets the availability status for the specified resource.
func (ac *azureClient) GetByResource(ctx context.Context, resourceURI string) (armresourcehealth.AvailabilityStatus, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "resourcehealth.AzureClient.GetByResource")
	defer done()

	opts := &armresourcehealth.AvailabilityStatusesClientGetByResourceOptions{}
	resp, err := ac.availabilityStatuses.GetByResource(ctx, resourceURI, opts)
	if err != nil {
		return armresourcehealth.AvailabilityStatus{}, err
	}
	return resp.AvailabilityStatus, nil
}
