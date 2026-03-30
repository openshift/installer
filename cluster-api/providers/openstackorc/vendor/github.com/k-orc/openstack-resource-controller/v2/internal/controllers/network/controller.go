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

package network

import (
	"context"
	"errors"

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

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=networks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=networks/status,verbs=get;update;patch

const controllerName = "network"

var (
	projectDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.NetworkList, *orcv1alpha1.Project](
		"spec.resource.projectRef",
		func(network *orcv1alpha1.Network) []string {
			resource := network.Spec.Resource
			if resource == nil || resource.ProjectRef == nil {
				return nil
			}
			return []string{string(*resource.ProjectRef)}
		},
		finalizer, externalObjectFieldOwner,
	)

	projectImportDependency = dependency.NewDependency[*orcv1alpha1.NetworkList, *orcv1alpha1.Project](
		"spec.import.filter.projectRef",
		func(network *orcv1alpha1.Network) []string {
			resource := network.Spec.Import
			if resource == nil || resource.Filter == nil || resource.Filter.ProjectRef == nil {
				return nil
			}
			return []string{string(*resource.Filter.ProjectRef)}
		},
	)
)

type networkReconcilerConstructor struct {
	scopeFactory scope.Factory
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return networkReconcilerConstructor{
		scopeFactory: scopeFactory,
	}
}

func (networkReconcilerConstructor) GetName() string {
	return controllerName
}

// SetupWithManager sets up the controller with the Manager.
func (c networkReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)
	k8sClient := mgr.GetClient()

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
		For(&orcv1alpha1.Network{}).
		Watches(&orcv1alpha1.Project{}, projectWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Project{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Project{}, projectImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Project{})),
		)

	if err := errors.Join(
		projectDependency.AddToManager(ctx, mgr),
		projectImportDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, mgr.GetClient(), builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, mgr.GetClient(), c.scopeFactory, networkHelperFactory{}, networkStatusWriter{})
	return builder.Complete(&r)
}
