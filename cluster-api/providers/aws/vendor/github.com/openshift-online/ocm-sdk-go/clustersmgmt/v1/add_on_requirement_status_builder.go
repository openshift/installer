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

// AddOnRequirementStatusBuilder contains the data and logic needed to build 'add_on_requirement_status' objects.
//
// Representation of an add-on requirement status.
type AddOnRequirementStatusBuilder struct {
	bitmap_   uint32
	errorMsgs []string
	fulfilled bool
}

// NewAddOnRequirementStatus creates a new builder of 'add_on_requirement_status' objects.
func NewAddOnRequirementStatus() *AddOnRequirementStatusBuilder {
	return &AddOnRequirementStatusBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnRequirementStatusBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ErrorMsgs sets the value of the 'error_msgs' attribute to the given values.
func (b *AddOnRequirementStatusBuilder) ErrorMsgs(values ...string) *AddOnRequirementStatusBuilder {
	b.errorMsgs = make([]string, len(values))
	copy(b.errorMsgs, values)
	b.bitmap_ |= 1
	return b
}

// Fulfilled sets the value of the 'fulfilled' attribute to the given value.
func (b *AddOnRequirementStatusBuilder) Fulfilled(value bool) *AddOnRequirementStatusBuilder {
	b.fulfilled = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnRequirementStatusBuilder) Copy(object *AddOnRequirementStatus) *AddOnRequirementStatusBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.errorMsgs != nil {
		b.errorMsgs = make([]string, len(object.errorMsgs))
		copy(b.errorMsgs, object.errorMsgs)
	} else {
		b.errorMsgs = nil
	}
	b.fulfilled = object.fulfilled
	return b
}

// Build creates a 'add_on_requirement_status' object using the configuration stored in the builder.
func (b *AddOnRequirementStatusBuilder) Build() (object *AddOnRequirementStatus, err error) {
	object = new(AddOnRequirementStatus)
	object.bitmap_ = b.bitmap_
	if b.errorMsgs != nil {
		object.errorMsgs = make([]string, len(b.errorMsgs))
		copy(object.errorMsgs, b.errorMsgs)
	}
	object.fulfilled = b.fulfilled
	return
}
