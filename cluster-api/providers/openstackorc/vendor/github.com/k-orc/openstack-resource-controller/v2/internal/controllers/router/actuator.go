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

package router

import (
	"context"
	"fmt"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	osclients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/tags"
)

type (
	osResourceT = routers.Router

	createResourceActuator    = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator    = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	reconcileResourceActuator = interfaces.ReconcileResourceActuator[orcObjectPT, osResourceT]
	resourceReconciler        = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory             = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
	routerIterator            = iter.Seq2[*osResourceT, error]
)

type routerActuator struct {
	osClient osclients.NetworkClient
}

type routerCreateActuator struct {
	routerActuator
	k8sClient client.Client
}

var _ createResourceActuator = routerCreateActuator{}
var _ deleteResourceActuator = routerActuator{}

func (routerActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator routerActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	router, err := actuator.osClient.GetRouter(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return router, nil
}

func (actuator routerActuator) ListOSResourcesForAdoption(ctx context.Context, obj *orcv1alpha1.Router) (routerIterator, bool) {
	if obj.Spec.Resource == nil {
		return nil, false
	}

	listOpts := routers.ListOpts{Name: getResourceName(obj)}
	return actuator.osClient.ListRouter(ctx, listOpts), true
}

func (actuator routerCreateActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
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

	listOpts := routers.ListOpts{
		Name:        string(ptr.Deref(filter.Name, "")),
		Description: string(ptr.Deref(filter.Description, "")),
		ProjectID:   ptr.Deref(project.Status.ID, ""),
		Tags:        tags.Join(filter.Tags),
		TagsAny:     tags.Join(filter.TagsAny),
		NotTags:     tags.Join(filter.NotTags),
		NotTagsAny:  tags.Join(filter.NotTagsAny),
	}

	return actuator.osClient.ListRouter(ctx, listOpts), nil
}

func (actuator routerCreateActuator) CreateResource(ctx context.Context, obj *orcv1alpha1.Router) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}

	var reconcileStatus progress.ReconcileStatus

	gatewayInfo := &routers.GatewayInfo{}
	if len(resource.ExternalGateways) > 0 {
		var externalGW *orcv1alpha1.Network
		// Fetch dependencies and ensure they have our finalizer
		externalGW, reconcileStatus = externalGWDep.GetDependency(
			ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.Network) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if externalGW != nil {
			gatewayInfo.NetworkID = ptr.Deref(externalGW.Status.ID, "")
		}
	}

	var projectID string
	if resource.ProjectRef != nil {
		project, projectDepRS := projectDependency.GetDependency(
			ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.Project) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(projectDepRS)
		if project != nil {
			projectID = ptr.Deref(project.Status.ID, "")
		}
	}

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	createOpts := routers.CreateOpts{
		Name:         getResourceName(obj),
		Description:  string(ptr.Deref(resource.Description, "")),
		AdminStateUp: resource.AdminStateUp,
		Distributed:  resource.Distributed,
		GatewayInfo:  gatewayInfo,
		ProjectID:    projectID,
	}

	if len(resource.AvailabilityZoneHints) > 0 {
		createOpts.AvailabilityZoneHints = make([]string, len(resource.AvailabilityZoneHints))
		for i := range resource.AvailabilityZoneHints {
			createOpts.AvailabilityZoneHints[i] = string(resource.AvailabilityZoneHints[i])
		}
	}

	osResource, err := actuator.osClient.CreateRouter(ctx, &createOpts)

	// We should require the spec to be updated before retrying a create which returned a conflict
	if orcerrors.IsConflict(err) {
		err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
	}

	if err != nil {
		return nil, progress.WrapError(err)
	}
	return osResource, nil
}

func (actuator routerActuator) DeleteResource(ctx context.Context, _ orcObjectPT, router *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteRouter(ctx, router.ID))
}

func (actuator routerActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := &routers.UpdateOpts{}

	handleNameUpdate(updateOpts, obj, osResource)
	handleDescriptionUpdate(updateOpts, resource, osResource)
	handleAdminStateUpUpdate(updateOpts, resource, osResource)

	needsUpdate, err := needsUpdate(updateOpts)
	if err != nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err))
	}
	if !needsUpdate {
		log.V(logging.Debug).Info("No changes")
		return nil
	}

	updateOpts.RevisionNumber = &osResource.RevisionNumber

	_, err = actuator.osClient.UpdateRouter(ctx, osResource.ID, updateOpts)

	if orcerrors.IsConflict(err) {
		err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
	}
	if err != nil {
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts routers.UpdateOptsBuilder) (bool, error) {
	updateOptsMap, err := updateOpts.ToRouterUpdateMap()
	if err != nil {
		return false, err
	}

	routerUpdateMap, ok := updateOptsMap["router"].(map[string]any)
	if !ok {
		routerUpdateMap = make(map[string]any)
	}

	return len(routerUpdateMap) > 0, nil
}

func handleNameUpdate(updateOpts *routers.UpdateOpts, obj orcObjectPT, osResource *osResourceT) {
	name := getResourceName(obj)
	if osResource.Name != name {
		updateOpts.Name = name
	}
}

func handleDescriptionUpdate(updateOpts *routers.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	description := string(ptr.Deref(resource.Description, ""))
	if osResource.Description != description {
		updateOpts.Description = &description
	}
}

func handleAdminStateUpUpdate(updateOpts *routers.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	// Default is true
	AdminStateUp := ptr.Deref(resource.AdminStateUp, true)
	if osResource.AdminStateUp != AdminStateUp {
		updateOpts.AdminStateUp = &AdminStateUp
	}
}

var _ reconcileResourceActuator = routerActuator{}

func (actuator routerActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		tags.ReconcileTags[orcObjectPT, osResourceT](orcObject.Spec.Resource.Tags, osResource.Tags, tags.NewNeutronTagReplacer(actuator.osClient, "routers", osResource.ID)),
		actuator.updateResource,
	}, nil
}

type routerHelperFactory struct{}

var _ helperFactory = routerHelperFactory{}

func (routerHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return routerAdapter{obj}
}

func (routerHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newCreateActuator(ctx, orcObject, controller)
}

func (routerHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.Router, controller interfaces.ResourceController) (routerActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return routerActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return routerActuator{}, nil
	}
	osClient, err := clientScope.NewNetworkClient()
	if err != nil {
		return routerActuator{}, nil
	}

	return routerActuator{
		osClient: osClient,
	}, nil
}

func newCreateActuator(ctx context.Context, orcObject *orcv1alpha1.Router, controller interfaces.ResourceController) (routerCreateActuator, progress.ReconcileStatus) {
	routerActuator, reconcileStatus := newActuator(ctx, orcObject, controller)
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return routerCreateActuator{}, reconcileStatus
	}

	return routerCreateActuator{
		routerActuator: routerActuator,
		k8sClient:      controller.GetK8sClient(),
	}, nil
}
