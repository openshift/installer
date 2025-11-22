package aws

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"slices"
	"sort"

	ec2v2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/dns"
)

type resourceRequirements struct {
	minimumVCpus  int64
	minimumMemory int64
}

var controlPlaneReq = resourceRequirements{
	minimumVCpus:  4,
	minimumMemory: 16384,
}

var computeReq = resourceRequirements{
	minimumVCpus:  2,
	minimumMemory: 8192,
}

// Validate executes platform-specific validation.
func Validate(ctx context.Context, meta *Metadata, config *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	if config.Platform.AWS == nil {
		return errors.New(field.Required(field.NewPath("platform", "aws"), "AWS validation requires an AWS platform configuration").Error())
	}

	allErrs = append(allErrs, validateAMI(ctx, meta, config)...)
	allErrs = append(allErrs, validatePublicIpv4Pool(ctx, meta, field.NewPath("platform", "aws", "publicIpv4PoolId"), config)...)
	allErrs = append(allErrs, validatePlatform(ctx, meta, field.NewPath("platform", "aws"), config)...)

	if awstypes.IsPublicOnlySubnetsEnabled() {
		logrus.Warnln("Public-only subnets install. Please be warned this is not supported")
		if config.Publish == types.InternalPublishingStrategy {
			allErrs = append(allErrs, field.Invalid(field.NewPath("publish"), config.Publish, "cluster cannot be private with public subnets"))
		}
	}

	if config.ControlPlane != nil {
		arch := string(config.ControlPlane.Architecture)
		pool := &awstypes.MachinePool{}
		pool.Set(config.AWS.DefaultMachinePlatform)
		pool.Set(config.ControlPlane.Platform.AWS)
		allErrs = append(allErrs, validateMachinePool(ctx, meta, field.NewPath("controlPlane", "platform", "aws"), config.Platform.AWS, pool, controlPlaneReq, "", arch)...)
	}

	var archSeen string
	for idx, compute := range config.Compute {
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Name == types.MachinePoolEdgeRoleName {
			if len(config.Platform.AWS.VPC.Subnets) == 0 {
				if compute.Platform.AWS == nil {
					allErrs = append(allErrs, field.Required(fldPath.Child("platform", "aws"), "edge compute pools are only supported on the AWS platform"))
				}
			}
		}

		arch := string(compute.Architecture)
		if arch == "" {
			arch = string(config.ControlPlane.Architecture)
		}
		switch {
		case archSeen == "":
			archSeen = arch
		case arch != archSeen:
			allErrs = append(allErrs, field.Invalid(fldPath.Child("architecture"), arch, "all compute machine pools must be of the same architecture"))
		default:
			// compute machine pools have the same arch so far
		}
		pool := &awstypes.MachinePool{}
		pool.Set(config.AWS.DefaultMachinePlatform)
		pool.Set(compute.Platform.AWS)
		allErrs = append(allErrs, validateMachinePool(ctx, meta, fldPath.Child("platform", "aws"), config.Platform.AWS, pool, computeReq, compute.Name, arch)...)
	}
	return allErrs.ToAggregate()
}

func validatePlatform(ctx context.Context, meta *Metadata, fldPath *field.Path, config *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	platform := config.Platform.AWS

	allErrs = append(allErrs, validateServiceEndpoints(fldPath.Child("serviceEndpoints"), platform.Region, platform.ServiceEndpoints)...)

	// Fail fast when service endpoints are invalid to avoid long timeouts.
	if len(allErrs) > 0 {
		return allErrs
	}

	if len(platform.VPC.Subnets) > 0 {
		allErrs = append(allErrs, validateSubnets(ctx, meta, fldPath.Child("vpc").Child("subnets"), config)...)
		allErrs = append(allErrs, validateSharedVPC(ctx, meta, fldPath.Child("vpc").Child("subnets"))...)
	}
	if platform.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, validateMachinePool(ctx, meta, fldPath.Child("defaultMachinePlatform"), platform, platform.DefaultMachinePlatform, controlPlaneReq, "", "")...)
	}
	return allErrs
}

func validateAMI(ctx context.Context, meta *Metadata, config *types.InstallConfig) field.ErrorList {
	// accept AMI from the rhcos stream metadata
	if rhcos.AMIRegions(config.ControlPlane.Architecture).Has(config.Platform.AWS.Region) {
		return nil
	}

	// accept AMI specified at the platform level
	if config.Platform.AWS.AMIID != "" {
		return nil
	}

	// accept AMI specified for the default machine platform
	if config.Platform.AWS.DefaultMachinePlatform != nil {
		if config.Platform.AWS.DefaultMachinePlatform.AMIID != "" {
			return nil
		}
	}

	// accept AMIs specified specifically for each machine pool
	controlPlaneHasAMISpecified := false
	if config.ControlPlane != nil && config.ControlPlane.Platform.AWS != nil {
		controlPlaneHasAMISpecified = config.ControlPlane.Platform.AWS.AMIID != ""
	}

	computesHaveAMISpecified := true
	for _, c := range config.Compute {
		if c.Replicas != nil && *c.Replicas == 0 {
			continue
		}
		if c.Platform.AWS == nil || c.Platform.AWS.AMIID == "" {
			computesHaveAMISpecified = false
		}
	}

	if controlPlaneHasAMISpecified && computesHaveAMISpecified {
		return nil
	}

	// accept AMI that can be copied from us-east-1 if the region is in the standard AWS partition
	regions, err := meta.Regions(ctx)
	if err != nil {
		return field.ErrorList{field.InternalError(field.NewPath("platform", "aws", "region"), fmt.Errorf("failed to get list of regions: %w", err))}
	}
	if sets.New(regions...).Has(config.Platform.AWS.Region) {
		defaultEndpoint, err := GetDefaultServiceEndpoint(ctx, ec2v2.ServiceID, EndpointOptions{Region: config.Platform.AWS.Region, UseFIPS: false})
		if err != nil {
			return field.ErrorList{field.InternalError(field.NewPath("platform", "aws", "region"), fmt.Errorf("failed to resolve ec2 endpoint"))}
		}
		if defaultEndpoint.PartitionID == endpoints.AwsPartitionID {
			return nil
		}
	}

	// fail validation since we do not have an AMI to use
	return field.ErrorList{field.Required(field.NewPath("platform", "aws", "amiID"), "AMI must be provided")}
}

