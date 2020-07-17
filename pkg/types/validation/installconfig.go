package validation

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	dockerref "github.com/containers/image/docker/reference"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	awsvalidation "github.com/openshift/installer/pkg/types/aws/validation"
	"github.com/openshift/installer/pkg/types/azure"
	azurevalidation "github.com/openshift/installer/pkg/types/azure/validation"
	"github.com/openshift/installer/pkg/types/baremetal"
	baremetalvalidation "github.com/openshift/installer/pkg/types/baremetal/validation"
	"github.com/openshift/installer/pkg/types/gcp"
	gcpvalidation "github.com/openshift/installer/pkg/types/gcp/validation"
	"github.com/openshift/installer/pkg/types/libvirt"
	libvirtvalidation "github.com/openshift/installer/pkg/types/libvirt/validation"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackvalidation "github.com/openshift/installer/pkg/types/openstack/validation"
	"github.com/openshift/installer/pkg/types/ovirt"
	ovirtvalidation "github.com/openshift/installer/pkg/types/ovirt/validation"
	"github.com/openshift/installer/pkg/types/vsphere"
	vspherevalidation "github.com/openshift/installer/pkg/types/vsphere/validation"
	"github.com/openshift/installer/pkg/validate"
)

const (
	masterPoolName = "master"
)

// ValidateInstallConfig checks that the specified install config is valid.
func ValidateInstallConfig(c *types.InstallConfig, openStackValidValuesFetcher openstackvalidation.ValidValuesFetcher) field.ErrorList {
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
	if c.SSHKey != "" {
		if err := validate.SSHPublicKey(c.SSHKey); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("sshKey"), c.SSHKey, err.Error()))
		}
	}
	if c.AdditionalTrustBundle != "" {
		if err := validate.CABundle(c.AdditionalTrustBundle); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("additionalTrustBundle"), c.AdditionalTrustBundle, err.Error()))
		}
	}
	nameErr := validate.ClusterName(c.ObjectMeta.Name)
	if c.Platform.GCP != nil || c.Platform.Azure != nil {
		nameErr = validate.ClusterName1035(c.ObjectMeta.Name)
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
	} else {
		allErrs = append(allErrs, field.Required(field.NewPath("networking"), "networking is required"))
	}
	allErrs = append(allErrs, validatePlatform(&c.Platform, field.NewPath("platform"), openStackValidValuesFetcher, c.Networking, c)...)
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
		allErrs = append(allErrs, validateProxy(c.Proxy, field.NewPath("proxy"))...)
	}
	allErrs = append(allErrs, validateImageContentSources(c.ImageContentSources, field.NewPath("imageContentSources"))...)
	if _, ok := validPublishingStrategies[c.Publish]; !ok {
		allErrs = append(allErrs, field.NotSupported(field.NewPath("publish"), c.Publish, validPublishingStrategyValues))
	}
	allErrs = append(allErrs, validateCloudCredentialsMode(c.CredentialsMode, field.NewPath("credentialsMode"), c.Platform.Name())...)

	return allErrs
}

// ipAddressTypeByField is a map of field path to whether they request IPv4 or IPv6.
type ipAddressTypeByField map[string]struct{ IPv4, IPv6 bool }

// ipByField is a map of field path to the net.IPs in sorted order.
type ipByField map[string][]net.IP

// inferIPVersionFromInstallConfig infers the user's desired ip version from the networking config.
// Presence field names match the field path of the struct within the Networking type. This function
// assumes a valid install config.
func inferIPVersionFromInstallConfig(n *types.Networking) (hasIPv4, hasIPv6 bool, presence ipAddressTypeByField, addresses ipByField) {
	if n == nil {
		return
	}
	addresses = make(ipByField)
	for _, network := range n.MachineNetwork {
		addresses["machineNetwork"] = append(addresses["machineNetwork"], network.CIDR.IP)
	}
	for _, network := range n.ServiceNetwork {
		addresses["serviceNetwork"] = append(addresses["serviceNetwork"], network.IP)
	}
	for _, network := range n.ClusterNetwork {
		addresses["clusterNetwork"] = append(addresses["clusterNetwork"], network.CIDR.IP)
	}
	presence = make(ipAddressTypeByField)
	for k, ips := range addresses {
		for _, ip := range ips {
			has := presence[k]
			if ip.To4() != nil {
				has.IPv4 = true
				if k == "serviceNetwork" {
					hasIPv4 = true
				}
			} else {
				has.IPv6 = true
				if k == "serviceNetwork" {
					hasIPv6 = true
				}
			}
			presence[k] = has
		}
	}
	return
}

