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

package addressscope

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/addressscopes"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	"github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
)

// OpenStack resource types
type (
	osResourceT = addressscopes.AddressScope

	createResourceActuator = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	resourceReconciler     = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory          = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

type addressscopeActuator struct {
	osClient  osclients.AddressScopeClient
	k8sClient client.Client
}

var _ createResourceActuator = addressscopeActuator{}
var _ deleteResourceActuator = addressscopeActuator{}

func (addressscopeActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator addressscopeActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	resource, err := actuator.osClient.GetAddressScope(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return resource, nil
}

func (actuator addressscopeActuator) ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*osResourceT, error], bool) {
	resourceSpec := orcObject.Spec.Resource
	if resourceSpec == nil {
		return nil, false
	}

	// Resolve the project ID from ProjectRef if set.
	var projectID string
	if resourceSpec.ProjectRef != nil {
		project, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, orcObject.Namespace, resourceSpec.ProjectRef, "Project",
			func(dep *orcv1alpha1.Project) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := rs.NeedsReschedule(); needsReschedule {
			return nil, false
		}
		projectID = ptr.Deref(project.Status.ID, "")
	}

	listOpts := addressscopes.ListOpts{
		Name:      getResourceName(orcObject),
		IPVersion: int(resourceSpec.IPVersion),
		ProjectID: projectID,
	}

	return actuator.osClient.ListAddressScopes(ctx, listOpts), true
}

func (actuator addressscopeActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	var reconcileStatus progress.ReconcileStatus

	project, rs := dependency.FetchDependency[*orcv1alpha1.Project](
		ctx, actuator.k8sClient, obj.Namespace,
		filter.ProjectRef, "Project",
		orcv1alpha1.IsAvailable,
	)
	reconcileStatus = reconcileStatus.WithReconcileStatus(rs)

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	listOpts := addressscopes.ListOpts{
		Name:      string(ptr.Deref(filter.Name, "")),
		ProjectID: ptr.Deref(project.Status.ID, ""),
		IPVersion: int(filter.IPVersion),
		Shared:    filter.Shared,
	}

	return actuator.osClient.ListAddressScopes(ctx, listOpts), reconcileStatus
}

func (actuator addressscopeActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}
	var reconcileStatus progress.ReconcileStatus

	var projectID string
	if resource.ProjectRef != nil {
		project, projectDepRS := projectDependency.GetDependency(
			ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(projectDepRS)
		if project != nil {
			projectID = ptr.Deref(project.Status.ID, "")
		}
	}
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	createOpts := addressscopes.CreateOpts{
		Name:      getResourceName(obj),
		ProjectID: projectID,
		IPVersion: int(resource.IPVersion),
	}

	if resource.Shared != nil {
		createOpts.Shared = *resource.Shared
	}

	osResource, err := actuator.osClient.CreateAddressScope(ctx, createOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a conflict
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator addressscopeActuator) DeleteResource(ctx context.Context, _ orcObjectPT, resource *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteAddressScope(ctx, resource.ID))
}

func (actuator addressscopeActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := addressscopes.UpdateOpts{}

	handleNameUpdate(&updateOpts, obj, osResource)
	handleSharedUpdate(&updateOpts, resource, osResource)

	needsUpdate, err := needsUpdate(updateOpts)
	if err != nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err))
	}
	if !needsUpdate {
		log.V(logging.Debug).Info("No changes")
		return nil
	}

	_, err = actuator.osClient.UpdateAddressScope(ctx, osResource.ID, updateOpts)

	if err != nil {
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
		}
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts addressscopes.UpdateOpts) (bool, error) {
	updateOptsMap, err := updateOpts.ToAddressScopeUpdateMap()
	if err != nil {
		return false, err
	}

	updateMap, ok := updateOptsMap["address_scope"].(map[string]any)
	if !ok {
		updateMap = make(map[string]any)
	}

	return len(updateMap) > 0, nil
}

func handleNameUpdate(updateOpts *addressscopes.UpdateOpts, obj orcObjectPT, osResource *osResourceT) {
	name := getResourceName(obj)
	if osResource.Name != name {
		updateOpts.Name = &name
	}
}

func handleSharedUpdate(updateOpts *addressscopes.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	shared := ptr.Deref(resource.Shared, false)
	if shared != osResource.Shared {
		updateOpts.Shared = &shared
	}
}

func (actuator addressscopeActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		actuator.updateResource,
	}, nil
}

type addressscopeHelperFactory struct{}

var _ helperFactory = addressscopeHelperFactory{}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.AddressScope, controller interfaces.ResourceController) (addressscopeActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return addressscopeActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return addressscopeActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewAddressScopeClient()
	if err != nil {
		return addressscopeActuator{}, progress.WrapError(err)
	}

	return addressscopeActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}

func (addressscopeHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return addressscopeAdapter{obj}
}

func (addressscopeHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (addressscopeHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}
