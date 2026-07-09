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

package endpoint

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/endpoints"
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
	osResourceT = endpoints.Endpoint

	createResourceActuator = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	resourceReconciler     = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory          = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

type endpointActuator struct {
	osClient  osclients.EndpointClient
	k8sClient client.Client
}

var _ createResourceActuator = endpointActuator{}
var _ deleteResourceActuator = endpointActuator{}

func (endpointActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator endpointActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	resource, err := actuator.osClient.GetEndpoint(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return resource, nil
}

func (actuator endpointActuator) ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*osResourceT, error], bool) {
	resourceSpec := orcObject.Spec.Resource
	if resourceSpec == nil {
		return nil, false
	}

	service, _ := serviceDependency.GetDependency(
		ctx, actuator.k8sClient, orcObject, orcv1alpha1.IsAvailable,
	)

	if service == nil {
		return nil, false
	}

	filters := []osclients.ResourceFilter[osResourceT]{
		func(e *endpoints.Endpoint) bool {
			return e.URL == resourceSpec.URL
		},
	}

	listOpts := endpoints.ListOpts{
		Availability: gophercloud.Availability(resourceSpec.Interface),
		ServiceID:    ptr.Deref(service.Status.ID, ""),
	}

	return actuator.listOsResources(ctx, listOpts, filters), true
}

func (actuator endpointActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	var reconcileStatus progress.ReconcileStatus

	service, rs := dependency.FetchDependency[*orcv1alpha1.Service](
		ctx, actuator.k8sClient, obj.Namespace,
		filter.ServiceRef, "Service",
		orcv1alpha1.IsAvailable,
	)
	reconcileStatus = reconcileStatus.WithReconcileStatus(rs)

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	var resourceFilters []osclients.ResourceFilter[osResourceT]
	if filter.URL != "" {
		resourceFilters = append(resourceFilters, func(e *endpoints.Endpoint) bool {
			return e.URL == filter.URL
		})
	}

	listOpts := endpoints.ListOpts{
		ServiceID:    ptr.Deref(service.Status.ID, ""),
		Availability: gophercloud.Availability(filter.Interface),
	}

	return actuator.listOsResources(ctx, listOpts, resourceFilters), nil
}

func (actuator endpointActuator) listOsResources(ctx context.Context, listOpts endpoints.ListOpts, filter []osclients.ResourceFilter[osResourceT]) iter.Seq2[*osResourceT, error] {
	endpoints := actuator.osClient.ListEndpoints(ctx, listOpts)
	return osclients.Filter(endpoints, filter...)
}

func (actuator endpointActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}
	var reconcileStatus progress.ReconcileStatus

	var serviceID string
	service, serviceDepRS := serviceDependency.GetDependency(
		ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
	)

	reconcileStatus = reconcileStatus.WithReconcileStatus(serviceDepRS)
	if service != nil {
		serviceID = ptr.Deref(service.Status.ID, "")
	}
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}
	createOpts := endpoints.CreateOpts{
		Availability: gophercloud.Availability(resource.Interface),
		Description:  ptr.Deref(resource.Description, ""),
		Enabled:      resource.Enabled,
		ServiceID:    serviceID,
		URL:          resource.URL,
	}

	osResource, err := actuator.osClient.CreateEndpoint(ctx, createOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a conflict
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator endpointActuator) DeleteResource(ctx context.Context, _ orcObjectPT, resource *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteEndpoint(ctx, resource.ID))
}

func (actuator endpointActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := endpoints.UpdateOpts{}

	handleEnabledUpdate(&updateOpts, resource, osResource)
	handleURLUpdate(&updateOpts, resource, osResource)
	handleInterfaceUpdate(&updateOpts, resource, osResource)

	needsUpdate, err := needsUpdate(updateOpts)
	if err != nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err))
	}
	if !needsUpdate {
		log.V(logging.Debug).Info("No changes")
		return nil
	}

	_, err = actuator.osClient.UpdateEndpoint(ctx, osResource.ID, updateOpts)

	if err != nil {
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
		}
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts endpoints.UpdateOpts) (bool, error) {
	updateOptsMap, err := updateOpts.ToEndpointUpdateMap()
	if err != nil {
		return false, err
	}

	updateMap, ok := updateOptsMap["endpoint"].(map[string]any)
	if !ok {
		updateMap = make(map[string]any)
	}

	return len(updateMap) > 0, nil
}

func handleURLUpdate(updateOpts *endpoints.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	url := resource.URL
	if osResource.URL != url {
		updateOpts.URL = url
	}
}

func handleInterfaceUpdate(updateOpts *endpoints.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	endpointInterface := gophercloud.Availability(resource.Interface)
	if osResource.Availability != endpointInterface {
		updateOpts.Availability = endpointInterface
	}
}

func handleEnabledUpdate(updateOpts *endpoints.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	enabled := resource.Enabled
	if enabled != nil && osResource.Enabled != *enabled {
		updateOpts.Enabled = enabled
	}
}

func (actuator endpointActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		actuator.updateResource,
	}, nil
}

type endpointHelperFactory struct{}

var _ helperFactory = endpointHelperFactory{}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.Endpoint, controller interfaces.ResourceController) (endpointActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return endpointActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return endpointActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewEndpointClient()
	if err != nil {
		return endpointActuator{}, progress.WrapError(err)
	}

	return endpointActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}

func (endpointHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return endpointAdapter{obj}
}

func (endpointHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (endpointHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}
