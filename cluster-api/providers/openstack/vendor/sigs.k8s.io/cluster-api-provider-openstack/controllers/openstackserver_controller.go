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
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	ipamv1 "sigs.k8s.io/cluster-api/api/ipam/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
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
	"github.com/k-orc/openstack-resource-controller/v2/pkg/predicates"

	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/compute"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
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
	Recorder         record.EventRecorder
	WatchFilterValue string
	ScopeFactory     scope.Factory
	CaCertificates   []byte // PEM encoded ca certificates.

	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackservers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=ipaddressclaims;ipaddressclaims/status,verbs=get;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=ipaddresses;ipaddresses/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=images,verbs=get;list;watch

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

	clientScope, err := r.ScopeFactory.NewClientScopeFromObject(ctx, r.Client, r.CaCertificates, log, openStackServer)
	if err != nil {
		return reconcile.Result{}, err
	}
	scope := scope.NewWithLogger(clientScope, log)

	scope.Logger().Info("Reconciling OpenStackServer")

	cluster, err := getClusterFromMetadata(ctx, r.Client, openStackServer.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster != nil {
		if annotations.IsPaused(cluster, openStackServer) {
			scope.Logger().Info("OpenStackServer linked to a Cluster that is paused. Won't reconcile")
			return reconcile.Result{}, nil
		}
	}

	patchHelper, err := patch.NewHelper(openStackServer, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	defer func() {
		// Propagate terminal errors
		terminalError := &capoerrors.TerminalError{}
		if errors.As(reterr, &terminalError) {
			v1beta1conditions.MarkFalse(openStackServer, infrav1.InstanceReadyCondition, terminalError.Reason, clusterv1beta1.ConditionSeverityError, "%s", terminalError.Message)
		}

		if err := patchServer(ctx, patchHelper, openStackServer); err != nil {
			result = ctrl.Result{}
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	if !openStackServer.DeletionTimestamp.IsZero() {
		// When moving a cluster, we need to populate the server status with the resources
		// that were in another object's status.
		// This is because the status is not persisted across CAPI resources moves.
		if openStackServer.Status.Resolved == nil || openStackServer.Status.Resources == nil {
			if _, err := r.reconcileNormal(ctx, scope, openStackServer); err != nil {
				return ctrl.Result{}, err
			}
		}
		return reconcile.Result{}, r.reconcileDelete(scope, openStackServer)
	}

	return r.reconcileNormal(ctx, scope, openStackServer)
}

func patchServer(ctx context.Context, patchHelper *patch.Helper, openStackServer *infrav1alpha1.OpenStackServer, options ...patch.Option) error {
	// Always update the readyCondition by summarizing the state of other conditions.
	applicableConditions := []clusterv1beta1.ConditionType{
		infrav1.InstanceReadyCondition,
	}

	v1beta1conditions.SetSummary(openStackServer, v1beta1conditions.WithConditions(applicableConditions...))

	// Patch the object, ignoring conflicts on the conditions owned by this controller.
	// Also, if requested, we are adding additional options like e.g. Patch ObservedGeneration when issuing the
	// patch at the end of the reconcile loop.
	options = append(options,
		patch.WithOwnedConditions{Conditions: []string{
			clusterv1.ReadyCondition,
			string(infrav1.InstanceReadyCondition),
		}},
	)
	v1beta1conditions.SetSummary(openStackServer,
		v1beta1conditions.WithConditions(
			infrav1.InstanceReadyCondition,
		),
	)

	return patchHelper.Patch(ctx, openStackServer, options...)
}

func (r *OpenStackServerReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	const imageRefPath = "spec.image.imageRef.name"

	log := ctrl.LoggerFrom(ctx)

	// Index servers by referenced image
	if err := mgr.GetFieldIndexer().IndexField(ctx, &infrav1alpha1.OpenStackServer{}, imageRefPath, func(obj client.Object) []string {
		server, ok := obj.(*infrav1alpha1.OpenStackServer)
		if !ok {
			return nil
		}
		if server.Spec.Image.ImageRef == nil {
			return nil
		}
		return []string{server.Spec.Image.ImageRef.Name}
	}); err != nil {
		return fmt.Errorf("adding servers by image index: %w", err)
	}

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1alpha1.OpenStackServer{}).
		Watches(&orcv1alpha1.Image{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
				log = log.WithValues("watch", "Image")

				k8sClient := mgr.GetClient()
				serverList := &infrav1alpha1.OpenStackServerList{}
				if err := k8sClient.List(ctx, serverList, client.InNamespace(obj.GetNamespace()), client.MatchingFields{imageRefPath: obj.GetName()}); err != nil {
					log.Error(err, "listing OpenStackServers")
					return nil
				}

				requests := make([]reconcile.Request, len(serverList.Items))
				for i := range serverList.Items {
					server := &serverList.Items[i]
					request := &requests[i]

					request.Name = server.Name
					request.Namespace = server.Namespace
				}
				return requests
			}),
			builder.WithPredicates(predicates.NewBecameAvailable(mgr.GetLogger(), &orcv1alpha1.Image{})),
		).
		Watches(
			&ipamv1.IPAddressClaim{},
			handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &infrav1alpha1.OpenStackServer{}),
		).
		Complete(r)
}

