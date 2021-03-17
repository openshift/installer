package aws

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	awsvalidation "github.com/openshift/installer/pkg/types/aws/validation"
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
		allErrs = append(allErrs, validateMachinePool(ctx, meta, field.NewPath("controlPlane", "platform", "aws"), config.Platform.AWS, config.ControlPlane.Platform.AWS, controlPlaneReq)...)
	}
	for idx, compute := range config.Compute {
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.AWS != nil {
			allErrs = append(allErrs, validateMachinePool(ctx, meta, fldPath.Child("platform", "aws"), config.Platform.AWS, compute.Platform.AWS, computeReq)...)
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
		allErrs = append(allErrs, validateMachinePool(ctx, meta, fldPath.Child("defaultMachinePlatform"), platform, platform.DefaultMachinePlatform, controlPlaneReq)...)
	}
	return allErrs
}

func validateAMI(ctx context.Context, config *types.InstallConfig) field.ErrorList {
	// accept AMI from the rhcos stream metadata
	switch config.ControlPlane.Architecture {
	case types.ArchitectureAMD64:
		if sets.NewString(rhcos.AMIRegionsX86_64...).Has(config.Platform.AWS.Region) {
			return nil
		}
	case types.ArchitectureARM64:
		if sets.NewString(rhcos.AMIRegionsAARCH64...).Has(config.Platform.AWS.Region) {
			return nil
		}
	default:
		return field.ErrorList{field.NotSupported(field.NewPath("controlPlane", "architecture"), config.ControlPlane.Architecture, awsvalidation.ValidArchitectureValues)}
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

	allErrs = append(allErrs, validateSubnetCIDR(fldPath, privateSubnets, privateSubnetsIdx, networking.MachineNetwork)...)
	allErrs = append(allErrs, validateSubnetCIDR(fldPath, publicSubnets, publicSubnetsIdx, networking.MachineNetwork)...)
	allErrs = append(allErrs, validateDuplicateSubnetZones(fldPath, privateSubnets, privateSubnetsIdx, "private")...)
	allErrs = append(allErrs, validateDuplicateSubnetZones(fldPath, publicSubnets, publicSubnetsIdx, "public")...)

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

func validateMachinePool(ctx context.Context, meta *Metadata, fldPath *field.Path, platform *awstypes.Platform, pool *awstypes.MachinePool, req resourceRequirements) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(pool.Zones) > 0 {
		availableZones := sets.String{}
		if len(platform.Subnets) > 0 {
			privateSubnets, err := meta.PrivateSubnets(ctx)
			if err != nil {
				return append(allErrs, field.InternalError(fldPath, err))
			}
			for _, subnet := range privateSubnets {
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
	// For each provided service endpoint, verify we can resolve and connect with net.Dial.
	for id, service := range services {
		// Ignore e2e.local from unit tests.
		if service.URL == "e2e.local" {
			continue
		}
		URL, err := url.Parse(service.URL)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(id).Child("url"), service.URL, err.Error()))
			continue
		}
		port := URL.Port()
		if port == "" {
			port = "https"
		}
		conn, err := net.Dial("tcp", net.JoinHostPort(URL.Hostname(), port))
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(id).Child("url"), service.URL, err.Error()))
			continue
		}
		conn.Close()
	}
	if _, partitionFound := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), region); partitionFound {
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
func ValidateForProvisioning(session *session.Session, ic *types.InstallConfig, metadata *Metadata) error {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateExistingHostedZone(session, ic, metadata)...)
	return allErrs.ToAggregate()
}

func validateExistingHostedZone(session *session.Session, ic *types.InstallConfig, metadata *Metadata) field.ErrorList {
	if ic.AWS.HostedZone == "" {
		return nil
	}

	// validate that the hosted zone exists
	hostedZonePath := field.NewPath("aws", "hostedZone")
	client := route53.New(session)
	zone, err := client.GetHostedZone(&route53.GetHostedZoneInput{Id: aws.String(ic.AWS.HostedZone)})
	if err != nil {
		return field.ErrorList{
			field.Invalid(hostedZonePath, ic.AWS.HostedZone, "cannot find hosted zone"),
		}
	}

	allErrs := field.ErrorList{}

	// validate that the hosted zone is associated with the VPC containing the existing subnets for the cluster
	vpcID, err := metadata.VPC(context.TODO())
	if err == nil {
		if !isHostedZoneAssociatedWithVPC(zone, vpcID) {
			allErrs = append(allErrs, field.Invalid(hostedZonePath, ic.AWS.HostedZone, "hosted zone is not associated with the VPC"))
		}
	} else {
		allErrs = append(allErrs, field.Invalid(hostedZonePath, ic.AWS.HostedZone, "no VPC found"))
	}

	dottedClusterDomain := ic.ClusterDomain() + "."

	// validate that the domain of the hosted zone is the cluster domain or a parent of the cluster domain
	if !isHostedZoneDomainParentOfClusterDomain(zone.HostedZone, dottedClusterDomain) {
		allErrs = append(allErrs, field.Invalid(hostedZonePath, ic.AWS.HostedZone,
			fmt.Sprintf("hosted zone domain %q is not a parent of the cluster domain %q", *zone.HostedZone.Name, dottedClusterDomain)))
	}

	// validate that the hosted zone does not already have any record sets for the cluster domain
	var problematicRecords []string
	if err := client.ListResourceRecordSetsPages(
		&route53.ListResourceRecordSetsInput{HostedZoneId: zone.HostedZone.Id},
		func(out *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			for _, recordSet := range out.ResourceRecordSets {
				name := aws.StringValue(recordSet.Name)
				// skip record sets that are not sub-domains of the cluster domain. Such record sets may exist for
				// hosted zones that are used for other clusters or other purposes.
				if !strings.HasSuffix(name, dottedClusterDomain) {
					continue
				}
				// skip record sets that are the cluster domain. Record sets for the cluster domain are fine. If the
				// hosted zone has the name of the cluster domain, then there will be NS and SOA record sets for the
				// cluster domain.
				if len(name) == len(dottedClusterDomain) {
					continue
				}
				problematicRecords = append(problematicRecords, fmt.Sprintf("%s (%s)", name, aws.StringValue(recordSet.Type)))
			}
			return !lastPage
		},
	); err != nil {
		allErrs = append(allErrs, field.InternalError(hostedZonePath,
			errors.Wrapf(err, "could not list record sets for hosted zone %q", ic.AWS.HostedZone)))
	}
	if len(problematicRecords) > 0 {
		detail := fmt.Sprintf(
			"hosted zone already has record sets for the domain of the cluster: [%s]",
			strings.Join(problematicRecords, ", "),
		)
		allErrs = append(allErrs, field.Invalid(hostedZonePath, ic.AWS.HostedZone, detail))
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

func isHostedZoneDomainParentOfClusterDomain(hostedZone *route53.HostedZone, dottedClusterDomain string) bool {
	if *hostedZone.Name == dottedClusterDomain {
		return true
	}
	return strings.HasSuffix(dottedClusterDomain, "."+*hostedZone.Name)
}
