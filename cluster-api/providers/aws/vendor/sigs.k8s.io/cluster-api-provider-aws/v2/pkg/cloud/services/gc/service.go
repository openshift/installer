/*
Copyright 2022 The Kubernetes Authors.

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

package gc

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// Service is used to perform operations against a tenant/workload/child cluster.
type Service struct {
	scope                 cloud.ClusterScoper
	elbClient             elbiface.ELBAPI
	elbv2Client           elbv2iface.ELBV2API
	resourceTaggingClient resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
	ec2Client             ec2iface.EC2API
	cleanupFuncs          ResourceCleanupFuncs
	collectFuncs          ResourceCollectFuncs
}

// NewService creates a new Service.
func NewService(clusterScope cloud.ClusterScoper, opts ...ServiceOption) *Service {
	svc := &Service{
		scope:                 clusterScope,
		elbClient:             scope.NewELBClient(clusterScope, clusterScope, clusterScope, clusterScope.InfraCluster()),
		elbv2Client:           scope.NewELBv2Client(clusterScope, clusterScope, clusterScope, clusterScope.InfraCluster()),
		resourceTaggingClient: scope.NewResourgeTaggingClient(clusterScope, clusterScope, clusterScope, clusterScope.InfraCluster()),
		ec2Client:             scope.NewEC2Client(clusterScope, clusterScope, clusterScope, clusterScope.InfraCluster()),
		cleanupFuncs:          ResourceCleanupFuncs{},
		collectFuncs:          ResourceCollectFuncs{},
	}
	addDefaultCleanupFuncs(svc)

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func addDefaultCleanupFuncs(s *Service) {
	s.cleanupFuncs = []ResourceCleanupFunc{
		s.deleteLoadBalancers,
		s.deleteTargetGroups,
		s.deleteSecurityGroups,
	}
}

func addDefaultCollectFuncs(s *Service) {
	s.collectFuncs = []ResourceCollectFunc{
		s.defaultGetResources,
	}
}

func addAlternativeCollectFuncs(s *Service) {
	s.collectFuncs = []ResourceCollectFunc{
		s.getProviderOwnedLoadBalancers,
		s.getProviderOwnedLoadBalancersV2,
		s.getProviderOwnedTargetgroups,
		s.getProviderOwnedSecurityGroups,
	}
}

// AWSResource represents a resource in AWS.
type AWSResource struct {
	ARN  *arn.ARN
	Tags map[string]string
}

// ResourceCleanupFunc is a function type to cleaning up resources for a specific AWS service type.
type ResourceCleanupFunc func(ctx context.Context, resources []*AWSResource) error

// ResourceCleanupFuncs is a collection of ResourceCleanupFunc.
type ResourceCleanupFuncs []ResourceCleanupFunc

// Execute will execute all the defined clean up functions against the aws resources.
func (fn ResourceCleanupFuncs) Execute(ctx context.Context, resources []*AWSResource) error {
	for _, f := range fn {
		if err := f(ctx, resources); err != nil {
			return err
		}
	}

	return nil
}

// ResourceCollectFunc is a function type to collect resources for a specific AWS service type.
type ResourceCollectFunc func(ctx context.Context) ([]*AWSResource, error)

// ResourceCollectFuncs is a collection of ResourceCollectFunc.
type ResourceCollectFuncs []ResourceCollectFunc

// Execute will execute all the defined collect functions against the aws resources.
func (fn ResourceCollectFuncs) Execute(ctx context.Context) ([]*AWSResource, error) {
	var resources []*AWSResource
	for _, f := range fn {
		rs, err := f(ctx)
		if err != nil {
			return nil, err
		}

		resources = append(resources, rs...)
	}

	return resources, nil
}
