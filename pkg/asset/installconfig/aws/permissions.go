// Package aws collects AWS-specific configuration.
package aws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	ccaws "github.com/openshift/cloud-credential-operator/pkg/aws"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
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

	// PermissionDeleteSharedNetworking is a set of permissions required when the installer destroys resources from a shared-network cluster.
	PermissionDeleteSharedNetworking PermissionGroup = "delete-shared-networking"

	// PermissionCreateInstanceRole is a set of permissions required when the installer creates instance roles.
	PermissionCreateInstanceRole PermissionGroup = "create-instance-role"

	// PermissionDeleteSharedInstanceRole is a set of permissions required when the installer destroys resources from a
	// cluster with user-supplied IAM roles for instances.
	PermissionDeleteSharedInstanceRole PermissionGroup = "delete-shared-instance-role"

	// PermissionCreateInstanceProfile is a set of permission required when the installer creates instance profiles.
	PermissionCreateInstanceProfile PermissionGroup = "create-instance-profile"

	// PermissionDeleteSharedInstanceProfile is a set of permissions required when the installer destroys resources from
	// a cluster with user-supplied IAM instance profiles for instances.
	PermissionDeleteSharedInstanceProfile PermissionGroup = "delete-shared-instance-profile"

	// PermissionCreateHostedZone is a set of permissions required when the installer creates a route53 hosted zone.
	PermissionCreateHostedZone PermissionGroup = "create-hosted-zone"

	// PermissionDeleteHostedZone is a set of permissions required when the installer destroys a route53 hosted zone.
	PermissionDeleteHostedZone PermissionGroup = "delete-hosted-zone"

	// PermissionKMSEncryptionKeys is an additional set of permissions required when the installer uses user provided kms encryption keys.
	PermissionKMSEncryptionKeys PermissionGroup = "kms-encryption-keys"

	// PermissionPublicIpv4Pool is an additional set of permissions required when the installer uses public IPv4 pools.
	PermissionPublicIpv4Pool PermissionGroup = "public-ipv4-pool"

	// PermissionDeleteIgnitionObjects is a permission set required when `preserveBootstrapIgnition` is not set.
	PermissionDeleteIgnitionObjects PermissionGroup = "delete-ignition-objects"
)

