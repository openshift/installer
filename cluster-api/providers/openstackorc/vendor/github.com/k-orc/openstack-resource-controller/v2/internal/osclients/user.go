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
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type UserClient interface {
	ListUsers(ctx context.Context, listOpts users.ListOptsBuilder) iter.Seq2[*users.User, error]
	CreateUser(ctx context.Context, opts users.CreateOptsBuilder) (*users.User, error)
	DeleteUser(ctx context.Context, resourceID string) error
	GetUser(ctx context.Context, resourceID string) (*users.User, error)
	UpdateUser(ctx context.Context, id string, opts users.UpdateOptsBuilder) (*users.User, error)
}

type userClient struct{ client *gophercloud.ServiceClient }

// NewUserClient returns a new OpenStack client.
func NewUserClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (UserClient, error) {
	client, err := openstack.NewIdentityV3(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user service client: %v", err)
	}

	return &userClient{client}, nil
}

func (c userClient) ListUsers(ctx context.Context, listOpts users.ListOptsBuilder) iter.Seq2[*users.User, error] {
	pager := users.List(c.client, listOpts)
	return func(yield func(*users.User, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(users.ExtractUsers, yield))
	}
}

func (c userClient) CreateUser(ctx context.Context, opts users.CreateOptsBuilder) (*users.User, error) {
	return users.Create(ctx, c.client, opts).Extract()
}

func (c userClient) DeleteUser(ctx context.Context, resourceID string) error {
	return users.Delete(ctx, c.client, resourceID).ExtractErr()
}

func (c userClient) GetUser(ctx context.Context, resourceID string) (*users.User, error) {
	return users.Get(ctx, c.client, resourceID).Extract()
}

func (c userClient) UpdateUser(ctx context.Context, id string, opts users.UpdateOptsBuilder) (*users.User, error) {
	return users.Update(ctx, c.client, id, opts).Extract()
}

type userErrorClient struct{ error }

// NewUserErrorClient returns a UserClient in which every method returns the given error.
func NewUserErrorClient(e error) UserClient {
	return userErrorClient{e}
}

func (e userErrorClient) ListUsers(_ context.Context, _ users.ListOptsBuilder) iter.Seq2[*users.User, error] {
	return func(yield func(*users.User, error) bool) {
		yield(nil, e.error)
	}
}

func (e userErrorClient) CreateUser(_ context.Context, _ users.CreateOptsBuilder) (*users.User, error) {
	return nil, e.error
}

func (e userErrorClient) DeleteUser(_ context.Context, _ string) error {
	return e.error
}

func (e userErrorClient) GetUser(_ context.Context, _ string) (*users.User, error) {
	return nil, e.error
}

func (e userErrorClient) UpdateUser(_ context.Context, _ string, _ users.UpdateOptsBuilder) (*users.User, error) {
	return nil, e.error
}
