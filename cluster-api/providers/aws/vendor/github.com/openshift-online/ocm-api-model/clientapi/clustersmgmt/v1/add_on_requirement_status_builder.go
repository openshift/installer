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

// Representation of an add-on requirement status.
type AddOnRequirementStatusBuilder struct {
	fieldSet_ []bool
	errorMsgs []string
	fulfilled bool
}

// NewAddOnRequirementStatus creates a new builder of 'add_on_requirement_status' objects.
func NewAddOnRequirementStatus() *AddOnRequirementStatusBuilder {
	return &AddOnRequirementStatusBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnRequirementStatusBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// ErrorMsgs sets the value of the 'error_msgs' attribute to the given values.
func (b *AddOnRequirementStatusBuilder) ErrorMsgs(values ...string) *AddOnRequirementStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.errorMsgs = make([]string, len(values))
	copy(b.errorMsgs, values)
	b.fieldSet_[0] = true
	return b
}

// Fulfilled sets the value of the 'fulfilled' attribute to the given value.
func (b *AddOnRequirementStatusBuilder) Fulfilled(value bool) *AddOnRequirementStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.fulfilled = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnRequirementStatusBuilder) Copy(object *AddOnRequirementStatus) *AddOnRequirementStatusBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.errorMsgs != nil {
		object.errorMsgs = make([]string, len(b.errorMsgs))
		copy(object.errorMsgs, b.errorMsgs)
	}
	object.fulfilled = b.fulfilled
	return
}