var permissions = map[PermissionGroup][]string{
	// Base set of permissions required for cluster creation
	PermissionCreateBase: {
		// EC2 related perms
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
		"ec2:DescribeSecurityGroupRules",
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
		"elasticloadbalancing:SetSecurityGroups",

		// IAM related perms
		"iam:GetInstanceProfile",
		"iam:GetRole",
		"iam:GetRolePolicy",
		"iam:GetUser",
		"iam:ListInstanceProfilesForRole",
		"iam:ListRoles",
		"iam:ListUsers",
		"iam:PassRole",
		"iam:SimulatePrincipalPolicy",
		"iam:TagInstanceProfile",
		"iam:TagRole",

		// Route53 related perms
		"route53:ChangeResourceRecordSets",
		"route53:ChangeTagsForResource",
		"route53:GetChange",
		"route53:GetHostedZone",
		"route53:ListHostedZones",
		"route53:ListHostedZonesByName",
		"route53:ListResourceRecordSets",
		"route53:ListTagsForResource",
		"route53:UpdateHostedZoneComment",

		// S3 related perms
		"s3:CreateBucket",
		"s3:GetAccelerateConfiguration",
		"s3:GetBucketAcl",
		"s3:GetBucketCors",
		"s3:GetBucketLocation",
		"s3:GetBucketLogging",
		"s3:GetBucketObjectLockConfiguration",
		"s3:GetBucketPolicy",
		"s3:GetBucketRequestPayment",
		"s3:GetBucketTagging",
		"s3:GetBucketVersioning",
		"s3:GetBucketWebsite",
		"s3:GetEncryptionConfiguration",
		"s3:GetLifecycleConfiguration",
		"s3:GetReplicationConfiguration",
		"s3:ListBucket",
		"s3:PutBucketAcl",
		"s3:PutBucketPolicy",
		"s3:PutBucketTagging",
		"s3:PutEncryptionConfiguration",

		// More S3 (would be nice to limit 'Resource' to just the bucket we actually interact with...)
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
		"ec2:DeletePlacementGroup",
		"ec2:DeleteTags",
		"ec2:DeleteVolume",
		"elasticloadbalancing:DeleteTargetGroup",
		"elasticloadbalancing:DescribeTargetGroups",
		"iam:DeleteAccessKey",
		"iam:DeleteUser",
		"iam:ListAttachedRolePolicies",
		"iam:ListInstanceProfiles",
		"iam:ListRolePolicies",
		"iam:ListUserPolicies",
		"s3:DeleteBucket",
		"s3:DeleteObject",
		"s3:ListBucketVersions",
		"tag:GetResources",
	},
	// Permissions required for creating network resources
	PermissionCreateNetworking: {
		"ec2:AllocateAddress",
		"ec2:AssociateDhcpOptions",
		"ec2:AssociateAddress",
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
		"ec2:ReleaseAddress",
		"ec2:ReplaceRouteTableAssociation",
	},
	// Permissions required for deleting a cluster with shared network resources
	PermissionDeleteSharedNetworking: {
		"tag:UnTagResources",
	},
	// Permissions required for creating an instance role
	PermissionCreateInstanceRole: {
		"iam:CreateRole",
		"iam:DeleteRole",
		"iam:DeleteRolePolicy",
		"iam:PutRolePolicy",
	},
	// Permissions required for deleting a cluster with shared instance roles
	PermissionDeleteSharedInstanceRole: {
		"iam:UntagRole",
	},
	// Permissions required for creating an instance profile
	PermissionCreateInstanceProfile: {
		"iam:AddRoleToInstanceProfile",
		"iam:CreateInstanceProfile",
		"iam:DeleteInstanceProfile",
		"iam:RemoveRoleFromInstanceProfile",
	},
	// Permissions required for deleting a cluster with shared instance profiles
	PermissionDeleteSharedInstanceProfile: {
		"iam:UntagInstanceProfile",
	},
	PermissionCreateHostedZone: {
		"route53:CreateHostedZone",
	},
	PermissionDeleteHostedZone: {
		"route53:DeleteHostedZone",
	},
	PermissionKMSEncryptionKeys: {
		"kms:Decrypt",
		"kms:Encrypt",
		"kms:GenerateDataKey",
		"kms:GenerateDataKeyWithoutPlainText",
		"kms:DescribeKey",
		"kms:RevokeGrant",
		"kms:CreateGrant",
		"kms:ListGrants",
	},
	PermissionPublicIpv4Pool: {
		// Needed to check the IP pools during install-config validation
		"ec2:DescribePublicIpv4Pools",
		// Needed by terraform because of bootstrap EIP created
		"ec2:DisassociateAddress",
	},
	PermissionDeleteIgnitionObjects: {
		// Needed by terraform during the bootstrap destroy stage.
		"s3:DeleteBucket",
		// Needed by capa which always deletes the ignition objects once the VMs are up.
		"s3:DeleteObject",
	},
}

