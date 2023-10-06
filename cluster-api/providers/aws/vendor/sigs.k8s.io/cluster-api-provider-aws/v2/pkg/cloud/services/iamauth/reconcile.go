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

package iamauth

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
)

// ReconcileIAMAuthenticator is used to create the aws-iam-authenticator in a cluster.
func (s *Service) ReconcileIAMAuthenticator(ctx context.Context) error {
	s.scope.Info("Reconciling aws-iam-authenticator configuration", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()))

	remoteClient, err := s.scope.RemoteClient()
	if err != nil {
		s.scope.Error(err, "getting client for remote cluster")
		return fmt.Errorf("getting client for remote cluster: %w", err)
	}

	authBackend, err := NewBackend(s.backend, remoteClient)
	if err != nil {
		return fmt.Errorf("getting aws-iam-authenticator backend: %w", err)
	}
	nodeRoles, err := s.getRolesForWorkers(ctx)
	if err != nil {
		s.scope.Error(err, "getting roles for remote workers")
		return fmt.Errorf("getting roles for remote workers: %w", err)
	}
	for roleName := range nodeRoles {
		roleARN, err := s.getARNForRole(roleName)
		if err != nil {
			return fmt.Errorf("failed to get ARN for role %s: %w", roleARN, err)
		}
		nodesRoleMapping := ekscontrolplanev1.RoleMapping{
			RoleARN: roleARN,
			KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
				UserName: EC2NodeUserName,
				Groups:   NodeGroups,
			},
		}
		s.scope.Debug("Mapping node IAM role", "iam-role", nodesRoleMapping.RoleARN, "user", nodesRoleMapping.UserName)
		if err := authBackend.MapRole(nodesRoleMapping); err != nil {
			return fmt.Errorf("mapping iam node role: %w", err)
		}
	}

	s.scope.Debug("Mapping additional IAM roles and users")
	iamCfg := s.scope.IAMAuthConfig()
	for _, roleMapping := range iamCfg.RoleMappings {
		s.scope.Debug("Mapping IAM role", "iam-role", roleMapping.RoleARN, "user", roleMapping.UserName)
		if err := authBackend.MapRole(roleMapping); err != nil {
			return fmt.Errorf("mapping iam role: %w", err)
		}
	}

	for _, userMapping := range iamCfg.UserMappings {
		s.scope.Debug("Mapping IAM user", "iam-user", userMapping.UserARN, "user", userMapping.UserName)
		if err := authBackend.MapUser(userMapping); err != nil {
			return fmt.Errorf("mapping iam user: %w", err)
		}
	}

	s.scope.Info("Reconciled aws-iam-authenticator configuration", "cluster", klog.KRef("", s.scope.Name()))

	return nil
}

func (s *Service) getARNForRole(role string) (string, error) {
	input := &iam.GetRoleInput{
		RoleName: aws.String(role),
	}
	out, err := s.IAMClient.GetRole(input)
	if err != nil {
		return "", errors.Wrap(err, "unable to get role")
	}
	return aws.StringValue(out.Role.Arn), nil
}

func (s *Service) getRolesForWorkers(ctx context.Context) (map[string]struct{}, error) {
	allRoles := map[string]struct{}{}
	if err := s.getRolesForMachineDeployments(ctx, allRoles); err != nil {
		return nil, fmt.Errorf("failed to get roles from machine deployments %w", err)
	}
	if err := s.getRolesForMachinePools(ctx, allRoles); err != nil {
		return nil, fmt.Errorf("failed to get roles from machine pools %w", err)
	}
	return allRoles, nil
}

func (s *Service) getRolesForMachineDeployments(ctx context.Context, allRoles map[string]struct{}) error {
	deploymentList := &clusterv1.MachineDeploymentList{}
	selectors := []client.ListOption{
		client.InNamespace(s.scope.Namespace()),
		client.MatchingLabels{
			clusterv1.ClusterNameLabel: s.scope.Name(),
		},
	}
	err := s.client.List(ctx, deploymentList, selectors...)
	if err != nil {
		return fmt.Errorf("failed to list machine deployments for cluster %s/%s: %w", s.scope.Namespace(), s.scope.Name(), err)
	}

	for _, deployment := range deploymentList.Items {
		ref := deployment.Spec.Template.Spec.InfrastructureRef
		if ref.Kind != "AWSMachineTemplate" {
			continue
		}
		awsMachineTemplate := &infrav1.AWSMachineTemplate{}
		err := s.client.Get(ctx, client.ObjectKey{
			Name:      ref.Name,
			Namespace: s.scope.Namespace(),
		}, awsMachineTemplate)
		if err != nil {
			return fmt.Errorf("failed to get AWSMachine %s/%s: %w", ref.Namespace, ref.Name, err)
		}
		instanceProfile := awsMachineTemplate.Spec.Template.Spec.IAMInstanceProfile
		if _, ok := allRoles[instanceProfile]; !ok && instanceProfile != "" {
			allRoles[instanceProfile] = struct{}{}
		}
	}
	return nil
}

func (s *Service) getRolesForMachinePools(ctx context.Context, allRoles map[string]struct{}) error {
	machinePoolList := &expclusterv1.MachinePoolList{}
	selectors := []client.ListOption{
		client.InNamespace(s.scope.Namespace()),
		client.MatchingLabels{
			clusterv1.ClusterNameLabel: s.scope.Name(),
		},
	}
	err := s.client.List(ctx, machinePoolList, selectors...)
	if err != nil {
		return fmt.Errorf("failed to list machine pools for cluster %s/%s: %w", s.scope.Namespace(), s.scope.Name(), err)
	}
	for _, pool := range machinePoolList.Items {
		ref := pool.Spec.Template.Spec.InfrastructureRef
		switch ref.Kind {
		case "AWSMachinePool":
			if err := s.getRolesForAWSMachinePool(ctx, ref, allRoles); err != nil {
				return err
			}
		case "AWSManagedMachinePool":
			if err := s.getRolesForAWSManagedMachinePool(ctx, ref, allRoles); err != nil {
				return err
			}
		default:
		}
	}
	return nil
}

func (s *Service) getRolesForAWSMachinePool(ctx context.Context, ref corev1.ObjectReference, allRoles map[string]struct{}) error {
	awsMachinePool := &expinfrav1.AWSMachinePool{}
	err := s.client.Get(ctx, client.ObjectKey{
		Name:      ref.Name,
		Namespace: s.scope.Namespace(),
	}, awsMachinePool)
	if err != nil {
		return fmt.Errorf("failed to get AWSMachine %s/%s: %w", ref.Namespace, ref.Name, err)
	}
	instanceProfile := awsMachinePool.Spec.AWSLaunchTemplate.IamInstanceProfile
	if _, ok := allRoles[instanceProfile]; !ok && instanceProfile != "" {
		allRoles[instanceProfile] = struct{}{}
	}
	return nil
}

func (s *Service) getRolesForAWSManagedMachinePool(ctx context.Context, ref corev1.ObjectReference, allRoles map[string]struct{}) error {
	awsManagedMachinePool := &expinfrav1.AWSManagedMachinePool{}
	err := s.client.Get(ctx, client.ObjectKey{
		Name:      ref.Name,
		Namespace: s.scope.Namespace(),
	}, awsManagedMachinePool)
	if err != nil {
		return fmt.Errorf("failed to get AWSMachine %s/%s: %w", ref.Namespace, ref.Name, err)
	}
	instanceProfile := awsManagedMachinePool.Spec.RoleName
	if _, ok := allRoles[instanceProfile]; !ok && instanceProfile != "" {
		allRoles[instanceProfile] = struct{}{}
	}
	return nil
}
