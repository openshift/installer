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
	"strings"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	v1beta2conditions "sigs.k8s.io/cluster-api/util/conditions/v1beta2"
	"sigs.k8s.io/cluster-api/util/finalizers"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/paused"
	"sigs.k8s.io/cluster-api/util/predicates"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
)

// IBMPowerVSClusterReconciler reconciles a IBMPowerVSCluster object.
type IBMPowerVSClusterReconciler struct {
	client.Client
	Recorder        record.EventRecorder
	ServiceEndpoint []endpoints.ServiceEndpoint
	Scheme          *runtime.Scheme

	// WatchFilterValue is the label value used to filter events prior to reconciliation.
	WatchFilterValue string

	ClientFactory scope.ClientFactory
}

type powerVSCluster struct {
	cluster *infrav1beta2.IBMPowerVSCluster
	mu      sync.Mutex
}

type reconcileResult struct {
	reconcile.Result
	error
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters/status,verbs=get;update;patch

// Reconcile implements controller runtime Reconciler interface and handles reconcileation logic for IBMPowerVSCluster.
func (r *IBMPowerVSClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	log.Info("Reconciling IBMPowerVSCluster")
	defer log.Info("Finished reconciling IBMPowerVSCluster")

	// Fetch the IBMPowerVSCluster instance.
	ibmPowerVSCluster := &infrav1beta2.IBMPowerVSCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, ibmPowerVSCluster); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("IBMPowerVSCluster not found")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get IBMPowerVSCluster: %w", err)
	}

	// Add finalizer first if not set to avoid the race condition between init and delete.
	if finalizerAdded, err := finalizers.EnsureFinalizer(ctx, r.Client, ibmPowerVSCluster, infrav1beta2.IBMPowerVSClusterFinalizer); err != nil || finalizerAdded {
		return ctrl.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, ibmPowerVSCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get cluster for IBMPowerVSCluster: %w", err)
	}
	if cluster == nil {
		log.Info("Waiting for cluster controller to set OwnerRef on IBMPowerVSCluster")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("Cluster", klog.KObj(cluster))
	ctx = ctrl.LoggerInto(ctx, log)

	if isPaused, requeue, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, ibmPowerVSCluster); err != nil || isPaused || requeue {
		return ctrl.Result{}, err
	}

	// Create the scope.
	clusterScope, err := scope.NewPowerVSClusterScope(scope.PowerVSClusterScopeParams{
		Client:            r.Client,
		Cluster:           cluster,
		IBMPowerVSCluster: ibmPowerVSCluster,
		ServiceEndpoint:   r.ServiceEndpoint,
		ClientFactory:     r.ClientFactory,
	})

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create IBMPowerVSCluster scope: %w", err)
	}

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(ibmPowerVSCluster, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the IBMPowerVSCluster object and status after each reconciliation.
	defer func() {
		if err := patchIBMPowerVSCluster(ctx, patchHelper, ibmPowerVSCluster); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	// Handle deleted clusters.
	if !ibmPowerVSCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, clusterScope)
	}

	return r.reconcile(ctx, clusterScope)
}

