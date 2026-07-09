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
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/groups"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type GroupClient interface {
	ListGroups(ctx context.Context, listOpts groups.ListOptsBuilder) iter.Seq2[*groups.Group, error]
	CreateGroup(ctx context.Context, opts groups.CreateOptsBuilder) (*groups.Group, error)
	DeleteGroup(ctx context.Context, resourceID string) error
	GetGroup(ctx context.Context, resourceID string) (*groups.Group, error)
	UpdateGroup(ctx context.Context, id string, opts groups.UpdateOptsBuilder) (*groups.Group, error)
}

type groupClient struct{ client *gophercloud.ServiceClient }

// NewGroupClient returns a new OpenStack client.
func NewGroupClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (GroupClient, error) {
	client, err := openstack.NewIdentityV3(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create group service client: %v", err)
	}

	return &groupClient{client}, nil
}

func (c groupClient) ListGroups(ctx context.Context, listOpts groups.ListOptsBuilder) iter.Seq2[*groups.Group, error] {
	pager := groups.List(c.client, listOpts)
	return func(yield func(*groups.Group, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(groups.ExtractGroups, yield))
	}
}

func (c groupClient) CreateGroup(ctx context.Context, opts groups.CreateOptsBuilder) (*groups.Group, error) {
	return groups.Create(ctx, c.client, opts).Extract()
}

func (c groupClient) DeleteGroup(ctx context.Context, resourceID string) error {
	return groups.Delete(ctx, c.client, resourceID).ExtractErr()
}

func (c groupClient) GetGroup(ctx context.Context, resourceID string) (*groups.Group, error) {
	return groups.Get(ctx, c.client, resourceID).Extract()
}

func (c groupClient) UpdateGroup(ctx context.Context, id string, opts groups.UpdateOptsBuilder) (*groups.Group, error) {
	return groups.Update(ctx, c.client, id, opts).Extract()
}

type groupErrorClient struct{ error }

// NewGroupErrorClient returns a GroupClient in which every method returns the given error.
func NewGroupErrorClient(e error) GroupClient {
	return groupErrorClient{e}
}

func (e groupErrorClient) ListGroups(_ context.Context, _ groups.ListOptsBuilder) iter.Seq2[*groups.Group, error] {
	return func(yield func(*groups.Group, error) bool) {
		yield(nil, e.error)
	}
}

func (e groupErrorClient) CreateGroup(_ context.Context, _ groups.CreateOptsBuilder) (*groups.Group, error) {
	return nil, e.error
}

func (e groupErrorClient) DeleteGroup(_ context.Context, _ string) error {
	return e.error
}

func (e groupErrorClient) GetGroup(_ context.Context, _ string) (*groups.Group, error) {
	return nil, e.error
}

func (e groupErrorClient) UpdateGroup(_ context.Context, _ string, _ groups.UpdateOptsBuilder) (*groups.Group, error) {
	return nil, e.error
}
