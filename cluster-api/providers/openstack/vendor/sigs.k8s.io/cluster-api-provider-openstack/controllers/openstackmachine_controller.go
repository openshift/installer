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

package controllers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/compute"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/loadbalancer"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/names"
)

// OpenStackMachineReconciler reconciles a OpenStackMachine object.
type OpenStackMachineReconciler struct {
	Client           client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
	ScopeFactory     scope.Factory
	CaCertificates   []byte // PEM encoded ca certificates.
}

const (
	waitForClusterInfrastructureReadyDuration = 15 * time.Second
	waitForInstanceBecomeActiveToReconcile    = 60 * time.Second
	waitForBuildingInstanceToReconcile        = 10 * time.Second
)

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=ipaddressclaims;ipaddressclaims/status,verbs=get;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=ipaddresses;ipaddresses/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

func (r *OpenStackMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the OpenStackMachine instance.
	openStackMachine := &infrav1.OpenStackMachine{}
	err := r.Client.Get(ctx, req.NamespacedName, openStackMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	log = log.WithValues("openStackMachine", openStackMachine.Name)
	log.V(4).Info("Reconciling OpenStackMachine")

	// Fetch the Machine.
	machine, err := util.GetOwnerMachine(ctx, r.Client, openStackMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		log.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("machine", machine.Name)

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		log.Info("Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	if annotations.IsPaused(cluster, openStackMachine) {
		log.Info("OpenStackMachine or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	infraCluster, err := r.getInfraCluster(ctx, cluster, openStackMachine)
	if err != nil {
		return ctrl.Result{}, errors.New("error getting infra provider cluster")
	}
	if infraCluster == nil {
		log.Info("OpenStackCluster is not ready yet")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("openStackCluster", infraCluster.Name)

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(openStackMachine, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always patch the openStackMachine when exiting this function so we can persist any OpenStackMachine changes.
	defer func() {
		if err := patchMachine(ctx, patchHelper, openStackMachine, machine); err != nil {
			result = ctrl.Result{}
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	clientScope, err := r.ScopeFactory.NewClientScopeFromMachine(ctx, r.Client, openStackMachine, infraCluster, r.CaCertificates, log)
	if err != nil {
		return reconcile.Result{}, err
	}
	scope := scope.NewWithLogger(clientScope, log)

	clusterResourceName := fmt.Sprintf("%s-%s", cluster.Namespace, cluster.Name)

	// Handle deleted machines
	if !openStackMachine.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(scope, clusterResourceName, infraCluster, machine, openStackMachine)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, scope, clusterResourceName, infraCluster, machine, openStackMachine)
}

func resolveMachineResources(scope *scope.WithLogger, clusterResourceName string, openStackCluster *infrav1.OpenStackCluster, openStackMachine *infrav1.OpenStackMachine, machine *clusterv1.Machine) (bool, error) {
	resolved := openStackMachine.Status.Resolved
	if resolved == nil {
		resolved = &infrav1.ResolvedMachineSpec{}
		openStackMachine.Status.Resolved = resolved
	}
	// Resolve and store resources
	return compute.ResolveMachineSpec(scope,
		&openStackMachine.Spec, resolved,
		clusterResourceName, openStackMachine.Name,
		openStackCluster, getManagedSecurityGroup(openStackCluster, machine))
}

func adoptMachineResources(scope *scope.WithLogger, openStackMachine *infrav1.OpenStackMachine) error {
	resources := openStackMachine.Status.Resources
	if resources == nil {
		resources = &infrav1.MachineResources{}
		openStackMachine.Status.Resources = resources
	}

	// Adopt any existing resources
	return compute.AdoptMachineResources(scope, openStackMachine.Status.Resolved, resources)
}

func patchMachine(ctx context.Context, patchHelper *patch.Helper, openStackMachine *infrav1.OpenStackMachine, machine *clusterv1.Machine, options ...patch.Option) error {
	// Always update the readyCondition by summarizing the state of other conditions.
	applicableConditions := []clusterv1.ConditionType{
		infrav1.InstanceReadyCondition,
	}

	if util.IsControlPlaneMachine(machine) {
		applicableConditions = append(applicableConditions, infrav1.APIServerIngressReadyCondition)
	}

	conditions.SetSummary(openStackMachine,
		conditions.WithConditions(applicableConditions...),
	)

	// Patch the object, ignoring conflicts on the conditions owned by this controller.
	// Also, if requested, we are adding additional options like e.g. Patch ObservedGeneration when issuing the
	// patch at the end of the reconcile loop.
	options = append(options,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
			infrav1.InstanceReadyCondition,
			infrav1.APIServerIngressReadyCondition,
		}},
	)
	return patchHelper.Patch(ctx, openStackMachine, options...)
}

func (r *OpenStackMachineReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.OpenStackMachine{}).
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrav1.GroupVersion.WithKind("OpenStackMachine"))),
		).
		Watches(
			&infrav1.OpenStackCluster{},
			handler.EnqueueRequestsFromMapFunc(r.OpenStackClusterToOpenStackMachines(ctx)),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), r.WatchFilterValue)).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(r.requeueOpenStackMachinesForUnpausedCluster(ctx)),
			builder.WithPredicates(predicates.ClusterUnpausedAndInfrastructureReady(ctrl.LoggerFrom(ctx))),
		).
		Watches(
			&ipamv1.IPAddressClaim{},
			handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &infrav1.OpenStackMachine{}),
		).
		Complete(r)
}

