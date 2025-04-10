/*
Copyright 2018 The Kubernetes Authors.

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

package scope

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/availabilitysets"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/disks"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/inboundnatrules"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/networkinterfaces"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/publicips"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/resourceskus"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/roleassignments"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/virtualmachineimages"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/virtualmachines"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/vmextensions"
	"sigs.k8s.io/cluster-api-provider-azure/feature"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/futures"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	Client       client.Client
	ClusterScope azure.ClusterScoper
	Machine      *clusterv1.Machine
	AzureMachine *infrav1.AzureMachine
	Cache        *MachineCache
	SKUCache     SKUCacher
}

// NewMachineScope creates a new MachineScope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	if params.Client == nil {
		return nil, errors.New("client is required when creating a MachineScope")
	}
	if params.Machine == nil {
		return nil, errors.New("machine is required when creating a MachineScope")
	}
	if params.AzureMachine == nil {
		return nil, errors.New("azure machine is required when creating a MachineScope")
	}

	helper, err := patch.NewHelper(params.AzureMachine, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &MachineScope{
		client:        params.Client,
		Machine:       params.Machine,
		AzureMachine:  params.AzureMachine,
		patchHelper:   helper,
		ClusterScoper: params.ClusterScope,
		cache:         params.Cache,
		skuCache:      params.SKUCache,
	}, nil
}

// MachineScope defines a scope defined around a machine and its cluster.
type MachineScope struct {
	client      client.Client
	patchHelper *patch.Helper

	azure.ClusterScoper
	Machine      *clusterv1.Machine
	AzureMachine *infrav1.AzureMachine
	cache        *MachineCache
	skuCache     SKUCacher
}

// SKUCacher fetches a SKU from its cache.
type SKUCacher interface {
	Get(context.Context, string, resourceskus.ResourceType) (resourceskus.SKU, error)
}

// MachineCache stores common machine information so we don't have to hit the API multiple times within the same reconcile loop.
type MachineCache struct {
	BootstrapData      string
	VMImage            *infrav1.Image
	VMSKU              resourceskus.SKU
	availabilitySetSKU resourceskus.SKU
}

// InitMachineCache sets cached information about the machine to be used in the scope.
func (m *MachineScope) InitMachineCache(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "azure.MachineScope.InitMachineCache")
	defer done()

	if m.cache == nil {
		var err error
		m.cache = &MachineCache{}

		m.cache.BootstrapData, err = m.GetBootstrapData(ctx)
		if err != nil {
			return err
		}

		m.cache.VMImage, err = m.GetVMImage(ctx)
		if err != nil {
			return err
		}

		skuCache := m.skuCache
		if skuCache == nil {
			cache, err := resourceskus.GetCache(m, m.Location())
			if err != nil {
				return err
			}
			skuCache = cache
		}

		m.cache.VMSKU, err = skuCache.Get(ctx, m.AzureMachine.Spec.VMSize, resourceskus.VirtualMachines)
		if err != nil {
			return errors.Wrapf(err, "failed to get VM SKU %s in compute api", m.AzureMachine.Spec.VMSize)
		}

		m.cache.availabilitySetSKU, err = skuCache.Get(ctx, string(armcompute.AvailabilitySetSKUTypesAligned), resourceskus.AvailabilitySets)
		// Resource SKU API for availability sets may not be available in Azure Stack environments.
		if err != nil && !strings.EqualFold(m.CloudEnvironment(), "HybridEnvironment") {
			return errors.Wrapf(err, "failed to get availability set SKU %s in compute api", string(armcompute.AvailabilitySetSKUTypesAligned))
		}
	}

	return nil
}

// VMSpec returns the VM spec.
func (m *MachineScope) VMSpec() azure.ResourceSpecGetter {
	spec := &virtualmachines.VMSpec{
		Name:                       m.Name(),
		Location:                   m.Location(),
		ExtendedLocation:           m.ExtendedLocation(),
		ResourceGroup:              m.NodeResourceGroup(),
		ClusterName:                m.ClusterName(),
		Role:                       m.Role(),
		NICIDs:                     m.NICIDs(),
		SSHKeyData:                 m.AzureMachine.Spec.SSHPublicKey,
		Size:                       m.AzureMachine.Spec.VMSize,
		OSDisk:                     m.AzureMachine.Spec.OSDisk,
		DataDisks:                  m.AzureMachine.Spec.DataDisks,
		AvailabilitySetID:          m.AvailabilitySetID(),
		Zone:                       m.AvailabilityZone(),
		Identity:                   m.AzureMachine.Spec.Identity,
		UserAssignedIdentities:     m.AzureMachine.Spec.UserAssignedIdentities,
		SpotVMOptions:              m.AzureMachine.Spec.SpotVMOptions,
		SecurityProfile:            m.AzureMachine.Spec.SecurityProfile,
		DiagnosticsProfile:         m.AzureMachine.Spec.Diagnostics,
		DisableExtensionOperations: ptr.Deref(m.AzureMachine.Spec.DisableExtensionOperations, false),
		AdditionalTags:             m.AdditionalTags(),
		AdditionalCapabilities:     m.AzureMachine.Spec.AdditionalCapabilities,
		CapacityReservationGroupID: m.GetCapacityReservationGroupID(),
		ProviderID:                 m.ProviderID(),
	}
	if m.cache != nil {
		spec.SKU = m.cache.VMSKU
		spec.Image = m.cache.VMImage
		spec.BootstrapData = m.cache.BootstrapData
	}
	return spec
}

// TagsSpecs returns the tags for the AzureMachine.
func (m *MachineScope) TagsSpecs() []azure.TagsSpec {
	return []azure.TagsSpec{
		{
			Scope:      azure.VMID(m.SubscriptionID(), m.NodeResourceGroup(), m.Name()),
			Tags:       m.AdditionalTags(),
			Annotation: azure.VMTagsLastAppliedAnnotation,
		},
	}
}

// PublicIPSpecs returns the public IP specs.
func (m *MachineScope) PublicIPSpecs() []azure.ResourceSpecGetter {
	var specs []azure.ResourceSpecGetter
	if m.AzureMachine.Spec.AllocatePublicIP {
		specs = append(specs, &publicips.PublicIPSpec{
			Name:             azure.GenerateNodePublicIPName(m.Name()),
			ResourceGroup:    m.NodeResourceGroup(),
			ClusterName:      m.ClusterName(),
			DNSName:          "",    // Set to default value
			IsIPv6:           false, // Set to default value
			Location:         m.Location(),
			ExtendedLocation: m.ExtendedLocation(),
			FailureDomains:   m.FailureDomains(),
			AdditionalTags:   m.ClusterScoper.AdditionalTags(),
		})
	}
	return specs
}

// InboundNatSpecs returns the inbound NAT specs.
func (m *MachineScope) InboundNatSpecs() []azure.ResourceSpecGetter {
	// The existing inbound NAT rules are needed in order to find an available SSH port for each new inbound NAT rule.
	if m.Role() == infrav1.ControlPlane {
		spec := &inboundnatrules.InboundNatSpec{
			Name:                      m.Name(),
			ResourceGroup:             m.NodeResourceGroup(),
			LoadBalancerName:          m.APIServerLBName(),
			FrontendIPConfigurationID: nil,
		}
		if frontEndIPs := m.APIServerLB().FrontendIPs; len(frontEndIPs) > 0 {
			ipConfig := frontEndIPs[0].Name
			id := azure.FrontendIPConfigID(m.SubscriptionID(), m.NodeResourceGroup(), m.APIServerLBName(), ipConfig)
			spec.FrontendIPConfigurationID = ptr.To(id)
		}

		return []azure.ResourceSpecGetter{spec}
	}
	return []azure.ResourceSpecGetter{}
}

// NICSpecs returns the network interface specs.
func (m *MachineScope) NICSpecs() []azure.ResourceSpecGetter {
	nicSpecs := []azure.ResourceSpecGetter{}

	// For backwards compatibility we need to ensure the NIC Name does not change on existing machines
	// created prior to multiple NIC support
	isMultiNIC := len(m.AzureMachine.Spec.NetworkInterfaces) > 1

	for i := 0; i < len(m.AzureMachine.Spec.NetworkInterfaces); i++ {
		isPrimary := i == 0
		nicName := azure.GenerateNICName(m.Name(), isMultiNIC, i)
		nicSpecs = append(nicSpecs, m.BuildNICSpec(nicName, m.AzureMachine.Spec.NetworkInterfaces[i], isPrimary))
	}
	return nicSpecs
}

// BuildNICSpec takes a NetworkInterface from the AzureMachineSpec and returns a NICSpec for use by the networkinterfaces service.
func (m *MachineScope) BuildNICSpec(nicName string, infrav1NetworkInterface infrav1.NetworkInterface, primaryNetworkInterface bool) *networkinterfaces.NICSpec {
	spec := &networkinterfaces.NICSpec{
		Name:                  nicName,
		ResourceGroup:         m.NodeResourceGroup(),
		Location:              m.Location(),
		ExtendedLocation:      m.ExtendedLocation(),
		SubscriptionID:        m.SubscriptionID(),
		MachineName:           m.Name(),
		VNetName:              m.Vnet().Name,
		VNetResourceGroup:     m.Vnet().ResourceGroup,
		AcceleratedNetworking: infrav1NetworkInterface.AcceleratedNetworking,
		IPv6Enabled:           m.IsIPv6Enabled(),
		EnableIPForwarding:    m.AzureMachine.Spec.EnableIPForwarding,
		SubnetName:            infrav1NetworkInterface.SubnetName,
		AdditionalTags:        m.AdditionalTags(),
		ClusterName:           m.ClusterName(),
		IPConfigs:             []networkinterfaces.IPConfig{},
	}

	if m.cache != nil {
		spec.SKU = &m.cache.VMSKU
	}

	for i := 0; i < infrav1NetworkInterface.PrivateIPConfigs; i++ {
		spec.IPConfigs = append(spec.IPConfigs, networkinterfaces.IPConfig{})
	}

	if primaryNetworkInterface {
		spec.DNSServers = m.AzureMachine.Spec.DNSServers

		if m.Role() == infrav1.ControlPlane {
			spec.PublicLBName = m.OutboundLBName(m.Role())
			spec.PublicLBAddressPoolName = m.OutboundPoolName(m.Role())
			if m.IsAPIServerPrivate() {
				spec.InternalLBName = m.APIServerLBName()
				spec.InternalLBAddressPoolName = m.APIServerLBPoolName()
			} else {
				if feature.Gates.Enabled(feature.APIServerILB) {
					spec.InternalLBName = m.APIServerLBName() + "-internal"
					spec.InternalLBAddressPoolName = m.APIServerLBPoolName() + "-internal"
				}
				spec.PublicLBNATRuleName = m.Name()
				spec.PublicLBAddressPoolName = m.APIServerLBPoolName()
			}
		}

		if m.Role() == infrav1.Node && m.AzureMachine.Spec.AllocatePublicIP {
			spec.PublicIPName = azure.GenerateNodePublicIPName(m.Name())
		}
		// If the NAT gateway is not enabled and node has no public IP, then the NIC needs to reference the LB to get outbound traffic.
		if m.Role() == infrav1.Node && !m.Subnet().IsNatGatewayEnabled() && !m.AzureMachine.Spec.AllocatePublicIP {
			spec.PublicLBName = m.OutboundLBName(m.Role())
			spec.PublicLBAddressPoolName = m.OutboundPoolName(m.Role())
		}
	}

	return spec
}

// NICIDs returns the NIC resource IDs.
func (m *MachineScope) NICIDs() []string {
	nicspecs := m.NICSpecs()
	nicIDs := make([]string, len(nicspecs))
	for i, nic := range nicspecs {
		nicIDs[i] = azure.NetworkInterfaceID(m.SubscriptionID(), nic.ResourceGroupName(), nic.ResourceName())
	}

	return nicIDs
}

// DiskSpecs returns the disk specs.
func (m *MachineScope) DiskSpecs() []azure.ResourceSpecGetter {
	diskSpecs := make([]azure.ResourceSpecGetter, 1+len(m.AzureMachine.Spec.DataDisks))
	diskSpecs[0] = &disks.DiskSpec{
		Name:          azure.GenerateOSDiskName(m.Name()),
		ResourceGroup: m.NodeResourceGroup(),
	}

	for i, dd := range m.AzureMachine.Spec.DataDisks {
		diskSpecs[i+1] = &disks.DiskSpec{
			Name:          azure.GenerateDataDiskName(m.Name(), dd.NameSuffix),
			ResourceGroup: m.NodeResourceGroup(),
		}
	}
	return diskSpecs
}

// RoleAssignmentSpecs returns the role assignment specs.
func (m *MachineScope) RoleAssignmentSpecs(principalID *string) []azure.ResourceSpecGetter {
	roles := make([]azure.ResourceSpecGetter, 1)
	if m.HasSystemAssignedIdentity() {
		roles[0] = &roleassignments.RoleAssignmentSpec{
			Name:             m.SystemAssignedIdentityName(),
			MachineName:      m.Name(),
			ResourceType:     azure.VirtualMachine,
			ResourceGroup:    m.NodeResourceGroup(),
			Scope:            m.SystemAssignedIdentityScope(),
			RoleDefinitionID: m.SystemAssignedIdentityDefinitionID(),
			PrincipalID:      principalID,
			PrincipalType:    armauthorization.PrincipalTypeServicePrincipal,
		}
		return roles
	}
	return []azure.ResourceSpecGetter{}
}

// RoleAssignmentResourceType returns the role assignment resource type.
func (m *MachineScope) RoleAssignmentResourceType() string {
	return azure.VirtualMachine
}

// HasSystemAssignedIdentity returns true if the azure machine has
// system assigned identity.
func (m *MachineScope) HasSystemAssignedIdentity() bool {
	return m.AzureMachine.Spec.Identity == infrav1.VMIdentitySystemAssigned
}

// VMExtensionSpecs returns the VM extension specs.
func (m *MachineScope) VMExtensionSpecs() []azure.ResourceSpecGetter {
	if ptr.Deref(m.AzureMachine.Spec.DisableExtensionOperations, false) {
		return []azure.ResourceSpecGetter{}
	}

	var extensionSpecs = []azure.ResourceSpecGetter{}
	for _, extension := range m.AzureMachine.Spec.VMExtensions {
		extensionSpecs = append(extensionSpecs, &vmextensions.VMExtensionSpec{
			ExtensionSpec: azure.ExtensionSpec{
				Name:              extension.Name,
				VMName:            m.Name(),
				Publisher:         extension.Publisher,
				Version:           extension.Version,
				Settings:          extension.Settings,
				ProtectedSettings: extension.ProtectedSettings,
			},
			ResourceGroup: m.NodeResourceGroup(),
			Location:      m.Location(),
		})
	}

	cpuArchitectureType, _ := m.cache.VMSKU.GetCapability(resourceskus.CPUArchitectureType)
	bootstrapExtensionSpec := azure.GetBootstrappingVMExtension(m.AzureMachine.Spec.OSDisk.OSType, m.CloudEnvironment(), m.Name(), cpuArchitectureType)

	if bootstrapExtensionSpec != nil {
		extensionSpecs = append(extensionSpecs, &vmextensions.VMExtensionSpec{
			ExtensionSpec: *bootstrapExtensionSpec,
			ResourceGroup: m.NodeResourceGroup(),
			Location:      m.Location(),
		})
	}

	return extensionSpecs
}

// Subnet returns the machine's subnet.
func (m *MachineScope) Subnet() infrav1.SubnetSpec {
	for _, subnet := range m.Subnets() {
		if subnet.Name == m.AzureMachine.Spec.NetworkInterfaces[0].SubnetName {
			return subnet
		}
	}

	return infrav1.SubnetSpec{}
}

// AvailabilityZone returns the AzureMachine Availability Zone.
// Priority for selecting the AZ is
//  1. Machine.Spec.FailureDomain
//  2. AzureMachine.Spec.FailureDomain (This is to support deprecated AZ)
//  3. No AZ
func (m *MachineScope) AvailabilityZone() string {
	if m.Machine.Spec.FailureDomain != nil {
		return *m.Machine.Spec.FailureDomain
	}
	// Deprecated: to support old clients
	if m.AzureMachine.Spec.FailureDomain != nil {
		return *m.AzureMachine.Spec.FailureDomain
	}

	return ""
}

// Name returns the AzureMachine name.
func (m *MachineScope) Name() string {
	if id := m.GetVMID(); id != "" {
		return id
	}
	// Windows Machine names cannot be longer than 15 chars
	if m.AzureMachine.Spec.OSDisk.OSType == azure.WindowsOS && len(m.AzureMachine.Name) > 15 {
		return strings.TrimSuffix(m.AzureMachine.Name[0:9], "-") + "-" + m.AzureMachine.Name[len(m.AzureMachine.Name)-5:]
	}
	return m.AzureMachine.Name
}

// Namespace returns the namespace name.
func (m *MachineScope) Namespace() string {
	return m.AzureMachine.Namespace
}

// IsControlPlane returns true if the machine is a control plane.
func (m *MachineScope) IsControlPlane() bool {
	return util.IsControlPlaneMachine(m.Machine)
}

// Role returns the machine role from the labels.
func (m *MachineScope) Role() string {
	if util.IsControlPlaneMachine(m.Machine) {
		return infrav1.ControlPlane
	}
	return infrav1.Node
}

// GetVMID returns the AzureMachine instance id by parsing the scope's providerID.
func (m *MachineScope) GetVMID() string {
	resourceID, err := azureutil.ParseResourceID(m.ProviderID())
	if err != nil {
		return ""
	}
	return resourceID.Name
}

// ProviderID returns the AzureMachine providerID from the spec.
func (m *MachineScope) ProviderID() string {
	return ptr.Deref(m.AzureMachine.Spec.ProviderID, "")
}

// AvailabilitySetSpec returns the availability set spec for this machine if available.
func (m *MachineScope) AvailabilitySetSpec() azure.ResourceSpecGetter {
	availabilitySetName, ok := m.AvailabilitySet()
	if !ok {
		return nil
	}

	spec := &availabilitysets.AvailabilitySetSpec{
		Name:             availabilitySetName,
		ResourceGroup:    m.NodeResourceGroup(),
		ClusterName:      m.ClusterName(),
		Location:         m.Location(),
		CloudEnvironment: m.CloudEnvironment(),
		SKU:              nil,
		AdditionalTags:   m.AdditionalTags(),
	}

	if m.cache != nil {
		spec.SKU = &m.cache.availabilitySetSKU
	}

	return spec
}

// AvailabilitySet returns the availability set for this machine if available.
func (m *MachineScope) AvailabilitySet() (string, bool) {
	// AvailabilitySet service is not supported on EdgeZone currently.
	// AvailabilitySet cannot be used with Spot instances.
	if !m.AvailabilitySetEnabled() || m.AzureMachine.Spec.SpotVMOptions != nil || m.ExtendedLocation() != nil ||
		m.AzureMachine.Spec.FailureDomain != nil || m.Machine.Spec.FailureDomain != nil {
		return "", false
	}

	if m.IsControlPlane() {
		return azure.GenerateAvailabilitySetName(m.ClusterName(), azure.ControlPlaneNodeGroup), true
	}

	// get machine deployment name from labels for machines that maybe part of a machine deployment.
	if mdName, ok := m.Machine.Labels[clusterv1.MachineDeploymentNameLabel]; ok {
		return azure.GenerateAvailabilitySetName(m.ClusterName(), mdName), true
	}

	// if machine deployment name label is not available, use machine set name.
	if msName, ok := m.Machine.Labels[clusterv1.MachineSetNameLabel]; ok {
		return azure.GenerateAvailabilitySetName(m.ClusterName(), msName), true
	}

	return "", false
}

// AvailabilitySetID returns the availability set for this machine, or "" if there is no availability set.
func (m *MachineScope) AvailabilitySetID() string {
	var asID string
	if asName, ok := m.AvailabilitySet(); ok {
		asID = azure.AvailabilitySetID(m.SubscriptionID(), m.NodeResourceGroup(), asName)
	}
	return asID
}

// SystemAssignedIdentityName returns the role assignment name for the system assigned identity.
func (m *MachineScope) SystemAssignedIdentityName() string {
	if m.AzureMachine.Spec.SystemAssignedIdentityRole != nil {
		return m.AzureMachine.Spec.SystemAssignedIdentityRole.Name
	}
	return ""
}

// SystemAssignedIdentityScope returns the scope for the system assigned identity.
func (m *MachineScope) SystemAssignedIdentityScope() string {
	if m.AzureMachine.Spec.SystemAssignedIdentityRole != nil {
		return m.AzureMachine.Spec.SystemAssignedIdentityRole.Scope
	}
	return ""
}

// SystemAssignedIdentityDefinitionID returns the role definition id for the system assigned identity.
func (m *MachineScope) SystemAssignedIdentityDefinitionID() string {
	if m.AzureMachine.Spec.SystemAssignedIdentityRole != nil {
		return m.AzureMachine.Spec.SystemAssignedIdentityRole.DefinitionID
	}
	return ""
}

// SetProviderID sets the AzureMachine providerID in spec.
func (m *MachineScope) SetProviderID(v string) {
	m.AzureMachine.Spec.ProviderID = ptr.To(v)
}

// VMState returns the AzureMachine VM state.
func (m *MachineScope) VMState() infrav1.ProvisioningState {
	if m.AzureMachine.Status.VMState != nil {
		return *m.AzureMachine.Status.VMState
	}
	return ""
}

// SetVMState sets the AzureMachine VM state.
func (m *MachineScope) SetVMState(v infrav1.ProvisioningState) {
	m.AzureMachine.Status.VMState = &v
}

// SetReady sets the AzureMachine Ready Status to true.
func (m *MachineScope) SetReady() {
	m.AzureMachine.Status.Ready = true
}

// SetNotReady sets the AzureMachine Ready Status to false.
func (m *MachineScope) SetNotReady() {
	m.AzureMachine.Status.Ready = false
}

// SetFailureMessage sets the AzureMachine status failure message.
func (m *MachineScope) SetFailureMessage(v error) {
	m.AzureMachine.Status.FailureMessage = ptr.To(v.Error())
}

// SetFailureReason sets the AzureMachine status failure reason.
func (m *MachineScope) SetFailureReason(v string) {
	m.AzureMachine.Status.FailureReason = &v
}

// SetConditionFalse sets the specified AzureMachine condition to false.
func (m *MachineScope) SetConditionFalse(conditionType clusterv1.ConditionType, reason string, severity clusterv1.ConditionSeverity, message string) {
	conditions.MarkFalse(m.AzureMachine, conditionType, reason, severity, message)
}

// SetAnnotation sets a key value annotation on the AzureMachine.
func (m *MachineScope) SetAnnotation(key, value string) {
	if m.AzureMachine.Annotations == nil {
		m.AzureMachine.Annotations = map[string]string{}
	}
	m.AzureMachine.Annotations[key] = value
}

// AnnotationJSON returns a map[string]interface from a JSON annotation.
func (m *MachineScope) AnnotationJSON(annotation string) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	jsonAnnotation := m.AzureMachine.GetAnnotations()[annotation]
	if jsonAnnotation == "" {
		return out, nil
	}
	err := json.Unmarshal([]byte(jsonAnnotation), &out)
	if err != nil {
		return out, err
	}
	return out, nil
}

// UpdateAnnotationJSON updates the `annotation` with
// `content`. `content` in this case should be a `map[string]interface{}`
// suitable for turning into JSON. This `content` map will be marshalled into a
// JSON string before being set as the given `annotation`.
func (m *MachineScope) UpdateAnnotationJSON(annotation string, content map[string]interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}
	m.SetAnnotation(annotation, string(b))
	return nil
}

// SetAddresses sets the Azure address status.
func (m *MachineScope) SetAddresses(addrs []corev1.NodeAddress) {
	m.AzureMachine.Status.Addresses = addrs
}

// PatchObject persists the machine spec and status.
func (m *MachineScope) PatchObject(ctx context.Context) error {
	conditions.SetSummary(m.AzureMachine)

	return m.patchHelper.Patch(
		ctx,
		m.AzureMachine,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
			infrav1.VMRunningCondition,
			infrav1.AvailabilitySetReadyCondition,
			infrav1.NetworkInterfaceReadyCondition,
		}})
}

// Close the MachineScope by updating the machine spec, machine status.
func (m *MachineScope) Close(ctx context.Context) error {
	return m.PatchObject(ctx)
}

// AdditionalTags merges AdditionalTags from the scope's AzureCluster and AzureMachine. If the same key is present in both,
// the value from AzureMachine takes precedence.
func (m *MachineScope) AdditionalTags() infrav1.Tags {
	tags := make(infrav1.Tags)
	// Start with the cluster-wide tags...
	tags.Merge(m.ClusterScoper.AdditionalTags())
	// ... and merge in the Machine's
	tags.Merge(m.AzureMachine.Spec.AdditionalTags)
	// Set the cloud provider tag
	tags[infrav1.ClusterAzureCloudProviderTagKey(m.ClusterName())] = string(infrav1.ResourceLifecycleOwned)

	return tags
}

// GetBootstrapData returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName.
func (m *MachineScope) GetBootstrapData(ctx context.Context) (string, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scope.MachineScope.GetBootstrapData")
	defer done()

	if m.Machine.Spec.Bootstrap.DataSecretName == nil {
		return "", errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}
	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: m.Namespace(), Name: *m.Machine.Spec.Bootstrap.DataSecretName}
	if err := m.client.Get(ctx, key, secret); err != nil {
		return "", errors.Wrapf(err, "failed to retrieve bootstrap data secret for AzureMachine %s/%s", m.Namespace(), m.Name())
	}

	value, ok := secret.Data["value"]
	if !ok {
		return "", errors.New("error retrieving bootstrap data: secret value key is missing")
	}
	return base64.StdEncoding.EncodeToString(value), nil
}

// GetVMImage returns the image from the machine configuration, or a default one.
func (m *MachineScope) GetVMImage(ctx context.Context) (*infrav1.Image, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "scope.MachineScope.GetVMImage")
	defer done()

	// Use custom Marketplace image, Image ID or a Shared Image Gallery image if provided
	if m.AzureMachine.Spec.Image != nil {
		return m.AzureMachine.Spec.Image, nil
	}

	svc, err := virtualmachineimages.New(m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create virtualmachineimages service")
	}

	if m.AzureMachine.Spec.OSDisk.OSType == azure.WindowsOS {
		runtime := m.AzureMachine.Annotations["runtime"]
		windowsServerVersion := m.AzureMachine.Annotations["windowsServerVersion"]
		log.Info("No image specified for machine, using default Windows Image", "machine", m.AzureMachine.GetName(), "runtime", runtime, "windowsServerVersion", windowsServerVersion)
		return svc.GetDefaultWindowsImage(ctx, m.Location(), ptr.Deref(m.Machine.Spec.Version, ""), runtime, windowsServerVersion)
	}

	log.Info("No image specified for machine, using default Linux Image", "machine", m.AzureMachine.GetName())
	return svc.GetDefaultLinuxImage(ctx, m.Location(), ptr.Deref(m.Machine.Spec.Version, ""))
}

// SetSubnetName defaults the AzureMachine subnet name to the name of one the subnets with the machine role when there is only one of them.
// Note: this logic exists only for purposes of ensuring backwards compatibility for old clusters created without the `subnetName` field being
// set, and should be removed in the future when this field is no longer optional.
func (m *MachineScope) SetSubnetName() error {
	if m.AzureMachine.Spec.NetworkInterfaces[0].SubnetName == "" {
		subnetName := ""
		subnets := m.Subnets()
		var subnetCount int
		clusterSubnetName := ""
		for _, subnet := range subnets {
			if string(subnet.Role) == m.Role() {
				subnetCount++
				subnetName = subnet.Name
			}
			if subnet.Role == infrav1.SubnetCluster {
				clusterSubnetName = subnet.Name
			}
		}

		if subnetName == "" && clusterSubnetName != "" {
			subnetName = clusterSubnetName
			subnetCount = 1
		}

		if subnetCount == 0 || subnetCount > 1 || subnetName == "" {
			return errors.New("a subnet name must be specified when no subnets are specified or more than 1 subnet of the same role exist")
		}

		m.AzureMachine.Spec.NetworkInterfaces[0].SubnetName = subnetName
	}

	return nil
}

// SetLongRunningOperationState will set the future on the AzureMachine status to allow the resource to continue
// in the next reconciliation.
func (m *MachineScope) SetLongRunningOperationState(future *infrav1.Future) {
	futures.Set(m.AzureMachine, future)
}

// GetLongRunningOperationState will get the future on the AzureMachine status.
func (m *MachineScope) GetLongRunningOperationState(name, service, futureType string) *infrav1.Future {
	return futures.Get(m.AzureMachine, name, service, futureType)
}

// DeleteLongRunningOperationState will delete the future from the AzureMachine status.
func (m *MachineScope) DeleteLongRunningOperationState(name, service, futureType string) {
	futures.Delete(m.AzureMachine, name, service, futureType)
}

// UpdateDeleteStatus updates a condition on the AzureMachine status after a DELETE operation.
func (m *MachineScope) UpdateDeleteStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkFalse(m.AzureMachine, condition, infrav1.DeletedReason, clusterv1.ConditionSeverityInfo, "%s successfully deleted", service)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(m.AzureMachine, condition, infrav1.DeletingReason, clusterv1.ConditionSeverityInfo, "%s deleting", service)
	default:
		conditions.MarkFalse(m.AzureMachine, condition, infrav1.DeletionFailedReason, clusterv1.ConditionSeverityError, "%s failed to delete. err: %s", service, err.Error())
	}
}

// UpdatePutStatus updates a condition on the AzureMachine status after a PUT operation.
func (m *MachineScope) UpdatePutStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkTrue(m.AzureMachine, condition)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(m.AzureMachine, condition, infrav1.CreatingReason, clusterv1.ConditionSeverityInfo, "%s creating or updating", service)
	default:
		conditions.MarkFalse(m.AzureMachine, condition, infrav1.FailedReason, clusterv1.ConditionSeverityError, "%s failed to create or update. err: %s", service, err.Error())
	}
}

// UpdatePatchStatus updates a condition on the AzureMachine status after a PATCH operation.
func (m *MachineScope) UpdatePatchStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkTrue(m.AzureMachine, condition)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(m.AzureMachine, condition, infrav1.UpdatingReason, clusterv1.ConditionSeverityInfo, "%s updating", service)
	default:
		conditions.MarkFalse(m.AzureMachine, condition, infrav1.FailedReason, clusterv1.ConditionSeverityError, "%s failed to update. err: %s", service, err.Error())
	}
}

// GetCapacityReservationGroupID returns the CapacityReservationGroupID from the spec if the
// value is assigned, or else returns an empty string.
func (m *MachineScope) GetCapacityReservationGroupID() string {
	return ptr.Deref(m.AzureMachine.Spec.CapacityReservationGroupID, "")
}
