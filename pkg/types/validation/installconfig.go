package validation

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	dockerref "github.com/containers/image/v5/docker/reference"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	utilsnet "k8s.io/utils/net"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/api/features"
	operv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/hostcrypt"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	awsvalidation "github.com/openshift/installer/pkg/types/aws/validation"
	"github.com/openshift/installer/pkg/types/azure"
	azurevalidation "github.com/openshift/installer/pkg/types/azure/validation"
	"github.com/openshift/installer/pkg/types/baremetal"
	baremetalvalidation "github.com/openshift/installer/pkg/types/baremetal/validation"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/featuregates"
	"github.com/openshift/installer/pkg/types/gcp"
	gcpvalidation "github.com/openshift/installer/pkg/types/gcp/validation"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	ibmcloudvalidation "github.com/openshift/installer/pkg/types/ibmcloud/validation"
	"github.com/openshift/installer/pkg/types/nutanix"
	nutanixvalidation "github.com/openshift/installer/pkg/types/nutanix/validation"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackvalidation "github.com/openshift/installer/pkg/types/openstack/validation"
	"github.com/openshift/installer/pkg/types/ovirt"
	ovirtvalidation "github.com/openshift/installer/pkg/types/ovirt/validation"
	"github.com/openshift/installer/pkg/types/powervs"
	powervsvalidation "github.com/openshift/installer/pkg/types/powervs/validation"
	"github.com/openshift/installer/pkg/types/vsphere"
	vspherevalidation "github.com/openshift/installer/pkg/types/vsphere/validation"
	"github.com/openshift/installer/pkg/validate"
	"github.com/openshift/installer/pkg/version"
)

// hostCryptBypassedAnnotation is set if the host crypt check was bypassed via environment variable.
const hostCryptBypassedAnnotation = "install.openshift.io/hostcrypt-check-bypassed"

// list of known plugins that require hostPrefix to be set
var pluginsUsingHostPrefix = sets.NewString(string(operv1.NetworkTypeOVNKubernetes))

// ValidateInstallConfig checks that the specified install config is valid.
//
//nolint:gocyclo
func ValidateInstallConfig(c *types.InstallConfig, usingAgentMethod bool) field.ErrorList {
	allErrs := field.ErrorList{}
	if c.TypeMeta.APIVersion == "" {
		return field.ErrorList{field.Required(field.NewPath("apiVersion"), "install-config version required")}
	}
	switch v := c.APIVersion; v {
	case types.InstallConfigVersion:
		// Current version
	default:
		return field.ErrorList{field.Invalid(field.NewPath("apiVersion"), c.TypeMeta.APIVersion, fmt.Sprintf("install-config version must be %q", types.InstallConfigVersion))}
	}

	if c.FIPS {
		allErrs = append(allErrs, validateFIPSconfig(c)...)
	} else if c.SSHKey != "" {
		if err := validate.SSHPublicKey(c.SSHKey); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("sshKey"), c.SSHKey, err.Error()))
		}
	}

	if c.AdditionalTrustBundle != "" {
		if err := validate.CABundle(c.AdditionalTrustBundle); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("additionalTrustBundle"), c.AdditionalTrustBundle, err.Error()))
		}
	}
	if c.AdditionalTrustBundlePolicy != "" {
		if err := validateAdditionalCABundlePolicy(c); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("additionalTrustBundlePolicy"), c.AdditionalTrustBundlePolicy, err.Error()))
		}
	}
	nameErr := validate.ClusterName(c.ObjectMeta.Name)
	if c.Platform.GCP != nil || c.Platform.Azure != nil {
		nameErr = validate.ClusterName1035(c.ObjectMeta.Name)
	}
	if c.Platform.VSphere != nil || c.Platform.BareMetal != nil || c.Platform.OpenStack != nil || c.Platform.Nutanix != nil {
		nameErr = validate.OnPremClusterName(c.ObjectMeta.Name)
	}
	if nameErr != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata", "name"), c.ObjectMeta.Name, nameErr.Error()))
	}
	baseDomainErr := validate.DomainName(c.BaseDomain, true)
	if baseDomainErr != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("baseDomain"), c.BaseDomain, baseDomainErr.Error()))
	}
	if nameErr == nil && baseDomainErr == nil {
		clusterDomain := c.ClusterDomain()
		if err := validate.DomainName(clusterDomain, true); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("baseDomain"), clusterDomain, err.Error()))
		}
	}
	if c.Networking != nil {
		allErrs = append(allErrs, validateNetworking(c.Networking, field.NewPath("networking"))...)
		allErrs = append(allErrs, validateNetworkingIPVersion(c.Networking, &c.Platform)...)
		allErrs = append(allErrs, validateNetworkingClusterNetworkMTU(c, field.NewPath("networking", "clusterNetworkMTU"))...)
		allErrs = append(allErrs, validateVIPsForPlatform(c.Networking, &c.Platform, usingAgentMethod, field.NewPath("platform"))...)
		allErrs = append(allErrs, validateOVNKubernetesConfig(c.Networking, field.NewPath("networking"))...)
	} else {
		allErrs = append(allErrs, field.Required(field.NewPath("networking"), "networking is required"))
	}
	allErrs = append(allErrs, validatePlatform(&c.Platform, usingAgentMethod, field.NewPath("platform"), c.Networking, c)...)
	if c.ControlPlane != nil {
		allErrs = append(allErrs, validateControlPlane(&c.Platform, c.ControlPlane, field.NewPath("controlPlane"))...)
	} else {
		allErrs = append(allErrs, field.Required(field.NewPath("controlPlane"), "controlPlane is required"))
	}

	if c.Arbiter != nil {
		if c.EnabledFeatureGates().Enabled(features.FeatureGateHighlyAvailableArbiter) {
			allErrs = append(allErrs, validateArbiter(&c.Platform, c.Arbiter, c.ControlPlane, field.NewPath("arbiter"))...)
		} else {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("arbiter"), fmt.Sprintf("%s feature must be enabled in order to use arbiter cluster deployment", features.FeatureGateHighlyAvailableArbiter)))
		}
	}
	multiArchEnabled := types.MultiArchFeatureGateEnabled(c.Platform.Name(), c.EnabledFeatureGates())
	allErrs = append(allErrs, validateCompute(&c.Platform, c.ControlPlane, c.Compute, field.NewPath("compute"), multiArchEnabled)...)

	releaseArch, err := version.ReleaseArchitecture()
	if err != nil {
		allErrs = append(allErrs, field.InternalError(nil, err))
	} else {
		allErrs = append(allErrs, validateReleaseArchitecture(c.ControlPlane, c.Compute, types.Architecture(releaseArch))...)
	}

	if err := validate.ImagePullSecret(c.PullSecret); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("pullSecret"), c.PullSecret, err.Error()))
	}
	if c.Proxy != nil {
		allErrs = append(allErrs, validateProxy(c.Proxy, c, field.NewPath("proxy"))...)
	}
	allErrs = append(allErrs, validateImageContentSources(c.DeprecatedImageContentSources, field.NewPath("imageContentSources"))...)
	if _, ok := validPublishingStrategies[c.Publish]; !ok {
		allErrs = append(allErrs, field.NotSupported(field.NewPath("publish"), c.Publish, validPublishingStrategyValues))
	}
	allErrs = append(allErrs, validateImageDigestSources(c.ImageDigestSources, field.NewPath("imageDigestSources"))...)
	if _, ok := validPublishingStrategies[c.Publish]; !ok {
		allErrs = append(allErrs, field.NotSupported(field.NewPath("publish"), c.Publish, validPublishingStrategyValues))
	}
	if len(c.DeprecatedImageContentSources) != 0 && len(c.ImageDigestSources) != 0 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("imageContentSources"), c.Publish, "cannot set imageContentSources and imageDigestSources at the same time"))
	}
	if len(c.DeprecatedImageContentSources) != 0 {
		logrus.Warningln("imageContentSources is deprecated, please use ImageDigestSources")
	}
	allErrs = append(allErrs, validateCloudCredentialsMode(c.CredentialsMode, field.NewPath("credentialsMode"), c.Platform)...)
	if c.Capabilities != nil {
		allErrs = append(allErrs, validateCapabilities(c.Capabilities, field.NewPath("capabilities"))...)
	}

	if c.Publish == types.InternalPublishingStrategy {
		switch platformName := c.Platform.Name(); platformName {
		case aws.Name, azure.Name, gcp.Name, ibmcloud.Name, powervs.Name:
		default:
			allErrs = append(allErrs, field.Invalid(field.NewPath("publish"), c.Publish, fmt.Sprintf("Internal publish strategy is not supported on %q platform", platformName)))
		}
	}

	if c.Publish == types.MixedPublishingStrategy {
		switch platformName := c.Platform.Name(); platformName {
		case azure.Name:
		default:
			allErrs = append(allErrs, field.Invalid(field.NewPath("publish"), c.Publish, fmt.Sprintf("mixed publish strategy is not supported on %q platform", platformName)))
		}
		if c.OperatorPublishingStrategy == nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("publish"), c.Publish, "please specify the operator publishing strategy for mixed publish strategy"))
		}
	} else if c.OperatorPublishingStrategy != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("operatorPublishingStrategy"), c.Publish, "operator publishing strategy is only allowed with mixed publishing strategy installs"))
	}

	if c.OperatorPublishingStrategy != nil {
		acceptedValues := sets.New[string]("Internal", "External")
		if c.OperatorPublishingStrategy.APIServer == "" {
			c.OperatorPublishingStrategy.APIServer = "External"
		}
		if c.OperatorPublishingStrategy.Ingress == "" {
			c.OperatorPublishingStrategy.Ingress = "External"
		}
		if !acceptedValues.Has(c.OperatorPublishingStrategy.APIServer) {
			allErrs = append(allErrs, field.NotSupported(field.NewPath("apiserver"), c.OperatorPublishingStrategy.APIServer, sets.List(acceptedValues)))
		}
		if !acceptedValues.Has(c.OperatorPublishingStrategy.Ingress) {
			allErrs = append(allErrs, field.NotSupported(field.NewPath("ingress"), c.OperatorPublishingStrategy.Ingress, sets.List(acceptedValues)))
		}
		if c.OperatorPublishingStrategy.APIServer == "Internal" && c.OperatorPublishingStrategy.Ingress == "Internal" {
			allErrs = append(allErrs, field.Invalid(field.NewPath("publish"), c.OperatorPublishingStrategy.APIServer, "cannot set both fields to internal in a mixed cluster, use publish internal instead"))
		}
	}

	if c.Capabilities != nil {
		capSet := c.Capabilities.BaselineCapabilitySet
		if capSet == "" {
			capSet = configv1.ClusterVersionCapabilitySetCurrent
		}
		enabledCaps := sets.New[configv1.ClusterVersionCapability](configv1.ClusterVersionCapabilitySets[capSet]...)
		enabledCaps.Insert(c.Capabilities.AdditionalEnabledCapabilities...)

		if c.Capabilities.BaselineCapabilitySet == configv1.ClusterVersionCapabilitySetNone {
			enabledCaps := sets.New[configv1.ClusterVersionCapability](c.Capabilities.AdditionalEnabledCapabilities...)
			if enabledCaps.Has(configv1.ClusterVersionCapabilityMarketplace) && !enabledCaps.Has(configv1.ClusterVersionCapabilityOperatorLifecycleManager) {
				allErrs = append(allErrs, field.Invalid(field.NewPath("additionalEnabledCapabilities"), c.Capabilities.AdditionalEnabledCapabilities,
					"the marketplace capability requires the OperatorLifecycleManager capability"))
			}
			if c.Platform.BareMetal != nil && !enabledCaps.Has(configv1.ClusterVersionCapabilityBaremetal) {
				allErrs = append(allErrs, field.Invalid(field.NewPath("additionalEnabledCapabilities"), c.Capabilities.AdditionalEnabledCapabilities,
					"platform baremetal requires the baremetal capability"))
			}
		}

		if enabledCaps.Has(configv1.ClusterVersionCapabilityMarketplace) && !enabledCaps.Has(configv1.ClusterVersionCapabilityOperatorLifecycleManager) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("additionalEnabledCapabilities"), c.Capabilities.AdditionalEnabledCapabilities,
				"the marketplace capability requires the OperatorLifecycleManager capability"))
		}

		if !enabledCaps.Has(configv1.ClusterVersionCapabilityCloudCredential) {
			// check if platform is cloud
			if c.None == nil && c.BareMetal == nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("capabilities"), c.Capabilities,
					"disabling CloudCredential capability available only for baremetal platforms"))
			}
		}

		if !enabledCaps.Has(configv1.ClusterVersionCapabilityCloudControllerManager) {
			if c.None == nil && c.BareMetal == nil && c.External == nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("capabilities"), c.Capabilities,
					"disabling CloudControllerManager is only supported on the Baremetal, None, or External platform with cloudControllerManager value none"))
			}
			if c.External != nil && c.External.CloudControllerManager == external.CloudControllerManagerTypeExternal {
				allErrs = append(allErrs, field.Invalid(field.NewPath("capabilities"), c.Capabilities,
					"disabling CloudControllerManager on External platform supported only with cloudControllerManager value none"))
			}
		}

		if !enabledCaps.Has(configv1.ClusterVersionCapabilityIngress) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("capabilities"), c.Capabilities,
				"the Ingress capability is required"))
		}
	}

	allErrs = append(allErrs, ValidateFeatureSet(c)...)

	return allErrs
}

