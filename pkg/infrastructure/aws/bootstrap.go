package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
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
	s3Client := s3.New(session)
	s3Tags := make([]*s3.Tag, 0, len(bootstrapInput.additionalTags))
	for k, v := range bootstrapInput.additionalTags {
		k := k
		v := v
		s3Tags = append(s3Tags, &s3.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	// Create an S3 bucket for ignition
	bucketName := fmt.Sprintf("%v-bootstrap", bootstrapInput.clusterID)
	createBucketInput := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	_, err := s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bootstrapInput.ignitionBucket),
	})
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "NotFound" {
			_, err := s3Client.CreateBucket(createBucketInput)
			if err != nil {
				return fmt.Errorf("error creating S3 bucket: %w", err)
			}
			l.WithField("name", bucketName).Infoln("Created s3 bucket")
		} else {
			return err
		}
	} else {
		l.WithField("name", bucketName).Infoln("s3 bucket already exists")
	}

	_, err = s3Client.PutBucketTagging(&s3.PutBucketTaggingInput{
		Bucket: aws.String(bucketName),
		Tagging: &s3.Tagging{
			TagSet: append(s3Tags,
				&s3.Tag{
					Key:   aws.String("Name"),
					Value: aws.String(bucketName),
				},
			),
		},
	})
	if err != nil {
		l.WithError(err).Debugln("failed to tag s3 bucket")
		return fmt.Errorf("failed to tag s3 bucket: %w", err)
	}

	// Upload the bootstrap.ign file to the S3 bucket
	uploadObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(bootstrapInput.ignitionBucket),
		Key:    aws.String("bootstrap.ign"),
		Body:   strings.NewReader(bootstrapInput.ignitionContent),
	}

	_, err = s3Client.PutObject(uploadObjectInput)
	if err != nil {
		return fmt.Errorf("error uploading object to S3: %w", err)
	}

	l.WithField("bucket name", bootstrapInput.ignitionBucket).Infoln("Uploaded bootstrap.ign to S3 bucket")

	// S3 Object tagging supports only up to 10 tags
	// FIXME: make sure that the "owned" tag is not excluded
	if len(s3Tags) > 9 {
		l.Infoln("S3 object accepts up to 10 tags so \"owner\" tag might be missing out of the first 9 tags")
		s3Tags = s3Tags[:9]
	}
	_, err = s3Client.PutObjectTagging(&s3.PutObjectTaggingInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("bootstrap.ign"),
		Tagging: &s3.Tagging{
			TagSet: append(s3Tags, &s3.Tag{
				Key:   aws.String("Name"),
				Value: aws.String(bucketName),
			}),
		},
	})
	if err != nil {
		l.WithError(err).Debugln("failed to tag s3 bucket object")
		return fmt.Errorf("failed to tag ignition s3 bucket object: %w", err)
	}

	iamClient := iam.New(session)
	instanceProfileARN, err := CreateBootstrapInstanceProfile(l, iamClient, bootstrapInput.clusterID, bootstrapInput.additionalTags)
	if err != nil {
		return err
	}

	ec2Tags := make([]*ec2.Tag, 0, len(bootstrapInput.additionalTags))
	for k, v := range bootstrapInput.additionalTags {
		k := k
		v := v
		ec2Tags = append(ec2Tags, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	ec2Client := ec2.New(session)
	bootstrapInstanceOptions := instanceOptions{
		name:                     "bootstrap",
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
		iamInstanceProfileARN:    *instanceProfileARN,
		additionalEC2Tags:        ec2Tags,
	}
	instance, err := createEC2Instance(l, ec2Client, bootstrapInput.clusterID, bootstrapInstanceOptions)
	if err != nil {
		return err
	}

	elbClient := elbv2.New(session)
	for _, targetGroupARN := range bootstrapInput.targetGroupARNs {
		_, err = elbClient.RegisterTargets(&elbv2.RegisterTargetsInput{
			TargetGroupArn: aws.String(targetGroupARN),
			Targets: []*elbv2.TargetDescription{
				{
					Id: instance.PrivateIpAddress,
				},
			},
		})
		if err != nil {
			return err
		}
		l.WithField("id", aws.StringValue(instance.PrivateIpAddress)).WithField("target group", targetGroupARN).Infoln("Target registered")
	}
	l.Infof("Target registered")

	return nil
}

func CreateBootstrapInstanceProfile(l *logrus.Logger, client iamiface.IAMAPI, infraID string, additionalTags map[string]string) (*string, error) {
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

	iamTags := make([]*iam.Tag, 0, len(additionalTags))
	for k, v := range additionalTags {
		k := k
		v := v
		iamTags = append(iamTags, &iam.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	profileName := fmt.Sprintf("%s-bootstrap-profile", infraID)
	roleName := fmt.Sprintf("%s-bootstrap-role", infraID)
	role, err := existingRole(client, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		_, err := client.CreateRole(&iam.CreateRoleInput{
			AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
			Path:                     aws.String("/"),
			RoleName:                 aws.String(roleName),
			Tags: append(iamTags, &iam.Tag{
				Key:   aws.String("Name"),
				Value: aws.String(roleName),
			}),
		})
		if err != nil {
			return nil, fmt.Errorf("cannot create worker role: %w", err)
		}
		l.WithField("name", roleName).Infoln("Created role")
	} else {
		l.WithField("name", roleName).Infoln("Found existing role")
	}
	instanceProfile, err := existingInstanceProfile(client, profileName)
	if err != nil {
		return nil, err
	}
	if instanceProfile == nil {
		result, err := client.CreateInstanceProfile(&iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			Path:                aws.String("/"),
			Tags: append(iamTags, &iam.Tag{
				Key:   aws.String("Name"),
				Value: aws.String(profileName),
			}),
		})
		if err != nil {
			return nil, fmt.Errorf("cannot create instance profile: %w", err)
		}
		instanceProfile = result.InstanceProfile
		l.WithField("name", profileName).Infoln("Created instance profile")
	} else {
		l.WithField("name", profileName).Infoln("Found existing instance profile")
	}
	hasRole := false
	for _, role := range instanceProfile.Roles {
		if aws.StringValue(role.RoleName) == roleName {
			hasRole = true
		}
	}
	if !hasRole {
		_, err = client.AddRoleToInstanceProfile(&iam.AddRoleToInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			RoleName:            aws.String(roleName),
		})
		if err != nil {
			return nil, fmt.Errorf("cannot add role to instance profile: %w", err)
		}
		l.WithField("role", roleName).WithField("profile", profileName).Infoln("Added role to instance profile")
	}
	rolePolicyName := fmt.Sprintf("%s-policy", profileName)
	hasPolicy, err := existingRolePolicy(client, roleName, rolePolicyName)
	if err != nil {
		return nil, err
	}
	if !hasPolicy {
		_, err = client.PutRolePolicy(&iam.PutRolePolicyInput{
			PolicyName:     aws.String(rolePolicyName),
			PolicyDocument: aws.String(bootstrapPolicy),
			RoleName:       aws.String(roleName),
		})
		if err != nil {
			return nil, fmt.Errorf("cannot create profile policy: %w", err)
		}
		l.WithField("name", rolePolicyName).Infoln("Created role policy")
	}

	waitContext, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var lastError error
	wait.Until(func() {
		instanceProfile, err := existingInstanceProfile(client, profileName)
		if err != nil {
			lastError = err
		}
		if instanceProfile != nil {
			l.WithField("name", profileName).Infoln("Instance profile was created and exists")
			lastError = nil
			cancel()
		}
	}, 2*time.Second, waitContext.Done())
	waitErr := waitContext.Err()
	if waitErr != nil {
		if errors.Is(waitErr, context.DeadlineExceeded) {
			return nil, fmt.Errorf("waiting for profile to exist process timed out: %w", lastError)
		}
	}

	// We sleep here otherwise got an error when creating the ec2 instance referencing the profile.
	time.Sleep(10 * time.Second)
	return instanceProfile.Arn, nil
}

func existingRole(client iamiface.IAMAPI, roleName string) (*iam.Role, error) {
	result, err := client.GetRole(&iam.GetRoleInput{RoleName: aws.String(roleName)})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("cannot get existing role: %w", err)
	}
	return result.Role, nil
}

