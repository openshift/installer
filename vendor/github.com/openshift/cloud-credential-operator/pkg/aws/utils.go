package aws

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	minterv1 "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"

	"k8s.io/apimachinery/pkg/runtime"
)

// credMintingActions is a list of AWS verbs needed to run in the mode where the
// cloud-credential-operator can mint new creds to satisfy CredentialRequest CRDs
var (
	credMintingActions = []string{
		"iam:CreateAccessKey",
		"iam:CreateUser",
		"iam:DeleteAccessKey",
		"iam:DeleteUser",
		"iam:DeleteUserPolicy",
		"iam:GetUser",
		"iam:GetUserPolicy",
		"iam:ListAccessKeys",
		"iam:PutUserPolicy",
		"iam:TagUser",
		"iam:SimulatePrincipalPolicy", // needed so we can verify the above list of course
	}

	credPassthroughActions = []string{
		// so we can query whether we have the below list of creds
		"iam:GetUser",
		"iam:SimulatePrincipalPolicy",

		// openshift-ingress
		"elasticloadbalancing:DescribeLoadBalancers",
		"route53:ListHostedZones",
		"route53:ChangeResourceRecordSets",
		"tag:GetResources",

		// openshift-image-registry
		"s3:CreateBucket",
		"s3:DeleteBucket",
		"s3:PutBucketTagging",
		"s3:GetBucketTagging",
		"s3:PutEncryptionConfiguration",
		"s3:GetEncryptionConfiguration",
		"s3:PutLifecycleConfiguration",
		"s3:GetLifecycleConfiguration",
		"s3:GetBucketLocation",
		"s3:ListBucket",
		"s3:GetObject",
		"s3:PutObject",
		"s3:DeleteObject",
		"s3:ListBucketMultipartUploads",
		"s3:AbortMultipartUpload",

		// openshift-cluster-api
		"ec2:DescribeImages",
		"ec2:DescribeVpcs",
		"ec2:DescribeSubnets",
		"ec2:DescribeAvailabilityZones",
		"ec2:DescribeSecurityGroups",
		"ec2:RunInstances",
		"ec2:DescribeInstances",
		"ec2:TerminateInstances",
		"elasticloadbalancing:RegisterInstancesWithLoadBalancer",
		"elasticloadbalancing:DescribeLoadBalancers",
		"elasticloadbalancing:DescribeTargetGroups",
		"elasticloadbalancing:RegisterTargets",
		"ec2:DescribeVpcs",
		"ec2:DescribeSubnets",
		"ec2:DescribeAvailabilityZones",
		"ec2:DescribeSecurityGroups",
		"ec2:RunInstances",
		"ec2:DescribeInstances",
		"ec2:TerminateInstances",
		"elasticloadbalancing:RegisterInstancesWithLoadBalancer",
		"elasticloadbalancing:DescribeLoadBalancers",
		"elasticloadbalancing:DescribeTargetGroups",
		"elasticloadbalancing:RegisterTargets",

		// iam-ro
		"iam:GetUser",
		"iam:GetUserPolicy",
		"iam:ListAccessKeys",
	}

	credentailRequestScheme = runtime.NewScheme()
)

func init() {
	if err := minterv1.AddToScheme(credentailRequestScheme); err != nil {
		panic(err)
	}
}

// SimulateParams captures any additional details that should be used
// when simulating permissions.
type SimulateParams struct {
	Region string
}

// CheckCloudCredCreation will see whether we have enough permissions to create new sub-creds
func CheckCloudCredCreation(ctx context.Context, awsClient Client, logger log.FieldLogger) (bool, error) {
	// Empty SimulateParams{} b/c creating IAM users and assigning policies
	// are all IAM API alls which are not region-specific
	return CheckPermissionsAgainstActions(ctx, awsClient, credMintingActions, &SimulateParams{}, logger)
}

// getClientDetails will return the *iam.User associated with the provided client's credentials,
// a boolean indicating whether the user is the 'root' account, and any error encountered
// while trying to gather the info.
func getClientDetails(ctx context.Context, awsClient Client) (*iamtypes.User, bool, error) {
	rootUser := false

	user, err := awsClient.GetUser(ctx, &iam.GetUserInput{})
	if err != nil {
		return nil, rootUser, fmt.Errorf("error querying username: %v", err)
	}

	// Detect whether the AWS account's root user is being used
	parsed, err := arn.Parse(*user.User.Arn)
	if err != nil {
		return nil, rootUser, fmt.Errorf("error parsing user's ARN: %v", err)
	}
	if parsed.AccountID == *user.User.UserId {
		rootUser = true
	}

	return user.User, rootUser, nil
}

