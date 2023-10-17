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

package scope

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/groups"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/managedclusters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/privateendpoints"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/subnets"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/virtualnetworks"
	"sigs.k8s.io/cluster-api-provider-azure/util/futures"
	"sigs.k8s.io/cluster-api-provider-azure/util/maps"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/secret"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const resourceHealthWarningInitialGracePeriod = 1 * time.Hour

// ManagedControlPlaneScopeParams defines the input parameters used to create a new managed
// control plane.
type ManagedControlPlaneScopeParams struct {
	AzureClients
	Client              client.Client
	Cluster             *clusterv1.Cluster
	ControlPlane        *infrav1.AzureManagedControlPlane
	ManagedMachinePools []ManagedMachinePool
	Cache               *ManagedControlPlaneCache
	VnetDescriber       VnetDescriber
}

// VnetDescriber answers whether a virtual network is managed or not.
type VnetDescriber interface {
	IsManaged(context.Context) (bool, error)
}

// NewManagedControlPlaneScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewManagedControlPlaneScope(ctx context.Context, params ManagedControlPlaneScopeParams) (*ManagedControlPlaneScope, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scope.NewManagedControlPlaneScope")
	defer done()

	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}

	if params.ControlPlane == nil {
		return nil, errors.New("failed to generate new scope from nil ControlPlane")
	}

	if params.ControlPlane.Spec.IdentityRef == nil {
		if err := params.AzureClients.setCredentials(params.ControlPlane.Spec.SubscriptionID, params.ControlPlane.Spec.AzureEnvironment); err != nil {
			return nil, errors.Wrap(err, "failed to create Azure session")
		}
	} else {
		credentialsProvider, err := NewManagedControlPlaneCredentialsProvider(ctx, params.Client, params.ControlPlane)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init credentials provider")
		}

		if err := params.AzureClients.setCredentialsWithProvider(ctx, params.ControlPlane.Spec.SubscriptionID, params.ControlPlane.Spec.AzureEnvironment, credentialsProvider); err != nil {
			return nil, errors.Wrap(err, "failed to configure azure settings and credentials for Identity")
		}
	}

	if params.Cache == nil {
		params.Cache = &ManagedControlPlaneCache{}
	}

	helper, err := patch.NewHelper(params.ControlPlane, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &ManagedControlPlaneScope{
		Client:              params.Client,
		AzureClients:        params.AzureClients,
		Cluster:             params.Cluster,
		ControlPlane:        params.ControlPlane,
		ManagedMachinePools: params.ManagedMachinePools,
		patchHelper:         helper,
		cache:               params.Cache,
		VnetDescriber:       params.VnetDescriber,
	}, nil
}

// ManagedControlPlaneScope defines the basic context for an actuator to operate upon.
type ManagedControlPlaneScope struct {
	Client         client.Client
	patchHelper    *patch.Helper
	kubeConfigData []byte
	cache          *ManagedControlPlaneCache

	AzureClients
	Cluster             *clusterv1.Cluster
	ControlPlane        *infrav1.AzureManagedControlPlane
	ManagedMachinePools []ManagedMachinePool
	VnetDescriber       VnetDescriber
}

// ManagedControlPlaneCache stores ManagedControlPlane data locally so we don't have to hit the API multiple times within the same reconcile loop.
type ManagedControlPlaneCache struct {
	isVnetManaged *bool
}

// GetClient returns the controller-runtime client.
func (s *ManagedControlPlaneScope) GetClient() client.Client {
	return s.Client
}

// ResourceGroup returns the managed control plane's resource group.
func (s *ManagedControlPlaneScope) ResourceGroup() string {
	if s.ControlPlane == nil {
		return ""
	}
	return s.ControlPlane.Spec.ResourceGroupName
}

// NodeResourceGroup returns the managed control plane's node resource group.
func (s *ManagedControlPlaneScope) NodeResourceGroup() string {
	if s.ControlPlane == nil {
		return ""
	}
	return s.ControlPlane.Spec.NodeResourceGroupName
}

// ClusterName returns the managed control plane's name.
func (s *ManagedControlPlaneScope) ClusterName() string {
	return s.Cluster.Name
}

// Location returns the managed control plane's Azure location, or an empty string.
func (s *ManagedControlPlaneScope) Location() string {
	if s.ControlPlane == nil {
		return ""
	}
	return s.ControlPlane.Spec.Location
}

// ExtendedLocation has not been implemented for AzureManagedControlPlane.
func (s *ManagedControlPlaneScope) ExtendedLocation() *infrav1.ExtendedLocationSpec {
	return nil
}

