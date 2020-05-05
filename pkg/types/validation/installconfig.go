package validation

import (
        "archive/tar"
	"compress/gzip"
	"fmt"
        "gopkg.in/yaml.v2"
        "io"
        "io/ioutil"
	"log"
	"net"
	"os"
	"sort"
        "strconv"
	"strings"

	dockerref "github.com/containers/image/docker/reference"
	"github.com/pkg/errors"
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
	"github.com/openshift/installer/pkg/types/vsphere"
	vspherevalidation "github.com/openshift/installer/pkg/types/vsphere/validation"
	"github.com/openshift/installer/pkg/validate"
)

const (
	masterPoolName = "master"
)

type AciContainersConfig struct {
        Data ConfigData `yaml:data,omitempty`
}

type ConfigData struct {
        HostConfig string `yaml:"host-agent-config"`
}

type HostConfigMap struct {
        ServiceVLAN int    `yaml:"service-vlan"`
        InfraVLAN   int    `yaml:"aci-infra-vlan"`
        KubeApiVLAN int    `yaml:"kubeapi-vlan"`
        PodSubnet   string `yaml:"pod-subnet"`
        NodeSubnet  string `yaml:"node-subnet"`
}

type ClusterConfig03 struct {
	ApiVersion string     `yaml:"apiVersion"`
        Kind       string     `yaml:"kind"`
        Metadata   MetaEntry  `yaml:"metadata,omitempty"`
	Spec       SpecEntry  `yaml:"spec,omitempty"`
}

type MetaEntry struct {
	Name	string `yaml:"name"`
}

type SpecEntry struct {
	Multus		bool				`yaml:"disableMultiNetwork"`
        ClusterNetwork	[]ClusterEntry 			`yaml:"clusterNetwork,omitempty"`  
        DefaultNetwork  DefaultNetEntry			`yaml:"defaultNetwork,omitempty"`
        NetworkType	string				`yaml:"networkType,omitempty"`
        ServiceNetwork	[]string			`yaml:"serviceNetwork,omitempty"`
}

type ClusterEntry struct {
        CIDR		string	`yaml:"cidr"`
	HostPrefix	int32	`yaml:"hostPrefix"`
}

type DefaultNetEntry struct {
	Type	string	`yaml:"type"`
}


// ClusterDomain returns the cluster domain for a cluster with the specified
// base domain and cluster name.
func ClusterDomain(baseDomain, clusterName string) string {
	return fmt.Sprintf("%s.%s", clusterName, baseDomain)
}

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
	if c.Platform.GCP != nil {
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
		clusterDomain := ClusterDomain(c.BaseDomain, c.ObjectMeta.Name)
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

	tarField := field.NewPath("ProvisionTar")
        r, err := os.Open(c.Platform.OpenStack.AciNetExt.ProvisionTar)
        if err != nil {
		allErrs = append(allErrs, field.Invalid(tarField, c.Platform.OpenStack.AciNetExt.ProvisionTar, err.Error()))
    	} else {
                config, err := ExtractTarGz(r)
                if err != nil {
                        allErrs = append(allErrs, field.Invalid(tarField.Child("Unmarshal"),
                                c.Platform.OpenStack.AciNetExt.ProvisionTar, err.Error()))
                } else {
			c.Platform.OpenStack.AciNetExt.KubeApiVLAN = strconv.Itoa(config.KubeApiVLAN)
			c.Platform.OpenStack.AciNetExt.InfraVLAN = strconv.Itoa(config.InfraVLAN)
			c.Platform.OpenStack.AciNetExt.ServiceVLAN = strconv.Itoa(config.ServiceVLAN)

                        // Validate against values from install config
			machineCIDR := &c.Networking.MachineNetwork[0].CIDR
			clusterNetworkCIDR := &c.Networking.ClusterNetwork[0].CIDR
			nodeDiff := DiffSubnets(config.NodeSubnet, machineCIDR)
                        if nodeDiff != nil {
				option := UserPrompt(nodeDiff.String(), machineCIDR, "node_subnet", "machineNetworkCIDR")
				if (option == true) {
					cidrValue, _ := ipnet.ParseCIDR(nodeDiff.String())
					c.Networking.MachineNetwork[0].CIDR = *cidrValue
					log.Print("Setting machineCIDR to " + nodeDiff.String())
				} else {
                                	allErrs = append(allErrs, field.Invalid(field.NewPath("machineNetworkCIDR"),
                                        	c.Networking.DeprecatedMachineCIDR.String(), "node_subnet in acc-provision input(" + nodeDiff.String() + ") has to be the same as machineNetwork CIDR in install-config.yaml(" + machineCIDR.String() + ")"))
				}
                        }
			clusterDiff := DiffSubnets(config.PodSubnet, clusterNetworkCIDR)
                        if clusterDiff != nil {
				option := UserPrompt(clusterDiff.String(), clusterNetworkCIDR, "pod_subnet", "clusterNetworkCIDR")
				if (option == true) {
					parsedCIDR, _ := ipnet.ParseCIDR(clusterDiff.String())
					c.Networking.ClusterNetwork[0].CIDR = *parsedCIDR
					log.Print("Setting clusterNetwork CIDR to " + clusterDiff.String())
				} else {
                                	allErrs = append(allErrs, field.Invalid(field.NewPath("clusterNetworkCIDR"),
                                        	clusterNetworkCIDR.String(), "pod_subnet in acc-provision input(" + clusterDiff.String() + ") has to be the same as clusterNetwork:cidr in install-config.yaml(" + clusterNetworkCIDR.String() + ")"))
				}
                        }
		}
	}

	allErrs = append(allErrs, validatePlatform(&c.Platform, field.NewPath("platform"), openStackValidValuesFetcher, c.Networking, c)...)
	if c.ControlPlane != nil {
		allErrs = append(allErrs, validateControlPlane(&c.Platform, c.ControlPlane, field.NewPath("controlPlane"))...)
	} else {
		allErrs = append(allErrs, field.Required(field.NewPath("controlPlane"), "controlPlane is required"))
	}
	allErrs = append(allErrs, validateCompute(&c.Platform, c.Compute, field.NewPath("compute"))...)
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

	return allErrs
}

