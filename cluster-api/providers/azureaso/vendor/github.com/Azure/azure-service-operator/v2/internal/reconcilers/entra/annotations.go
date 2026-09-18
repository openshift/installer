/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package entra

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// getEntraID retrieves the Entra ID from the annotations of the object.
// Returns the ID and true if found, empty string and false if not.
func getEntraID(obj v1.Object) (string, bool) {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		return "", false
	}

	result, ok := annotations[genruntime.ResourceIDAnnotation]
	return result, ok
}

// setEntraID stored the Entra ID in the annotations of the object.
func setEntraID(obj v1.Object, id string) {
	genruntime.AddAnnotation(obj, genruntime.ResourceIDAnnotation, id)
}
