package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// InstanceType holds metadata for an instance type.
type InstanceType struct {
	DefaultVCpus int64
	MemInMiB     int64
	Arches       []ec2types.ArchitectureType
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
		for _, info := range page.InstanceTypes {
			types[string(info.InstanceType)] = InstanceType{
				DefaultVCpus: int64(aws.ToInt32(info.VCpuInfo.DefaultVCpus)),
				MemInMiB:     aws.ToInt64(info.MemoryInfo.SizeInMiB),
				Arches:       info.ProcessorInfo.SupportedArchitectures,
			}
		}
	}

	return types, nil
}
