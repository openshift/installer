package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/sirupsen/logrus"
)

type controlPlaneInput struct {
	clusterID       string
	userData        string
	amiID           string
	instanceType    string
	subnetIDs       []string
	securityGroupID string
	volumeType      string
	volumeSize      int64
	volumeIOPS      int64
	encrypted       bool
	kmsKeyID        string
	targetGroupARNs []string
	replicas        int
}

func createControlPlaneResources(l *logrus.Logger, session *session.Session, controlPlaneInput *controlPlaneInput) error {
	iamClient := iam.New(session)
	instanceProfileARN, err := CreateControlPlaneInstanceProfile(l, iamClient, controlPlaneInput.clusterID)
	if err != nil {
		return err
	}

	baseInstanceOptions := instanceOptions{
		amiID:                    controlPlaneInput.amiID,
		instanceType:             controlPlaneInput.instanceType,
		userData:                 controlPlaneInput.userData,
		securityGroupIDs:         []string{controlPlaneInput.securityGroupID},
		associatePublicIPAddress: false,
		volumeType:               controlPlaneInput.volumeType,
		volumeSize:               controlPlaneInput.volumeSize,
		volumeIOPS:               controlPlaneInput.volumeIOPS,
		encrypted:                controlPlaneInput.encrypted,
		kmsKeyID:                 controlPlaneInput.kmsKeyID,
		iamInstanceProfileARN:    *instanceProfileARN,
	}

	for i := 0; i < controlPlaneInput.replicas; i++ {
		ec2Client := ec2.New(session)

		options := baseInstanceOptions
		options.name = fmt.Sprintf("master-%d", i)
		options.subnetID = controlPlaneInput.subnetIDs[i]

		instance, err := createEC2Instance(l, ec2Client, controlPlaneInput.clusterID, options)
		if err != nil {
			return err
		}

		elbClient := elbv2.New(session)
		for _, targetGroupARN := range controlPlaneInput.targetGroupARNs {
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
	}

	return nil
}

func CreateControlPlaneInstanceProfile(l *logrus.Logger, client iamiface.IAMAPI, infraID string) (*string, error) {
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
		policy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:AttachVolume",
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:CreateSecurityGroup",
        "ec2:CreateTags",
        "ec2:CreateVolume",
        "ec2:DeleteSecurityGroup",
        "ec2:DeleteVolume",
        "ec2:Describe*",
        "ec2:DetachVolume",
        "ec2:ModifyInstanceAttribute",
        "ec2:ModifyVolume",
        "ec2:RevokeSecurityGroupIngress",
        "elasticloadbalancing:AddTags",
        "elasticloadbalancing:AttachLoadBalancerToSubnets",
        "elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
        "elasticloadbalancing:CreateListener",
        "elasticloadbalancing:CreateLoadBalancer",
        "elasticloadbalancing:CreateLoadBalancerPolicy",
        "elasticloadbalancing:CreateLoadBalancerListeners",
        "elasticloadbalancing:CreateTargetGroup",
        "elasticloadbalancing:ConfigureHealthCheck",
        "elasticloadbalancing:DeleteListener",
        "elasticloadbalancing:DeleteLoadBalancer",
        "elasticloadbalancing:DeleteLoadBalancerListeners",
        "elasticloadbalancing:DeleteTargetGroup",
        "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
        "elasticloadbalancing:DeregisterTargets",
        "elasticloadbalancing:Describe*",
        "elasticloadbalancing:DetachLoadBalancerFromSubnets",
        "elasticloadbalancing:ModifyListener",
        "elasticloadbalancing:ModifyLoadBalancerAttributes",
        "elasticloadbalancing:ModifyTargetGroup",
        "elasticloadbalancing:ModifyTargetGroupAttributes",
        "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
        "elasticloadbalancing:RegisterTargets",
        "elasticloadbalancing:SetLoadBalancerPoliciesForBackendServer",
        "elasticloadbalancing:SetLoadBalancerPoliciesOfListener",
        "kms:DescribeKey"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}`
	)

	profileName := fmt.Sprintf("%s-master-profile", infraID)
	roleName := fmt.Sprintf("%s-master-role", infraID)
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
			PolicyDocument: aws.String(policy),
			RoleName:       aws.String(roleName),
		})
		if err != nil {
			return nil, fmt.Errorf("cannot create profile policy: %w", err)
		}
		l.Info("Created role policy", "name", rolePolicyName)
	}

	// We sleep here otherwise got an error when creating the ec2 instance referencing the profile.
	time.Sleep(10 * time.Second)
	return instanceProfile.Arn, nil
}
