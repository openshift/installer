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
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type RoleAssignmentClient interface {
	ListRoleAssignments(ctx context.Context, listOpts roles.ListAssignmentsOpts) iter.Seq2[*roles.RoleAssignment, error]
	AssignRole(ctx context.Context, roleID string, opts roles.AssignOpts) error
	UnassignRole(ctx context.Context, roleID string, opts roles.UnassignOpts) error
}

type roleassignmentClient struct{ client *gophercloud.ServiceClient }

// NewRoleAssignmentClient returns a new OpenStack Identity client for role assignments.
func NewRoleAssignmentClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (RoleAssignmentClient, error) {
	client, err := openstack.NewIdentityV3(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create role assignment service client: %v", err)
	}

	return &roleassignmentClient{client}, nil
}

func (c roleassignmentClient) ListRoleAssignments(ctx context.Context, listOpts roles.ListAssignmentsOpts) iter.Seq2[*roles.RoleAssignment, error] {
	pager := roles.ListAssignments(c.client, listOpts)
	return func(yield func(*roles.RoleAssignment, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(roles.ExtractRoleAssignments, yield))
	}
}

func (c roleassignmentClient) AssignRole(ctx context.Context, roleID string, opts roles.AssignOpts) error {
	return roles.Assign(ctx, c.client, roleID, opts).ExtractErr()
}

func (c roleassignmentClient) UnassignRole(ctx context.Context, roleID string, opts roles.UnassignOpts) error {
	return roles.Unassign(ctx, c.client, roleID, opts).ExtractErr()
}

type roleassignmentErrorClient struct{ error }

// NewRoleAssignmentErrorClient returns a RoleAssignmentClient in which every method returns the given error.
func NewRoleAssignmentErrorClient(e error) RoleAssignmentClient {
	return roleassignmentErrorClient{e}
}

func (e roleassignmentErrorClient) ListRoleAssignments(_ context.Context, _ roles.ListAssignmentsOpts) iter.Seq2[*roles.RoleAssignment, error] {
	return func(yield func(*roles.RoleAssignment, error) bool) {
		yield(nil, e.error)
	}
}

func (e roleassignmentErrorClient) AssignRole(_ context.Context, _ string, _ roles.AssignOpts) error {
	return e.error
}

func (e roleassignmentErrorClient) UnassignRole(_ context.Context, _ string, _ roles.UnassignOpts) error {
	return e.error
}
