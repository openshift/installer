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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// FeatureToggleBuilder contains the data and logic needed to build 'feature_toggle' objects.
type FeatureToggleBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	enabled bool
}

// NewFeatureToggle creates a new builder of 'feature_toggle' objects.
func NewFeatureToggle() *FeatureToggleBuilder {
	return &FeatureToggleBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *FeatureToggleBuilder) Link(value bool) *FeatureToggleBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *FeatureToggleBuilder) ID(value string) *FeatureToggleBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *FeatureToggleBuilder) HREF(value string) *FeatureToggleBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *FeatureToggleBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *FeatureToggleBuilder) Enabled(value bool) *FeatureToggleBuilder {
	b.enabled = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *FeatureToggleBuilder) Copy(object *FeatureToggle) *FeatureToggleBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.enabled = object.enabled
	return b
}

// Build creates a 'feature_toggle' object using the configuration stored in the builder.
func (b *FeatureToggleBuilder) Build() (object *FeatureToggle, err error) {
	object = new(FeatureToggle)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.enabled = b.enabled
	return
}
