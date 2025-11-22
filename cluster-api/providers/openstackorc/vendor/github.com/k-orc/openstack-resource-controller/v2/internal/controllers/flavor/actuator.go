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

package flavor

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	generic "github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	osclients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
)

// OpenStack resource types
type (
	osResourceT = flavors.Flavor

	createResourceActuator = generic.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator = generic.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	helperFactory          = generic.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

type flavorClient interface {
	GetFlavor(context.Context, string) (*flavors.Flavor, error)
	ListFlavors(context.Context, flavors.ListOptsBuilder) iter.Seq2[*flavors.Flavor, error]
	CreateFlavor(context.Context, flavors.CreateOptsBuilder) (*flavors.Flavor, error)
	DeleteFlavor(context.Context, string) error
}

type flavorActuator struct {
	osClient flavorClient
}

var _ createResourceActuator = flavorActuator{}
var _ deleteResourceActuator = flavorActuator{}

func (flavorActuator) GetResourceID(osResource *flavors.Flavor) string {
	return osResource.ID
}

func (actuator flavorActuator) GetOSResourceByID(ctx context.Context, id string) (*flavors.Flavor, progress.ReconcileStatus) {
	flavor, err := actuator.osClient.GetFlavor(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return flavor, nil
}

func (actuator flavorActuator) ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*flavors.Flavor, error], bool) {
	resourceSpec := orcObject.Spec.Resource
	if resourceSpec == nil {
		return nil, false
	}

	var filters []osclients.ResourceFilter[osResourceT]
	listOpts := flavors.ListOpts{}

	filters = append(filters,
		func(f *flavors.Flavor) bool {
			name := getResourceName(orcObject)
			// Compare non-pointer values
			return f.Name == name &&
				f.RAM == int(resourceSpec.RAM) &&
				f.VCPUs == int(resourceSpec.Vcpus) &&
				f.Disk == int(resourceSpec.Disk) &&
				f.Swap == int(resourceSpec.Swap) &&
				f.Ephemeral == int(resourceSpec.Ephemeral)
		},
	)

	if resourceSpec.Description != nil {
		filters = append(filters, func(f *flavors.Flavor) bool {
			return f.Description == *resourceSpec.Description
		})
	}

	// We can select on isPublic server-side, so add it to listOpts
	if resourceSpec.IsPublic != nil {
		if *resourceSpec.IsPublic {
			listOpts.AccessType = flavors.PublicAccess
		} else {
			listOpts.AccessType = flavors.PrivateAccess
		}
	}

	return actuator.listOSResources(ctx, filters, &listOpts), true
}

func (actuator flavorActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	var filters []osclients.ResourceFilter[osResourceT]

	if filter.Name != nil {
		filters = append(filters, func(f *flavors.Flavor) bool { return f.Name == string(*filter.Name) })
	}

	if filter.RAM != nil {
		filters = append(filters, func(f *flavors.Flavor) bool { return f.RAM == int(*filter.RAM) })
	}

	if filter.Vcpus != nil {
		filters = append(filters, func(f *flavors.Flavor) bool { return f.VCPUs == int(*filter.Vcpus) })
	}

	if filter.Disk != nil {
		filters = append(filters, func(f *flavors.Flavor) bool { return f.Disk == int(*filter.Disk) })
	}

	return actuator.listOSResources(ctx, filters, &flavors.ListOpts{}), nil
}

func (actuator flavorActuator) listOSResources(ctx context.Context, filters []osclients.ResourceFilter[osResourceT], listOpts flavors.ListOptsBuilder) iter.Seq2[*flavors.Flavor, error] {
	flavors := actuator.osClient.ListFlavors(ctx, listOpts)
	return osclients.Filter(flavors, filters...)
}

func (actuator flavorActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*flavors.Flavor, progress.ReconcileStatus) {
	resource := obj.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}

	createOpts := flavors.CreateOpts{
		Name:        getResourceName(obj),
		RAM:         int(resource.RAM),
		VCPUs:       int(resource.Vcpus),
		Disk:        ptr.To(int(resource.Disk)),
		Swap:        ptr.To(int(resource.Swap)),
		IsPublic:    resource.IsPublic,
		Ephemeral:   ptr.To(int(resource.Ephemeral)),
		Description: ptr.Deref(resource.Description, ""),
	}

	osResource, err := actuator.osClient.CreateFlavor(ctx, createOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a conflict
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator flavorActuator) DeleteResource(ctx context.Context, _ orcObjectPT, flavor *flavors.Flavor) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteFlavor(ctx, flavor.ID))
}

type flavorHelperFactory struct{}

var _ helperFactory = flavorHelperFactory{}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.Flavor, controller generic.ResourceController) (flavorActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return flavorActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return flavorActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewComputeClient()
	if err != nil {
		return flavorActuator{}, progress.WrapError(err)
	}

	return flavorActuator{
		osClient: osClient,
	}, nil
}

func (flavorHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return flavorAdapter{obj}
}

func (flavorHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller generic.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (flavorHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller generic.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}
