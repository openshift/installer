/*
Copyright The ORC Authors.

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

package user

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	"github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/applyconfigs"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

// OpenStack resource types
type (
	osResourceT = users.User

	createResourceActuator = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	resourceReconciler     = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory          = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

type userActuator struct {
	osClient  osclients.UserClient
	k8sClient client.Client
}

var _ createResourceActuator = userActuator{}
var _ deleteResourceActuator = userActuator{}

func (userActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator userActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	resource, err := actuator.osClient.GetUser(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return resource, nil
}

func (actuator userActuator) ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*osResourceT, error], bool) {
	resourceSpec := orcObject.Spec.Resource
	if resourceSpec == nil {
		return nil, false
	}

	// Resolve the domain ID from DomainRef if set. Without the domain
	// ID, adoption could match a user in the wrong domain.
	var domainID string
	if resourceSpec.DomainRef != nil {
		domain, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, orcObject.Namespace, resourceSpec.DomainRef, "Domain",
			func(dep *orcv1alpha1.Domain) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := rs.NeedsReschedule(); needsReschedule {
			return nil, false
		}
		domainID = ptr.Deref(domain.Status.ID, "")
	}

	listOpts := users.ListOpts{
		Name:     getResourceName(orcObject),
		DomainID: domainID,
	}

	return actuator.osClient.ListUsers(ctx, listOpts), true
}

func (actuator userActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	var reconcileStatus progress.ReconcileStatus

	domain, rs := dependency.FetchDependency[*orcv1alpha1.Domain](
		ctx, actuator.k8sClient, obj.Namespace,
		filter.DomainRef, "Domain",
		orcv1alpha1.IsAvailable,
	)
	reconcileStatus = reconcileStatus.WithReconcileStatus(rs)

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	listOpts := users.ListOpts{
		Name:     string(ptr.Deref(filter.Name, "")),
		DomainID: ptr.Deref(domain.Status.ID, ""),
	}

	return actuator.osClient.ListUsers(ctx, listOpts), reconcileStatus
}

func (actuator userActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}
	var reconcileStatus progress.ReconcileStatus

	var domainID string
	if resource.DomainRef != nil {
		domain, domainDepRS := domainDependency.GetDependency(
			ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(domainDepRS)
		if domain != nil {
			domainID = ptr.Deref(domain.Status.ID, "")
		}
	}

	var defaultProjectID string
	if resource.DefaultProjectRef != nil {
		project, projectDepRS := projectDependency.GetDependency(
			ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(projectDepRS)
		if project != nil {
			defaultProjectID = ptr.Deref(project.Status.ID, "")
		}
	}

	var password string
	if resource.PasswordRef != nil {
		secret, secretReconcileStatus := dependency.FetchDependency(
			ctx, actuator.k8sClient, obj.Namespace,
			resource.PasswordRef, "Secret",
			func(*corev1.Secret) bool { return true },
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(secretReconcileStatus)
		if secretReconcileStatus == nil {
			passwordBytes, ok := secret.Data["password"]
			if !ok {
				reconcileStatus = reconcileStatus.WithReconcileStatus(
					progress.NewReconcileStatus().WithProgressMessage("Password secret does not contain \"password\" key"))
			} else {
				password = string(passwordBytes)
			}
		}
	}

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}
	createOpts := users.CreateOpts{
		Name:             getResourceName(obj),
		Description:      ptr.Deref(resource.Description, ""),
		DomainID:         domainID,
		Enabled:          resource.Enabled,
		DefaultProjectID: defaultProjectID,
		Password:         password,
	}

	osResource, err := actuator.osClient.CreateUser(ctx, createOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a conflict
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator userActuator) DeleteResource(ctx context.Context, _ orcObjectPT, resource *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteUser(ctx, resource.ID))
}

func (actuator userActuator) reconcilePassword(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil || resource.PasswordRef == nil {
		return nil
	}

	currentRef := string(*resource.PasswordRef)
	var lastAppliedRef string
	if obj.Status.Resource != nil {
		lastAppliedRef = obj.Status.Resource.AppliedPasswordRef
	}

	if lastAppliedRef == currentRef {
		return nil
	}

	// Read the password from the referenced Secret
	secret, secretRS := dependency.FetchDependency(
		ctx, actuator.k8sClient, obj.Namespace,
		resource.PasswordRef, "Secret",
		func(*corev1.Secret) bool { return true },
	)
	if secretRS != nil {
		return secretRS
	}

	passwordBytes, ok := secret.Data["password"]
	if !ok {
		return progress.NewReconcileStatus().WithProgressMessage("Password secret does not contain \"password\" key")
	}
	password := string(passwordBytes)

	// Only call UpdateUser if this is not the first reconcile after creation.
	// CreateResource already set the initial password.
	if lastAppliedRef != "" {
		log.V(logging.Info).Info("Updating password")
		_, err := actuator.osClient.UpdateUser(ctx, osResource.ID, users.UpdateOpts{
			Password: password,
		})

		if err != nil {
			if !orcerrors.IsRetryable(err) {
				err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
			}
			return progress.WrapError(err)
		}
	}

	// Update the lastAppliedPasswordRef status field via a MergePatch.
	// MergePatch sets only the specified fields without claiming SSA
	// ownership, so the main SSA status update won't remove this field.
	statusApply := orcapplyconfigv1alpha1.UserResourceStatus().
		WithAppliedPasswordRef(currentRef)
	applyConfig := orcapplyconfigv1alpha1.User(obj.Name, obj.Namespace).
		WithUID(obj.UID).
		WithStatus(orcapplyconfigv1alpha1.UserStatus().
			WithResource(statusApply))
	if err := actuator.k8sClient.Status().Patch(ctx, obj,
		applyconfigs.Patch(types.MergePatchType, applyConfig)); err != nil {
		return progress.WrapError(err)
	}

	if lastAppliedRef != "" {
		return progress.NeedsRefresh()
	}
	return nil
}

func (actuator userActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := users.UpdateOpts{}

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

	_, err = actuator.osClient.UpdateUser(ctx, osResource.ID, updateOpts)

	if err != nil {
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
		}
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts users.UpdateOpts) (bool, error) {
	updateOptsMap, err := updateOpts.ToUserUpdateMap()
	if err != nil {
		return false, err
	}

	updateMap, ok := updateOptsMap["user"].(map[string]any)
	if !ok {
		updateMap = make(map[string]any)
	}

	return len(updateMap) > 0, nil
}

func handleNameUpdate(updateOpts *users.UpdateOpts, obj orcObjectPT, osResource *osResourceT) {
	name := getResourceName(obj)
	if osResource.Name != name {
		updateOpts.Name = name
	}
}

func handleDescriptionUpdate(updateOpts *users.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	description := ptr.Deref(resource.Description, "")
	if osResource.Description != description {
		updateOpts.Description = &description
	}
}

func handleEnabledUpdate(updateOpts *users.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	enabled := ptr.Deref(resource.Enabled, true)
	if osResource.Enabled != enabled {
		updateOpts.Enabled = &enabled
	}
}

func (actuator userActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		actuator.reconcilePassword,
		actuator.updateResource,
	}, nil
}

type userHelperFactory struct{}

var _ helperFactory = userHelperFactory{}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.User, controller interfaces.ResourceController) (userActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return userActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return userActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewUserClient()
	if err != nil {
		return userActuator{}, progress.WrapError(err)
	}

	return userActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}

func (userHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return userAdapter{obj}
}

func (userHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (userHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}