// CheckPermissionsUsingQueryClient will use queryClient to query whether the credentials in targetClient can perform the actions
// listed in the statementEntries. queryClient will need iam:GetUser and iam:SimulatePrincipalPolicy
func CheckPermissionsUsingQueryClient(ctx context.Context, queryClient, targetClient Client, statementEntries []minterv1.StatementEntry,
	params *SimulateParams, logger log.FieldLogger) (bool, error) {
	targetUser, isRoot, err := getClientDetails(ctx, targetClient)
	if err != nil {
		return false, fmt.Errorf("error gathering AWS credentials details: %v", err)
	}
	if isRoot {
		// warn about using the root creds, and just return that the creds are good enough
		logger.Warn("Using the AWS account root user is not recommended: https://docs.aws.amazon.com/general/latest/gr/managing-aws-access-keys.html")
		return true, nil
	}

	allowList := []string{}
	for _, statement := range statementEntries {
		allowList = append(allowList, statement.Action...)
	}

	input := &iam.SimulatePrincipalPolicyInput{
		PolicySourceArn: targetUser.Arn,
		ActionNames:     allowList,
		ContextEntries:  []iamtypes.ContextEntry{},
	}

	if params != nil {
		if params.Region != "" {
			input.ContextEntries = append(input.ContextEntries, iamtypes.ContextEntry{
				ContextKeyName:   awssdk.String("aws:RequestedRegion"),
				ContextKeyType:   iamtypes.ContextKeyTypeEnumString,
				ContextKeyValues: []string{params.Region},
			})
		}
	}

	// Either all actions are allowed and we'll return 'true', or it's a failure
	allClear := true

	paginator := iam.NewSimulatePrincipalPolicyPaginator(queryClient, input)
	for paginator.HasMorePages() {
		response, err := paginator.NextPage(ctx)
		if err != nil {
			return false, fmt.Errorf("error simulating policy: %v", err)
		}

		for _, result := range response.EvaluationResults {
			if result.EvalDecision != iamtypes.PolicyEvaluationDecisionTypeAllowed {
				// Don't bail out after the first failure, so we can log the full list
				// of failed/denied actions
				logger.WithField("action", *result.EvalActionName).Warning("Action not allowed with tested creds")
				allClear = false
			}
		}
	}

	if !allClear {
		logger.Warningf("Tested creds not able to perform all requested actions")
		return false, nil
	}

	return true, nil

}

// CheckPermissionsAgainstStatementList will test to see whether the list of actions in the provided
// list of StatementEntries can work with the credentials used by the passed-in awsClient
func CheckPermissionsAgainstStatementList(ctx context.Context, awsClient Client, statementEntries []minterv1.StatementEntry,
	params *SimulateParams, logger log.FieldLogger) (bool, error) {
	return CheckPermissionsUsingQueryClient(ctx, awsClient, awsClient, statementEntries, params, logger)
}

// CheckPermissionsAgainstActions will take the static list of Actions to check whether the provided
// awsClient creds have sufficient permissions to perform the actions.
// Will return true/false indicating whether the permissions are sufficient.
func CheckPermissionsAgainstActions(ctx context.Context, awsClient Client, actionList []string, params *SimulateParams, logger log.FieldLogger) (bool, error) {
	statementList := []minterv1.StatementEntry{
		{
			Action:   actionList,
			Resource: "*",
			Effect:   "Allow",
		},
	}

	return CheckPermissionsAgainstStatementList(ctx, awsClient, statementList, params, logger)
}

// CheckCloudCredPassthrough will see if the provided creds are good enough to pass through
// to other components as-is based on the static list of permissions needed by the various
// users of CredentialsRequests
// TODO: move away from static list (to dynamic passthrough validation?)
func CheckCloudCredPassthrough(ctx context.Context, awsClient Client, params *SimulateParams, logger log.FieldLogger) (bool, error) {
	return CheckPermissionsAgainstActions(ctx, awsClient, credPassthroughActions, params, logger)
}
