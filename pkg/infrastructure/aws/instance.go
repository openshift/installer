package aws

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

var errInstanceNotCreated = errors.New("instance was not created")

const errInstanceIDNotFound = "InvalidInstanceID.NotFound"

type instanceInputOptions struct {
	infraID            string
	name               string
	amiID              string
	instanceType       string
	subnetID           string
	userData           string
	kmsKeyID           string
	iamRole            string
	instanceProfileARN string
	volumeType         string
	metadataAuth       string
	volumeSize         int64
	volumeIOPS         int64
	isEncrypted        bool
	associatePublicIP  bool
	securityGroupIds   []string
	targetGroupARNs    []string
	tags               map[string]string
}

func ensureInstance(ctx context.Context, logger logrus.FieldLogger, ec2Client ec2iface.EC2API, elbClient elbv2iface.ELBV2API, input *instanceInputOptions) (*ec2.Instance, error) {
	l := logger.WithField("name", input.name)
	filters := ec2Filters(input.infraID, input.name)
	createdOrFoundMsg := "Found existing instance"
	instance, err := existingInstance(ctx, ec2Client, filters)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, fmt.Errorf("failed to find instance: %w", err)
		}
		createdOrFoundMsg = "Created instance"
		instance, err = createInstance(ctx, ec2Client, input)
		if err != nil {
			return nil, err
		}
	}
	l = l.WithField("id", aws.StringValue(instance.InstanceId))
	l.Infoln(createdOrFoundMsg)

	// wait for the instance to get an IP address since we need to return it.
	err = wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (bool, error) {
			l.Debugln("Waiting for instance to be created and acquire IP address")
			res, err := ec2Client.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{
				InstanceIds: []*string{instance.InstanceId},
			})
			if err != nil {
				var awsErr awserr.Error
				if errors.As(err, &awsErr) && strings.EqualFold(awsErr.Code(), errInstanceIDNotFound) {
					return false, nil
				}
				return true, err
			}
			// Should not happen but let's be safe
			if len(res.Reservations) == 0 || len(res.Reservations[0].Instances) == 0 {
				return false, nil
			}
			instance = res.Reservations[0].Instances[0]
			if instance.PrivateIpAddress == nil {
				return false, nil
			}
			if input.associatePublicIP && instance.PublicIpAddress == nil {
				return false, nil
			}
			return true, nil
		},
	)
	if err != nil {
		l.WithError(err).Infoln("failed to wait for instance to acquire IP address")
		// We dont return an error here since the installation can still
		// proceed. However, we might not be able to gather instance logs in
		// case of bootstrap failure
	}

	for _, targetGroup := range input.targetGroupARNs {
		_, err = elbClient.RegisterTargetsWithContext(ctx, &elbv2.RegisterTargetsInput{
			TargetGroupArn: aws.String(targetGroup),
			Targets: []*elbv2.TargetDescription{
				{Id: instance.PrivateIpAddress},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to register target group (%s): %w", targetGroup, err)
		}
	}
	l.Infoln("Target groups registered")

	return instance, nil
}

func existingInstance(ctx context.Context, client ec2iface.EC2API, filters []*ec2.Filter) (*ec2.Instance, error) {
	res, err := client.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{Filters: filters})
	if err != nil {
		return nil, err
	}
	if len(res.Reservations) > 0 && len(res.Reservations[0].Instances) > 0 {
		return res.Reservations[0].Instances[0], nil
	}

	return nil, errNotFound
}

