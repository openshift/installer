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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

// ImageOverrides holds the lists of available images per cloud provider.
type ImageOverridesBuilder struct {
	fieldSet_ []bool
	id        string
	href      string
	aws       []*AMIOverrideBuilder
	gcp       []*GCPImageOverrideBuilder
}

// NewImageOverrides creates a new builder of 'image_overrides' objects.
func NewImageOverrides() *ImageOverridesBuilder {
	return &ImageOverridesBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ImageOverridesBuilder) Link(value bool) *ImageOverridesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ImageOverridesBuilder) ID(value string) *ImageOverridesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ImageOverridesBuilder) HREF(value string) *ImageOverridesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ImageOverridesBuilder) Empty() bool {
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

// AWS sets the value of the 'AWS' attribute to the given values.
func (b *ImageOverridesBuilder) AWS(values ...*AMIOverrideBuilder) *ImageOverridesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.aws = make([]*AMIOverrideBuilder, len(values))
	copy(b.aws, values)
	b.fieldSet_[3] = true
	return b
}

// GCP sets the value of the 'GCP' attribute to the given values.
func (b *ImageOverridesBuilder) GCP(values ...*GCPImageOverrideBuilder) *ImageOverridesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.gcp = make([]*GCPImageOverrideBuilder, len(values))
	copy(b.gcp, values)
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ImageOverridesBuilder) Copy(object *ImageOverrides) *ImageOverridesBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.aws != nil {
		b.aws = make([]*AMIOverrideBuilder, len(object.aws))
		for i, v := range object.aws {
			b.aws[i] = NewAMIOverride().Copy(v)
		}
	} else {
		b.aws = nil
	}
	if object.gcp != nil {
		b.gcp = make([]*GCPImageOverrideBuilder, len(object.gcp))
		for i, v := range object.gcp {
			b.gcp[i] = NewGCPImageOverride().Copy(v)
		}
	} else {
		b.gcp = nil
	}
	return b
}

// Build creates a 'image_overrides' object using the configuration stored in the builder.
func (b *ImageOverridesBuilder) Build() (object *ImageOverrides, err error) {
	object = new(ImageOverrides)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.aws != nil {
		object.aws = make([]*AMIOverride, len(b.aws))
		for i, v := range b.aws {
			object.aws[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.gcp != nil {
		object.gcp = make([]*GCPImageOverride, len(b.gcp))
		for i, v := range b.gcp {
			object.gcp[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