func (r *IBMPowerVSClusterReconciler) reconcile(ctx context.Context, clusterScope *scope.PowerVSClusterScope) (ctrl.Result, error) { //nolint:gocyclo
	log := ctrl.LoggerFrom(ctx)
	// check for annotation set for cluster resource and decide on proceeding with infra creation.
	// do not proceed further if "powervs.cluster.x-k8s.io/create-infra=true" annotation is not set.
	if !scope.CheckCreateInfraAnnotation(*clusterScope.IBMPowerVSCluster) {
		log.V(3).Info("IBMPowerVSCluster has no create infrastructure annotation, setting cluster status to ready")
		clusterScope.IBMPowerVSCluster.Status.Ready = true
		return ctrl.Result{}, nil
	}

	// validate PER availability for the PowerVS zone, proceed further only if PowerVS zone support PER.
	// more information about PER can be found here: https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-per
	if err := clusterScope.IsPowerVSZoneSupportsPER(); err != nil {
		return reconcile.Result{}, fmt.Errorf("error checking PER capability for PowerVS zone: %w", err)
	}

	// reconcile resource group
	log.Info("Reconciling resource group")
	if err := clusterScope.ReconcileResourceGroup(ctx); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to reconcile resource group: %w", err)
	}

	powerVSCluster := &powerVSCluster{
		cluster: clusterScope.IBMPowerVSCluster,
	}

	var wg sync.WaitGroup
	ch := make(chan reconcileResult)

	// reconcile PowerVS resources
	wg.Add(1)
	go r.reconcilePowerVSResources(ctx, clusterScope, powerVSCluster, ch, &wg)

	// reconcile VPC
	wg.Add(1)
	go r.reconcileVPCResources(ctx, clusterScope, powerVSCluster, ch, &wg)

	// wait for above reconcile to complete and close the channel
	go func() {
		wg.Wait()
		close(ch)
	}()

	var requeue bool
	var errList []error
	// receive return values from the channel and decide the requeue
	for val := range ch {
		if val.Requeue {
			requeue = true
		}
		if val.error != nil {
			errList = append(errList, val.error)
		}
	}

	if requeue && len(errList) > 1 {
		return ctrl.Result{RequeueAfter: 30 * time.Second}, kerrors.NewAggregate(errList)
	} else if requeue {
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	} else if len(errList) > 1 {
		return ctrl.Result{}, kerrors.NewAggregate(errList)
	}

	// reconcile Transit Gateway
	log.Info("Reconciling transit gateway")
	if requeue, err := clusterScope.ReconcileTransitGateway(ctx); err != nil {
		conditions.MarkFalse(powerVSCluster.cluster, infrav1beta2.TransitGatewayReadyCondition, infrav1beta2.TransitGatewayReconciliationFailedReason, capiv1beta1.ConditionSeverityError, "%s", err.Error())
		v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
			Type:    infrav1beta2.TransitGatewayReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1beta2.TransitGatewayNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		return reconcile.Result{}, fmt.Errorf("failed to reconcile transit gateway: %w", err)
	} else if requeue {
		log.Info("Creating a transit gateway is pending, requeuing")
		return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
	}
	conditions.MarkTrue(powerVSCluster.cluster, infrav1beta2.TransitGatewayReadyCondition)
	v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
		Type:   infrav1beta2.TransitGatewayReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1beta2.TransitGatewayReadyV1Beta2Reason,
	})

	// reconcile COSInstance
	if clusterScope.IBMPowerVSCluster.Spec.Ignition != nil {
		log.Info("Reconciling COS service instance")
		if err := clusterScope.ReconcileCOSInstance(ctx); err != nil {
			conditions.MarkFalse(powerVSCluster.cluster, infrav1beta2.COSInstanceReadyCondition, infrav1beta2.COSInstanceReconciliationFailedReason, capiv1beta1.ConditionSeverityError, "%s", err.Error())
			v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
				Type:    infrav1beta2.COSInstanceReadyV1Beta2Condition,
				Status:  metav1.ConditionFalse,
				Reason:  infrav1beta2.COSInstanceNotReadyV1Beta2Reason,
				Message: err.Error(),
			})
			return reconcile.Result{}, fmt.Errorf("failed to reconcile COS instance: %w", err)
		}
		conditions.MarkTrue(powerVSCluster.cluster, infrav1beta2.COSInstanceReadyCondition)
		v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
			Type:   infrav1beta2.COSInstanceReadyV1Beta2Condition,
			Status: metav1.ConditionTrue,
			Reason: infrav1beta2.COSInstanceReadyV1Beta2Reason,
		})
	}

	var networkReady, loadBalancerReady bool
	for _, cond := range clusterScope.IBMPowerVSCluster.Status.Conditions {
		if cond.Type == infrav1beta2.NetworkReadyCondition && cond.Status == corev1.ConditionTrue {
			networkReady = true
		}
		if cond.Type == infrav1beta2.LoadBalancerReadyCondition && cond.Status == corev1.ConditionTrue {
			loadBalancerReady = true
		}
	}

	if !networkReady || !loadBalancerReady {
		log.Info("Network or LoadBalancer still not ready, requeuing")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Getting load balancer host")
	hostName, err := clusterScope.GetPublicLoadBalancerHostName()
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to fetch public loadbalancer: %w", err)
	}
	if hostName == nil || *hostName == "" {
		log.Info("LoadBalancer hostname is not yet available, requeuing")
		return reconcile.Result{RequeueAfter: time.Minute}, nil
	}

	// update cluster object with load balancer host name
	clusterScope.IBMPowerVSCluster.Spec.ControlPlaneEndpoint.Host = *hostName
	clusterScope.IBMPowerVSCluster.Spec.ControlPlaneEndpoint.Port = clusterScope.APIServerPort()
	clusterScope.IBMPowerVSCluster.Status.Ready = true
	return ctrl.Result{}, nil
}

