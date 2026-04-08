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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck
	"sigs.k8s.io/cluster-api/util"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"         //nolint:staticcheck
	v1beta2conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2" //nolint:staticcheck
	v1beta1patch "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/patch"                   //nolint:staticcheck
	"sigs.k8s.io/cluster-api/util/deprecated/v1beta1/paused"
	"sigs.k8s.io/cluster-api/util/finalizers"
	"sigs.k8s.io/cluster-api/util/predicates"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
)

// IBMVPCClusterReconciler reconciles a IBMVPCCluster object.
type IBMVPCClusterReconciler struct {
	client.Client
	Log             logr.Logger
	Recorder        record.EventRecorder
	ServiceEndpoint []endpoints.ServiceEndpoint
	Scheme          *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

// Reconcile implements controller runtime Reconciler interface and handles reconcileation logic for IBMVPCCluster.
func (r *IBMVPCClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	log.Info("Reconciling IBMVPCCluster")
	defer log.Info("Finished reconciling IBMVPCCluster")

	// Fetch the IBMVPCCluster instance.
	ibmVPCCluster := &infrav1.IBMVPCCluster{}
	err := r.Get(ctx, req.NamespacedName, ibmVPCCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("IBMVPCCluster not found")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Determine whether the Cluster is designed for extended Infrastructure support, implemented in a separate path.
	if ibmVPCCluster.Spec.Network != nil {
		return r.reconcileV2(ctx, req)
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, ibmVPCCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	// Add finalizer first if not set to avoid the race condition between init and delete.
	if finalizerAdded, err := finalizers.EnsureFinalizer(ctx, r.Client, ibmVPCCluster, infrav1.ClusterFinalizer); err != nil || finalizerAdded {
		return ctrl.Result{}, err
	}

	log = log.WithValues("Cluster", klog.KObj(cluster))
	ctx = ctrl.LoggerInto(ctx, log)

	if isPaused, requeue, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, ibmVPCCluster); err != nil || isPaused || requeue {
		return ctrl.Result{}, err
	}

	clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:          r.Client,
		Cluster:         cluster,
		IBMVPCCluster:   ibmVPCCluster,
		ServiceEndpoint: r.ServiceEndpoint,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	// Initialize the patch helper.
	patchHelper, err := v1beta1patch.NewHelper(ibmVPCCluster, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to initialize patch helper: %w", err)
	}

	// Always attempt to Patch the IBMVPCCluster object and status after each reconciliation.
	defer func() {
		if err := patchIBMVPCCluster(ctx, patchHelper, ibmVPCCluster); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	// Handle deleted clusters.
	if !ibmVPCCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, clusterScope)
	}

	return r.reconcile(ctx, clusterScope)
}

func (r *IBMVPCClusterReconciler) reconcileV2(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx).WithValues("controller", "IBMVPCCluster")

	// Fetch the IBMVPCCluster instance.
	ibmVPCCluster := &infrav1.IBMVPCCluster{}
	err := r.Get(ctx, req.NamespacedName, ibmVPCCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Add finalizer first if not set to avoid the race condition between init and delete.
	if finalizerAdded, err := finalizers.EnsureFinalizer(ctx, r.Client, ibmVPCCluster, infrav1.ClusterFinalizer); err != nil || finalizerAdded {
		return ctrl.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, ibmVPCCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	if isPaused, requeue, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, ibmVPCCluster); err != nil || isPaused || requeue {
		return ctrl.Result{}, err
	}

	clusterScope, err := scope.NewVPCClusterScope(scope.VPCClusterScopeParams{
		Client:          r.Client,
		Logger:          log,
		Cluster:         cluster,
		IBMVPCCluster:   ibmVPCCluster,
		ServiceEndpoint: r.ServiceEndpoint,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	// Initialize the patch helper.
	patchHelper, err := v1beta1patch.NewHelper(ibmVPCCluster, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the IBMVPCCluster object and status after each reconciliation.
	defer func() {
		if err := patchIBMVPCCluster(ctx, patchHelper, ibmVPCCluster); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	// Handle deleted clusters.
	if !ibmVPCCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDeleteV2(clusterScope)
	}

	return r.reconcileCluster(ctx, clusterScope)
}

func (r *IBMVPCClusterReconciler) reconcile(ctx context.Context, clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx).WithValues("controller", "IBMVPCCluster")
	// If the IBMVPCCluster doesn't have our finalizer, add it.
	if controllerutil.AddFinalizer(clusterScope.IBMVPCCluster, infrav1.ClusterFinalizer) {
		return ctrl.Result{}, nil
	}

	if clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint.Host != "" {
		loadBalancerEndpoint, err := clusterScope.GetLoadBalancerByHostname(clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint.Host)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error when retrieving load balancer with specified hostname: %w", err)
		}

		if loadBalancerEndpoint == nil {
			return ctrl.Result{}, fmt.Errorf("no loadBalancer found with hostname - %s", clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint.Host)
		}
		r.reconcileLBState(ctx, clusterScope, loadBalancerEndpoint)
	}

	log.Info("Reconciling VPC")
	vpc, err := clusterScope.CreateVPC()
	if err != nil {
		v1beta1conditions.MarkFalse(clusterScope.IBMVPCCluster, infrav1.VPCReadyCondition, infrav1.VPCReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
			Type:    infrav1.VPCReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VPCNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		return ctrl.Result{}, fmt.Errorf("failed to reconcile VPC for IBMVPCCluster %s: %w", klog.KObj(clusterScope.IBMVPCCluster), err)
	}
	if vpc != nil {
		clusterScope.IBMVPCCluster.Status.VPC = infrav1.VPC{
			ID:   *vpc.ID,
			Name: *vpc.Name,
		}
		v1beta1conditions.MarkTrue(clusterScope.IBMVPCCluster, infrav1.VPCReadyCondition)
		v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
			Type:   infrav1.VPCReadyV1Beta2Condition,
			Status: metav1.ConditionTrue,
			Reason: infrav1.VPCReadyV1Beta2Reason,
		})
		log.Info("Reconciliation of VPC complete")
	}

	if clusterScope.IBMVPCCluster.Status.Subnet.ID == nil {
		log.Info("Reconciling VPC Subnets")
		subnet, err := clusterScope.CreateSubnet()
		if err != nil {
			v1beta1conditions.MarkFalse(clusterScope.IBMVPCCluster, infrav1.VPCSubnetReadyCondition, infrav1.VPCSubnetReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
			v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
				Type:    infrav1.VPCSubnetReadyV1Beta2Condition,
				Status:  metav1.ConditionFalse,
				Reason:  infrav1.VPCSubnetNotReadyV1Beta2Reason,
				Message: err.Error(),
			})
			return ctrl.Result{}, fmt.Errorf("failed to reconcile Subnet for IBMVPCCluster %s/%s: %w", clusterScope.IBMVPCCluster.Namespace, clusterScope.IBMVPCCluster.Name, err)
		}
		if subnet != nil {
			clusterScope.IBMVPCCluster.Status.Subnet = infrav1.Subnet{
				Ipv4CidrBlock: subnet.Ipv4CIDRBlock,
				Name:          subnet.Name,
				ID:            subnet.ID,
				Zone:          subnet.Zone.Name,
			}
			v1beta1conditions.MarkTrue(clusterScope.IBMVPCCluster, infrav1.VPCSubnetReadyCondition)
			v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
				Type:   infrav1.VPCSubnetReadyV1Beta2Condition,
				Status: metav1.ConditionTrue,
				Reason: infrav1.VPCSubnetReadyV1Beta2Reason,
			})
			log.Info("Reconciliation of VPC Subnets complete")
		}
	}

	if clusterScope.IBMVPCCluster.Spec.ControlPlaneLoadBalancer != nil && clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint.Host == "" {
		log.Info("Reconciling Load Balancers")
		loadBalancer, err := r.getOrCreate(clusterScope)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to reconcile Control Plane LoadBalancer for IBMVPCCluster %s/%s: %w", clusterScope.IBMVPCCluster.Namespace, clusterScope.IBMVPCCluster.Name, err)
		}

		if loadBalancer != nil {
			clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint.Host = *loadBalancer.Hostname
			r.reconcileLBState(ctx, clusterScope, loadBalancer)
			log.Info("Reconciliation of Load Balancers complete")
		}
	}

	// Requeue after 1 minute if cluster is not ready to update status of the cluster properly.
	if !clusterScope.IsReady() {
		log.Info("Cluster is not yet ready")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}
	return ctrl.Result{}, nil
}

