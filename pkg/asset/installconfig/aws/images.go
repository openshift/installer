package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// ImageInfo holds metadata for an AMI.
type ImageInfo struct {
	BootMode string
}

// images retrieves image metadata for the specified AMI ID in the given region.
func images(ctx context.Context, session *session.Session, region string, amiID string) (ImageInfo, error) {
	client := ec2.New(session, aws.NewConfig().WithRegion(region))

	imageOutput, err := client.DescribeImagesWithContext(ctx, &ec2.DescribeImagesInput{
		ImageIds: []*string{aws.String(amiID)},
	})
	if err != nil {
		return ImageInfo{}, fmt.Errorf("fetching images: %w", err)
	}

	if len(imageOutput.Images) == 0 {
		return ImageInfo{}, fmt.Errorf("AMI %s not found", amiID)
	}

	image := imageOutput.Images[0]
	return ImageInfo{
		BootMode: aws.StringValue(image.BootMode),
	}, nil
}
