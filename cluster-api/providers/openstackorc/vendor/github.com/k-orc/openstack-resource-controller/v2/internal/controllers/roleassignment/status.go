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
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

type roleassignmentStatusWriter struct{}

type objectApplyT = orcapplyconfigv1alpha1.RoleAssignmentApplyConfiguration
type statusApplyT = orcapplyconfigv1alpha1.RoleAssignmentStatusApplyConfiguration

var _ interfaces.ResourceStatusWriter[*orcv1alpha1.RoleAssignment, *osResourceT, *objectApplyT, *statusApplyT] = roleassignmentStatusWriter{}

func (roleassignmentStatusWriter) GetApplyConfig(name, namespace string) *objectApplyT {
	return orcapplyconfigv1alpha1.RoleAssignment(name, namespace)
}

// ResourceAvailableStatus returns the availability status of the role assignment.
// Role assignments don't have Status.ID, so availability is based on osResource
// presence and status component fields.
func (roleassignmentStatusWriter) ResourceAvailableStatus(orcObject *orcv1alpha1.RoleAssignment, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource != nil {
		return metav1.ConditionTrue, nil
	}

	// If we previously observed component IDs but can't fetch the resource now,
	// report Unknown since we can't confirm availability.
	if orcObject.Status.Resource != nil &&
		(orcObject.Status.Resource.RoleID != "" ||
			orcObject.Status.Resource.UserID != "" ||
			orcObject.Status.Resource.GroupID != "" ||
			orcObject.Status.Resource.ProjectID != "" ||
			orcObject.Status.Resource.DomainID != "") {
		return metav1.ConditionUnknown, nil
	}

	return metav1.ConditionFalse, nil
}

// ApplyResourceStatus writes the role assignment component IDs to status.
func (roleassignmentStatusWriter) ApplyResourceStatus(_ logr.Logger, osResource *osResourceT, statusApply *statusApplyT) {
	resourceStatus := orcapplyconfigv1alpha1.RoleAssignmentResourceStatus()

	if osResource.Role.ID != "" {
		resourceStatus.WithRoleID(osResource.Role.ID)
	}
	if osResource.User.ID != "" {
		resourceStatus.WithUserID(osResource.User.ID)
	}
	if osResource.Group.ID != "" {
		resourceStatus.WithGroupID(osResource.Group.ID)
	}
	if osResource.Scope.Project.ID != "" {
		resourceStatus.WithProjectID(osResource.Scope.Project.ID)
	}
	if osResource.Scope.Domain.ID != "" {
		resourceStatus.WithDomainID(osResource.Scope.Domain.ID)
	}

	statusApply.WithResource(resourceStatus)
}