func (r *IBMVPCClusterReconciler) reconcileCluster(ctx context.Context, clusterScope *scope.VPCClusterScope) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	// If the IBMVPCCluster doesn't have our finalizer, add it.
	if controllerutil.AddFinalizer(clusterScope.IBMVPCCluster, infrav1.ClusterFinalizer) {
		return ctrl.Result{}, nil
	}

	// Reconcile the cluster's VPC.
	log.Info("Reconciling VPC")
	if requeue, err := clusterScope.ReconcileVPC(ctx); err != nil {
		log.Error(err, "failed to reconcile VPC")
		v1beta1conditions.MarkFalse(clusterScope.IBMVPCCluster, infrav1.VPCReadyCondition, infrav1.VPCReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
			Type:    infrav1.VPCReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VPCNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		return reconcile.Result{}, err
	} else if requeue {
		log.Info("VPC creation is pending, requeuing")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}
	log.Info("Reconciliation of VPC complete")
	v1beta1conditions.MarkTrue(clusterScope.IBMVPCCluster, infrav1.VPCReadyCondition)
	v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
		Type:   infrav1.VPCReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1.VPCReadyV1Beta2Reason,
	})

	// Reconcile the cluster's VPC Custom Image.
	log.Info("Reconciling VPC Custom Image")
	if requeue, err := clusterScope.ReconcileVPCCustomImage(ctx); err != nil {
		log.Error(err, "failed to reconcile VPC Custom Image")
		v1beta1conditions.MarkFalse(clusterScope.IBMVPCCluster, infrav1.ImageReadyCondition, infrav1.ImageReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
			Type:    infrav1.VPCImageReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VPCImageNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		return reconcile.Result{}, err
	} else if requeue {
		log.Info("VPC Custom Image creation is pending, requeueing")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}
	log.Info("Reconciliation of VPC Custom Image complete")
	v1beta1conditions.MarkTrue(clusterScope.IBMVPCCluster, infrav1.ImageReadyCondition)
	v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
		Type:   infrav1.VPCImageReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1.VPCImageReadyV1Beta2Reason,
	})

	// Reconcile the cluster's VPC Subnets.
	log.Info("Reconciling VPC Subnets")
	if requeue, err := clusterScope.ReconcileSubnets(ctx); err != nil {
		log.Error(err, "failed to reconcile VPC Subnets")
		v1beta1conditions.MarkFalse(clusterScope.IBMVPCCluster, infrav1.VPCSubnetReadyCondition, infrav1.VPCSubnetReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
			Type:    infrav1.VPCSubnetReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VPCSubnetNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		return reconcile.Result{}, err
	} else if requeue {
		log.Info("VPC Subnets creation is pending, requeueing")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}
	log.Info("Reconciliation of VPC Subnets complete")
	v1beta1conditions.MarkTrue(clusterScope.IBMVPCCluster, infrav1.VPCSubnetReadyCondition)
	v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
		Type:   infrav1.VPCSubnetReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1.VPCSubnetReadyV1Beta2Reason,
	})

	// Reconcile the cluster's Security Groups (and Security Group Rules)
	log.Info("Reconciling Security Groups")
	if requeue, err := clusterScope.ReconcileSecurityGroups(ctx); err != nil {
		log.Error(err, "failed to reconcile Security Groups")
		v1beta1conditions.MarkFalse(clusterScope.IBMVPCCluster, infrav1.VPCSecurityGroupReadyCondition, infrav1.VPCSecurityGroupReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
			Type:    infrav1.VPCSecurityGroupReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VPCSecurityGroupNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		return reconcile.Result{}, err
	} else if requeue {
		log.Info("Security Groups creation is pending, requeueing")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}
	log.Info("Reconciliation of Security Groups complete")
	v1beta1conditions.MarkTrue(clusterScope.IBMVPCCluster, infrav1.VPCSecurityGroupReadyCondition)
	v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
		Type:   infrav1.VPCSecurityGroupReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1.VPCSecurityGroupReadyV1Beta2Reason,
	})

	// Reconcile the cluster's Load Balancers
	log.Info("Reconciling Load Balancers")
	if requeue, err := clusterScope.ReconcileLoadBalancers(ctx); err != nil {
		log.Error(err, "failed to reconcile Load Balancers")
		v1beta1conditions.MarkFalse(clusterScope.IBMVPCCluster, infrav1.LoadBalancerReadyCondition, infrav1.LoadBalancerReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
			Type:    infrav1.VPCLoadBalancerReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VPCLoadBalancerNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		return reconcile.Result{}, err
	} else if requeue {
		log.Info("Load Balancers creation is pending, requeueing")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}
	log.Info("Reconciliation of Load Balancers complete")
	v1beta1conditions.MarkTrue(clusterScope.IBMVPCCluster, infrav1.LoadBalancerReadyCondition)
	v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
		Type:   infrav1.VPCLoadBalancerReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1.VPCLoadBalancerReadyV1Beta2Reason,
	})

	// Collect cluster's Load Balancer hostname for spec.
	hostName, err := clusterScope.GetLoadBalancerHostName()
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("error retrieving load balancer hostname: %w", err)
	} else if hostName == nil || *hostName == "" {
		log.Info("No Load Balancer hostname found, requeueing")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}

	// Mark cluster as ready.
	clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint.Host = *hostName
	clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint.Port = clusterScope.GetAPIServerPort()
	clusterScope.IBMVPCCluster.Status.Ready = true
	log.Info("cluster infrastructure is now ready for cluster", "clusterName", clusterScope.IBMVPCCluster.Name)
	return ctrl.Result{}, nil
}

