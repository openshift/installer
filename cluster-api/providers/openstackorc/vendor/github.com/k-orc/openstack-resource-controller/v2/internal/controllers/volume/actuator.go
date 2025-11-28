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

package volume

import (
	"context"
	"iter"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
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
	osResourceT = volumes.Volume

	createResourceActuator = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	resourceReconciler     = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory          = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

const (
	// The frequency to poll when waiting for the resource to become available
	volumeAvailablePollingPeriod = 15 * time.Second
	// The frequency to poll when waiting for the resource to be deleted
	volumeDeletingPollingPeriod = 15 * time.Second
)

type volumeActuator struct {
	osClient  osclients.VolumeClient
	k8sClient client.Client
}

var _ createResourceActuator = volumeActuator{}
var _ deleteResourceActuator = volumeActuator{}

func (volumeActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator volumeActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	resource, err := actuator.osClient.GetVolume(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return resource, nil
}

func (actuator volumeActuator) ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*osResourceT, error], bool) {
	resourceSpec := orcObject.Spec.Resource
	if resourceSpec == nil {
		return nil, false
	}

	var filters []osclients.ResourceFilter[osResourceT]

	// NOTE: The API doesn't allow filtering by description or size
	// we'll have to do it client-side.
	if resourceSpec.Description != nil {
		filters = append(filters, func(f *volumes.Volume) bool {
			return f.Description == *resourceSpec.Description
		})
	}
	filters = append(filters, func(f *volumes.Volume) bool {
		return f.Size == int(resourceSpec.Size)
	})

	metadata := make(map[string]string)
	for _, m := range resourceSpec.Metadata {
		metadata[m.Name] = m.Value
	}

	listOpts := volumes.ListOpts{
		Name:     getResourceName(orcObject),
		Metadata: metadata,
	}

	return actuator.listOSResources(ctx, filters, listOpts), true
}

func (actuator volumeActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	var filters []osclients.ResourceFilter[osResourceT]

	// NOTE: The API doesn't allow filtering by description or size
	// we'll have to do it client-side.
	if filter.Description != nil {
		filters = append(filters, func(f *volumes.Volume) bool {
			return f.Description == *filter.Description
		})
	}
	if filter.Size != nil {
		filters = append(filters, func(f *volumes.Volume) bool {
			return f.Size == int(*filter.Size)
		})
	}

	listOpts := volumes.ListOpts{
		Name: string(ptr.Deref(filter.Name, "")),
	}

	return actuator.listOSResources(ctx, filters, listOpts), nil
}

func (actuator volumeActuator) listOSResources(ctx context.Context, filters []osclients.ResourceFilter[osResourceT], listOpts volumes.ListOptsBuilder) iter.Seq2[*volumes.Volume, error] {
	volumes := actuator.osClient.ListVolumes(ctx, listOpts)
	return osclients.Filter(volumes, filters...)
}

func (actuator volumeActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}
	var reconcileStatus progress.ReconcileStatus

	var volumetypeID string
	if resource.VolumeTypeRef != nil {
		volumetype, volumetypeDepRS := volumetypeDependency.GetDependency(
			ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.VolumeType) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(volumetypeDepRS)
		if volumetype != nil {
			volumetypeID = ptr.Deref(volumetype.Status.ID, "")
		}
	}

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	metadata := make(map[string]string)
	for _, m := range resource.Metadata {
		metadata[m.Name] = m.Value
	}

	createOpts := volumes.CreateOpts{
		Name:        getResourceName(obj),
		Description: ptr.Deref(resource.Description, ""),
		Size:        int(resource.Size),
		Metadata:    metadata,
		VolumeType:  volumetypeID,
	}

	osResource, err := actuator.osClient.CreateVolume(ctx, createOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a conflict
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator volumeActuator) DeleteResource(ctx context.Context, _ orcObjectPT, resource *osResourceT) progress.ReconcileStatus {
	if resource.Status == VolumeStatusDeleting {
		return progress.WaitingOnOpenStack(progress.WaitingOnReady, volumeDeletingPollingPeriod)
	}

	// FIXME(mandre) Make this optional
	deleteOpts := volumes.DeleteOpts{
		Cascade: false,
	}

	return progress.WrapError(actuator.osClient.DeleteVolume(ctx, resource.ID, deleteOpts))
}

func (actuator volumeActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := volumes.UpdateOpts{}

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

	_, err = actuator.osClient.UpdateVolume(ctx, osResource.ID, updateOpts)

	// We should require the spec to be updated before retrying an update which returned a conflict
	if orcerrors.IsConflict(err) {
		err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
	}

	if err != nil {
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts volumes.UpdateOpts) (bool, error) {
	updateOptsMap, err := updateOpts.ToVolumeUpdateMap()
	if err != nil {
		return false, err
	}

	updateMap, ok := updateOptsMap["volume"].(map[string]any)
	if !ok {
		updateMap = make(map[string]any)
	}

	return len(updateMap) > 0, nil
}

func handleNameUpdate(updateOpts *volumes.UpdateOpts, obj orcObjectPT, osResource *osResourceT) {
	name := getResourceName(obj)
	if osResource.Name != name {
		updateOpts.Name = &name
	}
}

func handleDescriptionUpdate(updateOpts *volumes.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	description := ptr.Deref(resource.Description, "")
	if osResource.Description != description {
		updateOpts.Description = &description
	}
}

func (actuator volumeActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		actuator.updateResource,
	}, nil
}

type volumeHelperFactory struct{}

var _ helperFactory = volumeHelperFactory{}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.Volume, controller interfaces.ResourceController) (volumeActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return volumeActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return volumeActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewVolumeClient()
	if err != nil {
		return volumeActuator{}, progress.WrapError(err)
	}

	return volumeActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}

func (volumeHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return volumeAdapter{obj}
}

func (volumeHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (volumeHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}
