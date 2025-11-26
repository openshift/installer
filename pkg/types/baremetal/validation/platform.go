package validation

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/go-playground/validator/v10"
	"github.com/metal3-io/baremetal-operator/pkg/hardwareutils/bmc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/common"
	"github.com/openshift/installer/pkg/validate"
)

type interfaceValidatorFactory func(string) (func(string) error, error)

var interfaceValidator interfaceValidatorFactory = libvirtInterfaceValidator

func validateIPinMachineCIDR(vip string, n *types.Networking) error {
	var networks []string

	for _, network := range n.MachineNetwork {
		if network.CIDR.Contains(net.ParseIP(vip)) {
			return nil
		}
		networks = append(networks, network.CIDR.String())
	}

	return fmt.Errorf("IP expected to be in one of the machine networks: %s", strings.Join(networks, ","))
}

func validateCIDRSize(p *baremetal.Platform) error {
	provisionNetworkCIDR := &p.ProvisioningNetworkCIDR.IPNet
	if p.ProvisioningNetwork == baremetal.ManagedProvisioningNetwork {
		cidrSize, _ := provisionNetworkCIDR.Mask.Size()
		if cidrSize < 64 && provisionNetworkCIDR.IP.To4() == nil && provisionNetworkCIDR.IP.To16() != nil {
			return fmt.Errorf("provisioningNetworkCIDR mask must be greater than or equal to 64 for IPv6 networks")
		}
	}
	return nil
}

func validateIPNotinMachineCIDR(ip string, n *types.Networking) error {
	for _, network := range n.MachineNetwork {
		if network.CIDR.Contains(net.ParseIP(ip)) {
			return fmt.Errorf("the IP must not be in one of the machine networks")
		}
	}
	return nil
}

func validateNoOverlapMachineCIDR(target *net.IPNet, n *types.Networking) error {
	allIPv4 := ipnet.MustParseCIDR("0.0.0.0/0")
	allIPv6 := ipnet.MustParseCIDR("::/0")
	netIsIPv6 := target.IP.To4() == nil

	for _, machineCIDR := range n.MachineNetwork {
		machineCIDRisIPv6 := machineCIDR.CIDR.IP.To4() == nil

		// Only compare if both are the same IP version
		if netIsIPv6 == machineCIDRisIPv6 {
			var err error
			if netIsIPv6 {
				err = cidr.VerifyNoOverlap(
					[]*net.IPNet{
						target,
						&machineCIDR.CIDR.IPNet,
					},
					&allIPv6.IPNet,
				)
			} else {
				err = cidr.VerifyNoOverlap(
					[]*net.IPNet{
						target,
						&machineCIDR.CIDR.IPNet,
					},
					&allIPv4.IPNet,
				)
			}

			if err != nil {
				return errors.Wrap(err, "cannot overlap with machine network")
			}
		}
	}

	return nil
}

func validateOSImageURI(uri string) error {
	// Check for valid URI and sha256 checksum part of the URL
	parsedURL, err := url.ParseRequestURI(uri)
	if err != nil {
		return fmt.Errorf("the URI provided: %s is invalid", uri)
	}
	if parsedURL.Scheme == "http" || parsedURL.Scheme == "https" {
		var sha256Checksum string
		if sha256Checksums, ok := parsedURL.Query()["sha256"]; ok {
			sha256Checksum = sha256Checksums[0]
		}
		if sha256Checksum == "" {
			return fmt.Errorf("the sha256 parameter in the %s URI is missing", uri)
		}
		if len(sha256Checksum) != 64 {
			return fmt.Errorf("the sha256 parameter in the %s URI is invalid", uri)
		}
	} else {
		return fmt.Errorf("the URI provided: %s must begin with http/https", uri)
	}
	return nil
}

// ValidateNTPServers checks list of NTP servers strings are valid IPs or domain names.
func ValidateNTPServers(servers []string, fldPath *field.Path) (allErrs field.ErrorList) {
	for i, server := range servers {
		if ipErr := validate.IP(server); ipErr != nil {
			if domainErr := validate.DomainName(server, true); domainErr != nil {
				allErrs = append(allErrs, field.Invalid(fldPath, servers[i], "NTP server is not a valid IP or domain name"))
			}
		}
	}
	return
}

