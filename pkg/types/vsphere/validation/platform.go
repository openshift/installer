package validation

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *vsphere.Platform, agentBasedInstallation bool, fldPath *field.Path, c *types.InstallConfig) field.ErrorList {
	isLegacyUpi := false
	// This is to cover existing UPI non-zonal case
	// where neither network or cluster is required.
	// In 4.13 we will warn for this, in later releases this
	// should be removed.

	if p.DeprecatedNetwork == "" && p.DeprecatedCluster == "" && p.DeprecatedVCenter != "" {
		isLegacyUpi = true
	}

	allErrs := field.ErrorList{}
	// diskType is optional, but if provided should pass validation
	if len(p.DiskType) != 0 {
		allErrs = append(allErrs, validateDiskType(p, fldPath)...)
	}

	if !agentBasedInstallation {
		if len(p.VCenters) == 0 {
			return append(allErrs, field.Required(fldPath.Child("vcenters"), "must be defined"))
		}
		allErrs = append(allErrs, validateVCenters(p, fldPath.Child("vcenters"))...)

		if len(p.FailureDomains) == 0 {
			return append(allErrs, field.Required(fldPath.Child("failureDomains"), "must be defined"))
		}
		allErrs = append(allErrs, validateFailureDomains(p, fldPath.Child("failureDomains"), isLegacyUpi)...)

		// Validate hosts if configured for static IP
		if p.Hosts != nil {
			allErrs = append(allErrs, validateHosts(p, c, fldPath.Child("hosts"))...)
		}
	}

	// Platform fields only allowed in TechPreviewNoUpgrade
	if c.FeatureSet != configv1.TechPreviewNoUpgrade {
		if c.VSphere.LoadBalancer != nil {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("loadBalancer"), "load balancer is not supported in this feature set"))
		}
	}

	if c.VSphere.LoadBalancer != nil {
		if !validateLoadBalancer(c.VSphere.LoadBalancer.Type) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("loadBalancer", "type"), c.VSphere.LoadBalancer.Type, "invalid load balancer type"))
		}
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

func validateFailureDomains(p *vsphere.Platform, fldPath *field.Path, isLegacyUpi bool) field.ErrorList {
	allErrs := field.ErrorList{}
	topologyFld := fldPath.Child("topology")
	var associatedVCenter *vsphere.VCenter
	for _, failureDomain := range p.FailureDomains {
		if len(failureDomain.Name) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("name"), "must specify the name"))
		}
		if len(failureDomain.Server) > 0 {
			for i, vcenter := range p.VCenters {
				if vcenter.Server == failureDomain.Server {
					associatedVCenter = &p.VCenters[i]
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
		} else {
			datastore := failureDomain.Topology.Datastore

			datastorePathRegexp := regexp.MustCompile(`^/(.*?)/datastore/(.*?)$`)
			datastorePathParts := datastorePathRegexp.FindStringSubmatch(datastore)
			if len(datastorePathParts) < 3 {
				return append(allErrs, field.Invalid(topologyFld.Child("datastore"), datastore, "full path of datastore must be provided in format /<datacenter/datastore/<datastore>"))
			}

			if !strings.Contains(failureDomain.Topology.Datastore, failureDomain.Topology.Datacenter) {
				return append(allErrs, field.Invalid(topologyFld.Child("datastore"), failureDomain.Topology.Datastore, "the datastore defined does not exist in the correct datacenter"))
			}
		}

		if len(failureDomain.Topology.Networks) == 0 {
			if isLegacyUpi {
				logrus.Warn("network field empty is now deprecated, in later releases this field will be required.")
			} else {
				allErrs = append(allErrs, field.Required(topologyFld.Child("networks"), "must specify a network"))
			}
		}

		// Folder in failuredomain is optional
		if len(failureDomain.Topology.Folder) != 0 {
			folderPathRegexp := regexp.MustCompile(`^/(.*?)/vm/(.*?)$`)
			folderPathParts := folderPathRegexp.FindStringSubmatch(failureDomain.Topology.Folder)
			if len(folderPathParts) < 3 {
				return append(allErrs, field.Invalid(topologyFld.Child("folder"), failureDomain.Topology.Folder, "full path of folder must be provided in format /<datacenter>/vm/<folder>"))
			}

			if !strings.Contains(failureDomain.Topology.Folder, failureDomain.Topology.Datacenter) {
				return append(allErrs, field.Invalid(topologyFld.Child("folder"), failureDomain.Topology.Folder, "the folder defined does not exist in the correct datacenter"))
			}
		}

		if len(failureDomain.Topology.ComputeCluster) == 0 {
			if isLegacyUpi {
				logrus.Warn("cluster field empty is not deprecated, in later releases this field will be required.")
			} else {
				allErrs = append(allErrs, field.Required(topologyFld.Child("computeCluster"), "must specify a computeCluster"))
			}
		} else {
			computeCluster := failureDomain.Topology.ComputeCluster
			clusterPathRegexp := regexp.MustCompile(`^/(.*?)/host/(.*?)$`)
			clusterPathParts := clusterPathRegexp.FindStringSubmatch(computeCluster)
			if len(clusterPathParts) < 3 {
				return append(allErrs, field.Invalid(topologyFld.Child("computeCluster"), computeCluster, "full path of compute cluster must be provided in format /<datacenter>/host/<cluster>"))
			}
			datacenterName := clusterPathParts[1]

			if len(failureDomain.Topology.Datacenter) != 0 && datacenterName != failureDomain.Topology.Datacenter {
				return append(allErrs, field.Invalid(topologyFld.Child("computeCluster"), computeCluster, fmt.Sprintf("compute cluster must be in datacenter %s", failureDomain.Topology.Datacenter)))
			}
		}

		if len(failureDomain.Topology.ResourcePool) != 0 {
			resourcePool := failureDomain.Topology.ResourcePool
			resourcePoolRegexp := regexp.MustCompile(`^\/(.*?)\/host\/(.*?)\/(.*?)$`)
			resourcePoolPathParts := resourcePoolRegexp.FindStringSubmatch(resourcePool)
			if len(resourcePoolPathParts) < 4 {
				return append(allErrs, field.Invalid(topologyFld.Child("resourcePool"), resourcePool, "full path of resource pool must be provided in format /<datacenter>/host/<cluster>/..."))
			}
			datacenterName := resourcePoolPathParts[1]
			clusterName := resourcePoolPathParts[2]
			if len(failureDomain.Topology.Datacenter) != 0 && datacenterName != failureDomain.Topology.Datacenter {
				return append(allErrs, field.Invalid(topologyFld.Child("resourcePool"), resourcePool, fmt.Sprintf("resource pool must be in datacenter %s", failureDomain.Topology.Datacenter)))
			}
			if len(failureDomain.Topology.ComputeCluster) != 0 && !strings.Contains(failureDomain.Topology.ComputeCluster, clusterName) {
				return append(allErrs, field.Invalid(topologyFld.Child("resourcePool"), resourcePool, fmt.Sprintf("resource pool must be in compute cluster %s", failureDomain.Topology.ComputeCluster)))
			}
		}
	}

	return allErrs
}

