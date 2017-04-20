package cloudforms

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
)

type AwsBucket struct {
	sess           *session.Session
	hostedZoneName string
}

func NewAwsBucket(sess *session.Session,
	hostedZoneName string) *AwsBucket {

	return &AwsBucket{
		sess:           sess,
		hostedZoneName: hostedZoneName,
	}
}

// Bucket returns the genrated unique bucket name in S3
func (a *AwsBucket) Bucket() string {
	return uniqueS3Bucket(a.sess, a.hostedZoneName)
}

// Url returns the complete S3 link to the object
func (a *AwsBucket) Url(filename string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", a.Bucket(), filename)
}

// Upload puts the filename and data in the bucket
func (a *AwsBucket) Upload(filename string, contents []byte) error {
	_, err := uploadS3(a.sess, a.Bucket(), filename, contents)
	return err
}

// Remove removes the specified filename from the S3 bucket
func (a *AwsBucket) Remove(filename string) error {
	return deleteS3(a.sess, a.Bucket(), filename)
}
