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

// AddonStatusConditionBuilder contains the data and logic needed to build 'addon_status_condition' objects.
//
// Representation of an addon status condition type.
type AddonStatusConditionBuilder struct {
	bitmap_     uint32
	message     string
	reason      string
	statusType  AddonStatusConditionType
	statusValue AddonStatusConditionValue
}

// NewAddonStatusCondition creates a new builder of 'addon_status_condition' objects.
func NewAddonStatusCondition() *AddonStatusConditionBuilder {
	return &AddonStatusConditionBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonStatusConditionBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Message sets the value of the 'message' attribute to the given value.
func (b *AddonStatusConditionBuilder) Message(value string) *AddonStatusConditionBuilder {
	b.message = value
	b.bitmap_ |= 1
	return b
}

// Reason sets the value of the 'reason' attribute to the given value.
func (b *AddonStatusConditionBuilder) Reason(value string) *AddonStatusConditionBuilder {
	b.reason = value
	b.bitmap_ |= 2
	return b
}

// StatusType sets the value of the 'status_type' attribute to the given value.
//
// Representation of an addon status condition type field.
func (b *AddonStatusConditionBuilder) StatusType(value AddonStatusConditionType) *AddonStatusConditionBuilder {
	b.statusType = value
	b.bitmap_ |= 4
	return b
}

// StatusValue sets the value of the 'status_value' attribute to the given value.
//
// Representation of an addon status condition value field.
func (b *AddonStatusConditionBuilder) StatusValue(value AddonStatusConditionValue) *AddonStatusConditionBuilder {
	b.statusValue = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonStatusConditionBuilder) Copy(object *AddonStatusCondition) *AddonStatusConditionBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.message = object.message
	b.reason = object.reason
	b.statusType = object.statusType
	b.statusValue = object.statusValue
	return b
}

// Build creates a 'addon_status_condition' object using the configuration stored in the builder.
func (b *AddonStatusConditionBuilder) Build() (object *AddonStatusCondition, err error) {
	object = new(AddonStatusCondition)
	object.bitmap_ = b.bitmap_
	object.message = b.message
	object.reason = b.reason
	object.statusType = b.statusType
	object.statusValue = b.statusValue
	return
}