func (r *IBMPowerVSClusterReconciler) reconcilePowerVSResources(ctx context.Context, clusterScope *scope.PowerVSClusterScope, powerVSCluster *powerVSCluster, ch chan reconcileResult, wg *sync.WaitGroup) {
	defer wg.Done()

	log := ctrl.LoggerFrom(ctx)
	log = log.WithName("powervs")

	log.Info("Reconciling PowerVS resources")
	defer log.Info("Finished Reconciling PowerVS resources")

	// reconcile PowerVS service instance
	log.Info("Reconciling PowerVS service instance")
	if requeue, err := clusterScope.ReconcilePowerVSServiceInstance(ctx); err != nil {
		powerVSCluster.updateCondition(capiv1beta1.Condition{
			Status:   corev1.ConditionFalse,
			Type:     infrav1beta2.ServiceInstanceReadyCondition,
			Reason:   infrav1beta2.ServiceInstanceReconciliationFailedReason,
			Severity: capiv1beta1.ConditionSeverityError,
			Message:  err.Error(),
		})
		//TODO: When we completely transition into v1beta2 api's update the conditions with lock
		v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
			Type:    infrav1beta2.WorkspaceReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1beta2.WorkspaceNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		ch <- reconcileResult{reconcile.Result{}, fmt.Errorf("failed to reconcile PowerVS service instance: %w", err)}
		return
	} else if requeue {
		log.Info("PowerVS service instance creation is pending, requeuing")
		ch <- reconcileResult{reconcile.Result{Requeue: true}, nil}
		return
	}
	powerVSCluster.updateCondition(capiv1beta1.Condition{
		Status: corev1.ConditionTrue,
		Type:   infrav1beta2.ServiceInstanceReadyCondition,
	})
	v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
		Type:   infrav1beta2.WorkspaceReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1beta2.WorkspaceReadyV1Beta2Reason,
	})

	clusterScope.IBMPowerVSClient.WithClients(powervs.ServiceOptions{CloudInstanceID: clusterScope.GetServiceInstanceID()})

	// reconcile network
	log.Info("Reconciling network")
	if networkActive, err := clusterScope.ReconcileNetwork(ctx); err != nil {
		powerVSCluster.updateCondition(capiv1beta1.Condition{
			Status:   corev1.ConditionFalse,
			Type:     infrav1beta2.NetworkReadyCondition,
			Reason:   infrav1beta2.NetworkReconciliationFailedReason,
			Severity: capiv1beta1.ConditionSeverityError,
			Message:  err.Error(),
		})
		v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
			Type:    infrav1beta2.NetworkReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1beta2.NetworkNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		ch <- reconcileResult{reconcile.Result{}, fmt.Errorf("failed to reconcile network: %w", err)}
		return
	} else if networkActive {
		powerVSCluster.updateCondition(capiv1beta1.Condition{
			Status: corev1.ConditionTrue,
			Type:   infrav1beta2.NetworkReadyCondition,
		})
		v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
			Type:   infrav1beta2.NetworkReadyV1Beta2Condition,
			Status: metav1.ConditionTrue,
			Reason: infrav1beta2.NetworkReadyV1Beta2Reason,
		})
		return
	}
	// Do not want to block the reconciliation of other resources like setting up TG and COS, so skipping the requeue and only logging the info.
	log.Info("PowerVS network creation is pending")
}

