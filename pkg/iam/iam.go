// Package iam checks that the AWS user currently logged in
// has sufficient permissions to run the installer for both
// `create` and `destroy` operations.
package iam

import (
	"encoding/json"
	"net/url"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sirupsen/logrus"
)

// AccountDetails stores the AWS account-wide groups and policies.
type AccountDetails struct {
	users    []*iam.UserDetail
	groups   []*iam.GroupDetail
	policies []*iam.ManagedPolicyDetail
}

// actions is a slice containing the names of AWS IAM permissions
// (actions) required to run the installer.
var actions = []string{
	"ec2:AllocateAddress",
	"ec2:AssociateRouteTable",
	"ec2:AttachInternetGateway",
	"ec2:AuthorizeSecurityGroupEgress",
	"ec2:AuthorizeSecurityGroupIngress",
	"ec2:CreateInternetGateway",
	"ec2:CreateNatGateway",
	"ec2:CreateRoute",
	"ec2:CreateRouteTable",
	"ec2:CreateSecurityGroup",
	"ec2:CreateSubnet",
	"ec2:CreateTags",
	"ec2:CreateVpc",
	"ec2:DeleteInternetGateway",
	"ec2:DeleteNatGateway",
	"ec2:DeleteNetworkInterface",
	"ec2:DeleteRoute",
	"ec2:DeleteRouteTable",
	"ec2:DeleteSecurityGroup",
	"ec2:DeleteSubnet",
	"ec2:DeleteVpc",
	"ec2:DescribeAccountAttributes",
	"ec2:DescribeAddresses",
	"ec2:DescribeAvailabilityZones",
	"ec2:DescribeImages",
	"ec2:DescribeInstanceAttribute",
	"ec2:DescribeInstances",
	"ec2:DescribeInternetGateways",
	"ec2:DescribeNatGateways",
	"ec2:DescribeNetworkAcls",
	"ec2:DescribeNetworkInterfaces",
	"ec2:DescribeRegions",
	"ec2:DescribeRouteTables",
	"ec2:DescribeSecurityGroups",
	"ec2:DescribeSubnets",
	"ec2:DescribeTags",
	"ec2:DescribeVolumes",
	"ec2:DescribeVpcAttribute",
	"ec2:DescribeVpcClassicLink",
	"ec2:DescribeVpcClassicLinkDnsSupport",
	"ec2:DescribeVpcs",
	"ec2:DetachInternetGateway",
	"ec2:DisassociateRouteTable",
	"ec2:ModifyNetworkInterfaceAttribute",
	"ec2:ModifyVpcAttribute",
	"ec2:ReleaseAddress",
	"ec2:ReplaceRouteTableAssociation",
	"ec2:RevokeSecurityGroupEgress",
	"ec2:RevokeSecurityGroupIngress",
	"ec2:RunInstances",
	"ec2:TerminateInstances",
	"elasticloadbalancing:AddTags",
	"elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
	"elasticloadbalancing:AttachLoadBalancerToSubnets",
	"elasticloadbalancing:ConfigureHealthCheck",
	"elasticloadbalancing:CreateLoadBalancer",
	"elasticloadbalancing:CreateLoadBalancerListeners",
	"elasticloadbalancing:DeleteLoadBalancer",
	"elasticloadbalancing:DeleteLoadBalancerListeners",
	"elasticloadbalancing:DescribeLoadBalancerAttributes",
	"elasticloadbalancing:DescribeLoadBalancers",
	"elasticloadbalancing:DescribeTags",
	"elasticloadbalancing:ModifyLoadBalancerAttributes",
	"elasticloadbalancing:RegisterInstancesWithLoadBalancer",
	"iam:AddRoleToInstanceProfile",
	"iam:CreateInstanceProfile",
	"iam:CreateRole",
	"iam:DeleteInstanceProfile",
	"iam:DeleteRole",
	"iam:DeleteRolePolicy",
	"iam:GetAccountAuthorizationDetails",
	"iam:GetInstanceProfile",
	"iam:GetRole",
	"iam:GetUser",
	"iam:ListAttachedGroupPolicies",
	"iam:ListGroupsForUser",
	"iam:ListGroupPolicies",
	"iam:ListInstanceProfiles",
	"iam:ListRolePolicies",
	"iam:PassRole",
	"iam:PutRolePolicy",
	"iam:RemoveRoleFromInstanceProfile",
	"route53:ChangeResourceRecordSets",
	"route53:ChangeTagsForResource",
	"route53:CreateHostedZone",
	"route53:DeleteHostedZone",
	"route53:GetChange",
	"route53:GetHostedZone",
	"route53:ListHostedZones",
	"route53:ListResourceRecordSets",
	"route53:ListTagsForResource",
	"route53:UpdateHostedZoneComment",
	"s3:CreateBucket",
	"s3:DeleteBucket",
	"s3:DeleteObject",
	"s3:GetAccelerateConfiguration",
	"s3:GetBucketCORS",
	"s3:GetBucketLocation",
	"s3:GetBucketLogging",
	"s3:GetBucketRequestPayment",
	"s3:GetBucketTagging",
	"s3:GetBucketVersioning",
	"s3:GetBucketWebsite",
	"s3:GetEncryptionConfiguration",
	"s3:GetLifecycleConfiguration",
	"s3:GetObject",
	"s3:GetObjectTagging",
	"s3:GetReplicationConfiguration",
	"s3:HeadBucket",
	"s3:ListAllMyBuckets",
	"s3:ListBucket",
	"s3:PutBucketAcl",
	"s3:PutBucketTagging",
	"s3:PutObject",
	"s3:PutObjectAcl",
	"s3:PutObjectTagging",
}

