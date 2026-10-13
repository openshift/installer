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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

type WifCredentialRequestBuilder struct {
	fieldSet_           []bool
	secretRef           *WifSecretRefBuilder
	serviceAccountNames []string
}

// NewWifCredentialRequest creates a new builder of 'wif_credential_request' objects.
func NewWifCredentialRequest() *WifCredentialRequestBuilder {
	return &WifCredentialRequestBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifCredentialRequestBuilder) Empty() bool {
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

// SecretRef sets the value of the 'secret_ref' attribute to the given value.
func (b *WifCredentialRequestBuilder) SecretRef(value *WifSecretRefBuilder) *WifCredentialRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.secretRef = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// ServiceAccountNames sets the value of the 'service_account_names' attribute to the given values.
func (b *WifCredentialRequestBuilder) ServiceAccountNames(values ...string) *WifCredentialRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.serviceAccountNames = make([]string, len(values))
	copy(b.serviceAccountNames, values)
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifCredentialRequestBuilder) Copy(object *WifCredentialRequest) *WifCredentialRequestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.secretRef != nil {
		b.secretRef = NewWifSecretRef().Copy(object.secretRef)
	} else {
		b.secretRef = nil
	}
	if object.serviceAccountNames != nil {
		b.serviceAccountNames = make([]string, len(object.serviceAccountNames))
		copy(b.serviceAccountNames, object.serviceAccountNames)
	} else {
		b.serviceAccountNames = nil
	}
	return b
}

// Build creates a 'wif_credential_request' object using the configuration stored in the builder.
func (b *WifCredentialRequestBuilder) Build() (object *WifCredentialRequest, err error) {
	object = new(WifCredentialRequest)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.secretRef != nil {
		object.secretRef, err = b.secretRef.Build()
		if err != nil {
			return
		}
	}
	if b.serviceAccountNames != nil {
		object.serviceAccountNames = make([]string, len(b.serviceAccountNames))
		copy(object.serviceAccountNames, b.serviceAccountNames)
	}
	return
}
