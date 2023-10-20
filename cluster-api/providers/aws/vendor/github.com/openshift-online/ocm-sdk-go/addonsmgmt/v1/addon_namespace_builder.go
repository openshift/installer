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

// AddonNamespaceBuilder contains the data and logic needed to build 'addon_namespace' objects.
//
// Representation of an addon namespace object.
type AddonNamespaceBuilder struct {
	bitmap_     uint32
	annotations map[string]string
	labels      map[string]string
	name        string
	enabled     bool
}

// NewAddonNamespace creates a new builder of 'addon_namespace' objects.
func NewAddonNamespace() *AddonNamespaceBuilder {
	return &AddonNamespaceBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonNamespaceBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Annotations sets the value of the 'annotations' attribute to the given value.
func (b *AddonNamespaceBuilder) Annotations(value map[string]string) *AddonNamespaceBuilder {
	b.annotations = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddonNamespaceBuilder) Enabled(value bool) *AddonNamespaceBuilder {
	b.enabled = value
	b.bitmap_ |= 2
	return b
}

// Labels sets the value of the 'labels' attribute to the given value.
func (b *AddonNamespaceBuilder) Labels(value map[string]string) *AddonNamespaceBuilder {
	b.labels = value
	if value != nil {
		b.bitmap_ |= 4
	} else {
		b.bitmap_ &^= 4
	}
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddonNamespaceBuilder) Name(value string) *AddonNamespaceBuilder {
	b.name = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonNamespaceBuilder) Copy(object *AddonNamespace) *AddonNamespaceBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if len(object.annotations) > 0 {
		b.annotations = map[string]string{}
		for k, v := range object.annotations {
			b.annotations[k] = v
		}
	} else {
		b.annotations = nil
	}
	b.enabled = object.enabled
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

// Build creates a 'addon_namespace' object using the configuration stored in the builder.
func (b *AddonNamespaceBuilder) Build() (object *AddonNamespace, err error) {
	object = new(AddonNamespace)
	object.bitmap_ = b.bitmap_
	if b.annotations != nil {
		object.annotations = make(map[string]string)
		for k, v := range b.annotations {
			object.annotations[k] = v
		}
	}
	object.enabled = b.enabled
	if b.labels != nil {
		object.labels = make(map[string]string)
		for k, v := range b.labels {
			object.labels[k] = v
		}
	}
	object.name = b.name
	return
}
