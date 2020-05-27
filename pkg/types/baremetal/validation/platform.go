package validation

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/validate"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// dynamicValidator is a function that validates certain fields in the platform.
type dynamicValidator func(*baremetal.Platform, *field.Path) field.ErrorList

// dynamicValidators is an array of dynamicValidator functions. This array can be added to by an init function, and
// is intended to be used for validations that require dependencies not built with the default tags, e.g. libvirt
// libraries.
var dynamicValidators []dynamicValidator

func validateIPinMachineCIDR(vip string, n *types.Networking) error {
	for _, network := range n.MachineNetwork {
		if network.CIDR.Contains(net.ParseIP(vip)) {
			return nil
		}
	}
	return fmt.Errorf("the virtual IP is expected to be in one of the machine networks")
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

// validateHosts checks that hosts have all required fields set with appropriate values
func validateHosts(hosts []*baremetal.Host, fldPath *field.Path) field.ErrorList {
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
		err := validate.Struct(host)
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

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *baremetal.Platform, n *types.Networking, fldPath *field.Path, c *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	if err := validate.URI(p.LibvirtURI); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("libvirtURI"), p.LibvirtURI, err.Error()))
	}

	if err := validate.IP(p.ClusterProvisioningIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningHostIP"), p.ClusterProvisioningIP, err.Error()))
	}

	if err := validate.IP(p.BootstrapProvisioningIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapProvisioningIP"), p.BootstrapProvisioningIP, err.Error()))
	}

	if p.ProvisioningNetworkCIDR != nil {
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
	}

	if p.ProvisioningDHCPRange != "" {
		dhcpRange := strings.Split(p.ProvisioningDHCPRange, ",")
		if len(dhcpRange) != 2 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningDHCPRange"), p.ProvisioningDHCPRange, "provisioning dhcp range should be in format: start_ip,end_ip"))
		} else {
			for _, ip := range dhcpRange {
				// Ensure IP is valid
				if err := validate.IP(ip); err != nil {
					allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningDHCPRange"), p.ProvisioningDHCPRange, fmt.Sprintf("%s: %s", ip, err.Error())))
				}

				// Validate IP is in the provisioning network
				if p.ProvisioningNetworkCIDR != nil && !p.ProvisioningNetworkCIDR.Contains(net.ParseIP(ip)) {
					allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningDHCPRange"), p.ProvisioningDHCPRange, fmt.Sprintf("%q is not in the provisioning network", ip)))
				}
			}
		}
	}

	// Make sure the provisioning interface is set.  Very little we can do to validate this  as it's not on this machine.
	if p.ProvisioningNetworkInterface == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningNetworkInterface"), p.ProvisioningNetworkInterface, "no provisioning network interface is configured, please set this value to be the interface on the provisioning network on your cluster's baremetal hosts"))
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

	if err := validateIPNotinMachineCIDR(p.ClusterProvisioningIP, n); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningHostIP"), p.ClusterProvisioningIP, err.Error()))
	}
	if err := validateIPNotinMachineCIDR(p.BootstrapProvisioningIP, n); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapHostIP"), p.BootstrapProvisioningIP, err.Error()))
	}

	if err := validateHostsCount(p.Hosts, c); err != nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("Hosts"), err.Error()))
	}

	allErrs = append(allErrs, validateOSImages(p, fldPath)...)

	allErrs = append(allErrs, validateHosts(p.Hosts, fldPath)...)

	for _, validator := range dynamicValidators {
		allErrs = append(allErrs, validator(p, fldPath)...)
	}

	return allErrs
}
