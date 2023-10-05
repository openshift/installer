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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/eks/identityprovider"
)

func (s *Service) reconcileIdentityProvider(ctx context.Context) error {
	s.scope.Info("reconciling oidc identity provider")
	if s.scope.ControlPlane.Spec.OIDCIdentityProviderConfig == nil {
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

	s.scope.V(2).Info("creating oidc provider plan", "desired", desired, "current", current)

	identityProviderPlan := identityprovider.NewPlan(clusterName, current, desired, s.EKSClient, s.scope.Logger)

	procedures, err := identityProviderPlan.Create(ctx)
	if err != nil {
		s.scope.Error(err, "failed creating eks identity provider plan")
		return fmt.Errorf("creating eks identity provider plan: %w", err)
	}
	s.scope.V(2).Info("computed EKS identity provider plan", "numprocs", len(procedures))

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

	if latest != nil {
		s.scope.ControlPlane.Status.IdentityProviderStatus = ekscontrolplanev1.IdentityProviderStatus{
			ARN:    aws.StringValue(latest.IdentityProviderConfigArn),
			Status: aws.StringValue(latest.Status),
		}

		err := s.scope.PatchObject()
		if err != nil {
			return errors.Wrap(err, "updating identity provider status")
		}
	}

	return nil
}

func (s *Service) getAssociatedIdentityProvider(ctx context.Context, clusterName string) (*identityprovider.OidcIdentityProviderConfig, error) {
	list, err := s.EKSClient.ListIdentityProviderConfigsWithContext(ctx, &eks.ListIdentityProviderConfigsInput{
		ClusterName: aws.String(clusterName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "listing identity provider configs")
	}

	// these is only one identity provider

	if len(list.IdentityProviderConfigs) == 0 {
		return nil, nil
	}

	providerconfig, err := s.EKSClient.DescribeIdentityProviderConfigWithContext(ctx, &eks.DescribeIdentityProviderConfigInput{
		ClusterName:            aws.String(clusterName),
		IdentityProviderConfig: list.IdentityProviderConfigs[0],
	})

	if err != nil {
		return nil, errors.Wrap(err, "describing identity provider config")
	}

	config := providerconfig.IdentityProviderConfig.Oidc

	return &identityprovider.OidcIdentityProviderConfig{
		ClientID:                   *config.ClientId,
		GroupsClaim:                config.GroupsClaim,
		GroupsPrefix:               config.GroupsPrefix,
		IdentityProviderConfigArn:  config.IdentityProviderConfigArn,
		IdentityProviderConfigName: *config.IdentityProviderConfigName,
		IssuerURL:                  *config.IssuerUrl,
		RequiredClaims:             config.RequiredClaims,
		Status:                     config.Status,
		Tags:                       converters.MapPtrToMap(config.Tags),
		UsernameClaim:              config.UsernameClaim,
		UsernamePrefix:             config.UsernamePrefix,
	}, nil
}
