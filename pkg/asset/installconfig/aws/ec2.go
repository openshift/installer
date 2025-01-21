package aws

import (
	"context"
	"fmt"
	"time"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	ec2v2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getEC2Client(cfg awsv2.Config) *ec2v2.Client {
	return ec2v2.NewFromConfig(cfg)
}

// GetRegions get all regions that are accessible.
func GetRegions(ctx context.Context, config awsv2.Config, allRegions bool) ([]string, error) {
	client := getEC2Client(config)

	output, err := client.DescribeRegions(ctx, &ec2v2.DescribeRegionsInput{AllRegions: aws.Bool(allRegions)})
	if err != nil {
		return nil, fmt.Errorf("failed to get all regions: %w", err)
	}

	regions := []string{}
	for _, region := range output.Regions {
		regions = append(regions, *region.RegionName)
	}
	return regions, nil
}

// DescribeSecurityGroups returns the list of ec2 Security Groups that contain the group id and vpc id.
func DescribeSecurityGroups(ctx context.Context, session *session.Session, securityGroupIDs []string, region string) ([]*ec2.SecurityGroup, error) {
	client := ec2.New(session, aws.NewConfig().WithRegion(region))

	sgIDPtrs := []*string{}
	for _, sgid := range securityGroupIDs {
		sgid := sgid
		sgIDPtrs = append(sgIDPtrs, &sgid)
	}

	cctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	sgOutput, err := client.DescribeSecurityGroupsWithContext(cctx, &ec2.DescribeSecurityGroupsInput{GroupIds: sgIDPtrs})
	if err != nil {
		return nil, err
	}
	return sgOutput.SecurityGroups, nil
}

// DescribePublicIpv4Pool returns the ec2 public IPv4 Pool attributes from the given ID.
func DescribePublicIpv4Pool(ctx context.Context, session *session.Session, region string, poolID string) (*ec2.PublicIpv4Pool, error) {
	client := ec2.New(session, aws.NewConfig().WithRegion(region))

	cctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	poolOutputs, err := client.DescribePublicIpv4PoolsWithContext(cctx, &ec2.DescribePublicIpv4PoolsInput{PoolIds: []*string{aws.String(poolID)}})
	if err != nil {
		return nil, err
	}
	if len(poolOutputs.PublicIpv4Pools) == 0 {
		return nil, fmt.Errorf("public IPv4 Pool not found: %s", poolID)
	}
	// it should not happen
	if len(poolOutputs.PublicIpv4Pools) > 1 {
		return nil, fmt.Errorf("more than one Public IPv4 Pool: %s", poolID)
	}
	return poolOutputs.PublicIpv4Pools[0], nil
}