func (r *OpenStackMachineReconciler) reconcileDelete(scope *scope.WithLogger, clusterResourceName string, openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine) (ctrl.Result, error) { //nolint:unparam
	scope.Logger().Info("Reconciling Machine delete")

	computeService, err := compute.NewService(scope)
	if err != nil {
		return ctrl.Result{}, err
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Nothing to do if the cluster is not ready because no machine resources were created.
	if !openStackCluster.Status.Ready || openStackCluster.Status.Network == nil {
		// The finalizer should not have been added yet in this case,
		// but the following handles the upgrade case.
		controllerutil.RemoveFinalizer(openStackMachine, infrav1.MachineFinalizer)
		return ctrl.Result{}, nil
	}

	// For machines created after v0.10, or any machine which has been
	// reconciled at least once by v0.10 or later, status.Resolved always
	// exists before any resources are created. We can therefore assume
	// that if it does not exist, no resources were created.
	//
	// There is an upgrade edge case where a machine may have been marked
	// deleted before upgrade but we are completing it after upgrade. For
	// this use case only we make a best effort to resolve resources before
	// continuing, but if we get an error we log it and continue anyway.
	// This has the potential to leak resources, but only in this specific
	// edge case. The alternative is to continue retrying until it succeeds,
	// but that risks never deleting a machine which cannot be resolved due
	// to a spec error.
	//
	// This code can and should be deleted in a future release when we are
	// sure that all machines have been reconciled at least by a v0.10 or
	// later controller.
	if _, err := resolveMachineResources(scope, clusterResourceName, openStackCluster, openStackMachine, machine); err != nil {
		// Return the error, but allow the resource to be removed anyway.
		controllerutil.RemoveFinalizer(openStackMachine, infrav1.MachineFinalizer)
		return ctrl.Result{}, err
	}

	// Check for any orphaned resources
	// N.B. Unlike resolveMachineResources, we must always look for orphaned resources in the delete path.
	if err := adoptMachineResources(scope, openStackMachine); err != nil {
		return ctrl.Result{}, fmt.Errorf("adopting machine resources: %w", err)
	}

	instanceStatus, err := getInstanceStatus(openStackMachine, computeService)
	if err != nil {
		return ctrl.Result{}, err
	}

	if util.IsControlPlaneMachine(machine) {
		if err := removeAPIServerEndpoint(scope, openStackCluster, openStackMachine, instanceStatus, clusterResourceName); err != nil {
			return ctrl.Result{}, err
		}
	}

	// If no instance was created we currently need to check for orphaned
	// volumes.
	if instanceStatus == nil {
		if err := computeService.DeleteVolumes(openStackMachine.Name, openStackMachine.Spec.RootVolume, openStackMachine.Spec.AdditionalBlockDevices); err != nil {
			return ctrl.Result{}, fmt.Errorf("delete volumes: %w", err)
		}
	} else {
		if err := computeService.DeleteInstance(openStackMachine, instanceStatus); err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceDeleteFailedReason, clusterv1.ConditionSeverityError, "Deleting instance failed: %v", err)
			return ctrl.Result{}, fmt.Errorf("delete instance: %w", err)
		}
	}

	trunkSupported, err := networkingService.IsTrunkExtSupported()
	if err != nil {
		return ctrl.Result{}, err
	}

	if openStackMachine.Status.Resources != nil {
		portsStatus := openStackMachine.Status.Resources.Ports
		for _, port := range portsStatus {
			if err := networkingService.DeleteInstanceTrunkAndPort(openStackMachine, port, trunkSupported); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to delete port %q: %w", port.ID, err)
			}
		}
	}

	if err := r.reconcileDeleteFloatingAddressFromPool(scope, openStackMachine); err != nil {
		return ctrl.Result{}, err
	}

	controllerutil.RemoveFinalizer(openStackMachine, infrav1.MachineFinalizer)
	scope.Logger().Info("Reconciled Machine delete successfully")
	return ctrl.Result{}, nil
}

func getInstanceStatus(openStackMachine *infrav1.OpenStackMachine, computeService *compute.Service) (*compute.InstanceStatus, error) {
	if openStackMachine.Status.InstanceID != nil {
		return computeService.GetInstanceStatus(*openStackMachine.Status.InstanceID)
	}
	return computeService.GetInstanceStatusByName(openStackMachine, openStackMachine.Name)
}