func existingInstanceProfile(client iamiface.IAMAPI, profileName string) (*iam.InstanceProfile, error) {
	result, err := client.GetInstanceProfile(&iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("cannot get existing instance profile: %w", err)
	}
	return result.InstanceProfile, nil
}

func existingRolePolicy(client iamiface.IAMAPI, roleName, policyName string) (bool, error) {
	result, err := client.GetRolePolicy(&iam.GetRolePolicyInput{
		RoleName:   aws.String(roleName),
		PolicyName: aws.String(policyName),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
				return false, nil
			}
		}
		return false, fmt.Errorf("cannot get existing role policy: %w", err)
	}
	return aws.StringValue(result.PolicyName) == policyName, nil
}

type instanceOptions struct {
	name                     string
	amiID                    string
	instanceType             string
	subnetID                 string
	userData                 string
	securityGroupIDs         []string
	associatePublicIPAddress bool
	volumeType               string
	volumeSize               int64
	volumeIOPS               int64
	encrypted                bool
	kmsKeyID                 string
	iamInstanceProfileARN    string
	additionalEC2Tags        []*ec2.Tag
}

var iopsInputPermittedTypes = [...]string{"gp3", "io1", "io2"}

// createEC2Instance creates an EC2 instance and returns its instance ID.
func createEC2Instance(l *logrus.Logger, ec2Client *ec2.EC2, clusterID string, options instanceOptions) (*ec2.Instance, error) {
	// Check if an instance exists.
	existingInstances, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(fmt.Sprintf("%s-%s", clusterID, options.name))},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	// If an instance already exists, return its instance ID.
	if len(existingInstances.Reservations) > 0 && len(existingInstances.Reservations[0].Instances) > 0 {
		l.WithField("id", aws.StringValue(existingInstances.Reservations[0].Instances[0].InstanceId)).Infoln("Instance already exists")
		return existingInstances.Reservations[0].Instances[0], nil
	}

	kmsKeyID := options.kmsKeyID
	// Get default KMS key ID
	if kmsKeyID == "" {
		resp, err := ec2Client.GetEbsDefaultKmsKeyId(&ec2.GetEbsDefaultKmsKeyIdInput{})
		if err != nil {
			return nil, fmt.Errorf("failed to get default KMS key: %w", err)
		}
		kmsKeyID = aws.StringValue(resp.KmsKeyId)
	}

	// InvalidParameterCombination: The parameter iops is not supported for gp2 volumes.
	var iops *int64
	if options.volumeIOPS > 0 {
		for _, permitted := range iopsInputPermittedTypes {
			if options.volumeType == permitted {
				iops = aws.Int64(options.volumeIOPS)
			}
		}
	}

	// Create a new EC2 instance.
	runResult, err := ec2Client.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(options.amiID),
		InstanceType: aws.String(options.instanceType),
		//SubnetId:     aws.String(subnetID),
		NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex:              aws.Int64(0),
				SubnetId:                 aws.String(options.subnetID),
				Groups:                   aws.StringSlice(options.securityGroupIDs),
				AssociatePublicIpAddress: aws.Bool(options.associatePublicIPAddress),
			},
		},
		UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(options.userData))),
		// InvalidParameterCombination: Network interfaces and an instance-level security groups may not be specified on the same request
		// SecurityGroupIds:  aws.StringSlice(options.securityGroupIDs),
		MinCount:          aws.Int64(1),
		MaxCount:          aws.Int64(1),
		TagSpecifications: ec2TagSpecifications("instance", fmt.Sprintf("%s-%s", clusterID, options.name), options.additionalEC2Tags),
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeType: aws.String(options.volumeType),
					VolumeSize: aws.Int64(options.volumeSize),
					Encrypted:  aws.Bool(options.encrypted),
					KmsKeyId:   aws.String(kmsKeyID),
					Iops:       iops,
				},
			},
		},
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Arn: aws.String(options.iamInstanceProfileARN),
		},
	})
	if err != nil {
		return nil, err
	}

	if len(runResult.Instances) > 0 {
		l.WithField("id", aws.StringValue(runResult.Instances[0].InstanceId)).Infoln("Created instance")
		return runResult.Instances[0], nil
	}

	return nil, fmt.Errorf("no instances were created")
}
