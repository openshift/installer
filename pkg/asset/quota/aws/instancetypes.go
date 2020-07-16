package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

// InstanceTypeInfo describes the instance type
type InstanceTypeInfo struct {
	Name string
	vCPU int64
}

// InstanceTypes returns information of the all the instance types available for a region.
// It returns a map of instance type name to it's information.
func InstanceTypes(ctx context.Context, sess *session.Session, region string) (map[string]InstanceTypeInfo, error) {
	ret := map[string]InstanceTypeInfo{}

	client := ec2.New(sess, aws.NewConfig().WithRegion(region))
	if err := client.DescribeInstanceTypesPagesWithContext(ctx,
		&ec2.DescribeInstanceTypesInput{},
		func(page *ec2.DescribeInstanceTypesOutput, lastPage bool) bool {
			for _, info := range page.InstanceTypes {
				ti := InstanceTypeInfo{Name: aws.StringValue(info.InstanceType)}
				if info.VCpuInfo == nil {
					continue
				}
				ti.vCPU = aws.Int64Value(info.VCpuInfo.DefaultVCpus)
				ret[ti.Name] = ti
			}
			return !lastPage
		}); err != nil {
		return nil, err
	}

	return ret, nil
}

// IsUnauthorizedOperation checks if the error is un authorized due to permission failure or lack of service availability.
func IsUnauthorizedOperation(err error) bool {
	if err == nil {
		return false
	}
	var awsErr awserr.Error
	if errors.As(err, &awsErr) {
		// see reference:
		// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/errors-overview.html#CommonErrors
		return awsErr.Code() == "UnauthorizedOperation" || awsErr.Code() == "AuthFailure" || awsErr.Code() == "Blocked"
	}
	return false
}
