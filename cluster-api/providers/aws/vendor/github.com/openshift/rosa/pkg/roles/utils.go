package roles

import (
	"fmt"

	errors "github.com/zgalor/weberr"

	"github.com/openshift/rosa/pkg/aws"
	awscb "github.com/openshift/rosa/pkg/aws/commandbuilder"
	"github.com/openshift/rosa/pkg/aws/tags"
	"github.com/openshift/rosa/pkg/rosa"
)

const policyDocumentBody = ` \
'{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": "sts:AssumeRole",
    "Resource": "%{shared_vpc_role_arn}"
  }
}'`

type ManualSharedVpcPolicyDetails struct {
	Command       string
	Name          string
	AlreadyExists bool
	Path          string
}

func GetHcpSharedVpcPolicyDetails(r *rosa.Runtime, roleArn string) (bool, string,
	string, error) {
	interpolatedPolicyDetails := aws.InterpolatePolicyDocument(r.Creator.Partition, policyDocumentBody,
		map[string]string{
			"shared_vpc_role_arn": roleArn,
		})

	roleName, err := aws.GetResourceIdFromARN(roleArn)
	if err != nil {
		return false, "", "", err
	}
	path, err := aws.GetPathFromARN(roleArn)
	if err != nil {
		return false, "", "", err
	}

	policyName := fmt.Sprintf(aws.AssumeRolePolicyPrefix, roleName)

	predictedPolicyArn := aws.GetPolicyArn(r.Creator.Partition, r.Creator.AccountID, policyName, path)

	existsQuery, _ := r.AWSClient.IsPolicyExists(predictedPolicyArn)

	var iamTags = map[string]string{
		tags.RedHatManaged: aws.TrueString,
		tags.HcpSharedVpc:  aws.TrueString,
	}

	createPolicy := awscb.NewIAMCommandBuilder().
		SetCommand(awscb.CreatePolicy).
		AddParam(awscb.PolicyName, policyName).
		AddParam(awscb.PolicyDocument, interpolatedPolicyDetails).
		AddTags(iamTags).
		AddParam(awscb.Path, path).
		Build()

	return existsQuery != nil, createPolicy, policyName, nil
}

func CheckIfRolesAreHcpSharedVpc(r *rosa.Runtime, roles []string) bool {
	isHcpSharedVpc := false
	for _, roleName := range roles {
		ptrRoleName := roleName
		attachedPolicies, err := r.AWSClient.GetPolicyDetailsFromRole(&ptrRoleName)
		if err != nil {
			r.Reporter.Errorf("Failed to get policy details for role '%s': %v", roleName, err)
		}
		for _, attachedPolicy := range attachedPolicies {
			rhManaged := false
			hcpSharedVpc := false
			for _, tag := range attachedPolicy.Policy.Tags {
				if *tag.Key == tags.RedHatManaged {
					rhManaged = true
				} else if *tag.Key == tags.HcpSharedVpc {
					hcpSharedVpc = true
				}
			}
			if rhManaged && hcpSharedVpc {
				isHcpSharedVpc = true
			}
		}
	}
	return isHcpSharedVpc
}

func GetPolicyDetailsByName(details map[string]ManualSharedVpcPolicyDetails, name string) (ManualSharedVpcPolicyDetails,
	error) {
	for _, detail := range details {
		if detail.Name == name {
			return detail, nil
		}
	}
	return ManualSharedVpcPolicyDetails{}, errors.Errorf("Policy %s not found", name)
}