// getAccountAuthDetails fetches lists of users, groups, and polices.
func getAccountAuthDetails(svc *iam.IAM) AccountDetails {

	var results AccountDetails

	// Initial API call to get account auth details
	input := &iam.GetAccountAuthorizationDetailsInput{}
	resp, err := svc.GetAccountAuthorizationDetails(input)
	if err != nil {
		logrus.Fatalf("Error fetching account authorization details: %s", err.Error())
	}

	results.users = append(results.users, resp.UserDetailList...)
	results.groups = append(results.groups, resp.GroupDetailList...)
	results.policies = append(results.policies, resp.Policies...)

	// If the above results were truncated, fetch again until
	// all results are recorded.
	logrus.Debugf("GetAccountAuthorizationDetails results truncated? %s", *resp.IsTruncated)
	for *resp.IsTruncated {
		logrus.Debugf("Fetching truncated page...")
		marker := *resp.Marker
		input = &iam.GetAccountAuthorizationDetailsInput{Marker: &marker}
		resp, err = svc.GetAccountAuthorizationDetails(input)
		if err != nil {
			logrus.Fatalf("Error fetching account authorization details: %s", err.Error())
		}
		results.users = append(results.users, resp.UserDetailList...)
		results.groups = append(results.groups, resp.GroupDetailList...)
		results.policies = append(results.policies, resp.Policies...)
	}

	return results
}

// getCurrentUserDetails returns details about the user making this request.
func getCurrentUserDetails(svc *iam.IAM, a AccountDetails) *iam.UserDetail {

	// Get the name of the user making this request.
	currentuser, err := svc.GetUser(&iam.GetUserInput{})
	if err != nil {
		logrus.Fatalf("Error running AWS IAM GetUser: %s", err)
	}
	currentusername := *currentuser.User.UserName

	// Look through accountAuthDetails for all users.
	// Narrow results down to current user details.
	var userdetails *iam.UserDetail
	for _, user := range a.users {
		if *user.UserName == currentusername {
			userdetails = user
			return userdetails
		}
	}

	return userdetails
}

// getPolicyDocFromUser inspects a UserDetail object and parses the
// AttachedManagedPolicies and UserPolicyList, to return a list of
// PolicyDocuments.
func getPolicyDocFromUser(svc *iam.IAM, user *iam.UserDetail, userPolicyDocs *[]string, a AccountDetails) []string {
	// Look through each policy ARN in this AWS account.
	// Find the ARNs specifically attached to this user.
	for _, userPolicy := range user.AttachedManagedPolicies {
		for _, managedPolicy := range a.policies {
			if *managedPolicy.Arn == *userPolicy.PolicyArn {
				logrus.Infof("Found attached managed user policy: %s", *managedPolicy.PolicyName)
				for _, v := range managedPolicy.PolicyVersionList {
					if *v.IsDefaultVersion {
						*userPolicyDocs = append(*userPolicyDocs, *v.Document)
					}
				}
			}
		}
	}
	// In addition to AttachedManagedPolicies above, also check the
	// UserPolicyList for inline policies attached to the user.
	for _, i := range user.UserPolicyList {
		*userPolicyDocs = append(*userPolicyDocs, *i.PolicyDocument)
		logrus.Infof("Found inline user policy:", *i.PolicyName)
	}
	return *userPolicyDocs
}

