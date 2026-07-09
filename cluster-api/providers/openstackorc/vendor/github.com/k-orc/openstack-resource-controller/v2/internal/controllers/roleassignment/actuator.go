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

package roleassignment

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
)

// OpenStack resource type
type osResourceT = roles.RoleAssignment

type roleassignmentActuator struct {
	osClient  osclients.RoleAssignmentClient
	k8sClient client.Client
}

// buildListOpts constructs a ListAssignmentsOpts from component IDs.
// Only non-empty fields are set, so this works for both exact queries
// (all fields populated) and partial filter queries.
func buildListOpts(roleID, userID, groupID, projectID, domainID string) roles.ListAssignmentsOpts {
	// Note: Don't set Effective parameter - it can cause issues with group assignments
	listOpts := roles.ListAssignmentsOpts{}

	if roleID != "" {
		listOpts.RoleID = roleID
	}
	if userID != "" {
		listOpts.UserID = userID
	}
	if groupID != "" {
		listOpts.GroupID = groupID
	}
	if projectID != "" {
		listOpts.ScopeProjectID = projectID
	}
	if domainID != "" {
		listOpts.ScopeDomainID = domainID
	}

	return listOpts
}

// GetResourceByComponents queries for the role assignment by its tuple (role, actor, scope).
// OpenStack doesn't assign IDs to role assignments - they're identified by this tuple.
// Exactly one of userID/groupID must be set, and exactly one of projectID/domainID must be set.
func (actuator roleassignmentActuator) GetResourceByComponents(
	ctx context.Context,
	roleID string,
	userID string,
	groupID string,
	projectID string,
	domainID string,
) (*osResourceT, progress.ReconcileStatus) {
	listOpts := buildListOpts(roleID, userID, groupID, projectID, domainID)

	// Query with exact filters - should return exactly one result
	osResource, err := atMostOne(actuator.osClient.ListRoleAssignments(ctx, listOpts),
		orcerrors.Terminal(orcv1alpha1.ConditionReasonUnrecoverableError,
			"found more than one matching role assignment for the same (role, actor, scope) tuple"))
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return osResource, nil
}

func (actuator roleassignmentActuator) ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*osResourceT, error], bool) {
	resourceSpec := orcObject.Spec.Resource
	if resourceSpec == nil {
		return nil, false
	}

	// Fetch all dependencies to build the exact filter
	var roleID, userID, groupID, projectID, domainID string

	// Role dependency (required)
	role, rs := dependency.FetchDependency(
		ctx, actuator.k8sClient, orcObject.Namespace, &resourceSpec.RoleRef, "Role",
		func(dep *orcv1alpha1.Role) bool {
			return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
		},
	)
	if needsReschedule, _ := rs.NeedsReschedule(); needsReschedule {
		return nil, false // Not ready
	}
	roleID = ptr.Deref(role.Status.ID, "")

	// Actor dependency (user XOR group)
	if resourceSpec.UserRef != nil {
		user, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, orcObject.Namespace, resourceSpec.UserRef, "User",
			func(dep *orcv1alpha1.User) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := rs.NeedsReschedule(); needsReschedule {
			return nil, false // Not ready
		}
		userID = ptr.Deref(user.Status.ID, "")
	} else {
		group, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, orcObject.Namespace, resourceSpec.GroupRef, "Group",
			func(dep *orcv1alpha1.Group) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := rs.NeedsReschedule(); needsReschedule {
			return nil, false // Not ready
		}
		groupID = ptr.Deref(group.Status.ID, "")
	}

	// Scope dependency (project XOR domain)
	if resourceSpec.ProjectRef != nil {
		project, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, orcObject.Namespace, resourceSpec.ProjectRef, "Project",
			func(dep *orcv1alpha1.Project) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := rs.NeedsReschedule(); needsReschedule {
			return nil, false // Not ready
		}
		projectID = ptr.Deref(project.Status.ID, "")
	} else {
		domain, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, orcObject.Namespace, resourceSpec.DomainRef, "Domain",
			func(dep *orcv1alpha1.Domain) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := rs.NeedsReschedule(); needsReschedule {
			return nil, false // Not ready
		}
		domainID = ptr.Deref(domain.Status.ID, "")
	}

	return actuator.osClient.ListRoleAssignments(ctx, buildListOpts(roleID, userID, groupID, projectID, domainID)), true
}

