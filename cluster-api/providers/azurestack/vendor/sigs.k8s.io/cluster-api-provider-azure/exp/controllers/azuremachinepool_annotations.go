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
)

// AnnotationJSON returns a map[string]interface from a JSON annotation.
// This method gets the given `annotation` from an `annotationReaderWriter` and unmarshalls it
// from a JSON string into a `map[string]interface{}`.
func (ampr *AzureMachinePoolReconciler) AnnotationJSON(rw annotationReaderWriter, annotation string) (map[string]interface{}, error) {
	out := map[string]interface{}{}

	jsonAnnotation := ampr.Annotation(rw, annotation)
	if jsonAnnotation == "" {
		return out, nil
	}

	err := json.Unmarshal([]byte(jsonAnnotation), &out)
	if err != nil {
		return out, err
	}

	return out, nil
}

// Annotation fetches the specific machine annotation.
func (ampr *AzureMachinePoolReconciler) Annotation(rw annotationReaderWriter, annotation string) string {
	return rw.GetAnnotations()[annotation]
}
