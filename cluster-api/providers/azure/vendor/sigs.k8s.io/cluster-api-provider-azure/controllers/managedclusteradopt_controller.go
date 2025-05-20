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
	asoresourcesv1 "github.com/Azure/azure-service-operator/v2/api/resources/v1api20200601"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	infrav1alpha "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const (
	adoptAnnotation      = "sigs.k8s.io/cluster-api-provider-azure-adopt"
	adoptAnnotationValue = "true"
)

// ManagedClusterAdoptReconciler adopts ASO ManagedCluster resources into a CAPI Cluster.
type ManagedClusterAdoptReconciler struct {
	client.Client
}

// SetupWithManager sets up the controller with the Manager.
func (r *ManagedClusterAdoptReconciler) SetupWithManager(_ context.Context, mgr ctrl.Manager, options controller.Options) error {
	_, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&asocontainerservicev1.ManagedCluster{}).
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

// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters,verbs=create

// Reconcile reconciles an AzureASOManagedCluster.
func (r *ManagedClusterAdoptReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, resultErr error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.ManagedClusterAdoptReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", "ManagedCluster"),
	)
	defer done()

	managedCluster := &asocontainerservicev1.ManagedCluster{}
	err := r.Get(ctx, req.NamespacedName, managedCluster)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if managedCluster.GetAnnotations()[adoptAnnotation] != adoptAnnotationValue {
		return ctrl.Result{}, nil
	}

	for _, owner := range managedCluster.GetOwnerReferences() {
		if owner.APIVersion == infrav1alpha.GroupVersion.Identifier() &&
			owner.Kind == infrav1alpha.AzureASOManagedControlPlaneKind {
			return ctrl.Result{}, nil
		}
	}

	log.Info("adopting")

	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: managedCluster.Namespace,
			Name:      managedCluster.Name,
		},
		Spec: clusterv1.ClusterSpec{
			InfrastructureRef: &corev1.ObjectReference{
				APIVersion: infrav1alpha.GroupVersion.Identifier(),
				Kind:       infrav1alpha.AzureASOManagedClusterKind,
				Name:       managedCluster.Name,
			},
			ControlPlaneRef: &corev1.ObjectReference{
				APIVersion: infrav1alpha.GroupVersion.Identifier(),
				Kind:       infrav1alpha.AzureASOManagedControlPlaneKind,
				Name:       managedCluster.Name,
			},
		},
	}
	err = r.Create(ctx, cluster)
	if client.IgnoreAlreadyExists(err) != nil {
		return ctrl.Result{}, err
	}

	resourceGroup := &asoresourcesv1.ResourceGroup{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: managedCluster.Namespace,
			Name:      managedCluster.Owner().Name,
		},
	}

	err = r.Get(ctx, client.ObjectKeyFromObject(resourceGroup), resourceGroup)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting ResourceGroup %s: %w", client.ObjectKeyFromObject(resourceGroup), err)
	}

	// filter down to what will be persisted in the AzureASOManagedCluster
	resourceGroup = &asoresourcesv1.ResourceGroup{
		TypeMeta: metav1.TypeMeta{
			APIVersion: asoresourcesv1.GroupVersion.Identifier(),
			Kind:       "ResourceGroup",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceGroup.Name,
		},
		Spec: resourceGroup.Spec,
	}

	asoManagedCluster := &infrav1alpha.AzureASOManagedCluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: managedCluster.Namespace,
			Name:      managedCluster.Name,
		},
		Spec: infrav1alpha.AzureASOManagedClusterSpec{
			AzureASOManagedClusterTemplateResourceSpec: infrav1alpha.AzureASOManagedClusterTemplateResourceSpec{
				Resources: []runtime.RawExtension{
					{Object: resourceGroup},
				},
			},
		},
	}
	err = r.Create(ctx, asoManagedCluster)
	if client.IgnoreAlreadyExists(err) != nil {
		return ctrl.Result{}, err
	}

	// agent pools are defined by AzureASOManagedMachinePools. Remove them from this resource.
	managedClusterBefore := managedCluster.DeepCopy()
	managedCluster.Spec.AgentPoolProfiles = nil
	err = r.Patch(ctx, managedCluster, client.MergeFrom(managedClusterBefore))
	if err != nil {
		return ctrl.Result{}, err
	}

	managedCluster = &asocontainerservicev1.ManagedCluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: asocontainerservicev1.GroupVersion.Identifier(),
			Kind:       "ManagedCluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: managedCluster.Name,
		},
		Spec: managedCluster.Spec,
	}

	asoManagedControlPlane := &infrav1alpha.AzureASOManagedControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      managedCluster.Name,
		},
		Spec: infrav1alpha.AzureASOManagedControlPlaneSpec{
			AzureASOManagedControlPlaneTemplateResourceSpec: infrav1alpha.AzureASOManagedControlPlaneTemplateResourceSpec{
				Resources: []runtime.RawExtension{
					{Object: managedCluster},
				},
			},
		},
	}
	err = r.Create(ctx, asoManagedControlPlane)
	if client.IgnoreAlreadyExists(err) != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
