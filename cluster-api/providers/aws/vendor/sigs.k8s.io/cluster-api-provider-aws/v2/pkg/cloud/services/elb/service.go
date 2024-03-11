/*
Copyright 2018 The Kubernetes Authors.

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

package elb

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope                 scope.ELBScope
	EC2Client             ec2iface.EC2API
	ELBClient             elbiface.ELBAPI
	ELBV2Client           elbv2iface.ELBV2API
	ResourceTaggingClient resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
}

// NewService returns a new service given the api clients.
func NewService(elbScope scope.ELBScope) *Service {
	return &Service{
		scope:                 elbScope,
		EC2Client:             scope.NewEC2Client(elbScope, elbScope, elbScope, elbScope.InfraCluster()),
		ELBClient:             scope.NewELBClient(elbScope, elbScope, elbScope, elbScope.InfraCluster()),
		ELBV2Client:           scope.NewELBv2Client(elbScope, elbScope, elbScope, elbScope.InfraCluster()),
		ResourceTaggingClient: scope.NewResourgeTaggingClient(elbScope, elbScope, elbScope, elbScope.InfraCluster()),
	}
}
