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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// ImageOverridesBuilder contains the data and logic needed to build 'image_overrides' objects.
//
// ImageOverrides holds the lists of available images per cloud provider.
type ImageOverridesBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	aws     []*AMIOverrideBuilder
	gcp     []*GCPImageOverrideBuilder
}

// NewImageOverrides creates a new builder of 'image_overrides' objects.
func NewImageOverrides() *ImageOverridesBuilder {
	return &ImageOverridesBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ImageOverridesBuilder) Link(value bool) *ImageOverridesBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ImageOverridesBuilder) ID(value string) *ImageOverridesBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ImageOverridesBuilder) HREF(value string) *ImageOverridesBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ImageOverridesBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// AWS sets the value of the 'AWS' attribute to the given values.
func (b *ImageOverridesBuilder) AWS(values ...*AMIOverrideBuilder) *ImageOverridesBuilder {
	b.aws = make([]*AMIOverrideBuilder, len(values))
	copy(b.aws, values)
	b.bitmap_ |= 8
	return b
}

// GCP sets the value of the 'GCP' attribute to the given values.
func (b *ImageOverridesBuilder) GCP(values ...*GCPImageOverrideBuilder) *ImageOverridesBuilder {
	b.gcp = make([]*GCPImageOverrideBuilder, len(values))
	copy(b.gcp, values)
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ImageOverridesBuilder) Copy(object *ImageOverrides) *ImageOverridesBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
