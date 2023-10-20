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

// AddOnNamespaceBuilder contains the data and logic needed to build 'add_on_namespace' objects.
type AddOnNamespaceBuilder struct {
	bitmap_     uint32
	id          string
	href        string
	annotations map[string]string
	labels      map[string]string
	name        string
}

// NewAddOnNamespace creates a new builder of 'add_on_namespace' objects.
func NewAddOnNamespace() *AddOnNamespaceBuilder {
	return &AddOnNamespaceBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AddOnNamespaceBuilder) Link(value bool) *AddOnNamespaceBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AddOnNamespaceBuilder) ID(value string) *AddOnNamespaceBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AddOnNamespaceBuilder) HREF(value string) *AddOnNamespaceBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnNamespaceBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Annotations sets the value of the 'annotations' attribute to the given value.
func (b *AddOnNamespaceBuilder) Annotations(value map[string]string) *AddOnNamespaceBuilder {
	b.annotations = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// Labels sets the value of the 'labels' attribute to the given value.
func (b *AddOnNamespaceBuilder) Labels(value map[string]string) *AddOnNamespaceBuilder {
	b.labels = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddOnNamespaceBuilder) Name(value string) *AddOnNamespaceBuilder {
	b.name = value
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnNamespaceBuilder) Copy(object *AddOnNamespace) *AddOnNamespaceBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if len(object.annotations) > 0 {
		b.annotations = map[string]string{}
		for k, v := range object.annotations {
			b.annotations[k] = v
		}
	} else {
		b.annotations = nil
	}
	if len(object.labels) > 0 {
		b.labels = map[string]string{}
		for k, v := range object.labels {
			b.labels[k] = v
		}
	} else {
		b.labels = nil
	}
	b.name = object.name
	return b
}

// Build creates a 'add_on_namespace' object using the configuration stored in the builder.
func (b *AddOnNamespaceBuilder) Build() (object *AddOnNamespace, err error) {
	object = new(AddOnNamespace)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.annotations != nil {
		object.annotations = make(map[string]string)
		for k, v := range b.annotations {
			object.annotations[k] = v
		}
	}
	if b.labels != nil {
		object.labels = make(map[string]string)
		for k, v := range b.labels {
			object.labels[k] = v
		}
	}
	object.name = b.name
	return
}
