package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	dnstypes "github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/validate"
	mapiutil "github.com/openshift/machine-api-provider-gcp/pkg/cloud/gcp/actuators/util"
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

const unknownArchitecture = ""

// Validate executes platform-specific validation.
func Validate(client API, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	if err := validate.GCPClusterName(ic.ObjectMeta.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("clusterName"), ic.ObjectMeta.Name, err.Error()))
	}

	allErrs = append(allErrs, validateProject(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateNetworkProject(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateRegion(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateZones(client, ic)...)
	allErrs = append(allErrs, validateNetworks(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateServiceEndpoints(client, ic, field.NewPath("platform").Child("gcp"))...)
	allErrs = append(allErrs, validateInstanceTypes(client, ic)...)
	allErrs = append(allErrs, ValidateCredentialMode(client, ic)...)
	allErrs = append(allErrs, validatePreexistingServiceAccount(client, ic)...)
	allErrs = append(allErrs, validateServiceAccountPresent(client, ic)...)
	allErrs = append(allErrs, validateMarketplaceImages(client, ic)...)
	allErrs = append(allErrs, validatePlatformKMSKeys(client, ic, field.NewPath("platform").Child("gcp"))...)

	if err := validateUserTags(client, ic.Platform.GCP.ProjectID, ic.Platform.GCP.UserTags); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("platform").Child("gcp").Child("userTags"), ic.Platform.GCP.UserTags, err.Error()))
	}

	return allErrs.ToAggregate()
}

func validateInstanceAndDiskType(fldPath *field.Path, diskType, instanceType, arch string) *field.Error {
	if instanceType == "" {
		// nothing to validate
		return nil
	}

	family, _, _ := strings.Cut(instanceType, "-")
	if family == "custom" {
		family = gcp.DefaultCustomInstanceType
	}
	diskTypes, ok := gcp.InstanceTypeToDiskTypeMap[family]
	if !ok {
		return field.NotFound(fldPath.Child("type"), family)
	}

	acceptedArmFamilies := sets.New("c4a", "t2a")
	if arch == types.ArchitectureARM64 && !acceptedArmFamilies.Has(family) {
		return field.NotSupported(fldPath.Child("type"), family, sets.List(acceptedArmFamilies))
	}

	if diskType != "" {
		if !sets.New(diskTypes...).Has(diskType) {
			return field.Invalid(
				fldPath.Child("diskType"),
				diskType,
				fmt.Sprintf("%s instance requires one of the following disk types: %v", instanceType, diskTypes),
			)
		}
	}
	return nil
}

// ValidateInstanceType ensures the instance type has sufficient Vcpu and Memory.
func ValidateInstanceType(client API, fieldPath *field.Path, project, region string, zones []string, diskType string, instanceType string, req resourceRequirements, arch string) field.ErrorList {
	allErrs := field.ErrorList{}

	typeMeta, typeZones, err := client.GetMachineTypeWithZones(context.TODO(), project, region, instanceType)
	if err != nil {
		if _, ok := err.(*googleapi.Error); ok {
			return append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, err.Error()))
		}
		return append(allErrs, field.InternalError(nil, err))
	}

	if fieldErr := validateInstanceAndDiskType(fieldPath, diskType, instanceType, arch); fieldErr != nil {
		return append(allErrs, fieldErr)
	}

	userZones := sets.New(zones...)
	if len(userZones) == 0 {
		userZones = typeZones
	}
	if diff := userZones.Difference(typeZones); len(diff) > 0 {
		errMsg := fmt.Sprintf("instance type not available in zones: %v", sets.List(diff))
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, errMsg))
	}

	if typeMeta.GuestCpus < req.minimumVCpus {
		errMsg := fmt.Sprintf("instance type does not meet minimum resource requirements of %d vCPUs", req.minimumVCpus)
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, errMsg))
	}
	if typeMeta.MemoryMb < req.minimumMemory {
		errMsg := fmt.Sprintf("instance type does not meet minimum resource requirements of %d MB Memory", req.minimumMemory)
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, errMsg))
	}

	if arch != unknownArchitecture {
		if typeArch := mapiutil.CPUArchitecture(instanceType); string(typeArch) != arch {
			errMsg := fmt.Sprintf("instance type architecture %s does not match specified architecture %s", typeArch, arch)
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, errMsg))
		}
	}

	return allErrs
}