// ipAddressType indicates the address types provided for a given field
type ipAddressType struct {
	IPv4    bool
	IPv6    bool
	Primary corev1.IPFamily
}

// ipAddressTypeByField is a map of field path to ipAddressType
type ipAddressTypeByField map[string]ipAddressType

// ipNetByField is a map of field path to the IPNets
type ipNetByField map[string][]ipnet.IPNet

// inferIPVersionFromInstallConfig infers the user's desired ip version from the networking config.
// Presence field names match the field path of the struct within the Networking type. This function
// assumes a valid install config.
func inferIPVersionFromInstallConfig(n *types.Networking) (hasIPv4, hasIPv6 bool, presence ipAddressTypeByField, addresses ipNetByField) {
	if n == nil {
		return
	}
	addresses = make(ipNetByField)
	for _, network := range n.MachineNetwork {
		addresses["machineNetwork"] = append(addresses["machineNetwork"], network.CIDR)
	}
	for _, network := range n.ServiceNetwork {
		addresses["serviceNetwork"] = append(addresses["serviceNetwork"], network)
	}
	for _, network := range n.ClusterNetwork {
		addresses["clusterNetwork"] = append(addresses["clusterNetwork"], network.CIDR)
	}
	presence = make(ipAddressTypeByField)
	for k, ipnets := range addresses {
		for i, ipnet := range ipnets {
			has := presence[k]
			if ipnet.IP.To4() != nil {
				has.IPv4 = true
				if i == 0 {
					has.Primary = corev1.IPv4Protocol
				}
				if k == "serviceNetwork" {
					hasIPv4 = true
				}
			} else {
				has.IPv6 = true
				if i == 0 {
					has.Primary = corev1.IPv6Protocol
				}
				if k == "serviceNetwork" {
					hasIPv6 = true
				}
			}
			presence[k] = has
		}
	}
	return
}

func ipnetworksToStrings(networks []ipnet.IPNet) []string {
	var diag []string
	for _, sn := range networks {
		diag = append(diag, sn.String())
	}
	return diag
}

// validateNetworkingIPVersion checks parameters for consistency when the user
// requests single-stack IPv6 or dual-stack modes.
func validateNetworkingIPVersion(n *types.Networking, p *types.Platform) field.ErrorList {
	var allErrs field.ErrorList

	hasIPv4, hasIPv6, presence, addresses := inferIPVersionFromInstallConfig(n)

	switch {
	case hasIPv4 && hasIPv6:
		if len(n.ServiceNetwork) != 2 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("networking", "serviceNetwork"), strings.Join(ipnetworksToStrings(n.ServiceNetwork), ", "), "when installing dual-stack IPv4/IPv6 you must provide two service networks, one for each IP address type"))
		}

		allowV6Primary := false
		experimentalDualStackEnabled, _ := strconv.ParseBool(os.Getenv("OPENSHIFT_INSTALL_EXPERIMENTAL_DUAL_STACK"))
		switch {
		case p.Azure != nil && experimentalDualStackEnabled:
			logrus.Warnf("Using experimental Azure dual-stack support")
		case p.BareMetal != nil:
			// We now support ipv6-primary dual stack on baremetal
			allowV6Primary = true
		case p.VSphere != nil:
			// as well as on vSphere
			allowV6Primary = true
		case p.OpenStack != nil:
			allowV6Primary = true
		case p.Ovirt != nil:
		case p.Nutanix != nil:
		case p.None != nil:
		case p.External != nil:
		default:
			allErrs = append(allErrs, field.Invalid(field.NewPath("networking"), "DualStack", "dual-stack IPv4/IPv6 is not supported for this platform, specify only one type of address"))
		}
		for k, v := range presence {
			switch {
			case v.IPv4 && !v.IPv6:
				allErrs = append(allErrs, field.Invalid(field.NewPath("networking", k), strings.Join(ipnetworksToStrings(addresses[k]), ", "), "dual-stack IPv4/IPv6 requires an IPv6 network in this list"))
			case !v.IPv4 && v.IPv6:
				allErrs = append(allErrs, field.Invalid(field.NewPath("networking", k), strings.Join(ipnetworksToStrings(addresses[k]), ", "), "dual-stack IPv4/IPv6 requires an IPv4 network in this list"))
			}

			// FIXME: we should allow either all-networks-IPv4Primary or
			// all-networks-IPv6Primary, but the latter currently causes
			// confusing install failures, so block it.
			if !allowV6Primary && v.IPv4 && v.IPv6 && v.Primary != corev1.IPv4Protocol {
				allErrs = append(allErrs, field.Invalid(field.NewPath("networking", k), strings.Join(ipnetworksToStrings(addresses[k]), ", "), "IPv4 addresses must be listed before IPv6 addresses"))
			}
		}

	case hasIPv6:
		switch {
		case p.BareMetal != nil:
		case p.VSphere != nil:
		case p.OpenStack != nil:
		case p.Ovirt != nil:
		case p.Nutanix != nil:
		case p.None != nil:
		case p.External != nil:
		case p.Azure != nil && p.Azure.CloudName == azure.StackCloud:
			allErrs = append(allErrs, field.Invalid(field.NewPath("networking"), "IPv6", "Azure Stack does not support IPv6"))
		default:
			allErrs = append(allErrs, field.Invalid(field.NewPath("networking"), "IPv6", "single-stack IPv6 is not supported for this platform"))
		}

	case hasIPv4:
		if len(n.ServiceNetwork) > 1 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("networking", "serviceNetwork"), strings.Join(ipnetworksToStrings(n.ServiceNetwork), ", "), "only one service network can be specified"))
		}

	default:
		// we should have a validation error for no specified machineNetwork, serviceNetwork, or clusterNetwork
	}

	return allErrs
}

