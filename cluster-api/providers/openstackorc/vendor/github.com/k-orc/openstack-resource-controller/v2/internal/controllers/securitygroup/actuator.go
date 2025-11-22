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
	"context"
	"errors"
	"fmt"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	osclients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/tags"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/ptr"
	"k8s.io/utils/set"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	osResourceT = groups.SecGroup

	createResourceActuator    = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator    = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	reconcileResourceActuator = interfaces.ReconcileResourceActuator[orcObjectPT, osResourceT]
	resourceReconciler        = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory             = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
	securityGroupIterator     = iter.Seq2[*osResourceT, error]
)

type securityGroupActuator struct {
	osClient  osclients.NetworkClient
	k8sClient client.Client
}

var _ createResourceActuator = securityGroupActuator{}
var _ deleteResourceActuator = securityGroupActuator{}

func (actuator securityGroupActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator securityGroupActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	secGroup, err := actuator.osClient.GetSecGroup(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return secGroup, nil
}

func (actuator securityGroupActuator) ListOSResourcesForAdoption(ctx context.Context, obj *orcv1alpha1.SecurityGroup) (securityGroupIterator, bool) {
	if obj.Spec.Resource == nil {
		return nil, false
	}

	listOpts := groups.ListOpts{Name: getResourceName(obj)}
	return actuator.osClient.ListSecGroup(ctx, listOpts), true
}

func (actuator securityGroupActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	var reconcileStatus progress.ReconcileStatus

	project := &orcv1alpha1.Project{}
	if filter.ProjectRef != nil {
		projectKey := client.ObjectKey{Name: string(*filter.ProjectRef), Namespace: obj.Namespace}
		if err := actuator.k8sClient.Get(ctx, projectKey, project); err != nil {
			if apierrors.IsNotFound(err) {
				reconcileStatus = reconcileStatus.WithReconcileStatus(
					progress.WaitingOnObject("Project", projectKey.Name, progress.WaitingOnCreation))
			} else {
				reconcileStatus = reconcileStatus.WithReconcileStatus(
					progress.WrapError(fmt.Errorf("fetching project %s: %w", projectKey.Name, err)))
			}
		} else {
			if !orcv1alpha1.IsAvailable(project) || project.Status.ID == nil {
				reconcileStatus = reconcileStatus.WithReconcileStatus(
					progress.WaitingOnObject("Project", projectKey.Name, progress.WaitingOnReady))
			}
		}
	}

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	listOpts := groups.ListOpts{
		Name:        string(ptr.Deref(filter.Name, "")),
		Description: string(ptr.Deref(filter.Description, "")),
		ProjectID:   ptr.Deref(project.Status.ID, ""),
		Tags:        tags.Join(filter.Tags),
		TagsAny:     tags.Join(filter.TagsAny),
		NotTags:     tags.Join(filter.NotTags),
		NotTagsAny:  tags.Join(filter.NotTagsAny),
	}
	return actuator.osClient.ListSecGroup(ctx, listOpts), nil
}

func (actuator securityGroupActuator) CreateResource(ctx context.Context, obj *orcv1alpha1.SecurityGroup) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}

	var projectID string
	if resource.ProjectRef != nil {
		project, reconcileStatus := projectDependency.GetDependency(
			ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.Project) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
			return nil, reconcileStatus
		}
		projectID = ptr.Deref(project.Status.ID, "")
	}

	createOpts := groups.CreateOpts{
		Name:        getResourceName(obj),
		Description: string(ptr.Deref(resource.Description, "")),
		Stateful:    resource.Stateful,
		ProjectID:   projectID,
	}

	// FIXME(mandre) The security group inherits the default security group
	// rules. This could be a problem when we implement `update` if ORC
	// does not takes these rules into account.
	osResource, err := actuator.osClient.CreateSecGroup(ctx, &createOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a conflict
		if orcerrors.IsConflict(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator securityGroupActuator) DeleteResource(ctx context.Context, _ *orcv1alpha1.SecurityGroup, osResource *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteSecGroup(ctx, osResource.ID))
}

var _ reconcileResourceActuator = securityGroupActuator{}

func (actuator securityGroupActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		tags.ReconcileTags[orcObjectPT, osResourceT](orcObject.Spec.Resource.Tags, osResource.Tags, tags.NewNeutronTagReplacer(actuator.osClient, "security-groups", osResource.ID)),
		actuator.updateRules,
		actuator.updateResource,
	}, nil
}

func (actuator securityGroupActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := groups.UpdateOpts{RevisionNumber: &osResource.RevisionNumber}

	handleNameUpdate(&updateOpts, obj, osResource)
	handleDescriptionUpdate(&updateOpts, resource, osResource)

	needsUpdate, err := needsUpdate(updateOpts)
	if err != nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err))
	}
	if !needsUpdate {
		log.V(logging.Debug).Info("No changes")
		return nil
	}

	_, err = actuator.osClient.UpdateSecGroup(ctx, osResource.ID, updateOpts)

	if orcerrors.IsConflict(err) {
		err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
	}
	if err != nil {
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts groups.UpdateOptsBuilder) (bool, error) {
	updateOptsMap, err := updateOpts.ToSecGroupUpdateMap()
	if err != nil {
		return false, err
	}

	secGroupUpdateMap, ok := updateOptsMap["security_group"].(map[string]any)
	if !ok {
		secGroupUpdateMap = make(map[string]any)
	}

	// Revision number is not returned in the output of updateOpts.ToSecGroupUpdateMap()
	// so nothing to drop here

	return len(secGroupUpdateMap) > 0, nil
}

