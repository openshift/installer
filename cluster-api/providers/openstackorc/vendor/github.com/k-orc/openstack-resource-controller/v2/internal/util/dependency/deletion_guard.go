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
	"fmt"
	"strings"

	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/finalizers"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// A deletion guard is a controller which prevents the deletion of objects that objects of another type depend on.
//
// Example: Subnet depends on Network
//
// We add a deletion guard to Network that prevents the Network from being
// deleted if it is still in use by any Subnet. It is added by the Subnet
// controller, but it is a separate controller which reconciles Network objects.

func addDeletionGuard[objTP ObjectType[objT], objT any, depTP ObjectType[depT], depT any](
	mgr ctrl.Manager, finalizer string, fieldOwner client.FieldOwner,
	getDepRefsFromObject func(client.Object) []string,
	getObjectsFromDep func(context.Context, client.Client, depTP) ([]objT, error),
	overrideDependencyName *string,
) error {
	var depSpecimen depTP = new(depT)
	var objSpecimen objTP = new(objT)

	scheme := mgr.GetScheme()
	depKind, err := getObjectKind(depSpecimen, scheme)
	if err != nil {
		return err
	}
	objKind, err := getObjectKind(objSpecimen, scheme)
	if err != nil {
		return err
	}

	dependencyName := ptr.Deref(overrideDependencyName, strings.ToLower(depKind))
	controllerName := dependencyName + "_deletion_guard_for_" + strings.ToLower(objKind)

	// deletionGuard reconciles the dependency object
	// If the dependency is marked deleted, we remove the finalizer only when there are no objects referencing it
	deletionGuard := reconcile.Func(func(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
		log := ctrl.LoggerFrom(ctx, "name", req.Name, "namespace", req.Namespace)
		log.V(logging.Verbose).Info("Reconciling deletion guard")

		k8sClient := mgr.GetClient()

		var dep depTP = new(depT)
		err := k8sClient.Get(ctx, req.NamespacedName, dep)
		if err != nil {
			if apierrors.IsNotFound(err) {
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, err
		}

		// Nothing to do if the object isn't marked deleted
		// NOTE: we also try to arrange our triggers so we won't be reconciled in this case
		if dep.GetDeletionTimestamp().IsZero() {
			log.V(logging.Verbose).Info("Dependency is not marked deleted")
			return ctrl.Result{}, nil
		}

		log.V(logging.Debug).Info("Handling delete")

		refObjects, err := getObjectsFromDep(ctx, k8sClient, dep)
		if err != nil {
			return reconcile.Result{}, nil
		}

		depUID := dep.GetUID()
		depOwns := func(obj objTP) bool {
			owners := obj.GetOwnerReferences()
			for i := range owners {
				owner := &owners[i]
				if owner.Kind == depKind && owner.UID == depUID {
					return true
				}
			}
			return false
		}

		// Don't proceed if there are any referring objects, except owners of this object.
		// We don't block the deletion of the object which created us, because
		// that would cause a deadlock.
		for i := range refObjects {
			refObject := &refObjects[i]
			if !depOwns(refObject) {
				log.V(logging.Verbose).Info("Waiting for dependencies", "dependencies", len(refObjects))
				return ctrl.Result{}, nil
			}
		}

		log.V(logging.Verbose).Info("Removing finalizer")
		patch := finalizers.RemoveFinalizerPatch(dep)
		return ctrl.Result{}, k8sClient.Patch(ctx, dep, patch, client.ForceOwnership, fieldOwner)
	})

	// Register deletionGuard with the manager as a reconciler of the
	// dependency.  We also watch for referring objects, but we're only
	// interested in deletion events.  We need to ensure that if the depdency
	// object is marked deleted we will continue to call deletionGuard every
	// time a referring object is deleted so that we will eventually be called
	// when the last dependent object is deleted and we can remove the
	// dependency.
	err = builder.ControllerManagedBy(mgr).
		For(depSpecimen,
			// Only reconcile objects which are marked deleted and have our finalizer
			builder.WithPredicates(predicate.NewPredicateFuncs(func(obj client.Object) bool {
				if obj.GetDeletionTimestamp().IsZero() {
					return false
				}
				for _, objFinalizer := range obj.GetFinalizers() {
					if objFinalizer == finalizer {
						return true
					}
				}
				return false
			}))).
		Watches(objSpecimen,
			handler.Funcs{
				DeleteFunc: func(ctx context.Context, evt event.TypedDeleteEvent[client.Object], q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
					for _, guarded := range getDepRefsFromObject(evt.Object) {
						q.Add(reconcile.Request{
							NamespacedName: types.NamespacedName{
								Namespace: evt.Object.GetNamespace(),
								Name:      guarded,
							},
						})
					}
				},
			},
		).
		Named(controllerName).
		Complete(deletionGuard)

	if err != nil {
		return fmt.Errorf("failed to construct %s deletion guard for %s controller: %w", depKind, objKind, err)
	}

	return nil
}

func getObjectKind(obj runtime.Object, scheme *runtime.Scheme) (string, error) {
	gvks, _, err := scheme.ObjectKinds(obj)
	if err != nil {
		return "", fmt.Errorf("looking up GVK for guarded object %T: %w", obj, err)
	}
	if len(gvks) == 0 {
		return "", fmt.Errorf("no registered kind for guarded object %T", obj)
	}

	return gvks[0].Kind, nil
}
