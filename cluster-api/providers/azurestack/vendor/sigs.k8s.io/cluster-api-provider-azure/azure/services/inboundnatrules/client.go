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

package inboundnatrules

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

// client wraps go-sdk.
type client interface {
	List(context.Context, string, string) (result []armnetwork.InboundNatRule, err error)
}

// azureClient contains the Azure go-sdk Client.
type azureClient struct {
	inboundnatrules *armnetwork.InboundNatRulesClient
	apiCallTimeout  time.Duration
}

var _ client = (*azureClient)(nil)

// newClient creates a new inbound NAT rules client from an authorizer.
func newClient(auth azure.Authorizer, apiCallTimeout time.Duration) (*azureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create inboundnatrules client options")
	}
	factory, err := armnetwork.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armnetwork client factory")
	}
	return &azureClient{factory.NewInboundNatRulesClient(), apiCallTimeout}, nil
}

// Get gets the specified inbound NAT rules.
func (ac *azureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "inboundnatrules.azureClient.Get")
	defer done()

	resp, err := ac.inboundnatrules.Get(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.InboundNatRule, nil
}

// List returns all inbound NAT rules on a load balancer.
func (ac *azureClient) List(ctx context.Context, resourceGroupName, lbName string) (result []armnetwork.InboundNatRule, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "inboundnatrules.azureClient.List")
	defer done()

	var natRules []armnetwork.InboundNatRule
	pager := ac.inboundnatrules.NewListPager(resourceGroupName, lbName, nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return natRules, errors.Wrap(err, "could not iterate inbound NAT rules")
		}
		for _, natRule := range nextResult.Value {
			natRules = append(natRules, *natRule)
		}
	}

	return natRules, nil
}

// CreateOrUpdateAsync creates or updates an inbound NAT rule asynchronously.
// It sends a PUT request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *azureClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armnetwork.InboundNatRulesClientCreateOrUpdateResponse], err error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "inboundnatrules.azureClient.CreateOrUpdateAsync")
	defer done()

	natRule, ok := parameters.(armnetwork.InboundNatRule)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armnetwork.InboundNatRule", parameters)
	}

	opts := &armnetwork.InboundNatRulesClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken}
	log.V(4).Info("sending request", "resumeToken", resumeToken)
	poller, err = ac.inboundnatrules.BeginCreateOrUpdate(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), natRule, opts)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ac.apiCallTimeout)
	defer cancel()

	pollOpts := &runtime.PollUntilDoneOptions{Frequency: async.DefaultPollerFrequency}
	resp, err := poller.PollUntilDone(ctx, pollOpts)
	if err != nil {
		// If an error occurs, return the poller.
		// This means the long-running operation didn't finish in the specified timeout.
		return nil, poller, err
	}

	// if the operation completed, return a nil poller
	return resp.InboundNatRule, nil, err
}

// DeleteAsync deletes an inbound NAT rule asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *azureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armnetwork.InboundNatRulesClientDeleteResponse], err error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "inboundnatrules.azureClient.DeleteAsync")
	defer done()

	opts := &armnetwork.InboundNatRulesClientBeginDeleteOptions{ResumeToken: resumeToken}
	log.V(4).Info("sending request", "resumeToken", resumeToken)
	poller, err = ac.inboundnatrules.BeginDelete(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), opts)
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
