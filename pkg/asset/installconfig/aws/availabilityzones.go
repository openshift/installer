package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

// availabilityZones retrieves a list of availability zones for the given region.
func availabilityZones(ctx context.Context, session *session.Session, region string) ([]string, error) {
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
		return nil, errors.Wrap(err, "fetching availability zones")
	}

	// This finds local zones and allows us to filter them out. It does so by
	// using the fact that local zones have a predictable format based on the
	// region, for example a local zone in us-west-2 will look like
	// "us-west-2-las-1a". This is calculated relative to the region name as not
	// all regions follow the same format with 2 hyphens.
	// This is a temporary fix until we update our aws-sdk-go package to a
	// version that adds the ZoneType field to the AvailabilityZone struct
	numberOfHyphensInRegion := strings.Count(region, "-")

	zones := []string{}
	for _, zone := range resp.AvailabilityZones {
		numberOfHyphensInZone := strings.Count(*zone.ZoneName, "-")
		if numberOfHyphensInZone != numberOfHyphensInRegion+2 {
			zones = append(zones, *zone.ZoneName)
		}
	}

	if len(zones) == 0 {
		return nil, errors.Errorf("no available zones in %s", region)
	}

	return zones, nil
}
