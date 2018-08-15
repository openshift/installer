package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ghodss/yaml"
	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
	"github.com/pkg/errors"
)

func zoneRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "aws/zones",
		RebuildHelper: zoneRebuilder,
	}

	parents, err := asset.GetParents(ctx, getByName, "aws/region")
	if err != nil {
		return nil, err
	}

	region := aws.String(string(parents["aws/region"].Data))
	ssn := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: region,
		},
	}))

	resp, err := ec2.New(ssn).DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("region-name"),
				Values: []*string{region},
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "describe availability zones")
	}

	zones := []string{}
	for _, zone := range resp.AvailabilityZones {
		zones = append(zones, *zone.ZoneName)
	}

	asset.Data, err = yaml.Marshal(zones)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	installerassets.Rebuilders["aws/zones"] = zoneRebuilder
}
