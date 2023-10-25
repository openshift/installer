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

package bastionhosts

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
	bastionhosts *armnetwork.BastionHostsClient
}

// newClient creates a new bastion hosts client from an authorizer.
func newClient(auth azure.Authorizer) (*azureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bastionhosts client options")
	}
	factory, err := armnetwork.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armnetwork client factory")
	}
	return &azureClient{factory.NewBastionHostsClient()}, nil
}

// Get gets the specified bastion host.
func (ac *azureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "bastionhosts.azureClient.Get")
	defer done()

	resp, err := ac.bastionhosts.Get(ctx, spec.ResourceGroupName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.BastionHost, nil
}

// CreateOrUpdateAsync creates or updates a bastion host asynchronously.
// It sends a PUT request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *azureClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armnetwork.BastionHostsClientCreateOrUpdateResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "bastionhosts.azureClient.CreateOrUpdateAsync")
	defer done()

	host, ok := parameters.(armnetwork.BastionHost)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armnetwork.BastionHost", parameters)
	}

	opts := &armnetwork.BastionHostsClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken}
	poller, err = ac.bastionhosts.BeginCreateOrUpdate(ctx, spec.ResourceGroupName(), spec.ResourceName(), host, opts)
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
	return resp.BastionHost, nil, err
}

// DeleteAsync deletes a bastion host asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *azureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armnetwork.BastionHostsClientDeleteResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "bastionhosts.azureClient.Delete")
	defer done()

	opts := &armnetwork.BastionHostsClientBeginDeleteOptions{ResumeToken: resumeToken}
	poller, err = ac.bastionhosts.BeginDelete(ctx, spec.ResourceGroupName(), spec.ResourceName(), opts)
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
