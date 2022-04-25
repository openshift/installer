package gcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	gcpclient "github.com/openshift/installer/pkg/client/gcp"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/validate"
)

type resourceRequirements struct {
	minimumVCpus  int64
	minimumMemory int64
}

var controlPlaneReq = resourceRequirements{
	minimumVCpus:  4,
	minimumMemory: 15360,
}

var computeReq = resourceRequirements{
	minimumVCpus:  2,
	minimumMemory: 7680,
}

// Validate executes platform-specific validation.
func Validate(client gcpclient.API, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	if err := validate.GCPClusterName(ic.ObjectMeta.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("clusterName"), ic.ObjectMeta.Name, err.Error()))
	}

	allErrs = append(allErrs, validateProject(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateNetworks(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateInstanceTypes(client, ic)...)

	return allErrs.ToAggregate()
}

// ValidateInstanceType ensures the instance type has sufficient Vcpu and Memory.
func ValidateInstanceType(client gcpclient.API, fieldPath *field.Path, project, zone, instanceType string, req resourceRequirements) field.ErrorList {
	allErrs := field.ErrorList{}

	typeMeta, err := client.GetMachineType(context.TODO(), project, zone, instanceType)
	if err != nil {
		if _, ok := err.(*googleapi.Error); ok {
			return append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, err.Error()))
		}
		return append(allErrs, field.InternalError(nil, err))
	}

	if typeMeta.GuestCpus < req.minimumVCpus {
		errMsg := fmt.Sprintf("instance type does not meet minimum resource requirements of %d vCPUs", req.minimumVCpus)
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, errMsg))
	}
	if typeMeta.MemoryMb < req.minimumMemory {
		errMsg := fmt.Sprintf("instance type does not meet minimum resource requirements of %d MB Memory", req.minimumMemory)
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, errMsg))
	}

	return allErrs
}

// validateInstanceTypes checks that the user-provided instance types are valid.
func validateInstanceTypes(client gcpclient.API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	// Get list of zones in region
	zones, err := client.GetZones(context.TODO(), ic.GCP.ProjectID, fmt.Sprintf("region eq .*%s", ic.GCP.Region))
	if err != nil {
		return append(allErrs, field.InternalError(nil, err))
	}

	// Default requirements need to be sufficient to support Control Plane instances.
	defaultInstanceReq := controlPlaneReq

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.GCP != nil && ic.ControlPlane.Platform.GCP.InstanceType != "" {
		// Default requirements can be relaxed when the controlPlane type is set explicitly.
		defaultInstanceReq = computeReq

		allErrs = append(allErrs, ValidateInstanceType(client, field.NewPath("controlPlane", "platform", "gcp"), ic.GCP.ProjectID, zones[0].Name,
			ic.ControlPlane.Platform.GCP.InstanceType, controlPlaneReq)...)
	}

	if ic.Platform.GCP.DefaultMachinePlatform != nil && ic.Platform.GCP.DefaultMachinePlatform.InstanceType != "" {
		allErrs = append(allErrs, ValidateInstanceType(client, field.NewPath("platform", "gcp", "defaultMachinePlatform"), ic.GCP.ProjectID, zones[0].Name,
			ic.Platform.GCP.DefaultMachinePlatform.InstanceType, defaultInstanceReq)...)
	}

	for idx, compute := range ic.Compute {
		fieldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.GCP != nil && compute.Platform.GCP.InstanceType != "" {
			allErrs = append(allErrs, ValidateInstanceType(client, fieldPath.Child("platform", "gcp"), ic.GCP.ProjectID, zones[0].Name,
				compute.Platform.GCP.InstanceType, computeReq)...)
		}
	}

	return allErrs
}

// ValidatePreExitingPublicDNS ensure no pre-existing DNS record exists in the public
// DNS zone for cluster's Kubernetes API.
func ValidatePreExitingPublicDNS(client gcpclient.API, ic *types.InstallConfig) error {
	// If this is an internal cluster, this check is not necessary
	if ic.Publish == types.InternalPublishingStrategy {
		return nil
	}

	record := fmt.Sprintf("api.%s.", strings.TrimSuffix(ic.ClusterDomain(), "."))

	zone, err := client.GetPublicDNSZone(context.TODO(), ic.Platform.GCP.ProjectID, ic.BaseDomain)
	if err != nil {
		var gErr *googleapi.Error
		if errors.As(err, &gErr) {
			if gErr.Code == http.StatusNotFound {
				return field.NotFound(field.NewPath("baseDomain"), fmt.Sprintf("DNS Zone (%s/%s)", ic.Platform.GCP.ProjectID, ic.BaseDomain))
			}
		}
		return field.InternalError(field.NewPath("baseDomain"), err)
	}

	rrSets, err := client.GetRecordSets(context.TODO(), ic.Platform.GCP.ProjectID, zone.Name)
	if err != nil {
		return field.InternalError(field.NewPath("baseDomain"), err)
	}

	for _, r := range rrSets {
		if strings.EqualFold(r.Name, record) {
			return field.Invalid(field.NewPath("metadata", "name"), ic.ObjectMeta.Name, fmt.Sprintf("record %s already exists in DNS Zone (%s/%s) and might be in use by another cluster, please remove it to continue", record, ic.Platform.GCP.ProjectID, zone.Name))
		}
	}
	return nil
}

