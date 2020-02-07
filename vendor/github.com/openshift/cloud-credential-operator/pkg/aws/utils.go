package aws

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	minterv1 "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/iam"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
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
		"s3:HeadBucket",
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
	credentialRequestCodec  = serializer.NewCodecFactory(credentailRequestScheme)
)

const (
	infrastructureConfigName = "cluster"
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
func CheckCloudCredCreation(awsClient Client, logger log.FieldLogger) (bool, error) {
	// Empty SimulateParams{} b/c creating IAM users and assigning policies
	// are all IAM API alls which are not region-specific
	return CheckPermissionsAgainstActions(awsClient, credMintingActions, &SimulateParams{}, logger)
}

// getClientDetails will return the *iam.User associated with the provided client's credentials,
// a boolean indicating whether the user is the 'root' account, and any error encountered
// while trying to gather the info.
func getClientDetails(awsClient Client) (*iam.User, bool, error) {
	rootUser := false

	user, err := awsClient.GetUser(nil)
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
func CheckPermissionsUsingQueryClient(queryClient, targetClient Client, statementEntries []minterv1.StatementEntry,
	params *SimulateParams, logger log.FieldLogger) (bool, error) {
	targetUser, isRoot, err := getClientDetails(targetClient)
	if err != nil {
		return false, fmt.Errorf("error gathering AWS credentials details: %v", err)
	}
	if isRoot {
		// warn about using the root creds, and just return that the creds are good enough
		logger.Warn("Using the AWS account root user is not recommended: https://docs.aws.amazon.com/general/latest/gr/managing-aws-access-keys.html")
		return true, nil
	}

	allowList := []*string{}
	for _, statement := range statementEntries {
		for _, action := range statement.Action {
			allowList = append(allowList, aws.String(action))
		}
	}

	input := &iam.SimulatePrincipalPolicyInput{
		PolicySourceArn: targetUser.Arn,
		ActionNames:     allowList,
		ContextEntries:  []*iam.ContextEntry{},
	}

	if params != nil {
		if params.Region != "" {
			input.ContextEntries = append(input.ContextEntries, &iam.ContextEntry{
				ContextKeyName:   aws.String("aws:RequestedRegion"),
				ContextKeyType:   aws.String("string"),
				ContextKeyValues: []*string{aws.String(params.Region)},
			})
		}
	}

	// Either all actions are allowed and we'll return 'true', or it's a failure
	allClear := true

	err = queryClient.SimulatePrincipalPolicyPages(input, func(response *iam.SimulatePolicyResponse, lastPage bool) bool {

		for _, result := range response.EvaluationResults {
			if *result.EvalDecision != "allowed" {
				// Don't bail out after the first failure, so we can log the full list
				// of failed/denied actions
				logger.WithField("action", *result.EvalActionName).Warning("Action not allowed with tested creds")
				allClear = false
			}
		}
		return !lastPage
	})
	if err != nil {
		return false, fmt.Errorf("error simulating policy: %v", err)
	}

	if !allClear {
		logger.Warningf("Tested creds not able to perform all requested actions")
		return false, nil
	}

	return true, nil

}

// CheckPermissionsAgainstStatementList will test to see whether the list of actions in the provided
// list of StatementEntries can work with the credentials used by the passed-in awsClient
func CheckPermissionsAgainstStatementList(awsClient Client, statementEntries []minterv1.StatementEntry,
	params *SimulateParams, logger log.FieldLogger) (bool, error) {
	return CheckPermissionsUsingQueryClient(awsClient, awsClient, statementEntries, params, logger)
}

// CheckPermissionsAgainstActions will take the static list of Actions to check whether the provided
// awsClient creds have sufficient permissions to perform the actions.
// Will return true/false indicating whether the permissions are sufficient.
func CheckPermissionsAgainstActions(awsClient Client, actionList []string, params *SimulateParams, logger log.FieldLogger) (bool, error) {
	statementList := []minterv1.StatementEntry{
		{
			Action:   actionList,
			Resource: "*",
			Effect:   "Allow",
		},
	}

	return CheckPermissionsAgainstStatementList(awsClient, statementList, params, logger)
}

// CheckCloudCredPassthrough will see if the provided creds are good enough to pass through
// to other components as-is based on the static list of permissions needed by the various
// users of CredentialsRequests
// TODO: move away from static list (to dynamic passthrough validation?)
func CheckCloudCredPassthrough(awsClient Client, params *SimulateParams, logger log.FieldLogger) (bool, error) {
	return CheckPermissionsAgainstActions(awsClient, credPassthroughActions, params, logger)
}

func readCredentialRequest(cr []byte) (*minterv1.CredentialsRequest, error) {

	newObj, err := runtime.Decode(credentialRequestCodec.UniversalDecoder(minterv1.SchemeGroupVersion), cr)
	if err != nil {
		return nil, fmt.Errorf("error decoding credentialrequest: %v", err)
	}
	return newObj.(*minterv1.CredentialsRequest), nil
}

func getCredentialRequestStatements(crBytes []byte) ([]minterv1.StatementEntry, error) {
	statementList := []minterv1.StatementEntry{}

	awsCodec, err := minterv1.NewCodec()
	if err != nil {
		return statementList, fmt.Errorf("error creating credentialrequest codec: %v", err)
	}

	cr, err := readCredentialRequest(crBytes)
	if err != nil {
		return statementList, err
	}

	awsSpec := minterv1.AWSProviderSpec{}
	err = awsCodec.DecodeProviderSpec(cr.Spec.ProviderSpec, &awsSpec)
	if err != nil {
		return statementList, fmt.Errorf("error decoding spec.ProviderSpec: %v", err)
	}

	statementList = append(statementList, awsSpec.StatementEntries...)

	return statementList, nil
}
