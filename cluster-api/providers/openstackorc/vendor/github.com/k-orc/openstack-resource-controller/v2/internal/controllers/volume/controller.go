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

package volume

import (
	"context"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	applyconfigmetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/reconciler"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	"github.com/k-orc/openstack-resource-controller/v2/internal/scope"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/applyconfigs"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/credentials"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	orcstrings "github.com/k-orc/openstack-resource-controller/v2/internal/util/strings"
	applyconfigv1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/pkg/predicates"
)

const controllerName = "volume"

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=volumes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=volumes/status,verbs=get;update;patch

type volumeReconcilerConstructor struct {
	scopeFactory scope.Factory
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return volumeReconcilerConstructor{scopeFactory: scopeFactory}
}

func (volumeReconcilerConstructor) GetName() string {
	return controllerName
}

var volumetypeDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.VolumeList, *orcv1alpha1.VolumeType](
	"spec.resource.volumeTypeRef",
	func(volume *orcv1alpha1.Volume) []string {
		resource := volume.Spec.Resource
		if resource == nil || resource.VolumeTypeRef == nil {
			return nil
		}
		return []string{string(*resource.VolumeTypeRef)}
	},
	finalizer, externalObjectFieldOwner,
)

// serverToVolumeMapFunc creates a mapping function that reconciles volumes when:
// - a volume ID appears in server status but the volume doesn't have attachment info for that server
// - a volume has attachment info for a server, but the server no longer lists that volume
func serverToVolumeMapFunc(ctx context.Context, k8sClient client.Client) handler.MapFunc {
	log := ctrl.LoggerFrom(ctx)

	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		server, ok := obj.(*orcv1alpha1.Server)
		if !ok {
			log.Info("serverToVolumeMapFunc got unexpected object type",
				"got", fmt.Sprintf("%T", obj),
				"expected", fmt.Sprintf("%T", &orcv1alpha1.Server{}))
			return nil
		}

		// Get the server's ID and volume IDs from status
		serverStatus := server.Status.Resource
		if serverStatus == nil {
			return nil
		}

		serverID := ptr.Deref(server.Status.ID, "")
		if serverID == "" {
			// Server doesn't have an ID yet, nothing to reconcile
			return nil
		}

		// Build a set of volume IDs attached to this server according to server status
		serverVolumeIDs := make(map[string]struct{})
		for i := range serverStatus.Volumes {
			volumeID := serverStatus.Volumes[i].ID
			if volumeID != "" {
				serverVolumeIDs[volumeID] = struct{}{}
			}
		}

		// List all volumes in the same namespace
		volumeList := &orcv1alpha1.VolumeList{}
		if err := k8sClient.List(ctx, volumeList, client.InNamespace(server.Namespace)); err != nil {
			log.Error(err, "failed to list volumes", "namespace", server.Namespace)
			return nil
		}

		requests := []reconcile.Request{}

		for i := range volumeList.Items {
			volume := &volumeList.Items[i]
			volumeStatus := volume.Status.Resource
			if volumeStatus == nil {
				continue
			}

			volumeID := ptr.Deref(volume.Status.ID, "")
			if volumeID == "" {
				continue
			}

			shouldReconcile := false
			var reason string

			// Volume ID is in server's status, but volume doesn't have attachment info for this server
			if _, volumeInServerStatus := serverVolumeIDs[volumeID]; volumeInServerStatus {
				hasAttachment := false
				for j := range volumeStatus.Attachments {
					if volumeStatus.Attachments[j].ServerID == serverID {
						hasAttachment = true
						break
					}
				}
				if !hasAttachment {
					shouldReconcile = true
					reason = "Server attached volume but volume status not updated"
					log.V(logging.Verbose).Info("volume needs reconciliation: listed in server status but no attachment info",
						"volume", client.ObjectKeyFromObject(volume),
						"server", client.ObjectKeyFromObject(server))
				}
			}

			// Volume has attachment info for this server, but server no longer lists this volume
			if !shouldReconcile {
				for j := range volumeStatus.Attachments {
					if volumeStatus.Attachments[j].ServerID == serverID {
						// Volume thinks it's attached to this server
						if _, stillAttached := serverVolumeIDs[volumeID]; !stillAttached {
							shouldReconcile = true
							reason = "Server detached volume but volume status not updated"
							log.V(logging.Verbose).Info("volume needs reconciliation: has attachment info but not in server status",
								"volume", client.ObjectKeyFromObject(volume),
								"server", client.ObjectKeyFromObject(server))
						}
						break
					}
				}
			}

			if shouldReconcile {
				// Update the volume's Progressing condition to trigger reconciliation
				volumeApply := applyconfigv1.Volume(volume.Name, volume.Namespace).
					WithStatus(applyconfigv1.VolumeStatus().
						WithConditions(applyconfigmetav1.Condition().
							WithType(orcv1alpha1.ConditionProgressing).
							WithStatus(metav1.ConditionTrue).
							WithReason(orcv1alpha1.ConditionReasonProgressing).
							WithMessage(reason).
							WithObservedGeneration(volume.Generation),
						),
					)

				if err := k8sClient.Status().Patch(ctx, volume, applyconfigs.Patch(types.ApplyPatchType, volumeApply), orcstrings.GetSSAFieldOwner(controllerName), client.ForceOwnership); err != nil {
					log.Error(err, "failed to update volume progressing status",
						"volume", client.ObjectKeyFromObject(volume),
						"server", client.ObjectKeyFromObject(server))
				}

				// Also add to reconcile requests
				requests = append(requests, reconcile.Request{
					NamespacedName: client.ObjectKeyFromObject(volume),
				})
			}
		}

		return requests
	}
}

// SetupWithManager sets up the controller with the Manager.
func (c volumeReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)
	k8sClient := mgr.GetClient()

	volumetypeWatchEventHandler, err := volumetypeDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		Watches(&orcv1alpha1.VolumeType{}, volumetypeWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.VolumeType{})),
		).
		Watches(&orcv1alpha1.Server{}, handler.EnqueueRequestsFromMapFunc(serverToVolumeMapFunc(ctx, k8sClient)),
			builder.WithPredicates(predicates.NewServerVolumesChanged(log)),
		).
		For(&orcv1alpha1.Volume{})

	if err := errors.Join(
		volumetypeDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, mgr.GetClient(), builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, mgr.GetClient(), c.scopeFactory, volumeHelperFactory{}, volumeStatusWriter{})
	return builder.Complete(&r)
}
