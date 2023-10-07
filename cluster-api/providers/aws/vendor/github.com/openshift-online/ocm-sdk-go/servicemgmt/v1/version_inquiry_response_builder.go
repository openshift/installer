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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

// VersionInquiryResponseBuilder contains the data and logic needed to build 'version_inquiry_response' objects.
type VersionInquiryResponseBuilder struct {
	bitmap_ uint32
	version string
}

// NewVersionInquiryResponse creates a new builder of 'version_inquiry_response' objects.
func NewVersionInquiryResponse() *VersionInquiryResponseBuilder {
	return &VersionInquiryResponseBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *VersionInquiryResponseBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Version sets the value of the 'version' attribute to the given value.
func (b *VersionInquiryResponseBuilder) Version(value string) *VersionInquiryResponseBuilder {
	b.version = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *VersionInquiryResponseBuilder) Copy(object *VersionInquiryResponse) *VersionInquiryResponseBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.version = object.version
	return b
}

// Build creates a 'version_inquiry_response' object using the configuration stored in the builder.
func (b *VersionInquiryResponseBuilder) Build() (object *VersionInquiryResponse, err error) {
	object = new(VersionInquiryResponse)
	object.bitmap_ = b.bitmap_
	object.version = b.version
	return
}
