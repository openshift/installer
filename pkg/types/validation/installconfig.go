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

	dockerref "github.com/containers/image/docker/reference"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	utilsnet "k8s.io/utils/net"

	configv1 "github.com/openshift/api/config/v1"
	operv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/hostcrypt"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/alibabacloud"
	alibabacloudvalidation "github.com/openshift/installer/pkg/types/alibabacloud/validation"
	"github.com/openshift/installer/pkg/types/aws"
	awsvalidation "github.com/openshift/installer/pkg/types/aws/validation"
	"github.com/openshift/installer/pkg/types/azure"
	azurevalidation "github.com/openshift/installer/pkg/types/azure/validation"
	"github.com/openshift/installer/pkg/types/baremetal"
	baremetalvalidation "github.com/openshift/installer/pkg/types/baremetal/validation"
	"github.com/openshift/installer/pkg/types/featuregates"
	"github.com/openshift/installer/pkg/types/gcp"
	gcpvalidation "github.com/openshift/installer/pkg/types/gcp/validation"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	ibmcloudvalidation "github.com/openshift/installer/pkg/types/ibmcloud/validation"
	"github.com/openshift/installer/pkg/types/libvirt"
	libvirtvalidation "github.com/openshift/installer/pkg/types/libvirt/validation"
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
)

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
		allErrs = append(allErrs, validateNetworking(c.Networking, c.IsSingleNodeOpenShift(), field.NewPath("networking"))...)
		allErrs = append(allErrs, validateNetworkingIPVersion(c.Networking, &c.Platform)...)
		allErrs = append(allErrs, validateNetworkingForPlatform(c.Networking, &c.Platform, field.NewPath("networking"))...)
		allErrs = append(allErrs, validateNetworkingClusterNetworkMTU(c, field.NewPath("networking", "clusterNetworkMTU"))...)
		allErrs = append(allErrs, validateVIPsForPlatform(c.Networking, &c.Platform, field.NewPath("platform"))...)
	} else {
		allErrs = append(allErrs, field.Required(field.NewPath("networking"), "networking is required"))
	}
	allErrs = append(allErrs, validatePlatform(&c.Platform, usingAgentMethod, field.NewPath("platform"), c.Networking, c)...)
	if c.ControlPlane != nil {
		allErrs = append(allErrs, validateControlPlane(&c.Platform, c.ControlPlane, field.NewPath("controlPlane"))...)
	} else {
		allErrs = append(allErrs, field.Required(field.NewPath("controlPlane"), "controlPlane is required"))
	}
	allErrs = append(allErrs, validateCompute(&c.Platform, c.ControlPlane, c.Compute, field.NewPath("compute"))...)
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
		logrus.Warningln("imageContentSources is deprecated, please use ImageDigestSource")
	}
	allErrs = append(allErrs, validateCloudCredentialsMode(c.CredentialsMode, field.NewPath("credentialsMode"), c.Platform)...)
	if c.Capabilities != nil {
		allErrs = append(allErrs, validateCapabilities(c.Capabilities, field.NewPath("capabilities"))...)
	}

	if c.Publish == types.InternalPublishingStrategy {
		switch platformName := c.Platform.Name(); platformName {
		case aws.Name, azure.Name, gcp.Name, alibabacloud.Name, ibmcloud.Name, powervs.Name:
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
			if enabledCaps.Has(configv1.ClusterVersionCapabilityBaremetal) && !enabledCaps.Has(configv1.ClusterVersionCapabilityMachineAPI) {
				allErrs = append(allErrs, field.Invalid(field.NewPath("additionalEnabledCapabilities"), c.Capabilities.AdditionalEnabledCapabilities,
					"the baremetal capability requires the MachineAPI capability"))
			}
			if enabledCaps.Has(configv1.ClusterVersionCapabilityMarketplace) && !enabledCaps.Has(configv1.ClusterVersionCapabilityOperatorLifecycleManager) {
				allErrs = append(allErrs, field.Invalid(field.NewPath("additionalEnabledCapabilities"), c.Capabilities.AdditionalEnabledCapabilities,
					"the marketplace capability requires the OperatorLifecycleManager capability"))
			}
			if c.Platform.BareMetal != nil && !enabledCaps.Has(configv1.ClusterVersionCapabilityBaremetal) {
				allErrs = append(allErrs, field.Invalid(field.NewPath("additionalEnabledCapabilities"), c.Capabilities.AdditionalEnabledCapabilities,
					"platform baremetal requires the baremetal capability"))
			}
		}

		if !enabledCaps.Has(configv1.ClusterVersionCapabilityCloudCredential) {
			// check if platform is cloud
			if c.None == nil && c.BareMetal == nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("capabilities"), c.Capabilities,
					"disabling CloudCredential capability available only for baremetal platforms"))
			}
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

