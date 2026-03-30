/*
Copyright 2024 The ORC Authors.

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

package servergroup

import (
	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

type servergroupStatusWriter struct{}

type objectApplyT = orcapplyconfigv1alpha1.ServerGroupApplyConfiguration
type statusApplyT = orcapplyconfigv1alpha1.ServerGroupStatusApplyConfiguration

var _ interfaces.ResourceStatusWriter[*orcv1alpha1.ServerGroup, *servergroups.ServerGroup, *objectApplyT, *statusApplyT] = servergroupStatusWriter{}

func (servergroupStatusWriter) GetApplyConfig(name, namespace string) *objectApplyT {
	return orcapplyconfigv1alpha1.ServerGroup(name, namespace)
}

func (servergroupStatusWriter) ResourceAvailableStatus(orcObject *orcv1alpha1.ServerGroup, osResource *servergroups.ServerGroup) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}

	// ServerGroup is available as soon as it exists
	return metav1.ConditionTrue, nil
}

func (servergroupStatusWriter) ApplyResourceStatus(_ logr.Logger, osResource *servergroups.ServerGroup, statusApply *statusApplyT) {
	status := orcapplyconfigv1alpha1.ServerGroupResourceStatus().
		WithName(osResource.Name).
		WithProjectID(osResource.ProjectID).
		WithUserID(osResource.UserID).
		WithPolicy(ptr.Deref(osResource.Policy, ""))

	if osResource.Rules != nil {
		rules := orcapplyconfigv1alpha1.ServerGroupRulesStatusApplyConfiguration{
			MaxServerPerHost: ptr.To(int32(osResource.Rules.MaxServerPerHost)),
		}
		status.WithRules(&rules)
	}

	statusApply.WithResource(status)
}
