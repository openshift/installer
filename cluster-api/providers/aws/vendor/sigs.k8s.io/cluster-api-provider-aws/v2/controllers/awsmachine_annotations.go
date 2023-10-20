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
	"encoding/json"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

// updateMachineAnnotationJSON updates the `annotation` on `machine` with
// `content`. `content` in this case should be a `map[string]interface{}`
// suitable for turning into JSON. This `content` map will be marshalled into a
// JSON string before being set as the given `annotation`.
func (r *AWSMachineReconciler) updateMachineAnnotationJSON(machine *infrav1.AWSMachine, annotation string, content map[string]interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}

	r.updateMachineAnnotation(machine, annotation, string(b))
	return nil
}

// updateMachineAnnotation updates the `annotation` on the given `machine` with
// `content`.
func (r *AWSMachineReconciler) updateMachineAnnotation(machine *infrav1.AWSMachine, annotation, content string) {
	// Get the annotations
	annotations := machine.GetAnnotations()

	// Set our annotation to the given content.
	if annotations != nil {
		annotations[annotation] = content
	}

	// Update the machine object with these annotations
	machine.SetAnnotations(annotations)
}

// Returns a map[string]interface from a JSON annotation.
// This method gets the given `annotation` from the `machine` and unmarshalls it
// from a JSON string into a `map[string]interface{}`.
func (r *AWSMachineReconciler) machineAnnotationJSON(machine *infrav1.AWSMachine, annotation string) (map[string]interface{}, error) {
	out := map[string]interface{}{}

	jsonAnnotation := r.machineAnnotation(machine, annotation)
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
func (r *AWSMachineReconciler) machineAnnotation(machine *infrav1.AWSMachine, annotation string) string {
	return machine.GetAnnotations()[annotation]
}
