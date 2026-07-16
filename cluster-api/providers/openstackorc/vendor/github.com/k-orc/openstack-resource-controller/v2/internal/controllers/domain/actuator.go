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

package domain

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	"github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
)

// OpenStack resource types
type (
	osResourceT = domains.Domain

	createResourceActuator = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	resourceReconciler     = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory          = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

type domainActuator struct {
	osClient  osclients.DomainClient
	k8sClient client.Client
}

var _ createResourceActuator = domainActuator{}
var _ deleteResourceActuator = domainActuator{}

func (domainActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator domainActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	resource, err := actuator.osClient.GetDomain(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return resource, nil
}

func (actuator domainActuator) ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*osResourceT, error], bool) {
	resourceSpec := orcObject.Spec.Resource
	if resourceSpec == nil {
		return nil, false
	}

	listOpts := domains.ListOpts{
		Name: getResourceName(orcObject),
	}

	return actuator.osClient.ListDomains(ctx, listOpts), true
}

func (actuator domainActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {

	listOpts := domains.ListOpts{
		Name:    string(ptr.Deref(filter.Name, "")),
		Enabled: filter.Enabled,
	}

	return actuator.osClient.ListDomains(ctx, listOpts), nil
}

func (actuator domainActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}
	createOpts := domains.CreateOpts{
		Name:        getResourceName(obj),
		Description: ptr.Deref(resource.Description, ""),
		Enabled:     resource.Enabled,
	}

	osResource, err := actuator.osClient.CreateDomain(ctx, createOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a conflict
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator domainActuator) DeleteResource(ctx context.Context, _ orcObjectPT, resource *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteDomain(ctx, resource.ID))
}

func (actuator domainActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := domains.UpdateOpts{}

	handleNameUpdate(&updateOpts, obj, osResource)
	handleDescriptionUpdate(&updateOpts, resource, osResource)
	handleEnabledUpdate(&updateOpts, resource, osResource)

	needsUpdate, err := needsUpdate(updateOpts)
	if err != nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err))
	}
	if !needsUpdate {
		log.V(logging.Debug).Info("No changes")
		return nil
	}

	_, err = actuator.osClient.UpdateDomain(ctx, osResource.ID, updateOpts)

	if err != nil {
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
		}
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts domains.UpdateOpts) (bool, error) {
	updateOptsMap, err := updateOpts.ToDomainUpdateMap()
	if err != nil {
		return false, err
	}

	updateMap, ok := updateOptsMap["domain"].(map[string]any)
	if !ok {
		updateMap = make(map[string]any)
	}

	return len(updateMap) > 0, nil
}

func handleNameUpdate(updateOpts *domains.UpdateOpts, obj orcObjectPT, osResource *osResourceT) {
	name := getResourceName(obj)
	if osResource.Name != name {
		updateOpts.Name = name
	}
}

func handleDescriptionUpdate(updateOpts *domains.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	description := ptr.Deref(resource.Description, "")
	if osResource.Description != description {
		updateOpts.Description = &description
	}
}

func handleEnabledUpdate(updateOpts *domains.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	Enabled := ptr.Deref(resource.Enabled, true)
	if osResource.Enabled != Enabled {
		updateOpts.Enabled = &Enabled
	}
}

func (actuator domainActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		actuator.updateResource,
	}, nil
}

type domainHelperFactory struct{}

var _ helperFactory = domainHelperFactory{}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.Domain, controller interfaces.ResourceController) (domainActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return domainActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return domainActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewDomainClient()
	if err != nil {
		return domainActuator{}, progress.WrapError(err)
	}

	return domainActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}

func (domainHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return domainAdapter{obj}
}

func (domainHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (domainHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}