func validatePublicIpv4Pool(ctx context.Context, meta *Metadata, fldPath *field.Path, config *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	if config.Platform.AWS.PublicIpv4Pool == "" {
		return nil
	}
	poolID := config.Platform.AWS.PublicIpv4Pool
	if config.Publish != types.ExternalPublishingStrategy {
		return append(allErrs, field.Invalid(fldPath, poolID, fmt.Errorf("publish strategy %s can't be used with custom Public IPv4 Pools", config.Publish).Error()))
	}

	// Pool validations
	// Resources claiming Public IPv4 from Pool in regular 'External' installations:
	// 1* for Bootsrtap
	// N*Zones for NAT Gateways
	// N*Zones for API LB
	// N*Zones for Ingress LB
	allzones, err := meta.AvailabilityZones(ctx)
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, err))
	}
	totalPublicIPRequired := int64(1 + (len(allzones) * 3))

	sess, err := meta.Session(ctx)
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, fmt.Errorf("unable to retrieve aws session: %w", err)))
	}

	publicIpv4Pool, err := DescribePublicIpv4Pool(ctx, sess, config.Platform.AWS.Region, poolID)
	if err != nil {
		return append(allErrs, field.Invalid(fldPath, poolID, err.Error()))
	}

	got := aws.Int64Value(publicIpv4Pool.TotalAvailableAddressCount)
	if got < totalPublicIPRequired {
		err = fmt.Errorf("required a minimum of %d Public IPv4 IPs available in the pool %s, got %d", totalPublicIPRequired, poolID, got)
		return append(allErrs, field.InternalError(fldPath, err))
	}

	return nil
}

// subnetData holds a subnet information collected from install config and AWS API for validations.
type subnetData struct {
	// The subnet index in the install config.
	Idx int
	// The subnet assigned roles in the install config.
	Roles []awstypes.SubnetRole
	// The subnet metadata from AWS API.
	Subnet
}

// subnetDataGroups is a collection of subnet information
// grouped by subnet type (i.e. public, private, and edge) and indexed by subnetIDs for validations.
type subnetDataGroups struct {
	Public  map[string]subnetData
	Private map[string]subnetData
	Edge    map[string]subnetData
	// A convenient alias that contains all information for all subnets.
	All map[string]subnetData
}

// Converts subnetGroups (i.e. provided subnets) to subnetDataGroups to include additional information
// from the install-config such as index and roles for validations.
func (sdg *subnetDataGroups) From(ctx context.Context, meta *Metadata, providedSubnets []awstypes.Subnet) error {
	if sdg.Private == nil {
		sdg.Private = make(map[string]subnetData)
	}
	if sdg.Public == nil {
		sdg.Public = make(map[string]subnetData)
	}
	if sdg.Edge == nil {
		sdg.Edge = make(map[string]subnetData)
	}
	if sdg.All == nil {
		sdg.All = make(map[string]subnetData)
	}

	subnets, err := meta.Subnets(ctx)
	if err != nil {
		return err
	}

	for idx, subnet := range providedSubnets {
		var subnetDataGroup map[string]subnetData
		var subnetMeta Subnet

		if awsSubnet, ok := subnets.Private[string(subnet.ID)]; ok {
			subnetDataGroup = sdg.Private
			subnetMeta = awsSubnet
		}

		if awsSubnet, ok := subnets.Public[string(subnet.ID)]; ok {
			subnetDataGroup = sdg.Public
			subnetMeta = awsSubnet
		}

		if awsSubnet, ok := subnets.Edge[string(subnet.ID)]; ok {
			subnetDataGroup = sdg.Edge
			subnetMeta = awsSubnet
		}

		if subnetDataGroup == nil {
			// Should not occur but safe against panics
			continue
		}

		subnetData := subnetData{
			Subnet: subnetMeta,
			Idx:    idx,
			Roles:  subnet.Roles,
		}
		subnetDataGroup[string(subnet.ID)] = subnetData
		sdg.All[string(subnet.ID)] = subnetData
	}

	return nil
}

