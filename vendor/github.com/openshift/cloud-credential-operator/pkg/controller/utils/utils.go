package utils

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	minterv1 "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	ccaws "github.com/openshift/cloud-credential-operator/pkg/aws"

	"github.com/aws/aws-sdk-go/aws"
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

func init() {
	if err := minterv1.AddToScheme(credentailRequestScheme); err != nil {
		panic(err)
	}
}

// CheckCloudCredCreation will see whether we have enough permissions to create new sub-creds
func CheckCloudCredCreation(awsClient ccaws.Client, logger log.FieldLogger) (bool, error) {
	return CheckPermissionsAgainstActions(awsClient, credMintingActions, logger)
}

// CheckPermissionsUsingQueryClient will use queryClient to query whether the credentials in targetClient can perform the actions
// listed in the statementEntries. queryClient will need iam:GetUser and iam:SimulatePrincipalPolicy
func CheckPermissionsUsingQueryClient(queryClient, targetClient ccaws.Client, statementEntries []minterv1.StatementEntry, logger log.FieldLogger) (bool, error) {
	targetUsername, err := targetClient.GetUser(nil)
	if err != nil {
		return false, fmt.Errorf("error querying current username: %v", err)
	}

	allowList := []*string{}
	for _, statement := range statementEntries {
		for _, action := range statement.Action {
			allowList = append(allowList, aws.String(action))
		}
	}

	results, err := queryClient.SimulatePrincipalPolicy(&iam.SimulatePrincipalPolicyInput{
		PolicySourceArn: targetUsername.User.Arn,
		ActionNames:     allowList,
	})
	if err != nil {
		return false, fmt.Errorf("error simulating policy: %v", err)
	}

	// Either they are all allowed and we return 'true', or it's a failure
	allClear := true
	for _, result := range results.EvaluationResults {
		if *result.EvalDecision != "allowed" {
			// Don't return on the first failure, so we can log the full list
			// of failed/denied actions
			logger.WithField("action", *result.EvalActionName).Warning("Action not allowed with tested creds")
			allClear = false
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
func CheckPermissionsAgainstStatementList(awsClient ccaws.Client, statementEntries []minterv1.StatementEntry, logger log.FieldLogger) (bool, error) {
	return CheckPermissionsUsingQueryClient(awsClient, awsClient, statementEntries, logger)
}

// CheckPermissionsAgainstActions will take the static list of Actions to check whether the provided
// awsClient creds have sufficient permissions to perform the actions.
// Will return true/false indicating whether the permissions are sufficient.
func CheckPermissionsAgainstActions(awsClient ccaws.Client, actionList []string, logger log.FieldLogger) (bool, error) {
	statementList := []minterv1.StatementEntry{
		{
			Action:   actionList,
			Resource: "*",
			Effect:   "Allow",
		},
	}

	return CheckPermissionsAgainstStatementList(awsClient, statementList, logger)
}

// CheckCloudCredPassthrough will see if the provided creds are good enough to pass through
// to other components as-is based on the static list of permissions needed by the various
// users of CredentialsRequests
// TODO: move away from static list (to dynamic passthrough validation?)
func CheckCloudCredPassthrough(awsClient ccaws.Client, logger log.FieldLogger) (bool, error) {
	return CheckPermissionsAgainstActions(awsClient, credPassthroughActions, logger)
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

	awsSpec, err := awsCodec.DecodeProviderSpec(cr.Spec.ProviderSpec, &minterv1.AWSProviderSpec{})
	if err != nil {
		return statementList, fmt.Errorf("error decoding spec.ProviderSpec: %v", err)
	}

	statementList = append(statementList, awsSpec.StatementEntries...)

	return statementList, nil
}
