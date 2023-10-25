/*
Copyright 2019 The Kubernetes Authors.

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

package subnets

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// AzureClient contains the Azure go-sdk Client.
type AzureClient struct {
	subnets *armnetwork.SubnetsClient
}

// NewClient creates a new subnets client from an authorizer.
func NewClient(auth azure.Authorizer) (*AzureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create subnets client options")
	}
	factory, err := armnetwork.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armnetwork client factory")
	}
	return &AzureClient{factory.NewSubnetsClient()}, nil
}

// Get gets the specified subnet by virtual network and resource group.
func (ac *AzureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "subnets.AzureClient.Get")
	defer done()

	resp, err := ac.subnets.Get(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.Subnet, nil
}

// CreateOrUpdateAsync creates or updates a subnet asynchronously.
// It sends a PUT request to Azure and if accepted without error, the func will return a poller which can be used to track the ongoing
// progress of the operation.
func (ac *AzureClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armnetwork.SubnetsClientCreateOrUpdateResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "subnets.AzureClient.CreateOrUpdateAsync")
	defer done()

	subnet, ok := parameters.(armnetwork.Subnet)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armnetwork.Subnet", parameters)
	}

	opts := &armnetwork.SubnetsClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken}
	poller, err = ac.subnets.BeginCreateOrUpdate(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), subnet, opts)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureCallTimeout)
	defer cancel()

	pollOpts := &runtime.PollUntilDoneOptions{Frequency: async.DefaultPollerFrequency}
	resp, err := poller.PollUntilDone(ctx, pollOpts)
	if err != nil {
		// if an error occurs, return the poller.
		// this means the long-running operation didn't finish in the specified timeout.
		return nil, poller, err
	}

	// if the operation completed, return a nil poller
	return resp.Subnet, nil, err
}

// DeleteAsync deletes a subnet asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a poller which can be used to track the ongoing
// progress of the operation.
func (ac *AzureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armnetwork.SubnetsClientDeleteResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "subnets.AzureClient.DeleteAsync")
	defer done()

	opts := &armnetwork.SubnetsClientBeginDeleteOptions{ResumeToken: resumeToken}
	poller, err = ac.subnets.BeginDelete(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureCallTimeout)
	defer cancel()

	pollOpts := &runtime.PollUntilDoneOptions{Frequency: async.DefaultPollerFrequency}
	_, err = poller.PollUntilDone(ctx, pollOpts)
	if err != nil {
		// if an error occurs, return the Poller.
		// this means the long-running operation didn't finish in the specified timeout.
		return poller, err
	}

	// if the operation completed, return a nil poller.
	return nil, err
}
