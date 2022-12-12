package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"

	typesaws "github.com/openshift/installer/pkg/types/aws"
)

// describeAvailabilityZones retrieves a list of all zones for the given region.
func describeAvailabilityZones(ctx context.Context, session *session.Session, region string) ([]*ec2.AvailabilityZone, error) {
	client := ec2.New(session, aws.NewConfig().WithRegion(region))
	resp, err := client.DescribeAvailabilityZonesWithContext(ctx, &ec2.DescribeAvailabilityZonesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("region-name"),
				Values: []*string{aws.String(region)},
			},
			{
				Name:   aws.String("state"),
				Values: []*string{aws.String("available")},
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "fetching zones")
	}

	return resp.AvailabilityZones, nil
}

// availabilityZones retrieves a list of zones type 'availability-zone' for the region.
func availabilityZones(ctx context.Context, session *session.Session, region string) ([]string, error) {
	azs, err := describeAvailabilityZones(ctx, session, region)
	if err != nil {
		return nil, errors.Wrap(err, "fetching availability zones")
	}
	zones := []string{}
	for _, zone := range azs {
		if *zone.ZoneType == typesaws.AvailabilityZoneType {
			zones = append(zones, *zone.ZoneName)
		}
	}

	if len(zones) == 0 {
		return nil, errors.Errorf("no available zones in %s", region)
	}

	return zones, nil
}