// validateDHCPRange ensures the provided range is valid, and that provisioning service IP's do not overlap.
func validateDHCPRange(p *baremetal.Platform, fldPath *field.Path) (allErrs field.ErrorList) {
	dhcpRange := strings.Split(p.ProvisioningDHCPRange, ",")

	if len(dhcpRange) != 2 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningDHCPRange"), p.ProvisioningDHCPRange, "provisioning DHCP range should be in format: start_ip,end_ip"))
		return
	}

	for _, ip := range dhcpRange {
		// Ensure IP is valid
		if err := validate.IP(ip); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningDHCPRange"), p.ProvisioningDHCPRange, fmt.Sprintf("%s: %s", ip, err.Error())))
			return
		}

		// Validate IP is in the provisioning network
		if p.ProvisioningNetworkCIDR != nil && !p.ProvisioningNetworkCIDR.Contains(net.ParseIP(ip)) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningDHCPRange"), p.ProvisioningDHCPRange, fmt.Sprintf("%q is not in the provisioning network", ip)))
			return
		}
	}

	// Validate VIP's are not within the DHCP Range
	start := net.ParseIP(dhcpRange[0])
	end := net.ParseIP(dhcpRange[1])

	if start != nil && end != nil {
		// Validate ClusterProvisioningIP is not in DHCP range
		if clusterProvisioningIP := net.ParseIP(p.ClusterProvisioningIP); clusterProvisioningIP != nil && bytes.Compare(clusterProvisioningIP, start) >= 0 && bytes.Compare(clusterProvisioningIP, end) <= 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterProvisioningIP"), p.ClusterProvisioningIP, fmt.Sprintf("%q overlaps with the allocated DHCP range", p.ClusterProvisioningIP)))
		}

		// Validate BootstrapProvisioningIP is not in DHCP range
		if bootstrapProvisioningIP := net.ParseIP(p.BootstrapProvisioningIP); bootstrapProvisioningIP != nil && bytes.Compare(bootstrapProvisioningIP, start) >= 0 && bytes.Compare(bootstrapProvisioningIP, end) <= 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapProvisioningIP"), p.BootstrapProvisioningIP, fmt.Sprintf("%q overlaps with the allocated DHCP range", p.BootstrapProvisioningIP)))
		}
	}

	return
}

// validateHostsBase validates the hosts based on a filtering function
func validateHostsBase(hosts []*baremetal.Host, fldPath *field.Path, filter validator.FilterFunc) field.ErrorList {
	return common.ValidateUniqueAndRequiredFields(hosts, fldPath, filter)
}

// filterHostsBMC is a function to control whether to filter BMC details of Hosts
func filterHostsBMC(ns []byte) bool {
	return bytes.Contains(ns, []byte(".BMC"))
}

// validateHostsWithoutBMC utilizes the filter function to disable BMC checking while validating hosts
func validateHostsWithoutBMC(hosts []*baremetal.Host, fldPath *field.Path) field.ErrorList {
	return validateHostsBase(hosts, fldPath, filterHostsBMC)
}

// validateHostsBMCOnly utilizes the filter function to only perform validation on BMC part of the hosts
func validateHostsBMCOnly(hosts []*baremetal.Host, fldPath *field.Path) field.ErrorList {
	return validateHostsBase(hosts, fldPath, func(ns []byte) bool {
		return !filterHostsBMC(ns)
	})
}

func validateOSImages(p *baremetal.Platform, fldPath *field.Path) field.ErrorList {
	var errs field.ErrorList

	fields := map[string]string{
		"bootstrapOSImage": p.BootstrapOSImage,
		"clusterOSImage":   p.ClusterOSImage,
	}

	for fieldName, url := range fields {
		if url == "" {
			continue
		}
		if err := validateOSImageURI(url); err != nil {
			errs = append(errs,
				field.Invalid(fldPath.Child(fieldName), url, err.Error()))
		} else if res, err := http.Head(url); err != nil || res.StatusCode != http.StatusOK /* #nosec G107 */ {
			errs = append(errs,
				field.NotFound(fldPath.Child(fieldName), url))
		}
	}
	return errs
}

