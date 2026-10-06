/*
Copyright 2026 The Kubernetes Authors.

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

// Package ocmrole provides OCM role creation and management functionality.
// This code is extracted from github.com/openshift/rosa/pkg/ocmrole to avoid
// the Go 1.25.8 dependency requirement while CAPA targets Go 1.24.0.
package ocmrole

import (
	"fmt"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	rosaaws "github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/reporter"
	rosacli "github.com/openshift/rosa/pkg/rosa"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa/awshelpers"
)

// RoleProfile defines the type of OCM role to create.
type RoleProfile string

// OCM role profile types.
const (
	ProfileStandard  RoleProfile = "standard"
	ProfileAdmin     RoleProfile = "admin"
	ProfileNoConsole RoleProfile = "no-console"
)

const (
	// TrueString is the string representation of boolean true used in IAM tags.
	TrueString = "true"
)

// GetOrCreateOCMRole gets an existing OCM role or creates it if it doesn't exist (idempotent operation).
//
// Behavior:
//   - If the role exists with the correct profile: returns it immediately (created=false)
//   - If the role exists with wrong profile: returns an error
//   - If the role doesn't exist: creates it with the specified configuration (created=true)
//
// When checking existing roles, this function performs self-healing for policy/tag mismatches
// (e.g., admin policy attached but tag missing).
//
// Profile Mismatch Handling:
// If a role exists but with an incompatible profile (e.g., role is admin but standard was requested),
// this function returns an error describing the mismatch. Upgrading or downgrading role profiles is
// the caller's responsibility.
//
// This function does NOT link the role to OCM organization (caller should use OCMClient.LinkOrgToRole after).
//
// Returns:
//   - roleARN: ARN of the role (whether it existed or was created)
//   - orgID: OCM organization ID
//   - created: true if the role was created by this call, false if it already existed
//   - error: nil on success, or an error describing the issue
func GetOrCreateOCMRole(
	runtime *rosacli.Runtime,
	log logger.Wrapper,
	prefix string,
	profile RoleProfile,
	permissionsBoundary string,
	path string,
) (string, string, bool, error) {
	// Validate inputs
	if runtime == nil {
		return "", "", false, fmt.Errorf("runtime cannot be nil")
	}
	if runtime.AWSClient == nil {
		return "", "", false, fmt.Errorf("AWS client cannot be nil")
	}
	if runtime.Creator == nil {
		return "", "", false, fmt.Errorf("creator cannot be nil")
	}
	if runtime.Reporter == nil {
		return "", "", false, fmt.Errorf("reporter cannot be nil")
	}
	if runtime.OCMClient == nil {
		return "", "", false, fmt.Errorf("OCM client cannot be nil")
	}
	if log == nil {
		return "", "", false, fmt.Errorf("logger cannot be nil")
	}
	if prefix == "" {
		return "", "", false, fmt.Errorf("prefix cannot be empty")
	}
	if profile != ProfileStandard && profile != ProfileAdmin && profile != ProfileNoConsole {
		return "", "", false, fmt.Errorf("profile must be one of: %s, %s, %s", ProfileStandard, ProfileAdmin, ProfileNoConsole)
	}

	awsClient := runtime.AWSClient
	creator := runtime.Creator
	rep := runtime.Reporter

	if path == "" {
		path = "/"
	}

	// Get current OCM organization
	orgID, externalID, err := runtime.OCMClient.GetCurrentOrganization()
	if err != nil {
		return "", "", false, fmt.Errorf("failed to get organization account: %w", err)
	}

	roleName := rosaaws.GetOCMRoleName(prefix, awshelpers.OCMRoleType, externalID)

	policies, err := runtime.OCMClient.GetPolicies("OCMRole")
	if err != nil {
		return "", "", false, fmt.Errorf("failed to get OCM policies: %w", err)
	}

	// Check if role already exists (will self-heal missing policies)
	roleARN, exists, err := checkRoleExists(awsClient, creator, log, roleName, profile, path, policies, runtime.Reporter)
	if err != nil {
		return "", "", false, err
	}
	if exists {
		return roleARN, orgID, false, nil
	}

	// Role doesn't exist - need to create it

	env := rosa.GetOCMClientEnv(runtime.OCMClient)

	// Validate no-console policy exists before creating role
	// This prevents orphaned IAM roles when the policy is missing or malformed
	if profile == ProfileNoConsole {
		filename := fmt.Sprintf("sts_%s_permission_policy", awshelpers.OCMNoConsoleRolePolicyFile)
		policy, ok := policies[filename]
		if !ok || policy.Details() == "" {
			return "", "", false, fmt.Errorf("the no-console OCM role profile is not yet enabled for your Organization")
		}
	}

	roleARN, err = createRole(awsClient, creator, log, prefix, roleName, path, permissionsBoundary,
		orgID, env, profile, policies, rep)
	if err != nil {
		return "", "", false, err
	}

	return roleARN, orgID, true, nil
}

// checkRoleExists checks if a role exists and validates it matches the requested profile.
// If the role exists but is missing policies, it will self-heal by attaching the missing policies.
func checkRoleExists(awsClient rosaaws.Client, creator *rosaaws.Creator, log logger.Wrapper, roleName string, profile RoleProfile, rolePath string, policies map[string]*cmv1.AWSSTSPolicy, rep reporter.Logger) (string, bool, error) {
	exists, roleARN, err := awsClient.CheckRoleExists(roleName)
	if err != nil {
		return "", false, err
	}
	if !exists {
		return "", false, nil
	}

	isExistingRoleAdmin, err := awsClient.IsAdminRole(roleName)
	if err != nil {
		return "", true, err
	}
	isExistingRoleNoConsole, err := IsNoConsoleRole(awsClient, roleName)
	if err != nil {
		return "", true, err
	}

	log.Warn("Role already exists", "roleName", roleName)

	partition := creator.Partition
	accountID := creator.AccountID

	switch profile {
	case ProfileStandard:
		if isExistingRoleAdmin {
			return roleARN, true, fmt.Errorf("the existing role is an admin role. To remove admin capabilities please delete the admin policy and the '%s' tag", awshelpers.TagAdminRole)
		}
		if isExistingRoleNoConsole {
			return roleARN, true, fmt.Errorf("the existing role is a no-console role. To use standard permissions please delete the role and recreate it")
		}

		// Self-heal: verify standard policy is attached, attach if missing
		attachedPolicies, err := awsClient.ListAttachedRolePolicies(roleName)
		if err != nil {
			return "", true, err
		}

		expectedPolicyARN := rosaaws.GetPolicyArnWithSuffix(partition, accountID, roleName, rolePath)
		if !slices.Contains(attachedPolicies, expectedPolicyARN) {
			log.Info("Standard policy missing - attaching", "roleName", roleName, "policyARN", expectedPolicyARN)
			filename := fmt.Sprintf("sts_%s_permission_policy", awshelpers.OCMRolePolicyFile)
			policyDetail := rosaaws.GetPolicyDetails(policies, filename)
			err = createPermissionPolicy(awsClient, log, expectedPolicyARN, map[string]string{}, roleName, rolePath, policyDetail, rep)
			if err != nil {
				return "", true, fmt.Errorf("failed to attach standard policy: %w", err)
			}
		}

		return roleARN, true, nil

	case ProfileAdmin:
		if isExistingRoleNoConsole {
			return roleARN, true, fmt.Errorf("the existing role is a no-console role. To use admin permissions please delete the role and recreate it")
		}

		// Self-heal: verify both standard and admin policies are attached
		attachedPolicies, err := awsClient.ListAttachedRolePolicies(roleName)
		if err != nil {
			return "", true, err
		}

		// Check and attach standard policy if missing
		expectedStandardPolicyARN := rosaaws.GetPolicyArnWithSuffix(partition, accountID, roleName, rolePath)
		if !slices.Contains(attachedPolicies, expectedStandardPolicyARN) {
			log.Info("Standard policy missing for admin role - attaching", "roleName", roleName, "policyARN", expectedStandardPolicyARN)
			filename := fmt.Sprintf("sts_%s_permission_policy", awshelpers.OCMRolePolicyFile)
			policyDetail := rosaaws.GetPolicyDetails(policies, filename)
			err = createPermissionPolicy(awsClient, log, expectedStandardPolicyARN, map[string]string{}, roleName, rolePath, policyDetail, rep)
			if err != nil {
				return "", true, fmt.Errorf("failed to attach standard policy: %w", err)
			}
		}

		// Check and attach admin policy if missing
		expectedAdminPolicyARN := rosaaws.GetAdminPolicyARN(partition, accountID, roleName, rolePath)
		if !slices.Contains(attachedPolicies, expectedAdminPolicyARN) {
			log.Info("Admin policy missing - attaching", "roleName", roleName, "policyARN", expectedAdminPolicyARN)
			filename := fmt.Sprintf("sts_%s_permission_policy", awshelpers.OCMAdminRolePolicyFile)
			iamTags := map[string]string{awshelpers.TagAdminRole: TrueString}
			policyDetail := rosaaws.GetPolicyDetails(policies, filename)
			err = createPermissionPolicy(awsClient, log, expectedAdminPolicyARN, iamTags, roleName, rolePath, policyDetail, rep)
			if err != nil {
				return "", true, fmt.Errorf("failed to attach admin policy: %w", err)
			}
		}

		// Self-heal: add admin tag if missing
		if !isExistingRoleAdmin {
			log.Info("Admin tag missing - adding", "roleName", roleName)
			err = awsClient.AddRoleTag(roleName, awshelpers.TagAdminRole, TrueString)
			if err != nil {
				return "", true, fmt.Errorf("failed to add admin role tag: %w", err)
			}
		}

		return roleARN, true, nil

	case ProfileNoConsole:
		if isExistingRoleAdmin {
			return roleARN, true, fmt.Errorf("the existing role is an admin role. To use no-console permissions please delete the role and recreate it")
		}

		// Self-heal: verify no-console policy is attached, attach if missing
		attachedPolicies, err := awsClient.ListAttachedRolePolicies(roleName)
		if err != nil {
			return "", true, err
		}

		expectedNoConsolePolicyARN := awshelpers.GetNoConsolePolicyARN(partition, accountID, roleName, rolePath)
		if !slices.Contains(attachedPolicies, expectedNoConsolePolicyARN) {
			log.Info("No-console policy missing - attaching", "roleName", roleName, "policyARN", expectedNoConsolePolicyARN)
			filename := fmt.Sprintf("sts_%s_permission_policy", awshelpers.OCMNoConsoleRolePolicyFile)
			iamTags := map[string]string{awshelpers.TagNoConsoleRole: TrueString}
			policyDetail := rosaaws.GetPolicyDetails(policies, filename)
			err = createPermissionPolicy(awsClient, log, expectedNoConsolePolicyARN, iamTags, roleName, rolePath, policyDetail, rep)
			if err != nil {
				return "", true, fmt.Errorf("failed to attach no-console policy: %w", err)
			}
		}

		// Self-heal: add no-console tag if missing
		if !isExistingRoleNoConsole {
			log.Info("No-console tag missing - adding", "roleName", roleName)
			err = awsClient.AddRoleTag(roleName, awshelpers.TagNoConsoleRole, TrueString)
			if err != nil {
				return "", true, fmt.Errorf("failed to add no-console role tag: %w", err)
			}
		}

		return roleARN, true, nil

	default:
		// Should never reach here if validation is done at boundaries
		return "", false, fmt.Errorf("invalid profile: %s (must be one of: %s, %s, %s)",
			profile, ProfileStandard, ProfileAdmin, ProfileNoConsole)
	}
}

// IsNoConsoleRole checks if a role has the no-console tag.
// This is a helper for older rosa versions that don't have IsNoConsoleRole on the Client interface.
func IsNoConsoleRole(awsClient rosaaws.Client, roleName string) (bool, error) {
	role, err := awsClient.GetRoleByName(roleName)
	if err != nil {
		return false, err
	}

	for _, tag := range role.Tags {
		if aws.ToString(tag.Key) == awshelpers.TagNoConsoleRole && aws.ToString(tag.Value) == TrueString {
			return true, nil
		}
	}

	return false, nil
}

// createRole creates a new OCM role with the specified configuration.
func createRole(awsClient rosaaws.Client, creator *rosaaws.Creator, log logger.Wrapper, prefix string, roleName string, rolePath string,
	permissionsBoundary string, orgID string, env string, profile RoleProfile,
	policies map[string]*cmv1.AWSSTSPolicy, rep reporter.Logger,
) (string, error) {
	partition := creator.Partition
	accountID := creator.AccountID

	var policyARN string

	if profile != ProfileNoConsole {
		policyARN = rosaaws.GetPolicyArnWithSuffix(partition, accountID, roleName, rolePath)
	}

	// Build trust policy
	filename := fmt.Sprintf("sts_%s_trust_policy", awshelpers.OCMRolePolicyFile)
	policyDetail := rosaaws.GetPolicyDetails(policies, filename)
	policy := rosaaws.InterpolatePolicyDocument(partition, policyDetail, map[string]string{
		"partition":           partition,
		"aws_account_id":      rosaaws.GetJumpAccount(env),
		"ocm_organization_id": orgID,
	})

	// Build IAM tags
	iamTags := map[string]string{
		awshelpers.TagRolePrefix:    prefix,
		awshelpers.TagRoleType:      awshelpers.OCMRoleType,
		awshelpers.TagEnvironment:   env,
		awshelpers.TagRedHatManaged: TrueString,
	}

	// Verify profile is valid before creating any AWS resources
	if profile != ProfileStandard && profile != ProfileAdmin && profile != ProfileNoConsole {
		return "", fmt.Errorf("invalid profile: %s (must be one of: %s, %s, %s)",
			profile, ProfileStandard, ProfileAdmin, ProfileNoConsole)
	}

	log.Debug("Creating role", "roleName", roleName)

	roleARN, err := awsClient.EnsureRole(rep, roleName, policy, permissionsBoundary, "", iamTags, rolePath, false)
	if err != nil {
		return "", err
	}
	log.Info("Created role", "roleName", roleName, "roleARN", roleARN)

	switch profile {
	case ProfileStandard:
		filename := fmt.Sprintf("sts_%s_permission_policy", awshelpers.OCMRolePolicyFile)
		policyDetail := rosaaws.GetPolicyDetails(policies, filename)
		err := createPermissionPolicy(awsClient, log, policyARN, iamTags, roleName, rolePath, policyDetail, rep)
		if err != nil {
			return "", err
		}

	case ProfileAdmin:
		// standard policy first
		filename := fmt.Sprintf("sts_%s_permission_policy", awshelpers.OCMRolePolicyFile)
		policyDetail := rosaaws.GetPolicyDetails(policies, filename)
		err := createPermissionPolicy(awsClient, log, policyARN, iamTags, roleName, rolePath, policyDetail, rep)
		if err != nil {
			return "", err
		}

		// create and attach the admin policy to the role
		filename = fmt.Sprintf("sts_%s_permission_policy", awshelpers.OCMAdminRolePolicyFile)
		policyARN = rosaaws.GetAdminPolicyARN(partition, accountID, roleName, rolePath)
		iamTags[awshelpers.TagAdminRole] = TrueString
		policyDetail = rosaaws.GetPolicyDetails(policies, filename)
		err = createPermissionPolicy(awsClient, log, policyARN, iamTags, roleName, rolePath, policyDetail, rep)
		if err != nil {
			return "", err
		}

		// tag role with admin tag
		err = awsClient.AddRoleTag(roleName, awshelpers.TagAdminRole, TrueString)
		if err != nil {
			return "", err
		}

	case ProfileNoConsole:
		filename := fmt.Sprintf("sts_%s_permission_policy", awshelpers.OCMNoConsoleRolePolicyFile)

		// create and attach the no-console policy to the role
		policyARN = awshelpers.GetNoConsolePolicyARN(partition, accountID, roleName, rolePath)
		iamTags[awshelpers.TagNoConsoleRole] = TrueString
		policyDetail := rosaaws.GetPolicyDetails(policies, filename)
		err := createPermissionPolicy(awsClient, log, policyARN, iamTags, roleName, rolePath, policyDetail, rep)
		if err != nil {
			return "", err
		}

		// tag role with no-console tag
		err = awsClient.AddRoleTag(roleName, awshelpers.TagNoConsoleRole, TrueString)
		if err != nil {
			return "", err
		}

	default:
		return "", fmt.Errorf("invalid profile: %s (must be one of: %s, %s, %s)",
			profile, ProfileStandard, ProfileAdmin, ProfileNoConsole)
	}

	return roleARN, nil
}

// createPermissionPolicy creates and attaches a customer-managed permission policy to an IAM role.
func createPermissionPolicy(awsClient rosaaws.Client, log logger.Wrapper, policyARN string,
	iamTags map[string]string, roleName string, rolePath string, policyDetail string, rep reporter.Logger,
) error {
	log.Debug("Creating permission policy", "policyARN", policyARN)
	var err error
	policyARN, err = awsClient.EnsurePolicy(policyARN, policyDetail, "", iamTags, rolePath)
	if err != nil {
		return err
	}

	log.Debug("Attaching permission policy to role", "roleName", roleName, "policyARN", policyARN)
	err = awsClient.AttachRolePolicy(rep, roleName, policyARN)
	if err != nil {
		return err
	}

	return nil
}
