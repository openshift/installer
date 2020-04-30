// Package aws collects AWS-specific configuration.
package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	ccaws "github.com/openshift/cloud-credential-operator/pkg/aws"
)

// PermissionGroup is the group of permissions needed by cluster creation, operation, or teardown.
type PermissionGroup string

const (
	// PermissionCreateBase is a base set of permissions required in all installs where the installer creates resources.
	PermissionCreateBase PermissionGroup = "create-base"

	// PermissionDeleteBase is a base set of permissions required in all installs where the installer deletes resources.
	PermissionDeleteBase PermissionGroup = "delete-base"

	// PermissionCreateNetworking is an additional set of permissions required when the installer creates networking resources.
	PermissionCreateNetworking PermissionGroup = "create-networking"

	// PermissionDeleteNetworking is a set of permissions required when the installer destroys networking resources.
	PermissionDeleteNetworking PermissionGroup = "delete-networking"
)

var permissions = map[PermissionGroup][]string{
	// Base set of permissions required for cluster creation
	PermissionCreateBase: {
		// EC2 related perms
		"ec2:AllocateAddress",
		"ec2:AssociateAddress",
		"ec2:AuthorizeSecurityGroupEgress",
		"ec2:AuthorizeSecurityGroupIngress",
		"ec2:CopyImage",
		"ec2:CreateNetworkInterface",
		"ec2:AttachNetworkInterface",
		"ec2:CreateSecurityGroup",
		"ec2:CreateTags",
		"ec2:CreateVolume",
		"ec2:DeleteSecurityGroup",
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
		"ec2:DescribeNetworkInterfaces",
		"ec2:DescribePrefixLists",
		"ec2:DescribeRegions",
		"ec2:DescribeRouteTables",
		"ec2:DescribeSecurityGroups",
		"ec2:DescribeSubnets",
		"ec2:DescribeTags",
		"ec2:DescribeVolumes",
		"ec2:DescribeVpcAttribute",
		"ec2:DescribeVpcClassicLink",
		"ec2:DescribeVpcClassicLinkDnsSupport",
		"ec2:DescribeVpcEndpoints",
		"ec2:DescribeVpcs",
		"ec2:GetEbsDefaultKmsKeyId",
		"ec2:ModifyInstanceAttribute",
		"ec2:ModifyNetworkInterfaceAttribute",
		"ec2:ReleaseAddress",
		"ec2:RevokeSecurityGroupEgress",
		"ec2:RevokeSecurityGroupIngress",
		"ec2:RunInstances",
		"ec2:TerminateInstances",

		// ELB related perms
		"elasticloadbalancing:AddTags",
		"elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
		"elasticloadbalancing:AttachLoadBalancerToSubnets",
		"elasticloadbalancing:ConfigureHealthCheck",
		"elasticloadbalancing:CreateListener",
		"elasticloadbalancing:CreateLoadBalancer",
		"elasticloadbalancing:CreateLoadBalancerListeners",
		"elasticloadbalancing:CreateTargetGroup",
		"elasticloadbalancing:DeleteLoadBalancer",
		"elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
		"elasticloadbalancing:DeregisterTargets",
		"elasticloadbalancing:DescribeInstanceHealth",
		"elasticloadbalancing:DescribeListeners",
		"elasticloadbalancing:DescribeLoadBalancerAttributes",
		"elasticloadbalancing:DescribeLoadBalancers",
		"elasticloadbalancing:DescribeTags",
		"elasticloadbalancing:DescribeTargetGroupAttributes",
		"elasticloadbalancing:DescribeTargetHealth",
		"elasticloadbalancing:ModifyLoadBalancerAttributes",
		"elasticloadbalancing:ModifyTargetGroup",
		"elasticloadbalancing:ModifyTargetGroupAttributes",
		"elasticloadbalancing:RegisterInstancesWithLoadBalancer",
		"elasticloadbalancing:RegisterTargets",
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
		"route53:CreateHostedZone",
		"route53:DeleteHostedZone",
		"route53:GetChange",
		"route53:GetHostedZone",
		"route53:ListHostedZones",
		"route53:ListHostedZonesByName",
		"route53:ListResourceRecordSets",
		"route53:ListTagsForResource",
		"route53:UpdateHostedZoneComment",

		// S3 related perms
		"s3:CreateBucket",
		"s3:DeleteBucket",
		"s3:GetAccelerateConfiguration",
		"s3:GetBucketAcl",
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

		// More S3 (would be nice to limit 'Resource' to just the bucket we actually interact with...)
		"s3:DeleteObject",
		"s3:GetObject",
		"s3:GetObjectAcl",
		"s3:GetObjectTagging",
		"s3:GetObjectVersion",
		"s3:PutObject",
		"s3:PutObjectAcl",
		"s3:PutObjectTagging",
	},
	// Permissions required for deleting base cluster resources
	PermissionDeleteBase: {
		"autoscaling:DescribeAutoScalingGroups",
		"ec2:DeleteNetworkInterface",
		"ec2:DeleteVolume",
		"elasticloadbalancing:DeleteTargetGroup",
		"elasticloadbalancing:DescribeTargetGroups",
		"iam:DeleteAccessKey",
		"iam:DeleteUser",
		"iam:ListInstanceProfiles",
		"iam:ListRolePolicies",
		"iam:ListUserPolicies",
		"s3:DeleteObject",
		"s3:ListBucketVersions",
		"tag:GetResources",
	},
	// Permissions required for creating network resources
	PermissionCreateNetworking: {
		"ec2:AssociateDhcpOptions",
		"ec2:AssociateRouteTable",
		"ec2:AttachInternetGateway",
		"ec2:CreateDhcpOptions",
		"ec2:CreateInternetGateway",
		"ec2:CreateNatGateway",
		"ec2:CreateRoute",
		"ec2:CreateRouteTable",
		"ec2:CreateSubnet",
		"ec2:CreateVpc",
		"ec2:CreateVpcEndpoint",
		"ec2:ModifySubnetAttribute",
		"ec2:ModifyVpcAttribute",
	},
	// Permissions required for deleting network resources
	PermissionDeleteNetworking: {
		"ec2:DeleteDhcpOptions",
		"ec2:DeleteInternetGateway",
		"ec2:DeleteNatGateway",
		"ec2:DeleteRoute",
		"ec2:DeleteRouteTable",
		"ec2:DeleteSubnet",
		"ec2:DeleteVpc",
		"ec2:DeleteVpcEndpoints",
		"ec2:DetachInternetGateway",
		"ec2:DisassociateRouteTable",
		"ec2:ReplaceRouteTableAssociation",
	},
}

// ValidateCreds will try to create an AWS session, and also verify that the current credentials
// are sufficient to perform an installation, and that they can be used for cluster runtime
// as either capable of creating new credentials for components that interact with the cloud or
// being able to be passed through as-is to the components that need cloud credentials
func ValidateCreds(ssn *session.Session, groups []PermissionGroup, region string) error {
	// Compile a list of permissions based on the permission groups provided
	requiredPermissions := []string{}
	for _, group := range groups {
		groupPerms, ok := permissions[group]
		if !ok {
			return errors.Errorf("unable to access permissions group %s", group)
		}
		requiredPermissions = append(requiredPermissions, groupPerms...)
	}

	client, err := ccaws.NewClientFromIAMClient(iam.New(ssn))
	if err != nil {
		return errors.Wrap(err, "failed to create client for permission check")
	}

	sParams := &ccaws.SimulateParams{
		Region: region,
	}

	// Check whether we can do an installation
	logger := logrus.StandardLogger()
	canInstall, err := ccaws.CheckPermissionsAgainstActions(client, requiredPermissions, sParams, logger)
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