// validateSubnets ensures BYO subnets are valid.
func validateSubnets(ctx context.Context, meta *Metadata, fldPath *field.Path, config *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	networking := config.Networking
	providedSubnets := config.AWS.VPC.Subnets
	publish := config.Publish

	subnetDataGroups := subnetDataGroups{}
	if err := subnetDataGroups.From(ctx, meta, providedSubnets); err != nil {
		return append(allErrs, field.Invalid(fldPath, providedSubnets, err.Error()))
	}

	publicOnlySubnet := awstypes.IsPublicOnlySubnetsEnabled()

	if publicOnlySubnet && len(subnetDataGroups.Public) == 0 {
		allErrs = append(allErrs, field.Required(fldPath, "public subnets are required for a public-only subnets cluster"))
	}

	if !publicOnlySubnet && len(subnetDataGroups.Private) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath, providedSubnets, "no private subnets found"))
	}

	if publish == types.InternalPublishingStrategy && len(subnetDataGroups.Public) > 0 {
		logrus.Warnf("public subnets should not be provided when publish is set to %s", types.InternalPublishingStrategy)
	}

	subnetsWithRole := make(map[awstypes.SubnetRoleType][]subnetData)
	for _, subnet := range providedSubnets {
		for _, role := range subnet.Roles {
			subnetsWithRole[role.Type] = append(subnetsWithRole[role.Type], subnetDataGroups.All[string(subnet.ID)])
		}
	}

	allErrs = append(allErrs, validateSharedSubnets(ctx, meta, fldPath)...)
	allErrs = append(allErrs, validateSubnetCIDR(fldPath, subnetDataGroups.Private, networking.MachineNetwork)...)
	allErrs = append(allErrs, validateSubnetCIDR(fldPath, subnetDataGroups.Public, networking.MachineNetwork)...)

	if len(subnetsWithRole) > 0 {
		allErrs = append(allErrs, validateSubnetRoles(fldPath, subnetsWithRole, subnetDataGroups, config)...)
	} else {
		allErrs = append(allErrs, validateUntaggedSubnets(ctx, fldPath, meta, subnetDataGroups)...)
		allErrs = append(allErrs, validateDuplicateSubnetZones(fldPath, subnetDataGroups.Private, "private")...)
		allErrs = append(allErrs, validateDuplicateSubnetZones(fldPath, subnetDataGroups.Public, "public")...)
		allErrs = append(allErrs, validateDuplicateSubnetZones(fldPath, subnetDataGroups.Edge, "edge")...)
	}

	privateZones := sets.New[string]()
	publicZones := sets.New[string]()
	for _, subnet := range subnetDataGroups.Private {
		privateZones.Insert(subnet.Zone.Name)
	}
	for _, subnet := range subnetDataGroups.Public {
		publicZones.Insert(subnet.Zone.Name)
	}
	if publish == types.ExternalPublishingStrategy && !publicZones.IsSuperset(privateZones) {
		errMsg := fmt.Sprintf("No public subnet provided for zones %s", sets.List(privateZones.Difference(publicZones)))
		allErrs = append(allErrs, field.Invalid(fldPath, providedSubnets, errMsg))
	}

	return allErrs
}

func validateMachinePool(ctx context.Context, meta *Metadata, fldPath *field.Path, platform *awstypes.Platform, pool *awstypes.MachinePool, req resourceRequirements, poolName string, arch string) field.ErrorList {
	var err error
	allErrs := field.ErrorList{}

	// Pool's specific validation.
	// Edge Compute Pool / AWS Local Zones:
	// - is valid when installing in existing VPC; or
	// - is valid in new VPC when Local Zone name is defined
	if poolName == types.MachinePoolEdgeRoleName {
		if len(platform.VPC.Subnets) > 0 {
			edgeSubnets, err := meta.EdgeSubnets(ctx)
			if err != nil {
				errMsg := fmt.Sprintf("%s pool. %v", poolName, err.Error())
				return append(allErrs, field.Invalid(field.NewPath("platform", "aws", "vpc", "subnets"), platform.VPC.Subnets, errMsg))
			}
			if len(edgeSubnets) == 0 {
				return append(allErrs, field.Required(fldPath, "the provided subnets must include valid subnets for the specified edge zones"))
			}
		} else {
			if pool.Zones == nil || len(pool.Zones) == 0 {
				return append(allErrs, field.Required(fldPath, "zone is required when using edge machine pools"))
			}
			for _, zone := range pool.Zones {
				err := validateZoneLocal(ctx, meta, fldPath.Child("zones"), zone)
				if err != nil {
					allErrs = append(allErrs, err)
				}
			}
			if len(allErrs) > 0 {
				return allErrs
			}
		}
	}

	if pool.Zones != nil && len(pool.Zones) > 0 {
		availableZones := sets.New[string]()
		diffErrMsgPrefix := "One or more zones are unavailable"
		if len(platform.VPC.Subnets) > 0 {
			diffErrMsgPrefix = "No subnets provided for zones"
			var subnets Subnets
			if poolName == types.MachinePoolEdgeRoleName {
				subnets, err = meta.EdgeSubnets(ctx)
			} else {
				subnets, err = meta.PrivateSubnets(ctx)
			}

			if err != nil {
				return append(allErrs, field.InternalError(fldPath, err))
			}
			for _, subnet := range subnets {
				availableZones.Insert(subnet.Zone.Name)
			}
		} else {
			var allzones []string
			if poolName == types.MachinePoolEdgeRoleName {
				allzones, err = meta.EdgeZones(ctx)
			} else {
				allzones, err = meta.AvailabilityZones(ctx)
			}
			if err != nil {
				return append(allErrs, field.InternalError(fldPath, err))
			}
			availableZones.Insert(allzones...)
		}

		if diff := sets.New[string](pool.Zones...).Difference(availableZones); diff.Len() > 0 {
			errMsg := fmt.Sprintf("%s %s", diffErrMsgPrefix, sets.List(diff))
			allErrs = append(allErrs, field.Invalid(fldPath.Child("zones"), pool.Zones, errMsg))
		}
	}
	if pool.InstanceType != "" {
		instanceTypes, err := meta.InstanceTypes(ctx)
		if err != nil {
			return append(allErrs, field.InternalError(fldPath, err))
		}
		if typeMeta, ok := instanceTypes[pool.InstanceType]; ok {
			if typeMeta.DefaultVCpus < req.minimumVCpus {
				errMsg := fmt.Sprintf("instance type does not meet minimum resource requirements of %d vCPUs", req.minimumVCpus)
				allErrs = append(allErrs, field.Invalid(fldPath.Child("type"), pool.InstanceType, errMsg))
			}
			if typeMeta.MemInMiB < req.minimumMemory {
				errMsg := fmt.Sprintf("instance type does not meet minimum resource requirements of %d MiB Memory", req.minimumMemory)
				allErrs = append(allErrs, field.Invalid(fldPath.Child("type"), pool.InstanceType, errMsg))
			}
			instanceArches := translateEC2Arches(typeMeta.Arches)
			// `arch` might not be specified (e.g, defaultMachinePool)
			if len(arch) > 0 && !instanceArches.Has(arch) {
				errMsg := fmt.Sprintf("instance type supported architectures %s do not match specified architecture %s", sets.List(instanceArches), arch)
				allErrs = append(allErrs, field.Invalid(fldPath.Child("type"), pool.InstanceType, errMsg))
			}
		} else {
			errMsg := fmt.Sprintf("instance type %s not found", pool.InstanceType)
			allErrs = append(allErrs, field.Invalid(fldPath.Child("type"), pool.InstanceType, errMsg))
		}
	}

	if len(pool.AdditionalSecurityGroupIDs) > 0 {
		allErrs = append(allErrs, validateSecurityGroupIDs(ctx, meta, fldPath.Child("additionalSecurityGroupIDs"), platform, pool)...)
	}

	if len(pool.IAMProfile) > 0 {
		if len(pool.IAMRole) > 0 {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("iamRole"), "cannot be used with iamProfile"))
		}
		if err := validateInstanceProfile(ctx, meta, fldPath.Child("iamProfile"), pool); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	allErrs = append(allErrs, validateHostPlacement(ctx, meta, fldPath, pool)...)

	return allErrs
}

