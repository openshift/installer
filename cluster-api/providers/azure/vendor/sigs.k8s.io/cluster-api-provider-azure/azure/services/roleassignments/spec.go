/*
Copyright 2019 The Kubernetes Authors.

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

package roleassignments

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
)

// RoleAssignmentSpec defines the specification for a role assignment.
type RoleAssignmentSpec struct {
	Name             string
	MachineName      string
	ResourceGroup    string
	ResourceType     string
	PrincipalID      *string
	RoleDefinitionID string
	Scope            string
}

// ResourceName returns the name of the role assignment.
func (s *RoleAssignmentSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *RoleAssignmentSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName returns the scope for role assignment.
// TODO: Consider renaming the function for better readability (@sonasingh46).
func (s *RoleAssignmentSpec) OwnerResourceName() string {
	return s.Scope
}

// Parameters returns the parameters for the RoleAssignmentSpec.
func (s *RoleAssignmentSpec) Parameters(ctx context.Context, existing interface{}) (interface{}, error) {
	if existing != nil {
		if _, ok := existing.(armauthorization.RoleAssignment); !ok {
			return nil, errors.Errorf("%T is not an armauthorization.RoleAssignment", existing)
		}
		// RoleAssignmentSpec already exists
		return nil, nil
	}
	return armauthorization.RoleAssignmentCreateParameters{
		Properties: &armauthorization.RoleAssignmentProperties{
			PrincipalID:      s.PrincipalID,
			RoleDefinitionID: ptr.To(s.RoleDefinitionID),
		},
	}, nil
}
