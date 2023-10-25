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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

// ObjectReferenceBuilder contains the data and logic needed to build 'object_reference' objects.
//
// representation of object reference/subscription
type ObjectReferenceBuilder struct {
	bitmap_ uint32
	href    string
	id      string
	kind    string
}

// NewObjectReference creates a new builder of 'object_reference' objects.
func NewObjectReference() *ObjectReferenceBuilder {
	return &ObjectReferenceBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ObjectReferenceBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Href sets the value of the 'href' attribute to the given value.
func (b *ObjectReferenceBuilder) Href(value string) *ObjectReferenceBuilder {
	b.href = value
	b.bitmap_ |= 1
	return b
}

// Id sets the value of the 'id' attribute to the given value.
func (b *ObjectReferenceBuilder) Id(value string) *ObjectReferenceBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// Kind sets the value of the 'kind' attribute to the given value.
func (b *ObjectReferenceBuilder) Kind(value string) *ObjectReferenceBuilder {
	b.kind = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ObjectReferenceBuilder) Copy(object *ObjectReference) *ObjectReferenceBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.href = object.href
	b.id = object.id
	b.kind = object.kind
	return b
}

// Build creates a 'object_reference' object using the configuration stored in the builder.
func (b *ObjectReferenceBuilder) Build() (object *ObjectReference, err error) {
	object = new(ObjectReference)
	object.bitmap_ = b.bitmap_
	object.href = b.href
	object.id = b.id
	object.kind = b.kind
	return
}
