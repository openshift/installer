package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"

	typesaws "github.com/openshift/installer/pkg/types/aws"
)

// Zones stores the map of Zone attributes indexed by Zone Name.
type Zones map[string]*Zone

// Zone stores the Availability or Local Zone attributes used to set machine attributes, and to
// feed VPC resources as a source for for terraform variables.
type Zone struct {

	// Name is the availability, local or wavelength zone name.
	Name string

	// ZoneType is the type of subnet's availability zone.
	// The valid values are availability-zone and local-zone.
	Type string

	// ZoneGroupName is the AWS zone group name.
	// For Availability Zones, this parameter has the same value as the Region name.
	//
	// For Local Zones, the name of the associated group, for example us-west-2-lax-1.
	GroupName string

	// ParentZoneName is the name of the zone that handles some of the Local Zone
	// control plane operations, such as API calls.
	ParentZoneName string

	// PreferredInstanceType is the offered instance type on the subnet's zone.
	// It's used for the edge pools which does not offer the same type across different zone groups.
	PreferredInstanceType string
}

// describeAvailabilityZones retrieves a list of all zones for the given region.
func describeAvailabilityZones(ctx context.Context, session *session.Session, region string, zones []string) ([]*ec2.AvailabilityZone, error) {
	client := ec2.New(session, aws.NewConfig().WithRegion(region))
	input := &ec2.DescribeAvailabilityZonesInput{
		AllAvailabilityZones: aws.Bool(true),
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
	}
	if len(zones) > 0 {
		for _, zone := range zones {
			input.ZoneNames = append(input.ZoneNames, aws.String(zone))
		}
	}
	resp, err := client.DescribeAvailabilityZonesWithContext(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "fetching zones")
	}

	return resp.AvailabilityZones, nil
}

// zonesByType retrieves a list of zones by a given ZoneType attribute within the region.
// ZoneType can be availability-zone, local-zone or wavelength-zone.
func zonesByType(ctx context.Context, session *session.Session, region string, zoneType string) ([]string, error) {
	azs, err := describeAvailabilityZones(ctx, session, region, []string{})
	if err != nil {
		return nil, errors.Wrapf(err, "fetching %s", zoneType)
	}
	zones := []string{}
	for _, zone := range azs {
		if aws.StringValue(zone.ZoneType) == zoneType {
			zones = append(zones, aws.StringValue(zone.ZoneName))
		}
	}

	if len(zones) == 0 {
		return nil, errors.Errorf("no zones with type %s in %s", zoneType, region)
	}

	return zones, nil
}

// availabilityZones retrieves a list of zones type 'availability-zone' for the region.
func availabilityZones(ctx context.Context, session *session.Session, region string) ([]string, error) {
	return zonesByType(ctx, session, region, typesaws.AvailabilityZoneType)
}

// localZones retrieves a list of zones type 'local-zone' for the region.
func localZones(ctx context.Context, session *session.Session, region string) ([]string, error) {
	return zonesByType(ctx, session, region, typesaws.LocalZoneType)
}

// describeFilteredZones retrieves a list of all zones for the given region.
func describeFilteredZones(ctx context.Context, session *session.Session, region string, zones []string) ([]*ec2.AvailabilityZone, error) {
	azs, err := describeAvailabilityZones(ctx, session, region, zones)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching %s", zones)
	}

	return azs, nil
}