func validateHostsName(hosts []*baremetal.Host, fldPath *field.Path) (errors field.ErrorList) {
	for idx, host := range hosts {
		validationMessages := validation.IsDNS1123Subdomain(host.Name)
		if len(validationMessages) != 0 {
			msg := strings.Join(validationMessages, "\n")
			errors = append(errors, field.Invalid(fldPath.Index(idx).Child("name"), host.Name, msg))
		}
	}
	return
}

// ensure that the number of hosts is enough to cover the ControlPlane
// and Compute replicas. Hosts without role will be considered eligible
// for the ControlPlane or Compute requirements.
func validateHostsCount(hosts []*baremetal.Host, installConfig *types.InstallConfig) error {

	numRequiredMasters := int64(0)
	if installConfig.ControlPlane != nil && installConfig.ControlPlane.Replicas != nil {
		numRequiredMasters += *installConfig.ControlPlane.Replicas
	}

	numRequiredWorkers := int64(0)
	for _, worker := range installConfig.Compute {
		if worker.Replicas != nil {
			numRequiredWorkers += *worker.Replicas
		}
	}

	numRequiredArbiters := int64(0)
	if installConfig.Arbiter != nil && installConfig.Arbiter.Replicas != nil {
		numRequiredArbiters += *installConfig.Arbiter.Replicas
	}

	numMasters := int64(0)
	numArbiters := int64(0)
	numWorkers := int64(0)

	for _, h := range hosts {
		if h.IsMaster() {
			numMasters++
		} else if h.IsArbiter() {
			numArbiters++
		} else if h.IsWorker() {
			numWorkers++
		} else {
			logrus.Warn(fmt.Sprintf("Host %s hasn't any role configured", h.Name))
			if numMasters < numRequiredMasters {
				numMasters++
			} else if numArbiters < numRequiredArbiters {
				numArbiters++
			} else if numWorkers < numRequiredWorkers {
				numWorkers++
			}
		}
	}

	if numMasters < numRequiredMasters {
		return fmt.Errorf("not enough hosts found (%v) to support all the configured ControlPlane replicas (%v)", numMasters, numRequiredMasters)
	}

	if numArbiters < numRequiredArbiters {
		return fmt.Errorf("not enough hosts found (%v) to support all the configured Arbiter replicas (%v)", numArbiters, numRequiredArbiters)
	}

	if numWorkers < numRequiredWorkers {
		return fmt.Errorf("not enough hosts found (%v) to support all the configured Compute replicas (%v)", numWorkers, numRequiredWorkers)
	}

	return nil
}

func validateMTUIsInteger(nmstateYAML []byte, fldPath *field.Path) field.ErrorList {
	var config map[string]interface{}
	if err := yaml.Unmarshal(nmstateYAML, &config); err != nil {
		return field.ErrorList{
			field.Invalid(fldPath, string(nmstateYAML), fmt.Sprintf("failed to unmarshal NMState config: %v", err)),
		}
	}

	interfaces, ok := config["interfaces"].([]interface{})
	if !ok {
		return nil // no interfaces, nothing to check
	}

	var allErrs field.ErrorList
	for idx, iface := range interfaces {
		ifaceMap, ok := iface.(map[string]interface{})
		if !ok {
			continue
		}

		if mtu, exists := ifaceMap["mtu"]; exists {
			switch v := mtu.(type) {
			case int, int64, float64: // yaml unmarshals numbers as float64
				// ok
			case string:
				// check if string is actually numeric
				if _, err := strconv.Atoi(v); err != nil {
					allErrs = append(allErrs,
						field.Invalid(
							fldPath.Child("interfaces").Index(idx).Child("mtu"),
							v,
							"mtu must be an integer",
						),
					)
				} else {
					allErrs = append(allErrs,
						field.Invalid(
							fldPath.Child("interfaces").Index(idx).Child("mtu"),
							v,
							"mtu must be an integer (not quoted string)",
						),
					)
				}
			default:
				allErrs = append(allErrs,
					field.Invalid(
						fldPath.Child("interfaces").Index(idx).Child("mtu"),
						v,
						fmt.Sprintf("mtu must be an integer, got %T", v),
					),
				)
			}
		}
	}
	return allErrs
}

