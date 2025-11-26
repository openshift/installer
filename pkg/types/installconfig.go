package types

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	configv1 "github.com/openshift/api/config/v1"
	features "github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/featuregates"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervc"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
)

const (
	// InstallConfigVersion is the version supported by this package.
	// If you bump this, you must also update the list of convertable values in
	// pkg/types/conversion/installconfig.go
	InstallConfigVersion = "v1"
)

var (
	// PlatformNames is a slice with all the visibly-supported
	// platform names in alphabetical order. This is the list of
	// platforms presented to the user in the interactive wizard.
	PlatformNames = []string{
		aws.Name,
		azure.Name,
		baremetal.Name,
		gcp.Name,
		ibmcloud.Name,
		nutanix.Name,
		openstack.Name,
		powervc.Name,
		powervs.Name,
		vsphere.Name,
	}
	// HiddenPlatformNames is a slice with all the
	// hidden-but-supported platform names. This list isn't presented
	// to the user in the interactive wizard.
	HiddenPlatformNames = []string{
		external.Name,
		none.Name,
	}

	// SCOS is a setting to enable CentOS Stream CoreOS-only modifications
	SCOS = false
)

// PublishingStrategy is a strategy for how various endpoints for the cluster are exposed.
// +kubebuilder:validation:Enum="";External;Internal;Mixed
type PublishingStrategy string

const (
	// ExternalPublishingStrategy exposes endpoints for the cluster to the Internet.
	ExternalPublishingStrategy PublishingStrategy = "External"
	// InternalPublishingStrategy exposes the endpoints for the cluster to the private network only.
	InternalPublishingStrategy PublishingStrategy = "Internal"
	// MixedPublishingStrategy allows for the api server and the ingress to be configured individually for exposure to
	// private network or Internet.
	MixedPublishingStrategy PublishingStrategy = "Mixed"
)

// PolicyType is for usage polices that are applied to additionalTrustBundle.
// +kubebuilder:validation:Enum="";Proxyonly;Always
type PolicyType string

const (
	// PolicyProxyOnly  enables use of AdditionalTrustBundle when http/https proxy is configured.
	PolicyProxyOnly PolicyType = "Proxyonly"
	// PolicyAlways ignores all conditions and uses AdditionalTrustBundle.
	PolicyAlways PolicyType = "Always"
)

//go:generate go run ../../vendor/sigs.k8s.io/controller-tools/cmd/controller-gen crd:crdVersions=v1 paths=. output:dir=../../data/data/
//go:generate go run ../../vendor/k8s.io/code-generator/cmd/deepcopy-gen --output-file zz_generated.deepcopy.go ./...

