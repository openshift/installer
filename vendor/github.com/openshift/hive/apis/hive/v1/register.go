// NOTE: Boilerplate only.  Ignore this file.

// Package v1 contains API Schema definitions for the hive v1 API group
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/openshift/hive/apis/hive
// +k8s:defaulter-gen=TypeMeta
// +groupName=hive.openshift.io
package v1

import (
	"github.com/openshift/hive/apis/scheme"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// HiveAPIGroup is the group that all hive objects belong to in the API server.
	HiveAPIGroup = "hive.openshift.io"

	// HiveAPIVersion is the api version that all hive objects are currently at.
	HiveAPIVersion = "v1"

	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: HiveAPIGroup, Version: HiveAPIVersion}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}

	// AddToScheme is a shortcut for SchemeBuilder.AddToScheme
	AddToScheme = SchemeBuilder.AddToScheme
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
