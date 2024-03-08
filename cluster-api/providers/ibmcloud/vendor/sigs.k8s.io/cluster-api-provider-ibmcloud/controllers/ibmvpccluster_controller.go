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
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
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
	log := r.Log.WithValues("ibmvpccluster", req.NamespacedName)

	// Fetch the IBMVPCCluster instance.
	ibmCluster := &infrav1beta2.IBMVPCCluster{}
	err := r.Get(ctx, req.NamespacedName, ibmCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, ibmCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:          r.Client,
		Logger:          log,
		Cluster:         cluster,
		IBMVPCCluster:   ibmCluster,
		ServiceEndpoint: r.ServiceEndpoint,
	})

	// Always close the scope when exiting this function so we can persist any GCPMachine changes.
	defer func() {
		if err := clusterScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Handle deleted clusters.
	if !ibmCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(clusterScope)
	}

	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}
	return r.reconcile(clusterScope)
}

func (r *IBMVPCClusterReconciler) reconcile(clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	if !controllerutil.ContainsFinalizer(clusterScope.IBMVPCCluster, infrav1beta2.ClusterFinalizer) {
		controllerutil.AddFinalizer(clusterScope.IBMVPCCluster, infrav1beta2.ClusterFinalizer)
		return ctrl.Result{}, nil
	}

	vpc, err := clusterScope.CreateVPC()
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to reconcile VPC for IBMVPCCluster %s/%s", clusterScope.IBMVPCCluster.Namespace, clusterScope.IBMVPCCluster.Name)
	}
	if vpc != nil {
		clusterScope.IBMVPCCluster.Status.VPC = infrav1beta2.VPC{
			ID:   *vpc.ID,
			Name: *vpc.Name,
		}
	}

	if clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint.Host == "" && clusterScope.IBMVPCCluster.Spec.ControlPlaneLoadBalancer == nil {
		fip, err := clusterScope.ReserveFIP()
		if err != nil {
			return ctrl.Result{}, errors.Wrapf(err, "failed to reconcile Control Plane Endpoint for IBMVPCCluster %s/%s", clusterScope.IBMVPCCluster.Namespace, clusterScope.IBMVPCCluster.Name)
		}

		if fip != nil {
			clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint = capiv1beta1.APIEndpoint{
				Host: *fip.Address,
				Port: clusterScope.APIServerPort(),
			}

			clusterScope.IBMVPCCluster.Status.VPCEndpoint = infrav1beta2.VPCEndpoint{
				Address: fip.Address,
				FIPID:   fip.ID,
			}
		}
		clusterScope.SetReady()
	}

	if clusterScope.IBMVPCCluster.Status.Subnet.ID == nil {
		subnet, err := clusterScope.CreateSubnet()
		if err != nil {
			return ctrl.Result{}, errors.Wrapf(err, "failed to reconcile Subnet for IBMVPCCluster %s/%s", clusterScope.IBMVPCCluster.Namespace, clusterScope.IBMVPCCluster.Name)
		}
		if subnet != nil {
			clusterScope.IBMVPCCluster.Status.Subnet = infrav1beta2.Subnet{
				Ipv4CidrBlock: subnet.Ipv4CIDRBlock,
				Name:          subnet.Name,
				ID:            subnet.ID,
				Zone:          subnet.Zone.Name,
			}
		}
	}

	if clusterScope.IBMVPCCluster.Spec.ControlPlaneLoadBalancer != nil {
		loadBalancer, err := r.getOrCreate(clusterScope)
		if err != nil {
			return ctrl.Result{}, errors.Wrapf(err, "failed to reconcile Control Plane LoadBalancer for IBMVPCCluster %s/%s", clusterScope.IBMVPCCluster.Namespace, clusterScope.IBMVPCCluster.Name)
		}

		if loadBalancer != nil {
			clusterScope.SetLoadBalancerID(loadBalancer.ID)
			clusterScope.Logger.V(3).Info("LoadBalancerID - " + clusterScope.GetLoadBalancerID())
			clusterScope.SetLoadBalancerState(*loadBalancer.ProvisioningStatus)
			clusterScope.Logger.V(3).Info("LoadBalancerState - " + string(clusterScope.GetLoadBalancerState()))

			clusterScope.IBMVPCCluster.Spec.ControlPlaneEndpoint = capiv1beta1.APIEndpoint{
				Host: *loadBalancer.Hostname,
				Port: clusterScope.APIServerPort(),
			}

			clusterScope.IBMVPCCluster.Status.VPCEndpoint = infrav1beta2.VPCEndpoint{
				Address: loadBalancer.Hostname,
				LBID:    loadBalancer.ID,
			}

			switch clusterScope.GetLoadBalancerState() {
			case infrav1beta2.VPCLoadBalancerStateCreatePending:
				clusterScope.Logger.V(3).Info("LoadBalancer is in create state")
				clusterScope.SetNotReady()
				conditions.MarkFalse(clusterScope.IBMVPCCluster, infrav1beta2.LoadBalancerReadyCondition, string(infrav1beta2.VPCLoadBalancerStateCreatePending), capiv1beta1.ConditionSeverityInfo, *loadBalancer.OperatingStatus)
			case infrav1beta2.VPCLoadBalancerStateActive:
				clusterScope.Logger.V(3).Info("LoadBalancer is in active state")
				clusterScope.SetReady()
				conditions.MarkTrue(clusterScope.IBMVPCCluster, infrav1beta2.LoadBalancerReadyCondition)
			default:
				clusterScope.Logger.V(3).Info("LoadBalancer state is undefined", "state", clusterScope.GetLoadBalancerState(), "loadbalancer-id", clusterScope.GetLoadBalancerID())
				clusterScope.SetNotReady()
				conditions.MarkUnknown(clusterScope.IBMVPCCluster, infrav1beta2.LoadBalancerReadyCondition, *loadBalancer.ProvisioningStatus, "")
			}
		}
	}

	// Requeue after 1 minute if cluster is not ready to update status of the cluster properly.
	if !clusterScope.IsReady() {
		clusterScope.Info("Cluster is not yet ready")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}
	return ctrl.Result{}, nil
}

