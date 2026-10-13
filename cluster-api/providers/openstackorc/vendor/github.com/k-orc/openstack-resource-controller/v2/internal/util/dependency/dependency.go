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

package dependency

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/finalizers"
)

// NewDependency returns a new Dependency, which can perform tasks necessary to manage a dependency between 2 object types. The 2 object types are:
//   - Object: this is the 'source' object.
//   - Dependency: this is the object that a 'source' object may depend on.
//
// For example, a Port may depend on a Subnet, because it references one or more
// Subnets in its Addresses. In this case 'Object' is the Port, and 'Dependency'
// is the subnet.
//
// NewDependency has several type parameters, but only the first 2 are required as all the rest can be inferred. The 2 required parameters are:
//   - pointer to the List type of Object
//   - pointer to the Dependency type
//
// NewDependency takes the following arguments:
//   - indexName: a name representing the path to the Dependency reference in Object.
//   - getDependencyRefs: a function that takes a pointer to Object and returns a slice of strings containing the names of Dependencies
//
// Taking the Port -> Subnet example, the required type parameter is:
//   - *PortList: pointer to the list type of Port
//
// and the arguments are:
//   - indexName: "spec.resource.addresses[].subnetRef" - a symbolic path to the subnet reference in a Port
//   - getDependencyRefs: func(object *Port) []string{ ... returns a slice containing all subnetRefs in this Port's addresses ... }
func NewDependency[
	objectListTP ObjectListType[objectListT, objectT],
	depTP DependencyType[depT],

	objectTP ObjectType[objectT],
	objectT any, objectListT any, depT any,
](indexName string, getDependencyRefs func(objectTP) []string) Dependency[objectTP, objectListTP, depTP, objectT, objectListT, depT] {
	return Dependency[objectTP, objectListTP, depTP, objectT, objectListT, depT]{
		indexName:         indexName,
		getDependencyRefs: getDependencyRefs,
	}
}

type deletionGuardConfig struct {
	overrideDependencyName *string
}

type deletionGuardOpt = func(*deletionGuardConfig)

func OverrideDependencyName(name string) deletionGuardOpt {
	return func(opts *deletionGuardConfig) {
		opts.overrideDependencyName = &name
	}
}

// NewDeletionGuardDependency returns a Dependency which can additionally create a deletion guard for the dependency. See NewDependency for a discussion of the base functionality.
//
// In addition to the arguments required by NewDependency, NewDeletionGuardDependency requires:
// - finalizer: the string to add to Finalizers in objects that we depend on
// - fieldOwner: a client.FieldOwner identifying this controller when adding a finalizer to objects we depend on
func NewDeletionGuardDependency[
	objectListTP ObjectListType[objectListT, objectT],
	depTP DependencyType[depT],

	objectTP ObjectType[objectT],
	objectT any, objectListT any, depT any,
](indexName string, getDependencyRefs func(objectTP) []string, finalizer string, fieldOwner client.FieldOwner, opts ...deletionGuardOpt) DeletionGuardDependency[objectTP, objectListTP, depTP, objectT, objectListT, depT] {
	config := deletionGuardConfig{}
	for _, opt := range opts {
		opt(&config)
	}

	return DeletionGuardDependency[objectTP, objectListTP, depTP, objectT, objectListT, depT]{
		Dependency:             NewDependency[objectListTP, depTP](indexName, getDependencyRefs),
		finalizer:              finalizer,
		fieldOwner:             fieldOwner,
		overrideDependencyName: config.overrideDependencyName,
	}
}

type Dependency[
	objectTP ObjectType[objectT],
	objectListTP ObjectListType[objectListT, objectT],
	depTP DependencyType[depT],

	objectT any, objectListT any, depT any,
] struct {
	indexName         string
	getDependencyRefs func(objectTP) []string
}

type DeletionGuardDependency[
	objectTP ObjectType[objectT],
	objectListTP ObjectListType[objectListT, objectT],
	depTP DependencyType[depT],

	objectT any, objectListT any, depT any,
] struct {
	Dependency[objectTP, objectListTP, depTP, objectT, objectListT, depT]

	finalizer              string
	fieldOwner             client.FieldOwner
	overrideDependencyName *string
}

type ObjectType[objectT any] interface {
	*objectT
	client.Object
}

type ObjectListType[objectListT any, objectT any] interface {
	client.ObjectList
	*objectListT

	GetItems() []objectT
}

type DependencyType[depT any] interface {
	*depT
	client.Object
}

// GetObjectsForDependency returns a slice of all Objects which depend on the given Dependency.
func (d *Dependency[_, objectListTP, depTP, objectT, objectListT, _]) GetObjectsForDependency(ctx context.Context, k8sClient client.Client, dep depTP) ([]objectT, error) {
	var objectList objectListTP = new(objectListT)
	if err := k8sClient.List(ctx, objectList, client.InNamespace(dep.GetNamespace()), client.MatchingFields{d.indexName: dep.GetName()}); err != nil {
		return nil, err
	}
	return objectList.GetItems(), nil
}

// addIndexer adds the required field indexer for this dependency to a manager
// Called by AddToManager
func (d *Dependency[objectTP, _, _, objectT, _, _]) addIndexer(ctx context.Context, mgr ctrl.Manager) error {
	return mgr.GetFieldIndexer().IndexField(ctx, objectTP(new(objectT)), d.indexName, func(cObj client.Object) []string {
		obj, ok := cObj.(objectTP)
		if !ok {
			return nil
		}

		return d.getDependencyRefs(obj)
	})
}

func (d *Dependency[_, _, _, _, _, _]) AddToManager(ctx context.Context, mgr ctrl.Manager) error {
	return d.addIndexer(ctx, mgr)
}

