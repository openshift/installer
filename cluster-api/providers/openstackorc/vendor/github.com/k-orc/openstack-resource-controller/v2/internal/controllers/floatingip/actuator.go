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

package floatingip

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	osclients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/tags"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	osResourceT               = floatingips.FloatingIP
	createResourceActuator    = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator    = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	reconcileResourceActuator = interfaces.ReconcileResourceActuator[orcObjectPT, osResourceT]
	resourceReconciler        = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory             = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
	floatingipIterator        = iter.Seq2[*osResourceT, error]
)

type floatingipActuator struct {
	osClient  osclients.NetworkClient
	k8sClient client.Client
}

var _ createResourceActuator = floatingipActuator{}
var _ deleteResourceActuator = floatingipActuator{}

func (floatingipActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator floatingipActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	floatingip, err := actuator.osClient.GetFloatingIP(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return floatingip, nil
}

func (actuator floatingipActuator) ListOSResourcesForAdoption(ctx context.Context, obj *orcv1alpha1.FloatingIP) (floatingipIterator, bool) {
	resource := obj.Spec.Resource
	if resource == nil {
		return nil, false
	}
	// we only support adoption of floatingips by IP as they don't have name
	if resource.FloatingIP == nil {
		return nil, false
	}

	// Resolve the floating network ID from either FloatingNetworkRef or
	// FloatingSubnetRef. Exactly one of these must be set per API
	// validation. Without the network ID, adoption could match a floating
	// IP on the wrong network.
	var floatingNetworkID string
	if resource.FloatingNetworkRef != nil {
		network, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, obj.Namespace, resource.FloatingNetworkRef, "Network",
			func(dep *orcv1alpha1.Network) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := rs.NeedsReschedule(); needsReschedule {
			return nil, false
		}
		floatingNetworkID = ptr.Deref(network.Status.ID, "")
	} else if resource.FloatingSubnetRef != nil {
		subnet, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, obj.Namespace, resource.FloatingSubnetRef, "Subnet",
			func(dep *orcv1alpha1.Subnet) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil && dep.Status.Resource != nil
			},
		)
		if needsReschedule, _ := rs.NeedsReschedule(); needsReschedule {
			return nil, false
		}
		floatingNetworkID = subnet.Status.Resource.NetworkID
	}

	// Resolve the project ID from ProjectRef if set.
	var projectID string
	if resource.ProjectRef != nil {
		project, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, obj.Namespace, resource.ProjectRef, "Project",
			func(dep *orcv1alpha1.Project) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := rs.NeedsReschedule(); needsReschedule {
			return nil, false
		}
		projectID = ptr.Deref(project.Status.ID, "")
	}

	listOpts := floatingips.ListOpts{
		FloatingIP:        string(ptr.Deref(resource.FloatingIP, "")),
		FloatingNetworkID: floatingNetworkID,
		ProjectID:         projectID,
		Tags:              tags.Join(resource.Tags),
	}
	return actuator.osClient.ListFloatingIP(ctx, listOpts), true
}

func (actuator floatingipActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	var reconcileStatus progress.ReconcileStatus

	network, rs := dependency.FetchDependency[*orcv1alpha1.Network](
		ctx, actuator.k8sClient, obj.Namespace, filter.FloatingNetworkRef, "Network",
		orcv1alpha1.IsAvailable,
	)
	reconcileStatus = reconcileStatus.WithReconcileStatus(rs)

	port, rs := dependency.FetchDependency[*orcv1alpha1.Port](
		ctx, actuator.k8sClient, obj.Namespace, filter.PortRef, "Port",
		orcv1alpha1.IsAvailable,
	)
	reconcileStatus = reconcileStatus.WithReconcileStatus(rs)

	project, rs := dependency.FetchDependency[*orcv1alpha1.Project](
		ctx, actuator.k8sClient, obj.Namespace, filter.ProjectRef, "Project",
		orcv1alpha1.IsAvailable,
	)
	reconcileStatus = reconcileStatus.WithReconcileStatus(rs)

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	listOpts := floatingips.ListOpts{
		FloatingIP:        string(ptr.Deref(filter.FloatingIP, "")),
		PortID:            ptr.Deref(port.Status.ID, ""),
		FloatingNetworkID: ptr.Deref(network.Status.ID, ""),
		ProjectID:         ptr.Deref(project.Status.ID, ""),
		Description:       string(ptr.Deref(filter.Description, "")),
		Tags:              tags.Join(filter.Tags),
		TagsAny:           tags.Join(filter.TagsAny),
		NotTags:           tags.Join(filter.NotTags),
		NotTagsAny:        tags.Join(filter.NotTagsAny),
		Status:            filter.Status,
	}

	return actuator.osClient.ListFloatingIP(ctx, listOpts), nil
}

