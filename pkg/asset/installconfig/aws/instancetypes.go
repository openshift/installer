package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
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
	Networking   Networking
	Features     []string
}

// instanceTypes retrieves a list of instance types for the given region.
func instanceTypes(ctx context.Context, client *ec2.Client) (map[string]InstanceType, error) {
	types := map[string]InstanceType{}

	paginator := ec2.NewDescribeInstanceTypesPaginator(client, &ec2.DescribeInstanceTypesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list instance types: %w", err)
		}

		for _, sdkTypeInfo := range page.InstanceTypes {
			typeInfo := InstanceType{
				DefaultVCpus: int64(aws.ToInt32(sdkTypeInfo.VCpuInfo.DefaultVCpus)),
				MemInMiB:     aws.ToInt64(sdkTypeInfo.MemoryInfo.SizeInMiB),
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

			types[string(sdkTypeInfo.InstanceType)] = typeInfo
		}
	}

	return types, nil
}
