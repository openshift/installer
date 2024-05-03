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

package loadbalancers

import (
	"context"

	k8scloud "github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/filter"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"google.golang.org/api/compute/v1"

	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
)

type addressesInterface interface {
	Get(ctx context.Context, key *meta.Key, options ...k8scloud.Option) (*compute.Address, error)
	Insert(ctx context.Context, key *meta.Key, obj *compute.Address, options ...k8scloud.Option) error
	Delete(ctx context.Context, key *meta.Key, options ...k8scloud.Option) error
}

type backendservicesInterface interface {
	Get(ctx context.Context, key *meta.Key, options ...k8scloud.Option) (*compute.BackendService, error)
	Insert(ctx context.Context, key *meta.Key, obj *compute.BackendService, options ...k8scloud.Option) error
	Update(ctx context.Context, key *meta.Key, obj *compute.BackendService, options ...k8scloud.Option) error
	Delete(ctx context.Context, key *meta.Key, options ...k8scloud.Option) error
}

type forwardingrulesInterface interface {
	Get(ctx context.Context, key *meta.Key, options ...k8scloud.Option) (*compute.ForwardingRule, error)
	Insert(ctx context.Context, key *meta.Key, obj *compute.ForwardingRule, options ...k8scloud.Option) error
	Delete(ctx context.Context, key *meta.Key, options ...k8scloud.Option) error
}

type healthchecksInterface interface {
	Get(ctx context.Context, key *meta.Key, options ...k8scloud.Option) (*compute.HealthCheck, error)
	Insert(ctx context.Context, key *meta.Key, obj *compute.HealthCheck, options ...k8scloud.Option) error
	Delete(ctx context.Context, key *meta.Key, options ...k8scloud.Option) error
}

type instancegroupsInterface interface {
	Get(ctx context.Context, key *meta.Key, options ...k8scloud.Option) (*compute.InstanceGroup, error)
	List(ctx context.Context, zone string, fl *filter.F, options ...k8scloud.Option) ([]*compute.InstanceGroup, error)
	Insert(ctx context.Context, key *meta.Key, obj *compute.InstanceGroup, options ...k8scloud.Option) error
	Delete(ctx context.Context, key *meta.Key, options ...k8scloud.Option) error
}

type targettcpproxiesInterface interface {
	Get(ctx context.Context, key *meta.Key, options ...k8scloud.Option) (*compute.TargetTcpProxy, error)
	Insert(ctx context.Context, key *meta.Key, obj *compute.TargetTcpProxy, options ...k8scloud.Option) error
	Delete(ctx context.Context, key *meta.Key, options ...k8scloud.Option) error
}

// Scope is an interfaces that hold used methods.
type Scope interface {
	cloud.Cluster
	AddressSpec() *compute.Address
	BackendServiceSpec() *compute.BackendService
	ForwardingRuleSpec() *compute.ForwardingRule
	HealthCheckSpec() *compute.HealthCheck
	InstanceGroupSpec(zone string) *compute.InstanceGroup
	TargetTCPProxySpec() *compute.TargetTcpProxy
}

// Service implements loadbalancers reconciler.
type Service struct {
	scope            Scope
	addresses        addressesInterface
	backendservices  backendservicesInterface
	forwardingrules  forwardingrulesInterface
	healthchecks     healthchecksInterface
	instancegroups   instancegroupsInterface
	targettcpproxies targettcpproxiesInterface
}

var _ cloud.Reconciler = &Service{}

// New returns Service from given scope.
func New(scope Scope) *Service {
	return &Service{
		scope:            scope,
		addresses:        scope.Cloud().GlobalAddresses(),
		backendservices:  scope.Cloud().BackendServices(),
		forwardingrules:  scope.Cloud().GlobalForwardingRules(),
		healthchecks:     scope.Cloud().HealthChecks(),
		instancegroups:   scope.Cloud().InstanceGroups(),
		targettcpproxies: scope.Cloud().TargetTcpProxies(),
	}
}
