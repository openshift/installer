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
	"sort"

	asocontainerservicev1 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001"
	asocontainerservicev1preview "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231102preview"
	asocontainerservicev1hub "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20240901/storage"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api/util/secret"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/agentpools"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aso"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	"sigs.k8s.io/cluster-api-provider-azure/util/versions"
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

	// NetworkPolicy used for building Kubernetes network. Possible values include: 'azure', 'calico', 'cilium'.
	NetworkPolicy string

	// NetworkDataplane used for building Kubernetes network. Possible values include: 'azure', 'cilium'.
	NetworkDataplane *infrav1.NetworkDataplaneType

	// OutboundType used for building Kubernetes network. Possible values include: 'loadBalancer', 'managedNATGateway', 'userAssignedNATGateway', 'userDefinedRouting'.
	OutboundType *infrav1.ManagedControlPlaneOutboundType

	// SSHPublicKey is a string literal containing an ssh public key. Will autogenerate and discard if not provided.
	SSHPublicKey string

	// GetAllAgentPools is a function that returns the list of agent pool specifications in this cluster.
	GetAllAgentPools func() ([]azure.ASOResourceSpecGetter[genruntime.MetaObject], error)

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

	// AutoUpgradeProfile defines auto upgrade configuration.
	AutoUpgradeProfile *ManagedClusterAutoUpgradeProfile

	// SecurityProfile defines the security profile for the cluster.
	SecurityProfile *ManagedClusterSecurityProfile

	// Patches are extra patches to be applied to the ASO resource.
	Patches []string

	// Preview enables the preview API version.
	Preview bool
}

