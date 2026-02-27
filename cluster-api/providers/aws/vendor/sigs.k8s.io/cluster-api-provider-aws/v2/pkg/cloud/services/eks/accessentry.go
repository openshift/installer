/*
Copyright 2025 The Kubernetes Authors.

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

package eks

import (
	"context"
	"slices"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

func (s *Service) reconcileAccessEntries(ctx context.Context) error {
	if len(s.scope.ControlPlane.Spec.AccessEntries) == 0 {
		s.scope.Info("no access entries defined, skipping reconcile")
		return nil
	}

	if s.scope.ControlPlane.Spec.AccessConfig == nil ||
		s.scope.ControlPlane.Spec.AccessConfig.AuthenticationMode == "" ||
		(s.scope.ControlPlane.Spec.AccessConfig.AuthenticationMode != ekscontrolplanev1.EKSAuthenticationModeAPI &&
			s.scope.ControlPlane.Spec.AccessConfig.AuthenticationMode != ekscontrolplanev1.EKSAuthenticationModeAPIAndConfigMap) {
		s.scope.Info("access mode is not api or api_and_config_map, skipping reconcile")
		return nil
	}

	managedAccessEntries, err := s.getManagedAccessEntries(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list existing access entries")
	}

	for _, accessEntry := range s.scope.ControlPlane.Spec.AccessEntries {
		if _, exists := managedAccessEntries[accessEntry.PrincipalARN]; exists {
			if err := s.updateAccessEntry(ctx, accessEntry); err != nil {
				return errors.Wrapf(err, "failed to update access entry for principal %s", accessEntry.PrincipalARN)
			}
			delete(managedAccessEntries, accessEntry.PrincipalARN)
		} else {
			if err := s.createAccessEntry(ctx, accessEntry); err != nil {
				return errors.Wrapf(err, "failed to create access entry for principal %s", accessEntry.PrincipalARN)
			}
		}
	}

	for principalArn := range managedAccessEntries {
		if err := s.deleteAccessEntry(ctx, principalArn); err != nil {
			return errors.Wrapf(err, "failed to delete access entry for principal %s", principalArn)
		}
	}

	record.Event(s.scope.ControlPlane, "SuccessfulReconcileAccessEntries", "Reconciled access entries")
	return nil
}

func (s *Service) getManagedAccessEntries(ctx context.Context) (map[string]bool, error) {
	existingAccessEntries := make(map[string]bool)
	var nextToken *string

	clusterName := s.scope.KubernetesClusterName()
	managedTag := infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())

	for {
		input := &eks.ListAccessEntriesInput{
			ClusterName: &clusterName,
			NextToken:   nextToken,
		}

		output, err := s.EKSClient.ListAccessEntries(ctx, input)
		if err != nil {
			return nil, errors.Wrap(err, "failed to list access entries")
		}

		for _, principalArn := range output.AccessEntries {
			describeOutput, err := s.EKSClient.DescribeAccessEntry(ctx, &eks.DescribeAccessEntryInput{
				ClusterName:  &clusterName,
				PrincipalArn: &principalArn,
			})
			if err != nil {
				s.scope.Error(err, "failed to describe access entry", "principalARN", principalArn)
				continue
			}

			if describeOutput.AccessEntry.Tags != nil {
				if _, managed := describeOutput.AccessEntry.Tags[managedTag]; managed {
					existingAccessEntries[principalArn] = true
				}
			}
		}

		if output.NextToken == nil {
			break
		}

		nextToken = output.NextToken
	}

	return existingAccessEntries, nil
}

func (s *Service) createAccessEntry(ctx context.Context, accessEntry ekscontrolplanev1.AccessEntry) error {
	clusterName := s.scope.KubernetesClusterName()

	additionalTags := s.scope.AdditionalTags()
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)
	tags := make(map[string]string)
	for k, v := range additionalTags {
		tags[k] = v
	}

	createInput := &eks.CreateAccessEntryInput{
		ClusterName:  &clusterName,
		PrincipalArn: &accessEntry.PrincipalARN,
		Tags:         tags,
	}

	if len(accessEntry.KubernetesGroups) > 0 {
		createInput.KubernetesGroups = accessEntry.KubernetesGroups
	}

	if accessEntry.Type != "" {
		createInput.Type = accessEntry.Type.APIValue()
	}

	if accessEntry.Username != "" {
		createInput.Username = &accessEntry.Username
	}

	if _, err := s.EKSClient.CreateAccessEntry(ctx, createInput); err != nil {
		return errors.Wrapf(err, "failed to create access entry for principal %s", accessEntry.PrincipalARN)
	}

	if err := s.reconcileAccessPolicies(ctx, accessEntry); err != nil {
		return errors.Wrapf(err, "failed to reconcile access policies for principal %s", accessEntry.PrincipalARN)
	}

	return nil
}

func (s *Service) updateAccessEntry(ctx context.Context, accessEntry ekscontrolplanev1.AccessEntry) error {
	clusterName := s.scope.KubernetesClusterName()
	describeInput := &eks.DescribeAccessEntryInput{
		ClusterName:  &clusterName,
		PrincipalArn: &accessEntry.PrincipalARN,
	}

	describeOutput, err := s.EKSClient.DescribeAccessEntry(ctx, describeInput)
	if err != nil {
		return errors.Wrapf(err, "failed to describe access entry for principal %s", accessEntry.PrincipalARN)
	}

	// EKS requires recreate when changing type or removing username
	existingUsername := ""
	if describeOutput.AccessEntry.Username != nil {
		existingUsername = *describeOutput.AccessEntry.Username
	}

	if *accessEntry.Type.APIValue() != *describeOutput.AccessEntry.Type || accessEntry.Username != existingUsername {
		if err = s.deleteAccessEntry(ctx, accessEntry.PrincipalARN); err != nil {
			return errors.Wrapf(err, "failed to delete access entry for principal %s during recreation", accessEntry.PrincipalARN)
		}

		if err = s.createAccessEntry(ctx, accessEntry); err != nil {
			return errors.Wrapf(err, "failed to recreate access entry for principal %s", accessEntry.PrincipalARN)
		}
		return nil
	}

	slices.Sort(accessEntry.KubernetesGroups)
	slices.Sort(describeOutput.AccessEntry.KubernetesGroups)

	updateInput := &eks.UpdateAccessEntryInput{
		ClusterName:  &clusterName,
		PrincipalArn: &accessEntry.PrincipalARN,
	}

	if !slices.Equal(accessEntry.KubernetesGroups, describeOutput.AccessEntry.KubernetesGroups) {
		updateInput.KubernetesGroups = accessEntry.KubernetesGroups
		if _, err := s.EKSClient.UpdateAccessEntry(ctx, updateInput); err != nil {
			return errors.Wrapf(err, "failed to update access entry for principal %s", accessEntry.PrincipalARN)
		}
	}

	if err := s.reconcileAccessPolicies(ctx, accessEntry); err != nil {
		return errors.Wrapf(err, "failed to reconcile access policies for principal %s", accessEntry.PrincipalARN)
	}

	return nil
}

func (s *Service) deleteAccessEntry(ctx context.Context, principalArn string) error {
	clusterName := s.scope.KubernetesClusterName()

	if _, err := s.EKSClient.DeleteAccessEntry(ctx, &eks.DeleteAccessEntryInput{
		ClusterName:  &clusterName,
		PrincipalArn: &principalArn,
	}); err != nil {
		return errors.Wrapf(err, "failed to delete access entry for principal %s", principalArn)
	}

	return nil
}

func (s *Service) reconcileAccessPolicies(ctx context.Context, accessEntry ekscontrolplanev1.AccessEntry) error {
	if accessEntry.Type == ekscontrolplanev1.AccessEntryTypeEC2Linux || accessEntry.Type == ekscontrolplanev1.AccessEntryTypeEC2Windows {
		s.scope.Info("Skipping access policy reconciliation for EC2 access type", "principalARN", accessEntry.PrincipalARN)
		return nil
	}

	existingPolicies, err := s.getExistingAccessPolicies(ctx, accessEntry.PrincipalARN)
	if err != nil {
		return errors.Wrapf(err, "failed to get existing access policies for principal %s", accessEntry.PrincipalARN)
	}

	clusterName := s.scope.KubernetesClusterName()

	for _, policy := range accessEntry.AccessPolicies {
		input := &eks.AssociateAccessPolicyInput{
			ClusterName:  &clusterName,
			PrincipalArn: &accessEntry.PrincipalARN,
			PolicyArn:    &policy.PolicyARN,
			AccessScope: &ekstypes.AccessScope{
				Type: ekstypes.AccessScopeType(policy.AccessScope.Type),
			},
		}

		if policy.AccessScope.Type == "namespace" && len(policy.AccessScope.Namespaces) > 0 {
			input.AccessScope.Namespaces = policy.AccessScope.Namespaces
		}

		if _, err := s.EKSClient.AssociateAccessPolicy(ctx, input); err != nil {
			return errors.Wrapf(err, "failed to associate access policy %s", policy.PolicyARN)
		}

		delete(existingPolicies, policy.PolicyARN)
	}

	for policyARN := range existingPolicies {
		if _, err := s.EKSClient.DisassociateAccessPolicy(ctx, &eks.DisassociateAccessPolicyInput{
			ClusterName:  &clusterName,
			PrincipalArn: &accessEntry.PrincipalARN,
			PolicyArn:    &policyARN,
		}); err != nil {
			return errors.Wrapf(err, "failed to disassociate access policy %s", policyARN)
		}
	}

	return nil
}

func (s *Service) getExistingAccessPolicies(ctx context.Context, principalARN string) (map[string]ekstypes.AssociatedAccessPolicy, error) {
	existingPolicies := map[string]ekstypes.AssociatedAccessPolicy{}
	var nextToken *string
	clusterName := s.scope.KubernetesClusterName()

	for {
		input := &eks.ListAssociatedAccessPoliciesInput{
			ClusterName:  &clusterName,
			PrincipalArn: &principalARN,
			NextToken:    nextToken,
		}

		output, err := s.EKSClient.ListAssociatedAccessPolicies(ctx, input)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list associated access policies for principal %s", principalARN)
		}

		for _, policy := range output.AssociatedAccessPolicies {
			existingPolicies[*policy.PolicyArn] = policy
		}

		if output.NextToken == nil {
			break
		}

		nextToken = output.NextToken
	}

	return existingPolicies, nil
}
