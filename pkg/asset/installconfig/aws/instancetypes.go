package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/pkg/errors"
)

// InstanceType holds metadata for an instance type.
type InstanceType struct {
	DefaultVCpus int32
	MemInMiB     int64
}

// instanceTypes retrieves a list of instance types for the configured region.
func instanceTypes(ctx context.Context, config aws.Config) (map[string]InstanceType, error) {
	types := map[string]InstanceType{}

	client := ec2.NewFromConfig(config)
	instancePages := ec2.NewDescribeInstanceTypesPaginator(client, nil)
	for instancePages.HasMorePages() {
		page, err := instancePages.NextPage(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "fetching instance types")
		}
		for _, info := range page.InstanceTypes {
			types[string(info.InstanceType)] = InstanceType{
				DefaultVCpus: aws.ToInt32(info.VCpuInfo.DefaultVCpus),
				MemInMiB:     aws.ToInt64(info.MemoryInfo.SizeInMiB),
			}
		}
	}

	return types, nil
}