// ManagedClusterAutoUpgradeProfile auto upgrade profile for a managed cluster.
type ManagedClusterAutoUpgradeProfile struct {
	// UpgradeChannel defines the channel for auto upgrade configuration.
	UpgradeChannel *infrav1.UpgradeChannel
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
	ManagedOutboundIPs *int

	// OutboundIPPrefixes are the desired outbound IP Prefix resources for the cluster load balancer.
	OutboundIPPrefixes []string

	// OutboundIPs are the desired outbound IP resources for the cluster load balancer.
	OutboundIPs []string

	// AllocatedOutboundPorts are the desired number of allocated SNAT ports per VM. Allowed values must be in the range of 0 to 64000 (inclusive). The default value is 0 which results in Azure dynamically allocating ports.
	AllocatedOutboundPorts *int

	// IdleTimeoutInMinutes  are the desired outbound flow idle timeout in minutes. Allowed values must be in the range of 4 to 120 (inclusive). The default value is 30 minutes.
	IdleTimeoutInMinutes *int
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

// ManagedClusterSecurityProfile defines the security profile for the cluster.
type ManagedClusterSecurityProfile struct {
	// AzureKeyVaultKms defines Azure Key Vault key management service settings for the security profile.
	AzureKeyVaultKms *AzureKeyVaultKms

	// Defender defines Microsoft Defender settings for the security profile.
	Defender *ManagedClusterSecurityProfileDefender

	// ImageCleaner settings for the security profile.
	ImageCleaner *ManagedClusterSecurityProfileImageCleaner

	// Workloadidentity enables Kubernetes applications to access Azure cloud resources securely with Azure AD.
	WorkloadIdentity *ManagedClusterSecurityProfileWorkloadIdentity
}

// ManagedClusterSecurityProfileDefender defines Microsoft Defender settings for the security profile.
type ManagedClusterSecurityProfileDefender struct {
	// LogAnalyticsWorkspaceResourceID is the ID of the Log Analytics workspace that has to be associated with Microsoft Defender.
	// When Microsoft Defender is enabled, this field is required and must be a valid workspace resource ID.
	LogAnalyticsWorkspaceResourceID *string

	// SecurityMonitoring profile defines the Microsoft Defender threat detection for Cloud settings for the security profile.
	SecurityMonitoring *ManagedClusterSecurityProfileDefenderSecurityMonitoring
}

// ManagedClusterSecurityProfileDefenderSecurityMonitoring settings for the security profile threat detection.
type ManagedClusterSecurityProfileDefenderSecurityMonitoring struct {
	// Enabled enables Defender threat detection
	Enabled *bool
}

// ManagedClusterSecurityProfileImageCleaner removes unused images from nodes, freeing up disk space and helping to reduce attack surface area.
type ManagedClusterSecurityProfileImageCleaner struct {
	// Enabled enables Image Cleaner on AKS cluster.
	Enabled *bool

	// Image Cleaner scanning interval in hours.
	IntervalHours *int
}

// ManagedClusterSecurityProfileWorkloadIdentity defines Workload identity settings for the security profile.
type ManagedClusterSecurityProfileWorkloadIdentity struct {
	// Enabled enables workload identity.
	Enabled *bool
}

// AzureKeyVaultKms Azure Key Vault key management service settings for the security profile.
type AzureKeyVaultKms struct {
	// Enabled enables Azure Key Vault key management service. The default is false.
	Enabled *bool

	// KeyID defines the Identifier of Azure Key Vault key.
	// When Azure Key Vault key management service is enabled, this field is required and must be a valid key identifier.
	KeyID *string

	// KeyVaultNetworkAccess defines the network access of key vault.
	// The possible values are Public and Private.
	// Public means the key vault allows public access from all networks.
	// Private means the key vault disables public access and enables private link. The default value is Public.
	KeyVaultNetworkAccess *infrav1.KeyVaultNetworkAccessTypes

	// KeyVaultResourceID is the Resource ID of key vault. When keyVaultNetworkAccess is Private, this field is required and must be a valid resource ID.
	KeyVaultResourceID *string
}

// Constants for managed cluster configuration.
const (
	// identityTypeMSI is the identity type for managed service identity.
	identityTypeMSI = "msi"
	// dnsIPLastOctet is the default last octet for the DNS service IP.
	dnsIPLastOctet = byte(10)
)

// buildAutoScalerProfile builds the AutoScalerProfile for the ManagedClusterProperties.
func buildAutoScalerProfile(autoScalerProfile *AutoScalerProfile) *asocontainerservicev1hub.ManagedClusterProperties_AutoScalerProfile {
	if autoScalerProfile == nil {
		return nil
	}

	mcAutoScalerProfile := &asocontainerservicev1hub.ManagedClusterProperties_AutoScalerProfile{
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
		mcAutoScalerProfile.Expander = ptr.To(string(asocontainerservicev1.ManagedClusterProperties_AutoScalerProfile_Expander(*autoScalerProfile.Expander)))
	}

	return mcAutoScalerProfile
}

// getManagedClusterVersion gets the desired managed k8s version.
// If autoupgrade channels is set to patch, stable or rapid, clusters can be upgraded to higher version by AKS.
// If autoupgrade is triggered, existing kubernetes version will be higher than the user desired kubernetes version.
// CAPZ should honour the upgrade and it should not downgrade to the lower desired version.
func (s *ManagedClusterSpec) getManagedClusterVersion(existing *asocontainerservicev1hub.ManagedCluster) string {
	if existing == nil || existing.Status.CurrentKubernetesVersion == nil {
		return s.Version
	}
	return versions.GetHigherK8sVersion(s.Version, *existing.Status.CurrentKubernetesVersion)
}

// ResourceRef implements azure.ASOResourceSpecGetter.
func (s *ManagedClusterSpec) ResourceRef() genruntime.MetaObject {
	if s.Preview {
		return &asocontainerservicev1preview.ManagedCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: azure.GetNormalizedKubernetesName(s.Name),
			},
		}
	}
	return &asocontainerservicev1.ManagedCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: azure.GetNormalizedKubernetesName(s.Name),
		},
	}
}

// Parameters returns the parameters for the managed clusters.
func (s *ManagedClusterSpec) Parameters(ctx context.Context, existingObj genruntime.MetaObject) (params genruntime.MetaObject, err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "managedclusters.Service.Parameters")
	defer done()

	existing, err := s.convertExistingToHub(existingObj)
	if err != nil {
		return nil, err
	}

	managedCluster := s.initializeManagedCluster(existing)

	s.configureBasicProperties(managedCluster, existing)

	if err := s.configureNetworkProfile(managedCluster); err != nil {
		return nil, err
	}

	if err := s.configureLinuxProfile(managedCluster); err != nil {
		return nil, err
	}

	s.configureOperatorSpec(managedCluster)
	s.configureAADProfile(managedCluster)
	s.configureAddonProfiles(managedCluster)
	s.configureSKU(managedCluster)
	s.configureAPIServerAccessProfile(managedCluster)

	if err := s.configureIdentity(managedCluster); err != nil {
		return nil, err
	}

	s.configureHTTPProxyConfig(managedCluster)
	s.configureOIDCIssuerProfile(managedCluster)
	s.configureAutoUpgradeProfile(managedCluster)
	s.configureSecurityProfile(managedCluster)

	if err := s.configureAgentPoolProfiles(ctx, managedCluster); err != nil {
		return nil, err
	}

	return s.convertToTargetVersion(managedCluster)
}

