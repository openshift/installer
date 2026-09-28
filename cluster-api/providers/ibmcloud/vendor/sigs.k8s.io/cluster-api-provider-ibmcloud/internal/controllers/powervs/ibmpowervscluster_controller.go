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

package powervs

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
	"k8s.io/utils/ptr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	deprecatedv1beta1conditions "sigs.k8s.io/cluster-api/util/conditions/deprecated/v1beta1"
	"sigs.k8s.io/cluster-api/util/finalizers"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/paused"
	"sigs.k8s.io/cluster-api/util/predicates"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/endpoints"
	powervsscope "sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/scope/powervs"
)

const (
	deprecatedConditionsField = "conditions"
	deprecatedStatus          = "deprecated"
	statusField               = "status"
	v1beta2Version            = "v1beta2"
)

// IBMPowerVSClusterReconciler reconciles a IBMPowerVSCluster object.
type IBMPowerVSClusterReconciler struct {
	client.Client
	Recorder        record.EventRecorder
	ServiceEndpoint []endpoints.ServiceEndpoint
	Scheme          *runtime.Scheme

	// WatchFilterValue is the label value used to filter events prior to reconciliation.
	WatchFilterValue string

	ClientBuilder powervsscope.ClientBuilder
}

// componentResult holds the outcome of a concurrent component reconciliation.
type componentResult struct {
	requeue    bool
	err        error
	conditions []metav1.Condition
	legacy     []*clusterv1.Condition
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=create;delete;get;list;watch

// Reconcile implements controller runtime Reconciler interface and handles reconcileation logic for IBMPowerVSCluster.
func (r *IBMPowerVSClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	log.Info("Reconciling IBMPowerVSCluster")
	defer log.Info("Finished reconciling IBMPowerVSCluster")

	// Fetch the IBMPowerVSCluster instance.
	ibmPowerVSCluster := &infrav1.IBMPowerVSCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, ibmPowerVSCluster); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("IBMPowerVSCluster not found")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get IBMPowerVSCluster: %w", err)
	}

	// Add finalizer first if not set to avoid the race condition between init and delete.
	if finalizerAdded, err := finalizers.EnsureFinalizer(ctx, r.Client, ibmPowerVSCluster, infrav1.IBMPowerVSClusterFinalizer); err != nil || finalizerAdded {
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
	clusterScope, err := powervsscope.NewPowerVSClusterScope(ctx, powervsscope.ClusterScopeParams{
		Client:            r.Client,
		Cluster:           cluster,
		IBMPowerVSCluster: ibmPowerVSCluster,
		ServiceEndpoint:   r.ServiceEndpoint,
		ClientBuilder:     r.ClientBuilder,
		Recorder:          r.Recorder,
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

func (r *IBMPowerVSClusterReconciler) reconcile(ctx context.Context, clusterScope *powervsscope.ClusterScope) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// If it's VirtualIP, do the minimal logic and return early.
	if clusterScope.IBMPowerVSCluster.Spec.Topology == infrav1.PowerVSVirtualIPTopology {
		log.Info("Reconciling in VirtualIP topology mode")
		clusterScope.IBMPowerVSCluster.Status.Initialization.Provisioned = ptr.To(true)
		return ctrl.Result{}, nil
	}

	// Otherwise, assume LoadBalancer and proceed with the heavy VPC/TG logic.
	if clusterScope.IBMPowerVSCluster.Spec.Topology != infrav1.PowerVSLoadBalancerTopology {
		return ctrl.Result{}, fmt.Errorf("unknown topology: %q", clusterScope.IBMPowerVSCluster.Spec.Topology)
	}

	// validate PER availability for the PowerVS zone, proceed further only if PowerVS zone support PER.
	// more information about PER can be found here: https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-per
	if err := clusterScope.ValidateZoneSupportsPER(ctx); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to validate PER capability for PowerVS zone: %w", err)
	}

	// reconcile resource group
	log.Info("Reconciling resource group")
	if err := clusterScope.ReconcileResourceGroup(ctx); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to reconcile resource group: %w", err)
	}

	var wg sync.WaitGroup
	ch := make(chan componentResult, 2)

	wg.Go(func() {
		ch <- r.reconcilePowerVSResources(ctx, clusterScope)
	})

	wg.Go(func() {
		ch <- r.reconcileVPCResources(ctx, clusterScope)
	})

	wg.Wait()
	close(ch)

	var errList []error
	var needsRequeue bool

	for res := range ch {
		if res.err != nil {
			errList = append(errList, res.err)
		}
		if res.requeue {
			needsRequeue = true
		}

		for i := range res.conditions {
			conditions.Set(clusterScope.IBMPowerVSCluster, res.conditions[i])
			deprecatedv1beta1conditions.Set(clusterScope.IBMPowerVSCluster, res.legacy[i])
		}
	}

	if len(errList) > 0 {
		return ctrl.Result{}, kerrors.NewAggregate(errList)
	}
	if needsRequeue {
		log.Info("PowerVS or VPC infrastructure components are still provisioning, requeuing before proceeding to Transit Gateway")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Reconciling transit gateway")
	if requeue, err := clusterScope.ReconcileTransitGateway(ctx); err != nil {
		condition, legacyCondition := r.buildConditions(infrav1.TransitGatewayReadyCondition, infrav1.TransitGatewayReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.TransitGatewayNotReadyReason, infrav1.TransitGatewayReconciliationFailedV1Beta2Reason, err.Error())
		conditions.Set(clusterScope.IBMPowerVSCluster, condition)
		deprecatedv1beta1conditions.Set(clusterScope.IBMPowerVSCluster, legacyCondition)
		return reconcile.Result{}, fmt.Errorf("failed to reconcile transit gateway: %w", err)
	} else if requeue {
		log.Info("Creating a transit gateway is pending, requeuing")
		return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	condition, legacyCondition := r.buildConditions(infrav1.TransitGatewayReadyCondition, infrav1.TransitGatewayReadyV1Beta2Condition, metav1.ConditionTrue, infrav1.TransitGatewayReadyReason, "", "")
	conditions.Set(clusterScope.IBMPowerVSCluster, condition)
	deprecatedv1beta1conditions.Set(clusterScope.IBMPowerVSCluster, legacyCondition)

	if clusterScope.IBMPowerVSCluster.Spec.COSInstance.Type != "" {
		log.Info("Reconciling COS service instance")
		if err := clusterScope.ReconcileCOSInstance(ctx); err != nil {
			condition, legacyCondition := r.buildConditions(infrav1.COSInstanceReadyCondition, infrav1.COSInstanceReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.COSInstanceNotReadyReason, infrav1.COSInstanceReconciliationFailedV1Beta2Reason, err.Error())
			conditions.Set(clusterScope.IBMPowerVSCluster, condition)
			deprecatedv1beta1conditions.Set(clusterScope.IBMPowerVSCluster, legacyCondition)
			return reconcile.Result{}, fmt.Errorf("failed to reconcile COS instance: %w", err)
		}
		condition, legacyCondition := r.buildConditions(infrav1.COSInstanceReadyCondition, infrav1.COSInstanceReadyV1Beta2Condition, metav1.ConditionTrue, infrav1.COSInstanceReadyReason, "", "")
		conditions.Set(clusterScope.IBMPowerVSCluster, condition)
		deprecatedv1beta1conditions.Set(clusterScope.IBMPowerVSCluster, legacyCondition)
	}

	if !conditions.IsTrue(clusterScope.IBMPowerVSCluster, infrav1.NetworkReadyCondition) ||
		!conditions.IsTrue(clusterScope.IBMPowerVSCluster, infrav1.VPCLoadBalancerReadyCondition) {
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

	clusterScope.IBMPowerVSCluster.Spec.ControlPlaneEndpoint.Host = *hostName
	clusterScope.IBMPowerVSCluster.Spec.ControlPlaneEndpoint.Port = clusterScope.APIServerPort()
	clusterScope.IBMPowerVSCluster.Status.Initialization.Provisioned = ptr.To(true)

	return ctrl.Result{}, nil
}

func (r *IBMPowerVSClusterReconciler) reconcilePowerVSResources(ctx context.Context, clusterScope *powervsscope.ClusterScope) componentResult {
	log := ctrl.LoggerFrom(ctx).WithName("powervs")
	res := componentResult{}

	log.Info("Reconciling PowerVS resources")
	defer log.Info("Finished Reconciling PowerVS resources")

	log.Info("Reconciling PowerVS workspace")
	if requeue, err := clusterScope.ReconcileWorkspace(ctx); err != nil {
		condition, legacyCondition := r.buildConditions(infrav1.WorkspaceReadyCondition, infrav1.ServiceInstanceReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.WorkspaceNotReadyReason, infrav1.ServiceInstanceReconciliationFailedV1Beta2Reason, err.Error())
		res.conditions = append(res.conditions, condition)
		res.legacy = append(res.legacy, legacyCondition)
		res.err = fmt.Errorf("failed to reconcile PowerVS workspace: %w", err)
		return res
	} else if requeue {
		log.Info("PowerVS workspace creation is pending")
		res.requeue = true
		return res
	}
	condition, legacyCondition := r.buildConditions(infrav1.WorkspaceReadyCondition, infrav1.ServiceInstanceReadyV1Beta2Condition, metav1.ConditionTrue, infrav1.WorkspaceReadyReason, "", "")
	res.conditions = append(res.conditions, condition)
	res.legacy = append(res.legacy, legacyCondition)

	log.Info("Reconciling network")
	if requeue, err := clusterScope.ReconcileNetwork(ctx); err != nil {
		condition, legacyCondition := r.buildConditions(infrav1.NetworkReadyCondition, infrav1.NetworkReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.NetworkNotReadyReason, infrav1.NetworkReconciliationFailedV1Beta2Reason, err.Error())
		res.conditions = append(res.conditions, condition)
		res.legacy = append(res.legacy, legacyCondition)
		res.err = fmt.Errorf("failed to reconcile network: %w", err)
		return res
	} else if requeue {
		log.Info("PowerVS network creation is pending")
	} else {
		condition, legacyCondition = r.buildConditions(infrav1.NetworkReadyCondition, infrav1.NetworkReadyV1Beta2Condition, metav1.ConditionTrue, infrav1.NetworkReadyReason, "", "")
		res.conditions = append(res.conditions, condition)
		res.legacy = append(res.legacy, legacyCondition)
	}

	return res
}

func (r *IBMPowerVSClusterReconciler) reconcileVPCResources(ctx context.Context, clusterScope *powervsscope.ClusterScope) componentResult {
	log := ctrl.LoggerFrom(ctx).WithName("vpc")
	res := componentResult{}

	log.Info("Reconciling VPC")
	if requeue, err := clusterScope.ReconcileVPC(ctx); err != nil {
		condition, legacyCondition := r.buildConditions(infrav1.VPCReadyCondition, infrav1.VPCReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.VPCNotReadyReason, infrav1.VPCReconciliationFailedV1Beta2Reason, err.Error())
		res.conditions = append(res.conditions, condition)
		res.legacy = append(res.legacy, legacyCondition)
		res.err = fmt.Errorf("failed to reconcile VPC: %w", err)
		return res
	} else if requeue {
		log.Info("VPC creation is pending")
		res.requeue = true
		return res
	}
	condition, legacyCondition := r.buildConditions(infrav1.VPCReadyCondition, infrav1.VPCReadyV1Beta2Condition, metav1.ConditionTrue, infrav1.VPCReadyReason, "", "")
	res.conditions = append(res.conditions, condition)
	res.legacy = append(res.legacy, legacyCondition)

	// reconcile VPC Subnet
	log.Info("Reconciling VPC subnets")
	if requeue, err := clusterScope.ReconcileVPCSubnets(ctx); err != nil {
		condition, legacyCondition := r.buildConditions(infrav1.VPCSubnetReadyCondition, infrav1.VPCSubnetReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.VPCSubnetNotReadyReason, infrav1.VPCSubnetReconciliationFailedV1Beta2Reason, err.Error())
		res.conditions = append(res.conditions, condition)
		res.legacy = append(res.legacy, legacyCondition)
		res.err = fmt.Errorf("failed to reconcile VPC subnets: %w", err)
		return res
	} else if requeue {
		log.Info("VPC subnet creation is pending")
		res.requeue = true
		return res
	}
	condition, legacyCondition = r.buildConditions(infrav1.VPCSubnetReadyCondition, infrav1.VPCSubnetReadyV1Beta2Condition, metav1.ConditionTrue, infrav1.VPCSubnetReadyReason, "", "")
	res.conditions = append(res.conditions, condition)
	res.legacy = append(res.legacy, legacyCondition)

	// reconcile VPC security group
	log.Info("Reconciling VPC security group")
	if err := clusterScope.ReconcileVPCSecurityGroups(ctx); err != nil {
		condition, legacyCondition := r.buildConditions(infrav1.VPCSecurityGroupReadyCondition, infrav1.VPCSecurityGroupReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.VPCSecurityGroupReconciliationFailedReason, infrav1.VPCSecurityGroupReconciliationFailedV1Beta2Reason, err.Error())
		res.conditions = append(res.conditions, condition)
		res.legacy = append(res.legacy, legacyCondition)
		res.err = fmt.Errorf("failed to reconcile VPC security groups: %w", err)
		return res
	}
	condition, legacyCondition = r.buildConditions(infrav1.VPCSecurityGroupReadyCondition, infrav1.VPCSecurityGroupReadyV1Beta2Condition, metav1.ConditionTrue, infrav1.VPCSecurityGroupReadyReason, "", "")
	res.conditions = append(res.conditions, condition)
	res.legacy = append(res.legacy, legacyCondition)

	// reconcile LoadBalancer
	log.Info("Reconciling VPC load balancers")
	if loadBalancerReady, err := clusterScope.ReconcileLoadBalancers(ctx); err != nil {
		condition, legacyCondition := r.buildConditions(infrav1.VPCLoadBalancerReadyCondition, infrav1.LoadBalancerReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.VPCLoadBalancerNotReadyReason, infrav1.LoadBalancerReconciliationFailedV1Beta2Reason, err.Error())
		res.conditions = append(res.conditions, condition)
		res.legacy = append(res.legacy, legacyCondition)
		res.err = fmt.Errorf("failed to reconcile VPC load balancers: %w", err)
		return res
	} else if !loadBalancerReady {
		log.Info("VPC load balancer creation is pending")
		// Not blocking here.
	} else {
		condition, legacyCondition = r.buildConditions(infrav1.VPCLoadBalancerReadyCondition, infrav1.LoadBalancerReadyV1Beta2Condition, metav1.ConditionTrue, infrav1.VPCLoadBalancerReadyReason, "", "")
		res.conditions = append(res.conditions, condition)
		res.legacy = append(res.legacy, legacyCondition)
	}

	return res
}

func (r *IBMPowerVSClusterReconciler) reconcileDelete(ctx context.Context, clusterScope *powervsscope.ClusterScope) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	log.Info("Reconciling IBMPowerVSCluster delete ")
	defer log.Info("Finished reconciling IBMPowerVSCluster delete")

	cluster := clusterScope.IBMPowerVSCluster

	if result, err := r.deleteIBMPowerVSImage(ctx, clusterScope); err != nil || !result.IsZero() {
		return result, err
	}

	// Check the cluster topology to decide if we need to proceed with VPC/TransitGateway infra deletion.
	if clusterScope.IBMPowerVSCluster.Spec.Topology != infrav1.PowerVSLoadBalancerTopology {
		log.Info("IBMPowerVSCluster is not in LoadBalancer topology mode, skipping advanced infra deletion and removing finalizer")
		controllerutil.RemoveFinalizer(cluster, infrav1.IBMPowerVSClusterFinalizer)
		return ctrl.Result{}, nil
	}

	var allErrs []error

	log.Info("Deleting transit gateway")
	conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1.TransitGatewayReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrav1.TransitGatewayDeletingReason,
	})
	if requeue, err := clusterScope.DeleteTransitGateway(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete transit gateway: %w", err))
	} else if requeue {
		log.Info("Transit gateway deletion is pending, requeuing")
		return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	log.Info("Deleting VPC load balancer")
	conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1.VPCLoadBalancerReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrav1.VPCLoadBalancerDeletingReason,
	})
	if requeue, err := clusterScope.DeleteLoadBalancer(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete VPC load balancer: %w", err))
	} else if requeue {
		log.Info("VPC load balancer deletion is pending, requeuing")
		return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	log.Info("Deleting VPC security group")
	conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1.VPCSecurityGroupReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrav1.VPCSecurityGroupDeletingReason,
	})
	if err := clusterScope.DeleteVPCSecurityGroups(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete VPC security group: %w", err))
	}

	log.Info("Deleting VPC subnet")
	conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1.VPCSubnetReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrav1.VPCSubnetDeletingReason,
	})
	if requeue, err := clusterScope.DeleteVPCSubnets(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete VPC subnet: %w", err))
	} else if requeue {
		log.Info("VPC subnet deletion is pending, requeuing")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}

	log.Info("Deleting VPC")
	conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1.VPCReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrav1.VPCDeletingReason,
	})
	if requeue, err := clusterScope.DeleteVPC(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete VPC: %w", err))
	} else if requeue {
		log.Info("VPC deletion is pending, requeuing")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}

	log.Info("Deleting DHCP server")
	conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1.NetworkReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrav1.NetworkDeletingReason,
	})
	if err := clusterScope.DeleteDHCPServer(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete DHCP server: %w", err))
	}

	log.Info("Deleting PowerVS workspace")
	conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
		Type:   infrav1.WorkspaceReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrav1.WorkspaceDeletingReason,
	})
	if requeue, err := clusterScope.DeleteWorkspace(ctx); err != nil {
		allErrs = append(allErrs, fmt.Errorf("failed to delete PowerVS workspace: %w", err))
	} else if requeue {
		log.Info("PowerVS workspace deletion is pending, requeuing")
		return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	if clusterScope.IBMPowerVSCluster.Spec.COSInstance.Type != "" {
		conditions.Set(clusterScope.IBMPowerVSCluster, metav1.Condition{
			Type:   infrav1.COSInstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrav1.COSInstanceDeletingReason,
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
	controllerutil.RemoveFinalizer(cluster, infrav1.IBMPowerVSClusterFinalizer)
	return ctrl.Result{}, nil
}

func (r *IBMPowerVSClusterReconciler) deleteIBMPowerVSImage(ctx context.Context, clusterScope *powervsscope.ClusterScope) (ctrl.Result, error) {
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
	ibmPowerVSImages infrav1.IBMPowerVSImageList
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
func (r *IBMPowerVSClusterReconciler) listDescendants(ctx context.Context, cluster *infrav1.IBMPowerVSCluster) (clusterDescendants, error) {
	var descendants clusterDescendants

	listOptions := []client.ListOption{
		client.InNamespace(cluster.Namespace),
		client.MatchingLabels(map[string]string{clusterv1.ClusterNameLabel: cluster.Name}),
	}

	if err := r.Client.List(ctx, &descendants.ibmPowerVSImages, listOptions...); err != nil {
		return descendants, fmt.Errorf("failed to list IBMPowerVSImages for cluster %s/%s: %w", cluster.Namespace, cluster.Name, err)
	}

	return descendants, nil
}

// filterOwnedDescendants returns an array of runtime.Objects containing only those descendants that have the cluster
// as an owner reference.
func (c *clusterDescendants) filterOwnedDescendants(cluster *infrav1.IBMPowerVSCluster) ([]client.Object, error) {
	var ownedDescendants []client.Object
	eachFunc := func(o runtime.Object) error {
		obj := o.(client.Object)
		acc, err := meta.Accessor(obj)
		if err != nil {
			return nil //nolint:nilerr // We don't want to exit the EachListItem loop, just continue
		}

		if util.IsOwnedByObject(acc, cluster, cluster.GroupVersionKind().GroupKind()) {
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

// buildConditions generates both a metav1.Condition and a legacy v1beta1 Condition simultaneously.
// reason is used for the modern condition; legacyReason is used for the legacy condition.
func (r *IBMPowerVSClusterReconciler) buildConditions(condType string, legacyType clusterv1.ConditionType, status metav1.ConditionStatus, reason, legacyReason, msg string) (metav1.Condition, *clusterv1.Condition) {
	cond := metav1.Condition{
		Type:    condType,
		Status:  status,
		Reason:  reason,
		Message: msg,
	}
	legacy := &clusterv1.Condition{
		Type:    legacyType,
		Status:  corev1.ConditionStatus(status),
		Reason:  legacyReason,
		Message: msg,
	}
	if status == metav1.ConditionFalse {
		legacy.Severity = clusterv1.ConditionSeverityError
	}
	return cond, legacy
}

// patchIBMPowerVSCluster updates the IBMPowerVSCluster and its status on the API server.
func patchIBMPowerVSCluster(ctx context.Context, patchHelper *patch.Helper, ibmPowerVSCluster *infrav1.IBMPowerVSCluster) error {
	// We don't need to set VPC/LoadBalancer conditions for an IBMPowerVSCluster
	// unless it is explicitly using the LoadBalancer topology.
	if ibmPowerVSCluster.Spec.Topology != infrav1.PowerVSLoadBalancerTopology {
		if err := patchHelper.Patch(ctx, ibmPowerVSCluster); err != nil {
			return fmt.Errorf("error patching IBMPowerVSCluster: %w", err)
		}
		return nil
	}

	if err := conditions.SetSummaryCondition(ibmPowerVSCluster, ibmPowerVSCluster, infrav1.IBMPowerVSClusterReadyCondition,
		conditions.ForConditionTypes{
			infrav1.WorkspaceReadyCondition,
			infrav1.NetworkReadyCondition,
			infrav1.VPCReadyCondition,
			infrav1.VPCSubnetReadyCondition,
			infrav1.VPCSecurityGroupReadyCondition,
			infrav1.VPCLoadBalancerReadyCondition,
			infrav1.TransitGatewayReadyCondition,
			infrav1.COSInstanceReadyCondition,
		},
		conditions.IgnoreTypesIfMissing{
			infrav1.COSInstanceReadyCondition,
		},
		// Using a custom merge strategy to override reasons applied during merge.
		conditions.CustomMergeStrategy{
			MergeStrategy: conditions.DefaultMergeStrategy(
				// Use custom reasons.
				conditions.ComputeReasonFunc(conditions.GetDefaultComputeMergeReasonFunc(
					infrav1.IBMPowerVSClusterReadyReason,
					infrav1.IBMPowerVSClusterReadyUnknownReason,
					infrav1.IBMPowerVSClusterReadyReason,
				)),
			),
		},
	); err != nil {
		return fmt.Errorf("failed to set %s condition: %w", infrav1.IBMPowerVSClusterReadyCondition, err)
	}

	return patchHelper.Patch(ctx, ibmPowerVSCluster,
		patch.WithOwnedConditions{Conditions: []string{
			clusterv1.PausedCondition,
			infrav1.IBMPowerVSClusterReadyCondition,
			infrav1.WorkspaceReadyCondition,
			infrav1.NetworkReadyCondition,
			infrav1.VPCReadyCondition,
			infrav1.VPCSubnetReadyCondition,
			infrav1.VPCSecurityGroupReadyCondition,
			infrav1.TransitGatewayReadyCondition,
			infrav1.COSInstanceReadyCondition,
		}}, patch.Clusterv1ConditionsFieldPath{statusField, deprecatedStatus, v1beta2Version, deprecatedConditionsField},
	)
}

// SetupWithManager creates a new IBMPowerVSCluster controller for a manager.
func (r *IBMPowerVSClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	if r.ClientBuilder == nil {
		r.ClientBuilder = powervsscope.ProdClientBuilder{}
	}
	predicateLog := ctrl.LoggerFrom(ctx).WithValues("controller", "ibmpowervscluster")
	err := ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.IBMPowerVSCluster{}).
		WithEventFilter(predicates.ResourceHasFilterLabel(r.Scheme, predicateLog, r.WatchFilterValue)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(r.Scheme, predicateLog)).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(ctx, infrav1.GroupVersion.WithKind(ibmPowerVSClusterKind), mgr.GetClient(), &infrav1.IBMPowerVSCluster{})),
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