func translateEC2Arches(arches []string) sets.Set[string] {
	res := sets.New[string]()
	for _, arch := range arches {
		switch arch {
		case ec2.ArchitectureTypeX8664:
			res.Insert(types.ArchitectureAMD64)
		case ec2.ArchitectureTypeArm64:
			res.Insert(types.ArchitectureARM64)
		default:
			continue
		}
	}
	return res
}

func validateHostPlacement(ctx context.Context, meta *Metadata, fldPath *field.Path, pool *awstypes.MachinePool) field.ErrorList {
	allErrs := field.ErrorList{}

	if pool.HostPlacement == nil {
		return allErrs
	}

	if pool.HostPlacement.Affinity != nil && *pool.HostPlacement.Affinity == awstypes.HostAffinityDedicatedHost {
		placementPath := fldPath.Child("hostPlacement")
		if pool.HostPlacement.DedicatedHost != nil {
			configuredHosts := pool.HostPlacement.DedicatedHost

			// Check to see if all configured hosts exist
			foundHosts, err := meta.DedicatedHosts(ctx)
			if err != nil {
				allErrs = append(allErrs, field.InternalError(placementPath.Child("dedicatedHost"), err))
			} else {
				// Check the returned configured hosts to see if the dedicated hosts defined in install-config exists.
				for idx, host := range configuredHosts {
					dhPath := placementPath.Child("dedicatedHost").Index(idx)

					// Is host in AWS?
					foundHost, ok := foundHosts[host.ID]
					if !ok {
						errMsg := fmt.Sprintf("dedicated host %s not found", host.ID)
						allErrs = append(allErrs, field.Invalid(dhPath, host, errMsg))
						continue
					}

					// Is host valid for pools region and zone config?
					if !slices.Contains(pool.Zones, foundHost.Zone) {
						errMsg := fmt.Sprintf("dedicated host %s is not available in pool's zone list", host.ID)
						allErrs = append(allErrs, field.Invalid(dhPath, host, errMsg))
					}

					// If user configured the zone for the dedicated host, let's check to make sure its correct
					if host.Zone != "" && host.Zone != foundHost.Zone {
						errMsg := fmt.Sprintf("dedicated host was configured with zone %v but expected zone %v", host.Zone, foundHost.Zone)
						allErrs = append(allErrs, field.Invalid(dhPath.Child("zone"), host, errMsg))
					}
				}
			}
		}
	}

	return allErrs
}

func validateSecurityGroupIDs(ctx context.Context, meta *Metadata, fldPath *field.Path, platform *awstypes.Platform, pool *awstypes.MachinePool) field.ErrorList {
	allErrs := field.ErrorList{}

	vpc, err := meta.VPCID(ctx)
	if err != nil {
		errMsg := fmt.Sprintf("could not determine cluster VPC: %s", err.Error())
		return append(allErrs, field.Invalid(fldPath, vpc, errMsg))
	}

	session, err := meta.Session(ctx)
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, fmt.Errorf("unable to retrieve aws session: %w", err)))
	}

	securityGroups, err := DescribeSecurityGroups(ctx, session, pool.AdditionalSecurityGroupIDs, platform.Region)
	if err != nil {
		return append(allErrs, field.Invalid(fldPath, pool.AdditionalSecurityGroupIDs, err.Error()))
	}

	for _, sg := range securityGroups {
		sgVpcID := *sg.VpcId
		if sgVpcID != vpc {
			errMsg := fmt.Sprintf("sg %s is associated with vpc %s not the provided vpc %s", *sg.GroupId, sgVpcID, vpc)
			allErrs = append(allErrs, field.Invalid(fldPath, sgVpcID, errMsg))
		}
	}

	return allErrs
}

