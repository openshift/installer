/*
Copyright 2024 The Kubernetes Authors.

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

	asocontainerservicev1 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	infrav1alpha "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// AgentPoolAdoptReconciler adopts ASO ManagedClustersAgentPool resources into a CAPI Cluster.
type AgentPoolAdoptReconciler struct {
	client.Client
}

// SetupWithManager sets up the controller with the Manager.
func (r *AgentPoolAdoptReconciler) SetupWithManager(_ context.Context, mgr ctrl.Manager, options controller.Options) error {
	_, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&asocontainerservicev1.ManagedClustersAgentPool{}).
		WithEventFilter(predicate.Funcs{
			UpdateFunc: func(ev event.UpdateEvent) bool {
				return ev.ObjectOld.GetAnnotations()[adoptAnnotation] != ev.ObjectNew.GetAnnotations()[adoptAnnotation]
			},
			DeleteFunc: func(_ event.DeleteEvent) bool { return false },
		}).
		Build(r)
	if err != nil {
		return err
	}

	return nil
}

// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools,verbs=create

// Reconcile reconciles an AzureASOManagedCluster.
func (r *AgentPoolAdoptReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, resultErr error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AgentPoolAdoptReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", "ManagedCluster"),
	)
	defer done()

	agentPool := &asocontainerservicev1.ManagedClustersAgentPool{}
	err := r.Get(ctx, req.NamespacedName, agentPool)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if agentPool.GetAnnotations()[adoptAnnotation] != adoptAnnotationValue {
		return ctrl.Result{}, nil
	}

	for _, owner := range agentPool.GetOwnerReferences() {
		if owner.APIVersion == infrav1alpha.GroupVersion.Identifier() &&
			owner.Kind == infrav1alpha.AzureASOManagedMachinePoolKind {
			return ctrl.Result{}, nil
		}
	}

	log.Info("adopting")

	namespace := agentPool.Namespace

	// filter down to what will be persisted in the AzureASOManagedMachinePool
	agentPool = &asocontainerservicev1.ManagedClustersAgentPool{
		TypeMeta: metav1.TypeMeta{
			APIVersion: asocontainerservicev1.GroupVersion.Identifier(),
			Kind:       "ManagedClustersAgentPool",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: agentPool.Name,
		},
		Spec: agentPool.Spec,
	}

	var replicas *int32
	if agentPool.Spec.Count != nil {
		replicas = ptr.To(int32(*agentPool.Spec.Count))
		agentPool.Spec.Count = nil
	}

	managedCluster := &asocontainerservicev1.ManagedCluster{}
	if agentPool.Owner() == nil {
		return ctrl.Result{}, fmt.Errorf("agent pool %s/%s has no owner", namespace, agentPool.Name)
	}
	managedClusterKey := client.ObjectKey{
		Namespace: namespace,
		Name:      agentPool.Owner().Name,
	}
	err = r.Get(ctx, managedClusterKey, managedCluster)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get ManagedCluster %s: %w", managedClusterKey, err)
	}
	var managedControlPlaneOwner *metav1.OwnerReference
	for _, owner := range managedCluster.GetOwnerReferences() {
		if owner.APIVersion == infrav1alpha.GroupVersion.Identifier() &&
			owner.Kind == infrav1alpha.AzureASOManagedControlPlaneKind &&
			owner.Name == agentPool.Owner().Name {
			managedControlPlaneOwner = ptr.To(owner)
			break
		}
	}
	if managedControlPlaneOwner == nil {
		return ctrl.Result{}, fmt.Errorf("ManagedCluster %s is not owned by any AzureASOManagedControlPlane", managedClusterKey)
	}
	asoManagedControlPlane := &infrav1alpha.AzureASOManagedControlPlane{}
	managedControlPlaneKey := client.ObjectKey{
		Namespace: namespace,
		Name:      managedControlPlaneOwner.Name,
	}
	err = r.Get(ctx, managedControlPlaneKey, asoManagedControlPlane)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get AzureASOManagedControlPlane %s: %w", managedControlPlaneKey, err)
	}
	clusterName := asoManagedControlPlane.Labels[clusterv1.ClusterNameLabel]

	asoManagedMachinePool := &infrav1alpha.AzureASOManagedMachinePool{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      agentPool.Name,
		},
		Spec: infrav1alpha.AzureASOManagedMachinePoolSpec{
			AzureASOManagedMachinePoolTemplateResourceSpec: infrav1alpha.AzureASOManagedMachinePoolTemplateResourceSpec{
				Resources: []runtime.RawExtension{
					{Object: agentPool},
				},
			},
		},
	}

	machinePool := &expv1.MachinePool{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      agentPool.Name,
		},
		Spec: expv1.MachinePoolSpec{
			ClusterName: clusterName,
			Replicas:    replicas,
			Template: clusterv1.MachineTemplateSpec{
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To(""),
					},
					ClusterName: clusterName,
					InfrastructureRef: corev1.ObjectReference{
						APIVersion: infrav1alpha.GroupVersion.Identifier(),
						Kind:       infrav1alpha.AzureASOManagedMachinePoolKind,
						Name:       asoManagedMachinePool.Name,
					},
				},
			},
		},
	}

	if ptr.Deref(agentPool.Spec.EnableAutoScaling, false) {
		machinePool.Annotations = map[string]string{
			clusterv1.ReplicasManagedByAnnotation: infrav1alpha.ReplicasManagedByAKS,
		}
	}

	err = r.Create(ctx, machinePool)
	if client.IgnoreAlreadyExists(err) != nil {
		return ctrl.Result{}, err
	}

	err = r.Create(ctx, asoManagedMachinePool)
	if client.IgnoreAlreadyExists(err) != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