// ensure that the NetworkConfig field contains a valid Yaml string
func validateNetworkConfig(hosts []*baremetal.Host, fldPath *field.Path) (errors field.ErrorList) {
	for idx, host := range hosts {
		if host.NetworkConfig != nil {
			networkConfig := make(map[string]interface{})
			err := yaml.Unmarshal(host.NetworkConfig.Raw, &networkConfig)
			if err != nil {
				errors = append(errors, field.Invalid(fldPath.Index(idx).Child("networkConfig"), host.NetworkConfig, fmt.Sprintf("Not a valid yaml: %s", err.Error())))
			}

			errors = append(errors,
				validateMTUIsInteger(host.NetworkConfig.Raw, fldPath.Index(idx).Child("networkConfig"))...,
			)
		}
	}
	return
}

// ensure that the bootMode field contains a valid value
func validateBootMode(hosts []*baremetal.Host, fldPath *field.Path) (errors field.ErrorList) {
	for idx, host := range hosts {
		switch host.BootMode {
		case "", baremetal.UEFI, baremetal.Legacy:
		case baremetal.UEFISecureBoot:
			accessDetails, err := bmc.NewAccessDetails(host.BMC.Address, host.BMC.DisableCertificateVerification)
			if err == nil && !accessDetails.SupportsSecureBoot() {
				msg := fmt.Sprintf("driver %s does not support UEFI secure boot", accessDetails.Driver())
				errors = append(errors, field.Invalid(fldPath.Index(idx).Child("bootMode"), host.BootMode, msg))
			}
			// if access details cannot be constructed, this should be reported elsewhere
		default:
			valid := []string{string(baremetal.UEFI), string(baremetal.UEFISecureBoot), string(baremetal.Legacy)}
			errors = append(errors, field.NotSupported(fldPath.Index(idx).Child("bootMode"), host.BootMode, valid))
		}
	}
	return
}

// ValidateHostRootDeviceHints checks that a rootDeviceHints field contains no
// invalid values.
func ValidateHostRootDeviceHints(rdh *baremetal.RootDeviceHints, fldPath *field.Path) (errors field.ErrorList) {
	if rdh == nil || rdh.DeviceName == "" {
		return
	}
	devField := fldPath.Child("deviceName")
	subpath := strings.TrimPrefix(rdh.DeviceName, "/dev/")
	if rdh.DeviceName == subpath {
		errors = append(errors, field.Invalid(devField, rdh.DeviceName,
			"Device Name of root device hint must be a /dev/ path"))
	}

	subpath = strings.TrimPrefix(subpath, "disk/by-path/")
	if strings.Contains(subpath, "/") {
		errors = append(errors, field.Invalid(devField, rdh.DeviceName,
			"Device Name of root device hint must be path in /dev/ or /dev/disk/by-path/"))
	}
	return
}

// ensure that none of the rootDeviceHints fields contain invalid values.
func validateRootDeviceHints(hosts []*baremetal.Host, fldPath *field.Path) (errors field.ErrorList) {
	for idx, host := range hosts {
		if host == nil || host.RootDeviceHints == nil {
			continue
		}
		errors = append(errors, ValidateHostRootDeviceHints(host.RootDeviceHints, fldPath.Index(idx).Child("rootDeviceHints"))...)
	}
	return
}

// validateProvisioningNetworkDisabledSupported validates hosts bmc address support provisioning network is disabled
func validateProvisioningNetworkDisabledSupported(hosts []*baremetal.Host, fldPath *field.Path) (errors field.ErrorList) {
	for idx, host := range hosts {
		accessDetails, err := bmc.NewAccessDetails(host.BMC.Address, host.BMC.DisableCertificateVerification)
		if err != nil {
			errors = append(errors, field.Invalid(fldPath.Index(idx).Child("bmc"), host.BMC.Address, err.Error()))
		} else if accessDetails.RequiresProvisioningNetwork() {
			msg := fmt.Sprintf("driver %s requires provisioning network", accessDetails.Driver())
			errors = append(errors, field.Invalid(fldPath.Index(idx).Child("bmc"), host.BMC.Address, msg))
		}
	}

	return
}

