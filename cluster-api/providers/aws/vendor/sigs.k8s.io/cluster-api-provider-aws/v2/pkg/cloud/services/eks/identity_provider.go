/*
Copyright 2021 The Kubernetes Authors.

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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/pkg/errors"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks/identityprovider"
)

func (s *Service) reconcileIdentityProvider(ctx context.Context) error {
	s.scope.Info("reconciling oidc identity provider")
	if s.scope.OIDCIdentityProviderConfig() == nil {
		s.scope.Info("no oidc provider config, skipping reconcile")
		return nil
	}

	clusterName := s.scope.KubernetesClusterName()
	current, err := s.getAssociatedIdentityProvider(ctx, clusterName)
	if err != nil {
		return errors.Wrap(err, "unable to list associated identity providers")
	}

	desired := converters.ConvertSDKToIdentityProvider(s.scope.OIDCIdentityProviderConfig())

	if desired == nil && current == nil {
		s.scope.Info("no identity provider required or installed, no action needed")
		return nil
	}

	s.scope.Debug("creating oidc provider plan", "desired", desired, "current", current)
	procedures, err := identityprovider.
		NewPlan(clusterName, current, desired, s.EKSClient, s.scope).
		Create(ctx)
	if err != nil {
		s.scope.Error(err, "failed creating eks identity provider plan")
		return fmt.Errorf("creating eks identity provider plan: %w", err)
	}

	// nothing will be done, we can leave
	if len(procedures) == 0 {
		return nil
	}

	s.scope.Debug("computed EKS identity provider plan", "numprocs", len(procedures))

	// Perform required operations
	for _, procedure := range procedures {
		s.scope.Info("Executing identity provider procedure", "name", procedure.Name())
		if err := procedure.Do(ctx); err != nil {
			s.scope.Error(err, "failed executing identity provider procedure", "name", procedure.Name())
			return fmt.Errorf("%s: %w", procedure.Name(), err)
		}
	}

	latest, err := s.getAssociatedIdentityProvider(ctx, clusterName)
	if err != nil {
		return errors.Wrap(err, "getting associated identity provider")
	}

	if latest == nil {
		return nil
	}

	// don't patch if arn/status is the same
	if latest.IdentityProviderConfigArn == s.scope.ControlPlane.Status.IdentityProviderStatus.ARN &&
		latest.Status == s.scope.ControlPlane.Status.IdentityProviderStatus.Status {
		return nil
	}

	// idp status has changed, patch the control plane
	s.scope.ControlPlane.Status.IdentityProviderStatus = ekscontrolplanev1.IdentityProviderStatus{
		ARN:    latest.IdentityProviderConfigArn,
		Status: latest.Status,
	}

	if err := s.scope.PatchObject(); err != nil {
		return errors.Wrap(err, "updating identity provider status")
	}

	return nil
}

func (s *Service) getAssociatedIdentityProvider(ctx context.Context, clusterName string) (*identityprovider.OidcIdentityProviderConfig, error) {
	list, err := s.EKSClient.ListIdentityProviderConfigs(ctx, &eks.ListIdentityProviderConfigsInput{
		ClusterName: aws.String(clusterName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "listing identity provider configs")
	}

	// these is only one identity provider

	if len(list.IdentityProviderConfigs) == 0 {
		return nil, nil
	}

	providerconfig, err := s.EKSClient.DescribeIdentityProviderConfig(ctx, &eks.DescribeIdentityProviderConfigInput{
		ClusterName:            aws.String(clusterName),
		IdentityProviderConfig: &list.IdentityProviderConfigs[0],
	})

	if err != nil {
		return nil, errors.Wrap(err, "describing identity provider config")
	}

	config := providerconfig.IdentityProviderConfig.Oidc

	return &identityprovider.OidcIdentityProviderConfig{
		ClientID:                   aws.ToString(config.ClientId),
		GroupsClaim:                aws.ToString(config.GroupsClaim),
		GroupsPrefix:               aws.ToString(config.GroupsPrefix),
		IdentityProviderConfigArn:  aws.ToString(config.IdentityProviderConfigArn),
		IdentityProviderConfigName: aws.ToString(config.IdentityProviderConfigName),
		IssuerURL:                  aws.ToString(config.IssuerUrl),
		RequiredClaims:             config.RequiredClaims,
		Status:                     string(config.Status),
		Tags:                       converters.MapPtrToMap(config.Tags),
		UsernameClaim:              aws.ToString(config.UsernameClaim),
		UsernamePrefix:             aws.ToString(config.UsernamePrefix),
	}, nil
}
