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

package virtualmachines

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

type (
	// AzureClient contains the Azure go-sdk Client.
	AzureClient struct {
		virtualmachines *armcompute.VirtualMachinesClient
		apiCallTimeout  time.Duration
	}

	// Client provides operations on Azure virtual machine resources.
	Client interface {
		Get(context.Context, azure.ResourceSpecGetter) (interface{}, error)
		CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armcompute.VirtualMachinesClientCreateOrUpdateResponse], err error)
		DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armcompute.VirtualMachinesClientDeleteResponse], err error)
	}
)

var _ Client = &AzureClient{}

// NewClient creates a VMs client from an authorizer.
func NewClient(auth azure.Authorizer, apiCallTimeout time.Duration) (*AzureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create virtualmachines client options")
	}
	factory, err := armcompute.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armcompute client factory")
	}
	return &AzureClient{factory.NewVirtualMachinesClient(), apiCallTimeout}, nil
}

// Get retrieves information about the model view or the instance view of a virtual machine.
func (ac *AzureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "virtualmachines.AzureClient.Get")
	defer done()

	resp, err := ac.virtualmachines.Get(ctx, spec.ResourceGroupName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.VirtualMachine, nil
}

// CreateOrUpdateAsync creates or updates a virtual machine asynchronously.
// It sends a PUT request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *AzureClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armcompute.VirtualMachinesClientCreateOrUpdateResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "virtualmachines.AzureClient.CreateOrUpdate")
	defer done()

	vm, ok := parameters.(armcompute.VirtualMachine)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armcompute.VirtualMachine", parameters)
	}

	opts := &armcompute.VirtualMachinesClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken}
	poller, err = ac.virtualmachines.BeginCreateOrUpdate(ctx, spec.ResourceGroupName(), spec.ResourceName(), vm, opts)
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
	return resp.VirtualMachine, nil, err
}

// DeleteAsync deletes a virtual machine asynchronously. DeleteAsync sends a DELETE
// request to Azure and if accepted without error, the func will return a Poller which can be used to track the ongoing
// progress of the operation.
func (ac *AzureClient) DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[armcompute.VirtualMachinesClientDeleteResponse], err error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "virtualmachines.AzureClient.Delete")
	defer done()

	forceDelete := ptr.To(true)
	opts := &armcompute.VirtualMachinesClientBeginDeleteOptions{ResumeToken: resumeToken, ForceDeletion: forceDelete}
	poller, err = ac.virtualmachines.BeginDelete(ctx, spec.ResourceGroupName(), spec.ResourceName(), opts)
	if err != nil {
		if azure.BadRequest(err) {
			log.Info("Failed to Begin VM Delete with Force Deletion, retrying without the force flag")
			opts.ForceDeletion = ptr.To(false)
			poller, err = ac.virtualmachines.BeginDelete(ctx, spec.ResourceGroupName(), spec.ResourceName(), opts)
		}
		if err != nil {
			return nil, err
		}
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
