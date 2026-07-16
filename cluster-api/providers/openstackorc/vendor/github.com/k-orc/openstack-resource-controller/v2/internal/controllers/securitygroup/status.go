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

package securitygroup

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

type objectApplyPT = *orcapplyconfigv1alpha1.SecurityGroupApplyConfiguration
type statusApplyPT = *orcapplyconfigv1alpha1.SecurityGroupStatusApplyConfiguration

type securityGroupStatusWriter struct{}

var _ interfaces.ResourceStatusWriter[orcObjectPT, *osResourceT, objectApplyPT, statusApplyPT] = securityGroupStatusWriter{}

func (securityGroupStatusWriter) GetApplyConfig(name, namespace string) objectApplyPT {
	return orcapplyconfigv1alpha1.SecurityGroup(name, namespace)
}

func (securityGroupStatusWriter) ResourceAvailableStatus(orcObject orcObjectPT, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}

	resourceSpec := orcObject.Spec.Resource
	if resourceSpec != nil && resourceSpec.Rules != nil {
		// Make sure specified security group rules exist in resource

		resourceStatus := orcObject.Status.Resource
		if resourceStatus == nil || resourceStatus.Rules == nil {
			return metav1.ConditionFalse, progress.WaitingOnOpenStack(progress.WaitingOnReady, securityGroupAvailablePollingPeriod)
		}

		if len(resourceSpec.Rules) != len(resourceStatus.Rules) {
			return metav1.ConditionFalse, progress.WaitingOnOpenStack(progress.WaitingOnReady, securityGroupAvailablePollingPeriod)
		}

		if len(resourceSpec.Rules) != len(osResource.Rules) {
			return metav1.ConditionFalse, progress.WaitingOnOpenStack(progress.WaitingOnReady, securityGroupAvailablePollingPeriod)
		}
	}

	return metav1.ConditionTrue, nil
}

func (securityGroupStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osResourceT, statusApply statusApplyPT) {
	securitygroupResourceStatus := orcapplyconfigv1alpha1.SecurityGroupResourceStatus().
		WithName(osResource.Name).
		WithProjectID(osResource.ProjectID).
		WithTags(osResource.Tags...).
		WithStateful(osResource.Stateful).
		WithCreatedAt(metav1.NewTime(osResource.CreatedAt)).
		WithUpdatedAt(metav1.NewTime(osResource.UpdatedAt))

	if osResource.Description != "" {
		securitygroupResourceStatus.WithDescription(osResource.Description)
	}
	for i := range osResource.Rules {
		rule := &osResource.Rules[i]

		ruleStatus := orcapplyconfigv1alpha1.SecurityGroupRuleStatus().
			WithID(osResource.Rules[i].ID).
			WithDescription(osResource.Rules[i].Description).
			WithDirection(osResource.Rules[i].Direction).
			WithRemoteGroupID(osResource.Rules[i].RemoteGroupID).
			WithRemoteIPPrefix(osResource.Rules[i].RemoteIPPrefix).
			WithProtocol(osResource.Rules[i].Protocol).
			WithEthertype(osResource.Rules[i].EtherType)

		if rule.PortRangeMin != 0 || rule.PortRangeMax != 0 {
			ruleStatus.WithPortRange(orcapplyconfigv1alpha1.PortRangeStatus().
				WithMin(int32(rule.PortRangeMin)).
				WithMax(int32(rule.PortRangeMax)))
		}

		securitygroupResourceStatus.WithRules(ruleStatus)
	}

	statusApply.WithResource(securitygroupResourceStatus)
}
