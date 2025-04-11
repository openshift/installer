/*
Copyright 2024 The Kubernetes Authors.

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
	"errors"
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/api/v1alpha1"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/pkg/clients/applyconfiguration/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/pkg/utils/ssa"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/clients"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/orc"
)

//+kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=images,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=images/status,verbs=get;update;patch

func (r *orcImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	orcImage := &orcv1alpha1.Image{}
	err := r.client.Get(ctx, req.NamespacedName, orcImage)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !orcImage.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, orcImage)
	}

	return r.reconcileNormal(ctx, orcImage)
}

func (r *orcImageReconciler) getImageClient(ctx context.Context, orcImage *orcv1alpha1.Image) (clients.ImageClient, error) {
	log := ctrl.LoggerFrom(ctx)

	clientScope, err := r.scopeFactory.NewClientScopeFromObject(ctx, r.client, r.caCertificates, log, orc.IdentityRefProvider(orcImage))
	if err != nil {
		return nil, err
	}
	return clientScope.NewImageClient()
}

func (r *orcImageReconciler) reconcileNormal(ctx context.Context, orcImage *orcv1alpha1.Image) (_ ctrl.Result, err error) {
	log := ctrl.LoggerFrom(ctx)
	log.V(3).Info("Reconciling image")

	if !controllerutil.ContainsFinalizer(orcImage, Finalizer) {
		return ctrl.Result{}, r.setFinalizer(ctx, orcImage)
	}

	var statusOpts []updateStatusOpt
	addStatus := func(opt updateStatusOpt) {
		statusOpts = append(statusOpts, opt)
	}

	// Ensure we always update status
	defer func() {
		if err != nil {
			addStatus(withError(err))
		}

		err = errors.Join(err, r.updateStatus(ctx, orcImage, statusOpts...))

		var terminalError *capoerrors.TerminalError
		if errors.As(err, &terminalError) {
			log.Error(err, "not scheduling further reconciles for terminal error")
			err = nil
		}
	}()

	imageClient, err := r.getImageClient(ctx, orcImage)
	if err != nil {
		return ctrl.Result{}, err
	}

	glanceImage, err := func() (*images.Image, error) {
		if orcImage.Status.ID != nil {
			log.V(4).Info("Fetching existing glance image", "ID", *orcImage.Status.ID)

			image, err := imageClient.GetImage(*orcImage.Status.ID)
			if capoerrors.IsNotFound(err) {
				// An image we previously referenced has been deleted unexpectedly. We can't recover from this.
				err = capoerrors.Terminal(orcv1alpha1.OpenStackConditionReasonUnrecoverableError, "image has been deleted from glance")
			}
			return image, err
		}

		if orcImage.Spec.Import != nil {
			log.V(4).Info("Importing existing glance image")

			if orcImage.Spec.Import.ID != nil {
				image, err := imageClient.GetImage(*orcImage.Spec.Import.ID)
				if capoerrors.IsNotFound(err) {
					// We assume that an image imported by ID must already exist. It's a terminal error if it doesn't.
					err = capoerrors.Terminal(orcv1alpha1.OpenStackConditionReasonUnrecoverableError, "referenced image does not exist in glance")
				}
				return image, err
			}

			listOpts := listOptsFromImportFilter(orcImage.Spec.Import.Filter)
			return getGlanceImageFromList(ctx, listOpts, imageClient)

			// TODO: When we support 'import and manage' we need to implement
			// setting spec.resource from the discovered glance image here.
		}

		log.V(4).Info("Checking for previously created image")

		listOpts := listOptsFromCreation(orcImage)
		return getGlanceImageFromList(ctx, listOpts, imageClient)
	}()
	if err != nil {
		return ctrl.Result{}, err
	}

	if glanceImage == nil {
		if orcImage.Spec.Import != nil {
			log.V(3).Info("Image does not yet exist")
			addStatus(withProgressMessage("Waiting for glance image to be created externally"))

			return ctrl.Result{RequeueAfter: waitForGlanceImageStatusUpdate}, err
		}

		glanceImage, err = createImage(ctx, orcImage, imageClient)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	addStatus(withGlanceImage(glanceImage))

	if orcImage.Status.ID == nil {
		if err := r.setStatusID(ctx, orcImage, glanceImage.ID); err != nil {
			return ctrl.Result{}, err
		}
	}

	log = log.WithValues("ID", glanceImage.ID)
	ctx = ctrl.LoggerInto(ctx, log)

	log.V(4).Info("Got glance image", "status", glanceImage.Status)

	switch glanceImage.Status {
	// Cases where we're not going to take any action until the next resync
	case images.ImageStatusActive, images.ImageStatusDeactivated:
		return ctrl.Result{}, nil

	// Content is being saved. Check back in a minute
	// "importing" is seen during web-download
	// "saving" is seen while uploading, but might be seen because our upload failed and glance hasn't reset yet.
	case images.ImageStatusImporting, images.ImageStatusSaving:
		addStatus(withProgressMessage(downloadingMessage("Glance is downloading image content", orcImage)))
		return ctrl.Result{RequeueAfter: waitForGlanceImageStatusUpdate}, nil

	// Newly created image, waiting for upload, or... previous upload was interrupted and has now reset
	case images.ImageStatusQueued:
		// Don't attempt image creation if we're not managing the image
		if orcImage.Spec.ManagementPolicy == orcv1alpha1.ManagementPolicyUnmanaged {
			addStatus(withProgressMessage("Waiting for glance image content to be uploaded externally"))

			return ctrl.Result{
				RequeueAfter: waitForGlanceImageStatusUpdate,
			}, err
		}

		if ptr.Deref(orcImage.Status.DownloadAttempts, 0) >= maxDownloadAttempts {
			return ctrl.Result{}, capoerrors.Terminal(orcv1alpha1.OpenStackConditionReasonInvalidConfiguration, fmt.Sprintf("Unable to download content after %d attempts", maxDownloadAttempts))
		}

		canWebDownload, err := r.canWebDownload(ctx, orcImage, imageClient)
		if err != nil {
			return ctrl.Result{}, err
		}

		if canWebDownload {
			// We frequently hit a race with glance here. There is a
			// delay after doing an import before glance updates the
			// status from queued, meaning we frequently attempt to
			// start a second import. Although the status isn't
			// updated yet, glance still returns a 409 error when
			// this happens due to the existing task. This is
			// harmless.

			err := r.webDownload(ctx, orcImage, imageClient, glanceImage)
			if err != nil {
				return ctrl.Result{}, err
			}

			// Don't increment DownloadAttempts unless webDownload returned success
			addStatus(withIncrementDownloadAttempts())

			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, r.uploadImageContent(ctx, orcImage, imageClient, glanceImage)
		}

	// Error cases
	case images.ImageStatusKilled:
		return ctrl.Result{}, capoerrors.Terminal(orcv1alpha1.OpenStackConditionReasonUnrecoverableError, "a glance error occurred while saving image content")
	case images.ImageStatusDeleted, images.ImageStatusPendingDelete:
		return ctrl.Result{}, capoerrors.Terminal(orcv1alpha1.OpenStackConditionReasonUnrecoverableError, "image status is deleting")
	default:
		return ctrl.Result{}, errors.New("unknown image status: " + string(glanceImage.Status))
	}
}

func (r *orcImageReconciler) reconcileDelete(ctx context.Context, orcImage *orcv1alpha1.Image) (_ ctrl.Result, err error) {
	log := ctrl.LoggerFrom(ctx)
	log.V(3).Info("Reconciling image delete")

	var statusOpts []updateStatusOpt
	addStatus := func(opt updateStatusOpt) {
		statusOpts = append(statusOpts, opt)
	}

	deleted := false
	defer func() {
		// No point updating status after removing the finalizer
		if !deleted {
			if err != nil {
				addStatus(withError(err))
			}
			err = errors.Join(err, r.updateStatus(ctx, orcImage, statusOpts...))
		}
	}()

	// We won't delete the resource for an unmanaged object, or if onDelete is detach
	if orcImage.Spec.ManagementPolicy == orcv1alpha1.ManagementPolicyUnmanaged || orcImage.Spec.ManagedOptions.GetOnDelete() == orcv1alpha1.OnDeleteDetach {
		logPolicy := []any{"managementPolicy", orcImage.Spec.ManagementPolicy}
		if orcImage.Spec.ManagementPolicy == orcv1alpha1.ManagementPolicyManaged {
			logPolicy = append(logPolicy, "onDelete", orcImage.Spec.ManagedOptions.GetOnDelete())
		}
		log.V(4).Info("Not deleting Glance image due to policy", logPolicy...)
	} else {
		imageClient, err := r.getImageClient(ctx, orcImage)
		if err != nil {
			return ctrl.Result{}, err
		}

		var glanceImage *images.Image
		glanceImage, err = getGlanceImage(ctx, orcImage, imageClient)
		if err != nil && !capoerrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		addStatus(withGlanceImage(glanceImage))

		// Delete any returned glance image, but don't clear the finalizer until getGlanceImage() returns nothing
		if glanceImage != nil {
			log.V(4).Info("Deleting image", "id", glanceImage.ID)
			err := imageClient.DeleteImage(ctx, glanceImage.ID)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
		}

		log.V(4).Info("Image is deleted")
	}

	deleted = true

	// Clear the finalizer
	applyConfig := orcapplyconfigv1alpha1.Image(orcImage.Name, orcImage.Namespace).WithUID(orcImage.UID)
	return ctrl.Result{}, r.client.Patch(ctx, orcImage, ssa.ApplyConfigPatch(applyConfig), client.ForceOwnership, ssaFieldOwner(SSAFinalizerTxn))
}

// getGlanceImage returns the glance image associated with an ORC Image, or nil if none was found.
// If Status.ImageID is set, it returns this image, or an error if it does not exist.
// Otherwise it looks for an existing image with the expected name. It returns nil if none exists.
func getGlanceImage(ctx context.Context, orcImage *orcv1alpha1.Image, imageClient clients.ImageClient) (*images.Image, error) {
	log := ctrl.LoggerFrom(ctx)

	log.V(4).Info("Looking for existing glance image to adopt")

	// Check for existing image by name in case we're adopting or failed to write to status
	imageName := getImageName(orcImage)
	glanceImages, err := imageClient.ListImages(images.ListOpts{Name: imageName})
	if err != nil {
		return nil, err
	}
	switch {
	case len(glanceImages) == 1:
		image := &glanceImages[0]
		log.V(3).Info("Adopting existing glance image", "imageID", image.ID)
		return image, nil
	case len(glanceImages) > 1:
		return nil, capoerrors.Terminal(orcv1alpha1.OpenStackConditionReasonInvalidConfiguration, "found multiple images with name "+imageName)
	}

	return nil, nil
}

// getImageName returns the name of the glance image we should use.
func getImageName(orcImage *orcv1alpha1.Image) string {
	if orcImage.Spec.Resource.Name != "" {
		return orcImage.Spec.Resource.Name
	}
	return orcImage.Name
}

func listOptsFromImportFilter(filter *orcv1alpha1.ImageFilter) images.ListOptsBuilder {
	return images.ListOpts{Name: ptr.Deref(filter.Name, "")}
}

// listOptsFromCreation returns a listOpts which will return the image which
// would have been created from the current spec and hopefully no other image.
// Its purpose is to automatically adopt an image that we created but failed to
// write to status.id.
func listOptsFromCreation(orcImage *orcv1alpha1.Image) images.ListOptsBuilder {
	return images.ListOpts{Name: getImageName(orcImage)}
}

func getGlanceImageFromList(_ context.Context, listOpts images.ListOptsBuilder, imageClient clients.ImageClient) (*images.Image, error) {
	glanceImages, err := imageClient.ListImages(listOpts)
	if err != nil {
		return nil, err
	}

	if len(glanceImages) == 1 {
		return &glanceImages[0], nil
	}

	// No image found
	if len(glanceImages) == 0 {
		return nil, nil
	}

	// Multiple images found
	return nil, capoerrors.Terminal(orcv1alpha1.OpenStackConditionReasonInvalidConfiguration, fmt.Sprintf("Expected to find exactly one image to import. Found %d", len(glanceImages)))
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

// createImage creates a glance image for an ORC Image.
func createImage(ctx context.Context, orcImage *orcv1alpha1.Image, imageClient clients.ImageClient) (*images.Image, error) {
	if orcImage.Spec.ManagementPolicy == orcv1alpha1.ManagementPolicyUnmanaged {
		// Should have been caught by API validation
		return nil, capoerrors.Terminal(orcv1alpha1.OpenStackConditionReasonInvalidConfiguration, "Not creating unmanaged resource")
	}

	log := ctrl.LoggerFrom(ctx)
	log.V(3).Info("Creating image")

	resource := orcImage.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, capoerrors.Terminal(orcv1alpha1.OpenStackConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set")
	}

	if resource.Content == nil {
		// Should have been caught by API validation
		return nil, capoerrors.Terminal(orcv1alpha1.OpenStackConditionReasonInvalidConfiguration, "Creation requested, but spec.resource.content is not set")
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
			minDisk = *properties.MinDiskGB
		}
		if properties.MinMemoryMB != nil {
			minMemory = *properties.MinMemoryMB
		}

		if err := glancePropertiesFromStruct(properties.Hardware, additionalProperties); err != nil {
			return nil, capoerrors.Terminal(orcv1alpha1.OpenStackConditionReasonUnrecoverableError, "programming error", err)
		}
	}

	var visibility *images.ImageVisibility
	if resource.Visibility != nil {
		visibility = ptr.To(images.ImageVisibility(*resource.Visibility))
	}

	image, err := imageClient.CreateImage(ctx, &images.CreateOpts{
		Name:            getImageName(orcImage),
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
	if capoerrors.IsConflict(err) {
		err = capoerrors.Terminal(orcv1alpha1.OpenStackConditionReasonInvalidConfiguration, "invalid configuration creating image: "+err.Error(), err)
	}

	return image, err
}

func downloadingMessage(msg string, orcImage *orcv1alpha1.Image) string {
	if ptr.Deref(orcImage.Status.DownloadAttempts, 0) > 1 {
		return fmt.Sprintf("%s: attempt %d", msg, *orcImage.Status.DownloadAttempts)
	}
	return msg
}
