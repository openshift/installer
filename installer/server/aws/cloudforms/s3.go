package cloudforms

import (
	"fmt"
	"strings"

	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// The bucket name is tectonic.<region>.<domain> as it has to be globally unique
// across all AWS accounts (a domain is typically owned by one user) as well
// as per region. The object name is <cluster-name>/cloudforms.template, which
// helps keeping a finite number of buckets, regardless of the number of
// deployments.
func uniqueS3Bucket(session *session.Session, hostedZoneName string) string {
	return fmt.Sprintf(
		"tectonic.%s.%s",
		aws.StringValue(session.Config.Region),
		hostedZoneName,
	)
}

// uploadS3 puts an object into the specified S3 bucket.
func uploadS3(sess *session.Session, bucket, name string, body []byte) (string, error) {
	svc := s3.New(sess)

	// Try to create the S3 bucket and ignore the error if it was already
	// existing (BucketAlreadyExists / BucketAlreadyOwnedByYou).
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String("private"),
	})
	if err != nil && !strings.HasPrefix(err.Error(), "BucketAlready") {
		return "", err
	}

	// PUT the file in S3.
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(name),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(body),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, name), nil
}

// deleteS3 deletes an object from S3.
func deleteS3(sess *session.Session, bucket, name string) error {
	_, err := s3.New(sess).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	return err
}
