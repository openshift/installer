package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

// InstanceType holds metadata for an instance type.
type InstanceType struct {
	DefaultVCpus int64
	MemInMiB     int64
}

// instanceTypes retrieves a list of instance types for the given region.
func instanceTypes(ctx context.Context, session *session.Session, region string) (map[string]InstanceType, error) {
	types := map[string]InstanceType{}

	client := ec2.New(session, aws.NewConfig().WithRegion(region))
	if err := client.DescribeInstanceTypesPagesWithContext(ctx,
		&ec2.DescribeInstanceTypesInput{},
		func(page *ec2.DescribeInstanceTypesOutput, lastPage bool) bool {
			for _, info := range page.InstanceTypes {
				types[*info.InstanceType] = InstanceType{
					DefaultVCpus: aws.Int64Value(info.VCpuInfo.DefaultVCpus),
					MemInMiB:     aws.Int64Value(info.MemoryInfo.SizeInMiB),
				}
			}
			return !lastPage
		}); err != nil {
		return nil, errors.Wrap(err, "fetching instance types")
	}

	return types, nil
}