func validateServiceAccountPresent(client API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.GCP.NetworkProjectID != "" {
		creds := client.GetCredentials()
		if creds != nil && creds.JSON == nil {
			if ic.ControlPlane.Platform.GCP != nil && ic.ControlPlane.Platform.GCP.ServiceAccount == "" {
				errMsg := "service account must be provided when authentication credentials do not provide a service account"
				allErrs = append(allErrs, field.Required(field.NewPath("controlPlane").Child("platform").Child("gcp").Child("serviceAccount"), errMsg))
			}
		}
	}
	return allErrs
}

// DefaultInstanceTypeForArch returns the appropriate instance type based on the target architecture.
func DefaultInstanceTypeForArch(arch types.Architecture) string {
	if arch == types.ArchitectureARM64 {
		return "t2a-standard-4"
	}
	return "n2-standard-4"
}

// validateInstanceTypes checks that the user-provided instance types are valid.
func validateInstanceTypes(client API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	defaultInstanceType := ""
	defaultDiskType := gcp.PDSSD
	defaultZones := []string{}

	// Default requirements need to be sufficient to support Control Plane instances.
	defaultInstanceReq := controlPlaneReq
	if ic.ControlPlane != nil && ic.ControlPlane.Platform.GCP != nil && ic.ControlPlane.Platform.GCP.InstanceType != "" {
		// Default requirements can be relaxed when the controlPlane type is set explicitly.
		defaultInstanceReq = computeReq
	}

	if ic.GCP.DefaultMachinePlatform != nil {
		defaultZones = ic.GCP.DefaultMachinePlatform.Zones
		defaultInstanceType = ic.GCP.DefaultMachinePlatform.InstanceType
		if ic.GCP.DefaultMachinePlatform.DiskType != "" {
			defaultDiskType = ic.GCP.DefaultMachinePlatform.DiskType
		}

		if ic.GCP.DefaultMachinePlatform.InstanceType != "" {
			allErrs = append(allErrs,
				ValidateInstanceType(
					client,
					field.NewPath("platform", "gcp", "defaultMachinePlatform"),
					ic.GCP.ProjectID,
					ic.GCP.Region,
					ic.GCP.DefaultMachinePlatform.Zones,
					defaultDiskType,
					ic.GCP.DefaultMachinePlatform.InstanceType,
					defaultInstanceReq,
					unknownArchitecture,
				)...)
		}
	}

	zones := defaultZones
	instanceType := defaultInstanceType
	arch := types.ArchitectureAMD64
	cpDiskType := defaultDiskType
	if ic.ControlPlane != nil {
		arch = string(ic.ControlPlane.Architecture)
		if instanceType == "" {
			instanceType = DefaultInstanceTypeForArch(ic.ControlPlane.Architecture)
		}
		if ic.ControlPlane.Platform.GCP != nil {
			if ic.ControlPlane.Platform.GCP.InstanceType != "" {
				instanceType = ic.ControlPlane.Platform.GCP.InstanceType
			}
			if len(ic.ControlPlane.Platform.GCP.Zones) > 0 {
				zones = ic.ControlPlane.Platform.GCP.Zones
			}
			if ic.ControlPlane.Platform.GCP.DiskType != "" {
				cpDiskType = ic.ControlPlane.Platform.GCP.DiskType
			}
		}
	}

	// The IOPS minimum Control plane requirements are not met for pd-standard machines.
	if cpDiskType == "pd-standard" {
		allErrs = append(allErrs,
			field.NotSupported(field.NewPath("controlPlane", "type"),
				cpDiskType,
				sets.List(gcp.ControlPlaneSupportedDisks)),
		)
	} else {
		allErrs = append(allErrs,
			ValidateInstanceType(
				client,
				field.NewPath("controlPlane", "platform", "gcp"),
				ic.GCP.ProjectID,
				ic.GCP.Region,
				zones,
				cpDiskType,
				instanceType,
				controlPlaneReq,
				arch,
			)...)
	}

	for idx, compute := range ic.Compute {
		fieldPath := field.NewPath("compute").Index(idx)
		zones := defaultZones
		instanceType := defaultInstanceType
		diskType := defaultDiskType
		if instanceType == "" {
			instanceType = DefaultInstanceTypeForArch(compute.Architecture)
		}
		if diskType == "" {
			diskType = gcp.PDSSD
		}
		arch := compute.Architecture
		if compute.Platform.GCP != nil {
			if compute.Platform.GCP.InstanceType != "" {
				instanceType = compute.Platform.GCP.InstanceType
			}
			if len(compute.Platform.GCP.Zones) > 0 {
				zones = compute.Platform.GCP.Zones
			}
		}

		if compute.Platform.GCP != nil && compute.Platform.GCP.DiskType != "" {
			diskType = compute.Platform.GCP.DiskType
		}

		allErrs = append(allErrs,
			ValidateInstanceType(
				client,
				fieldPath.Child("platform", "gcp"),
				ic.GCP.ProjectID,
				ic.GCP.Region,
				zones,
				diskType,
				instanceType,
				computeReq,
				string(arch),
			)...)
	}

	return allErrs
}

