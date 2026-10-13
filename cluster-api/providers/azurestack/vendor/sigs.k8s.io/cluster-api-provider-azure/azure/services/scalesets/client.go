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

package scalesets

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

// Client wraps go-sdk.
type Client interface {
	Get(context.Context, azure.ResourceSpecGetter) (interface{}, error)
	List(context.Context, string) ([]armcompute.VirtualMachineScaleSet, error)
	ListInstances(context.Context, string, string) ([]armcompute.VirtualMachineScaleSetVM, error)

	CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armcompute.VirtualMachineScaleSetsClientCreateOrUpdateResponse], err error)
	DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armcompute.VirtualMachineScaleSetsClientDeleteResponse], err error)
}

// AzureClient contains the Azure go-sdk Client.
type AzureClient struct {
	scalesetvms    *armcompute.VirtualMachineScaleSetVMsClient
	scalesets      *armcompute.VirtualMachineScaleSetsClient
	apiCallTimeout time.Duration
}

var _ Client = &AzureClient{}

// NewClient creates a new VMSS client from an authorizer.
func NewClient(auth azure.Authorizer, apiCallTimeout time.Duration) (*AzureClient, error) {
	scaleSetVMsClient, err := newVirtualMachineScaleSetVMsClient(auth)
	if err != nil {
		return nil, err
	}
	scaleSetsClient, err := newVirtualMachineScaleSetsClient(auth)
	if err != nil {
		return nil, err
	}
	return &AzureClient{
		scalesetvms:    scaleSetVMsClient,
		scalesets:      scaleSetsClient,
		apiCallTimeout: apiCallTimeout,
	}, nil
}

// newVirtualMachineScaleSetVMsClient creates a vmss VM client from an authorizer.
func newVirtualMachineScaleSetVMsClient(auth azure.Authorizer) (*armcompute.VirtualMachineScaleSetVMsClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create scalesetvms client options")
	}
	factory, err := armcompute.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armcompute client factory")
	}
	return factory.NewVirtualMachineScaleSetVMsClient(), nil
}

// newVirtualMachineScaleSetsClient creates a vmss client from an authorizer.
func newVirtualMachineScaleSetsClient(auth azure.Authorizer) (*armcompute.VirtualMachineScaleSetsClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create scalesets client options")
	}
	factory, err := armcompute.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armcompute client factory")
	}
	return factory.NewVirtualMachineScaleSetsClient(), nil
}

// ListInstances retrieves information about the model views of a virtual machine scale set.
func (ac *AzureClient) ListInstances(ctx context.Context, resourceGroupName string, resourceName string) ([]armcompute.VirtualMachineScaleSetVM, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scalesets.AzureClient.ListInstances")
	defer done()

	var instances []armcompute.VirtualMachineScaleSetVM
	pager := ac.scalesetvms.NewListPager(resourceGroupName, resourceName, nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "could not iterate scalesetvms")
		}
		for _, scaleSetVM := range nextResult.Value {
			instances = append(instances, *scaleSetVM)
		}
	}

	return instances, nil
}

// List returns all scale sets in a resource group.
func (ac *AzureClient) List(ctx context.Context, resourceGroupName string) ([]armcompute.VirtualMachineScaleSet, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scalesets.AzureClient.List")
	defer done()

	var scaleSets []armcompute.VirtualMachineScaleSet
	pager := ac.scalesets.NewListPager(resourceGroupName, nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return scaleSets, errors.Wrap(err, "could not iterate scalesets")
		}
		for _, scaleSet := range nextResult.Value {
			scaleSets = append(scaleSets, *scaleSet)
		}
	}

	return scaleSets, nil
}

// Get retrieves information about the model view of a virtual machine scale set.
func (ac *AzureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (interface{}, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scalesets.AzureClient.Get")
	defer done()

	resp, err := ac.scalesets.Get(ctx, spec.ResourceGroupName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.VirtualMachineScaleSet, nil
}

// CreateOrUpdateAsync creates or updates a virtual machine scale set asynchronously.
// It sends a PUT request to Azure and if accepted without error, the func will return a poller which can be used to track the ongoing
// progress of the operation.
func (ac *AzureClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armcompute.VirtualMachineScaleSetsClientCreateOrUpdateResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scalesets.AzureClient.CreateOrUpdateAsync")
	defer done()

	scaleset, ok := parameters.(armcompute.VirtualMachineScaleSet)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armcompute.VirtualMachineScaleSet", parameters)
	}

	opts := &armcompute.VirtualMachineScaleSetsClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken}
	poller, err = ac.scalesets.BeginCreateOrUpdate(ctx, spec.ResourceGroupName(), spec.ResourceName(), scaleset, opts)
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
	return resp.VirtualMachineScaleSet, nil, err
}

// DeleteAsync is the operation to delete a virtual machine scale set asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a poller which can be used to track the ongoing
// progress of the operation.
//
// Parameters:
//
//	spec - The ResourceSpecGetter containing used for name and resource group of the virtual machine scale set.
func (ac *AzureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armcompute.VirtualMachineScaleSetsClientDeleteResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scalesets.AzureClient.DeleteAsync")
	defer done()

	opts := &armcompute.VirtualMachineScaleSetsClientBeginDeleteOptions{ResumeToken: resumeToken}
	poller, err = ac.scalesets.BeginDelete(ctx, spec.ResourceGroupName(), spec.ResourceName(), opts)
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
