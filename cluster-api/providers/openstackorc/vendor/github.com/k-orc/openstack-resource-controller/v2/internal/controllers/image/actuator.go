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

package image

import (
	"context"
	"fmt"
	"iter"
	"reflect"
	"slices"

	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type (
	osResourceT = images.Image

	createResourceActuator    = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator    = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	reconcileResourceActuator = interfaces.ReconcileResourceActuator[orcObjectPT, osResourceT]
	resourceReconciler        = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory             = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
	imageIterator             = iter.Seq2[*osResourceT, error]
)

type imageActuator struct {
	osClient  osclients.ImageClient
	k8sClient client.Client
}

var _ createResourceActuator = imageActuator{}
var _ deleteResourceActuator = imageActuator{}

func (imageActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator imageActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	image, err := actuator.osClient.GetImage(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return image, nil
}

func (actuator imageActuator) ListOSResourcesForAdoption(ctx context.Context, obj orcObjectPT) (imageIterator, bool) {
	if obj.Spec.Resource == nil {
		return nil, false
	}

	listOpts := images.ListOpts{
		Name: getResourceName(obj),
	}
	if obj.Spec.Resource.Visibility != nil {
		listOpts.Visibility = images.ImageVisibility(*obj.Spec.Resource.Visibility)
	}

	if len(obj.Spec.Resource.Tags) > 0 {
		listOpts.Tags = make([]string, len(obj.Spec.Resource.Tags))
		for i := range obj.Spec.Resource.Tags {
			listOpts.Tags[i] = string(obj.Spec.Resource.Tags[i])
		}
	}

	existingImage := actuator.osClient.ListImages(ctx, listOpts)
	return existingImage, true
}

func (actuator imageActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (imageIterator, progress.ReconcileStatus) {
	listOpts := images.ListOpts{
		Name: string(ptr.Deref(filter.Name, "")),
	}
	if filter.Visibility != nil {
		listOpts.Visibility = images.ImageVisibility(*filter.Visibility)
	}

	if len(filter.Tags) > 0 {
		listOpts.Tags = make([]string, len(filter.Tags))
		for i := range filter.Tags {
			listOpts.Tags[i] = string(filter.Tags[i])
		}
	}

	return actuator.osClient.ListImages(ctx, listOpts), nil
}

func (actuator imageActuator) CreateResource(ctx context.Context, obj *orcv1alpha1.Image) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}

	if resource.Content == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource.content is not set"))
	}

	tags := make([]string, len(resource.Tags))
	for i := range resource.Tags {
		tags[i] = string(resource.Tags[i])
	}
	// Sort tags before creation to simplify comparisons
	slices.Sort(tags)

	var minDisk, minMemory int
	properties := resource.Properties
	additionalProperties := map[string]string{}
	if properties != nil {
		if properties.MinDiskGB != nil {
			minDisk = int(*properties.MinDiskGB)
		}
		if properties.MinMemoryMB != nil {
			minMemory = int(*properties.MinMemoryMB)
		}

		if err := glancePropertiesFromStruct(properties.Hardware, additionalProperties); err != nil {
			return nil, progress.WrapError(orcerrors.Terminal(orcv1alpha1.ConditionReasonUnrecoverableError, "programming error", err))
		}
		if err := glancePropertiesFromStruct(properties.OperatingSystem, additionalProperties); err != nil {
			return nil, progress.WrapError(orcerrors.Terminal(orcv1alpha1.ConditionReasonUnrecoverableError, "programming error", err))
		}
		if properties.Architecture != nil {
			additionalProperties["architecture"] = *properties.Architecture
		}
		if properties.HypervisorType != nil {
			additionalProperties["hypervisor_type"] = *properties.HypervisorType
		}
	}

	var visibility *images.ImageVisibility
	if resource.Visibility != nil {
		visibility = ptr.To(images.ImageVisibility(*resource.Visibility))
	}

	image, err := actuator.osClient.CreateImage(ctx, &images.CreateOpts{
		Name:            getResourceName(obj),
		Visibility:      visibility,
		Tags:            tags,
		ContainerFormat: string(resource.Content.ContainerFormat),
		DiskFormat:      (string)(resource.Content.DiskFormat),
		MinDisk:         minDisk,
		MinRAM:          minMemory,
		Protected:       resource.Protected,
		Properties:      additionalProperties,
	})

	// We should require the spec to be updated before retrying a create which returned a conflict
	if orcerrors.IsConflict(err) {
		err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating image: "+err.Error(), err)
	}

	if err != nil {
		return nil, progress.WrapError(err)
	}
	return image, nil
}

