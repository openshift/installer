package powervs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

// Validate executes platform specific validation/
func Validate(config *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	if config.Platform.PowerVS == nil {
		allErrs = append(allErrs, field.Required(field.NewPath("platform", "powervs"), "Power VS Validation requires a Power VS platform configuration."))
	} else {
		if config.ControlPlane != nil {
			fldPath := field.NewPath("controlPlane")
			allErrs = append(allErrs, validateMachinePool(fldPath, config.ControlPlane)...)
		}
		for idx, compute := range config.Compute {
			fldPath := field.NewPath("compute").Index(idx)
			allErrs = append(allErrs, validateMachinePool(fldPath, &compute)...)
		}
		// Machine pool CIDR check
		for i := range config.Networking.MachineNetwork {
			// Each machine pool CIDR must have 24 significant bits (/24)
			if bits, _ := config.Networking.MachineNetwork[i].CIDR.Mask.Size(); bits != 24 {
				// If not, create an error displaying the CIDR in the install config vs the expectation (/24)
				fldPath := field.NewPath("Networking")
				allErrs = append(allErrs, field.Invalid(fldPath.Child("MachineNetwork").Child("CIDR"), (&config.Networking.MachineNetwork[i].CIDR).String(), "Machine Pool CIDR must be /24."))
			}
		}
	}
	return allErrs.ToAggregate()
}

func validateMachinePool(fldPath *field.Path, machinePool *types.MachinePool) field.ErrorList {
	allErrs := field.ErrorList{}
	if machinePool.Architecture != "ppc64le" {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("architecture"), machinePool.Architecture, []string{"ppc64le"}))
	}
	return allErrs
}

// ValidatePERAvailability ensures the target datacenter has PER enabled.
func ValidatePERAvailability(client API, ic *types.InstallConfig) error {
	capabilities, err := client.GetDatacenterCapabilities(context.TODO(), ic.PowerVS.Zone)
	if err != nil {
		return err
	}
	const per = "power-edge-router"
	perAvail, ok := capabilities[per]
	if !ok {
		return fmt.Errorf("%s capability unknown at: %s", per, ic.PowerVS.Zone)
	}
	if !perAvail {
		return fmt.Errorf("%s is not available at: %s", per, ic.PowerVS.Zone)
	}

	return nil
}

// ValidatePreExistingDNS ensures no pre-existing DNS record exists in the CIS
// DNS zone or IBM DNS zone for cluster's Kubernetes API.
func ValidatePreExistingDNS(client API, ic *types.InstallConfig, metadata MetadataAPI) error {
	allErrs := field.ErrorList{}

	fldPath := field.NewPath("baseDomain")
	if ic.Publish == types.ExternalPublishingStrategy {
		allErrs = append(allErrs, validatePreExistingPublicDNS(fldPath, client, ic, metadata)...)
	} else {
		allErrs = append(allErrs, validatePreExistingPrivateDNS(fldPath, client, ic, metadata)...)
	}

	return allErrs.ToAggregate()
}

func validatePreExistingPublicDNS(fldPath *field.Path, client API, ic *types.InstallConfig, metadata MetadataAPI) field.ErrorList {
	allErrs := field.ErrorList{}
	// Get CIS CRN
	crn, err := metadata.CISInstanceCRN(context.TODO())
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, err))
	}

	// Get CIS zone ID by name
	zoneID, err := client.GetDNSZoneIDByName(context.TODO(), ic.BaseDomain, types.ExternalPublishingStrategy)
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, err))
	}

	// Search for existing records
	recordNames := [...]string{fmt.Sprintf("api.%s", ic.ClusterDomain()), fmt.Sprintf("api-int.%s", ic.ClusterDomain())}
	for _, recordName := range recordNames {
		records, err := client.GetDNSRecordsByName(context.TODO(), crn, zoneID, recordName, types.ExternalPublishingStrategy)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath, err))
		}

		// DNS record exists
		if len(records) != 0 {
			allErrs = append(allErrs, field.Duplicate(fldPath, fmt.Sprintf("record %s already exists in CIS zone (%s) and might be in use by another cluster, please remove it to continue", recordName, zoneID)))
		}
	}
	return allErrs
}

func validatePreExistingPrivateDNS(fldPath *field.Path, client API, ic *types.InstallConfig, metadata MetadataAPI) field.ErrorList {
	allErrs := field.ErrorList{}
	// Get DNS CRN
	crn, err := metadata.DNSInstanceCRN(context.TODO())
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, err))
	}

	// Get CIS zone ID by name
	zoneID, err := client.GetDNSZoneIDByName(context.TODO(), ic.BaseDomain, types.InternalPublishingStrategy)
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, err))
	}

	// Search for existing records
	recordNames := [...]string{fmt.Sprintf("api-int.%s", ic.ClusterDomain())}
	for _, recordName := range recordNames {
		records, err := client.GetDNSRecordsByName(context.TODO(), crn, zoneID, recordName, types.InternalPublishingStrategy)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath, err))
		}

		// DNS record exists
		if len(records) != 0 {
			allErrs = append(allErrs, field.Duplicate(fldPath, fmt.Sprintf("record %s already exists in DNS zone (%s) and might be in use by another cluster, please remove it to continue", recordName, zoneID)))
		}
	}
	return allErrs
}

