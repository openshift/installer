/*
Copyright 2021 The Kubernetes Authors.

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

// Package controlleridentitycreator provides a way to reconcile AWSClusterControllerIdentity instance.
package controlleridentitycreator

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// AWSControllerIdentityReconciler reconciles a AWSClusterControllerIdentity object.
type AWSControllerIdentityReconciler struct {
	client.Client
	Log              logr.Logger
	Endpoints        []scope.ServiceEndpoint
	WatchFilterValue string
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsclustercontrolleridentities,verbs=get;list;watch;create

func (r *AWSControllerIdentityReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.FromContext(ctx)

	var identityRef *infrav1.AWSIdentityReference

	// Fetch the AWSCluster instance
	awsCluster := &infrav1.AWSCluster{}
	clusterFound := true
	err := r.Get(ctx, req.NamespacedName, awsCluster)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}
		log.Trace("AWSCluster not found, trying AWSManagedControlPlane")
		clusterFound = false
	} else {
		log.Trace("Found identityRef on AWSCluster")
		identityRef = awsCluster.Spec.IdentityRef
	}

	// If AWSCluster is not found, check if AWSManagedControlPlane is used.
	if !clusterFound && feature.Gates.Enabled(feature.EKS) {
		awsControlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
		if err := r.Client.Get(ctx, req.NamespacedName, awsControlPlane); err != nil {
			if apierrors.IsNotFound(err) {
				log.Trace("AWSManagedMachinePool not found, no identityRef so no action taken")
				return ctrl.Result{}, nil
			}
			return reconcile.Result{}, err
		}
		log.Trace("Found identityRef on AWSManagedControlPlane")
		identityRef = awsControlPlane.Spec.IdentityRef
	}

	log = log.WithValues("cluster", klog.KObj(awsCluster))
	if identityRef == nil {
		log.Info("IdentityRef is nil, skipping reconciliation")
		return ctrl.Result{Requeue: true}, nil
	}

	// If identity type is not AWSClusterControllerIdentity, then no need to create AWSClusterControllerIdentity singleton.
	if identityRef.Kind == infrav1.ClusterRoleIdentityKind ||
		identityRef.Kind == infrav1.ClusterStaticIdentityKind {
		log.Trace("Cluster does not use AWSClusterControllerIdentity as identityRef, skipping new instance creation")
		return ctrl.Result{}, nil
	}

	// Fetch the AWSClusterControllerIdentity instance
	controllerIdentity := &infrav1.AWSClusterControllerIdentity{}
	err = r.Get(ctx, types.NamespacedName{Name: infrav1.AWSClusterControllerIdentityName}, controllerIdentity)
	// If AWSClusterControllerIdentity instance already exists, then do not update it.
	if err == nil {
		return ctrl.Result{}, nil
	}
	if apierrors.IsNotFound(err) {
		log.Info("AWSClusterControllerIdentity instance not found, creating a new instance")
		// Fetch the AWSClusterControllerIdentity instance
		controllerIdentity = &infrav1.AWSClusterControllerIdentity{
			TypeMeta: metav1.TypeMeta{
				APIVersion: infrav1.GroupVersion.String(),
				Kind:       string(infrav1.ControllerIdentityKind),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: infrav1.AWSClusterControllerIdentityName,
			},
			Spec: infrav1.AWSClusterControllerIdentitySpec{
				AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
					AllowedNamespaces: &infrav1.AllowedNamespaces{},
				},
			},
		}
		if err := r.Create(ctx, controllerIdentity); err != nil {
			if apierrors.IsAlreadyExists(err) {
				return reconcile.Result{}, nil
			}
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	return reconcile.Result{}, err
}

func (r *AWSControllerIdentityReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	controller := ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.AWSCluster{}).
		Named("awscontrolleridentity").
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), logger.FromContext(ctx).GetLogger(), r.WatchFilterValue))

	if feature.Gates.Enabled(feature.EKS) {
		controller.Watches(
			&ekscontrolplanev1.AWSManagedControlPlane{},
			handler.EnqueueRequestsFromMapFunc(r.managedControlPlaneMap),
		)
	}

	return controller.Complete(r)
}

func (r *AWSControllerIdentityReconciler) managedControlPlaneMap(_ context.Context, o client.Object) []ctrl.Request {
	managedControlPlane, ok := o.(*ekscontrolplanev1.AWSManagedControlPlane)
	if !ok {
		klog.Errorf("Expected a managedControlPlane but got a %T", o)
	}

	return []ctrl.Request{
		{
			NamespacedName: types.NamespacedName{
				Name:      managedControlPlane.Name,
				Namespace: managedControlPlane.Namespace,
			},
		},
	}
}
