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

package agentpools

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// azureClient contains the Azure go-sdk Client.
type azureClient struct {
	agentpools *armcontainerservice.AgentPoolsClient
}

// newClient creates a new agentpools client from an authorizer.
func newClient(scope AgentPoolScope) (*azureClient, error) {
	var headers map[string]string
	if customHeaders, ok := scope.AgentPoolSpec().(azure.ResourceSpecGetterWithHeaders); ok {
		headers = customHeaders.CustomHeaders()
	}
	opts, err := azure.ARMClientOptions(scope.CloudEnvironment(), azure.CustomPutPatchHeaderPolicy{Headers: headers})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create agentpools client options")
	}
	factory, err := armcontainerservice.NewClientFactory(scope.SubscriptionID(), scope.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armcontainerservice client factory")
	}
	return &azureClient{factory.NewAgentPoolsClient()}, nil
}

// Get gets an agent pool.
func (ac *azureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "agentpools.azureClient.Get")
	defer done()

	resp, err := ac.agentpools.Get(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.AgentPool, nil
}

// CreateOrUpdateAsync creates or updates an agent pool asynchronously.
// It sends a PUT request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *azureClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (
	result interface{}, poller *runtime.Poller[armcontainerservice.AgentPoolsClientCreateOrUpdateResponse], err error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "agentpools.azureClient.CreateOrUpdate")
	defer done()

	var agentPool armcontainerservice.AgentPool
	if parameters != nil {
		ap, ok := parameters.(armcontainerservice.AgentPool)
		if !ok {
			return nil, nil, errors.Errorf("%T is not an armcontainerservice.AgentPool", parameters)
		}
		agentPool = ap
	}

	opts := &armcontainerservice.AgentPoolsClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken}
	log.V(4).Info("sending request", "resumeToken", resumeToken)
	poller, err = ac.agentpools.BeginCreateOrUpdate(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), agentPool, opts)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureCallTimeout)
	defer cancel()

	pollOpts := &runtime.PollUntilDoneOptions{Frequency: async.DefaultPollerFrequency}
	resp, err := poller.PollUntilDone(ctx, pollOpts)
	if err != nil {
		// If an error occurs, return the poller.
		return nil, poller, err
	}

	// if the operation completed, return a nil poller
	return resp.AgentPool, nil, err
}

// DeleteAsync deletes an agent pool asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *azureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (
	poller *runtime.Poller[armcontainerservice.AgentPoolsClientDeleteResponse], err error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "agentpools.azureClient.DeleteAsync")
	defer done()

	opts := &armcontainerservice.AgentPoolsClientBeginDeleteOptions{ResumeToken: resumeToken}
	log.V(4).Info("sending request", "resumeToken", resumeToken)
	poller, err = ac.agentpools.BeginDelete(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureCallTimeout)
	defer cancel()

	pollOpts := &runtime.PollUntilDoneOptions{Frequency: async.DefaultPollerFrequency}
	_, err = poller.PollUntilDone(ctx, pollOpts)
	if err != nil {
		// If an error occurs, return the poller.
		// This means the long-running operation didn't finish in the specified timeout.
		return poller, err
	}

	// if the operation completed, return a nil poller.
	return nil, err
}
