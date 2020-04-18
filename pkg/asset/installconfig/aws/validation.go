package aws

import (
	"context"
	"fmt"
	"net"
	"sort"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/pkg/errors"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// Validate executes platform-specific validation.
func Validate(ctx context.Context, meta *Metadata, config *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	if config.Platform.AWS == nil {
		return errors.New(field.Required(field.NewPath("platform", "aws"), "AWS validation requires an AWS platform configuration").Error())
	}
	allErrs = append(allErrs, validatePlatform(ctx, meta, field.NewPath("platform", "aws"), config.Platform.AWS, config.Networking, config.Publish)...)

	if config.ControlPlane != nil && config.ControlPlane.Platform.AWS != nil {
		allErrs = append(allErrs, validateMachinePool(ctx, meta, field.NewPath("controlPlane", "platform", "aws"), config.Platform.AWS, config.ControlPlane.Platform.AWS)...)
	}
	for idx, compute := range config.Compute {
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.AWS != nil {
			allErrs = append(allErrs, validateMachinePool(ctx, meta, fldPath.Child("platform", "aws"), config.Platform.AWS, compute.Platform.AWS)...)
		}
	}
	return allErrs.ToAggregate()
}

func validatePlatform(ctx context.Context, meta *Metadata, fldPath *field.Path, platform *awstypes.Platform, networking *types.Networking, publish types.PublishingStrategy) field.ErrorList {
	allErrs := field.ErrorList{}

	if !isAWSSDKRegion(platform.Region) && platform.AMIID == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("amiID"), "AMI must be provided"))
	}

	if len(platform.Subnets) > 0 {
		allErrs = append(allErrs, validateSubnets(ctx, meta, fldPath.Child("subnets"), platform.Subnets, networking, publish)...)
	}
	if err := validateServiceEndpoints(fldPath.Child("serviceEndpoints"), platform.Region, platform.ServiceEndpoints); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("serviceEndpoints"), platform.ServiceEndpoints, err.Error()))
	}
	if platform.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, validateMachinePool(ctx, meta, fldPath.Child("defaultMachinePlatform"), platform, platform.DefaultMachinePlatform)...)
	}
	return allErrs
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

func validateMachinePool(ctx context.Context, meta *Metadata, fldPath *field.Path, platform *awstypes.Platform, pool *awstypes.MachinePool) field.ErrorList {
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

func validateServiceEndpoints(fldPath *field.Path, region string, services []awstypes.ServiceEndpoint) error {
	if isAWSSDKRegion(region) {
		return nil
	}

	resolver := newAWSResolver(region, services)
	var errs []error
	for _, service := range requiredServices {
		_, err := resolver.EndpointFor(service, region, endpoints.StrictMatchingOption)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to find endpoint for service %q", service))
		}
	}
	return utilerrors.NewAggregate(errs)
}

func isAWSSDKRegion(region string) bool {
	for _, partition := range endpoints.DefaultPartitions() {
		for _, partitionRegion := range partition.Regions() {
			if region == partitionRegion.ID() {
				return true
			}
		}
	}
	return false
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
