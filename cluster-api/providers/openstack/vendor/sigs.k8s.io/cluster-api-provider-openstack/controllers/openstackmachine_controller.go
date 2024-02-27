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
	"reflect"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
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
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
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
)

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

func (r *OpenStackMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
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
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	scope, err := r.ScopeFactory.NewClientScopeFromMachine(ctx, r.Client, openStackMachine, r.CaCertificates, log)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Handle deleted machines
	if !openStackMachine.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(scope, cluster, infraCluster, machine, openStackMachine)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, scope, cluster, infraCluster, machine, openStackMachine)
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
		For(
			&infrav1.OpenStackMachine{},
			builder.WithPredicates(
				predicate.Funcs{
					// Avoid reconciling if the event triggering the reconciliation is related to incremental status updates
					UpdateFunc: func(e event.UpdateEvent) bool {
						oldMachine := e.ObjectOld.(*infrav1.OpenStackMachine).DeepCopy()
						newMachine := e.ObjectNew.(*infrav1.OpenStackMachine).DeepCopy()
						oldMachine.Status = infrav1.OpenStackMachineStatus{}
						newMachine.Status = infrav1.OpenStackMachineStatus{}
						oldMachine.ObjectMeta.ResourceVersion = ""
						newMachine.ObjectMeta.ResourceVersion = ""
						return !reflect.DeepEqual(oldMachine, newMachine)
					},
				},
			),
		).
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
		Complete(r)
}

func (r *OpenStackMachineReconciler) reconcileDelete(scope scope.Scope, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine) (ctrl.Result, error) { //nolint:unparam
	scope.Logger().Info("Reconciling Machine delete")

	clusterName := fmt.Sprintf("%s-%s", cluster.ObjectMeta.Namespace, cluster.Name)

	computeService, err := compute.NewService(scope)
	if err != nil {
		return ctrl.Result{}, err
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return ctrl.Result{}, err
	}

	if openStackCluster.Spec.APIServerLoadBalancer.Enabled {
		loadBalancerService, err := loadbalancer.NewService(scope)
		if err != nil {
			return ctrl.Result{}, err
		}

		err = loadBalancerService.DeleteLoadBalancerMember(openStackCluster, machine, openStackMachine, clusterName)
		if err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.APIServerIngressReadyCondition, infrav1.LoadBalancerMemberErrorReason, clusterv1.ConditionSeverityWarning, "Machine could not be removed from load balancer: %v", err)
			return ctrl.Result{}, err
		}
	}

	instanceStatus, err := computeService.GetInstanceStatusByName(openStackMachine, openStackMachine.Name)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !openStackCluster.Spec.APIServerLoadBalancer.Enabled && util.IsControlPlaneMachine(machine) && openStackCluster.Spec.APIServerFloatingIP == "" {
		if instanceStatus != nil {
			instanceNS, err := instanceStatus.NetworkStatus()
			if err != nil {
				openStackMachine.SetFailure(
					capierrors.UpdateMachineError,
					fmt.Errorf("get network status for OpenStack instance %s with ID %s: %v", instanceStatus.Name(), instanceStatus.ID(), err),
				)
				return ctrl.Result{}, nil
			}

			addresses := instanceNS.Addresses()
			for _, address := range addresses {
				if address.Type == corev1.NodeExternalIP {
					if err = networkingService.DeleteFloatingIP(openStackMachine, address.Address); err != nil {
						conditions.MarkFalse(openStackMachine, infrav1.APIServerIngressReadyCondition, infrav1.FloatingIPErrorReason, clusterv1.ConditionSeverityError, "Deleting floating IP failed: %v", err)
						return ctrl.Result{}, fmt.Errorf("delete floating IP %q: %w", address.Address, err)
					}
				}
			}
		}
	}

	instanceSpec := machineToInstanceSpec(openStackCluster, machine, openStackMachine, "")

	if err := computeService.DeleteInstance(openStackCluster, openStackMachine, instanceStatus, instanceSpec); err != nil {
		conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceDeleteFailedReason, clusterv1.ConditionSeverityError, "Deleting instance failed: %v", err)
		return ctrl.Result{}, fmt.Errorf("delete instance: %w", err)
	}

	controllerutil.RemoveFinalizer(openStackMachine, infrav1.MachineFinalizer)
	scope.Logger().Info("Reconciled Machine delete successfully")
	return ctrl.Result{}, nil
}

