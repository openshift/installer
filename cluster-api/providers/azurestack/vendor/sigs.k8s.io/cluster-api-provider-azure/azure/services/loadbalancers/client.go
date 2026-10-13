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

package loadbalancers

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// azureClient contains the Azure go-sdk Client.
type azureClient struct {
	loadbalancers  *armnetwork.LoadBalancersClient
	auth           azure.Authorizer
	apiCallTimeout time.Duration
}

// newClient creates a new load balancer client from an authorizer.
func newClient(auth azure.Authorizer, apiCallTimeout time.Duration) (*azureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get load balancer client options")
	}

	factory, err := armnetwork.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armnetwork client factory")
	}
	return &azureClient{factory.NewLoadBalancersClient(), auth, apiCallTimeout}, nil
}

// Get gets the specified load balancer.
func (ac *azureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "loadbalancers.azureClient.Get")
	defer done()

	resp, err := ac.loadbalancers.Get(ctx, spec.ResourceGroupName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}

	return resp.LoadBalancer, nil
}

// CreateOrUpdateAsync creates or updates a load balancer asynchronously.
// It sends a PUT request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *azureClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armnetwork.LoadBalancersClientCreateOrUpdateResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "loadbalancers.azureClient.CreateOrUpdate")
	defer done()

	loadBalancer, ok := parameters.(armnetwork.LoadBalancer)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armnetwork.LoadBalancer", parameters)
	}

	var extraPolicies []policy.Policy
	if loadBalancer.Etag != nil {
		extraPolicies = append(extraPolicies, azure.CustomPutPatchHeaderPolicy{
			Headers: map[string]string{
				"If-Match": *loadBalancer.Etag,
			},
		})
	}

	// Create a new client that knows how to add etag headers to the request.
	clientOpts, err := azure.ARMClientOptions(ac.auth.CloudEnvironment(), ac.auth.BaseURI(), extraPolicies...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create loadbalancer client options")
	}

	factory, err := armnetwork.NewClientFactory(ac.auth.SubscriptionID(), ac.auth.Token(), clientOpts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create armnetwork client factory")
	}

	client := factory.NewLoadBalancersClient()
	opts := &armnetwork.LoadBalancersClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken}
	poller, err = client.BeginCreateOrUpdate(ctx, spec.ResourceGroupName(), spec.ResourceName(), loadBalancer, opts)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ac.apiCallTimeout)
	defer cancel()

	pollOpts := &runtime.PollUntilDoneOptions{Frequency: async.DefaultPollerFrequency}
	resp, err := poller.PollUntilDone(ctx, pollOpts)
	if err != nil {
		// if an error occurs, return the poller.
		// This means the long-running operation didn't finish in the specified timeout.
		return nil, poller, err
	}

	// if the operation completed, return nil poller.
	return resp.LoadBalancer, nil, err
}

// DeleteAsync deletes a load balancer asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *azureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armnetwork.LoadBalancersClientDeleteResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "loadbalancers.azureClient.Delete")
	defer done()

	opts := &armnetwork.LoadBalancersClientBeginDeleteOptions{ResumeToken: resumeToken}
	poller, err = ac.loadbalancers.BeginDelete(ctx, spec.ResourceGroupName(), spec.ResourceName(), opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ac.apiCallTimeout)
	defer cancel()

	pollOpts := &runtime.PollUntilDoneOptions{Frequency: async.DefaultPollerFrequency}
	_, err = poller.PollUntilDone(ctx, pollOpts)
	if err != nil {
		// if error occurs, return the poller.
		// this means the long-running operation didn't finish in the specified timeout.
		return poller, err
	}
	// if the operation completed, return nil poller.
	return nil, err
}
