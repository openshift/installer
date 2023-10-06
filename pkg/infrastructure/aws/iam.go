package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/sirupsen/logrus"
)

type iamInput struct {
	clusterID string
}

func createIAMResources(l *logrus.Logger, session *session.Session, iamInput *iamInput) error {
	iamClient := iam.New(session)
	_, err := CreateComputeInstanceProfile(l, iamClient, iamInput.clusterID)
	if err != nil {
		return err
	}

	return nil
}

func CreateComputeInstanceProfile(l *logrus.Logger, client iamiface.IAMAPI, infraID string) (*string, error) {
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
		"ec2:DescribeInstances",
		"ec2:DescribeRegions"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}`
	)

	profileName := fmt.Sprintf("%s-worker-profile", infraID)
	roleName := fmt.Sprintf("%s-worker-role", infraID)
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
	// time.Sleep(10 * time.Second)
	return instanceProfile.Arn, nil
}