func validateNetworking(n *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(n.NetworkType) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("networkType"), "network provider type required"))
	}

	// NOTE(dulek): We're hardcoding "Kuryr" here as the plan is to remove it from the API very soon. We can remove
	//              this check once some more general validation of the supported NetworkTypes is in place.
	if strings.EqualFold(n.NetworkType, "Kuryr") {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("networkType"), n.NetworkType, "networkType Kuryr is not supported on OpenShift later than 4.14"))
	}

	if strings.EqualFold(n.NetworkType, string(operv1.NetworkTypeOpenShiftSDN)) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("networkType"), n.NetworkType, "networkType OpenShiftSDN is not supported, please use OVNKubernetes"))
	}

	if len(n.MachineNetwork) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("machineNetwork"), "at least one machine network is required"))
	}
	for i, mn := range n.MachineNetwork {
		allErrs = append(allErrs, validateMachineNetwork(n, &mn, i, fldPath.Child("machineNetwork").Index(i))...)
	}

	if len(n.ServiceNetwork) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("serviceNetwork"), "a service network is required"))
	}
	for i, sn := range n.ServiceNetwork {
		allErrs = append(allErrs, validateServiceNetwork(n, &sn, i, fldPath.Child("serviceNetwork").Index(i))...)
	}

	if len(n.ClusterNetwork) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("clusterNetwork"), "cluster network required"))
	}
	for i, cn := range n.ClusterNetwork {
		allErrs = append(allErrs, validateClusterNetwork(n, &cn, i, fldPath.Child("clusterNetwork").Index(i))...)
	}

	return allErrs
}

func validateMachineNetwork(n *types.Networking, mn *types.MachineNetworkEntry, idx int, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if err := validate.SubnetCIDR(&mn.CIDR.IPNet); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, mn.CIDR.String(), err.Error()))
		return allErrs // CIDR value is invalid, so no further validation.
	}

	if validate.DoCIDRsOverlap(&mn.CIDR.IPNet, validate.DockerBridgeSubnet) {
		logrus.Warnf("%s: %s overlaps with default Docker Bridge subnet", fldPath, mn.CIDR.String())
	}

	allErrs = append(allErrs, validateNetworkNotOverlapDefaultOVNSubnets(n, &mn.CIDR.IPNet, fldPath)...)

	for i, subNetwork := range n.MachineNetwork[0:idx] {
		if validate.DoCIDRsOverlap(&mn.CIDR.IPNet, &subNetwork.CIDR.IPNet) {
			allErrs = append(allErrs, field.Invalid(fldPath, mn.CIDR.String(), fmt.Sprintf("machine network must not overlap with machine network %d", i)))
		}
	}

	return allErrs
}

func validateServiceNetwork(n *types.Networking, sn *ipnet.IPNet, idx int, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if err := validate.ServiceSubnetCIDR(&sn.IPNet); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, sn.String(), err.Error()))
		return allErrs // CIDR value is invalid, so no further validation.
	}

	if validate.DoCIDRsOverlap(&sn.IPNet, validate.DockerBridgeSubnet) {
		logrus.Warnf("%s: %s overlaps with default Docker Bridge subnet", fldPath, sn.String())
	}

	allErrs = append(allErrs, validateNetworkNotOverlapDefaultOVNSubnets(n, &sn.IPNet, fldPath)...)

	for _, mn := range n.MachineNetwork {
		if validate.DoCIDRsOverlap(&sn.IPNet, &mn.CIDR.IPNet) {
			allErrs = append(allErrs, field.Invalid(fldPath, sn.String(), "service network must not overlap with any of the machine networks"))
		}
	}
	for i, snn := range n.ServiceNetwork[0:idx] {
		if validate.DoCIDRsOverlap(&sn.IPNet, &snn.IPNet) {
			allErrs = append(allErrs, field.Invalid(fldPath, sn.String(), fmt.Sprintf("service network must not overlap with service network %d", i)))
		}
	}
	return allErrs
}

func validateClusterNetwork(n *types.Networking, cn *types.ClusterNetworkEntry, idx int, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if err := validate.SubnetCIDR(&cn.CIDR.IPNet); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("cidr"), cn.CIDR.IPNet.String(), err.Error()))
		return allErrs // CIDR value is invalid, so no further validation.
	}

	if validate.DoCIDRsOverlap(&cn.CIDR.IPNet, validate.DockerBridgeSubnet) {
		logrus.Warnf("%s: %s overlaps with default Docker Bridge subnet", fldPath.Index(idx), cn.CIDR.String())
	}

	allErrs = append(allErrs, validateNetworkNotOverlapDefaultOVNSubnets(n, &cn.CIDR.IPNet, fldPath.Child("cidr"))...)

	for _, network := range n.MachineNetwork {
		if validate.DoCIDRsOverlap(&cn.CIDR.IPNet, &network.CIDR.IPNet) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("cidr"), cn.CIDR.String(), "cluster network must not overlap with any of the machine networks"))
		}
	}
	for i, sn := range n.ServiceNetwork {
		if validate.DoCIDRsOverlap(&cn.CIDR.IPNet, &sn.IPNet) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("cidr"), cn.CIDR.String(), fmt.Sprintf("cluster network must not overlap with service network %d", i)))
		}
	}
	for i, acn := range n.ClusterNetwork[0:idx] {
		if validate.DoCIDRsOverlap(&cn.CIDR.IPNet, &acn.CIDR.IPNet) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("cidr"), cn.CIDR.String(), fmt.Sprintf("cluster network must not overlap with cluster network %d", i)))
		}
	}

	// ignore hostPrefix if the plugin does not use it and has it unset
	if pluginsUsingHostPrefix.Has(n.NetworkType) || (cn.HostPrefix != 0) {
		if cn.HostPrefix < 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("hostPrefix"), cn.HostPrefix, "hostPrefix must be positive"))
		} else if ones, bits := cn.CIDR.Mask.Size(); cn.HostPrefix < int32(ones) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("hostPrefix"), cn.HostPrefix, "cluster network host subnetwork prefix must not be larger size than CIDR "+cn.CIDR.String()))
		} else if bits == 32 {
			// setting different value for clusternetwork CIDR host prefix is not allowed
			// we only need to check IPv4 as IPv6 prefix must be 64
			for _, acn := range n.ClusterNetwork[0:idx] {
				if acn.CIDR.IP.To4() != nil && cn.HostPrefix != acn.HostPrefix {
					allErrs = append(allErrs, field.Invalid(fldPath.Child("hostPrefix"), cn.HostPrefix, "cluster network host subnetwork prefix must be the same value for IPv4 networks"))
					break
				}
			}
		} else if bits == 128 && cn.HostPrefix != 64 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("hostPrefix"), cn.HostPrefix, "cluster network host subnetwork prefix must be 64 for IPv6 networks"))
		}
	}
	return allErrs
}

