/*
Copyright 2020 The Kubernetes Authors.

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

// Package eks provides a service to reconcile EKS control plane and nodegroups.
package eks

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/klog/v2"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// ReconcileControlPlane reconciles a EKS control plane.
func (s *Service) ReconcileControlPlane(ctx context.Context) error {
	s.scope.Debug("Reconciling EKS control plane", "cluster", klog.KRef(s.scope.Cluster.Namespace, s.scope.Cluster.Name))

	// Control Plane IAM Role
	if err := s.reconcileControlPlaneIAMRole(); err != nil {
		conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1.IAMControlPlaneRolesReadyCondition, ekscontrolplanev1.IAMControlPlaneRolesReconciliationFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
		return err
	}
	conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.IAMControlPlaneRolesReadyCondition)

	// EKS Cluster
	if err := s.reconcileCluster(ctx); err != nil {
		conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneReadyCondition, ekscontrolplanev1.EKSControlPlaneReconciliationFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
		return err
	}
	conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneReadyCondition)

	// EKS Addons
	if err := s.reconcileAddons(ctx); err != nil {
		conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1.EKSAddonsConfiguredCondition, ekscontrolplanev1.EKSAddonsConfiguredFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
		return errors.Wrap(err, "failed reconciling eks addons")
	}
	conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSAddonsConfiguredCondition)

	// EKS Identity Provider
	if err := s.reconcileIdentityProvider(ctx); err != nil {
		conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1.EKSIdentityProviderConfiguredCondition, ekscontrolplanev1.EKSIdentityProviderConfiguredFailedReason, clusterv1.ConditionSeverityWarning, "%s", err.Error())
		return errors.Wrap(err, "failed reconciling eks identity provider")
	}
	conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSIdentityProviderConfiguredCondition)

	s.scope.Debug("Reconcile EKS control plane completed successfully")
	return nil
}

// DeleteControlPlane deletes the EKS control plane.
func (s *Service) DeleteControlPlane() (err error) {
	s.scope.Debug("Deleting EKS control plane")

	// EKS Cluster
	if err := s.deleteCluster(); err != nil {
		return err
	}

	// Control Plane IAM role
	if err := s.deleteControlPlaneIAMRole(); err != nil {
		return err
	}

	// OIDC Provider
	if err := s.deleteOIDCProvider(); err != nil {
		return err
	}

	s.scope.Debug("Delete EKS control plane completed successfully")
	return nil
}

// ReconcilePool is the entrypoint for ManagedMachinePool reconciliation.
func (s *NodegroupService) ReconcilePool(ctx context.Context) error {
	s.scope.Debug("Reconciling EKS nodegroup")

	if err := s.reconcileNodegroupIAMRole(); err != nil {
		conditions.MarkFalse(
			s.scope.ManagedMachinePool,
			expinfrav1.IAMNodegroupRolesReadyCondition,
			expinfrav1.IAMNodegroupRolesReconciliationFailedReason,
			clusterv1.ConditionSeverityError,
			"%s",
			err.Error(),
		)
		return err
	}
	conditions.MarkTrue(s.scope.ManagedMachinePool, expinfrav1.IAMNodegroupRolesReadyCondition)

	if err := s.reconcileNodegroup(ctx); err != nil {
		conditions.MarkFalse(
			s.scope.ManagedMachinePool,
			expinfrav1.EKSNodegroupReadyCondition,
			expinfrav1.EKSNodegroupReconciliationFailedReason,
			clusterv1.ConditionSeverityError,
			"%s",
			err.Error(),
		)
		return err
	}
	conditions.MarkTrue(s.scope.ManagedMachinePool, expinfrav1.EKSNodegroupReadyCondition)

	return nil
}

// ReconcilePoolDelete is the entrypoint for ManagedMachinePool deletion
// reconciliation.
func (s *NodegroupService) ReconcilePoolDelete() error {
	s.scope.Debug("Reconciling deletion of EKS nodegroup")

	eksNodegroupName := s.scope.NodegroupName()

	ng, err := s.describeNodegroup()
	if err != nil {
		if awserrors.IsNotFound(err) {
			s.scope.Trace("EKS nodegroup does not exist")
			return nil
		}
		return errors.Wrap(err, "failed to describe EKS nodegroup")
	}
	if ng == nil {
		return nil
	}

	if err := s.deleteNodegroupAndWait(); err != nil {
		return errors.Wrap(err, "failed to delete nodegroup")
	}

	if err := s.deleteNodegroupIAMRole(); err != nil {
		return errors.Wrap(err, "failed to delete nodegroup IAM role")
	}

	record.Eventf(s.scope.ManagedMachinePool, "SuccessfulDeleteEKSNodegroup", "Deleted EKS nodegroup %s", eksNodegroupName)

	return nil
}
