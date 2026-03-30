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

package volumetype

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
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
	osResourceT = volumetypes.VolumeType

	createResourceActuator = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	resourceReconciler     = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory          = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

type volumetypeActuator struct {
	osClient  osclients.VolumeTypeClient
	k8sClient client.Client
}

var _ createResourceActuator = volumetypeActuator{}
var _ deleteResourceActuator = volumetypeActuator{}

func (volumetypeActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator volumetypeActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	resource, err := actuator.osClient.GetVolumeType(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return resource, nil
}

func (actuator volumetypeActuator) ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*osResourceT, error], bool) {
	resourceSpec := orcObject.Spec.Resource
	if resourceSpec == nil {
		return nil, false
	}

	var filters []osclients.ResourceFilter[osResourceT]

	// NOTE: The API doesn't allow filtering by name or description, we'll have to do it client-side.
	filters = append(filters,
		func(f *volumetypes.VolumeType) bool {
			name := getResourceName(orcObject)
			// Compare non-pointer values
			return f.Name == name
		},
	)
	if resourceSpec.Description != nil {
		filters = append(filters, func(f *volumetypes.VolumeType) bool {
			return f.Description == *resourceSpec.Description
		})
	}

	isPublic := volumetypes.VisibilityDefault
	if resourceSpec.IsPublic != nil {
		if *resourceSpec.IsPublic {
			isPublic = volumetypes.VisibilityPublic
		} else {
			isPublic = volumetypes.VisibilityPrivate
		}
	}

	listOpts := volumetypes.ListOpts{
		IsPublic: isPublic,
	}

	return actuator.listOSResources(ctx, filters, listOpts), true
}

func (actuator volumetypeActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	var filters []osclients.ResourceFilter[osResourceT]

	// NOTE: The API doesn't allow filtering by name or description, we'll have to do it client-side.
	if filter.Name != nil {
		filters = append(filters, func(f *volumetypes.VolumeType) bool {
			return f.Name == string(*filter.Name)
		})
	}
	if filter.Description != nil {
		filters = append(filters, func(f *volumetypes.VolumeType) bool {
			return f.Description == *filter.Description
		})
	}

	isPublic := volumetypes.VisibilityDefault
	if filter.IsPublic != nil {
		if *filter.IsPublic {
			isPublic = volumetypes.VisibilityPublic
		} else {
			isPublic = volumetypes.VisibilityPrivate
		}
	}

	listOpts := volumetypes.ListOpts{
		IsPublic: isPublic,
	}

	return actuator.listOSResources(ctx, filters, listOpts), nil
}

func (actuator volumetypeActuator) listOSResources(ctx context.Context, filters []osclients.ResourceFilter[osResourceT], listOpts volumetypes.ListOptsBuilder) iter.Seq2[*volumetypes.VolumeType, error] {
	volumetypes := actuator.osClient.ListVolumeTypes(ctx, listOpts)
	return osclients.Filter(volumetypes, filters...)
}

func (actuator volumetypeActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}

	extraSpecs := make(map[string]string)
	for _, spec := range resource.ExtraSpecs {
		extraSpecs[spec.Name] = spec.Value
	}

	createOpts := volumetypes.CreateOpts{
		Name:        getResourceName(obj),
		Description: ptr.Deref(resource.Description, ""),
		IsPublic:    resource.IsPublic,
		ExtraSpecs:  extraSpecs,
	}

	osResource, err := actuator.osClient.CreateVolumeType(ctx, createOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a conflict
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator volumetypeActuator) DeleteResource(ctx context.Context, _ orcObjectPT, resource *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteVolumeType(ctx, resource.ID))
}

func (actuator volumetypeActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := volumetypes.UpdateOpts{}

	handleNameUpdate(&updateOpts, obj, osResource)
	handleDescriptionUpdate(&updateOpts, resource, osResource)
	handleIsPublicUpdate(&updateOpts, resource, osResource)

	needsUpdate, err := needsUpdate(updateOpts)
	if err != nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err))
	}
	if !needsUpdate {
		log.V(logging.Debug).Info("No changes")
		return nil
	}

	_, err = actuator.osClient.UpdateVolumeType(ctx, osResource.ID, updateOpts)

	// We should require the spec to be updated before retrying an update which returned a conflict
	if orcerrors.IsConflict(err) {
		err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
	}

	if err != nil {
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts volumetypes.UpdateOpts) (bool, error) {
	updateOptsMap, err := updateOpts.ToVolumeTypeUpdateMap()
	if err != nil {
		return false, err
	}

	updateMap, ok := updateOptsMap["volume_type"].(map[string]any)
	if !ok {
		updateMap = make(map[string]any)
	}

	return len(updateMap) > 0, nil
}

func handleNameUpdate(updateOpts *volumetypes.UpdateOpts, obj orcObjectPT, osResource *osResourceT) {
	name := getResourceName(obj)
	if osResource.Name != name {
		updateOpts.Name = &name
	}
}

func handleDescriptionUpdate(updateOpts *volumetypes.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	description := ptr.Deref(resource.Description, "")
	if osResource.Description != description {
		updateOpts.Description = &description
	}
}

func handleIsPublicUpdate(updateOpts *volumetypes.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	// Default is true
	isPublic := ptr.Deref(resource.IsPublic, true)
	if osResource.IsPublic != isPublic {
		updateOpts.IsPublic = &isPublic
	}
}

func (actuator volumetypeActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		actuator.updateResource,
	}, nil
}

type volumetypeHelperFactory struct{}

var _ helperFactory = volumetypeHelperFactory{}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.VolumeType, controller interfaces.ResourceController) (volumetypeActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return volumetypeActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return volumetypeActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewVolumeTypeClient()
	if err != nil {
		return volumetypeActuator{}, progress.WrapError(err)
	}

	return volumetypeActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}

func (volumetypeHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return volumetypeAdapter{obj}
}

func (volumetypeHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (volumetypeHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}
