/*
Copyright 2024 The Kubernetes Authors.

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

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/events"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	ipamv1 "sigs.k8s.io/cluster-api/api/ipam/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	conditions "sigs.k8s.io/cluster-api/util/conditions"
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

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"

	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/orc"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/names"
)

const (
	SpecHashAnnotation = "infrastructure.cluster.x-k8s.io/spec-hash"
)

// OpenStackServerReconciler reconciles a OpenStackServer object.
type OpenStackServerReconciler struct {
	Client           client.Client
	Recorder         events.EventRecorder
	WatchFilterValue string
	ScopeFactory     scope.Factory
	CaCertificates   []byte // PEM encoded ca certificates.

	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackservers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=ipaddressclaims;ipaddressclaims/status,verbs=get;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=ipaddresses;ipaddresses/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=servers;ports;volumes;images;flavors;keypairs;servergroups;networks;subnets;securitygroups;trunks;volumetypes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackclusteridentities,verbs=get;list;watch

func (r *OpenStackServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the OpenStackServer instance.
	openStackServer := &infrav1alpha1.OpenStackServer{}
	err := r.Client.Get(ctx, req.NamespacedName, openStackServer)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	patchHelper, err := patch.NewHelper(openStackServer, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	defer func() {
		// Propagate terminal errors
		terminalError := &capoerrors.TerminalError{}
		if errors.As(reterr, &terminalError) {
			conditions.Set(openStackServer, metav1.Condition{
				Type:    infrav1.InstanceReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  terminalError.Reason,
				Message: terminalError.Message,
			})
		}

		if err := patchServer(ctx, patchHelper, openStackServer); err != nil {
			result = ctrl.Result{}
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	logger := log.WithValues("OpenStackServer", klog.KObj(openStackServer))
	logger.Info("Reconciling OpenStackServer")

	cluster, err := getClusterFromMetadata(ctx, r.Client, openStackServer.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster != nil {
		if annotations.IsPaused(cluster, openStackServer) {
			logger.Info("OpenStackServer linked to a Cluster that is paused. Won't reconcile")
			return reconcile.Result{}, nil
		}
	}

	// Handle deleted servers
	if !openStackServer.DeletionTimestamp.IsZero() {
		return reconcile.Result{}, r.reconcileDelete(ctx, openStackServer)
	}

	// Handle non-deleted servers
	return r.reconcileNormal(ctx, openStackServer)
}

func patchServer(ctx context.Context, patchHelper *patch.Helper, openStackServer *infrav1alpha1.OpenStackServer, options ...patch.Option) error {
	options = append(options,
		patch.WithOwnedConditions{Conditions: []string{
			clusterv1.ReadyCondition,
			infrav1.InstanceReadyCondition,
			infrav1.FloatingAddressFromPoolReadyCondition,
		}},
	)

	return patchHelper.Patch(ctx, openStackServer, options...)
}

func (r *OpenStackServerReconciler) reconcileDelete(ctx context.Context, openStackServer *infrav1alpha1.OpenStackServer) error {
	log := ctrl.LoggerFrom(ctx).WithValues("OpenStackServer", klog.KObj(openStackServer))
	log.Info("Reconciling OpenStackServer delete")

	orcReconciler := &orc.Reconciler{Client: r.Client, Scheme: r.Scheme}
	done, err := orcReconciler.DeleteORCResources(ctx, openStackServer)
	if err != nil {
		return err
	}
	if !done {
		// Waiting for ORC Server deletion. The Owns() watch on
		// orcv1alpha1.Server will re-trigger the controller.
		log.Info("Waiting for ORC Server to be deleted")
		return nil
	}

	// Handle floating IP cleanup.
	if err := r.reconcileDeleteFloatingAddressFromPool(ctx, openStackServer); err != nil {
		return err
	}

	controllerutil.RemoveFinalizer(openStackServer, infrav1alpha1.OpenStackServerFinalizer)
	log.Info("Reconciled OpenStackServer deleted successfully")
	return nil
}

func IsServerTerminalError(server *infrav1alpha1.OpenStackServer) bool {
	if server.Status.InstanceState != nil && *server.Status.InstanceState == infrav1.InstanceStateError {
		return true
	}
	return false
}

func (r *OpenStackServerReconciler) reconcileNormal(ctx context.Context, openStackServer *infrav1alpha1.OpenStackServer) (_ ctrl.Result, reterr error) {
	// If the OpenStackServer is in an error state, return early.
	if IsServerTerminalError(openStackServer) {
		log := ctrl.LoggerFrom(ctx).WithValues("OpenStackServer", klog.KObj(openStackServer))
		log.Info("Not reconciling OpenStackServer in error state. See openStackServer.status or previously logged error for details")
		return ctrl.Result{}, nil
	}

	log := ctrl.LoggerFrom(ctx).WithValues("OpenStackServer", klog.KObj(openStackServer))
	log.Info("Reconciling OpenStackServer")

	// Add finalizer. We requeue so we never create resources unless we
	// have observed that the finalizer was successfully written.
	if controllerutil.AddFinalizer(openStackServer, infrav1alpha1.OpenStackServerFinalizer) {
		return ctrl.Result{}, nil
	}

	orcReconciler := &orc.Reconciler{Client: r.Client, Scheme: r.Scheme}
	orcResult, err := orcReconciler.Reconcile(ctx, openStackServer)
	if err != nil {
		// Terminal errors are caught by the defer block in the parent
		// Reconcile method, which sets the InstanceReadyCondition.
		return ctrl.Result{}, err
	}

	// Track ORC server name in status.
	if openStackServer.Status.Resources == nil {
		openStackServer.Status.Resources = &infrav1alpha1.ServerResources{}
	}
	openStackServer.Status.Resources.ORCServerName = orcResult.ORCServerName

	if !orcResult.Done {
		conditions.Set(openStackServer, metav1.Condition{
			Type:    infrav1.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.InstanceNotReadyReason,
			Message: "Waiting for ORC resources to become available",
		})
		// No explicit requeue — the Owns() watches on all ORC types
		// will re-trigger the controller when any sub-resource changes.
		return ctrl.Result{}, nil
	}

	// ORC Server is Available — update status from the ORC result.
	openStackServer.Status.InstanceID = ptr.To(orcResult.ServerID)
	openStackServer.Status.InstanceState = &orcResult.ServerState
	openStackServer.Status.Addresses = orcResult.Addresses

	// Handle floating IP — the IPAM claim flow stays in CAPO.
	floatingAddressClaim, waitingForFloatingAddress, err := r.reconcileFloatingAddressFromPool(ctx, openStackServer)
	if err != nil || waitingForFloatingAddress {
		return ctrl.Result{}, err
	}

	if floatingAddressClaim != nil {
		// Create an OpenStack scope only for floating IP association
		// (the only remaining Gophercloud usage).
		clientScope, err := r.ScopeFactory.NewClientScopeFromObject(ctx, r.Client, r.CaCertificates, log, openStackServer)
		if err != nil {
			conditions.Set(openStackServer, metav1.Condition{
				Type:    infrav1.FloatingAddressFromPoolReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  infrav1.FloatingAddressFromPoolErrorReason,
				Message: fmt.Sprintf("Failed to create scope for floating IP: %v", err),
			})
			return ctrl.Result{}, err
		}
		s := scope.NewWithLogger(clientScope, log)
		networkingService, err := networking.NewService(s)
		if err != nil {
			return ctrl.Result{}, err
		}

		if err := r.associateIPAddressFromIPAddressClaim(ctx, openStackServer, orcResult.ServerID, orcResult.Addresses, floatingAddressClaim, networkingService); err != nil {
			conditions.Set(openStackServer, metav1.Condition{
				Type:    infrav1.FloatingAddressFromPoolReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  infrav1.FloatingAddressFromPoolErrorReason,
				Message: fmt.Sprintf("Failed while associating ip from pool: %v", err),
			})
			return ctrl.Result{}, err
		}
		conditions.Set(openStackServer, metav1.Condition{
			Type:   infrav1.FloatingAddressFromPoolReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: infrav1.ReadyConditionReason,
		})
	}

	// Handle instance state.
	switch orcResult.ServerState {
	case infrav1.InstanceStateActive:
		log.Info("Server instance state is ACTIVE", "id", orcResult.ServerID)
		conditions.Set(openStackServer, metav1.Condition{
			Type:   infrav1.InstanceReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: infrav1.ReadyConditionReason,
		})
		// Set the Ready field for v1alpha1 compatibility with predicates
		openStackServer.Status.Ready = true
	case infrav1.InstanceStateError:
		log.Info("Server instance state is ERROR", "id", orcResult.ServerID)
		conditions.Set(openStackServer, metav1.Condition{
			Type:   infrav1.InstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrav1.InstanceStateErrorReason,
		})
		return ctrl.Result{}, nil
	case infrav1.InstanceStateDeleted:
		// we should avoid further actions for DELETED VM
		log.Info("Server instance state is DELETED, no actions")
		conditions.Set(openStackServer, metav1.Condition{
			Type:   infrav1.InstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrav1.InstanceDeletedReason,
		})
		return ctrl.Result{}, nil
	default:
		// Other states (BUILD, SHUTOFF, migrating, etc.) — wait for
		// ORC to update the status. No explicit requeue needed because
		// the Owns() watch on orcv1alpha1.Server will re-trigger us.
		log.Info("Waiting for instance to become ACTIVE", "id", orcResult.ServerID, "status", orcResult.ServerState)
		conditions.Set(openStackServer, metav1.Condition{
			Type:    infrav1.InstanceReadyCondition,
			Status:  metav1.ConditionUnknown,
			Reason:  infrav1.InstanceNotReadyReason,
			Message: fmt.Sprintf("Instance state is not handled: %s", orcResult.ServerState),
		})
		return ctrl.Result{}, nil
	}

	log.Info("Reconciled Server create successfully")
	return ctrl.Result{}, nil
}

func (r *OpenStackServerReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1alpha1.OpenStackServer{}).
		// Watch all ORC resource types this controller creates.
		// When any owned ORC object changes, the controller is
		// re-triggered for the owning OpenStackServer.
		Owns(&orcv1alpha1.Server{}).
		Owns(&orcv1alpha1.Port{}).
		Owns(&orcv1alpha1.Volume{}).
		Owns(&orcv1alpha1.Image{}).
		Owns(&orcv1alpha1.Flavor{}).
		Owns(&orcv1alpha1.KeyPair{}).
		Owns(&orcv1alpha1.ServerGroup{}).
		Owns(&orcv1alpha1.Network{}).
		Owns(&orcv1alpha1.Subnet{}).
		Owns(&orcv1alpha1.SecurityGroup{}).
		Owns(&orcv1alpha1.Trunk{}).
		Owns(&orcv1alpha1.VolumeType{}).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(r.requeueOpenStackServersForCluster(ctx)),
			builder.WithPredicates(predicates.ClusterPausedTransitionsOrInfrastructureProvisioned(mgr.GetScheme(), log)),
		).
		Watches(
			&ipamv1.IPAddressClaim{},
			handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &infrav1alpha1.OpenStackServer{}),
		).
		Complete(r)
}

// ── Floating IP helpers (minimal Gophercloud usage) ─────────────────

// reconcileFloatingAddressFromPool reconciles the floating IP address from the pool.
// It returns the IPAddressClaim and a boolean indicating if we are still waiting.
func (r *OpenStackServerReconciler) reconcileFloatingAddressFromPool(ctx context.Context, openStackServer *infrav1alpha1.OpenStackServer) (*ipamv1.IPAddressClaim, bool, error) {
	if openStackServer.Spec.FloatingIPPoolRef == nil {
		return nil, false, nil
	}
	claim, err := r.getOrCreateIPAddressClaimForFloatingAddress(ctx, openStackServer)
	if err != nil {
		conditions.Set(openStackServer, metav1.Condition{
			Type:    infrav1.FloatingAddressFromPoolReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.FloatingAddressFromPoolErrorReason,
			Message: fmt.Sprintf("Failed to reconcile floating IP claims: %v", err),
		})
		return nil, true, err
	}
	if claim.Status.AddressRef.Name == "" {
		r.Recorder.Eventf(openStackServer, nil, corev1.EventTypeNormal, "WaitingForIPAddressClaim", "WaitingForIPAddressClaim", "Waiting for IPAddressClaim %s/%s to be allocated", claim.Namespace, claim.Name)
		return claim, true, nil
	}
	conditions.Set(openStackServer, metav1.Condition{
		Type:   infrav1.FloatingAddressFromPoolReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrav1.ReadyConditionReason,
	})
	return claim, false, nil
}

// getOrCreateIPAddressClaimForFloatingAddress creates IPAddressClaim for the FloatingAddressFromPool if it does not exist yet.
func (r *OpenStackServerReconciler) getOrCreateIPAddressClaimForFloatingAddress(ctx context.Context, openStackServer *infrav1alpha1.OpenStackServer) (*ipamv1.IPAddressClaim, error) {
	log := ctrl.LoggerFrom(ctx)

	poolRef := openStackServer.Spec.FloatingIPPoolRef
	claimName := names.GetFloatingAddressClaimName(openStackServer.Name)
	claim := &ipamv1.IPAddressClaim{}

	err := r.Client.Get(ctx, client.ObjectKey{Namespace: openStackServer.Namespace, Name: claimName}, claim)
	if err == nil {
		return claim, nil
	} else if client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	claim = &ipamv1.IPAddressClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      claimName,
			Namespace: openStackServer.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: openStackServer.APIVersion,
					Kind:       openStackServer.Kind,
					Name:       openStackServer.Name,
					UID:        openStackServer.UID,
				},
			},
			Finalizers: []string{infrav1.IPClaimMachineFinalizer},
			Labels:     map[string]string{},
		},
		Spec: ipamv1.IPAddressClaimSpec{
			PoolRef: ipamv1.IPPoolReference{
				Name:     poolRef.Name,
				Kind:     poolRef.Kind,
				APIGroup: *poolRef.APIGroup,
			},
		},
	}

	if openStackServer.Labels[clusterv1.ClusterNameLabel] != "" {
		claim.Labels[clusterv1.ClusterNameLabel] = openStackServer.Labels[clusterv1.ClusterNameLabel]
	}

	if err := r.Client.Create(ctx, claim); err != nil {
		return nil, err
	}

	r.Recorder.Eventf(openStackServer, nil, corev1.EventTypeNormal, "CreatingIPAddressClaim", "CreatingIPAddressClaim", "Creating IPAddressClaim %s/%s", claim.Namespace, claim.Name)
	log.Info("Created IPAddressClaim", "name", claim.Name)
	return claim, nil
}

// associateIPAddressFromIPAddressClaim associates a floating IP from an
// IPAM claim with the server. This is the only code path that still
// uses Gophercloud (via the networking service) for Neutron API calls.
func (r *OpenStackServerReconciler) associateIPAddressFromIPAddressClaim(ctx context.Context, openStackServer *infrav1alpha1.OpenStackServer, serverID string, serverAddresses []corev1.NodeAddress, claim *ipamv1.IPAddressClaim, networkingService *networking.Service) error {
	address := &ipamv1.IPAddress{}
	addressKey := client.ObjectKey{Namespace: openStackServer.Namespace, Name: claim.Status.AddressRef.Name}

	if err := r.Client.Get(ctx, addressKey, address); err != nil {
		return err
	}

	// Check if already associated
	for _, addr := range serverAddresses {
		if addr.Address == address.Spec.Address {
			conditions.Set(openStackServer, metav1.Condition{
				Type:   infrav1.FloatingAddressFromPoolReadyCondition,
				Status: metav1.ConditionTrue,
				Reason: infrav1.ReadyConditionReason,
			})
			return nil
		}
	}

	fip, err := networkingService.GetFloatingIP(address.Spec.Address)
	if err != nil {
		return err
	}

	if fip == nil {
		conditions.Set(openStackServer, metav1.Condition{
			Type:    infrav1.FloatingAddressFromPoolReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.FloatingAddressFromPoolErrorReason,
			Message: "floating IP does not exist",
		})
		return fmt.Errorf("floating IP %q does not exist", address.Spec.Address)
	}

	port, err := networkingService.GetPortForExternalNetwork(serverID, fip.FloatingNetworkID)
	if err != nil {
		return fmt.Errorf("get port for floating IP %q: %w", fip.FloatingIP, err)
	}

	if port == nil {
		conditions.Set(openStackServer, metav1.Condition{
			Type:    infrav1.FloatingAddressFromPoolReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.FloatingAddressFromPoolErrorReason,
			Message: fmt.Sprintf("Can't find port for floating IP %q on external network %s", fip.FloatingIP, fip.FloatingNetworkID),
		})
		return fmt.Errorf("port for floating IP %q on network %s does not exist", fip.FloatingIP, fip.FloatingNetworkID)
	}

	if err = networkingService.AssociateFloatingIP(openStackServer, fip, port.ID); err != nil {
		return err
	}
	conditions.Set(openStackServer, metav1.Condition{
		Type:   infrav1.FloatingAddressFromPoolReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrav1.ReadyConditionReason,
	})
	return nil
}

func (r *OpenStackServerReconciler) reconcileDeleteFloatingAddressFromPool(ctx context.Context, openStackServer *infrav1alpha1.OpenStackServer) error {
	log := ctrl.LoggerFrom(ctx).WithValues("openStackServer", openStackServer.Name)
	log.Info("Reconciling delete floating address from pool")
	if openStackServer.Spec.FloatingIPPoolRef == nil {
		return nil
	}
	claimName := names.GetFloatingAddressClaimName(openStackServer.Name)
	claim := &ipamv1.IPAddressClaim{}
	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: openStackServer.Namespace, Name: claimName}, claim); err != nil {
		return client.IgnoreNotFound(err)
	}

	controllerutil.RemoveFinalizer(claim, infrav1.IPClaimMachineFinalizer)
	return r.Client.Update(ctx, claim)
}

// ── Predicates & helpers ────────────────────────────────────────────

// OpenStackServerReconcileComplete returns a predicate that determines if a OpenStackServer has finished reconciling.
func OpenStackServerReconcileComplete(log logr.Logger) predicate.Funcs {
	log = log.WithValues("predicate", "OpenStackServerReconcileComplete")

	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			log = log.WithValues("eventType", "create")

			server, ok := e.Object.(*infrav1alpha1.OpenStackServer)
			if !ok {
				log.V(4).Info("Expected OpenStackServer", "type", fmt.Sprintf("%T", e.Object))
				return false
			}
			log = log.WithValues("OpenStackServer", klog.KObj(server))

			if server.Status.Ready || IsServerTerminalError(server) {
				log.V(5).Info("OpenStackServer finished reconciling, allowing further processing")
				return true
			}
			log.V(5).Info("OpenStackServer is still reconciling, blocking further processing")
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			log := log.WithValues("eventType", "update")

			oldServer, ok := e.ObjectOld.(*infrav1alpha1.OpenStackServer)
			if !ok {
				log.V(4).Info("Expected OpenStackServer", "type", fmt.Sprintf("%T", e.ObjectOld))
				return false
			}
			log = log.WithValues("OpenStackServer", klog.KObj(oldServer))

			newServer, ok := e.ObjectNew.(*infrav1alpha1.OpenStackServer)
			if !ok {
				log.V(4).Info("Expected OpenStackServer (new)", "type", fmt.Sprintf("%T", e.ObjectNew))
				return false
			}

			oldFinished := oldServer.Status.Ready || IsServerTerminalError(oldServer)
			newFinished := newServer.Status.Ready || IsServerTerminalError(newServer)
			if !oldFinished && newFinished {
				log.V(5).Info("OpenStackServer finished reconciling, allowing further processing")
				return true
			}

			log.V(4).Info("OpenStackServer is still reconciling, blocking further processing")
			return false
		},
		DeleteFunc:  func(event.DeleteEvent) bool { return true },
		GenericFunc: func(event.GenericEvent) bool { return false },
	}
}

// getClusterFromMetadata returns the Cluster object (if present) using the object metadata.
func getClusterFromMetadata(ctx context.Context, c client.Client, obj metav1.ObjectMeta) (*clusterv1.Cluster, error) {
	if obj.Labels[clusterv1.ClusterNameLabel] == "" {
		return nil, nil
	}
	return util.GetClusterByName(ctx, c, obj.Namespace, obj.Labels[clusterv1.ClusterNameLabel])
}

// requeueOpenStackServersForCluster returns a handler.MapFunc that watches for
// Cluster changes and triggers reconciliation of all OpenStackServers in that cluster.
func (r *OpenStackServerReconciler) requeueOpenStackServersForCluster(ctx context.Context) handler.MapFunc {
	log := ctrl.LoggerFrom(ctx)
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		c, ok := o.(*clusterv1.Cluster)
		if !ok {
			panic(fmt.Sprintf("Expected a Cluster but got a %T", o))
		}

		log := log.WithValues("objectMapper", "clusterToOpenStackServer", "namespace", c.Namespace, "cluster", c.Name)

		if !c.DeletionTimestamp.IsZero() {
			log.V(4).Info("Cluster has a deletion timestamp, skipping mapping.")
			return nil
		}

		serverList := &infrav1alpha1.OpenStackServerList{}
		if err := r.Client.List(
			ctx,
			serverList,
			client.InNamespace(c.Namespace),
			client.MatchingLabels{clusterv1.ClusterNameLabel: c.Name},
		); err != nil {
			log.Error(err, "Failed to list OpenStackServers for cluster")
			return nil
		}

		requests := make([]ctrl.Request, 0, len(serverList.Items))
		for i := range serverList.Items {
			server := &serverList.Items[i]
			requests = append(requests, ctrl.Request{
				NamespacedName: client.ObjectKey{
					Namespace: server.Namespace,
					Name:      server.Name,
				},
			})
			log.V(5).Info("Queueing OpenStackServer for reconciliation", "server", server.Name)
		}

		return requests
	}
}
