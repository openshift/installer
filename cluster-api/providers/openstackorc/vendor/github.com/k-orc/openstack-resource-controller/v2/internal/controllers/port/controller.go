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

package port

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
	applyconfigv1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/pkg/predicates"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/reconciler"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	"github.com/k-orc/openstack-resource-controller/v2/internal/scope"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/applyconfigs"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/credentials"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	orcstrings "github.com/k-orc/openstack-resource-controller/v2/internal/util/strings"
)

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=ports,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=ports/status,verbs=get;update;patch

const controllerName = "port"

var (
	networkDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.PortList, *orcv1alpha1.Network](
		"spec.resource.networkRef",
		func(port *orcv1alpha1.Port) []string {
			resource := port.Spec.Resource
			if resource == nil {
				return nil
			}
			return []string{string(resource.NetworkRef)}
		},
		finalizer, externalObjectFieldOwner,
	)

	networkImportDependency = dependency.NewDependency[*orcv1alpha1.PortList, *orcv1alpha1.Network](
		"spec.import.filter.networkRef",
		func(port *orcv1alpha1.Port) []string {
			resource := port.Spec.Import
			if resource == nil || resource.Filter == nil {
				return nil
			}
			return []string{string(resource.Filter.NetworkRef)}
		},
	)

	subnetDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.PortList, *orcv1alpha1.Subnet](
		"spec.resource.addresses[].subnetRef",
		func(port *orcv1alpha1.Port) []string {
			if port.Spec.Resource == nil {
				return nil
			}
			subnets := make([]string, len(port.Spec.Resource.Addresses))
			for i := range port.Spec.Resource.Addresses {
				subnets[i] = string(port.Spec.Resource.Addresses[i].SubnetRef)
			}
			return subnets
		},
		finalizer, externalObjectFieldOwner,
	)

	securityGroupDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.PortList, *orcv1alpha1.SecurityGroup](
		"spec.resource.securityGroupRefs",
		func(port *orcv1alpha1.Port) []string {
			if port.Spec.Resource == nil {
				return nil
			}
			securityGroups := make([]string, len(port.Spec.Resource.SecurityGroupRefs))
			for i := range port.Spec.Resource.SecurityGroupRefs {
				securityGroups[i] = string(port.Spec.Resource.SecurityGroupRefs[i])
			}
			return securityGroups
		},
		finalizer, externalObjectFieldOwner,
	)

	projectDependency = dependency.NewDeletionGuardDependency[*orcv1alpha1.PortList, *orcv1alpha1.Project](
		"spec.resource.projectRef",
		func(port *orcv1alpha1.Port) []string {
			resource := port.Spec.Resource
			if resource == nil || resource.ProjectRef == nil {
				return nil
			}
			return []string{string(*resource.ProjectRef)}
		},
		finalizer, externalObjectFieldOwner,
	)

	projectImportDependency = dependency.NewDependency[*orcv1alpha1.PortList, *orcv1alpha1.Project](
		"spec.import.filter.projectRef",
		func(port *orcv1alpha1.Port) []string {
			resource := port.Spec.Import
			if resource == nil || resource.Filter == nil || resource.Filter.ProjectRef == nil {
				return nil
			}
			return []string{string(*resource.Filter.ProjectRef)}
		},
	)
)

