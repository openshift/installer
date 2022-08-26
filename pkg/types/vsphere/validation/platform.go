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

	if len(p.VCenters) > 0 {
		allErrs = append(allErrs, validateMultiVCenter(p, fldPath)...)
	}

	// If all VIPs are empty, skip IP validation.  All VIPs are required to be defined together.
	if strings.Join([]string{p.APIVIP, p.IngressVIP}, "") != "" {
		allErrs = append(allErrs, validateVIPs(p, fldPath)...)
	}

	return allErrs
}

func validateMultiVCenter(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateVCenters(p, fldPath.Child("vcenters"))...)
	if len(allErrs) > 0 {
		// if vcenters fails validation, this will cascade to failureDomains and deploymentZones
		return allErrs
	}

	if len(p.FailureDomains) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("failureDomains"), "must be defined if vcenters is defined"))
	}
	if len(p.DeploymentZones) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("deploymentZones"), "must be defined if vcenters is defined"))
	}
	if len(allErrs) > 0 {
		// if failureDomains and deploymentZones don't exist, this will cascade to checks to follow
		return allErrs
	}

	allErrs = append(allErrs, validateFailureDomains(p, fldPath.Child("failureDomains"))...)
	if len(allErrs) > 0 {
		// if failureDomains fails validation, this will cascade to deploymentZones
		return allErrs
	}

	// DeploymentZones is optional, but if defined should pass validation
	if len(p.DeploymentZones) != 0 {
		allErrs = append(allErrs, validateDeploymentZones(p, fldPath.Child("deploymentZones"))...)
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

func validateFailureDomain(failureDomain *vsphere.FailureDomainCoordinate, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(failureDomain.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), "must specify the name"))
	}

	switch failureDomain.Type {
	case vsphere.ComputeClusterFailureDomain:
	case vsphere.DatacenterFailureDomain:
	case vsphere.HostGroupFailureDomain:
		allErrs = append(allErrs, field.Invalid(fldPath.Child("type"), failureDomain.Type, "is not supported"))
	default:
		allErrs = append(allErrs, field.Invalid(fldPath.Child("type"), failureDomain.Type, "must be ComputeCluster or Datacenter"))
	}

	if len(failureDomain.TagCategory) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("tagCategory"), "must specify a tag category"))
	}
	return allErrs
}

func validateFailureDomains(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	topologyFld := fldPath.Child("topology")
	for _, failureDomain := range p.FailureDomains {
		if len(failureDomain.Name) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("name"), "must specify the name"))
		}

		allErrs = append(allErrs, validateFailureDomain(&failureDomain.Region, fldPath.Child("region"))...)
		allErrs = append(allErrs, validateFailureDomain(&failureDomain.Zone, fldPath.Child("zone"))...)

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

		if failureDomain.Topology.Hosts != nil {
			hosts := failureDomain.Topology.Hosts
			if len(hosts.VMGroupName) == 0 {
				allErrs = append(allErrs, field.Required(topologyFld.Child("hosts").Child("vmGroupName"), "must specify the vmGroupName"))
			}
			if len(hosts.HostGroupName) == 0 {
				allErrs = append(allErrs, field.Required(topologyFld.Child("hosts").Child("hostGroupName"), "must specify the hostGroupName"))
			}
		}
	}

	return allErrs
}

func validateDeploymentZones(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for _, deploymentZone := range p.DeploymentZones {
		var vCenter *vsphere.VCenter
		var failureDomain *vsphere.FailureDomain

		if len(deploymentZone.Server) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("server"), "must specify the hostname or IP address of a defined vCenter server"))
			return allErrs
		}

		for _, testVcenter := range p.VCenters {
			if testVcenter.Server == deploymentZone.Server {
				vCenter = &testVcenter
				break
			}
		}

		if vCenter == nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("server"), deploymentZone.Server, "server does not exist in vcenters"))
			return allErrs
		}

		if len(deploymentZone.Name) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("name"), "must specify the name"))
		}

		if len(deploymentZone.ControlPlane) != 0 {
			switch deploymentZone.ControlPlane {
			case vsphere.Allowed:
				break
			case vsphere.NotAllowed:
				break
			default:
				allErrs = append(allErrs, field.Invalid(fldPath.Child("controlPlane"), deploymentZone.ControlPlane, "valid values are Allowed and NotAllowed"))
			}
		}

		if len(deploymentZone.FailureDomain) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("failureDomain"), "must specify the failureDomain name"))
			return allErrs
		}

		for _, testFailureDomain := range p.FailureDomains {
			if testFailureDomain.Name == deploymentZone.FailureDomain {
				failureDomain = &testFailureDomain
				break
			}
		}
		if failureDomain == nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("failureDomain"), deploymentZone.FailureDomain, "does not exist in failureDomains"))
			return allErrs
		}

		if deploymentZone.PlacementConstraint.Folder != "" {
			prefix := "^\\/(.*?)\\/vm\\/(.*?)"
			match, _ := regexp.MatchString(prefix, deploymentZone.PlacementConstraint.Folder)
			if match == false {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("placementConstraint", "folder"), deploymentZone.PlacementConstraint.Folder, fmt.Sprintf("full path of folder must be provided in format /<datacenter>/vm/<folder>")))
			}
		}

		if deploymentZone.PlacementConstraint.ResourcePool != "" {
			prefix := "^\\/(.*?)\\/host\\/(.*?)\\/Resources"
			match, _ := regexp.MatchString(prefix, deploymentZone.PlacementConstraint.ResourcePool)
			if match == false {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("placementConstraint", "resourcePool"), deploymentZone.PlacementConstraint.ResourcePool, fmt.Sprintf("full path of resource pool must be provided in format /<datacenter>/host/<cluster>/Resources/<resource-pool>")))
			}
		}

		if len(failureDomain.Topology.Datacenter) > 0 {
			datacenterInVCenter := false
			for _, datacenter := range vCenter.Datacenters {
				if datacenter == failureDomain.Topology.Datacenter {
					datacenterInVCenter = true
					break
				}
			}
			if datacenterInVCenter == false {
				allErrs = append(allErrs, field.Invalid(fldPath.Root().Child("failureDomains", "topology", "datacenter"), failureDomain.Topology.Datacenter, fmt.Sprintf("datacenter %s in failure domain topology does not exist in associated vCenter", failureDomain.Topology.Datacenter)))
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

	allErrs = append(allErrs, validateVIPs(p, fldPath)...)
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