func validateSubnetCIDR(fldPath *field.Path, subnetDataGroup map[string]subnetData, networks []types.MachineNetworkEntry) field.ErrorList {
	allErrs := field.ErrorList{}
	for id, subnetData := range subnetDataGroup {
		fp := fldPath.Index(subnetData.Idx)
		cidr, _, err := net.ParseCIDR(subnetData.CIDR)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fp, id, err.Error()))
			continue
		}
		allErrs = append(allErrs, validateMachineNetworksContainIP(fp, networks, id, cidr)...)
	}
	return allErrs
}

func validateMachineNetworksContainIP(fldPath *field.Path, networks []types.MachineNetworkEntry, subnetName string, ip net.IP) field.ErrorList {
	for _, network := range networks {
		if network.CIDR.Contains(ip) {
			return nil
		}
	}
	return field.ErrorList{field.Invalid(fldPath, subnetName, fmt.Sprintf("subnet's CIDR range start %s is outside of the specified machine networks", ip))}
}

func validateDuplicateSubnetZones(fldPath *field.Path, subnetDataGroup map[string]subnetData, typ string) field.ErrorList {
	subnetIDs := make([]string, 0)
	for id := range subnetDataGroup {
		subnetIDs = append(subnetIDs, id)
	}
	sort.Strings(subnetIDs)

	allErrs := field.ErrorList{}
	zones := map[string]string{}
	for _, id := range subnetIDs {
		subnetData := subnetDataGroup[id]
		if conflictingSubnet, ok := zones[subnetData.Zone.Name]; ok {
			errMsg := fmt.Sprintf("%s subnet %s is also in zone %s", typ, conflictingSubnet, subnetData.Zone.Name)
			allErrs = append(allErrs, field.Invalid(fldPath.Index(subnetData.Idx), id, errMsg))
		} else {
			zones[subnetData.Zone.Name] = id
		}
	}
	return allErrs
}