func removeAPIServerEndpoint(scope *scope.WithLogger, openStackCluster *infrav1.OpenStackCluster, openStackMachine *infrav1.OpenStackMachine, instanceStatus *compute.InstanceStatus, clusterResourceName string) error {
	if openStackCluster.Spec.APIServerLoadBalancer.IsEnabled() {
		loadBalancerService, err := loadbalancer.NewService(scope)
		if err != nil {
			return err
		}

		err = loadBalancerService.DeleteLoadBalancerMember(openStackCluster, openStackMachine, clusterResourceName)
		if err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.APIServerIngressReadyCondition, infrav1.LoadBalancerMemberErrorReason, clusterv1.ConditionSeverityWarning, "Machine could not be removed from load balancer: %v", err)
			return err
		}
		return nil
	}

	// XXX(mdbooth): This looks wrong to me. Surely we should only ever
	// disassociate the floating IP here. I would expect the API server
	// floating IP to be created and deleted with the cluster. And if the
	// delete behaviour is correct, we leak it if the instance was
	// previously deleted.
	if openStackCluster.Spec.APIServerFloatingIP == nil && instanceStatus != nil {
		instanceNS, err := instanceStatus.NetworkStatus()
		if err != nil {
			openStackMachine.SetFailure(
				capierrors.UpdateMachineError,
				fmt.Errorf("get network status for OpenStack instance %s with ID %s: %v", instanceStatus.Name(), instanceStatus.ID(), err),
			)
			return nil
		}

		networkingService, err := networking.NewService(scope)
		if err != nil {
			return err
		}

		addresses := instanceNS.Addresses()
		for _, address := range addresses {
			if address.Type == corev1.NodeExternalIP {
				if err = networkingService.DeleteFloatingIP(openStackMachine, address.Address); err != nil {
					conditions.MarkFalse(openStackMachine, infrav1.APIServerIngressReadyCondition, infrav1.FloatingIPErrorReason, clusterv1.ConditionSeverityError, "Deleting floating IP failed: %v", err)
					return fmt.Errorf("delete floating IP %q: %w", address.Address, err)
				}
			}
		}
	}

	return nil
}

// GetPortIDs returns a list of port IDs from a list of PortStatus.
func GetPortIDs(ports []infrav1.PortStatus) []string {
	portIDs := make([]string, len(ports))
	for i, port := range ports {
		portIDs[i] = port.ID
	}
	return portIDs
}

// reconcileFloatingAddressFromPool reconciles the floating IP address from the pool.
// It returns the IPAddressClaim and a boolean indicating if the IPAddressClaim is ready.
func (r *OpenStackMachineReconciler) reconcileFloatingAddressFromPool(ctx context.Context, scope *scope.WithLogger, openStackMachine *infrav1.OpenStackMachine, openStackCluster *infrav1.OpenStackCluster) (*ipamv1.IPAddressClaim, bool, error) {
	if openStackMachine.Spec.FloatingIPPoolRef == nil {
		return nil, false, nil
	}
	var claim *ipamv1.IPAddressClaim
	claim, err := r.getOrCreateIPAddressClaimForFloatingAddress(ctx, scope, openStackMachine, openStackCluster)
	if err != nil {
		conditions.MarkFalse(openStackMachine, infrav1.FloatingAddressFromPoolReadyCondition, infrav1.FloatingAddressFromPoolErrorReason, clusterv1.ConditionSeverityInfo, "Failed to reconcile floating IP claims: %v", err)
		return nil, true, err
	}
	if claim.Status.AddressRef.Name == "" {
		r.Recorder.Eventf(openStackMachine, corev1.EventTypeNormal, "WaitingForIPAddressClaim", "Waiting for IPAddressClaim %s/%s to be allocated", claim.Namespace, claim.Name)
		return claim, true, nil
	}
	conditions.MarkTrue(openStackMachine, infrav1.FloatingAddressFromPoolReadyCondition)
	return claim, false, nil
}

