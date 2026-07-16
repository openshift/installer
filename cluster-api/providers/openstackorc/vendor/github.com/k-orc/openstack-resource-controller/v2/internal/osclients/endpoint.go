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
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/endpoints"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type EndpointClient interface {
	ListEndpoints(ctx context.Context, listOpts endpoints.ListOptsBuilder) iter.Seq2[*endpoints.Endpoint, error]
	CreateEndpoint(ctx context.Context, opts endpoints.CreateOptsBuilder) (*endpoints.Endpoint, error)
	DeleteEndpoint(ctx context.Context, resourceID string) error
	GetEndpoint(ctx context.Context, resourceID string) (*endpoints.Endpoint, error)
	UpdateEndpoint(ctx context.Context, id string, opts endpoints.UpdateOptsBuilder) (*endpoints.Endpoint, error)
}

type endpointClient struct{ client *gophercloud.ServiceClient }

// NewEndpointClient returns a new OpenStack client.
func NewEndpointClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (EndpointClient, error) {
	client, err := openstack.NewIdentityV3(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint service client: %v", err)
	}

	return &endpointClient{client}, nil
}

func (c endpointClient) ListEndpoints(ctx context.Context, listOpts endpoints.ListOptsBuilder) iter.Seq2[*endpoints.Endpoint, error] {
	pager := endpoints.List(c.client, listOpts)
	return func(yield func(*endpoints.Endpoint, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(endpoints.ExtractEndpoints, yield))
	}
}

func (c endpointClient) CreateEndpoint(ctx context.Context, opts endpoints.CreateOptsBuilder) (*endpoints.Endpoint, error) {
	return endpoints.Create(ctx, c.client, opts).Extract()
}

func (c endpointClient) DeleteEndpoint(ctx context.Context, resourceID string) error {
	return endpoints.Delete(ctx, c.client, resourceID).ExtractErr()
}

func (c endpointClient) GetEndpoint(ctx context.Context, resourceID string) (*endpoints.Endpoint, error) {
	return endpoints.Get(ctx, c.client, resourceID).Extract()
}

func (c endpointClient) UpdateEndpoint(ctx context.Context, id string, opts endpoints.UpdateOptsBuilder) (*endpoints.Endpoint, error) {
	return endpoints.Update(ctx, c.client, id, opts).Extract()
}

type endpointErrorClient struct{ error }

// NewEndpointErrorClient returns a EndpointClient in which every method returns the given error.
func NewEndpointErrorClient(e error) EndpointClient {
	return endpointErrorClient{e}
}

func (e endpointErrorClient) ListEndpoints(_ context.Context, _ endpoints.ListOptsBuilder) iter.Seq2[*endpoints.Endpoint, error] {
	return func(yield func(*endpoints.Endpoint, error) bool) {
		yield(nil, e.error)
	}
}

func (e endpointErrorClient) CreateEndpoint(_ context.Context, _ endpoints.CreateOptsBuilder) (*endpoints.Endpoint, error) {
	return nil, e.error
}

func (e endpointErrorClient) DeleteEndpoint(_ context.Context, _ string) error {
	return e.error
}

func (e endpointErrorClient) GetEndpoint(_ context.Context, _ string) (*endpoints.Endpoint, error) {
	return nil, e.error
}

func (e endpointErrorClient) UpdateEndpoint(_ context.Context, _ string, _ endpoints.UpdateOptsBuilder) (*endpoints.Endpoint, error) {
	return nil, e.error
}
