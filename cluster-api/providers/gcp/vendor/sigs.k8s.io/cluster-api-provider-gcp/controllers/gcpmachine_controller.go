/*
Copyright 2019 The Kubernetes Authors.

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

// Package controllers implements controller types.
package controllers

import (
	"context"
	"time"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/compute/instances"
	"sigs.k8s.io/cluster-api-provider-gcp/util/reconciler"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/predicates"
	"sigs.k8s.io/cluster-api/util/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// GCPMachineReconciler reconciles a GCPMachine object.
type GCPMachineReconciler struct {
	client.Client
	ReconcileTimeout time.Duration
	WatchFilterValue string
}

// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=gcpmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=gcpmachines/status,verbs=get;update;patch

func (r *GCPMachineReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)
	c, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.GCPMachine{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), ctrl.LoggerFrom(ctx), r.WatchFilterValue)).
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrav1.GroupVersion.WithKind("GCPMachine"))),
		).
		Watches(
			&infrav1.GCPCluster{},
			handler.EnqueueRequestsFromMapFunc(r.GCPClusterToGCPMachines(ctx)),
		).
		Build(r)
	if err != nil {
		return errors.Wrap(err, "error creating controller")
	}

	clusterToObjectFunc, err := util.ClusterToTypedObjectsMapper(r.Client, &infrav1.GCPMachineList{}, mgr.GetScheme())
	if err != nil {
		return errors.Wrap(err, "failed to create mapper for Cluster to GCPMachines")
	}

	// Add a watch on clusterv1.Cluster object for unpause & ready notifications.
	if err := c.Watch(
		source.Kind[client.Object](mgr.GetCache(), &clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(clusterToObjectFunc),
			predicates.ClusterPausedTransitionsOrInfrastructureReady(mgr.GetScheme(), log),
		)); err != nil {
		return errors.Wrap(err, "failed adding a watch for ready clusters")
	}

	return nil
}

// GCPClusterToGCPMachines is a handler.ToRequestsFunc to be used to enqeue requests for reconciliation
// of GCPMachines.
func (r *GCPMachineReconciler) GCPClusterToGCPMachines(ctx context.Context) handler.MapFunc {
	log := ctrl.LoggerFrom(ctx)
	return func(mapCtx context.Context, o client.Object) []ctrl.Request {
		result := []ctrl.Request{}

		c, ok := o.(*infrav1.GCPCluster)
		if !ok {
			log.Error(errors.Errorf("expected a GCPCluster but got a %T", o), "failed to get GCPMachine for GCPCluster")
			return nil
		}

		cluster, err := util.GetOwnerCluster(mapCtx, r.Client, c.ObjectMeta)
		switch {
		case apierrors.IsNotFound(err) || cluster == nil:
			return result
		case err != nil:
			log.Error(err, "failed to get owning cluster")
			return result
		}

		labels := map[string]string{clusterv1.ClusterNameLabel: cluster.Name}
		machineList := &clusterv1.MachineList{}
		if err := r.List(mapCtx, machineList, client.InNamespace(c.Namespace), client.MatchingLabels(labels)); err != nil {
			log.Error(err, "failed to list Machines")
			return nil
		}
		for _, m := range machineList.Items {
			if m.Spec.InfrastructureRef.Name == "" {
				continue
			}
			name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.InfrastructureRef.Name}
			result = append(result, ctrl.Request{NamespacedName: name})
		}

		return result
	}
}

func (r *GCPMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultedLoopTimeout(r.ReconcileTimeout))
	defer cancel()

	log := ctrl.LoggerFrom(ctx)
	gcpMachine := &infrav1.GCPMachine{}
	err := r.Get(ctx, req.NamespacedName, gcpMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	machine, err := util.GetOwnerMachine(ctx, r.Client, gcpMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		log.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("machine", machine.Name)
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		log.Info("Machine is missing cluster label or cluster does not exist")

		return ctrl.Result{}, nil
	}

	if annotations.IsPaused(cluster, gcpMachine) {
		log.Info("GCPMachine or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)
	gcpCluster := &infrav1.GCPCluster{}
	gcpClusterKey := client.ObjectKey{
		Namespace: gcpMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, gcpClusterKey, gcpCluster); err != nil {
		log.Info("GCPCluster is not available yet")
		return ctrl.Result{}, nil
	}

	// Create the cluster scope
	clusterScope, err := scope.NewClusterScope(ctx, scope.ClusterScopeParams{
		Client:     r.Client,
		Cluster:    cluster,
		GCPCluster: gcpCluster,
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	// Create the machine scope
	machineScope, err := scope.NewMachineScope(scope.MachineScopeParams{
		Client:        r.Client,
		Machine:       machine,
		GCPMachine:    gcpMachine,
		ClusterGetter: clusterScope,
	})
	if err != nil {
		return ctrl.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}

	// Always close the scope when exiting this function so we can persist any GCPMachine changes.
	defer func() {
		if err := machineScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Handle deleted machines
	if !gcpMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(ctx, machineScope)
	}

	// Handle non-deleted machines
	return r.reconcile(ctx, machineScope)
}

func (r *GCPMachineReconciler) reconcile(ctx context.Context, machineScope *scope.MachineScope) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling GCPMachine")

	controllerutil.AddFinalizer(machineScope.GCPMachine, infrav1.MachineFinalizer)
	if err := machineScope.PatchObject(); err != nil {
		return ctrl.Result{}, err
	}

	if err := instances.New(machineScope).Reconcile(ctx); err != nil {
		log.Error(err, "Error reconciling instance resources")
		record.Warnf(machineScope.GCPMachine, "GCPMachineReconcile", "Reconcile error - %v", err)
		return ctrl.Result{}, err
	}

	instanceState := *machineScope.GetInstanceStatus()
	switch instanceState {
	case infrav1.InstanceStatusProvisioning, infrav1.InstanceStatusStaging:
		log.Info("GCPMachine instance is pending", "instance-id", *machineScope.GetInstanceID())
		record.Eventf(machineScope.GCPMachine, "GCPMachineReconcile", "GCPMachine instance is pending - instance-id: %s", *machineScope.GetInstanceID())
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	case infrav1.InstanceStatusRunning:
		log.Info("GCPMachine instance is running", "instance-id", *machineScope.GetInstanceID())
		record.Eventf(machineScope.GCPMachine, "GCPMachineReconcile", "GCPMachine instance is running - instance-id: %s", *machineScope.GetInstanceID())
		record.Event(machineScope.GCPMachine, "GCPMachineReconcile", "Reconciled")
		machineScope.SetReady()
		return ctrl.Result{}, nil
	default:
		machineScope.SetFailureReason("UpdateError")
		machineScope.SetFailureMessage(errors.Errorf("GCPMachine instance state %s is unexpected", instanceState))
		return ctrl.Result{Requeue: true}, nil
	}
}

func (r *GCPMachineReconciler) reconcileDelete(ctx context.Context, machineScope *scope.MachineScope) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling Delete GCPMachine")

	if err := instances.New(machineScope).Delete(ctx); err != nil {
		log.Error(err, "Error deleting instance resources")
		return err
	}

	controllerutil.RemoveFinalizer(machineScope.GCPMachine, infrav1.MachineFinalizer)
	record.Event(machineScope.GCPMachine, "GCPMachineReconcile", "Reconciled")
	return nil
}