// validateHostsBMCForFencing validates BMC addresses are RedFish-compatible for fencing (TNF).
// This is required for Two-Node Fencing configurations where Pacemaker needs to fence nodes
// via RedFish-compatible BMCs.
func validateHostsBMCForFencing(hosts []*baremetal.Host, installConfig *types.InstallConfig, fldPath *field.Path) field.ErrorList {
	errors := field.ErrorList{}

	// Only validate if this is a TNF cluster (2 control plane replicas with fencing enabled)
	if installConfig.ControlPlane == nil || installConfig.ControlPlane.Replicas == nil {
		return errors
	}

	// TNF requires exactly 2 control plane nodes and fencing credentials
	isTNF := *installConfig.ControlPlane.Replicas == 2 &&
		installConfig.ControlPlane.Fencing != nil &&
		len(installConfig.ControlPlane.Fencing.Credentials) > 0

	if !isTNF {
		return errors
	}

	// For TNF clusters, validate that control plane BMC addresses are RedFish-compatible
	for idx, host := range hosts {
		if !host.IsMaster() {
			continue // Only validate control plane hosts
		}

		// Use the shared RedFish BMC validation function from types/common package
		if validationErrs := common.ValidateRedfishBMCAddress(host.BMC.Address, fldPath.Index(idx).Child("bmc").Child("address")); len(validationErrs) > 0 {
			errors = append(errors, validationErrs...)
		}
	}

	return errors
}

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *baremetal.Platform, agentBasedInstallation bool, n *types.Networking, fldPath *field.Path, c *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	provisioningNetwork := sets.NewString(string(baremetal.ManagedProvisioningNetwork),
		string(baremetal.UnmanagedProvisioningNetwork),
		string(baremetal.DisabledProvisioningNetwork))

	if !provisioningNetwork.Has(string(p.ProvisioningNetwork)) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("provisioningNetwork"), p.ProvisioningNetwork, provisioningNetwork.List()))
	}

	if p.BootstrapProvisioningIP != "" {
		if err := validate.IP(p.BootstrapProvisioningIP); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapProvisioningIP"), p.BootstrapProvisioningIP, err.Error()))
		}
	}

	if p.ClusterProvisioningIP != "" {
		if err := validate.IP(p.ClusterProvisioningIP); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterProvisioningIP"), p.ClusterProvisioningIP, err.Error()))
		}
	}

	enabledCaps := c.GetEnabledCapabilities()
	if !agentBasedInstallation && enabledCaps.Has(configv1.ClusterVersionCapabilityMachineAPI) && p.Hosts == nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("hosts"), p.Hosts, "bare metal hosts are missing"))
	}

	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}

	if !agentBasedInstallation && enabledCaps.Has(configv1.ClusterVersionCapabilityMachineAPI) {
		if err := ValidateHosts(p, fldPath, c); err != nil {
			allErrs = append(allErrs, err...)
		}
	}

	if c.BareMetal.LoadBalancer != nil {
		if !validateLoadBalancer(c.BareMetal.LoadBalancer.Type) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("loadBalancer", "type"), c.BareMetal.LoadBalancer.Type, "invalid load balancer type"))
		}
	}

	return allErrs
}