func (actuator roleassignmentActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	var reconcileStatus progress.ReconcileStatus

	// Build ListAssignmentsOpts from filter references
	var roleID, userID, groupID, projectID, domainID string

	if filter.RoleRef != nil {
		role, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, obj.Namespace, filter.RoleRef, "Role",
			func(dep *orcv1alpha1.Role) bool { return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil },
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(rs)
		if role != nil && role.Status.ID != nil {
			roleID = *role.Status.ID
		}
	}

	if filter.UserRef != nil {
		user, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, obj.Namespace, filter.UserRef, "User",
			func(dep *orcv1alpha1.User) bool { return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil },
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(rs)
		if user != nil && user.Status.ID != nil {
			userID = *user.Status.ID
		}
	}

	if filter.GroupRef != nil {
		group, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, obj.Namespace, filter.GroupRef, "Group",
			func(dep *orcv1alpha1.Group) bool { return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil },
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(rs)
		if group != nil && group.Status.ID != nil {
			groupID = *group.Status.ID
		}
	}

	if filter.ProjectRef != nil {
		project, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, obj.Namespace, filter.ProjectRef, "Project",
			func(dep *orcv1alpha1.Project) bool { return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil },
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(rs)
		if project != nil && project.Status.ID != nil {
			projectID = *project.Status.ID
		}
	}

	if filter.DomainRef != nil {
		domain, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, obj.Namespace, filter.DomainRef, "Domain",
			func(dep *orcv1alpha1.Domain) bool { return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil },
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(rs)
		if domain != nil && domain.Status.ID != nil {
			domainID = *domain.Status.ID
		}
	}

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	return actuator.osClient.ListRoleAssignments(ctx, buildListOpts(roleID, userID, groupID, projectID, domainID)), nil
}

func (actuator roleassignmentActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}
	var reconcileStatus progress.ReconcileStatus

	// Fetch role dependency (required)
	role, roleDepRS := roleDependency.GetDependency(
		ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.Role) bool {
			return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
		},
	)
	reconcileStatus = reconcileStatus.WithReconcileStatus(roleDepRS)
	var roleID string
	if role != nil {
		roleID = ptr.Deref(role.Status.ID, "")
	}

	// Fetch actor dependency (user XOR group)
	var userID, groupID string
	if resource.UserRef != nil {
		user, userDepRS := userDependency.GetDependency(
			ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.User) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(userDepRS)
		if user != nil {
			userID = ptr.Deref(user.Status.ID, "")
		}
	} else {
		group, groupDepRS := groupDependency.GetDependency(
			ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.Group) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(groupDepRS)
		if group != nil {
			groupID = ptr.Deref(group.Status.ID, "")
		}
	}

	// Fetch scope dependency (project XOR domain)
	var projectID, domainID string
	if resource.ProjectRef != nil {
		project, projectDepRS := projectDependency.GetDependency(
			ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.Project) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(projectDepRS)
		if project != nil {
			projectID = ptr.Deref(project.Status.ID, "")
		}
	} else {
		domain, domainDepRS := domainDependency.GetDependency(
			ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.Domain) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(domainDepRS)
		if domain != nil {
			domainID = ptr.Deref(domain.Status.ID, "")
		}
	}

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	// Build AssignOpts
	assignOpts := roles.AssignOpts{
		UserID:    userID,
		GroupID:   groupID,
		ProjectID: projectID,
		DomainID:  domainID,
	}

	// Assign the role (idempotent - returns 204 even if already exists)
	err := actuator.osClient.AssignRole(ctx, roleID, assignOpts)
	if err != nil {
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating role assignment: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	// Verify the assignment was created by listing with exact filters
	osResource, verifyErr := atMostOne(actuator.osClient.ListRoleAssignments(ctx, buildListOpts(roleID, userID, groupID, projectID, domainID)),
		orcerrors.Terminal(orcv1alpha1.ConditionReasonUnrecoverableError,
			"found more than one matching role assignment after creation"))
	if verifyErr != nil {
		return nil, progress.WrapError(verifyErr)
	}
	if osResource == nil {
		// This shouldn't happen - we just assigned it
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonUnrecoverableError,
				"role assignment succeeded but could not be found in OpenStack"))
	}
	return osResource, nil
}

func (actuator roleassignmentActuator) DeleteResource(ctx context.Context, _ orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	// Build UnassignOpts from the osResource
	unassignOpts := roles.UnassignOpts{
		UserID:    osResource.User.ID,
		GroupID:   osResource.Group.ID,
		ProjectID: osResource.Scope.Project.ID,
		DomainID:  osResource.Scope.Domain.ID,
	}

	return progress.WrapError(actuator.osClient.UnassignRole(ctx, osResource.Role.ID, unassignOpts))
}
