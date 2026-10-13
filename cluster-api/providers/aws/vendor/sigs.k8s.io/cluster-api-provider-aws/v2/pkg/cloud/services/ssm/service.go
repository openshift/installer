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

package ssm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope     cloud.ClusterScoper
	SSMClient SSMAPI
}

// SSMAPI defines the interface for interacting with AWS SSM Parameter Store.
type SSMAPI interface {
	PutParameter(ctx context.Context, input *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
	DeleteParameter(ctx context.Context, input *ssm.DeleteParameterInput, optFns ...func(*ssm.Options)) (*ssm.DeleteParameterOutput, error)
	GetParameter(ctx context.Context, input *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	// Add more methods as needed
}

// Ensure ssm.Client satisfies the SSMAPI interface.
var _ SSMAPI = &ssm.Client{}

// NewService creates a new Service for managing secrets in AWS SSM.
func NewService(secretsScope cloud.ClusterScoper) *Service {
	return &Service{
		scope:     secretsScope,
		SSMClient: scope.NewSSMClient(secretsScope, secretsScope, secretsScope, secretsScope.InfraCluster()),
	}
}
