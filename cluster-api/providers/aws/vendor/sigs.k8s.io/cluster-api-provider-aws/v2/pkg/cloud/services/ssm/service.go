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
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope     cloud.ClusterScoper
	SSMClient ssmiface.SSMAPI
}

// NewService returns a new service given the api clients.
func NewService(secretsScope cloud.ClusterScoper) *Service {
	return &Service{
		scope:     secretsScope,
		SSMClient: scope.NewSSMClient(secretsScope, secretsScope, secretsScope, secretsScope.InfraCluster()),
	}
}
