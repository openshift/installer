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

package port

import (
	"context"
	"fmt"
	"iter"
	"slices"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsbinding"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portstrustedvif"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	osclients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/tags"
)

type (
	osResourceT = osclients.PortExt

	createResourceActuator    = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator    = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	reconcileResourceActuator = interfaces.ReconcileResourceActuator[orcObjectPT, osResourceT]
	resourceReconciler        = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory             = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
	portIterator              = iter.Seq2[*osResourceT, error]
)

const (
	// The frequency to poll when waiting for a server to become active
	serverBuildPollingPeriod = 15 * time.Second
)

// resolveHostID resolves the actual host ID string to use for port binding.
// It handles both direct ID specification and server reference.
// Returns the resolved host ID and a reconcile status (for waiting on dependencies).
func resolveHostID(
	ctx context.Context,
	k8sClient client.Client,
	obj orcObjectPT,
	hostIDSpec *orcv1alpha1.HostID,
) (string, progress.ReconcileStatus) {
	if hostIDSpec == nil {
		return "", nil
	}

	// Direct ID specification
	if hostIDSpec.ID != "" {
		return hostIDSpec.ID, nil
	}

	// Server reference - fetch the server and extract its hostID
	if hostIDSpec.ServerRef != "" {
		server, serverDepRS := dependency.FetchDependency(
			ctx, k8sClient, obj.Namespace, &hostIDSpec.ServerRef, "Server",
			func(dep *orcv1alpha1.Server) bool {
				return orcv1alpha1.IsAvailable(dep) &&
					dep.Status.Resource != nil &&
					dep.Status.Resource.HostID != ""
			},
		)
		if needsReschedule, _ := serverDepRS.NeedsReschedule(); needsReschedule {
			return "", serverDepRS
		}
		if server != nil && server.Status.Resource != nil {
			return server.Status.Resource.HostID, nil
		}
	}

	return "", nil
}

type portActuator struct {
	osClient  osclients.NetworkClient
	k8sClient client.Client
}

var _ createResourceActuator = portActuator{}
var _ deleteResourceActuator = portActuator{}

func (portActuator) GetResourceID(osResource *osResourceT) string {
	return osResource.ID
}