// ExtendedLocationName has not been implemented for AzureManagedControlPlane.
func (s *ManagedControlPlaneScope) ExtendedLocationName() string {
	return ""
}

// ExtendedLocationType has not been implemented for AzureManagedControlPlane.
func (s *ManagedControlPlaneScope) ExtendedLocationType() string {
	return ""
}

// AvailabilitySetEnabled is always false for a managed control plane.
func (s *ManagedControlPlaneScope) AvailabilitySetEnabled() bool {
	return false // not applicable for a managed control plane
}

// AdditionalTags returns AdditionalTags from the ControlPlane spec.
func (s *ManagedControlPlaneScope) AdditionalTags() infrav1.Tags {
	tags := make(infrav1.Tags)
	if s.ControlPlane.Spec.AdditionalTags != nil {
		tags = s.ControlPlane.Spec.AdditionalTags.DeepCopy()
	}
	return tags
}

// SubscriptionID returns the Azure client Subscription ID.
func (s *ManagedControlPlaneScope) SubscriptionID() string {
	return s.AzureClients.SubscriptionID()
}

// BaseURI returns the Azure ResourceManagerEndpoint.
func (s *ManagedControlPlaneScope) BaseURI() string {
	return s.AzureClients.ResourceManagerEndpoint
}

// Authorizer returns the Azure client Authorizer.
func (s *ManagedControlPlaneScope) Authorizer() autorest.Authorizer {
	return s.AzureClients.Authorizer
}

// PatchObject persists the cluster configuration and status.
func (s *ManagedControlPlaneScope) PatchObject(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scope.ManagedControlPlaneScope.PatchObject")
	defer done()

	conditions.SetSummary(s.ControlPlane)

	return s.patchHelper.Patch(
		ctx,
		s.ControlPlane,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
			infrav1.ResourceGroupReadyCondition,
			infrav1.VNetReadyCondition,
			infrav1.SubnetsReadyCondition,
			infrav1.ManagedClusterRunningCondition,
			infrav1.AgentPoolsReadyCondition,
			infrav1.AzureResourceAvailableCondition,
		}})
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ManagedControlPlaneScope) Close(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scope.ManagedControlPlaneScope.Close")
	defer done()

	return s.PatchObject(ctx)
}

// Vnet returns the cluster Vnet.
func (s *ManagedControlPlaneScope) Vnet() *infrav1.VnetSpec {
	return &infrav1.VnetSpec{
		ResourceGroup: s.ControlPlane.Spec.VirtualNetwork.ResourceGroup,
		Name:          s.ControlPlane.Spec.VirtualNetwork.Name,
		VnetClassSpec: infrav1.VnetClassSpec{
			CIDRBlocks: []string{s.ControlPlane.Spec.VirtualNetwork.CIDRBlock},
		},
	}
}

// GroupSpec returns the resource group spec.
func (s *ManagedControlPlaneScope) GroupSpec() azure.ASOResourceSpecGetter {
	return &groups.GroupSpec{
		Name:           s.ResourceGroup(),
		Namespace:      s.Cluster.Namespace,
		Location:       s.Location(),
		ClusterName:    s.ClusterName(),
		AdditionalTags: s.AdditionalTags(),
		Owner:          *metav1.NewControllerRef(s.ControlPlane, infrav1.GroupVersion.WithKind("AzureManagedControlPlane")),
	}
}

// VNetSpec returns the virtual network spec.
func (s *ManagedControlPlaneScope) VNetSpec() azure.ResourceSpecGetter {
	return &virtualnetworks.VNetSpec{
		ResourceGroup:  s.Vnet().ResourceGroup,
		Name:           s.Vnet().Name,
		CIDRs:          s.Vnet().CIDRBlocks,
		Location:       s.Location(),
		ClusterName:    s.ClusterName(),
		AdditionalTags: s.AdditionalTags(),
	}
}

// ControlPlaneRouteTable returns the cluster controlplane routetable.
func (s *ManagedControlPlaneScope) ControlPlaneRouteTable() infrav1.RouteTable {
	return infrav1.RouteTable{}
}

// NodeRouteTable returns the cluster node routetable.
func (s *ManagedControlPlaneScope) NodeRouteTable() infrav1.RouteTable {
	return infrav1.RouteTable{}
}

// NodeNatGateway returns the cluster node NAT gateway.
func (s *ManagedControlPlaneScope) NodeNatGateway() infrav1.NatGateway {
	return infrav1.NatGateway{}
}