func (r *IBMVPCClusterReconciler) reconcileDelete(clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	// check if still have existing VSIs.
	listVSIOpts := &vpcv1.ListInstancesOptions{
		VPCID: &clusterScope.IBMVPCCluster.Status.VPC.ID,
	}
	vsis, _, err := clusterScope.IBMVPCClient.ListInstances(listVSIOpts)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "Error when listing VSIs when tried to delete subnet ")
	}
	// skip deleting other resources if still have vsis running.
	if *vsis.TotalCount != int64(0) {
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	if clusterScope.IBMVPCCluster.Spec.ControlPlaneLoadBalancer != nil {
		deleted, err := clusterScope.DeleteLoadBalancer()
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "failed to delete loadBalancer")
		}
		// skip deleting other resources if still have loadBalancers running.
		if deleted {
			return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
		}
	}

	if err := clusterScope.DeleteSubnet(); err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to delete subnet")
	}

	if clusterScope.IBMVPCCluster.Spec.ControlPlaneLoadBalancer == nil {
		if err := clusterScope.DeleteFloatingIP(); err != nil {
			return ctrl.Result{}, errors.Wrap(err, "failed to delete floatingIP")
		}
	}

	if err := clusterScope.DeleteVPC(); err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to delete VPC")
	}
	controllerutil.RemoveFinalizer(clusterScope.IBMVPCCluster, infrav1beta2.ClusterFinalizer)
	return ctrl.Result{}, nil
}

func (r *IBMVPCClusterReconciler) getOrCreate(clusterScope *scope.ClusterScope) (*vpcv1.LoadBalancer, error) {
	loadBalancer, err := clusterScope.CreateLoadBalancer()
	return loadBalancer, err
}

// SetupWithManager creates a new IBMVPCCluster controller for a manager.
func (r *IBMVPCClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1beta2.IBMVPCCluster{}).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(ctrl.LoggerFrom(context.TODO()))).
		Complete(r)
}