// createIPAddressClaim creates IPAddressClaim for the FloatingAddressFromPool if it does not exist yet.
func (r *OpenStackMachineReconciler) getOrCreateIPAddressClaimForFloatingAddress(ctx context.Context, scope *scope.WithLogger, openStackMachine *infrav1.OpenStackMachine, openStackCluster *infrav1.OpenStackCluster) (*ipamv1.IPAddressClaim, error) {
	var err error

	poolRef := openStackMachine.Spec.FloatingIPPoolRef
	claimName := names.GetFloatingAddressClaimName(openStackMachine.Name)
	claim := &ipamv1.IPAddressClaim{}

	err = r.Client.Get(ctx, client.ObjectKey{Namespace: openStackMachine.Namespace, Name: claimName}, claim)
	if err == nil {
		return claim, nil
	} else if client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	claim = &ipamv1.IPAddressClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      claimName,
			Namespace: openStackMachine.Namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: openStackCluster.Labels[clusterv1.ClusterNameLabel],
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: openStackMachine.APIVersion,
					Kind:       openStackMachine.Kind,
					Name:       openStackMachine.Name,
					UID:        openStackMachine.UID,
				},
			},
			Finalizers: []string{infrav1.IPClaimMachineFinalizer},
		},
		Spec: ipamv1.IPAddressClaimSpec{
			PoolRef: *poolRef,
		},
	}

	if err := r.Client.Create(ctx, claim); err != nil {
		return nil, err
	}

	r.Recorder.Eventf(openStackMachine, corev1.EventTypeNormal, "CreatingIPAddressClaim", "Creating IPAddressClaim %s/%s", claim.Namespace, claim.Name)
	scope.Logger().Info("Created IPAddressClaim", "name", claim.Name)
	return claim, nil
}

func (r *OpenStackMachineReconciler) associateIPAddressFromIPAddressClaim(ctx context.Context, scope *scope.WithLogger, openStackMachine *infrav1.OpenStackMachine, instanceStatus *compute.InstanceStatus, instanceNS *compute.InstanceNetworkStatus, claim *ipamv1.IPAddressClaim) error {
	address := &ipamv1.IPAddress{}
	addressKey := client.ObjectKey{Namespace: openStackMachine.Namespace, Name: claim.Status.AddressRef.Name}

	if err := r.Client.Get(ctx, addressKey, address); err != nil {
		return err
	}

	instanceAddresses := instanceNS.Addresses()
	for _, instanceAddress := range instanceAddresses {
		if instanceAddress.Address == address.Spec.Address {
			conditions.MarkTrue(openStackMachine, infrav1.FloatingAddressFromPoolReadyCondition)
			return nil
		}
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return err
	}

	fip, err := networkingService.GetFloatingIP(address.Spec.Address)
	if err != nil {
		return err
	}

	if fip == nil {
		conditions.MarkFalse(openStackMachine, infrav1.FloatingAddressFromPoolReadyCondition, infrav1.FloatingAddressFromPoolErrorReason, clusterv1.ConditionSeverityError, "floating IP does not exist")
		return fmt.Errorf("floating IP %q does not exist", address.Spec.Address)
	}

	port, err := networkingService.GetPortForExternalNetwork(instanceStatus.ID(), fip.FloatingNetworkID)
	if err != nil {
		return fmt.Errorf("get port for floating IP %q: %w", fip.FloatingIP, err)
	}

	if port == nil {
		conditions.MarkFalse(openStackMachine, infrav1.FloatingAddressFromPoolReadyCondition, infrav1.FloatingAddressFromPoolErrorReason, clusterv1.ConditionSeverityError, "Can't find port for floating IP %q on external network %s", fip.FloatingIP, fip.FloatingNetworkID)
		return fmt.Errorf("port for floating IP %q on network %s does not exist", fip.FloatingIP, fip.FloatingNetworkID)
	}

	if err = networkingService.AssociateFloatingIP(openStackMachine, fip, port.ID); err != nil {
		return err
	}
	conditions.MarkTrue(openStackMachine, infrav1.FloatingAddressFromPoolReadyCondition)
	return nil
}

func (r *OpenStackMachineReconciler) reconcileDeleteFloatingAddressFromPool(scope *scope.WithLogger, openStackMachine *infrav1.OpenStackMachine) error {
	log := scope.Logger().WithValues("openStackMachine", openStackMachine.Name)
	log.Info("Reconciling Machine delete floating address from pool")
	if openStackMachine.Spec.FloatingIPPoolRef == nil {
		return nil
	}
	claimName := names.GetFloatingAddressClaimName(openStackMachine.Name)
	claim := &ipamv1.IPAddressClaim{}
	if err := r.Client.Get(context.Background(), client.ObjectKey{Namespace: openStackMachine.Namespace, Name: claimName}, claim); err != nil {
		return client.IgnoreNotFound(err)
	}

	controllerutil.RemoveFinalizer(claim, infrav1.IPClaimMachineFinalizer)
	return r.Client.Update(context.Background(), claim)
}