// ValidateHosts returns an error if the Hosts are not valid.
func ValidateHosts(p *baremetal.Platform, fldPath *field.Path, c *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	fldPath = fldPath.Child("hosts")
	if err := validateHostsCount(p.Hosts, c); err != nil {
		allErrs = append(allErrs, field.Required(fldPath, err.Error()))
	}
	allErrs = append(allErrs, validateHostsWithoutBMC(p.Hosts, fldPath)...)
	allErrs = append(allErrs, validateBootMode(p.Hosts, fldPath)...)
	allErrs = append(allErrs, validateRootDeviceHints(p.Hosts, fldPath)...)
	allErrs = append(allErrs, validateNetworkConfig(p.Hosts, fldPath)...)

	allErrs = append(allErrs, validateHostsName(p.Hosts, fldPath)...)

	// Validate BMC addresses for TNF (Two-Node Fencing) clusters
	allErrs = append(allErrs, validateHostsBMCForFencing(p.Hosts, c, fldPath)...)

	return allErrs
}

// validateLoadBalancer returns an error if the load balancer is not valid.
func validateLoadBalancer(lbType configv1.PlatformLoadBalancerType) bool {
	switch lbType {
	case configv1.LoadBalancerTypeOpenShiftManagedDefault, configv1.LoadBalancerTypeUserManaged:
		return true
	default:
		return false
	}
}

// ValidateProvisioning checks that provisioning network requirements specified is valid.
func ValidateProvisioning(p *baremetal.Platform, n *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validateProvisioningBootstrapAndImages(p, n, fldPath)...)

	allErrs = append(allErrs, ValidateProvisioningNetworking(p, n, fldPath)...)

	allErrs = append(allErrs, validateProvisioningBootstrapNetworking(p, fldPath)...)

	return allErrs
}

// validateProvisioningBootstrapAndImagechecks that provisioning settings and images required for bootstrap are valid.
func validateProvisioningBootstrapAndImages(p *baremetal.Platform, n *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	switch p.ProvisioningNetwork {
	case baremetal.DisabledProvisioningNetwork:
		// If set, ensure bootstrapProvisioningIP is in one of the machine networks
		if p.BootstrapProvisioningIP != "" {
			if err := validateIPinMachineCIDR(p.BootstrapProvisioningIP, n); err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapProvisioningIP"), p.BootstrapProvisioningIP, fmt.Sprintf("provisioning network is disabled, %s", err.Error())))
			}
		}
	default:
		// Ensure bootstrapProvisioningIP is in the provisioningNetworkCIDR
		if !p.ProvisioningNetworkCIDR.Contains(net.ParseIP(p.BootstrapProvisioningIP)) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapProvisioningIP"), p.BootstrapProvisioningIP, fmt.Sprintf("%q is not in the provisioning network", p.BootstrapProvisioningIP)))
		}

		if err := validateIPNotinMachineCIDR(p.BootstrapProvisioningIP, n); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapProvisioningIP"), p.BootstrapProvisioningIP, err.Error()))
		}

		if err := validate.MAC(p.ExternalMACAddress); p.ExternalMACAddress != "" && err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("externalMACAddress"), p.ExternalMACAddress, err.Error()))
		}

		if err := validate.MAC(p.ProvisioningMACAddress); p.ProvisioningMACAddress != "" && err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningMACAddress"), p.ProvisioningMACAddress, err.Error()))
		}

		if p.ProvisioningMACAddress != "" && strings.EqualFold(p.ProvisioningMACAddress, p.ExternalMACAddress) {
			allErrs = append(allErrs, field.Duplicate(fldPath.Child("provisioningMACAddress"), "provisioning and external MAC addresses may not be identical"))
		}
	}

	if err := validate.URI(p.LibvirtURI); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("libvirtURI"), p.LibvirtURI, err.Error()))
	}

	allErrs = append(allErrs, validateOSImages(p, fldPath)...)

	return allErrs
}

