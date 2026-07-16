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

package applicationcredential

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

const controllerName = "applicationcredential"

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=applicationcredentials,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=applicationcredentials/status,verbs=get;update;patch

var (
	// We don't need a deletion guard on the application credential secret because it's only
	// used on creation.
	secretDependency = dependency.NewDependency[*orcv1alpha1.ApplicationCredentialList, *corev1.Secret](
		"spec.resource.secretRef",
		func(applicationcredential *orcv1alpha1.ApplicationCredential) []string {
			resource := applicationcredential.Spec.Resource
			if resource == nil {
				return nil
			}

			return []string{string(resource.SecretRef)}
		},
	)

	roleDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.ApplicationCredentialList, *orcv1alpha1.Role](
		"spec.resource.roleRefs",
		func(applicationcredential *orcv1alpha1.ApplicationCredential) []string {
			resource := applicationcredential.Spec.Resource
			if resource == nil {
				return nil
			}

			roles := make([]string, len(resource.RoleRefs))
			for i := range resource.RoleRefs {
				roles[i] = string(resource.RoleRefs[i])
			}
			return roles
		},
		finalizer, externalObjectFieldOwner,
	)

	serviceDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.ApplicationCredentialList, *orcv1alpha1.Service](
		"spec.resource.accessRules[].serviceRef",
		func(applicationcredential *orcv1alpha1.ApplicationCredential) []string {
			resource := applicationcredential.Spec.Resource
			if resource == nil {
				return nil
			}

			services := make([]string, 0)
			for i := range resource.AccessRules {
				if resource.AccessRules[i].ServiceRef == nil {
					continue
				}
				services = append(services, string(*resource.AccessRules[i].ServiceRef))
			}
			return services
		},
		finalizer, externalObjectFieldOwner,
	)
)

type applicationcredentialReconcilerConstructor struct {
	scopeFactory        scope.Factory
	defaultResyncPeriod time.Duration
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return &applicationcredentialReconcilerConstructor{scopeFactory: scopeFactory}
}

func (applicationcredentialReconcilerConstructor) GetName() string {
	return controllerName
}

func (c *applicationcredentialReconcilerConstructor) SetDefaultResyncPeriod(d time.Duration) {
	c.defaultResyncPeriod = d
}

var userDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.ApplicationCredentialList, *orcv1alpha1.User](
	"spec.resource.userRef",
	func(applicationcredential *orcv1alpha1.ApplicationCredential) []string {
		resource := applicationcredential.Spec.Resource
		if resource == nil {
			return nil
		}
		return []string{string(resource.UserRef)}
	},
	finalizer, externalObjectFieldOwner,
)

var userImportDependency = dependency.NewDependency[*orcv1alpha1.ApplicationCredentialList, *orcv1alpha1.User](
	"spec.import.filter.userRef",
	func(applicationcredential *orcv1alpha1.ApplicationCredential) []string {
		resource := applicationcredential.Spec.Import
		if resource == nil || resource.Filter == nil {
			return nil
		}
		return []string{string(resource.Filter.UserRef)}
	},
)

// SetupWithManager sets up the controller with the Manager.
func (c *applicationcredentialReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)
	k8sClient := mgr.GetClient()

	userWatchEventHandler, err := userDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	userImportWatchEventHandler, err := userImportDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	secretWatchEventHandler, err := secretDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	roleWatchEventHandler, err := roleDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	serviceWatchEventHandler, err := serviceDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		Watches(&orcv1alpha1.User{}, userWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.User{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.User{}, userImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.User{})),
		).
		// XXX: This is a general watch on secrets. A general watch on secrets
		// is undesirable because:
		// - It requires problematic RBAC
		// - Secrets are arbitrarily large, and we don't want to cache their contents
		//
		// These will require separate solutions. For the latter we should
		// probably use a MetadataOnly watch only secrets.
		Watches(&corev1.Secret{}, secretWatchEventHandler).
		Watches(&orcv1alpha1.Role{}, roleWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Role{})),
		).
		Watches(&orcv1alpha1.Service{}, serviceWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Service{})),
		).
		For(&orcv1alpha1.ApplicationCredential{})

	if err := errors.Join(
		userDependency.AddToManager(ctx, mgr),
		userImportDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		secretDependency.AddToManager(ctx, mgr),
		roleDependency.AddToManager(ctx, mgr),
		serviceDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, mgr.GetClient(), builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, mgr.GetClient(), c.scopeFactory, applicationcredentialHelperFactory{}, applicationcredentialStatusWriter{}, c.defaultResyncPeriod)
	return builder.Complete(&r)
}
