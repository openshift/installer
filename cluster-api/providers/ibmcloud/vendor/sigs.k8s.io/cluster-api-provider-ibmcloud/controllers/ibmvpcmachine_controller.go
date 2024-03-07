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

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
)

// IBMVPCMachineReconciler reconciles a IBMVPCMachine object.
type IBMVPCMachineReconciler struct {
	client.Client
	Log             logr.Logger
	Recorder        record.EventRecorder
	ServiceEndpoint []endpoints.ServiceEndpoint
	Scheme          *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

// Reconcile implements controller runtime Reconciler interface and handles reconcileation logic for IBMVPCMachine.
func (r *IBMVPCMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := r.Log.WithValues("ibmvpcmachine", req.NamespacedName)

	// Fetch the IBMVPCMachine instance.

	ibmVpcMachine := &infrav1beta2.IBMVPCMachine{}
	err := r.Get(ctx, req.NamespacedName, ibmVpcMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	// Fetch the Machine.
	machine, err := util.GetOwnerMachine(ctx, r.Client, ibmVpcMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		log.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, ibmVpcMachine.ObjectMeta)
	if err != nil {
		log.Info("Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	ibmCluster := &infrav1beta2.IBMVPCCluster{}
	ibmVpcClusterName := client.ObjectKey{
		Namespace: ibmVpcMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, ibmVpcClusterName, ibmCluster); err != nil {
		log.Info("IBMVPCCluster is not available yet")
		return ctrl.Result{}, nil
	}

	// Create the machine scope.
	machineScope, err := scope.NewMachineScope(scope.MachineScopeParams{
		Client:          r.Client,
		Logger:          log,
		Cluster:         cluster,
		IBMVPCCluster:   ibmCluster,
		Machine:         machine,
		IBMVPCMachine:   ibmVpcMachine,
		ServiceEndpoint: r.ServiceEndpoint,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	// Always close the scope when exiting this function, so we can persist any IBMVPCMachine changes.
	defer func() {
		if machineScope != nil {
			if err := machineScope.Close(); err != nil && reterr == nil {
				reterr = err
			}
		}
	}()

	// Handle deleted machines.
	if !ibmVpcMachine.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(machineScope)
	}

	// Handle non-deleted machines.
	return r.reconcileNormal(machineScope)
}

// SetupWithManager creates a new IBMVPCMachine controller for a manager.
func (r *IBMVPCMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1beta2.IBMVPCMachine{}).
		Complete(r)
}

func (r *IBMVPCMachineReconciler) reconcileNormal(machineScope *scope.MachineScope) (ctrl.Result, error) {
	if controllerutil.AddFinalizer(machineScope.IBMVPCMachine, infrav1beta2.MachineFinalizer) {
		return ctrl.Result{}, nil
	}

	// Make sure bootstrap data is available and populated.
	if machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		machineScope.Info("Bootstrap data secret reference is not yet available")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	if machineScope.IBMVPCCluster.Status.Subnet.ID != nil {
		machineScope.IBMVPCMachine.Spec.PrimaryNetworkInterface = infrav1beta2.NetworkInterface{
			Subnet: *machineScope.IBMVPCCluster.Status.Subnet.ID,
		}
	}

	instance, err := r.getOrCreate(machineScope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile VSI for IBMVPCMachine %s/%s: %w", machineScope.IBMVPCMachine.Namespace, machineScope.IBMVPCMachine.Name, err)
	}

	if instance != nil {
		machineScope.IBMVPCMachine.Status.InstanceID = *instance.ID
		machineScope.IBMVPCMachine.Status.Addresses = []corev1.NodeAddress{
			{
				Type:    corev1.NodeInternalIP,
				Address: *instance.PrimaryNetworkInterface.PrimaryIP.Address,
			},
		}
		_, ok := machineScope.IBMVPCMachine.Labels[capiv1beta1.MachineControlPlaneNameLabel]
		if err = machineScope.SetProviderID(instance.ID); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to set provider id IBMVPCMachine %s/%s: %w", machineScope.IBMVPCMachine.Namespace, machineScope.IBMVPCMachine.Name, err)
		}
		if ok {
			if instance.PrimaryNetworkInterface.PrimaryIP.Address == nil || *instance.PrimaryNetworkInterface.PrimaryIP.Address == "0.0.0.0" {
				return ctrl.Result{}, fmt.Errorf("invalid primary ip address")
			}
			internalIP := instance.PrimaryNetworkInterface.PrimaryIP.Address
			port := int64(machineScope.APIServerPort())
			poolMember, err := machineScope.CreateVPCLoadBalancerPoolMember(internalIP, port)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to bind port %d to control plane %s/%s: %w", port, machineScope.IBMVPCMachine.Namespace, machineScope.IBMVPCMachine.Name, err)
			}
			if poolMember != nil && *poolMember.ProvisioningStatus != string(infrav1beta2.VPCLoadBalancerStateActive) {
				return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
			}
		}
		machineScope.IBMVPCMachine.Status.Ready = true
	}

	return ctrl.Result{}, nil
}

func (r *IBMVPCMachineReconciler) getOrCreate(scope *scope.MachineScope) (*vpcv1.Instance, error) {
	instance, err := scope.CreateMachine()
	return instance, err
}

func (r *IBMVPCMachineReconciler) reconcileDelete(scope *scope.MachineScope) (_ ctrl.Result, reterr error) {
	scope.Info("Handling deleted IBMVPCMachine")

	if _, ok := scope.IBMVPCMachine.Labels[capiv1beta1.MachineControlPlaneNameLabel]; ok {
		if err := scope.DeleteVPCLoadBalancerPoolMember(); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to delete loadBalancer pool member: %w", err)
		}
	}

	if err := scope.DeleteMachine(); err != nil {
		scope.Info("error deleting IBMVPCMachine")
		return ctrl.Result{}, fmt.Errorf("error deleting IBMVPCMachine %s/%s: %w", scope.IBMVPCMachine.Namespace, scope.IBMVPCMachine.Spec.Name, err)
	}

	defer func() {
		if reterr == nil {
			// VSI is deleted so remove the finalizer.
			controllerutil.RemoveFinalizer(scope.IBMVPCMachine, infrav1beta2.MachineFinalizer)
		}
	}()

	return ctrl.Result{}, nil
}
