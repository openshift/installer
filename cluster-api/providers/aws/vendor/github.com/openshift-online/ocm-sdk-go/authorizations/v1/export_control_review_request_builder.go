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

package v1 // github.com/openshift-online/ocm-sdk-go/authorizations/v1

// ExportControlReviewRequestBuilder contains the data and logic needed to build 'export_control_review_request' objects.
type ExportControlReviewRequestBuilder struct {
	bitmap_         uint32
	accountUsername string
}

// NewExportControlReviewRequest creates a new builder of 'export_control_review_request' objects.
func NewExportControlReviewRequest() *ExportControlReviewRequestBuilder {
	return &ExportControlReviewRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ExportControlReviewRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AccountUsername sets the value of the 'account_username' attribute to the given value.
func (b *ExportControlReviewRequestBuilder) AccountUsername(value string) *ExportControlReviewRequestBuilder {
	b.accountUsername = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ExportControlReviewRequestBuilder) Copy(object *ExportControlReviewRequest) *ExportControlReviewRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.accountUsername = object.accountUsername
	return b
}

// Build creates a 'export_control_review_request' object using the configuration stored in the builder.
func (b *ExportControlReviewRequestBuilder) Build() (object *ExportControlReviewRequest, err error) {
	object = new(ExportControlReviewRequest)
	object.bitmap_ = b.bitmap_
	object.accountUsername = b.accountUsername
	return
}
