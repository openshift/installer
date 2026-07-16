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

package group

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

const controllerName = "group"

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=groups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=groups/status,verbs=get;update;patch

type groupReconcilerConstructor struct {
	scopeFactory        scope.Factory
	defaultResyncPeriod time.Duration
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return &groupReconcilerConstructor{scopeFactory: scopeFactory}
}

func (groupReconcilerConstructor) GetName() string {
	return controllerName
}

func (c *groupReconcilerConstructor) SetDefaultResyncPeriod(d time.Duration) {
	c.defaultResyncPeriod = d
}

var domainDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.GroupList, *orcv1alpha1.Domain](
	"spec.resource.domainRef",
	func(group *orcv1alpha1.Group) []string {
		resource := group.Spec.Resource
		if resource == nil || resource.DomainRef == nil {
			return nil
		}
		return []string{string(*resource.DomainRef)}
	},
	finalizer, externalObjectFieldOwner,
)

var domainImportDependency = dependency.NewDependency[*orcv1alpha1.GroupList, *orcv1alpha1.Domain](
	"spec.import.filter.domainRef",
	func(group *orcv1alpha1.Group) []string {
		resource := group.Spec.Import
		if resource == nil || resource.Filter == nil || resource.Filter.DomainRef == nil {
			return nil
		}
		return []string{string(*resource.Filter.DomainRef)}
	},
)

// SetupWithManager sets up the controller with the Manager.
func (c *groupReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)
	k8sClient := mgr.GetClient()

	domainWatchEventHandler, err := domainDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	domainImportWatchEventHandler, err := domainImportDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		Watches(&orcv1alpha1.Domain{}, domainWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Domain{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Domain{}, domainImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Domain{})),
		).
		For(&orcv1alpha1.Group{})

	if err := errors.Join(
		domainDependency.AddToManager(ctx, mgr),
		domainImportDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, mgr.GetClient(), builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, mgr.GetClient(), c.scopeFactory, groupHelperFactory{}, groupStatusWriter{}, c.defaultResyncPeriod)
	return builder.Complete(&r)
}