func validateProject(client gcpclient.API, ic *types.InstallConfig, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.GCP.ProjectID != "" {
		projects, err := client.GetProjects(context.TODO())
		if err != nil {
			return append(allErrs, field.InternalError(fieldPath.Child("project"), err))
		}
		if _, found := projects[ic.GCP.ProjectID]; !found {
			return append(allErrs, field.Invalid(fieldPath.Child("project"), ic.GCP.ProjectID, "invalid project ID"))
		}
	}

	return allErrs
}

// validateNetworks checks that the user-provided VPC is in the project and the provided subnets are valid.
func validateNetworks(client gcpclient.API, ic *types.InstallConfig, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.GCP.Network != "" {
		_, err := client.GetNetwork(context.TODO(), ic.GCP.Network, ic.GCP.ProjectID)
		if err != nil {
			return append(allErrs, field.Invalid(fieldPath.Child("network"), ic.GCP.Network, err.Error()))
		}

		subnets, err := client.GetSubnetworks(context.TODO(), ic.GCP.Network, ic.GCP.ProjectID, ic.GCP.Region)
		if err != nil {
			return append(allErrs, field.Invalid(fieldPath.Child("network"), ic.GCP.Network, "failed to retrieve subnets"))
		}

		allErrs = append(allErrs, validateSubnet(client, ic, fieldPath.Child("computeSubnet"), subnets, ic.GCP.ComputeSubnet)...)
		allErrs = append(allErrs, validateSubnet(client, ic, fieldPath.Child("controlPlaneSubnet"), subnets, ic.GCP.ControlPlaneSubnet)...)
	}

	return allErrs
}

func validateSubnet(client gcpclient.API, ic *types.InstallConfig, fieldPath *field.Path, subnets []*compute.Subnetwork, name string) field.ErrorList {
	allErrs := field.ErrorList{}

	subnet, errMsg := findSubnet(subnets, name, ic.GCP.Network, ic.GCP.Region)
	if subnet == nil {
		return append(allErrs, field.Invalid(fieldPath, name, errMsg))
	}

	subnetIP, _, err := net.ParseCIDR(subnet.IpCidrRange)
	if err != nil {
		return append(allErrs, field.Invalid(fieldPath, name, "unable to parse subnet CIDR"))
	}

	allErrs = append(allErrs, validateMachineNetworksContainIP(fieldPath, ic.Networking.MachineNetwork, name, subnetIP)...)
	return allErrs
}

// findSubnet checks that the subnets are in the provided VPC and region.
func findSubnet(subnets []*compute.Subnetwork, userSubnet, network, region string) (*compute.Subnetwork, string) {
	for _, vpcSubnet := range subnets {
		if userSubnet == vpcSubnet.Name {
			return vpcSubnet, ""
		}
	}
	return nil, fmt.Sprintf("could not find subnet %s in network %s and region %s", userSubnet, network, region)
}

func validateMachineNetworksContainIP(fldPath *field.Path, networks []types.MachineNetworkEntry, subnetName string, ip net.IP) field.ErrorList {
	for _, network := range networks {
		if network.CIDR.Contains(ip) {
			return nil
		}
	}
	return field.ErrorList{field.Invalid(fldPath, subnetName, fmt.Sprintf("subnet CIDR range start %s is outside of the specified machine networks", ip))}
}

//ValidateEnabledServices gets all the enabled services for a project and validate if any of the required services are not enabled.
func ValidateEnabledServices(ctx context.Context, client gcpclient.API, project string) error {
	services := sets.NewString("compute.googleapis.com",
		"cloudapis.googleapis.com",
		"cloudresourcemanager.googleapis.com",
		"dns.googleapis.com",
		"iam.googleapis.com",
		"iamcredentials.googleapis.com",
		"servicemanagement.googleapis.com",
		"serviceusage.googleapis.com",
		"storage-api.googleapis.com",
		"storage-component.googleapis.com")
	projectServices, err := client.GetEnabledServices(ctx, project)

	if err != nil {
		var gErr *googleapi.Error
		if errors.As(err, &gErr) {
			if gErr.Code == http.StatusForbidden {
				logrus.Warn("Permission denied. Unable to fetch enabled services for project.")
				return nil
			}
		}
		return err
	}

	if remaining := services.Difference(sets.NewString(projectServices...)); remaining.Len() > 0 {
		return fmt.Errorf("the following required services are not enabled in this project: %s",
			strings.Join(remaining.List(), ","))
	}
	return nil
}
