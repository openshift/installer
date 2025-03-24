/*
Copyright 2021 The Kubernetes Authors.

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

package vnetpeerings

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// AzureClient contains the Azure go-sdk Client.
type AzureClient struct {
	peerings       *armnetwork.VirtualNetworkPeeringsClient
	apiCallTimeout time.Duration
}

// NewClient creates a new virtual network peerings client from an authorizer.
func NewClient(auth azure.Authorizer, apiCallTimeout time.Duration) (*AzureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create vnetpeerings client options")
	}
	factory, err := armnetwork.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armnetwork client factory")
	}
	return &AzureClient{factory.NewVirtualNetworkPeeringsClient(), apiCallTimeout}, nil
}

// Get gets the specified virtual network peering by the peering name, virtual network, and resource group.
func (ac *AzureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "vnetpeerings.AzureClient.Get")
	defer done()

	resp, err := ac.peerings.Get(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.VirtualNetworkPeering, nil
}

// CreateOrUpdateAsync creates or updates a virtual network peering asynchronously.
// It sends a PUT request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *AzureClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armnetwork.VirtualNetworkPeeringsClientCreateOrUpdateResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "vnetpeerings.AzureClient.CreateOrUpdateAsync")
	defer done()

	peering, ok := parameters.(armnetwork.VirtualNetworkPeering)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armnetwork.VirtualNetworkPeering", parameters)
	}

	opts := &armnetwork.VirtualNetworkPeeringsClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken}
	poller, err = ac.peerings.BeginCreateOrUpdate(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), peering, opts)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ac.apiCallTimeout)
	defer cancel()

	pollOpts := &runtime.PollUntilDoneOptions{Frequency: async.DefaultPollerFrequency}
	resp, err := poller.PollUntilDone(ctx, pollOpts)
	if err != nil {
		// if an error occurs, return the poller.
		// this means the long-running operation didn't finish in the specified timeout.
		return nil, poller, err
	}

	// if the operation completed, return a nil poller
	return resp.VirtualNetworkPeering, nil, err
}

// DeleteAsync deletes a virtual network peering asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a Future which can be used to track the ongoing
// progress of the operation.
func (ac *AzureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armnetwork.VirtualNetworkPeeringsClientDeleteResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "vnetpeerings.AzureClient.DeleteAsync")
	defer done()

	opts := &armnetwork.VirtualNetworkPeeringsClientBeginDeleteOptions{ResumeToken: resumeToken}
	poller, err = ac.peerings.BeginDelete(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ac.apiCallTimeout)
	defer cancel()

	pollOpts := &runtime.PollUntilDoneOptions{Frequency: async.DefaultPollerFrequency}
	_, err = poller.PollUntilDone(ctx, pollOpts)
	if err != nil {
		// if an error occurs, return the poller.
		// this means the long-running operation didn't finish in the specified timeout.
		return poller, err
	}

	// if the operation completed, return a nil poller.
	return nil, err
}
