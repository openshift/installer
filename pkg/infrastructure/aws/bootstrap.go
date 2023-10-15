package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
)

type bootstrapInput struct {
	clusterID                string
	ignitionBucket           string
	ignitionContent          string
	userData                 string
	amiID                    string
	instanceType             string
	subnetID                 string
	securityGroupIDs         []string
	associatePublicIPAddress bool
	volumeType               string
	volumeSize               int64
	volumeIOPS               int64
	encrypted                bool
	kmsKeyID                 string
	targetGroupARNs          []string
	additionalTags           map[string]string
}

func createBootstrapResources(l *logrus.Logger, session *session.Session, bootstrapInput *bootstrapInput) error {
	if err := setupIgnition(l, session, bootstrapInput.ignitionBucket, bootstrapInput.ignitionContent, bootstrapInput.additionalTags); err != nil {
		return fmt.Errorf("failed to create ignition resources: %w", err)
	}

	instanceProfileARN, err := CreateBootstrapInstanceProfile(l, session, bootstrapInput.clusterID, bootstrapInput.additionalTags)
	if err != nil {
		return fmt.Errorf("failed to create bootstrap instance profile: %w", err)
	}

	ec2Client := ec2.New(session)
	bootstrapInstanceOptions := instanceOptions{
		name:                     fmt.Sprintf("%s-bootstrap", bootstrapInput.clusterID),
		amiID:                    bootstrapInput.amiID,
		instanceType:             bootstrapInput.instanceType,
		subnetID:                 bootstrapInput.subnetID,
		userData:                 bootstrapInput.userData,
		securityGroupIDs:         bootstrapInput.securityGroupIDs,
		associatePublicIPAddress: bootstrapInput.associatePublicIPAddress,
		volumeType:               bootstrapInput.volumeType,
		volumeSize:               bootstrapInput.volumeSize,
		volumeIOPS:               bootstrapInput.volumeIOPS,
		encrypted:                bootstrapInput.encrypted,
		kmsKeyID:                 bootstrapInput.kmsKeyID,
		iamInstanceProfileARN:    instanceProfileARN,
		additionalEC2Tags:        ec2CreateTags(bootstrapInput.additionalTags),
	}
	instance, err := createInstance(l, ec2Client, bootstrapInstanceOptions)
	if err != nil {
		return fmt.Errorf("failed to create bootstrap instance: %w", err)
	}

	if err := RegisterTargetGroups(l, session, bootstrapInput.targetGroupARNs, aws.StringValue(instance.PrivateIpAddress)); err != nil {
		return fmt.Errorf("failed to register bootstrap target groups: %w", err)
	}

	return nil
}

func RegisterTargetGroups(l *logrus.Logger, session *session.Session, targetGroupARNs []string, ipAddress string) error {
	elbClient := elbv2.New(session)
	for _, targetGroupARN := range targetGroupARNs {
		_, err := elbClient.RegisterTargets(&elbv2.RegisterTargetsInput{
			TargetGroupArn: aws.String(targetGroupARN),
			Targets: []*elbv2.TargetDescription{
				{
					Id: aws.String(ipAddress),
				},
			},
		})
		if err != nil {
			return err
		}
		l.WithField("IP addr", ipAddress).WithField("target group", targetGroupARN).Infoln("Target registered")
	}
	l.Infoln("Target Groups registered")
	return nil
}

func setupIgnition(l *logrus.Logger, session *session.Session, bucketName string, ignitionContent string, tags map[string]string) error {
	s3Client := s3.New(session)
	logger := l.WithField("bucket name", bucketName)

	exists, err := s3BucketExists(s3Client, bucketName)
	if err != nil {
		return fmt.Errorf("could not determine if ignition bucket exits: %w", err)
	}

	// Create an S3 bucket for ignition
	if !exists {
		if err := s3CreateBucket(s3Client, bucketName); err != nil {
			logger.WithError(err).Errorln("could not create ignition bucket")
			return fmt.Errorf("failed to create ignition bucket: %w", err)
		}
		logger.Infoln("Created s3 bucket")
	} else {
		logger.Infoln("s3 bucket already exists")
	}

	s3Tags := s3CreateTags(tags)
	if err := s3TagBucket(s3Client, bucketName, s3Tags); err != nil {
		logger.WithError(err).Errorln("could not tag ignition bucket")
		return fmt.Errorf("failed to tag ignition bucket: %w", err)
	}

	ignitionKey := "bootstrap.ign"
	// Upload the bootstrap.ign file to the S3 bucket
	if err := s3BucketPutObject(s3Client, bucketName, ignitionKey, ignitionContent); err != nil {
		logger.WithError(err).Errorln("could not upload ignition to bucket")
		return fmt.Errorf("failed to upload %s: %w", ignitionKey, err)
	}
	logger.Infoln("Uploaded bootstrap.ign to S3 bucket")

	// S3 Object tagging supports only up to 10 tags
	// FIXME: make sure that the "owned" tag is not excluded
	if len(s3Tags) > 9 {
		logger.Infoln("S3 object accepts up to 10 tags so \"owner\" tag might be missing out of the first 9 tags")
		s3Tags = s3Tags[:9]
	}
	if err := s3BucketTagObject(s3Client, bucketName, ignitionKey, s3Tags); err != nil {
		logger.WithError(err).Errorln("could not tag ignition object")
		return fmt.Errorf("failed to tag ignition object: %w", err)
	}

	return nil
}

func CreateBootstrapInstanceProfile(l *logrus.Logger, session *session.Session, infraID string, additionalTags map[string]string) (string, error) {
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

	return createInstanceProfile(l, session, fmt.Sprintf("%s-bootstrap", infraID), assumeRolePolicy, bootstrapPolicy, additionalTags)
}