func validateNetworking(n *types.Networking, singleNodeOpenShift bool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if n.NetworkType == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("networkType"), "network provider type required"))
	}

	// NOTE(dulek): We're hardcoding "Kuryr" here as the plan is to remove it from the API very soon. We can remove
	//              this check once some more general validation of the supported NetworkTypes is in place.
	if n.NetworkType == "Kuryr" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("networkType"), n.NetworkType, "networkType Kuryr is not supported on OpenShift later than 4.14"))
	}

	if n.NetworkType == string(operv1.NetworkTypeOpenShiftSDN) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("networkType"), n.NetworkType, "networkType OpenShiftSDN is deprecated, please use OVNKubernetes"))
	}

	if len(n.MachineNetwork) > 0 {
		for i, network := range n.MachineNetwork {
			if err := validate.SubnetCIDR(&network.CIDR.IPNet); err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("machineNetwork").Index(i), network.CIDR.String(), err.Error()))
			}
			for j, subNetwork := range n.MachineNetwork[0:i] {
				if validate.DoCIDRsOverlap(&network.CIDR.IPNet, &subNetwork.CIDR.IPNet) {
					allErrs = append(allErrs, field.Invalid(fldPath.Child("machineNetwork").Index(i), network.CIDR.String(), fmt.Sprintf("machine network must not overlap with machine network %d", j)))
				}
			}
		}
	} else {
		allErrs = append(allErrs, field.Required(fldPath.Child("machineNetwork"), "at least one machine network is required"))
	}

	for i, sn := range n.ServiceNetwork {
		if err := validate.ServiceSubnetCIDR(&sn.IPNet); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("serviceNetwork").Index(i), sn.String(), err.Error()))
		}
		for _, network := range n.MachineNetwork {
			if validate.DoCIDRsOverlap(&sn.IPNet, &network.CIDR.IPNet) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("serviceNetwork").Index(i), sn.String(), "service network must not overlap with any of the machine networks"))
			}
		}
		for j, snn := range n.ServiceNetwork[0:i] {
			if validate.DoCIDRsOverlap(&sn.IPNet, &snn.IPNet) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("serviceNetwork").Index(i), sn.String(), fmt.Sprintf("service network must not overlap with service network %d", j)))
			}
		}
	}
	if len(n.ServiceNetwork) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("serviceNetwork"), "a service network is required"))
	}

	for i, cn := range n.ClusterNetwork {
		allErrs = append(allErrs, validateClusterNetwork(n, &cn, i, fldPath.Child("clusterNetwork").Index(i))...)
	}
	if len(n.ClusterNetwork) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("clusterNetwork"), "cluster network required"))
	}
	return allErrs
}

func validateNetworkingForPlatform(n *types.Networking, platform *types.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	switch {
	case platform.Libvirt != nil:
		errMsg := "overlaps with default Docker Bridge subnet"
		for idx, mn := range n.MachineNetwork {
			if validate.DoCIDRsOverlap(&mn.CIDR.IPNet, validate.DockerBridgeCIDR) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("machineNewtork").Index(idx), mn.CIDR.String(), errMsg))
			}
		}
		for idx, sn := range n.ServiceNetwork {
			if validate.DoCIDRsOverlap(&sn.IPNet, validate.DockerBridgeCIDR) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("serviceNetwork").Index(idx), sn.String(), errMsg))
			}
		}
		for idx, cn := range n.ClusterNetwork {
			if validate.DoCIDRsOverlap(&cn.CIDR.IPNet, validate.DockerBridgeCIDR) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterNetwork").Index(idx), cn.CIDR.String(), errMsg))
			}
		}
	default:
		warningMsgFmt := "%s: %s overlaps with default Docker Bridge subnet"
		for idx, mn := range n.MachineNetwork {
			if validate.DoCIDRsOverlap(&mn.CIDR.IPNet, validate.DockerBridgeCIDR) {
				logrus.Warnf(warningMsgFmt, fldPath.Child("machineNetwork").Index(idx), mn.CIDR.String())
			}
		}
		for idx, sn := range n.ServiceNetwork {
			if validate.DoCIDRsOverlap(&sn.IPNet, validate.DockerBridgeCIDR) {
				logrus.Warnf(warningMsgFmt, fldPath.Child("serviceNetwork").Index(idx), sn.String())
			}
		}
		for idx, cn := range n.ClusterNetwork {
			if validate.DoCIDRsOverlap(&cn.CIDR.IPNet, validate.DockerBridgeCIDR) {
				logrus.Warnf(warningMsgFmt, fldPath.Child("clusterNetwork").Index(idx), cn.CIDR.String())
			}
		}
	}
	return allErrs
}

