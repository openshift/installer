/*
Copyright 2019 The Kubernetes Authors.

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

package converters

import (
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
)

// MapToTags converts a map[string]*string into a infrav1.Tags.
func MapToTags(src map[string]*string) infrav1.Tags {
	if src == nil {
		return nil
	}

	tags := make(infrav1.Tags, len(src))

	for k, v := range src {
		tags[k] = ptr.Deref(v, "")
	}

	return tags
}

// TagsToMap converts infrav1.Tags into a map[string]*string.
func TagsToMap(src infrav1.Tags) map[string]*string {
	return azure.StringMapPtr(src)
}