// serverToPortMapFunc creates a mapping function that reconciles ports when:
// - a port ID appears in server status but the port doesn't have attachment info for that server
// - a port has attachment info for a server, but the server no longer lists that port
func serverToPortMapFunc(ctx context.Context, k8sClient client.Client) handler.MapFunc {
	log := ctrl.LoggerFrom(ctx)

	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		server, ok := obj.(*orcv1alpha1.Server)
		if !ok {
			log.Info("serverToPortMapFunc got unexpected object type",
				"got", fmt.Sprintf("%T", obj),
				"expected", fmt.Sprintf("%T", &orcv1alpha1.Server{}))
			return nil
		}

		// Get the server's ID and port IDs from status
		serverStatus := server.Status.Resource
		if serverStatus == nil {
			return nil
		}

		serverID := ptr.Deref(server.Status.ID, "")
		if serverID == "" {
			// Server doesn't have an ID yet, nothing to reconcile
			return nil
		}

		// Build a set of port IDs attached to this server according to server status
		serverPortIDs := make(map[string]struct{})
		for i := range serverStatus.Interfaces {
			portID := serverStatus.Interfaces[i].PortID
			if portID != "" {
				serverPortIDs[portID] = struct{}{}
			}
		}

		// List all ports in the same namespace
		portList := &orcv1alpha1.PortList{}
		if err := k8sClient.List(ctx, portList, client.InNamespace(server.Namespace)); err != nil {
			log.Error(err, "failed to list ports", "namespace", server.Namespace)
			return nil
		}

		requests := []reconcile.Request{}

		for i := range portList.Items {
			port := &portList.Items[i]
			portStatus := port.Status.Resource
			if portStatus == nil {
				continue
			}

			portID := ptr.Deref(port.Status.ID, "")
			if portID == "" {
				continue
			}

			shouldReconcile := false
			var reason string

			// Port ID is in server's status, but port doesn't have attachment info for this server
			if _, portInServerStatus := serverPortIDs[portID]; portInServerStatus {
				if portStatus.DeviceID != serverID {
					shouldReconcile = true
					reason = "Server attached port but port status not updated"
					log.V(logging.Verbose).Info("port needs reconciliation: listed in server status but deviceID not set",
						"port", client.ObjectKeyFromObject(port),
						"server", client.ObjectKeyFromObject(server))
				}
			}

			// Port has attachment info for this server, but server no longer lists this port
			if !shouldReconcile {
				if portStatus.DeviceID == serverID {
					// Port thinks it's attached to this server
					if _, stillAttached := serverPortIDs[portID]; !stillAttached {
						shouldReconcile = true
						reason = "Server detached port but port status not updated"
						log.V(logging.Verbose).Info("port needs reconciliation: has deviceID set but not in server status",
							"port", client.ObjectKeyFromObject(port),
							"server", client.ObjectKeyFromObject(server))
					}
				}
			}

			if shouldReconcile {
				// Update the port's Progressing condition to trigger reconciliation
				portApply := applyconfigv1.Port(port.Name, port.Namespace).
					WithStatus(applyconfigv1.PortStatus().
						WithConditions(applyconfigmetav1.Condition().
							WithType(orcv1alpha1.ConditionProgressing).
							WithStatus(metav1.ConditionTrue).
							WithReason(orcv1alpha1.ConditionReasonProgressing).
							WithMessage(reason).
							WithObservedGeneration(port.Generation),
						),
					)

				if err := k8sClient.Status().Patch(ctx, port, applyconfigs.Patch(types.ApplyPatchType, portApply), orcstrings.GetSSAFieldOwner(controllerName), client.ForceOwnership); err != nil {
					log.Error(err, "failed to update port progressing status",
						"port", client.ObjectKeyFromObject(port),
						"server", client.ObjectKeyFromObject(server))
				}

				// Also add to reconcile requests
				requests = append(requests, reconcile.Request{
					NamespacedName: client.ObjectKeyFromObject(port),
				})
			}
		}

		return requests
	}
}

type portReconcilerConstructor struct {
	scopeFactory scope.Factory
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return portReconcilerConstructor{scopeFactory: scopeFactory}
}

func (portReconcilerConstructor) GetName() string {
	return controllerName
}

// SetupWithManager sets up the controller with the Manager.
func (c portReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := mgr.GetLogger().WithValues("controller", controllerName)
	k8sClient := mgr.GetClient()

	networkWatchEventHandler, err := networkDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	networkImportWatchEventHandler, err := networkImportDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	subnetWatchEventHandler, err := subnetDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	securityGroupWatchEventHandler, err := securityGroupDependency.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

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
		For(&orcv1alpha1.Port{}).
		Watches(&orcv1alpha1.Network{}, networkWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Network{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Network{}, networkImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Network{})),
		).
		Watches(&orcv1alpha1.Subnet{}, subnetWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Subnet{})),
		).
		Watches(&orcv1alpha1.SecurityGroup{}, securityGroupWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.SecurityGroup{})),
		).
		Watches(&orcv1alpha1.Project{}, projectWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Project{})),
		).
		// A second watch is necessary because we need a different handler that omits deletion guards
		Watches(&orcv1alpha1.Project{}, projectImportWatchEventHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Project{})),
		).
		Watches(&orcv1alpha1.Server{}, handler.EnqueueRequestsFromMapFunc(serverToPortMapFunc(ctx, k8sClient)),
			builder.WithPredicates(predicates.NewServerInterfacesChanged(log)),
		)

	if err := errors.Join(
		networkDependency.AddToManager(ctx, mgr),
		networkImportDependency.AddToManager(ctx, mgr),
		subnetDependency.AddToManager(ctx, mgr),
		securityGroupDependency.AddToManager(ctx, mgr),
		projectDependency.AddToManager(ctx, mgr),
		projectImportDependency.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, k8sClient, builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, k8sClient, c.scopeFactory, portHelperFactory{}, portStatusWriter{})
	return builder.Complete(&r)
}
