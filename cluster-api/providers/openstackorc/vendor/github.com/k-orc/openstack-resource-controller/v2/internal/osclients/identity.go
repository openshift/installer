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
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type IdentityClient interface {
	CreateProject(ctx context.Context, opts projects.CreateOptsBuilder) (*projects.Project, error)
	GetProject(ctx context.Context, id string) (*projects.Project, error)
	DeleteProject(ctx context.Context, id string) error
	ListProjects(ctx context.Context, opts projects.ListOptsBuilder) iter.Seq2[*projects.Project, error]
	UpdateProject(ctx context.Context, id string, opts projects.UpdateOptsBuilder) (*projects.Project, error)
}

type identityClient struct{ client *gophercloud.ServiceClient }

// NewIdentityClient returns a new cinder client.
func NewIdentityClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (IdentityClient, error) {
	identity, err := openstack.NewIdentityV3(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create identity service client: %v", err)
	}

	return &identityClient{identity}, nil
}

func (c identityClient) ListProjects(ctx context.Context, opts projects.ListOptsBuilder) iter.Seq2[*projects.Project, error] {
	pager := projects.List(c.client, opts)
	return func(yield func(*projects.Project, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(projects.ExtractProjects, yield))
	}
}

func (c identityClient) CreateProject(ctx context.Context, opts projects.CreateOptsBuilder) (*projects.Project, error) {
	return projects.Create(ctx, c.client, opts).Extract()
}

func (c identityClient) DeleteProject(ctx context.Context, id string) error {
	return projects.Delete(ctx, c.client, id).ExtractErr()
}

func (c identityClient) GetProject(ctx context.Context, id string) (*projects.Project, error) {
	return projects.Get(ctx, c.client, id).Extract()
}

func (c identityClient) UpdateProject(ctx context.Context, id string, opts projects.UpdateOptsBuilder) (*projects.Project, error) {
	return projects.Update(ctx, c.client, id, opts).Extract()
}

type identityErrorClient struct{ error }

// NewProjectErrorClient returns a IdentityClient in which every method returns the given error.
func NewIdentityErrorClient(e error) IdentityClient {
	return identityErrorClient{e}
}

func (e identityErrorClient) ListProjects(_ context.Context, _ projects.ListOptsBuilder) iter.Seq2[*projects.Project, error] {
	return func(yield func(*projects.Project, error) bool) {
		yield(nil, e.error)
	}
}

func (e identityErrorClient) CreateProject(_ context.Context, _ projects.CreateOptsBuilder) (*projects.Project, error) {
	return nil, e.error
}

func (e identityErrorClient) DeleteProject(_ context.Context, _ string) error {
	return e.error
}

func (e identityErrorClient) GetProject(_ context.Context, _ string) (*projects.Project, error) {
	return nil, e.error
}

func (e identityErrorClient) UpdateProject(_ context.Context, _ string, _ projects.UpdateOptsBuilder) (*projects.Project, error) {
	return nil, e.error
}
