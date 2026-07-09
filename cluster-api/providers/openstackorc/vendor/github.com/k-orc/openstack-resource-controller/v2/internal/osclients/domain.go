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
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type DomainClient interface {
	ListDomains(ctx context.Context, listOpts domains.ListOptsBuilder) iter.Seq2[*domains.Domain, error]
	CreateDomain(ctx context.Context, opts domains.CreateOptsBuilder) (*domains.Domain, error)
	DeleteDomain(ctx context.Context, resourceID string) error
	GetDomain(ctx context.Context, resourceID string) (*domains.Domain, error)
	UpdateDomain(ctx context.Context, id string, opts domains.UpdateOptsBuilder) (*domains.Domain, error)
}

type domainClient struct{ client *gophercloud.ServiceClient }

// NewDomainClient returns a new OpenStack client.
func NewDomainClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (DomainClient, error) {
	client, err := openstack.NewIdentityV3(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create domain service client: %v", err)
	}

	return &domainClient{client}, nil
}

func (c domainClient) ListDomains(ctx context.Context, listOpts domains.ListOptsBuilder) iter.Seq2[*domains.Domain, error] {
	pager := domains.List(c.client, listOpts)
	return func(yield func(*domains.Domain, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(domains.ExtractDomains, yield))
	}
}

func (c domainClient) CreateDomain(ctx context.Context, opts domains.CreateOptsBuilder) (*domains.Domain, error) {
	return domains.Create(ctx, c.client, opts).Extract()
}

func (c domainClient) DeleteDomain(ctx context.Context, resourceID string) error {
	return domains.Delete(ctx, c.client, resourceID).ExtractErr()
}

func (c domainClient) GetDomain(ctx context.Context, resourceID string) (*domains.Domain, error) {
	return domains.Get(ctx, c.client, resourceID).Extract()
}

func (c domainClient) UpdateDomain(ctx context.Context, id string, opts domains.UpdateOptsBuilder) (*domains.Domain, error) {
	return domains.Update(ctx, c.client, id, opts).Extract()
}

type domainErrorClient struct{ error }

// NewDomainErrorClient returns a DomainClient in which every method returns the given error.
func NewDomainErrorClient(e error) DomainClient {
	return domainErrorClient{e}
}

func (e domainErrorClient) ListDomains(_ context.Context, _ domains.ListOptsBuilder) iter.Seq2[*domains.Domain, error] {
	return func(yield func(*domains.Domain, error) bool) {
		yield(nil, e.error)
	}
}

func (e domainErrorClient) CreateDomain(_ context.Context, _ domains.CreateOptsBuilder) (*domains.Domain, error) {
	return nil, e.error
}

func (e domainErrorClient) DeleteDomain(_ context.Context, _ string) error {
	return e.error
}

func (e domainErrorClient) GetDomain(_ context.Context, _ string) (*domains.Domain, error) {
	return nil, e.error
}

func (e domainErrorClient) UpdateDomain(_ context.Context, _ string, _ domains.UpdateOptsBuilder) (*domains.Domain, error) {
	return nil, e.error
}
