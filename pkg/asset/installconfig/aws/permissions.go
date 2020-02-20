// Package aws collects AWS-specific configuration.
package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	ccaws "github.com/openshift/cloud-credential-operator/pkg/aws"
	"github.com/openshift/installer/pkg/version"
)

var installPermissions = []string{
	// EC2 related perms
	"ec2:AllocateAddress",
	"ec2:AssociateAddress",
	"ec2:AssociateDhcpOptions",
	"ec2:AssociateRouteTable",
	"ec2:AttachInternetGateway",
	"ec2:AuthorizeSecurityGroupEgress",
	"ec2:AuthorizeSecurityGroupIngress",
	"ec2:CopyImage",
	"ec2:CreateDhcpOptions",
	"ec2:CreateInternetGateway",
	"ec2:CreateNatGateway",
	"ec2:CreateNetworkInterface",
	"ec2:CreateRoute",
	"ec2:CreateRouteTable",
	"ec2:CreateSecurityGroup",
	"ec2:CreateSubnet",
	"ec2:CreateTags",
	"ec2:CreateVpc",
	"ec2:CreateVpcEndpoint",
	"ec2:CreateVolume",
	"ec2:DeleteSnapshot",
	"ec2:DeregisterImage",
	"ec2:DescribeAccountAttributes",
	"ec2:DescribeAddresses",
	"ec2:DescribeAvailabilityZones",
	"ec2:DescribeDhcpOptions",
	"ec2:DescribeImages",
	"ec2:DescribeInstanceAttribute",
	"ec2:DescribeInstanceCreditSpecifications",
	"ec2:DescribeInstances",
	"ec2:DescribeInternetGateways",
	"ec2:DescribeKeyPairs",
	"ec2:DescribeNatGateways",
	"ec2:DescribeNetworkAcls",
	"ec2:DescribePrefixLists",
	"ec2:DescribeRegions",
	"ec2:DescribeRouteTables",
	"ec2:DescribeSecurityGroups",
	"ec2:DescribeSubnets",
	"ec2:DescribeTags",
	"ec2:DescribeVpcEndpoints",
	"ec2:DescribeVpcs",
	"ec2:DescribeVpcAttribute",
	"ec2:DescribeVolumes",
	"ec2:DescribeVpcClassicLink",
	"ec2:DescribeVpcClassicLinkDnsSupport",
	"ec2:ModifyInstanceAttribute",
	"ec2:ModifySubnetAttribute",
	"ec2:ModifyVpcAttribute",
	"ec2:RevokeSecurityGroupEgress",
	"ec2:RunInstances",
	"ec2:TerminateInstances",
	"ec2:DeleteDhcpOptions",
	"ec2:DeleteRoute",
	"ec2:RevokeSecurityGroupIngress",
	"ec2:DisassociateRouteTable",
	"ec2:ReplaceRouteTableAssociation",
	"ec2:DeleteRouteTable",
	"ec2:DeleteSubnet",
	"ec2:DescribeNetworkInterfaces",
	"ec2:ModifyNetworkInterfaceAttribute",
	"ec2:DeleteNatGateway",
	"ec2:DeleteSecurityGroup",
	"ec2:DetachInternetGateway",
	"ec2:DeleteInternetGateway",
	"ec2:ReleaseAddress",
	"ec2:DeleteVpc",

	// ELB related perms
	"elasticloadbalancing:AddTags",
	"elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
	"elasticloadbalancing:AttachLoadBalancerToSubnets",
	"elasticloadbalancing:CreateListener",
	"elasticloadbalancing:CreateLoadBalancer",
	"elasticloadbalancing:CreateLoadBalancerListeners",
	"elasticloadbalancing:CreateTargetGroup",
	"elasticloadbalancing:ConfigureHealthCheck",
	"elasticloadbalancing:DeleteLoadBalancer",
	"elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
	"elasticloadbalancing:DeregisterTargets",
	"elasticloadbalancing:DescribeInstanceHealth",
	"elasticloadbalancing:DescribeListeners",
	"elasticloadbalancing:DescribeLoadBalancers",
	"elasticloadbalancing:DescribeLoadBalancerAttributes",
	"elasticloadbalancing:DescribeTags",
	"elasticloadbalancing:DescribeTargetGroupAttributes",
	"elasticloadbalancing:DescribeTargetHealth",
	"elasticloadbalancing:ModifyLoadBalancerAttributes",
	"elasticloadbalancing:ModifyTargetGroup",
	"elasticloadbalancing:ModifyTargetGroupAttributes",
	"elasticloadbalancing:RegisterTargets",
	"elasticloadbalancing:RegisterInstancesWithLoadBalancer",
	"elasticloadbalancing:SetLoadBalancerPoliciesOfListener",

	// IAM related perms
	"iam:AddRoleToInstanceProfile",
	"iam:CreateInstanceProfile",
	"iam:CreateRole",
	"iam:DeleteInstanceProfile",
	"iam:DeleteRole",
	"iam:DeleteRolePolicy",
	"iam:GetInstanceProfile",
	"iam:GetRole",
	"iam:GetRolePolicy",
	"iam:GetUser",
	"iam:ListInstanceProfilesForRole",
	"iam:ListRoles",
	"iam:ListUsers",
	"iam:PassRole",
	"iam:PutRolePolicy",
	"iam:RemoveRoleFromInstanceProfile",
	"iam:SimulatePrincipalPolicy",
	"iam:TagRole",

	// Route53 related perms
	"route53:ChangeResourceRecordSets",
	"route53:ChangeTagsForResource",
	"route53:GetChange",
	"route53:GetHostedZone",
	"route53:CreateHostedZone",
	"route53:DeleteHostedZone",
	"route53:ListHostedZones",
	"route53:ListHostedZonesByName",
	"route53:ListResourceRecordSets",
	"route53:ListTagsForResource",
	"route53:UpdateHostedZoneComment",

	// S3 related perms
	"s3:CreateBucket",
	"s3:DeleteBucket",
	"s3:GetAccelerateConfiguration",
	"s3:GetBucketCors",
	"s3:GetBucketLocation",
	"s3:GetBucketLogging",
	"s3:GetBucketObjectLockConfiguration",
	"s3:GetBucketReplication",
	"s3:GetBucketRequestPayment",
	"s3:GetBucketTagging",
	"s3:GetBucketVersioning",
	"s3:GetBucketWebsite",
	"s3:GetEncryptionConfiguration",
	"s3:GetLifecycleConfiguration",
	"s3:GetReplicationConfiguration",
	"s3:ListBucket",
	"s3:PutBucketAcl",
	"s3:PutBucketTagging",
	"s3:PutEncryptionConfiguration",

	// More S3 (would be nice to limit 'Resource' to just the bucket we actualy interact with...)
	"s3:PutObject",
	"s3:PutObjectAcl",
	"s3:PutObjectTagging",
	"s3:GetObject",
	"s3:GetObjectAcl",
	"s3:GetObjectTagging",
	"s3:GetObjectVersion",
	"s3:DeleteObject",

	// Uninstall-specific perms
	"autoscaling:DescribeAutoScalingGroups",
	"ec2:DeleteNetworkInterface",
	"ec2:DeleteVolume",
	"ec2:DeleteVpcEndpoints",
	"elasticloadbalancing:DescribeTargetGroups",
	"elasticloadbalancing:DeleteTargetGroup",
	"iam:ListInstanceProfiles",
	"iam:ListRolePolicies",
	"iam:ListUserPolicies",
	"tag:GetResources",
}

