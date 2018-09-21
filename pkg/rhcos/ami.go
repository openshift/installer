package rhcos

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	// DefaultChannel is the default RHCOS channel for the cluster.
	DefaultChannel = "tested"
)

// AMI calculates a Red Hat CoreOS AMI.
func AMI(ctx context.Context, channel, region string) (ami string, err error) {
	if channel != DefaultChannel {
		return "", fmt.Errorf("channel %q is not yet supported", channel)
	}

	ssn := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(region),
		},
	}))

	svc := ec2.New(ssn)

	result, err := svc.DescribeImagesWithContext(ctx, &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("name"),
				Values: aws.StringSlice([]string{"rhcos*"}),
			},
			{
				Name:   aws.String("architecture"),
				Values: aws.StringSlice([]string{"x86_64"}),
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: aws.StringSlice([]string{"hvm"}),
			},
			{
				Name:   aws.String("image-type"),
				Values: aws.StringSlice([]string{"machine"}),
			},
			{
				Name:   aws.String("owner-id"),
				Values: aws.StringSlice([]string{"531415883065"}),
			},
			{
				Name:   aws.String("state"),
				Values: aws.StringSlice([]string{"available"}),
			},
		},
	})
	if err != nil {
		return "", err
	}

	var image *ec2.Image
	var created time.Time
	for _, nextImage := range result.Images {
		if nextImage.ImageId == nil || nextImage.CreationDate == nil {
			continue
		}
		nextCreated, err := time.Parse(time.RFC3339, *nextImage.CreationDate)
		if err != nil {
			return "", err
		}

		if image == nil || nextCreated.After(created) {
			image = nextImage
			created = nextCreated
		}
	}

	if image == nil {
		return "", fmt.Errorf("no RHCOS AMIs found in %s", region)
	}

	return *image.ImageId, nil
}
