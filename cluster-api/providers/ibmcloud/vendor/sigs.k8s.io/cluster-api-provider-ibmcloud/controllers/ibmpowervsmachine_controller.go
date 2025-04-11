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
	"github.com/pkg/errors"

	"github.com/IBM-Cloud/power-go-client/power/models"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
	capibmrecord "sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/record"
)

// IBMPowerVSMachineReconciler reconciles a IBMPowerVSMachine object.
type IBMPowerVSMachineReconciler struct {
	client.Client
	Log             logr.Logger
	Recorder        record.EventRecorder
	ServiceEndpoint []endpoints.ServiceEndpoint
	Scheme          *runtime.Scheme
}

// dhcpCacheStore is a cache store to hold the Power VS VM DHCP IP.
var dhcpCacheStore cache.Store

func init() {
	dhcpCacheStore = powervs.InitialiseDHCPCacheStore()
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachines/status,verbs=get;update;patch

// Reconcile implements controller runtime Reconciler interface and handles reconcileation logic for IBMPowerVSMachine.
func (r *IBMPowerVSMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	ibmPowerVSMachine := &infrav1beta2.IBMPowerVSMachine{}
	err := r.Get(ctx, req.NamespacedName, ibmPowerVSMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch the Machine.
	machine, err := util.GetOwnerMachine(ctx, r.Client, ibmPowerVSMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		log.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}
	log = log.WithValues("ibmPowerVSMachine", machine.Name)

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, ibmPowerVSMachine.ObjectMeta)
	if err != nil {
		log.Info("Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	ibmCluster := &infrav1beta2.IBMPowerVSCluster{}
	ibmPowerVSClusterName := client.ObjectKey{
		Namespace: ibmPowerVSMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, ibmPowerVSClusterName, ibmCluster); err != nil {
		log.Info("IBMPowerVSCluster is not available yet")
		return ctrl.Result{}, nil
	}

	var ibmPowerVSImage *infrav1beta2.IBMPowerVSImage
	if ibmPowerVSMachine.Spec.ImageRef != nil {
		ibmPowerVSImage = &infrav1beta2.IBMPowerVSImage{}
		ibmPowerVSImageName := client.ObjectKey{
			Namespace: ibmPowerVSMachine.Namespace,
			Name:      ibmPowerVSMachine.Spec.ImageRef.Name,
		}
		if err := r.Client.Get(ctx, ibmPowerVSImageName, ibmPowerVSImage); err != nil {
			log.Info("IBMPowerVSImage is not available yet", "IBMPowerVSImage", klog.KObj(ibmPowerVSImage))
			return ctrl.Result{}, nil
		}
	}

	// Create the machine scope.
	machineScope, err := scope.NewPowerVSMachineScope(scope.PowerVSMachineScopeParams{
		Client:            r.Client,
		Logger:            log,
		Cluster:           cluster,
		IBMPowerVSCluster: ibmCluster,
		Machine:           machine,
		IBMPowerVSMachine: ibmPowerVSMachine,
		IBMPowerVSImage:   ibmPowerVSImage,
		ServiceEndpoint:   r.ServiceEndpoint,
		DHCPIPCacheStore:  dhcpCacheStore,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	// Always close the scope when exiting this function so we can persist any IBMPowerVSMachine changes.
	defer func() {
		if machineScope != nil {
			if err := machineScope.Close(); err != nil && reterr == nil {
				reterr = err
			}
		}
	}()

	// Handle deleted machines.
	if !ibmPowerVSMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(machineScope)
	}

	// Handle non-deleted machines.
	return r.reconcileNormal(machineScope)
}

func (r *IBMPowerVSMachineReconciler) reconcileDelete(scope *scope.PowerVSMachineScope) (_ ctrl.Result, reterr error) {
	scope.Info("Handling deleted IBMPowerVSMachine")

	defer func() {
		if reterr == nil {
			// VSI is deleted so remove the finalizer.
			controllerutil.RemoveFinalizer(scope.IBMPowerVSMachine, infrav1beta2.IBMPowerVSMachineFinalizer)
		}
	}()

	if scope.IBMPowerVSMachine.Status.InstanceID == "" {
		scope.Info("InstanceID is not yet set, hence not invoking the PowerVS API to delete the instance")
		return ctrl.Result{}, nil
	}
	if err := scope.DeleteMachine(); err != nil {
		scope.Info("error deleting IBMPowerVSMachine")
		return ctrl.Result{}, fmt.Errorf("error deleting IBMPowerVSMachine %v: %w", klog.KObj(scope.IBMPowerVSMachine), err)
	}
	if err := scope.DeleteMachineIgnition(); err != nil {
		scope.Info("error deleting IBMPowerVSMachine ignition")
		return ctrl.Result{}, fmt.Errorf("error deleting IBMPowerVSMachine ignition %v: %w", klog.KObj(scope.IBMPowerVSMachine), err)
	}
	// Remove the cached VM IP
	err := scope.DHCPIPCacheStore.Delete(powervs.VMip{Name: scope.IBMPowerVSMachine.Name})
	if err != nil {
		scope.Error(err, "failed to delete the VM entry from DHCP cache store", "VM", scope.IBMPowerVSMachine.Name)
	}
	return ctrl.Result{}, nil
}

func (r *IBMPowerVSMachineReconciler) getOrCreate(scope *scope.PowerVSMachineScope) (*models.PVMInstanceReference, error) {
	instance, err := scope.CreateMachine()
	return instance, err
}

// handleLoadBalancerPoolMemberConfiguration handles loadbalancer pool member creation flow.
func (r *IBMPowerVSMachineReconciler) handleLoadBalancerPoolMemberConfiguration(machineScope *scope.PowerVSMachineScope) (ctrl.Result, error) {
	poolMember, err := machineScope.CreateVPCLoadBalancerPoolMember()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create loadbalancer pool member %s: %w", machineScope.IBMPowerVSMachine.Name, err)
	}
	if poolMember != nil && *poolMember.ProvisioningStatus != string(infrav1beta2.VPCLoadBalancerStateActive) {
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	return ctrl.Result{}, nil
}

func (r *IBMPowerVSMachineReconciler) reconcileNormal(machineScope *scope.PowerVSMachineScope) (ctrl.Result, error) {
	machineScope.Info("Reconciling IBMPowerVSMachine")

	if !machineScope.Cluster.Status.InfrastructureReady {
		machineScope.Info("Cluster infrastructure is not ready yet")
		conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1beta2.InstanceReadyCondition, infrav1beta2.WaitingForClusterInfrastructureReason, capiv1beta1.ConditionSeverityInfo, "")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	if machineScope.IBMPowerVSImage != nil {
		if !machineScope.IBMPowerVSImage.Status.Ready {
			machineScope.Info("IBMPowerVSImage is not ready yet")
			conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1beta2.InstanceReadyCondition, infrav1beta2.WaitingForIBMPowerVSImageReason, capiv1beta1.ConditionSeverityInfo, "")
			return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
		}
	}

	// Make sure bootstrap data is available and populated.
	if machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		machineScope.Info("Bootstrap data secret reference is not yet available")
		conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1beta2.InstanceReadyCondition, infrav1beta2.WaitingForBootstrapDataReason, capiv1beta1.ConditionSeverityInfo, "")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	if controllerutil.AddFinalizer(machineScope.IBMPowerVSMachine, infrav1beta2.IBMPowerVSMachineFinalizer) {
		return ctrl.Result{}, nil
	}

	ins, err := r.getOrCreate(machineScope)
	if err != nil {
		machineScope.Error(err, "Unable to create instance")
		conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1beta2.InstanceReadyCondition, infrav1beta2.InstanceProvisionFailedReason, capiv1beta1.ConditionSeverityError, "%s", err.Error())
		return ctrl.Result{}, fmt.Errorf("failed to reconcile VSI for IBMPowerVSMachine %s/%s: %w", machineScope.IBMPowerVSMachine.Namespace, machineScope.IBMPowerVSMachine.Name, err)
	}

	if ins != nil {
		instance, err := machineScope.IBMPowerVSClient.GetInstance(*ins.PvmInstanceID)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err := machineScope.SetProviderID(*ins.PvmInstanceID); err != nil {
			return ctrl.Result{}, errors.Wrapf(err, "failed to set provider id")
		}
		machineScope.SetInstanceID(instance.PvmInstanceID)
		machineScope.SetAddresses(instance)
		machineScope.SetHealth(instance.Health)
		machineScope.SetInstanceState(instance.Status)
		switch machineScope.GetInstanceState() {
		case infrav1beta2.PowerVSInstanceStateBUILD:
			machineScope.SetNotReady()
			conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1beta2.InstanceReadyCondition, infrav1beta2.InstanceNotReadyReason, capiv1beta1.ConditionSeverityWarning, "")
		case infrav1beta2.PowerVSInstanceStateSHUTOFF:
			machineScope.SetNotReady()
			conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1beta2.InstanceReadyCondition, infrav1beta2.InstanceStoppedReason, capiv1beta1.ConditionSeverityError, "")
			return ctrl.Result{}, nil
		case infrav1beta2.PowerVSInstanceStateACTIVE:
			machineScope.SetReady()
			conditions.MarkTrue(machineScope.IBMPowerVSMachine, infrav1beta2.InstanceReadyCondition)
		case infrav1beta2.PowerVSInstanceStateERROR:
			msg := ""
			if instance.Fault != nil {
				msg = instance.Fault.Details
			}
			machineScope.SetNotReady()
			machineScope.SetFailureReason(infrav1beta2.UpdateMachineError)
			machineScope.SetFailureMessage(msg)
			conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1beta2.InstanceReadyCondition, infrav1beta2.InstanceErroredReason, capiv1beta1.ConditionSeverityError, "%s", msg)
			capibmrecord.Warnf(machineScope.IBMPowerVSMachine, "FailedBuildInstance", "Failed to build the instance - %s", msg)
			return ctrl.Result{}, nil
		default:
			machineScope.SetNotReady()
			machineScope.Info("PowerVS instance state is undefined", "state", *instance.Status, "instance-id", machineScope.GetInstanceID())
			conditions.MarkUnknown(machineScope.IBMPowerVSMachine, infrav1beta2.InstanceReadyCondition, "", "")
		}
	} else {
		machineScope.SetNotReady()
		conditions.MarkUnknown(machineScope.IBMPowerVSMachine, infrav1beta2.InstanceReadyCondition, infrav1beta2.InstanceStateUnknownReason, "")
	}
	// Requeue after 2 minute if machine is not ready to update status of the machine properly.
	if !machineScope.IsReady() {
		return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
	}

	if machineScope.IBMPowerVSCluster.Spec.VPC == nil || machineScope.IBMPowerVSCluster.Spec.VPC.Region == nil {
		machineScope.Info("Skipping configuring machine to loadbalancer as VPC is not set")
		return ctrl.Result{}, nil
	}

	// Register instance with load balancer
	machineScope.Info("updating loadbalancer for machine", "name", machineScope.IBMPowerVSMachine.Name)
	internalIP := machineScope.GetMachineInternalIP()
	if internalIP == "" {
		machineScope.Info("Unable to update the LoadBalancer, Machine internal IP not yet set", "machineName", machineScope.IBMPowerVSMachine.Name)
		return ctrl.Result{}, nil
	}

	if util.IsControlPlaneMachine(machineScope.Machine) {
		machineScope.Info("Configuring loadbalancer configuration for control plane machine", "machineName", machineScope.IBMPowerVSMachine.Name)
		return r.handleLoadBalancerPoolMemberConfiguration(machineScope)
	}
	machineScope.Info("skipping loadbalancer configuration for worker machine", "machineName", machineScope.IBMPowerVSMachine.Name)

	return ctrl.Result{}, nil
}

