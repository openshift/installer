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
	securityGroupID          string
	associatePublicIPAddress bool
	volumeType               string
	volumeSize               int64
	volumeIOPS               int64
	encrypted                bool
	kmsKeyID                 string
	targetGroupARNs          []string
}

func createBootstrapResources(l *logrus.Logger, session *session.Session, bootstrapInput *bootstrapInput) error {
	s3Client := s3.New(session)
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
			l.Infof("Created s3 bucket: %v", bucketName)
		} else {
			return err
		}
	} else {
		l.Infof("s3 bucket already exists: %v", bucketName)
	}

	_, err = s3Client.PutBucketTagging(&s3.PutBucketTaggingInput{
		Bucket: aws.String(bucketName),
		Tagging: &s3.Tagging{
			TagSet: []*s3.Tag{
				{
					Key:   aws.String(clusterTag(bootstrapInput.clusterID)),
					Value: aws.String(clusterTagValue),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String(bucketName),
				},
			},
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

	l.Infof("Uploaded bootstrap.ign to S3 bucket: %s\n", bootstrapInput.ignitionBucket)

	_, err = s3Client.PutObjectTagging(&s3.PutObjectTaggingInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("bootstrap.ign"),
		Tagging: &s3.Tagging{
			TagSet: []*s3.Tag{
				{
					Key:   aws.String(clusterTag(bootstrapInput.clusterID)),
					Value: aws.String(clusterTagValue),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String(bucketName),
				},
			},
		},
	})
	if err != nil {
		l.WithError(err).Debugln("failed to tag s3 bucket object")
		return fmt.Errorf("failed to tag ignition s3 bucket object: %w", err)
	}

	iamClient := iam.New(session)
	instanceProfileARN, err := CreateBootstrapInstanceProfile(l, iamClient, bootstrapInput.clusterID)
	if err != nil {
		return err
	}

	ec2Client := ec2.New(session)
	instance, err := createEC2Instance(l, ec2Client, bootstrapInput.amiID, bootstrapInput.instanceType, bootstrapInput.subnetID, bootstrapInput.userData, bootstrapInput.securityGroupID,
		bootstrapInput.associatePublicIPAddress, bootstrapInput.clusterID, bootstrapInput.volumeType, bootstrapInput.volumeSize, bootstrapInput.volumeIOPS,
		bootstrapInput.encrypted, bootstrapInput.kmsKeyID, *instanceProfileARN, "bootstrap")
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
		l.Infof("Target registered id: %v, targetGroup: %v", *instance.PrivateIpAddress, targetGroupARN)
	}
	l.Infof("Target registered")

	return nil
}

func CreateBootstrapInstanceProfile(l *logrus.Logger, client iamiface.IAMAPI, infraID string) (*string, error) {
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
			Tags: []*iam.Tag{
				{
					Key:   aws.String(clusterTag(infraID)),
					Value: aws.String(clusterTagValue),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String(roleName),
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("cannot create worker role: %w", err)
		}
		l.Info("Created role", "name", roleName)
	} else {
		l.Info("Found existing role", "name", roleName)
	}
	instanceProfile, err := existingInstanceProfile(client, profileName)
	if err != nil {
		return nil, err
	}
	if instanceProfile == nil {
		result, err := client.CreateInstanceProfile(&iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			Path:                aws.String("/"),
			Tags: []*iam.Tag{
				{
					Key:   aws.String(clusterTag(infraID)),
					Value: aws.String(clusterTagValue),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String(profileName),
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("cannot create instance profile: %w", err)
		}
		instanceProfile = result.InstanceProfile
		l.Info("Created instance profile", "name", profileName)
	} else {
		l.Info("Found existing instance profile", "name", profileName)
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
		l.Info("Added role to instance profile", "role", roleName, "profile", profileName)
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
		l.Info("Created role policy", "name", rolePolicyName)
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
			l.Infof("Instance profile was created and exists: %s", profileName)
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

// createEC2Instance creates an EC2 instance and returns its instance ID.
func createEC2Instance(l *logrus.Logger, ec2Client *ec2.EC2, amiID, instanceType, subnetID, userData, securityGroupID string,
	associatePublicIPAddress bool, clusterID, volumeType string, volumeSize, volumeIOPS int64,
	encrypted bool, kmsKeyID, iamInstanceProfileARN, instanceName string) (*ec2.Instance, error) {

	// Check if an instance exists.
	existingInstances, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(fmt.Sprintf("%s-%s", clusterID, instanceName))},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	// If an instance already exists, return its instance ID.
	if len(existingInstances.Reservations) > 0 && len(existingInstances.Reservations[0].Instances) > 0 {
		l.Infof("Instance already exists: %v", *existingInstances.Reservations[0].Instances[0].InstanceId)
		return existingInstances.Reservations[0].Instances[0], nil
	}

	// Create a new EC2 instance.
	runResult, err := ec2Client.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(amiID),
		InstanceType: aws.String(instanceType),
		//SubnetId:     aws.String(subnetID),
		NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex: aws.Int64(0),
				SubnetId:    aws.String(subnetID),
				// TODO (alberto): Parameterize this.
				Groups:                   []*string{aws.String(securityGroupID)},
				AssociatePublicIpAddress: aws.Bool(associatePublicIPAddress),
			},
		},
		UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(userData))),
		// TODO(alberto): create dedicated SGs.
		// SecurityGroupIds:         []*string{aws.String(securityGroupID)},
		MinCount: aws.Int64(1),
		MaxCount: aws.Int64(1),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(fmt.Sprintf("%s-%s", clusterID, instanceName)),
					},
					{
						Key:   aws.String(clusterTag(clusterID)),
						Value: aws.String(clusterTagValue),
					},
				},
			},
		},
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeType: aws.String(volumeType),
					VolumeSize: aws.Int64(volumeSize),
					// TODO(alberto): Parameterize this.
					Encrypted: aws.Bool(false),
					//Encrypted:  aws.Bool(encrypted),
					//KmsKeyId:   aws.String(kmsKeyID),
					//Iops:       aws.Int64(volumeIOPS),
				},
			},
		},
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Arn: aws.String(iamInstanceProfileARN),
		},
	})
	if err != nil {
		return nil, err
	}

	if len(runResult.Instances) > 0 {
		l.Infof("Created instance: %v", *runResult.Instances[0].InstanceId)
		return runResult.Instances[0], nil
	}

	return nil, fmt.Errorf("no instances were created")
}
