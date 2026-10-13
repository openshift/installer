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

package network

import (
	"context"
	"fmt"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/mtu"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	"github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/tags"
)

type (
	osResourceT = osclients.NetworkExt

	createResourceActuator    = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator    = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	reconcileResourceActuator = interfaces.ReconcileResourceActuator[orcObjectPT, osResourceT]
	resourceReconciler        = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory             = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

type networkActuator struct {
	osClient  osclients.NetworkClient
	k8sClient client.Client
}

var _ createResourceActuator = networkActuator{}
var _ deleteResourceActuator = networkActuator{}
var _ reconcileResourceActuator = networkActuator{}

func (networkActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator networkActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	network, err := actuator.osClient.GetNetwork(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return network, nil
}

func (actuator networkActuator) ListOSResourcesForAdoption(ctx context.Context, obj orcObjectPT) (iter.Seq2[*osResourceT, error], bool) {
	if obj.Spec.Resource == nil {
		return nil, false
	}

	listOpts := networks.ListOpts{Name: getResourceName(obj)}
	return actuator.osClient.ListNetwork(ctx, listOpts), true
}

func (actuator networkActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
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

	listOpts := networks.ListOpts{
		Name:        string(ptr.Deref(filter.Name, "")),
		Description: string(ptr.Deref(filter.Description, "")),
		ProjectID:   ptr.Deref(project.Status.ID, ""),
		Tags:        tags.Join(filter.Tags),
		TagsAny:     tags.Join(filter.TagsAny),
		NotTags:     tags.Join(filter.NotTags),
		NotTagsAny:  tags.Join(filter.NotTagsAny),
	}

	return actuator.osClient.ListNetwork(ctx, listOpts), nil
}

func (actuator networkActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}

	var projectID string
	if resource.ProjectRef != nil {
		project, reconcileStatus := projectDependency.GetDependency(
			ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.Project) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
			return nil, reconcileStatus
		}
		projectID = ptr.Deref(project.Status.ID, "")
	}

	var createOpts networks.CreateOptsBuilder
	{
		createOptsBase := networks.CreateOpts{
			Name:         getResourceName(obj),
			Description:  string(ptr.Deref(resource.Description, "")),
			AdminStateUp: resource.AdminStateUp,
			Shared:       resource.Shared,
			ProjectID:    projectID,
		}

		if len(resource.AvailabilityZoneHints) > 0 {
			createOptsBase.AvailabilityZoneHints = make([]string, len(resource.AvailabilityZoneHints))
			for i := range resource.AvailabilityZoneHints {
				createOptsBase.AvailabilityZoneHints[i] = string(resource.AvailabilityZoneHints[i])
			}
		}
		createOpts = createOptsBase
	}

	if resource.DNSDomain != nil {
		createOpts = &dns.NetworkCreateOptsExt{
			CreateOptsBuilder: createOpts,
			DNSDomain:         string(*resource.DNSDomain),
		}
	}

	if resource.MTU != nil {
		createOpts = &mtu.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			MTU:               int(*resource.MTU),
		}
	}

	if resource.PortSecurityEnabled != nil {
		createOpts = &portsecurity.NetworkCreateOptsExt{
			CreateOptsBuilder:   createOpts,
			PortSecurityEnabled: resource.PortSecurityEnabled,
		}
	}

	if resource.External != nil {
		createOpts = &external.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			External:          resource.External,
		}
	}

	osResource, err := actuator.osClient.CreateNetwork(ctx, createOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a conflict
		if orcerrors.IsConflict(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator networkActuator) DeleteResource(ctx context.Context, _ orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteNetwork(ctx, osResource.ID))
}

func (actuator networkActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		tags.ReconcileTags[orcObjectPT, osResourceT](orcObject.Spec.Resource.Tags, osResource.Tags, tags.NewNeutronTagReplacer(actuator.osClient, "networks", osResource.ID)),
		actuator.updateResource,
	}, nil
}

