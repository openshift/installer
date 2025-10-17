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

// Package securitygroup provides a service to manage AWS security group resources.
package securitygroup

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/common"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope     scope.SGScope
	roles     []infrav1.SecurityGroupRole
	EC2Client common.EC2API
}

// NewService returns a new service given the api clients with a defined
// set of roles.
func NewService(sgScope scope.SGScope, roles []infrav1.SecurityGroupRole) *Service {
	return &Service{
		scope:     sgScope,
		roles:     roles,
		EC2Client: scope.NewEC2Client(sgScope, sgScope, sgScope, sgScope.InfraCluster()),
	}
}
