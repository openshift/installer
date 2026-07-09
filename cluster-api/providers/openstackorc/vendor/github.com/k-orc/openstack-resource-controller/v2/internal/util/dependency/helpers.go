/*
Copyright 2025 The ORC Authors.

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

package dependency

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
)

// FetchDependency fetches a resource by name and checks if it's ready.
// Unlike GetDependency on DeletionGuardDependency, this doesn't add finalizers
// and is suitable for one-off lookups like resolving refs in import filters.
//
// Always returns an object (empty struct if not found/ready/error) for safe field access.
//
// Returns:
//   - The fetched object (empty struct if name is nil, not found, not ready, or on error)
//   - ReconcileStatus indicating wait state or error (nil only if name is nil or object is ready)
func FetchDependency[TP DependencyType[T], T any](
	ctx context.Context,
	k8sClient client.Client,
	namespace string,
	name *orcv1alpha1.KubernetesNameRef,
	kind string,
	isReady func(TP) bool,
) (TP, progress.ReconcileStatus) {
	var obj TP = new(T)

	if ptr.Deref(name, "") == "" {
		return obj, nil
	}

	objectKey := client.ObjectKey{Name: string(*name), Namespace: namespace}

	if err := k8sClient.Get(ctx, objectKey, obj); err != nil {
		if apierrors.IsNotFound(err) {
			return obj, progress.NewReconcileStatus().WaitingOnObject(kind, string(*name), progress.WaitingOnCreation)
		}
		return obj, progress.WrapError(fmt.Errorf("fetching %s %s: %w", kind, string(*name), err))
	}

	if !isReady(obj) {
		return obj, progress.NewReconcileStatus().WaitingOnObject(kind, string(*name), progress.WaitingOnReady)
	}

	return obj, nil
}