func (actuator floatingipActuator) CreateResource(ctx context.Context, obj *orcv1alpha1.FloatingIP) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}
	var reconcileStatus progress.ReconcileStatus

	var networkID string
	if resource.FloatingNetworkRef != nil {
		// Fetch dependencies and ensure they have our finalizer
		network, networkDepRS := networkDep.GetDependency(
			ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(networkDepRS)
		if network != nil {
			networkID = ptr.Deref(network.Status.ID, "")
		}
	}

	var subnetID string
	// If we have a subnet (i.e. we don't have FloatingNetworkRef), we need to fetch it to get its ID and the network ID (as it's required by gophercloud)
	if resource.FloatingSubnetRef != nil {
		// Fetch dependencies and ensure they have our finalizer
		subnet, subnetDepRS := subnetDep.GetDependency(
			ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(subnetDepRS)
		if subnet != nil {
			subnetID = ptr.Deref(subnet.Status.ID, "")
			networkID = subnet.Status.Resource.NetworkID
		}
	}

	var portID string
	if resource.PortRef != nil {
		// Fetch dependencies and ensure they have our finalizer
		port, portDepRS := portDep.GetDependency(
			ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(portDepRS)
		if port != nil {
			portID = ptr.Deref(port.Status.ID, "")
		}
	}

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

	createOpts := floatingips.CreateOpts{
		Description:       string(ptr.Deref(resource.Description, "")),
		FloatingNetworkID: networkID,
		SubnetID:          subnetID,
		PortID:            portID,
		FloatingIP:        string(ptr.Deref(resource.FloatingIP, "")),
		FixedIP:           string(ptr.Deref(resource.FixedIP, "")),
		ProjectID:         projectID,
	}

	osResource, err := actuator.osClient.CreateFloatingIP(ctx, &createOpts)

	if err != nil {
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}
	return osResource, nil
}

func (actuator floatingipActuator) DeleteResource(ctx context.Context, _ orcObjectPT, floatingip *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteFloatingIP(ctx, floatingip.ID))
}

func (actuator floatingipActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := &floatingips.UpdateOpts{}

	handleDescriptionUpdate(updateOpts, resource, osResource)

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

	_, err = actuator.osClient.UpdateFloatingIP(ctx, osResource.ID, updateOpts)

	if err != nil {
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
		}
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts floatingips.UpdateOptsBuilder) (bool, error) {
	updateOptsMap, err := updateOpts.ToFloatingIPUpdateMap()
	if err != nil {
		return false, err
	}

	floatingIPUpdateMap, ok := updateOptsMap["floatingip"].(map[string]any)
	if !ok {
		floatingIPUpdateMap = make(map[string]any)
	}

	return len(floatingIPUpdateMap) > 0, nil
}

func handleDescriptionUpdate(updateOpts *floatingips.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	description := string(ptr.Deref(resource.Description, ""))
	if osResource.Description != description {
		updateOpts.Description = &description
	}
}

var _ reconcileResourceActuator = floatingipActuator{}

func (actuator floatingipActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		tags.ReconcileTags[orcObjectPT, osResourceT](orcObject.Spec.Resource.Tags, osResource.Tags, tags.NewNeutronTagReplacer(actuator.osClient, "floatingips", osResource.ID)),
		actuator.updateResource,
	}, nil
}

type floatingipHelperFactory struct{}

var _ helperFactory = floatingipHelperFactory{}

func (floatingipHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return floatingipAdapter{obj}
}

func (floatingipHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (floatingipHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.FloatingIP, controller interfaces.ResourceController) (floatingipActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return floatingipActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return floatingipActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewNetworkClient()
	if err != nil {
		return floatingipActuator{}, progress.WrapError(err)
	}

	return floatingipActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}
