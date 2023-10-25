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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

// ErrorBuilder contains the data and logic needed to build 'error' objects.
//
// Definition of a Web RCA error.
type ErrorBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	code    string
	reason  string
}

// NewError creates a new builder of 'error' objects.
func NewError() *ErrorBuilder {
	return &ErrorBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ErrorBuilder) Link(value bool) *ErrorBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ErrorBuilder) ID(value string) *ErrorBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ErrorBuilder) HREF(value string) *ErrorBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ErrorBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Code sets the value of the 'code' attribute to the given value.
func (b *ErrorBuilder) Code(value string) *ErrorBuilder {
	b.code = value
	b.bitmap_ |= 8
	return b
}

// Reason sets the value of the 'reason' attribute to the given value.
func (b *ErrorBuilder) Reason(value string) *ErrorBuilder {
	b.reason = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ErrorBuilder) Copy(object *Error) *ErrorBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.code = object.code
	b.reason = object.reason
	return b
}

// Build creates a 'error' object using the configuration stored in the builder.
func (b *ErrorBuilder) Build() (object *Error, err error) {
	object = new(Error)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.code = b.code
	object.reason = b.reason
	return
}
