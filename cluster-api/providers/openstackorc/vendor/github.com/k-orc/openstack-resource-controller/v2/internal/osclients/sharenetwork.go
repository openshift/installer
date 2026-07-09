/*
Copyright The ORC Authors.

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

package osclients

import (
	"context"
	"fmt"
	"iter"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharenetworks"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type ShareNetworkClient interface {
	ListShareNetworks(ctx context.Context, listOpts sharenetworks.ListOptsBuilder) iter.Seq2[*sharenetworks.ShareNetwork, error]
	CreateShareNetwork(ctx context.Context, opts sharenetworks.CreateOptsBuilder) (*sharenetworks.ShareNetwork, error)
	DeleteShareNetwork(ctx context.Context, resourceID string) error
	GetShareNetwork(ctx context.Context, resourceID string) (*sharenetworks.ShareNetwork, error)
	UpdateShareNetwork(ctx context.Context, id string, opts sharenetworks.UpdateOptsBuilder) (*sharenetworks.ShareNetwork, error)
}

type sharenetworkClient struct{ client *gophercloud.ServiceClient }

// NewShareNetworkClient returns a new OpenStack client.
func NewShareNetworkClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (ShareNetworkClient, error) {
	client, err := openstack.NewSharedFileSystemV2(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create sharenetwork service client: %v", err)
	}

	return &sharenetworkClient{client}, nil
}

func (c sharenetworkClient) ListShareNetworks(ctx context.Context, listOpts sharenetworks.ListOptsBuilder) iter.Seq2[*sharenetworks.ShareNetwork, error] {
	pager := sharenetworks.ListDetail(c.client, listOpts)
	return func(yield func(*sharenetworks.ShareNetwork, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(sharenetworks.ExtractShareNetworks, yield))
	}
}

func (c sharenetworkClient) CreateShareNetwork(ctx context.Context, opts sharenetworks.CreateOptsBuilder) (*sharenetworks.ShareNetwork, error) {
	return sharenetworks.Create(ctx, c.client, opts).Extract()
}

func (c sharenetworkClient) DeleteShareNetwork(ctx context.Context, resourceID string) error {
	return sharenetworks.Delete(ctx, c.client, resourceID).ExtractErr()
}

func (c sharenetworkClient) GetShareNetwork(ctx context.Context, resourceID string) (*sharenetworks.ShareNetwork, error) {
	return sharenetworks.Get(ctx, c.client, resourceID).Extract()
}

func (c sharenetworkClient) UpdateShareNetwork(ctx context.Context, id string, opts sharenetworks.UpdateOptsBuilder) (*sharenetworks.ShareNetwork, error) {
	return sharenetworks.Update(ctx, c.client, id, opts).Extract()
}

type sharenetworkErrorClient struct{ error }

// NewShareNetworkErrorClient returns a ShareNetworkClient in which every method returns the given error.
func NewShareNetworkErrorClient(e error) ShareNetworkClient {
	return sharenetworkErrorClient{e}
}

func (e sharenetworkErrorClient) ListShareNetworks(_ context.Context, _ sharenetworks.ListOptsBuilder) iter.Seq2[*sharenetworks.ShareNetwork, error] {
	return func(yield func(*sharenetworks.ShareNetwork, error) bool) {
		yield(nil, e.error)
	}
}

func (e sharenetworkErrorClient) CreateShareNetwork(_ context.Context, _ sharenetworks.CreateOptsBuilder) (*sharenetworks.ShareNetwork, error) {
	return nil, e.error
}

func (e sharenetworkErrorClient) DeleteShareNetwork(_ context.Context, _ string) error {
	return e.error
}

func (e sharenetworkErrorClient) GetShareNetwork(_ context.Context, _ string) (*sharenetworks.ShareNetwork, error) {
	return nil, e.error
}

func (e sharenetworkErrorClient) UpdateShareNetwork(_ context.Context, _ string, _ sharenetworks.UpdateOptsBuilder) (*sharenetworks.ShareNetwork, error) {
	return nil, e.error
}
