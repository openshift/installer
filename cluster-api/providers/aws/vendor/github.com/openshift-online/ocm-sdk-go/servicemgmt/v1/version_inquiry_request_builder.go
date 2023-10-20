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

// VersionInquiryRequestBuilder contains the data and logic needed to build 'version_inquiry_request' objects.
type VersionInquiryRequestBuilder struct {
	bitmap_     uint32
	serviceType string
}

// NewVersionInquiryRequest creates a new builder of 'version_inquiry_request' objects.
func NewVersionInquiryRequest() *VersionInquiryRequestBuilder {
	return &VersionInquiryRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *VersionInquiryRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ServiceType sets the value of the 'service_type' attribute to the given value.
func (b *VersionInquiryRequestBuilder) ServiceType(value string) *VersionInquiryRequestBuilder {
	b.serviceType = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *VersionInquiryRequestBuilder) Copy(object *VersionInquiryRequest) *VersionInquiryRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.serviceType = object.serviceType
	return b
}

// Build creates a 'version_inquiry_request' object using the configuration stored in the builder.
func (b *VersionInquiryRequestBuilder) Build() (object *VersionInquiryRequest, err error) {
	object = new(VersionInquiryRequest)
	object.bitmap_ = b.bitmap_
	object.serviceType = b.serviceType
	return
}