// convertExistingToHub converts the existing object to the hub version for manipulation.
func (s *ManagedClusterSpec) convertExistingToHub(existingObj genruntime.MetaObject) (*asocontainerservicev1hub.ManagedCluster, error) {
	if existingObj == nil {
		return nil, nil
	}

	hub := &asocontainerservicev1hub.ManagedCluster{}
	if err := existingObj.(conversion.Convertible).ConvertTo(hub); err != nil {
		return nil, err
	}
	return hub, nil
}

// initializeManagedCluster creates a new managed cluster or returns the existing one.
func (s *ManagedClusterSpec) initializeManagedCluster(existing *asocontainerservicev1hub.ManagedCluster) *asocontainerservicev1hub.ManagedCluster {
	if existing != nil {
		return existing
	}

	return &asocontainerservicev1hub.ManagedCluster{
		Spec: asocontainerservicev1hub.ManagedCluster_Spec{
			Tags: infrav1.Build(infrav1.BuildParams{
				Lifecycle:   infrav1.ResourceLifecycleOwned,
				ClusterName: s.ClusterName,
				Name:        ptr.To(s.Name),
				Role:        ptr.To(infrav1.CommonRole),
			}),
		},
	}
}

// configureBasicProperties sets the basic properties of the managed cluster.
func (s *ManagedClusterSpec) configureBasicProperties(managedCluster *asocontainerservicev1hub.ManagedCluster, existing *asocontainerservicev1hub.ManagedCluster) {
	managedCluster.Spec.AzureName = s.Name
	managedCluster.Spec.Owner = &genruntime.KnownResourceReference{
		Name: azure.GetNormalizedKubernetesName(s.ResourceGroup),
	}
	managedCluster.Spec.Identity = &asocontainerservicev1hub.ManagedClusterIdentity{
		Type: ptr.To(string(asocontainerservicev1.ManagedClusterIdentity_Type_SystemAssigned)),
	}
	managedCluster.Spec.Location = &s.Location
	managedCluster.Spec.NodeResourceGroup = &s.NodeResourceGroup
	managedCluster.Spec.EnableRBAC = ptr.To(true)
	managedCluster.Spec.DnsPrefix = s.DNSPrefix

	if kubernetesVersion := s.getManagedClusterVersion(existing); kubernetesVersion != "" {
		managedCluster.Spec.KubernetesVersion = &kubernetesVersion
	}

	managedCluster.Spec.ServicePrincipalProfile = &asocontainerservicev1hub.ManagedClusterServicePrincipalProfile{
		ClientId: ptr.To(identityTypeMSI),
	}
	managedCluster.Spec.AutoScalerProfile = buildAutoScalerProfile(s.AutoScalerProfile)
}

// configureNetworkProfile configures the network profile for the managed cluster.
func (s *ManagedClusterSpec) configureNetworkProfile(managedCluster *asocontainerservicev1hub.ManagedCluster) error {
	managedCluster.Spec.NetworkProfile = &asocontainerservicev1hub.ContainerServiceNetworkProfile{
		NetworkPlugin:   azure.AliasOrNil[string](ptr.To(s.NetworkPlugin)),
		LoadBalancerSku: azure.AliasOrNil[string](ptr.To(s.LoadBalancerSKU)),
		NetworkPolicy:   azure.AliasOrNil[string](ptr.To(s.NetworkPolicy)),
	}

	if s.NetworkDataplane != nil {
		managedCluster.Spec.NetworkProfile.NetworkDataplane = ptr.To(string(asocontainerservicev1.ContainerServiceNetworkProfile_NetworkDataplane(*s.NetworkDataplane)))
	}

	if s.NetworkPluginMode != nil {
		managedCluster.Spec.NetworkProfile.NetworkPluginMode = ptr.To(string(asocontainerservicev1.ContainerServiceNetworkProfile_NetworkPluginMode(*s.NetworkPluginMode)))
	}

	if s.PodCIDR != "" {
		managedCluster.Spec.NetworkProfile.PodCidr = &s.PodCIDR
	}

	if err := s.configureServiceCIDR(managedCluster); err != nil {
		return err
	}

	if s.LoadBalancerProfile != nil {
		managedCluster.Spec.NetworkProfile.LoadBalancerProfile = s.GetLoadBalancerProfile()
	}

	if s.OutboundType != nil {
		managedCluster.Spec.NetworkProfile.OutboundType = ptr.To(string(asocontainerservicev1.ContainerServiceNetworkProfile_OutboundType(*s.OutboundType)))
	}

	return nil
}