func handleNameUpdate(updateOpts *groups.UpdateOpts, obj orcObjectPT, osResource *osResourceT) {
	name := getResourceName(obj)
	if osResource.Name != name {
		updateOpts.Name = name
	}
}

func handleDescriptionUpdate(updateOpts *groups.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	description := string(ptr.Deref(resource.Description, ""))
	if osResource.Description != description {
		updateOpts.Description = &description
	}
}

func rulesMatch(orcRule *orcv1alpha1.SecurityGroupRule, osRule *rules.SecGroupRule) bool {
	// Don't compare description if it's not set in the spec
	if orcRule.Description != nil && string(*orcRule.Description) != osRule.Description {
		return false
	}

	// Don't compare direction if it's not set in the spec.
	// TODO check what we get from neutron in this field if we didn't set it in the spec
	if orcRule.Direction != nil && string(*orcRule.Direction) != osRule.Direction {
		return false
	}

	// Always compare RemoteIPPrefix. If unset in ORC it must be empty in OpenStack
	if string(ptr.Deref(orcRule.RemoteIPPrefix, "")) != osRule.RemoteIPPrefix {
		return false
	}

	// Always compare protocol. Unset == "" from gophercloud
	if string(ptr.Deref(orcRule.Protocol, "")) != osRule.Protocol {
		return false
	}

	if string(orcRule.Ethertype) != osRule.EtherType {
		return false
	}

	if orcRule.PortRange == nil {
		if osRule.PortRangeMin != 0 || osRule.PortRangeMax != 0 {
			return false
		}
	} else {
		if int(orcRule.PortRange.Min) != osRule.PortRangeMin || int(orcRule.PortRange.Max) != osRule.PortRangeMax {
			return false
		}
	}

	return true
}

func (actuator securityGroupActuator) updateRules(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	resource := orcObject.Spec.Resource
	if resource == nil {
		return nil
	}

	var projectID string
	if resource.ProjectRef != nil {
		project, reconcileStatus := projectDependency.GetDependency(
			ctx, actuator.k8sClient, orcObject, func(dep *orcv1alpha1.Project) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
			return reconcileStatus
		}
		projectID = ptr.Deref(project.Status.ID, "")
	}

	matchedRuleIDs := set.New[string]()
	allRuleIDS := set.New[string]()
	var createRules []*orcv1alpha1.SecurityGroupRule

orcRules:
	for i := range resource.Rules {
		orcRule := &resource.Rules[i]
		for j := range osResource.Rules {
			osRule := &osResource.Rules[j]

			if rulesMatch(orcRule, osRule) {
				matchedRuleIDs.Insert(osRule.ID)
				continue orcRules
			}
		}
		createRules = append(createRules, orcRule)
	}

	for i := range osResource.Rules {
		allRuleIDS.Insert(osResource.Rules[i].ID)
	}
	deleteRuleIDs := allRuleIDS.Difference(matchedRuleIDs)

	ruleCreateOpts := make([]rules.CreateOpts, len(createRules))
	for i := range createRules {
		ruleCreateOpts[i] = rules.CreateOpts{
			SecGroupID:     osResource.ID,
			Description:    string(ptr.Deref(createRules[i].Description, "")),
			Direction:      rules.RuleDirection(ptr.Deref(createRules[i].Direction, "")),
			RemoteIPPrefix: string(ptr.Deref(createRules[i].RemoteIPPrefix, "")),
			Protocol:       rules.RuleProtocol(ptr.Deref(createRules[i].Protocol, "")),
			EtherType:      rules.RuleEtherType(createRules[i].Ethertype),
			ProjectID:      projectID,
		}
		if createRules[i].PortRange != nil {
			ruleCreateOpts[i].PortRangeMin = int(resource.Rules[i].PortRange.Min)
			ruleCreateOpts[i].PortRangeMax = int(resource.Rules[i].PortRange.Max)
		}
	}

	var err error
	if len(ruleCreateOpts) > 0 {
		if _, createErr := actuator.osClient.CreateSecGroupRules(ctx, ruleCreateOpts); createErr != nil {
			// We should require the spec to be updated before retrying a create which returned a conflict
			if orcerrors.IsRetryable(createErr) {
				createErr = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+createErr.Error(), createErr)
			} else {
				createErr = fmt.Errorf("creating security group rules: %w", createErr)
			}

			err = errors.Join(err, createErr)
		}
	}

	for _, id := range deleteRuleIDs.UnsortedList() {
		if deleteErr := actuator.osClient.DeleteSecGroupRule(ctx, id); deleteErr != nil {
			err = errors.Join(err, fmt.Errorf("deleting security group rule %s: %w", id, deleteErr))
		}
	}

	if err != nil {
		return progress.WrapError(err)
	}

	// If we added or removed any rules above, schedule another reconcile so we can observe the updated security group
	if len(ruleCreateOpts) > 0 || len(deleteRuleIDs) > 0 {
		return progress.NeedsRefresh()
	}

	return nil
}

type securityGroupHelperFactory struct{}

var _ helperFactory = securityGroupHelperFactory{}

func (securityGroupHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return securitygroupAdapter{obj}
}

func (securityGroupHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (securityGroupHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.SecurityGroup, controller interfaces.ResourceController) (securityGroupActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)
	k8sClient := controller.GetK8sClient()

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, k8sClient, orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return securityGroupActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, k8sClient, log, orcObject)
	if err != nil {
		return securityGroupActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewNetworkClient()
	if err != nil {
		return securityGroupActuator{}, progress.WrapError(err)
	}

	return securityGroupActuator{
		osClient:  osClient,
		k8sClient: k8sClient,
	}, nil
}
