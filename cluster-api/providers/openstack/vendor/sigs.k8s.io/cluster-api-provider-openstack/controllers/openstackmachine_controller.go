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
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/compute"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/loadbalancer"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
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

	clientScope, err := r.ScopeFactory.NewClientScopeFromObject(ctx, r.Client, r.CaCertificates, log, openStackMachine, infraCluster)
	if err != nil {
		return reconcile.Result{}, err
	}
	scope := scope.NewWithLogger(clientScope, log)

	clusterResourceName := fmt.Sprintf("%s-%s", cluster.Namespace, cluster.Name)

	// Handle deleted machines
	if !openStackMachine.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, scope, clusterResourceName, infraCluster, machine, openStackMachine)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, scope, clusterResourceName, infraCluster, machine, openStackMachine)
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
	log := ctrl.LoggerFrom(ctx)

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.OpenStackMachine{}).
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrav1.SchemeGroupVersion.WithKind("OpenStackMachine"))),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), r.WatchFilterValue)).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(r.requeueOpenStackMachinesForUnpausedCluster(ctx)),
			builder.WithPredicates(predicates.ClusterUnpausedAndInfrastructureReady(log)),
		).
		// NOTE: we don't watch OpenStackCluster here, even though the
		// OpenStackMachine controller directly requires values from
		// OpenStackCluster. The reason is that we are already observing Cluster
		// with the ClusterUnpausedAndInfrastructureReady predicate. The only
		// fields in OpenStackCluster we are interested in are dependent on
		// InfrastructureReady, so we don't need to watch both.
		Watches(
			&ipamv1.IPAddressClaim{},
			handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &infrav1.OpenStackMachine{}),
		).
		Watches(
			&infrav1alpha1.OpenStackServer{},
			handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &infrav1.OpenStackMachine{}),
			builder.WithPredicates(OpenStackServerReconcileComplete(log)),
		).
		Complete(r)
}

func (r *OpenStackMachineReconciler) reconcileDelete(ctx context.Context, scope *scope.WithLogger, clusterResourceName string, openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine) (ctrl.Result, error) { //nolint:unparam
	scope.Logger().Info("Reconciling Machine delete")

	computeService, err := compute.NewService(scope)
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

	machineServer, err := r.getMachineServer(ctx, openStackMachine)
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	var instanceStatus *compute.InstanceStatus
	if machineServer != nil && machineServer.Status.InstanceID != nil {
		instanceStatus, err = computeService.GetInstanceStatus(*machineServer.Status.InstanceID)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if util.IsControlPlaneMachine(machine) {
		if err := removeAPIServerEndpoint(scope, openStackCluster, openStackMachine, instanceStatus, clusterResourceName); err != nil {
			return ctrl.Result{}, err
		}
	}

	if machineServer != nil {
		scope.Logger().Info("Deleting server", "name", machineServer.Name)
		if err := r.Client.Delete(ctx, machineServer); err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceDeleteFailedReason, clusterv1.ConditionSeverityError, "Deleting instance failed: %v", err)
			return ctrl.Result{}, err
		}
		// If the server was found, we need to wait for it to be deleted before
		// removing the OpenStackMachine finalizer.
		scope.Logger().Info("Waiting for server to be deleted before removing finalizer")
		return ctrl.Result{}, nil
	}

	controllerutil.RemoveFinalizer(openStackMachine, infrav1.MachineFinalizer)
	scope.Logger().Info("Reconciled Machine delete successfully")
	return ctrl.Result{}, nil
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

	if controllerutil.AddFinalizer(openStackMachine, infrav1.MachineFinalizer) {
		return ctrl.Result{}, nil
	}

	scope.Logger().Info("Reconciling Machine")

	machineServer, waitingForServer, err := r.reconcileMachineServer(ctx, scope, openStackMachine, openStackCluster, machine)
	if err != nil || waitingForServer {
		return ctrl.Result{}, err
	}

	computeService, err := compute.NewService(scope)
	if err != nil {
		return ctrl.Result{}, err
	}

	// instanceStatus is required for the API server load balancer and floating IP reconciliation
	// when Octavia is enabled.
	var instanceStatus *compute.InstanceStatus
	if instanceStatus, err = computeService.GetInstanceStatus(*machineServer.Status.InstanceID); err != nil {
		return ctrl.Result{}, err
	}

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

	result := r.reconcileMachineState(scope, openStackMachine, machine, machineServer)
	if result != nil {
		return *result, nil
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

// reconcileMachineState updates the conditions of the OpenStackMachine instance based on the instance state
// and sets the ProviderID and Ready fields when the instance is active.
// It returns a reconcile request if the instance is not yet active.
func (r *OpenStackMachineReconciler) reconcileMachineState(scope *scope.WithLogger, openStackMachine *infrav1.OpenStackMachine, machine *clusterv1.Machine, openStackServer *infrav1alpha1.OpenStackServer) *ctrl.Result {
	switch *openStackServer.Status.InstanceState {
	case infrav1.InstanceStateActive:
		scope.Logger().Info("Machine instance state is ACTIVE", "id", openStackServer.Status.InstanceID)
		conditions.MarkTrue(openStackMachine, infrav1.InstanceReadyCondition)

		// Set properties required by CAPI machine controller
		var region string
		if openStackMachine.Spec.IdentityRef != nil {
			region = openStackMachine.Spec.IdentityRef.Region
		}
		openStackMachine.Spec.ProviderID = ptr.To(fmt.Sprintf("openstack://%s/%s", region, *openStackServer.Status.InstanceID))
		openStackMachine.Status.InstanceID = openStackServer.Status.InstanceID
		openStackMachine.Status.Ready = true
	case infrav1.InstanceStateError:
		// If the machine has a NodeRef then it must have been working at some point,
		// so the error could be something temporary.
		// If not, it is more likely a configuration error so we set failure and never retry.
		scope.Logger().Info("Machine instance state is ERROR", "id", openStackServer.Status.InstanceID)
		if machine.Status.NodeRef == nil {
			err := fmt.Errorf("instance state %v is unexpected", openStackServer.Status.InstanceState)
			openStackMachine.SetFailure(capierrors.UpdateMachineError, err)
		}
		conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceStateErrorReason, clusterv1.ConditionSeverityError, "")
		return &ctrl.Result{}
	case infrav1.InstanceStateDeleted:
		// we should avoid further actions for DELETED VM
		scope.Logger().Info("Machine instance state is DELETED, no actions")
		conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceDeletedReason, clusterv1.ConditionSeverityError, "")
		return &ctrl.Result{}
	case infrav1.InstanceStateBuild, infrav1.InstanceStateUndefined:
		scope.Logger().Info("Waiting for instance to become ACTIVE", "id", openStackServer.Status.InstanceID, "status", openStackServer.Status.InstanceState)
		return &ctrl.Result{RequeueAfter: waitForBuildingInstanceToReconcile}
	default:
		// The other state is normal (for example, migrating, shutoff) but we don't want to proceed until it's ACTIVE
		// due to potential conflict or unexpected actions
		scope.Logger().Info("Waiting for instance to become ACTIVE", "id", openStackServer.Status.InstanceID, "status", openStackServer.Status.InstanceState)
		conditions.MarkUnknown(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceNotReadyReason, "Instance state is not handled: %v", openStackServer.Status.InstanceState)

		return &ctrl.Result{RequeueAfter: waitForInstanceBecomeActiveToReconcile}
	}
	return nil
}