func validateNetworkNotOverlapDefaultOVNSubnets(n *types.Networking, network *net.IPNet, fldPath *field.Path) field.ErrorList {
	if !strings.EqualFold(n.NetworkType, string(operv1.NetworkTypeOVNKubernetes)) {
		return nil
	}

	allErrs := field.ErrorList{}

	// getOVNSubnet returns the *net.IPNet for each type of subnet that will be used by OVNKubernetes
	// and whether it is user-defined in the install-config.
	getOVNSubnet := func(defaultSubnet *net.IPNet) (*net.IPNet, bool) {
		if n.OVNKubernetesConfig == nil {
			return defaultSubnet, false
		}

		ovnConfig := n.OVNKubernetesConfig

		// Since each subnet has a unique non-overlapping CIDR,
		// we can use that to distinguish the type of subnet without having to define extra constants.
		switch defaultSubnet.String() {
		case validate.OVNIPv4JoinSubnet.String():
			if ovnConfig.IPv4 != nil && ovnConfig.IPv4.InternalJoinSubnet != nil {
				return &ovnConfig.IPv4.InternalJoinSubnet.IPNet, true
			}
		default:
		}
		return defaultSubnet, false
	}

	// We only check against OVNKubernetes default subnets.
	// Any overrides of default subnets is validated in func validateOVNKubernetesConfig.
	subnetsCheck := func(joinSubnet, transitSubnet, masqueradeSubnet *net.IPNet) {
		// Join subnet
		if ovnsubnet, configured := getOVNSubnet(joinSubnet); !configured && validate.DoCIDRsOverlap(network, ovnsubnet) {
			allErrs = append(allErrs, field.Invalid(fldPath, network.String(), fmt.Sprintf("must not overlap with OVNKubernetes default internal subnet %s", ovnsubnet.String())))
		}

		// Transit subnet
		if ovnsubnet, configured := getOVNSubnet(transitSubnet); !configured && validate.DoCIDRsOverlap(network, ovnsubnet) {
			allErrs = append(allErrs, field.Invalid(fldPath, network.String(), fmt.Sprintf("must not overlap with OVNKubernetes default transit subnet %s", ovnsubnet.String())))
		}

		// Masquerade subnet
		if ovnsubnet, configured := getOVNSubnet(masqueradeSubnet); !configured && validate.DoCIDRsOverlap(network, ovnsubnet) {
			allErrs = append(allErrs, field.Invalid(fldPath, network.String(), fmt.Sprintf("must not overlap with OVNKubernetes default masquerade subnet %s", ovnsubnet.String())))
		}
	}

	if network.IP.To4() != nil {
		subnetsCheck(validate.OVNIPv4JoinSubnet, validate.OVNIPv4TransitSubnet, validate.OVNIPv4MasqueradeSubnet)
	} else {
		subnetsCheck(validate.OVNIPv6JoinSubnet, validate.OVNIPv6TransitSubnet, validate.OVNIPv6MasqueradeSubnet)
	}

	return allErrs
}

func validateNetworkingClusterNetworkMTU(c *types.InstallConfig, fldPath *field.Path) field.ErrorList {
	// higherLimitMTUVPC is the MTU limit for AWS VPC.
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/network_mtu.html#jumbo_frame_instances
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/network_mtu.html
	const higherLimitMTUVPC uint32 = uint32(9001)

	// lowerLimitMTUVPC is the lower limit to prevent users setting too low values impacting in the
	// cluster network performance. Tested values with 1100 decreases 70% in the network performance
	// in AWS deployments:
	const lowerLimitMTUVPC uint32 = uint32(1000)

	// higherLimitMTUEdge defines the maximium generally supported MTU in AWS Local and Wavelength Zones.
	// Mostly AWS Local or Wavelength zones have limited MTU between those and in the Region.
	// It is required to raise a warning message when the user-defined MTU is higher than general supported.
	// https://docs.aws.amazon.com/local-zones/latest/ug/how-local-zones-work.html#considerations
	// https://docs.aws.amazon.com/wavelength/latest/developerguide/how-wavelengths-work.html
	const higherLimitMTUEdge uint32 = uint32(1300)

	// MTU overhead for the network plugin OVNKubernetes.
	// https://docs.openshift.com/container-platform/4.14/networking/changing-cluster-network-mtu.html#mtu-value-selection_changing-cluster-network-mtu
	const minOverheadOVN uint32 = uint32(100)

	allErrs := field.ErrorList{}

	if c.Networking == nil {
		return nil
	}

	if c.Networking.ClusterNetworkMTU == 0 {
		return nil
	}

	if c.Platform.Name() != aws.Name {
		return append(allErrs, field.Invalid(fldPath, int(c.Networking.ClusterNetworkMTU), "cluster network MTU is allowed only in AWS deployments"))
	}

	network := c.NetworkType
	mtu := c.Networking.ClusterNetworkMTU

	// Calculating the MTU limits considering the base overhead for each network plugin.
	limitEdgeOVNKubernetes := higherLimitMTUEdge - minOverheadOVN
	limitOVNKubernetes := higherLimitMTUVPC - minOverheadOVN

	if mtu > higherLimitMTUVPC {
		return append(allErrs, field.Invalid(fldPath, int(mtu), fmt.Sprintf("cluster network MTU exceeds the maximum value of %d", higherLimitMTUVPC)))
	}

	// Prevent too low MTU values.
	// Tests in AWS Local Zones with MTU of 1100 decreased the network
	// performance in 70%. The check protects the cluster stability from
	// user defining too lower numbers.
	// https://issues.redhat.com/browse/OCPBUGS-11098
	if mtu < lowerLimitMTUVPC {
		return append(allErrs, field.Invalid(fldPath, int(mtu), fmt.Sprintf("cluster network MTU is lower than the minimum value of %d", lowerLimitMTUVPC)))
	}

	hasEdgePool := false
	warnEdgePool := false
	for _, compute := range c.Compute {
		if compute.Name == types.MachinePoolEdgeRoleName {
			hasEdgePool = true
			break
		}
	}

	if network != string(operv1.NetworkTypeOVNKubernetes) {
		return append(allErrs, field.Invalid(fldPath, int(mtu), fmt.Sprintf("cluster network MTU is not valid with network plugin %s", network)))
	}

	if mtu > limitOVNKubernetes {
		return append(allErrs, field.Invalid(fldPath, int(mtu), fmt.Sprintf("cluster network MTU exceeds the maximum value with the network plugin %s of %d", network, limitOVNKubernetes)))
	}
	if hasEdgePool && (mtu > limitEdgeOVNKubernetes) {
		warnEdgePool = true
	}
	if warnEdgePool {
		logrus.Warnf("networking.ClusterNetworkMTU exceeds the maximum value generally supported by AWS Local or Wavelength zones. Please ensure all AWS Zones defined in the edge compute pool accepts the MTU %d bytes between nodes (EC2) in the zone and in the Region.", mtu)
	}

	return allErrs
}

func validateOVNKubernetesConfig(n *types.Networking, fldPath *field.Path) field.ErrorList {
	if n.OVNKubernetesConfig == nil {
		return nil
	}

	allErrs := field.ErrorList{}

	if !strings.EqualFold(n.NetworkType, string(operv1.NetworkTypeOVNKubernetes)) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("networkType"), n.NetworkType, "ovnKubernetesConfig may only be specified with the OVNKubernetes networkType"))
	}

	allErrs = append(allErrs, validateOVNIPv4InternalJoinSubnet(n, fldPath.Child("ovnKubernetesConfig", "ipv4", "internalJoinSubnet"))...)
	return allErrs
}

func validateOVNIPv4InternalJoinSubnet(n *types.Networking, fldPath *field.Path) field.ErrorList {
	if ipv4 := n.OVNKubernetesConfig.IPv4; ipv4 == nil || ipv4.InternalJoinSubnet == nil {
		return nil
	}

	allErrs := field.ErrorList{}
	ipv4JoinNet := n.OVNKubernetesConfig.IPv4.InternalJoinSubnet

	if err := validate.SubnetCIDR(&ipv4JoinNet.IPNet); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, ipv4JoinNet.IPNet.String(), err.Error()))
		return allErrs // CIDR value is invalid, so we cannot perform further validation.
	}

	for _, net := range n.ClusterNetwork {
		if validate.DoCIDRsOverlap(&ipv4JoinNet.IPNet, &net.CIDR.IPNet) {
			errMsg := fmt.Sprintf("must not overlap with clusterNetwork %s", net.CIDR.String())
			allErrs = append(allErrs, field.Invalid(fldPath, ipv4JoinNet.String(), errMsg))
		}
	}

	for _, net := range n.MachineNetwork {
		if validate.DoCIDRsOverlap(&ipv4JoinNet.IPNet, &net.CIDR.IPNet) {
			errMsg := fmt.Sprintf("must not overlap with machineNetwork %s", net.CIDR.String())
			allErrs = append(allErrs, field.Invalid(fldPath, ipv4JoinNet.String(), errMsg))
		}
	}

	for _, net := range n.ServiceNetwork {
		if validate.DoCIDRsOverlap(&ipv4JoinNet.IPNet, &net.IPNet) {
			errMsg := fmt.Sprintf("must not overlap with serviceNetwork %s", net.String())
			allErrs = append(allErrs, field.Invalid(fldPath, ipv4JoinNet.String(), errMsg))
		}
	}

	if largeEnough, err := isV4NodeSubnetLargeEnough(n.ClusterNetwork, ipv4JoinNet); err == nil && !largeEnough {
		errMsg := `ipv4InternalJoinSubnet is not large enough for the maximum number of nodes which can be supported by ClusterNetwork`
		allErrs = append(allErrs, field.Invalid(fldPath, ipv4JoinNet.String(), errMsg))
	} else if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath, err))
	}

	return allErrs
}

