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

package availabilitysets

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// AzureClient contains the Azure go-sdk Client.
type AzureClient struct {
	availabilitySets *armcompute.AvailabilitySetsClient
}

// NewClient creates a new availability sets client from an authorizer.
func NewClient(auth azure.Authorizer) (*AzureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create availabilitysets client options")
	}
	factory, err := armcompute.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armcompute client factory")
	}
	return &AzureClient{factory.NewAvailabilitySetsClient()}, nil
}

// Get gets an availability set.
func (ac *AzureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "availabilitysets.AzureClient.Get")
	defer done()

	resp, err := ac.availabilitySets.Get(ctx, spec.ResourceGroupName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.AvailabilitySet, nil
}

// CreateOrUpdateAsync creates or updates an availability set asynchronously.
// It sends a PUT request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *AzureClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, _resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armcompute.AvailabilitySetsClientCreateOrUpdateResponse], err error) { //nolint:revive // keeping _resumeToken for understanding purposes
	ctx, _, done := tele.StartSpanWithLogger(ctx, "availabilitySets.AzureClient.CreateOrUpdateAsync")
	defer done()

	availabilitySet, ok := parameters.(armcompute.AvailabilitySet)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armcompute.AvailabilitySet", parameters)
	}

	// Note: there is no async `BeginCreateOrUpdate` implementation for availability sets, so this func will never return a poller.
	resp, err := ac.availabilitySets.CreateOrUpdate(ctx, spec.ResourceGroupName(), spec.ResourceName(), availabilitySet, nil)
	if err != nil {
		return nil, nil, err
	}

	// if the operation completed, return a nil poller
	return resp.AvailabilitySet, nil, err
}

// DeleteAsync deletes a availability set asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *AzureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, _resumeToken string) (poller *runtime.Poller[armcompute.AvailabilitySetsClientDeleteResponse], err error) { //nolint:revive // keeping _resumeToken for understanding purposes
	ctx, _, done := tele.StartSpanWithLogger(ctx, "availabilitysets.AzureClient.DeleteAsync")
	defer done()

	// Note: there is no async `BeginDelete` implementation for availability sets, so this func will never return a poller.
	_, err = ac.availabilitySets.Delete(ctx, spec.ResourceGroupName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}

	// if the operation completed, return a nil poller.
	return nil, err
}
