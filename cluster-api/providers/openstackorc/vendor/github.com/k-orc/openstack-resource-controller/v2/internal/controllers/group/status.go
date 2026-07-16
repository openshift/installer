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

package group

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

type groupStatusWriter struct{}

type objectApplyT = orcapplyconfigv1alpha1.GroupApplyConfiguration
type statusApplyT = orcapplyconfigv1alpha1.GroupStatusApplyConfiguration

var _ interfaces.ResourceStatusWriter[*orcv1alpha1.Group, *osResourceT, *objectApplyT, *statusApplyT] = groupStatusWriter{}

func (groupStatusWriter) GetApplyConfig(name, namespace string) *objectApplyT {
	return orcapplyconfigv1alpha1.Group(name, namespace)
}

func (groupStatusWriter) ResourceAvailableStatus(orcObject *orcv1alpha1.Group, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}
	return metav1.ConditionTrue, nil
}

func (groupStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osResourceT, statusApply *statusApplyT) {
	resourceStatus := orcapplyconfigv1alpha1.GroupResourceStatus().
		WithDomainID(osResource.DomainID).
		WithName(osResource.Name)

	if osResource.Description != "" {
		resourceStatus.WithDescription(osResource.Description)
	}

	statusApply.WithResource(resourceStatus)
}
