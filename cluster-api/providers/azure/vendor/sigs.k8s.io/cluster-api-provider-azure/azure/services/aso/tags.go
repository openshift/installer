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
	"encoding/json"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/tags"
	"sigs.k8s.io/cluster-api-provider-azure/util/maps"
)

// tagsLastAppliedAnnotation is the key for the annotation which tracks the AdditionalTags.
// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
// for annotation formatting rules.
const tagsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-azure-last-applied-tags"

// reconcileTags modifies parameters in place to update its tags and its last-applied annotation.
func reconcileTags[T genruntime.MetaObject](t TagsGetterSetter[T], existing T, resourceExists bool, parameters T) error {
	var existingTags infrav1.Tags
	lastAppliedTags := map[string]interface{}{}

	if resourceExists {
		lastAppliedTagsJSON := existing.GetAnnotations()[tagsLastAppliedAnnotation]
		if lastAppliedTagsJSON != "" {
			err := json.Unmarshal([]byte(lastAppliedTagsJSON), &lastAppliedTags)
			if err != nil {
				return errors.Wrapf(err, "failed to unmarshal JSON from %s annotation", tagsLastAppliedAnnotation)
			}
		}

		existingTags = t.GetDesiredTags(existing)
	}

	existingTagsMap := converters.TagsToMap(existingTags)
	_, createdOrUpdated, deleted, newAnnotation := tags.TagsChanged(lastAppliedTags, t.GetAdditionalTags(), existingTagsMap)
	newTags := maps.Merge(maps.Merge(existingTags, t.GetDesiredTags(parameters)), createdOrUpdated)
	for k := range deleted {
		delete(newTags, k)
	}
	if len(newTags) == 0 {
		newTags = nil
	}
	t.SetTags(parameters, newTags)

	// We also need to update the annotation even if nothing changed to
	// ensure it's set immediately following resource creation.
	newAnnotationJSON, err := json.Marshal(newAnnotation)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal JSON to %s annotation", tagsLastAppliedAnnotation)
	}
	parameters.SetAnnotations(maps.Merge(parameters.GetAnnotations(), map[string]string{
		tagsLastAppliedAnnotation: string(newAnnotationJSON),
	}))

	return nil
}
