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
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// Service defines the specs for a service.
type Service struct {
	scope     scope.IAMAuthScope
	backend   BackendType
	client    client.Client
	IAMClient iamiface.IAMAPI
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
