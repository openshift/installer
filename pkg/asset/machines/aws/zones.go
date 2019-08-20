package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	awsutil "github.com/openshift/installer/pkg/asset/installconfig/aws"
)

// AvailabilityZones retrieves a list of availability zones for the given region.
func AvailabilityZones(region string, customEndpoints map[string]string) ([]string, error) {
	ec2Client, err := ec2Client(region, customEndpoints)
	if err != nil {
		return nil, err
	}
	zones, err := fetchAvailabilityZones(ec2Client, region)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch availability zones: %v", err)
	}
	return zones, nil
}

func ec2Client(region string, customEndpoints map[string]string) (*ec2.EC2, error) {
	ssn, err := awsutil.GetSession()
	if err != nil {
		return nil, err
	}
	awsConfig := aws.NewConfig()
	if endpoint, ok := customEndpoints["ec2"]; ok {
		awsConfig = awsConfig.WithEndpoint(endpoint)
	}
	client := ec2.New(ssn, awsConfig.WithRegion(region))
	return client, nil
}

func fetchAvailabilityZones(client *ec2.EC2, region string) ([]string, error) {
	req := &ec2.DescribeAvailabilityZonesInput{
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