func (actuator imageActuator) DeleteResource(ctx context.Context, _ orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteImage(ctx, osResource.ID))
}

func (actuator imageActuator) UpdateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := images.UpdateOpts{}

	updateOpts = handleNameUpdate(updateOpts, obj, osResource)
	updateOpts = handleVisibilityUpdate(updateOpts, resource, osResource)
	updateOpts = handleProtectedUpdate(updateOpts, resource, osResource)
	updateOpts = handleTagsUpdate(updateOpts, resource, osResource)

	if !needsUpdate(updateOpts) {
		log.V(logging.Debug).Info("No changes")
		return nil
	}

	_, err := actuator.osClient.UpdateImage(ctx, osResource.ID, updateOpts)

	if orcerrors.IsConflict(err) {
		err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
	}
	if err != nil {
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts images.UpdateOpts) bool {
	return len(updateOpts) > 0
}

func handleNameUpdate(updateOpts images.UpdateOpts, obj orcObjectPT, osResource *osResourceT) images.UpdateOpts {
	name := getResourceName(obj)

	if osResource.Name != name {
		patch := images.ReplaceImageName{
			NewName: name,
		}
		updateOpts = append(updateOpts, patch)
	}
	return updateOpts
}

func handleVisibilityUpdate(updateOpts images.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) images.UpdateOpts {
	desiredVisibility := resource.Visibility

	if desiredVisibility == nil {
		return updateOpts
	}

	visValue := images.ImageVisibility(*desiredVisibility)

	if osResource.Visibility != visValue {
		patch := images.UpdateVisibility{
			Visibility: visValue,
		}
		updateOpts = append(updateOpts, patch)
	}

	return updateOpts
}

func handleProtectedUpdate(updateOpts images.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) images.UpdateOpts {
	protectedValue := ptr.Deref(resource.Protected, false)

	if osResource.Protected != protectedValue {
		patch := images.ReplaceImageProtected{
			NewProtected: protectedValue,
		}
		updateOpts = append(updateOpts, patch)
	}

	return updateOpts
}

func handleTagsUpdate(updateOpts images.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) images.UpdateOpts {
	DesiredTags := []string{}
	if resource.Tags != nil {
		DesiredTags = make([]string, len(resource.Tags))
		for i, tag := range resource.Tags {
			DesiredTags[i] = string(tag)
		}
	}

	slices.Sort(DesiredTags)
	slices.Sort(osResource.Tags)

	if !slices.Equal(DesiredTags, osResource.Tags) {
		patch := images.ReplaceImageTags{
			NewTags: DesiredTags,
		}
		updateOpts = append(updateOpts, patch)
	}

	return updateOpts
}

// glancePropertiesFromStruct populates a properties struct using field values and glance tags defined on the given struct
// glance tags are defined in the API.
func glancePropertiesFromStruct(propStruct interface{}, properties map[string]string) error {
	sp := reflect.ValueOf(propStruct)
	if sp.Kind() != reflect.Pointer {
		return fmt.Errorf("glancePropertiesFromStruct expects pointer to struct, got %T", propStruct)
	}
	if sp.IsZero() {
		return nil
	}

	s := sp.Elem()
	st := s.Type()
	if st.Kind() != reflect.Struct {
		return fmt.Errorf("glancePropertiesFromStruct expects pointer to struct, got %T", propStruct)
	}

	for i := range st.NumField() {
		field := st.Field(i)
		glanceTag, ok := field.Tag.Lookup(orcv1alpha1.GlanceTag)
		if !ok {
			panic(fmt.Errorf("glance tag not defined for field %s on struct %T", field.Name, st.Name))
		}

		value := s.Field(i)
		if value.Kind() == reflect.Pointer {
			if value.IsZero() {
				continue
			}
			value = value.Elem()
		}

		// Gophercloud takes only strings, but values may not be
		// strings. Value.String() prints semantic information for
		// non-strings, but Sprintf does what we want.
		properties[glanceTag] = fmt.Sprintf("%v", value)
	}

	return nil
}

var _ reconcileResourceActuator = imageActuator{}

func (actuator imageActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		actuator.handleUpload,
		actuator.UpdateResource,
	}, nil
}

