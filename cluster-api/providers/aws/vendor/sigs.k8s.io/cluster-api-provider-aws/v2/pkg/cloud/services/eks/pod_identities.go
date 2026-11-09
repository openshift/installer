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

package eks

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	awserrors "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

func (s *Service) reconcilePodIdentityAssociations(ctx context.Context) error {
	s.scope.Info("Reconciling EKS Pod Identities")

	// If EKSPodIdentityAssociationConfiguredCondition is present, we've at least attempted
	// to add associations before and still need to make this API call in case there are
	// orphans to clean up.
	if len(s.scope.ControlPlane.Spec.PodIdentityAssociations) == 0 &&
		!v1beta1conditions.Has(s.scope.ControlPlane, ekscontrolplanev1.EKSPodIdentityAssociationConfiguredCondition) {
		s.scope.Debug("no eks pod identity associations desired or previously configured, skipping reconciliation")
		return nil
	}

	eksClusterName := s.scope.KubernetesClusterName()

	// Get existing eks pod identities on the cluster
	managedAssociations, err := s.getManagedPodIdentityAssociations(ctx, eksClusterName)
	if err != nil {
		s.scope.Error(err, "failed listing eks pod identity associations")
		return fmt.Errorf("listing eks pod identity associations: %w", err)
	}

	if len(managedAssociations) == 0 && len(s.scope.ControlPlane.Spec.PodIdentityAssociations) == 0 {
		s.scope.Debug("no eks pod identities found, skipping reconciliation")
		return nil
	}

	for _, assoc := range s.scope.ControlPlane.Spec.PodIdentityAssociations {
		namespacedName := fmt.Sprintf("%s/%s", assoc.ServiceAccountNamespace, assoc.ServiceAccountName)
		if assocID, exists := managedAssociations[namespacedName]; exists {
			if err := s.updatePodIdentityAssociation(ctx, assocID, assoc); err != nil {
				return errors.Wrapf(err, "failed to update pod identity association for %s", namespacedName)
			}
			delete(managedAssociations, namespacedName)
		} else {
			if err := s.createPodIdentityAssociation(ctx, assoc); err != nil {
				return errors.Wrapf(err, "failed to create pod identity association for %s", namespacedName)
			}
		}
	}

	for namespacedName, assocID := range managedAssociations {
		if err := s.deletePodIdentityAssociation(ctx, assocID); err != nil {
			return errors.Wrapf(err, "failed to delete pod identity association for %s", namespacedName)
		}
	}

	record.Eventf(s.scope.ControlPlane, "SuccessfulReconcileEKSClusterPodIdentityAssociations", "Reconciled Pod Identity associations for EKS Cluster %s", s.scope.KubernetesClusterName())
	s.scope.Info("Reconcile EKS pod identity associations completed successfully")

	return nil
}

func (s *Service) getManagedPodIdentityAssociations(ctx context.Context, eksClusterName string) (map[string]string, error) {
	associations := make(map[string]string)
	var nextToken *string

	managedTag := v1beta2.ClusterAWSCloudProviderTagKey(eksClusterName)

	for {
		input := &eks.ListPodIdentityAssociationsInput{
			ClusterName: &eksClusterName,
			NextToken:   nextToken,
		}

		output, err := s.EKSClient.ListPodIdentityAssociations(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("listing eks pod identity associations: %w", err)
		}

		for _, assoc := range output.Associations {
			describeOutput, err := s.EKSClient.DescribePodIdentityAssociation(ctx, &eks.DescribePodIdentityAssociationInput{
				AssociationId: assoc.AssociationId,
				ClusterName:   assoc.ClusterName,
			})
			if err != nil {
				smithyErr := awserrors.ParseSmithyError(err)
				notFoundErr := &ekstypes.ResourceNotFoundException{}
				if smithyErr.ErrorCode() == notFoundErr.ErrorCode() {
					s.scope.Debug("pod identity association no longer exists, skipping", "associationID", *assoc.AssociationId)
					continue
				}
				return nil, fmt.Errorf("describing eks pod identity association %s: %w", *assoc.AssociationId, err)
			}

			if describeOutput.Association.Tags != nil {
				if _, ok := describeOutput.Association.Tags[managedTag]; ok {
					namespacedName := fmt.Sprintf("%s/%s", *describeOutput.Association.Namespace, *describeOutput.Association.ServiceAccount)
					associations[namespacedName] = *assoc.AssociationId
				}
			}
		}

		if output.NextToken == nil {
			break
		}

		nextToken = output.NextToken
	}

	return associations, nil
}