func (actuator portActuator) GetOSResourceByID(ctx context.Context, id string) (*osResourceT, progress.ReconcileStatus) {
	port, err := actuator.osClient.GetPort(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return port, nil
}

func (actuator portActuator) ListOSResourcesForAdoption(ctx context.Context, obj *orcv1alpha1.Port) (portIterator, bool) {
	resource := obj.Spec.Resource
	if resource == nil {
		return nil, false
	}

	// Resolve the network ID from NetworkRef. Without the network ID,
	// adoption could match a port on the wrong network.
	network, rs := dependency.FetchDependency(
		ctx, actuator.k8sClient, obj.Namespace, &resource.NetworkRef, "Network",
		func(dep *orcv1alpha1.Network) bool {
			return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
		},
	)
	if needsReschedule, _ := rs.NeedsReschedule(); needsReschedule {
		return nil, false
	}

	// Resolve the project ID from ProjectRef if set.
	var projectID string
	if resource.ProjectRef != nil {
		project, rs := dependency.FetchDependency(
			ctx, actuator.k8sClient, obj.Namespace, resource.ProjectRef, "Project",
			func(dep *orcv1alpha1.Project) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := rs.NeedsReschedule(); needsReschedule {
			return nil, false
		}
		projectID = ptr.Deref(project.Status.ID, "")
	}

	listOpts := ports.ListOpts{
		Name:       getResourceName(obj),
		NetworkID:  ptr.Deref(network.Status.ID, ""),
		MACAddress: resource.MACAddress,
		ProjectID:  projectID,
	}
	return actuator.osClient.ListPort(ctx, listOpts), true
}

func (actuator portActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	var reconcileStatus progress.ReconcileStatus

	network, rs := dependency.FetchDependency[*orcv1alpha1.Network](
		ctx, actuator.k8sClient, obj.Namespace, &filter.NetworkRef, "Network",
		orcv1alpha1.IsAvailable,
	)
	reconcileStatus = reconcileStatus.WithReconcileStatus(rs)

	project, rs := dependency.FetchDependency[*orcv1alpha1.Project](
		ctx, actuator.k8sClient, obj.Namespace, filter.ProjectRef, "Project",
		orcv1alpha1.IsAvailable,
	)
	reconcileStatus = reconcileStatus.WithReconcileStatus(rs)

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	listOpts := ports.ListOpts{
		Name:         string(ptr.Deref(filter.Name, "")),
		Description:  string(ptr.Deref(filter.Description, "")),
		NetworkID:    ptr.Deref(network.Status.ID, ""),
		ProjectID:    ptr.Deref(project.Status.ID, ""),
		Tags:         tags.Join(filter.Tags),
		TagsAny:      tags.Join(filter.TagsAny),
		NotTags:      tags.Join(filter.NotTags),
		NotTagsAny:   tags.Join(filter.NotTagsAny),
		AdminStateUp: filter.AdminStateUp,
		MACAddress:   filter.MACAddress,
	}

	return actuator.osClient.ListPort(ctx, listOpts), nil
}

func (actuator portActuator) CreateResource(ctx context.Context, obj *orcv1alpha1.Port) (*osResourceT, progress.ReconcileStatus) {
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}

	// Fetch all dependencies and ensure they have our finalizer
	network, networkDepRS := networkDependency.GetDependency(
		ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
	)
	subnetMap, subnetDepRS := subnetDependency.GetDependencies(
		ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
	)
	secGroupMap, secGroupDepRS := securityGroupDependency.GetDependencies(
		ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
	)
	reconcileStatus := progress.NewReconcileStatus().
		WithReconcileStatus(networkDepRS).
		WithReconcileStatus(subnetDepRS).
		WithReconcileStatus(secGroupDepRS)

	var projectID string
	if resource.ProjectRef != nil {
		project, projectDepRS := projectDependency.GetDependency(
			ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
		)
		reconcileStatus = reconcileStatus.WithReconcileStatus(projectDepRS)
		if project != nil {
			projectID = ptr.Deref(project.Status.ID, "")
		}
	}

	// Resolve hostID if specified
	var resolvedHostID string
	if resource.HostID != nil {
		var hostIDReconcileStatus progress.ReconcileStatus
		resolvedHostID, hostIDReconcileStatus = resolveHostID(ctx, actuator.k8sClient, obj, resource.HostID)
		reconcileStatus = reconcileStatus.WithReconcileStatus(hostIDReconcileStatus)
	}

	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return nil, reconcileStatus
	}

	var valueSpecs *map[string]string
	if len(resource.ValueSpecs) > 0 {
		vs := make(map[string]string, len(resource.ValueSpecs))
		for _, valueSpec := range resource.ValueSpecs {
			vs[valueSpec.Key] = *valueSpec.Value
		}
		valueSpecs = &vs
	}

	createOpts := ports.CreateOpts{
		NetworkID:             *network.Status.ID,
		Name:                  getResourceName(obj),
		Description:           string(ptr.Deref(resource.Description, "")),
		ProjectID:             projectID,
		AdminStateUp:          resource.AdminStateUp,
		MACAddress:            resource.MACAddress,
		ValueSpecs:            valueSpecs,
		PropagateUplinkStatus: resource.PropagateUplinkStatus,
	}

	if len(resource.AllowedAddressPairs) > 0 {
		if resource.PortSecurity == orcv1alpha1.PortSecurityDisabled {
			return nil, progress.WrapError(
				orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "AllowedAddressPairs cannot be set when PortSecurity is disabled"))
		}
		createOpts.AllowedAddressPairs = make([]ports.AddressPair, len(resource.AllowedAddressPairs))
		for i := range resource.AllowedAddressPairs {
			createOpts.AllowedAddressPairs[i].IPAddress = string(resource.AllowedAddressPairs[i].IP)
			if resource.AllowedAddressPairs[i].MAC != nil {
				createOpts.AllowedAddressPairs[i].MACAddress = string(*resource.AllowedAddressPairs[i].MAC)
			}
		}
	}

	// We explicitly disable creation of IP addresses by passing an empty
	// value whenever the user does not specify addresses
	fixedIPs := make([]ports.IP, len(resource.Addresses))
	for i := range resource.Addresses {
		subnetName := string(resource.Addresses[i].SubnetRef)
		subnet, ok := subnetMap[subnetName]
		if !ok {
			// Programming error
			return nil, progress.WrapError(fmt.Errorf("subnet %s was not returned by GetDependencies", subnetName))
		}
		fixedIPs[i].SubnetID = *subnet.Status.ID

		if resource.Addresses[i].IP != nil {
			fixedIPs[i].IPAddress = string(*resource.Addresses[i].IP)
		}
	}
	createOpts.FixedIPs = fixedIPs

	// We explicitly disable default security groups by passing an empty
	// value whenever the user does not specifies security groups
	securityGroups := make([]string, len(resource.SecurityGroupRefs))
	if len(securityGroups) > 0 && resource.PortSecurity == orcv1alpha1.PortSecurityDisabled {
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "SecurityGroupRefs cannot be set when PortSecurity is disabled"))
	}
	for i := range resource.SecurityGroupRefs {
		secGroupName := string(resource.SecurityGroupRefs[i])
		secGroup, ok := secGroupMap[secGroupName]
		if !ok {
			// Programming error
			return nil, progress.WrapError(fmt.Errorf("security group %s was not returned by GetDependencies", secGroupName))
		}
		securityGroups[i] = *secGroup.Status.ID
	}
	createOpts.SecurityGroups = &securityGroups

	portsBindingOpts := portsbinding.CreateOptsExt{
		CreateOptsBuilder: createOpts,
		VNICType:          resource.VNICType,
		HostID:            resolvedHostID,
	}

	portSecurityOpts := portsecurity.PortCreateOptsExt{
		CreateOptsBuilder: portsBindingOpts,
	}
	switch resource.PortSecurity {
	case orcv1alpha1.PortSecurityEnabled:
		portSecurityOpts.PortSecurityEnabled = ptr.To(true)
	case orcv1alpha1.PortSecurityDisabled:
		portSecurityOpts.PortSecurityEnabled = ptr.To(false)
	case orcv1alpha1.PortSecurityInherit:
		// do nothing
	default:
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, fmt.Sprintf("Invalid value %s", resource.PortSecurity)))
	}

	portTrustedOpts := portstrustedvif.PortCreateOptsExt{
		CreateOptsBuilder: portSecurityOpts,
	}
	if resource.TrustedVIF != nil {
		portTrustedOpts.PortTrustedVIF = resource.TrustedVIF
	}

	osResource, err := actuator.osClient.CreatePort(ctx, &portTrustedOpts)
	if err != nil {
		// We should require the spec to be updated before retrying a create which returned a non-retryable error
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
		}
		return nil, progress.WrapError(err)
	}

	return osResource, nil
}