func (r *IBMPowerVSClusterReconciler) reconcileVPCResources(ctx context.Context, clusterScope *scope.PowerVSClusterScope, powerVSCluster *powerVSCluster, ch chan reconcileResult, wg *sync.WaitGroup) {
	defer wg.Done()

	log := ctrl.LoggerFrom(ctx)
	log = log.WithName("vpc")

	log.Info("Reconciling VPC")
	defer log.Info("Finished VPC reconciliation")

	if requeue, err := clusterScope.ReconcileVPC(ctx); err != nil {
		powerVSCluster.updateCondition(capiv1beta1.Condition{
			Status:   corev1.ConditionFalse,
			Type:     infrav1beta2.VPCReadyCondition,
			Reason:   infrav1beta2.VPCReconciliationFailedReason,
			Severity: capiv1beta1.ConditionSeverityError,
			Message:  err.Error(),
		})
		v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
			Type:    infrav1beta2.VPCReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1beta2.VPCNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		ch <- reconcileResult{reconcile.Result{}, fmt.Errorf("failed to reconcile VPC: %w", err)}
		return
	} else if requeue {
		log.Info("VPC creation is pending, requeuing")
		ch <- reconcileResult{reconcile.Result{Requeue: true}, nil}
		return
	}
	powerVSCluster.updateCondition(capiv1beta1.Condition{
		Status: corev1.ConditionTrue,
		Type:   infrav1beta2.VPCReadyCondition,
	})
	v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
		Type:   infrav1beta2.VPCReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1beta2.VPCReadyV1Beta2Reason,
	})

	// reconcile VPC Subnet
	log.Info("Reconciling VPC subnets")
	if requeue, err := clusterScope.ReconcileVPCSubnets(ctx); err != nil {
		powerVSCluster.updateCondition(capiv1beta1.Condition{
			Status:   corev1.ConditionFalse,
			Type:     infrav1beta2.VPCSubnetReadyCondition,
			Reason:   infrav1beta2.VPCSubnetReconciliationFailedReason,
			Severity: capiv1beta1.ConditionSeverityError,
			Message:  err.Error(),
		})
		v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
			Type:    infrav1beta2.VPCSubnetReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1beta2.VPCSubnetNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		ch <- reconcileResult{reconcile.Result{}, fmt.Errorf("failed to reconcile VPC subnets: %w", err)}
		return
	} else if requeue {
		log.Info("VPC subnet creation is pending, requeuing")
		ch <- reconcileResult{reconcile.Result{Requeue: true}, nil}
		return
	}
	powerVSCluster.updateCondition(capiv1beta1.Condition{
		Status: corev1.ConditionTrue,
		Type:   infrav1beta2.VPCSubnetReadyCondition,
	})
	v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
		Type:   infrav1beta2.VPCSubnetReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1beta2.VPCSubnetReadyV1Beta2Reason,
	})

	// reconcile VPC security group
	log.Info("Reconciling VPC security group")
	if err := clusterScope.ReconcileVPCSecurityGroups(ctx); err != nil {
		powerVSCluster.updateCondition(capiv1beta1.Condition{
			Status:   corev1.ConditionFalse,
			Type:     infrav1beta2.VPCSecurityGroupReadyCondition,
			Reason:   infrav1beta2.VPCSecurityGroupReconciliationFailedReason,
			Severity: capiv1beta1.ConditionSeverityError,
			Message:  err.Error(),
		})
		v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
			Type:    infrav1beta2.VPCSecurityGroupReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1beta2.VPCSecurityGroupNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		ch <- reconcileResult{reconcile.Result{}, fmt.Errorf("failed to reconcile VPC security groups: %w", err)}
		return
	}
	powerVSCluster.updateCondition(capiv1beta1.Condition{
		Status: corev1.ConditionTrue,
		Type:   infrav1beta2.VPCSecurityGroupReadyCondition,
	})
	v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
		Type:   infrav1beta2.VPCSecurityGroupReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1beta2.VPCSecurityGroupReadyV1Beta2Reason,
	})

	// reconcile LoadBalancer
	log.Info("Reconciling VPC load balancers")
	if loadBalancerReady, err := clusterScope.ReconcileLoadBalancers(ctx); err != nil {
		powerVSCluster.updateCondition(capiv1beta1.Condition{
			Status:   corev1.ConditionFalse,
			Type:     infrav1beta2.LoadBalancerReadyCondition,
			Reason:   infrav1beta2.LoadBalancerReconciliationFailedReason,
			Severity: capiv1beta1.ConditionSeverityError,
			Message:  err.Error(),
		})
		v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
			Type:    infrav1beta2.VPCLoadBalancerReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1beta2.VPCLoadBalancerNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		ch <- reconcileResult{reconcile.Result{}, fmt.Errorf("failed to reconcile VPC load balancers: %w", err)}
		return
	} else if loadBalancerReady {
		powerVSCluster.updateCondition(capiv1beta1.Condition{
			Status: corev1.ConditionTrue,
			Type:   infrav1beta2.LoadBalancerReadyCondition,
		})
		v1beta2conditions.Set(powerVSCluster.cluster, metav1.Condition{
			Type:   infrav1beta2.VPCLoadBalancerReadyV1Beta2Condition,
			Status: metav1.ConditionTrue,
			Reason: infrav1beta2.VPCLoadBalancerReadyV1Beta2Reason,
		})
		return
	}
	// Do not want to block the reconciliation of other resources like setting up TG and COS, so skipping the requeue and only logging the info.
	log.Info("VPC load balancer creation is pending")
}

