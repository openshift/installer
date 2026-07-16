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

package user

import (
	"context"
	"errors"
	"time"

	corev1 "k8s.io/api/core/v1"
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

const controllerName = "user"

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=users,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=users/status,verbs=get;update;patch

type userReconcilerConstructor struct {
	scopeFactory        scope.Factory
	defaultResyncPeriod time.Duration
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return &userReconcilerConstructor{scopeFactory: scopeFactory}
}

func (userReconcilerConstructor) GetName() string {
	return controllerName
}

func (c *userReconcilerConstructor) SetDefaultResyncPeriod(d time.Duration) {
	c.defaultResyncPeriod = d
}

var domainDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.UserList, *orcv1alpha1.Domain](
	"spec.resource.domainRef",
	func(user *orcv1alpha1.User) []string {
		resource := user.Spec.Resource
		if resource == nil || resource.DomainRef == nil {
			return nil
		}
		return []string{string(*resource.DomainRef)}
	},
	finalizer, externalObjectFieldOwner,
)

var projectDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.UserList, *orcv1alpha1.Project](
	"spec.resource.defaultProjectRef",
	func(user *orcv1alpha1.User) []string {
		resource := user.Spec.Resource
		if resource == nil || resource.DefaultProjectRef == nil {
			return nil
		}
		return []string{string(*resource.DefaultProjectRef)}
	},
	finalizer, externalObjectFieldOwner,
)

var domainImportDependency = dependency.NewDependency[*orcv1alpha1.UserList, *orcv1alpha1.Domain](
	"spec.import.filter.domainRef",
	func(user *orcv1alpha1.User) []string {
		resource := user.Spec.Import
		if resource == nil || resource.Filter == nil || resource.Filter.DomainRef == nil {
			return nil
		}
		return []string{string(*resource.Filter.DomainRef)}
	},
)

var passwordDependency = dependency.NewDependency[*orcv1alpha1.UserList, *corev1.Secret](
	"spec.resource.passwordRef",
	func(user *orcv1alpha1.User) []string {
		resource := user.Spec.Resource
		if resource == nil || resource.PasswordRef == nil {
			return nil
		}
		return []string{string(*resource.PasswordRef)}
	},
)

// SetupWithManager sets up the controller with the Manager.
func (c *userReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)
	k8sClient := mgr.GetClient()

	domainWatchEventHandler, err := domainDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	projectWatchEventHandler, err := projectDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	domainImportWatchEventHandler, err := domainImportDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	passwordWatchEventHandler, err := passwordDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&orcv1alpha1.User{}).
		Watches(&orcv1alpha1.Domain{}, domainWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Domain{})),
		).
		Watches(&orcv1alpha1.Project{}, projectWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Project{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Domain{}, domainImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Domain{})),
		).
		// XXX: This is a general watch on secrets. A general watch on secrets
		// is undesirable because:
		// - It requires problematic RBAC
		// - Secrets are arbitrarily large, and we don't want to cache their contents
		//
		// These will require separate solutions. For the latter we should
		// probably use a MetadataOnly watch on secrets.
		Watches(&corev1.Secret{}, passwordWatchEventHandler)

	if err := errors.Join(
		domainDependency.AddToManager(ctx, mgr),
		projectDependency.AddToManager(ctx, mgr),
		domainImportDependency.AddToManager(ctx, mgr),
		passwordDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, mgr.GetClient(), builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, mgr.GetClient(), c.scopeFactory, userHelperFactory{}, userStatusWriter{}, c.defaultResyncPeriod)
	return builder.Complete(&r)
}
