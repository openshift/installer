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

package publicips

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
	publicips      *armnetwork.PublicIPAddressesClient
	apiCallTimeout time.Duration
}

// NewClient creates a new public IP client from an authorizer.
func NewClient(auth azure.Authorizer, apiCallTimeout time.Duration) (*AzureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create publicips client options")
	}

	factory, err := armnetwork.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create publicips client factory")
	}
	return &AzureClient{factory.NewPublicIPAddressesClient(), apiCallTimeout}, nil
}

// Get gets the specified public IP address in a specified resource group.
func (ac *AzureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "publicips.AzureClient.Get")
	defer done()

	resp, err := ac.publicips.Get(ctx, spec.ResourceGroupName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.PublicIPAddress, nil
}

// CreateOrUpdateAsync creates or updates a static or dynamic public IP address.
// It sends a PUT request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *AzureClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armnetwork.PublicIPAddressesClientCreateOrUpdateResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "publicips.AzureClient.CreateOrUpdate")
	defer done()

	publicip, ok := parameters.(armnetwork.PublicIPAddress)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armnetwork.PublicIPAddress", parameters)
	}

	opts := &armnetwork.PublicIPAddressesClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken}
	poller, err = ac.publicips.BeginCreateOrUpdate(ctx, spec.ResourceGroupName(), spec.ResourceName(), publicip, opts)
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

	// if the operation completed, return a nil poller.
	return resp.PublicIPAddress, nil, err
}

// DeleteAsync deletes the specified public IP address asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *AzureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armnetwork.PublicIPAddressesClientDeleteResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "publicips.AzureClient.DeleteAsync")
	defer done()

	opts := &armnetwork.PublicIPAddressesClientBeginDeleteOptions{ResumeToken: resumeToken}
	poller, err = ac.publicips.BeginDelete(ctx, spec.ResourceGroupName(), spec.ResourceName(), opts)
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
