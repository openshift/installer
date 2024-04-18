/*
Copyright 2023 The Kubernetes Authors.

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

package aso

import (
	"context"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Reconciler is a generic interface used to perform reconciliation of Azure resources backed by ASO.
type Reconciler[T genruntime.MetaObject] interface {
	CreateOrUpdateResource(ctx context.Context, spec azure.ASOResourceSpecGetter[T], serviceName string) (result T, err error)
	DeleteResource(ctx context.Context, resource T, serviceName string) (err error)
	PauseResource(ctx context.Context, resource T, serviceName string) (err error)
}

// TagsGetterSetter represents an object that supports tags.
type TagsGetterSetter[T genruntime.MetaObject] interface {
	GetAdditionalTags() infrav1.Tags
	GetDesiredTags(resource T) infrav1.Tags
	SetTags(resource T, tags infrav1.Tags)
}

// Patcher supplies extra patches to be applied to an ASO resource.
type Patcher interface {
	ExtraPatches() []string
}

// Scope represents the common functionality related to all scopes needed for ASO services.
type Scope interface {
	azure.AsyncStatusUpdater
	GetClient() client.Client
	ClusterName() string
	ASOOwner() client.Object
}
