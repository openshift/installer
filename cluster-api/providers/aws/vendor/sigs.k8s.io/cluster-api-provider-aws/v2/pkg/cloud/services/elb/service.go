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

// Package elb provides a service for managing AWS load balancers.
package elb

import (
	"context"
	"time"

	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	rgapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/common"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/network"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope                 scope.ELBScope
	EC2Client             common.EC2API
	ELBClient             ELBAPI
	ELBV2Client           ELBV2API
	ResourceTaggingClient ResourceGroupsTaggingAPIAPI
	netService            *network.Service
}

// ELBAPI is the subset of the AWS ELB API used by CAPA.
type ELBAPI interface {
	// Subset of AWS ELB API
	AddTags(ctx context.Context, params *elb.AddTagsInput, optFns ...func(*elb.Options)) (*elb.AddTagsOutput, error)
	ApplySecurityGroupsToLoadBalancer(ctx context.Context, params *elb.ApplySecurityGroupsToLoadBalancerInput, optFns ...func(*elb.Options)) (*elb.ApplySecurityGroupsToLoadBalancerOutput, error)
	AttachLoadBalancerToSubnets(ctx context.Context, params *elb.AttachLoadBalancerToSubnetsInput, optFns ...func(*elb.Options)) (*elb.AttachLoadBalancerToSubnetsOutput, error)
	ConfigureHealthCheck(ctx context.Context, params *elb.ConfigureHealthCheckInput, optFns ...func(*elb.Options)) (*elb.ConfigureHealthCheckOutput, error)
	CreateLoadBalancer(ctx context.Context, params *elb.CreateLoadBalancerInput, optFns ...func(*elb.Options)) (*elb.CreateLoadBalancerOutput, error)
	DeleteLoadBalancer(ctx context.Context, params *elb.DeleteLoadBalancerInput, optFns ...func(*elb.Options)) (*elb.DeleteLoadBalancerOutput, error)
	DeregisterInstancesFromLoadBalancer(ctx context.Context, params *elb.DeregisterInstancesFromLoadBalancerInput, optFns ...func(*elb.Options)) (*elb.DeregisterInstancesFromLoadBalancerOutput, error)
	DescribeLoadBalancerAttributes(ctx context.Context, params *elb.DescribeLoadBalancerAttributesInput, optFns ...func(*elb.Options)) (*elb.DescribeLoadBalancerAttributesOutput, error)
	DescribeLoadBalancers(ctx context.Context, params *elb.DescribeLoadBalancersInput, optFns ...func(*elb.Options)) (*elb.DescribeLoadBalancersOutput, error)
	DescribeTags(ctx context.Context, params *elb.DescribeTagsInput, optFns ...func(*elb.Options)) (*elb.DescribeTagsOutput, error)
	ModifyLoadBalancerAttributes(ctx context.Context, params *elb.ModifyLoadBalancerAttributesInput, optFns ...func(*elb.Options)) (*elb.ModifyLoadBalancerAttributesOutput, error)
	RemoveTags(ctx context.Context, params *elb.RemoveTagsInput, optFns ...func(*elb.Options)) (*elb.RemoveTagsOutput, error)
	RegisterInstancesWithLoadBalancer(ctx context.Context, params *elb.RegisterInstancesWithLoadBalancerInput, optFns ...func(*elb.Options)) (*elb.RegisterInstancesWithLoadBalancerOutput, error)

	// CAPA Custom function
	DescribeLoadBalancersPages(ctx context.Context, input *elb.DescribeLoadBalancersInput, fn func(*elb.DescribeLoadBalancersOutput)) error
}

