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

// Representation of Red Hat's Terms and Conditions for using OpenShift Dedicated and Amazon Red Hat OpenShift [Terms]
// review requests.
type TermsReviewRequestBuilder struct {
	fieldSet_          []bool
	accountUsername    string
	eventCode          string
	siteCode           string
	checkOptionalTerms bool
}

// NewTermsReviewRequest creates a new builder of 'terms_review_request' objects.
func NewTermsReviewRequest() *TermsReviewRequestBuilder {
	return &TermsReviewRequestBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *TermsReviewRequestBuilder) Empty() bool {
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

// AccountUsername sets the value of the 'account_username' attribute to the given value.
func (b *TermsReviewRequestBuilder) AccountUsername(value string) *TermsReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.accountUsername = value
	b.fieldSet_[0] = true
	return b
}

// CheckOptionalTerms sets the value of the 'check_optional_terms' attribute to the given value.
func (b *TermsReviewRequestBuilder) CheckOptionalTerms(value bool) *TermsReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.checkOptionalTerms = value
	b.fieldSet_[1] = true
	return b
}

// EventCode sets the value of the 'event_code' attribute to the given value.
func (b *TermsReviewRequestBuilder) EventCode(value string) *TermsReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.eventCode = value
	b.fieldSet_[2] = true
	return b
}

// SiteCode sets the value of the 'site_code' attribute to the given value.
func (b *TermsReviewRequestBuilder) SiteCode(value string) *TermsReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.siteCode = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *TermsReviewRequestBuilder) Copy(object *TermsReviewRequest) *TermsReviewRequestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.accountUsername = object.accountUsername
	b.checkOptionalTerms = object.checkOptionalTerms
	b.eventCode = object.eventCode
	b.siteCode = object.siteCode
	return b
}

// Build creates a 'terms_review_request' object using the configuration stored in the builder.
func (b *TermsReviewRequestBuilder) Build() (object *TermsReviewRequest, err error) {
	object = new(TermsReviewRequest)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.accountUsername = b.accountUsername
	object.checkOptionalTerms = b.checkOptionalTerms
	object.eventCode = b.eventCode
	object.siteCode = b.siteCode
	return
}
