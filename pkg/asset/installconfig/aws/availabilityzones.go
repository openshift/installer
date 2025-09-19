package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	typesaws "github.com/openshift/installer/pkg/types/aws"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
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
func describeAvailabilityZones(ctx context.Context, client *ec2.Client, region string, zones []string) ([]ec2types.AvailabilityZone, error) {
	input := &ec2.DescribeAvailabilityZonesInput{
		AllAvailabilityZones: aws.Bool(true),
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("region-name"),
				Values: []string{region},
			},
			{
				Name:   aws.String("state"),
				Values: []string{"available"},
			},
		},
	}

	if len(zones) > 0 {
		input.ZoneNames = zones
	}

	resp, err := client.DescribeAvailabilityZones(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get availability zones: %w", err)
	}

	return resp.AvailabilityZones, nil
}

// filterZonesByType retrieves a list of zones by a given ZoneType attribute within the region.
// ZoneType can be availability-zone, local-zone or wavelength-zone.
func filterZonesByType(ctx context.Context, client *ec2.Client, region, zoneType string) ([]string, error) {
	azs, err := describeAvailabilityZones(ctx, client, region, []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to filter zones by type %s: %w", zoneType, err)
	}
	zones := []string{}
	for _, zone := range azs {
		if aws.ToString(zone.ZoneType) == zoneType {
			zones = append(zones, aws.ToString(zone.ZoneName))
		}
	}

	return zones, nil
}

// availabilityZones retrieves a list of zones type 'availability-zone' in the region.
func availabilityZones(ctx context.Context, client *ec2.Client, region string) ([]string, error) {
	zones, err := filterZonesByType(ctx, client, region, typesaws.AvailabilityZoneType)
	if err != nil {
		return nil, err
	}
	if len(zones) == 0 {
		return nil, fmt.Errorf("no zones with type availability-zone in %s", region)
	}
	return awsdefaults.SupportedZones(zones), nil
}

// edgeZones retrieves a list of zones type 'local-zone' and 'wavelength-zone' in the region.
func edgeZones(ctx context.Context, client *ec2.Client, region string) ([]string, error) {
	localZones, err := filterZonesByType(ctx, client, region, typesaws.LocalZoneType)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Local Zone names: %w", err)
	}

	wavelengthZones, err := filterZonesByType(ctx, client, region, typesaws.WavelengthZoneType)
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
func describeFilteredZones(ctx context.Context, client *ec2.Client, region string, zones []string) ([]ec2types.AvailabilityZone, error) {
	azs, err := describeAvailabilityZones(ctx, client, region, zones)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", zones, err)
	}

	return azs, nil
}