func (r *IBMVPCClusterReconciler) reconcileDelete(ctx context.Context, clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	// check if still have existing VSIs.
	listVSIOpts := &vpcv1.ListInstancesOptions{
		VPCID: &clusterScope.IBMVPCCluster.Status.VPC.ID,
	}
	vsis, _, err := clusterScope.IBMVPCClient.ListInstances(listVSIOpts)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error when listing VSIs when tried to delete subnet: %w", err)
	}
	// skip deleting other resources if still have vsis running.
	if *vsis.TotalCount != int64(0) {
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	// skip load balancer deletion if a pre-created load balancer is being set as the controlplane endpoint.
	if clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint.Host != "" && clusterScope.IBMVPCCluster.Spec.ControlPlaneLoadBalancer == nil {
		return handleFinalizerRemoval(clusterScope)
	}

	if clusterScope.IBMVPCCluster.Spec.ControlPlaneLoadBalancer != nil {
		loadBalancer, err := clusterScope.GetLoadBalancerByHostname(clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint.Host)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error when retrieving load balancer with specified hostname: %w", err)
		}

		if loadBalancer == nil && (string(clusterScope.GetLoadBalancerState()) != string(infrav1.VPCLoadBalancerStateDeletePending)) {
			return handleFinalizerRemoval(clusterScope)
		}
		if loadBalancer != nil {
			clusterScope.SetLoadBalancerState(*loadBalancer.ProvisioningStatus)
			if *loadBalancer.Name != clusterScope.IBMVPCCluster.Spec.ControlPlaneLoadBalancer.Name {
				return handleFinalizerRemoval(clusterScope)
			}

			log.Info("Deleting VPC load balancer")
			v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
				Type:   infrav1.VPCLoadBalancerReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.VPCLoadBalancerDeletingV1Beta2Reason,
			})
			deleted, err := clusterScope.DeleteLoadBalancer()
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to delete loadBalancer: %w", err)
			}
			// Skip deleting other resources if still have loadBalancers running.
			if deleted {
				return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
			}
		}
	}

	log.Info("Deleting VPC subnet")
	v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
		Type:   infrav1.VPCSubnetReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: infrav1.VPCSubnetDeletingV1Beta2Reason,
	})
	if err := clusterScope.DeleteSubnet(ctx); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to delete subnet: %w", err)
	}

	log.Info("Deleting VPC")
	v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
		Type:   infrav1.VPCReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: infrav1.VPCDeletingV1Beta2Reason,
	})
	if err := clusterScope.DeleteVPC(); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to delete VPC: %w", err)
	}

	log.Info("IBMVPCCluster deletion completed")
	return handleFinalizerRemoval(clusterScope)
}

