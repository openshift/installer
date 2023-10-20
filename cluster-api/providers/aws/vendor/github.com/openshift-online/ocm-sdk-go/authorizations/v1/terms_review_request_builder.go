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

// TermsReviewRequestBuilder contains the data and logic needed to build 'terms_review_request' objects.
//
// Representation of Red Hat's Terms and Conditions for using OpenShift Dedicated and Amazon Red Hat OpenShift [Terms]
// review requests.
type TermsReviewRequestBuilder struct {
	bitmap_            uint32
	accountUsername    string
	eventCode          string
	siteCode           string
	checkOptionalTerms bool
}

// NewTermsReviewRequest creates a new builder of 'terms_review_request' objects.
func NewTermsReviewRequest() *TermsReviewRequestBuilder {
	return &TermsReviewRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *TermsReviewRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AccountUsername sets the value of the 'account_username' attribute to the given value.
func (b *TermsReviewRequestBuilder) AccountUsername(value string) *TermsReviewRequestBuilder {
	b.accountUsername = value
	b.bitmap_ |= 1
	return b
}

// CheckOptionalTerms sets the value of the 'check_optional_terms' attribute to the given value.
func (b *TermsReviewRequestBuilder) CheckOptionalTerms(value bool) *TermsReviewRequestBuilder {
	b.checkOptionalTerms = value
	b.bitmap_ |= 2
	return b
}

// EventCode sets the value of the 'event_code' attribute to the given value.
func (b *TermsReviewRequestBuilder) EventCode(value string) *TermsReviewRequestBuilder {
	b.eventCode = value
	b.bitmap_ |= 4
	return b
}

// SiteCode sets the value of the 'site_code' attribute to the given value.
func (b *TermsReviewRequestBuilder) SiteCode(value string) *TermsReviewRequestBuilder {
	b.siteCode = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *TermsReviewRequestBuilder) Copy(object *TermsReviewRequest) *TermsReviewRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.accountUsername = object.accountUsername
	b.checkOptionalTerms = object.checkOptionalTerms
	b.eventCode = object.eventCode
	b.siteCode = object.siteCode
	return b
}

// Build creates a 'terms_review_request' object using the configuration stored in the builder.
func (b *TermsReviewRequestBuilder) Build() (object *TermsReviewRequest, err error) {
	object = new(TermsReviewRequest)
	object.bitmap_ = b.bitmap_
	object.accountUsername = b.accountUsername
	object.checkOptionalTerms = b.checkOptionalTerms
	object.eventCode = b.eventCode
	object.siteCode = b.siteCode
	return
}
