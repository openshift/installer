package roles

import (
	"fmt"

	common "github.com/openshift-online/ocm-common/pkg/aws/validations"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	awscb "github.com/openshift/rosa/pkg/aws/commandbuilder"
	"github.com/openshift/rosa/pkg/aws/tags"
)

type ManualCommandsForMissingOperatorRolesInput struct {
	ClusterID                string
	OperatorRolePolicyPrefix string
	Operator                 *cmv1.STSOperator
	RoleName                 string
	Filename                 string
	RolePath                 string
	PolicyARN                string
	ManagedPolicies          bool
}

func ManualCommandsForMissingOperatorRole(input ManualCommandsForMissingOperatorRolesInput) []string {
	commands := make([]string, 0)
	iamTags := map[string]string{
		tags.ClusterID:         input.ClusterID,
		tags.RolePrefix:        input.OperatorRolePolicyPrefix,
		tags.OperatorNamespace: input.Operator.Namespace(),
		tags.OperatorName:      input.Operator.Name(),
		tags.RedHatManaged:     "true",
	}
	if input.ManagedPolicies {
		iamTags[common.ManagedPolicies] = "true"
	}

	createRole := awscb.NewIAMCommandBuilder().
		SetCommand(awscb.CreateRole).
		AddParam(awscb.RoleName, input.RoleName).
		AddParam(awscb.AssumeRolePolicyDocument, fmt.Sprintf("file://%s", input.Filename)).
		AddTags(iamTags).
		AddParam(awscb.Path, input.RolePath).
		Build()
	attachRolePolicy := awscb.NewIAMCommandBuilder().
		SetCommand(awscb.AttachRolePolicy).
		AddParam(awscb.RoleName, input.RoleName).
		AddParam(awscb.PolicyArn, input.PolicyARN).
		Build()
	commands = append(commands, createRole, attachRolePolicy)
	return commands
}

type ManualCommandsForUpgradeOperatorRolePolicyInput struct {
	PolicyExists             bool
	OperatorRolePolicyPrefix string
	Operator                 *cmv1.STSOperator
	CredRequest              string
	OperatorPolicyPath       string
	PolicyARN                string
	DefaultPolicyVersion     string
	PolicyName               string
	OperatorRoleName         string
	FileName                 string
}

func ManualCommandsForUpgradeOperatorRolePolicy(input ManualCommandsForUpgradeOperatorRolePolicyInput) []string {
	commands := make([]string, 0)
	attachRolePolicy := awscb.NewIAMCommandBuilder().
		SetCommand(awscb.AttachRolePolicy).
		AddParam(awscb.RoleName, input.OperatorRoleName).
		AddParam(awscb.PolicyArn, input.PolicyARN).
		Build()
	if !input.PolicyExists {
		iamTags := map[string]string{
			common.OpenShiftVersion: input.DefaultPolicyVersion,
			tags.RolePrefix:         input.OperatorRolePolicyPrefix,
			tags.OperatorNamespace:  input.Operator.Namespace(),
			tags.OperatorName:       input.Operator.Name(),
			tags.RedHatManaged:      "true",
		}
		createPolicy := awscb.NewIAMCommandBuilder().
			SetCommand(awscb.CreatePolicy).
			AddParam(awscb.PolicyName, input.PolicyName).
			AddParam(awscb.PolicyDocument, fmt.Sprintf("file://%s", input.FileName)).
			AddTags(iamTags).
			AddParam(awscb.Path, input.OperatorPolicyPath).
			Build()
		commands = append(commands, createPolicy)
		if input.OperatorRoleName != "" {
			commands = append(commands, attachRolePolicy)
		}
	} else {
		policyTags := map[string]string{
			common.OpenShiftVersion: input.DefaultPolicyVersion,
		}

		createPolicyVersion := awscb.NewIAMCommandBuilder().
			SetCommand(awscb.CreatePolicyVersion).
			AddParam(awscb.PolicyArn, input.PolicyARN).
			AddParam(awscb.PolicyDocument, fmt.Sprintf("file://%s", input.FileName)).
			AddParamNoValue(awscb.SetAsDefault).
			Build()

		tagPolicy := awscb.NewIAMCommandBuilder().
			SetCommand(awscb.TagPolicy).
			AddTags(policyTags).
			AddParam(awscb.PolicyArn, input.PolicyARN).
			Build()
		if input.OperatorRoleName != "" {
			commands = append(commands, attachRolePolicy)
		}
		commands = append(commands, createPolicyVersion, tagPolicy)
	}
	return commands
}

type ManualCommandsForUpgradeAccountRolePolicyInput struct {
	DefaultPolicyVersion string
	RoleName             string
	PolicyExists         bool
	Prefix               string
	File                 string
	PolicyName           string
	AccountPolicyPath    string
	PolicyARN            string
}

func ManualCommandsForUpgradeAccountRolePolicy(input ManualCommandsForUpgradeAccountRolePolicyInput) []string {
	commands := make([]string, 0)
	iamRoleTags := map[string]string{
		common.OpenShiftVersion: input.DefaultPolicyVersion,
	}

	tagRole := awscb.NewIAMCommandBuilder().
		SetCommand(awscb.TagRole).
		AddTags(iamRoleTags).
		AddParam(awscb.RoleName, input.RoleName).
		Build()

	attachRolePolicy := awscb.NewIAMCommandBuilder().
		SetCommand(awscb.AttachRolePolicy).
		AddParam(awscb.RoleName, input.RoleName).
		AddParam(awscb.PolicyArn, input.PolicyARN).
		Build()
	if !input.PolicyExists {
		iamTags := map[string]string{
			common.OpenShiftVersion: input.DefaultPolicyVersion,
			tags.RolePrefix:         input.Prefix,
			tags.RoleType:           input.File,
			tags.RedHatManaged:      "true",
		}
		createPolicy := awscb.NewIAMCommandBuilder().
			SetCommand(awscb.CreatePolicy).
			AddParam(awscb.PolicyName, input.PolicyName).
			AddParam(awscb.PolicyDocument, fmt.Sprintf("file://sts_%s_permission_policy.json", input.File)).
			AddTags(iamTags).
			AddParam(awscb.Path, input.AccountPolicyPath).
			Build()
		commands = append(commands, createPolicy, attachRolePolicy, tagRole)
	} else {
		createPolicyVersion := awscb.NewIAMCommandBuilder().
			SetCommand(awscb.CreatePolicyVersion).
			AddParam(awscb.PolicyArn, input.PolicyARN).
			AddParam(awscb.PolicyDocument, fmt.Sprintf("file://sts_%s_permission_policy.json", input.File)).
			AddParamNoValue(awscb.SetAsDefault).
			Build()

		tagPolicies := awscb.NewIAMCommandBuilder().
			SetCommand(awscb.TagPolicy).
			AddTags(iamRoleTags).
			AddParam(awscb.PolicyArn, input.PolicyARN).
			Build()
		commands = append(commands, attachRolePolicy, createPolicyVersion, tagPolicies, tagRole)
	}
	return commands
}

type ManualCommandsForDetachRolePolicyInput struct {
	RoleName  string
	PolicyARN string
}

func ManualCommandsForDetachRolePolicy(input ManualCommandsForDetachRolePolicyInput) string {
	return awscb.NewIAMCommandBuilder().
		SetCommand(awscb.DetachRolePolicy).
		AddParam(awscb.RoleName, input.RoleName).
		AddParam(awscb.PolicyArn, input.PolicyARN).
		Build()
}