func (r *OpenStackMachineReconciler) reconcileNormal(ctx context.Context, scope scope.Scope, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine) (_ ctrl.Result, reterr error) {
	// If the OpenStackMachine is in an error state, return early.
	if openStackMachine.Status.FailureReason != nil || openStackMachine.Status.FailureMessage != nil {
		scope.Logger().Info("Not reconciling machine in failed state. See openStackMachine.status.failureReason, openStackMachine.status.failureMessage, or previously logged error for details")
		return ctrl.Result{}, nil
	}

	// If the OpenStackMachine doesn't have our finalizer, add it.
	if controllerutil.AddFinalizer(openStackMachine, infrav1.MachineFinalizer) {
		// Register the finalizer immediately to avoid orphaning OpenStack resources on delete
		return ctrl.Result{}, nil
	}

	if !cluster.Status.InfrastructureReady {
		scope.Logger().Info("Cluster infrastructure is not ready yet, requeuing machine")
		conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.WaitingForClusterInfrastructureReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{RequeueAfter: waitForClusterInfrastructureReadyDuration}, nil
	}

	// Make sure bootstrap data is available and populated.
	if machine.Spec.Bootstrap.DataSecretName == nil {
		scope.Logger().Info("Bootstrap data secret reference is not yet available")
		conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.WaitingForBootstrapDataReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}
	userData, err := r.getBootstrapData(ctx, machine, openStackMachine)
	if err != nil {
		return ctrl.Result{}, err
	}
	scope.Logger().Info("Reconciling Machine")

	clusterName := fmt.Sprintf("%s-%s", cluster.ObjectMeta.Namespace, cluster.Name)

	computeService, err := compute.NewService(scope)
	if err != nil {
		return ctrl.Result{}, err
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return ctrl.Result{}, err
	}

	instanceStatus, err := r.getOrCreate(scope.Logger(), cluster, openStackCluster, machine, openStackMachine, computeService, userData)
	if err != nil {
		// Conditions set in getOrCreate
		return ctrl.Result{}, err
	}

	// TODO(sbueringer) From CAPA: TODO(ncdc): move this validation logic into a validating webhook (for us: create validation logic in webhook)

	openStackMachine.Spec.ProviderID = pointer.String(fmt.Sprintf("openstack:///%s", instanceStatus.ID()))
	openStackMachine.Spec.InstanceID = pointer.String(instanceStatus.ID())

	state := instanceStatus.State()
	openStackMachine.Status.InstanceState = &state

	instanceNS, err := instanceStatus.NetworkStatus()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("get network status: %w", err)
	}

	addresses := instanceNS.Addresses()
	openStackMachine.Status.Addresses = addresses

	switch instanceStatus.State() {
	case infrav1.InstanceStateActive:
		scope.Logger().Info("Machine instance state is ACTIVE", "id", instanceStatus.ID())
		conditions.MarkTrue(openStackMachine, infrav1.InstanceReadyCondition)
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

	if openStackCluster.Spec.APIServerLoadBalancer.Enabled {
		err = r.reconcileLoadBalancerMember(scope, openStackCluster, openStackMachine, instanceNS, clusterName)
		if err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.APIServerIngressReadyCondition, infrav1.LoadBalancerMemberErrorReason, clusterv1.ConditionSeverityError, "Reconciling load balancer member failed: %v", err)
			return ctrl.Result{}, fmt.Errorf("reconcile load balancer member: %w", err)
		}
	} else if !openStackCluster.Spec.DisableAPIServerFloatingIP {
		floatingIPAddress := openStackCluster.Spec.ControlPlaneEndpoint.Host
		if openStackCluster.Spec.APIServerFloatingIP != "" {
			floatingIPAddress = openStackCluster.Spec.APIServerFloatingIP
		}
		fp, err := networkingService.GetOrCreateFloatingIP(openStackMachine, openStackCluster, clusterName, floatingIPAddress)
		if err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.APIServerIngressReadyCondition, infrav1.FloatingIPErrorReason, clusterv1.ConditionSeverityError, "Floating IP cannot be obtained or created: %v", err)
			return ctrl.Result{}, fmt.Errorf("get or create floating IP %q: %w", floatingIPAddress, err)
		}
		port, err := computeService.GetManagementPort(openStackCluster, instanceStatus)
		if err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.APIServerIngressReadyCondition, infrav1.FloatingIPErrorReason, clusterv1.ConditionSeverityError, "Obtaining management port for control plane machine failed: %v", err)
			return ctrl.Result{}, fmt.Errorf("get management port for control plane machine: %w", err)
		}

		if fp.PortID != "" {
			scope.Logger().Info("Floating IP already associated to a port", "id", fp.ID, "fixedIP", fp.FixedIP, "portID", port.ID)
		} else {
			err = networkingService.AssociateFloatingIP(openStackMachine, fp, port.ID)
			if err != nil {
				conditions.MarkFalse(openStackMachine, infrav1.APIServerIngressReadyCondition, infrav1.FloatingIPErrorReason, clusterv1.ConditionSeverityError, "Associating floating IP failed: %v", err)
				return ctrl.Result{}, fmt.Errorf("associate floating IP %q with port %q: %w", fp.FloatingIP, port.ID, err)
			}
		}
	}
	conditions.MarkTrue(openStackMachine, infrav1.APIServerIngressReadyCondition)

	scope.Logger().Info("Reconciled Machine create successfully")
	return ctrl.Result{}, nil
}