// validateDiskType checks that the specified diskType is valid.
func validateDiskType(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	validDiskTypes := sets.NewString(string(vsphere.DiskTypeThin), string(vsphere.DiskTypeThick), string(vsphere.DiskTypeEagerZeroedThick))
	if !validDiskTypes.Has(string(p.DiskType)) {
		errMsg := fmt.Sprintf("diskType must be one of %v", validDiskTypes.List())
		allErrs = append(allErrs, field.Invalid(fldPath.Child("diskType"), p.DiskType, errMsg))
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

// validateHosts.
func validateHosts(p *vsphere.Platform, installConfig *types.InstallConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// Validate hosts counts match desired replicas
	allErrs = append(allErrs, validateHostsCount(p.Hosts, installConfig, fldPath)...)

	// Iterate through hosts
	for _, host := range p.Hosts {
		// Check Role (bootstrap, control-plane, compute)
		allErrs = append(allErrs, validateHostRole(host, fldPath)...)

		// Check failure domain (must exist in failure domains)
		if host.FailureDomain != "" {
			allErrs = append(allErrs, validateHostFailureDomain(host, p.FailureDomains, fldPath)...)
		}

		// Check networking
		if host.NetworkDevice == nil {
			allErrs = append(allErrs, field.Required(fldPath.Child("networkDevice"), "must specify networkDevice configuration"))
		} else {
			allErrs = append(allErrs, validateHostNetworking(host.NetworkDevice, fldPath)...)
		}
	}

	return allErrs
}

// validateHostRole returns error if the host role is invalid.
func validateHostRole(host *vsphere.Host, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	validHostRoles := sets.NewString(vsphere.BootstrapRole, vsphere.ControlPlaneRole, vsphere.ComputeRole)
	if !validHostRoles.Has(host.Role) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("role"), host.Role, validHostRoles.List()))
	}
	return allErrs
}

