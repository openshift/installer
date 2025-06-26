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

// Package iamauth provides a way to interact with AWS IAM.
package iamauth

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// Service defines the specs for a service.
type Service struct {
	scope     scope.IAMAuthScope
	backend   BackendType
	client    client.Client
	IAMClient IAMAPI
}

// IAMAPI defines the interface for IAM operations.
type IAMAPI interface {
	GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error)
	CreateRole(ctx context.Context, params *iam.CreateRoleInput, optFns ...func(*iam.Options)) (*iam.CreateRoleOutput, error)
	DeleteRole(ctx context.Context, params *iam.DeleteRoleInput, optFns ...func(*iam.Options)) (*iam.DeleteRoleOutput, error)
	UpdateAssumeRolePolicy(ctx context.Context, params *iam.UpdateAssumeRolePolicyInput, optFns ...func(*iam.Options)) (*iam.UpdateAssumeRolePolicyOutput, error)
	GetPolicy(ctx context.Context, params *iam.GetPolicyInput, optFns ...func(*iam.Options)) (*iam.GetPolicyOutput, error)
	AttachRolePolicy(ctx context.Context, params *iam.AttachRolePolicyInput, optFns ...func(*iam.Options)) (*iam.AttachRolePolicyOutput, error)
	DetachRolePolicy(ctx context.Context, params *iam.DetachRolePolicyInput, optFns ...func(*iam.Options)) (*iam.DetachRolePolicyOutput, error)
	ListAttachedRolePolicies(ctx context.Context, params *iam.ListAttachedRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedRolePoliciesOutput, error)
	CreateOpenIDConnectProvider(ctx context.Context, params *iam.CreateOpenIDConnectProviderInput, optFns ...func(*iam.Options)) (*iam.CreateOpenIDConnectProviderOutput, error)
	GetOpenIDConnectProvider(ctx context.Context, params *iam.GetOpenIDConnectProviderInput, optFns ...func(*iam.Options)) (*iam.GetOpenIDConnectProviderOutput, error)
	DeleteOpenIDConnectProvider(ctx context.Context, params *iam.DeleteOpenIDConnectProviderInput, optFns ...func(*iam.Options)) (*iam.DeleteOpenIDConnectProviderOutput, error)
	ListOpenIDConnectProviders(ctx context.Context, params *iam.ListOpenIDConnectProvidersInput, optFns ...func(*iam.Options)) (*iam.ListOpenIDConnectProvidersOutput, error)
	TagOpenIDConnectProvider(ctx context.Context, params *iam.TagOpenIDConnectProviderInput, optFns ...func(*iam.Options)) (*iam.TagOpenIDConnectProviderOutput, error)
	TagRole(ctx context.Context, params *iam.TagRoleInput, optFns ...func(*iam.Options)) (*iam.TagRoleOutput, error)
	UntagRole(ctx context.Context, params *iam.UntagRoleInput, optFns ...func(*iam.Options)) (*iam.UntagRoleOutput, error)
}

// NewService will create a new Service object.
func NewService(iamScope scope.IAMAuthScope, backend BackendType, client client.Client) *Service {
	return &Service{
		scope:     iamScope,
		backend:   backend,
		client:    client,
		IAMClient: scope.NewIAMClient(iamScope, iamScope, iamScope, iamScope.InfraCluster()),
	}
}
