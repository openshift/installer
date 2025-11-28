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

package subnet

import (
	"context"
	"fmt"
	"iter"
	"slices"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	"github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/tags"
)

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=subnets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=subnets/status,verbs=get;update;patch

// OpenStack resource types
type (
	osResourceT = subnets.Subnet

	createResourceActuator    = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator    = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	reconcileResourceActuator = interfaces.ReconcileResourceActuator[orcObjectPT, osResourceT]
	resourceReconciler        = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory             = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
)

type subnetActuator struct {
	osClient  osclients.NetworkClient
	k8sClient client.Client
}

var _ createResourceActuator = subnetActuator{}
var _ deleteResourceActuator = subnetActuator{}

func (subnetActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator subnetActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	subnet, err := actuator.osClient.GetSubnet(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return subnet, nil
}

func (actuator subnetActuator) ListOSResourcesForAdoption(ctx context.Context, obj orcObjectPT) (iter.Seq2[*osResourceT, error], bool) {
	if obj.Spec.Resource == nil {
		return nil, false
	}
	listOpts := subnets.ListOpts{Name: getResourceName(obj)}
	return actuator.osClient.ListSubnet(ctx, listOpts), true
}

func (actuator subnetActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	var reconcileStatus progress.ReconcileStatus

	network := &orcv1alpha1.Network{}
	if filter.NetworkRef != "" {
		networkKey := client.ObjectKey{Name: string(filter.NetworkRef), Namespace: obj.Namespace}
		if err := actuator.k8sClient.Get(ctx, networkKey, network); err != nil {
			if apierrors.IsNotFound(err) {
				reconcileStatus = reconcileStatus.WithReconcileStatus(
					progress.WaitingOnObject("Network", networkKey.Name, progress.WaitingOnCreation))
			} else {
				reconcileStatus = reconcileStatus.WithReconcileStatus(
					progress.WrapError(fmt.Errorf("fetching network %s: %w", networkKey.Name, err)))
			}
		} else {
			if !orcv1alpha1.IsAvailable(network) || network.Status.ID == nil {
				reconcileStatus = reconcileStatus.WithReconcileStatus(
					progress.WaitingOnObject("Network", networkKey.Name, progress.WaitingOnReady))
			}
		}
	}

	project := &orcv1alpha1.Project{}
	if filter.ProjectRef != nil {
		projectKey := client.ObjectKey{Name: string(*filter.ProjectRef), Namespace: obj.Namespace}
		if err := actuator.k8sClient.Get(ctx, projectKey, project); err != nil {
			if apierrors.IsNotFound(err) {
				reconcileStatus = reconcileStatus.WithReconcileStatus(
					progress.WaitingOnObject("Project", projectKey.Name, progress.WaitingOnCreation))
			} else {
				reconcileStatus = reconcileStatus.WithReconcileStatus(
					progress.WrapError(fmt.Errorf("fetching project %s: %w", projectKey.Name, err)))
			}
		} else {
			if !orcv1alpha1.IsAvailable(project) || project.Status.ID == nil {
				reconcileStatus = reconcileStatus.WithReconcileStatus(
					progress.WaitingOnObject("Project", projectKey.Name, progress.WaitingOnReady))
			}
		}
	}

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	listOpts := subnets.ListOpts{
		Name:        string(ptr.Deref(filter.Name, "")),
		Description: string(ptr.Deref(filter.Description, "")),
		NetworkID:   ptr.Deref(network.Status.ID, ""),
		ProjectID:   ptr.Deref(project.Status.ID, ""),
		IPVersion:   int(ptr.Deref(filter.IPVersion, 0)),
		GatewayIP:   string(ptr.Deref(filter.GatewayIP, "")),
		CIDR:        string(ptr.Deref(filter.CIDR, "")),
		Tags:        tags.Join(filter.Tags),
		TagsAny:     tags.Join(filter.TagsAny),
		NotTags:     tags.Join(filter.NotTags),
		NotTagsAny:  tags.Join(filter.NotTagsAny),
	}
	if filter.IPv6 != nil {
		listOpts.IPv6AddressMode = string(ptr.Deref(filter.IPv6.AddressMode, ""))
		listOpts.IPv6RAMode = string(ptr.Deref(filter.IPv6.RAMode, ""))
	}

	return actuator.osClient.ListSubnet(ctx, listOpts), nil
}

func (actuator subnetActuator) CreateResource(ctx context.Context, obj orcObjectPT) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}

	network, reconcileStatus := networkDependency.GetDependency(
		ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.Network) bool {
			return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
		},
	)

	if resource.RouterRef != nil {
		_, routerDepRS := routerDependency.GetDependency(
			ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.Router) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(routerDepRS)
	}

	var projectID string
	if resource.ProjectRef != nil {
		project, projectDepRS := projectDependency.GetDependency(
			ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.Project) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(projectDepRS)
		if project != nil {
			projectID = ptr.Deref(project.Status.ID, "")
		}
	}

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	createOpts := subnets.CreateOpts{
		NetworkID:         *network.Status.ID,
		CIDR:              string(resource.CIDR),
		Name:              getResourceName(obj),
		Description:       string(ptr.Deref(resource.Description, "")),
		IPVersion:         gophercloud.IPVersion(resource.IPVersion),
		EnableDHCP:        resource.EnableDHCP,
		DNSPublishFixedIP: resource.DNSPublishFixedIP,
		ProjectID:         projectID,
	}

	if len(resource.AllocationPools) > 0 {
		createOpts.AllocationPools = make([]subnets.AllocationPool, len(resource.AllocationPools))
		for i := range resource.AllocationPools {
			createOpts.AllocationPools[i].Start = string(resource.AllocationPools[i].Start)
			createOpts.AllocationPools[i].End = string(resource.AllocationPools[i].End)
		}
	}

	if resource.Gateway != nil {
		switch resource.Gateway.Type {
		case orcv1alpha1.SubnetGatewayTypeAutomatic:
			// Nothing to do
		case orcv1alpha1.SubnetGatewayTypeNone:
			createOpts.GatewayIP = ptr.To("")
		case orcv1alpha1.SubnetGatewayTypeIP:
			fallthrough
		default:
			createOpts.GatewayIP = (*string)(resource.Gateway.IP)
		}
	}

	if len(resource.DNSNameservers) > 0 {
		createOpts.DNSNameservers = make([]string, len(resource.DNSNameservers))
		for i := range resource.DNSNameservers {
			createOpts.DNSNameservers[i] = string(resource.DNSNameservers[i])
		}
	}

	if len(resource.HostRoutes) > 0 {
		createOpts.HostRoutes = make([]subnets.HostRoute, len(resource.HostRoutes))
		for i := range resource.HostRoutes {
			createOpts.HostRoutes[i].DestinationCIDR = string(resource.HostRoutes[i].Destination)
			createOpts.HostRoutes[i].NextHop = string(resource.HostRoutes[i].NextHop)
		}
	}

	if resource.IPv6 != nil {
		createOpts.IPv6AddressMode = string(ptr.Deref(resource.IPv6.AddressMode, ""))
		createOpts.IPv6RAMode = string(ptr.Deref(resource.IPv6.RAMode, ""))
	}

	osResource, err := actuator.osClient.CreateSubnet(ctx, &createOpts)

	// We should require the spec to be updated before retrying a create which returned a conflict
	if orcerrors.IsConflict(err) {
		err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
	}

	if err != nil {
		return nil, progress.WrapError(err)
	}
	return osResource, nil
}