func createInstance(ctx context.Context, client ec2iface.EC2API, input *instanceInputOptions) (*ec2.Instance, error) {
	kmsKeyID := input.kmsKeyID
	if len(kmsKeyID) == 0 {
		kmsKey, err := client.GetEbsDefaultKmsKeyIdWithContext(ctx, &ec2.GetEbsDefaultKmsKeyIdInput{})
		if err != nil {
			return nil, fmt.Errorf("failed to get default KMS key: %w", err)
		}
		kmsKeyID = aws.StringValue(kmsKey.KmsKeyId)
	}

	// The iops parameter is not support if the volume type is gp2
	var iops *int64
	if input.volumeIOPS > 0 && input.volumeType != ec2.VolumeTypeGp2 {
		iops = aws.Int64(input.volumeIOPS)
	}

	tags := mergeTags(input.tags, map[string]string{"Name": input.name})
	volTags := mergeTags(input.tags, map[string]string{"Name": fmt.Sprintf("%s-vol", input.name)})
	httpTokens := input.metadataAuth
	if len(httpTokens) == 0 {
		httpTokens = "optional"
	}
	res, err := client.RunInstancesWithContext(ctx, &ec2.RunInstancesInput{
		ImageId:      aws.String(input.amiID),
		InstanceType: aws.String(input.instanceType),
		NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex:              aws.Int64(0),
				SubnetId:                 aws.String(input.subnetID),
				Groups:                   aws.StringSlice(input.securityGroupIds),
				AssociatePublicIpAddress: aws.Bool(input.associatePublicIP),
			},
		},
		MetadataOptions: &ec2.InstanceMetadataOptionsRequest{
			HttpEndpoint: aws.String("enabled"),
			HttpTokens:   aws.String(httpTokens),
		},
		UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(input.userData))),
		// InvalidParameterCombination: Network interfaces and an instance-level security groups may not be specified on the same request
		// SecurityGroupIds:  aws.StringSlice(options.securityGroupIDs),
		MinCount: aws.Int64(1),
		MaxCount: aws.Int64(1),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags:         ec2Tags(tags),
			},
			{
				ResourceType: aws.String("volume"),
				Tags:         ec2Tags(volTags),
			},
		},
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeType: aws.String(input.volumeType),
					VolumeSize: aws.Int64(input.volumeSize),
					Encrypted:  aws.Bool(input.isEncrypted),
					KmsKeyId:   aws.String(kmsKeyID),
					Iops:       iops,
				},
			},
		},
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Arn: aws.String(input.instanceProfileARN),
		},
	})
	if err != nil {
		return nil, err
	}
	if len(res.Instances) > 0 {
		return res.Instances[0], nil
	}

	return nil, errInstanceNotCreated
}

type instanceProfileOptions struct {
	namePrefix       string
	roleName         string
	assumeRolePolicy string
	policyDocument   string
	tags             map[string]string
}

func createInstanceProfile(ctx context.Context, logger logrus.FieldLogger, client iamiface.IAMAPI, input *instanceProfileOptions) (*iam.InstanceProfile, error) {
	useExistingRole := len(input.roleName) > 0
	roleName := input.roleName
	if !useExistingRole {
		roleName = fmt.Sprintf("%s-role", input.namePrefix)
		_, err := ensureRole(ctx, logger, client, roleName, input.assumeRolePolicy, input.tags)
		if err != nil {
			return nil, err
		}
	} else {
		logger.WithField("name", roleName).Infoln("Using user-supplied role")
	}

	profileName := fmt.Sprintf("%s-profile", input.namePrefix)
	profile, err := ensureProfile(ctx, logger, client, profileName, input.tags)
	if err != nil {
		return nil, err
	}

	hasRole := false
	for _, role := range profile.Roles {
		if aws.StringValue(role.RoleName) == roleName {
			hasRole = true
			break
		}
	}
	if !hasRole {
		_, err := client.AddRoleToInstanceProfileWithContext(ctx, &iam.AddRoleToInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			RoleName:            aws.String(roleName),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add role (%s) to instance profile (%s): %w", roleName, profileName, err)
		}
		logger.WithFields(logrus.Fields{
			"role":    roleName,
			"profile": profileName,
		}).Infoln("Added role to instance profile")
	} else {
		logger.WithFields(logrus.Fields{
			"role":    roleName,
			"profile": profileName,
		}).Infoln("Role already added to instance profile")
	}

	if !useExistingRole {
		rolePolicyName := fmt.Sprintf("%s-policy", profileName)
		err = existingRolePolicy(ctx, client, roleName, rolePolicyName)
		if err != nil {
			if !errors.Is(err, errNotFound) {
				return nil, fmt.Errorf("failed to get role policy: %w", err)
			}
			_, err = client.PutRolePolicyWithContext(ctx, &iam.PutRolePolicyInput{
				PolicyName:     aws.String(rolePolicyName),
				PolicyDocument: aws.String(input.policyDocument),
				RoleName:       aws.String(roleName),
			})
		}
		if err != nil {
			return nil, fmt.Errorf("failed to create role policy: %w", err)
		}
	}

	// We sleep here otherwise got an error when creating the ec2 instance
	// referencing the profile.
	time.Sleep(10 * time.Second)

	return profile, nil
}