// IBMPowerVSClusterToIBMPowerVSMachines is a handler.ToRequestsFunc to be used to enqeue requests for reconciliation
// of IBMPowerVSMachines.
func (r *IBMPowerVSMachineReconciler) IBMPowerVSClusterToIBMPowerVSMachines(ctx context.Context) handler.MapFunc {
	log := ctrl.LoggerFrom(ctx)
	return func(mapCtx context.Context, o client.Object) []ctrl.Request {
		result := []ctrl.Request{}

		c, ok := o.(*infrav1beta2.IBMPowerVSCluster)
		if !ok {
			log.Error(errors.Errorf("expected a IBMPowerVSCluster but got a %T", o), "failed to get IBMPowerVSMachines for IBMPowerVSCluster")
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

		labels := map[string]string{capiv1beta1.ClusterNameLabel: cluster.Name}
		machineList := &capiv1beta1.MachineList{}
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

// SetupWithManager creates a new IBMVPCMachine controller for a manager.
func (r *IBMPowerVSMachineReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	controller, err := ctrl.NewControllerManagedBy(mgr).
		For(&infrav1beta2.IBMPowerVSMachine{}).
		WithEventFilter(predicates.ResourceNotPaused(r.Scheme, ctrl.LoggerFrom(ctx))).
		Watches(
			&capiv1beta1.Machine{},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrav1beta2.GroupVersion.WithKind("IBMPowerVSMachine"))),
		).
		Watches(
			&infrav1beta2.IBMPowerVSCluster{},
			handler.EnqueueRequestsFromMapFunc(r.IBMPowerVSClusterToIBMPowerVSMachines(ctx)),
		).
		Build(r)
	if err != nil {
		return errors.Wrap(err, "error creating controller")
	}

	clusterToObjectFunc, err := util.ClusterToTypedObjectsMapper(r.Client, &infrav1beta2.IBMPowerVSMachineList{}, mgr.GetScheme())
	if err != nil {
		return errors.Wrap(err, "failed to create mapper for Cluster to IBMPowerVSMachines")
	}
	// Add a watch on capiv1beta1.Cluster object for unpause & ready notifications.
	if err := controller.Watch(
		source.Kind[client.Object](mgr.GetCache(), &capiv1beta1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(clusterToObjectFunc),
			predicates.ClusterPausedTransitionsOrInfrastructureReady(r.Scheme, ctrl.LoggerFrom(ctx)),
		)); err != nil {
		return errors.Wrap(err, "failed adding a watch for ready clusters")
	}
	return nil
}