func validateClusterNetwork(n *types.Networking, cn *types.ClusterNetworkEntry, idx int, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if err := validate.SubnetCIDR(&cn.CIDR.IPNet); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("cidr"), cn.CIDR.IPNet.String(), err.Error()))
	}
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
	if cn.HostPrefix < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("hostPrefix"), cn.HostPrefix, "hostPrefix must be positive"))
	}
	// ignore hostPrefix if the plugin does not use it and has it unset
	if pluginsUsingHostPrefix.Has(n.NetworkType) || (cn.HostPrefix != 0) {
		if ones, bits := cn.CIDR.Mask.Size(); cn.HostPrefix < int32(ones) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("hostPrefix"), cn.HostPrefix, "cluster network host subnetwork prefix must not be larger size than CIDR "+cn.CIDR.String()))
		} else if bits == 128 && cn.HostPrefix != 64 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("hostPrefix"), cn.HostPrefix, "cluster network host subnetwork prefix must be 64 for IPv6 networks"))
		}
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

func validateComputeEdge(platform *types.Platform, pName string, fldPath *field.Path, pfld *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if platform.Name() != aws.Name {
		allErrs = append(allErrs, field.NotSupported(pfld.Child("name"), pName, []string{types.MachinePoolComputeRoleName}))
	}

	return allErrs
}

func validateCompute(platform *types.Platform, control *types.MachinePool, pools []types.MachinePool, fldPath *field.Path) field.ErrorList {
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
		if control != nil && control.Architecture != p.Architecture {
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
func validateVIPsForPlatform(network *types.Networking, platform *types.Platform, fldPath *field.Path) field.ErrorList {
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

		allErrs = append(allErrs, validateAPIAndIngressVIPs(virtualIPs, newVIPsFields, false, false, lbType, network, fldPath.Child(vsphere.Name))...)
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

					if apiVIPNet.Equal(ingressVIPNet) {
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
	if platform.AlibabaCloud != nil {
		validate(alibabacloud.Name, platform.AlibabaCloud, func(f *field.Path) field.ErrorList {
			return alibabacloudvalidation.ValidatePlatform(platform.AlibabaCloud, network, f)
		})
	}
	if platform.AWS != nil {
		validate(aws.Name, platform.AWS, func(f *field.Path) field.ErrorList {
			return awsvalidation.ValidatePlatform(platform.AWS, c.CredentialsMode, f)
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
	if platform.Libvirt != nil {
		validate(libvirt.Name, platform.Libvirt, func(f *field.Path) field.ErrorList { return libvirtvalidation.ValidatePlatform(platform.Libvirt, f) })
	}
	if platform.OpenStack != nil {
		validate(openstack.Name, platform.OpenStack, func(f *field.Path) field.ErrorList {
			return openstackvalidation.ValidatePlatform(platform.OpenStack, network, f, c)
		})
	}
	if platform.PowerVS != nil {
		validate(powervs.Name, platform.PowerVS, func(f *field.Path) field.ErrorList { return powervsvalidation.ValidatePlatform(platform.PowerVS, f) })
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
		alibabacloud.Name: {types.ManualCredentialsMode},
		aws.Name:          {types.MintCredentialsMode, types.PassthroughCredentialsMode, types.ManualCredentialsMode},
		azure.Name:        allowedAzureModes,
		gcp.Name:          {types.MintCredentialsMode, types.PassthroughCredentialsMode, types.ManualCredentialsMode},
		ibmcloud.Name:     {types.ManualCredentialsMode},
		powervs.Name:      {types.ManualCredentialsMode},
		nutanix.Name:      {types.ManualCredentialsMode},
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
		logrus.Warnf("%v", err)
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

	if _, ok := configv1.FeatureSets[c.FeatureSet]; !ok {
		sortedFeatureSets := func() []string {
			v := []string{}
			for n := range configv1.FeatureSets {
				v = append(v, string(n))
			}
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
