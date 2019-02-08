package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	awsutil "github.com/openshift/installer/pkg/asset/installconfig/aws"
)

var cache map[string][]string

// AvailabilityZones retrieves a list of availability zones for the given region.
func AvailabilityZones(region string) ([]string, error) {
	if cache == nil {
		cache = map[string][]string{}
	} else if zones, ok := cache[region]; ok {
		return zones, nil
	}

	ec2Client, err := ec2Client(region)
	if err != nil {
		return nil, err
	}
	zones, err := fetchAvailabilityZones(ec2Client, region)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch availability zones: %v", err)
	}
	cache[region] = zones
	return zones, nil
}

func ec2Client(region string) (*ec2.EC2, error) {
	ssn, err := awsutil.GetSession()
	if err != nil {
		return nil, err
	}

	client := ec2.New(ssn, aws.NewConfig().WithRegion(region))
	return client, nil
}

func fetchAvailabilityZones(client *ec2.EC2, region string) ([]string, error) {
	zoneFilter := &ec2.Filter{
		Name:   aws.String("region-name"),
		Values: []*string{aws.String(region)},
	}
	req := &ec2.DescribeAvailabilityZonesInput{
		Filters: []*ec2.Filter{zoneFilter},
	}
	resp, err := client.DescribeAvailabilityZones(req)
	if err != nil {
		return nil, err
	}
	zones := []string{}
	for _, zone := range resp.AvailabilityZones {
		zones = append(zones, *zone.ZoneName)
	}
	return zones, nil
}