func (actuator portActuator) DeleteResource(ctx context.Context, _ *orcv1alpha1.Port, osResource *osResourceT) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeletePort(ctx, osResource.ID))
}

var _ reconcileResourceActuator = portActuator{}

func (actuator portActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		actuator.checkAttachedServer,
		tags.ReconcileTags[orcObjectPT, osResourceT](orcObject.Spec.Resource.Tags, osResource.Tags, tags.NewNeutronTagReplacer(actuator.osClient, "ports", osResource.ID)),
		actuator.updateResource,
	}, nil
}

func (actuator portActuator) checkAttachedServer(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)

	// If the port is attached to a device, check if it's a server in BUILD status
	if osResource.DeviceID == "" {
		return nil
	}

	// List all servers in the namespace to find the one with matching ID
	serverList := &orcv1alpha1.ServerList{}
	if err := actuator.k8sClient.List(ctx, serverList, client.InNamespace(obj.Namespace)); err != nil {
		log.Error(err, "failed to list servers", "namespace", obj.Namespace)
		// Don't block port reconciliation if we can't list servers
		return nil
	}

	// Find server with matching ID
	for i := range serverList.Items {
		server := &serverList.Items[i]
		if server.Status.ID != nil && *server.Status.ID == osResource.DeviceID {
			// Check if server is in BUILD status
			if server.Status.Resource != nil && server.Status.Resource.Status == "BUILD" {
				log.V(logging.Verbose).Info("Port is attached to server in BUILD status, waiting",
					"port", obj.Name,
					"server", server.Name,
					"serverStatus", server.Status.Resource.Status)
				return progress.NewReconcileStatus().WaitingOnOpenStack(progress.WaitingOnReady, serverBuildPollingPeriod)
			}
			// Server found and not in BUILD status, continue reconciliation
			break
		}
	}

	return nil
}

