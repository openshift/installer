package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/sirupsen/logrus"
)

type controlPlaneInputOptions struct {
	instanceInputOptions
	nReplicas         int
	privateSubnetIDs  []string
	zoneToSubnetMap   map[string]string
	availabilityZones []string
}

type controlPlaneOutput struct {
	controlPlaneIPs []string
}

func createControlPlaneResources(ctx context.Context, logger logrus.FieldLogger, ec2Client ec2iface.EC2API, iamClient iamiface.IAMAPI, elbClient elbv2iface.ELBV2API, input *controlPlaneInputOptions) (*controlPlaneOutput, error) {
	profileName := fmt.Sprintf("%s-master", input.infraID)
	instanceProfile, err := createControlPlaneInstanceProfile(ctx, logger, iamClient, profileName, input.iamRole, input.tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create control plane instance profile: %w", err)
	}

	instanceIPs := make([]string, 0, input.nReplicas)
	for i := 0; i < input.nReplicas; i++ {
		options := input.instanceInputOptions
		options.name = fmt.Sprintf("%s-master-%d", input.infraID, i)
		// Choose appropriate subnet according to zone
		zoneIdx := i % len(input.availabilityZones)
		options.subnetID = input.zoneToSubnetMap[input.availabilityZones[zoneIdx]]
		options.instanceProfileARN = aws.StringValue(instanceProfile.Arn)

		instance, err := ensureInstance(ctx, logger, ec2Client, elbClient, &options)
		if err != nil {
			return nil, fmt.Errorf("failed to create control plane (%s): %w", options.name, err)
		}
		instanceIPs = append(instanceIPs, aws.StringValue(instance.PrivateIpAddress))
	}
	logger.Infoln("Created control plane instances")

	return &controlPlaneOutput{controlPlaneIPs: instanceIPs}, nil
}

func createControlPlaneInstanceProfile(ctx context.Context, logger logrus.FieldLogger, client iamiface.IAMAPI, name string, roleName string, tags map[string]string) (*iam.InstanceProfile, error) {
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

	profileInput := &instanceProfileOptions{
		namePrefix:       name,
		roleName:         roleName,
		assumeRolePolicy: assumeRolePolicy,
		policyDocument:   policy,
		tags:             tags,
	}
	return createInstanceProfile(ctx, logger, client, profileInput)
}