func (actuator subnetActuator) DeleteResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	// Delete any RouterInterface first, as this would prevent deletion of the subnet
	routerInterface, err := getRouterInterface(ctx, actuator.k8sClient, obj)
	if err != nil {
		return progress.WrapError(err)
	}

	if routerInterface != nil {
		// We will be reconciled again when it's gone
		if routerInterface.GetDeletionTimestamp().IsZero() {
			if err := actuator.k8sClient.Delete(ctx, routerInterface); err != nil {
				return progress.WrapError(err)
			}
		}
		return progress.WaitingOnObject("RouterInterface", routerInterface.GetName(), progress.WaitingOnDeletion)
	}

	return progress.WrapError(actuator.osClient.DeleteSubnet(ctx, osResource.ID))
}

func (actuator subnetActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	updateOpts := subnets.UpdateOpts{RevisionNumber: &osResource.RevisionNumber}

	handleNameUpdate(&updateOpts, obj, osResource)
	handleDescriptionUpdate(&updateOpts, resource, osResource)
	handleAllocationPoolsUpdate(&updateOpts, resource, osResource)
	handleHostRoutesUpdate(&updateOpts, resource, osResource)
	handleDNSNameserversUpdate(&updateOpts, resource, osResource)
	handleEnableDHCPUpdate(&updateOpts, resource, osResource)
	handleGatewayUpdate(&updateOpts, resource, osResource)

	// Note that we didn't make dnsPublishFixedIP mutable as it could constantly try to reconcile in some environments
	// as seen in https://github.com/k-orc/openstack-resource-controller/issues/189

	needsUpdate, err := needsUpdate(updateOpts)
	if err != nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err))
	}
	if !needsUpdate {
		log.V(logging.Debug).Info("No changes")
		return nil
	}

	_, err = actuator.osClient.UpdateSubnet(ctx, osResource.ID, updateOpts)

	// We should require the spec to be updated before retrying an update which returned a conflict
	if orcerrors.IsConflict(err) {
		err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
	}

	if err != nil {
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts subnets.UpdateOpts) (bool, error) {
	updateOptsMap, err := updateOpts.ToSubnetUpdateMap()
	if err != nil {
		return false, err
	}

	subnetUpdateMap, ok := updateOptsMap["subnet"].(map[string]any)
	if !ok {
		subnetUpdateMap = make(map[string]any)
	}

	// Revision number is not returned in the output of updateOpts.ToSubnetUpdateMap()
	// so nothing to drop here

	return len(subnetUpdateMap) > 0, nil
}

