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

package virtualnetworks

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

// azureClient contains the Azure go-sdk Client.
type azureClient struct {
	virtualnetworks *armnetwork.VirtualNetworksClient
}

// newClient creates a new virtual networks client from an authorizer.
func newClient(auth azure.Authorizer) (*azureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create virtualnetworks client options")
	}
	factory, err := armnetwork.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armnetwork client factory")
	}
	return &azureClient{factory.NewVirtualNetworksClient()}, nil
}

// Get gets the specified virtual network.
func (ac *azureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "virtualnetworks.azureClient.Get")
	defer done()

	resp, err := ac.virtualnetworks.Get(ctx, spec.ResourceGroupName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.VirtualNetwork, nil
}

// CreateOrUpdateAsync creates or updates a virtual network in the specified resource group asynchronously.
// It sends a PUT request to Azure and if accepted without error, the func will return a poller which can be used to track the ongoing
// progress of the operation.
func (ac *azureClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armnetwork.VirtualNetworksClientCreateOrUpdateResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "virtualnetworks.azureClient.CreateOrUpdateAsync")
	defer done()

	vn, ok := parameters.(armnetwork.VirtualNetwork)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armnetwork.VirtualNetwork", parameters)
	}

	opts := &armnetwork.VirtualNetworksClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken}
	poller, err = ac.virtualnetworks.BeginCreateOrUpdate(ctx, spec.ResourceGroupName(), spec.ResourceName(), vn, opts)
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
	return resp.VirtualNetwork, nil, err
}

// DeleteAsync deletes a virtual network asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a poller which can be used to track the ongoing
// progress of the operation.
func (ac *azureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armnetwork.VirtualNetworksClientDeleteResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "virtualnetworks.azureClient.DeleteAsync")
	defer done()

	opts := &armnetwork.VirtualNetworksClientBeginDeleteOptions{ResumeToken: resumeToken}
	poller, err = ac.virtualnetworks.BeginDelete(ctx, spec.ResourceGroupName(), spec.ResourceName(), opts)
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