func (actuator imageActuator) handleUpload(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)

	switch osResource.Status {

	// Cases where we're not going to take any action until the next resync
	case images.ImageStatusActive, images.ImageStatusDeactivated:
		return progress.WrapError(setDownloadingStatus(ctx, false, "Data saved", orcv1alpha1.ConditionReasonSuccess, metav1.ConditionFalse, orcObject, actuator.k8sClient))

	// Content is being saved. Check back in a minute
	// "importing" is seen during web-download
	// "saving" is seen while uploading, but might be seen because our upload failed and glance hasn't reset yet.
	case images.ImageStatusImporting, images.ImageStatusSaving:
		return progress.NewReconcileStatus().
			WithProgressMessage("Glance is downloading image content").
			WithRequeue(externalUpdatePollingPeriod)

	// Newly created image, waiting for upload, or... previous upload was interrupted and has now reset
	case images.ImageStatusQueued:
		// Don't attempt image creation if we're not managing the image
		if orcObject.Spec.ManagementPolicy == orcv1alpha1.ManagementPolicyUnmanaged {
			return progress.NewReconcileStatus().
				WithProgressMessage("Waiting for glance image content to be uploaded externally").
				WithRequeue(externalUpdatePollingPeriod)
		}

		// Initialize download status
		if orcObject.Status.DownloadAttempts == nil {
			err := setDownloadingStatus(ctx, false, "Starting image upload", orcv1alpha1.ConditionReasonProgressing, metav1.ConditionTrue, orcObject, actuator.k8sClient)
			if err != nil {
				return progress.WrapError(err)
			}

			return progress.NewReconcileStatus().
				WithProgressMessage("Starting image upload")
		}

		if ptr.Deref(orcObject.Status.DownloadAttempts, 0) >= maxDownloadAttempts {
			return progress.WrapError(
				orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, fmt.Sprintf("Unable to download content after %d attempts", maxDownloadAttempts)))
		}

		canWebDownload, err := actuator.canWebDownload(ctx, orcObject)
		if err != nil {
			return progress.WrapError(err)
		}

		if canWebDownload {
			// We frequently hit a race with glance here. There is a
			// delay after doing an import before glance updates the
			// status from queued, meaning we frequently attempt to
			// start a second import. Although the status isn't
			// updated yet, glance still returns a 409 error when
			// this happens due to the existing task. This is
			// harmless.

			err := actuator.webDownload(ctx, orcObject, osResource)
			if err != nil {
				return progress.WrapError(err)
			}

			// Don't increment DownloadAttempts unless webDownload returned success
			err = setDownloadingStatus(ctx, true, "Web download in progress", orcv1alpha1.ConditionReasonProgressing, metav1.ConditionTrue, orcObject, actuator.k8sClient)
			if err != nil {
				return progress.WrapError(err)
			}

			return progress.WaitingOnOpenStack(progress.WaitingOnReady, externalUpdatePollingPeriod)
		} else {
			err := actuator.uploadImageContent(ctx, orcObject, osResource)
			if err != nil {
				return progress.WrapError(err)
			}
			return progress.NeedsRefresh()
		}

	// Error cases
	case images.ImageStatusKilled:
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonUnrecoverableError, "a glance error occurred while saving image content"))

	case images.ImageStatusDeleted, images.ImageStatusPendingDelete:
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonUnrecoverableError, "image status is deleting"))

	default:
		log.V(logging.Verbose).Info("Waiting for OpenStack resource to be ACTIVE")
		return progress.WaitingOnOpenStack(progress.WaitingOnReady, externalUpdatePollingPeriod)
	}
}

type imageHelperFactory struct{}

var _ helperFactory = imageHelperFactory{}

func (imageHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return imageAdapter{obj}
}

func (imageHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, controller, orcObject)
}

func (imageHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, controller, orcObject)
}

func newActuator(ctx context.Context, controller interfaces.ResourceController, orcObject *orcv1alpha1.Image) (imageActuator, progress.ReconcileStatus) {
	if orcObject == nil {
		return imageActuator{}, progress.WrapError(fmt.Errorf("orcObject may not be nil"))
	}

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return imageActuator{}, reconcileStatus
	}

	log := ctrl.LoggerFrom(ctx)
	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return imageActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewImageClient()
	if err != nil {
		return imageActuator{}, progress.WrapError(err)
	}

	return imageActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}