// SubnetSpecs returns the subnets specs.
func (s *ManagedControlPlaneScope) SubnetSpecs() []azure.ResourceSpecGetter {
	return []azure.ResourceSpecGetter{
		&subnets.SubnetSpec{
			Name:              s.NodeSubnet().Name,
			ResourceGroup:     s.ResourceGroup(),
			SubscriptionID:    s.SubscriptionID(),
			CIDRs:             s.NodeSubnet().CIDRBlocks,
			VNetName:          s.Vnet().Name,
			VNetResourceGroup: s.Vnet().ResourceGroup,
			IsVNetManaged:     s.IsVnetManaged(),
			Role:              infrav1.SubnetNode,
			ServiceEndpoints:  s.NodeSubnet().ServiceEndpoints,
		},
	}
}

// Subnets returns the subnets specs.
func (s *ManagedControlPlaneScope) Subnets() infrav1.Subnets {
	return infrav1.Subnets{}
}

// NodeSubnet returns the cluster node subnet.
func (s *ManagedControlPlaneScope) NodeSubnet() infrav1.SubnetSpec {
	return infrav1.SubnetSpec{
		SubnetClassSpec: infrav1.SubnetClassSpec{
			CIDRBlocks:       []string{s.ControlPlane.Spec.VirtualNetwork.Subnet.CIDRBlock},
			Name:             s.ControlPlane.Spec.VirtualNetwork.Subnet.Name,
			ServiceEndpoints: s.ControlPlane.Spec.VirtualNetwork.Subnet.ServiceEndpoints,
			PrivateEndpoints: s.ControlPlane.Spec.VirtualNetwork.Subnet.PrivateEndpoints,
		},
	}
}

// SetSubnet sets the passed subnet spec into the scope.
// This is not used when using a managed control plane.
func (s *ManagedControlPlaneScope) SetSubnet(_ infrav1.SubnetSpec) {
	// no-op
}

// UpdateSubnetCIDRs updates the subnet CIDRs for the subnet with the same name.
// This is not used when using a managed control plane.
func (s *ManagedControlPlaneScope) UpdateSubnetCIDRs(_ string, _ []string) {
	// no-op
}

// UpdateSubnetID updates the subnet ID for the subnet with the same name.
// This is not used when using a managed control plane.
func (s *ManagedControlPlaneScope) UpdateSubnetID(_ string, _ string) {
	// no-op
}

// ControlPlaneSubnet returns the cluster control plane subnet.
func (s *ManagedControlPlaneScope) ControlPlaneSubnet() infrav1.SubnetSpec {
	return infrav1.SubnetSpec{}
}

// NodeSubnets returns the subnets with the node role.
func (s *ManagedControlPlaneScope) NodeSubnets() []infrav1.SubnetSpec {
	return []infrav1.SubnetSpec{
		{
			SubnetClassSpec: infrav1.SubnetClassSpec{
				CIDRBlocks:       []string{s.ControlPlane.Spec.VirtualNetwork.Subnet.CIDRBlock},
				Name:             s.ControlPlane.Spec.VirtualNetwork.Subnet.Name,
				ServiceEndpoints: s.ControlPlane.Spec.VirtualNetwork.Subnet.ServiceEndpoints,
				PrivateEndpoints: s.ControlPlane.Spec.VirtualNetwork.Subnet.PrivateEndpoints,
			},
		},
	}
}

// Subnet returns the subnet with the provided name.
func (s *ManagedControlPlaneScope) Subnet(name string) infrav1.SubnetSpec {
	subnet := infrav1.SubnetSpec{}
	if name == s.ControlPlane.Spec.VirtualNetwork.Subnet.Name {
		subnet.Name = s.ControlPlane.Spec.VirtualNetwork.Subnet.Name
		subnet.CIDRBlocks = []string{s.ControlPlane.Spec.VirtualNetwork.Subnet.CIDRBlock}
		subnet.ServiceEndpoints = s.ControlPlane.Spec.VirtualNetwork.Subnet.ServiceEndpoints
		subnet.PrivateEndpoints = s.ControlPlane.Spec.VirtualNetwork.Subnet.PrivateEndpoints
	}

	return subnet
}

// IsIPv6Enabled returns true if a cluster is ipv6 enabled.
// Currently always false as managed control planes do not currently implement ipv6.
func (s *ManagedControlPlaneScope) IsIPv6Enabled() bool {
	return false
}