// configureServiceCIDR configures the service CIDR and DNS service IP.
func (s *ManagedClusterSpec) configureServiceCIDR(managedCluster *asocontainerservicev1hub.ManagedCluster) error {
	if s.ServiceCIDR == "" {
		return nil
	}

	managedCluster.Spec.NetworkProfile.ServiceCidr = &s.ServiceCIDR
	managedCluster.Spec.NetworkProfile.DnsServiceIP = s.DNSServiceIP

	if s.DNSServiceIP == nil {
		dnsIP, err := s.calculateDNSServiceIP()
		if err != nil {
			return err
		}
		managedCluster.Spec.NetworkProfile.DnsServiceIP = &dnsIP
	}

	return nil
}

// calculateDNSServiceIP calculates the DNS service IP from the service CIDR.
func (s *ManagedClusterSpec) calculateDNSServiceIP() (string, error) {
	ip, _, err := net.ParseCIDR(s.ServiceCIDR)
	if err != nil {
		return "", fmt.Errorf("failed to parse service cidr: %w", err)
	}
	// Set the last octet of the IP to .10 to ensure a valid DNS IP within the service CIDR.
	// Reference: https://golang.org/src/net/ip.go#L48
	ip[15] = dnsIPLastOctet
	return ip.String(), nil
}

// configureLinuxProfile configures the Linux profile with SSH settings.
func (s *ManagedClusterSpec) configureLinuxProfile(managedCluster *asocontainerservicev1hub.ManagedCluster) error {
	if s.SSHPublicKey == "" {
		return nil
	}

	decodedSSHPublicKey, err := base64.StdEncoding.DecodeString(s.SSHPublicKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode SSHPublicKey")
	}

	managedCluster.Spec.LinuxProfile = &asocontainerservicev1hub.ContainerServiceLinuxProfile{
		AdminUsername: ptr.To(azure.DefaultAKSUserName),
		Ssh: &asocontainerservicev1hub.ContainerServiceSshConfiguration{
			PublicKeys: []asocontainerservicev1hub.ContainerServiceSshPublicKey{
				{
					KeyData: ptr.To(string(decodedSSHPublicKey)),
				},
			},
		},
	}

	return nil
}

// configureOperatorSpec configures the operator spec for kubeconfig secrets.
func (s *ManagedClusterSpec) configureOperatorSpec(managedCluster *asocontainerservicev1hub.ManagedCluster) {
	managedCluster.Spec.OperatorSpec = &asocontainerservicev1hub.ManagedClusterOperatorSpec{
		Secrets: &asocontainerservicev1hub.ManagedClusterOperatorSecrets{
			AdminCredentials: &genruntime.SecretDestination{
				Name: adminKubeconfigSecretName(s.ClusterName),
				Key:  secret.KubeconfigDataName,
			},
		},
	}
}

// configureAADProfile configures the Azure Active Directory profile.
func (s *ManagedClusterSpec) configureAADProfile(managedCluster *asocontainerservicev1hub.ManagedCluster) {
	if s.AADProfile == nil {
		return
	}

	managedCluster.Spec.AadProfile = &asocontainerservicev1hub.ManagedClusterAADProfile{
		Managed:             &s.AADProfile.Managed,
		EnableAzureRBAC:     &s.AADProfile.EnableAzureRBAC,
		AdminGroupObjectIDs: s.AADProfile.AdminGroupObjectIDs,
	}

	if s.DisableLocalAccounts != nil {
		managedCluster.Spec.DisableLocalAccounts = s.DisableLocalAccounts
	}

	if ptr.Deref(s.DisableLocalAccounts, false) {
		// Admin credentials cannot be fetched when local accounts are disabled.
		managedCluster.Spec.OperatorSpec.Secrets.AdminCredentials = nil
	}

	if s.AADProfile.Managed {
		managedCluster.Spec.OperatorSpec.Secrets.UserCredentials = &genruntime.SecretDestination{
			Name: userKubeconfigSecretName(s.ClusterName),
			Key:  secret.KubeconfigDataName,
		}
	}
}

// configureAddonProfiles configures the addon profiles for the managed cluster.
func (s *ManagedClusterSpec) configureAddonProfiles(managedCluster *asocontainerservicev1hub.ManagedCluster) {
	if len(s.AddonProfiles) == 0 {
		return
	}

	managedCluster.Spec.AddonProfiles = make(map[string]asocontainerservicev1hub.ManagedClusterAddonProfile, len(s.AddonProfiles))
	for _, item := range s.AddonProfiles {
		addonProfile := asocontainerservicev1hub.ManagedClusterAddonProfile{
			Enabled: &item.Enabled,
		}
		if item.Config != nil {
			addonProfile.Config = item.Config
		}
		managedCluster.Spec.AddonProfiles[item.Name] = addonProfile
	}
}

