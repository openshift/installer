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

package applicationcredential

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

type applicationcredentialStatusWriter struct{}

type objectApplyT = orcapplyconfigv1alpha1.ApplicationCredentialApplyConfiguration
type statusApplyT = orcapplyconfigv1alpha1.ApplicationCredentialStatusApplyConfiguration

var _ interfaces.ResourceStatusWriter[*orcv1alpha1.ApplicationCredential, *osResourceT, *objectApplyT, *statusApplyT] = applicationcredentialStatusWriter{}

func (applicationcredentialStatusWriter) GetApplyConfig(name, namespace string) *objectApplyT {
	return orcapplyconfigv1alpha1.ApplicationCredential(name, namespace)
}

func (applicationcredentialStatusWriter) ResourceAvailableStatus(orcObject *orcv1alpha1.ApplicationCredential, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}
	return metav1.ConditionTrue, nil
}

func (applicationcredentialStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osResourceT, statusApply *statusApplyT) {
	resourceStatus := orcapplyconfigv1alpha1.ApplicationCredentialResourceStatus().
		WithName(osResource.Name).
		WithUnrestricted(osResource.Unrestricted).
		WithProjectID(osResource.ProjectID)

	if !osResource.ExpiresAt.IsZero() {
		resourceStatus.WithExpiresAt(metav1.NewTime(osResource.ExpiresAt))
	}

	if osResource.Description != "" {
		resourceStatus.WithDescription(osResource.Description)
	}

	for i := range osResource.Roles {
		roleStatus := orcapplyconfigv1alpha1.ApplicationCredentialRoleStatus().
			WithID(osResource.Roles[i].ID).
			WithName(osResource.Roles[i].Name)

		if osResource.Roles[i].DomainID != "" {
			roleStatus.WithDomainID(osResource.Roles[i].DomainID)
		}

		resourceStatus.WithRoles(roleStatus)
	}

	for i := range osResource.AccessRules {
		accessRuleStatus := orcapplyconfigv1alpha1.ApplicationCredentialAccessRuleStatus().
			WithID(osResource.AccessRules[i].ID).
			WithPath(osResource.AccessRules[i].Path).
			WithMethod(osResource.AccessRules[i].Method).
			WithService(osResource.AccessRules[i].Service)

		resourceStatus.WithAccessRules(accessRuleStatus)
	}

	statusApply.WithResource(resourceStatus)
}