func (r *OpenStackMachineReconciler) reconcileNormal(ctx context.Context, scope *scope.WithLogger, clusterResourceName string, openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine) (_ ctrl.Result, reterr error) {
	var err error

	// If the OpenStackMachine is in an error state, return early.
	if openStackMachine.Status.FailureReason != nil || openStackMachine.Status.FailureMessage != nil {
		scope.Logger().Info("Not reconciling machine in failed state. See openStackMachine.status.failureReason, openStackMachine.status.failureMessage, or previously logged error for details")
		return ctrl.Result{}, nil
	}

	if !openStackCluster.Status.Ready {
		scope.Logger().Info("Cluster infrastructure is not ready yet, re-queuing machine")
		conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.WaitingForClusterInfrastructureReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{RequeueAfter: waitForClusterInfrastructureReadyDuration}, nil
	}

	// Make sure bootstrap data is available and populated.
	if machine.Spec.Bootstrap.DataSecretName == nil {
		scope.Logger().Info("Bootstrap data secret reference is not yet available")
		conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.WaitingForBootstrapDataReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	changed, err := resolveMachineResources(scope, clusterResourceName, openStackCluster, openStackMachine, machine)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Also add the finalizer when writing resolved resources so we can start creating resources on the next reconcile.
	if controllerutil.AddFinalizer(openStackMachine, infrav1.MachineFinalizer) {
		changed = true
	}

	// We requeue if we either added the finalizer or resolved machine
	// resources. This means that we never create any resources unless we
	// have observed that the finalizer and resolved machine resources were
	// successfully written in a previous transaction. This in turn means
	// that in the delete path we can be sure that if there are no resolved
	// resources then no resources were created.
	if changed {
		scope.Logger().V(6).Info("Machine resources updated, requeuing")
		return ctrl.Result{}, nil
	}

	// Check for orphaned resources previously created but not written to the status
	if err := adoptMachineResources(scope, openStackMachine); err != nil {
		return ctrl.Result{}, fmt.Errorf("adopting machine resources: %w", err)
	}

	scope.Logger().Info("Reconciling Machine")
	userData, err := r.getBootstrapData(ctx, machine, openStackMachine)
	if err != nil {
		return ctrl.Result{}, err
	}

	computeService, err := compute.NewService(scope)
	if err != nil {
		return ctrl.Result{}, err
	}

	floatingAddressClaim, waitingForFloatingAddress, err := r.reconcileFloatingAddressFromPool(ctx, scope, openStackMachine, openStackCluster)
	if err != nil || waitingForFloatingAddress {
		return ctrl.Result{}, err
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = getOrCreateMachinePorts(openStackMachine, networkingService)
	if err != nil {
		return ctrl.Result{}, err
	}
	portIDs := GetPortIDs(openStackMachine.Status.Resources.Ports)

	instanceStatus, err := r.getOrCreateInstance(scope.Logger(), openStackCluster, machine, openStackMachine, computeService, userData, portIDs)
	if err != nil || instanceStatus == nil {
		// Conditions set in getOrCreateInstance
		return ctrl.Result{}, err
	}

	state := instanceStatus.State()
	openStackMachine.Status.InstanceID = ptr.To(instanceStatus.ID())
	openStackMachine.Status.InstanceState = &state

	instanceNS, err := instanceStatus.NetworkStatus()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("get network status: %w", err)
	}

	addresses := instanceNS.Addresses()

	// For OpenShift and likely other systems, a node joins the cluster if the CSR generated by kubelet with the node name is approved.
	// The approval happens if the Machine InternalDNS matches the node name. The in-tree provider used the server name for the node name.
	// Let's add the server name as the InternalDNS to keep getting CSRs of nodes upgraded from in-tree provider approved.
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeInternalDNS,
		Address: instanceStatus.Name(),
	})
	openStackMachine.Status.Addresses = addresses

	if floatingAddressClaim != nil {
		if err := r.associateIPAddressFromIPAddressClaim(ctx, scope, openStackMachine, instanceStatus, instanceNS, floatingAddressClaim); err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.FloatingAddressFromPoolReadyCondition, infrav1.FloatingAddressFromPoolErrorReason, clusterv1.ConditionSeverityError, "Failed while associating ip from pool: %v", err)
			return ctrl.Result{}, err
		}
		conditions.MarkTrue(openStackMachine, infrav1.FloatingAddressFromPoolReadyCondition)
	}

	switch instanceStatus.State() {
	case infrav1.InstanceStateActive:
		scope.Logger().Info("Machine instance state is ACTIVE", "id", instanceStatus.ID())
		conditions.MarkTrue(openStackMachine, infrav1.InstanceReadyCondition)

		// Set properties required by CAPI machine controller
		openStackMachine.Spec.ProviderID = ptr.To(fmt.Sprintf("openstack:///%s", instanceStatus.ID()))
		openStackMachine.Status.Ready = true
	case infrav1.InstanceStateError:
		// If the machine has a NodeRef then it must have been working at some point,
		// so the error could be something temporary.
		// If not, it is more likely a configuration error so we set failure and never retry.
		scope.Logger().Info("Machine instance state is ERROR", "id", instanceStatus.ID())
		if machine.Status.NodeRef == nil {
			err = fmt.Errorf("instance state %q is unexpected", instanceStatus.State())
			openStackMachine.SetFailure(capierrors.UpdateMachineError, err)
		}
		conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceStateErrorReason, clusterv1.ConditionSeverityError, "")
		return ctrl.Result{}, nil
	case infrav1.InstanceStateDeleted:
		// we should avoid further actions for DELETED VM
		scope.Logger().Info("Machine instance state is DELETED, no actions")
		conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceDeletedReason, clusterv1.ConditionSeverityError, "")
		return ctrl.Result{}, nil
	case infrav1.InstanceStateBuild, infrav1.InstanceStateUndefined:
		scope.Logger().Info("Waiting for instance to become ACTIVE", "id", instanceStatus.ID(), "status", instanceStatus.State())
		return ctrl.Result{RequeueAfter: waitForBuildingInstanceToReconcile}, nil
	default:
		// The other state is normal (for example, migrating, shutoff) but we don't want to proceed until it's ACTIVE
		// due to potential conflict or unexpected actions
		scope.Logger().Info("Waiting for instance to become ACTIVE", "id", instanceStatus.ID(), "status", instanceStatus.State())
		conditions.MarkUnknown(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceNotReadyReason, "Instance state is not handled: %s", instanceStatus.State())

		return ctrl.Result{RequeueAfter: waitForInstanceBecomeActiveToReconcile}, nil
	}

	if !util.IsControlPlaneMachine(machine) {
		scope.Logger().Info("Not a Control plane machine, no floating ip reconcile needed, Reconciled Machine create successfully")
		return ctrl.Result{}, nil
	}

	err = r.reconcileAPIServerLoadBalancer(scope, openStackCluster, openStackMachine, instanceStatus, instanceNS, clusterResourceName)
	if err != nil {
		return ctrl.Result{}, err
	}

	conditions.MarkTrue(openStackMachine, infrav1.APIServerIngressReadyCondition)
	scope.Logger().Info("Reconciled Machine create successfully")
	return ctrl.Result{}, nil
}

