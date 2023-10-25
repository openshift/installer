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

// PlanBuilder contains the data and logic needed to build 'plan' objects.
type PlanBuilder struct {
	bitmap_  uint32
	id       string
	href     string
	category string
	name     string
	type_    string
}

// NewPlan creates a new builder of 'plan' objects.
func NewPlan() *PlanBuilder {
	return &PlanBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *PlanBuilder) Link(value bool) *PlanBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *PlanBuilder) ID(value string) *PlanBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *PlanBuilder) HREF(value string) *PlanBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *PlanBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Category sets the value of the 'category' attribute to the given value.
func (b *PlanBuilder) Category(value string) *PlanBuilder {
	b.category = value
	b.bitmap_ |= 8
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *PlanBuilder) Name(value string) *PlanBuilder {
	b.name = value
	b.bitmap_ |= 16
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *PlanBuilder) Type(value string) *PlanBuilder {
	b.type_ = value
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *PlanBuilder) Copy(object *Plan) *PlanBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.category = object.category
	b.name = object.name
	b.type_ = object.type_
	return b
}

// Build creates a 'plan' object using the configuration stored in the builder.
func (b *PlanBuilder) Build() (object *Plan, err error) {
	object = new(Plan)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.category = b.category
	object.name = b.name
	object.type_ = b.type_
	return
}