func handleNameUpdate(updateOpts *subnets.UpdateOpts, obj orcObjectPT, osResource *osResourceT) {
	name := getResourceName(obj)
	if osResource.Name != name {
		updateOpts.Name = &name
	}
}

func handleDescriptionUpdate(updateOpts *subnets.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	description := string(ptr.Deref(resource.Description, ""))
	if osResource.Description != description {
		updateOpts.Description = &description
	}
}

func handleAllocationPoolsUpdate(updateOpts *subnets.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	missingAllocationPool := false
	allocationPools := make([]subnets.AllocationPool, len(resource.AllocationPools))
	for i := range resource.AllocationPools {
		allocationPools[i].Start = string(resource.AllocationPools[i].Start)
		allocationPools[i].End = string(resource.AllocationPools[i].End)

		found := false
		for _, pool := range osResource.AllocationPools {
			if pool.Start == allocationPools[i].Start &&
				pool.End == allocationPools[i].End {
				found = true
				break
			}
		}
		if !found {
			missingAllocationPool = true
		}
	}

	extraAllocationPool := false
	for i := range osResource.AllocationPools {
		found := false
		for _, pool := range allocationPools {
			if pool.Start == osResource.AllocationPools[i].Start &&
				pool.End == osResource.AllocationPools[i].End {
				found = true
				break
			}
		}
		if !found {
			extraAllocationPool = true
		}
	}

	// If the spec doesn't set allocation pools, we'll get a default one and should not try to update it.
	if len(resource.AllocationPools) > 0 && (missingAllocationPool || extraAllocationPool) {
		updateOpts.AllocationPools = allocationPools
	}
}

func handleHostRoutesUpdate(updateOpts *subnets.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	missingHostRoute := false
	hostRoutes := make([]subnets.HostRoute, len(resource.HostRoutes))
	for i := range resource.HostRoutes {
		hostRoutes[i].DestinationCIDR = string(resource.HostRoutes[i].Destination)
		hostRoutes[i].NextHop = string(resource.HostRoutes[i].NextHop)

		found := false
		for _, route := range osResource.HostRoutes {
			if route.DestinationCIDR == hostRoutes[i].DestinationCIDR &&
				route.NextHop == hostRoutes[i].NextHop {
				found = true
				break
			}
		}
		if !found {
			missingHostRoute = true
		}
	}

	extraHostRoute := false
	for i := range osResource.HostRoutes {
		found := false
		for _, route := range hostRoutes {
			if route.DestinationCIDR == osResource.HostRoutes[i].DestinationCIDR &&
				route.NextHop == osResource.HostRoutes[i].NextHop {
				found = true
				break
			}
		}
		if !found {
			extraHostRoute = true
		}
	}

	if missingHostRoute || extraHostRoute {
		updateOpts.HostRoutes = &hostRoutes
	}
}

func handleDNSNameserversUpdate(updateOpts *subnets.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	nameservers := make([]string, len(resource.DNSNameservers))
	for i := range resource.DNSNameservers {
		nameservers[i] = string(resource.DNSNameservers[i])
	}

	// Let's not bother about potential duplicate entries: they will be rejected by neutron API
	if !slices.Equal(osResource.DNSNameservers, nameservers) {
		updateOpts.DNSNameservers = &nameservers
	}
}

func handleEnableDHCPUpdate(updateOpts *subnets.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	// Default is true
	enableDHCP := ptr.Deref(resource.EnableDHCP, true)
	if osResource.EnableDHCP != enableDHCP {
		updateOpts.EnableDHCP = &enableDHCP
	}
}

func handleGatewayUpdate(updateOpts *subnets.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	if resource.Gateway != nil {
		switch resource.Gateway.Type {
		case orcv1alpha1.SubnetGatewayTypeAutomatic:
			// Nothing to do
		case orcv1alpha1.SubnetGatewayTypeNone:
			if osResource.GatewayIP != "" {
				updateOpts.GatewayIP = ptr.To("")
			}
		default:
			if osResource.GatewayIP != string(ptr.Deref(resource.Gateway.IP, "")) {
				updateOpts.GatewayIP = (*string)(resource.Gateway.IP)
			}
		}
	}
}

