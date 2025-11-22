/*
Copyright 2024 The ORC Authors.

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

package routerinterface

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	osclients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/finalizers"
	orcstrings "github.com/k-orc/openstack-resource-controller/v2/internal/util/strings"
)

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=routerinterfaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=routerinterfaces/status,verbs=get;update;patch

func (r *orcRouterInterfaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	router := &orcv1alpha1.Router{}
	if err := r.client.Get(ctx, req.NamespacedName, router); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// We still want to reconcile router interface for routers that are scheduled for deletion
	// so that we can clear the finalizer
	if router.Status.ID == nil || (!orcv1alpha1.IsAvailable(router) && router.GetDeletionTimestamp().IsZero()) {
		log.V(logging.Verbose).Info("Not reconciling interfaces for not-Available router")
		return ctrl.Result{}, nil
	}

	routerInterfaces, err := routerDependency.GetObjectsForDependency(ctx, r.client, router)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("fetching router interfaces: %w", err)
	}

	// We don't need to query neutron for ports if there are no interfaces to reconcile
	if len(routerInterfaces) == 0 {
		return ctrl.Result{}, nil
	}

	// If there are interfaces, the router should have our finalizer
	if err := dependency.EnsureFinalizer(ctx, r.client, router, finalizer, fieldOwner); err != nil {
		return ctrl.Result{}, fmt.Errorf("writing finalizer: %w", err)
	}

	listOpts := ports.ListOpts{
		DeviceID: *router.Status.ID,
	}

	networkClient, err := r.getNetworkClient(ctx, router)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("getting network client: %w", err)
	}

	routerInterfacePortIterator := networkClient.ListPort(ctx, &listOpts)
	// We're going to iterate over all interfaces multiple times, so pull them all in to a slice
	var routerInterfacePorts []osclients.PortExt
	for port, err := range routerInterfacePortIterator {
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("fetching router interface ports: %w", err)
		}
		switch port.DeviceOwner {
		case "network:router_interface", "network:ha_router_replicated_interface", "network:router_interface_distributed":
			routerInterfacePorts = append(routerInterfacePorts, *port)
		default:
			log.V(logging.Debug).Info("ignoring port with unexpected device owner",
				"deviceID", *router.Status.ID,
				"deviceOwner", port.DeviceOwner,
				"portID", port.ID)
			continue
		}
	}

	var reconcileStatus progress.ReconcileStatus
	for i := range routerInterfaces {
		routerInterface := &routerInterfaces[i]
		log = log.WithValues("name", routerInterface.Name)

		var ifReconcileStatus progress.ReconcileStatus
		if routerInterface.GetDeletionTimestamp().IsZero() {
			ifReconcileStatus = r.reconcileNormal(ctx, log, router, routerInterface, routerInterfacePorts, networkClient)
		} else {
			ifReconcileStatus = r.reconcileDelete(ctx, log, router, routerInterface, routerInterfacePorts, networkClient)
		}

		// Don't aggregate terminal errors because we don't return them to controller runtime
		err := ifReconcileStatus.GetError()
		var terminalError *orcerrors.TerminalError
		if !errors.As(err, &terminalError) {
			reconcileStatus = reconcileStatus.WithReconcileStatus(ifReconcileStatus)
		}

	}
	return reconcileStatus.Return(log)
}

func (r *orcRouterInterfaceReconciler) getNetworkClient(ctx context.Context, obj orcv1alpha1.CloudCredentialsRefProvider) (osclients.NetworkClient, error) {
	log := ctrl.LoggerFrom(ctx)

	clientScope, err := r.scopeFactory.NewClientScopeFromObject(ctx, r.client, log, obj)
	if err != nil {
		return nil, err
	}
	return clientScope.NewNetworkClient()
}

func (r *orcRouterInterfaceReconciler) reconcileNormal(ctx context.Context, log logr.Logger, router *orcv1alpha1.Router, routerInterface *orcv1alpha1.RouterInterface, routerInterfacePorts []osclients.PortExt, networkClient osclients.NetworkClient) (reconcileStatus progress.ReconcileStatus) {
	log.V(logging.Verbose).Info("Reconciling router interface")

	var osResource *osclients.PortExt

	// Ensure we always update status
	defer func() {
		reconcileStatus = reconcileStatus.WithReconcileStatus(
			r.updateStatus(ctx, routerInterface, osResource, reconcileStatus))

		if needsReschedule, _ := reconcileStatus.NeedsReschedule(); !needsReschedule && osResource != nil {
			log.V(logging.Status).Info("Router interface is available")
		}

		// Don't return a terminal error because we don't aggregate them
		err := reconcileStatus.GetError()
		var terminalError *orcerrors.TerminalError
		if errors.As(err, &terminalError) {
			reconcileStatus = nil
		}
	}()

	switch routerInterface.Spec.Type {
	case orcv1alpha1.RouterInterfaceTypeSubnet:
		osResource, reconcileStatus = r.reconcileNormalSubnet(ctx, log, router, routerInterface, routerInterfacePorts, networkClient)

	default:
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, fmt.Sprintf("Invalid type %s", routerInterface.Spec.Type)))
	}
	return reconcileStatus
}

func (r *orcRouterInterfaceReconciler) createRouterInterface(ctx context.Context, log logr.Logger, router *orcv1alpha1.Router, routerInterface *orcv1alpha1.RouterInterface, createOpts routers.AddInterfaceOptsBuilder, networkClient osclients.NetworkClient) progress.ReconcileStatus {
	// Add finalizer immediately before creating a resource
	// Adding the finalizer only when creating a resource means we don't add
	// it until all dependent resources are available, which means we don't
	// have to handle unavailable dependencies in the delete flow
	if err := dependency.EnsureFinalizer(ctx, r.client, routerInterface, finalizer, orcstrings.GetSSAFieldOwnerWithTxn(controllerName, orcstrings.SSATransactionFinalizer)); err != nil {
		return progress.WrapError(
			fmt.Errorf("setting finalizer for %s: %w", client.ObjectKeyFromObject(routerInterface), err))
	}

	info, err := networkClient.AddRouterInterface(ctx, *router.Status.ID, createOpts)
	if err != nil {
		return progress.WrapError(
			fmt.Errorf("adding router interface: %w", err))
	}
	log.V(logging.Debug).Info("added router interface", "id", info.ID, "portID", info.PortID)

	// We're going to have to poll the interface port anyway, so rather than fetching it here we just schedule the next poll and we'll fetch it next time
	return progress.WaitingOnOpenStack(progress.WaitingOnReady, portStatusPollingPeriod)
}

func findPortBySubnetID(routerInterfacePorts []osclients.PortExt, subnetID string) *osclients.PortExt {
	for i := range routerInterfacePorts {
		routerInterfacePort := &routerInterfacePorts[i]

		for j := range routerInterfacePort.FixedIPs {
			fixedIP := &routerInterfacePort.FixedIPs[j]
			if fixedIP.SubnetID == subnetID {
				return routerInterfacePort
			}
		}
	}
	return nil
}

func (r *orcRouterInterfaceReconciler) reconcileNormalSubnet(ctx context.Context, log logr.Logger, router *orcv1alpha1.Router, routerInterface *orcv1alpha1.RouterInterface, routerInterfacePorts []osclients.PortExt, networkClient osclients.NetworkClient) (*osclients.PortExt, progress.ReconcileStatus) {
	subnet := &orcv1alpha1.Subnet{}
	subnetKey := client.ObjectKey{
		Namespace: routerInterface.Namespace,
		Name:      string(ptr.Deref(routerInterface.Spec.SubnetRef, "")),
	}
	if err := r.client.Get(ctx, subnetKey, subnet); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, progress.WaitingOnObject("Subnet", subnetKey.Name, progress.WaitingOnCreation)
		}
		return nil, progress.WrapError(fmt.Errorf("fetching subnet %s: %w", subnetKey, err))
	}
	log = log.WithValues("subnet", subnet.Name)

	// Don't reconcile for a subnet which has been deleted
	if !subnet.GetDeletionTimestamp().IsZero() {
		log.V(logging.Verbose).Info("Not reconciling interface for deleted subnet")
		return nil, progress.WaitingOnObject("Subnet", subnetKey.Name, progress.WaitingOnReady)
	}

	// Ensure the dependent subnet has our finalizer
	if err := dependency.EnsureFinalizer(ctx, r.client, subnet, finalizer, fieldOwner); err != nil {
		return nil, progress.WrapError(fmt.Errorf("adding finalizer to subnet: %w", err))
	}

	if subnet.Status.ID == nil {
		// We don't wait on Available here, but the subnet won't be available until this interface is up
		return nil, progress.WaitingOnObject("Subnet", subnetKey.Name, progress.WaitingOnReady)
	}
	subnetID := *subnet.Status.ID
	log = log.WithValues("subnetID", subnetID)

	// Port already exists for this subnet
	port := findPortBySubnetID(routerInterfacePorts, subnetID)
	if port != nil {
		log.V(logging.Debug).Info("found existing port", "portID", port.ID)
		return port, nil
	}

	createOpts := &routers.AddInterfaceOpts{SubnetID: subnetID}
	return nil, r.createRouterInterface(ctx, log, router, routerInterface, createOpts, networkClient)
}

func (r *orcRouterInterfaceReconciler) reconcileDelete(ctx context.Context, log logr.Logger, router *orcv1alpha1.Router, routerInterface *orcv1alpha1.RouterInterface, routerInterfacePorts []osclients.PortExt, networkClient osclients.NetworkClient) (reconcileStatus progress.ReconcileStatus) {
	var foundFinalizer bool
	for _, f := range routerInterface.GetFinalizers() {
		if f == finalizer {
			foundFinalizer = true
		} else {
			reconcileStatus = reconcileStatus.WaitingOnFinalizer(f)
		}
	}

	// Cleanup not required if our finalizer is not present
	if !foundFinalizer {
		log.V(logging.Verbose).Info("Not reconciling delete of router interface without finalizer")
		return reconcileStatus
	}

	if needsReschedule, err := reconcileStatus.NeedsReschedule(); needsReschedule {
		if err == nil {
			log.V(logging.Verbose).Info("Deferring resource cleanup due to remaining external finalizers")
		}
		return reconcileStatus
	}

	log.V(logging.Verbose).Info("Reconciling router interface delete")

	var osResource *osclients.PortExt
	deleted := false
	defer func() {
		// No point updating status after removing the finalizer
		if !deleted {
			reconcileStatus = reconcileStatus.WithReconcileStatus(
				r.updateStatus(ctx, routerInterface, osResource, reconcileStatus))
		}
	}()

	var deleteOpts routers.RemoveInterfaceOptsBuilder
	switch routerInterface.Spec.Type {
	case orcv1alpha1.RouterInterfaceTypeSubnet:
		osResource, deleteOpts, reconcileStatus = r.reconcileDeleteSubnet(ctx, log, routerInterface, routerInterfacePorts)
	default:
		reconcileStatus = reconcileStatus.WithError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, fmt.Sprintf("Invalid type %s", routerInterface.Spec.Type)))
	}
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return reconcileStatus
	}

	if deleteOpts != nil {
		log.V(logging.Debug).Info("Deleting router interface")
		_, err := networkClient.RemoveRouterInterface(ctx, *router.Status.ID, deleteOpts)
		if err != nil {
			return progress.WrapError(
				fmt.Errorf("removing router interface: %w", err))
		}

		return progress.WaitingOnOpenStack(progress.WaitingOnDeletion, portStatusPollingPeriod)
	}

	// Clear the finalizer
	deleted = true
	log.V(logging.Info).Info("Router interface deleted")
	return progress.WrapError(
		r.client.Patch(ctx, routerInterface, finalizers.RemoveFinalizerPatch(routerInterface), client.ForceOwnership, orcstrings.GetSSAFieldOwnerWithTxn(controllerName, orcstrings.SSATransactionFinalizer)))
}

func (r *orcRouterInterfaceReconciler) reconcileDeleteSubnet(ctx context.Context, log logr.Logger, routerInterface *orcv1alpha1.RouterInterface, routerInterfacePorts []osclients.PortExt) (*osclients.PortExt, routers.RemoveInterfaceOptsBuilder, progress.ReconcileStatus) {
	subnet := &orcv1alpha1.Subnet{}
	subnetKey := client.ObjectKey{
		Namespace: routerInterface.Namespace,
		Name:      string(ptr.Deref(routerInterface.Spec.SubnetRef, "")),
	}
	if err := r.client.Get(ctx, subnetKey, subnet); err != nil {
		if apierrors.IsNotFound(err) {
			// This should not happen unless something external messed with our
			// finalizers. We can't continue in this case because we don't know
			// the subnet ID so we can't check if the interface has been
			// removed. We will be automatically reconciled again if the subnet
			// is recreated.
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonUnrecoverableError, "Subnet has been deleted unexpectedly")
		} else {
			err = fmt.Errorf("fetching subnet %s: %w", subnetKey, err)
		}

		return nil, nil, progress.WrapError(err)
	}
	if subnet.Status.ID == nil {
		// This is unrecoverable on delete. We shouldn't have added a finalizer
		// unless the the subnet was ready, so something is wrong.
		return nil, nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonUnrecoverableError, "Subnet ID is not set"))
	}
	subnetID := *subnet.Status.ID
	log = log.WithValues("subnetID", subnetID)
	log.V(logging.Debug).Info("Found subnet")

	routerInterfacePort := findPortBySubnetID(routerInterfacePorts, subnetID)
	if routerInterfacePort == nil {
		log.V(logging.Debug).Info("No port found deleting router interface. Assuming already deleted.")
		return nil, nil, nil
	}

	log.V(logging.Debug).Info("Will delete router interface", "portID", routerInterfacePort.ID)
	return routerInterfacePort, &routers.RemoveInterfaceOpts{SubnetID: subnetID}, nil
}