func (r *OpenStackMachineReconciler) reconcileAPIServerLoadBalancer(scope *scope.WithLogger, openStackCluster *infrav1.OpenStackCluster, openStackMachine *infrav1.OpenStackMachine, instanceStatus *compute.InstanceStatus, instanceNS *compute.InstanceNetworkStatus, clusterResourceName string) error {
	scope.Logger().Info("Reconciling APIServerLoadBalancer")
	computeService, err := compute.NewService(scope)
	if err != nil {
		return err
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return err
	}

	if openStackCluster.Spec.APIServerLoadBalancer.IsEnabled() {
		err = r.reconcileLoadBalancerMember(scope, openStackCluster, openStackMachine, instanceNS, clusterResourceName)
		if err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.APIServerIngressReadyCondition, infrav1.LoadBalancerMemberErrorReason, clusterv1.ConditionSeverityError, "Reconciling load balancer member failed: %v", err)
			return fmt.Errorf("reconcile load balancer member: %w", err)
		}
	} else if !ptr.Deref(openStackCluster.Spec.DisableAPIServerFloatingIP, false) {
		var floatingIPAddress *string
		switch {
		case openStackCluster.Spec.ControlPlaneEndpoint != nil && openStackCluster.Spec.ControlPlaneEndpoint.IsValid():
			floatingIPAddress = &openStackCluster.Spec.ControlPlaneEndpoint.Host
		case openStackCluster.Spec.APIServerFloatingIP != nil:
			floatingIPAddress = openStackCluster.Spec.APIServerFloatingIP
		}
		fp, err := networkingService.GetOrCreateFloatingIP(openStackMachine, openStackCluster, clusterResourceName, floatingIPAddress)
		if err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.APIServerIngressReadyCondition, infrav1.FloatingIPErrorReason, clusterv1.ConditionSeverityError, "Floating IP cannot be obtained or created: %v", err)
			return fmt.Errorf("get or create floating IP %v: %w", floatingIPAddress, err)
		}
		port, err := computeService.GetManagementPort(openStackCluster, instanceStatus)
		if err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.APIServerIngressReadyCondition, infrav1.FloatingIPErrorReason, clusterv1.ConditionSeverityError, "Obtaining management port for control plane machine failed: %v", err)
			return fmt.Errorf("get management port for control plane machine: %w", err)
		}

		if fp.PortID != "" {
			scope.Logger().Info("Floating IP already associated to a port", "id", fp.ID, "fixedIP", fp.FixedIP, "portID", port.ID)
		} else {
			err = networkingService.AssociateFloatingIP(openStackMachine, fp, port.ID)
			if err != nil {
				conditions.MarkFalse(openStackMachine, infrav1.APIServerIngressReadyCondition, infrav1.FloatingIPErrorReason, clusterv1.ConditionSeverityError, "Associating floating IP failed: %v", err)
				return fmt.Errorf("associate floating IP %q with port %q: %w", fp.FloatingIP, port.ID, err)
			}
		}
	}
	conditions.MarkTrue(openStackMachine, infrav1.APIServerIngressReadyCondition)
	return nil
}

