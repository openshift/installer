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
	"github.com/openshift/installer/pkg/types/gcp"
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
func Validate(client API, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	if err := validate.GCPClusterName(ic.ObjectMeta.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("clusterName"), ic.ObjectMeta.Name, err.Error()))
	}

	allErrs = append(allErrs, validateProject(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateNetworkProject(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateRegion(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateNetworks(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateZoneProjects(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateManagedZones(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateInstanceTypes(client, ic)...)
	allErrs = append(allErrs, validateCredentialMode(client, ic)...)
	allErrs = append(allErrs, validateMarketplaceImages(client, ic)...)

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

// validateZoneProjects will validate the public and private zone projects when provided
func validateZoneProjects(client API, ic *types.InstallConfig, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	projects, err := client.GetProjects(context.TODO())
	if err != nil {
		return append(allErrs, field.InternalError(fieldPath.Child("project"), err))
	}

	// If the PublicZoneProject is empty, the value will default to ProjectID, and it won't be checked here
	if ic.GCP.PublicDNSZone != nil && ic.GCP.PublicDNSZone.ProjectID != "" {
		if _, found := projects[ic.GCP.PublicDNSZone.ProjectID]; !found {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("PublicDNSZone").Child("ProjectID"), ic.GCP.PublicDNSZone.ProjectID, "invalid public zone project"))
		}
	}

	if ic.GCP.PrivateDNSZone != nil && ic.GCP.PrivateDNSZone.ProjectID != "" {
		if _, found := projects[ic.GCP.PrivateDNSZone.ProjectID]; !found {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("PrivateDNSZone").Child("ProjectID"), ic.GCP.PrivateDNSZone.ProjectID, "invalid private zone project"))
		}
	}

	return allErrs
}

// findProject finds the correct project to use during installation. If the project id is
// provided in the zone use the project id, otherwise use the default project.
func findProject(zone *gcp.DNSZone, defaultProject string) string {
	if zone != nil && zone.ProjectID != "" {
		return zone.ProjectID
	}
	return defaultProject
}

// findDNSZone finds a zone in a project. If a project is provided in the zone, the project
// is used otherwise the default project is used.
func findDNSZone(client API, zone *gcp.DNSZone, project string) (*dns.ManagedZone, error) {
	returnedZone, err := client.GetDNSZoneByName(context.TODO(), project, zone.ID)
	if err != nil {
		return nil, err
	}
	return returnedZone, nil
}

// validateManagedZones will validate the public and private managed zones if they exist.
func validateManagedZones(client API, ic *types.InstallConfig, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.GCP.PublicDNSZone != nil && ic.GCP.PublicDNSZone.ID != "" {
		project := findProject(ic.GCP.PublicDNSZone, ic.GCP.ProjectID)
		returnedZone, err := findDNSZone(client, ic.Platform.GCP.PublicDNSZone, project)
		if err != nil {
			switch {
			case IsNotFound(err):
				allErrs = append(allErrs, field.NotFound(field.NewPath("baseDomain"), errors.Wrapf(err, "dns zone (%s/%s)", project, ic.BaseDomain).Error()))
			case IsForbidden(err):
				errMsg := errors.Wrapf(err, "unable to fetch public dns zone information: %s", ic.BaseDomain).Error()
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("publicManagedZone"), ic.GCP.PublicDNSZone.ID, errMsg))
			default:
				allErrs = append(allErrs, field.InternalError(field.NewPath("baseDomain"), err))
			}
		} else {
			// verify that the managed zone exists in the BaseDomain - Trim both values just in case
			if !strings.EqualFold(strings.TrimSuffix(returnedZone.DnsName, "."), strings.TrimSuffix(ic.BaseDomain, ".")) {
				errMsg := fmt.Sprintf("publicDNSZone does not exist in baseDomain %s", ic.BaseDomain)
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("publicDNSZone").Child("id"), ic.Platform.GCP.PublicDNSZone.ID, errMsg))
			}
		}
	}

	if ic.GCP.PrivateDNSZone != nil && ic.GCP.PrivateDNSZone.ID != "" {
		project := findProject(ic.GCP.PrivateDNSZone, ic.GCP.ProjectID)
		returnedZone, err := findDNSZone(client, ic.Platform.GCP.PrivateDNSZone, project)
		if err != nil {
			switch {
			case IsNotFound(err):
				allErrs = append(allErrs, field.NotFound(field.NewPath("clusterDomain"), errors.Wrapf(err, "dns zone (%s/%s)", project, ic.ClusterDomain()).Error()))
			case IsForbidden(err):
				errMsg := errors.Wrapf(err, "unable to fetch private dns zone information: %s", ic.ClusterDomain()).Error()
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("privateDNSZone").Child("id"), ic.GCP.PrivateDNSZone.ID, errMsg))
			default:
				allErrs = append(allErrs, field.InternalError(field.NewPath("clusterDomain"), err))
			}
		} else {
			if !strings.EqualFold(strings.TrimSuffix(returnedZone.DnsName, "."), strings.TrimSuffix(ic.ClusterDomain(), ".")) {
				errMsg := fmt.Sprintf("dns zone %s did not match expected %s", returnedZone.DnsName, ic.ClusterDomain())
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("privateManagedZone"), ic.GCP.PrivateDNSZone.ID, errMsg))
			}
		}
	}

	return allErrs
}

