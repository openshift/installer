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

package scalesetvms

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// client wraps go-sdk.
type client interface {
	Get(context.Context, azure.ResourceSpecGetter) (interface{}, error)
	CreateOrUpdateAsync(context.Context, azure.ResourceSpecGetter, string, interface{}) (interface{}, *runtime.Poller[armcompute.VirtualMachineScaleSetVMsClientUpdateResponse], error)
	DeleteAsync(context.Context, azure.ResourceSpecGetter, string) (*runtime.Poller[armcompute.VirtualMachineScaleSetVMsClientDeleteResponse], error)
}

// azureClient contains the Azure go-sdk Client.
type azureClient struct {
	scalesetvms    *armcompute.VirtualMachineScaleSetVMsClient
	apiCallTimeout time.Duration
}

var _ client = &azureClient{}

// newClient creates a VMSS client from an authorizer.
func newClient(auth azure.Authorizer, apiCallTimeout time.Duration) (*azureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create scalesetvms client options")
	}
	factory, err := armcompute.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armcompute client factory")
	}
	return &azureClient{factory.NewVirtualMachineScaleSetVMsClient(), apiCallTimeout}, nil
}

// Get retrieves the Virtual Machine Scale Set Virtual Machine.
func (ac *azureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (interface{}, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scalesetvms.azureClient.Get")
	defer done()

	resp, err := ac.scalesetvms.Get(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.VirtualMachineScaleSetVM, nil
}

// CreateOrUpdateAsync is a dummy implementation to fulfill the async.Reconciler interface.
func (ac *azureClient) CreateOrUpdateAsync(ctx context.Context, _ azure.ResourceSpecGetter, _ string, _ interface{}) (result interface{}, poller *runtime.Poller[armcompute.VirtualMachineScaleSetVMsClientUpdateResponse], err error) {
	_, _, done := tele.StartSpanWithLogger(ctx, "scalesets.AzureClient.CreateOrUpdateAsync")
	defer done()

	return nil, nil, nil
}

// DeleteAsync deletes a virtual machine scale set instance asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *azureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armcompute.VirtualMachineScaleSetVMsClientDeleteResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scalesetvms.AzureClient.DeleteAsync")
	defer done()

	opts := &armcompute.VirtualMachineScaleSetVMsClientBeginDeleteOptions{ResumeToken: resumeToken}
	poller, err = ac.scalesetvms.BeginDelete(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName(), opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ac.apiCallTimeout)
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
