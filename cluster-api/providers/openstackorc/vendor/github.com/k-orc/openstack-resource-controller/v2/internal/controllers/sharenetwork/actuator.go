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

package sharenetwork

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharenetworks"
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
	osResourceT = sharenetworks.ShareNetwork

	createResourceActuator = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	resourceReconciler     = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory          = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

type sharenetworkActuator struct {
	osClient  osclients.ShareNetworkClient
	k8sClient client.Client
}

var _ createResourceActuator = sharenetworkActuator{}
var _ deleteResourceActuator = sharenetworkActuator{}

func (sharenetworkActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator sharenetworkActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	resource, err := actuator.osClient.GetShareNetwork(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return resource, nil
}

func (actuator sharenetworkActuator) ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*osResourceT, error], bool) {
	resourceSpec := orcObject.Spec.Resource
	if resourceSpec == nil {
		return nil, false
	}

	listOpts := sharenetworks.ListOpts{
		Name: getResourceName(orcObject),
	}

	return actuator.osClient.ListShareNetworks(ctx, listOpts), true
}

func (actuator sharenetworkActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	listOpts := sharenetworks.ListOpts{
		Name:        string(ptr.Deref(filter.Name, "")),
		Description: ptr.Deref(filter.Description, ""),
	}

	return actuator.osClient.ListShareNetworks(ctx, listOpts), nil
}

func (actuator sharenetworkActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}
	var reconcileStatus progress.ReconcileStatus

	var networkID string
	if resource.NetworkRef != nil {
		network, networkDepRS := networkDependency.GetDependency(
			ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(networkDepRS)
		if network != nil {
			networkID = ptr.Deref(network.Status.ID, "")
		}
	}

	var subnetID string
	if resource.SubnetRef != nil {
		subnet, subnetDepRS := subnetDependency.GetDependency(
			ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(subnetDepRS)
		if subnet != nil {
			subnetID = ptr.Deref(subnet.Status.ID, "")
		}
	}
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}
	createOpts := sharenetworks.CreateOpts{
		Name:            getResourceName(obj),
		Description:     ptr.Deref(resource.Description, ""),
		NeutronNetID:    networkID,
		NeutronSubnetID: subnetID,
	}

	osResource, err := actuator.osClient.CreateShareNetwork(ctx, createOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a conflict
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator sharenetworkActuator) DeleteResource(ctx context.Context, _ orcObjectPT, resource *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteShareNetwork(ctx, resource.ID))
}

func (actuator sharenetworkActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := sharenetworks.UpdateOpts{}

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

	_, err = actuator.osClient.UpdateShareNetwork(ctx, osResource.ID, updateOpts)
	if err != nil {
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
		}
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts sharenetworks.UpdateOptsBuilder) (bool, error) {
	updateOptsMap, err := updateOpts.ToShareNetworkUpdateMap()
	if err != nil {
		return false, err
	}

	updateMap, ok := updateOptsMap["share_network"].(map[string]any)
	if !ok {
		updateMap = make(map[string]any)
	}

	return len(updateMap) > 0, nil
}

func handleNameUpdate(updateOpts *sharenetworks.UpdateOpts, obj orcObjectPT, osResource *osResourceT) {
	name := getResourceName(obj)
	if osResource.Name != name {
		updateOpts.Name = &name
	}
}

func handleDescriptionUpdate(updateOpts *sharenetworks.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	description := ptr.Deref(resource.Description, "")
	if osResource.Description != description {
		updateOpts.Description = &description
	}
}

func (actuator sharenetworkActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		actuator.updateResource,
	}, nil
}

type sharenetworkHelperFactory struct{}

var _ helperFactory = sharenetworkHelperFactory{}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.ShareNetwork, controller interfaces.ResourceController) (sharenetworkActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return sharenetworkActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return sharenetworkActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewShareNetworkClient()
	if err != nil {
		return sharenetworkActuator{}, progress.WrapError(err)
	}

	return sharenetworkActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}

func (sharenetworkHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return sharenetworkAdapter{obj}
}

func (sharenetworkHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (sharenetworkHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}