func DiffSubnets(sub1 string, sub2 *ipnet.IPNet) *net.IPNet {
        // Returns first subnet if the subnets are different
        _, net1, _ := net.ParseCIDR(sub1)
        if net1.String() != sub2.String() {
                return net1
	}
        return nil
}

func UserPrompt(sub1 string, sub2 *ipnet.IPNet, item1 string, item2 string) bool {
	var option string
	log.Print("There's a discrepancy between " + item1 + "(" + sub1 + ") in acc-provision input and " + item2 + "(" + sub2.String() + ") in install-config.yaml")
	log.Print("Enter Y to use acc-provision value, or N to exit installer and fix acc-provision tar")
	fmt.Scanln(&option)
	var op bool
	if (option == "y" || option == "Y") {
		op = true
	}
	return op
}

func ExtractTarGz(gzipStream io.Reader) (HostConfigMap, error) {
	config := HostConfigMap{}
        uncompressedStream, err := gzip.NewReader(gzipStream)
        if err != nil {
		return config, err
        }

        tarReader := tar.NewReader(uncompressedStream)

        for true {
                header, err := tarReader.Next()

                if err == io.EOF {
                        break
                }

                if err != nil {
			return config, err
                }

                switch header.Typeflag {
                case tar.TypeReg:
                        temp, err := ioutil.ReadAll(tarReader)
			if err != nil {
				return config, err
			}

			// Unmarshal acc configmap to get acc-provision values
                        if strings.Contains(header.Name, "aci-containers-config") {
                                t := AciContainersConfig{}
                                err = yaml.Unmarshal(temp, &t)
                                if err != nil {
					return config, err
                                }
                                err = yaml.Unmarshal([]byte(t.Data.HostConfig), &config)
                                if err != nil {
					return config, err
                                }
				err = checkForParsedConfigValues(config)
				if err != nil {
                                        return config, err
                                }
                        }
                default:
			return config, errors.New("Unsupported file type in tar")
		}

        }
        return config, nil
}

func checkForParsedConfigValues(config HostConfigMap) error {
	if config.PodSubnet == "" || config.NodeSubnet == "" ||
		config.KubeApiVLAN == 0 || config.ServiceVLAN == 0 ||
			config.InfraVLAN == 0 {
				return errors.New("One or more values missing from acc-provision tar configmap")
	}
	return nil
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

func validateCompute(platform *types.Platform, pools []types.MachinePool, fldPath *field.Path) field.ErrorList {
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
			return baremetalvalidation.ValidatePlatform(platform.BareMetal, network, f)
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