func validateControlPlane(platform *types.Platform, pool *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if pool.Name != types.MachinePoolControlPlaneRoleName {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("name"), pool.Name, []string{types.MachinePoolControlPlaneRoleName}))
	}
	if pool.Replicas != nil && *pool.Replicas == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("replicas"), pool.Replicas, "number of control plane replicas must be positive"))
	}
	allErrs = append(allErrs, ValidateMachinePool(platform, pool, fldPath)...)
	return allErrs
}

func validateArbiter(platform *types.Platform, arbiterPool, masterPool *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if platform != nil && platform.BareMetal == nil {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("platform"), platform.Name(), []string{baremetal.Name}))
	}
	if arbiterPool.Name != types.MachinePoolArbiterRoleName {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("name"), arbiterPool.Name, []string{types.MachinePoolArbiterRoleName}))
	}
	if arbiterPool.Replicas != nil && *arbiterPool.Replicas == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("replicas"), arbiterPool.Replicas, "number of arbiter replicas must be positive"))
	}
	if masterPool.Replicas == nil || *masterPool.Replicas < 2 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("replicas"), masterPool.Replicas, "number of controlPlane replicas must be at least 2 for arbiter deployments"))
	}
	allErrs = append(allErrs, ValidateMachinePool(platform, arbiterPool, fldPath)...)
	return allErrs
}

func validateComputeEdge(platform *types.Platform, pName string, fldPath *field.Path, pfld *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if platform.Name() != aws.Name {
		allErrs = append(allErrs, field.NotSupported(pfld.Child("name"), pName, []string{types.MachinePoolComputeRoleName}))
	}

	return allErrs
}

func validateCompute(platform *types.Platform, control *types.MachinePool, pools []types.MachinePool, fldPath *field.Path, isMultiArchEnabled bool) field.ErrorList {
	allErrs := field.ErrorList{}
	poolNames := map[string]bool{}
	for i, p := range pools {
		poolFldPath := fldPath.Index(i)
		switch p.Name {
		case types.MachinePoolComputeRoleName:
		case types.MachinePoolEdgeRoleName:
			allErrs = append(allErrs, validateComputeEdge(platform, p.Name, poolFldPath, poolFldPath)...)
		default:
			allErrs = append(allErrs, field.NotSupported(poolFldPath.Child("name"), p.Name, []string{types.MachinePoolComputeRoleName, types.MachinePoolEdgeRoleName}))
		}

		if poolNames[p.Name] {
			allErrs = append(allErrs, field.Duplicate(poolFldPath.Child("name"), p.Name))
		}
		poolNames[p.Name] = true
		if control != nil && control.Architecture != p.Architecture && !isMultiArchEnabled {
			allErrs = append(allErrs, field.Invalid(poolFldPath.Child("architecture"), p.Architecture, "heteregeneous multi-arch is not supported; compute pool architecture must match control plane"))
		}
		allErrs = append(allErrs, ValidateMachinePool(platform, &p, poolFldPath)...)
	}
	return allErrs
}

// vips defines the VIPs to validate
type vips struct {
	API     []string
	Ingress []string
}

// vipFields defines the field names to which validation errors for each VIP
// type should be assigned to
type vipFields struct {
	APIVIPs     string
	IngressVIPs string
}

// validateVIPsForPlatform validates the VIPs (for API and Ingress) for the
// given platform
func validateVIPsForPlatform(network *types.Networking, platform *types.Platform, usingAgentMethod bool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	virtualIPs := vips{}
	newVIPsFields := vipFields{
		APIVIPs:     "apiVIPs",
		IngressVIPs: "ingressVIPs",
	}

	var lbType configv1.PlatformLoadBalancerType

	switch {
	case platform.BareMetal != nil:
		virtualIPs = vips{
			API:     platform.BareMetal.APIVIPs,
			Ingress: platform.BareMetal.IngressVIPs,
		}

		if platform.BareMetal.LoadBalancer != nil {
			lbType = platform.BareMetal.LoadBalancer.Type
		}

		allErrs = append(allErrs, validateAPIAndIngressVIPs(virtualIPs, newVIPsFields, true, true, lbType, network, fldPath.Child(baremetal.Name))...)
	case platform.Nutanix != nil:
		allErrs = append(allErrs, ensureIPv4IsFirstInDualStackSlice(&platform.Nutanix.APIVIPs, fldPath.Child(nutanix.Name, newVIPsFields.APIVIPs))...)
		allErrs = append(allErrs, ensureIPv4IsFirstInDualStackSlice(&platform.Nutanix.IngressVIPs, fldPath.Child(nutanix.Name, newVIPsFields.IngressVIPs))...)

		virtualIPs = vips{
			API:     platform.Nutanix.APIVIPs,
			Ingress: platform.Nutanix.IngressVIPs,
		}

		if platform.Nutanix.LoadBalancer != nil {
			lbType = platform.Nutanix.LoadBalancer.Type
		}

		allErrs = append(allErrs, validateAPIAndIngressVIPs(virtualIPs, newVIPsFields, false, false, lbType, network, fldPath.Child(nutanix.Name))...)
	case platform.OpenStack != nil:
		virtualIPs = vips{
			API:     platform.OpenStack.APIVIPs,
			Ingress: platform.OpenStack.IngressVIPs,
		}

		if platform.OpenStack.LoadBalancer != nil {
			lbType = platform.OpenStack.LoadBalancer.Type
		}

		allErrs = append(allErrs, validateAPIAndIngressVIPs(virtualIPs, newVIPsFields, true, true, lbType, network, fldPath.Child(openstack.Name))...)
	case platform.VSphere != nil:
		virtualIPs = vips{
			API:     platform.VSphere.APIVIPs,
			Ingress: platform.VSphere.IngressVIPs,
		}

		if platform.VSphere.LoadBalancer != nil {
			lbType = platform.VSphere.LoadBalancer.Type
		}

		vipIsRequired, reqVIPinMachineCIDR := usingAgentMethod, usingAgentMethod
		allErrs = append(allErrs, validateAPIAndIngressVIPs(virtualIPs, newVIPsFields, vipIsRequired, reqVIPinMachineCIDR, lbType, network, fldPath.Child(vsphere.Name))...)
	case platform.Ovirt != nil:
		allErrs = append(allErrs, ensureIPv4IsFirstInDualStackSlice(&platform.Ovirt.APIVIPs, fldPath.Child(ovirt.Name, newVIPsFields.APIVIPs))...)
		allErrs = append(allErrs, ensureIPv4IsFirstInDualStackSlice(&platform.Ovirt.IngressVIPs, fldPath.Child(ovirt.Name, newVIPsFields.IngressVIPs))...)

		newVIPsFields = vipFields{
			APIVIPs:     "api_vips",
			IngressVIPs: "ingress_vips",
		}
		virtualIPs = vips{
			API:     platform.Ovirt.APIVIPs,
			Ingress: platform.Ovirt.IngressVIPs,
		}

		if platform.Ovirt.LoadBalancer != nil {
			lbType = platform.Ovirt.LoadBalancer.Type
		}

		allErrs = append(allErrs, validateAPIAndIngressVIPs(virtualIPs, newVIPsFields, true, true, lbType, network, fldPath.Child(ovirt.Name))...)
	default:
		//no vips to validate on this platform
	}

	return allErrs
}

func ensureIPv4IsFirstInDualStackSlice(vips *[]string, fldPath *field.Path) field.ErrorList {
	errList := field.ErrorList{}
	isDualStack, err := utilsnet.IsDualStackIPStrings(*vips)
	if err != nil {
		errList = append(errList, field.Invalid(fldPath, vips, err.Error()))
		return errList
	}

	if isDualStack {
		if len(*vips) == 2 {
			if utilsnet.IsIPv4String((*vips)[1]) && utilsnet.IsIPv6String((*vips)[0]) {
				(*vips)[0], (*vips)[1] = (*vips)[1], (*vips)[0]
			}
		} else {
			errList = append(errList, field.Invalid(fldPath, vips, "wrong number of VIPs given. Expecting 2 VIPs for dual stack"))
			return errList
		}
	}

	return errList
}