// configureSKU configures the SKU for the managed cluster.
func (s *ManagedClusterSpec) configureSKU(managedCluster *asocontainerservicev1hub.ManagedCluster) {
	if s.SKU == nil {
		return
	}

	tierName := asocontainerservicev1.ManagedClusterSKU_Tier(s.SKU.Tier)
	managedCluster.Spec.Sku = &asocontainerservicev1hub.ManagedClusterSKU{
		Name: ptr.To(string(asocontainerservicev1.ManagedClusterSKU_Name("Base"))),
		Tier: ptr.To(string(tierName)),
	}
}

// configureAPIServerAccessProfile configures the API server access profile.
func (s *ManagedClusterSpec) configureAPIServerAccessProfile(managedCluster *asocontainerservicev1hub.ManagedCluster) {
	if s.APIServerAccessProfile == nil {
		return
	}

	managedCluster.Spec.ApiServerAccessProfile = &asocontainerservicev1hub.ManagedClusterAPIServerAccessProfile{
		EnablePrivateCluster:           s.APIServerAccessProfile.EnablePrivateCluster,
		PrivateDNSZone:                 s.APIServerAccessProfile.PrivateDNSZone,
		EnablePrivateClusterPublicFQDN: s.APIServerAccessProfile.EnablePrivateClusterPublicFQDN,
	}

	if s.APIServerAccessProfile.AuthorizedIPRanges != nil {
		managedCluster.Spec.ApiServerAccessProfile.AuthorizedIPRanges = s.APIServerAccessProfile.AuthorizedIPRanges
	}
}

// configureIdentity configures the identity settings for the managed cluster.
func (s *ManagedClusterSpec) configureIdentity(managedCluster *asocontainerservicev1hub.ManagedCluster) error {
	if s.Identity != nil {
		identity, err := getIdentity(s.Identity)
		if err != nil {
			return errors.Wrapf(err, "Identity is not valid: %s", err)
		}
		managedCluster.Spec.Identity = identity
	}

	if s.KubeletUserAssignedIdentity != "" {
		managedCluster.Spec.IdentityProfile = map[string]asocontainerservicev1hub.UserAssignedIdentity{
			kubeletIdentityKey: {
				ResourceReference: &genruntime.ResourceReference{
					ARMID: s.KubeletUserAssignedIdentity,
				},
			},
		}
	}

	return nil
}

// configureHTTPProxyConfig configures the HTTP proxy settings.
func (s *ManagedClusterSpec) configureHTTPProxyConfig(managedCluster *asocontainerservicev1hub.ManagedCluster) {
	if s.HTTPProxyConfig == nil {
		return
	}

	managedCluster.Spec.HttpProxyConfig = &asocontainerservicev1hub.ManagedClusterHTTPProxyConfig{
		HttpProxy:  s.HTTPProxyConfig.HTTPProxy,
		HttpsProxy: s.HTTPProxyConfig.HTTPSProxy,
		TrustedCa:  s.HTTPProxyConfig.TrustedCA,
	}

	if s.HTTPProxyConfig.NoProxy != nil {
		managedCluster.Spec.HttpProxyConfig.NoProxy = s.HTTPProxyConfig.NoProxy
	}
}

// configureOIDCIssuerProfile configures the OIDC issuer profile.
func (s *ManagedClusterSpec) configureOIDCIssuerProfile(managedCluster *asocontainerservicev1hub.ManagedCluster) {
	if s.OIDCIssuerProfile == nil {
		return
	}

	managedCluster.Spec.OidcIssuerProfile = &asocontainerservicev1hub.ManagedClusterOIDCIssuerProfile{
		Enabled: s.OIDCIssuerProfile.Enabled,
	}

	if ptr.Deref(s.OIDCIssuerProfile.Enabled, false) {
		managedCluster.Spec.OperatorSpec.ConfigMaps = &asocontainerservicev1hub.ManagedClusterOperatorConfigMaps{
			OIDCIssuerProfile: &genruntime.ConfigMapDestination{
				Name: oidcIssuerURLConfigMapName(s.ClusterName),
				Key:  oidcIssuerProfileURL,
			},
		}
	}
}

// configureAutoUpgradeProfile configures the auto upgrade profile.
func (s *ManagedClusterSpec) configureAutoUpgradeProfile(managedCluster *asocontainerservicev1hub.ManagedCluster) {
	if s.AutoUpgradeProfile == nil {
		return
	}

	managedCluster.Spec.AutoUpgradeProfile = &asocontainerservicev1hub.ManagedClusterAutoUpgradeProfile{
		UpgradeChannel: (*string)(s.AutoUpgradeProfile.UpgradeChannel),
	}
}

