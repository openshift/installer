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

// TermsReviewResponseBuilder contains the data and logic needed to build 'terms_review_response' objects.
//
// Representation of Red Hat's Terms and Conditions for using OpenShift Dedicated and Amazon Red Hat OpenShift [Terms]
// review response.
type TermsReviewResponseBuilder struct {
	bitmap_        uint32
	accountId      string
	organizationID string
	redirectUrl    string
	termsAvailable bool
	termsRequired  bool
}

// NewTermsReviewResponse creates a new builder of 'terms_review_response' objects.
func NewTermsReviewResponse() *TermsReviewResponseBuilder {
	return &TermsReviewResponseBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *TermsReviewResponseBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AccountId sets the value of the 'account_id' attribute to the given value.
func (b *TermsReviewResponseBuilder) AccountId(value string) *TermsReviewResponseBuilder {
	b.accountId = value
	b.bitmap_ |= 1
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *TermsReviewResponseBuilder) OrganizationID(value string) *TermsReviewResponseBuilder {
	b.organizationID = value
	b.bitmap_ |= 2
	return b
}

// RedirectUrl sets the value of the 'redirect_url' attribute to the given value.
func (b *TermsReviewResponseBuilder) RedirectUrl(value string) *TermsReviewResponseBuilder {
	b.redirectUrl = value
	b.bitmap_ |= 4
	return b
}

// TermsAvailable sets the value of the 'terms_available' attribute to the given value.
func (b *TermsReviewResponseBuilder) TermsAvailable(value bool) *TermsReviewResponseBuilder {
	b.termsAvailable = value
	b.bitmap_ |= 8
	return b
}

// TermsRequired sets the value of the 'terms_required' attribute to the given value.
func (b *TermsReviewResponseBuilder) TermsRequired(value bool) *TermsReviewResponseBuilder {
	b.termsRequired = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *TermsReviewResponseBuilder) Copy(object *TermsReviewResponse) *TermsReviewResponseBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.accountId = object.accountId
	b.organizationID = object.organizationID
	b.redirectUrl = object.redirectUrl
	b.termsAvailable = object.termsAvailable
	b.termsRequired = object.termsRequired
	return b
}

// Build creates a 'terms_review_response' object using the configuration stored in the builder.
func (b *TermsReviewResponseBuilder) Build() (object *TermsReviewResponse, err error) {
	object = new(TermsReviewResponse)
	object.bitmap_ = b.bitmap_
	object.accountId = b.accountId
	object.organizationID = b.organizationID
	object.redirectUrl = b.redirectUrl
	object.termsAvailable = b.termsAvailable
	object.termsRequired = b.termsRequired
	return
}
