/*
Copyright 2020 The Kubernetes Authors.

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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/controllers"
	infrav1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	kubeadmv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// AzureClusterToAzureMachinePoolsMapper creates a mapping handler to transform AzureClusters into AzureMachinePools. The transform
// requires AzureCluster to map to the owning Cluster, then from the Cluster, collect the MachinePools belonging to the cluster,
// then finally projecting the infrastructure reference to the AzureMachinePool.
func AzureClusterToAzureMachinePoolsMapper(ctx context.Context, c client.Client, scheme *runtime.Scheme, log logr.Logger) (handler.MapFunc, error) {
	gvk, err := apiutil.GVKForObject(new(infrav1exp.AzureMachinePool), scheme)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find GVK for AzureMachinePool")
	}

	return func(ctx context.Context, o client.Object) []ctrl.Request {
		ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultMappingTimeout)
		defer cancel()

		azCluster, ok := o.(*infrav1.AzureCluster)
		if !ok {
			log.Error(errors.Errorf("expected an AzureCluster, got %T instead", o), "failed to map AzureCluster")
			return nil
		}

		log := log.WithValues("AzureCluster", azCluster.Name, "Namespace", azCluster.Namespace)

		// Don't handle deleted AzureClusters
		if !azCluster.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(4).Info("AzureCluster has a deletion timestamp, skipping mapping.")
			return nil
		}

		clusterName, ok := controllers.GetOwnerClusterName(azCluster.ObjectMeta)
		if !ok {
			log.V(4).Info("unable to get the owner cluster")
			return nil
		}

		machineList := &expv1.MachinePoolList{}
		machineList.SetGroupVersionKind(gvk)
		// list all of the requested objects within the cluster namespace with the cluster name label
		if err := c.List(ctx, machineList, client.InNamespace(azCluster.Namespace), client.MatchingLabels{clusterv1.ClusterNameLabel: clusterName}); err != nil {
			log.V(4).Info(fmt.Sprintf("unable to list machine pools in cluster %s", clusterName))
			return nil
		}

		mapFunc := MachinePoolToInfrastructureMapFunc(gvk, log)
		var results []ctrl.Request
		for _, machine := range machineList.Items {
			m := machine
			azureMachines := mapFunc(ctx, &m)
			results = append(results, azureMachines...)
		}

		return results
	}, nil
}

// AzureManagedControlPlaneToAzureMachinePoolsMapper creates a mapping handler to transform AzureManagedControlPlanes into AzureMachinePools. The transform
// requires AzureManagedControlPlane to map to the owning Cluster, then from the Cluster, collect the MachinePools belonging to the cluster,
// then finally projecting the infrastructure reference to the AzureMachinePool.
func AzureManagedControlPlaneToAzureMachinePoolsMapper(ctx context.Context, c client.Client, scheme *runtime.Scheme, log logr.Logger) (handler.MapFunc, error) {
	gvk, err := apiutil.GVKForObject(new(infrav1exp.AzureMachinePool), scheme)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find GVK for AzureMachinePool")
	}

	return func(ctx context.Context, o client.Object) []ctrl.Request {
		ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultMappingTimeout)
		defer cancel()

		azControlPlane, ok := o.(*infrav1.AzureManagedControlPlane)
		if !ok {
			log.Error(errors.Errorf("expected an AzureManagedControlPlane, got %T instead", o), "failed to map AzureManagedControlPlane")
			return nil
		}

		log := log.WithValues("AzureManagedControlPlane", azControlPlane.Name, "Namespace", azControlPlane.Namespace)

		// Don't handle deleted AzureManagedControlPlane
		if !azControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(4).Info("AzureManagedControlPlane has a deletion timestamp, skipping mapping.")
			return nil
		}

		clusterName, ok := controllers.GetOwnerClusterName(azControlPlane.ObjectMeta)
		if !ok {
			log.V(4).Info("unable to get the owner cluster")
			return nil
		}

		machineList := &expv1.MachinePoolList{}
		machineList.SetGroupVersionKind(gvk)
		// list all of the requested objects within the cluster namespace with the cluster name label
		if err := c.List(ctx, machineList, client.InNamespace(azControlPlane.Namespace), client.MatchingLabels{clusterv1.ClusterNameLabel: clusterName}); err != nil {
			log.V(4).Info(fmt.Sprintf("unable to list machine pools in cluster %s", clusterName))
			return nil
		}

		mapFunc := MachinePoolToInfrastructureMapFunc(gvk, log)
		var results []ctrl.Request
		for _, machine := range machineList.Items {
			m := machine
			azureMachines := mapFunc(ctx, &m)
			results = append(results, azureMachines...)
		}

		return results
	}, nil
}

// AzureMachinePoolMachineMapper creates a mapping handler to transform AzureMachinePoolMachine to AzureMachinePools.
func AzureMachinePoolMachineMapper(scheme *runtime.Scheme, log logr.Logger) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		gvk, err := apiutil.GVKForObject(new(infrav1exp.AzureMachinePool), scheme)
		if err != nil {
			log.Error(errors.WithStack(err), "failed to find GVK for AzureMachinePool")
			return nil
		}

		azureMachinePoolMachine, ok := o.(*infrav1exp.AzureMachinePoolMachine)
		if !ok {
			log.Error(errors.Errorf("expected an AzureCluster, got %T instead", o), "failed to map AzureMachinePoolMachine")
			return nil
		}

		log := log.WithValues("AzureMachinePoolMachine", azureMachinePoolMachine.Name, "Namespace", azureMachinePoolMachine.Namespace)
		for _, ref := range azureMachinePoolMachine.OwnerReferences {
			if ref.Kind != gvk.Kind {
				continue
			}

			gv, err := schema.ParseGroupVersion(ref.APIVersion)
			if err != nil {
				log.Error(errors.WithStack(err), "unable to parse group version", "APIVersion", ref.APIVersion)
				return nil
			}

			if gv.Group == gvk.Group {
				return []ctrl.Request{
					{
						NamespacedName: types.NamespacedName{
							Name:      ref.Name,
							Namespace: azureMachinePoolMachine.Namespace,
						},
					},
				}
			}
		}

		return nil
	}
}

// MachinePoolToInfrastructureMapFunc returns a handler.MapFunc that watches for
// MachinePool events and returns reconciliation requests for an infrastructure provider object.
func MachinePoolToInfrastructureMapFunc(gvk schema.GroupVersionKind, log logr.Logger) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []reconcile.Request {
		m, ok := o.(*expv1.MachinePool)
		if !ok {
			log.V(4).Info("attempt to map incorrect type", "type", fmt.Sprintf("%T", o))
			return nil
		}

		gk := gvk.GroupKind()
		ref := m.Spec.Template.Spec.InfrastructureRef
		// Return early if the GroupKind doesn't match what we expect.
		infraGK := ref.GroupVersionKind().GroupKind()
		if gk != infraGK {
			log.V(4).Info("gk does not match", "gk", gk, "infraGK", infraGK)
			return nil
		}

		return []reconcile.Request{
			{
				NamespacedName: client.ObjectKey{
					Namespace: m.Namespace,
					Name:      ref.Name,
				},
			},
		}
	}
}

// AzureClusterToAzureMachinePoolsFunc is a handler.MapFunc to be used to enqueue
// requests for reconciliation of AzureMachinePools.
func AzureClusterToAzureMachinePoolsFunc(ctx context.Context, c client.Client, log logr.Logger) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []reconcile.Request {
		ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultMappingTimeout)
		defer cancel()

		ac, ok := o.(*infrav1.AzureCluster)
		if !ok {
			log.Error(errors.Errorf("expected a AzureCluster but got a %T", o), "failed to get AzureCluster")
			return nil
		}
		logWithValues := log.WithValues("AzureCluster", ac.Name, "Namespace", ac.Namespace)

		cluster, err := util.GetOwnerCluster(ctx, c, ac.ObjectMeta)
		switch {
		case apierrors.IsNotFound(err) || cluster == nil:
			logWithValues.V(4).Info("owning cluster not found")
			return nil
		case err != nil:
			logWithValues.Error(err, "failed to get owning cluster")
			return nil
		}

		labels := map[string]string{clusterv1.ClusterNameLabel: cluster.Name}
		ampl := &infrav1exp.AzureMachinePoolList{}
		if err := c.List(ctx, ampl, client.InNamespace(ac.Namespace), client.MatchingLabels(labels)); err != nil {
			logWithValues.Error(err, "failed to list AzureMachinePools")
			return nil
		}

		var result []reconcile.Request
		for _, m := range ampl.Items {
			result = append(result, reconcile.Request{
				NamespacedName: client.ObjectKey{
					Namespace: m.Namespace,
					Name:      m.Name,
				},
			})
		}

		return result
	}
}

// AzureMachinePoolToAzureMachinePoolMachines maps an AzureMachinePool to its child AzureMachinePoolMachines through
// Cluster and MachinePool labels.
func AzureMachinePoolToAzureMachinePoolMachines(ctx context.Context, c client.Client, log logr.Logger) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []reconcile.Request {
		ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultMappingTimeout)
		defer cancel()

		amp, ok := o.(*infrav1exp.AzureMachinePool)
		if !ok {
			log.Error(errors.Errorf("expected a AzureMachinePool but got a %T", o), "failed to get AzureMachinePool")
			return nil
		}
		logWithValues := log.WithValues("AzureMachinePool", amp.Name, "Namespace", amp.Namespace)

		labels := map[string]string{
			clusterv1.ClusterNameLabel:      amp.Labels[clusterv1.ClusterNameLabel],
			infrav1exp.MachinePoolNameLabel: amp.Name,
		}
		ampml := &infrav1exp.AzureMachinePoolMachineList{}
		if err := c.List(ctx, ampml, client.InNamespace(amp.Namespace), client.MatchingLabels(labels)); err != nil {
			logWithValues.Error(err, "failed to list AzureMachinePoolMachines")
			return nil
		}

		logWithValues.Info("mapping from AzureMachinePool", "count", len(ampml.Items))
		var result []reconcile.Request
		for _, m := range ampml.Items {
			result = append(result, reconcile.Request{
				NamespacedName: client.ObjectKey{
					Namespace: m.Namespace,
					Name:      m.Name,
				},
			})
		}

		return result
	}
}

// MachinePoolModelHasChanged predicates any events based on changes to the AzureMachinePool model.
func MachinePoolModelHasChanged(logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			log := logger.WithValues("predicate", "MachinePoolModelHasChanged", "eventType", "update")

			oldAmp, ok := e.ObjectOld.(*infrav1exp.AzureMachinePool)
			if !ok {
				log.V(4).Info("Expected AzureMachinePool", "type", e.ObjectOld.GetObjectKind().GroupVersionKind().String())
				return false
			}
			log = log.WithValues("namespace", oldAmp.Namespace, "azureMachinePool", oldAmp.Name)

			newAmp := e.ObjectNew.(*infrav1exp.AzureMachinePool)

			// if any of these are not equal, run the update
			shouldUpdate := !cmp.Equal(oldAmp.Spec.Identity, newAmp.Spec.Identity) ||
				!cmp.Equal(oldAmp.Spec.Template, newAmp.Spec.Template) ||
				!cmp.Equal(oldAmp.Spec.UserAssignedIdentities, newAmp.Spec.UserAssignedIdentities) ||
				!cmp.Equal(oldAmp.Status.ProvisioningState, newAmp.Status.ProvisioningState)

			// if shouldUpdate {
			log.Info("machine pool predicate", "shouldUpdate", shouldUpdate)
			//}
			return shouldUpdate
		},
		CreateFunc:  func(e event.CreateEvent) bool { return false },
		DeleteFunc:  func(e event.DeleteEvent) bool { return false },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}
}

// MachinePoolMachineHasStateOrVersionChange predicates any events based on changes to the AzureMachinePoolMachine status
// relevant for the AzureMachinePool controller.
func MachinePoolMachineHasStateOrVersionChange(logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			log := logger.WithValues("predicate", "MachinePoolModelHasChanged", "eventType", "update")

			oldAmp, ok := e.ObjectOld.(*infrav1exp.AzureMachinePoolMachine)
			if !ok {
				log.V(4).Info("Expected AzureMachinePoolMachine", "type", e.ObjectOld.GetObjectKind().GroupVersionKind().String())
				return false
			}
			log = log.WithValues("namespace", oldAmp.Namespace, "machinePoolMachine", oldAmp.Name)

			newAmp := e.ObjectNew.(*infrav1exp.AzureMachinePoolMachine)

			// if any of these are not equal, run the update
			shouldUpdate := oldAmp.Status.LatestModelApplied != newAmp.Status.LatestModelApplied ||
				oldAmp.Status.Version != newAmp.Status.Version ||
				oldAmp.Status.ProvisioningState != newAmp.Status.ProvisioningState ||
				oldAmp.Status.Ready != newAmp.Status.Ready

			if shouldUpdate {
				log.Info("machine pool machine predicate", "shouldUpdate", shouldUpdate)
			}
			return shouldUpdate
		},
		CreateFunc:  func(e event.CreateEvent) bool { return false },
		DeleteFunc:  func(e event.DeleteEvent) bool { return false },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}
}

// KubeadmConfigToInfrastructureMapFunc returns a handler.ToRequestsFunc that watches for KubeadmConfig events and returns.
func KubeadmConfigToInfrastructureMapFunc(ctx context.Context, c client.Client, log logr.Logger) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []reconcile.Request {
		ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultMappingTimeout)
		defer cancel()

		kc, ok := o.(*kubeadmv1.KubeadmConfig)
		if !ok {
			log.V(4).Info("attempt to map incorrect type", "type", fmt.Sprintf("%T", o))
			return nil
		}

		mpKey := client.ObjectKey{
			Namespace: kc.Namespace,
			Name:      kc.Name,
		}

		// fetch MachinePool to get reference
		mp := &expv1.MachinePool{}
		if err := c.Get(ctx, mpKey, mp); err != nil {
			if !apierrors.IsNotFound(err) {
				log.Error(err, "failed to fetch MachinePool for KubeadmConfig")
			}
			return []reconcile.Request{}
		}

		ref := mp.Spec.Template.Spec.Bootstrap.ConfigRef
		if ref == nil {
			log.V(4).Info("fetched MachinePool has no Bootstrap.ConfigRef")
			return []reconcile.Request{}
		}
		sameKind := ref.Kind != o.GetObjectKind().GroupVersionKind().Kind
		sameName := ref.Name == kc.Name
		sameNamespace := ref.Namespace == kc.Namespace
		if !sameKind || !sameName || !sameNamespace {
			log.V(4).Info("Bootstrap.ConfigRef does not match",
				"sameKind", sameKind,
				"ref kind", ref.Kind,
				"other kind", o.GetObjectKind().GroupVersionKind().Kind,
				"sameName", sameName,
				"sameNamespace", sameNamespace,
			)
			return []reconcile.Request{}
		}

		key := client.ObjectKey{
			Namespace: kc.Namespace,
			Name:      kc.Name,
		}
		log.V(4).Info("adding KubeadmConfig to watch", "key", key)

		return []reconcile.Request{
			{
				NamespacedName: key,
			},
		}
	}
}