// ValidateCreds will try to create an AWS session, and also verify that the current credentials
// are sufficient to perform an installation, and that they can be used for cluster runtime
// as either capable of creating new credentials for components that interact with the cloud or
// being able to be passed through as-is to the components that need cloud credentials
func ValidateCreds(ssn *session.Session, region string) error {
	creds, err := ssn.Config.Credentials.Get()
	if err != nil {
		return errors.Wrap(err, "getting creds from session")
	}

	client, err := ccaws.NewClient([]byte(creds.AccessKeyID), []byte(creds.SecretAccessKey), fmt.Sprintf("OpenShift/4.x Installer/%s", version.Raw))
	if err != nil {
		return errors.Wrap(err, "initialize cloud-credentials client")
	}

	sParams := &ccaws.SimulateParams{
		Region: region,
	}

	// Check whether we can do an installation
	logger := logrus.StandardLogger()
	canInstall, err := ccaws.CheckPermissionsAgainstActions(client, installPermissions, sParams, logger)
	if err != nil {
		return errors.Wrap(err, "checking install permissions")
	}
	if !canInstall {
		return errors.New("current credentials insufficient for performing cluster installation")
	}

	// Check whether we can mint new creds for cluster services needing to interact with the cloud
	canMint, err := ccaws.CheckCloudCredCreation(client, logger)
	if err != nil {
		return errors.Wrap(err, "mint credentials check")
	}
	if canMint {
		return nil
	}

	// Check whether we can use the current credentials in passthrough mode to satisfy
	// cluster services needing to interact with the cloud
	canPassthrough, err := ccaws.CheckCloudCredPassthrough(client, sParams, logger)
	if err != nil {
		return errors.Wrap(err, "passthrough credentials check")
	}
	if canPassthrough {
		return nil
	}

	return errors.New("AWS credentials cannot be used to either create new creds or use as-is")
}