// validateSubnetRoles ensures BYO subnets have valid roles assigned if roles are provided.
func validateSubnetRoles(fldPath *field.Path, subnetsWithRole map[awstypes.SubnetRoleType][]subnetData, subnetDataGroups subnetDataGroups, config *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	supportedRoles := []awstypes.SubnetRoleType{
		awstypes.ClusterNodeSubnetRole,
		awstypes.EdgeNodeSubnetRole,
		awstypes.BootstrapNodeSubnetRole,
		awstypes.IngressControllerLBSubnetRole,
		awstypes.ControlPlaneExternalLBSubnetRole,
		awstypes.ControlPlaneInternalLBSubnetRole,
	}

	// Subnets of the same role must be in different AZs.
	// Especially, IngressControllerLB subnets must be in different AZs as required by AWS CCM.
	for _, role := range supportedRoles {
		snZones := make(map[string]string)
		for _, subnetData := range subnetsWithRole[role] {
			if conflictingSubnet, ok := snZones[subnetData.Zone.Name]; ok {
				allErrs = append(allErrs, field.Invalid(fldPath.Index(subnetData.Idx), subnetData.ID,
					fmt.Sprintf("subnets %s and %s have role %s and are both in zone %s", conflictingSubnet, subnetData.ID, role, subnetData.Zone.Name)))
			} else {
				snZones[subnetData.Zone.Name] = subnetData.ID
			}
		}
	}

	// BootstrapNode subnets must be assigned to public subnets
	// in external cluster.
	for _, bstrSubnet := range subnetsWithRole[awstypes.BootstrapNodeSubnetRole] {
		// We validate edge subnets in subsequent validations.
		if _, ok := subnetDataGroups.Edge[bstrSubnet.ID]; ok {
			continue
		}
		if config.Publish == types.ExternalPublishingStrategy && !bstrSubnet.Public {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(bstrSubnet.Idx), bstrSubnet.ID,
				fmt.Sprintf("subnet %s has role %s, but is private, expected to be public", bstrSubnet.ID, awstypes.BootstrapNodeSubnetRole)))
		}
	}

	// ClusterNode subnets must be assigned to private subnets
	// unless cluster is public-only.
	for _, cnSubnet := range subnetsWithRole[awstypes.ClusterNodeSubnetRole] {
		// We validate edge subnets in subsequent validations.
		if _, ok := subnetDataGroups.Edge[cnSubnet.ID]; ok {
			continue
		}
		if cnSubnet.Public && !awstypes.IsPublicOnlySubnetsEnabled() {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(cnSubnet.Idx), cnSubnet.ID,
				fmt.Sprintf("subnet %s has role %s, but is public, expected to be private", cnSubnet.ID, awstypes.ClusterNodeSubnetRole)))
		}
	}

	// Type of ControlPlaneLB subnets must match its scope:
	// - ControlPlaneInternalLB subnets must be private
	// - ControlPlaneExternalLB subnets must be public.
	// Private cluster must not have ControlPlaneExternalLB subnets (i.e. statically validated in pkg/types/aws/validation/platform.go).
	for _, ctrlPSubnet := range subnetsWithRole[awstypes.ControlPlaneInternalLBSubnetRole] {
		// We validate edge subnets in subsequent validations.
		if _, ok := subnetDataGroups.Edge[ctrlPSubnet.ID]; ok {
			continue
		}
		if ctrlPSubnet.Public && !awstypes.IsPublicOnlySubnetsEnabled() {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(ctrlPSubnet.Idx), ctrlPSubnet.ID,
				fmt.Sprintf("subnet %s has role %s, but is public, expected to be private", ctrlPSubnet.ID, awstypes.ControlPlaneInternalLBSubnetRole)))
		}
	}
	for _, ctrlPSubnet := range subnetsWithRole[awstypes.ControlPlaneExternalLBSubnetRole] {
		// We validate edge subnets in subsequent validations.
		if _, ok := subnetDataGroups.Edge[ctrlPSubnet.ID]; ok {
			continue
		}
		if !ctrlPSubnet.Public {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(ctrlPSubnet.Idx), ctrlPSubnet.ID,
				fmt.Sprintf("subnet %s has role %s, but is private, expected to be public", ctrlPSubnet.ID, awstypes.ControlPlaneExternalLBSubnetRole)))
		}
	}

	// Type of IngressControllerLB subnets must match cluster scope:
	// - In public cluster, only public IngressControllerLB subnets is allowed.
	// - In private cluster, only private IngressControllerLB subnets is allowed.
	for _, ingressSubnet := range subnetsWithRole[awstypes.IngressControllerLBSubnetRole] {
		// We validate edge subnets in subsequent validations.
		if _, ok := subnetDataGroups.Edge[ingressSubnet.ID]; ok {
			continue
		}

		if ingressSubnet.Public != config.PublicIngress() {
			subnetType := "private"
			if ingressSubnet.Public {
				subnetType = "public"
			}
			allErrs = append(allErrs, field.Invalid(fldPath.Index(ingressSubnet.Idx), ingressSubnet.ID,
				fmt.Sprintf("subnet %s has role %s and is %s, which is not allowed when publish is set to %s", ingressSubnet.ID, awstypes.IngressControllerLBSubnetRole, subnetType, config.Publish)))
		}
	}

	// AZs of LB subnets match AZs of ClusterNode subnets.
	lbRoles := []awstypes.SubnetRoleType{
		awstypes.ControlPlaneInternalLBSubnetRole,
		awstypes.IngressControllerLBSubnetRole,
	}
	if config.PublicAPI() {
		lbRoles = append(lbRoles, awstypes.ControlPlaneExternalLBSubnetRole)
	}
	for _, role := range lbRoles {
		allErrs = append(allErrs, validateLBSubnetAZMatchClusterNodeAZ(fldPath, subnetDataGroups, role, subnetsWithRole[role], subnetsWithRole[awstypes.ClusterNodeSubnetRole])...)
	}

	// EdgeNode subnets must be subnets in Local or Wavelength Zones.
	for _, edgeSubnet := range subnetsWithRole[awstypes.EdgeNodeSubnetRole] {
		if _, ok := subnetDataGroups.Edge[edgeSubnet.ID]; !ok {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(edgeSubnet.Idx), edgeSubnet.ID,
				fmt.Sprintf("subnet %s has role %s, but is not in a Local or WaveLength Zone", edgeSubnet.ID, awstypes.EdgeNodeSubnetRole)))
		}
	}

	// Subnets that are in Local or Wavelength Zones must only have EdgeNode role.
	for _, edgeSubnet := range subnetDataGroups.Edge {
		for _, role := range edgeSubnet.Roles {
			if role.Type != awstypes.EdgeNodeSubnetRole {
				allErrs = append(allErrs, field.Invalid(fldPath.Index(edgeSubnet.Idx), edgeSubnet.ID,
					fmt.Sprintf("subnet %s must only be assigned role %s since it is in a Local or WaveLength Zone", edgeSubnet.ID, awstypes.EdgeNodeSubnetRole)))
				break
			}
		}
	}

	return allErrs
}

// validateUntaggedSubnets ensures there are no additional untagged subnets in the BYO VPC.
// An untagged subnet is a subnet without tag kubernetes.io/cluster/<cluster-id>.
// Untagged subnets may be selected by the CCM, leading to various bugs, RFEs, and support cases. See:
// - https://issues.redhat.com/browse/OCPBUGS-17432.
// - https://issues.redhat.com/browse/RFE-2816.
func validateUntaggedSubnets(ctx context.Context, fldPath *field.Path, meta *Metadata, subnetDataGroups subnetDataGroups) field.ErrorList {
	allErrs := field.ErrorList{}

	vpcSubnets, err := meta.VPCSubnets(ctx)
	if err != nil {
		return append(allErrs, field.Invalid(fldPath, meta.ProvidedSubnets, err.Error()))
	}

	untaggedSubnetIDs := make([]string, 0)
	for _, subnet := range mergeSubnets(vpcSubnets.Public, vpcSubnets.Private, vpcSubnets.Edge) {
		// We only check other subnets in the VPC that are not provided in the install-config.
		if _, ok := subnetDataGroups.All[subnet.ID]; !ok && !subnet.Tags.HasTagKeyPrefix(TagNameKubernetesClusterPrefix) {
			untaggedSubnetIDs = append(untaggedSubnetIDs, subnet.ID)
		}
	}
	sort.Strings(untaggedSubnetIDs)

	if len(untaggedSubnetIDs) > 0 {
		errMsg := fmt.Sprintf("additional subnets %v without tag prefix %s are found in vpc %s of provided subnets. %s", untaggedSubnetIDs, TagNameKubernetesClusterPrefix, vpcSubnets.VpcID,
			fmt.Sprintf("Please add a tag %s to those subnets to exclude them from cluster installation or explicitly assign roles in the install-config to provided subnets", TagNameKubernetesUnmanaged))
		allErrs = append(allErrs, field.Forbidden(fldPath, errMsg))
	}

	return allErrs
}