// ValidateCreds will try to create an AWS session, and also verify that the current credentials
// are sufficient to perform an installation, and that they can be used for cluster runtime
// as either capable of creating new credentials for components that interact with the cloud or
// being able to be passed through as-is to the components that need cloud credentials
func ValidateCreds(ssn *session.Session, groups []PermissionGroup, region string) error {
	requiredPermissions, err := PermissionsList(groups)
	if err != nil {
		return err
	}

	client := ccaws.NewClientFromSession(ssn)

	sParams := &ccaws.SimulateParams{
		Region: region,
	}

	// Check whether we can do an installation
	logger := logrus.StandardLogger()
	canInstall, err := ccaws.CheckPermissionsAgainstActions(client, requiredPermissions, sParams, logger)
	if err != nil {
		return fmt.Errorf("checking install permissions: %w", err)
	}
	if !canInstall {
		return errors.New("current credentials insufficient for performing cluster installation")
	}

	// Check whether we can mint new creds for cluster services needing to interact with the cloud
	canMint, err := ccaws.CheckCloudCredCreation(client, logger)
	if err != nil {
		return fmt.Errorf("mint credentials check: %w", err)
	}
	if canMint {
		return nil
	}

	// Check whether we can use the current credentials in passthrough mode to satisfy
	// cluster services needing to interact with the cloud
	canPassthrough, err := ccaws.CheckCloudCredPassthrough(client, sParams, logger)
	if err != nil {
		return fmt.Errorf("passthrough credentials check: %w", err)
	}
	if canPassthrough {
		return nil
	}

	return errors.New("AWS credentials cannot be used to either create new creds or use as-is")
}

// RequiredPermissionGroups returns a set of required permissions for a given cluster configuration.
func RequiredPermissionGroups(ic *types.InstallConfig) []PermissionGroup {
	permissionGroups := []PermissionGroup{PermissionCreateBase}
	usingExistingVPC := len(ic.AWS.Subnets) != 0
	usingExistingPrivateZone := len(ic.AWS.HostedZone) != 0

	if !usingExistingVPC {
		permissionGroups = append(permissionGroups, PermissionCreateNetworking)
	}

	if !usingExistingPrivateZone {
		permissionGroups = append(permissionGroups, PermissionCreateHostedZone)
	}

	if includesKMSEncryptionKey(ic) {
		logrus.Debugf("Adding %s to the group of permissions", PermissionKMSEncryptionKeys)
		permissionGroups = append(permissionGroups, PermissionKMSEncryptionKeys)
	}

	// Add delete permissions for non-C2S installs.
	if !aws.IsSecretRegion(ic.AWS.Region) {
		permissionGroups = append(permissionGroups, PermissionDeleteBase)
		if usingExistingVPC {
			permissionGroups = append(permissionGroups, PermissionDeleteSharedNetworking)
		} else {
			permissionGroups = append(permissionGroups, PermissionDeleteNetworking)
		}
		if !usingExistingPrivateZone {
			permissionGroups = append(permissionGroups, PermissionDeleteHostedZone)
		}
	}

	if ic.AWS.PublicIpv4Pool != "" {
		permissionGroups = append(permissionGroups, PermissionPublicIpv4Pool)
	}

	if !ic.AWS.BestEffortDeleteIgnition {
		permissionGroups = append(permissionGroups, PermissionDeleteIgnitionObjects)
	}

	if includesCreateInstanceRole(ic) {
		permissionGroups = append(permissionGroups, PermissionCreateInstanceRole)
	}

	if includesExistingInstanceRole(ic) {
		permissionGroups = append(permissionGroups, PermissionDeleteSharedInstanceRole)
	}

	if includesExistingInstanceProfile(ic) {
		permissionGroups = append(permissionGroups, PermissionDeleteSharedInstanceProfile)
	}

	if includesCreateInstanceProfile(ic) {
		permissionGroups = append(permissionGroups, PermissionCreateInstanceProfile)
	}

	return permissionGroups
}

// PermissionsList compiles a list of permissions based on the permission groups provided.
func PermissionsList(required []PermissionGroup) ([]string, error) {
	requiredPermissions := sets.New[string]()
	for _, group := range required {
		groupPerms, ok := permissions[group]
		if !ok {
			return nil, fmt.Errorf("unable to access permissions group %s", group)
		}
		requiredPermissions.Insert(groupPerms...)
	}

	return sets.List(requiredPermissions), nil
}