func validatePreexistingServiceAccount(client API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.ControlPlane.Platform.GCP != nil && ic.ControlPlane.Platform.GCP.ServiceAccount != "" {
		fldPath := field.NewPath("controlPlane").Child("platform").Child("gcp").Child("serviceAccount")

		// The service account is required for resources in the host project.
		serviceAccount, err := client.GetServiceAccount(context.Background(), ic.GCP.ProjectID, ic.ControlPlane.Platform.GCP.ServiceAccount)
		if err != nil {
			return append(allErrs, field.InternalError(fldPath, err))
		}
		if serviceAccount == "" {
			return append(allErrs, field.NotFound(fldPath, ic.ControlPlane.Platform.GCP.ServiceAccount))
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
	if ic.Platform.GCP.UserProvisionedDNS == dnstypes.UserProvisionedDNSEnabled {
		return nil
	}

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

func validateServiceEndpoints(_ API, ic *types.InstallConfig, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// attempt to resolve all the custom (overridden) endpoints. If any are not reachable,
	// then the installation should fail not skip the endpoint use.
	for id, serviceEndpoint := range ic.GCP.ServiceEndpoints {
		if _, err := url.Parse(serviceEndpoint.URL); err != nil {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("serviceEndpoints").Index(id), serviceEndpoint.URL, fmt.Sprintf("failed to parse service endpoint url: %v", err)))
		} else if _, err := http.Head(serviceEndpoint.URL); err != nil {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("serviceEndpoints").Index(id), serviceEndpoint.URL, fmt.Sprintf("error connecting to endpoint: %v", err)))
		}
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
		"storage-component.googleapis.com",
		"file.googleapis.com")
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

// ValidateCredentialMode The presence of `authorized_user` in the credentials indicates that no service account
// was used for authentication and requires Manual credential mode.
func ValidateCredentialMode(client API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	creds := client.GetCredentials()

	if creds.JSON != nil {
		var credsMap map[string]interface{}
		err := json.Unmarshal(creds.JSON, &credsMap)
		if err != nil {
			return append(allErrs, field.Invalid(field.NewPath("credentials").Child("JSON"), creds.JSON, "failed to unmarshal JSON credentials"))
		}

		credsType, found := credsMap["type"]
		if !found {
			return append(allErrs, field.NotFound(field.NewPath("credentials").Child("JSON").Child("type"), "failed to find credentials type"))
		}

		if credsType.(string) == string(gcp.AuthorizedUserMode) && ic.CredentialsMode != types.ManualCredentialsMode {
			errMsg := "environmental authentication is only supported with Manual credentials mode"
			return append(allErrs, field.Forbidden(field.NewPath("credentialsMode"), errMsg))
		}
	} else if creds.JSON == nil && ic.CredentialsMode != types.ManualCredentialsMode {
		errMsg := "Manual credentials mode needs to be enabled to use environmental authentication"
		return append(allErrs, field.Forbidden(field.NewPath("credentialsMode"), errMsg))
	}
	return allErrs
}

