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

// Package instancetemplates implements reconciliation for instanceTemplate GCP resources.
package instancetemplates

import (
	"context"

	k8scloud "github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/filter"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	compute "google.golang.org/api/compute/v1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
)

type instancetemplatesInterface interface {
	Get(ctx context.Context, key *meta.Key, options ...k8scloud.Option) (*compute.InstanceTemplate, error)
	List(ctx context.Context, fl *filter.F, options ...k8scloud.Option) ([]*compute.InstanceTemplate, error)
	Insert(ctx context.Context, key *meta.Key, obj *compute.InstanceTemplate, options ...k8scloud.Option) error
	Delete(ctx context.Context, key *meta.Key, options ...k8scloud.Option) error
}

// Scope is an interfaces that hold used methods.
type Scope interface {
	Cloud() cloud.Cloud

	// ClusterName returns the name of the cluster.
	ClusterName() string

	// InstanceTemplateResource returns the desired instanceTemplate
	InstanceTemplateResource(ctx context.Context) (*compute.InstanceTemplate, error)

	// BaseInstanceTemplateResourceName returns the base instanceTemplate selfLink
	BaseInstanceTemplateResourceName() (*meta.Key, error)
}

// Service implements managed instance groups reconciler.
type Service struct {
	scope             Scope
	instanceTemplates instancetemplatesInterface
}

// var _ cloud.Reconciler = &Service{}

// New returns Service from given scope.
func New(scope Scope) *Service {
	cloudScope := scope.Cloud()

	return &Service{
		scope:             scope,
		instanceTemplates: cloudScope.InstanceTemplates(),
	}
}
