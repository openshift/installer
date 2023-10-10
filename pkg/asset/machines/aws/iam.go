package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/pkg/errors"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
)

var (
	policies = map[string]*iamv1.PolicyDocument{
		"master": {
			Version: "2012-10-17",
			Statement: []iamv1.StatementEntry{
				{
					Effect: "Allow",
					Action: []string{
						"autoscaling:DescribeAutoScalingGroups",
						"autoscaling:DescribeInstanceRefreshes",
						"ec2:AllocateAddress",
						"ec2:AssignIpv6Addresses",
						"ec2:AssignPrivateIpAddresses",
						"ec2:AssociateRouteTable",
						"ec2:AttachInternetGateway",
						"ec2:AttachNetworkInterface",
						"ec2:AttachVolume",
						"ec2:AuthorizeSecurityGroupIngress",
						"ec2:CreateEgressOnlyInternetGateway",
						"ec2:CreateInternetGateway",
						"ec2:CreateLaunchTemplate",
						"ec2:CreateLaunchTemplateVersion",
						"ec2:CreateNatGateway",
						"ec2:CreateNetworkInterface",
						"ec2:CreateRoute",
						"ec2:CreateRouteTable",
						"ec2:CreateSecurityGroup",
						"ec2:CreateSubnet",
						"ec2:CreateTags",
						"ec2:CreateVolume",
						"ec2:CreateVpc",
						"ec2:DeleteEgressOnlyInternetGateway",
						"ec2:DeleteInternetGateway",
						"ec2:DeleteLaunchTemplate",
						"ec2:DeleteLaunchTemplateVersions",
						"ec2:DeleteNatGateway",
						"ec2:DeleteRouteTable",
						"ec2:DeleteSecurityGroup",
						"ec2:DeleteSubnet",
						"ec2:DeleteTags",
						"ec2:DeleteVolume",
						"ec2:DeleteVpc",
						"ec2:Describe*",
						"ec2:DescribeAccountAttributes",
						"ec2:DescribeAddresses",
						"ec2:DescribeAvailabilityZones",
						"ec2:DescribeEgressOnlyInternetGateways",
						"ec2:DescribeImages",
						"ec2:DescribeInstances",
						"ec2:DescribeInstanceTypes",
						"ec2:DescribeInternetGateways",
						"ec2:DescribeKeyPairs",
						"ec2:DescribeLaunchTemplates",
						"ec2:DescribeLaunchTemplateVersions",
						"ec2:DescribeNatGateways",
						"ec2:DescribeNetworkInterfaceAttribute",
						"ec2:DescribeNetworkInterfaces",
						"ec2:DescribeRouteTables",
						"ec2:DescribeSecurityGroups",
						"ec2:DescribeSubnets",
						"ec2:DescribeTags",
						"ec2:DescribeVolumes",
						"ec2:DescribeVpcAttribute",
						"ec2:DescribeVpcs",
						"ec2:DetachInternetGateway",
						"ec2:DetachNetworkInterface",
						"ec2:DetachVolume",
						"ec2:DisassociateAddress",
						"ec2:DisassociateRouteTable",
						"ec2:ModifyInstanceAttribute",
						"ec2:ModifyInstanceMetadataOptions",
						"ec2:ModifyNetworkInterfaceAttribute",
						"ec2:ModifySubnetAttribute",
						"ec2:ModifyVolume",
						"ec2:ModifyVpcAttribute",
						"ec2:ReleaseAddress",
						"ec2:ReplaceRoute",
						"ec2:RevokeSecurityGroupIngress",
						"ec2:RunInstances",
						"ec2:TerminateInstances",
						"ec2:UnassignPrivateIpAddresses",
						"elasticloadbalancing:AddTags",
						"elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
						"elasticloadbalancing:AttachLoadBalancerToSubnets",
						"elasticloadbalancing:ConfigureHealthCheck",
						"elasticloadbalancing:CreateListener",
						"elasticloadbalancing:CreateLoadBalancer",
						"elasticloadbalancing:CreateLoadBalancerListeners",
						"elasticloadbalancing:CreateLoadBalancerPolicy",
						"elasticloadbalancing:CreateTargetGroup",
						"elasticloadbalancing:DeleteListener",
						"elasticloadbalancing:DeleteLoadBalancer",
						"elasticloadbalancing:DeleteLoadBalancerListeners",
						"elasticloadbalancing:DeleteTargetGroup",
						"elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
						"elasticloadbalancing:DeregisterTargets",
						"elasticloadbalancing:Describe*",
						"elasticloadbalancing:DescribeListeners",
						"elasticloadbalancing:DescribeLoadBalancerAttributes",
						"elasticloadbalancing:DescribeLoadBalancers",
						"elasticloadbalancing:DescribeTags",
						"elasticloadbalancing:DescribeTargetGroups",
						"elasticloadbalancing:DescribeTargetHealth",
						"elasticloadbalancing:DetachLoadBalancerFromSubnets",
						"elasticloadbalancing:ModifyListener",
						"elasticloadbalancing:ModifyLoadBalancerAttributes",
						"elasticloadbalancing:ModifyTargetGroup",
						"elasticloadbalancing:ModifyTargetGroupAttributes",
						"elasticloadbalancing:RegisterInstancesWithLoadBalancer",
						"elasticloadbalancing:RegisterTargets",
						"elasticloadbalancing:RemoveTags",
						"elasticloadbalancing:SetLoadBalancerPoliciesForBackendServer",
						"elasticloadbalancing:SetLoadBalancerPoliciesOfListener",
						"elasticloadbalancing:SetSubnets",
						"kms:DescribeKey",
						"tag:GetResources",
					},
					Resource: iamv1.Resources{
						"*",
					},
				},
				{
					Effect: "Allow",
					Action: []string{
						"autoscaling:CreateAutoScalingGroup",
						"autoscaling:UpdateAutoScalingGroup",
						"autoscaling:CreateOrUpdateTags",
						"autoscaling:StartInstanceRefresh",
						"autoscaling:DeleteAutoScalingGroup",
						"autoscaling:DeleteTags",
					},
					Resource: iamv1.Resources{
						"arn:*:autoscaling:*:*:autoScalingGroup:*:autoScalingGroupName/*",
					},
				},
				{
					Effect: "Allow",
					Action: []string{
						"iam:CreateServiceLinkedRole",
					},
					Resource: iamv1.Resources{
						"arn:*:iam::*:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling",
					},
					Condition: iamv1.Conditions{
						"StringLike": map[string]string{
							"iam:AWSServiceName": "autoscaling.amazonaws.com",
						},
					},
				},
				{
					Effect: "Allow",
					Action: []string{
						"iam:CreateServiceLinkedRole",
					},
					Resource: iamv1.Resources{
						"arn:*:iam::*:role/aws-service-role/elasticloadbalancing.amazonaws.com/AWSServiceRoleForElasticLoadBalancing",
					},
					Condition: iamv1.Conditions{
						"StringLike": map[string]string{
							"iam:AWSServiceName": "elasticloadbalancing.amazonaws.com",
						},
					},
				},
				{
					Effect: "Allow",
					Action: []string{
						"iam:CreateServiceLinkedRole",
					},
					Resource: iamv1.Resources{
						"arn:*:iam::*:role/aws-service-role/spot.amazonaws.com/AWSServiceRoleForEC2Spot",
					},
					Condition: iamv1.Conditions{
						"StringLike": map[string]string{
							"iam:AWSServiceName": "spot.amazonaws.com",
						},
					},
				},
				{
					Effect: "Allow",
					Action: []string{
						"s3:CreateBucket",
						"s3:DeleteBucket",
						"s3:PutObject",
						"s3:DeleteObject",
						"s3:PutBucketPolicy",
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

// This is here for now, not sure where it should really go.
func PutIAMRoles(clusterID string, ic *installconfig.InstallConfig) error {
	// Create the IAM Role with the aws sdk.
	// https://docs.aws.amazon.com/sdk-for-go/api/service/iam/#IAM.CreateRole
	session, err := ic.AWS.Session(context.TODO())
	if err != nil {
		return errors.Wrap(err, "failed to load AWS session")
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
					"Service": []string{
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
		return errors.Wrap(err, "failed to marshal assume policy")
	}
	createRoleInput := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(string(assumePolicyBytes)),
		Tags:                     []*iam.Tag{clusterOwnedIAMTag},
	}

	for _, role := range []string{"master", "worker"} {
		roleName := aws.String(fmt.Sprintf("%s-%s-role", clusterID, role))
		profileName := aws.String(fmt.Sprintf("%s-%s-profile", clusterID, role))
		policyName := aws.String(fmt.Sprintf("%s-%s-policy", clusterID, role))

		createRoleInput.RoleName = roleName
		if _, err := svc.GetRole(&iam.GetRoleInput{RoleName: roleName}); err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() != iam.ErrCodeNoSuchEntityException {
				return errors.Wrapf(err, "failed to get %s role", role)
			}
			// If the role does not exist, create it.
			if _, err := svc.CreateRole(createRoleInput); err != nil {
				return errors.Wrapf(err, "failed to create %s role", role)
			}
			time.Sleep(10 * time.Second)
			if err := svc.WaitUntilRoleExists(&iam.GetRoleInput{RoleName: roleName}); err != nil {
				return errors.Wrapf(err, "failed to wait for %s role to exist", role)
			}
		}

		// Put the policy inline.
		b, err := json.Marshal(policies[role])
		if err != nil {
			return errors.Wrapf(err, "failed to marshal %s policy", role)
		}
		if _, err := svc.PutRolePolicy(&iam.PutRolePolicyInput{
			PolicyDocument: aws.String(string(b)),
			PolicyName:     policyName,
			RoleName:       roleName,
		}); err != nil {
			return errors.Wrapf(err, "failed to create inline policy for role %s ", role)
		}

		if _, err := svc.GetInstanceProfile(&iam.GetInstanceProfileInput{InstanceProfileName: roleName}); err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() != iam.ErrCodeNoSuchEntityException {
				return errors.Wrapf(err, "failed to get %s instance profile", role)
			}
			// If the profile does not exist, create it.
			if _, err := svc.CreateInstanceProfile(&iam.CreateInstanceProfileInput{
				InstanceProfileName: profileName,
				Tags:                []*iam.Tag{clusterOwnedIAMTag},
			}); err != nil {
				return errors.Wrapf(err, "failed to create %s instance profile", role)
			}
			time.Sleep(10 * time.Second)
			if err := svc.WaitUntilInstanceProfileExists(&iam.GetInstanceProfileInput{InstanceProfileName: profileName}); err != nil {
				return errors.Wrapf(err, "failed to wait for %s role to exist", role)
			}

			// Finally, attach the role to the profile.
			if _, err := svc.AddRoleToInstanceProfile(&iam.AddRoleToInstanceProfileInput{
				InstanceProfileName: profileName,
				RoleName:            roleName,
			}); err != nil {
				return errors.Wrapf(err, "failed to add %s role to instance profile", role)
			}
		}
	}

	return nil
}
