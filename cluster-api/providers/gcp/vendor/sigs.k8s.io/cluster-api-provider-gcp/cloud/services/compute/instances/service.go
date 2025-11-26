/*
Copyright 2021 The Kubernetes Authors.

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

package instances

import (
	"context"

	"github.com/go-logr/logr"

	k8scloud "github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/filter"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"google.golang.org/api/compute/v1"

	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
)

type instancesInterface interface {
	Get(ctx context.Context, key *meta.Key, options ...k8scloud.Option) (*compute.Instance, error)
	Insert(ctx context.Context, key *meta.Key, obj *compute.Instance, options ...k8scloud.Option) error
	Delete(ctx context.Context, key *meta.Key, options ...k8scloud.Option) error
}

type instancegroupsInterface interface {
	AddInstances(ctx context.Context, key *meta.Key, req *compute.InstanceGroupsAddInstancesRequest, options ...k8scloud.Option) error
	ListInstances(ctx context.Context, key *meta.Key, req *compute.InstanceGroupsListInstancesRequest, fl *filter.F, options ...k8scloud.Option) ([]*compute.InstanceWithNamedPorts, error)
	RemoveInstances(ctx context.Context, key *meta.Key, req *compute.InstanceGroupsRemoveInstancesRequest, options ...k8scloud.Option) error
}

// Scope is an interfaces that hold used methods.
type Scope interface {
	cloud.Machine
	InstanceSpec(log logr.Logger) *compute.Instance
}

// Service implements instances reconciler.
type Service struct {
	scope          Scope
	instances      instancesInterface
	instancegroups instancegroupsInterface
}

var _ cloud.Reconciler = &Service{}

// New returns Service from given scope.
func New(scope Scope) *Service {
	return &Service{
		scope:          scope,
		instances:      scope.Cloud().Instances(),
		instancegroups: scope.Cloud().InstanceGroups(),
	}
}