// validateAPIAndIngressVIPs validates the API and Ingress VIPs
//
//nolint:gocyclo
func validateAPIAndIngressVIPs(vips vips, fieldNames vipFields, vipIsRequired, reqVIPinMachineCIDR bool, lbType configv1.PlatformLoadBalancerType, n *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(vips.API) == 0 {
		if vipIsRequired {
			allErrs = append(allErrs, field.Required(fldPath.Child(fieldNames.APIVIPs), "must specify at least one VIP for the API"))
		}
	} else if len(vips.API) <= 2 {
		for _, vip := range vips.API {
			if err := validate.IP(vip); err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child(fieldNames.APIVIPs), vip, err.Error()))
			}

			// When using user-managed loadbalancer we do not require API and Ingress VIP to be different as well as
			// we allow them to be from outside the machine network CIDR.
			if lbType != configv1.LoadBalancerTypeUserManaged {
				for _, ingressVIP := range vips.Ingress {
					apiVIPNet := net.ParseIP(vip)
					ingressVIPNet := net.ParseIP(ingressVIP)
					if apiVIPNet != nil && apiVIPNet.Equal(ingressVIPNet) {
						allErrs = append(allErrs, field.Invalid(fldPath.Child(fieldNames.APIVIPs), vip, "VIP for API must not be one of the Ingress VIPs"))
					}
				}

				if err := ValidateIPinMachineCIDR(vip, n); reqVIPinMachineCIDR && err != nil {
					allErrs = append(allErrs, field.Invalid(fldPath.Child(fieldNames.APIVIPs), vip, err.Error()))
				}
			}
		}

		if len(vips.Ingress) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child(fieldNames.IngressVIPs), "must specify VIP for ingress, when VIP for API is set"))
		}

		if len(vips.API) == 1 {
			hasIPv4, hasIPv6, presence, _ := inferIPVersionFromInstallConfig(n)

			apiVIPIPFamily := corev1.IPv4Protocol
			if utilsnet.IsIPv6String(vips.API[0]) {
				apiVIPIPFamily = corev1.IPv6Protocol
			}

			if hasIPv4 && hasIPv6 && apiVIPIPFamily != presence["machineNetwork"].Primary {
				allErrs = append(allErrs, field.Invalid(fldPath.Child(fieldNames.APIVIPs), vips.API[0], "VIP for the API must be of the same IP family with machine network's primary IP Family for dual-stack IPv4/IPv6"))
			}
		} else if len(vips.API) == 2 {
			if isDualStack, _ := utilsnet.IsDualStackIPStrings(vips.API); !isDualStack {
				allErrs = append(allErrs, field.Invalid(fldPath.Child(fieldNames.APIVIPs), vips.API, "If two API VIPs are given, one must be an IPv4 address, the other an IPv6"))
			}
		}
	} else {
		allErrs = append(allErrs, field.TooMany(fldPath.Child(fieldNames.APIVIPs), len(vips.API), 2))
	}

	if len(vips.Ingress) == 0 {
		if vipIsRequired {
			allErrs = append(allErrs, field.Required(fldPath.Child(fieldNames.IngressVIPs), "must specify at least one VIP for the Ingress"))
		}
	} else if len(vips.Ingress) <= 2 {
		for _, vip := range vips.Ingress {
			if err := validate.IP(vip); err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child(fieldNames.IngressVIPs), vip, err.Error()))
			}

			// When using user-managed loadbalancer we do not require API and Ingress VIP to be different as well as
			// we allow them to be from outside the machine network CIDR.
			if lbType != configv1.LoadBalancerTypeUserManaged {
				if err := ValidateIPinMachineCIDR(vip, n); reqVIPinMachineCIDR && err != nil {
					allErrs = append(allErrs, field.Invalid(fldPath.Child(fieldNames.IngressVIPs), vip, err.Error()))
				}
			}
		}

		if len(vips.API) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child(fieldNames.APIVIPs), "must specify VIP for API, when VIP for ingress is set"))
		}

		if len(vips.Ingress) == 1 {
			hasIPv4, hasIPv6, presence, _ := inferIPVersionFromInstallConfig(n)

			ingressVIPIPFamily := corev1.IPv4Protocol
			if utilsnet.IsIPv6String(vips.Ingress[0]) {
				ingressVIPIPFamily = corev1.IPv6Protocol
			}

			if hasIPv4 && hasIPv6 && ingressVIPIPFamily != presence["machineNetwork"].Primary {
				allErrs = append(allErrs, field.Invalid(fldPath.Child(fieldNames.IngressVIPs), vips.Ingress[0], "VIP for the Ingress must be of the same IP family with machine network's primary IP Family for dual-stack IPv4/IPv6"))
			}
		} else if len(vips.Ingress) == 2 {
			if isDualStack, _ := utilsnet.IsDualStackIPStrings(vips.Ingress); !isDualStack {
				allErrs = append(allErrs, field.Invalid(fldPath.Child(fieldNames.IngressVIPs), vips.Ingress, "If two Ingress VIPs are given, one must be an IPv4 address, the other an IPv6"))
			}
		}
	} else {
		allErrs = append(allErrs, field.TooMany(fldPath.Child(fieldNames.IngressVIPs), len(vips.Ingress), 2))
	}

	return allErrs
}

// ValidateIPinMachineCIDR confirms if the specified VIP is in the machine CIDR.
func ValidateIPinMachineCIDR(vip string, n *types.Networking) error {
	var networks []string

	for _, network := range n.MachineNetwork {
		if network.CIDR.Contains(net.ParseIP(vip)) {
			return nil
		}
		networks = append(networks, network.CIDR.String())
	}

	return fmt.Errorf("IP expected to be in one of the machine networks: %s", strings.Join(networks, ","))
}

func validatePlatform(platform *types.Platform, usingAgentMethod bool, fldPath *field.Path, network *types.Networking, c *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	activePlatform := platform.Name()
	platforms := make([]string, len(types.PlatformNames))
	copy(platforms, types.PlatformNames)
	platforms = append(platforms, types.HiddenPlatformNames...)
	sort.Strings(platforms)
	i := sort.SearchStrings(platforms, activePlatform)
	if i == len(platforms) || platforms[i] != activePlatform {
		allErrs = append(allErrs, field.Invalid(fldPath, activePlatform, fmt.Sprintf("must specify one of the platforms (%s)", strings.Join(platforms, ", "))))
	}
	validate := func(n string, value interface{}, validation func(*field.Path) field.ErrorList) {
		if n != activePlatform {
			allErrs = append(allErrs, field.Invalid(fldPath, activePlatform, fmt.Sprintf("must only specify a single type of platform; cannot use both %q and %q", activePlatform, n)))
		}
		allErrs = append(allErrs, validation(fldPath.Child(n))...)
	}
	if platform.AWS != nil {
		validate(aws.Name, platform.AWS, func(f *field.Path) field.ErrorList {
			return awsvalidation.ValidatePlatform(platform.AWS, c.Publish, c.CredentialsMode, f)
		})
	}
	if platform.Azure != nil {
		validate(azure.Name, platform.Azure, func(f *field.Path) field.ErrorList {
			return azurevalidation.ValidatePlatform(platform.Azure, c.Publish, f, c)
		})
	}
	if platform.GCP != nil {
		validate(gcp.Name, platform.GCP, func(f *field.Path) field.ErrorList { return gcpvalidation.ValidatePlatform(platform.GCP, f, c) })
	}
	if platform.IBMCloud != nil {
		validate(ibmcloud.Name, platform.IBMCloud, func(f *field.Path) field.ErrorList { return ibmcloudvalidation.ValidatePlatform(platform.IBMCloud, f) })
	}
	if platform.OpenStack != nil {
		validate(openstack.Name, platform.OpenStack, func(f *field.Path) field.ErrorList {
			return openstackvalidation.ValidatePlatform(platform.OpenStack, network, f, c)
		})
	}
	if platform.PowerVS != nil {
		if c.SSHKey == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("sshKey"), "sshKey is required"))
		}
		validate(powervs.Name, platform.PowerVS, func(f *field.Path) field.ErrorList {
			return powervsvalidation.ValidatePlatform(platform.PowerVS, f)
		})
	}
	if platform.VSphere != nil {
		validate(vsphere.Name, platform.VSphere, func(f *field.Path) field.ErrorList {
			return vspherevalidation.ValidatePlatform(platform.VSphere, usingAgentMethod, f, c)
		})
	}
	if platform.BareMetal != nil {
		validate(baremetal.Name, platform.BareMetal, func(f *field.Path) field.ErrorList {
			return baremetalvalidation.ValidatePlatform(platform.BareMetal, usingAgentMethod, network, f, c)
		})
	}
	if platform.Ovirt != nil {
		validate(ovirt.Name, platform.Ovirt, func(f *field.Path) field.ErrorList {
			return ovirtvalidation.ValidatePlatform(platform.Ovirt, f, c)
		})
	}
	if platform.Nutanix != nil {
		validate(nutanix.Name, platform.Nutanix, func(f *field.Path) field.ErrorList {
			return nutanixvalidation.ValidatePlatform(platform.Nutanix, f, c)
		})
	}
	return allErrs
}