func (r *IBMPowerVSClusterReconciler) reconcileDelete(ctx context.Context, clusterScope *scope.PowerVSClusterScope) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	log.Info("Reconciling IBMPowerVSCluster delete ")
	defer log.Info("Finished reconciling IBMPowerVSCluster delete")

	cluster := clusterScope.IBMPowerVSCluster

	if result, err := r.deleteIBMPowerVSImage(ctx, clusterScope); err != nil || !result.IsZero() {
		return result, err
	}

	// check for annotation set for cluster resource and decide on proceeding with infra deletion.
	if !scope.CheckCreateInfraAnnotation(*clusterScope.IBMPowerVSCluster) {
		log.Info("IBMPowerVSCluster has no infra annotation, removing finalizer")
		controllerutil.RemoveFinalizer(cluster, infrav1beta2.IBMPowerVSClusterFinalizer)
		return ctrl.Result{}, nil
	}

	var allErrs []error
	clusterScope.IBMPowerVSClient.WithClients(powervs.ServiceOptions{CloudInstanceID: clusterScope.GetServiceInstanceID()})

	log.Info("Deleting transit gateway")
	v1beta2conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1beta2.TransitGatewayReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: infrav1beta2.TransitGatewayDeletingV1Beta2Reason,
	})
	if requeue, err := clusterScope.DeleteTransitGateway(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete transit gateway: %w", err))
	} else if requeue {
		log.Info("Transit gateway deletion is pending, requeuing")
		return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	log.Info("Deleting VPC load balancer")
	v1beta2conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1beta2.VPCLoadBalancerReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: infrav1beta2.VPCLoadBalancerDeletingV1Beta2Reason,
	})
	if requeue, err := clusterScope.DeleteLoadBalancer(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete VPC load balancer: %w", err))
	} else if requeue {
		log.Info("VPC load balancer deletion is pending, requeuing")
		return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	log.Info("Deleting VPC security group")
	v1beta2conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1beta2.VPCSecurityGroupReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: infrav1beta2.VPCSecurityGroupDeletingV1Beta2Reason,
	})
	if err := clusterScope.DeleteVPCSecurityGroups(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete VPC security group: %w", err))
	}

	log.Info("Deleting VPC subnet")
	v1beta2conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1beta2.VPCSubnetReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: infrav1beta2.VPCSubnetDeletingV1Beta2Reason,
	})
	if requeue, err := clusterScope.DeleteVPCSubnet(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete VPC subnet: %w", err))
	} else if requeue {
		log.Info("VPC subnet deletion is pending, requeuing")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}

	log.Info("Deleting VPC")
	v1beta2conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1beta2.VPCReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: infrav1beta2.VPCDeletingV1Beta2Reason,
	})
	if requeue, err := clusterScope.DeleteVPC(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete VPC: %w", err))
	} else if requeue {
		log.Info("VPC deletion is pending, requeuing")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}

	log.Info("Deleting DHCP server")
	v1beta2conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1beta2.NetworkReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: infrav1beta2.NetworkDeletingV1Beta2Reason,
	})
	if err := clusterScope.DeleteDHCPServer(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete DHCP server: %w", err))
	}

	log.Info("Deleting PowerVS service instance")
	v1beta2conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1beta2.WorkspaceReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: infrav1beta2.WorkspaceDeletingV1Beta2Reason,
	})
	if requeue, err := clusterScope.DeleteServiceInstance(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete PowerVS service instance: %w", err))
	} else if requeue {
		log.Info("PowerVS service instance deletion is pending, requeuing")
		return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	if clusterScope.IBMPowerVSCluster.Spec.Ignition != nil {
		v1beta2conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
			Type:   infrav1beta2.COSInstanceReadyV1Beta2Condition,
			Status: metav1.ConditionFalse,
			Reason: infrav1beta2.COSInstanceDeletingV1Beta2Reason,
		})
		log.Info("Deleting COS service instance")
		if err := clusterScope.DeleteCOSInstance(ctx); err != nil {
			allErrs = append(allErrs, fmt.Errorf("failed to delete COS service instance: %w", err))
		}
	}

	if len(allErrs) > 0 {
		return ctrl.Result{}, kerrors.NewAggregate(allErrs)
	}

	log.Info("IBMPowerVSCluster deletion completed")
	controllerutil.RemoveFinalizer(cluster, infrav1beta2.IBMPowerVSClusterFinalizer)
	return ctrl.Result{}, nil
}