func (r *OpenStackMachineReconciler) getMachineServer(ctx context.Context, openStackMachine *infrav1.OpenStackMachine) (*infrav1alpha1.OpenStackServer, error) {
	machineServer := &infrav1alpha1.OpenStackServer{}
	machineServerName := client.ObjectKey{
		Namespace: openStackMachine.Namespace,
		Name:      openStackMachine.Name,
	}
	err := r.Client.Get(ctx, machineServerName, machineServer)
	if err != nil {
		return nil, err
	}
	return machineServer, nil
}

// openStackMachineSpecToOpenStackServerSpec converts an OpenStackMachineSpec to an OpenStackServerSpec.
// It returns the OpenStackServerSpec object and an error if there is any.
func openStackMachineSpecToOpenStackServerSpec(openStackMachineSpec *infrav1.OpenStackMachineSpec, identityRef infrav1.OpenStackIdentityReference, tags []string, failureDomain string, userDataRef *corev1.LocalObjectReference, defaultSecGroup *string, defaultNetworkID string) *infrav1alpha1.OpenStackServerSpec {
	openStackServerSpec := &infrav1alpha1.OpenStackServerSpec{
		AdditionalBlockDevices:            openStackMachineSpec.AdditionalBlockDevices,
		ConfigDrive:                       openStackMachineSpec.ConfigDrive,
		Flavor:                            openStackMachineSpec.Flavor,
		IdentityRef:                       identityRef,
		Image:                             openStackMachineSpec.Image,
		RootVolume:                        openStackMachineSpec.RootVolume,
		ServerMetadata:                    openStackMachineSpec.ServerMetadata,
		SSHKeyName:                        openStackMachineSpec.SSHKeyName,
		ServerGroup:                       openStackMachineSpec.ServerGroup,
		SchedulerHintAdditionalProperties: openStackMachineSpec.SchedulerHintAdditionalProperties,
	}

	if len(tags) > 0 {
		openStackServerSpec.Tags = tags
	}

	if failureDomain != "" {
		openStackServerSpec.AvailabilityZone = &failureDomain
	}

	if userDataRef != nil {
		openStackServerSpec.UserDataRef = userDataRef
	}

	if openStackMachineSpec.Trunk {
		openStackServerSpec.Trunk = ptr.To(true)
	}

	if openStackMachineSpec.FloatingIPPoolRef != nil {
		openStackServerSpec.FloatingIPPoolRef = openStackMachineSpec.FloatingIPPoolRef
	}

	// If not ports are provided we create one.
	// Ports must have a network so if none is provided we use the default network.
	serverPorts := openStackMachineSpec.Ports
	if len(openStackMachineSpec.Ports) == 0 {
		serverPorts = make([]infrav1.PortOpts, 1)
	}
	for i := range serverPorts {
		if serverPorts[i].Network == nil {
			serverPorts[i].Network = &infrav1.NetworkParam{
				ID: &defaultNetworkID,
			}
		}
		if len(serverPorts[i].SecurityGroups) == 0 && defaultSecGroup != nil {
			serverPorts[i].SecurityGroups = []infrav1.SecurityGroupParam{
				{
					ID: defaultSecGroup,
				},
			}
		}
		if len(openStackMachineSpec.SecurityGroups) > 0 {
			serverPorts[i].SecurityGroups = append(serverPorts[i].SecurityGroups, openStackMachineSpec.SecurityGroups...)
		}
	}
	openStackServerSpec.Ports = serverPorts

	return openStackServerSpec
}

