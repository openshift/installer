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
	"errors"

	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/pkg/predicates"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/reconciler"
	"github.com/k-orc/openstack-resource-controller/v2/internal/scope"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/credentials"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
)

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=floatingips,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=floatingips/status,verbs=get;update;patch

type floatingipReconcilerConstructor struct {
	scopeFactory scope.Factory
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return floatingipReconcilerConstructor{scopeFactory: scopeFactory}
}

func (floatingipReconcilerConstructor) GetName() string {
	return controllerName
}

const controllerName = "floatingip"

var (
	networkDep = dependency.NewDeletionGuardDependency[*orcv1alpha1.FloatingIPList, *orcv1alpha1.Network](
		"spec.resource.floatingNetworkRef",
		func(floatingip *orcv1alpha1.FloatingIP) []string {
			resource := floatingip.Spec.Resource
			if resource == nil || resource.FloatingNetworkRef == nil {
				return nil
			}
			return []string{string(ptr.Deref(resource.FloatingNetworkRef, ""))}
		},
		finalizer, externalObjectFieldOwner,
	)

	networkImportDep = dependency.NewDependency[*orcv1alpha1.FloatingIPList, *orcv1alpha1.Network](
		"spec.import.filter.floatingNetworkRef",
		func(floatingip *orcv1alpha1.FloatingIP) []string {
			resource := floatingip.Spec.Import
			if resource == nil || resource.Filter == nil || resource.Filter.FloatingNetworkRef == nil {
				return nil
			}
			return []string{string(ptr.Deref(resource.Filter.FloatingNetworkRef, ""))}
		},
	)

	subnetDep = dependency.NewDeletionGuardDependency[*orcv1alpha1.FloatingIPList, *orcv1alpha1.Subnet](
		"spec.resource.floatingSubnetRef",
		func(floatingip *orcv1alpha1.FloatingIP) []string {
			resource := floatingip.Spec.Resource
			if resource == nil || resource.FloatingSubnetRef == nil {
				return nil
			}
			return []string{string(ptr.Deref(resource.FloatingSubnetRef, ""))}
		},
		finalizer, externalObjectFieldOwner,
	)

	portDep = dependency.NewDeletionGuardDependency[*orcv1alpha1.FloatingIPList, *orcv1alpha1.Port](
		"spec.resource.portRef",
		func(floatingip *orcv1alpha1.FloatingIP) []string {
			resource := floatingip.Spec.Resource
			if resource == nil || resource.PortRef == nil {
				return nil
			}
			return []string{string(*resource.PortRef)}
		},
		finalizer, externalObjectFieldOwner,
	)

	portImportDep = dependency.NewDependency[*orcv1alpha1.FloatingIPList, *orcv1alpha1.Port](
		"spec.import.filter.portRef",
		func(floatingip *orcv1alpha1.FloatingIP) []string {
			resource := floatingip.Spec.Import
			if resource == nil || resource.Filter == nil || resource.Filter.PortRef == nil {
				return nil
			}
			return []string{string(ptr.Deref(resource.Filter.PortRef, ""))}
		},
	)

	projectDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.FloatingIPList, *orcv1alpha1.Project](
		"spec.resource.projectRef",
		func(floatingip *orcv1alpha1.FloatingIP) []string {
			resource := floatingip.Spec.Resource
			if resource == nil || resource.ProjectRef == nil {
				return nil
			}
			return []string{string(*resource.ProjectRef)}
		},
		finalizer, externalObjectFieldOwner,
	)

	projectImportDependency = dependency.NewDependency[*orcv1alpha1.FloatingIPList, *orcv1alpha1.Project](
		"spec.import.filter.projectRef",
		func(floatingip *orcv1alpha1.FloatingIP) []string {
			resource := floatingip.Spec.Import
			if resource == nil || resource.Filter == nil || resource.Filter.ProjectRef == nil {
				return nil
			}
			return []string{string(*resource.Filter.ProjectRef)}
		},
	)
)

// SetupWithManager sets up the controller with the Manager.
func (c floatingipReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := mgr.GetLogger().WithValues("controller", controllerName)
	k8sClient := mgr.GetClient()

	networkHandler, err := networkDep.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	networkImportWatchEventHandler, err := networkImportDep.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	subnetHandler, err := subnetDep.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	portHandler, err := portDep.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	portImportWatchEventHandler, err := portImportDep.WatchEventHandler(log, k8sClient)
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
		For(&orcv1alpha1.FloatingIP{}).
		Watches(&orcv1alpha1.Network{}, networkHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Network{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Network{}, networkImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Network{})),
		).
		Watches(&orcv1alpha1.Subnet{}, subnetHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Subnet{})),
		).
		Watches(&orcv1alpha1.Port{}, portHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Port{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Port{}, portImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Port{})),
		).
		Watches(&orcv1alpha1.Project{}, projectWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Project{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Project{}, projectImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Project{})),
		)

	if err := errors.Join(
		networkDep.AddToManager(ctx, mgr),
		networkImportDep.AddToManager(ctx, mgr),
		subnetDep.AddToManager(ctx, mgr),
		portDep.AddToManager(ctx, mgr),
		portImportDep.AddToManager(ctx, mgr),
		projectDependency.AddToManager(ctx, mgr),
		projectImportDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, k8sClient, builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, k8sClient, c.scopeFactory, floatingipHelperFactory{}, floatingipStatusWriter{})
	return builder.Complete(&r)
}