func ensureRole(ctx context.Context, logger logrus.FieldLogger, client iamiface.IAMAPI, name string, assumeRolePolicy string, tags map[string]string) (*iam.Role, error) {
	createdOrFoundMsg := "Found existing role"
	role, err := existingRole(ctx, client, name)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, fmt.Errorf("failed to get existing role: %w", err)
		}
		createdOrFoundMsg = "Created role"
		role, err = createRole(ctx, client, name, assumeRolePolicy, tags)
		if err != nil {
			return nil, fmt.Errorf("failed to create role (%s): %w", name, err)
		}
	}
	logger.WithField("name", name).Infoln(createdOrFoundMsg)

	return role, nil
}

func existingRole(ctx context.Context, client iamiface.IAMAPI, name string) (*iam.Role, error) {
	res, err := client.GetRoleWithContext(ctx, &iam.GetRoleInput{
		RoleName: aws.String(name),
	})
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == iam.ErrCodeNoSuchEntityException {
			return nil, errNotFound
		}
		return nil, err
	}
	return res.Role, nil
}

func createRole(ctx context.Context, client iamiface.IAMAPI, name string, assumeRolePolicy string, tags map[string]string) (*iam.Role, error) {
	rtags := mergeTags(tags, map[string]string{"Name": name})
	res, err := client.CreateRoleWithContext(ctx, &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(name),
		Tags:                     iamTags(rtags),
	})
	if err != nil {
		return nil, err
	}
	return res.Role, nil
}

func ensureProfile(ctx context.Context, logger logrus.FieldLogger, client iamiface.IAMAPI, name string, tags map[string]string) (*iam.InstanceProfile, error) {
	createdOrFoundMsg := "Found existing instance profile"
	profile, err := existingProfile(ctx, client, name)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, fmt.Errorf("failed to get instance profile: %w", err)
		}
		createdOrFoundMsg = "Created instance profile"
		profile, err = createProfile(ctx, client, name, tags)
		if err != nil {
			return nil, fmt.Errorf("failed to create instance profile: %w", err)
		}
	}
	logger.WithField("name", name).Infoln(createdOrFoundMsg)

	return profile, nil
}

func existingProfile(ctx context.Context, client iamiface.IAMAPI, name string) (*iam.InstanceProfile, error) {
	res, err := client.GetInstanceProfileWithContext(ctx, &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	})
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == iam.ErrCodeNoSuchEntityException {
			return nil, errNotFound
		}
		return nil, err
	}
	return res.InstanceProfile, nil
}

func createProfile(ctx context.Context, client iamiface.IAMAPI, name string, tags map[string]string) (*iam.InstanceProfile, error) {
	ptags := mergeTags(tags, map[string]string{"Name": name})
	res, err := client.CreateInstanceProfileWithContext(ctx, &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(name),
		Path:                aws.String("/"),
		Tags:                iamTags(ptags),
	})
	if err != nil {
		return nil, err
	}

	waitCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var lastError error
	wait.UntilWithContext(
		waitCtx,
		func(ctx context.Context) {
			_, err := existingProfile(ctx, client, name)
			if err != nil {
				lastError = err
			} else {
				lastError = nil
				cancel()
			}
		},
		2*time.Second,
	)
	if err := waitCtx.Err(); err != nil {
		// Canceled error means that we found the instance profile
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("timed out waiting for instance profile to exist: %w", lastError)
		}
	}

	return res.InstanceProfile, nil
}

func existingRolePolicy(ctx context.Context, client iamiface.IAMAPI, roleName string, policyName string) error {
	res, err := client.GetRolePolicyWithContext(ctx, &iam.GetRolePolicyInput{
		RoleName:   aws.String(roleName),
		PolicyName: aws.String(policyName),
	})
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == iam.ErrCodeNoSuchEntityException {
			return errNotFound
		}
		return err
	}
	if aws.StringValue(res.PolicyName) != policyName {
		return errNotFound
	}

	return nil
}

func iamTags(tags map[string]string) []*iam.Tag {
	iamTags := make([]*iam.Tag, 0, len(tags))
	for k, v := range tags {
		k, v := k, v
		iamTags = append(iamTags, &iam.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return iamTags
}
