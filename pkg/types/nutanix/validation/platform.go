package validation

import (
	"fmt"
	"regexp"

	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
// nolint:gocyclo
func ValidatePlatform(p *nutanix.Platform, fldPath *field.Path, c *types.InstallConfig, usingAgentMethod bool) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(p.PrismCentral.Endpoint.Address) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismCentral").Child("endpoint").Child("address"),
			"must specify the Prism Central endpoint address"))
	} else {
		if err := validate.Host(p.PrismCentral.Endpoint.Address); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("prismCentral").Child("endpoint").Child("address"),
				p.PrismCentral.Endpoint.Address, "must be the domain name or IP address of the Prism Central"))
		}
	}

	if p.PrismCentral.Endpoint.Port < 1 || p.PrismCentral.Endpoint.Port > 65535 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("prismCentral").Child("endpoint").Child("port"),
			p.PrismCentral.Endpoint.Port, "The Prism Central endpoint port is invalid, must be in the range of 1 to 65535"))
	}

	if len(p.PrismCentral.Username) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismCentral").Child("username"),
			"must specify the Prism Central username"))
	}

	if len(p.PrismCentral.Password) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismCentral").Child("password"),
			"must specify the Prism Central password"))
	}

	// Currently we only support one Prism Element for an OpenShift cluster
	if len(p.PrismElements) != 1 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismElements"), "must specify one Prism Element"))
	}

	for _, pe := range p.PrismElements {
		if len(pe.UUID) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("prismElements").Child("uuid"),
				"must specify the Prism Element UUID"))
		}

		if len(pe.Endpoint.Address) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("prismElements").Child("endpoint").Child("address"),
				"must specify the Prism Element endpoint address"))
		} else {
			if err := validate.Host(pe.Endpoint.Address); err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("prismElements").Child("endpoint").Child("address"),
					pe.Endpoint.Address, "must be the domain name or IP address of the Prism Element (cluster)"))
			}
		}

		if pe.Endpoint.Port < 1 || pe.Endpoint.Port > 65535 {
			allErrs = append(allErrs, field.Required(fldPath.Child("prismElements").Child("endpoint").Child("port"),
				"The Prism Element endpoint port is invalid, must be in the range of 1 to 65535"))
		}
	}

	// validate subnets configuration
	if errs := validateSubnets(fldPath.Child("subnetUUIDs"), p.SubnetUUIDs); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	// For the agent-installer, the below fields are ignored. So we do not need to validate them.
	if usingAgentMethod {
		return allErrs
	}

	if c.Nutanix.LoadBalancer != nil {
		if !validateLoadBalancer(c.Nutanix.LoadBalancer.Type) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("loadBalancer", "type"), c.Nutanix.LoadBalancer.Type, "invalid load balancer type"))
		}
	}

	if c.Nutanix.DNSRecordsType == configv1.DNSRecordsTypeExternal && (c.Nutanix.LoadBalancer == nil || c.Nutanix.LoadBalancer.Type != configv1.LoadBalancerTypeUserManaged) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("dnsRecordsType"), c.Nutanix.DNSRecordsType, "external DNS records can only be configured with user-managed loadbalancers"))
	}

	// validate prismAPICallTimeout if configured
	if p.PrismAPICallTimeout != nil {
		if *p.PrismAPICallTimeout <= 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("prismAPICallTimeout"), *p.PrismAPICallTimeout, "must be a positive integer value"))
		}
	}

	// validate failureDomains if configured
	if len(p.FailureDomains) > 0 {
		pattern := "[a-z0-9]([-a-z0-9]*[a-z0-9])?"
		rexp, err := regexp.Compile(pattern)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("failureDomain", "name"), fmt.Errorf("fail to compile the pattern %q: %w", pattern, err)))
		} else {
			for _, fd := range p.FailureDomains {
				if !rexp.MatchString(fd.Name) {
					allErrs = append(allErrs, field.Invalid(fldPath.Child("failureDomain", "name"), fd.Name, fmt.Sprintf("failureDomain name should match the pattern %q.", pattern)))
				}

				if fd.PrismElement.UUID == "" {
					allErrs = append(allErrs, field.Required(fldPath.Child("failureDomain", "prismElement", "uuid"), "failureDomain prismElement uuid cannot be empty"))
				}

				// validate subnets configuration
				if errs := validateSubnets(fldPath.Child("failureDomain", "subnetUUIDs"), fd.SubnetUUIDs); len(errs) > 0 {
					allErrs = append(allErrs, errs...)
				}

				for _, sc := range fd.StorageContainers {
					if sc.ReferenceName == "" {
						allErrs = append(allErrs, field.Required(fldPath.Child("failureDomain", "storageContainers", "referenceName"), fmt.Sprintf("failureDomain %q: missing storageContainer referenceName", fd.Name)))
					}

					if sc.UUID == "" {
						allErrs = append(allErrs, field.Required(fldPath.Child("failureDomain", "storageContainers", "uuid"), fmt.Sprintf("failureDomain %q: missing storageContainer uuid", fd.Name)))
					}
				}

				for _, dsImg := range fd.DataSourceImages {
					if dsImg.ReferenceName == "" {
						allErrs = append(allErrs, field.Required(fldPath.Child("failureDomain", "dataSourceImages", "referenceName"), fmt.Sprintf("failureDomain %q: missing dataSourceImage referenceName", fd.Name)))
					}

					if dsImg.UUID == "" && dsImg.Name == "" {
						allErrs = append(allErrs, field.Required(fldPath.Child("failureDomain", "dataSourceImages"), fmt.Sprintf("failureDomain %q: both the dataSourceImage's uuid and name are empty, you need to configure one.", fd.Name)))
					}
				}
			}
		}
	}

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

// validateSubnets validates the input subnetUUIDs meet the configuration requirements.
func validateSubnets(fldPath *field.Path, subnetUUIDs []string) field.ErrorList {
	var errs field.ErrorList

	count := len(subnetUUIDs)
	switch {
	case count == 0 || subnetUUIDs[0] == "":
		errs = append(errs, field.Required(fldPath, "must specify at least one subnet"))
	case count > 32:
		errs = append(errs, field.TooMany(fldPath, count, 32))
	default:
		// check duplication
		visited := make(map[string]bool, 0)
		for _, uuid := range subnetUUIDs {
			if _, ok := visited[uuid]; ok {
				errs = append(errs, field.Invalid(fldPath, uuid, "should not configure duplicate value"))
			} else {
				visited[uuid] = true
			}
		}
	}

	return errs
}