func (s *Service) updatePodIdentityAssociation(ctx context.Context, assocID string, assoc ekscontrolplanev1.PodIdentityAssociation) error {
	clusterName := s.scope.KubernetesClusterName()
	describeInput := &eks.DescribePodIdentityAssociationInput{
		AssociationId: &assocID,
		ClusterName:   &clusterName,
	}

	describeOutput, err := s.EKSClient.DescribePodIdentityAssociation(ctx, describeInput)
	if err != nil {
		return fmt.Errorf("describing eks pod identity association: %w", err)
	}

	// EKS requires recreating the pod identity association if NS or service account name changes
	if *describeOutput.Association.Namespace != assoc.ServiceAccountNamespace || *describeOutput.Association.ServiceAccount != assoc.ServiceAccountName {
		if err := s.deletePodIdentityAssociation(ctx, assocID); err != nil {
			return errors.Wrapf(err, "failed to delete pod identity association %s for update", assocID)
		}

		if err := s.createPodIdentityAssociation(ctx, assoc); err != nil {
			return errors.Wrapf(err, "failed to recreate pod identity association for service account %s/%s", assoc.ServiceAccountNamespace, assoc.ServiceAccountName)
		}

		return nil
	}

	existingRoleArn := ""
	existingTargetRoleArn := ""

	if describeOutput.Association.RoleArn != nil {
		existingRoleArn = *describeOutput.Association.RoleArn
	}

	if describeOutput.Association.TargetRoleArn != nil {
		existingTargetRoleArn = *describeOutput.Association.TargetRoleArn
	}

	if existingRoleArn != assoc.RoleARN || existingTargetRoleArn != assoc.TargetRoleARN {
		updateInput := &eks.UpdatePodIdentityAssociationInput{
			AssociationId: &assocID,
			ClusterName:   &clusterName,
			RoleArn:       &assoc.RoleARN,
		}

		if assoc.TargetRoleARN != "" || existingTargetRoleArn != "" {
			updateInput.TargetRoleArn = &assoc.TargetRoleARN
		}

		if _, err := s.EKSClient.UpdatePodIdentityAssociation(ctx, updateInput); err != nil {
			return errors.Wrapf(err, "failed to update pod identity association %s", assocID)
		}
	}

	return nil
}

func (s *Service) createPodIdentityAssociation(ctx context.Context, assoc ekscontrolplanev1.PodIdentityAssociation) error {
	clusterName := s.scope.KubernetesClusterName()

	additionalTags := s.scope.AdditionalTags()
	additionalTags[v1beta2.ClusterAWSCloudProviderTagKey(s.scope.KubernetesClusterName())] = string(v1beta2.ResourceLifecycleOwned)

	input := &eks.CreatePodIdentityAssociationInput{
		ClusterName:    &clusterName,
		Namespace:      &assoc.ServiceAccountNamespace,
		RoleArn:        &assoc.RoleARN,
		ServiceAccount: &assoc.ServiceAccountName,
		Tags:           additionalTags,
	}

	if assoc.TargetRoleARN != "" {
		input.TargetRoleArn = &assoc.TargetRoleARN
	}

	if _, err := s.EKSClient.CreatePodIdentityAssociation(ctx, input); err != nil {
		return errors.Wrapf(err, "creating pod identity association for service account %s/%s", assoc.ServiceAccountNamespace, assoc.ServiceAccountName)
	}

	return nil
}

func (s *Service) deletePodIdentityAssociation(ctx context.Context, assocID string) error {
	clusterName := s.scope.KubernetesClusterName()

	input := &eks.DeletePodIdentityAssociationInput{
		AssociationId: &assocID,
		ClusterName:   &clusterName,
	}

	if _, err := s.EKSClient.DeletePodIdentityAssociation(ctx, input); err != nil {
		return errors.Wrapf(err, "failed to delete pod identity association %s", assocID)
	}

	return nil
}
