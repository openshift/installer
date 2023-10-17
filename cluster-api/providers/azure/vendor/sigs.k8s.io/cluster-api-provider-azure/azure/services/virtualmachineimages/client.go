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

package virtualmachineimages

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// Client is an interface for listing VM images.
type Client interface {
	List(ctx context.Context, location, publisher, offer, sku string) (armcompute.VirtualMachineImagesClientListResponse, error)
}

// AzureClient contains the Azure go-sdk Client.
type AzureClient struct {
	images *armcompute.VirtualMachineImagesClient
}

var _ Client = (*AzureClient)(nil)

// NewClient creates an AzureClient from an Authorizer.
func NewClient(auth azure.Authorizer) (*AzureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create virtualmachineimages client options")
	}
	computeClientFactory, err := armcompute.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armcompute client factory")
	}
	return &AzureClient{computeClientFactory.NewVirtualMachineImagesClient()}, nil
}

// List returns a VM image list response.
func (ac *AzureClient) List(ctx context.Context, location, publisher, offer, sku string) (armcompute.VirtualMachineImagesClientListResponse, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "virtualmachineimages.AzureClient.List")
	defer done()

	opts := &armcompute.VirtualMachineImagesClientListOptions{}
	return ac.images.List(ctx, location, publisher, offer, sku, opts)
}
