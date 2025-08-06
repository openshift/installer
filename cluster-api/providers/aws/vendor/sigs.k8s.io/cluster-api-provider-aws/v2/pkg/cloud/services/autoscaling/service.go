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

package asg

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/common"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the asg client.
type Service struct {
	scope     cloud.ClusterScoper
	ASGClient AutoScalingAPI
	EC2Client common.EC2API
}

// AutoScalingAPI is an interface for the AWS AutoScaling API client.
type AutoScalingAPI interface {
	CancelInstanceRefresh(ctx context.Context, params *autoscaling.CancelInstanceRefreshInput, optFns ...func(*autoscaling.Options)) (*autoscaling.CancelInstanceRefreshOutput, error)
	CreateAutoScalingGroup(ctx context.Context, params *autoscaling.CreateAutoScalingGroupInput, optFns ...func(*autoscaling.Options)) (*autoscaling.CreateAutoScalingGroupOutput, error)
	DeleteAutoScalingGroup(ctx context.Context, params *autoscaling.DeleteAutoScalingGroupInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DeleteAutoScalingGroupOutput, error)
	DescribeAutoScalingGroups(ctx context.Context, params *autoscaling.DescribeAutoScalingGroupsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
	UpdateAutoScalingGroup(ctx context.Context, params *autoscaling.UpdateAutoScalingGroupInput, optFns ...func(*autoscaling.Options)) (*autoscaling.UpdateAutoScalingGroupOutput, error)
	StartInstanceRefresh(ctx context.Context, params *autoscaling.StartInstanceRefreshInput, optFns ...func(*autoscaling.Options)) (*autoscaling.StartInstanceRefreshOutput, error)
	CreateOrUpdateTags(ctx context.Context, params *autoscaling.CreateOrUpdateTagsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.CreateOrUpdateTagsOutput, error)
	ResumeProcesses(ctx context.Context, params *autoscaling.ResumeProcessesInput, optFns ...func(*autoscaling.Options)) (*autoscaling.ResumeProcessesOutput, error)
	DeleteTags(ctx context.Context, params *autoscaling.DeleteTagsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DeleteTagsOutput, error)
	SuspendProcesses(ctx context.Context, params *autoscaling.SuspendProcessesInput, optFns ...func(*autoscaling.Options)) (*autoscaling.SuspendProcessesOutput, error)
	DescribeInstanceRefreshes(ctx context.Context, params *autoscaling.DescribeInstanceRefreshesInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeInstanceRefreshesOutput, error)
	DescribeLifecycleHooks(ctx context.Context, params *autoscaling.DescribeLifecycleHooksInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeLifecycleHooksOutput, error)
	PutLifecycleHook(ctx context.Context, params *autoscaling.PutLifecycleHookInput, optFns ...func(*autoscaling.Options)) (*autoscaling.PutLifecycleHookOutput, error)
	DeleteLifecycleHook(ctx context.Context, params *autoscaling.DeleteLifecycleHookInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DeleteLifecycleHookOutput, error)
}

var _ AutoScalingAPI = &autoscaling.Client{}

// NewService returns a new service given the asg api client.
func NewService(clusterScope cloud.ClusterScoper) *Service {
	return &Service{
		scope:     clusterScope,
		ASGClient: scope.NewASGClient(clusterScope, clusterScope, clusterScope, clusterScope.InfraCluster()),
		EC2Client: scope.NewEC2Client(clusterScope, clusterScope, clusterScope, clusterScope.InfraCluster()),
	}
}
