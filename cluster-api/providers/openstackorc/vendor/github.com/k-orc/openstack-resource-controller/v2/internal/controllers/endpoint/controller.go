/*
Copyright The ORC Authors.

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

package endpoint

import (
	"context"
	"errors"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/reconciler"
	"github.com/k-orc/openstack-resource-controller/v2/internal/scope"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/credentials"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	"github.com/k-orc/openstack-resource-controller/v2/pkg/predicates"
)

const controllerName = "endpoint"

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=endpoints/status,verbs=get;update;patch

type endpointReconcilerConstructor struct {
	scopeFactory        scope.Factory
	defaultResyncPeriod time.Duration
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return &endpointReconcilerConstructor{scopeFactory: scopeFactory}
}

func (endpointReconcilerConstructor) GetName() string {
	return controllerName
}

func (c *endpointReconcilerConstructor) SetDefaultResyncPeriod(d time.Duration) {
	c.defaultResyncPeriod = d
}

var serviceDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.EndpointList, *orcv1alpha1.Service](
	"spec.resource.serviceRef",
	func(endpoint *orcv1alpha1.Endpoint) []string {
		resource := endpoint.Spec.Resource
		if resource == nil {
			return nil
		}
		return []string{string(resource.ServiceRef)}
	},
	finalizer, externalObjectFieldOwner,
)

var serviceImportDependency = dependency.NewDependency[*orcv1alpha1.EndpointList, *orcv1alpha1.Service](
	"spec.import.filter.serviceRef",
	func(endpoint *orcv1alpha1.Endpoint) []string {
		resource := endpoint.Spec.Import
		if resource == nil || resource.Filter == nil || resource.Filter.ServiceRef == nil {
			return nil
		}
		return []string{string(*resource.Filter.ServiceRef)}
	},
)

// SetupWithManager sets up the controller with the Manager.
func (c *endpointReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)
	k8sClient := mgr.GetClient()

	serviceWatchEventHandler, err := serviceDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	serviceImportWatchEventHandler, err := serviceImportDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		Watches(&orcv1alpha1.Service{}, serviceWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Service{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Service{}, serviceImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Service{})),
		).
		For(&orcv1alpha1.Endpoint{})

	if err := errors.Join(
		serviceDependency.AddToManager(ctx, mgr),
		serviceImportDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, mgr.GetClient(), builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, mgr.GetClient(), c.scopeFactory, endpointHelperFactory{}, endpointStatusWriter{}, c.defaultResyncPeriod)
	return builder.Complete(&r)
}
