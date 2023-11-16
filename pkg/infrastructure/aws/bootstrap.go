package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/sirupsen/logrus"
)

const ignitionKey = "bootstrap.ign"

type bootstrapInputOptions struct {
	instanceInputOptions
	ignitionBucket  string
	ignitionContent string
}

type bootstrapOutput struct {
	privateIP string
	publicIP  string
}

func createBootstrapResources(ctx context.Context, logger logrus.FieldLogger, ec2Client ec2iface.EC2API, iamClient iamiface.IAMAPI, s3Client s3iface.S3API, elbClient elbv2iface.ELBV2API, input *bootstrapInputOptions) (*bootstrapOutput, error) {
	err := ensureIgnition(ctx, logger, s3Client, input.infraID, input.ignitionBucket, input.ignitionContent, input.tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create ignition resources: %w", err)
	}

	profileName := fmt.Sprintf("%s-bootstrap", input.infraID)
	instanceProfile, err := createBootstrapInstanceProfile(ctx, logger, iamClient, profileName, input.iamRole, input.tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create bootstrap instance profile: %w", err)
	}
	input.instanceProfileARN = aws.StringValue(instanceProfile.Arn)

	input.name = fmt.Sprintf("%s-bootstrap", input.infraID)
	instance, err := ensureInstance(ctx, logger, ec2Client, elbClient, &input.instanceInputOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create bootstrap instace: %w", err)
	}

	return &bootstrapOutput{
		privateIP: aws.StringValue(instance.PrivateIpAddress),
		publicIP:  aws.StringValue(instance.PublicIpAddress),
	}, nil
}

func ensureIgnition(ctx context.Context, logger logrus.FieldLogger, client s3iface.S3API, infraID string, bucket string, content string, tags map[string]string) error {
	err := ensureIgnitionBucket(ctx, logger, client, bucket, tags)
	if err != nil {
		return err
	}

	// Upload the bootstrap.ign file to the S3 bucket
	_, err = client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(ignitionKey),
		Body:                 strings.NewReader(content),
		ServerSideEncryption: aws.String(s3.ServerSideEncryptionAes256),
	})
	if err != nil {
		return fmt.Errorf("failed to upload %s to bucket: %w", ignitionKey, err)
	}
	logger.Infof("Uploaded %s to S3 bucket", ignitionKey)

	// S3 Object tagging supports only up to 10 tags
	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObjectTagging.html
	objTags := limitTags(tags, 8)
	objTags = mergeTags(objTags, map[string]string{
		clusterOwnedTag(infraID): ownedTagValue,
		"Name":                   bucket,
	})
	_, err = client.PutObjectTaggingWithContext(ctx, &s3.PutObjectTaggingInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(ignitionKey),
		Tagging: &s3.Tagging{
			TagSet: s3Tags(objTags),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to tag ignition object: %w", err)
	}
	logger.Infof("Tagged %s object", ignitionKey)

	return nil
}

func ensureIgnitionBucket(ctx context.Context, logger logrus.FieldLogger, client s3iface.S3API, name string, tags map[string]string) error {
	l := logger.WithField("name", name)
	createdOrFoundMsg := "Found existing ignition bucket"
	err := existingBucket(ctx, client, name)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return err
		}
		createdOrFoundMsg = "Created ignition bucket"
		_, err = client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("failed to create ignition bucket: %w", err)
		}
	}
	l.Infoln(createdOrFoundMsg)

	btags := mergeTags(tags, map[string]string{"Name": name})
	_, err = client.PutBucketTaggingWithContext(ctx, &s3.PutBucketTaggingInput{
		Bucket: aws.String(name),
		Tagging: &s3.Tagging{
			TagSet: s3Tags(btags),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to tag ignition bucket: %w", err)
	}
	l.Infoln("Tagged ignition bucket")

	return nil
}

func existingBucket(ctx context.Context, client s3iface.S3API, name string) error {
	_, err := client.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && (awsErr.Code() == s3.ErrCodeNoSuchBucket || strings.EqualFold(awsErr.Code(), "NotFound")) {
			return errNotFound
		}
		return err
	}
	return nil
}

func s3Tags(tags map[string]string) []*s3.Tag {
	stags := make([]*s3.Tag, 0, len(tags))
	for k, v := range tags {
		k, v := k, v
		stags = append(stags, &s3.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return stags
}

func limitTags(tags map[string]string, size int) map[string]string {
	curSize := 0
	resized := make(map[string]string, size)
	for k, v := range tags {
		if curSize > size {
			break
		}
		resized[k] = v
		curSize++
	}
	return resized
}

func createBootstrapInstanceProfile(ctx context.Context, logger logrus.FieldLogger, client iamiface.IAMAPI, name string, roleName string, tags map[string]string) (*iam.InstanceProfile, error) {
	const (
		assumeRolePolicy = `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
                "Service": "ec2.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}`
		bootstrapPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "ec2:Describe*",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": "ec2:AttachVolume",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": "ec2:DetachVolume",
      "Resource": "*"
    }
  ]
}`
	)

	profileInput := &instanceProfileOptions{
		namePrefix:       name,
		roleName:         roleName,
		assumeRolePolicy: assumeRolePolicy,
		policyDocument:   bootstrapPolicy,
		tags:             tags,
	}
	return createInstanceProfile(ctx, logger, client, profileInput)
}
