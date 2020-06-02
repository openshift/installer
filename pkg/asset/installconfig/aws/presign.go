package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// PresignedS3URL returns a presigned S3 URL for a bucket/object pair
func PresignedS3URL(session *session.Session, region string, bucket string, object string) (string, error) {
	client := s3.New(session, aws.NewConfig().WithRegion(region))
	req, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	})
	presignedURL, err := req.Presign(60 * time.Minute)
	if err != nil {
		return "", err
	}

	return presignedURL, nil
}