func getOrCreateMachinePorts(openStackMachine *infrav1.OpenStackMachine, networkingService *networking.Service) error {
	resolved := openStackMachine.Status.Resolved
	if resolved == nil {
		return errors.New("machine resolved is nil")
	}
	resources := openStackMachine.Status.Resources
	if resources == nil {
		return errors.New("machine resources is nil")
	}
	desiredPorts := resolved.Ports

	if len(desiredPorts) == len(resources.Ports) {
		return nil
	}

	if err := networkingService.CreatePorts(openStackMachine, desiredPorts, resources); err != nil {
		return fmt.Errorf("creating ports: %w", err)
	}

	return nil
}

func (r *OpenStackMachineReconciler) getOrCreateInstance(logger logr.Logger, openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine, computeService *compute.Service, userData string, portIDs []string) (*compute.InstanceStatus, error) {
	var instanceStatus *compute.InstanceStatus
	var err error
	if openStackMachine.Status.InstanceID != nil {
		instanceStatus, err = computeService.GetInstanceStatus(*openStackMachine.Status.InstanceID)
		if err != nil {
			logger.Info("Unable to get OpenStack instance", "name", openStackMachine.Name)
			conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.OpenStackErrorReason, clusterv1.ConditionSeverityError, err.Error())
			return nil, err
		}
	}
	if instanceStatus == nil {
		// Check if there is an existing instance with machine name, in case where instance ID would not have been stored in machine status
		instanceStatus, err = computeService.GetInstanceStatusByName(openStackMachine, openStackMachine.Name)
		if err != nil {
			logger.Info("Unable to get OpenStack instance by name", "name", openStackMachine.Name)
			conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceCreateFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return nil, err
		}
		if instanceStatus != nil {
			return instanceStatus, nil
		}
		if openStackMachine.Status.InstanceID != nil {
			logger.Info("Not reconciling machine in failed state. The previously existing OpenStack instance is no longer available")
			conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceNotFoundReason, clusterv1.ConditionSeverityError, "virtual machine no longer exists")
			openStackMachine.SetFailure(capierrors.UpdateMachineError, errors.New("virtual machine no longer exists"))
			return nil, nil
		}

		instanceSpec, err := machineToInstanceSpec(openStackCluster, machine, openStackMachine, userData)
		if err != nil {
			return nil, err
		}

		logger.Info("Machine does not exist, creating Machine", "name", openStackMachine.Name)
		instanceStatus, err = computeService.CreateInstance(openStackMachine, instanceSpec, portIDs)
		if err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceCreateFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return nil, fmt.Errorf("create OpenStack instance: %w", err)
		}
	}
	return instanceStatus, nil
}

func machineToInstanceSpec(openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine, userData string) (*compute.InstanceSpec, error) {
	resolved := openStackMachine.Status.Resolved
	if resolved == nil {
		return nil, errors.New("machine resolved is nil")
	}

	serverMetadata := make(map[string]string, len(openStackMachine.Spec.ServerMetadata))
	for i := range openStackMachine.Spec.ServerMetadata {
		key := openStackMachine.Spec.ServerMetadata[i].Key
		value := openStackMachine.Spec.ServerMetadata[i].Value
		serverMetadata[key] = value
	}

	instanceSpec := compute.InstanceSpec{
		Name:                   openStackMachine.Name,
		ImageID:                resolved.ImageID,
		Flavor:                 openStackMachine.Spec.Flavor,
		SSHKeyName:             openStackMachine.Spec.SSHKeyName,
		UserData:               userData,
		Metadata:               serverMetadata,
		ConfigDrive:            openStackMachine.Spec.ConfigDrive != nil && *openStackMachine.Spec.ConfigDrive,
		RootVolume:             openStackMachine.Spec.RootVolume,
		AdditionalBlockDevices: openStackMachine.Spec.AdditionalBlockDevices,
		ServerGroupID:          resolved.ServerGroupID,
		Trunk:                  openStackMachine.Spec.Trunk,
	}

	// Add the failure domain only if specified
	if machine.Spec.FailureDomain != nil {
		instanceSpec.FailureDomain = *machine.Spec.FailureDomain
	}

	instanceSpec.Tags = compute.InstanceTags(&openStackMachine.Spec, openStackCluster)

	return &instanceSpec, nil
}

// getManagedSecurityGroup returns the ID of the security group managed by the
// OpenStackCluster whether it's a control plane or a worker machine.
func getManagedSecurityGroup(openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine) *string {
	if openStackCluster.Spec.ManagedSecurityGroups == nil {
		return nil
	}

	if util.IsControlPlaneMachine(machine) {
		if openStackCluster.Status.ControlPlaneSecurityGroup != nil {
			return &openStackCluster.Status.ControlPlaneSecurityGroup.ID
		}
	} else {
		if openStackCluster.Status.WorkerSecurityGroup != nil {
			return &openStackCluster.Status.WorkerSecurityGroup.ID
		}
	}

	return nil
}

