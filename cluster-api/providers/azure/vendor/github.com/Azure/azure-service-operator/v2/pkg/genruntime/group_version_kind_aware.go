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
		result.Version = transformDeprecatedVersionToStable(result.Version)
		return result
	}

	// The GVK of our current object
	result := obj.GetObjectKind().GroupVersionKind()
	result.Version = transformDeprecatedVersionToStable(result.Version)
	return result
}

// TODO: Delete this function and its usages once we've reached a release sufficiently far from 2.4.0 (when beta versions were deprecated)
// transformDeprecatedVersionToStable ensures that we don't attempt to use a resource version that has been removed.
func transformDeprecatedVersionToStable(version string) string {
	result := version
	result = strings.Replace(result, "v1alpha1api", "v1api", 1)
	if version == "v1beta1" { // For handcrafted resources
		result = "v1"
	}
	result = strings.Replace(result, "v1beta", "v1api", 1)

	return result
}
