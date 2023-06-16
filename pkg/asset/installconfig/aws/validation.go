package aws

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
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
	allErrs = append(allErrs, validatePlatform(ctx, meta, field.NewPath("platform", "aws"), config.Platform.AWS, config.Networking, config.Publish)...)

	if config.ControlPlane != nil && config.ControlPlane.Platform.AWS != nil {
		allErrs = append(allErrs, validateMachinePool(ctx, meta, field.NewPath("controlPlane", "platform", "aws"), config.Platform.AWS, config.ControlPlane.Platform.AWS, controlPlaneReq, "")...)
	}

	for idx, compute := range config.Compute {
		fldPath := field.NewPath("compute").Index(idx)

		// Pool's specific validation.
		// Edge Compute Pool: AWS Local Zones is valid only when installing in existing VPC.
		if compute.Name == types.MachinePoolEdgeRoleName {
			if len(config.Platform.AWS.Subnets) == 0 {
				return errors.New(field.Required(fldPath, "invalid install config. edge machine pool is valid when installing in existing VPC").Error())
			}
			edgeSubnets, err := meta.EdgeSubnets(ctx)
			if err != nil {
				errMsg := fmt.Sprintf("%s pool. %v", compute.Name, err.Error())
				return errors.New(field.Invalid(field.NewPath("platform", "aws", "subnets"), config.Platform.AWS.Subnets, errMsg).Error())
			}
			if len(edgeSubnets) == 0 {
				return errors.New(field.Required(fldPath, "invalid install config. There is no valid subnets for edge machine pool").Error())
			}
		}

		if compute.Platform.AWS != nil {
			allErrs = append(allErrs, validateMachinePool(ctx, meta, fldPath.Child("platform", "aws"), config.Platform.AWS, compute.Platform.AWS, computeReq, compute.Name)...)
		}
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

	if len(platform.Subnets) > 0 {
		allErrs = append(allErrs, validateSubnets(ctx, meta, fldPath.Child("subnets"), platform.Subnets, networking, publish)...)
	}
	if platform.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, validateMachinePool(ctx, meta, fldPath.Child("defaultMachinePlatform"), platform, platform.DefaultMachinePlatform, controlPlaneReq, "")...)
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
	publicSubnetsIdx := map[string]int{}
	for idx, id := range subnets {
		if _, ok := publicSubnets[id]; ok {
			publicSubnetsIdx[id] = idx
		}
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

	privateZones := sets.NewString()
	publicZones := sets.NewString()
	for _, subnet := range privateSubnets {
		privateZones.Insert(subnet.Zone)
	}
	for _, subnet := range publicSubnets {
		publicZones.Insert(subnet.Zone)
	}
	if publish == types.ExternalPublishingStrategy && !publicZones.IsSuperset(privateZones) {
		errMsg := fmt.Sprintf("No public subnet provided for zones %s", privateZones.Difference(publicZones).List())
		allErrs = append(allErrs, field.Invalid(fldPath, subnets, errMsg))
	}

	return allErrs
}

func validateMachinePool(ctx context.Context, meta *Metadata, fldPath *field.Path, platform *awstypes.Platform, pool *awstypes.MachinePool, req resourceRequirements, poolName string) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(pool.Zones) > 0 {
		availableZones := sets.String{}
		if len(platform.Subnets) > 0 {
			var err error
			var subnets Subnets
			switch poolName {
			case types.MachinePoolEdgeRoleName:
				subnets, err = meta.EdgeSubnets(ctx)
			default:
				subnets, err = meta.PrivateSubnets(ctx)
			}
			if err != nil {
				return append(allErrs, field.InternalError(fldPath, err))
			}
			for _, subnet := range subnets {
				availableZones.Insert(subnet.Zone)
			}
		} else {
			allzones, err := meta.AvailabilityZones(ctx)
			if err != nil {
				return append(allErrs, field.InternalError(fldPath, err))
			}
			availableZones.Insert(allzones...)
		}

		if diff := sets.NewString(pool.Zones...).Difference(availableZones); diff.Len() > 0 {
			errMsg := fmt.Sprintf("No subnets provided for zones %s", diff.List())
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
		} else {
			errMsg := fmt.Sprintf("instance type %s not found", pool.InstanceType)
			allErrs = append(allErrs, field.Invalid(fldPath.Child("type"), pool.InstanceType, errMsg))
		}
	}

	if len(pool.AdditionalSecurityGroupIDs) > 0 {
		allErrs = append(allErrs, validateSecurityGroupIDs(ctx, meta, fldPath.Child("additionalSecurityGroupIDs"), platform, pool)...)
	}

	return allErrs
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

func validateSubnetCIDR(fldPath *field.Path, subnets map[string]Subnet, idxMap map[string]int, networks []types.MachineNetworkEntry) field.ErrorList {
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

func validateDuplicateSubnetZones(fldPath *field.Path, subnets map[string]Subnet, idxMap map[string]int, typ string) field.ErrorList {
	var keys []string
	for id := range subnets {
		keys = append(keys, id)
	}
	sort.Strings(keys)

	allErrs := field.ErrorList{}
	zones := map[string]string{}
	for _, id := range keys {
		subnet := subnets[id]
		if conflictingSubnet, ok := zones[subnet.Zone]; ok {
			errMsg := fmt.Sprintf("%s subnet %s is also in zone %s", typ, conflictingSubnet, subnet.Zone)
			allErrs = append(allErrs, field.Invalid(fldPath.Index(idxMap[id]), id, errMsg))
		} else {
			zones[subnet.Zone] = id
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
			errs = append(errs, errors.Wrapf(err, "failed to find endpoint for service %q", service))
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
			errMsg := errors.Wrapf(err, "unable to retrieve hosted zone").Error()
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