// IsVnetManaged returns true if the vnet is managed.
func (s *ManagedControlPlaneScope) IsVnetManaged() bool {
	if s.cache.isVnetManaged != nil {
		return ptr.Deref(s.cache.isVnetManaged, false)
	}
	// TODO refactor `IsVnetManaged` so that it is able to use an upstream context
	// see https://github.com/kubernetes-sigs/cluster-api-provider-azure/issues/2581
	ctx := context.Background()
	ctx, log, done := tele.StartSpanWithLogger(ctx, "scope.ManagedControlPlaneScope.IsVnetManaged")
	defer done()

	var vnetDescriber = s.VnetDescriber
	if vnetDescriber == nil {
		virtualNetworksSvc, err := virtualnetworks.New(s)
		if err != nil {
			log.Error(err, "failed to create virtualnetworks service")
			return false
		}
		vnetDescriber = virtualNetworksSvc
	}
	isManaged, err := vnetDescriber.IsManaged(ctx)
	if err != nil {
		log.Error(err, "Unable to determine if ManagedControlPlaneScope VNET is managed by capz", "AzureManagedCluster", s.ClusterName())
	}

	s.cache.isVnetManaged = ptr.To(isManaged)
	return isManaged
}

// APIServerLB returns the API Server LB spec.
func (s *ManagedControlPlaneScope) APIServerLB() *infrav1.LoadBalancerSpec {
	return nil // does not apply for AKS
}

// APIServerLBName returns the API Server LB name.
func (s *ManagedControlPlaneScope) APIServerLBName() string {
	return "" // does not apply for AKS
}

// APIServerLBPoolName returns the API Server LB backend pool name.
func (s *ManagedControlPlaneScope) APIServerLBPoolName(_ string) string {
	return "" // does not apply for AKS
}

// IsAPIServerPrivate returns true if the API Server LB is of type Internal.
// Currently always false as managed control planes do not currently implement private clusters.
func (s *ManagedControlPlaneScope) IsAPIServerPrivate() bool {
	return false
}

// OutboundLBName returns the name of the outbound LB.
// Note: for managed clusters, the outbound LB lifecycle is not managed.
func (s *ManagedControlPlaneScope) OutboundLBName(_ string) string {
	return "kubernetes"
}

// OutboundPoolName returns the outbound LB backend pool name.
func (s *ManagedControlPlaneScope) OutboundPoolName(_ string) string {
	return "aksOutboundBackendPool" // hard-coded in aks
}

// GetPrivateDNSZoneName returns the Private DNS Zone from the spec or generate it from cluster name.
// Currently always empty as managed control planes do not currently implement private clusters.
func (s *ManagedControlPlaneScope) GetPrivateDNSZoneName() string {
	return ""
}

// CloudProviderConfigOverrides returns the cloud provider config overrides for the cluster.
func (s *ManagedControlPlaneScope) CloudProviderConfigOverrides() *infrav1.CloudProviderConfigOverrides {
	return nil
}

// FailureDomains returns the failure domains for the cluster.
func (s *ManagedControlPlaneScope) FailureDomains() []*string {
	return []*string{}
}

// ManagedClusterAnnotations returns the annotations for the managed cluster.
func (s *ManagedControlPlaneScope) ManagedClusterAnnotations() map[string]string {
	return s.ControlPlane.Annotations
}