func ipSliceToStrings(ips []net.IP) []string {
	var s []string
	for _, ip := range ips {
		s = append(s, ip.String())
	}
	return s
}

func ipnetworksToStrings(networks []ipnet.IPNet) []string {
	var diag []string
	for _, sn := range networks {
		diag = append(diag, sn.String())
	}
	sort.Strings(diag)
	return diag
}

// validateNetworkingIPVersion checks parameters for consistency when the user
// requests single-stack IPv6 or dual-stack modes.
func validateNetworkingIPVersion(n *types.Networking, p *types.Platform) field.ErrorList {
	var allErrs field.ErrorList

	hasIPv4, hasIPv6, presence, addresses := inferIPVersionFromInstallConfig(n)

	switch {
	case hasIPv4 && hasIPv6:
		if n.NetworkType == "OpenShiftSDN" {
			allErrs = append(allErrs, field.Invalid(field.NewPath("networking", "networkType"), n.NetworkType, "dual-stack IPv4/IPv6 is not supported for this networking plugin"))
		}

		if len(n.ServiceNetwork) != 2 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("networking", "serviceNetwork"), strings.Join(ipnetworksToStrings(n.ServiceNetwork), ", "), "when installing dual-stack IPv4/IPv6 you must provide two service networks, one for each IP address type"))
		}

		switch {
		case p.Azure != nil:
		case p.BareMetal != nil:
		case p.None != nil:
		default:
			allErrs = append(allErrs, field.Invalid(field.NewPath("networking"), "DualStack", "dual-stack IPv4/IPv6 is not supported for this platform, specify only one type of address"))
		}
		for k, v := range presence {
			switch {
			case k == "machineNetwork" && p.AWS != nil:
				// AWS can default an ipv6 subnet
			case v.IPv4 && !v.IPv6:
				allErrs = append(allErrs, field.Invalid(field.NewPath("networking", k), strings.Join(ipSliceToStrings(addresses[k]), ", "), "dual-stack IPv4/IPv6 requires an IPv6 address in this list"))
			case !v.IPv4 && v.IPv6:
				allErrs = append(allErrs, field.Invalid(field.NewPath("networking", k), strings.Join(ipSliceToStrings(addresses[k]), ", "), "dual-stack IPv4/IPv6 requires an IPv4 address in this list"))
			}
		}

	case hasIPv6:
		if n.NetworkType == "OpenShiftSDN" {
			allErrs = append(allErrs, field.Invalid(field.NewPath("networking", "networkType"), n.NetworkType, "IPv6 is not supported for this networking plugin"))
		}

		switch {
		case p.BareMetal != nil:
		case p.None != nil:
		case p.Azure != nil && os.Getenv("OPENSHIFT_INSTALL_AZURE_EMULATE_SINGLESTACK_IPV6") == "true":
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
	if n.NetworkType == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("networkType"), "network provider type required"))
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
		if err := validate.SubnetCIDR(&sn.IPNet); err != nil {
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
	if ones, _ := cn.CIDR.Mask.Size(); cn.HostPrefix < int32(ones) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("hostPrefix"), cn.HostPrefix, "cluster network host subnetwork prefix must not be larger size than CIDR "+cn.CIDR.String()))
	}
	return allErrs
}

func validateControlPlane(platform *types.Platform, pool *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if pool.Name != masterPoolName {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("name"), pool.Name, []string{masterPoolName}))
	}
	if pool.Replicas != nil && *pool.Replicas == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("replicas"), pool.Replicas, "number of control plane replicas must be positive"))
	}
	allErrs = append(allErrs, ValidateMachinePool(platform, pool, fldPath)...)
	return allErrs
}