// InstallConfig is the configuration for an OpenShift install.
type InstallConfig struct {
	// +optional
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata"`

	// AdditionalTrustBundle is a PEM-encoded X.509 certificate bundle
	// that will be added to the nodes' trusted certificate store.
	//
	// +optional
	AdditionalTrustBundle string `json:"additionalTrustBundle,omitempty"`

	// AdditionalTrustBundlePolicy determines when to add the AdditionalTrustBundle
	// to the nodes' trusted certificate store. "Proxyonly" is the default.
	// The field can be set to following specified values.
	// "Proxyonly" : adds the AdditionalTrustBundle to nodes when http/https proxy is configured.
	// "Always" : always adds AdditionalTrustBundle.
	AdditionalTrustBundlePolicy PolicyType `json:"additionalTrustBundlePolicy,omitempty"`

	// SSHKey is the public Secure Shell (SSH) key to provide access to instances.
	// +optional
	SSHKey string `json:"sshKey,omitempty"`

	// BaseDomain is the base domain to which the cluster should belong.
	BaseDomain string `json:"baseDomain"`

	// Networking is the configuration for the pod network provider in
	// the cluster.
	*Networking `json:"networking,omitempty"`

	// ControlPlane is the configuration for the machines that comprise the
	// control plane.
	// +optional
	ControlPlane *MachinePool `json:"controlPlane,omitempty"`

	// Arbiter is the configuration for the machines that comprise the
	// arbiter nodes.
	// +optional
	Arbiter *MachinePool `json:"arbiter,omitempty"`

	// Compute is the configuration for the machines that comprise the
	// compute nodes.
	// +optional
	Compute []MachinePool `json:"compute,omitempty"`

	// Platform is the configuration for the specific platform upon which to
	// perform the installation.
	Platform `json:"platform"`

	// PullSecret is the secret to use when pulling images.
	PullSecret string `json:"pullSecret"`

	// Proxy defines the proxy settings for the cluster.
	// If unset, the cluster will not be configured to use a proxy.
	// +optional
	Proxy *Proxy `json:"proxy,omitempty"`

	// ImageContentSources lists sources/repositories for the release-image content.
	// The field is deprecated. Please use imageDigestSources.
	// +optional
	DeprecatedImageContentSources []ImageContentSource `json:"imageContentSources,omitempty"`

	// ImageDigestSources lists sources/repositories for the release-image content.
	// +optional
	ImageDigestSources []ImageDigestSource `json:"imageDigestSources,omitempty"`

	// Publish controls how the user facing endpoints of the cluster like the Kubernetes API, OpenShift routes etc. are exposed.
	// A "Mixed" strategy only applies to the "azure" platform, and requires "operatorPublishingStrategy" to be configured.
	// When no strategy is specified, the strategy is "External".
	//
	// +kubebuilder:default=External
	// +optional
	Publish PublishingStrategy `json:"publish,omitempty"`

	// OperatorPublishingStrategy controls the visibility of ingress and apiserver. Defaults to public.
	OperatorPublishingStrategy *OperatorPublishingStrategy `json:"operatorPublishingStrategy,omitempty"`

	// FIPS configures https://www.nist.gov/itl/fips-general-information
	//
	// +kubebuilder:default=false
	// +optional
	FIPS bool `json:"fips,omitempty"`

	// CPUPartitioning determines if a cluster should be setup for CPU workload partitioning at install time.
	// When this field is set the cluster will be flagged for CPU Partitioning allowing users to segregate workloads to
	// specific CPU Sets. This does not make any decisions on workloads it only configures the nodes to allow CPU Partitioning.
	// The "AllNodes" value will setup all nodes for CPU Partitioning, the default is "None".
	//
	// +kubebuilder:default="None"
	// +optional
	CPUPartitioning CPUPartitioningMode `json:"cpuPartitioningMode,omitempty"`

	// CredentialsMode is used to explicitly set the mode with which CredentialRequests are satisfied.
	//
	// If this field is set, then the installer will not attempt to query the cloud permissions before attempting
	// installation. If the field is not set or empty, then the installer will perform its normal verification that the
	// credentials provided are sufficient to perform an installation.
	//
	// There are three possible values for this field, but the valid values are dependent upon the platform being used.
	// "Mint": create new credentials with a subset of the overall permissions for each CredentialsRequest
	// "Passthrough": copy the credentials with all of the overall permissions for each CredentialsRequest
	// "Manual": CredentialsRequests must be handled manually by the user
	//
	// For each of the following platforms, the field can be set to the specified values. For all other platforms, the
	// field must not be set.
	// AWS: "Mint", "Passthrough", "Manual"
	// Azure: "Passthrough", "Manual"
	// AzureStack: "Manual"
	// GCP: "Mint", "Passthrough", "Manual"
	// IBMCloud: "Manual"
	// OpenStack: "Passthrough"
	// PowerVC: "Passthrough"
	// PowerVS: "Manual"
	// Nutanix: "Manual"
	// +optional
	CredentialsMode CredentialsMode `json:"credentialsMode,omitempty"`

	// BootstrapInPlace is the configuration for installing a single node
	// with bootstrap in place installation.
	BootstrapInPlace *BootstrapInPlace `json:"bootstrapInPlace,omitempty"`

	// Capabilities configures the installation of optional core cluster components.
	// +optional
	Capabilities *Capabilities `json:"capabilities,omitempty"`

	// FeatureSet enables features that are not part of the default feature set.
	// Valid values are "Default", "TechPreviewNoUpgrade" and "CustomNoUpgrade".
	// When omitted, the "Default" feature set is used.
	// +optional
	FeatureSet configv1.FeatureSet `json:"featureSet,omitempty"`

	// FeatureGates enables a set of custom feature gates.
	// May only be used in conjunction with FeatureSet "CustomNoUpgrade".
	// Features may be enabled or disabled by providing a true or false value for the feature gate.
	// E.g. "featureGates": ["FeatureGate1=true", "FeatureGate2=false"].
	// +optional
	FeatureGates []string `json:"featureGates,omitempty"`
}