func (actuator portActuator) updateResource(ctx context.Context, obj orcObjectPT, osResource *osResourceT) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Update requested, but spec.resource is not set"))
	}

	secGroupMap, secGroupDepRS := securityGroupDependency.GetDependencies(
		ctx, actuator.k8sClient, obj, orcv1alpha1.IsAvailable,
	)

	reconcileStatus := progress.NewReconcileStatus().
		WithReconcileStatus(secGroupDepRS)

	needsReschedule, _ := reconcileStatus.NeedsReschedule()
	if needsReschedule {
		return reconcileStatus
	}

	var updateOpts ports.UpdateOptsBuilder
	{
		baseUpdateOpts := &ports.UpdateOpts{
			RevisionNumber: &osResource.RevisionNumber,
		}
		handleNameUpdate(baseUpdateOpts, obj, osResource)
		handleDescriptionUpdate(baseUpdateOpts, resource, osResource)
		handleAllowedAddressPairsUpdate(baseUpdateOpts, resource, osResource)
		handleSecurityGroupRefsUpdate(baseUpdateOpts, resource, osResource, secGroupMap)
		handleAdminStateUpUpdate(baseUpdateOpts, resource, osResource)
		updateOpts = baseUpdateOpts
	}

	updateOpts = handlePortBindingUpdate(updateOpts, resource, osResource)
	updateOpts = handlePortSecurityUpdate(updateOpts, resource, osResource)
	updateOpts = handlePortTrustedVIFUpdate(updateOpts, resource, osResource)

	needsUpdate, err := needsUpdate(updateOpts)
	if err != nil {
		return progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err))
	}
	if !needsUpdate {
		log.V(logging.Debug).Info("No changes")
		return nil
	}

	_, err = actuator.osClient.UpdatePort(ctx, osResource.ID, updateOpts)

	if err != nil {
		if !orcerrors.IsRetryable(err) {
			err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration updating resource: "+err.Error(), err)
		}
		return progress.WrapError(err)
	}

	return progress.NeedsRefresh()
}

func needsUpdate(updateOpts ports.UpdateOptsBuilder) (bool, error) {
	updateOptsMap, err := updateOpts.ToPortUpdateMap()
	if err != nil {
		return false, err
	}

	portUpdateMap, ok := updateOptsMap["port"].(map[string]any)
	if !ok {
		portUpdateMap = make(map[string]any)
	}

	// Revision number is not returned in the output of updateOpts.ToPortUpdateMap()
	// so nothing to drop here

	return len(portUpdateMap) > 0, nil
}

func handleNameUpdate(updateOpts *ports.UpdateOpts, obj orcObjectPT, osResource *osResourceT) {
	name := getResourceName(obj)
	if osResource.Name != name {
		updateOpts.Name = &name
	}
}

func handleDescriptionUpdate(updateOpts *ports.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	description := string(ptr.Deref(resource.Description, ""))
	if osResource.Description != description {
		updateOpts.Description = &description
	}
}

func handleAllowedAddressPairsUpdate(updateOpts *ports.UpdateOpts, resource *orcv1alpha1.PortResourceSpec, osResource *osclients.PortExt) {
	desiredPairs := make([]ports.AddressPair, len(resource.AllowedAddressPairs))
	for i, pair := range resource.AllowedAddressPairs {
		desiredPairs[i].IPAddress = string(pair.IP)

		// The MAC address is optional. If it's nil in the spec, it will be an empty string
		// in the OpenStack API struct, which is the correct representation.
		if pair.MAC != nil {
			desiredPairs[i].MACAddress = string(*pair.MAC)
		}
	}

	missingPair := false
	for _, desired := range desiredPairs {
		found := false
		for _, actual := range osResource.AllowedAddressPairs {
			if actual.IPAddress == desired.IPAddress && actual.MACAddress == desired.MACAddress {
				found = true
				break
			}
		}
		if !found {
			missingPair = true
			break
		}
	}

	extraPair := false
	for _, actual := range osResource.AllowedAddressPairs {
		found := false
		for _, desired := range desiredPairs {
			if actual.IPAddress == desired.IPAddress && actual.MACAddress == desired.MACAddress {
				found = true
				break
			}
		}
		if !found {
			extraPair = true
			break
		}
	}

	if missingPair || extraPair {
		updateOpts.AllowedAddressPairs = &desiredPairs
	}
}