func (r *OpenStackMachineReconciler) reconcileLoadBalancerMember(scope *scope.WithLogger, openStackCluster *infrav1.OpenStackCluster, openStackMachine *infrav1.OpenStackMachine, instanceNS *compute.InstanceNetworkStatus, clusterResourceName string) error {
	ip := instanceNS.IP(openStackCluster.Status.Network.Name)
	loadbalancerService, err := loadbalancer.NewService(scope)
	if err != nil {
		return err
	}

	return loadbalancerService.ReconcileLoadBalancerMember(openStackCluster, openStackMachine, clusterResourceName, ip)
}

// OpenStackClusterToOpenStackMachines is a handler.ToRequestsFunc to be used to enqeue requests for reconciliation
// of OpenStackMachines.
func (r *OpenStackMachineReconciler) OpenStackClusterToOpenStackMachines(ctx context.Context) handler.MapFunc {
	log := ctrl.LoggerFrom(ctx)
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		c, ok := o.(*infrav1.OpenStackCluster)
		if !ok {
			panic(fmt.Sprintf("Expected a OpenStackCluster but got a %T", o))
		}

		log := log.WithValues("objectMapper", "openStackClusterToOpenStackMachine", "namespace", c.Namespace, "openStackCluster", c.Name)

		// Don't handle deleted OpenStackClusters
		if !c.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(4).Info("OpenStackClusters has a deletion timestamp, skipping mapping.")
			return nil
		}

		cluster, err := util.GetOwnerCluster(ctx, r.Client, c.ObjectMeta)
		switch {
		case apierrors.IsNotFound(err) || cluster == nil:
			log.V(4).Info("Cluster for OpenStackCluster not found, skipping mapping.")
			return nil
		case err != nil:
			log.Error(err, "Failed to get owning cluster, skipping mapping.")
			return nil
		}

		return r.requestsForCluster(ctx, log, cluster.Namespace, cluster.Name)
	}
}

func (r *OpenStackMachineReconciler) getBootstrapData(ctx context.Context, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine) (string, error) {
	if machine.Spec.Bootstrap.DataSecretName == nil {
		return "", errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: machine.Namespace, Name: *machine.Spec.Bootstrap.DataSecretName}
	if err := r.Client.Get(ctx, key, secret); err != nil {
		return "", fmt.Errorf("failed to retrieve bootstrap data secret for Openstack Machine %s/%s: %w", machine.Namespace, openStackMachine.Name, err)
	}

	value, ok := secret.Data["value"]
	if !ok {
		return "", errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return base64.StdEncoding.EncodeToString(value), nil
}

func (r *OpenStackMachineReconciler) requeueOpenStackMachinesForUnpausedCluster(ctx context.Context) handler.MapFunc {
	log := ctrl.LoggerFrom(ctx)
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		c, ok := o.(*clusterv1.Cluster)
		if !ok {
			panic(fmt.Sprintf("Expected a Cluster but got a %T", o))
		}

		log := log.WithValues("objectMapper", "clusterToOpenStackMachine", "namespace", c.Namespace, "cluster", c.Name)

		// Don't handle deleted clusters
		if !c.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(4).Info("Cluster has a deletion timestamp, skipping mapping.")
			return nil
		}

		return r.requestsForCluster(ctx, log, c.Namespace, c.Name)
	}
}

func (r *OpenStackMachineReconciler) requestsForCluster(ctx context.Context, log logr.Logger, namespace, name string) []ctrl.Request {
	labels := map[string]string{clusterv1.ClusterNameLabel: name}
	machineList := &clusterv1.MachineList{}
	if err := r.Client.List(ctx, machineList, client.InNamespace(namespace), client.MatchingLabels(labels)); err != nil {
		log.Error(err, "Failed to get owned Machines, skipping mapping.")
		return nil
	}

	result := make([]ctrl.Request, 0, len(machineList.Items))
	for _, m := range machineList.Items {
		if m.Spec.InfrastructureRef.Name != "" {
			result = append(result, ctrl.Request{NamespacedName: client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.InfrastructureRef.Name}})
		}
	}
	return result
}

func (r *OpenStackMachineReconciler) getInfraCluster(ctx context.Context, cluster *clusterv1.Cluster, openStackMachine *infrav1.OpenStackMachine) (*infrav1.OpenStackCluster, error) {
	openStackCluster := &infrav1.OpenStackCluster{}
	openStackClusterName := client.ObjectKey{
		Namespace: openStackMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, openStackClusterName, openStackCluster); err != nil {
		return nil, err
	}
	return openStackCluster, nil
}
