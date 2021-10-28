package validation

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(p.VCenter) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("vCenter"), "must specify the name of the vCenter"))
	}
	if len(p.Username) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("username"), "must specify the username"))
	}
	if len(p.Password) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("password"), "must specify the password"))
	}
	if len(p.Datacenter) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("datacenter"), "must specify the datacenter"))
	}
	if len(p.DefaultDatastore) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("defaultDatastore"), "must specify the default datastore"))
	}

	if len(p.VCenter) != 0 {
		if err := validate.Host(p.VCenter); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("vCenter"), p.VCenter, "must be the domain name or IP address of the vCenter"))
		}
	}

	// validate zoning
	if len(p.VCenters) != 0 {
		allErrs = append(allErrs, validateZoning(p, fldPath)...)
	}

	// If all VIPs are empty, skip IP validation.  All VIPs are required to be defined together.
	if strings.Join([]string{p.APIVIP, p.IngressVIP}, "") != "" {
		allErrs = append(allErrs, validateVIPs(p, fldPath)...)
	}

	// folder is optional, but if provided should pass validation
	if len(p.Folder) != 0 {
		allErrs = append(allErrs, validateFolder(p, fldPath)...)
	}

	// resource pool is optional, but if provided should pass validation
	if len(p.ResourcePool) != 0 {
		allErrs = append(allErrs, validateResourcePool(p, fldPath)...)
	}

	return allErrs
}

// ValidateForProvisioning checks that the specified platform is valid.
func ValidateForProvisioning(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(p.Cluster) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("cluster"), "must specify the cluster"))
	}

	if len(p.Network) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("network"), "must specify the network"))
	}

	allErrs = append(allErrs, validateVIPs(p, fldPath)...)
	return allErrs
}

func validateZoning(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	vCentersFieldPath := fldPath.Child("vcenters")
	regionsFieldPath := vCentersFieldPath.Child("regions")
	zonesFieldPath := regionsFieldPath.Child("zones")

	// Check if the VCenters slice has entries
	if len(p.VCenters) != 0 {
		// We currently do not support more than one vCenter
		if len(p.VCenters) > 1 {
			allErrs = append(allErrs, field.Required(vCentersFieldPath, "must specify a single VCenters entry."))
		}

		for _, v := range p.VCenters {
			if len(v.Server) != 0 {
				if err := validate.Host(v.Server); err != nil {
					allErrs = append(allErrs, field.Invalid(vCentersFieldPath.Child("server"), v.Server, "must be the domain name or IP address of the vCenter"))
				}
			}
			if v.Port != 0 {
				allErrs = append(allErrs, field.Required(vCentersFieldPath.Child("port"), "is currently not supported"))
			}
			if len(v.Server) == 0 {
				allErrs = append(allErrs, field.Required(vCentersFieldPath.Child("server"), "must specify the name of the vCenter"))
			}
			if len(v.User) == 0 {
				allErrs = append(allErrs, field.Required(vCentersFieldPath.Child("username"), "must specify the username"))
			}
			if len(v.Password) == 0 {
				allErrs = append(allErrs, field.Required(vCentersFieldPath.Child("password"), "must specify the password"))
			}
			if len(v.Regions) == 0 {
				allErrs = append(allErrs, field.Required(regionsFieldPath, "must specify the regions"))
			} else {
				for _, r := range v.Regions {
					if len(r.Datacenter) == 0 {
						allErrs = append(allErrs, field.Required(regionsFieldPath.Child("datacenter"), "must specify the datacenter"))
					}
					if len(r.Name) == 0 {
						allErrs = append(allErrs, field.Required(regionsFieldPath.Child("name"), "must specify the region name"))
					}
					if len(r.Zones) == 0 {
						allErrs = append(allErrs, field.Required(zonesFieldPath, "must specify the zones"))
					} else {
						for _, z := range r.Zones {
							if len(z.Datastore) == 0 {
								allErrs = append(allErrs, field.Required(zonesFieldPath.Child("datastore"), "must specify the datastore"))
							}
							if len(z.Network) == 0 {
								allErrs = append(allErrs, field.Required(zonesFieldPath.Child("network"), "must specify a network"))
							}
							if len(z.Cluster) == 0 {
								allErrs = append(allErrs, field.Required(zonesFieldPath.Child("cluster"), "must specify a cluster"))
							}
							if len(z.Name) == 0 {
								allErrs = append(allErrs, field.Required(zonesFieldPath.Child("name"), "must specify a zone name"))
							}
						}
					}
				}
			}
		}
	}

	return allErrs
}

// validateVIPs checks that all required VIPs are provided and are valid IP addresses.
func validateVIPs(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(p.APIVIP) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("apiVIP"), "must specify a VIP for the API"))
	} else if err := validate.IP(p.APIVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	}

	if len(p.IngressVIP) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("ingressVIP"), "must specify a VIP for Ingress"))
	} else if err := validate.IP(p.IngressVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
	}

	if len(p.APIVIP) != 0 && len(p.IngressVIP) != 0 && p.APIVIP == p.IngressVIP {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, "IPs for both API and Ingress should not be the same."))
	}

	return allErrs
}

// validateFolder checks that a provided folder is an absolute path in the correct datacenter.
func validateFolder(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	dc := p.Datacenter
	if len(dc) == 0 {
		dc = "<datacenter>"
	}
	expectedPrefix := fmt.Sprintf("/%s/vm/", dc)

	if !strings.HasPrefix(p.Folder, expectedPrefix) {
		errMsg := fmt.Sprintf("folder must be absolute path: expected prefix %s", expectedPrefix)
		allErrs = append(allErrs, field.Invalid(fldPath.Child("folder"), p.Folder, errMsg))
	}

	return allErrs
}

// validateResourcePool checks that a provided resource pool is an absolute path in the correct cluster.
func validateResourcePool(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	dc := p.Datacenter
	if len(dc) == 0 {
		dc = "<datacenter>"
	}
	cluster := p.Cluster
	if len(cluster) == 0 {
		cluster = "<cluster>"
	}
	expectedPrefix := fmt.Sprintf("/%s/host/%s/Resources/", dc, cluster)

	if !strings.HasPrefix(p.ResourcePool, expectedPrefix) {
		errMsg := fmt.Sprintf("resourcePool must be absolute path: expected prefix %s", expectedPrefix)
		allErrs = append(allErrs, field.Invalid(fldPath.Child("resourcePool"), p.ResourcePool, errMsg))
	}

	return allErrs
}
