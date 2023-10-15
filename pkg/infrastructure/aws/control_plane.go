package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
)

type controlPlaneInput struct {
	clusterID         string
	userData          string
	amiID             string
	instanceType      string
	subnetIDs         []string
	securityGroupIDs  []string
	volumeType        string
	volumeSize        int64
	volumeIOPS        int64
	encrypted         bool
	kmsKeyID          string
	targetGroupARNs   []string
	replicas          int
	availabilityZones []string
	additionalTags    map[string]string
	zoneToSubnetIDMap map[string]string
}

func createControlPlaneResources(l *logrus.Logger, session *session.Session, controlPlaneInput *controlPlaneInput) error {
	instanceProfileARN, err := CreateControlPlaneInstanceProfile(l, session, controlPlaneInput.clusterID, controlPlaneInput.additionalTags)
	if err != nil {
		return fmt.Errorf("failed to create master instance profile: %w", err)
	}

	baseInstanceOptions := instanceOptions{
		amiID:                    controlPlaneInput.amiID,
		instanceType:             controlPlaneInput.instanceType,
		userData:                 controlPlaneInput.userData,
		securityGroupIDs:         controlPlaneInput.securityGroupIDs,
		volumeType:               controlPlaneInput.volumeType,
		volumeSize:               controlPlaneInput.volumeSize,
		volumeIOPS:               controlPlaneInput.volumeIOPS,
		encrypted:                controlPlaneInput.encrypted,
		kmsKeyID:                 controlPlaneInput.kmsKeyID,
		iamInstanceProfileARN:    instanceProfileARN,
		additionalEC2Tags:        ec2CreateTags(controlPlaneInput.additionalTags),
		associatePublicIPAddress: false,
	}

	ec2Client := ec2.New(session)
	for i := 0; i < controlPlaneInput.replicas; i++ {
		options := baseInstanceOptions
		options.name = fmt.Sprintf("%s-master-%d", controlPlaneInput.clusterID, i)
		// Choose appropriate subnet according to the zone
		// We assume that len(AZs) == len(masters)
		zoneIdx := i % len(controlPlaneInput.availabilityZones)
		options.subnetID = controlPlaneInput.zoneToSubnetIDMap[controlPlaneInput.availabilityZones[zoneIdx]]

		instance, err := createInstance(l, ec2Client, options)
		if err != nil {
			return fmt.Errorf("failed to create master %s: %w", options.name, err)
		}

		if err := RegisterTargetGroups(l, session, controlPlaneInput.targetGroupARNs, aws.StringValue(instance.PrivateIpAddress)); err != nil {
			return fmt.Errorf("failed to register master target groups: %w", err)
		}
	}

	return nil
}

func CreateControlPlaneInstanceProfile(l *logrus.Logger, session *session.Session, infraID string, additionalTags map[string]string) (string, error) {
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

	return createInstanceProfile(l, session, fmt.Sprintf("%s-master", infraID), assumeRolePolicy, policy, additionalTags)
}
