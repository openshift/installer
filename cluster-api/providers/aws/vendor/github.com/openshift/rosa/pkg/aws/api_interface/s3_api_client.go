package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3ApiClient is an interface that defines the methods that we want to use
// from the Client type in the AWS SDK ("github.com/aws/aws-sdk-go-v2/service/s3")
// The aim is to only contain methods that are defined in the AWS SDK's S3
// Client.
// For the cases where logic is desired to be implemened combining S3 calls
// and other logic use the pkg/aws.Client type.
// If you need to use a method provided by the AWS SDK's S3 Client but it
// is not defined in this interface then it has to be added and all
// the types implementing this interface have to implement the new method.
// The reason this interface has been defined is so we can perform unit testing
// on methods that make use of the AWS S3 service.
//

type S3ApiClient interface {
	CreateBucket(ctx context.Context,
		params *s3.CreateBucketInput, optFns ...func(*s3.Options),
	) (*s3.CreateBucketOutput, error)

	DeleteBucket(ctx context.Context,
		params *s3.DeleteBucketInput, optFns ...func(*s3.Options),
	) (*s3.DeleteBucketOutput, error)
	DeleteObject(ctx context.Context,
		params *s3.DeleteObjectInput, optFns ...func(*s3.Options),
	) (*s3.DeleteObjectOutput, error)

	HeadBucket(context.Context,
		*s3.HeadBucketInput, ...func(*s3.Options),
	) (*s3.HeadBucketOutput, error)

	ListObjects(ctx context.Context,
		params *s3.ListObjectsInput, optFns ...func(*s3.Options),
	) (*s3.ListObjectsOutput, error)

	PutObject(ctx context.Context,
		params *s3.PutObjectInput, optFns ...func(*s3.Options),
	) (*s3.PutObjectOutput, error)

	PutBucketTagging(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options),
	) (*s3.PutBucketTaggingOutput, error)

	PutPublicAccessBlock(ctx context.Context, params *s3.PutPublicAccessBlockInput, optFns ...func(*s3.Options),
	) (*s3.PutPublicAccessBlockOutput, error)

	PutBucketPolicy(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options),
	) (*s3.PutBucketPolicyOutput, error)
}

// interface guard to ensure that all methods defined in the S3ApiClient
// interface are implemented by the real AWS S3 client. This interface
// guard should always compile
var _ S3ApiClient = (*s3.Client)(nil)