// validateHostsCount ensure that the number of hosts is enough to cover the
// ControlPlane and Compute replicas. Hosts without role will be considered
// eligible for the ControlPlane or Compute requirements.
func validateHostsCount(hosts []*vsphere.Host, installConfig *types.InstallConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	numRequiredControlPlane := int64(0)
	if installConfig.ControlPlane != nil && installConfig.ControlPlane.Replicas != nil {
		numRequiredControlPlane += *installConfig.ControlPlane.Replicas
	}

	numRequiredWorkers := int64(0)
	for _, worker := range installConfig.Compute {
		if worker.Replicas != nil {
			numRequiredWorkers += *worker.Replicas
		}
	}

	numControlPlane := int64(0)
	numWorkers := int64(0)
	numBootstrap := int64(0)
	for _, h := range hosts {
		switch {
		case h.IsControlPlane():
			numControlPlane++
		case h.IsCompute():
			numWorkers++
		case h.IsBootstrap():
			numBootstrap++
		}
	}

	if numBootstrap != 1 {
		allErrs = append(allErrs, field.Invalid(fldPath, "bootstrap", "a single host with the bootstrap role must be defined"))
	}

	if numControlPlane < numRequiredControlPlane {
		errMsg := fmt.Sprintf("not enough hosts found (%v) to support all the configured control plane replicas (%v)", numControlPlane, numRequiredControlPlane)
		allErrs = append(allErrs, field.Invalid(fldPath, "control-plane", errMsg))
	} else if numControlPlane > numRequiredControlPlane {
		errMsg := fmt.Sprintf("too many hosts found (%v) for the configured control plane replicas (%v)", numControlPlane, numRequiredControlPlane)
		allErrs = append(allErrs, field.Invalid(fldPath, "control-plane", errMsg))
	}

	if numWorkers < numRequiredWorkers {
		errMsg := fmt.Sprintf("not enough hosts found (%v) to support all the configured compute replicas (%v)", numWorkers, numRequiredWorkers)
		allErrs = append(allErrs, field.Invalid(fldPath, "compute", errMsg))
	} else if numWorkers > numRequiredWorkers {
		errMsg := fmt.Sprintf("too many hosts found (%v) for the configured compute replicas (%v)", numWorkers, numRequiredWorkers)
		allErrs = append(allErrs, field.Invalid(fldPath, "compute", errMsg))
	}

	return allErrs
}

// validateHostFailureDomain returns error if the FailureDomain is not found.
func validateHostFailureDomain(host *vsphere.Host, fds []vsphere.FailureDomain, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	fdFound := false
	for _, domain := range fds {
		if domain.Name == host.FailureDomain {
			fdFound = true
			break
		}
	}
	if !fdFound {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("failureDomain"), host.FailureDomain, "failure domain not found"))
	}
	return allErrs
}

// validateHostNetworking checks all fields related to networking for a host.  If any errors are found, they will
// be returned (invalid IP, IP required).
func validateHostNetworking(network *vsphere.NetworkDeviceSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// Check ip addresses
	if len(network.IPAddrs) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("ipAddrs"), "must specify a IP"))
	}
	for _, ip := range network.IPAddrs {
		allErrs = append(allErrs, validateIPWithCidr(ip, true, fldPath.Child("ipAddrs"))...)
	}

	// Check nameservers
	if len(network.Nameservers) > 3 {
		allErrs = append(allErrs, field.TooMany(fldPath.Child("nameservers"), len(network.Nameservers), 3))
	}
	for _, nameserver := range network.Nameservers {
		allErrs = append(allErrs, validateIP(nameserver, false, fldPath.Child("nameservers"))...)
	}

	// Check gateway
	allErrs = append(allErrs, validateIP(network.Gateway, false, fldPath.Child("gateway"))...)

	return allErrs
}

// validateIPWithCidr checks IP/CIDR value to see if it is valid.  If IP is required, an error will be returned if
// the IP is not specified.
func validateIPWithCidr(ip string, req bool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if ip == "" && req {
		allErrs = append(allErrs, field.Required(fldPath, "must specify a IP address with CIDR"))
	} else if ip != "" {
		if _, _, valErr := net.ParseCIDR(ip); valErr != nil {
			allErrs = append(allErrs, field.Invalid(fldPath, ip, valErr.Error()))
		}
	}
	return allErrs
}

// validateIP checks IP value to see if it is valid.  If IP is required, an error will be returned if
// the IP is not specified.
func validateIP(ip string, req bool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if ip == "" && req {
		allErrs = append(allErrs, field.Required(fldPath, "must specify a IP"))
	} else if ip != "" {
		if valErr := validate.IP(ip); valErr != nil {
			allErrs = append(allErrs, field.Invalid(fldPath, ip, valErr.Error()))
		}
	}
	return allErrs
}
