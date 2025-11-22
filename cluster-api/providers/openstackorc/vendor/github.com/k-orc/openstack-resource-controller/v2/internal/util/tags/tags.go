/*
Copyright 2024 The ORC Authors.

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

package tags

import (
	"context"
	"strings"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/tags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	"k8s.io/utils/set"
)

// StringTag is a constraint that includes all custom tag types defined in the API.
// This allows a single generic function to work with any of them.
type StringTag interface {
	orcv1alpha1.NeutronTag | orcv1alpha1.ServerTag | orcv1alpha1.ImageTag | orcv1alpha1.KeystoneTag
}

// Join joins a slice of tags into a comma-separated list of tags.
func Join[T StringTag](tags []T) string {
	var b strings.Builder
	for i := range tags {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(string(tags[i]))
	}
	return b.String()
}

func ReconcileTags[orcObjectPT, osResourceT any, T StringTag](
	specTags []T,
	observedTags []string,
	tagReplacer TagReplacer,
) interfaces.ResourceReconciler[orcObjectPT, osResourceT] {
	return func(ctx context.Context, _ orcObjectPT, _ *osResourceT) progress.ReconcileStatus {
		observedTagSet := set.New(observedTags...)
		specTagSet := set.New[string]()
		for i := range specTags {
			specTagSet.Insert(string(specTags[i]))
		}

		// Tags are in sync, nothing to do.
		if specTagSet.Equal(observedTagSet) {
			return nil
		}

		// Tags are out of sync, call the API to replace them.
		err := tagReplacer(ctx, specTagSet.SortedList())
		if err != nil {
			return progress.WrapError(err)
		}

		// If we updated the tags, we need another reconcile to refresh the resource status.
		return progress.NeedsRefresh()
	}
}

type TagReplacer func(ctx context.Context, tags []string) error

// NewNeutronTagReplacer returns a TagReplacer function for Neutron resources.
func NewNeutronTagReplacer(networkClient osclients.NetworkClient, resourceType, resourceID string) TagReplacer {
	return func(ctx context.Context, tagsToSet []string) error {
		opts := attributestags.ReplaceAllOpts{Tags: tagsToSet}
		_, err := networkClient.ReplaceAllAttributesTags(ctx, resourceType, resourceID, &opts)
		return err
	}
}

// NewServerTagReplacer returns a TagReplacer function for Nova Server resources.
func NewServerTagReplacer(computeClient osclients.ComputeClient, resourceID string) TagReplacer {
	return func(ctx context.Context, tagsToSet []string) error {
		opts := tags.ReplaceAllOpts{Tags: tagsToSet}
		_, err := computeClient.ReplaceAllServerAttributesTags(ctx, resourceID, &opts)
		return err
	}
}
