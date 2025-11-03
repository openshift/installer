/*
Copyright 2024 The Kubernetes Authors.

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

package subnet

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"k8s.io/klog/v2"

	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/subnetclient"
)

type Repository interface {
	CreateOrUpdate(ctx context.Context, rg string, vnetName string, subnetName string, subnet armnetwork.Subnet) error
	Get(ctx context.Context, rg string, vnetName string, subnetName string) (*armnetwork.Subnet, error)
}

type repo struct {
	SubnetsClient subnetclient.Interface
}

func NewRepo(subnetsClient subnetclient.Interface) (Repository, error) {
	return &repo{
		SubnetsClient: subnetsClient,
	}, nil
}

// CreateOrUpdateSubnet invokes az.SubnetClient.CreateOrUpdate with exponential backoff retry
func (az *repo) CreateOrUpdate(ctx context.Context, rg string, vnetName string, subnetName string, subnet armnetwork.Subnet) error {
	_, rerr := az.SubnetsClient.CreateOrUpdate(ctx, rg, vnetName, subnetName, subnet)
	klog.V(10).Infof("SubnetsClient.CreateOrUpdate(%s): end", subnetName)
	if rerr != nil {
		klog.Errorf("SubnetClient.CreateOrUpdate(%s) failed: %s", subnetName, rerr.Error())
		return rerr
	}

	return nil
}

func (az *repo) Get(ctx context.Context, rg string, vnetName string, subnetName string) (*armnetwork.Subnet, error) {
	subnet, err := az.SubnetsClient.Get(ctx, rg, vnetName, subnetName, nil)
	if err != nil {
		klog.Errorf("SubnetClient.Get(%s) failed: %s", subnetName, err.Error())
		return nil, err
	}
	return subnet, nil
}