// validateSharedVPC ensures the BYO VPC can be shared to install the new cluster.
// That is the VPC must not have have tag: kubernetes.io/cluster/<another-cluster-id>: owned.
func validateSharedVPC(ctx context.Context, meta *Metadata, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	vpc, err := meta.VPC(ctx)
	if err != nil {
		return append(allErrs, field.Invalid(fldPath, meta.ProvidedSubnets, err.Error()))
	}

	if vpc.Tags.HasClusterOwnedTag() {
		clusterIDs := vpc.Tags.GetOwnedClusterIDs()
		allErrs = append(allErrs, field.Forbidden(fldPath,
			fmt.Sprintf("VPC of subnets is owned by other clusters %v and cannot be used for new installations, another VPC must be created separately", clusterIDs)))
	}

	return allErrs
}

// validateSharedSubnets ensures the BYO subnets can be shared to install the new cluster.
// That is the subnets must not have have tag: kubernetes.io/cluster/<another-cluster-id>: owned.
func validateSharedSubnets(ctx context.Context, meta *Metadata, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	subnets, err := meta.Subnets(ctx)
	if err != nil {
		return append(allErrs, field.Invalid(fldPath, meta.ProvidedSubnets, err.Error()))
	}

	for id, subnet := range mergeSubnets(subnets.Private, subnets.Public, subnets.Edge) {
		if subnet.Tags.HasClusterOwnedTag() {
			clusterIDs := subnet.Tags.GetOwnedClusterIDs()
			allErrs = append(allErrs, field.Forbidden(fldPath, fmt.Sprintf("subnet %s is owned by other clusters %v and cannot be used for new installations, another subnet must be created separately", id, clusterIDs)))
		}
	}

	return allErrs
}

// validateLBSubnetAZMatchClusterNodeAZ ensures AZs of LB subnets match AZs of ClusterNode subnets.
// AWS load balancers will NOT register a node located in an AZ that is not enabled for the load balancer.
func validateLBSubnetAZMatchClusterNodeAZ(fldPath *field.Path, subnetDataGroups subnetDataGroups, lbType awstypes.SubnetRoleType, lbSubnets []subnetData, clusterNodeSubnets []subnetData) field.ErrorList {
	allErrs := field.ErrorList{}

	lbZoneSet := sets.New[string]()
	for _, subnet := range lbSubnets {
		// We validate edge subnets in another place.
		if _, ok := subnetDataGroups.Edge[subnet.ID]; ok {
			continue
		}
		lbZoneSet.Insert(subnet.Zone.Name)
	}

	nodeZoneSet := sets.New[string]()
	for _, subnet := range clusterNodeSubnets {
		// We validate edge subnets in another place.
		if _, ok := subnetDataGroups.Edge[subnet.ID]; ok {
			continue
		}
		nodeZoneSet.Insert(subnet.Zone.Name)
	}

	// If the nodes use an AZ that is not in load balancer enabled AZs,
	// the router pod might be scheduled to nodes that the load balancer cannot reach.
	if diffSet := nodeZoneSet.Difference(lbZoneSet); diffSet.Len() > 0 {
		allErrs = append(allErrs, field.Forbidden(fldPath, fmt.Sprintf("zones %v are not enabled for %s load balancers, nodes in those zones are unreachable", sets.List(diffSet), lbType)))
	}

	// If the load balancer includes an AZ that is not in node AZs,
	// there will be no nodes in that AZ for the load balancer to register (i.e. not in use)
	if diffSet := lbZoneSet.Difference(nodeZoneSet); diffSet.Len() > 0 {
		allErrs = append(allErrs, field.Forbidden(fldPath, fmt.Sprintf("zones %v are enabled for %s load balancers, but are not used by any nodes", sets.List(diffSet), lbType)))
	}

	return allErrs
}

func validateServiceEndpoints(fldPath *field.Path, region string, services []awstypes.ServiceEndpoint) field.ErrorList {
	allErrs := field.ErrorList{}
	// Validate the endpoint overrides for all provided services.
	// The following is the list of required services by the installer. When
	// an override is not provided for these services, the default endpoint will be used.
	//	"ec2", "elasticloadbalancing", "iam", "route53", "s3", "sts", "tagging",
	for id, service := range services {
		err := validateEndpointAccessibility(service.URL)
		if err != nil {
			logrus.Debugf("failed to access %s endpoint at %s", service.Name, service.URL)
			allErrs = append(allErrs, field.Invalid(fldPath.Index(id).Child("url"), service.URL, err.Error()))
		}
	}

	return allErrs
}