// configureSecurityProfile configures the security profile for the managed cluster.
func (s *ManagedClusterSpec) configureSecurityProfile(managedCluster *asocontainerservicev1hub.ManagedCluster) {
	if s.SecurityProfile == nil {
		return
	}

	securityProfile := &asocontainerservicev1hub.ManagedClusterSecurityProfile{}

	s.configureAzureKeyVaultKms(securityProfile)
	s.configureDefender(securityProfile)
	s.configureImageCleaner(securityProfile)
	s.configureWorkloadIdentity(securityProfile)

	managedCluster.Spec.SecurityProfile = securityProfile
}

// configureAzureKeyVaultKms configures the Azure Key Vault KMS settings.
func (s *ManagedClusterSpec) configureAzureKeyVaultKms(securityProfile *asocontainerservicev1hub.ManagedClusterSecurityProfile) {
	if s.SecurityProfile.AzureKeyVaultKms == nil {
		return
	}

	securityProfile.AzureKeyVaultKms = &asocontainerservicev1hub.AzureKeyVaultKms{
		Enabled: s.SecurityProfile.AzureKeyVaultKms.Enabled,
		KeyId:   s.SecurityProfile.AzureKeyVaultKms.KeyID,
	}

	if s.SecurityProfile.AzureKeyVaultKms.KeyVaultNetworkAccess != nil {
		keyVaultNetworkAccess := string(*s.SecurityProfile.AzureKeyVaultKms.KeyVaultNetworkAccess)
		securityProfile.AzureKeyVaultKms.KeyVaultNetworkAccess = ptr.To(string(asocontainerservicev1.AzureKeyVaultKms_KeyVaultNetworkAccess(keyVaultNetworkAccess)))
	}

	if s.SecurityProfile.AzureKeyVaultKms.KeyVaultResourceID != nil {
		securityProfile.AzureKeyVaultKms.KeyVaultResourceReference = &genruntime.ResourceReference{
			ARMID: *s.SecurityProfile.AzureKeyVaultKms.KeyVaultResourceID,
		}
	}
}

// configureDefender configures the Microsoft Defender settings.
func (s *ManagedClusterSpec) configureDefender(securityProfile *asocontainerservicev1hub.ManagedClusterSecurityProfile) {
	if s.SecurityProfile.Defender == nil {
		return
	}

	securityProfile.Defender = &asocontainerservicev1hub.ManagedClusterSecurityProfileDefender{
		LogAnalyticsWorkspaceResourceReference: &genruntime.ResourceReference{
			ARMID: *s.SecurityProfile.Defender.LogAnalyticsWorkspaceResourceID,
		},
	}

	if s.SecurityProfile.Defender.SecurityMonitoring != nil {
		securityProfile.Defender.SecurityMonitoring = &asocontainerservicev1hub.ManagedClusterSecurityProfileDefenderSecurityMonitoring{
			Enabled: s.SecurityProfile.Defender.SecurityMonitoring.Enabled,
		}
	}
}

// configureImageCleaner configures the image cleaner settings.
func (s *ManagedClusterSpec) configureImageCleaner(securityProfile *asocontainerservicev1hub.ManagedClusterSecurityProfile) {
	if s.SecurityProfile.ImageCleaner == nil {
		return
	}

	securityProfile.ImageCleaner = &asocontainerservicev1hub.ManagedClusterSecurityProfileImageCleaner{
		Enabled:       s.SecurityProfile.ImageCleaner.Enabled,
		IntervalHours: s.SecurityProfile.ImageCleaner.IntervalHours,
	}
}

// configureWorkloadIdentity configures the workload identity settings.
func (s *ManagedClusterSpec) configureWorkloadIdentity(securityProfile *asocontainerservicev1hub.ManagedClusterSecurityProfile) {
	if s.SecurityProfile.WorkloadIdentity == nil {
		return
	}

	securityProfile.WorkloadIdentity = &asocontainerservicev1hub.ManagedClusterSecurityProfileWorkloadIdentity{
		Enabled: s.SecurityProfile.WorkloadIdentity.Enabled,
	}
}