// ManagedClusterSpec returns the managed cluster spec.
func (s *ManagedControlPlaneScope) ManagedClusterSpec() azure.ResourceSpecGetter {
	managedClusterSpec := managedclusters.ManagedClusterSpec{
		Name:              s.ControlPlane.Name,
		ResourceGroup:     s.ControlPlane.Spec.ResourceGroupName,
		NodeResourceGroup: s.ControlPlane.Spec.NodeResourceGroupName,
		ClusterName:       s.ClusterName(),
		Location:          s.ControlPlane.Spec.Location,
		Tags:              s.ControlPlane.Spec.AdditionalTags,
		Headers:           maps.FilterByKeyPrefix(s.ManagedClusterAnnotations(), infrav1.CustomHeaderPrefix),
		Version:           strings.TrimPrefix(s.ControlPlane.Spec.Version, "v"),
		DNSServiceIP:      s.ControlPlane.Spec.DNSServiceIP,
		VnetSubnetID: azure.SubnetID(
			s.ControlPlane.Spec.SubscriptionID,
			s.VNetSpec().ResourceGroupName(),
			s.ControlPlane.Spec.VirtualNetwork.Name,
			s.ControlPlane.Spec.VirtualNetwork.Subnet.Name,
		),
		GetAllAgentPools:            s.GetAllAgentPoolSpecs,
		OutboundType:                s.ControlPlane.Spec.OutboundType,
		Identity:                    s.ControlPlane.Spec.Identity,
		KubeletUserAssignedIdentity: s.ControlPlane.Spec.KubeletUserAssignedIdentity,
		NetworkPluginMode:           s.ControlPlane.Spec.NetworkPluginMode,
	}

	if s.ControlPlane.Spec.SSHPublicKey != nil {
		managedClusterSpec.SSHPublicKey = *s.ControlPlane.Spec.SSHPublicKey
	}
	if s.ControlPlane.Spec.NetworkPlugin != nil {
		managedClusterSpec.NetworkPlugin = *s.ControlPlane.Spec.NetworkPlugin
	}
	if s.ControlPlane.Spec.NetworkPolicy != nil {
		managedClusterSpec.NetworkPolicy = *s.ControlPlane.Spec.NetworkPolicy
	}
	if s.ControlPlane.Spec.LoadBalancerSKU != nil {
		managedClusterSpec.LoadBalancerSKU = *s.ControlPlane.Spec.LoadBalancerSKU
	}

	if clusterNetwork := s.Cluster.Spec.ClusterNetwork; clusterNetwork != nil {
		if clusterNetwork.Services != nil && len(clusterNetwork.Services.CIDRBlocks) == 1 {
			managedClusterSpec.ServiceCIDR = clusterNetwork.Services.CIDRBlocks[0]
		}
		if clusterNetwork.Pods != nil && len(clusterNetwork.Pods.CIDRBlocks) == 1 {
			managedClusterSpec.PodCIDR = clusterNetwork.Pods.CIDRBlocks[0]
		}
	}

	if s.ControlPlane.Spec.AADProfile != nil {
		managedClusterSpec.AADProfile = &managedclusters.AADProfile{
			Managed:             s.ControlPlane.Spec.AADProfile.Managed,
			EnableAzureRBAC:     s.ControlPlane.Spec.AADProfile.Managed,
			AdminGroupObjectIDs: s.ControlPlane.Spec.AADProfile.AdminGroupObjectIDs,
		}
	}

	if s.ControlPlane.Spec.AddonProfiles != nil {
		for _, profile := range s.ControlPlane.Spec.AddonProfiles {
			managedClusterSpec.AddonProfiles = append(managedClusterSpec.AddonProfiles, managedclusters.AddonProfile{
				Name:    profile.Name,
				Enabled: profile.Enabled,
				Config:  profile.Config,
			})
		}
	}

	if s.ControlPlane.Spec.SKU != nil {
		managedClusterSpec.SKU = &managedclusters.SKU{
			Tier: string(s.ControlPlane.Spec.SKU.Tier),
		}
	}

	if s.ControlPlane.Spec.LoadBalancerProfile != nil {
		managedClusterSpec.LoadBalancerProfile = &managedclusters.LoadBalancerProfile{
			ManagedOutboundIPs:     s.ControlPlane.Spec.LoadBalancerProfile.ManagedOutboundIPs,
			OutboundIPPrefixes:     s.ControlPlane.Spec.LoadBalancerProfile.OutboundIPPrefixes,
			OutboundIPs:            s.ControlPlane.Spec.LoadBalancerProfile.OutboundIPs,
			AllocatedOutboundPorts: s.ControlPlane.Spec.LoadBalancerProfile.AllocatedOutboundPorts,
			IdleTimeoutInMinutes:   s.ControlPlane.Spec.LoadBalancerProfile.IdleTimeoutInMinutes,
		}
	}

	if s.ControlPlane.Spec.APIServerAccessProfile != nil {
		managedClusterSpec.APIServerAccessProfile = &managedclusters.APIServerAccessProfile{
			AuthorizedIPRanges:             s.ControlPlane.Spec.APIServerAccessProfile.AuthorizedIPRanges,
			EnablePrivateCluster:           s.ControlPlane.Spec.APIServerAccessProfile.EnablePrivateCluster,
			PrivateDNSZone:                 s.ControlPlane.Spec.APIServerAccessProfile.PrivateDNSZone,
			EnablePrivateClusterPublicFQDN: s.ControlPlane.Spec.APIServerAccessProfile.EnablePrivateClusterPublicFQDN,
		}
	}

	if s.ControlPlane.Spec.AutoScalerProfile != nil {
		managedClusterSpec.AutoScalerProfile = &managedclusters.AutoScalerProfile{
			BalanceSimilarNodeGroups:      (*string)(s.ControlPlane.Spec.AutoScalerProfile.BalanceSimilarNodeGroups),
			Expander:                      (*string)(s.ControlPlane.Spec.AutoScalerProfile.Expander),
			MaxEmptyBulkDelete:            s.ControlPlane.Spec.AutoScalerProfile.MaxEmptyBulkDelete,
			MaxGracefulTerminationSec:     s.ControlPlane.Spec.AutoScalerProfile.MaxGracefulTerminationSec,
			MaxNodeProvisionTime:          s.ControlPlane.Spec.AutoScalerProfile.MaxNodeProvisionTime,
			MaxTotalUnreadyPercentage:     s.ControlPlane.Spec.AutoScalerProfile.MaxTotalUnreadyPercentage,
			NewPodScaleUpDelay:            s.ControlPlane.Spec.AutoScalerProfile.NewPodScaleUpDelay,
			OkTotalUnreadyCount:           s.ControlPlane.Spec.AutoScalerProfile.OkTotalUnreadyCount,
			ScanInterval:                  s.ControlPlane.Spec.AutoScalerProfile.ScanInterval,
			ScaleDownDelayAfterAdd:        s.ControlPlane.Spec.AutoScalerProfile.ScaleDownDelayAfterAdd,
			ScaleDownDelayAfterDelete:     s.ControlPlane.Spec.AutoScalerProfile.ScaleDownDelayAfterDelete,
			ScaleDownDelayAfterFailure:    s.ControlPlane.Spec.AutoScalerProfile.ScaleDownDelayAfterFailure,
			ScaleDownUnneededTime:         s.ControlPlane.Spec.AutoScalerProfile.ScaleDownUnneededTime,
			ScaleDownUnreadyTime:          s.ControlPlane.Spec.AutoScalerProfile.ScaleDownUnreadyTime,
			ScaleDownUtilizationThreshold: s.ControlPlane.Spec.AutoScalerProfile.ScaleDownUtilizationThreshold,
			SkipNodesWithLocalStorage:     (*string)(s.ControlPlane.Spec.AutoScalerProfile.SkipNodesWithLocalStorage),
			SkipNodesWithSystemPods:       (*string)(s.ControlPlane.Spec.AutoScalerProfile.SkipNodesWithSystemPods),
		}
	}

	if s.ControlPlane.Spec.HTTPProxyConfig != nil {
		managedClusterSpec.HTTPProxyConfig = &managedclusters.HTTPProxyConfig{
			HTTPProxy:  s.ControlPlane.Spec.HTTPProxyConfig.HTTPProxy,
			HTTPSProxy: s.ControlPlane.Spec.HTTPProxyConfig.HTTPSProxy,
			NoProxy:    s.ControlPlane.Spec.HTTPProxyConfig.NoProxy,
			TrustedCA:  s.ControlPlane.Spec.HTTPProxyConfig.TrustedCA,
		}
	}

	if s.ControlPlane.Spec.OIDCIssuerProfile != nil {
		managedClusterSpec.OIDCIssuerProfile = &managedclusters.OIDCIssuerProfile{
			Enabled: s.ControlPlane.Spec.OIDCIssuerProfile.Enabled,
		}
	}

	return &managedClusterSpec
}

