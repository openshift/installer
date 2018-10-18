package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AvailabilityZones retrieves a list of availability zones for the given region.
func AvailabilityZones(region string) ([]string, error) {
	ec2Client := ec2Client(region)
	zones, err := fetchAvailabilityZones(ec2Client, region)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch availability zones: %v", err)
	}
	return zones, nil
}

func ec2Client(region string) *ec2.EC2 {
	ssn := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(region),
		},
	}))
	return ec2.New(ssn)
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
