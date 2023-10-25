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

// FlavourNodesBuilder contains the data and logic needed to build 'flavour_nodes' objects.
//
// Counts of different classes of nodes inside a flavour.
type FlavourNodesBuilder struct {
	bitmap_ uint32
	master  int
}

// NewFlavourNodes creates a new builder of 'flavour_nodes' objects.
func NewFlavourNodes() *FlavourNodesBuilder {
	return &FlavourNodesBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *FlavourNodesBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Master sets the value of the 'master' attribute to the given value.
func (b *FlavourNodesBuilder) Master(value int) *FlavourNodesBuilder {
	b.master = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *FlavourNodesBuilder) Copy(object *FlavourNodes) *FlavourNodesBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.master = object.master
	return b
}

// Build creates a 'flavour_nodes' object using the configuration stored in the builder.
func (b *FlavourNodesBuilder) Build() (object *FlavourNodes, err error) {
	object = new(FlavourNodes)
	object.bitmap_ = b.bitmap_
	object.master = b.master
	return
}