func (r *OpenStackServerReconciler) reconcileDelete(scope *scope.WithLogger, openStackServer *infrav1alpha1.OpenStackServer) error {
	scope.Logger().Info("Reconciling Server delete")

	computeService, err := compute.NewService(scope)
	if err != nil {
		return err
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return err
	}

	// Check for any orphaned resources
	// N.B. Unlike resolveServerResources, we must always look for orphaned resources in the delete path.
	if err := adoptServerResources(scope, openStackServer); err != nil {
		return fmt.Errorf("adopting server resources: %w", err)
	}

	instanceStatus, err := getServerStatus(openStackServer, computeService)
	if err != nil {
		return err
	}

	// If no instance was created we currently need to check for orphaned volumes.
	if instanceStatus == nil {
		if err := computeService.DeleteVolumes(openStackServer.Name, openStackServer.Spec.RootVolume, openStackServer.Spec.AdditionalBlockDevices); err != nil {
			return fmt.Errorf("delete volumes: %w", err)
		}
	} else {
		if err := computeService.DeleteInstance(openStackServer, instanceStatus); err != nil {
			v1beta1conditions.MarkFalse(openStackServer, infrav1.InstanceReadyCondition, infrav1.InstanceDeleteFailedReason, clusterv1beta1.ConditionSeverityError, "Deleting instance failed: %v", err)
			return fmt.Errorf("delete instance: %w", err)
		}
	}

	trunkSupported, err := networkingService.IsTrunkExtSupported()
	if err != nil {
		return err
	}

	if openStackServer.Status.Resources != nil {
		portsStatus := openStackServer.Status.Resources.Ports
		for _, port := range portsStatus {
			if err := networkingService.DeleteInstanceTrunkAndPort(openStackServer, port, trunkSupported); err != nil {
				return fmt.Errorf("failed to delete port %q: %w", port.ID, err)
			}
		}
	}

	if err := r.reconcileDeleteFloatingAddressFromPool(scope, openStackServer); err != nil {
		return err
	}

	controllerutil.RemoveFinalizer(openStackServer, infrav1alpha1.OpenStackServerFinalizer)
	scope.Logger().Info("Reconciled Server deleted successfully")
	return nil
}

func IsServerTerminalError(server *infrav1alpha1.OpenStackServer) bool {
	if server.Status.InstanceState != nil && *server.Status.InstanceState == infrav1.InstanceStateError {
		return true
	}
	return false
}

