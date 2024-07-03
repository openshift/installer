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

package networks

import (
	"context"

	k8scloud "github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"google.golang.org/api/compute/v1"

	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
)

type networksInterface interface {
	Get(ctx context.Context, key *meta.Key, options ...k8scloud.Option) (*compute.Network, error)
	Insert(ctx context.Context, key *meta.Key, obj *compute.Network, options ...k8scloud.Option) error
	Delete(ctx context.Context, key *meta.Key, options ...k8scloud.Option) error
}

type routersInterface interface {
	Get(ctx context.Context, key *meta.Key, options ...k8scloud.Option) (*compute.Router, error)
	Insert(ctx context.Context, key *meta.Key, obj *compute.Router, options ...k8scloud.Option) error
	Delete(ctx context.Context, key *meta.Key, options ...k8scloud.Option) error
}

// Scope is an interfaces that hold used methods.
type Scope interface {
	cloud.Cluster
	NetworkSpec() *compute.Network
	NatRouterSpec() *compute.Router
}

// Service implements networks reconciler.
type Service struct {
	scope    Scope
	networks networksInterface
	routers  routersInterface
}

var _ cloud.Reconciler = &Service{}

// New returns Service from given scope.
func New(scope Scope) *Service {
	scopeCloud := scope.Cloud()
	if scope.IsSharedVpc() {
		scopeCloud = scope.NetworkCloud()
	}

	return &Service{
		scope:    scope,
		networks: scopeCloud.Networks(),
		routers:  scopeCloud.Routers(),
	}
}
