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

package interfaces

import (
	"context"
	"iter"

	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
)

// ResourceHelperFactory is an interface defining constructors for objects
// required by the generic controller.
type ResourceHelperFactory[
	orcObjectPT interface {
		*orcObjectT
		client.Object
		orcv1alpha1.ObjectWithConditions
	}, orcObjectT any,
	resourceSpecT any, filterT any,
	osResourceT any,
] interface {
	// NewAPIObjectAdapter returns an APIObjectAdapter wrapping orcObject
	NewAPIObjectAdapter(orcObject orcObjectPT) APIObjectAdapter[orcObjectPT, resourceSpecT, filterT]

	// NewCreateActuator returns a CreateResourceActuator for the given
	// orcObject. If it is not able to return an actuator, it MUST return either
	// one or more ProgressStatuses, or an error. If returning ProgressStatuses,
	// these MUST ensure that the object will be reconciled again at an
	// appropriate time.
	NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller ResourceController) (CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT], progress.ReconcileStatus)

	// NewDeleteActuator returns a DeleteResourceActuator for the given
	// orcObject. If it is not able to return an actuator, it MUST return either
	// one or more ProgressStatuses, or an error. If returning ProgressStatuses,
	// these MUST ensure that the object will be reconciled again at an
	// appropriate time.
	//
	// Consider carefully whether a DeleteResourceActuator needs all the same
	// initialisation dependencies as a CreateResourceActuator. Consider that we
	// may want to delete a resource that is partially or not initialised, or
	// whose creation dependencies may no longer be in a healthy state.
	NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller ResourceController) (DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT], progress.ReconcileStatus)
}

type BaseResourceActuator[
	orcObjectPT interface {
		*orcObjectT
		client.Object
		orcv1alpha1.ObjectWithConditions
	}, orcObjectT any,
	osResourceT any,
] interface {
	// GetResourceID returns a string identifier for the OpenStack resource
	// managed by this actuator, typically the object's uuid. This is the value
	// stored in status.id.
	GetResourceID(osResource *osResourceT) string

	// GetOSResourceByID fetches this actuator's OpenStack resource by id.
	GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus)

	// ListOSResourcesForAdoption is used to prevent resource leaks in the event
	// that we create an OpenStack resource, but fail to write its ID to the
	// object's status. It returns a set of resources which match what the
	// actuator would have created for this object. Ideally it should attempt to
	// do this match as accurately as possible, remembering that OpenStack does
	// not prevent the creation of objects with duplicate names.
	//
	// It is called in the creation flow immediately before creating a resource,
	// and in the deletion flow when deleting an object which has a finalizer
	// but no status.id.
	//
	// It returns 2 values:
	// - an iterator over the matching OpenStack resources
	// - a boolean indicating whether adoption should be considered for this object
	//
	// For example, we must return false for an object with no resource spec.
	ListOSResourcesForAdoption(ctx context.Context, orcObject orcObjectPT) (iter.Seq2[*osResourceT, error], bool)
}

// CreateResourceActuator provides methods required by the generic controller
// during the create and update flows.
type CreateResourceActuator[
	orcObjectPT interface {
		*orcObjectT
		client.Object
		orcv1alpha1.ObjectWithConditions
	}, orcObjectT any,
	filterT any,
	osResourceT any,
] interface {
	BaseResourceActuator[orcObjectPT, orcObjectT, osResourceT]

	// ListOSResourcesForImport returns all OpenStack resources matching the
	// given resource import filter.
	ListOSResourcesForImport(ctx context.Context, orcObject orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus)

	// CreateResource creates an OpenStack resource for the current object. It
	// will return successfully at most once for a managed object, at the time
	// it is created. See `ResourceReconciler` for how to modify an existing
	// object.
	//
	// CreateResource MUST NOT perform any action which can fail after creating
	// the primary OpenStack resource. Once the OpenStack resource exists this
	// method will not be called again, so any failure after creation of the
	// OpenStack resource is not safe.
	//
	// If further initialisation of a resource is required after creation, for
	// example a Neutron object whose tags cannot be set during creation, this
	// must be done by a reconciler instead.
	//
	// CreateResource MAY also perform other actions prior to creating the
	// primary OpenStack resource, e.g. creating dependent ORC objects. These
	// must be created idempotently, and must be created BEFORE creating the
	// OpenStack resource.
	//
	// CreateResource does not need to check if the resource already exists, as
	// that is done before it is called.
	//
	// If CreateResource cannot create the resource it MUST return either one or
	// more ProgressStatuses, or an error. If returning ProgressStatuses, these
	// MUST be sufficient to ensure that the object will be reconciled again at
	// an appropriate time.
	CreateResource(ctx context.Context, orcObject orcObjectPT) (*osResourceT, progress.ReconcileStatus)
}