func (r *OpenStackServerReconciler) reconcileNormal(ctx context.Context, scope *scope.WithLogger, openStackServer *infrav1alpha1.OpenStackServer) (_ ctrl.Result, reterr error) {
	// If the OpenStackServer is in an error state, return early.
	if IsServerTerminalError(openStackServer) {
		scope.Logger().Info("Not reconciling server in error state. See openStackServer.status or previously logged error for details")
		return ctrl.Result{}, nil
	}

	scope.Logger().Info("Reconciling Server create")

	labels := openStackServer.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
		openStackServer.SetLabels(labels)
	}

	changed, resolveDone, err := compute.ResolveServerSpec(ctx, scope, r.Client, openStackServer)
	if err != nil || !resolveDone {
		return ctrl.Result{}, err
	}

	// Also add the finalizer when writing resolved resources so we can start creating resources on the next reconcile.
	if controllerutil.AddFinalizer(openStackServer, infrav1alpha1.OpenStackServerFinalizer) {
		changed = true
	}

	// We requeue if we either added the finalizer or resolved server
	// resources. This means that we never create any resources unless we
	// have observed that the finalizer and resolved server resources were
	// successfully written in a previous transaction. This in turn means
	// that in the delete path we can be sure that if there are no resolved
	// resources then no resources were created.
	if changed {
		scope.Logger().V(6).Info("Server resources updated, requeuing")
		return ctrl.Result{}, nil
	}

	// Check for orphaned resources previously created but not written to the status
	if err := adoptServerResources(scope, openStackServer); err != nil {
		return ctrl.Result{}, fmt.Errorf("adopting server resources: %w", err)
	}
	computeService, err := compute.NewService(scope)
	if err != nil {
		return ctrl.Result{}, err
	}
	networkingService, err := networking.NewService(scope)
	if err != nil {
		return ctrl.Result{}, err
	}

	floatingAddressClaim, waitingForFloatingAddress, err := r.reconcileFloatingAddressFromPool(ctx, scope, openStackServer)
	if err != nil || waitingForFloatingAddress {
		return ctrl.Result{}, err
	}

	err = getOrCreateServerPorts(openStackServer, networkingService)
	if err != nil {
		return ctrl.Result{}, err
	}
	portIDs := GetPortIDs(openStackServer.Status.Resources.Ports)

	instanceStatus, err := r.getOrCreateServer(ctx, scope.Logger(), openStackServer, computeService, portIDs)
	if err != nil || instanceStatus == nil {
		// Conditions set in getOrCreateInstance
		return ctrl.Result{}, err
	}

	instanceNS, err := instanceStatus.NetworkStatus()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("get network status: %w", err)
	}

	if floatingAddressClaim != nil {
		if err := r.associateIPAddressFromIPAddressClaim(ctx, openStackServer, instanceStatus, instanceNS, floatingAddressClaim, networkingService); err != nil {
			v1beta1conditions.MarkFalse(openStackServer, infrav1.FloatingAddressFromPoolReadyCondition, infrav1.FloatingAddressFromPoolErrorReason, clusterv1beta1.ConditionSeverityError, "Failed while associating ip from pool: %v", err)
			return ctrl.Result{}, err
		}
		v1beta1conditions.MarkTrue(openStackServer, infrav1.FloatingAddressFromPoolReadyCondition)
	}

	state := instanceStatus.State()
	openStackServer.Status.InstanceID = ptr.To(instanceStatus.ID())
	openStackServer.Status.InstanceState = &state

	switch instanceStatus.State() {
	case infrav1.InstanceStateActive:
		scope.Logger().Info("Server instance state is ACTIVE", "id", instanceStatus.ID())
		v1beta1conditions.MarkTrue(openStackServer, infrav1.InstanceReadyCondition)
		openStackServer.Status.Ready = true
	case infrav1.InstanceStateError:
		scope.Logger().Info("Server instance state is ERROR", "id", instanceStatus.ID())
		v1beta1conditions.MarkFalse(openStackServer, infrav1.InstanceReadyCondition, infrav1.InstanceStateErrorReason, clusterv1beta1.ConditionSeverityError, "")
		return ctrl.Result{}, nil
	case infrav1.InstanceStateDeleted:
		// we should avoid further actions for DELETED VM
		scope.Logger().Info("Server instance state is DELETED, no actions")
		v1beta1conditions.MarkFalse(openStackServer, infrav1.InstanceReadyCondition, infrav1.InstanceDeletedReason, clusterv1beta1.ConditionSeverityError, "")
		return ctrl.Result{}, nil
	case infrav1.InstanceStateBuild, infrav1.InstanceStateUndefined:
		scope.Logger().Info("Waiting for instance to become ACTIVE", "id", instanceStatus.ID(), "status", instanceStatus.State())
		return ctrl.Result{RequeueAfter: waitForBuildingInstanceToReconcile}, nil
	default:
		// The other state is normal (for example, migrating, shutoff) but we don't want to proceed until it's ACTIVE
		// due to potential conflict or unexpected actions
		scope.Logger().Info("Waiting for instance to become ACTIVE", "id", instanceStatus.ID(), "status", instanceStatus.State())
		v1beta1conditions.MarkUnknown(openStackServer, infrav1.InstanceReadyCondition, infrav1.InstanceNotReadyReason, "Instance state is not handled: %s", instanceStatus.State())

		return ctrl.Result{RequeueAfter: waitForInstanceBecomeActiveToReconcile}, nil
	}

	scope.Logger().Info("Reconciled Server create successfully")
	return ctrl.Result{}, nil
}