func validateZones(client API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	zones, err := client.GetZones(context.TODO(), ic.GCP.ProjectID, ic.GCP.Region)
	if err != nil {
		return append(allErrs, field.InternalError(nil, err))
	} else if len(zones) == 0 {
		return append(allErrs, field.InternalError(nil, fmt.Errorf("failed to fetch zones, this error usually occurs if the region is not found")))
	}

	projZones := sets.New[string]()
	for _, zone := range zones {
		projZones.Insert(zone.Name)
	}

	const errMsg = "zone(s) not found in region"

	if ic.Platform.GCP.DefaultMachinePlatform != nil {
		diff := sets.New(ic.Platform.GCP.DefaultMachinePlatform.Zones...).Difference(projZones)
		if len(diff) > 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("platform", "gcp", "defaultMachinePlatform", "zones"), sets.List(diff), errMsg))
		}
	}

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.GCP != nil {
		diff := sets.New(ic.ControlPlane.Platform.GCP.Zones...).Difference(projZones)
		if len(diff) > 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("controlPlane", "platform", "gcp", "zones"), sets.List(diff), errMsg))
		}
	}

	for idx, compute := range ic.Compute {
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.GCP != nil {
			diff := sets.New(compute.Platform.GCP.Zones...).Difference(projZones)
			if len(diff) > 0 {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("platform", "gcp", "zones"), sets.List(diff), errMsg))
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
	var (
		translateArchName = map[string]types.Architecture{
			"ARM64":  types.ArchitectureARM64,
			"X86_64": types.ArchitectureAMD64,
		}
	)

	if imageArch == "" || imageArch == unspecifiedArch {
		logrus.Warn(fmt.Sprintf("Boot image architecture is unspecified and might not be compatible with %s %s nodes", icArch, role))
	} else if translateArchName[imageArch] != icArch {
		return fmt.Sprintf("image architecture %s does not match %s node architecture %s", imageArch, role, icArch)
	}
	return ""
}

// validateUserTags check for existence and accessibility of user-defined tags and persists
// validated tags in-memory.
func validateUserTags(client API, projectID string, userTags []gcp.UserTag) error {
	return NewTagManager(client).validateAndPersistUserTags(context.Background(), projectID, userTags)
}

// validatePlatformKMSKeys checks for encryption keys for all the machine pools. The encryption key rings are
// checked against the API for validity/availability.
func validatePlatformKMSKeys(client API, ic *types.InstallConfig, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	cp := ic.ControlPlane
	validatedControlPlaneKey := false
	if cp != nil && cp.Platform.GCP != nil && cp.Platform.GCP.EncryptionKey != nil && cp.Platform.GCP.EncryptionKey.KMSKey != nil {
		if _, err := client.GetKeyRing(context.TODO(), cp.Platform.GCP.OSDisk.EncryptionKey.KMSKey.KeyRing); err != nil {
			return append(allErrs, field.Invalid(fieldPath.Child("controlPlane").Child("encryptionKey").Child("kmsKey").Child("keyRing"),
				cp.Platform.GCP.OSDisk.EncryptionKey.KMSKey.KeyRing,
				err.Error(),
			))
		}
		validatedControlPlaneKey = true
	}

	validatedComputeKeys := false
	for _, mp := range ic.Compute {
		if mp.Platform.GCP != nil && mp.Platform.GCP.EncryptionKey != nil && mp.Platform.GCP.EncryptionKey.KMSKey != nil {
			if _, err := client.GetKeyRing(context.TODO(), mp.Platform.GCP.OSDisk.EncryptionKey.KMSKey.KeyRing); err != nil {
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("compute").Child("encryptionKey").Child("kmsKey").Child("keyRing"),
					mp.Platform.GCP.OSDisk.EncryptionKey.KMSKey.KeyRing,
					err.Error(),
				))
			} else {
				validatedComputeKeys = true
			}
		}
	}

	defaultMp := ic.GCP.DefaultMachinePlatform
	if defaultMp != nil && defaultMp.EncryptionKey != nil && defaultMp.EncryptionKey.KMSKey != nil {
		if _, err := client.GetKeyRing(context.TODO(), defaultMp.EncryptionKey.KMSKey.KeyRing); err != nil {
			if validatedControlPlaneKey && (validatedComputeKeys && len(allErrs) == 0) {
				logrus.Warn("defaultMachinePool.encryptionKey.KMSKey.KeyRing is not valid, but compute and control plane key rings are valid")
			} else {
				return append(allErrs, field.Invalid(fieldPath.Child("defaultMachinePool").Child("encryptionKey").Child("kmsKey").Child("keyRing"),
					defaultMp.EncryptionKey.KMSKey.KeyRing,
					err.Error(),
				))
			}
		}
	}

	return allErrs
}