// ValidateCustomVPCSetup ensures optional VPC settings, if specified, are all legit.
func ValidateCustomVPCSetup(client API, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	var vpcRegion = ic.PowerVS.VPCRegion
	var vpcName = ic.PowerVS.VPCName
	var err error
	fldPath := field.NewPath("VPC")

	if vpcRegion != "" {
		if !powervstypes.ValidateVPCRegion(vpcRegion) {
			allErrs = append(allErrs, field.NotFound(fldPath.Child("vpcRegion"), vpcRegion))
		}
	} else {
		vpcRegion, err = powervstypes.VPCRegionForPowerVSRegion(ic.PowerVS.Region)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("region"), nil, ic.PowerVS.Region))
		}
	}

	if vpcName != "" {
		allErrs = append(allErrs, findVPCInRegion(client, vpcName, vpcRegion, fldPath)...)
		allErrs = append(allErrs, findSubnetInVPC(client, ic.PowerVS.VPCSubnets, vpcRegion, vpcName, fldPath)...)
	} else if len(ic.PowerVS.VPCSubnets) != 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("vpcSubnets"), nil, "invalid without vpcName"))
	}

	return allErrs.ToAggregate()
}

func findVPCInRegion(client API, name string, region string, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if name == "" {
		return allErrs
	}

	vpcs, err := client.GetVPCs(context.TODO(), region)
	if err != nil {
		return append(allErrs, field.InternalError(path.Child("vpcRegion"), err))
	}

	found := false
	for _, vpc := range vpcs {
		if *vpc.Name == name {
			found = true
			break
		}
	}
	if !found {
		allErrs = append(allErrs, field.NotFound(path.Child("vpcName"), name))
	}

	return allErrs
}

func findSubnetInVPC(client API, subnets []string, region string, name string, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(subnets) == 0 {
		return allErrs
	}

	subnet, err := client.GetSubnetByName(context.TODO(), subnets[0], region)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(path.Child("vpcSubnets"), err))
	} else if *subnet.VPC.Name != name {
		allErrs = append(allErrs, field.Invalid(path.Child("vpcSubnets"), nil, "not attached to VPC"))
	}

	return allErrs
}

// ValidateResourceGroup validates the resource group in our install config.
func ValidateResourceGroup(client API, ic *types.InstallConfig) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()

	resourceGroups, err := client.ListResourceGroups(ctx)
	if err != nil {
		return fmt.Errorf("failed to list resourceGroups: %w", err)
	}

	switch ic.PowerVS.PowerVSResourceGroup {
	case "":
		return errors.New("platform:powervs:powervsresourcegroup is empty")
	case "Default":
		found := false
		for _, resourceGroup := range resourceGroups.Resources {
			if resourceGroup.Default != nil && *resourceGroup.Default {
				found = true
				ic.PowerVS.PowerVSResourceGroup = *resourceGroup.Name
				break
			}
		}
		if !found {
			return errors.New("platform:powervs:powervsresourcegroup is default but no default exists")
		}
	default:
		found := false
		for _, resourceGroup := range resourceGroups.Resources {
			if *resourceGroup.Name == ic.PowerVS.PowerVSResourceGroup {
				found = true
				break
			}
		}
		if !found {
			return errors.New("platform:powervs:powervsresourcegroup has an invalid name")
		}
	}

	return nil
}

// ValidateSystemTypeForZone checks if the specified sysType is available in the target zone.
func ValidateSystemTypeForZone(client API, ic *types.InstallConfig) error {
	var (
		availableOnes []string
		err           error
	)

	if ic.ControlPlane == nil || ic.ControlPlane.Platform.PowerVS == nil || ic.ControlPlane.Platform.PowerVS.SysType == "" {
		return nil
	}
	availableOnes, err = client.GetDatacenterSupportedSystems(context.Background(), ic.PowerVS.Zone)
	if err != nil {
		// Fallback to hardcoded list
		availableOnes, err = powervstypes.AvailableSysTypes(ic.PowerVS.Region, ic.PowerVS.Zone)
		if err != nil {
			return fmt.Errorf("failed to obtain available SysTypes for: %s", ic.PowerVS.Zone)
		}
	}
	requested := ic.ControlPlane.Platform.PowerVS.SysType
	found := false
	for i := range availableOnes {
		if requested == availableOnes[i] {
			found = true
			break
		}
	}
	if found {
		return nil
	}
	return fmt.Errorf("%s is not available in: %s, these are %v", requested, ic.PowerVS.Zone, availableOnes)
}

// ValidateServiceInstance validates the optional service instance GUID in our install config.
func ValidateServiceInstance(client API, ic *types.InstallConfig) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()

	serviceInstances, err := client.ListServiceInstances(ctx)
	if err != nil {
		return err
	}

	switch ic.PowerVS.ServiceInstanceGUID {
	case "":
		return nil
	default:
		found := false
		for _, serviceInstance := range serviceInstances {
			guid := strings.SplitN(serviceInstance, " ", 2)[1]
			if guid == ic.PowerVS.ServiceInstanceGUID {
				found = true
				break
			}
		}
		if !found {
			return errors.New("platform:powervs:serviceInstanceGUID has an invalid guid")
		}
	}

	return nil
}

// ValidateTransitGateway validates the optional transit gateway name in our install config.
func ValidateTransitGateway(client API, ic *types.InstallConfig) error {
	var (
		id  string
		err error
	)

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()

	if len(ic.PowerVS.TGName) > 0 {
		id, err = client.TransitGatewayID(ctx, ic.PowerVS.TGName)
		if err != nil {
			return err
		}
		if id == "" {
			return errors.New("platform:powervs:tgName has an invalid name")
		}
	}

	return nil
}