// adoptServerResources adopts the OpenStack resources for the server.
func adoptServerResources(scope *scope.WithLogger, openStackServer *infrav1alpha1.OpenStackServer) error {
	resources := openStackServer.Status.Resources
	if resources == nil {
		resources = &infrav1alpha1.ServerResources{}
		openStackServer.Status.Resources = resources
	}

	// Adopt any existing resources
	return compute.AdoptServerResources(scope, openStackServer.Status.Resolved, resources)
}

func getOrCreateServerPorts(openStackServer *infrav1alpha1.OpenStackServer, networkingService *networking.Service) error {
	resolved := openStackServer.Status.Resolved
	if resolved == nil {
		return errors.New("server status resolved is nil")
	}
	resources := openStackServer.Status.Resources
	if resources == nil {
		return errors.New("server status resources is nil")
	}
	desiredPorts := resolved.Ports

	if err := networkingService.EnsurePorts(openStackServer, desiredPorts, resources); err != nil {
		return fmt.Errorf("creating ports: %w", err)
	}

	return nil
}

// getOrCreateServer gets or creates a server instance and returns the instance status, or an error.
func (r *OpenStackServerReconciler) getOrCreateServer(ctx context.Context, logger logr.Logger, openStackServer *infrav1alpha1.OpenStackServer, computeService *compute.Service, portIDs []string) (*compute.InstanceStatus, error) {
	var instanceStatus *compute.InstanceStatus
	var err error

	if openStackServer.Status.InstanceID != nil {
		instanceStatus, err = computeService.GetInstanceStatus(*openStackServer.Status.InstanceID)
		if err != nil || instanceStatus == nil {
			logger.Info("Unable to get OpenStack instance", "name", openStackServer.Name, "id", *openStackServer.Status.InstanceID)
			var msg string
			var reason string
			if err != nil {
				msg = err.Error()
				reason = infrav1.OpenStackErrorReason
			} else {
				msg = infrav1.ServerUnexpectedDeletedMessage
				reason = infrav1.InstanceNotFoundReason
			}
			v1beta1conditions.MarkFalse(openStackServer, infrav1.InstanceReadyCondition, reason, clusterv1beta1.ConditionSeverityError, "%s", msg)
			return nil, err
		}
	}
	if instanceStatus == nil {
		// Check if there is an existing instance with machine name, in case where instance ID would not have been stored in machine status
		instanceStatus, err := computeService.GetInstanceStatusByName(openStackServer, openStackServer.Name)
		if err != nil {
			logger.Error(err, "Failed to get instance by name", "name", openStackServer.Name)
			return nil, err
		}
		if instanceStatus != nil {
			logger.Info("Server already exists", "name", openStackServer.Name, "id", instanceStatus.ID())
			return instanceStatus, nil
		}

		logger.Info("Server does not exist, creating Server", "name", openStackServer.Name)
		instanceSpec, err := r.serverToInstanceSpec(ctx, openStackServer)
		if err != nil {
			return nil, err
		}
		instanceSpec.Name = openStackServer.Name
		instanceStatus, err = computeService.CreateInstance(openStackServer, instanceSpec, portIDs)
		if err != nil {
			v1beta1conditions.MarkFalse(openStackServer, infrav1.InstanceReadyCondition, infrav1.InstanceCreateFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
			openStackServer.Status.InstanceState = &infrav1.InstanceStateError
			return nil, fmt.Errorf("create OpenStack instance: %w", err)
		}
		// We reached a point where a server was created with no error but we can't predict its state yet which is why we don't update conditions yet.
		// The actual state of the server is checked in the next reconcile loops.
		return instanceStatus, nil
	}
	return instanceStatus, nil
}

func (r *OpenStackServerReconciler) getUserDataSecretValue(ctx context.Context, namespace, secretName string) (string, error) {
	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: namespace, Name: secretName}
	if err := r.Client.Get(ctx, key, secret); err != nil {
		return "", fmt.Errorf("failed to get secret %s/%s: %w", namespace, secretName, err)
	}

	value, ok := secret.Data["value"]
	if !ok {
		return "", fmt.Errorf("secret %s/%s does not contain userData", namespace, secretName)
	}

	return base64.StdEncoding.EncodeToString(value), nil
}