func (r *IBMVPCClusterReconciler) reconcileDeleteV2(clusterScope *scope.VPCClusterScope) (ctrl.Result, error) { //nolint:unparam
	clusterScope.Info("Delete cluster is not implemented for reconcile v2")
	controllerutil.RemoveFinalizer(clusterScope.IBMVPCCluster, infrav1.ClusterFinalizer)
	return ctrl.Result{}, nil
}

func (r *IBMVPCClusterReconciler) getOrCreate(clusterScope *scope.ClusterScope) (*vpcv1.LoadBalancer, error) {
	loadBalancer, err := clusterScope.CreateLoadBalancer()
	return loadBalancer, err
}

func handleFinalizerRemoval(clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	controllerutil.RemoveFinalizer(clusterScope.IBMVPCCluster, infrav1.ClusterFinalizer)
	return ctrl.Result{}, nil
}

func (r *IBMVPCClusterReconciler) reconcileLBState(ctx context.Context, clusterScope *scope.ClusterScope, loadBalancer *vpcv1.LoadBalancer) {
	log := ctrl.LoggerFrom(ctx)
	if clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint.Port == 0 {
		clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint.Port = clusterScope.APIServerPort()
	}

	clusterScope.SetLoadBalancerID(loadBalancer.ID)
	log.V(3).Info("LoadBalancerID - " + clusterScope.GetLoadBalancerID())
	clusterScope.SetLoadBalancerAddress(loadBalancer.Hostname)
	clusterScope.SetLoadBalancerState(*loadBalancer.ProvisioningStatus)
	log.V(3).Info("LoadBalancerState - " + string(clusterScope.GetLoadBalancerState()))

	switch clusterScope.GetLoadBalancerState() {
	case infrav1.VPCLoadBalancerStateCreatePending:
		log.V(3).Info("LoadBalancer is in create state")
		clusterScope.SetNotReady()
		v1beta1conditions.MarkFalse(clusterScope.IBMVPCCluster, infrav1.LoadBalancerReadyCondition, string(infrav1.VPCLoadBalancerStateCreatePending), clusterv1beta1.ConditionSeverityInfo, "%s", *loadBalancer.OperatingStatus)
		v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
			Type:    infrav1.VPCLoadBalancerReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VPCLoadBalancerNotReadyV1Beta2Reason,
			Message: "VPC load balancer is in creating state",
		})
	case infrav1.VPCLoadBalancerStateActive:
		log.V(3).Info("LoadBalancer is in active state")
		clusterScope.SetReady()
		v1beta1conditions.MarkTrue(clusterScope.IBMVPCCluster, infrav1.LoadBalancerReadyCondition)
		v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
			Type:    infrav1.VPCLoadBalancerReadyV1Beta2Condition,
			Status:  metav1.ConditionTrue,
			Reason:  infrav1.VPCLoadBalancerReadyV1Beta2Reason,
			Message: "VPC load balancer is in active state",
		})
	default:
		log.V(3).Info("LoadBalancer state is undefined", "state", clusterScope.GetLoadBalancerState(), "loadbalancerID", clusterScope.GetLoadBalancerID())
		clusterScope.SetNotReady()
		v1beta1conditions.MarkUnknown(clusterScope.IBMVPCCluster, infrav1.LoadBalancerReadyCondition, *loadBalancer.ProvisioningStatus, "")
		v1beta2conditions.Set(clusterScope.IBMVPCCluster, metav1.Condition{
			Type:    infrav1.VPCLoadBalancerReadyV1Beta2Condition,
			Status:  metav1.ConditionUnknown,
			Reason:  infrav1.VPCLoadBalancerNotReadyV1Beta2Reason,
			Message: fmt.Sprintf("VPC load balancer is in an unknown state: %s", *loadBalancer.ProvisioningStatus),
		})
	}
}