// ClusterDomain returns the DNS domain that all records for a cluster must belong to.
func (c *InstallConfig) ClusterDomain() string {
	return fmt.Sprintf("%s.%s", c.ObjectMeta.Name, strings.TrimSuffix(c.BaseDomain, "."))
}

// IsSCOS returns true if CentOs Stream CoreOS-only modifications are enabled
func (c *InstallConfig) IsSCOS() bool {
	return SCOS
}

// IsOKD returns true if community-only modifications are enabled
func (c *InstallConfig) IsOKD() bool {
	return c.IsSCOS()
}

// IsSingleNodeOpenShift returns true if the install-config has been configured for
// bootstrapInPlace
func (c *InstallConfig) IsSingleNodeOpenShift() bool {
	return c.BootstrapInPlace != nil
}

// IsArbiterEnabled returns if arbiter is enabled based off of the install-config arbiter machine pool.
func (c *InstallConfig) IsArbiterEnabled() bool {
	return c.Arbiter != nil &&
		c.Arbiter.Replicas != nil &&
		*c.Arbiter.Replicas > 0
}

// CPUPartitioningMode defines how the nodes should be setup for partitioning the CPU Sets.
// +kubebuilder:validation:Enum=None;AllNodes
type CPUPartitioningMode string

const (
	// CPUPartitioningNone means that no CPU Partitioning is on in this cluster infrastructure.
	CPUPartitioningNone CPUPartitioningMode = "None"
	// CPUPartitioningAllNodes means that all nodes are configured with CPU Partitioning in this cluster.
	CPUPartitioningAllNodes CPUPartitioningMode = "AllNodes"
)

// Platform is the configuration for the specific platform upon which to perform
// the installation. Only one of the platform configuration should be set.
type Platform struct {
	// AWS is the configuration used when installing on AWS.
	// +optional
	AWS *aws.Platform `json:"aws,omitempty"`

	// Azure is the configuration used when installing on Azure.
	// +optional
	Azure *azure.Platform `json:"azure,omitempty"`

	// BareMetal is the configuration used when installing on bare metal.
	// +optional
	BareMetal *baremetal.Platform `json:"baremetal,omitempty"`

	// GCP is the configuration used when installing on Google Cloud Platform.
	// +optional
	GCP *gcp.Platform `json:"gcp,omitempty"`

	// IBMCloud is the configuration used when installing on IBM Cloud.
	// +optional
	IBMCloud *ibmcloud.Platform `json:"ibmcloud,omitempty"`

	// None is the empty configuration used when installing on an unsupported
	// platform.
	None *none.Platform `json:"none,omitempty"`

	// External is the configuration used when installing on
	// an external cloud provider.
	External *external.Platform `json:"external,omitempty"`

	// OpenStack is the configuration used when installing on OpenStack.
	// +optional
	OpenStack *openstack.Platform `json:"openstack,omitempty"`

	// PowerVC is the configuration used when installing on Power VC.
	// +optional
	PowerVC *powervc.Platform `json:"powervc,omitempty"`

	// PowerVS is the configuration used when installing on Power VS.
	// +optional
	PowerVS *powervs.Platform `json:"powervs,omitempty"`

	// VSphere is the configuration used when installing on vSphere.
	// +optional
	VSphere *vsphere.Platform `json:"vsphere,omitempty"`

	// Ovirt is the configuration used when installing on oVirt.
	// +optional
	Ovirt *ovirt.Platform `json:"ovirt,omitempty"`

	// Nutanix is the configuration used when installing on Nutanix.
	// +optional
	Nutanix *nutanix.Platform `json:"nutanix,omitempty"`
}

// OperatorPublishingStrategy is used to control the visibility of the components which can be used to have a mix of public
// and private resources.
type OperatorPublishingStrategy struct {
	// Ingress sets the visibility of the created dns resources.
	// +kubebuilder:validation:Enum="";External;Internal
	// +kubebuilder:default=External
	// +optional
	Ingress string `json:"ingress,omitempty"`

	// APIServer sets the visibility of the load balancers servicing the APIserver.
	// +kubebuilder:validation:Enum="";External;Internal
	// +kubebuilder:default=External
	// +optional
	APIServer string `json:"apiserver,omitempty"`
}

