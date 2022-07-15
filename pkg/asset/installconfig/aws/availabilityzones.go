package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"
)

// availabilityZones retrieves a list of availability zones for the given region.
func availabilityZones(ctx context.Context, cfg aws.Config, region string) ([]string, error) {
	client := ec2.NewFromConfig(cfg)
	resp, err := client.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("region-name"),
				Values: []string{region},
			},
			{
				Name:   aws.String("state"),
				Values: []string{"available"},
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "fetching availability zones")
	}

	zones := []string{}
	for _, zone := range resp.AvailabilityZones {
		if aws.ToString(zone.ZoneType) == "availability-zone" {
			zones = append(zones, aws.ToString(zone.ZoneName))
		}
	}

	if len(zones) == 0 {
		return nil, errors.Errorf("no available zones in %s", region)
	}

	return zones, nil
}
