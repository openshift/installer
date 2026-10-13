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

package privatedns

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// azureVirtualNetworkLinksClient contains the Azure go-sdk Client for virtual network links.
type azureVirtualNetworkLinksClient struct {
	vnetlinks      *armprivatedns.VirtualNetworkLinksClient
	apiCallTimeout time.Duration
}

// newVirtualNetworkLinksClient creates a virtual network links client from an authorizer.
func newVirtualNetworkLinksClient(auth azure.Authorizer, apiCallTimeout time.Duration) (*azureVirtualNetworkLinksClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create virtualnetworkslink client options")
	}
	factory, err := armprivatedns.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armprivatedns client factory")
	}
	return &azureVirtualNetworkLinksClient{factory.NewVirtualNetworkLinksClient(), apiCallTimeout}, nil
}

// Get gets the specified virtual network link.
func (avc *azureVirtualNetworkLinksClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "privatedns.azureVirtualNetworkLinksClient.Get")
	defer done()

	resp, err := avc.vnetlinks.Get(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.VirtualNetworkLink, nil
}

// CreateOrUpdateAsync creates or updates a virtual network link asynchronously.
// It sends a PUT request to Azure and if accepted without error, the func will return a poller which can be used to track the ongoing
// progress of the operation.
func (avc *azureVirtualNetworkLinksClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armprivatedns.VirtualNetworkLinksClientCreateOrUpdateResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "privatedns.azureVirtualNetworkLinksClient.CreateOrUpdateAsync")
	defer done()

	link, ok := parameters.(armprivatedns.VirtualNetworkLink)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armprivatedns.VirtualNetworkLink", parameters)
	}

	opts := &armprivatedns.VirtualNetworkLinksClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken}
	poller, err = avc.vnetlinks.BeginCreateOrUpdate(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), link, opts)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, avc.apiCallTimeout)
	defer cancel()

	pollOpts := &runtime.PollUntilDoneOptions{Frequency: async.DefaultPollerFrequency}
	resp, err := poller.PollUntilDone(ctx, pollOpts)
	if err != nil {
		// if an error occurs, return the poller.
		// this means the long-running operation didn't finish in the specified timeout.
		return nil, poller, err
	}

	// if the operation completed, return a nil poller
	return resp.VirtualNetworkLink, nil, err
}

// DeleteAsync deletes a virtual network link asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a poller which can be used to track the ongoing
// progress of the operation.
func (avc *azureVirtualNetworkLinksClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armprivatedns.VirtualNetworkLinksClientDeleteResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "privatedns.azureVirtualNetworkLinksClient.DeleteAsync")
	defer done()

	opts := &armprivatedns.VirtualNetworkLinksClientBeginDeleteOptions{ResumeToken: resumeToken}
	poller, err = avc.vnetlinks.BeginDelete(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, avc.apiCallTimeout)
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