func validateProxy(p *types.Proxy, c *types.InstallConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.HTTPProxy == "" && p.HTTPSProxy == "" {
		allErrs = append(allErrs, field.Required(fldPath, "must include httpProxy or httpsProxy"))
	}

	if p.HTTPProxy != "" {
		allErrs = append(allErrs, validateURI(p.HTTPProxy, fldPath.Child("httpProxy"), []string{"http"})...)
		if c.Networking != nil {
			allErrs = append(allErrs, validateIPProxy(p.HTTPProxy, c.Networking, fldPath.Child("httpProxy"))...)
		}
	}
	if p.HTTPSProxy != "" {
		allErrs = append(allErrs, validateURI(p.HTTPSProxy, fldPath.Child("httpsProxy"), []string{"http", "https"})...)
		if c.Networking != nil {
			allErrs = append(allErrs, validateIPProxy(p.HTTPSProxy, c.Networking, fldPath.Child("httpsProxy"))...)
		}
	}
	if p.NoProxy != "" && p.NoProxy != "*" {
		if strings.Contains(p.NoProxy, " ") {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("noProxy"), p.NoProxy, fmt.Sprintf("noProxy must not have spaces")))
		}
		for idx, v := range strings.Split(p.NoProxy, ",") {
			v = strings.TrimSpace(v)
			errDomain := validate.NoProxyDomainName(v)
			_, _, errCIDR := net.ParseCIDR(v)
			ip := net.ParseIP(v)
			if errDomain != nil && errCIDR != nil && ip == nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("noProxy"), p.NoProxy, fmt.Sprintf(
					"each element of noProxy must be a IP, CIDR or domain without wildcard characters, which is violated by element %d %q", idx, v)))
			}
		}
	}

	return allErrs
}

func validateImageContentSources(groups []types.ImageContentSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for gidx, group := range groups {
		groupf := fldPath.Index(gidx)
		if err := validateNamedRepository(group.Source); err != nil {
			allErrs = append(allErrs, field.Invalid(groupf.Child("source"), group.Source, err.Error()))
		}

		for midx, mirror := range group.Mirrors {
			if err := validateNamedRepository(mirror); err != nil {
				allErrs = append(allErrs, field.Invalid(groupf.Child("mirrors").Index(midx), mirror, err.Error()))
				continue
			}
		}
	}
	return allErrs
}

func validateImageDigestSources(groups []types.ImageDigestSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for gidx, group := range groups {
		groupf := fldPath.Index(gidx)
		if err := validateNamedRepository(group.Source); err != nil {
			allErrs = append(allErrs, field.Invalid(groupf.Child("source"), group.Source, err.Error()))
		}

		for midx, mirror := range group.Mirrors {
			if err := validateNamedRepository(mirror); err != nil {
				allErrs = append(allErrs, field.Invalid(groupf.Child("mirrors").Index(midx), mirror, err.Error()))
				continue
			}
		}
	}
	return allErrs
}

func validateNamedRepository(r string) error {
	ref, err := dockerref.ParseNamed(r)
	if err != nil {
		// If a mirror name is provided without the named reference,
		// then the name is not considered canonical and will cause
		// an error. e.g. registry.lab.redhat.com:5000 will result
		// in an error. Instead we will check whether the input is
		// a valid hostname as a workaround.
		if err == dockerref.ErrNameNotCanonical {
			// If the hostname string contains a port, lets attempt
			// to split them
			host, _, err := net.SplitHostPort(r)
			if err != nil {
				host = r
			}
			if err = validate.Host(host); err != nil {
				return errors.Wrap(err, "the repository provided is invalid")
			}
			return nil
		}
		return errors.Wrap(err, "failed to parse")
	}
	if !dockerref.IsNameOnly(ref) {
		return errors.New("must be repository--not reference")
	}
	return nil
}

var (
	validPublishingStrategies = map[types.PublishingStrategy]struct{}{
		types.ExternalPublishingStrategy: {},
		types.InternalPublishingStrategy: {},
		types.MixedPublishingStrategy:    {},
	}

	validPublishingStrategyValues = func() []string {
		v := make([]string, 0, len(validPublishingStrategies))
		for m := range validPublishingStrategies {
			v = append(v, string(m))
		}
		sort.Strings(v)
		return v
	}()
)

func validateCloudCredentialsMode(mode types.CredentialsMode, fldPath *field.Path, platform types.Platform) field.ErrorList {
	if mode == "" {
		return nil
	}
	allErrs := field.ErrorList{}

	allowedAzureModes := []types.CredentialsMode{types.PassthroughCredentialsMode, types.ManualCredentialsMode}
	if platform.Azure != nil && platform.Azure.CloudName == azure.StackCloud {
		allowedAzureModes = []types.CredentialsMode{types.ManualCredentialsMode}
	}

	// validPlatformCredentialsModes is a map from the platform name to a slice of credentials modes that are valid
	// for the platform. If a platform name is not in the map, then the credentials mode cannot be set for that platform.
	validPlatformCredentialsModes := map[string][]types.CredentialsMode{
		aws.Name:      {types.MintCredentialsMode, types.PassthroughCredentialsMode, types.ManualCredentialsMode},
		azure.Name:    allowedAzureModes,
		gcp.Name:      {types.MintCredentialsMode, types.PassthroughCredentialsMode, types.ManualCredentialsMode},
		ibmcloud.Name: {types.ManualCredentialsMode},
		powervs.Name:  {types.ManualCredentialsMode},
		nutanix.Name:  {types.ManualCredentialsMode},
	}
	if validModes, ok := validPlatformCredentialsModes[platform.Name()]; ok {
		validModesSet := sets.NewString()
		for _, m := range validModes {
			validModesSet.Insert(string(m))
		}
		if !validModesSet.Has(string(mode)) {
			allErrs = append(allErrs, field.NotSupported(fldPath, mode, validModesSet.List()))
		}
	} else {
		allErrs = append(allErrs, field.Invalid(fldPath, mode, fmt.Sprintf("cannot be set when using the %q platform", platform.Name())))
	}
	return allErrs
}

// validateURI checks if the given url is of the right format. It also checks if the scheme of the uri
// provided is within the list of accepted schema provided as part of the input.
func validateURI(uri string, fldPath *field.Path, schemes []string) field.ErrorList {
	parsed, err := url.ParseRequestURI(uri)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, uri, err.Error())}
	}
	for _, scheme := range schemes {
		if scheme == parsed.Scheme {
			return nil
		}
	}
	return field.ErrorList{field.NotSupported(fldPath, parsed.Scheme, schemes)}
}

// validateIPProxy checks if the given proxy string is an IP and if so checks the service and
// cluster networks and returns error if the IP belongs in them. Returns nil if the proxy is
// not an IP address.
func validateIPProxy(proxy string, n *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	parsed, err := url.ParseRequestURI(proxy)
	if err != nil {
		return allErrs
	}

	proxyIP := net.ParseIP(parsed.Hostname())
	if proxyIP == nil {
		return nil
	}

	for _, network := range n.ClusterNetwork {
		if network.CIDR.Contains(proxyIP) {
			allErrs = append(allErrs, field.Invalid(fldPath, proxy, "proxy value is part of the cluster networks"))
			break
		}
	}

	for _, network := range n.ServiceNetwork {
		if network.Contains(proxyIP) {
			allErrs = append(allErrs, field.Invalid(fldPath, proxy, "proxy value is part of the service networks"))
			break
		}
	}

	if ovnCfg := n.OVNKubernetesConfig; ovnCfg != nil && ovnCfg.IPv4 != nil && ovnCfg.IPv4.InternalJoinSubnet != nil && ovnCfg.IPv4.InternalJoinSubnet.Contains(proxyIP) {
		allErrs = append(allErrs, field.Invalid(fldPath, proxy, "proxy value is part of the ovn-kubernetes IPv4 InternalJoinSubnet"))
	}
	return allErrs
}

