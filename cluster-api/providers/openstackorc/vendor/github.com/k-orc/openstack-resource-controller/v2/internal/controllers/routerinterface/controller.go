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

package routerinterface

import (
	"context"
	"errors"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/pkg/predicates"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/scope"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	orcstrings "github.com/k-orc/openstack-resource-controller/v2/internal/util/strings"
)

const (
	// The time to wait between polling a port for ACTIVE status
	portStatusPollingPeriod = 1 * time.Second
)

type routerInterfaceReconcilerConstructor struct {
	scopeFactory scope.Factory
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return routerInterfaceReconcilerConstructor{scopeFactory: scopeFactory}
}

func (routerInterfaceReconcilerConstructor) GetName() string {
	return controllerName
}

// orcRouterInterfaceReconciler reconciles an ORC Subnet.
type orcRouterInterfaceReconciler struct {
	client       client.Client
	scopeFactory scope.Factory
}

const controllerName = "routerinterface"

var (
	finalizer  = orcstrings.GetFinalizerName(controllerName)
	fieldOwner = orcstrings.GetSSAFieldOwner(controllerName)

	routerDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.RouterInterfaceList, *orcv1alpha1.Router](
		"spec.routerRef",
		func(routerIf *orcv1alpha1.RouterInterface) []string {
			return []string{string(routerIf.Spec.RouterRef)}
		},
		finalizer, fieldOwner,
	)

	subnetDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.RouterInterfaceList, *orcv1alpha1.Subnet](
		"spec.subnetRef",
		func(routerIf *orcv1alpha1.RouterInterface) []string {
			if routerIf.Spec.SubnetRef == nil {
				return nil
			}
			return []string{string(*routerIf.Spec.SubnetRef)}
		},
		finalizer, fieldOwner,
	)
)

// SetupWithManager sets up the controller with the Manager.
func (c routerInterfaceReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := mgr.GetLogger().WithValues("controller", controllerName)

	if err := errors.Join(
		routerDependency.AddToManager(ctx, mgr),
		subnetDependency.AddToManager(ctx, mgr),
	); err != nil {
		return err
	}

	k8sClient := mgr.GetClient()

	// NOTE: RouterInterface can't use the watch event handler from its
	// dependencies because it reconciles Routers, not RouterInterfaces.

	reconciler := orcRouterInterfaceReconciler{
		client:       k8sClient,
		scopeFactory: c.scopeFactory,
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&orcv1alpha1.Router{}, builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Router{}))).
		Named(controllerName).
		Watches(&orcv1alpha1.RouterInterface{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
				log := ctrl.LoggerFrom(ctx).WithValues("watch", "RouterInterface", "name", obj.GetName(), "namespace", obj.GetNamespace())
				routerInterface, ok := obj.(*orcv1alpha1.RouterInterface)
				if !ok {
					log.Info("Watch got unexpected object type", "type", fmt.Sprintf("%T", obj))
					return nil
				}
				return []reconcile.Request{
					{NamespacedName: types.NamespacedName{Namespace: routerInterface.Namespace, Name: string(routerInterface.Spec.RouterRef)}},
				}
			}),
		).
		Watches(&orcv1alpha1.Subnet{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
				log := ctrl.LoggerFrom(ctx).WithValues("watch", "Subnet", "name", obj.GetName(), "namespace", obj.GetNamespace())

				subnet, ok := obj.(*orcv1alpha1.Subnet)
				if !ok {
					log.Info("Watch got unexpected object type", "type", fmt.Sprintf("%T", obj))
					return nil
				}

				routerIfs, err := subnetDependency.GetObjectsForDependency(ctx, k8sClient, subnet)
				if err != nil {
					log.Error(err, "fetching routers for subnet")
					return nil
				}

				var requests []reconcile.Request
				for i := range routerIfs {
					routerIf := &routerIfs[i]
					if routerIf.Spec.SubnetRef != nil {
						requests = append(requests, reconcile.Request{
							NamespacedName: types.NamespacedName{Name: string(routerIf.Spec.RouterRef), Namespace: subnet.Namespace},
						})
					}
				}
				return requests
			}),
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Subnet{})),
		).
		WithOptions(options).
		Complete(&reconciler)
}
