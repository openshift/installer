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

// azureZonesClient contains the Azure go-sdk Client for private dns zone.
type azureZonesClient struct {
	privatezones   *armprivatedns.PrivateZonesClient
	apiCallTimeout time.Duration
}

// newPrivateZonesClient creates a private zones client from an authorizer.
func newPrivateZonesClient(auth azure.Authorizer, apiCallTimeout time.Duration) (*azureZonesClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create privatezones client options")
	}
	factory, err := armprivatedns.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armprivatedns client factory")
	}
	return &azureZonesClient{factory.NewPrivateZonesClient(), apiCallTimeout}, nil
}

// Get gets the specified private dns zone.
func (azc *azureZonesClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "privatedns.azureZonesClient.Get")
	defer done()

	resp, err := azc.privatezones.Get(ctx, spec.ResourceGroupName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.PrivateZone, nil
}

// CreateOrUpdateAsync creates or updates a private dns zone asynchronously.
// It sends a PUT request to Azure and if accepted without error, the func will return a poller which can be used to track the ongoing
// progress of the operation.
func (azc *azureZonesClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armprivatedns.PrivateZonesClientCreateOrUpdateResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "privatedns.azureZonesClient.CreateOrUpdateAsync")
	defer done()

	zone, ok := parameters.(armprivatedns.PrivateZone)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armprivatedns.PrivateZone", parameters)
	}

	opts := &armprivatedns.PrivateZonesClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken}
	poller, err = azc.privatezones.BeginCreateOrUpdate(ctx, spec.ResourceGroupName(), spec.ResourceName(), zone, opts)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, azc.apiCallTimeout)
	defer cancel()

	pollOpts := &runtime.PollUntilDoneOptions{Frequency: async.DefaultPollerFrequency}
	resp, err := poller.PollUntilDone(ctx, pollOpts)
	if err != nil {
		// if an error occurs, return the poller.
		// this means the long-running operation didn't finish in the specified timeout.
		return nil, poller, err
	}

	// if the operation completed, return a nil poller
	return resp.PrivateZone, nil, err
}

// DeleteAsync deletes a private dns zone asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a poller which can be used to track the ongoing
// progress of the operation.
func (azc *azureZonesClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armprivatedns.PrivateZonesClientDeleteResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "privatedns.azureZonesClient.DeleteAsync")
	defer done()

	opts := &armprivatedns.PrivateZonesClientBeginDeleteOptions{ResumeToken: resumeToken}
	poller, err = azc.privatezones.BeginDelete(ctx, spec.ResourceGroupName(), spec.ResourceName(), opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, azc.apiCallTimeout)
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