func (update *powerVSCluster) updateCondition(condition capiv1beta1.Condition) {
	update.mu.Lock()
	defer update.mu.Unlock()
	conditions.Set(update.cluster, &condition)
}

func (r *IBMPowerVSClusterReconciler) deleteIBMPowerVSImage(ctx context.Context, clusterScope *scope.PowerVSClusterScope) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	cluster := clusterScope.IBMPowerVSCluster
	descendants, err := r.listDescendants(ctx, cluster)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list descendants: %w", err)
	}

	// since we are avoiding using cache for IBMPowerVSCluster the Type meta of the retrieved object will be empty
	// explicitly setting here to filter children
	if gvk := cluster.GetObjectKind().GroupVersionKind(); gvk.Empty() {
		gvk, err := r.GroupVersionKindFor(cluster)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to get GVK of cluster: %w", err)
		}
		cluster.SetGroupVersionKind(gvk)
	}

	children, err := descendants.filterOwnedDescendants(cluster)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to filter owned descendants: %w", err)
	}

	if len(children) > 0 {
		log.Info("Cluster still has children - deleting them first", "count", len(children))

		var errs []error

		for _, child := range children {
			if !child.GetDeletionTimestamp().IsZero() {
				// Don't handle deleted child.
				continue
			}
			gvk := child.GetObjectKind().GroupVersionKind().String()

			log.Info("Deleting child object", "gvk", gvk, "name", child.GetName())
			if err := r.Client.Delete(ctx, child); err != nil {
				err = fmt.Errorf("error deleting child object %s: %w", child.GetName(), err)
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return ctrl.Result{}, kerrors.NewAggregate(errs)
		}
	}

	if descendantCount := descendants.length(); descendantCount > 0 {
		indirect := descendantCount - len(children)
		log.Info("Cluster still has descendants - need to requeue", "descendants", descendants.descendantNames(), "indirectDescendantsCount", indirect)
		// Requeue so we can check the next time to see if there are still any descendants left.
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}
	return ctrl.Result{}, nil
}