// ValidateProvisioningNetworking checks that provisioning network requirements specified is valid.
func ValidateProvisioningNetworking(p *baremetal.Platform, n *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	switch p.ProvisioningNetwork {
	// If we do not have a provisioning network, provisioning services
	// will be run on the external network. Users must provide IP's on the
	// machine networks to host those services.
	case baremetal.DisabledProvisioningNetwork:
		allErrs = validateProvisioningNetworkDisabledSupported(p.Hosts, fldPath.Child("hosts"))

		// If set, ensure clusterProvisioningIP is in one of the machine networks
		if p.ClusterProvisioningIP != "" {
			if err := validateIPinMachineCIDR(p.ClusterProvisioningIP, n); err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterProvisioningIP"), p.ClusterProvisioningIP, fmt.Sprintf("provisioning network is disabled, %s", err.Error())))
			}
		}
	default:
		// Ensure provisioningNetworkCIDR mask is >= 64 for managed ipv6 networks due to a dnsmasq limitation
		if err := validateCIDRSize(p); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningNetworkCIDR"), p.ProvisioningNetworkCIDR.String(), err.Error()))
		}

		// Ensure provisioningNetworkCIDR doesn't overlap with any machine network
		if err := validateNoOverlapMachineCIDR(&p.ProvisioningNetworkCIDR.IPNet, n); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningNetworkCIDR"), p.ProvisioningNetworkCIDR.String(), err.Error()))
		}

		// Ensure clusterProvisioningIP is in the provisioningNetworkCIDR
		if !p.ProvisioningNetworkCIDR.Contains(net.ParseIP(p.ClusterProvisioningIP)) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterProvisioningIP"), p.ClusterProvisioningIP, fmt.Sprintf("%q is not in the provisioning network", p.ClusterProvisioningIP)))
		}

		// Ensure provisioningNetworkCIDR does not have any host bits set
		expectedIP := p.ProvisioningNetworkCIDR.IP.Mask(p.ProvisioningNetworkCIDR.Mask)
		expectedLen, _ := p.ProvisioningNetworkCIDR.Mask.Size()
		if !p.ProvisioningNetworkCIDR.IP.Equal(expectedIP) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningNetworkCIDR"), p.ProvisioningNetworkCIDR,
				fmt.Sprintf("provisioningNetworkCIDR has host bits set, expected %s/%d", expectedIP, expectedLen)))
		}

		if len(p.AdditionalNTPServers) > 0 {
			allErrs = append(allErrs, ValidateNTPServers(p.AdditionalNTPServers, fldPath)...)
		}

		if p.ProvisioningDHCPRange != "" {
			allErrs = append(allErrs, validateDHCPRange(p, fldPath)...)
		}
	}

	allErrs = append(allErrs, validateHostsBMCOnly(p.Hosts, fldPath.Child("hosts"))...)

	return allErrs
}

// validateProvisioningBootstrapNetworking checks that provisioning network requirements specified is valid for the bootstrap VM.
func validateProvisioningBootstrapNetworking(p *baremetal.Platform, fldPath *field.Path) field.ErrorList {
	errorList := field.ErrorList{}

	if interfaceValidator != nil {
		findInterface, err := interfaceValidator(p.LibvirtURI)
		if err != nil {
			errorList = append(errorList, field.InternalError(fldPath.Child("libvirtURI"), err))
			return errorList
		}

		if err := findInterface(p.ExternalBridge); err != nil {
			errorList = append(errorList, field.Invalid(fldPath.Child("externalBridge"), p.ExternalBridge, err.Error()))
		}

		if err := findInterface(p.ProvisioningBridge); p.ProvisioningNetwork != baremetal.DisabledProvisioningNetwork && err != nil {
			errorList = append(errorList, field.Invalid(fldPath.Child("provisioningBridge"), p.ProvisioningBridge, err.Error()))
		}

	}

	if err := validate.MAC(p.ExternalMACAddress); p.ExternalMACAddress != "" && err != nil {
		errorList = append(errorList, field.Invalid(fldPath.Child("externalMACAddress"), p.ExternalMACAddress, err.Error()))
	}

	if err := validate.MAC(p.ProvisioningMACAddress); p.ProvisioningMACAddress != "" && err != nil {
		errorList = append(errorList, field.Invalid(fldPath.Child("provisioningMACAddress"), p.ProvisioningMACAddress, err.Error()))
	}

	if p.ProvisioningMACAddress != "" && strings.EqualFold(p.ProvisioningMACAddress, p.ExternalMACAddress) {
		errorList = append(errorList, field.Duplicate(fldPath.Child("provisioningMACAddress"), "provisioning and external MAC addresses may not be identical"))
	}

	return errorList
}
