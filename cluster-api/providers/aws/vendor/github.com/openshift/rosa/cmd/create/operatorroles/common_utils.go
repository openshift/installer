package operatorroles

import (
	"fmt"

	awsCommonUtils "github.com/openshift-online/ocm-common/pkg/aws/utils"
	errors "github.com/zgalor/weberr"

	"github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/aws/tags"
	"github.com/openshift/rosa/pkg/helper"
	"github.com/openshift/rosa/pkg/rosa"
)

const assumePolicyAction = "sts:AssumeRole"

const policyDetails = "{\n  \"Version\": \"2012-10-17\",\n  \"Statement\": {\n    \"Effect\": \"Allow\",\n    " +
	"\"Action\": \"sts:AssumeRole\",\n    \"Resource\": \"%{shared_vpc_role_arn}\"\n  }\n}\n"

func computePolicyARN(creator aws.Creator,
	prefix string, namespace string, name string, path string) string {
	if prefix == "" {
		prefix = aws.DefaultPrefix
	}
	policy := fmt.Sprintf("%s-%s-%s", prefix, namespace, name)
	policy = awsCommonUtils.TruncateRoleName(policy)
	if path != "" {
		return fmt.Sprintf("arn:%s:iam::%s:policy%s%s", creator.Partition, creator.AccountID, path, policy)
	}
	return fmt.Sprintf("arn:%s:iam::%s:policy/%s", creator.Partition, creator.AccountID, policy)
}

func validateIngressOperatorPolicyOverride(r *rosa.Runtime, policyArn string, sharedVpcRoleArn string,
	installerRolePrefix string) error {
	_, err := r.AWSClient.IsPolicyExists(policyArn)
	policyExists := err == nil
	if !policyExists {
		return nil
	}

	policyDocument, err := r.AWSClient.GetDefaultPolicyDocument(policyArn)
	if err != nil {
		return err
	}

	document, err := aws.ParsePolicyDocument(policyDocument)
	if err != nil {
		return err
	}

	for _, statement := range document.Statement {
		if statement.Action == assumePolicyAction && statement.Effect == "Allow" {
			// The policy associated with the installer role. In the case it contains a different shared VPC role ARN,
			// don't override it.
			if statement.Resource != sharedVpcRoleArn {
				return errors.UserErrorf("Policy with ARN '%s' contains '%s' with an unexpected shared VPC role ARN "+
					"[Expected: '%s', Provided '%s'].\n"+
					"The policy is associated with the installer role with the prefix '%s'.\n"+
					"To create operator roles with shared VPC role ARN '%s', please provide a different value for "+
					"'--role-arn'", policyArn, assumePolicyAction, statement.Resource, sharedVpcRoleArn,
					installerRolePrefix, sharedVpcRoleArn)
			}
		}
	}

	return nil
}

func getHcpSharedVpcPolicy(r *rosa.Runtime, roleArn string, defaultPolicyVersion string) (string, error) {
	interpolatedPolicyDetails := aws.InterpolatePolicyDocument(r.Creator.Partition, policyDetails, map[string]string{
		"shared_vpc_role_arn": roleArn,
	})
	userProvidedRoleName, err := aws.GetResourceIdFromARN(roleArn)
	if err != nil {
		return "", err
	}

	path, err := aws.GetPathFromARN(roleArn)
	if err != nil {
		return "", err
	}

	policyName := fmt.Sprintf(aws.AssumeRolePolicyPrefix, userProvidedRoleName)
	policy := aws.GetPolicyArn(r.Creator.Partition, r.Creator.AccountID, policyName, path)

	iamTags := map[string]string{
		tags.RedHatManaged: helper.True,
		tags.HcpSharedVpc:  helper.True,
	}

	policyArn, err := r.AWSClient.EnsurePolicy(policy, interpolatedPolicyDetails, defaultPolicyVersion,
		iamTags, path)
	if err != nil {
		return policyArn, err
	}
	return policyArn, nil
}
