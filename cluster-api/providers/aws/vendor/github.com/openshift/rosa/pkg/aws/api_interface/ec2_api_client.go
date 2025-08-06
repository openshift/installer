package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// Ec2ApiClient is an interface that defines the methods that we want to use
// from the Client type in the AWS SDK (github.com/aws/aws-sdk-go-v2/service/ec2)
// The AIM is to only contain methods that are defined in the AWS SDK's EC2
// Client.
// For the cases where logic is desired to be implemened combining EC2 calls and
// other logic use the pkg/aws.Client type.
// If you need to use a method provided by the AWS SDK's EC2 Client but it
// is not defined in this interface then it has to be added and all
// the types implementing this interface have to implement the new method.
// The reason this interface has been defined is so we can perform unit testing
// on methods that make use of the AWS EC2 service.
//

type Ec2ApiClient interface {
	DescribeSecurityGroups(ctx context.Context, params *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options),
	) (*ec2.DescribeSecurityGroupsOutput, error)

	DescribeVpcAttribute(ctx context.Context, params *ec2.DescribeVpcAttributeInput, optFns ...func(*ec2.Options),
	) (*ec2.DescribeVpcAttributeOutput, error)

	DescribeAvailabilityZones(ctx context.Context,
		params *ec2.DescribeAvailabilityZonesInput, optFns ...func(*ec2.Options),
	) (*ec2.DescribeAvailabilityZonesOutput, error)

	DescribeRouteTables(ctx context.Context, params *ec2.DescribeRouteTablesInput, optFns ...func(*ec2.Options),
	) (*ec2.DescribeRouteTablesOutput, error)

	DescribeSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options),
	) (*ec2.DescribeSubnetsOutput, error)

	DescribeInstanceTypeOfferings(ctx context.Context,
		params *ec2.DescribeInstanceTypeOfferingsInput, optFns ...func(*ec2.Options),
	) (*ec2.DescribeInstanceTypeOfferingsOutput, error)

	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options),
	) (*ec2.DescribeInstancesOutput, error)
}

// interface guard to ensure that all methods defined in the Ec2ApiClient
// interface are implemented by the real AWS EC2 client. This interface
// guard should always compile
var _ Ec2ApiClient = (*ec2.Client)(nil)
