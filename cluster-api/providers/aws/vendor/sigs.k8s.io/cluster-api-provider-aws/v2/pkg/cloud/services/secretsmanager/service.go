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

// Package secretsmanager provides a way to interact with AWS Secrets Manager.
package secretsmanager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope                cloud.ClusterScoper
	SecretsManagerClient SecretsManagerAPI
}

// SecretsManagerAPI is the subset of the AWS Secrets Manager API that is used by CAPA.
type SecretsManagerAPI interface {
	CreateSecret(ctx context.Context, params *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error)
	DeleteSecret(ctx context.Context, params *secretsmanager.DeleteSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.DeleteSecretOutput, error)
}

var _ SecretsManagerAPI = &secretsmanager.Client{}

// NewService returns a new service given the api clients.
func NewService(secretsScope cloud.ClusterScoper) *Service {
	return &Service{
		scope:                secretsScope,
		SecretsManagerClient: scope.NewSecretsManagerClient(secretsScope, secretsScope, secretsScope, secretsScope.InfraCluster()),
	}
}
