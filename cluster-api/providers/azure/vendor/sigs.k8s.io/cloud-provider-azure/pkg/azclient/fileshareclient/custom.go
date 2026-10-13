/*
Copyright 2023 The Kubernetes Authors.

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

package fileshareclient

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	armstorage "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"

	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/metrics"
)

func (client *Client) Create(ctx context.Context, resourceGroupName string, resourceName string, parentResourceName string, resource armstorage.FileShare, expand *string) (*armstorage.FileShare, error) {
	resp, err := client.FileSharesClient.Create(ctx, resourceGroupName, resourceName, parentResourceName, resource, &armstorage.FileSharesClientCreateOptions{
		Expand: expand,
	})
	if err != nil {
		return nil, err
	}
	return &resp.FileShare, nil
}

func (client *Client) Update(ctx context.Context, resourceGroupName string, resourceName string, parentResourceName string, resource armstorage.FileShare) (*armstorage.FileShare, error) {
	resp, err := client.FileSharesClient.Update(ctx, resourceGroupName, resourceName, parentResourceName, resource, nil)
	if err != nil {
		return nil, err
	}
	return &resp.FileShare, nil
}

// Delete deletes a FileShare by name.
func (client *Client) Delete(ctx context.Context, resourceGroupName string, parentResourceName string, resourceName string, expand *string) error {
	_, err := client.FileSharesClient.Delete(ctx, resourceGroupName, parentResourceName, resourceName, &armstorage.FileSharesClientDeleteOptions{
		Include: expand,
	})
	return err
}

const ListOperationName = "FileSharesClient.List"

// List gets a list of FileShare in the resource group.
func (client *Client) List(ctx context.Context, resourceGroupName string, accountName string, option *armstorage.FileSharesClientListOptions) (result []*armstorage.FileShareItem, err error) {
	metricsCtx := metrics.BeginARMRequest(client.subscriptionID, resourceGroupName, "FileShare", "list")
	defer func() { metricsCtx.Observe(ctx, err) }()
	ctx, endSpan := runtime.StartSpan(ctx, ListOperationName, client.tracer, nil)
	defer endSpan(err)
	pager := client.FileSharesClient.NewListPager(resourceGroupName, accountName, option)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, nextResult.Value...)
	}
	return result, nil
}

const GetOperationName = "FileSharesClient.Get"

// Get gets the FileShare
func (client *Client) Get(ctx context.Context, resourceGroupName string, accountName string, fileshareName string, option *armstorage.FileSharesClientGetOptions) (result *armstorage.FileShare, err error) {
	metricsCtx := metrics.BeginARMRequest(client.subscriptionID, resourceGroupName, "FileShare", "get")
	defer func() { metricsCtx.Observe(ctx, err) }()
	ctx, endSpan := runtime.StartSpan(ctx, GetOperationName, client.tracer, nil)
	defer endSpan(err)
	resp, err := client.FileSharesClient.Get(ctx, resourceGroupName, accountName, fileshareName, option)
	if err != nil {
		return nil, err
	}
	//handle statuscode
	return &resp.FileShare, nil
}
