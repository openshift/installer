package aws

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
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
	allErrs = append(allErrs, validateAMI(ctx, config)...)
	allErrs = append(allErrs, validatePublicIpv4Pool(ctx, meta, field.NewPath("platform", "aws", "publicIpv4PoolId"), config)...)
	allErrs = append(allErrs, validatePlatform(ctx, meta, field.NewPath("platform", "aws"), config.Platform.AWS, config.Networking, config.Publish)...)

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
			if len(config.Platform.AWS.DeprecatedSubnets) == 0 {
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

func validatePlatform(ctx context.Context, meta *Metadata, fldPath *field.Path, platform *awstypes.Platform, networking *types.Networking, publish types.PublishingStrategy) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validateServiceEndpoints(fldPath.Child("serviceEndpoints"), platform.Region, platform.ServiceEndpoints)...)

	// Fail fast when service endpoints are invalid to avoid long timeouts.
	if len(allErrs) > 0 {
		return allErrs
	}

	if len(platform.DeprecatedSubnets) > 0 {
		allErrs = append(allErrs, validateSubnets(ctx, meta, fldPath.Child("subnets"), platform.DeprecatedSubnets, networking, publish)...)
	} else if awstypes.IsPublicOnlySubnetsEnabled() {
		allErrs = append(allErrs, field.Required(fldPath.Child("subnets"), "subnets must be specified for public-only subnets clusters"))
	}
	if platform.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, validateMachinePool(ctx, meta, fldPath.Child("defaultMachinePlatform"), platform, platform.DefaultMachinePlatform, controlPlaneReq, "", "")...)
	}
	return allErrs
}

