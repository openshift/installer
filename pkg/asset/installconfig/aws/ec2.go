package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// GetRegions get all regions that are accessible.
func GetRegions(ctx context.Context, client *ec2.Client) ([]string, error) {
	output, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{AllRegions: aws.Bool(true)})
	if err != nil {
		return nil, fmt.Errorf("failed to get all regions: %w", err)
	}

	regions := []string{}
	for _, region := range output.Regions {
		regions = append(regions, aws.ToString(region.RegionName))
	}
	return regions, nil
}

// DescribeSecurityGroups returns the list of ec2 Security Groups that contain the group id and vpc id.
func DescribeSecurityGroups(ctx context.Context, client *ec2.Client, securityGroupIDs []string) ([]ec2types.SecurityGroup, error) {
	cctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	sgOutput, err := client.DescribeSecurityGroups(cctx, &ec2.DescribeSecurityGroupsInput{GroupIds: securityGroupIDs})
	if err != nil {
		return nil, err
	}
	return sgOutput.SecurityGroups, nil
}

// DescribePublicIpv4Pool returns the ec2 public IPv4 Pool attributes from the given ID.
func DescribePublicIpv4Pool(ctx context.Context, client *ec2.Client, poolID string) (ec2types.PublicIpv4Pool, error) {
	cctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	poolOutputs, err := client.DescribePublicIpv4Pools(cctx, &ec2.DescribePublicIpv4PoolsInput{PoolIds: []string{poolID}})
	if err != nil {
		return ec2types.PublicIpv4Pool{}, err
	}
	if len(poolOutputs.PublicIpv4Pools) == 0 {
		return ec2types.PublicIpv4Pool{}, fmt.Errorf("public IPv4 Pool not found: %s", poolID)
	}
	// it should not happen
	if len(poolOutputs.PublicIpv4Pools) > 1 {
		return ec2types.PublicIpv4Pool{}, fmt.Errorf("more than one Public IPv4 Pool: %s", poolID)
	}
	return poolOutputs.PublicIpv4Pools[0], nil
}
