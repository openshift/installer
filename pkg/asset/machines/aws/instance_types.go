package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
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

	sess, err := meta.Session(ctx)
	if err != nil {
		return types[0], err
	}

	found, err := getInstanceTypeZoneInfo(ctx, sess, meta.Region, types, zones)
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
	sess, err := meta.Session(ctx)
	if err != nil {
		return zones, err
	}

	types := []string{instanceType}
	found, err := getInstanceTypeZoneInfo(ctx, sess, meta.Region, types, zones)
	if err != nil {
		return zones, err
	}

	return found[instanceType].Intersection(sets.NewString(zones...)).List(), nil
}

func getInstanceTypeZoneInfo(ctx context.Context, session *session.Session, region string, types []string, zones []string) (map[string]sets.String, error) {
	found := map[string]sets.String{}

	client := ec2.New(session, aws.NewConfig().WithRegion(region))
	resp, err := client.DescribeInstanceTypeOfferingsWithContext(ctx, &ec2.DescribeInstanceTypeOfferingsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("location"),
				Values: aws.StringSlice(zones),
			},
			{
				Name:   aws.String("instance-type"),
				Values: aws.StringSlice(types),
			},
		},
		LocationType: aws.String("availability-zone"),
	})
	if err != nil {
		return found, err
	}

	// iterate through the offerings and create a map of instance type keys to location values
	for _, offering := range resp.InstanceTypeOfferings {
		f, ok := found[aws.StringValue(offering.InstanceType)]
		if !ok {
			f = sets.NewString()
			found[aws.StringValue(offering.InstanceType)] = f
		}
		f.Insert(aws.StringValue(offering.Location))
	}
	return found, nil
}
