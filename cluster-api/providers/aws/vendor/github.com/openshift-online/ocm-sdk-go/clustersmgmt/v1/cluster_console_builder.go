/*
Copyright (c) 2020 Red Hat, Inc.

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

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// ClusterConsoleBuilder contains the data and logic needed to build 'cluster_console' objects.
//
// Information about the console of a cluster.
type ClusterConsoleBuilder struct {
	bitmap_ uint32
	url     string
}

// NewClusterConsole creates a new builder of 'cluster_console' objects.
func NewClusterConsole() *ClusterConsoleBuilder {
	return &ClusterConsoleBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterConsoleBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// URL sets the value of the 'URL' attribute to the given value.
func (b *ClusterConsoleBuilder) URL(value string) *ClusterConsoleBuilder {
	b.url = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterConsoleBuilder) Copy(object *ClusterConsole) *ClusterConsoleBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.url = object.url
	return b
}

// Build creates a 'cluster_console' object using the configuration stored in the builder.
func (b *ClusterConsoleBuilder) Build() (object *ClusterConsole, err error) {
	object = new(ClusterConsole)
	object.bitmap_ = b.bitmap_
	object.url = b.url
	return
}