var _ reconcileResourceActuator = subnetActuator{}

func (actuator subnetActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		tags.ReconcileTags[orcObjectPT, osResourceT](orcObject.Spec.Resource.Tags, osResource.Tags, tags.NewNeutronTagReplacer(actuator.osClient, "subnets", osResource.ID)),
		actuator.ensureRouterInterface,
		actuator.updateResource,
	}, nil
}

func (actuator subnetActuator) ensureRouterInterface(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	var err error

	routerInterface, err := getRouterInterface(ctx, actuator.k8sClient, orcObject)
	if err != nil {
		return progress.WrapError(err)
	}
	if routerInterfaceMatchesSpec(routerInterface, orcObject.Name, orcObject.Spec.Resource) {
		// Nothing to do
		return nil
	}

	// If it doesn't match we should delete any existing interface
	if routerInterface != nil {
		if routerInterface.GetDeletionTimestamp().IsZero() {
			if err := actuator.k8sClient.Delete(ctx, routerInterface); err != nil {
				return progress.WrapError(fmt.Errorf("deleting RouterInterface %s: %w", client.ObjectKeyFromObject(routerInterface), err))
			}
		}
		return progress.WaitingOnObject("routerinterface", routerInterface.Name, progress.WaitingOnDeletion)
	}

	// Otherwise create it
	routerInterface = &orcv1alpha1.RouterInterface{}
	routerInterface.Name = getRouterInterfaceName(orcObject)
	routerInterface.Namespace = orcObject.Namespace
	routerInterface.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion:         orcObject.APIVersion,
			Kind:               orcObject.Kind,
			Name:               orcObject.Name,
			UID:                orcObject.UID,
			BlockOwnerDeletion: ptr.To(true),
		},
	}
	routerInterface.Spec = orcv1alpha1.RouterInterfaceSpec{
		Type:      orcv1alpha1.RouterInterfaceTypeSubnet,
		RouterRef: *orcObject.Spec.Resource.RouterRef,
		SubnetRef: ptr.To(orcv1alpha1.KubernetesNameRef(orcObject.Name)),
	}

	if err := actuator.k8sClient.Create(ctx, routerInterface); err != nil {
		return progress.WrapError(fmt.Errorf("creating RouterInterface %s: %w", client.ObjectKeyFromObject(orcObject), err))
	}
	return progress.WaitingOnObject("routerinterface", routerInterface.Name, progress.WaitingOnReady)
}

func getRouterInterfaceName(orcObject *orcv1alpha1.Subnet) string {
	return orcObject.Name + "-subnet"
}

func routerInterfaceMatchesSpec(routerInterface *orcv1alpha1.RouterInterface, objectName string, resource *orcv1alpha1.SubnetResourceSpec) bool {
	// No routerRef -> there should be no routerInterface
	if resource.RouterRef == nil {
		return routerInterface == nil
	}

	// The router interface should:
	// * Exist
	// * Be of Subnet type
	// * Reference this subnet
	// * Reference the router in our spec

	if routerInterface == nil {
		return false
	}

	if routerInterface.Spec.Type != orcv1alpha1.RouterInterfaceTypeSubnet {
		return false
	}

	if string(ptr.Deref(routerInterface.Spec.SubnetRef, "")) != objectName {
		return false
	}

	return routerInterface.Spec.RouterRef == *resource.RouterRef
}

// getRouterInterface returns the router interface for this subnet, identified by its name
// returns nil for routerinterface without returning an error if the routerinterface does not exist
func getRouterInterface(ctx context.Context, k8sClient client.Client, orcObject *orcv1alpha1.Subnet) (*orcv1alpha1.RouterInterface, error) {
	routerInterface := &orcv1alpha1.RouterInterface{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: getRouterInterfaceName(orcObject), Namespace: orcObject.GetNamespace()}, routerInterface)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("fetching RouterInterface: %w", err)
	}

	return routerInterface, nil
}

type subnetHelperFactory struct{}

var _ helperFactory = subnetHelperFactory{}

func (subnetHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return subnetAdapter{obj}
}

func (subnetHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, controller, orcObject)
}

func (subnetHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, controller, orcObject)
}

func newActuator(ctx context.Context, controller interfaces.ResourceController, orcObject *orcv1alpha1.Subnet) (subnetActuator, progress.ReconcileStatus) {
	if orcObject == nil {
		return subnetActuator{}, progress.WrapError(fmt.Errorf("orcObject may not be nil"))
	}

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return subnetActuator{}, reconcileStatus
	}

	log := ctrl.LoggerFrom(ctx)
	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return subnetActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewNetworkClient()
	if err != nil {
		return subnetActuator{}, progress.WrapError(err)
	}

	return subnetActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}