func (r *OpenStackServerReconciler) serverToInstanceSpec(ctx context.Context, openStackServer *infrav1alpha1.OpenStackServer) (*compute.InstanceSpec, error) {
	resolved := openStackServer.Status.Resolved
	if resolved == nil {
		return nil, errors.New("server resolved is nil")
	}

	serverMetadata := make(map[string]string, len(openStackServer.Spec.ServerMetadata))
	for i := range openStackServer.Spec.ServerMetadata {
		key := openStackServer.Spec.ServerMetadata[i].Key
		value := openStackServer.Spec.ServerMetadata[i].Value
		serverMetadata[key] = value
	}

	instanceSpec := &compute.InstanceSpec{
		AdditionalBlockDevices:        openStackServer.Spec.AdditionalBlockDevices,
		ConfigDrive:                   openStackServer.Spec.ConfigDrive != nil && *openStackServer.Spec.ConfigDrive,
		FlavorID:                      resolved.FlavorID,
		ImageID:                       resolved.ImageID,
		Metadata:                      serverMetadata,
		Name:                          openStackServer.Name,
		RootVolume:                    openStackServer.Spec.RootVolume,
		SSHKeyName:                    openStackServer.Spec.SSHKeyName,
		ServerGroupID:                 resolved.ServerGroupID,
		Tags:                          openStackServer.Spec.Tags,
		Trunk:                         openStackServer.Spec.Trunk != nil && *openStackServer.Spec.Trunk,
		SchedulerAdditionalProperties: openStackServer.Spec.SchedulerHintAdditionalProperties,
	}

	if openStackServer.Spec.UserDataRef != nil {
		userData, err := r.getUserDataSecretValue(ctx, openStackServer.Namespace, openStackServer.Spec.UserDataRef.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get user data secret value")
		}
		instanceSpec.UserData = userData
	}

	if openStackServer.Spec.AvailabilityZone != nil {
		instanceSpec.FailureDomain = *openStackServer.Spec.AvailabilityZone
	}

	return instanceSpec, nil
}

