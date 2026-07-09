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

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=servers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=servers/status,verbs=get;update;patch

type serverReconcilerConstructor struct {
	scopeFactory        scope.Factory
	defaultResyncPeriod time.Duration
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return &serverReconcilerConstructor{scopeFactory: scopeFactory}
}

func (serverReconcilerConstructor) GetName() string {
	return controllerName
}

func (c *serverReconcilerConstructor) SetDefaultResyncPeriod(d time.Duration) {
	c.defaultResyncPeriod = d
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
			if resource == nil || resource.ImageRef == nil {
				return nil
			}

			return []string{string(*resource.ImageRef)}
		},
		finalizer, externalObjectFieldOwner,
	)

	// bootVolumeDependency handles the boot volume specified in bootVolume for boot-from-volume.
	// This volume is attached at server creation time as the root disk.
	// deletion guard is in place because the server cannot boot without its root volume.
	// OverrideDependencyName is used to avoid conflict with volumeDependency which also
	// creates a Volume deletion guard for Server.
	bootVolumeDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.ServerList, *orcv1alpha1.Volume](
		"spec.resource.bootVolume.volumeRef",
		func(server *orcv1alpha1.Server) []string {
			resource := server.Spec.Resource
			if resource == nil || resource.BootVolume == nil {
				return nil
			}
			return []string{string(resource.BootVolume.VolumeRef)}
		},
		finalizer, externalObjectFieldOwner,
		dependency.OverrideDependencyName("bootvolume"),
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
		"spec.resource.schedulerHints.serverGroupRef",
		func(server *orcv1alpha1.Server) []string {
			resource := server.Spec.Resource
			if resource == nil || resource.SchedulerHints == nil || resource.SchedulerHints.ServerGroupRef == nil {
				return nil
			}

			return []string{string(*resource.SchedulerHints.ServerGroupRef)}
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

	// We don't need a deletion guard on the keypair because it's only
	// used on creation. The keypair reference is injected during server boot.
	keypairDependency = dependency.NewDependency[*orcv1alpha1.ServerList, *orcv1alpha1.KeyPair](
		"spec.resource.keypairRef",
		func(server *orcv1alpha1.Server) []string {
			resource := server.Spec.Resource
			if resource == nil || resource.KeypairRef == nil {
				return nil
			}

			return []string{string(*resource.KeypairRef)}
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

	// No deletion guard for server references in scheduler hints, because they
	// are only used on creation for placement decisions
	sameHostServerRefDependency = dependency.NewDependency[*orcv1alpha1.ServerList, *orcv1alpha1.Server](
		"spec.resource.schedulerHints.sameHostServerRefs",
		func(server *orcv1alpha1.Server) []string {
			resource := server.Spec.Resource
			if resource == nil || resource.SchedulerHints == nil {
				return nil
			}

			refs := make([]string, 0, len(resource.SchedulerHints.SameHostServerRefs))
			for _, ref := range resource.SchedulerHints.SameHostServerRefs {
				refs = append(refs, string(ref))
			}
			return refs
		},
	)

	differentHostServerRefDependency = dependency.NewDependency[*orcv1alpha1.ServerList, *orcv1alpha1.Server](
		"spec.resource.schedulerHints.differentHostServerRefs",
		func(server *orcv1alpha1.Server) []string {
			resource := server.Spec.Resource
			if resource == nil || resource.SchedulerHints == nil {
				return nil
			}

			refs := make([]string, 0, len(resource.SchedulerHints.DifferentHostServerRefs))
			for _, ref := range resource.SchedulerHints.DifferentHostServerRefs {
				refs = append(refs, string(ref))
			}
			return refs
		},
	)
)

// SetupWithManager sets up the controller with the Manager.
func (c *serverReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
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
	keypairWatchEventHandler, err := keypairDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}
	bootVolumeWatchEventHandler, err := bootVolumeDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}
	sameHostServerRefWatchEventHandler, err := sameHostServerRefDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}
	differentHostServerRefWatchEventHandler, err := differentHostServerRefDependency.WatchEventHandler(log, k8sClient)
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
		Watches(&orcv1alpha1.Volume{}, bootVolumeWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Volume{})),
		).
		Watches(&orcv1alpha1.KeyPair{}, keypairWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.KeyPair{})),
		).
		Watches(&orcv1alpha1.Server{}, sameHostServerRefWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Server{})),
		).
		Watches(&orcv1alpha1.Server{}, differentHostServerRefWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Server{})),
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
		bootVolumeDependency.AddToManager(ctx, mgr),
		keypairDependency.AddToManager(ctx, mgr),
		sameHostServerRefDependency.AddToManager(ctx, mgr),
		differentHostServerRefDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, k8sClient, builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, k8sClient, c.scopeFactory, serverHelperFactory{}, serverStatusWriter{}, c.defaultResyncPeriod)
	return builder.Complete(&r)
}
