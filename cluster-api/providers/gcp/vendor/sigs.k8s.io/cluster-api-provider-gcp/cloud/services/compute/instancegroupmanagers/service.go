/*
Copyright 2025 The Kubernetes Authors.

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

// Package instancegroupmanagers implements reconciliation for instanceGroupManager GCP resources.
package instancegroupmanagers

import (
	"context"

	k8scloud "github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	compute "google.golang.org/api/compute/v1"

	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
)

type instanceGroupManagersClient interface {
	Get(ctx context.Context, key *meta.Key, options ...k8scloud.Option) (*compute.InstanceGroupManager, error)
	Insert(ctx context.Context, key *meta.Key, obj *compute.InstanceGroupManager, options ...k8scloud.Option) error
	Delete(ctx context.Context, key *meta.Key, options ...k8scloud.Option) error
	Resize(context.Context, *meta.Key, int64, ...k8scloud.Option) error
	SetInstanceTemplate(context.Context, *meta.Key, *compute.InstanceGroupManagersSetInstanceTemplateRequest, ...k8scloud.Option) error
}

// Scope is an interfaces that hold used methods.
type Scope interface {
	Cloud() cloud.Cloud

	// InstanceGroupManagerResource returns the desired instanceGroupManager
	InstanceGroupManagerResource(instanceTemplateKey *meta.Key) (*compute.InstanceGroupManager, error)

	// InstanceGroupManagerResourceName returns the instanceGroupManager selfLink
	InstanceGroupManagerResourceName() (*meta.Key, error)
}

// Service implements managed instance groups reconciler.
type Service struct {
	scope                 Scope
	instanceGroupManagers instanceGroupManagersClient
	instanceGroups        k8scloud.InstanceGroups
}

// var _ cloud.Reconciler = &Service{}

// New returns Service from given scope.
func New(scope Scope) *Service {
	cloudScope := scope.Cloud()

	return &Service{
		scope:                 scope,
		instanceGroupManagers: cloudScope.InstanceGroupManagers(),
		instanceGroups:        cloudScope.InstanceGroups(),
	}
}