func getServerStatus(openStackServer *infrav1alpha1.OpenStackServer, computeService *compute.Service) (*compute.InstanceStatus, error) {
	if openStackServer.Status.InstanceID != nil {
		return computeService.GetInstanceStatus(*openStackServer.Status.InstanceID)
	}
	return computeService.GetInstanceStatusByName(openStackServer, openStackServer.Name)
}

// getClusterFromMetadata returns the Cluster object (if present) using the object metadata.
// This function was copied from the cluster-api project but manages errors differently.
func getClusterFromMetadata(ctx context.Context, c client.Client, obj metav1.ObjectMeta) (*clusterv1.Cluster, error) {
	// If the object is unlabeled, return early with no error.
	// It's fine for this object to not be part of a cluster.
	if obj.Labels[clusterv1.ClusterNameLabel] == "" {
		return nil, nil
	}
	// At this point, the object has a cluster name label so we should be able to find the cluster
	// and return an error if we can't.
	return util.GetClusterByName(ctx, c, obj.Namespace, obj.Labels[clusterv1.ClusterNameLabel])
}

// reconcileFloatingAddressFromPool reconciles the floating IP address from the pool.
// It returns the IPAddressClaim and a boolean indicating if the IPAddressClaim is ready.
func (r *OpenStackServerReconciler) reconcileFloatingAddressFromPool(ctx context.Context, scope *scope.WithLogger, openStackServer *infrav1alpha1.OpenStackServer) (*ipamv1.IPAddressClaim, bool, error) {
	if openStackServer.Spec.FloatingIPPoolRef == nil {
		return nil, false, nil
	}
	var claim *ipamv1.IPAddressClaim
	claim, err := r.getOrCreateIPAddressClaimForFloatingAddress(ctx, scope, openStackServer)
	if err != nil {
		v1beta1conditions.MarkFalse(openStackServer, infrav1.FloatingAddressFromPoolReadyCondition, infrav1.FloatingAddressFromPoolErrorReason, clusterv1beta1.ConditionSeverityInfo, "Failed to reconcile floating IP claims: %v", err)
		return nil, true, err
	}
	if claim.Status.AddressRef.Name == "" {
		r.Recorder.Eventf(openStackServer, corev1.EventTypeNormal, "WaitingForIPAddressClaim", "Waiting for IPAddressClaim %s/%s to be allocated", claim.Namespace, claim.Name)
		return claim, true, nil
	}
	v1beta1conditions.MarkTrue(openStackServer, infrav1.FloatingAddressFromPoolReadyCondition)
	return claim, false, nil
}

// createIPAddressClaim creates IPAddressClaim for the FloatingAddressFromPool if it does not exist yet.
func (r *OpenStackServerReconciler) getOrCreateIPAddressClaimForFloatingAddress(ctx context.Context, scope *scope.WithLogger, openStackServer *infrav1alpha1.OpenStackServer) (*ipamv1.IPAddressClaim, error) {
	var err error

	poolRef := openStackServer.Spec.FloatingIPPoolRef
	claimName := names.GetFloatingAddressClaimName(openStackServer.Name)
	claim := &ipamv1.IPAddressClaim{}

	err = r.Client.Get(ctx, client.ObjectKey{Namespace: openStackServer.Namespace, Name: claimName}, claim)
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

	// If the OpenStackServer has a ClusterNameLabel, set it on the IPAddressClaim as well.
	// This is useful for garbage collection of IPAddressClaims when a Cluster is deleted.
	if openStackServer.Labels[clusterv1.ClusterNameLabel] != "" {
		claim.Labels[clusterv1.ClusterNameLabel] = openStackServer.Labels[clusterv1.ClusterNameLabel]
	}

	if err := r.Client.Create(ctx, claim); err != nil {
		return nil, err
	}

	r.Recorder.Eventf(openStackServer, corev1.EventTypeNormal, "CreatingIPAddressClaim", "Creating IPAddressClaim %s/%s", claim.Namespace, claim.Name)
	scope.Logger().Info("Created IPAddressClaim", "name", claim.Name)
	return claim, nil
}

