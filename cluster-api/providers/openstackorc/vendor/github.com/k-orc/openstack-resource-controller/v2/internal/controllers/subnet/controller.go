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

package subnet

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/reconciler"
	"github.com/k-orc/openstack-resource-controller/v2/internal/scope"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/credentials"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	"github.com/k-orc/openstack-resource-controller/v2/pkg/predicates"
)

type subnetReconcilerConstructor struct {
	scopeFactory scope.Factory
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return subnetReconcilerConstructor{scopeFactory: scopeFactory}
}

func (subnetReconcilerConstructor) GetName() string {
	return controllerName
}

const controllerName = "subnet"

var (
	networkDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.SubnetList, *orcv1alpha1.Network](
		"spec.resource.networkRef",
		func(subnet *orcv1alpha1.Subnet) []string {
			resource := subnet.Spec.Resource
			if resource == nil {
				return nil
			}
			return []string{string(resource.NetworkRef)}
		},
		finalizer, externalObjectFieldOwner,
	)

	networkImportDependency = dependency.NewDependency[*orcv1alpha1.SubnetList, *orcv1alpha1.Network](
		"spec.import.filter.networkRef",
		func(subnet *orcv1alpha1.Subnet) []string {
			resource := subnet.Spec.Import
			if resource == nil || resource.Filter == nil {
				return nil
			}
			return []string{string(resource.Filter.NetworkRef)}
		},
	)

	routerDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.SubnetList, *orcv1alpha1.Router](
		"spec.resource.routerRef",
		func(subnet *orcv1alpha1.Subnet) []string {
			resource := subnet.Spec.Resource
			if resource == nil || resource.RouterRef == nil {
				return nil
			}
			return []string{string(*resource.RouterRef)}
		},
		finalizer, externalObjectFieldOwner,
	)

	projectDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.SubnetList, *orcv1alpha1.Project](
		"spec.resource.projectRef",
		func(subnet *orcv1alpha1.Subnet) []string {
			resource := subnet.Spec.Resource
			if resource == nil || resource.ProjectRef == nil {
				return nil
			}
			return []string{string(*resource.ProjectRef)}
		},
		finalizer, externalObjectFieldOwner,
	)

	projectImportDependency = dependency.NewDependency[*orcv1alpha1.SubnetList, *orcv1alpha1.Project](
		"spec.import.filter.projectRef",
		func(subnet *orcv1alpha1.Subnet) []string {
			resource := subnet.Spec.Import
			if resource == nil || resource.Filter == nil || resource.Filter.ProjectRef == nil {
				return nil
			}
			return []string{string(*resource.Filter.ProjectRef)}
		},
	)
)

// SetupWithManager sets up the controller with the Manager.
func (c subnetReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	controllerName := c.GetName()
	log := mgr.GetLogger().WithValues("controller", controllerName)
	k8sClient := mgr.GetClient()

	networkWatchEventHandler, err := networkDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	networkImportWatchEventHandler, err := networkImportDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	routerWatchEventHandler, err := routerDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	projectWatchEventHandler, err := projectDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	projectImportWatchEventHandler, err := projectImportDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&orcv1alpha1.Subnet{}).
		Watches(&orcv1alpha1.Network{}, networkWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Network{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Network{}, networkImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Network{})),
		).
		Watches(&orcv1alpha1.Router{}, routerWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Router{})),
		).
		Watches(&orcv1alpha1.Project{}, projectWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Project{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Project{}, projectImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Project{})),
		).
		Watches(&orcv1alpha1.RouterInterface{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
				log := log.WithValues("watch", "RouterInterface", "name", obj.GetName(), "namespace", obj.GetNamespace())
				routerInterface, ok := obj.(*orcv1alpha1.RouterInterface)
				if !ok {
					log.Info("Watch got unexpected object type", "type", fmt.Sprintf("%T", obj))
					return nil
				}
				subnetRef := routerInterface.Spec.SubnetRef
				if subnetRef == nil {
					return nil
				}
				return []reconcile.Request{
					{NamespacedName: types.NamespacedName{Namespace: routerInterface.Namespace, Name: string(*subnetRef)}},
				}
			}),
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.RouterInterface{})),
		)

	if err := errors.Join(
		networkDependency.AddToManager(ctx, mgr),
		networkImportDependency.AddToManager(ctx, mgr),
		routerDependency.AddToManager(ctx, mgr),
		projectDependency.AddToManager(ctx, mgr),
		projectImportDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, k8sClient, builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, k8sClient, c.scopeFactory, subnetHelperFactory{}, subnetStatusWriter{})
	return builder.Complete(&r)
}