// configureAgentPoolProfiles configures the agent pool profiles for initial cluster creation.
func (s *ManagedClusterSpec) configureAgentPoolProfiles(ctx context.Context, managedCluster *asocontainerservicev1hub.ManagedCluster) error {
	// Only include AgentPoolProfiles during initial cluster creation.
	// Agent pools are managed solely by the AzureManagedMachinePool controller thereafter.
	managedCluster.Spec.AgentPoolProfiles = nil

	if managedCluster.Status.AgentPoolProfiles != nil {
		return nil
	}

	agentPoolProfiles, err := s.buildAgentPoolProfiles(ctx)
	if err != nil {
		return err
	}

	// Keep a consistent order to prevent unnecessary updates.
	sort.Slice(agentPoolProfiles, func(i, j int) bool {
		return ptr.Deref(agentPoolProfiles[i].Name, "") < ptr.Deref(agentPoolProfiles[j].Name, "")
	})

	managedCluster.Spec.AgentPoolProfiles = agentPoolProfiles
	return nil
}

// buildAgentPoolProfiles builds the agent pool profiles from the agent pool specs.
func (s *ManagedClusterSpec) buildAgentPoolProfiles(ctx context.Context) ([]asocontainerservicev1hub.ManagedClusterAgentPoolProfile, error) {
	agentPoolSpecs, err := s.GetAllAgentPools()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get agent pool specs for managed cluster %s", s.Name)
	}

	scheme := runtime.NewScheme()
	if err := asocontainerservicev1.AddToScheme(scheme); err != nil {
		return nil, errors.Wrap(err, "error constructing scheme")
	}
	if err := asocontainerservicev1preview.AddToScheme(scheme); err != nil {
		return nil, errors.Wrap(err, "error constructing scheme")
	}

	var profiles []asocontainerservicev1hub.ManagedClusterAgentPoolProfile
	for _, agentPoolSpec := range agentPoolSpecs {
		profile, err := s.buildSingleAgentPoolProfile(ctx, scheme, agentPoolSpec)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// buildSingleAgentPoolProfile builds a single agent pool profile from its spec.
func (s *ManagedClusterSpec) buildSingleAgentPoolProfile(ctx context.Context, scheme *runtime.Scheme, agentPoolSpec azure.ASOResourceSpecGetter[genruntime.MetaObject]) (asocontainerservicev1hub.ManagedClusterAgentPoolProfile, error) {
	agentPool, err := aso.PatchedParameters(ctx, scheme, agentPoolSpec, nil)
	if err != nil {
		return asocontainerservicev1hub.ManagedClusterAgentPoolProfile{}, errors.Wrapf(err, "failed to get agent pool parameters for managed cluster %s", s.Name)
	}

	agentPoolSpecTyped := agentPoolSpec.(*agentpools.AgentPoolSpec)
	hubAgentPool := &asocontainerservicev1hub.ManagedClustersAgentPool{}

	if s.Preview {
		agentPoolTyped := agentPool.(*asocontainerservicev1preview.ManagedClustersAgentPool)
		if err := agentPoolTyped.ConvertTo(hubAgentPool); err != nil {
			return asocontainerservicev1hub.ManagedClusterAgentPoolProfile{}, err
		}
	} else {
		agentPoolTyped := agentPool.(*asocontainerservicev1.ManagedClustersAgentPool)
		if err := agentPoolTyped.ConvertTo(hubAgentPool); err != nil {
			return asocontainerservicev1hub.ManagedClusterAgentPoolProfile{}, err
		}
	}

	hubAgentPool.Spec.AzureName = agentPoolSpecTyped.AzureName
	return converters.AgentPoolToManagedClusterAgentPoolProfile(hubAgentPool), nil
}

// convertToTargetVersion converts the managed cluster to the target API version.
func (s *ManagedClusterSpec) convertToTargetVersion(managedCluster *asocontainerservicev1hub.ManagedCluster) (genruntime.MetaObject, error) {
	if s.Preview {
		prev := &asocontainerservicev1preview.ManagedCluster{}
		if err := prev.ConvertFrom(managedCluster); err != nil {
			return nil, err
		}
		return prev, nil
	}

	stable := &asocontainerservicev1.ManagedCluster{}
	if err := stable.ConvertFrom(managedCluster); err != nil {
		return nil, err
	}

	return stable, nil
}

// GetLoadBalancerProfile returns an asocontainerservicev1.ManagedClusterLoadBalancerProfile from the
// information present in ManagedClusterSpec.LoadBalancerProfile.
func (s *ManagedClusterSpec) GetLoadBalancerProfile() (loadBalancerProfile *asocontainerservicev1hub.ManagedClusterLoadBalancerProfile) {
	loadBalancerProfile = &asocontainerservicev1hub.ManagedClusterLoadBalancerProfile{
		AllocatedOutboundPorts: s.LoadBalancerProfile.AllocatedOutboundPorts,
		IdleTimeoutInMinutes:   s.LoadBalancerProfile.IdleTimeoutInMinutes,
	}
	if s.LoadBalancerProfile.ManagedOutboundIPs != nil {
		loadBalancerProfile.ManagedOutboundIPs = &asocontainerservicev1hub.ManagedClusterLoadBalancerProfile_ManagedOutboundIPs{Count: s.LoadBalancerProfile.ManagedOutboundIPs}
	}
	if len(s.LoadBalancerProfile.OutboundIPPrefixes) > 0 {
		loadBalancerProfile.OutboundIPPrefixes = &asocontainerservicev1hub.ManagedClusterLoadBalancerProfile_OutboundIPPrefixes{
			PublicIPPrefixes: convertToResourceReferences(s.LoadBalancerProfile.OutboundIPPrefixes),
		}
	}
	if len(s.LoadBalancerProfile.OutboundIPs) > 0 {
		loadBalancerProfile.OutboundIPs = &asocontainerservicev1hub.ManagedClusterLoadBalancerProfile_OutboundIPs{
			PublicIPs: convertToResourceReferences(s.LoadBalancerProfile.OutboundIPs),
		}
	}
	return
}

func convertToResourceReferences(resources []string) []asocontainerservicev1hub.ResourceReference {
	resourceReferences := make([]asocontainerservicev1hub.ResourceReference, len(resources))
	for i := range resources {
		resourceReferences[i] = asocontainerservicev1hub.ResourceReference{
			Reference: &genruntime.ResourceReference{
				ARMID: resources[i],
			},
		}
	}
	return resourceReferences
}

func getIdentity(identity *infrav1.Identity) (managedClusterIdentity *asocontainerservicev1hub.ManagedClusterIdentity, err error) {
	if identity.Type == "" {
		return
	}

	managedClusterIdentity = &asocontainerservicev1hub.ManagedClusterIdentity{
		Type: ptr.To(string(asocontainerservicev1.ManagedClusterIdentity_Type(identity.Type))),
	}
	if ptr.Deref(managedClusterIdentity.Type, "") == string(asocontainerservicev1.ManagedClusterIdentity_Type_UserAssigned) {
		if identity.UserAssignedIdentityResourceID == "" {
			err = errors.Errorf("Identity is set to \"UserAssigned\" but no UserAssignedIdentityResourceID is present")
			return
		}
		managedClusterIdentity.UserAssignedIdentities = []asocontainerservicev1hub.UserAssignedIdentityDetails{
			{
				Reference: genruntime.ResourceReference{
					ARMID: identity.UserAssignedIdentityResourceID,
				},
			},
		}
	}
	return
}

func adminKubeconfigSecretName(clusterName string) string {
	return secret.Name(clusterName+"-aso", secret.Kubeconfig)
}

func oidcIssuerURLConfigMapName(clusterName string) string {
	return secret.Name(clusterName+"-aso", "oidc-issuer-profile")
}

func userKubeconfigSecretName(clusterName string) string {
	return secret.Name(clusterName+"-user-aso", secret.Kubeconfig)
}

// WasManaged implements azure.ASOResourceSpecGetter.
func (s *ManagedClusterSpec) WasManaged(_ genruntime.MetaObject) bool {
	// CAPZ has never supported BYO managed clusters.
	return true
}

var _ aso.TagsGetterSetter[genruntime.MetaObject] = (*ManagedClusterSpec)(nil)

// GetAdditionalTags implements aso.TagsGetterSetter.
func (s *ManagedClusterSpec) GetAdditionalTags() infrav1.Tags {
	return s.Tags
}

// GetDesiredTags implements aso.TagsGetterSetter.
func (s *ManagedClusterSpec) GetDesiredTags(resource genruntime.MetaObject) infrav1.Tags {
	if s.Preview {
		return resource.(*asocontainerservicev1preview.ManagedCluster).Spec.Tags
	}
	return resource.(*asocontainerservicev1.ManagedCluster).Spec.Tags
}

// SetTags implements aso.TagsGetterSetter.
func (s *ManagedClusterSpec) SetTags(resource genruntime.MetaObject, tags infrav1.Tags) {
	if s.Preview {
		resource.(*asocontainerservicev1preview.ManagedCluster).Spec.Tags = tags
		return
	}
	resource.(*asocontainerservicev1.ManagedCluster).Spec.Tags = tags
}

var _ aso.Patcher = (*ManagedClusterSpec)(nil)

// ExtraPatches implements aso.Patcher.
func (s *ManagedClusterSpec) ExtraPatches() []string {
	return s.Patches
}
