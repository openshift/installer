/*
Copyright 2025 The ORC Authors.

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

package keypair

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	osclients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
)

// OpenStack resource types
type (
	osResourceT = keypairs.KeyPair

	createResourceActuator = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	resourceReconciler     = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory          = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

type KeyPairClient interface {
	GetKeyPair(context.Context, string) (*osResourceT, error)
	ListKeyPairs(context.Context, keypairs.ListOptsBuilder) iter.Seq2[*osResourceT, error]
	CreateKeyPair(context.Context, keypairs.CreateOptsBuilder) (*osResourceT, error)
	DeleteKeyPair(context.Context, string) error
}

type keypairActuator struct {
	osClient KeyPairClient
}

var _ createResourceActuator = keypairActuator{}
var _ deleteResourceActuator = keypairActuator{}

func (keypairActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.Name
}

func (actuator keypairActuator) GetOSResourceByID(ctx context.Context, name string) (*osResourceT, progress.ReconcileStatus) {
	// For Keypairs, ID is the name
	resource, err := actuator.osClient.GetKeyPair(ctx, name)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return resource, nil
}

func (actuator keypairActuator) ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*osResourceT, error], bool) {
	resourceSpec := orcObject.Spec.Resource
	if resourceSpec == nil {
		return nil, false
	}

	// Filter by the expected resource name to avoid adopting wrong keypairs.
	// The OpenStack Keypairs API does not support server-side filtering by name,
	// so we must use client-side filtering.
	filters := []osclients.ResourceFilter[osResourceT]{
		func(kp *keypairs.KeyPair) bool {
			return kp.Name == getResourceName(orcObject)
		},
	}

	return actuator.listOSResources(ctx, filters, keypairs.ListOpts{}), true
}

func (actuator keypairActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	// The OpenStack Keypairs API does not support server-side filtering,
	// so client-side filtering is required for all fields.
	var filters []osclients.ResourceFilter[osResourceT]

	if filter.Name != nil {
		filters = append(filters, func(kp *keypairs.KeyPair) bool {
			return kp.Name == string(*filter.Name)
		})
	}

	return actuator.listOSResources(ctx, filters, keypairs.ListOpts{}), nil
}

func (actuator keypairActuator) listOSResources(ctx context.Context, filters []osclients.ResourceFilter[osResourceT], listOpts keypairs.ListOptsBuilder) iter.Seq2[*osResourceT, error] {
	keypairs := actuator.osClient.ListKeyPairs(ctx, listOpts)
	return osclients.Filter(keypairs, filters...)
}

func (actuator keypairActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource

	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}

	createOpts := keypairs.CreateOpts{
		Name:      getResourceName(obj),
		Type:      ptr.Deref(resource.Type, ""),
		PublicKey: resource.PublicKey,
	}

	osResource, err := actuator.osClient.CreateKeyPair(ctx, createOpts)
	if err != nil {
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration,
				"invalid configuration creating Keypair: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator keypairActuator) DeleteResource(ctx context.Context, _ orcObjectPT, Keypair *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteKeyPair(ctx, Keypair.Name))
}

func (actuator keypairActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	// Keypairs are immutable - no update reconcilers needed
	return []resourceReconciler{}, nil
}

type keypairHelperFactory struct{}

var _ helperFactory = keypairHelperFactory{}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.KeyPair, controller interfaces.ResourceController) (keypairActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return keypairActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return keypairActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewKeyPairClient()
	if err != nil {
		return keypairActuator{}, progress.WrapError(err)
	}

	return keypairActuator{
		osClient: osClient,
	}, nil
}

func (keypairHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return keypairAdapter{obj}
}

func (keypairHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func (keypairHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}