// WatchEventHandler returns an EventHandler which maps a Dependency to all Objects which depend on it
func (d *Dependency[objectTP, _, depTP, _, _, depT]) WatchEventHandler(log logr.Logger, k8sClient client.Client) (handler.EventHandler, error) {
	depKind, err := getObjectKind(depTP(new(depT)), k8sClient.Scheme())
	if err != nil {
		return nil, err
	}

	log = log.WithValues("watch", depKind)
	return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
		log := log.WithValues("name", obj.GetName(), "namespace", obj.GetNamespace())

		dependency, ok := obj.(depTP)
		if !ok {
			log.Info("Watch got unexpected object type", "type", fmt.Sprintf("%T", obj))
			return nil
		}

		objects, err := d.GetObjectsForDependency(ctx, k8sClient, dependency)
		if err != nil {
			log.Error(err, fmt.Sprintf("listing %s", depKind))
			return nil
		}
		requests := make([]reconcile.Request, len(objects))
		for i := range objects {
			var object objectTP = &objects[i]
			request := &requests[i]

			request.Name = object.GetName()
			request.Namespace = object.GetNamespace()
		}
		return requests
	}), nil
}

// addDeletionGuard adds a deletion guard controller to the given manager appropriate for this dependency
// Called by AddToManager
func (d *DeletionGuardDependency[objectTP, _, _, _, _, _]) addDeletionGuard(mgr ctrl.Manager) error {
	getDependencyRefsForClientObject := func(cObj client.Object) []string {
		obj, ok := cObj.(objectTP)
		if !ok {
			return nil
		}
		return d.getDependencyRefs(obj)
	}

	return addDeletionGuard[objectTP](mgr, d.finalizer, d.fieldOwner, getDependencyRefsForClientObject, d.GetObjectsForDependency, d.overrideDependencyName)
}

// GetDependencies returns the dependencies of the given object, ensuring that all returned dependencies have the required finalizer. It returns:
// - a map of name -> object containing all objects which exist and are ready
// - a list of progressStatus for all dependencies which are not yet ready
// - an error
//
// Dependencies are filtered by the readyFilter argument. Dependencies which are not ready will be in progressStatus but not in the returned object map.
func (d *DeletionGuardDependency[objectTP, _, depTP, _, _, depT]) GetDependencies(ctx context.Context, k8sClient client.Client, obj objectTP, readyFilter func(depTP) bool) (map[string]depTP, progress.ReconcileStatus) {
	depKind, err := getObjectKind(depTP(new(depT)), k8sClient.Scheme())
	if err != nil {
		return nil, progress.WrapError(err)
	}

	var reconcileStatus progress.ReconcileStatus
	depsMap := make(map[string]depTP)
	for _, depRef := range d.getDependencyRefs(obj) {
		var dep depTP = new(depT)

		if depErr := k8sClient.Get(ctx, types.NamespacedName{Name: depRef, Namespace: obj.GetNamespace()}, dep); depErr != nil {
			if apierrors.IsNotFound(depErr) {
				reconcileStatus = reconcileStatus.WaitingOnObject(depKind, depRef, progress.WaitingOnCreation)
			} else {
				reconcileStatus = reconcileStatus.WithError(depErr)
			}

			continue
		}

		if readyFilter(dep) {
			// Don't add the finalizer until the dependency is ready. This makes
			// it easier to delete incorrectly created objects which never
			// became ready.
			if depErr := EnsureFinalizer(ctx, k8sClient, dep, d.finalizer, d.fieldOwner); depErr != nil {
				reconcileStatus = reconcileStatus.WithError(depErr)
				continue
			}

			depsMap[depRef] = dep
		} else {
			reconcileStatus = reconcileStatus.WaitingOnObject(depKind, depRef, progress.WaitingOnReady)
		}
	}

	return depsMap, reconcileStatus
}

// GetDependency is a convenience wrapper around GetDependencies when the caller only expects a single result.
func (d *DeletionGuardDependency[objectTP, _, depTP, _, _, depT]) GetDependency(ctx context.Context, k8sClient client.Client, obj objectTP, readyFilter func(depTP) bool) (depTP, progress.ReconcileStatus) {
	depsMap, reconcileStatus := d.GetDependencies(ctx, k8sClient, obj, readyFilter)
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}
	if len(depsMap) > 1 {
		// Programming error
		return nil, progress.WrapError(fmt.Errorf("GetDependencies returned multiple dependencies, expected one"))
	}
	for _, dep := range depsMap {
		return dep, nil
	}
	// Programming error
	return nil, progress.WrapError(fmt.Errorf("GetDependencies returned empty depsMap, progressStatus, and error"))
}

func (d *DeletionGuardDependency[objectTP, objectListTP, depTP, objectT, objectListT, depT]) AddToManager(ctx context.Context, mgr ctrl.Manager) error {
	return errors.Join(
		d.addIndexer(ctx, mgr),
		d.addDeletionGuard(mgr),
	)
}

// EnsureFinalizer adds a finalizer to the given object if it is not already present. It does nothing if the finalizer is already present.
func EnsureFinalizer(ctx context.Context, k8sClient client.Client, obj client.Object, finalizer string, fieldOwner client.FieldOwner) error {
	if slices.Contains(obj.GetFinalizers(), finalizer) {
		return nil
	}

	log := ctrl.LoggerFrom(ctx)
	log.V(logging.Verbose).Info("Adding finalizer", "objectName", obj.GetName(), "objectKind", obj.GetObjectKind().GroupVersionKind().Kind)
	patch := finalizers.SetFinalizerPatch(obj, finalizer)
	return k8sClient.Patch(ctx, obj, patch, client.ForceOwnership, fieldOwner)
}
