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
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type RoleClient interface {
	ListRoles(ctx context.Context, listOpts roles.ListOptsBuilder) iter.Seq2[*roles.Role, error]
	CreateRole(ctx context.Context, opts roles.CreateOptsBuilder) (*roles.Role, error)
	DeleteRole(ctx context.Context, resourceID string) error
	GetRole(ctx context.Context, resourceID string) (*roles.Role, error)
	UpdateRole(ctx context.Context, id string, opts roles.UpdateOptsBuilder) (*roles.Role, error)
}

type roleClient struct{ client *gophercloud.ServiceClient }

// NewRoleClient returns a new OpenStack client.
func NewRoleClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (RoleClient, error) {
	client, err := openstack.NewIdentityV3(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create role service client: %v", err)
	}

	return &roleClient{client}, nil
}

func (c roleClient) ListRoles(ctx context.Context, listOpts roles.ListOptsBuilder) iter.Seq2[*roles.Role, error] {
	pager := roles.List(c.client, listOpts)
	return func(yield func(*roles.Role, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(roles.ExtractRoles, yield))
	}
}

func (c roleClient) CreateRole(ctx context.Context, opts roles.CreateOptsBuilder) (*roles.Role, error) {
	return roles.Create(ctx, c.client, opts).Extract()
}

func (c roleClient) DeleteRole(ctx context.Context, resourceID string) error {
	return roles.Delete(ctx, c.client, resourceID).ExtractErr()
}

func (c roleClient) GetRole(ctx context.Context, resourceID string) (*roles.Role, error) {
	return roles.Get(ctx, c.client, resourceID).Extract()
}

func (c roleClient) UpdateRole(ctx context.Context, id string, opts roles.UpdateOptsBuilder) (*roles.Role, error) {
	return roles.Update(ctx, c.client, id, opts).Extract()
}

type roleErrorClient struct{ error }

// NewRoleErrorClient returns a RoleClient in which every method returns the given error.
func NewRoleErrorClient(e error) RoleClient {
	return roleErrorClient{e}
}

func (e roleErrorClient) ListRoles(_ context.Context, _ roles.ListOptsBuilder) iter.Seq2[*roles.Role, error] {
	return func(yield func(*roles.Role, error) bool) {
		yield(nil, e.error)
	}
}

func (e roleErrorClient) CreateRole(_ context.Context, _ roles.CreateOptsBuilder) (*roles.Role, error) {
	return nil, e.error
}

func (e roleErrorClient) DeleteRole(_ context.Context, _ string) error {
	return e.error
}

func (e roleErrorClient) GetRole(_ context.Context, _ string) (*roles.Role, error) {
	return nil, e.error
}

func (e roleErrorClient) UpdateRole(_ context.Context, _ string, _ roles.UpdateOptsBuilder) (*roles.Role, error) {
	return nil, e.error
}