// Name returns a string representation of the platform (e.g. "aws" if
// AWS is non-nil).  It returns an empty string if no platform is
// configured.
func (p *Platform) Name() string {
	switch {
	case p == nil:
		return ""
	case p.AWS != nil:
		return aws.Name
	case p.Azure != nil:
		return azure.Name
	case p.BareMetal != nil:
		return baremetal.Name
	case p.GCP != nil:
		return gcp.Name
	case p.IBMCloud != nil:
		return ibmcloud.Name
	case p.None != nil:
		return none.Name
	case p.External != nil:
		return external.Name
	// The PowerVC check needs to be performed before the OpenStack check
	case p.PowerVC != nil:
		return powervc.Name
	case p.OpenStack != nil:
		return openstack.Name
	case p.VSphere != nil:
		return vsphere.Name
	case p.Ovirt != nil:
		return ovirt.Name
	case p.PowerVS != nil:
		return powervs.Name
	case p.Nutanix != nil:
		return nutanix.Name
	default:
		return ""
	}
}

// Networking defines the pod network provider in the cluster.
type Networking struct {
	// NetworkType is the type of network to install.
	// The default value is OVNKubernetes.
	//
	// +kubebuilder:default=OVNKubernetes
	// +optional
	NetworkType string `json:"networkType,omitempty"`

	// MachineNetwork is the list of IP address pools for machines.
	// This field replaces MachineCIDR, and if set MachineCIDR must
	// be empty or match the first entry in the list.
	// Default is 10.0.0.0/16 for all platforms other than Power VS.
	// For Power VS, the default is 192.168.0.0/24.
	//
	// +optional
	MachineNetwork []MachineNetworkEntry `json:"machineNetwork,omitempty"`

	// ClusterNetwork is the list of IP address pools for pods.
	// Default is 10.128.0.0/14 and a host prefix of /23.
	//
	// +optional
	ClusterNetwork []ClusterNetworkEntry `json:"clusterNetwork,omitempty"`

	// ServiceNetwork is the list of IP address pools for services.
	// Default is 172.30.0.0/16.
	// NOTE: currently only one entry is supported.
	//
	// +kubebuilder:validation:MaxItems=1
	// +optional
	ServiceNetwork []ipnet.IPNet `json:"serviceNetwork,omitempty"`

	// ClusterNetworkMTU is the Maximum Transmit (MTU) Unit size in bytes to allocate to the cluster network.
	// For example, 1200 would set the MTU of the entire overlay network. If the deployment does
	// not require changes in the network plugin, leave it unset and the MTU will be calculated
	// automatically based on the host network MTU.
	// +optional
	ClusterNetworkMTU uint32 `json:"clusterNetworkMTU,omitempty"`

	// OVNKubernetesConfig provides configuration for ovn-kubernetes as the default
	// pod network when NetworkType is set to OVNKubernetes.
	OVNKubernetesConfig *OVNKubernetesConfig `json:"ovnKubernetesConfig,omitempty"`

	// Deprecated types, scheduled to be removed

	// Deprecated way to configure an IP address pool for machines.
	// Replaced by MachineNetwork which allows for multiple pools.
	// +optional
	DeprecatedMachineCIDR *ipnet.IPNet `json:"machineCIDR,omitempty"`

	// Deprecated name for NetworkType
	// +optional
	DeprecatedType string `json:"type,omitempty"`

	// Deprecated way to configure an IP address pool for services.
	// Replaced by ServiceNetwork which allows for multiple pools.
	// +optional
	DeprecatedServiceCIDR *ipnet.IPNet `json:"serviceCIDR,omitempty"`

	// Deprecated name for ClusterNetwork
	// +optional
	DeprecatedClusterNetworks []ClusterNetworkEntry `json:"clusterNetworks,omitempty"`
}

// MachineNetworkEntry is a single IP address block for node IP blocks.
type MachineNetworkEntry struct {
	// CIDR is the IP block address pool for machines within the cluster.
	CIDR ipnet.IPNet `json:"cidr"`
}