// DeleteResourceActuator provides methods required by the generic controller
// during the delete flow.
type DeleteResourceActuator[
	orcObjectPT interface {
		*orcObjectT
		client.Object
		orcv1alpha1.ObjectWithConditions
	}, orcObjectT any,
	osResourceT any,
] interface {
	BaseResourceActuator[orcObjectPT, orcObjectT, osResourceT]

	// DeleteResource deletes the OpenStack resource owned by the current
	// object.
	//
	// The delete flow does not succeed until an attempt to get the resource
	// returns no results, so DeleteResource does not need to verify that
	// deletion has completed.
	//
	// DeleteResource SHOULD NOT return an error if the resource no longer exists.
	//
	// DeleteResource MAY also perform other actions prior to deleting the
	// OpenStack resource, e.g. deleting dependent ORC objects. These actions
	// must be performed idempotently, and must be performed BEFORE deleting the
	// OpenStack resource.
	//
	// If DeleteResource cannot delete the resource it MUST return either one or
	// more ProgressStatuses, or an error. If returning ProgressStatuses, these
	// MUST be sufficient to ensure that the objet will be reconciled again at
	// an appropriate time.
	DeleteResource(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT) progress.ReconcileStatus
}

// ResourceReconciler is a function which reconciles an object after creation.
//
// ResourceReconcilers are called on every non-delete reconciliation when the
// resource exists, including the reconciliation which created the resource.
//
// A ResourceReconciler may return one or more ProgressStatuses, and/or an
// error. Both errors and ProgressStatuses returned by ResourceReconcilers are
// aggregated before being passed to the ResourceStatusWriter.
//
// In addition to informing the Progressing condition in the object's status, a
// ProgressStatus returned by a ResourceReconciler may be used to cause the
// controller to poll, for example because the resource has not yet reached an
// ACTIVE status.
type ResourceReconciler[orcObjectPT, osResourceT any] func(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT) progress.ReconcileStatus

type ReconcileResourceActuator[orcObjectPT, osResourceT any] interface {
	// GetResourceReconcilers returns zero or more ResourceReconcilers to be executed during the current reconcile.
	//
	// All ResourceReconcilers returned will be executed in the order returned.
	// They will all be passed the orcObject and osResource passed to
	// GetResourceReconcilers. Note therefore that any state changes performed
	// by earlier ResourceReconcilers may not be reflected in the objects passed
	// to later ones.
	//
	// Failure of a ResourceReconciler does not prevent execution of later
	// ResourceReconcilers. All ResourceReconcilers will be executed, and their
	// ProgressStatuses and errors aggregated.
	//
	// NOTE: Contrary to the typical Go idiom, GetResourceReconcilers may return
	// both valid results, and an error. In this case, all returned
	// ResourceReconcilers will still be executed. The error returned by
	// GetResourceReconcilers itself will be aggregated with those returned by
	// the ResourceReconcilers. An example situation in which
	// GetResourceReconcilers itself might fail is if it fetched a list of
	// objects and returned a separate ResourceReconciler for each of them.
	GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller ResourceController) ([]ResourceReconciler[orcObjectPT, osResourceT], progress.ReconcileStatus)
}