type clusterDescendants struct {
	ibmPowerVSImages infrav1beta2.IBMPowerVSImageList
}

// length returns the number of descendants.
func (c *clusterDescendants) length() int {
	return len(c.ibmPowerVSImages.Items)
}

func (c *clusterDescendants) descendantNames() string {
	descendants := make([]string, 0)
	ibmPowerVSImageNames := make([]string, len(c.ibmPowerVSImages.Items))
	for i, ibmPowerVSImage := range c.ibmPowerVSImages.Items {
		ibmPowerVSImageNames[i] = ibmPowerVSImage.Name
	}
	if len(ibmPowerVSImageNames) > 0 {
		descendants = append(descendants, "IBM Powervs Images: "+strings.Join(ibmPowerVSImageNames, ","))
	}

	return strings.Join(descendants, ";")
}

// listDescendants returns a list of all IBMPowerVSImages for the cluster.
func (r *IBMPowerVSClusterReconciler) listDescendants(ctx context.Context, cluster *infrav1beta2.IBMPowerVSCluster) (clusterDescendants, error) {
	var descendants clusterDescendants

	listOptions := []client.ListOption{
		client.InNamespace(cluster.Namespace),
		client.MatchingLabels(map[string]string{capiv1beta1.ClusterNameLabel: cluster.Name}),
	}

	if err := r.Client.List(ctx, &descendants.ibmPowerVSImages, listOptions...); err != nil {
		return descendants, fmt.Errorf("failed to list IBMPowerVSImages for cluster %s/%s: %w", cluster.Namespace, cluster.Name, err)
	}

	return descendants, nil
}

// filterOwnedDescendants returns an array of runtime.Objects containing only those descendants that have the cluster
// as an owner reference.
func (c *clusterDescendants) filterOwnedDescendants(cluster *infrav1beta2.IBMPowerVSCluster) ([]client.Object, error) {
	var ownedDescendants []client.Object
	eachFunc := func(o runtime.Object) error {
		obj := o.(client.Object)
		acc, err := meta.Accessor(obj)
		if err != nil {
			return nil //nolint:nilerr // We don't want to exit the EachListItem loop, just continue
		}

		if util.IsOwnedByObject(acc, cluster) {
			ownedDescendants = append(ownedDescendants, obj)
		}

		return nil
	}

	lists := []client.ObjectList{
		&c.ibmPowerVSImages,
	}

	for _, list := range lists {
		if err := meta.EachListItem(list, eachFunc); err != nil {
			return nil, fmt.Errorf("error finding owned descendants of cluster %s/%s: %w", cluster.Namespace, cluster.Name, err)
		}
	}

	return ownedDescendants, nil
}