// GetAllAgentPoolSpecs gets a slice of azure.AgentPoolSpec for the list of agent pools.
func (s *ManagedControlPlaneScope) GetAllAgentPoolSpecs() ([]azure.ResourceSpecGetter, error) {
	var (
		ammps           = make([]azure.ResourceSpecGetter, 0, len(s.ManagedMachinePools))
		foundSystemPool = false
	)
	for _, pool := range s.ManagedMachinePools {
		// TODO: this should be in a webhook: https://github.com/kubernetes-sigs/cluster-api/issues/6040
		if pool.MachinePool != nil && pool.MachinePool.Spec.Template.Spec.Version != nil {
			version := *pool.MachinePool.Spec.Template.Spec.Version
			if semver.Compare(version, s.ControlPlane.Spec.Version) > 0 {
				return nil, errors.New("MachinePool version cannot be greater than the AzureManagedControlPlane version")
			}
		}

		if pool.InfraMachinePool != nil && pool.InfraMachinePool.Spec.Mode == string(infrav1.NodePoolModeSystem) {
			foundSystemPool = true
		}

		ammp := buildAgentPoolSpec(s.ControlPlane, pool.MachinePool, pool.InfraMachinePool, pool.InfraMachinePool.Annotations)
		ammps = append(ammps, ammp)
	}

	if !foundSystemPool {
		return nil, errors.New("failed to fetch azuremanagedMachine pool with mode:System, require at least 1 system node pool")
	}

	return ammps, nil
}

