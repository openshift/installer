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

package servergroup

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	generic "github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	osclients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
)

// OpenStack resource types
type (
	osResourceT = servergroups.ServerGroup

	createResourceActuator = generic.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator = generic.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	helperFactory          = generic.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

type servergroupClient interface {
	GetServerGroup(context.Context, string) (*servergroups.ServerGroup, error)
	ListServerGroups(context.Context, servergroups.ListOptsBuilder) iter.Seq2[*servergroups.ServerGroup, error]
	CreateServerGroup(context.Context, servergroups.CreateOptsBuilder) (*servergroups.ServerGroup, error)
	DeleteServerGroup(context.Context, string) error
}

type servergroupActuator struct {
	osClient servergroupClient
}

var _ createResourceActuator = servergroupActuator{}
var _ deleteResourceActuator = servergroupActuator{}

func (servergroupActuator) GetResourceID(osResource *servergroups.ServerGroup) string {
	return osResource.ID
}

func (actuator servergroupActuator) GetOSResourceByID(ctx context.Context, id string) (*servergroups.ServerGroup, progress.ReconcileStatus) {
	servergroup, err := actuator.osClient.GetServerGroup(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return servergroup, nil
}

func (actuator servergroupActuator) ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*servergroups.ServerGroup, error], bool) {
	resourceSpec := orcObject.Spec.Resource
	if resourceSpec == nil {
		return nil, false
	}

	var filters []osclients.ResourceFilter[osResourceT]
	listOpts := servergroups.ListOpts{}

	filters = append(filters,
		func(f *servergroups.ServerGroup) bool {
			name := getResourceName(orcObject)
			return f.Name == name
		},
	)

	return actuator.listOSResources(ctx, filters, &listOpts), true
}

func (actuator servergroupActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	var filters []osclients.ResourceFilter[osResourceT]

	if filter.Name != nil {
		filters = append(filters, func(f *servergroups.ServerGroup) bool { return f.Name == string(*filter.Name) })
	}

	return actuator.listOSResources(ctx, filters, &servergroups.ListOpts{}), nil
}

func (actuator servergroupActuator) listOSResources(ctx context.Context, filters []osclients.ResourceFilter[osResourceT], listOpts servergroups.ListOptsBuilder) iter.Seq2[*servergroups.ServerGroup, error] {
	servergroups := actuator.osClient.ListServerGroups(ctx, listOpts)
	return osclients.Filter(servergroups, filters...)
}

func (actuator servergroupActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*servergroups.ServerGroup, progress.ReconcileStatus) {
	resource := obj.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}

	var rules *servergroups.Rules
	if resource.Rules != nil {
		rules = &servergroups.Rules{
			MaxServerPerHost: int(resource.Rules.MaxServerPerHost),
		}
	}

	createOpts := servergroups.CreateOpts{
		Name:   getResourceName(obj),
		Policy: string(resource.Policy),
		Rules:  rules,
	}

	osResource, err := actuator.osClient.CreateServerGroup(ctx, createOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a conflict
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator servergroupActuator) DeleteResource(ctx context.Context, _ orcObjectPT, servergroup *servergroups.ServerGroup) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteServerGroup(ctx, servergroup.ID))
}

type servergroupHelperFactory struct{}

var _ helperFactory = servergroupHelperFactory{}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.ServerGroup, controller generic.ResourceController) (servergroupActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return servergroupActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return servergroupActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewComputeClient()
	if err != nil {
		return servergroupActuator{}, progress.WrapError(err)
	}

	return servergroupActuator{
		osClient: osClient,
	}, nil
}

func (servergroupHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return servergroupAdapter{obj}
}

func (servergroupHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller generic.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (servergroupHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller generic.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}
