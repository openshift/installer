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

package server

import (
	"context"
	"errors"

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

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=servers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=servers/status,verbs=get;update;patch

type serverReconcilerConstructor struct {
	scopeFactory scope.Factory
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return serverReconcilerConstructor{scopeFactory: scopeFactory}
}

func (serverReconcilerConstructor) GetName() string {
	return controllerName
}

const controllerName = "server"

var (
	// No deletion guard for flavor, because flavors can be safely deleted while
	// referenced by a server
	flavorDependency = dependency.NewDependency[*orcv1alpha1.ServerList, *orcv1alpha1.Flavor](
		"spec.resource.flavorRef",
		func(server *orcv1alpha1.Server) []string {
			resource := server.Spec.Resource
			if resource == nil {
				return nil
			}

			return []string{string(resource.FlavorRef)}
		},
	)

	// Image can sometimes, but not always (e.g. when an RBD-backed image has
	// been cloned), be safely deleted while referenced by a server. We just
	// prevent it always.
	imageDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.ServerList, *orcv1alpha1.Image](
		"spec.resource.imageRef",
		func(server *orcv1alpha1.Server) []string {
			resource := server.Spec.Resource
			if resource == nil {
				return nil
			}

			return []string{string(resource.ImageRef)}
		},
		finalizer, externalObjectFieldOwner,
	)

	portDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.ServerList, *orcv1alpha1.Port](
		"spec.resource.ports",
		func(server *orcv1alpha1.Server) []string {
			resource := server.Spec.Resource
			if resource == nil {
				return nil
			}

			refs := make([]string, 0, len(resource.Ports))
			for i := range resource.Ports {
				port := &resource.Ports[i]
				if port.PortRef != nil {
					refs = append(refs, string(*port.PortRef))
				}
			}
			return refs
		},
		finalizer, externalObjectFieldOwner,
	)

	// No deletion guard for server group, because server group can be safely deleted while
	// referenced by a server
	serverGroupDependency = dependency.NewDependency[*orcv1alpha1.ServerList, *orcv1alpha1.ServerGroup](
		"spec.resource.serverGroupRef",
		func(server *orcv1alpha1.Server) []string {
			resource := server.Spec.Resource
			if resource == nil || resource.ServerGroupRef == nil {
				return nil
			}

			return []string{string(*resource.ServerGroupRef)}
		},
	)

	// We don't need a deletion guard on the user-data secret because it's only
	// used on creation.
	userDataDependency = dependency.NewDependency[*orcv1alpha1.ServerList, *corev1.Secret](
		"spec.resource.userData.secretRef",
		func(server *orcv1alpha1.Server) []string {
			resource := server.Spec.Resource
			if resource == nil || resource.UserData == nil || resource.UserData.SecretRef == nil {
				return nil
			}

			return []string{string(*resource.UserData.SecretRef)}
		},
	)

	volumeDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.ServerList, *orcv1alpha1.Volume](
		"spec.resource.volumes",
		func(server *orcv1alpha1.Server) []string {
			resource := server.Spec.Resource
			if resource == nil {
				return nil
			}

			refs := make([]string, 0, len(resource.Volumes))
			for i := range resource.Volumes {
				volume := &resource.Volumes[i]
				refs = append(refs, string(volume.VolumeRef))
			}
			return refs
		},
		finalizer, externalObjectFieldOwner,
	)
)

// SetupWithManager sets up the controller with the Manager.
func (c serverReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := mgr.GetLogger().WithValues("controller", controllerName)
	k8sClient := mgr.GetClient()

	flavorWatchEventHandler, err := flavorDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}
	imageWatchEventHandler, err := imageDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}
	portWatchEventHandler, err := portDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}
	userDataWatchEventHandler, err := userDataDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}
	serverGroupWatchEventHandler, err := serverGroupDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}
	volumeWatchEventHandler, err := volumeDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&orcv1alpha1.Server{}).
		Watches(&orcv1alpha1.Flavor{}, flavorWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Flavor{})),
		).
		Watches(&orcv1alpha1.Image{}, imageWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Image{})),
		).
		Watches(&orcv1alpha1.Port{}, portWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Port{})),
		).
		Watches(&orcv1alpha1.ServerGroup{}, serverGroupWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.ServerGroup{})),
		).
		Watches(&orcv1alpha1.Volume{}, volumeWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Volume{})),
		).
		// XXX: This is a general watch on secrets. A general watch on secrets
		// is undesirable because:
		// - It requires problematic RBAC
		// - Secrets are arbitrarily large, and we don't want to cache their contents
		//
		// These will require separate solutions. For the latter we should
		// probably use a MetadataOnly watch only secrets.
		Watches(&corev1.Secret{}, userDataWatchEventHandler)

	if err := errors.Join(
		flavorDependency.AddToManager(ctx, mgr),
		imageDependency.AddToManager(ctx, mgr),
		portDependency.AddToManager(ctx, mgr),
		serverGroupDependency.AddToManager(ctx, mgr),
		userDataDependency.AddToManager(ctx, mgr),
		volumeDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, k8sClient, builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, k8sClient, c.scopeFactory, serverHelperFactory{}, serverStatusWriter{})
	return builder.Complete(&r)
}
