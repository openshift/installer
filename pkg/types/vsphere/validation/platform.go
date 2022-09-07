package validation

import (
	"fmt"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
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

	// folder is optional, but if provided should pass validation
	if len(p.Folder) != 0 {
		allErrs = append(allErrs, validateFolder(p, fldPath)...)
	}

	// resource pool is optional, but if provided should pass validation
	if len(p.ResourcePool) != 0 {
		allErrs = append(allErrs, validateResourcePool(p, fldPath)...)
	}

	// diskType is optional, but if provided should pass validation
	if len(p.DiskType) != 0 {
		allErrs = append(allErrs, validateDiskType(p, fldPath)...)
	}

	if len(p.FailureDomains) > 0 || len(p.VCenters) > 0 {
		allErrs = append(allErrs, validateMultiZone(p, fldPath)...)
	}

	return allErrs
}

func validateMultiZone(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(p.VCenters) == 0 {
		// if p.VCenters is empty, populate a single vCenter based on the legacy platform spec
		p.VCenters = append(p.VCenters, vsphere.VCenter{
			Server:      p.VCenter,
			Port:        443,
			Username:    p.Username,
			Password:    p.Password,
			Datacenters: []string{p.Datacenter},
		})

		// populate failure domains that dont explicitly define a server
		for idx, failureDomain := range p.FailureDomains {
			if len(failureDomain.Server) == 0 {
				p.FailureDomains[idx].Server = p.VCenter
			}
			if len(failureDomain.Topology.Datacenter) == 0 {
				p.FailureDomains[idx].Topology.Datacenter = p.Datacenter
			}
			if len(failureDomain.Topology.ComputeCluster) == 0 {
				p.FailureDomains[idx].Topology.ComputeCluster = fmt.Sprintf("/%s/host/%s", p.Datacenter, p.Cluster)
			}
			if len(failureDomain.Topology.ResourcePool) == 0 {
				p.FailureDomains[idx].Topology.ResourcePool = p.ResourcePool
			}
			if len(failureDomain.Topology.Networks) == 0 {
				p.FailureDomains[idx].Topology.Networks = []string{p.Network}
			}
			if len(failureDomain.Topology.Datastore) == 0 {
				p.FailureDomains[idx].Topology.Datastore = p.DefaultDatastore
			}
			if len(failureDomain.Topology.Folder) == 0 {
				p.FailureDomains[idx].Topology.Folder = p.Folder
			}
		}
	}

	allErrs = append(allErrs, validateVCenters(p, fldPath.Child("vcenters"))...)
	if len(allErrs) > 0 {
		// if vcenters fails validation, this will cascade to failureDomains and deploymentZones
		return allErrs
	}

	if len(p.FailureDomains) > 0 {
		allErrs = append(allErrs, validateFailureDomains(p, fldPath.Child("failureDomains"))...)
	} else if len(p.VCenters) > 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("failureDomains"), "must be defined if vcenters is defined"))
	}

	return allErrs
}

func validateVCenters(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(p.VCenters) > 1 {
		return field.ErrorList{field.TooMany(fldPath, len(p.VCenters), 1)}
	}

	for _, vCenter := range p.VCenters {
		if len(vCenter.Server) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("server"), "must be the domain name or IP address of the vCenter"))
		} else {
			if err := validate.Host(vCenter.Server); err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("server"), vCenter.Server, "must be the domain name or IP address of the vCenter"))
			}
		}
		if len(vCenter.Username) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("username"), "must specify the username"))
		}
		if len(vCenter.Password) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("password"), "must specify the password"))
		}
		if len(vCenter.Datacenters) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("datacenters"), "must specify at least one datacenter"))
		}
	}
	return allErrs
}

func validateFailureDomains(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	topologyFld := fldPath.Child("topology")
	var associatedVCenter *vsphere.VCenter
	for _, failureDomain := range p.FailureDomains {
		if len(failureDomain.Name) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("name"), "must specify the name"))
		}
		if len(failureDomain.Server) > 0 {
			for _, vcenter := range p.VCenters {
				if vcenter.Server == failureDomain.Server {
					associatedVCenter = &vcenter
					break
				}
			}
			if associatedVCenter == nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("server"), failureDomain.Server, "server does not exist in vcenters"))
			}
		} else {
			allErrs = append(allErrs, field.Required(fldPath.Child("server"), "must specify a vCenter server"))
		}

		if len(failureDomain.Zone) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("zone"), "must specify zone tag value"))
		}

		if len(failureDomain.Region) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("region"), "must specify region tag value"))
		}

		if len(failureDomain.Topology.Datacenter) == 0 {
			allErrs = append(allErrs, field.Required(topologyFld.Child("datacenter"), "must specify a datacenter"))
		}

		if len(failureDomain.Topology.Datastore) == 0 {
			allErrs = append(allErrs, field.Required(topologyFld.Child("datastore"), "must specify a datastore"))
		}

		if len(failureDomain.Topology.ComputeCluster) == 0 {
			allErrs = append(allErrs, field.Required(topologyFld.Child("computeCluster"), "must specify a computeCluster"))
		} else {
			computeCluster := failureDomain.Topology.ComputeCluster
			clusterPathRegexp := regexp.MustCompile("^\\/(.*?)\\/host\\/(.*?)$")
			clusterPathParts := clusterPathRegexp.FindStringSubmatch(computeCluster)
			if len(clusterPathParts) < 3 {
				return append(allErrs, field.Invalid(topologyFld.Child("computeCluster"), computeCluster, "full path of compute cluster must be provided in format /<datacenter>/host/<cluster>"))
			}
			datacenterName := clusterPathParts[1]

			if len(failureDomain.Topology.Datacenter) != 0 && datacenterName != failureDomain.Topology.Datacenter {
				return append(allErrs, field.Invalid(topologyFld.Child("computeCluster"), computeCluster, fmt.Sprintf("compute cluster must be in datacenter %s", failureDomain.Topology.Datacenter)))
			}
		}
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

// validateDiskType checks that the specified diskType is valid
func validateDiskType(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	validDiskTypes := sets.NewString(string(vsphere.DiskTypeThin), string(vsphere.DiskTypeThick), string(vsphere.DiskTypeEagerZeroedThick))
	if !validDiskTypes.Has(string(p.DiskType)) {
		errMsg := fmt.Sprintf("diskType must be one of %v", validDiskTypes.List())
		allErrs = append(allErrs, field.Invalid(fldPath.Child("diskType"), p.DiskType, errMsg))
	}

	return allErrs
}
