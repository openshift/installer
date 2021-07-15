package validation

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/validate"
)

// dynamicProvisioningValidator is a function that validates certain fields in the platform.
type dynamicProvisioningValidator func(*baremetal.Platform, *field.Path) field.ErrorList

// dynamicProvisioningValidators is an array of dynamicProvisioningValidator functions. This array can be added to by an init function, and
// is intended to be used for validations that require dependencies not built with the default tags, e.g. libvirt
// libraries.
var dynamicProvisioningValidators []dynamicProvisioningValidator

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

func validateOSImageURI(uris string) error {
	URIs := strings.Split(uris, ",")
	for _, uri := range URIs {
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
	}
	return nil
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
	hostErrs := field.ErrorList{}

	values := make(map[string]map[interface{}]struct{})

	//Initialize a new validator and register a custom validation rule for the tag `uniqueField`
	validate := validator.New()
	validate.RegisterValidation("uniqueField", func(fl validator.FieldLevel) bool {
		valueFound := false
		fieldName := fl.Parent().Type().Name() + "." + fl.FieldName()
		fieldValue := fl.Field().Interface()

		if fl.Field().Type().Comparable() {
			if _, present := values[fieldName]; !present {
				values[fieldName] = make(map[interface{}]struct{})
			}

			fieldValues := values[fieldName]
			if _, valueFound = fieldValues[fieldValue]; !valueFound {
				fieldValues[fieldValue] = struct{}{}
			}
		} else {
			panic(fmt.Sprintf("Cannot apply validation rule 'uniqueField' on field %s", fl.FieldName()))
		}

		return !valueFound
	})

	//Apply validations and translate errors
	fldPath = fldPath.Child("hosts")

	for idx, host := range hosts {
		err := validate.StructFiltered(host, filter)
		if err != nil {
			hostType := reflect.TypeOf(hosts).Elem().Elem().Name()
			for _, err := range err.(validator.ValidationErrors) {
				childName := fldPath.Index(idx).Child(err.Namespace()[len(hostType)+1:])
				switch err.Tag() {
				case "required":
					hostErrs = append(hostErrs, field.Required(childName, "missing "+err.Field()))
				case "uniqueField":
					hostErrs = append(hostErrs, field.Duplicate(childName, err.Value()))
				}
			}
		}
	}

	return hostErrs
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
	platformErrs := field.ErrorList{}

	validate := validator.New()

	customErrs := make(map[string]error)
	validate.RegisterValidation("osimageuri", func(fl validator.FieldLevel) bool {
		err := validateOSImageURI(fl.Field().String())
		if err != nil {
			customErrs[fl.FieldName()] = err
		}
		return err == nil
	})
	validate.RegisterValidation("urlexist", func(fl validator.FieldLevel) bool {
		if res, err := http.Head(fl.Field().String()); err == nil {
			return res.StatusCode == http.StatusOK
		}
		return false
	})
	err := validate.Struct(p)

	if err != nil {
		baseType := reflect.TypeOf(p).Elem().Name()
		for _, err := range err.(validator.ValidationErrors) {
			childName := fldPath.Child(err.Namespace()[len(baseType)+1:])
			switch err.Tag() {
			case "osimageuri":
				platformErrs = append(platformErrs, field.Invalid(childName, err.Value(), customErrs[err.Field()].Error()))
			case "urlexist":
				platformErrs = append(platformErrs, field.NotFound(childName, err.Value()))
			}
		}
	}

	return platformErrs
}

func validateHostsCount(hosts []*baremetal.Host, installConfig *types.InstallConfig) error {

	hostsNum := int64(len(hosts))
	counter := int64(0)

	for _, worker := range installConfig.Compute {
		if worker.Replicas != nil {
			counter += *worker.Replicas
		}
	}
	if installConfig.ControlPlane != nil && installConfig.ControlPlane.Replicas != nil {
		counter += *installConfig.ControlPlane.Replicas
	}
	if hostsNum < counter {
		return fmt.Errorf("not enough hosts found (%v) to support all the configured ControlPlane and Compute replicas (%v)", hostsNum, counter)
	}

	return nil
}

