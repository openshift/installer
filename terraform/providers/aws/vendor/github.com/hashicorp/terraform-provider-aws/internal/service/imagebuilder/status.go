package imagebuilder

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// statusImage fetches the Image and its Status
func statusImage(ctx context.Context, conn *imagebuilder.Imagebuilder, imageBuildVersionArn string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &imagebuilder.GetImageInput{
			ImageBuildVersionArn: aws.String(imageBuildVersionArn),
		}

		output, err := conn.GetImageWithContext(ctx, input)

		if err != nil {
			return nil, imagebuilder.ImageStatusPending, err
		}

		if output == nil || output.Image == nil || output.Image.State == nil {
			return nil, imagebuilder.ImageStatusPending, nil
		}

		status := aws.StringValue(output.Image.State.Status)

		if status == imagebuilder.ImageStatusFailed {
			return output.Image, status, fmt.Errorf("%s", aws.StringValue(output.Image.State.Reason))
		}

		return output.Image, status, nil
	}
}