// SetupWithManager creates a new IBMVPCCluster controller for a manager.
func (r *IBMVPCClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.IBMVPCCluster{}).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(r.Scheme, ctrl.LoggerFrom(ctx))).
		Complete(r)
}

// patchIBMVPCCluster updates the IBMVPCCluster and its status on the API server.
func patchIBMVPCCluster(ctx context.Context, patchHelper *v1beta1patch.Helper, ibmVPCCluster *infrav1.IBMVPCCluster) error {
	if err := v1beta2conditions.SetSummaryCondition(ibmVPCCluster, ibmVPCCluster, infrav1.IBMVPCClusterReadyV1Beta2Condition,
		v1beta2conditions.ForConditionTypes{
			infrav1.VPCReadyV1Beta2Condition,
			infrav1.VPCSubnetReadyV1Beta2Condition,
			infrav1.VPCLoadBalancerReadyV1Beta2Condition,
		},
		v1beta2conditions.IgnoreTypesIfMissing{
			infrav1.VPCSecurityGroupReadyV1Beta2Condition,
			infrav1.VPCImageReadyV1Beta2Condition,
		},
		// Using a custom merge strategy to override reasons applied during merge.
		v1beta2conditions.CustomMergeStrategy{
			MergeStrategy: v1beta2conditions.DefaultMergeStrategy(
				// Use custom reasons.
				v1beta2conditions.ComputeReasonFunc(v1beta2conditions.GetDefaultComputeMergeReasonFunc(
					infrav1.IBMVPCClusterNotReadyV1Beta2Reason,
					infrav1.IBMVPCClusterReadyUnknownV1Beta2Reason,
					infrav1.IBMVPCClusterReadyV1Beta2Reason,
				)),
			),
		},
	); err != nil {
		return fmt.Errorf("failed to set %s condition: %w", infrav1.IBMVPCClusterReadyV1Beta2Condition, err)
	}

	// Patch the IBMVPCCluster resource.
	return patchHelper.Patch(ctx, ibmVPCCluster, v1beta1patch.WithOwnedV1Beta2Conditions{Conditions: []string{
		infrav1.IBMVPCClusterReadyV1Beta2Condition,
		clusterv1beta1.PausedV1Beta2Condition,
		infrav1.VPCReadyV1Beta2Condition,
		infrav1.VPCSubnetReadyV1Beta2Condition,
		infrav1.VPCSecurityGroupReadyV1Beta2Condition,
		infrav1.VPCLoadBalancerReadyV1Beta2Condition,
		infrav1.VPCImageReadyV1Beta2Condition,
	}})
}