func validateZoneLocal(ctx context.Context, meta *Metadata, fldPath *field.Path, zoneName string) *field.Error {
	sess, err := meta.Session(ctx)
	if err != nil {
		return field.Invalid(fldPath, zoneName, fmt.Sprintf("unable to retrieve aws session: %s", err.Error()))
	}
	zones, err := describeFilteredZones(ctx, sess, meta.Region, []string{zoneName})
	if err != nil {
		return field.Invalid(fldPath, zoneName, fmt.Sprintf("unable to get describe zone: %s", err.Error()))
	}
	validZone := false
	for _, zone := range zones {
		if aws.StringValue(zone.ZoneName) == zoneName {
			switch aws.StringValue(zone.ZoneType) {
			case awstypes.LocalZoneType, awstypes.WavelengthZoneType:
			default:
				return field.Invalid(fldPath, zoneName, fmt.Sprintf("only zone type local-zone or wavelength-zone are valid in the edge machine pool: %s", aws.StringValue(zone.ZoneType)))
			}
			if aws.StringValue(zone.OptInStatus) != awstypes.ZoneOptInStatusOptedIn {
				return field.Invalid(fldPath, zoneName, fmt.Sprintf("zone group is not opted-in: %s", aws.StringValue(zone.GroupName)))
			}
			validZone = true
		}
	}
	if !validZone {
		return field.Invalid(fldPath, zoneName, fmt.Sprintf("invalid local zone name: %s", zoneName))
	}
	return nil
}

func validateEndpointAccessibility(endpointURL string) error {
	if _, err := url.Parse(endpointURL); err != nil {
		return fmt.Errorf("failed to parse service endpoint url: %w", err)
	}
	if _, err := http.Head(endpointURL); err != nil { //nolint:gosec
		return fmt.Errorf("failed to connect to service endpoint url: %w", err)
	}
	return nil
}

var requiredServices = []string{
	"ec2",
	"elasticloadbalancing",
	"iam",
	"route53",
	"s3",
	"sts",
	"tagging",
}

// ValidateForProvisioning validates if the install config is valid for provisioning the cluster.
func ValidateForProvisioning(client API, ic *types.InstallConfig, metadata *Metadata) error {
	if ic.Publish == types.InternalPublishingStrategy && ic.AWS.HostedZone == "" {
		return nil
	}

	if ic.AWS.UserProvisionedDNS == dns.UserProvisionedDNSEnabled {
		logrus.Debug("User Provisioned DNS enabled, skipping zone validation")
		return nil
	}

	var zoneName string
	var zonePath *field.Path
	var zone *route53.HostedZone

	allErrs := field.ErrorList{}
	r53cfg := GetR53ClientCfg(metadata.session, ic.AWS.HostedZoneRole)

	if ic.AWS.HostedZone != "" {
		zoneName = ic.AWS.HostedZone
		zonePath = field.NewPath("aws", "hostedZone")
		zoneOutput, err := client.GetHostedZone(zoneName, r53cfg)
		if err != nil {
			errMsg := fmt.Errorf("unable to retrieve hosted zone: %w", err).Error()
			return field.ErrorList{
				field.Invalid(zonePath, zoneName, errMsg),
			}.ToAggregate()
		}

		if errs := validateHostedZone(zoneOutput, zonePath, zoneName, metadata); len(errs) > 0 {
			allErrs = append(allErrs, errs...)
		}

		zone = zoneOutput.HostedZone
	} else {
		zoneName = ic.BaseDomain
		zonePath = field.NewPath("baseDomain")
		baseDomainOutput, err := client.GetBaseDomain(zoneName)
		if err != nil {
			return field.ErrorList{
				field.Invalid(zonePath, zoneName, "cannot find base domain"),
			}.ToAggregate()
		}

		zone = baseDomainOutput
	}

	if errs := client.ValidateZoneRecords(zone, zoneName, zonePath, ic, r53cfg); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	return allErrs.ToAggregate()
}

func validateHostedZone(hostedZoneOutput *route53.GetHostedZoneOutput, hostedZonePath *field.Path, hostedZoneName string, metadata *Metadata) field.ErrorList {
	allErrs := field.ErrorList{}

	// validate that the hosted zone is associated with the VPC containing the existing subnets for the cluster
	vpcID, err := metadata.VPCID(context.TODO())
	if err == nil {
		if !isHostedZoneAssociatedWithVPC(hostedZoneOutput, vpcID) {
			allErrs = append(allErrs, field.Invalid(hostedZonePath, hostedZoneName, "hosted zone is not associated with the VPC"))
		}
	} else {
		allErrs = append(allErrs, field.Invalid(hostedZonePath, hostedZoneName, "no VPC found"))
	}

	return allErrs
}

func isHostedZoneAssociatedWithVPC(hostedZone *route53.GetHostedZoneOutput, vpcID string) bool {
	if vpcID == "" {
		return false
	}
	for _, vpc := range hostedZone.VPCs {
		if aws.StringValue(vpc.VPCId) == vpcID {
			return true
		}
	}
	return false
}

func validateInstanceProfile(ctx context.Context, meta *Metadata, fldPath *field.Path, pool *awstypes.MachinePool) *field.Error {
	session, err := meta.Session(ctx)
	if err != nil {
		return field.InternalError(fldPath, fmt.Errorf("unable to retrieve aws session: %w", err))
	}
	client := iam.New(session)
	res, err := client.GetInstanceProfileWithContext(ctx, &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(pool.IAMProfile),
	})
	if err != nil {
		msg := fmt.Errorf("unable to retrieve instance profile: %w", err).Error()
		return field.Invalid(fldPath, pool.IAMProfile, msg)
	}
	if len(res.InstanceProfile.Roles) == 0 || res.InstanceProfile.Roles[0] == nil {
		return field.Invalid(fldPath, pool.IAMProfile, "no role attached to instance profile")
	}

	return nil
}
