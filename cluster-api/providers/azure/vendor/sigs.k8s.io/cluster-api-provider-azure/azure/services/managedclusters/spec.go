/*
Copyright 2022 The Kubernetes Authors.

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

package managedclusters

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// ManagedClusterSpec contains properties to create a managed cluster.
type ManagedClusterSpec struct {
	// Name is the name of this AKS Cluster.
	Name string

	// ResourceGroup is the name of the Azure resource group for this AKS Cluster.
	ResourceGroup string

	// NodeResourceGroup is the name of the Azure resource group containing IaaS VMs.
	NodeResourceGroup string

	// ClusterName is the name of the owning Cluster API Cluster resource.
	ClusterName string

	// VnetSubnetID is the Azure Resource ID for the subnet which should contain nodes.
	VnetSubnetID string

	// Location is a string matching one of the canonical Azure region names. Examples: "westus2", "eastus".
	Location string

	// Tags is a set of tags to add to this cluster.
	Tags map[string]string

	// Version defines the desired Kubernetes version.
	Version string

	// LoadBalancerSKU for the managed cluster. Possible values include: 'Standard', 'Basic'. Defaults to Standard.
	LoadBalancerSKU string

	// NetworkPlugin used for building Kubernetes network. Possible values include: 'azure', 'kubenet'. Defaults to azure.
	NetworkPlugin string

	// NetworkPluginMode is the mode the network plugin should use.
	NetworkPluginMode *infrav1.NetworkPluginMode

	// NetworkPolicy used for building Kubernetes network. Possible values include: 'calico', 'azure'.
	NetworkPolicy string

	// OutboundType used for building Kubernetes network. Possible values include: 'loadBalancer', 'managedNATGateway', 'userAssignedNATGateway', 'userDefinedRouting'.
	OutboundType *infrav1.ManagedControlPlaneOutboundType

	// SSHPublicKey is a string literal containing an ssh public key. Will autogenerate and discard if not provided.
	SSHPublicKey string

	// GetAllAgentPools is a function that returns the list of agent pool specifications in this cluster.
	GetAllAgentPools func() ([]azure.ResourceSpecGetter, error)

	// PodCIDR is the CIDR block for IP addresses distributed to pods
	PodCIDR string

	// ServiceCIDR is the CIDR block for IP addresses distributed to services
	ServiceCIDR string

	// DNSServiceIP is an IP address assigned to the Kubernetes DNS service
	DNSServiceIP *string

	// AddonProfiles are the profiles of managed cluster add-on.
	AddonProfiles []AddonProfile

	// AADProfile is Azure Active Directory configuration to integrate with AKS, for aad authentication.
	AADProfile *AADProfile

	// SKU is the SKU of the AKS to be provisioned.
	SKU *SKU

	// LoadBalancerProfile is the profile of the cluster load balancer.
	LoadBalancerProfile *LoadBalancerProfile

	// APIServerAccessProfile is the access profile for AKS API server.
	APIServerAccessProfile *APIServerAccessProfile

	// Headers is the list of headers to add to the HTTP requests to update this resource.
	Headers map[string]string

	// AutoScalerProfile is the parameters to be applied to the cluster-autoscaler when enabled.
	AutoScalerProfile *AutoScalerProfile

	// Identity is the AKS control plane Identity configuration
	Identity *infrav1.Identity

	// KubeletUserAssignedIdentity is the user-assigned identity for kubelet to authenticate to ACR.
	KubeletUserAssignedIdentity string

	// HTTPProxyConfig is the HTTP proxy configuration for the cluster.
	HTTPProxyConfig *HTTPProxyConfig

	// OIDCIssuerProfile is the OIDC issuer profile of the Managed Cluster.
	OIDCIssuerProfile *OIDCIssuerProfile

	// DNSPrefix allows the user to customize dns prefix.
	DNSPrefix *string
	// DisableLocalAccounts disables getting static credentials for this cluster when set. Expected to only be used for AAD clusters.
	DisableLocalAccounts *bool
}

// HTTPProxyConfig is the HTTP proxy configuration for the cluster.
type HTTPProxyConfig struct {
	// HTTPProxy is the HTTP proxy server endpoint to use.
	HTTPProxy *string `json:"httpProxy,omitempty"`

	// HTTPSProxy is the HTTPS proxy server endpoint to use.
	HTTPSProxy *string `json:"httpsProxy,omitempty"`

	// NoProxy is the endpoints that should not go through proxy.
	NoProxy []string `json:"noProxy,omitempty"`

	// TrustedCA is the Alternative CA cert to use for connecting to proxy servers.
	TrustedCA *string `json:"trustedCa,omitempty"`
}

// AADProfile is Azure Active Directory configuration to integrate with AKS, for aad authentication.
type AADProfile struct {
	// Managed defines whether to enable managed AAD.
	Managed bool

	// EnableAzureRBAC defines whether to enable Azure RBAC for Kubernetes authorization.
	EnableAzureRBAC bool

	// AdminGroupObjectIDs are the AAD group object IDs that will have admin role of the cluster.
	AdminGroupObjectIDs []string
}

// AddonProfile is the profile of a managed cluster add-on.
type AddonProfile struct {
	Name    string
	Config  map[string]string
	Enabled bool
}

// SKU is an AKS SKU.
type SKU struct {
	// Tier is the tier of a managed cluster SKU.
	Tier string
}

// LoadBalancerProfile is the profile of the cluster load balancer.
type LoadBalancerProfile struct {
	// Load balancer profile must specify at most one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs.
	// By default the AKS cluster automatically creates a public IP in the AKS-managed infrastructure resource group and assigns it to the load balancer outbound pool.
	// Alternatively, you can assign your own custom public IP or public IP prefix at cluster creation time.
	// See https://learn.microsoft.com/azure/aks/load-balancer-standard#provide-your-own-outbound-public-ips-or-prefixes

	// ManagedOutboundIPs are the desired managed outbound IPs for the cluster load balancer.
	ManagedOutboundIPs *int32

	// OutboundIPPrefixes are the desired outbound IP Prefix resources for the cluster load balancer.
	OutboundIPPrefixes []string

	// OutboundIPs are the desired outbound IP resources for the cluster load balancer.
	OutboundIPs []string

	// AllocatedOutboundPorts are the desired number of allocated SNAT ports per VM. Allowed values must be in the range of 0 to 64000 (inclusive). The default value is 0 which results in Azure dynamically allocating ports.
	AllocatedOutboundPorts *int32

	// IdleTimeoutInMinutes  are the desired outbound flow idle timeout in minutes. Allowed values must be in the range of 4 to 120 (inclusive). The default value is 30 minutes.
	IdleTimeoutInMinutes *int32
}

// APIServerAccessProfile is the access profile for AKS API server.
type APIServerAccessProfile struct {
	// AuthorizedIPRanges are the authorized IP Ranges to kubernetes API server.
	AuthorizedIPRanges []string
	// EnablePrivateCluster defines hether to create the cluster as a private cluster or not.
	EnablePrivateCluster *bool
	// PrivateDNSZone is the private dns zone for private clusters.
	PrivateDNSZone *string
	// EnablePrivateClusterPublicFQDN defines whether to create additional public FQDN for private cluster or not.
	EnablePrivateClusterPublicFQDN *bool
}

// AutoScalerProfile parameters to be applied to the cluster-autoscaler when enabled.
type AutoScalerProfile struct {
	// BalanceSimilarNodeGroups - Valid values are 'true' and 'false'
	BalanceSimilarNodeGroups *string
	// Expander - If not specified, the default is 'random'. See [expanders](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-are-expanders) for more information.
	Expander *string
	// MaxEmptyBulkDelete - The default is 10.
	MaxEmptyBulkDelete *string
	// MaxGracefulTerminationSec - The default is 600.
	MaxGracefulTerminationSec *string
	// MaxNodeProvisionTime - The default is '15m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	MaxNodeProvisionTime *string
	// MaxTotalUnreadyPercentage - The default is 45. The maximum is 100 and the minimum is 0.
	MaxTotalUnreadyPercentage *string
	// NewPodScaleUpDelay - For scenarios like burst/batch scale where you don't want CA to act before the kubernetes scheduler could schedule all the pods, you can tell CA to ignore unscheduled pods before they're a certain age. The default is '0s'. Values must be an integer followed by a unit ('s' for seconds, 'm' for minutes, 'h' for hours, etc).
	NewPodScaleUpDelay *string
	// OkTotalUnreadyCount - This must be an integer. The default is 3.
	OkTotalUnreadyCount *string
	// ScanInterval - The default is '10s'. Values must be an integer number of seconds.
	ScanInterval *string
	// ScaleDownDelayAfterAdd - The default is '10m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	ScaleDownDelayAfterAdd *string
	// ScaleDownDelayAfterDelete - The default is the scan-interval. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	ScaleDownDelayAfterDelete *string
	// ScaleDownDelayAfterFailure - The default is '3m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	ScaleDownDelayAfterFailure *string
	// ScaleDownUnneededTime - The default is '10m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	ScaleDownUnneededTime *string
	// ScaleDownUnreadyTime - The default is '20m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	ScaleDownUnreadyTime *string
	// ScaleDownUtilizationThreshold - The default is '0.5'.
	ScaleDownUtilizationThreshold *string
	// SkipNodesWithLocalStorage - The default is true.
	SkipNodesWithLocalStorage *string
	// SkipNodesWithSystemPods - The default is true.
	SkipNodesWithSystemPods *string
}

// OIDCIssuerProfile is the OIDC issuer profile of the Managed Cluster.
type OIDCIssuerProfile struct {
	// Enabled is whether the OIDC issuer is enabled.
	Enabled *bool
}

var _ azure.ResourceSpecGetterWithHeaders = (*ManagedClusterSpec)(nil)

// ResourceName returns the name of the AKS cluster.
func (s *ManagedClusterSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *ManagedClusterSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for managed clusters.
func (s *ManagedClusterSpec) OwnerResourceName() string {
	return "" // not applicable
}

// CustomHeaders returns custom headers to be added to the Azure API calls.
func (s *ManagedClusterSpec) CustomHeaders() map[string]string {
	return s.Headers
}

// buildAutoScalerProfile builds the AutoScalerProfile for the ManagedClusterProperties.
func buildAutoScalerProfile(autoScalerProfile *AutoScalerProfile) *armcontainerservice.ManagedClusterPropertiesAutoScalerProfile {
	if autoScalerProfile == nil {
		return nil
	}

	mcAutoScalerProfile := &armcontainerservice.ManagedClusterPropertiesAutoScalerProfile{
		BalanceSimilarNodeGroups:      autoScalerProfile.BalanceSimilarNodeGroups,
		MaxEmptyBulkDelete:            autoScalerProfile.MaxEmptyBulkDelete,
		MaxGracefulTerminationSec:     autoScalerProfile.MaxGracefulTerminationSec,
		MaxNodeProvisionTime:          autoScalerProfile.MaxNodeProvisionTime,
		MaxTotalUnreadyPercentage:     autoScalerProfile.MaxTotalUnreadyPercentage,
		NewPodScaleUpDelay:            autoScalerProfile.NewPodScaleUpDelay,
		OkTotalUnreadyCount:           autoScalerProfile.OkTotalUnreadyCount,
		ScanInterval:                  autoScalerProfile.ScanInterval,
		ScaleDownDelayAfterAdd:        autoScalerProfile.ScaleDownDelayAfterAdd,
		ScaleDownDelayAfterDelete:     autoScalerProfile.ScaleDownDelayAfterDelete,
		ScaleDownDelayAfterFailure:    autoScalerProfile.ScaleDownDelayAfterFailure,
		ScaleDownUnneededTime:         autoScalerProfile.ScaleDownUnneededTime,
		ScaleDownUnreadyTime:          autoScalerProfile.ScaleDownUnreadyTime,
		ScaleDownUtilizationThreshold: autoScalerProfile.ScaleDownUtilizationThreshold,
		SkipNodesWithLocalStorage:     autoScalerProfile.SkipNodesWithLocalStorage,
		SkipNodesWithSystemPods:       autoScalerProfile.SkipNodesWithSystemPods,
	}
	if autoScalerProfile.Expander != nil {
		mcAutoScalerProfile.Expander = ptr.To(armcontainerservice.Expander(*autoScalerProfile.Expander))
	}

	return mcAutoScalerProfile
}

// Parameters returns the parameters for the managed clusters.
//
//nolint:gocyclo // Function requires a lot of nil checks that raise complexity.
func (s *ManagedClusterSpec) Parameters(ctx context.Context, existing interface{}) (params interface{}, err error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "managedclusters.Service.Parameters")
	defer done()

	var decodedSSHPublicKey []byte
	if s.SSHPublicKey != "" {
		decodedSSHPublicKey, err = base64.StdEncoding.DecodeString(s.SSHPublicKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode SSHPublicKey")
		}
	}

	managedCluster := armcontainerservice.ManagedCluster{
		Identity: &armcontainerservice.ManagedClusterIdentity{
			Type: ptr.To(armcontainerservice.ResourceIdentityTypeSystemAssigned),
		},
		Location: &s.Location,
		Tags: converters.TagsToMap(infrav1.Build(infrav1.BuildParams{
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			ClusterName: s.ClusterName,
			Name:        ptr.To(s.Name),
			Role:        ptr.To(infrav1.CommonRole),
			Additional:  s.Tags,
		})),
		Properties: &armcontainerservice.ManagedClusterProperties{
			NodeResourceGroup: &s.NodeResourceGroup,
			EnableRBAC:        ptr.To(true),
			DNSPrefix:         s.DNSPrefix,
			KubernetesVersion: &s.Version,

			ServicePrincipalProfile: &armcontainerservice.ManagedClusterServicePrincipalProfile{
				ClientID: ptr.To("msi"),
			},
			AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{},
			NetworkProfile: &armcontainerservice.NetworkProfile{
				NetworkPlugin:   azure.AliasOrNil[armcontainerservice.NetworkPlugin](&s.NetworkPlugin),
				LoadBalancerSKU: azure.AliasOrNil[armcontainerservice.LoadBalancerSKU](&s.LoadBalancerSKU),
				NetworkPolicy:   azure.AliasOrNil[armcontainerservice.NetworkPolicy](&s.NetworkPolicy),
			},
		},
	}

	if decodedSSHPublicKey != nil {
		managedCluster.Properties.LinuxProfile = &armcontainerservice.LinuxProfile{
			AdminUsername: ptr.To(azure.DefaultAKSUserName),
			SSH: &armcontainerservice.SSHConfiguration{
				PublicKeys: []*armcontainerservice.SSHPublicKey{
					{
						KeyData: ptr.To(string(decodedSSHPublicKey)),
					},
				},
			},
		}
	}

	if s.NetworkPluginMode != nil {
		managedCluster.Properties.NetworkProfile.NetworkPluginMode = ptr.To(armcontainerservice.NetworkPluginMode(*s.NetworkPluginMode))
	}

	if s.PodCIDR != "" {
		managedCluster.Properties.NetworkProfile.PodCidr = &s.PodCIDR
	}

	if s.ServiceCIDR != "" {
		if s.DNSServiceIP == nil {
			managedCluster.Properties.NetworkProfile.ServiceCidr = &s.ServiceCIDR
			ip, _, err := net.ParseCIDR(s.ServiceCIDR)
			if err != nil {
				return nil, fmt.Errorf("failed to parse service cidr: %w", err)
			}
			// HACK: set the last octet of the IP to .10
			// This ensures the dns IP is valid in the service cidr without forcing the user
			// to specify it in both the Capi cluster and the Azure control plane.
			// https://golang.org/src/net/ip.go#L48
			ip[15] = byte(10)
			dnsIP := ip.String()
			managedCluster.Properties.NetworkProfile.DNSServiceIP = &dnsIP
		} else {
			managedCluster.Properties.NetworkProfile.DNSServiceIP = s.DNSServiceIP
		}
	}

	if s.AADProfile != nil {
		managedCluster.Properties.AADProfile = &armcontainerservice.ManagedClusterAADProfile{
			Managed:             &s.AADProfile.Managed,
			EnableAzureRBAC:     &s.AADProfile.EnableAzureRBAC,
			AdminGroupObjectIDs: azure.PtrSlice[string](&s.AADProfile.AdminGroupObjectIDs),
		}
		if s.DisableLocalAccounts != nil {
			managedCluster.Properties.DisableLocalAccounts = s.DisableLocalAccounts
		}
	}

	for i := range s.AddonProfiles {
		if managedCluster.Properties.AddonProfiles == nil {
			managedCluster.Properties.AddonProfiles = map[string]*armcontainerservice.ManagedClusterAddonProfile{}
		}
		item := s.AddonProfiles[i]
		addonProfile := &armcontainerservice.ManagedClusterAddonProfile{
			Enabled: &item.Enabled,
		}
		if item.Config != nil {
			addonProfile.Config = azure.StringMapPtr(item.Config)
		}
		managedCluster.Properties.AddonProfiles[item.Name] = addonProfile
	}

	if s.SKU != nil {
		tierName := armcontainerservice.ManagedClusterSKUTier(s.SKU.Tier)
		managedCluster.SKU = &armcontainerservice.ManagedClusterSKU{
			Name: ptr.To(armcontainerservice.ManagedClusterSKUName("Base")),
			Tier: ptr.To(tierName),
		}
	}

	if s.LoadBalancerProfile != nil {
		managedCluster.Properties.NetworkProfile.LoadBalancerProfile = s.GetLoadBalancerProfile()
	}

	if s.APIServerAccessProfile != nil {
		managedCluster.Properties.APIServerAccessProfile = &armcontainerservice.ManagedClusterAPIServerAccessProfile{
			EnablePrivateCluster:           s.APIServerAccessProfile.EnablePrivateCluster,
			PrivateDNSZone:                 s.APIServerAccessProfile.PrivateDNSZone,
			EnablePrivateClusterPublicFQDN: s.APIServerAccessProfile.EnablePrivateClusterPublicFQDN,
		}

		if s.APIServerAccessProfile.AuthorizedIPRanges != nil {
			managedCluster.Properties.APIServerAccessProfile.AuthorizedIPRanges = azure.PtrSlice[string](&s.APIServerAccessProfile.AuthorizedIPRanges)
		}
	}

	if s.OutboundType != nil {
		managedCluster.Properties.NetworkProfile.OutboundType = ptr.To(armcontainerservice.OutboundType(*s.OutboundType))
	}

	managedCluster.Properties.AutoScalerProfile = buildAutoScalerProfile(s.AutoScalerProfile)

	if s.Identity != nil {
		managedCluster.Identity, err = getIdentity(s.Identity)
		if err != nil {
			return nil, errors.Wrapf(err, "Identity is not valid: %s", err)
		}
	}

	if s.KubeletUserAssignedIdentity != "" {
		managedCluster.Properties.IdentityProfile = map[string]*armcontainerservice.UserAssignedIdentity{
			kubeletIdentityKey: {
				ResourceID: ptr.To(s.KubeletUserAssignedIdentity),
			},
		}
	}

	if s.HTTPProxyConfig != nil {
		managedCluster.Properties.HTTPProxyConfig = &armcontainerservice.ManagedClusterHTTPProxyConfig{
			HTTPProxy:  s.HTTPProxyConfig.HTTPProxy,
			HTTPSProxy: s.HTTPProxyConfig.HTTPSProxy,
			TrustedCa:  s.HTTPProxyConfig.TrustedCA,
		}

		if s.HTTPProxyConfig.NoProxy != nil {
			managedCluster.Properties.HTTPProxyConfig.NoProxy = azure.PtrSlice(&s.HTTPProxyConfig.NoProxy)
		}
	}

	if s.OIDCIssuerProfile != nil {
		managedCluster.Properties.OidcIssuerProfile = &armcontainerservice.ManagedClusterOIDCIssuerProfile{
			Enabled: s.OIDCIssuerProfile.Enabled,
		}
	}

	if existing != nil {
		existingMC, ok := existing.(armcontainerservice.ManagedCluster)
		if !ok {
			return nil, fmt.Errorf("%T is not an armcontainerservice.ManagedCluster", existing)
		}
		ps := *existingMC.Properties.ProvisioningState
		if ps != string(infrav1.Canceled) && ps != string(infrav1.Failed) && ps != string(infrav1.Succeeded) {
			return nil, azure.WithTransientError(errors.Errorf("Unable to update existing managed cluster in non-terminal state. Managed cluster must be in one of the following provisioning states: Canceled, Failed, or Succeeded. Actual state: %s", ps), 20*time.Second)
		}

		// Normalize the LoadBalancerProfile so the diff below doesn't get thrown off by AKS added properties.
		if managedCluster.Properties.NetworkProfile.LoadBalancerProfile == nil {
			// If our LoadBalancerProfile generated by the spec is nil, then don't worry about what AKS has added.
			existingMC.Properties.NetworkProfile.LoadBalancerProfile = nil
		} else {
			// If our LoadBalancerProfile generated by the spec is not nil, then remove the effective outbound IPs from
			// AKS.
			existingMC.Properties.NetworkProfile.LoadBalancerProfile.EffectiveOutboundIPs = nil
		}

		// Avoid changing agent pool profiles through AMCP and just use the existing agent pool profiles
		// AgentPool changes are managed through AMMP.
		managedCluster.Properties.AgentPoolProfiles = existingMC.Properties.AgentPoolProfiles

		// if the AuthorizedIPRanges is nil in the user-updated spec, but not nil in the existing spec, then
		// we need to set the AuthorizedIPRanges to empty array ([]*string{}) once so that the Azure API will
		// update the existing authorized IP ranges to nil.
		if !isAuthIPRangesNilOrEmpty(existingMC) && isAuthIPRangesNilOrEmpty(managedCluster) {
			log.V(4).Info("managed cluster spec has nil AuthorizedIPRanges, updating existing authorized IP ranges to an empty list")
			managedCluster.Properties.APIServerAccessProfile = &armcontainerservice.ManagedClusterAPIServerAccessProfile{
				AuthorizedIPRanges: []*string{},
			}
		}

		diff := computeDiffOfNormalizedClusters(managedCluster, existingMC)
		if diff == "" {
			log.V(4).Info("no changes found between user-updated spec and existing spec")
			return nil, nil
		}
		log.V(4).Info("found a diff between the desired spec and the existing managed cluster", "difference", diff)
	} else {
		// Add all agent pools to cluster spec that will be submitted to the API
		agentPoolSpecs, err := s.GetAllAgentPools()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get agent pool specs for managed cluster %s", s.Name)
		}

		for _, spec := range agentPoolSpecs {
			params, err := spec.Parameters(ctx, nil)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get agent pool parameters for managed cluster %s", s.Name)
			}
			agentPool, ok := params.(armcontainerservice.AgentPool)
			if !ok {
				return nil, fmt.Errorf("%T is not an armcontainerservice.AgentPool", agentPool)
			}
			agentPool.Name = ptr.To(spec.ResourceName())
			profile := converters.AgentPoolToManagedClusterAgentPoolProfile(agentPool)
			managedCluster.Properties.AgentPoolProfiles = append(managedCluster.Properties.AgentPoolProfiles, &profile)
		}
	}

	return managedCluster, nil
}

// GetLoadBalancerProfile returns an armcontainerservice.ManagedClusterLoadBalancerProfile from the
// information present in ManagedClusterSpec.LoadBalancerProfile.
func (s *ManagedClusterSpec) GetLoadBalancerProfile() (loadBalancerProfile *armcontainerservice.ManagedClusterLoadBalancerProfile) {
	loadBalancerProfile = &armcontainerservice.ManagedClusterLoadBalancerProfile{
		AllocatedOutboundPorts: s.LoadBalancerProfile.AllocatedOutboundPorts,
		IdleTimeoutInMinutes:   s.LoadBalancerProfile.IdleTimeoutInMinutes,
	}
	if s.LoadBalancerProfile.ManagedOutboundIPs != nil {
		loadBalancerProfile.ManagedOutboundIPs = &armcontainerservice.ManagedClusterLoadBalancerProfileManagedOutboundIPs{Count: s.LoadBalancerProfile.ManagedOutboundIPs}
	}
	if len(s.LoadBalancerProfile.OutboundIPPrefixes) > 0 {
		loadBalancerProfile.OutboundIPPrefixes = &armcontainerservice.ManagedClusterLoadBalancerProfileOutboundIPPrefixes{
			PublicIPPrefixes: convertToResourceReferences(s.LoadBalancerProfile.OutboundIPPrefixes),
		}
	}
	if len(s.LoadBalancerProfile.OutboundIPs) > 0 {
		loadBalancerProfile.OutboundIPs = &armcontainerservice.ManagedClusterLoadBalancerProfileOutboundIPs{
			PublicIPs: convertToResourceReferences(s.LoadBalancerProfile.OutboundIPs),
		}
	}
	return
}

func convertToResourceReferences(resources []string) []*armcontainerservice.ResourceReference {
	resourceReferences := make([]*armcontainerservice.ResourceReference, len(resources))
	for i := range resources {
		resourceReferences[i] = &armcontainerservice.ResourceReference{ID: &resources[i]}
	}
	return resourceReferences
}

func computeDiffOfNormalizedClusters(managedCluster armcontainerservice.ManagedCluster, existingMC armcontainerservice.ManagedCluster) string {
	// Normalize properties for the desired (CR spec) and existing managed
	// cluster, so that we check only those fields that were specified in
	// the initial CreateOrUpdate request and that can be modified.
	// Without comparing to normalized properties, we would always get a
	// difference in desired and existing, which would result in sending
	// unnecessary Azure API requests.
	propertiesNormalized := &armcontainerservice.ManagedClusterProperties{
		KubernetesVersion: managedCluster.Properties.KubernetesVersion,
		NetworkProfile:    &armcontainerservice.NetworkProfile{},
		AutoScalerProfile: &armcontainerservice.ManagedClusterPropertiesAutoScalerProfile{},
	}

	existingMCPropertiesNormalized := &armcontainerservice.ManagedClusterProperties{
		KubernetesVersion: existingMC.Properties.KubernetesVersion,
		NetworkProfile:    &armcontainerservice.NetworkProfile{},
		AutoScalerProfile: &armcontainerservice.ManagedClusterPropertiesAutoScalerProfile{},
	}

	if managedCluster.Properties.AADProfile != nil {
		propertiesNormalized.AADProfile = &armcontainerservice.ManagedClusterAADProfile{
			Managed:             managedCluster.Properties.AADProfile.Managed,
			EnableAzureRBAC:     managedCluster.Properties.AADProfile.EnableAzureRBAC,
			AdminGroupObjectIDs: managedCluster.Properties.AADProfile.AdminGroupObjectIDs,
		}
	}

	if existingMC.Properties.AADProfile != nil {
		existingMCPropertiesNormalized.AADProfile = &armcontainerservice.ManagedClusterAADProfile{
			Managed:             existingMC.Properties.AADProfile.Managed,
			EnableAzureRBAC:     existingMC.Properties.AADProfile.EnableAzureRBAC,
			AdminGroupObjectIDs: existingMC.Properties.AADProfile.AdminGroupObjectIDs,
		}
	}

	if existingMC.Properties.NetworkProfile != nil {
		existingMCPropertiesNormalized.NetworkProfile.LoadBalancerProfile = existingMC.Properties.NetworkProfile.LoadBalancerProfile

		existingMCPropertiesNormalized.NetworkProfile.NetworkPluginMode = existingMC.Properties.NetworkProfile.NetworkPluginMode
	}
	if managedCluster.Properties.NetworkProfile != nil {
		propertiesNormalized.NetworkProfile.LoadBalancerProfile = managedCluster.Properties.NetworkProfile.LoadBalancerProfile

		propertiesNormalized.NetworkProfile.NetworkPluginMode = managedCluster.Properties.NetworkProfile.NetworkPluginMode
		if propertiesNormalized.NetworkProfile.NetworkPluginMode == nil {
			propertiesNormalized.NetworkProfile.NetworkPluginMode = existingMCPropertiesNormalized.NetworkProfile.NetworkPluginMode
		}
	}

	if managedCluster.Properties.APIServerAccessProfile != nil && managedCluster.Properties.APIServerAccessProfile.AuthorizedIPRanges != nil {
		propertiesNormalized.APIServerAccessProfile = &armcontainerservice.ManagedClusterAPIServerAccessProfile{
			AuthorizedIPRanges: managedCluster.Properties.APIServerAccessProfile.AuthorizedIPRanges,
		}
	}

	if existingMC.Properties.APIServerAccessProfile != nil && existingMC.Properties.APIServerAccessProfile.AuthorizedIPRanges != nil {
		existingMCPropertiesNormalized.APIServerAccessProfile = &armcontainerservice.ManagedClusterAPIServerAccessProfile{
			AuthorizedIPRanges: existingMC.Properties.APIServerAccessProfile.AuthorizedIPRanges,
		}
	}

	if managedCluster.Properties.AutoScalerProfile != nil {
		propertiesNormalized.AutoScalerProfile = &armcontainerservice.ManagedClusterPropertiesAutoScalerProfile{
			BalanceSimilarNodeGroups:      managedCluster.Properties.AutoScalerProfile.BalanceSimilarNodeGroups,
			Expander:                      managedCluster.Properties.AutoScalerProfile.Expander,
			MaxEmptyBulkDelete:            managedCluster.Properties.AutoScalerProfile.MaxEmptyBulkDelete,
			MaxGracefulTerminationSec:     managedCluster.Properties.AutoScalerProfile.MaxGracefulTerminationSec,
			MaxNodeProvisionTime:          managedCluster.Properties.AutoScalerProfile.MaxNodeProvisionTime,
			MaxTotalUnreadyPercentage:     managedCluster.Properties.AutoScalerProfile.MaxTotalUnreadyPercentage,
			NewPodScaleUpDelay:            managedCluster.Properties.AutoScalerProfile.NewPodScaleUpDelay,
			OkTotalUnreadyCount:           managedCluster.Properties.AutoScalerProfile.OkTotalUnreadyCount,
			ScanInterval:                  managedCluster.Properties.AutoScalerProfile.ScanInterval,
			ScaleDownDelayAfterAdd:        managedCluster.Properties.AutoScalerProfile.ScaleDownDelayAfterAdd,
			ScaleDownDelayAfterDelete:     managedCluster.Properties.AutoScalerProfile.ScaleDownDelayAfterDelete,
			ScaleDownDelayAfterFailure:    managedCluster.Properties.AutoScalerProfile.ScaleDownDelayAfterFailure,
			ScaleDownUnneededTime:         managedCluster.Properties.AutoScalerProfile.ScaleDownUnneededTime,
			ScaleDownUnreadyTime:          managedCluster.Properties.AutoScalerProfile.ScaleDownUnreadyTime,
			ScaleDownUtilizationThreshold: managedCluster.Properties.AutoScalerProfile.ScaleDownUtilizationThreshold,
			SkipNodesWithLocalStorage:     managedCluster.Properties.AutoScalerProfile.SkipNodesWithLocalStorage,
			SkipNodesWithSystemPods:       managedCluster.Properties.AutoScalerProfile.SkipNodesWithSystemPods,
		}
	}

	if existingMC.Properties.AutoScalerProfile != nil {
		existingMCPropertiesNormalized.AutoScalerProfile = &armcontainerservice.ManagedClusterPropertiesAutoScalerProfile{
			BalanceSimilarNodeGroups:      existingMC.Properties.AutoScalerProfile.BalanceSimilarNodeGroups,
			Expander:                      existingMC.Properties.AutoScalerProfile.Expander,
			MaxEmptyBulkDelete:            existingMC.Properties.AutoScalerProfile.MaxEmptyBulkDelete,
			MaxGracefulTerminationSec:     existingMC.Properties.AutoScalerProfile.MaxGracefulTerminationSec,
			MaxNodeProvisionTime:          existingMC.Properties.AutoScalerProfile.MaxNodeProvisionTime,
			MaxTotalUnreadyPercentage:     existingMC.Properties.AutoScalerProfile.MaxTotalUnreadyPercentage,
			NewPodScaleUpDelay:            existingMC.Properties.AutoScalerProfile.NewPodScaleUpDelay,
			OkTotalUnreadyCount:           existingMC.Properties.AutoScalerProfile.OkTotalUnreadyCount,
			ScanInterval:                  existingMC.Properties.AutoScalerProfile.ScanInterval,
			ScaleDownDelayAfterAdd:        existingMC.Properties.AutoScalerProfile.ScaleDownDelayAfterAdd,
			ScaleDownDelayAfterDelete:     existingMC.Properties.AutoScalerProfile.ScaleDownDelayAfterDelete,
			ScaleDownDelayAfterFailure:    existingMC.Properties.AutoScalerProfile.ScaleDownDelayAfterFailure,
			ScaleDownUnneededTime:         existingMC.Properties.AutoScalerProfile.ScaleDownUnneededTime,
			ScaleDownUnreadyTime:          existingMC.Properties.AutoScalerProfile.ScaleDownUnreadyTime,
			ScaleDownUtilizationThreshold: existingMC.Properties.AutoScalerProfile.ScaleDownUtilizationThreshold,
			SkipNodesWithLocalStorage:     existingMC.Properties.AutoScalerProfile.SkipNodesWithLocalStorage,
			SkipNodesWithSystemPods:       existingMC.Properties.AutoScalerProfile.SkipNodesWithSystemPods,
		}
	}

	if managedCluster.Properties.IdentityProfile != nil {
		propertiesNormalized.IdentityProfile = map[string]*armcontainerservice.UserAssignedIdentity{
			kubeletIdentityKey: {
				ResourceID: managedCluster.Properties.IdentityProfile[kubeletIdentityKey].ResourceID,
			},
		}
	}

	if existingMC.Properties.IdentityProfile != nil {
		existingMCPropertiesNormalized.IdentityProfile = map[string]*armcontainerservice.UserAssignedIdentity{
			kubeletIdentityKey: {
				ResourceID: existingMC.Properties.IdentityProfile[kubeletIdentityKey].ResourceID,
			},
		}
	}

	// Once the AKS autoscaler has been updated it will always return values so we need to
	// respect those values even though the settings are now not being explicitly set by CAPZ.
	if existingMC.Properties.AutoScalerProfile != nil && managedCluster.Properties.AutoScalerProfile == nil {
		existingMCPropertiesNormalized.AutoScalerProfile = nil
		propertiesNormalized.AutoScalerProfile = nil
	}

	clusterNormalized := &armcontainerservice.ManagedCluster{
		Properties: propertiesNormalized,
	}
	existingMCClusterNormalized := &armcontainerservice.ManagedCluster{
		Properties: existingMCPropertiesNormalized,
	}

	if managedCluster.Identity != nil {
		clusterNormalized.Identity = &armcontainerservice.ManagedClusterIdentity{
			Type:                   managedCluster.Identity.Type,
			UserAssignedIdentities: managedCluster.Identity.UserAssignedIdentities,
		}
	}

	if existingMC.Identity != nil {
		existingMCClusterNormalized.Identity = &armcontainerservice.ManagedClusterIdentity{
			Type:                   existingMC.Identity.Type,
			UserAssignedIdentities: existingMC.Identity.UserAssignedIdentities,
		}

		// ClientID and PrincipalID are read-only and should not trigger a diff.
		for _, id := range existingMCClusterNormalized.Identity.UserAssignedIdentities {
			if id != nil {
				id.ClientID = nil
				id.PrincipalID = nil
			}
		}
	}

	if managedCluster.SKU != nil {
		clusterNormalized.SKU = managedCluster.SKU
	}
	if existingMC.SKU != nil {
		existingMCClusterNormalized.SKU = existingMC.SKU
	}

	if managedCluster.Properties.OidcIssuerProfile != nil {
		clusterNormalized.Properties.OidcIssuerProfile = &armcontainerservice.ManagedClusterOIDCIssuerProfile{
			Enabled: managedCluster.Properties.OidcIssuerProfile.Enabled,
		}
	}
	if existingMC.Properties.OidcIssuerProfile != nil {
		existingMCClusterNormalized.Properties.OidcIssuerProfile = &armcontainerservice.ManagedClusterOIDCIssuerProfile{
			Enabled: existingMC.Properties.OidcIssuerProfile.Enabled,
		}
	}

	if managedCluster.Properties.DisableLocalAccounts != nil {
		clusterNormalized.Properties.DisableLocalAccounts = managedCluster.Properties.DisableLocalAccounts
	}

	if existingMC.Properties.DisableLocalAccounts != nil {
		existingMCClusterNormalized.Properties.DisableLocalAccounts = existingMC.Properties.DisableLocalAccounts
	}

	diff := cmp.Diff(clusterNormalized, existingMCClusterNormalized)
	return diff
}

func getIdentity(identity *infrav1.Identity) (managedClusterIdentity *armcontainerservice.ManagedClusterIdentity, err error) {
	if identity.Type == "" {
		return
	}

	managedClusterIdentity = &armcontainerservice.ManagedClusterIdentity{
		Type: ptr.To(armcontainerservice.ResourceIdentityType(identity.Type)),
	}
	if ptr.Deref(managedClusterIdentity.Type, "") == armcontainerservice.ResourceIdentityTypeUserAssigned {
		if identity.UserAssignedIdentityResourceID == "" {
			err = errors.Errorf("Identity is set to \"UserAssigned\" but no UserAssignedIdentityResourceID is present")
			return
		}
		managedClusterIdentity.UserAssignedIdentities = map[string]*armcontainerservice.ManagedServiceIdentityUserAssignedIdentitiesValue{
			identity.UserAssignedIdentityResourceID: {},
		}
	}
	return
}

// isAuthIPRangesNilOrEmpty returns true if the managed cluster's APIServerAccessProfile or AuthorizedIPRanges is nil or if AuthorizedIPRanges is empty.
func isAuthIPRangesNilOrEmpty(managedCluster armcontainerservice.ManagedCluster) bool {
	return managedCluster.Properties.APIServerAccessProfile == nil ||
		managedCluster.Properties.APIServerAccessProfile.AuthorizedIPRanges == nil ||
		reflect.DeepEqual(managedCluster.Properties.APIServerAccessProfile.AuthorizedIPRanges, []*string{})
}
