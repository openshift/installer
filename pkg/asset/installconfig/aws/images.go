package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// ImageInfo holds metadata for an AMI.
type ImageInfo struct {
	BootMode string
}

// images retrieves image metadata for the specified AMI ID.
func images(ctx context.Context, client *ec2.Client, amiID string) (ImageInfo, error) {
	imageOutput, err := client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		ImageIds: []string{amiID},
	})
	if err != nil {
		return ImageInfo{}, fmt.Errorf("fetching images: %w", err)
	}

	if len(imageOutput.Images) == 0 {
		return ImageInfo{}, fmt.Errorf("AMI %s not found", amiID)
	}

	image := imageOutput.Images[0]
	return ImageInfo{
		BootMode: string(image.BootMode),
	}, nil
}