// validateFIPSconfig checks if the current install-config is compatible with FIPS standards
// and returns an error if it's not the case. As of this writing, only rsa or ecdsa algorithms are supported
// for ssh keys on FIPS.
func validateFIPSconfig(c *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	if c.SSHKey != "" {
		sshParsedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(c.SSHKey))
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("sshKey"), c.SSHKey, fmt.Sprintf("Fatal error trying to parse configured public key: %s", err)))
		} else {
			sshKeyType := sshParsedKey.Type()
			re := regexp.MustCompile(`^ecdsa-sha2-nistp\d{3}$|^ssh-rsa$`)
			if !re.MatchString(sshKeyType) {
				allErrs = append(allErrs, field.Invalid(field.NewPath("sshKey"), c.SSHKey, fmt.Sprintf("SSH key type %s unavailable when FIPS is enabled. Please use rsa or ecdsa.", sshKeyType)))
			}
		}
	}

	if err := hostcrypt.VerifyHostTargetState(c.FIPS); err != nil {
		if skip, ok := os.LookupEnv("OPENSHIFT_INSTALL_SKIP_HOSTCRYPT_VALIDATION"); ok && skip != "" {
			logrus.Warnf("%v", err)
			if c.Annotations == nil {
				c.Annotations = make(map[string]string)
			}
			c.Annotations[hostCryptBypassedAnnotation] = "true"
		} else {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("fips"), err.Error()))
		}
	}
	return allErrs
}

// validateCapabilities checks if additional, optional OpenShift components are specified in the
// install-config to be included in the installation.
func validateCapabilities(c *types.Capabilities, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allCapabilitySets := sets.NewString()
	allAvailableCapabilities := sets.NewString()
	// Create sets of all capability sets and *all* available capabilities across those capability sets
	for baselineSet, capabilities := range configv1.ClusterVersionCapabilitySets {
		allCapabilitySets.Insert(string(baselineSet))
		for _, capability := range capabilities {
			allAvailableCapabilities.Insert(string(capability))
		}
	}

	if !allCapabilitySets.Has(string(c.BaselineCapabilitySet)) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("baselineCapabilitySet"), c.BaselineCapabilitySet, allCapabilitySets.List()))
	}

	// Check to see the validity of additionalEnabledCapabilities specified by the user
	for i, capability := range c.AdditionalEnabledCapabilities {
		if !allAvailableCapabilities.Has(string(capability)) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("additionalEnabledCapabilities").Index(i), capability, allAvailableCapabilities.List()))
		}
	}
	return allErrs
}

func validateAdditionalCABundlePolicy(c *types.InstallConfig) error {
	switch c.AdditionalTrustBundlePolicy {
	case types.PolicyProxyOnly, types.PolicyAlways:
		return nil
	default:
		return fmt.Errorf("supported values \"Proxyonly\", \"Always\"")
	}
}

// ValidateFeatureSet returns an error if a gated feature is used without opting into the feature set.
func ValidateFeatureSet(c *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	clusterProfile := types.GetClusterProfileName()
	featureSets, ok := features.AllFeatureSets()[clusterProfile]
	if !ok {
		logrus.Warnf("no feature sets for cluster profile %q", clusterProfile)
	}
	if _, ok := featureSets[c.FeatureSet]; c.FeatureSet != configv1.CustomNoUpgrade && !ok {
		sortedFeatureSets := func() []string {
			v := []string{}
			for n := range features.AllFeatureSets()[clusterProfile] {
				v = append(v, string(n))
			}
			// Add CustomNoUpgrade since it is not part of features sets for profiles
			v = append(v, string(configv1.CustomNoUpgrade))
			sort.Strings(v)
			return v
		}()
		allErrs = append(allErrs, field.NotSupported(field.NewPath("featureSet"), c.FeatureSet, sortedFeatureSets))
	}

	if len(c.FeatureGates) > 0 {
		if c.FeatureSet != configv1.CustomNoUpgrade {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("featureGates"), "featureGates can only be used with the CustomNoUpgrade feature set"))
		}
		allErrs = append(allErrs, validateCustomFeatureGates(c)...)
	}

	// We can only accurately check gated features
	// if feature sets are correctly configured.
	if len(allErrs) == 0 {
		allErrs = append(allErrs, validateGatedFeatures(c)...)
	}

	return allErrs
}

// validateCustomFeatureGates checks that all provided custom features match the expected format.
// The expected format is <FeatureName>=<Enabled>.
func validateCustomFeatureGates(c *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	for i, rawFeature := range c.FeatureGates {
		featureParts := strings.Split(rawFeature, "=")
		if len(featureParts) != 2 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("featureGates").Index(i), rawFeature, "must match the format <feature-name>=<bool>"))
			continue
		}

		if _, err := strconv.ParseBool(featureParts[1]); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("featureGates").Index(i), rawFeature, "must match the format <feature-name>=<bool>, could not parse boolean value"))
		}
	}

	return allErrs
}

// validateGatedFeatures ensures that any gated features used in
// the install config are enabled.
func validateGatedFeatures(c *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	gatedFeatures := []featuregates.GatedInstallConfigFeature{}
	switch {
	case c.GCP != nil:
		gatedFeatures = append(gatedFeatures, gcpvalidation.GatedFeatures(c)...)
	case c.VSphere != nil:
		gatedFeatures = append(gatedFeatures, vspherevalidation.GatedFeatures(c)...)
	case c.AWS != nil:
		gatedFeatures = append(gatedFeatures, awsvalidation.GatedFeatures(c)...)
	}

	fg := c.EnabledFeatureGates()
	errMsgTemplate := "this field is protected by the %s feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set"

	fgCheck := func(c featuregates.GatedInstallConfigFeature) {
		if !fg.Enabled(c.FeatureGateName) && c.Condition {
			errMsg := fmt.Sprintf(errMsgTemplate, c.FeatureGateName)
			allErrs = append(allErrs, field.Forbidden(c.Field, errMsg))
		}
	}

	for _, gf := range gatedFeatures {
		fgCheck(gf)
	}

	return allErrs
}

// validateReleaseArchitecture ensures a compatible payload is used according to the desired architecture of the cluster.
func validateReleaseArchitecture(controlPlanePool *types.MachinePool, computePool []types.MachinePool, releaseArch types.Architecture) field.ErrorList {
	allErrs := field.ErrorList{}

	clusterArch := version.DefaultArch()
	if controlPlanePool != nil && controlPlanePool.Architecture != "" {
		clusterArch = controlPlanePool.Architecture
	}

	switch releaseArch {
	case "multi":
		// All good
	case "unknown":
		for _, p := range computePool {
			if p.Architecture != "" && clusterArch != p.Architecture {
				// Overriding release architecture is a must during dev/CI so just log a warning instead of erroring out
				logrus.Warnln("Could not determine release architecture for multi arch cluster configuration. Make sure the release is a multi architecture payload.")
				break
			}
		}
	default:
		if clusterArch != releaseArch {
			errMsg := fmt.Sprintf("cannot create %s controlPlane node from a single architecture %s release payload", clusterArch, releaseArch)
			allErrs = append(allErrs, field.Invalid(field.NewPath("controlPlane", "architecture"), clusterArch, errMsg))
		}
		for i, p := range computePool {
			poolFldPath := field.NewPath("compute").Index(i)
			if p.Architecture != "" && p.Architecture != releaseArch {
				errMsg := fmt.Sprintf("cannot create %s compute node from a single architecture %s release payload", p.Architecture, releaseArch)
				allErrs = append(allErrs, field.Invalid(poolFldPath.Child("architecture"), p.Architecture, errMsg))
			}
		}
	}

	return allErrs
}

// isV4NodeSubnetLargeEnough ensures the subnet is large enough for the maximum number of nodes supported by ClusterNetwork.
// This validation is performed by the cluster network operator: https://github.com/openshift/cluster-network-operator/blob/6b615be1447aa79252ddc73b10675b4638ae13f7/pkg/network/ovn_kubernetes.go#L1761.
// We need to duplicate it here to catch any issues with network customization prior to install.
func isV4NodeSubnetLargeEnough(cn []types.ClusterNetworkEntry, nodeSubnet *ipnet.IPNet) (bool, error) {
	var maxNodesNum int
	addrLen := 32
	for i, n := range cn {
		if utilsnet.IsIPv6CIDRString(n.CIDR.String()) {
			continue
		}

		mask, _ := n.CIDR.Mask.Size()
		if int(n.HostPrefix) < mask {
			return false, fmt.Errorf("cannot determine the number of nodes supported by cluster network %d due to invalid hostPrefix", i)
		}
		nodesNum := 1 << (int(n.HostPrefix) - mask)
		maxNodesNum += nodesNum
	}
	// We need to ensure each node can be assigned an IP address from the internal subnet
	intSubnetMask, _ := nodeSubnet.Mask.Size()

	// reserve one IP for the gw, one IP for network and one for broadcasting
	return maxNodesNum < (1<<(addrLen-intSubnetMask) - 3), nil
}