// ELBV2API is the subset of the AWS ELBV2 API used by CAPA.
type ELBV2API interface {
	// Subset of AWS ELBV2 API
	AddTags(ctx context.Context, params *elbv2.AddTagsInput, optFns ...func(*elbv2.Options)) (*elbv2.AddTagsOutput, error)
	CreateListener(ctx context.Context, params *elbv2.CreateListenerInput, optFns ...func(*elbv2.Options)) (*elbv2.CreateListenerOutput, error)
	CreateLoadBalancer(ctx context.Context, params *elbv2.CreateLoadBalancerInput, optFns ...func(*elbv2.Options)) (*elbv2.CreateLoadBalancerOutput, error)
	CreateTargetGroup(ctx context.Context, params *elbv2.CreateTargetGroupInput, optFns ...func(*elbv2.Options)) (*elbv2.CreateTargetGroupOutput, error)
	DeleteListener(ctx context.Context, params *elbv2.DeleteListenerInput, optFns ...func(*elbv2.Options)) (*elbv2.DeleteListenerOutput, error)
	DeleteLoadBalancer(ctx context.Context, params *elbv2.DeleteLoadBalancerInput, optFns ...func(*elbv2.Options)) (*elbv2.DeleteLoadBalancerOutput, error)
	DeleteTargetGroup(ctx context.Context, params *elbv2.DeleteTargetGroupInput, optFns ...func(*elbv2.Options)) (*elbv2.DeleteTargetGroupOutput, error)
	DeregisterTargets(ctx context.Context, params *elbv2.DeregisterTargetsInput, optFns ...func(*elbv2.Options)) (*elbv2.DeregisterTargetsOutput, error)
	DescribeListeners(ctx context.Context, params *elbv2.DescribeListenersInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeListenersOutput, error)
	DescribeLoadBalancerAttributes(ctx context.Context, params *elbv2.DescribeLoadBalancerAttributesInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeLoadBalancerAttributesOutput, error)
	DescribeLoadBalancers(ctx context.Context, params *elbv2.DescribeLoadBalancersInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeLoadBalancersOutput, error)
	DescribeTags(ctx context.Context, params *elbv2.DescribeTagsInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeTagsOutput, error)
	DescribeTargetGroups(ctx context.Context, params *elbv2.DescribeTargetGroupsInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeTargetGroupsOutput, error)
	DescribeTargetHealth(ctx context.Context, params *elbv2.DescribeTargetHealthInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeTargetHealthOutput, error)
	ModifyListener(ctx context.Context, params *elbv2.ModifyListenerInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyListenerOutput, error)
	ModifyLoadBalancerAttributes(ctx context.Context, params *elbv2.ModifyLoadBalancerAttributesInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyLoadBalancerAttributesOutput, error)
	ModifyTargetGroupAttributes(ctx context.Context, params *elbv2.ModifyTargetGroupAttributesInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyTargetGroupAttributesOutput, error)
	RegisterTargets(ctx context.Context, params *elbv2.RegisterTargetsInput, optFns ...func(*elbv2.Options)) (*elbv2.RegisterTargetsOutput, error)
	RemoveTags(ctx context.Context, params *elbv2.RemoveTagsInput, optFns ...func(*elbv2.Options)) (*elbv2.RemoveTagsOutput, error)
	SetSecurityGroups(ctx context.Context, params *elbv2.SetSecurityGroupsInput, optFns ...func(*elbv2.Options)) (*elbv2.SetSecurityGroupsOutput, error)
	SetSubnets(ctx context.Context, params *elbv2.SetSubnetsInput, optFns ...func(*elbv2.Options)) (*elbv2.SetSubnetsOutput, error)

	// CAPA Custom function
	DescribeLoadBalancersPages(ctx context.Context, input *elbv2.DescribeLoadBalancersInput, fn func(*elbv2.DescribeLoadBalancersOutput)) error
	WaitUntilLoadBalancerAvailable(ctx context.Context, input *elbv2.DescribeLoadBalancersInput, maxWait time.Duration) error
}

// ResourceGroupsTaggingAPIAPI is the subset of the AWS ResourceGroupsTaggingAPI API used by CAPA.
type ResourceGroupsTaggingAPIAPI interface {
	GetResources(ctx context.Context, params *rgapi.GetResourcesInput, optFns ...func(*rgapi.Options)) (*rgapi.GetResourcesOutput, error)

	GetResourcesPages(ctx context.Context, input *rgapi.GetResourcesInput, fn func(*rgapi.GetResourcesOutput)) error
}

// ELBClient is a wrapper over elb.Client for implementing custom methods of ELBAPI.
type ELBClient struct {
	*elb.Client
}

// ELBV2Client is a wrapper over elbv2.Client for implementing custom methods of ELBV2API.
type ELBV2Client struct {
	*elbv2.Client
}

// ResourceGroupsTaggingAPIClient is a wrapper over rgapi.Client for implementing custom methods of ResourceGroupsTaggingAPI.
type ResourceGroupsTaggingAPIClient struct {
	*rgapi.Client
}

// NewService returns a new service given the api clients.
func NewService(elbScope scope.ELBScope) *Service {
	return &Service{
		scope:     elbScope,
		EC2Client: scope.NewEC2Client(elbScope, elbScope, elbScope, elbScope.InfraCluster()),
		ELBClient: &ELBClient{
			Client: scope.NewELBClient(elbScope, elbScope, elbScope, elbScope.InfraCluster()),
		},
		ELBV2Client: &ELBV2Client{
			Client: scope.NewELBv2Client(elbScope, elbScope, elbScope, elbScope.InfraCluster()),
		},
		ResourceTaggingClient: &ResourceGroupsTaggingAPIClient{
			Client: scope.NewResourgeTaggingClient(elbScope, elbScope, elbScope, elbScope.InfraCluster()),
		},
		netService: network.NewService(elbScope.(scope.NetworkScope)),
	}
}
