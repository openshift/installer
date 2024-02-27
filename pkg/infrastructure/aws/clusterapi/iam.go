package clusterapi

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"

	"github.com/openshift/installer/pkg/asset/installconfig"
)

var (
	policies = map[string]*iamv1.PolicyDocument{
		"master": {
			Version: "2012-10-17",
			Statement: []iamv1.StatementEntry{
				{
					Effect: "Allow",
					Action: []string{
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
						"kms:DescribeKey",
					},
					Resource: iamv1.Resources{
						"*",
					},
				},
			},
		},
		"worker": {
			Version: "2012-10-17",
			Statement: []iamv1.StatementEntry{
				{
					Effect: "Allow",
					Action: iamv1.Actions{
						"ec2:DescribeInstances",
						"ec2:DescribeRegions",
					},
					Resource: iamv1.Resources{"*"},
				},
			},
		},
	}
)

// putIAMRoles creates the roles used by control-plane and compute nodes.
func putIAMRoles(clusterID string, ic *installconfig.InstallConfig) error {
	logrus.Infoln("Creating IAM roles for control-plane and compute nodes")
	// Create the IAM Role with the aws sdk.
	// https://docs.aws.amazon.com/sdk-for-go/api/service/iam/#IAM.CreateRole
	session, err := ic.AWS.Session(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load AWS session: %w", err)
	}
	svc := iam.New(session)

	// Create the IAM Roles for master and workers.
	clusterOwnedIAMTag := &iam.Tag{
		Key:   aws.String(fmt.Sprintf("kubernetes.io/cluster/%s", clusterID)),
		Value: aws.String("owned"),
	}
	assumePolicy := &iamv1.PolicyDocument{
		Version: "2012-10-17",
		Statement: iamv1.Statements{
			{
				Effect: "Allow",
				Principal: iamv1.Principals{
					iamv1.PrincipalService: []string{
						"ec2.amazonaws.com",
					},
				},
				Action: iamv1.Actions{
					"sts:AssumeRole",
				},
			},
		},
	}
	assumePolicyBytes, err := json.Marshal(assumePolicy)
	if err != nil {
		return fmt.Errorf("failed to marshal assume policy: %w", err)
	}

	for _, role := range []string{"master", "worker"} {
		roleName := aws.String(fmt.Sprintf("%s-%s-role", clusterID, role))
		if _, err := svc.GetRole(&iam.GetRoleInput{RoleName: roleName}); err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() != iam.ErrCodeNoSuchEntityException {
				return fmt.Errorf("failed to get %s role: %w", role, err)
			}
			// If the role does not exist, create it.
			createRoleInput := &iam.CreateRoleInput{
				RoleName:                 roleName,
				AssumeRolePolicyDocument: aws.String(string(assumePolicyBytes)),
				Tags:                     []*iam.Tag{clusterOwnedIAMTag},
			}
			if _, err := svc.CreateRole(createRoleInput); err != nil {
				return fmt.Errorf("failed to create %s role: %w", role, err)
			}
			time.Sleep(10 * time.Second)
			if err := svc.WaitUntilRoleExists(&iam.GetRoleInput{RoleName: roleName}); err != nil {
				return fmt.Errorf("failed to wait for %s role to exist: %w", role, err)
			}
		}

		// Put the policy inline.
		policyName := aws.String(fmt.Sprintf("%s-%s-policy", clusterID, role))
		b, err := json.Marshal(policies[role])
		if err != nil {
			return fmt.Errorf("failed to marshal %s policy: %w", role, err)
		}
		if _, err := svc.PutRolePolicy(&iam.PutRolePolicyInput{
			PolicyDocument: aws.String(string(b)),
			PolicyName:     policyName,
			RoleName:       roleName,
		}); err != nil {
			return fmt.Errorf("failed to create inline policy for role %s: %w", role, err)
		}

		profileName := aws.String(fmt.Sprintf("%s-%s-profile", clusterID, role))
		if _, err := svc.GetInstanceProfile(&iam.GetInstanceProfileInput{InstanceProfileName: profileName}); err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() != iam.ErrCodeNoSuchEntityException {
				return fmt.Errorf("failed to get %s instance profile: %w", role, err)
			}
			// If the profile does not exist, create it.
			if _, err := svc.CreateInstanceProfile(&iam.CreateInstanceProfileInput{
				InstanceProfileName: profileName,
				Tags:                []*iam.Tag{clusterOwnedIAMTag},
			}); err != nil {
				return fmt.Errorf("failed to create %s instance profile: %w", role, err)
			}
			time.Sleep(10 * time.Second)
			if err := svc.WaitUntilInstanceProfileExists(&iam.GetInstanceProfileInput{InstanceProfileName: profileName}); err != nil {
				return fmt.Errorf("failed to wait for %s role to exist: %w", role, err)
			}

			// Finally, attach the role to the profile.
			if _, err := svc.AddRoleToInstanceProfile(&iam.AddRoleToInstanceProfileInput{
				InstanceProfileName: profileName,
				RoleName:            roleName,
			}); err != nil {
				return fmt.Errorf("failed to add %s role to instance profile: %w", role, err)
			}
		}
	}

	return nil
}