// ensure that the bootMode field contains a valid value
func validateBootMode(hosts []*baremetal.Host, fldPath *field.Path) (errors field.ErrorList) {
	for idx, host := range hosts {
		switch host.BootMode {
		case "", baremetal.UEFI, baremetal.UEFISecureBoot, baremetal.Legacy:
		default:
			valid := []string{string(baremetal.UEFI), string(baremetal.UEFISecureBoot), string(baremetal.Legacy)}
			errors = append(errors, field.NotSupported(fldPath.Index(idx).Child("bootMode"), host.BootMode, valid))
		}
	}
	return
}

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *baremetal.Platform, n *types.Networking, fldPath *field.Path, c *types.InstallConfig) field.ErrorList {
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

	if p.Hosts == nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("hosts"), p.Hosts, "bare metal hosts are missing"))
	}

	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}

	if err := validate.IP(p.APIVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	}

	if err := validateIPinMachineCIDR(p.APIVIP, n); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	}

	if err := validate.IP(p.IngressVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
	}

	if err := validateIPinMachineCIDR(p.IngressVIP, n); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
	}

	if err := validateHostsCount(p.Hosts, c); err != nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("Hosts"), err.Error()))
	}

	allErrs = append(allErrs, validateHostsWithoutBMC(p.Hosts, fldPath)...)

	allErrs = append(allErrs, validateBootMode(p.Hosts, fldPath.Child("Hosts"))...)

	return allErrs
}

// ValidateProvisioning checks that provisioning network requirements specified is valid.
func ValidateProvisioning(p *baremetal.Platform, n *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	switch p.ProvisioningNetwork {
	// If we do not have a provisioning network, provisioning services
	// will be run on the external network. Users must provide IP's on the
	// machine networks to host those services.
	case baremetal.DisabledProvisioningNetwork:
		// If set, ensure bootstrapProvisioningIP is in one of the machine networks
		if p.BootstrapProvisioningIP != "" {
			if err := validateIPinMachineCIDR(p.BootstrapProvisioningIP, n); err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapProvisioningIP"), p.BootstrapProvisioningIP, fmt.Sprintf("provisioning network is disabled, %s", err.Error())))
			}
		}

		// If set, ensure clusterProvisioningIP is in one of the machine networks
		if p.ClusterProvisioningIP != "" {
			if err := validateIPinMachineCIDR(p.ClusterProvisioningIP, n); err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterProvisioningIP"), p.ClusterProvisioningIP, fmt.Sprintf("provisioning network is disabled, %s", err.Error())))
			}
		}
	default:
		// Ensure provisioningNetworkCIDR doesn't overlap with any machine network
		if err := validateNoOverlapMachineCIDR(&p.ProvisioningNetworkCIDR.IPNet, n); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningNetworkCIDR"), p.ProvisioningNetworkCIDR.String(), err.Error()))
		}

		// Ensure bootstrapProvisioningIP is in the provisioningNetworkCIDR
		if !p.ProvisioningNetworkCIDR.Contains(net.ParseIP(p.BootstrapProvisioningIP)) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapProvisioningIP"), p.BootstrapProvisioningIP, fmt.Sprintf("%q is not in the provisioning network", p.BootstrapProvisioningIP)))
		}

		// Ensure clusterProvisioningIP is in the provisioningNetworkCIDR
		if !p.ProvisioningNetworkCIDR.Contains(net.ParseIP(p.ClusterProvisioningIP)) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterProvisioningIP"), p.ClusterProvisioningIP, fmt.Sprintf("%q is not in the provisioning network", p.ClusterProvisioningIP)))
		}

		if err := validateIPNotinMachineCIDR(p.BootstrapProvisioningIP, n); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapProvisioningIP"), p.BootstrapProvisioningIP, err.Error()))
		}

		if p.ProvisioningDHCPRange != "" {
			allErrs = append(allErrs, validateDHCPRange(p, fldPath)...)
		}
		// Make sure the provisioning interface is set.  Very little we can do to validate this  as it's not on this machine.
		if p.ProvisioningNetworkInterface == "" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningNetworkInterface"), p.ProvisioningNetworkInterface, "no provisioning network interface is configured, please set this value to be the interface on the provisioning network on your cluster's baremetal hosts"))
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

	allErrs = append(allErrs, validateHostsBMCOnly(p.Hosts, fldPath)...)

	for _, validator := range dynamicProvisioningValidators {
		allErrs = append(allErrs, validator(p, fldPath)...)
	}

	return allErrs
}