// ClusterNetworkEntry is a single IP address block for pod IP blocks. IP blocks
// are allocated with size 2^HostSubnetLength.
type ClusterNetworkEntry struct {
	// CIDR is the IP block address pool.
	CIDR ipnet.IPNet `json:"cidr"`

	// HostPrefix is the prefix size to allocate to each node from the CIDR.
	// For example, 24 would allocate 2^8=256 adresses to each node. If this
	// field is not used by the plugin, it can be left unset.
	// When multiple CIDRs of the same family (i.e. IPv4/IPv6) are present,
	// their HostPrefix value must be the same.
	// +optional
	HostPrefix int32 `json:"hostPrefix,omitempty"`

	// The size of blocks to allocate from the larger pool.
	// This is the length in bits - so a 9 here will allocate a /23.
	// +optional
	DeprecatedHostSubnetLength int32 `json:"hostSubnetLength,omitempty"`
}

// OVNKubernetesConfig configures the ovn-kubernetes sdn plugin.
type OVNKubernetesConfig struct {
	// ipv4 allows users to configure IP settings for IPv4 connections. When omitted,
	// this means no opinions and the default configuration is used. Check individual
	// fields within ipv4 for details of default values.
	// +optional
	IPv4 *IPv4OVNKubernetesConfig `json:"ipv4,omitempty"`
}

// IPv4OVNKubernetesConfig is IPv4 configuration for the ovn-kubernetes sdn plugin.
type IPv4OVNKubernetesConfig struct {
	// internalJoinSubnet is a v4 subnet used internally by ovn-kubernetes in case the
	// default one is being already used by something else. It must not overlap with
	// any other subnet being used by OpenShift or by the node network. The size of the
	// subnet must be larger than the number of nodes. The value cannot be changed
	// after installation.
	// The current default value is 100.64.0.0/16
	// The subnet must be large enough to accommodate one IP per node in your cluster
	// The value must be in proper IPV4 CIDR format
	// +optional
	InternalJoinSubnet *ipnet.IPNet `json:"internalJoinSubnet,omitempty"`
}

// Proxy defines the proxy settings for the cluster.
// At least one of HTTPProxy or HTTPSProxy is required.
type Proxy struct {
	// HTTPProxy is the URL of the proxy for HTTP requests.
	// +optional
	HTTPProxy string `json:"httpProxy,omitempty"`

	// HTTPSProxy is the URL of the proxy for HTTPS requests.
	// +optional
	HTTPSProxy string `json:"httpsProxy,omitempty"`

	// NoProxy is a comma-separated list of domains and CIDRs for which the proxy should not be used.
	// +optional
	NoProxy string `json:"noProxy,omitempty"`
}

// ImageContentSource defines a list of sources/repositories that can be used to pull content.
// The field is deprecated. Please use imageDigestSources.
type ImageContentSource struct {
	// Source is the repository that users refer to, e.g. in image pull specifications.
	Source string `json:"source"`

	// Mirrors is one or more repositories that may also contain the same images.
	// +optional
	Mirrors []string `json:"mirrors,omitempty"`
}

// ImageDigestSource defines a list of sources/repositories that can be used to pull content.
type ImageDigestSource struct {
	// Source is the repository that users refer to, e.g. in image pull specifications.
	Source string `json:"source"`

	// Mirrors is one or more repositories that may also contain the same images.
	// +optional
	Mirrors []string `json:"mirrors,omitempty"`

	// SourcePolicy defines the fallback policy when there is a failure pulling an
	// image from the mirrors.
	// +optional
	SourcePolicy configv1.MirrorSourcePolicy `json:"sourcePolicy"`
}

// CredentialsMode is the mode by which CredentialsRequests will be satisfied.
// +kubebuilder:validation:Enum="";Mint;Passthrough;Manual
type CredentialsMode string

const (
	// ManualCredentialsMode indicates that cloud-credential-operator should not process any CredentialsRequests.
	ManualCredentialsMode CredentialsMode = "Manual"

	// MintCredentialsMode indicates that cloud-credential-operator should be creating users for each
	// CredentialsRequest.
	MintCredentialsMode CredentialsMode = "Mint"

	// PassthroughCredentialsMode indicates that cloud-credential-operator should just copy over the cluster's
	// cloud credentials for each CredentialsRequest.
	PassthroughCredentialsMode CredentialsMode = "Passthrough"
)

