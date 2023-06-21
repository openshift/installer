package gcp

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

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

var (
	apiRecordType = func(ic *types.InstallConfig) string {
		return fmt.Sprintf("api.%s.", strings.TrimSuffix(ic.ClusterDomain(), "."))
	}
	apiIntRecordName = func(ic *types.InstallConfig) string {
		return fmt.Sprintf("api-int.%s.", strings.TrimSuffix(ic.ClusterDomain(), "."))
	}
)

// Validate executes platform-specific validation.
func Validate(client API, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	if err := validate.GCPClusterName(ic.ObjectMeta.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("clusterName"), ic.ObjectMeta.Name, err.Error()))
	}

	allErrs = append(allErrs, validateProject(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateNetworkProject(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateRegion(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateNetworks(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateInstanceTypes(client, ic)...)
	allErrs = append(allErrs, validateCredentialMode(client, ic)...)

	return allErrs.ToAggregate()
}

// ValidateInstanceType ensures the instance type has sufficient Vcpu and Memory.
func ValidateInstanceType(client API, fieldPath *field.Path, project, zone, instanceType string, req resourceRequirements) field.ErrorList {
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
func validateInstanceTypes(client API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	// Get list of zones in region
	zones, err := client.GetZones(context.TODO(), ic.GCP.ProjectID, fmt.Sprintf("region eq .*%s", ic.GCP.Region))
	if err != nil {
		return append(allErrs, field.InternalError(nil, err))
	} else if len(zones) == 0 {
		return append(allErrs, field.InternalError(nil, fmt.Errorf("failed to fetch instance types, this error usually occurs if the region is not found")))
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

// ValidatePreExistingPublicDNS ensure no pre-existing DNS record exists in the public
// DNS zone for cluster's Kubernetes API. If a PublicDNSZone is provided, the provided
// zone is verified against the BaseDomain. If no zone is provided, the base domain is
// checked for any public zone that can be used.
func ValidatePreExistingPublicDNS(client API, ic *types.InstallConfig) *field.Error {
	// If this is an internal cluster, this check is not necessary
	if ic.Publish == types.InternalPublishingStrategy {
		return nil
	}

	zone, err := client.GetDNSZone(context.TODO(), ic.Platform.GCP.ProjectID, ic.BaseDomain, true)
	if err != nil {
		if IsNotFound(err) {
			return field.NotFound(field.NewPath("baseDomain"), fmt.Sprintf("Public DNS Zone (%s/%s)", ic.Platform.GCP.ProjectID, ic.BaseDomain))
		}
		return field.InternalError(field.NewPath("baseDomain"), err)
	}
	return checkRecordSets(client, ic, zone, []string{apiRecordType(ic)})
}

// ValidatePrivateDNSZone ensure no pre-existing DNS record exists in the private dns zone
// matching the name that will be used for this installation.
func ValidatePrivateDNSZone(client API, ic *types.InstallConfig) *field.Error {
	if ic.GCP.Network == "" || ic.GCP.NetworkProjectID == "" {
		return nil
	}

	zone, err := client.GetDNSZone(context.TODO(), ic.GCP.ProjectID, ic.ClusterDomain(), false)
	if err != nil {
		logrus.Debug("No private DNS Zone found")
		if IsNotFound(err) {
			return field.NotFound(field.NewPath("baseDomain"), fmt.Sprintf("Private DNS Zone (%s/%s)", ic.Platform.GCP.ProjectID, ic.BaseDomain))
		}
		return field.InternalError(field.NewPath("baseDomain"), err)
	}

	// Private Zone can be nil, check to see if it was found or not
	if zone != nil {
		return checkRecordSets(client, ic, zone, []string{apiRecordType(ic), apiIntRecordName(ic)})
	}
	return nil
}

func checkRecordSets(client API, ic *types.InstallConfig, zone *dns.ManagedZone, records []string) *field.Error {
	rrSets, err := client.GetRecordSets(context.TODO(), ic.GCP.ProjectID, zone.Name)
	if err != nil {
		return field.InternalError(field.NewPath("baseDomain"), err)
	}

	setOfReturnedRecords := sets.New[string]()
	for _, r := range rrSets {
		setOfReturnedRecords.Insert(r.Name)
	}
	preexistingRecords := sets.New[string](records...).Intersection(setOfReturnedRecords)

	if preexistingRecords.Len() > 0 {
		errMsg := fmt.Sprintf("record(s) %q already exists in DNS Zone (%s/%s) and might be in use by another cluster, please remove it to continue", sets.List(preexistingRecords), ic.GCP.ProjectID, zone.Name)
		return field.Invalid(field.NewPath("metadata", "name"), ic.ObjectMeta.Name, errMsg)
	}
	return nil
}

// ValidateForProvisioning validates that the install config is valid for provisioning the cluster.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	client, err := NewClient(context.TODO())
	if err != nil {
		return err
	}

	if err := ValidatePreExistingPublicDNS(client, ic); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := ValidatePrivateDNSZone(client, ic); err != nil {
		allErrs = append(allErrs, err)
	}

	return allErrs.ToAggregate()
}

func validateProject(client API, ic *types.InstallConfig, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.GCP.ProjectID != "" {
		_, err := client.GetProjectByID(context.TODO(), ic.GCP.ProjectID)
		if err != nil {
			if IsNotFound(err) {
				return append(allErrs, field.Invalid(fieldPath.Child("project"), ic.GCP.ProjectID, "invalid project ID"))
			}
			return append(allErrs, field.InternalError(fieldPath.Child("project"), err))
		}
	}

	return allErrs
}

func validateNetworkProject(client API, ic *types.InstallConfig, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.GCP.NetworkProjectID != "" {
		_, err := client.GetProjectByID(context.TODO(), ic.GCP.NetworkProjectID)
		if err != nil {
			if IsNotFound(err) {
				return append(allErrs, field.Invalid(fieldPath.Child("networkProjectID"), ic.GCP.NetworkProjectID, "invalid project ID"))
			}
			return append(allErrs, field.InternalError(fieldPath.Child("networkProjectID"), err))
		}
	}

	return allErrs
}

// validateNetworks checks that the user-provided VPC is in the project and the provided subnets are valid.
func validateNetworks(client API, ic *types.InstallConfig, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	networkProjectID := ic.GCP.NetworkProjectID
	if networkProjectID == "" {
		networkProjectID = ic.GCP.ProjectID
	}

	if ic.GCP.Network != "" {
		_, err := client.GetNetwork(context.TODO(), ic.GCP.Network, networkProjectID)
		if err != nil {
			return append(allErrs, field.Invalid(fieldPath.Child("network"), ic.GCP.Network, err.Error()))
		}

		subnets, err := client.GetSubnetworks(context.TODO(), ic.GCP.Network, networkProjectID, ic.GCP.Region)
		if err != nil {
			return append(allErrs, field.Invalid(fieldPath.Child("network"), ic.GCP.Network, "failed to retrieve subnets"))
		}

		allErrs = append(allErrs, validateSubnet(client, ic, fieldPath.Child("computeSubnet"), subnets, ic.GCP.ComputeSubnet)...)
		allErrs = append(allErrs, validateSubnet(client, ic, fieldPath.Child("controlPlaneSubnet"), subnets, ic.GCP.ControlPlaneSubnet)...)
	}

	return allErrs
}

func validateSubnet(client API, ic *types.InstallConfig, fieldPath *field.Path, subnets []*compute.Subnetwork, name string) field.ErrorList {
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

// ValidateEnabledServices gets all the enabled services for a project and validate if any of the required services are not enabled.
// also warns the user if optional services are not enabled.
func ValidateEnabledServices(ctx context.Context, client API, project string) error {
	requiredServices := sets.NewString("compute.googleapis.com",
		"cloudresourcemanager.googleapis.com",
		"dns.googleapis.com",
		"iam.googleapis.com",
		"iamcredentials.googleapis.com",
		"serviceusage.googleapis.com")
	optionalServices := sets.NewString("cloudapis.googleapis.com",
		"servicemanagement.googleapis.com",
		"deploymentmanager.googleapis.com",
		"storage-api.googleapis.com",
		"storage-component.googleapis.com")
	projectServices, err := client.GetEnabledServices(ctx, project)
	if err != nil {
		if IsForbidden(err) {
			return errors.Wrap(err, "unable to fetch enabled services for project. Make sure 'serviceusage.googleapis.com' is enabled")
		}
		return err
	}

	if remaining := requiredServices.Difference(sets.NewString(projectServices...)); remaining.Len() > 0 {
		return fmt.Errorf("the following required services are not enabled in this project: %s",
			strings.Join(remaining.List(), ","))
	}

	if remaining := optionalServices.Difference(sets.NewString(projectServices...)); remaining.Len() > 0 {
		logrus.Warnf("the following optional services are not enabled in this project: %s",
			strings.Join(remaining.List(), ","))
	}
	return nil
}

// ValidateProjectRegion determines whether the region is valid for the project
func validateRegion(client API, ic *types.InstallConfig, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	regionFound := false

	if ic.GCP.ProjectID != "" && ic.GCP.Region != "" {
		computeRegions, err := client.GetRegions(context.TODO(), ic.GCP.ProjectID)
		if err != nil {
			return append(allErrs, field.InternalError(fieldPath.Child("project"), err))
		} else if len(computeRegions) == 0 {
			return append(allErrs, field.Invalid(fieldPath.Child("project"), ic.GCP.ProjectID, "no regions found"))
		}

		for _, region := range computeRegions {
			if regionFound = region == ic.GCP.Region; regionFound {
				break
			}
		}
	}

	if !regionFound {
		return append(allErrs, field.Invalid(fieldPath.Child("region"), ic.GCP.Region, "invalid region"))
	}
	return nil
}

// ValidateCredentialMode checks whether the credential mode is
// compatible with the authentication mode.
func ValidateCredentialMode(client API, ic *types.InstallConfig) error {
	creds := client.GetCredentials()

	if creds.JSON == nil && ic.CredentialsMode != types.ManualCredentialsMode {
		errMsg := "environmental authentication is only supported with Manual credentials mode"
		return field.Forbidden(field.NewPath("credentialsMode"), errMsg)
	}

	return nil
}

func validateCredentialMode(client API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	creds := client.GetCredentials()
	if creds.JSON == nil {
		if ic.CredentialsMode == "" {
			logrus.Warn("Currently using GCP Environmental Authentication. Please set credentialsMode to manual, or provide a service account json file.")
		} else {
			if ic.CredentialsMode != "" && ic.CredentialsMode != types.ManualCredentialsMode {
				errMsg := "environmental authentication is only supported with Manual credentials mode"
				return append(allErrs, field.Forbidden(field.NewPath("credentialsMode"), errMsg))
			}
		}
	}

	return allErrs
}
