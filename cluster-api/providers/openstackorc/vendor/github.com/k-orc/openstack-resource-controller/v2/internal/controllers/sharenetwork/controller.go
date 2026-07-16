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

package sharenetwork

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

const controllerName = "sharenetwork"

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=sharenetworks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=sharenetworks/status,verbs=get;update;patch

type sharenetworkReconcilerConstructor struct {
	scopeFactory        scope.Factory
	defaultResyncPeriod time.Duration
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return &sharenetworkReconcilerConstructor{scopeFactory: scopeFactory}
}

func (sharenetworkReconcilerConstructor) GetName() string {
	return controllerName
}

func (c *sharenetworkReconcilerConstructor) SetDefaultResyncPeriod(d time.Duration) {
	c.defaultResyncPeriod = d
}

var networkDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.ShareNetworkList, *orcv1alpha1.Network](
	"spec.resource.networkRef",
	func(sharenetwork *orcv1alpha1.ShareNetwork) []string {
		resource := sharenetwork.Spec.Resource
		if resource == nil || resource.NetworkRef == nil {
			return nil
		}
		return []string{string(*resource.NetworkRef)}
	},
	finalizer, externalObjectFieldOwner,
)

var subnetDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.ShareNetworkList, *orcv1alpha1.Subnet](
	"spec.resource.subnetRef",
	func(sharenetwork *orcv1alpha1.ShareNetwork) []string {
		resource := sharenetwork.Spec.Resource
		if resource == nil || resource.SubnetRef == nil {
			return nil
		}
		return []string{string(*resource.SubnetRef)}
	},
	finalizer, externalObjectFieldOwner,
)

// SetupWithManager sets up the controller with the Manager.
func (c *sharenetworkReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)
	k8sClient := mgr.GetClient()

	networkWatchEventHandler, err := networkDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	subnetWatchEventHandler, err := subnetDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		Watches(&orcv1alpha1.Network{}, networkWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Network{})),
		).
		Watches(&orcv1alpha1.Subnet{}, subnetWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Subnet{})),
		).
		For(&orcv1alpha1.ShareNetwork{})

	if err := errors.Join(
		networkDependency.AddToManager(ctx, mgr),
		subnetDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, mgr.GetClient(), builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, mgr.GetClient(), c.scopeFactory, sharenetworkHelperFactory{}, sharenetworkStatusWriter{}, c.defaultResyncPeriod)
	return builder.Complete(&r)
}
