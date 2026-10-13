package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
)

// PreferredInstanceType returns a preferred instance type from the list of instance types provided in descending order of preference
// based on filters like the list of required availability zones.
func PreferredInstanceType(ctx context.Context, meta *awsconfig.Metadata, types []string, zones []string) (string, error) {
	if len(types) == 0 {
		return "", errors.New("at least one instance type required, empty instance types given")
	}

	client, err := awsconfig.NewEC2Client(ctx, awsconfig.EndpointOptions{
		Region:    meta.Region,
		Endpoints: meta.Services,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create EC2 client: %w", err)
	}

	found, err := getInstanceTypeZoneInfo(ctx, client, types, zones)
	if err != nil {
		return types[0], err
	}

	for _, t := range types {
		if found[t].HasAll(zones...) {
			return t, nil
		}
	}

	return types[0], errors.New("no instance type found for the zone constraint")
}

// FilterZonesBasedOnInstanceType return a filtered list of zones where the particular instance type is available. This is mainly necessary for ARM, where the instance type m6g is not
// available in all availability zones.
func FilterZonesBasedOnInstanceType(ctx context.Context, meta *awsconfig.Metadata, instanceType string, zones []string) ([]string, error) {
	client, err := awsconfig.NewEC2Client(ctx, awsconfig.EndpointOptions{
		Region:    meta.Region,
		Endpoints: meta.Services,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create EC2 client: %w", err)
	}

	types := []string{instanceType}
	found, err := getInstanceTypeZoneInfo(ctx, client, types, zones)
	if err != nil {
		return zones, err
	}

	return found[instanceType].Intersection(sets.New(zones...)).UnsortedList(), nil
}

func getInstanceTypeZoneInfo(ctx context.Context, client *ec2.Client, types []string, zones []string) (map[string]sets.Set[string], error) {
	found := map[string]sets.Set[string]{}
	resp, err := client.DescribeInstanceTypeOfferings(ctx, &ec2.DescribeInstanceTypeOfferingsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("location"),
				Values: zones,
			},
			{
				Name:   aws.String("instance-type"),
				Values: types,
			},
		},
		LocationType: ec2types.LocationTypeAvailabilityZone,
	})
	if err != nil {
		return found, err
	}

	// iterate through the offerings and create a map of instance type keys to location values
	for _, offering := range resp.InstanceTypeOfferings {
		f, ok := found[string(offering.InstanceType)]
		if !ok {
			f = sets.New[string]()
			found[string(offering.InstanceType)] = f
		}
		f.Insert(aws.ToString(offering.Location))
	}
	return found, nil
}