// BootstrapInPlace defines the configuration for bootstrap-in-place installation
type BootstrapInPlace struct {
	// InstallationDisk is the target disk drive for coreos-installer
	InstallationDisk string `json:"installationDisk"`
}

// Capabilities selects the managed set of optional, core cluster components.
type Capabilities struct {
	// baselineCapabilitySet selects an initial set of
	// optional capabilities to enable, which can be extended via
	// additionalEnabledCapabilities. The default is vCurrent.
	// +optional
	BaselineCapabilitySet configv1.ClusterVersionCapabilitySet `json:"baselineCapabilitySet,omitempty"`

	// additionalEnabledCapabilities extends the set of managed
	// capabilities beyond the baseline defined in
	// baselineCapabilitySet. The default is an empty set.
	// +optional
	AdditionalEnabledCapabilities []configv1.ClusterVersionCapability `json:"additionalEnabledCapabilities,omitempty"`
}

// GetEnabledCapabilities returns a set of enabled ClusterVersionCapabilities.
func (c *InstallConfig) GetEnabledCapabilities() sets.Set[configv1.ClusterVersionCapability] {
	enabledCaps := sets.Set[configv1.ClusterVersionCapability]{}
	if c.Capabilities == nil || c.Capabilities.BaselineCapabilitySet == "" {
		// when Capabilities and/or BaselineCapabilitySet is not specified, default is vCurrent
		baseSet := configv1.ClusterVersionCapabilitySets[configv1.ClusterVersionCapabilitySetCurrent]
		for _, cap := range baseSet {
			enabledCaps.Insert(cap)
		}
	}
	if c.Capabilities != nil {
		if c.Capabilities.BaselineCapabilitySet != "" {
			baseSet := configv1.ClusterVersionCapabilitySets[c.Capabilities.BaselineCapabilitySet]
			for _, cap := range baseSet {
				enabledCaps.Insert(cap)
			}
		}
		if c.Capabilities.AdditionalEnabledCapabilities != nil {
			for _, cap := range c.Capabilities.AdditionalEnabledCapabilities {
				enabledCaps.Insert(cap)
			}
		}
	}
	return enabledCaps
}

// WorkerMachinePool retrieves the worker MachinePool from InstallConfig.Compute
func (c *InstallConfig) WorkerMachinePool() *MachinePool {
	for _, machinePool := range c.Compute {
		switch machinePool.Name {
		case MachinePoolComputeRoleName, MachinePoolEdgeRoleName:
			return &machinePool
		}
	}

	return nil
}

// EnabledFeatureGates returns a FeatureGate that can be checked (using the Enabled function)
// to determine if a feature gate is enabled in the current feature sets.
func (c *InstallConfig) EnabledFeatureGates() featuregates.FeatureGate {
	var customFS *configv1.CustomFeatureGates

	if c.FeatureSet == configv1.CustomNoUpgrade {
		customFS = featuregates.GenerateCustomFeatures(c.FeatureGates)
	}

	clusterProfile := GetClusterProfileName()
	featureSets, ok := features.AllFeatureSets()[clusterProfile]
	if !ok {
		logrus.Warnf("no feature sets for cluster profile %q", clusterProfile)
	}
	fg := featuregates.FeatureGateFromFeatureSets(featureSets, c.FeatureSet, customFS)

	return fg
}

// PublicAPI indicates whether the API load balancer should be public
// by inspecting the cluster and operator publishing strategies.
func (c *InstallConfig) PublicAPI() bool {
	// When no strategy is specified, the strategy defaults to "External".
	if c.Publish == "" || c.Publish == ExternalPublishingStrategy {
		return true
	}

	if op := c.OperatorPublishingStrategy; op != nil && strings.EqualFold(op.APIServer, "External") {
		return true
	}
	return false
}

// PublicIngress indicates whether the Ingress load balancer should be public
// by inspecting the cluster and operator publishing strategies.
func (c *InstallConfig) PublicIngress() bool {
	// When no strategy is specified, the strategy defaults to "External".
	if c.Publish == "" || c.Publish == ExternalPublishingStrategy {
		return true
	}

	if op := c.OperatorPublishingStrategy; op != nil && strings.EqualFold(op.Ingress, "External") {
		return true
	}
	return false
}
