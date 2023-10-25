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

package disks

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// azureClient contains the Azure go-sdk Client.
type azureClient struct {
	disks *armcompute.DisksClient
}

// newClient creates a new disks client from an authorizer.
func newClient(auth azure.Authorizer) (*azureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create disks client options")
	}
	factory, err := armcompute.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armcompute client factory")
	}
	return &azureClient{factory.NewDisksClient()}, nil
}

// DeleteAsync deletes a disk asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *azureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armcompute.DisksClientDeleteResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "disks.azureClient.DeleteAsync")
	defer done()

	opts := &armcompute.DisksClientBeginDeleteOptions{ResumeToken: resumeToken}
	poller, err = ac.disks.BeginDelete(ctx, spec.ResourceGroupName(), spec.ResourceName(), opts)
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
