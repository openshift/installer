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

package service

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/services"
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
	osResourceT = services.Service

	createResourceActuator = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	resourceReconciler     = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory          = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

type serviceActuator struct {
	osClient  osclients.ServiceClient
	k8sClient client.Client
}

var _ createResourceActuator = serviceActuator{}
var _ deleteResourceActuator = serviceActuator{}

func (serviceActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator serviceActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	resource, err := actuator.osClient.GetService(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return resource, nil
}

func (actuator serviceActuator) ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*osResourceT, error], bool) {
	resourceSpec := orcObject.Spec.Resource
	if resourceSpec == nil {
		return nil, false
	}

	listOpts := services.ListOpts{
		Name:        getResourceName(orcObject),
		ServiceType: resourceSpec.Type,
	}

	return actuator.osClient.ListServices(ctx, listOpts), true
}

func (actuator serviceActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	listOpts := services.ListOpts{
		Name:        string(ptr.Deref(filter.Name, "")),
		ServiceType: ptr.Deref(filter.Type, ""),
	}

	return actuator.osClient.ListServices(ctx, listOpts), nil
}

func (actuator serviceActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}

	extra := map[string]any{
		"name":        getResourceName(obj),
		"description": resource.Description,
	}

	createOpts := services.CreateOpts{
		Type:    resource.Type,
		Enabled: resource.Enabled,
		Extra:   extra,
	}

	osResource, err := actuator.osClient.CreateService(ctx, createOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a conflict
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator serviceActuator) DeleteResource(ctx context.Context, _ orcObjectPT, resource *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteService(ctx, resource.ID))
}

func (actuator serviceActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := services.UpdateOpts{Extra: make(map[string]any)}

	handleNameUpdate(&updateOpts, obj, osResource)
	handleDescriptionUpdate(&updateOpts, resource, osResource)
	handleTypeUpdate(&updateOpts, resource, osResource)
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

	// NOTE(winiciusallan): we need to add Type before every update to avoid
	// gophercloud create UpdateOpts with type as empty value. for more
	// information https://github.com/gophercloud/gophercloud/issues/3553
	updateOpts.Type = resource.Type

	_, err = actuator.osClient.UpdateService(ctx, osResource.ID, updateOpts)

	if err != nil {
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
		}
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts services.UpdateOpts) (bool, error) {
	updateOptsMap, err := updateOpts.ToServiceUpdateMap()
	if err != nil {
		return false, err
	}

	updateMap, ok := updateOptsMap["service"].(map[string]any)
	if !ok {
		updateMap = make(map[string]any)
	}

	if serviceType, ok := updateMap["type"]; ok && serviceType == "" {
		delete(updateMap, "type")
	}

	return len(updateMap) > 0, nil
}

func handleNameUpdate(updateOpts *services.UpdateOpts, obj orcObjectPT, osResource *osResourceT) {
	name := getResourceName(obj)
	if osResource.Extra["name"] != "" {
		if osResource.Extra["name"] != name {
			updateOpts.Extra["name"] = name
		}
	}
}

func handleDescriptionUpdate(updateOpts *services.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	curr, ok := osResource.Extra["description"]

	if resource.Description == nil {
		if curr == nil {
			return
		}
		updateOpts.Extra["description"] = nil
		return
	}

	if !ok || curr != *resource.Description {
		updateOpts.Extra["description"] = *resource.Description
	}
}

func handleTypeUpdate(updateOpts *services.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	if osResource.Type != resource.Type {
		updateOpts.Type = resource.Type
	}
}

func handleEnabledUpdate(updateOpts *services.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	enabled := ptr.Deref(resource.Enabled, true)
	if osResource.Enabled != enabled {
		updateOpts.Enabled = &enabled
	}
}

func (actuator serviceActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		actuator.updateResource,
	}, nil
}

type serviceHelperFactory struct{}

var _ helperFactory = serviceHelperFactory{}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.Service, controller interfaces.ResourceController) (serviceActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return serviceActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return serviceActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewServiceClient()
	if err != nil {
		return serviceActuator{}, progress.WrapError(err)
	}

	return serviceActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}

func (serviceHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return serviceAdapter{obj}
}

func (serviceHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (serviceHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}
