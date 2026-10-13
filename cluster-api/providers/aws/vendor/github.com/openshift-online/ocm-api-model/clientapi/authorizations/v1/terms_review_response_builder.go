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
// review response.
type TermsReviewResponseBuilder struct {
	fieldSet_      []bool
	accountId      string
	organizationID string
	redirectUrl    string
	termsAvailable bool
	termsRequired  bool
}

// NewTermsReviewResponse creates a new builder of 'terms_review_response' objects.
func NewTermsReviewResponse() *TermsReviewResponseBuilder {
	return &TermsReviewResponseBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *TermsReviewResponseBuilder) Empty() bool {
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

// AccountId sets the value of the 'account_id' attribute to the given value.
func (b *TermsReviewResponseBuilder) AccountId(value string) *TermsReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.accountId = value
	b.fieldSet_[0] = true
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *TermsReviewResponseBuilder) OrganizationID(value string) *TermsReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.organizationID = value
	b.fieldSet_[1] = true
	return b
}

// RedirectUrl sets the value of the 'redirect_url' attribute to the given value.
func (b *TermsReviewResponseBuilder) RedirectUrl(value string) *TermsReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.redirectUrl = value
	b.fieldSet_[2] = true
	return b
}

// TermsAvailable sets the value of the 'terms_available' attribute to the given value.
func (b *TermsReviewResponseBuilder) TermsAvailable(value bool) *TermsReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.termsAvailable = value
	b.fieldSet_[3] = true
	return b
}

// TermsRequired sets the value of the 'terms_required' attribute to the given value.
func (b *TermsReviewResponseBuilder) TermsRequired(value bool) *TermsReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.termsRequired = value
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *TermsReviewResponseBuilder) Copy(object *TermsReviewResponse) *TermsReviewResponseBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.accountId = b.accountId
	object.organizationID = b.organizationID
	object.redirectUrl = b.redirectUrl
	object.termsAvailable = b.termsAvailable
	object.termsRequired = b.termsRequired
	return
}
