package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Networking describes the network settings for an instance type.
type Networking struct {
	// IPv6Supported indicates whether IPv6 is supported.
	IPv6Supported bool
}

// InstanceType holds metadata for an instance type.
type InstanceType struct {
	DefaultVCpus int64
	MemInMiB     int64
	Arches       []string
	Hypervisor   string
	Networking   Networking
	Features     []string
}

// getInstanceType returns metadata for the named instance type. If the type does
// not exist in the configured region.
func getInstanceType(ctx context.Context, client *ec2.Client, instanceType string) (InstanceType, error) {
	out, err := client.DescribeInstanceTypes(ctx, &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []ec2types.InstanceType{ec2types.InstanceType(instanceType)},
	})
	if err != nil {
		return InstanceType{}, fmt.Errorf("failed to get instance type %s details: %w", instanceType, err)
	}

	// A nonexistent type is reported as an InvalidInstanceType error above, so an
	// empty result here is an unexpected API response rather than a missing type.
	if len(out.InstanceTypes) == 0 {
		return InstanceType{}, fmt.Errorf("unexpected empty response describing instance type %s", instanceType)
	}

	sdkTypeInfo := out.InstanceTypes[0]
	typeInfo := InstanceType{
		DefaultVCpus: int64(aws.ToInt32(sdkTypeInfo.VCpuInfo.DefaultVCpus)),
		MemInMiB:     aws.ToInt64(sdkTypeInfo.MemoryInfo.SizeInMiB),
		Hypervisor:   string(sdkTypeInfo.Hypervisor),
	}

	for _, arch := range sdkTypeInfo.ProcessorInfo.SupportedArchitectures {
		typeInfo.Arches = append(typeInfo.Arches, string(arch))
	}

	if netInfo := sdkTypeInfo.NetworkInfo; netInfo != nil {
		typeInfo.Networking = Networking{
			IPv6Supported: aws.ToBool(netInfo.Ipv6Supported),
		}
	}

	for _, features := range sdkTypeInfo.ProcessorInfo.SupportedFeatures {
		typeInfo.Features = append(typeInfo.Features, string(features))
	}

	return typeInfo, nil
}
