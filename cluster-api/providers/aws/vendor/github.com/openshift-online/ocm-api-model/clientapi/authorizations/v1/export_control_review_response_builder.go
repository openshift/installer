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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/authorizations/v1

type ExportControlReviewResponseBuilder struct {
	fieldSet_  []bool
	restricted bool
}

// NewExportControlReviewResponse creates a new builder of 'export_control_review_response' objects.
func NewExportControlReviewResponse() *ExportControlReviewResponseBuilder {
	return &ExportControlReviewResponseBuilder{
		fieldSet_: make([]bool, 1),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ExportControlReviewResponseBuilder) Empty() bool {
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

// Restricted sets the value of the 'restricted' attribute to the given value.
func (b *ExportControlReviewResponseBuilder) Restricted(value bool) *ExportControlReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 1)
	}
	b.restricted = value
	b.fieldSet_[0] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ExportControlReviewResponseBuilder) Copy(object *ExportControlReviewResponse) *ExportControlReviewResponseBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.restricted = object.restricted
	return b
}

// Build creates a 'export_control_review_response' object using the configuration stored in the builder.
func (b *ExportControlReviewResponseBuilder) Build() (object *ExportControlReviewResponse, err error) {
	object = new(ExportControlReviewResponse)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.restricted = b.restricted
	return
}
