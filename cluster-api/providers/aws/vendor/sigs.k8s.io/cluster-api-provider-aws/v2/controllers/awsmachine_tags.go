/*
Copyright 2018 The Kubernetes Authors.

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

package controllers

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	service "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
)

const (
	// TagsLastAppliedAnnotation is the key for the machine object annotation
	// which tracks the AdditionalTags in the Machine Provider Config.
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// for annotation formatting rules.
	TagsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-aws-last-applied-tags"

	// VolumeTagsLastAppliedAnnotation is the key for the ebs volumes annotation
	// which tracks the AdditionalTags in the Machine Provider Config.
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// for annotation formatting rules.
	VolumeTagsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-last-applied-tags-on-volumes"
)

// Ensure that the tags of the machine are correct
// Returns bool, error
// Bool indicates if changes were made or not, allowing the caller to decide
// if the machine should be updated.
func (r *AWSMachineReconciler) ensureTags(svc service.EC2Interface, machine *infrav1.AWSMachine, instanceID *string, additionalTags map[string]string) (bool, error) {
	annotation, err := r.machineAnnotationJSON(machine, TagsLastAppliedAnnotation)
	if err != nil {
		return false, err
	}

	// Check if the instance tags were changed. If they were, update them.
	// It would be possible here to only send new/updated tags, but for the
	// moment we send everything, even if only a single tag was created or
	// updated.
	changed, created, deleted, newAnnotation := r.tagsChanged(annotation, additionalTags)
	if changed {
		err = svc.UpdateResourceTags(instanceID, created, deleted)
		if err != nil {
			return false, err
		}

		// We also need to update the annotation if anything changed.
		err = r.updateMachineAnnotationJSON(machine, TagsLastAppliedAnnotation, newAnnotation)
		if err != nil {
			return false, err
		}
	}

	return changed, nil
}

// Ensure that the tags of the volumes in the machine are correct
// Returns tags which are being created/updated/deleted and error.
func (r *AWSMachineReconciler) ensureVolumeTags(svc service.EC2Interface, volumeID *string, annotation map[string]interface{}, additionalTags map[string]string) (map[string]interface{}, error) {
	// Check if the volume tags were changed. If they were, update them.
	// It would be possible here to only send new/updated tags, but for the
	// moment we send everything, even if only a single tag was created or
	// updated.
	changed, created, deleted, subAnnotation := r.tagsChanged(annotation, additionalTags)
	if changed {
		err := svc.UpdateResourceTags(volumeID, created, deleted)
		if err != nil {
			return nil, err
		}
	}

	return subAnnotation, nil
}

// tagsChanged determines which tags to delete and which to add.
func (r *AWSMachineReconciler) tagsChanged(annotation map[string]interface{}, src map[string]string) (bool, map[string]string, map[string]string, map[string]interface{}) {
	// Bool tracking if we found any changed state.
	changed := false

	// Tracking for created/updated
	created := map[string]string{}

	// Tracking for tags that were deleted.
	deleted := map[string]string{}

	// The new annotation that we need to set if anything is created/updated.
	newAnnotation := map[string]interface{}{}

	// Loop over annotation, checking if entries are in src.
	// If an entry is present in annotation but not src, it has been deleted
	// since last time. We flag this in the deleted map.
	for t, v := range annotation {
		_, ok := src[t]

		// Entry isn't in src, it has been deleted.
		if !ok {
			// Cast v to a string here. This should be fine, tags are always
			// strings.
			deleted[t] = v.(string)
			changed = true
		}
	}

	// Loop over src, checking for entries in annotation.
	//
	// If an entry is in src, but not annotation, it has been created since
	// last time.
	//
	// If an entry is in both src and annotation, we compare their values, if
	// the value in src differs from that in annotation, the tag has been
	// updated since last time.
	for t, v := range src {
		av, ok := annotation[t]

		// Entries in the src always need to be noted in the newAnnotation. We
		// know they're going to be created or updated.
		newAnnotation[t] = v

		// Entry isn't in annotation, it's new.
		if !ok {
			created[t] = v
			newAnnotation[t] = v
			changed = true
			continue
		}

		// Entry is in annotation, has the value changed?
		if v != av {
			created[t] = v
			changed = true
		}

		// Entry existed in both src and annotation, and their values were
		// equal. Nothing to do.
	}

	// We made it through the loop, and everything that was in src, was also
	// in dst. Nothing changed.
	return changed, created, deleted, newAnnotation
}