// includesExistingInstanceRole checks if at least one BYO instance role is included in the install-config.
func includesExistingInstanceRole(installConfig *types.InstallConfig) bool {
	mpool := aws.MachinePool{}
	mpool.Set(installConfig.AWS.DefaultMachinePlatform)

	if mp := installConfig.ControlPlane; mp != nil {
		mpool.Set(mp.Platform.AWS)
	}

	for _, compute := range installConfig.Compute {
		mpool.Set(compute.Platform.AWS)
	}

	return len(mpool.IAMRole) > 0
}

// includesCreateInstanceRole checks if at least one instance role will be created by the installer.
// Note: instance profiles have a role attached to them.
func includesCreateInstanceRole(installConfig *types.InstallConfig) bool {
	{
		mpool := aws.MachinePool{}
		mpool.Set(installConfig.AWS.DefaultMachinePlatform)
		if mp := installConfig.ControlPlane; mp != nil {
			mpool.Set(mp.Platform.AWS)
		}
		if len(mpool.IAMRole) == 0 && len(mpool.IAMProfile) == 0 {
			return true
		}
	}

	for _, compute := range installConfig.Compute {
		mpool := aws.MachinePool{}
		mpool.Set(installConfig.AWS.DefaultMachinePlatform)
		mpool.Set(compute.Platform.AWS)
		if len(mpool.IAMRole) == 0 && len(mpool.IAMProfile) == 0 {
			return true
		}
	}

	if len(installConfig.Compute) > 0 {
		return false
	}

	// If compute stanza is not defined, we know it'll inherit the value from DefaultMachinePlatform
	mpool := aws.MachinePool{}
	mpool.Set(installConfig.AWS.DefaultMachinePlatform)
	return len(mpool.IAMRole) == 0 && len(mpool.IAMProfile) == 0
}

// includesKMSEncryptionKey checks if any KMS encryption keys are included in the install-config.
func includesKMSEncryptionKey(installConfig *types.InstallConfig) bool {
	mpool := aws.MachinePool{}
	mpool.Set(installConfig.AWS.DefaultMachinePlatform)

	if mp := installConfig.ControlPlane; mp != nil {
		mpool.Set(mp.Platform.AWS)
	}

	for _, compute := range installConfig.Compute {
		mpool.Set(compute.Platform.AWS)
	}

	return len(mpool.KMSKeyARN) > 0
}

// includesExistingInstanceProfile checks if at least one BYO instance profile is included in the install-config.
func includesExistingInstanceProfile(installConfig *types.InstallConfig) bool {
	mpool := aws.MachinePool{}
	mpool.Set(installConfig.AWS.DefaultMachinePlatform)

	if mp := installConfig.ControlPlane; mp != nil {
		mpool.Set(mp.Platform.AWS)
	}

	for _, compute := range installConfig.Compute {
		mpool.Set(compute.Platform.AWS)
	}

	return len(mpool.IAMProfile) > 0
}

// includesCreateInstanceProfile checks if at least one instance profile will be created by the Installer.
func includesCreateInstanceProfile(installConfig *types.InstallConfig) bool {
	{
		mpool := aws.MachinePool{}
		mpool.Set(installConfig.AWS.DefaultMachinePlatform)
		if mp := installConfig.ControlPlane; mp != nil {
			mpool.Set(mp.Platform.AWS)
		}
		if len(mpool.IAMProfile) == 0 {
			return true
		}
	}

	for _, compute := range installConfig.Compute {
		mpool := aws.MachinePool{}
		mpool.Set(installConfig.AWS.DefaultMachinePlatform)
		mpool.Set(compute.Platform.AWS)
		if len(mpool.IAMProfile) == 0 {
			return true
		}
	}

	if len(installConfig.Compute) > 0 {
		return false
	}

	// If compute stanza is not defined, we know it'll inherit the value from DefaultMachinePlatform
	mpool := aws.MachinePool{}
	mpool.Set(installConfig.AWS.DefaultMachinePlatform)
	return len(mpool.IAMProfile) == 0
}
