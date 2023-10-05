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

package controllers

import (
	"encoding/json"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
)

const (
	// TagsLastAppliedAnnotation is the key for the AWSMachinePool object annotation
	// which tracks the tags that the AWSMachinePool actuator is responsible
	// for. These are the tags that have been handled by the
	// AdditionalTags in the AWSMachinePool Provider Config.
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// for annotation formatting rules.
	TagsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-aws-last-applied-tags"
)

// Ensure that the tags of the AWSMachinePool are correct
// Returns bool, error
// Bool indicates if changes were made or not, allowing the caller to decide
// if the machine should be updated.
func (r *AWSMachinePoolReconciler) ensureTags(ec2svc services.EC2Interface, asgsvc services.ASGInterface, machinePool *expinfrav1.AWSMachinePool, launchTemplateID, asgName *string, additionalTags map[string]string) (bool, error) {
	annotation, err := r.machinePoolAnnotationJSON(machinePool, TagsLastAppliedAnnotation)
	if err != nil {
		return false, err
	}

	// Check if the instance tags were changed. If they were, update them.
	// It would be possible here to only send new/updated tags, but for the
	// moment we send everything, even if only a single tag was created or
	// upated.
	changed, created, deleted, newAnnotation := tagsChanged(annotation, additionalTags)
	if changed {
		err = ec2svc.UpdateResourceTags(launchTemplateID, created, deleted)
		if err != nil {
			return false, err
		}

		if err := asgsvc.UpdateResourceTags(asgName, created, deleted); err != nil {
			return false, err
		}

		// We also need to update the annotation if anything changed.
		err = r.updateMachinePoolAnnotationJSON(machinePool, TagsLastAppliedAnnotation, newAnnotation)
		if err != nil {
			return false, err
		}
	}

	return changed, nil
}

// tagsChanged determines which tags to delete and which to add.
func tagsChanged(annotation map[string]interface{}, src map[string]string) (bool, map[string]string, map[string]string, map[string]interface{}) {
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

// updateMachinePoolAnnotationJSON updates the `annotation` on `machinePool` with
// `content`. `content` in this case should be a `map[string]interface{}`
// suitable for turning into JSON. This `content` map will be marshalled into a
// JSON string before being set as the given `annotation`.
func (r *AWSMachinePoolReconciler) updateMachinePoolAnnotationJSON(machinePool *expinfrav1.AWSMachinePool, annotation string, content map[string]interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}

	r.updateMachinePoolAnnotation(machinePool, annotation, string(b))
	return nil
}

// updateMachinePoolAnnotation updates the `annotation` on the given `machinePool` with
// `content`.
func (r *AWSMachinePoolReconciler) updateMachinePoolAnnotation(machinePool *expinfrav1.AWSMachinePool, annotation, content string) {
	// Get the annotations
	annotations := machinePool.GetAnnotations()

	if annotations == nil {
		annotations = make(map[string]string)
	}

	// Set our annotation to the given content.
	annotations[annotation] = content

	// Update the machine object with these annotations
	machinePool.SetAnnotations(annotations)
}

// Returns a map[string]interface from a JSON annotation.
// This method gets the given `annotation` from the `machinePool` and unmarshalls it
// from a JSON string into a `map[string]interface{}`.
func (r *AWSMachinePoolReconciler) machinePoolAnnotationJSON(machinePool *expinfrav1.AWSMachinePool, annotation string) (map[string]interface{}, error) {
	out := map[string]interface{}{}

	jsonAnnotation := r.machinePoolAnnotation(machinePool, annotation)
	if len(jsonAnnotation) == 0 {
		return out, nil
	}

	err := json.Unmarshal([]byte(jsonAnnotation), &out)
	if err != nil {
		return out, err
	}

	return out, nil
}

// Fetches the specific machine annotation.
func (r *AWSMachinePoolReconciler) machinePoolAnnotation(machinePool *expinfrav1.AWSMachinePool, annotation string) string {
	return machinePool.GetAnnotations()[annotation]
}