// patchIBMPowerVSCluster updates the IBMPowerVSCluster and its status on the API server.
func patchIBMPowerVSCluster(ctx context.Context, patchHelper *patch.Helper, ibmPowerVSCluster *infrav1beta2.IBMPowerVSCluster) error {
	// we don't need to set any conditions for IBMPowerVSCluster without create infra annotation.
	if !scope.CheckCreateInfraAnnotation(*ibmPowerVSCluster) {
		if err := patchHelper.Patch(ctx, ibmPowerVSCluster); err != nil {
			return fmt.Errorf("error patching IBMPowerVSCluster: %w", err)
		}
		return nil
	}

	if err := v1beta2conditions.SetSummaryCondition(ibmPowerVSCluster, ibmPowerVSCluster, infrav1beta2.IBMPowerVSClusterReadyV1Beta2Condition,
		v1beta2conditions.ForConditionTypes{
			infrav1beta2.WorkspaceReadyV1Beta2Condition,
			infrav1beta2.NetworkReadyV1Beta2Condition,
			infrav1beta2.VPCReadyV1Beta2Condition,
			infrav1beta2.VPCSubnetReadyV1Beta2Condition,
			infrav1beta2.VPCSecurityGroupReadyV1Beta2Condition,
			infrav1beta2.VPCLoadBalancerReadyV1Beta2Condition,
			infrav1beta2.TransitGatewayReadyV1Beta2Condition,
			infrav1beta2.COSInstanceReadyV1Beta2Condition,
		},
		v1beta2conditions.IgnoreTypesIfMissing{
			infrav1beta2.COSInstanceReadyV1Beta2Condition,
		},
		// Using a custom merge strategy to override reasons applied during merge.
		v1beta2conditions.CustomMergeStrategy{
			MergeStrategy: v1beta2conditions.DefaultMergeStrategy(
				// Use custom reasons.
				v1beta2conditions.ComputeReasonFunc(v1beta2conditions.GetDefaultComputeMergeReasonFunc(
					infrav1beta2.IBMPowerVSClusterNotReadyV1Beta2Reason,
					infrav1beta2.IBMPowerVSClusterReadyUnknownV1Beta2Reason,
					infrav1beta2.IBMPowerVSClusterReadyV1Beta2Reason,
				)),
			),
		},
	); err != nil {
		return fmt.Errorf("failed to set %s condition: %w", infrav1beta2.IBMPowerVSClusterReadyV1Beta2Condition, err)
	}

	return patchHelper.Patch(ctx, ibmPowerVSCluster,
		patch.WithOwnedV1Beta2Conditions{Conditions: []string{
			capiv1beta1.PausedV1Beta2Condition,
			infrav1beta2.IBMPowerVSClusterReadyV1Beta2Condition,
			infrav1beta2.WorkspaceReadyV1Beta2Condition,
			infrav1beta2.NetworkReadyV1Beta2Condition,
			infrav1beta2.VPCReadyV1Beta2Condition,
			infrav1beta2.VPCSubnetReadyV1Beta2Condition,
			infrav1beta2.VPCSecurityGroupReadyV1Beta2Condition,
			infrav1beta2.TransitGatewayReadyV1Beta2Condition,
			infrav1beta2.COSInstanceReadyV1Beta2Condition,
		}},
	)
}

// SetupWithManager creates a new IBMPowerVSCluster controller for a manager.
func (r *IBMPowerVSClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	predicateLog := ctrl.LoggerFrom(ctx).WithValues("controller", "ibmpowervscluster")
	err := ctrl.NewControllerManagedBy(mgr).
		For(&infrav1beta2.IBMPowerVSCluster{}).
		WithEventFilter(predicates.ResourceHasFilterLabel(r.Scheme, predicateLog, r.WatchFilterValue)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(r.Scheme, predicateLog)).
		Watches(
			&capiv1beta1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(ctx, infrav1beta2.GroupVersion.WithKind("IBMPowerVSCluster"), mgr.GetClient(), &infrav1beta2.IBMPowerVSCluster{})),
			builder.WithPredicates(predicates.All(r.Scheme, predicateLog,
				predicates.ResourceIsChanged(r.Scheme, predicateLog),
				predicates.ClusterPausedTransitions(r.Scheme, predicateLog),
			)),
		).Complete(r)
	if err != nil {
		return fmt.Errorf("could not set up controller for IBMPowerVSCluster: %w", err)
	}
	return nil
}
