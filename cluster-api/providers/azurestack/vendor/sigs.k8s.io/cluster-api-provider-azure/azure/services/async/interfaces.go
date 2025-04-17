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

package async

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
)

// FutureScope stores and retrieves Futures and Conditions.
type FutureScope interface {
	azure.AsyncStatusUpdater
}

// Getter gets a resource.
type Getter interface {
	Get(ctx context.Context, spec azure.ResourceSpecGetter) (result interface{}, err error)
}

// TagsGetter is an interface that can get a tags resource.
type TagsGetter interface {
	GetAtScope(ctx context.Context, scope string) (result armresources.TagsResource, err error)
}

// Creator creates or updates a resource asynchronously.
type Creator[T any] interface {
	Getter
	CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[T], err error)
}

// Deleter deletes a resource asynchronously.
type Deleter[T any] interface {
	DeleteAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string) (poller *runtime.Poller[T], err error)
}

// Reconciler reconciles a resource.
type Reconciler interface {
	CreateOrUpdateResource(ctx context.Context, spec azure.ResourceSpecGetter, serviceName string) (result interface{}, err error)
	DeleteResource(ctx context.Context, spec azure.ResourceSpecGetter, serviceName string) (err error)
}