// getPolicyDocFromGroup inspects an IAM group and returns a list of
// PolicyDocuments.
func getPolicyDocFromGroup(svc *iam.IAM, user *iam.UserDetail, userPolicyDocs *[]string, a AccountDetails) []string {

	// If user is not a member of any groups, don't inspect groups.
	if len(user.GroupList) < 1 {
		return *userPolicyDocs
	}

	// Get details about the groups that the user belongs to.
	for _, g := range a.groups {
		for _, u := range user.GroupList {
			if *g.GroupName == *u {
				// Find all inline policies for the user's group(s).
				for _, policyInGroup := range g.GroupPolicyList {
					logrus.Infof("Found inline group policy:", *policyInGroup.PolicyName)
					*userPolicyDocs = append(*userPolicyDocs, *policyInGroup.PolicyDocument)
				}
				// Find all AttachedManagedPolicies for the user's group(s).
				if len(g.AttachedManagedPolicies) > 0 {
					for _, accountPolicy := range a.policies {
						for _, userPolicy := range g.AttachedManagedPolicies {
							if *accountPolicy.Arn == *userPolicy.PolicyArn {
								logrus.Infof("Found attached managed group policy:", *userPolicy.PolicyName)
								for _, v := range accountPolicy.PolicyVersionList {
									if *v.IsDefaultVersion {
										*userPolicyDocs = append(*userPolicyDocs, *v.Document)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return *userPolicyDocs
}

// getPermissionsFromPolicyDoc parses a PolicyDocument and extracts
// a list of effective actions granted. It only counts actions that
// are granted to all resources (listed as "Resource": "*" in the
// Policy Document).
func getPermissionsFromPolicyDoc(userPolicyDocs *[]string) []string {
	// permissions is the refined collection of all permissions
	// granted to this user through IAM.
	var permissions []string

	// Look through all Policy Documents associated with this user
	// and gather a list of all permissions granted to the user.
	logrus.Debugf("Found the following permissions applied to the current user:")
	for _, p := range *userPolicyDocs {
		u, err := url.QueryUnescape(p)
		logrus.Debugf(u)
		if err != nil {
			logrus.Fatalf("Error decoding policy document: %s", p)
		}

		userPolicy := []byte(u)

		// a stores the arbitrary json response received from AWS.
		// It contains the IAM policy documents related to this user.
		var a map[string]interface{}
		err = json.Unmarshal(userPolicy, &a)
		for _, statement := range a["Statement"].([]interface{}) {
			elem := statement.(map[string]interface{})
			if elem["Effect"] == "Allow" {
				var singleAction bool
				var allResources bool
				_, ok := elem["Resource"].(string)
				if ok {
					if elem["Resource"] == "*" {
						allResources = true
					}
				}
				_, ok = elem["Action"].(string)
				if ok {
					singleAction = true
				}
				if allResources && singleAction {
					singleact, ok := elem["Action"].(string)
					if ok {
						permissions = append(permissions, singleact)
					}
				} else if allResources {
					for _, action := range elem["Action"].([]interface{}) {
						act, ok := action.(string)
						if ok {
							permissions = append(permissions, act)
						} else {
							logrus.Fatalf("Error appending action %q to permissions list: %v", action, permissions)
						}
					}
				}
			}
		}
	}
	return permissions
}

// checkOffPermissionsByPrefix accepts a map of permissions
// to be used as a checklist. It then marks all permissions 'true'
// if they match the specified prefix.
func checkOffPermissionsByPrefix(checklist map[string]bool, prefix string) map[string]bool {
	if prefix == "*" {
		for k := range checklist {
			checklist[k] = true
		}
	}
	for k := range checklist {
		match, _ := regexp.MatchString(prefix, k)
		if match {
			checklist[k] = true
		}
	}
	return checklist
}

// RunPermissionsCheck ties together all the above helper functions
// to run a complete suite of AWS IAM permissions checks. This will
// inspect permissions attached directly to users and groups, as well
// as managed policies containing permissions granted to the user.
func RunPermissionsCheck() {

	// TODO: sessions should be shared whenever possible. This should
	// be passed in as an argument instead of created new.

	// Create an IAM service client.
	ssn := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{},
	}))
	svc := iam.New(ssn)

	// The above session creation will not return an error if user
	// fails to specify any credentials. Therefore, we must check
	// credentials before making any API calls.
	_, err := ssn.Config.Credentials.Get()
	if err != nil {
		logrus.Fatalf("Failed to find AWS credentials.")
	}

	accountAuthDetails := getAccountAuthDetails(svc)
	user := getCurrentUserDetails(svc, accountAuthDetails)
	logrus.Infof("Current IAM user: %s", *user.UserName)
	logrus.Debugf("User details:\n %s", user)

	// checklist is a map used as a checklist, containing all
	// permissions required to run both the installer and uninstaller.
	checklist := make(map[string]bool)
	for _, p := range actions {
		checklist[p] = false
	}

	// userPolicyDocs is a slice that stores the Policy Documents
	// (permissions details) associated with the current user.
	var userPolicyDocs []string

	getPolicyDocFromUser(svc, user, &userPolicyDocs, accountAuthDetails)
	getPolicyDocFromGroup(svc, user, &userPolicyDocs, accountAuthDetails)

	userPermissions := getPermissionsFromPolicyDoc(&userPolicyDocs)

	for _, permission := range userPermissions {
		checkOffPermissionsByPrefix(checklist, permission)
	}

	missing := false
	for permission, granted := range checklist {
		if !granted {
			missing = true
			logrus.Errorf("Missing permission: %s", permission)
		}
	}
	if missing {
		logrus.Fatalf("Please add the above permission(s) to the current IAM user and try again.")
	} else {
		logrus.Infof("\nIAM permissions check passed.")
	}
}
