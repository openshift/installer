/**
Copyright (c) 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aws

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	semver "github.com/hashicorp/go-version"
	common "github.com/openshift-online/ocm-common/pkg/aws/validations"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	errors "github.com/zgalor/weberr"

	"github.com/openshift/rosa/pkg/aws/tags"
	"github.com/openshift/rosa/pkg/helper"
)

var DefaultPrefix = "ManagedOpenShift"

type Operator struct {
	Name                string
	Namespace           string
	RoleARN             string
	ServiceAccountNames []string
	MinVersion          string
}

type AccountRole struct {
	Name string
	Flag string
}

type Role struct {
	RoleType      string `json:"RoleType,omitempty"`
	Version       string `json:"Version,omitempty"`
	RolePrefix    string `json:"RolePrefix,omitempty"`
	RoleName      string `json:"RoleName,omitempty"`
	RoleARN       string `json:"RoleARN,omitempty"`
	Linked        string `json:"Linked,omitempty"`
	Admin         string `json:"Admin,omitempty"`
	ManagedPolicy bool   `json:"ManagedPolicy,omitempty"`
	ClusterID     string `json:"ClusterID,omitempty"`
}

type OperatorRoleDetail struct {
	OperatorName      string   `json:"Name,omitempty"`
	OperatorNamespace string   `json:"Namespace,omitempty"`
	Version           string   `json:"Version,omitempty"`
	RoleName          string   `json:"RoleName,omitempty"`
	RoleARN           string   `json:"RoleARN,omitempty"`
	ClusterID         string   `json:"ClusterID,omitempty"`
	AttachedPolicies  []string `json:"Policy,omitempty"`
	ManagedPolicy     bool     `json:"ManagedPolicy,omitempty"`
}

type PolicyDetail struct {
	PolicyName string
	PolicyArn  string
	PolicyType string
}

type Policy struct {
	PolicyName     string         `json:"PolicyName,omitempty"`
	PolicyDocument PolicyDocument `json:"PolicyDocument,omitempty"`
}

const (
	InstallerAccountRole = "installer"

	InstallerAccountRoleType = "Installer"
	ControlPlaneAccountRole  = "instance_controlplane"

	ControlPlaneAccountRoleType = "Control plane"
	WorkerAccountRole           = "instance_worker"

	WorkerAccountRoleType = "Worker"

	SupportAccountRole = "support"

	SupportAccountRoleType = "Support"

	HCPInstallerRole = "installer"
	HCPWorkerRole    = "instance_worker"
	HCPSupportRole   = "support"

	OCMRole     = "OCM"
	OCMUserRole = "User"

	// AWS preferred suffix for ROSA related account roles - HCP only
	HCPSuffixPattern = "HCP-ROSA"

	IngressOperatorCloudCredentialsRoleType = "ingress_operator_cloud_credentials"
)

const (
	InstallerCoreKey        = "sts_installer_core_permission_policy"
	InstallerVPCKey         = "sts_installer_vpc_permission_policy"
	InstallerPrivateLinkKey = "sts_installer_privatelink_permission_policy"
)

var AccountRoles = map[string]AccountRole{
	InstallerAccountRole:    {Name: "Installer", Flag: "role-arn"},
	ControlPlaneAccountRole: {Name: "ControlPlane", Flag: "controlplane-iam-role"},
	WorkerAccountRole:       {Name: "Worker", Flag: "worker-iam-role"},
	SupportAccountRole:      {Name: "Support", Flag: "support-role-arn"},
}

var HCPAccountRoles = map[string]AccountRole{
	HCPInstallerRole: {Name: fmt.Sprintf("%s-Installer", HCPSuffixPattern), Flag: "role-arn"},
	HCPSupportRole:   {Name: fmt.Sprintf("%s-Support", HCPSuffixPattern), Flag: "support-role-arn"},
	HCPWorkerRole:    {Name: fmt.Sprintf("%s-Worker", HCPSuffixPattern), Flag: "worker-iam-role"},
}

var OCMUserRolePolicyFile = "ocm_user"
var OCMRolePolicyFile = "ocm"
var OCMAdminRolePolicyFile = "ocm_admin"

var roleTypeMap = map[string]string{
	InstallerAccountRole:    InstallerAccountRoleType,
	SupportAccountRole:      SupportAccountRoleType,
	ControlPlaneAccountRole: ControlPlaneAccountRoleType,
	WorkerAccountRole:       WorkerAccountRoleType,
}

func (c *awsClient) EnsureRole(name string, policy string, permissionsBoundary string,
	version string, tagList map[string]string, path string, managedPolicies bool) (string, error) {
	output, err := c.iamClient.GetRole(&iam.GetRoleInput{
		RoleName: aws.String(name),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return c.createRole(name, policy, permissionsBoundary, tagList, path)
			default:
				return "", err
			}
		}
	}

	if managedPolicies && !common.IsManagedRole(FromV1TagToV2Tag(output.Role.Tags)) {
		return "", fmt.Errorf("Role '%s' with unmanaged policies already exists", *output.Role.Arn)
	}

	outputPath, err := GetPathFromARN(aws.StringValue(output.Role.Arn))
	if err != nil {
		return "", err
	}
	if outputPath != path {
		return "", fmt.Errorf("Role with same name but different path exists. Existing role ARN: %s",
			*output.Role.Arn)
	}

	if permissionsBoundary != "" {
		_, err = c.iamClient.PutRolePermissionsBoundary(&iam.PutRolePermissionsBoundaryInput{
			RoleName:            aws.String(name),
			PermissionsBoundary: aws.String(permissionsBoundary),
		})
	} else if output.Role.PermissionsBoundary != nil {
		_, err = c.iamClient.DeleteRolePermissionsBoundary(&iam.DeleteRolePermissionsBoundaryInput{
			RoleName: aws.String(name),
		})
	}
	if err != nil {
		return "", err
	}

	role := output.Role
	roleArn := aws.StringValue(role.Arn)

	isCompatible, err := c.isRoleCompatible(name, version)
	if err != nil {
		return roleArn, err
	}

	policy, needsUpdate, err := updateAssumeRolePolicyPrincipals(policy, role)
	if err != nil {
		return roleArn, err
	}

	if needsUpdate || !isCompatible {
		_, err = c.iamClient.UpdateAssumeRolePolicy(&iam.UpdateAssumeRolePolicyInput{
			RoleName:       aws.String(name),
			PolicyDocument: aws.String(policy),
		})
		if err != nil {
			return roleArn, err
		}

		_, err = c.iamClient.TagRole(&iam.TagRoleInput{
			RoleName: aws.String(name),
			Tags:     getTags(tagList),
		})
		if err != nil {
			return roleArn, err
		}
	}

	return roleArn, nil
}

func (c *awsClient) ValidateRoleNameAvailable(name string) (err error) {
	_, err = c.iamClient.GetRole(&iam.GetRoleInput{
		RoleName: aws.String(name),
	})
	if err == nil {
		// If we found an existing role with this name we want to error
		return fmt.Errorf("A role named '%s' already exists. "+
			"Please delete the existing role, or provide a different prefix.\n"+
			"If you'd like to reuse the operator roles, please provide a "+
			"OIDC Configuration ID which has Issuer URL linked as the trusted relationship "+
			"of the chosen operator roles prefix.", name)
	}

	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case iam.ErrCodeNoSuchEntityException:
			// This is what we want
			return nil
		}
	}
	return fmt.Errorf("Error validating role name '%s': %v", name, err)
}

func (c *awsClient) createRole(name string, policy string, permissionsBoundary string,
	tagList map[string]string, path string) (string, error) {
	if !RoleNameRE.MatchString(name) {
		return "", fmt.Errorf("Role name is invalid")
	}
	createRoleInput := &iam.CreateRoleInput{
		RoleName:                 aws.String(name),
		AssumeRolePolicyDocument: aws.String(policy),
		Tags:                     getTags(tagList),
	}
	if path != "" {
		createRoleInput.Path = aws.String(path)
	}
	if permissionsBoundary != "" {
		createRoleInput.PermissionsBoundary = aws.String(permissionsBoundary)
	}
	output, err := c.iamClient.CreateRole(createRoleInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeEntityAlreadyExistsException:
				return "", nil
			}
		}
		return "", err
	}
	return aws.StringValue(output.Role.Arn), nil
}

func (c *awsClient) isRoleCompatible(name string, version string) (bool, error) {
	// Ignore if there is no version
	if version == "" {
		return true, nil
	}
	output, err := c.iamClient.ListRoleTags(&iam.ListRoleTagsInput{
		RoleName: aws.String(name),
	})
	if err != nil {
		return false, err
	}

	return c.hasCompatibleMajorMinorVersionTags(output.Tags, version)
}

func (c *awsClient) PutRolePolicy(roleName string, policyName string, policy string) error {
	_, err := c.iamClient.PutRolePolicy(&iam.PutRolePolicyInput{
		RoleName:       aws.String(roleName),
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(policy),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *awsClient) ForceEnsurePolicy(policyArn string, document string,
	version string, tagList map[string]string, path string) (string, error) {
	return c.ensurePolicyHelper(policyArn, document, version, tagList, path, true)
}

func (c *awsClient) EnsurePolicy(policyArn string, document string,
	version string, tagList map[string]string, path string) (string, error) {
	return c.ensurePolicyHelper(policyArn, document, version, tagList, path, false)
}

func (c *awsClient) ensurePolicyHelper(policyArn string, document string,
	version string, tagList map[string]string, path string, force bool) (string, error) {
	output, err := c.IsPolicyExists(policyArn)
	if err != nil {
		var policyArnLocal string
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				policyArnLocal, err = c.createPolicy(policyArn, document, tagList, path)
				if err != nil {
					if aerr, ok := err.(awserr.Error); ok {
						if aerr.Code() == iam.ErrCodeEntityAlreadyExistsException {
							return "", errors.Wrapf(err,
								fmt.Sprintf("Failed to create a policy with ARN '%s'", policyArn))
						}
					}
					return "", err
				}
				return policyArnLocal, nil
			default:
				return "", err
			}
		}
	}

	policyArn = aws.StringValue(output.Policy.Arn)

	isCompatible := false
	if !force {
		isCompatible, err = c.IsPolicyCompatible(policyArn, version)
		if err != nil {
			return policyArn, err
		}
	}

	if !isCompatible {
		// Since there is a limit to how many versions a policy can have, we delete all non-default
		// policy versions from the list, thus making space for the new one.
		err = c.deletePolicyVersions(policyArn)
		if err != nil {
			return policyArn, err
		}

		_, err = c.iamClient.CreatePolicyVersion(&iam.CreatePolicyVersionInput{
			PolicyArn:      aws.String(policyArn),
			PolicyDocument: aws.String(document),
			SetAsDefault:   aws.Bool(true),
		})
		if err != nil {
			return policyArn, err
		}

		_, err = c.iamClient.TagPolicy(&iam.TagPolicyInput{
			PolicyArn: aws.String(policyArn),
			Tags:      getTags(tagList),
		})
		if err != nil {
			return policyArn, err
		}
	}

	return policyArn, nil
}

func (c *awsClient) IsPolicyExists(policyArn string) (*iam.GetPolicyOutput, error) {
	output, err := c.iamClient.GetPolicy(&iam.GetPolicyInput{
		PolicyArn: aws.String(policyArn),
	})
	return output, err
}

func (c *awsClient) IsRolePolicyExists(roleName string, policyName string) (*iam.GetRolePolicyOutput, error) {
	output, err := c.iamClient.GetRolePolicy(&iam.GetRolePolicyInput{
		PolicyName: aws.String(policyName),
		RoleName:   aws.String(roleName),
	})
	return output, err
}

func (c *awsClient) createPolicy(policyArn string, document string, tagList map[string]string,
	path string) (string, error) {
	policyName, err := GetResourceIdFromARN(policyArn)
	if err != nil {
		return "", err
	}
	createPolicyInput := &iam.CreatePolicyInput{
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(document),
		Tags:           getTags(tagList),
	}
	if path != "" {
		createPolicyInput.Path = aws.String(path)
	}

	output, err := c.iamClient.CreatePolicy(createPolicyInput)

	if err != nil {
		return "", err
	}
	return aws.StringValue(output.Policy.Arn), nil
}

func (c *awsClient) IsPolicyCompatible(policyArn string, version string) (bool, error) {
	output, err := c.iamClient.ListPolicyTags(&iam.ListPolicyTagsInput{
		PolicyArn: aws.String(policyArn),
	})
	if err != nil {
		return false, err
	}

	return common.HasCompatibleVersionTags(FromV1TagToV2Tag(output.Tags), version)
}

func (c *awsClient) hasCompatibleMajorMinorVersionTags(iamTags []*iam.Tag, version string) (bool, error) {
	if len(iamTags) == 0 {
		return false, nil
	}
	for _, tag := range iamTags {
		if aws.StringValue(tag.Key) == common.OpenShiftVersion {
			if version == aws.StringValue(tag.Value) {
				return true, nil
			}

			upgradeVersion, err := semver.NewVersion(version)
			if err != nil {
				return false, err
			}

			currentVersion, err := semver.NewVersion(aws.StringValue(tag.Value))
			if err != nil {
				return false, err
			}

			upgradeVersionSegments := upgradeVersion.Segments64()
			c, err := semver.NewConstraint(fmt.Sprintf(">= %d.%d",
				upgradeVersionSegments[0], upgradeVersionSegments[1]))
			if err != nil {
				return false, err
			}
			return c.Check(currentVersion), nil
		}
	}
	return false, nil
}

func (c *awsClient) AttachRolePolicy(roleName string, policyARN string) error {
	_, err := c.iamClient.AttachRolePolicy(&iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(policyARN),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *awsClient) FindRoleARNs(roleType string, version string) ([]string, error) {
	roleARNs := []string{}
	roles, err := c.ListRoles()
	if err != nil {
		return roleARNs, err
	}
	for _, role := range roles {
		if !strings.Contains(aws.StringValue(role.RoleName), AccountRoles[roleType].Name) {
			continue
		}
		isValid, err := c.ValidateAccountRoleVersionCompatibility(*role.RoleName, roleType, version)
		if err != nil {
			return roleARNs, err
		}
		if !isValid {
			continue
		}
		roleARNs = append(roleARNs, aws.StringValue(role.Arn))
	}
	return roleARNs, nil
}

func (c *awsClient) FindRoleARNsClassic(roleType string, version string) ([]string, error) {
	roleARNs := []string{}
	roles, err := c.ListRoles()
	if err != nil {
		return roleARNs, err
	}
	for _, role := range roles {
		if !strings.Contains(aws.StringValue(role.RoleName), AccountRoles[roleType].Name) {
			continue
		}
		listRoleTagsOutput, err := c.iamClient.ListRoleTags(&iam.ListRoleTagsInput{
			RoleName: role.RoleName,
		})
		if err != nil {
			return roleARNs, err
		}
		isValid, err := validateAccountRoleVersionCompatibilityClassic(roleType, version, listRoleTagsOutput.Tags)
		if err != nil {
			return roleARNs, err
		}
		if !isValid {
			continue
		}
		roleARNs = append(roleARNs, aws.StringValue(role.Arn))
	}
	return roleARNs, nil
}

func (c *awsClient) FindRoleARNsHostedCp(roleType string, version string) ([]string, error) {
	roleARNs := []string{}
	roles, err := c.ListRoles()
	if err != nil {
		return roleARNs, err
	}
	for _, role := range roles {
		if !strings.Contains(aws.StringValue(role.RoleName), AccountRoles[roleType].Name) {
			continue
		}
		listRoleTagsOutput, err := c.iamClient.ListRoleTags(&iam.ListRoleTagsInput{
			RoleName: role.RoleName,
		})
		if err != nil {
			return roleARNs, err
		}
		isValid, err := validateAccountRoleVersionCompatibilityHostedCp(roleType, version, listRoleTagsOutput.Tags)
		if err != nil {
			return roleARNs, err
		}
		if !isValid {
			continue
		}
		roleARNs = append(roleARNs, aws.StringValue(role.Arn))
	}
	return roleARNs, nil
}

// FIXME: refactor similar calls to use this instead
func (c *awsClient) ValidateAccountRoleVersionCompatibility(roleName string, roleType string,
	minVersion string) (bool, error) {
	listRoleTagsOutput, err := c.iamClient.ListRoleTags(&iam.ListRoleTagsInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return false, err
	}

	return isAccountRoleVersionCompatible(listRoleTagsOutput.Tags, roleType, minVersion)
}

func validateAccountRoleVersionCompatibilityClassic(roleType string, minVersion string,
	tagList []*iam.Tag) (bool, error) {
	isCompatible, err := isAccountRoleVersionCompatible(tagList, roleType, minVersion)
	if err != nil {
		return false, err
	}
	if !isCompatible {
		return false, nil
	}

	// Account roles with HCP policies are not compatible with classic clusters
	if common.IamResourceHasTag(FromV1TagToV2Tag(tagList), tags.HypershiftPolicies, tags.True) {
		return false, nil
	}

	return true, nil
}

// FIXME: Remove when moving ROSA CLI to aws v2
func FromV1TagToV2Tag(tags []*iam.Tag) []iamtypes.Tag {
	v2TagList := []iamtypes.Tag{}
	for _, tag := range tags {
		v2TagList = append(v2TagList, iamtypes.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}
	return v2TagList
}

func validateAccountRoleVersionCompatibilityHostedCp(roleType string, minVersion string,
	tagsList []*iam.Tag) (bool, error) {
	isCompatible, err := isAccountRoleVersionCompatible(tagsList, roleType, minVersion)
	if err != nil {
		return false, err
	}
	if !isCompatible {
		return false, nil
	}

	// Only account roles with HCP managed policies are compatible with HCP clusters
	return common.IamResourceHasTag(FromV1TagToV2Tag(tagsList), tags.HypershiftPolicies, tags.True), nil
}

func isAccountRoleVersionCompatible(tagsList []*iam.Tag, roleType string,
	minVersion string) (bool, error) {
	skip := false
	isTagged := false
	for _, tag := range tagsList {
		tagValue := aws.StringValue(tag.Value)
		switch aws.StringValue(tag.Key) {
		case tags.RoleType:
			isTagged = true
			if tagValue != roleType {
				skip = true
				break
			}
		case common.OpenShiftVersion:
			isTagged = true

			if common.IamResourceHasTag(FromV1TagToV2Tag(tagsList), common.ManagedPolicies, tags.True) {
				// Managed policies will be up-to-date no need to check version tags
				break
			}

			if minVersion == "" {
				break
			}
			minExpectedVersion, err := semver.NewVersion(minVersion)
			if err != nil {
				skip = true
				break
			}
			policyVersion, err := semver.NewVersion(tagValue)
			if err != nil {
				skip = true
				break
			}
			if policyVersion.LessThan(minExpectedVersion) {
				skip = true
				break
			}
		}
	}
	if !isTagged || skip {
		return false, nil
	}

	return true, nil
}

func (c *awsClient) ListRoles() ([]*iam.Role, error) {
	roles := []*iam.Role{}
	err := c.iamClient.ListRolesPages(&iam.ListRolesInput{}, func(page *iam.ListRolesOutput, lastPage bool) bool {
		roles = append(roles, page.Roles...)
		return aws.BoolValue(page.IsTruncated)
	})
	return roles, err
}

func (c *awsClient) FindPolicyARN(operator Operator, version string) (string, error) {
	policies := []*iam.Policy{}
	err := c.iamClient.ListPoliciesPages(&iam.ListPoliciesInput{
		Scope: aws.String(iam.PolicyScopeTypeLocal),
	}, func(page *iam.ListPoliciesOutput, lastPage bool) bool {
		policies = append(policies, page.Policies...)
		return aws.BoolValue(page.IsTruncated)
	})
	if err != nil {
		return "", err
	}
	for _, policy := range policies {
		listPolicyTagsOutput, err := c.iamClient.ListPolicyTags(&iam.ListPolicyTagsInput{
			PolicyArn: policy.Arn,
		})
		if err != nil {
			return "", err
		}
		skip := false
		isTagged := false
		for _, tag := range listPolicyTagsOutput.Tags {
			tagValue := aws.StringValue(tag.Value)
			switch aws.StringValue(tag.Key) {
			case tags.OperatorNamespace:
				isTagged = true
				if tagValue != operator.Namespace {
					skip = true
					break
				}
			case tags.OperatorName:
				isTagged = true
				if tagValue != operator.Name {
					skip = true
					break
				}
			case common.OpenShiftVersion:
				isTagged = true
				if tagValue != version {
					skip = true
					break
				}
			}
		}
		if isTagged && !skip {
			return aws.StringValue(policy.Arn), nil
		}
	}
	return "", nil
}

func getTags(tagList map[string]string) []*iam.Tag {
	iamTags := []*iam.Tag{}
	for k, v := range tagList {
		iamTags = append(iamTags, &iam.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return iamTags
}

func IsOCMRole(roleName *string) bool {
	return strings.Contains(aws.StringValue(roleName), fmt.Sprintf("%s-Role", OCMRole))
}

// IsUserRole checks the role tags in addition to the role name, because the word 'user' is common
func (c *awsClient) IsUserRole(roleName *string) (bool, error) {
	if strings.Contains(aws.StringValue(roleName), OCMUserRole) {
		roleTags, err := c.iamClient.ListRoleTags(&iam.ListRoleTagsInput{
			RoleName: roleName,
		})
		if err != nil {
			return false, err
		}

		return common.IamResourceHasTag(FromV1TagToV2Tag(roleTags.Tags), tags.RoleType, OCMUserRole), nil
	}

	return false, nil
}

func (c *awsClient) ListUserRoles() ([]Role, error) {
	var userRoles []Role
	roles, err := c.ListRoles()
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		isUserRole, err := c.IsUserRole(role.RoleName)
		if err != nil {
			return nil, err
		}

		if isUserRole {
			var userRole Role
			userRole.RoleName = aws.StringValue(role.RoleName)
			userRole.RoleARN = aws.StringValue(role.Arn)

			userRoles = append(userRoles, userRole)
		}
	}

	return userRoles, nil
}

func (c *awsClient) ListOCMRoles() ([]Role, error) {
	var ocmRoles []Role
	roles, err := c.ListRoles()
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if IsOCMRole(role.RoleName) {
			var ocmRole Role
			ocmRole.RoleName = aws.StringValue(role.RoleName)
			ocmRole.RoleARN = aws.StringValue(role.Arn)

			roleTags, err := c.iamClient.ListRoleTags(&iam.ListRoleTagsInput{
				RoleName: role.RoleName,
			})
			if err != nil {
				return nil, err
			}
			v2Tags := FromV1TagToV2Tag(roleTags.Tags)
			if common.IamResourceHasTag(v2Tags, tags.AdminRole, tags.True) {
				ocmRole.Admin = "Yes"
			} else {
				ocmRole.Admin = "No"
			}
			if common.IamResourceHasTag(v2Tags, common.ManagedPolicies, tags.True) {
				ocmRole.ManagedPolicy = true
			}

			ocmRoles = append(ocmRoles, ocmRole)
		}
	}

	return ocmRoles, nil
}

func (c *awsClient) GetAccountRoleByArn(arn string) (*Role, error) {
	role, err := c.GetRoleByARN(arn)
	if err != nil {
		return nil, err
	}

	accountRole, err := c.mapToAccountRole("", role)

	if err != nil {
		return nil, err
	}

	return accountRole, nil
}

func (c *awsClient) mapToAccountRole(version string, role *iam.Role) (*Role, error) {
	if !checkIfAccountRole(role.RoleName) {
		return nil, nil
	}

	accountRole := Role{}

	listRoleTagsOutput, err := c.iamClient.ListRoleTags(&iam.ListRoleTagsInput{
		RoleName: role.RoleName,
	})
	if err != nil {
		return nil, err
	}

	for _, tag := range listRoleTagsOutput.Tags {
		switch aws.StringValue(tag.Key) {
		case tags.RoleType:
			accountRole.RoleType = roleTypeMap[aws.StringValue(tag.Value)]
		case common.OpenShiftVersion:
			tagValue := aws.StringValue(tag.Value)
			if version != "" && tagValue != version {
				return nil, nil
			}
			accountRole.Version = tagValue
		case common.ManagedPolicies:
			if aws.StringValue(tag.Value) == tags.True {
				accountRole.ManagedPolicy = true
			}
		}
	}

	accountRole.RoleName = aws.StringValue(role.RoleName)
	accountRole.RoleARN = aws.StringValue(role.Arn)

	return &accountRole, nil
}

func (c *awsClient) mapToAccountRoles(version string, roles []*iam.Role) ([]Role, error) {

	accountRoles := []Role{}
	for _, role := range roles {

		accountRole, err := c.mapToAccountRole(version, role)
		if err != nil {
			return accountRoles, err
		}

		if accountRole != nil {
			accountRoles = append(accountRoles, *accountRole)
		}
	}

	return accountRoles, nil
}

func (c *awsClient) ListAccountRoles(version string) ([]Role, error) {
	roles, err := c.ListRoles()
	if err != nil {
		return []Role{}, err
	}
	return c.mapToAccountRoles(version, roles)
}

func (c *awsClient) ListOperatorRoles(version string, targetClusterId string) (map[string][]OperatorRoleDetail, error) {
	operatorMap := map[string][]OperatorRoleDetail{}
	roles, err := c.ListRoles()
	if err != nil {
		return operatorMap, err
	}
	prefixOperatorRoleRE := regexp.MustCompile(`(?i)(?P<Prefix>[\w+=,.@-]+)-(openshift|kube-system)`)
	for _, role := range roles {
		operatorRole := OperatorRoleDetail{}
		matches := prefixOperatorRoleRE.FindStringSubmatch(*role.RoleName)
		if len(matches) == 0 {
			continue
		}
		prefixIndex := prefixOperatorRoleRE.SubexpIndex("Prefix")
		foundPrefix := strings.ToLower(matches[prefixIndex])
		if _, mapOk := operatorMap[foundPrefix]; !mapOk {
			operatorMap[foundPrefix] = []OperatorRoleDetail{}
		}
		listRoleTagsOutput, err := c.iamClient.ListRoleTags(&iam.ListRoleTagsInput{
			RoleName: role.RoleName,
		})
		if err != nil {
			return operatorMap, err
		}
		skip := false
		for _, tag := range listRoleTagsOutput.Tags {
			switch aws.StringValue(tag.Key) {
			case common.ManagedPolicies:
				if aws.StringValue(tag.Value) == tags.True {
					operatorRole.ManagedPolicy = true
				}
			case tags.ClusterID:
				tagValue := aws.StringValue(tag.Value)
				if targetClusterId != "" && tagValue != targetClusterId {
					skip = true
				}
				operatorRole.ClusterID = tagValue
			case tags.OperatorName:
				operatorRole.OperatorName = *tag.Value

			case tags.OperatorNamespace:
				operatorRole.OperatorNamespace = *tag.Value
			}
		}

		attachedPoliciesOutput, err := c.iamClient.ListAttachedRolePolicies(&iam.ListAttachedRolePoliciesInput{
			RoleName: role.RoleName,
		})
		if err != nil {
			return operatorMap, err
		}

		attachedPolicies := []string{}

		for _, policy := range attachedPoliciesOutput.AttachedPolicies {
			attachedPolicies = append(attachedPolicies, (aws.StringValue(policy.PolicyName)))
		}
		operatorRole.AttachedPolicies = attachedPolicies

		if skip {
			continue
		}

		if operatorRole.ManagedPolicy || len(attachedPoliciesOutput.AttachedPolicies) == 0 {
			operatorRole.RoleName = aws.StringValue(role.RoleName)
			operatorRole.RoleARN = aws.StringValue(role.Arn)
			operatorMap[foundPrefix] = append(operatorMap[foundPrefix], operatorRole)
			continue
		}

		for _, policy := range attachedPoliciesOutput.AttachedPolicies {
			listPolicyTagsOutput, err := c.iamClient.ListPolicyTags(&iam.ListPolicyTagsInput{
				PolicyArn: policy.PolicyArn,
			})
			if err != nil {
				return operatorMap, err
			}
			isTagged := false
			skip := false
			for _, tag := range listPolicyTagsOutput.Tags {
				switch aws.StringValue(tag.Key) {
				case common.OpenShiftVersion:
					tagValue := aws.StringValue(tag.Value)
					if version != "" && tagValue != version {
						skip = true
						break
					}
					isTagged = true
					operatorRole.Version = tagValue
				}
			}
			if isTagged && !skip {
				operatorRole.RoleName = aws.StringValue(role.RoleName)
				operatorRole.RoleARN = aws.StringValue(role.Arn)
				operatorMap[foundPrefix] = append(operatorMap[foundPrefix], operatorRole)
			}
		}
	}
	emptyListKeys := []string{}
	for key, list := range operatorMap {
		if len(list) == 0 {
			emptyListKeys = append(emptyListKeys, key)
		} else if targetClusterId != "" && list[0].ClusterID != targetClusterId {
			emptyListKeys = append(emptyListKeys, key)
		}
	}
	for _, key := range emptyListKeys {
		delete(operatorMap, key)
	}
	return operatorMap, nil
}

// Check if it is one of the ROSA account roles
func checkIfAccountRole(roleName *string) bool {
	for _, prefix := range AccountRoles {
		if strings.Contains(aws.StringValue(roleName), common.GetRoleName("", prefix.Name)) {
			return true
		}
	}
	return false
}

// Check if it is one of the ROSA account roles
func checkIfROSAOperatorRole(roleName *string, credRequest map[string]*cmv1.STSOperator) bool {
	for _, operatorRole := range credRequest {
		if strings.Contains(aws.StringValue(roleName), operatorRole.Namespace()) {
			return true
		}
	}
	return false
}

func (c *awsClient) DeleteOperatorRole(roleName string, managedPolicies bool) error {
	role := aws.String(roleName)
	policies, err := c.GetPolicies([]string{*role})
	if err != nil {
		return err
	}
	err = c.detachOperatorRolePolicies(role)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				fmt.Printf("Entity does not exist: %s", aerr)
				err = nil
			case iam.ErrCodeDeleteConflictException:
				fmt.Printf("Unable to detach operator role policy: %s", aerr)
				err = nil
			}
		}
		if err != nil {
			return err
		}
	}
	err = c.DeleteRole(*role)
	if err != nil {
		return err
	}
	if !managedPolicies {
		_, err = c.deletePolicies(policies[*role])
	}
	return err
}

func (c *awsClient) DeleteRole(role string) error {
	_, err := c.iamClient.DeleteRole(&iam.DeleteRoleInput{RoleName: aws.String(role)})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == iam.ErrCodeNoSuchEntityException {
				return fmt.Errorf("operator role '%s' does not exist. Skipping", role)
			}
		}
		return err
	}
	return nil
}

func (c *awsClient) GetInstanceProfilesForRole(r string) ([]string, error) {
	instanceProfiles := []string{}
	profiles, err := c.iamClient.ListInstanceProfilesForRole(&iam.ListInstanceProfilesForRoleInput{
		RoleName: aws.String(r),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return instanceProfiles, nil
			}
		}
		return nil, err
	}
	for _, profile := range profiles.InstanceProfiles {
		instanceProfiles = append(instanceProfiles, aws.StringValue(profile.InstanceProfileName))
	}
	return instanceProfiles, nil
}

func (c *awsClient) DeleteAccountRole(roleName string, managedPolicies bool) error {
	role := aws.String(roleName)
	err := c.DeleteInlineRolePolicies(aws.StringValue(role))
	if err != nil {
		return err
	}
	policyMap, err := c.GetPolicies([]string{*role})
	if err != nil {
		return err
	}
	err = c.detachAttachedRolePolicies(role)
	if err != nil {
		return err
	}
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				fmt.Printf("Entity does not exist: %s", aerr)
				err = nil
			case iam.ErrCodeDeleteConflictException:
				fmt.Printf("Unable to detach account role policy: %s", aerr)
				err = nil
			}
		}
		if err != nil {
			return err
		}
	}
	err = c.DeleteRole(*role)
	if err != nil {
		return err
	}
	if !managedPolicies {
		_, err = c.deletePolicies(policyMap[*role])
	}
	return err
}

func (c *awsClient) detachAttachedRolePolicies(role *string) error {
	attachedPoliciesOutput, err := c.iamClient.ListAttachedRolePolicies(&iam.ListAttachedRolePoliciesInput{
		RoleName: role,
	})
	if err != nil {
		return err
	}
	for _, policy := range attachedPoliciesOutput.AttachedPolicies {
		_, err = c.iamClient.DetachRolePolicy(&iam.DetachRolePolicyInput{
			PolicyArn: policy.PolicyArn,
			RoleName:  role,
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					continue
				}
			}
			return err
		}
	}

	return nil
}

func (c *awsClient) DeleteInlineRolePolicies(role string) error {
	listRolePolicyOutput, err := c.iamClient.ListRolePolicies(&iam.ListRolePoliciesInput{RoleName: aws.String(role)})
	if err != nil {
		return err
	}
	for _, policyName := range listRolePolicyOutput.PolicyNames {
		_, err = c.iamClient.DeleteRolePolicy(&iam.DeleteRolePolicyInput{
			PolicyName: policyName,
			RoleName:   aws.String(role),
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					continue
				}
			}
			return err
		}
	}

	return nil
}

func (c *awsClient) isPolicyAttachedToEntity(policyArn string) (bool, error) {
	policyOutput, err := c.iamClient.GetPolicy(&iam.GetPolicyInput{PolicyArn: aws.String(policyArn)})
	if err != nil {
		return false, err
	}

	if *policyOutput.Policy.AttachmentCount > 0 {
		return true, nil
	}

	return false, nil
}

func (c *awsClient) deletePolicies(policies []string) (*iam.DeletePolicyOutput, error) {
	var output *iam.DeletePolicyOutput

	for i := range policies {
		isAttached, err := c.isPolicyAttachedToEntity(policies[i])
		if err != nil {
			return output, err
		}
		if isAttached {
			continue
		}

		err = c.deletePolicyVersions(policies[i])
		if err != nil {
			return output, err
		}

		output, err = c.iamClient.DeletePolicy(&iam.DeletePolicyInput{PolicyArn: &policies[i]})
		if err != nil {
			return output, err
		}
	}
	return output, nil
}

// GetDefaultPolicyDocument gets a policy ARN and return a JSON policy document of the default policy version.
func (c *awsClient) GetDefaultPolicyDocument(policyArn string) (string, error) {
	versionId, err := c.getDefaultPolicyVersionId(policyArn)
	if err != nil {
		return "", err
	}

	policyVersionOutput, err := c.iamClient.GetPolicyVersion(&iam.GetPolicyVersionInput{
		VersionId: aws.String(versionId),
		PolicyArn: aws.String(policyArn),
	})
	if err != nil {
		return "", err
	}

	return url.QueryUnescape(aws.StringValue(policyVersionOutput.PolicyVersion.Document))
}

func (c *awsClient) getDefaultPolicyVersionId(policyArn string) (string, error) {
	policyVersionsOutput, err := c.iamClient.ListPolicyVersions(&iam.ListPolicyVersionsInput{
		PolicyArn: aws.String(policyArn),
	})
	if err != nil {
		return "", err
	}

	for _, version := range policyVersionsOutput.Versions {
		if aws.BoolValue(version.IsDefaultVersion) {
			return aws.StringValue(version.VersionId), nil
		}
	}

	// Shouldn't get here, each policy must have a default version.
	return "", errors.Errorf("Failed to find the default policy version for policy '%s'", policyArn)
}

func (c *awsClient) deletePolicyVersions(policyArn string) error {
	policyVersionsOutput, err := c.iamClient.ListPolicyVersions(&iam.ListPolicyVersionsInput{
		PolicyArn: aws.String(policyArn),
	})
	if err != nil {
		return err
	}

	for _, version := range policyVersionsOutput.Versions {
		if aws.BoolValue(version.IsDefaultVersion) {
			continue
		}
		_, err := c.iamClient.DeletePolicyVersion(&iam.DeletePolicyVersionInput{
			PolicyArn: aws.String(policyArn),
			VersionId: version.VersionId,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *awsClient) GetAttachedPolicy(role *string) ([]PolicyDetail, error) {
	policies := []PolicyDetail{}
	attachedPoliciesOutput, err := c.iamClient.ListAttachedRolePolicies(
		&iam.ListAttachedRolePoliciesInput{RoleName: role},
	)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				break
			default:
				return policies, err
			}
		} else {
			return policies, err
		}
	}

	for _, policy := range attachedPoliciesOutput.AttachedPolicies {
		policyDetail := PolicyDetail{
			PolicyName: aws.StringValue(policy.PolicyName),
			PolicyArn:  aws.StringValue(policy.PolicyArn),
			PolicyType: Attached,
		}
		policies = append(policies, policyDetail)
	}

	rolePolicyOutput, err := c.iamClient.ListRolePolicies(&iam.ListRolePoliciesInput{RoleName: role})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				break
			default:
				return policies, err
			}
		} else {
			return policies, err
		}
	}
	for _, policy := range rolePolicyOutput.PolicyNames {
		policyDetail := PolicyDetail{
			PolicyName: aws.StringValue(policy),
			PolicyType: Inline,
		}
		policies = append(policies, policyDetail)
	}

	return policies, nil
}

func (c *awsClient) detachOperatorRolePolicies(role *string) error {
	// get attached role policies as operator roles have managed policies
	policiesOutput, err := c.iamClient.ListAttachedRolePolicies(&iam.ListAttachedRolePoliciesInput{
		RoleName: role,
	})
	if err != nil {
		return err
	}
	for _, policy := range policiesOutput.AttachedPolicies {
		_, err := c.iamClient.DetachRolePolicy(&iam.DetachRolePolicyInput{PolicyArn: policy.PolicyArn, RoleName: role})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *awsClient) GetOperatorRolesFromAccountByClusterID(clusterID string,
	credRequest map[string]*cmv1.STSOperator) ([]string, error) {
	roleList := []string{}
	roles, err := c.ListRoles()
	if err != nil {
		return roleList, err
	}
	for _, role := range roles {
		if !checkIfROSAOperatorRole(role.RoleName, credRequest) {
			continue
		}
		listRoleTagsOutput, err := c.iamClient.ListRoleTags(&iam.ListRoleTagsInput{
			RoleName: role.RoleName,
		})
		if err != nil {
			return roleList, err
		}
		isTagged := false
		for _, tag := range listRoleTagsOutput.Tags {
			switch aws.StringValue(tag.Key) {
			case tags.ClusterID:
				if aws.StringValue(tag.Value) == clusterID {
					isTagged = true
					break
				}
			}
		}
		if isTagged {
			roleList = append(roleList, aws.StringValue(role.RoleName))
		}
	}
	return roleList, nil
}

func (c *awsClient) GetOperatorRolesFromAccountByPrefix(prefix string,
	credRequest map[string]*cmv1.STSOperator) ([]string, error) {
	roleList := []string{}
	roles, err := c.ListRoles()
	if err != nil {
		return roleList, err
	}
	prefixOperatorRoleRE := regexp.MustCompile(("(?i)" + fmt.Sprintf("(%s)-(openshift|kube-system)", prefix)))
	for _, role := range roles {
		if !checkIfROSAOperatorRole(role.RoleName, credRequest) {
			continue
		}
		if prefixOperatorRoleRE.MatchString(*role.RoleName) {
			roleList = append(roleList, aws.StringValue(role.RoleName))
		}
	}
	return roleList, nil
}

func (c *awsClient) GetPolicies(roles []string) (map[string][]string, error) {
	roleMap := make(map[string][]string)
	for _, role := range roles {
		policyArr := []string{}
		policiesOutput, err := c.iamClient.ListAttachedRolePolicies(&iam.ListAttachedRolePoliciesInput{
			RoleName: aws.String(role),
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					continue
				}
			}
			return roleMap, err
		}
		for _, policy := range policiesOutput.AttachedPolicies {
			policyArr = append(policyArr, aws.StringValue(policy.PolicyArn))
		}
		roleMap[role] = policyArr
	}
	return roleMap, nil
}

func (c *awsClient) GetAccountRolesForCurrentEnv(env string, accountID string) ([]Role, error) {
	roleList := []Role{}
	roles, err := c.ListRoles()
	if err != nil {
		return roleList, err
	}
	for _, role := range roles {
		if role.RoleName == nil {
			continue
		}
		if !strings.Contains(aws.StringValue(role.RoleName), ("Installer-Role")) {
			continue
		}
		policyDoc, err := getPolicyDocument(role.AssumeRolePolicyDocument)
		if err != nil {
			return roleList, err
		}
		statements := policyDoc.Statement
		for _, statement := range statements {
			awsPrincipal := statement.GetAWSPrincipals()
			if len(awsPrincipal) > 1 {
				break
			}
			for _, a := range awsPrincipal {
				str := strings.Split(a, ":")
				if len(str) > 4 {
					if str[4] == GetJumpAccount(env) {
						roles, err := c.buildRoles(aws.StringValue(role.RoleName), accountID)
						if err != nil {
							return roleList, err
						}
						roleList = append(roleList, roles...)
						break
					}
				}
			}
		}
	}
	return roleList, nil
}

func (c *awsClient) GetAccountRoleForCurrentEnv(env string, roleName string) (Role, error) {
	role := Role{}
	// This is done to ensure user did not provide invalid role before we check for installer role
	accountRoleResponse, err := c.iamClient.GetRole(&iam.GetRoleInput{RoleName: aws.String(roleName)})
	if err != nil {
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					return role, errors.NotFound.Errorf("Role '%s' not found", roleName)
				}
			}
		}
		return role, err
	}

	assumePolicyDoc := accountRoleResponse.Role.AssumeRolePolicyDocument
	if !strings.Contains(roleName, ("Installer-Role")) {
		installerRoleResponse, err := c.checkInstallerRoleExists(roleName)
		if err != nil {
			return role, err
		}
		if installerRoleResponse == nil {
			return Role{
				RoleARN:  aws.StringValue(accountRoleResponse.Role.Arn),
				RoleName: roleName,
			}, nil
		}
		assumePolicyDoc = installerRoleResponse.AssumeRolePolicyDocument
	}
	policyDoc, err := getPolicyDocument(assumePolicyDoc)
	if err != nil {
		return role, err
	}
	statements := policyDoc.Statement
	for _, statement := range statements {
		awsPrincipal := statement.GetAWSPrincipals()
		for _, a := range awsPrincipal {
			str := strings.Split(a, ":")
			if len(str) > 4 {
				if str[4] == GetJumpAccount(env) {
					r := Role{
						RoleARN:  aws.StringValue(accountRoleResponse.Role.Arn),
						RoleName: roleName,
					}
					return r, nil
				}
			}
		}
	}
	return role, nil
}

func (c *awsClient) checkInstallerRoleExists(roleName string) (*iam.Role, error) {
	rolePrefix := ""
	for _, prefix := range AccountRoles {
		p := fmt.Sprintf("%s-Role", prefix.Name)
		if strings.Contains(roleName, p) {
			rolePrefix = strings.Split(roleName, p)[0]
		}
	}
	installerRole := fmt.Sprintf("%s%s-Role", rolePrefix, "Installer")
	installerRoleResponse, err := c.iamClient.GetRole(&iam.GetRoleInput{RoleName: aws.String(installerRole)})
	//We try our best to determine the environment based on the trust policy in the installer
	//If the installer role is deleted we can assume that there is no cluster using the role
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return nil, nil
			default:
				return nil, err
			}
		}
		return nil, err
	}

	return installerRoleResponse.Role, nil
}

func (c *awsClient) GetAccountRoleForCurrentEnvWithPrefix(env string, rolePrefix string,
	accountRolesMap map[string]AccountRole) ([]Role, error) {
	roleList := []Role{}
	for _, prefix := range accountRolesMap {
		role, err := c.GetAccountRoleForCurrentEnv(env, fmt.Sprintf("%s-%s-Role", rolePrefix, prefix.Name))
		if err != nil {
			if errors.GetType(err) != errors.NotFound {
				return nil, err
			}
		}
		roleList = append(roleList, role)
	}
	return roleList, nil
}

func (c *awsClient) buildRoles(roleName string, accountID string) ([]Role, error) {
	roles := []Role{}
	rolePrefix := strings.Split(roleName, "-Installer-Role")[0]
	for _, prefix := range AccountRoles {
		roleName := fmt.Sprintf("%s-%s-Role", rolePrefix, prefix.Name)
		roleARN := GetRoleARN(accountID, roleName, "")

		if prefix.Name != "Installer" {
			_, err := c.iamClient.GetRole(&iam.GetRoleInput{RoleName: aws.String(roleName)})
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case iam.ErrCodeNoSuchEntityException:
						continue
					}
				}
				return roles, err
			}
		}
		role := Role{
			RoleARN:  roleARN,
			RoleName: roleName,
			RoleType: prefix.Name,
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (c *awsClient) GetAccountRolePolicies(roles []string) (map[string][]PolicyDetail, error) {
	roleMap := make(map[string][]PolicyDetail)
	for _, role := range roles {
		policies, err := c.GetAttachedPolicy(aws.String(role))
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					continue
				}
			}
			return roleMap, err
		}
		roleMap[role] = policies
	}
	return roleMap, nil
}

func (c *awsClient) GetOpenIDConnectProviderByClusterIdTag(clusterID string) (string, error) {
	providers, err := c.iamClient.ListOpenIDConnectProviders(&iam.ListOpenIDConnectProvidersInput{})
	if err != nil {
		return "", err
	}
	for _, provider := range providers.OpenIDConnectProviderList {
		providerValue := aws.StringValue(provider.Arn)
		connectProvider, err := c.iamClient.GetOpenIDConnectProvider(&iam.GetOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: provider.Arn,
		})
		if err != nil {
			return "", err
		}
		isTagged := false
		for _, providerTag := range connectProvider.Tags {
			switch aws.StringValue(providerTag.Key) {
			case tags.ClusterID:
				if aws.StringValue(providerTag.Value) == clusterID {
					isTagged = true
					break
				}
			}
		}
		if isTagged {
			return providerValue, nil
		}
		if strings.Contains(providerValue, clusterID) {
			return providerValue, nil
		}
	}
	return "", nil
}

func (c *awsClient) GetOpenIDConnectProviderByOidcEndpointUrl(oidcEndpointUrl string) (string, error) {
	providers, err := c.iamClient.ListOpenIDConnectProviders(&iam.ListOpenIDConnectProvidersInput{})
	if err != nil {
		return "", err
	}
	oidcEndpointUrl = strings.TrimPrefix(oidcEndpointUrl, fmt.Sprintf("%s://", helper.ProtocolHttps))
	for _, provider := range providers.OpenIDConnectProviderList {
		providerValue := aws.StringValue(provider.Arn)
		if err != nil {
			return "", err
		}
		providerName, err := GetResourceIdFromOidcProviderARN(providerValue)
		if err != nil {
			return "", err
		}
		if strings.Contains(providerName, oidcEndpointUrl) ||
			strings.Contains(oidcEndpointUrl, providerName) {
			return providerValue, nil
		}
	}
	return "", nil
}

func (c *awsClient) GetRoleARNPath(prefix string) (string, error) {
	for _, accountRole := range AccountRoles {
		roleName := fmt.Sprintf("%s-%s-Role", prefix, accountRole.Name)
		role, err := c.iamClient.GetRole(&iam.GetRoleInput{
			RoleName: aws.String(roleName),
		})
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return "", errors.NotFound.Errorf("Roles with the prefix'%s' not found", prefix)
			}
		}
		return GetPathFromARN(aws.StringValue(role.Role.Arn))
	}
	return "", nil
}

func (c *awsClient) IsUpgradedNeededForAccountRolePolicies(prefix string, version string) (bool, error) {
	for _, accountRole := range AccountRoles {
		roleName := fmt.Sprintf("%s-%s-Role", prefix, accountRole.Name)
		role, err := c.iamClient.GetRole(&iam.GetRoleInput{
			RoleName: aws.String(roleName),
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					return false, errors.NotFound.Errorf("Roles with the prefix '%s' not found", prefix)
				}
			}
			return false, err
		}
		isCompatible, err := c.validateRolePolicyUpgradeVersionCompatibility(aws.StringValue(role.Role.RoleName),
			version)

		if err != nil {
			return false, err
		}
		if !isCompatible {
			return true, nil
		}
	}
	return false, nil
}

func (c *awsClient) HasHostedCPPolicies(roleARN string) (bool, error) {
	if roleARN == "" {
		return false, nil
	}

	role, err := c.GetRoleByARN(roleARN)
	if err != nil {
		return false, err
	}

	return common.IamResourceHasTag(FromV1TagToV2Tag(role.Tags), tags.HypershiftPolicies, tags.True), nil
}

func (c *awsClient) HasManagedPolicies(roleARN string) (bool, error) {
	if roleARN == "" {
		return false, nil
	}

	role, err := c.GetRoleByARN(roleARN)
	if err != nil {
		return false, err
	}

	return common.IsManagedRole(FromV1TagToV2Tag(role.Tags)), nil
}

func (c *awsClient) IsUpgradedNeededForAccountRolePoliciesUsingCluster(
	cluster *cmv1.Cluster, version string) (bool, error) {
	for _, role := range AccountRoles {
		roleName, err := GetAccountRoleName(cluster, role.Name)
		if err != nil {
			return false, err
		}
		if roleName == "" {
			continue
		}
		isCompatible, err := c.validateRolePolicyUpgradeVersionCompatibility(aws.StringValue(&roleName), version)

		if err != nil {
			return false, err
		}
		if !isCompatible {
			return true, nil
		}
	}
	return false, nil
}

func (c *awsClient) UpdateTag(roleName string, defaultPolicyVersion string) error {
	return c.AddRoleTag(roleName, common.OpenShiftVersion, defaultPolicyVersion)
}

func (c *awsClient) AddRoleTag(roleName string, key string, value string) error {
	role, err := c.iamClient.GetRole(&iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return err
	}
	_, err = c.iamClient.TagRole(&iam.TagRoleInput{
		RoleName: role.Role.RoleName,
		Tags: []*iam.Tag{
			{
				Key:   aws.String(key),
				Value: aws.String(value),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *awsClient) IsUpgradedNeededForOperatorRolePoliciesUsingCluster(
	cluster *cmv1.Cluster,
	accountID string,
	version string,
	credRequests map[string]*cmv1.STSOperator,
	operatorRolePolicyPrefix string,
) (bool, error) {
	operatorRoles := cluster.AWS().STS().OperatorIAMRoles()
	generalPath, err := GetPathFromARN(operatorRoles[0].RoleARN())
	if err != nil {
		return true, err
	}
	for _, operator := range credRequests {
		operatorRoleARN := FindOperatorRoleBySTSOperator(operatorRoles, operator)
		if operatorRoleARN == "" {
			policyARN := GetOperatorPolicyARN(
				accountID,
				operatorRolePolicyPrefix,
				operator.Namespace(),
				operator.Name(),
				generalPath,
			)
			policyExistsAndUpToDate, err := c.checkPolicyExistsAndUpToDate(policyARN, version)
			return !policyExistsAndUpToDate, err
		}

		roleName, err := GetResourceIdFromARN(operatorRoleARN)
		if err != nil {
			return true, err
		}
		_, err = c.iamClient.GetRole(&iam.GetRoleInput{
			RoleName: aws.String(roleName),
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					return false, errors.NotFound.Errorf("Operator Role '%s' does not exists for the "+
						"cluster '%s'", roleName, cluster.ID())
				}
			}
			return true, err
		}
		isCompatible, err := c.validateRolePolicyUpgradeVersionCompatibility(roleName, version)
		if err != nil {
			return true, err
		}
		if !isCompatible {
			return true, nil
		}
	}
	return false, nil
}

func (c *awsClient) validateRolePolicyUpgradeVersionCompatibility(roleName string,
	version string) (bool, error) {
	attachedPolicies, err := c.GetAttachedPolicy(aws.String(roleName))
	if err != nil {
		return false, err
	}
	for _, attachedPolicy := range attachedPolicies {
		if attachedPolicy.PolicyType == Inline {
			continue
		}
		return c.isRolePolicyUpToDate(attachedPolicy.PolicyArn, version)
	}
	return false, nil
}

func (c *awsClient) IsUpgradedNeededForOperatorRolePoliciesUsingPrefix(prefix string, accountID string,
	version string, credRequests map[string]*cmv1.STSOperator, path string) (bool, error) {
	for _, operator := range credRequests {
		policyARN := GetOperatorPolicyARN(accountID, prefix, operator.Namespace(), operator.Name(), path)
		existsAndUpToDate, err := c.checkPolicyExistsAndUpToDate(policyARN, version)
		if err != nil {
			return false, err
		}
		if !existsAndUpToDate {
			return true, nil
		}
	}
	return false, nil
}

func (c *awsClient) checkPolicyExistsAndUpToDate(policyARN string, policyVersion string) (bool, error) {
	_, err := c.IsPolicyExists(policyARN)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return false, nil
			default:
				return false, err
			}
		}
	}
	isRoleUpToDate, err := c.isRolePolicyUpToDate(policyARN, policyVersion)
	return isRoleUpToDate, err
}

func (c *awsClient) isRolePolicyUpToDate(policyARN string, policyVersion string) (bool, error) {
	isCompatible, err := c.isRolePoliciesCompatibleForUpgrade(policyARN, policyVersion)
	if err != nil {
		return false, errors.Errorf("Failed to validate role policies : %v", err)
	}
	if !isCompatible {
		return false, nil
	}
	return true, nil
}

func (c *awsClient) isRolePoliciesCompatibleForUpgrade(policyARN string, version string) (bool, error) {
	policyTagOutput, err := c.iamClient.ListPolicyTags(&iam.ListPolicyTagsInput{
		PolicyArn: aws.String(policyARN),
	})
	if err != nil {
		return false, err
	}
	return c.hasCompatibleMajorMinorVersionTags(policyTagOutput.Tags, version)
}

func (c *awsClient) GetAccountRoleVersion(roleName string) (string, error) {
	role, err := c.iamClient.GetRole(&iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return "", err
	}
	_, version := GetTagValues(role.Role.Tags)
	return version, nil
}

func (c *awsClient) IsAdminRole(roleName string) (bool, error) {
	role, err := c.iamClient.GetRole(&iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return false, err
	}

	for _, tag := range role.Role.Tags {
		if aws.StringValue(tag.Key) == tags.AdminRole && aws.StringValue(tag.Value) == "true" {
			return true, nil
		}
	}

	return false, nil
}

func (c *awsClient) GetAccountRoleARN(prefix string, roleType string) (string, error) {
	output, err := c.iamClient.GetRole(&iam.GetRoleInput{
		RoleName: aws.String(common.GetRoleName(prefix, roleType)),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
				errorMessage := "Role with the prefix '%s' not found amongst Classic cluster roles"
				if strings.Contains(roleType, HCPSuffixPattern) {
					errorMessage = "Role with the prefix '%s' not found amongst HCP cluster roles"
				}
				return "", errors.NotFound.Errorf(errorMessage, prefix)
			}
		}

		return "", err
	}

	return aws.StringValue(output.Role.Arn), nil
}

func (c *awsClient) ValidateOperatorRolesManagedPolicies(cluster *cmv1.Cluster,
	operatorRoles map[string]*cmv1.STSOperator, policies map[string]*cmv1.AWSSTSPolicy, hostedCPPolicies bool) error {
	for key, operatorRole := range operatorRoles {
		roleName, exist := FindOperatorRoleNameBySTSOperator(cluster, operatorRole)
		if exist {
			err := c.validateManagedPolicy(policies, GetOperatorPolicyKey(key, hostedCPPolicies, false), roleName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *awsClient) ValidateAccountRolesManagedPolicies(prefix string, policies map[string]*cmv1.AWSSTSPolicy) error {
	for roleType, accountRole := range AccountRoles {
		roleName := common.GetRoleName(prefix, accountRole.Name)

		policyKeys := GetAccountRolePolicyKeys(roleType)
		for _, policyKey := range policyKeys {
			err := c.validateManagedPolicy(policies, policyKey, roleName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *awsClient) ValidateHCPAccountRolesManagedPolicies(prefix string,
	policies map[string]*cmv1.AWSSTSPolicy) error {
	for roleType, accountRole := range HCPAccountRoles {
		roleName := common.GetRoleName(prefix, accountRole.Name)

		policyKey := fmt.Sprintf("sts_hcp_%s_permission_policy", roleType)
		err := c.validateManagedPolicy(policies, policyKey, roleName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *awsClient) validateManagedPolicy(policies map[string]*cmv1.AWSSTSPolicy, policyKey string,
	roleName string) error {
	managedPolicyARN, err := GetManagedPolicyARN(policies, policyKey)
	if err != nil {
		return err
	}

	isPolicyAttached, err := c.isManagedPolicyAttached(roleName, managedPolicyARN)
	if err != nil {
		return err
	}
	if !isPolicyAttached {
		return fmt.Errorf("role '%s' is missing the attached managed policy '%s'", roleName, managedPolicyARN)
	}

	return nil
}

func (c *awsClient) isManagedPolicyAttached(roleName string, managedPolicyARN string) (bool, error) {
	policies, err := c.listRoleAttachedPolicies(roleName)
	if err != nil {
		return false, fmt.Errorf("failed to list role '%s' attached policies: %v", roleName, err)
	}

	for _, policy := range policies {
		if aws.StringValue(policy.PolicyArn) == managedPolicyARN {
			return true, nil
		}
	}

	return false, nil
}

func (c *awsClient) listRoleAttachedPolicies(roleName string) ([]*iam.AttachedPolicy, error) {
	attachedPoliciesOutput, err := c.iamClient.ListAttachedRolePolicies(
		&iam.ListAttachedRolePoliciesInput{RoleName: aws.String(roleName)},
	)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
				return []*iam.AttachedPolicy{}, errors.NotFound.Errorf("Role with name '%s' not found", roleName)
			}
		}

		return []*iam.AttachedPolicy{}, err
	}

	return attachedPoliciesOutput.AttachedPolicies, nil
}