func handleSecurityGroupRefsUpdate(updateOpts *ports.UpdateOpts, resource *resourceSpecT, osResource *osResourceT, secGroupMap map[string]*orcv1alpha1.SecurityGroup) {
	// Translate desired names → IDs
	desiredIDs := make([]string, len(resource.SecurityGroupRefs))
	for i, refName := range resource.SecurityGroupRefs {
		sg, ok := secGroupMap[string(refName)]
		if !ok || sg.Status.ID == nil {
			continue
		}
		desiredIDs[i] = *sg.Status.ID
	}
	currentIDs := make([]string, len(osResource.SecurityGroups))
	copy(currentIDs, osResource.SecurityGroups)

	slices.Sort(desiredIDs)
	slices.Sort(currentIDs)

	if !slices.Equal(desiredIDs, currentIDs) {
		updateOpts.SecurityGroups = &desiredIDs
	}
}

func handlePortBindingUpdate(updateOpts ports.UpdateOptsBuilder, resource *resourceSpecT, osResource *osResourceT) ports.UpdateOptsBuilder {
	if resource.VNICType != "" {
		if resource.VNICType != osResource.VNICType {
			updateOpts = &portsbinding.UpdateOptsExt{
				UpdateOptsBuilder: updateOpts,
				VNICType:          resource.VNICType,
			}
		}
	}

	return updateOpts
}

func handlePortSecurityUpdate(updateOpts ports.UpdateOptsBuilder, resource *resourceSpecT, osResource *osResourceT) ports.UpdateOptsBuilder {

	var desiredState *bool

	switch resource.PortSecurity {
	case orcv1alpha1.PortSecurityInherit:
		return updateOpts
	case orcv1alpha1.PortSecurityEnabled:
		desiredState = ptr.To(true)
	case orcv1alpha1.PortSecurityDisabled:
		desiredState = ptr.To(false)
	default:
		return updateOpts
	}

	if *desiredState != osResource.PortSecurityEnabled {
		updateOpts = &portsecurity.PortUpdateOptsExt{
			UpdateOptsBuilder:   updateOpts,
			PortSecurityEnabled: desiredState,
		}
	}

	return updateOpts
}

func handleAdminStateUpUpdate(updateOpts *ports.UpdateOpts, resource *resourceSpecT, osResource *osResourceT) {
	adminStateUp := resource.AdminStateUp
	if adminStateUp != nil {
		if *adminStateUp != osResource.AdminStateUp {
			updateOpts.AdminStateUp = adminStateUp
		}
	}
}

func handlePortTrustedVIFUpdate(updateOpts ports.UpdateOptsBuilder, resource *resourceSpecT, osResource *osResourceT) ports.UpdateOptsBuilder {
	trusted := resource.TrustedVIF
	if trusted != nil {
		if osResource.PortTrustedVIF == nil || *trusted != *osResource.PortTrustedVIF {
			updateOpts = portstrustedvif.PortUpdateOptsExt{
				UpdateOptsBuilder: updateOpts,
				PortTrustedVIF:    trusted,
			}
		}
	}

	return updateOpts
}

type portHelperFactory struct{}

var _ helperFactory = portHelperFactory{}

func (portHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return portAdapter{obj}
}

func (portHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, controller, orcObject)
}

func (portHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, controller, orcObject)
}

func newActuator(ctx context.Context, controller interfaces.ResourceController, orcObject *orcv1alpha1.Port) (portActuator, progress.ReconcileStatus) {
	if orcObject == nil {
		return portActuator{}, progress.WrapError(fmt.Errorf("orcObject may not be nil"))
	}

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return portActuator{}, reconcileStatus
	}

	log := ctrl.LoggerFrom(ctx)
	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return portActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewNetworkClient()
	if err != nil {
		return portActuator{}, progress.WrapError(err)
	}

	return portActuator{
		osClient:  osClient,
		k8sClient: controller.GetK8sClient(),
	}, nil
}
