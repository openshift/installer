package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	// PresignExpireDuration defines the expiration duration for the generated presign url.
	// Currently, this is used for bootstrap ignition.
	PresignExpireDuration = 60 * time.Minute
)

// PresignedS3URL returns a presigned S3 URL for a bucket/object pair
func PresignedS3URL(ctx context.Context, client *s3.Client, bucket string, object string) (string, error) {
	presignClient := s3.NewPresignClient(client)

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	}, func(po *s3.PresignOptions) {
		po.Expires = PresignExpireDuration
	})
	if err != nil {
		return "", fmt.Errorf("failed to get presigned url for object %s in bucket %s: %w", object, bucket, err)
	}

	return req.URL, nil
}
