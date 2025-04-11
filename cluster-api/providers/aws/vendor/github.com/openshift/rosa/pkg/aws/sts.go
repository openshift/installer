/*
Copyright (c) 2022 Red Hat, Inc.

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
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	awserr "github.com/openshift-online/ocm-common/pkg/aws/errors"
	common "github.com/openshift-online/ocm-common/pkg/aws/validations"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	awscbRoles "github.com/openshift/rosa/pkg/aws/commandbuilder/helper/roles"
	"github.com/openshift/rosa/pkg/aws/tags"
	rprtr "github.com/openshift/rosa/pkg/reporter"
)

func (c *awsClient) DeleteUserRole(roleName string) error {
	err := c.detachAttachedRolePolicies(aws.String(roleName))
	if err != nil {
		return err
	}

	err = c.deletePermissionsBoundary(roleName)
	if err != nil {
		return err
	}

	return c.DeleteRole(roleName)
}

func (c *awsClient) DeleteOCMRole(roleName string, managedPolicies bool) error {
	err := c.deleteOCMRolePolicies(roleName, managedPolicies)
	if err != nil {
		return err
	}

	err = c.deletePermissionsBoundary(roleName)
	if err != nil {
		return err
	}

	return c.DeleteRole(roleName)
}

func (c *awsClient) ValidateRoleARNAccountIDMatchCallerAccountID(roleARN string) error {
	creator, err := c.GetCreator()
	if err != nil {
		return fmt.Errorf("failed to get AWS creator: %v", err)
	}

	parsedARN, err := arn.Parse(roleARN)
	if err != nil {
		return err
	}

	if creator.AccountID != parsedARN.AccountID {
		return fmt.Errorf("role ARN '%s' doesn't match the user's account ID '%s'", roleARN, creator.AccountID)
	}

	return nil
}

func (c *awsClient) HasPermissionsBoundary(roleName string) (bool, error) {
	output, err := c.iamClient.GetRole(context.Background(), &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return false, err
	}

	return output.Role.PermissionsBoundary != nil, nil
}

func (c *awsClient) deletePermissionsBoundary(roleName string) error {
	output, err := c.iamClient.GetRole(context.Background(), &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return err
	}

	if output.Role.PermissionsBoundary != nil {
		_, err := c.iamClient.DeleteRolePermissionsBoundary(context.Background(), &iam.DeleteRolePermissionsBoundaryInput{
			RoleName: aws.String(roleName),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *awsClient) deleteOCMRolePolicies(roleName string, managedPolicies bool) error {
	policiesOutput, err := c.iamClient.ListAttachedRolePolicies(context.Background(), &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return err
	}

	for _, policy := range policiesOutput.AttachedPolicies {
		_, err := c.iamClient.DetachRolePolicy(context.Background(), &iam.DetachRolePolicyInput{
			PolicyArn: policy.PolicyArn,
			RoleName:  aws.String(roleName),
		})
		if err != nil {
			return err
		}

		if !managedPolicies {
			_, err = c.iamClient.DeletePolicy(context.Background(), &iam.DeletePolicyInput{PolicyArn: policy.PolicyArn})
			if err != nil {
				if awserr.IsDeleteConfictException(err) {
					continue
				}
				return err
			}
		}
	}

	return nil
}

func SortRolesByLinkedRole(roles []Role) {
	sort.SliceStable(roles, func(i, j int) bool {
		return roles[i].Linked == "Yes" && roles[j].Linked == "No"
	})
}

func UpgradeOperatorPolicies(reporter *rprtr.Object, awsClient Client, partition string, accountID string,
	prefix string, policies map[string]string, defaultPolicyVersion string,
	credRequests map[string]*cmv1.STSOperator, path string) error {
	for credrequest, operator := range credRequests {
		policyARN := GetOperatorPolicyARN(partition, accountID, prefix, operator.Namespace(), operator.Name(), path)
		filename := fmt.Sprintf("openshift_%s_policy", credrequest)
		policy := policies[filename]
		policyARN, err := awsClient.EnsurePolicy(policyARN, policy,
			defaultPolicyVersion, map[string]string{
				common.OpenShiftVersion: defaultPolicyVersion,
				tags.RolePrefix:         prefix,
				tags.RedHatManaged:      "true",
				tags.OperatorNamespace:  operator.Namespace(),
				tags.OperatorName:       operator.Name(),
			}, "")
		if err != nil {
			return err
		}
		reporter.Infof("Upgraded policy with ARN '%s' to version '%s'", policyARN, defaultPolicyVersion)
	}
	return nil
}

func BuildOperatorRoleCommands(prefix string, partition string, accountID string, awsClient Client,
	defaultPolicyVersion string, credRequests map[string]*cmv1.STSOperator, policyPath string,
	cluster *cmv1.Cluster) []string {
	commands := []string{}
	for credrequest, operator := range credRequests {
		policyARN := GetOperatorPolicyARN(
			partition,
			accountID,
			prefix,
			operator.Namespace(),
			operator.Name(),
			policyPath,
		)
		policyName := GetOperatorPolicyName(
			prefix,
			operator.Namespace(),
			operator.Name(),
		)
		_, err := awsClient.IsPolicyExists(policyARN)
		policyExists := err == nil
		isSharedVpc := cluster.AWS().PrivateHostedZoneRoleARN() != ""
		fileName := GetOperatorPolicyKey(credrequest, cluster.Hypershift().Enabled(), isSharedVpc)
		fileName = GetFormattedFileName(fileName)
		upgradePoliciesCommands := awscbRoles.ManualCommandsForUpgradeOperatorRolePolicy(
			awscbRoles.ManualCommandsForUpgradeOperatorRolePolicyInput{
				PolicyExists:             policyExists,
				OperatorRolePolicyPrefix: prefix,
				Operator:                 operator,
				CredRequest:              credrequest,
				OperatorPolicyPath:       policyPath,
				PolicyARN:                policyARN,
				DefaultPolicyVersion:     defaultPolicyVersion,
				PolicyName:               policyName,
				FileName:                 fileName,
			},
		)
		commands = append(commands, upgradePoliciesCommands...)
	}
	return commands
}

type OidcProviderOutput struct {
	Arn       string
	ClusterId string
}

func (c *awsClient) ListOidcProviders(targetClusterId string, config *cmv1.OidcConfig) ([]OidcProviderOutput, error) {
	providers := []OidcProviderOutput{}
	output, err := c.iamClient.ListOpenIDConnectProviders(context.Background(), &iam.ListOpenIDConnectProvidersInput{})
	if err != nil {
		return providers, err
	}
	for _, provider := range output.OpenIDConnectProviderList {
		if err != nil {
			return providers, err
		}
		isTruncated := true
		var marker *string
		for isTruncated {
			resp, err := c.iamClient.ListOpenIDConnectProviderTags(context.Background(), &iam.ListOpenIDConnectProviderTagsInput{
				OpenIDConnectProviderArn: provider.Arn,
				Marker:                   marker,
			})
			if err != nil {
				return providers, err
			}
			isTruncated = resp.IsTruncated
			marker = resp.Marker
			skip := true
			clusterId := ""
			for _, tag := range resp.Tags {
				switch *tag.Key {
				case tags.ClusterID:
					clusterId = *tag.Value
				case tags.RedHatManaged:
					skip = false
				}
			}
			if targetClusterId != "" {
				if targetClusterId != clusterId {
					skip = true
				} else {
					providers = append(providers, OidcProviderOutput{
						Arn:       *provider.Arn,
						ClusterId: clusterId,
					})
					return providers, nil
				}
			}
			if config != nil {
				resourceId, err := GetResourceIdFromOidcProviderARN(*provider.Arn)
				if err != nil {
					return nil, fmt.Errorf("unable to get resource ID from OIDC Provider's ARN. Error: '%v'", err)
				}
				if config == nil || !strings.Contains(config.IssuerUrl(), resourceId) {
					skip = true
				}
			}
			if skip {
				continue
			}
			providers = append(providers, OidcProviderOutput{
				Arn:       *provider.Arn,
				ClusterId: clusterId,
			})
		}
	}
	return providers, nil
}