func (r *OpenStackServerReconciler) associateIPAddressFromIPAddressClaim(ctx context.Context, openStackServer *infrav1alpha1.OpenStackServer, instanceStatus *compute.InstanceStatus, instanceNS *compute.InstanceNetworkStatus, claim *ipamv1.IPAddressClaim, networkingService *networking.Service) error {
	address := &ipamv1.IPAddress{}
	addressKey := client.ObjectKey{Namespace: openStackServer.Namespace, Name: claim.Status.AddressRef.Name}

	if err := r.Client.Get(ctx, addressKey, address); err != nil {
		return err
	}

	instanceAddresses := instanceNS.Addresses()
	for _, instanceAddress := range instanceAddresses {
		if instanceAddress.Address == address.Spec.Address {
			v1beta1conditions.MarkTrue(openStackServer, infrav1.FloatingAddressFromPoolReadyCondition)
			return nil
		}
	}

	fip, err := networkingService.GetFloatingIP(address.Spec.Address)
	if err != nil {
		return err
	}

	if fip == nil {
		v1beta1conditions.MarkFalse(openStackServer, infrav1.FloatingAddressFromPoolReadyCondition, infrav1.FloatingAddressFromPoolErrorReason, clusterv1beta1.ConditionSeverityError, "floating IP does not exist")
		return fmt.Errorf("floating IP %q does not exist", address.Spec.Address)
	}

	port, err := networkingService.GetPortForExternalNetwork(instanceStatus.ID(), fip.FloatingNetworkID)
	if err != nil {
		return fmt.Errorf("get port for floating IP %q: %w", fip.FloatingIP, err)
	}

	if port == nil {
		v1beta1conditions.MarkFalse(openStackServer, infrav1.FloatingAddressFromPoolReadyCondition, infrav1.FloatingAddressFromPoolErrorReason, clusterv1beta1.ConditionSeverityError, "Can't find port for floating IP %q on external network %s", fip.FloatingIP, fip.FloatingNetworkID)
		return fmt.Errorf("port for floating IP %q on network %s does not exist", fip.FloatingIP, fip.FloatingNetworkID)
	}

	if err = networkingService.AssociateFloatingIP(openStackServer, fip, port.ID); err != nil {
		return err
	}
	v1beta1conditions.MarkTrue(openStackServer, infrav1.FloatingAddressFromPoolReadyCondition)
	return nil
}

func (r *OpenStackServerReconciler) reconcileDeleteFloatingAddressFromPool(scope *scope.WithLogger, openStackServer *infrav1alpha1.OpenStackServer) error {
	log := scope.Logger().WithValues("openStackMachine", openStackServer.Name)
	log.Info("Reconciling Machine delete floating address from pool")
	if openStackServer.Spec.FloatingIPPoolRef == nil {
		return nil
	}
	claimName := names.GetFloatingAddressClaimName(openStackServer.Name)
	claim := &ipamv1.IPAddressClaim{}
	if err := r.Client.Get(context.Background(), client.ObjectKey{Namespace: openStackServer.Namespace, Name: claimName}, claim); err != nil {
		return client.IgnoreNotFound(err)
	}

	controllerutil.RemoveFinalizer(claim, infrav1.IPClaimMachineFinalizer)
	return r.Client.Update(context.Background(), claim)
}

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
				log.V(6).Info("OpenStackServer finished reconciling, allowing further processing")
				return true
			}
			log.V(6).Info("OpenStackServer is still reconciling, blocking further processing")
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
				log.V(6).Info("OpenStackServer finished reconciling, allowing further processing")
				return true
			}

			log.V(4).Info("OpenStackServer is still reconciling, blocking further processing")
			return false
		},
		DeleteFunc:  func(event.DeleteEvent) bool { return false },
		GenericFunc: func(event.GenericEvent) bool { return false },
	}
}
