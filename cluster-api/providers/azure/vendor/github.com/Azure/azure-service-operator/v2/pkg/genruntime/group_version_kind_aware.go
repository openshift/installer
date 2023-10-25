/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package genruntime

import (
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupVersionKindAware is implemented by resources that are aware of which version of the resource was originally
// specified. This allows us to interface with ARM using an API version specified by an end user.
type GroupVersionKindAware interface {
	// OriginalGVK returns the GroupVersionKind originally used to create the resource (regardless of any conversions)
	OriginalGVK() *schema.GroupVersionKind
}

// GetOriginalGVK gets the GVK the original GVK the object was created with.
func GetOriginalGVK(obj ARMMetaObject) schema.GroupVersionKind {
	// If our current resource is aware of its original GVK, use that for our result
	aware, ok := obj.(GroupVersionKindAware)
	if ok {
		result := *aware.OriginalGVK()
		result.Version = transformAlphaVersionToStable(result.Version)
		return result
	}

	// The GVK of our current object
	result := obj.GetObjectKind().GroupVersionKind()
	result.Version = transformAlphaVersionToStable(result.Version)
	return result
}

// TODO: Delete this function and its usages once we've reached v2.3.0 (or some version sufficiently far from v2.0.0)
// transformAlphaVersionToStable ensures that we don't attempt to use a resource version that has been removed.
func transformAlphaVersionToStable(version string) string {
	return strings.Replace(version, "v1alpha1api", "v1api", 1)
}