func (actuator networkActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	var updateOpts networks.UpdateOptsBuilder
	{
		updateOptsBase := networks.UpdateOpts{
			RevisionNumber: &osResource.RevisionNumber,
		}
		handleAdminStateUpUpdate(&updateOptsBase, resource, osResource)
		handleNameUpdate(&updateOptsBase, obj, osResource)
		handleDescriptionUpdate(&updateOptsBase, resource, osResource)
		handleSharedUpdate(&updateOptsBase, resource, osResource)
		updateOpts = updateOptsBase
	}

	updateOpts = handlePortSecurityEnabledUpdate(updateOpts, resource, osResource)
	updateOpts = handleMTUUpdate(updateOpts, resource, osResource)
	updateOpts = handleExternalUpdate(updateOpts, resource, osResource)

	needsUpdate, err := needsUpdate(updateOpts)
	if err != nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err))
	}
	if !needsUpdate {
		log.V(logging.Debug).Info("No changes")
		return nil
	}

	_, err = actuator.osClient.UpdateNetwork(ctx, osResource.ID, updateOpts)

	if orcerrors.IsConflict(err) {
		err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
	}
	if err != nil {
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts networks.UpdateOptsBuilder) (bool, error) {
	updateOptsMap, err := updateOpts.ToNetworkUpdateMap()
	if err != nil {
		return false, err
	}

	networkUpdateMap, ok := updateOptsMap["network"].(map[string]any)
	if !ok {
		networkUpdateMap = make(map[string]any)
	}

	// Revision number is not returned in the output of updateOpts.ToNetworkUpdateMap()
	// so nothing to drop here

	return len(networkUpdateMap) > 0, nil
}

func handleAdminStateUpUpdate(updateOpts *networks.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	// Default is true
	AdminStateUp := ptr.Deref(resource.AdminStateUp, true)
	if osResource.AdminStateUp != AdminStateUp {
		updateOpts.AdminStateUp = &AdminStateUp
	}
}

func handleNameUpdate(updateOpts *networks.UpdateOpts, obj orcObjectPT, osResource *osResourceT) {
	name := getResourceName(obj)
	if osResource.Name != name {
		updateOpts.Name = &name
	}
}

func handleDescriptionUpdate(updateOpts *networks.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	description := string(ptr.Deref(resource.Description, ""))
	if osResource.Description != description {
		updateOpts.Description = &description
	}
}

func handleSharedUpdate(updateOpts *networks.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	// Default is false
	Shared := ptr.Deref(resource.Shared, false)
	if osResource.Shared != Shared {
		updateOpts.Shared = &Shared
	}
}

func handlePortSecurityEnabledUpdate(updateOpts networks.UpdateOptsBuilder, resource *resourceSpecT, osResource *osResourceT) networks.UpdateOptsBuilder {
	if resource.PortSecurityEnabled != nil {
		if *resource.PortSecurityEnabled != osResource.PortSecurityEnabled {
			updateOpts = &portsecurity.NetworkUpdateOptsExt{
				UpdateOptsBuilder:   updateOpts,
				PortSecurityEnabled: resource.PortSecurityEnabled,
			}
		}
	}
	return updateOpts
}

func handleMTUUpdate(updateOpts networks.UpdateOptsBuilder, resource *resourceSpecT, osResource *osResourceT) networks.UpdateOptsBuilder {
	if resource.MTU != nil && int(*resource.MTU) != osResource.MTU {
		updateOpts = &mtu.UpdateOptsExt{
			UpdateOptsBuilder: updateOpts,
			MTU:               int(*resource.MTU),
		}
	}
	return updateOpts
}

func handleExternalUpdate(updateOpts networks.UpdateOptsBuilder, resource *resourceSpecT, osResource *osResourceT) networks.UpdateOptsBuilder {
	if resource.External != nil && *resource.External != osResource.External {
		updateOpts = &external.UpdateOptsExt{
			UpdateOptsBuilder: updateOpts,
			External:          resource.External,
		}
	}
	return updateOpts
}

type networkHelperFactory struct{}

var _ helperFactory = networkHelperFactory{}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.Network, controller interfaces.ResourceController) (networkActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return networkActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return networkActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewNetworkClient()
	if err != nil {
		return networkActuator{}, progress.WrapError(err)
	}

	return networkActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}

func (networkHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return networkAdapter{obj}
}

func (networkHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (networkHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}