// reconcileMachineServer reconciles the OpenStackServer object for the OpenStackMachine.
// It returns the OpenStackServer object and a boolean indicating if the OpenStackServer is ready.
func (r *OpenStackMachineReconciler) reconcileMachineServer(ctx context.Context, scope *scope.WithLogger, openStackMachine *infrav1.OpenStackMachine, openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine) (*infrav1alpha1.OpenStackServer, bool, error) {
	var server *infrav1alpha1.OpenStackServer
	server, err := r.getOrCreateMachineServer(ctx, openStackCluster, openStackMachine, machine)
	if err != nil {
		// If an error occurs while getting or creating the OpenStackServer,
		// we won't requeue the request so reconcileNormal can add conditions to the OpenStackMachine
		// and we can see the error in the logs.
		scope.Logger().Error(err, "Failed to get or create OpenStackServer")
		return server, false, err
	}
	if !server.Status.Ready {
		scope.Logger().Info("Waiting for OpenStackServer to be ready", "name", server.Name)
		return server, true, nil
	}
	return server, false, nil
}

// getOrCreateMachineServer gets or creates the OpenStackServer object for the OpenStackMachine.
// It returns the OpenStackServer object and an error if there is any.
func (r *OpenStackMachineReconciler) getOrCreateMachineServer(ctx context.Context, openStackCluster *infrav1.OpenStackCluster, openStackMachine *infrav1.OpenStackMachine, machine *clusterv1.Machine) (*infrav1alpha1.OpenStackServer, error) {
	if machine.Spec.Bootstrap.DataSecretName == nil {
		return nil, errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}
	userDataRef := &corev1.LocalObjectReference{
		Name: *machine.Spec.Bootstrap.DataSecretName,
	}

	var failureDomain string
	if machine.Spec.FailureDomain != nil {
		failureDomain = *machine.Spec.FailureDomain
	}
	machineServer, err := r.getMachineServer(ctx, openStackMachine)

	if client.IgnoreNotFound(err) != nil {
		return nil, err
	}
	if apierrors.IsNotFound(err) {
		// Use credentials from the machine object by default, falling back to cluster credentials.
		identityRef := func() infrav1.OpenStackIdentityReference {
			if openStackMachine.Spec.IdentityRef != nil {
				return *openStackMachine.Spec.IdentityRef
			}
			return openStackCluster.Spec.IdentityRef
		}()
		machineServerSpec := openStackMachineSpecToOpenStackServerSpec(&openStackMachine.Spec, identityRef, compute.InstanceTags(&openStackMachine.Spec, openStackCluster), failureDomain, userDataRef, getManagedSecurityGroup(openStackCluster, machine), openStackCluster.Status.Network.ID)
		machineServer = &infrav1alpha1.OpenStackServer{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					clusterv1.ClusterNameLabel: openStackCluster.Labels[clusterv1.ClusterNameLabel],
				},
				Name:      openStackMachine.Name,
				Namespace: openStackMachine.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: openStackMachine.APIVersion,
						Kind:       openStackMachine.Kind,
						Name:       openStackMachine.Name,
						UID:        openStackMachine.UID,
					},
				},
			},
			Spec: *machineServerSpec,
		}

		if err := r.Client.Create(ctx, machineServer); err != nil {
			return nil, fmt.Errorf("failed to create machine server: %w", err)
		}
	}
	return machineServer, nil
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

// getManagedSecurityGroup returns the ID of the security group managed by the
// OpenStackCluster whether it's a control plane or a worker machine.
func getManagedSecurityGroup(openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine) *string {
	if openStackCluster.Spec.ManagedSecurityGroups == nil {
		return nil
	}

	if machine == nil {
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
