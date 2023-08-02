package firehose

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindDeliveryStreamByName(ctx context.Context, conn *firehose.Firehose, name string) (*firehose.DeliveryStreamDescription, error) {
	input := &firehose.DescribeDeliveryStreamInput{
		DeliveryStreamName: aws.String(name),
	}

	output, err := conn.DescribeDeliveryStreamWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, firehose.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.DeliveryStreamDescription == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.DeliveryStreamDescription, nil
}

func FindDeliveryStreamEncryptionConfigurationByName(ctx context.Context, conn *firehose.Firehose, name string) (*firehose.DeliveryStreamEncryptionConfiguration, error) {
	output, err := FindDeliveryStreamByName(ctx, conn, name)

	if err != nil {
		return nil, err
	}

	if output.DeliveryStreamEncryptionConfiguration == nil {
		return nil, tfresource.NewEmptyResultError(nil)
	}

	return output.DeliveryStreamEncryptionConfiguration, nil
}