// ValidatePreExistingPrivateDNS ensure that the PrivateZone exists in the cluster domain
func ValidatePreExistingPrivateDNS(client API, ic *types.InstallConfig) error {
	record := fmt.Sprintf("api.%s.", strings.TrimSuffix(ic.ClusterDomain(), "."))
	project := findProject(ic.GCP.PrivateDNSZone, ic.GCP.ProjectID)

	if ic.GCP.PrivateDNSZone != nil && ic.GCP.PrivateDNSZone.ID != "" {
		zone, err := client.GetDNSZoneByName(context.TODO(), project, ic.GCP.PrivateDNSZone.ID)
		if err != nil {
			return fmt.Errorf("failed to find public zone in project %s", project)
		}

		rrSets, err := client.GetRecordSets(context.TODO(), project, zone.Name)
		if err != nil {
			return field.InternalError(field.NewPath("baseDomain"), err)
		}

		for _, r := range rrSets {
			if strings.EqualFold(r.Name, record) {
				errMsg := fmt.Sprintf("record %s already exists in DNS Zone (%s/%s) and might be in use by another cluster, please remove it to continue", record, project, zone.Name)
				return field.Invalid(field.NewPath("metadata", "name"), ic.ObjectMeta.Name, errMsg)
			}
		}
	}
	return nil
}

// ValidatePreExistingPublicDNS ensure no pre-existing DNS record exists in the public
// DNS zone for cluster's Kubernetes API. If a PublicDNSZone is provided, the provided
// zone is verified against the BaseDomain. If no zone is provided, the base domain is
// checked for any public zone that can be used.
func ValidatePreExistingPublicDNS(client API, ic *types.InstallConfig) error {
	// If this is an internal cluster, this check is not necessary
	if ic.Publish == types.InternalPublishingStrategy {
		return nil
	}

	record := fmt.Sprintf("api.%s.", strings.TrimSuffix(ic.ClusterDomain(), "."))
	project := findProject(ic.GCP.PublicDNSZone, ic.GCP.ProjectID)
	zoneName := ""

	if ic.GCP.PublicDNSZone != nil && ic.GCP.PublicDNSZone.ID != "" {
		zoneName = ic.GCP.PublicDNSZone.ID
	}

	if zoneName == "" {
		zone, err := client.GetPublicDNSZone(context.TODO(), project, ic.BaseDomain)
		if err != nil {
			return errors.Wrapf(err, "failed to find public zone in project %s", project)
		}

		zoneName = zone.Name
	}

	rrSets, err := client.GetRecordSets(context.TODO(), project, zoneName)
	if err != nil {
		return field.InternalError(field.NewPath("baseDomain"), err)
	}

	for _, r := range rrSets {
		if strings.EqualFold(r.Name, record) {
			errMsg := fmt.Sprintf("record %s already exists in DNS Zone (%s/%s) and might be in use by another cluster, please remove it to continue", record, project, zoneName)
			return field.Invalid(field.NewPath("metadata", "name"), ic.ObjectMeta.Name, errMsg)
		}
	}
	return nil
}

