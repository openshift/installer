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

package roleassignment

import (
	"context"
	"errors"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/scope"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/credentials"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	"github.com/k-orc/openstack-resource-controller/v2/pkg/predicates"
)

const controllerName = "roleassignment"

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=roleassignments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=roleassignments/status,verbs=get;update;patch

type roleassignmentReconcilerConstructor struct {
	scopeFactory        scope.Factory
	defaultResyncPeriod time.Duration
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return &roleassignmentReconcilerConstructor{scopeFactory: scopeFactory}
}

func (roleassignmentReconcilerConstructor) GetName() string {
	return controllerName
}

func (c *roleassignmentReconcilerConstructor) SetDefaultResyncPeriod(d time.Duration) {
	c.defaultResyncPeriod = d
}

var roleDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.RoleAssignmentList, *orcv1alpha1.Role](
	"spec.resource.roleRef",
	func(roleassignment *orcv1alpha1.RoleAssignment) []string {
		resource := roleassignment.Spec.Resource
		if resource == nil {
			return nil
		}
		return []string{string(resource.RoleRef)}
	},
	finalizer, externalObjectFieldOwner,
)

var userDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.RoleAssignmentList, *orcv1alpha1.User](
	"spec.resource.userRef",
	func(roleassignment *orcv1alpha1.RoleAssignment) []string {
		resource := roleassignment.Spec.Resource
		if resource == nil || resource.UserRef == nil {
			return nil
		}
		return []string{string(*resource.UserRef)}
	},
	finalizer, externalObjectFieldOwner,
)

var groupDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.RoleAssignmentList, *orcv1alpha1.Group](
	"spec.resource.groupRef",
	func(roleassignment *orcv1alpha1.RoleAssignment) []string {
		resource := roleassignment.Spec.Resource
		if resource == nil || resource.GroupRef == nil {
			return nil
		}
		return []string{string(*resource.GroupRef)}
	},
	finalizer, externalObjectFieldOwner,
)

var projectDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.RoleAssignmentList, *orcv1alpha1.Project](
	"spec.resource.projectRef",
	func(roleassignment *orcv1alpha1.RoleAssignment) []string {
		resource := roleassignment.Spec.Resource
		if resource == nil || resource.ProjectRef == nil {
			return nil
		}
		return []string{string(*resource.ProjectRef)}
	},
	finalizer, externalObjectFieldOwner,
)

var domainDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.RoleAssignmentList, *orcv1alpha1.Domain](
	"spec.resource.domainRef",
	func(roleassignment *orcv1alpha1.RoleAssignment) []string {
		resource := roleassignment.Spec.Resource
		if resource == nil || resource.DomainRef == nil {
			return nil
		}
		return []string{string(*resource.DomainRef)}
	},
	finalizer, externalObjectFieldOwner,
)

var roleImportDependency = dependency.NewDependency[*orcv1alpha1.RoleAssignmentList, *orcv1alpha1.Role](
	"spec.import.filter.roleRef",
	func(roleassignment *orcv1alpha1.RoleAssignment) []string {
		resource := roleassignment.Spec.Import
		if resource == nil || resource.Filter == nil || resource.Filter.RoleRef == nil {
			return nil
		}
		return []string{string(*resource.Filter.RoleRef)}
	},
)

var userImportDependency = dependency.NewDependency[*orcv1alpha1.RoleAssignmentList, *orcv1alpha1.User](
	"spec.import.filter.userRef",
	func(roleassignment *orcv1alpha1.RoleAssignment) []string {
		resource := roleassignment.Spec.Import
		if resource == nil || resource.Filter == nil || resource.Filter.UserRef == nil {
			return nil
		}
		return []string{string(*resource.Filter.UserRef)}
	},
)

var groupImportDependency = dependency.NewDependency[*orcv1alpha1.RoleAssignmentList, *orcv1alpha1.Group](
	"spec.import.filter.groupRef",
	func(roleassignment *orcv1alpha1.RoleAssignment) []string {
		resource := roleassignment.Spec.Import
		if resource == nil || resource.Filter == nil || resource.Filter.GroupRef == nil {
			return nil
		}
		return []string{string(*resource.Filter.GroupRef)}
	},
)

var projectImportDependency = dependency.NewDependency[*orcv1alpha1.RoleAssignmentList, *orcv1alpha1.Project](
	"spec.import.filter.projectRef",
	func(roleassignment *orcv1alpha1.RoleAssignment) []string {
		resource := roleassignment.Spec.Import
		if resource == nil || resource.Filter == nil || resource.Filter.ProjectRef == nil {
			return nil
		}
		return []string{string(*resource.Filter.ProjectRef)}
	},
)

var domainImportDependency = dependency.NewDependency[*orcv1alpha1.RoleAssignmentList, *orcv1alpha1.Domain](
	"spec.import.filter.domainRef",
	func(roleassignment *orcv1alpha1.RoleAssignment) []string {
		resource := roleassignment.Spec.Import
		if resource == nil || resource.Filter == nil || resource.Filter.DomainRef == nil {
			return nil
		}
		return []string{string(*resource.Filter.DomainRef)}
	},
)

// SetupWithManager sets up the controller with the Manager.
func (c roleassignmentReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)
	k8sClient := mgr.GetClient()

	roleWatchEventHandler, err := roleDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	userWatchEventHandler, err := userDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	groupWatchEventHandler, err := groupDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	projectWatchEventHandler, err := projectDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	domainWatchEventHandler, err := domainDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	roleImportWatchEventHandler, err := roleImportDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	userImportWatchEventHandler, err := userImportDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	groupImportWatchEventHandler, err := groupImportDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	projectImportWatchEventHandler, err := projectImportDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	domainImportWatchEventHandler, err := domainImportDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		Watches(&orcv1alpha1.Role{}, roleWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Role{})),
		).
		Watches(&orcv1alpha1.User{}, userWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.User{})),
		).
		Watches(&orcv1alpha1.Group{}, groupWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Group{})),
		).
		Watches(&orcv1alpha1.Project{}, projectWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Project{})),
		).
		Watches(&orcv1alpha1.Domain{}, domainWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Domain{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Role{}, roleImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Role{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.User{}, userImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.User{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Group{}, groupImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Group{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Project{}, projectImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Project{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Domain{}, domainImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Domain{})),
		).
		For(&orcv1alpha1.RoleAssignment{})

	if err := errors.Join(
		roleDependency.AddToManager(ctx, mgr),
		userDependency.AddToManager(ctx, mgr),
		groupDependency.AddToManager(ctx, mgr),
		projectDependency.AddToManager(ctx, mgr),
		domainDependency.AddToManager(ctx, mgr),
		roleImportDependency.AddToManager(ctx, mgr),
		userImportDependency.AddToManager(ctx, mgr),
		groupImportDependency.AddToManager(ctx, mgr),
		projectImportDependency.AddToManager(ctx, mgr),
		domainImportDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, mgr.GetClient(), builder, credentialsDependency),
	); err != nil {
		return err
	}

	// Custom reconciler for role assignments (relationships, not resources with IDs)
	reconciler := &roleassignmentReconciler{
		client:              mgr.GetClient(),
		scopeFactory:        c.scopeFactory,
		defaultResyncPeriod: c.defaultResyncPeriod,
	}
	return builder.Complete(reconciler)
}