func (r *OpenStackMachineReconciler) getOrCreate(logger logr.Logger, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine, computeService *compute.Service, userData string) (*compute.InstanceStatus, error) {
	instanceStatus, err := computeService.GetInstanceStatusByName(openStackMachine, openStackMachine.Name)
	if err != nil {
		logger.Info("Unable to get OpenStack instance", "name", openStackMachine.Name)
		conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.OpenStackErrorReason, clusterv1.ConditionSeverityError, err.Error())
		return nil, err
	}

	if instanceStatus == nil {
		instanceSpec := machineToInstanceSpec(openStackCluster, machine, openStackMachine, userData)
		logger.Info("Machine does not exist, creating Machine", "name", openStackMachine.Name)
		instanceStatus, err = computeService.CreateInstance(openStackMachine, openStackCluster, instanceSpec, cluster.Name, false)
		if err != nil {
			conditions.MarkFalse(openStackMachine, infrav1.InstanceReadyCondition, infrav1.InstanceCreateFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return nil, fmt.Errorf("create OpenStack instance: %w", err)
		}
	}

	return instanceStatus, nil
}

func machineToInstanceSpec(openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine, userData string) *compute.InstanceSpec {
	instanceSpec := compute.InstanceSpec{
		Name:          openStackMachine.Name,
		Image:         openStackMachine.Spec.Image,
		ImageUUID:     openStackMachine.Spec.ImageUUID,
		Flavor:        openStackMachine.Spec.Flavor,
		SSHKeyName:    openStackMachine.Spec.SSHKeyName,
		UserData:      userData,
		Metadata:      openStackMachine.Spec.ServerMetadata,
		ConfigDrive:   openStackMachine.Spec.ConfigDrive != nil && *openStackMachine.Spec.ConfigDrive,
		RootVolume:    openStackMachine.Spec.RootVolume,
		ServerGroupID: openStackMachine.Spec.ServerGroupID,
		Trunk:         openStackMachine.Spec.Trunk,
	}

	// Add the failure domain only if specified
	if machine.Spec.FailureDomain != nil {
		instanceSpec.FailureDomain = *machine.Spec.FailureDomain
	}

	machineTags := []string{}

	// Append machine specific tags
	machineTags = append(machineTags, openStackMachine.Spec.Tags...)

	// Append cluster scope tags
	machineTags = append(machineTags, openStackCluster.Spec.Tags...)

	// tags need to be unique or the "apply tags" call will fail.
	deduplicate := func(tags []string) []string {
		seen := make(map[string]struct{}, len(machineTags))
		unique := make([]string, 0, len(machineTags))
		for _, tag := range tags {
			if _, ok := seen[tag]; !ok {
				seen[tag] = struct{}{}
				unique = append(unique, tag)
			}
		}
		return unique
	}
	machineTags = deduplicate(machineTags)

	instanceSpec.Tags = machineTags

	instanceSpec.SecurityGroups = openStackMachine.Spec.SecurityGroups
	if openStackCluster.Spec.ManagedSecurityGroups {
		var managedSecurityGroup string
		if util.IsControlPlaneMachine(machine) && openStackCluster.Status.ControlPlaneSecurityGroup != nil {
			managedSecurityGroup = openStackCluster.Status.ControlPlaneSecurityGroup.ID
		} else if openStackCluster.Status.WorkerSecurityGroup != nil {
			managedSecurityGroup = openStackCluster.Status.WorkerSecurityGroup.ID
		}

		instanceSpec.SecurityGroups = append(instanceSpec.SecurityGroups, infrav1.SecurityGroupFilter{
			ID: managedSecurityGroup,
		})
	}

	instanceSpec.Ports = openStackMachine.Spec.Ports

	return &instanceSpec
}

func (r *OpenStackMachineReconciler) reconcileLoadBalancerMember(scope scope.Scope, openStackCluster *infrav1.OpenStackCluster, openStackMachine *infrav1.OpenStackMachine, instanceNS *compute.InstanceNetworkStatus, clusterName string) error {
	ip := instanceNS.IP(openStackCluster.Status.Network.Name)
	loadbalancerService, err := loadbalancer.NewService(scope)
	if err != nil {
		return err
	}

	return loadbalancerService.ReconcileLoadBalancerMember(openStackCluster, openStackMachine, clusterName, ip)
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
