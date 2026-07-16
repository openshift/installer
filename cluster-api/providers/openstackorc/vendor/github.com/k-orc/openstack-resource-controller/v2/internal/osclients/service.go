/*
Copyright 2025 The ORC Authors.

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
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/services"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type ServiceClient interface {
	ListServices(ctx context.Context, listOpts services.ListOptsBuilder) iter.Seq2[*services.Service, error]
	CreateService(ctx context.Context, opts services.CreateOptsBuilder) (*services.Service, error)
	DeleteService(ctx context.Context, resourceID string) error
	GetService(ctx context.Context, resourceID string) (*services.Service, error)
	UpdateService(ctx context.Context, id string, opts services.UpdateOptsBuilder) (*services.Service, error)
}

type serviceClient struct{ client *gophercloud.ServiceClient }

// NewServiceClient returns a new OpenStack client.
func NewServiceClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (ServiceClient, error) {
	client, err := openstack.NewIdentityV3(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create service service client: %v", err)
	}

	return &serviceClient{client}, nil
}

func (c serviceClient) ListServices(ctx context.Context, listOpts services.ListOptsBuilder) iter.Seq2[*services.Service, error] {
	pager := services.List(c.client, listOpts)
	return func(yield func(*services.Service, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(services.ExtractServices, yield))
	}
}

func (c serviceClient) CreateService(ctx context.Context, opts services.CreateOptsBuilder) (*services.Service, error) {
	return services.Create(ctx, c.client, opts).Extract()
}

func (c serviceClient) DeleteService(ctx context.Context, resourceID string) error {
	return services.Delete(ctx, c.client, resourceID).ExtractErr()
}

func (c serviceClient) GetService(ctx context.Context, resourceID string) (*services.Service, error) {
	return services.Get(ctx, c.client, resourceID).Extract()
}

func (c serviceClient) UpdateService(ctx context.Context, id string, opts services.UpdateOptsBuilder) (*services.Service, error) {
	return services.Update(ctx, c.client, id, opts).Extract()
}

type serviceErrorClient struct{ error }

// NewServiceErrorClient returns a ServiceClient in which every method returns the given error.
func NewServiceErrorClient(e error) ServiceClient {
	return serviceErrorClient{e}
}

func (e serviceErrorClient) ListServices(_ context.Context, _ services.ListOptsBuilder) iter.Seq2[*services.Service, error] {
	return func(yield func(*services.Service, error) bool) {
		yield(nil, e.error)
	}
}

func (e serviceErrorClient) CreateService(_ context.Context, _ services.CreateOptsBuilder) (*services.Service, error) {
	return nil, e.error
}

func (e serviceErrorClient) DeleteService(_ context.Context, _ string) error {
	return e.error
}

func (e serviceErrorClient) GetService(_ context.Context, _ string) (*services.Service, error) {
	return nil, e.error
}

func (e serviceErrorClient) UpdateService(_ context.Context, _ string, _ services.UpdateOptsBuilder) (*services.Service, error) {
	return nil, e.error
}
