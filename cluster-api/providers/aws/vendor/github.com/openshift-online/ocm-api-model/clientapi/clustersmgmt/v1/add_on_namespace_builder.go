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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

type AddOnNamespaceBuilder struct {
	fieldSet_   []bool
	id          string
	href        string
	annotations map[string]string
	labels      map[string]string
	name        string
}

// NewAddOnNamespace creates a new builder of 'add_on_namespace' objects.
func NewAddOnNamespace() *AddOnNamespaceBuilder {
	return &AddOnNamespaceBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AddOnNamespaceBuilder) Link(value bool) *AddOnNamespaceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AddOnNamespaceBuilder) ID(value string) *AddOnNamespaceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AddOnNamespaceBuilder) HREF(value string) *AddOnNamespaceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnNamespaceBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// Annotations sets the value of the 'annotations' attribute to the given value.
func (b *AddOnNamespaceBuilder) Annotations(value map[string]string) *AddOnNamespaceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.annotations = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// Labels sets the value of the 'labels' attribute to the given value.
func (b *AddOnNamespaceBuilder) Labels(value map[string]string) *AddOnNamespaceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.labels = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddOnNamespaceBuilder) Name(value string) *AddOnNamespaceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.name = value
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnNamespaceBuilder) Copy(object *AddOnNamespace) *AddOnNamespaceBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
