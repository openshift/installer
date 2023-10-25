/*
Copyright 2020 The Kubernetes Authors.

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

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "tags"

// TagScope defines the scope interface for a tags service.
type TagScope interface {
	azure.Authorizer
	ClusterName() string
	TagsSpecs() []azure.TagsSpec
	AnnotationJSON(string) (map[string]interface{}, error)
	UpdateAnnotationJSON(string, map[string]interface{}) error
}

// Service provides operations on Azure resources.
type Service struct {
	Scope TagScope
	client
}

// New creates a new service.
func New(scope TagScope) (*Service, error) {
	cli, err := NewClient(scope)
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope:  scope,
		client: cli,
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Some resource types are always assumed to be managed by CAPZ whether or not
// they have the canonical "owned" tag applied to most resources. The annotation
// key for those types should be listed here so their tags are always
// interpreted as managed.
var alwaysManagedAnnotations = map[string]struct{}{
	azure.ManagedClusterTagsLastAppliedAnnotation: {},
}

// Reconcile ensures tags are correct.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "tags.Service.Reconcile")
	defer done()

	for _, tagsSpec := range s.Scope.TagsSpecs() {
		existingTags, err := s.client.GetAtScope(ctx, tagsSpec.Scope)
		if err != nil {
			return errors.Wrap(err, "failed to get existing tags")
		}
		tags := make(map[string]*string)
		if existingTags.Properties != nil && existingTags.Properties.Tags != nil {
			tags = existingTags.Properties.Tags
		}

		if _, alwaysManaged := alwaysManagedAnnotations[tagsSpec.Annotation]; !alwaysManaged && !s.isResourceManaged(tags) {
			log.V(4).Info("Skipping tags reconcile for not managed resource")
			continue
		}

		lastAppliedTags, err := s.Scope.AnnotationJSON(tagsSpec.Annotation)
		if err != nil {
			return err
		}
		changed, createdOrUpdated, deleted, newAnnotation := TagsChanged(lastAppliedTags, tagsSpec.Tags, tags)
		if changed {
			log.V(2).Info("Updating tags")
			if len(createdOrUpdated) > 0 {
				createdOrUpdatedTags := make(map[string]*string)
				for k, v := range createdOrUpdated {
					createdOrUpdatedTags[k] = ptr.To(v)
				}

				if _, err := s.client.UpdateAtScope(ctx, tagsSpec.Scope, armresources.TagsPatchResource{Operation: ptr.To(armresources.TagsPatchOperationMerge), Properties: &armresources.Tags{Tags: createdOrUpdatedTags}}); err != nil {
					return errors.Wrap(err, "cannot update tags")
				}
			}

			if len(deleted) > 0 {
				deletedTags := make(map[string]*string)
				for k, v := range deleted {
					deletedTags[k] = ptr.To(v)
				}

				if _, err := s.client.UpdateAtScope(ctx, tagsSpec.Scope, armresources.TagsPatchResource{Operation: ptr.To(armresources.TagsPatchOperationDelete), Properties: &armresources.Tags{Tags: deletedTags}}); err != nil {
					return errors.Wrap(err, "cannot update tags")
				}
			}
			log.V(2).Info("successfully updated tags")
		}

		// We also need to update the annotation even if nothing changed to
		// ensure it's set immediately following resource creation.
		if err := s.Scope.UpdateAnnotationJSON(tagsSpec.Annotation, newAnnotation); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) isResourceManaged(tags map[string]*string) bool {
	return converters.MapToTags(tags).HasOwned(s.Scope.ClusterName())
}

// Delete is a no-op as the tags get deleted as part of VM deletion.
func (s *Service) Delete(ctx context.Context) error {
	_, _, done := tele.StartSpanWithLogger(ctx, "tags.Service.Delete")
	defer done()

	return nil
}

// TagsChanged determines which tags to delete and which to add.
func TagsChanged(lastAppliedTags map[string]interface{}, desiredTags map[string]string, currentTags map[string]*string) (change bool, createOrUpdates map[string]string, deletes map[string]string, annotation map[string]interface{}) {
	// Bool tracking if we found any changed state.
	changed := false

	// Tracking for created/updated
	createdOrUpdated := map[string]string{}

	// Tracking for tags that were deleted.
	deleted := map[string]string{}

	// The new annotation that we need to set if anything is created/updated.
	newAnnotation := map[string]interface{}{}

	// Loop over lastAppliedTags, checking if entries are in desiredTags.
	// If an entry is present in lastAppliedTags but not in desiredTags, it has been deleted
	// since last time. We flag this in the deleted map.
	for t, v := range lastAppliedTags {
		_, ok := desiredTags[t]

		// Entry isn't in desiredTags, it has been deleted.
		if !ok {
			// Cast v to a string here. This should be fine, tags are always
			// strings.
			deleted[t] = v.(string)
			changed = true
		}
	}

	// Loop over desiredTags, checking for entries in currentTags.
	//
	// If an entry is in desiredTags, but not currentTags, it has been created since
	// last time, or some external entity deleted it.
	//
	// If an entry is in both desiredTags and currentTags, we compare their values, if
	// the value in desiredTags differs from that in currentTags, the tag has been
	// updated since last time or some external entity modified it.
	for t, v := range desiredTags {
		av, ok := currentTags[t]

		// Entries in the desiredTags always need to be noted in the newAnnotation. We
		// know they're going to be created or updated.
		newAnnotation[t] = v

		// Entry isn't in desiredTags, it's new.
		if !ok {
			createdOrUpdated[t] = v
			newAnnotation[t] = v
			changed = true
			continue
		}

		// Entry is in desiredTags, has the value changed?
		if v != *av {
			createdOrUpdated[t] = v
			changed = true
		}

		// Entry existed in both desiredTags and desiredTags, and their values were
		// equal. Nothing to do.
	}

	// We made it through the loop, and everything that was in desiredTags, was also
	// in dst. Nothing changed.
	return changed, createdOrUpdated, deleted, newAnnotation
}

// IsManaged returns always returns true as CAPZ does not support BYO tags.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	return true, nil
}