// SetControlPlaneEndpoint sets a control plane endpoint.
func (s *ManagedControlPlaneScope) SetControlPlaneEndpoint(endpoint clusterv1.APIEndpoint) {
	s.ControlPlane.Spec.ControlPlaneEndpoint.Host = endpoint.Host
	s.ControlPlane.Spec.ControlPlaneEndpoint.Port = endpoint.Port
}

// MakeEmptyKubeConfigSecret creates an empty secret object that is used for storing kubeconfig secret data.
func (s *ManagedControlPlaneScope) MakeEmptyKubeConfigSecret() corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name(s.Cluster.Name, secret.Kubeconfig),
			Namespace: s.Cluster.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(s.ControlPlane, infrav1.GroupVersion.WithKind("AzureManagedControlPlane")),
			},
			Labels: map[string]string{clusterv1.ClusterNameLabel: s.Cluster.Name},
		},
	}
}

// GetKubeConfigData returns a []byte that contains kubeconfig.
func (s *ManagedControlPlaneScope) GetKubeConfigData() []byte {
	return s.kubeConfigData
}

// SetKubeConfigData sets kubeconfig data.
func (s *ManagedControlPlaneScope) SetKubeConfigData(kubeConfigData []byte) {
	s.kubeConfigData = kubeConfigData
}

// SetKubeletIdentity sets the ID of the user-assigned identity for kubelet if not already set.
func (s *ManagedControlPlaneScope) SetKubeletIdentity(id string) {
	s.ControlPlane.Spec.KubeletUserAssignedIdentity = id
}

// SetLongRunningOperationState will set the future on the AzureManagedControlPlane status to allow the resource to continue
// in the next reconciliation.
func (s *ManagedControlPlaneScope) SetLongRunningOperationState(future *infrav1.Future) {
	futures.Set(s.ControlPlane, future)
}

// GetLongRunningOperationState will get the future on the AzureManagedControlPlane status.
func (s *ManagedControlPlaneScope) GetLongRunningOperationState(name, service, futureType string) *infrav1.Future {
	return futures.Get(s.ControlPlane, name, service, futureType)
}

// DeleteLongRunningOperationState will delete the future from the AzureManagedControlPlane status.
func (s *ManagedControlPlaneScope) DeleteLongRunningOperationState(name, service, futureType string) {
	futures.Delete(s.ControlPlane, name, service, futureType)
}

// UpdateDeleteStatus updates a condition on the AzureManagedControlPlane status after a DELETE operation.
func (s *ManagedControlPlaneScope) UpdateDeleteStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkFalse(s.ControlPlane, condition, infrav1.DeletedReason, clusterv1.ConditionSeverityInfo, "%s successfully deleted", service)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(s.ControlPlane, condition, infrav1.DeletingReason, clusterv1.ConditionSeverityInfo, "%s deleting", service)
	default:
		conditions.MarkFalse(s.ControlPlane, condition, infrav1.DeletionFailedReason, clusterv1.ConditionSeverityError, "%s failed to delete. err: %s", service, err.Error())
	}
}

// UpdatePutStatus updates a condition on the AzureManagedControlPlane status after a PUT operation.
func (s *ManagedControlPlaneScope) UpdatePutStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkTrue(s.ControlPlane, condition)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(s.ControlPlane, condition, infrav1.CreatingReason, clusterv1.ConditionSeverityInfo, "%s creating or updating", service)
	default:
		conditions.MarkFalse(s.ControlPlane, condition, infrav1.FailedReason, clusterv1.ConditionSeverityError, "%s failed to create or update. err: %s", service, err.Error())
	}
}

// UpdatePatchStatus updates a condition on the AzureManagedControlPlane status after a PATCH operation.
func (s *ManagedControlPlaneScope) UpdatePatchStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkTrue(s.ControlPlane, condition)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(s.ControlPlane, condition, infrav1.UpdatingReason, clusterv1.ConditionSeverityInfo, "%s updating", service)
	default:
		conditions.MarkFalse(s.ControlPlane, condition, infrav1.FailedReason, clusterv1.ConditionSeverityError, "%s failed to update. err: %s", service, err.Error())
	}
}

// AnnotationJSON returns a map[string]interface from a JSON annotation.
func (s *ManagedControlPlaneScope) AnnotationJSON(annotation string) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	jsonAnnotation := s.ControlPlane.GetAnnotations()[annotation]
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
func (s *ManagedControlPlaneScope) UpdateAnnotationJSON(annotation string, content map[string]interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}
	s.SetAnnotation(annotation, string(b))
	return nil
}