func validateAMI(ctx context.Context, config *types.InstallConfig) field.ErrorList {
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
	if partition, partitionFound := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), config.Platform.AWS.Region); partitionFound {
		if partition.ID() == endpoints.AwsPartitionID {
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
		return append(allErrs, field.Invalid(fldPath, nil, fmt.Sprintf("unable to start a session: %s", err.Error())))
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

func validateSubnets(ctx context.Context, meta *Metadata, fldPath *field.Path, subnets []string, networking *types.Networking, publish types.PublishingStrategy) field.ErrorList {
	allErrs := field.ErrorList{}
	privateSubnets, err := meta.PrivateSubnets(ctx)
	if err != nil {
		return append(allErrs, field.Invalid(fldPath, subnets, err.Error()))
	}
	privateSubnetsIdx := map[string]int{}
	for idx, id := range subnets {
		if _, ok := privateSubnets[id]; ok {
			privateSubnetsIdx[id] = idx
		}
	}
	if len(privateSubnets) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath, subnets, "No private subnets found"))
	}

	publicSubnets, err := meta.PublicSubnets(ctx)
	if err != nil {
		return append(allErrs, field.Invalid(fldPath, subnets, err.Error()))
	}
	if publish == types.InternalPublishingStrategy && len(publicSubnets) > 0 {
		logrus.Warnf("Public subnets should not be provided when publish is set to %s", types.InternalPublishingStrategy)
	}
	publicSubnetsIdx := map[string]int{}
	for idx, id := range subnets {
		if _, ok := publicSubnets[id]; ok {
			publicSubnetsIdx[id] = idx
		}
	}
	if len(publicSubnets) == 0 && awstypes.IsPublicOnlySubnetsEnabled() {
		allErrs = append(allErrs, field.Required(fldPath, "public subnets are required for a public-only subnets cluster"))
	}

	edgeSubnets, err := meta.EdgeSubnets(ctx)
	if err != nil {
		return append(allErrs, field.Invalid(fldPath, subnets, err.Error()))
	}
	edgeSubnetsIdx := map[string]int{}
	for idx, id := range subnets {
		if _, ok := edgeSubnets[id]; ok {
			edgeSubnetsIdx[id] = idx
		}
	}

	allErrs = append(allErrs, validateSubnetCIDR(fldPath, privateSubnets, privateSubnetsIdx, networking.MachineNetwork)...)
	allErrs = append(allErrs, validateSubnetCIDR(fldPath, publicSubnets, publicSubnetsIdx, networking.MachineNetwork)...)
	allErrs = append(allErrs, validateDuplicateSubnetZones(fldPath, privateSubnets, privateSubnetsIdx, "private")...)
	allErrs = append(allErrs, validateDuplicateSubnetZones(fldPath, publicSubnets, publicSubnetsIdx, "public")...)
	allErrs = append(allErrs, validateDuplicateSubnetZones(fldPath, edgeSubnets, edgeSubnetsIdx, "edge")...)

	privateZones := sets.New[string]()
	publicZones := sets.New[string]()
	for _, subnet := range privateSubnets {
		privateZones.Insert(subnet.Zone.Name)
	}
	for _, subnet := range publicSubnets {
		publicZones.Insert(subnet.Zone.Name)
	}
	if publish == types.ExternalPublishingStrategy && !publicZones.IsSuperset(privateZones) {
		errMsg := fmt.Sprintf("No public subnet provided for zones %s", sets.List(privateZones.Difference(publicZones)))
		allErrs = append(allErrs, field.Invalid(fldPath, subnets, errMsg))
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
		if len(platform.DeprecatedSubnets) > 0 {
			edgeSubnets, err := meta.EdgeSubnets(ctx)
			if err != nil {
				errMsg := fmt.Sprintf("%s pool. %v", poolName, err.Error())
				return append(allErrs, field.Invalid(field.NewPath("subnets"), platform.DeprecatedSubnets, errMsg))
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
		if len(platform.DeprecatedSubnets) > 0 {
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

func validateSecurityGroupIDs(ctx context.Context, meta *Metadata, fldPath *field.Path, platform *awstypes.Platform, pool *awstypes.MachinePool) field.ErrorList {
	allErrs := field.ErrorList{}

	vpc, err := meta.VPC(ctx)
	if err != nil {
		errMsg := fmt.Sprintf("could not determine cluster VPC: %s", err.Error())
		return append(allErrs, field.Invalid(fldPath, vpc, errMsg))
	}

	securityGroups, err := DescribeSecurityGroups(ctx, meta.session, pool.AdditionalSecurityGroupIDs, platform.Region)
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

func validateSubnetCIDR(fldPath *field.Path, subnets Subnets, idxMap map[string]int, networks []types.MachineNetworkEntry) field.ErrorList {
	allErrs := field.ErrorList{}
	for id, v := range subnets {
		fp := fldPath.Index(idxMap[id])
		cidr, _, err := net.ParseCIDR(v.CIDR)
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

func validateDuplicateSubnetZones(fldPath *field.Path, subnets Subnets, idxMap map[string]int, typ string) field.ErrorList {
	var keys []string
	for id := range subnets {
		keys = append(keys, id)
	}
	sort.Strings(keys)

	allErrs := field.ErrorList{}
	zones := map[string]string{}
	for _, id := range keys {
		subnet := subnets[id]
		if conflictingSubnet, ok := zones[subnet.Zone.Name]; ok {
			errMsg := fmt.Sprintf("%s subnet %s is also in zone %s", typ, conflictingSubnet, subnet.Zone.Name)
			allErrs = append(allErrs, field.Invalid(fldPath.Index(idxMap[id]), id, errMsg))
		} else {
			zones[subnet.Zone.Name] = id
		}
	}
	return allErrs
}

func validateServiceEndpoints(fldPath *field.Path, region string, services []awstypes.ServiceEndpoint) field.ErrorList {
	allErrs := field.ErrorList{}
	ec2Endpoint := ""
	for id, service := range services {
		err := validateEndpointAccessibility(service.URL)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(id).Child("url"), service.URL, err.Error()))
			continue
		}
		if service.Name == ec2.ServiceName {
			ec2Endpoint = service.URL
		}
	}

	if partition, partitionFound := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), region); partitionFound {
		if _, ok := partition.Regions()[region]; !ok && ec2Endpoint == "" {
			err := validateRegion(region)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("region"), region, err.Error()))
			}
		}
		return allErrs
	}

	resolver := newAWSResolver(region, services)
	var errs []error
	for _, service := range requiredServices {
		_, err := resolver.EndpointFor(service, region, endpoints.StrictMatchingOption)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to find endpoint for service %q: %w", service, err))
		}
	}
	if err := utilerrors.NewAggregate(errs); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, services, err.Error()))
	}
	return allErrs
}

func validateRegion(region string) error {
	ses, err := GetSessionWithOptions(func(sess *session.Options) {
		sess.Config.Region = aws.String(region)
	})
	if err != nil {
		return err
	}
	ec2Session := ec2.New(ses)
	return validateEndpointAccessibility(ec2Session.Endpoint)
}

func validateZoneLocal(ctx context.Context, meta *Metadata, fldPath *field.Path, zoneName string) *field.Error {
	sess, err := meta.Session(ctx)
	if err != nil {
		return field.Invalid(fldPath, zoneName, fmt.Sprintf("unable to start a session: %s", err.Error()))
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
	// For each provided service endpoint, verify we can resolve and connect with net.Dial.
	// Ignore e2e.local from unit tests.
	if endpointURL == "e2e.local" {
		return nil
	}
	_, err := url.Parse(endpointURL)
	if err != nil {
		return err
	}
	_, err = http.Head(endpointURL)
	return err
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
	vpcID, err := metadata.VPC(context.TODO())
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
		return field.InternalError(fldPath, fmt.Errorf("unable to start a session: %w", err))
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
