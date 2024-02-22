package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

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
		return nil, fmt.Errorf("fetching zones: %w", err)
	}

	return resp.AvailabilityZones, nil
}

// filterZonesByType retrieves a list of zones by a given ZoneType attribute within the region.
// ZoneType can be availability-zone, local-zone or wavelength-zone.
func filterZonesByType(ctx context.Context, session *session.Session, region, zoneType string) ([]string, error) {
	azs, err := describeAvailabilityZones(ctx, session, region, []string{})
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", zoneType, err)
	}
	zones := []string{}
	for _, zone := range azs {
		if aws.StringValue(zone.ZoneType) == zoneType {
			zones = append(zones, aws.StringValue(zone.ZoneName))
		}
	}

	return zones, nil
}

// availabilityZones retrieves a list of zones type 'availability-zone' in the region.
func availabilityZones(ctx context.Context, session *session.Session, region string) ([]string, error) {
	zones, err := filterZonesByType(ctx, session, region, typesaws.AvailabilityZoneType)
	if err != nil {
		return nil, err
	}
	if len(zones) == 0 {
		return nil, fmt.Errorf("no zones with type availability-zone in %s", region)
	}
	return zones, nil
}

// edgeZones retrieves a list of zones type 'local-zone' and 'wavelength-zone' in the region.
func edgeZones(ctx context.Context, session *session.Session, region string) ([]string, error) {
	localZones, err := filterZonesByType(ctx, session, region, typesaws.LocalZoneType)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Local Zone names: %w", err)
	}

	wavelengthZones, err := filterZonesByType(ctx, session, region, typesaws.WavelengthZoneType)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Wavelength Zone names: %w", err)
	}
	edgeZones := make([]string, 0, len(localZones)+len(wavelengthZones))
	edgeZones = append(edgeZones, localZones...)
	edgeZones = append(edgeZones, wavelengthZones...)

	if len(edgeZones) == 0 {
		return nil, fmt.Errorf("unable to find zone types with local-zone or wavelength-zone in the region %s", region)
	}

	return edgeZones, nil
}

// describeFilteredZones retrieves a list of all zones for the given region.
func describeFilteredZones(ctx context.Context, session *session.Session, region string, zones []string) ([]*ec2.AvailabilityZone, error) {
	azs, err := describeAvailabilityZones(ctx, session, region, zones)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", zones, err)
	}

	return azs, nil
}
