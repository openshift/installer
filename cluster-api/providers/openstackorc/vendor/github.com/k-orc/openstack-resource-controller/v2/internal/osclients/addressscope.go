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
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/addressscopes"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type AddressScopeClient interface {
	ListAddressScopes(ctx context.Context, listOpts addressscopes.ListOptsBuilder) iter.Seq2[*addressscopes.AddressScope, error]
	CreateAddressScope(ctx context.Context, opts addressscopes.CreateOptsBuilder) (*addressscopes.AddressScope, error)
	DeleteAddressScope(ctx context.Context, resourceID string) error
	GetAddressScope(ctx context.Context, resourceID string) (*addressscopes.AddressScope, error)
	UpdateAddressScope(ctx context.Context, id string, opts addressscopes.UpdateOptsBuilder) (*addressscopes.AddressScope, error)
}

type addressscopeClient struct{ client *gophercloud.ServiceClient }

// NewAddressScopeClient returns a new OpenStack client.
func NewAddressScopeClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (AddressScopeClient, error) {
	client, err := openstack.NewNetworkV2(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create addressscope service client: %v", err)
	}

	return &addressscopeClient{client}, nil
}

func (c addressscopeClient) ListAddressScopes(ctx context.Context, listOpts addressscopes.ListOptsBuilder) iter.Seq2[*addressscopes.AddressScope, error] {
	pager := addressscopes.List(c.client, listOpts)
	return func(yield func(*addressscopes.AddressScope, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(addressscopes.ExtractAddressScopes, yield))
	}
}

func (c addressscopeClient) CreateAddressScope(ctx context.Context, opts addressscopes.CreateOptsBuilder) (*addressscopes.AddressScope, error) {
	return addressscopes.Create(ctx, c.client, opts).Extract()
}

func (c addressscopeClient) DeleteAddressScope(ctx context.Context, resourceID string) error {
	return addressscopes.Delete(ctx, c.client, resourceID).ExtractErr()
}

func (c addressscopeClient) GetAddressScope(ctx context.Context, resourceID string) (*addressscopes.AddressScope, error) {
	return addressscopes.Get(ctx, c.client, resourceID).Extract()
}

func (c addressscopeClient) UpdateAddressScope(ctx context.Context, id string, opts addressscopes.UpdateOptsBuilder) (*addressscopes.AddressScope, error) {
	return addressscopes.Update(ctx, c.client, id, opts).Extract()
}

type addressscopeErrorClient struct{ error }

// NewAddressScopeErrorClient returns a AddressScopeClient in which every method returns the given error.
func NewAddressScopeErrorClient(e error) AddressScopeClient {
	return addressscopeErrorClient{e}
}

func (e addressscopeErrorClient) ListAddressScopes(_ context.Context, _ addressscopes.ListOptsBuilder) iter.Seq2[*addressscopes.AddressScope, error] {
	return func(yield func(*addressscopes.AddressScope, error) bool) {
		yield(nil, e.error)
	}
}

func (e addressscopeErrorClient) CreateAddressScope(_ context.Context, _ addressscopes.CreateOptsBuilder) (*addressscopes.AddressScope, error) {
	return nil, e.error
}

func (e addressscopeErrorClient) DeleteAddressScope(_ context.Context, _ string) error {
	return e.error
}

func (e addressscopeErrorClient) GetAddressScope(_ context.Context, _ string) (*addressscopes.AddressScope, error) {
	return nil, e.error
}

func (e addressscopeErrorClient) UpdateAddressScope(_ context.Context, _ string, _ addressscopes.UpdateOptsBuilder) (*addressscopes.AddressScope, error) {
	return nil, e.error
}
