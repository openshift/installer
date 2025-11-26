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

// Package controllers provides a way to reconcile GKEConfig objects.
package controllers

import (
	"context"
	"time"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	infrav1exp "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	bootstrapv1exp "sigs.k8s.io/cluster-api-provider-gcp/exp/bootstrap/gke/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// GKEConfigReconciler reconciles a GKEConfig object.
type GKEConfigReconciler struct {
	client.Client
	WatchFilterValue string
	ReconcileTimeout time.Duration
}

// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=gkeconfigs,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=gkeconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools;clusters,verbs=get;list;watch

func (r *GKEConfigReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, option controller.Options) error {
	log := ctrl.LoggerFrom(ctx)

	b := ctrl.NewControllerManagedBy(mgr).
		For(&bootstrapv1exp.GKEConfig{}).
		WithOptions(option).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), log, r.WatchFilterValue)).
		Watches(
			&infrav1exp.GCPManagedMachinePool{},
			handler.EnqueueRequestsFromMapFunc(r.ManagedMachinePoolToGKEConfigMapFunc),
		)

	_, err := b.Build(r)
	if err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	return nil
}

func (r *GKEConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, rerr error) {
	log := ctrl.LoggerFrom(ctx)

	config := &bootstrapv1exp.GKEConfig{}
	if err := r.Get(ctx, req.NamespacedName, config); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get config")
		return ctrl.Result{}, err
	}
	log = log.WithValues("GKEConfig", config.GetName())

	machinePool, err := getOwnerMachinePool(ctx, r.Client, config.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to get owner MachinePool for GKEConfig", "GKEConfig", config.GetName())
		return ctrl.Result{}, err
	}

	if machinePool == nil {
		log.Info("No owner MachinePool found for GKEConfig", "GKEConfig", config.GetName())
		return ctrl.Result{}, nil
	}

	// fetch associated GCPManagedMachinePool
	gcpMP := &infrav1exp.GCPManagedMachinePool{}
	gcpMPKey := types.NamespacedName{
		Name:      machinePool.Spec.Template.Spec.InfrastructureRef.Name,
		Namespace: machinePool.Spec.Template.Spec.InfrastructureRef.Namespace,
	}
	if err := r.Get(ctx, gcpMPKey, gcpMP); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("GCPManagedMachinePool not found for MachinePool", "GCPManagedMachinePool", gcpMPKey)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// check if GCPManagedMachinePool is ready
	if !gcpMP.Status.Ready {
		log.Info("Waiting for GCPManagedMachinePool to be ready", "GKEConfig", config.GetName(), "GCPManagedMachinePool", gcpMPKey)
		return ctrl.Result{}, nil
	}

	// set GKEConfig as ready when GCPManagedMachinePool becomes ready
	config.Status.Ready = true
	if err := r.Status().Update(ctx, config); err != nil {
		log.Info("Failed to update GKEConfig status", "GKEConfig", config.GetName(), "error", err)
		return ctrl.Result{}, err
	}

	log.Info("Successfully reconciled GKEConfig", "GKEConfig", req.NamespacedName, "MachinePool", machinePool.GetName())

	return ctrl.Result{}, nil
}

// ManagedMachinePoolToGKEConfigMapFunc is a handler.ToRequestsFunc to be used to enqueue requests for
// GKEConfig reconciliation.
func (r *GKEConfigReconciler) ManagedMachinePoolToGKEConfigMapFunc(_ context.Context, o client.Object) []ctrl.Request {
	c, ok := o.(*infrav1exp.GCPManagedMachinePool)
	if !ok {
		klog.Errorf("Expected a Cluster but got a %T", o)
	}

	machinePool, err := getOwnerMachinePool(context.Background(), r.Client, c.ObjectMeta)
	if err != nil {
		klog.Errorf("Failed to get owner MachinePool for GCPManagedMachinePool %s/%s: %v", c.Namespace, c.Name, err)
		return nil
	}

	if machinePool == nil {
		klog.Infof("No owner MachinePool found for GCPManagedMachinePool %s/%s", c.Namespace, c.Name)
		return nil
	}

	return []ctrl.Request{
		{
			NamespacedName: client.ObjectKey{
				Name:      machinePool.Spec.Template.Spec.InfrastructureRef.Name,
				Namespace: machinePool.Spec.Template.Spec.InfrastructureRef.Namespace,
			},
		},
	}
}

func getOwnerMachinePool(ctx context.Context, c client.Client, obj metav1.ObjectMeta) (*expclusterv1.MachinePool, error) {
	for _, ref := range obj.OwnerReferences {
		if ref.Kind != "MachinePool" {
			continue
		}
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if gv.Group == expclusterv1.GroupVersion.Group {
			return getMachinePoolByName(ctx, c, obj.Namespace, ref.Name)
		}
	}

	return nil, nil
}

func getMachinePoolByName(ctx context.Context, c client.Client, namespace, name string) (*expclusterv1.MachinePool, error) {
	m := &expclusterv1.MachinePool{}
	key := client.ObjectKey{Name: name, Namespace: namespace}
	if err := c.Get(ctx, key, m); err != nil {
		return nil, err
	}

	return m, nil
}