func validateCompute(platform *types.Platform, control *types.MachinePool, pools []types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	poolNames := map[string]bool{}
	for i, p := range pools {
		poolFldPath := fldPath.Index(i)
		if p.Name != "worker" {
			allErrs = append(allErrs, field.NotSupported(poolFldPath.Child("name"), p.Name, []string{"worker"}))
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

func validatePlatform(platform *types.Platform, fldPath *field.Path, openStackValidValuesFetcher openstackvalidation.ValidValuesFetcher, network *types.Networking, c *types.InstallConfig) field.ErrorList {
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
		validate(aws.Name, platform.AWS, func(f *field.Path) field.ErrorList { return awsvalidation.ValidatePlatform(platform.AWS, f) })
	}
	if platform.Azure != nil {
		validate(azure.Name, platform.Azure, func(f *field.Path) field.ErrorList {
			return azurevalidation.ValidatePlatform(platform.Azure, c.Publish, f)
		})
	}
	if platform.GCP != nil {
		validate(gcp.Name, platform.GCP, func(f *field.Path) field.ErrorList { return gcpvalidation.ValidatePlatform(platform.GCP, f) })
	}
	if platform.Libvirt != nil {
		validate(libvirt.Name, platform.Libvirt, func(f *field.Path) field.ErrorList { return libvirtvalidation.ValidatePlatform(platform.Libvirt, f) })
	}
	if platform.OpenStack != nil {
		validate(openstack.Name, platform.OpenStack, func(f *field.Path) field.ErrorList {
			return openstackvalidation.ValidatePlatform(platform.OpenStack, network, f, openStackValidValuesFetcher, c)
		})
	}
	if platform.VSphere != nil {
		validate(vsphere.Name, platform.VSphere, func(f *field.Path) field.ErrorList { return vspherevalidation.ValidatePlatform(platform.VSphere, f) })
	}
	if platform.BareMetal != nil {
		validate(baremetal.Name, platform.BareMetal, func(f *field.Path) field.ErrorList {
			return baremetalvalidation.ValidatePlatform(platform.BareMetal, network, f, c)
		})
	}
	if platform.Ovirt != nil {
		validate(ovirt.Name, platform.Ovirt, func(f *field.Path) field.ErrorList {
			return ovirtvalidation.ValidatePlatform(platform.Ovirt, f)
		})
	}
	return allErrs
}

func validateProxy(p *types.Proxy, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.HTTPProxy == "" && p.HTTPSProxy == "" {
		allErrs = append(allErrs, field.Required(fldPath, "must include httpProxy or httpsProxy"))
	}
	if p.HTTPProxy != "" {
		if err := validate.URI(p.HTTPProxy); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("HTTPProxy"), p.HTTPProxy, err.Error()))
		}
	}
	if p.HTTPSProxy != "" {
		if err := validate.URI(p.HTTPSProxy); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("HTTPSProxy"), p.HTTPSProxy, err.Error()))
		}
	}
	if p.NoProxy != "" {
		for _, v := range strings.Split(p.NoProxy, ",") {
			v = strings.TrimSpace(v)
			errDomain := validate.NoProxyDomainName(v)
			_, _, errCIDR := net.ParseCIDR(v)
			if errDomain != nil && errCIDR != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("NoProxy"), v, "must be a CIDR or domain, without wildcard characters"))
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

func validateNamedRepository(r string) error {
	ref, err := dockerref.ParseNamed(r)
	if err != nil {
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

func validateCloudCredentialsMode(mode types.CredentialsMode, fldPath *field.Path, platform string) field.ErrorList {
	if mode == "" {
		return nil
	}
	allErrs := field.ErrorList{}
	// validPlatformCredentialsModes is a map from the platform name to a slice of credentials modes that are valid
	// for the platform. If a platform name is not in the map, then the credentials mode cannot be set for that platform.
	validPlatformCredentialsModes := map[string][]types.CredentialsMode{
		aws.Name:   {types.MintCredentialsMode, types.PassthroughCredentialsMode, types.ManualCredentialsMode},
		azure.Name: {types.MintCredentialsMode, types.PassthroughCredentialsMode},
		gcp.Name:   {types.MintCredentialsMode, types.PassthroughCredentialsMode},
	}
	if validModes, ok := validPlatformCredentialsModes[platform]; ok {
		validModesSet := sets.NewString()
		for _, m := range validModes {
			validModesSet.Insert(string(m))
		}
		if !validModesSet.Has(string(mode)) {
			allErrs = append(allErrs, field.NotSupported(fldPath, mode, validModesSet.List()))
		}
	} else {
		allErrs = append(allErrs, field.Invalid(fldPath, mode, fmt.Sprintf("cannot be set when using the %q platform", platform)))
	}
	return allErrs
}