func validateProject(client API, ic *types.InstallConfig, fieldPath *field.Path) field.ErrorList {
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

func validateNetworkProject(client API, ic *types.InstallConfig, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.GCP.NetworkProjectID != "" {
		projects, err := client.GetProjects(context.TODO())
		if err != nil {
			return append(allErrs, field.InternalError(fieldPath.Child("networkProjectID"), err))
		}
		if _, found := projects[ic.GCP.NetworkProjectID]; !found {
			return append(allErrs, field.Invalid(fieldPath.Child("networkProjectID"), ic.GCP.NetworkProjectID, "invalid project ID"))
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

func validateMarketplaceImages(client API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	const errorMessage string = "could not find the boot image: %v"
	var err error
	var defaultImage *compute.Image
	var defaultOsImage *gcp.OSImage

	if ic.GCP.DefaultMachinePlatform != nil && ic.GCP.DefaultMachinePlatform.OSImage != nil {
		defaultOsImage = ic.GCP.DefaultMachinePlatform.OSImage
		defaultImage, err = client.GetImage(context.TODO(), defaultOsImage.Name, defaultOsImage.Project)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("platform", "gcp", "defaultMachinePlatform", "osImage"), *defaultOsImage, fmt.Sprintf(errorMessage, err)))
		}
	}

	if ic.ControlPlane != nil {
		image := defaultImage
		osImage := defaultOsImage
		if ic.ControlPlane.Platform.GCP != nil && ic.ControlPlane.Platform.GCP.OSImage != nil {
			osImage = ic.ControlPlane.Platform.GCP.OSImage
			image, err = client.GetImage(context.TODO(), osImage.Name, osImage.Project)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("controlPlane", "platform", "gcp", "osImage"), *osImage, fmt.Sprintf(errorMessage, err)))
			}
		}
		if image != nil {
			if errMsg := checkArchitecture(image.Architecture, ic.ControlPlane.Architecture, "controlPlane"); errMsg != "" {
				allErrs = append(allErrs, field.Invalid(field.NewPath("controlPlane", "platform", "gcp", "osImage"), *osImage, errMsg))
			}
		}
	}

	for idx, compute := range ic.Compute {
		image := defaultImage
		osImage := defaultOsImage
		fieldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.GCP != nil && compute.Platform.GCP.OSImage != nil {
			osImage = compute.Platform.GCP.OSImage
			image, err = client.GetImage(context.TODO(), osImage.Name, osImage.Project)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("platform", "gcp", "osImage"), *osImage, fmt.Sprintf(errorMessage, err)))
			}
		}
		if image != nil {
			if errMsg := checkArchitecture(image.Architecture, compute.Architecture, "compute"); errMsg != "" {
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("platform", "gcp", "osImage"), *osImage, errMsg))
			}
		}
	}

	return allErrs
}

func checkArchitecture(imageArch string, icArch types.Architecture, role string) string {
	const unspecifiedArch string = "ARCHITECTURE_UNSPECIFIED"
	// The possible architecture names from image.Architecture are of type string hence we cannot directly obtain the possible values
	// In the docs the possible values are ARM64, X86_64, and ARCHITECTURE_UNSPECIFIED
	// There is no simple translation between the architecture values from Google and the architecture names used in the install config so a map is used

	translateArchName := map[string]types.Architecture{
		"ARM64":  types.ArchitectureARM64,
		"X86_64": types.ArchitectureAMD64,
	}

	if imageArch == "" || imageArch == unspecifiedArch {
		logrus.Warn(fmt.Sprintf("Boot image architecture is unspecified and might not be compatible with %s %s nodes", icArch, role))
	} else if translateArchName[imageArch] != icArch {
		return fmt.Sprintf("image architecture %s does not match %s node architecture %s", imageArch, role, icArch)
	}
	return ""
}