// SetAnnotation sets a key value annotation on the ControlPlane.
func (s *ManagedControlPlaneScope) SetAnnotation(key, value string) {
	if s.ControlPlane.Annotations == nil {
		s.ControlPlane.Annotations = map[string]string{}
	}
	s.ControlPlane.Annotations[key] = value
}

// TagsSpecs returns the tag specs for the ManagedControlPlane.
func (s *ManagedControlPlaneScope) TagsSpecs() []azure.TagsSpec {
	specs := []azure.TagsSpec{
		{
			Scope:      azure.ManagedClusterID(s.SubscriptionID(), s.ResourceGroup(), s.ManagedClusterSpec().ResourceName()),
			Tags:       s.AdditionalTags(),
			Annotation: azure.ManagedClusterTagsLastAppliedAnnotation,
		},
	}
	return specs
}

// AvailabilityStatusResource refers to the AzureManagedControlPlane.
func (s *ManagedControlPlaneScope) AvailabilityStatusResource() conditions.Setter {
	return s.ControlPlane
}

// AvailabilityStatusResourceURI constructs the ID of the underlying AKS resource.
func (s *ManagedControlPlaneScope) AvailabilityStatusResourceURI() string {
	return azure.ManagedClusterID(s.SubscriptionID(), s.ResourceGroup(), s.ControlPlane.Name)
}

// AvailabilityStatusFilter ignores the health metrics connection error that
// occurs on startup for every AKS cluster.
func (s *ManagedControlPlaneScope) AvailabilityStatusFilter(cond *clusterv1.Condition) *clusterv1.Condition {
	if time.Since(s.ControlPlane.CreationTimestamp.Time) < resourceHealthWarningInitialGracePeriod &&
		cond.Severity == clusterv1.ConditionSeverityWarning {
		return conditions.TrueCondition(infrav1.AzureResourceAvailableCondition)
	}
	return cond
}

// PrivateEndpointSpecs returns the private endpoint specs.
func (s *ManagedControlPlaneScope) PrivateEndpointSpecs() []azure.ResourceSpecGetter {
	privateEndpointSpecs := make([]azure.ResourceSpecGetter, len(s.ControlPlane.Spec.VirtualNetwork.Subnet.PrivateEndpoints))

	for _, privateEndpoint := range s.ControlPlane.Spec.VirtualNetwork.Subnet.PrivateEndpoints {
		privateEndpointSpec := &privateendpoints.PrivateEndpointSpec{
			Name:                       privateEndpoint.Name,
			ResourceGroup:              s.VNetSpec().ResourceGroupName(),
			Location:                   privateEndpoint.Location,
			CustomNetworkInterfaceName: privateEndpoint.CustomNetworkInterfaceName,
			PrivateIPAddresses:         privateEndpoint.PrivateIPAddresses,
			SubnetID: azure.SubnetID(
				s.ControlPlane.Spec.SubscriptionID,
				s.VNetSpec().ResourceGroupName(),
				s.ControlPlane.Spec.VirtualNetwork.Name,
				s.ControlPlane.Spec.VirtualNetwork.Subnet.Name,
			),
			ApplicationSecurityGroups: privateEndpoint.ApplicationSecurityGroups,
			ManualApproval:            privateEndpoint.ManualApproval,
			ClusterName:               s.ClusterName(),
			AdditionalTags:            s.AdditionalTags(),
		}

		for _, privateLinkServiceConnection := range privateEndpoint.PrivateLinkServiceConnections {
			pl := privateendpoints.PrivateLinkServiceConnection{
				PrivateLinkServiceID: privateLinkServiceConnection.PrivateLinkServiceID,
				Name:                 privateLinkServiceConnection.Name,
				RequestMessage:       privateLinkServiceConnection.RequestMessage,
				GroupIDs:             privateLinkServiceConnection.GroupIDs,
			}
			privateEndpointSpec.PrivateLinkServiceConnections = append(privateEndpointSpec.PrivateLinkServiceConnections, pl)
		}

		privateEndpointSpecs = append(privateEndpointSpecs, privateEndpointSpec)
	}

	return privateEndpointSpecs
}

// SetOIDCIssuerProfileStatus sets the status for the OIDC issuer profile config.
func (s *ManagedControlPlaneScope) SetOIDCIssuerProfileStatus(oidc *infrav1.OIDCIssuerProfileStatus) {
	s.ControlPlane.Status.OIDCIssuerProfile = oidc
}
