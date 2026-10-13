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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

type AccessTokenBuilder struct {
	fieldSet_ []bool
	auths     map[string]*AccessTokenAuthBuilder
}

// NewAccessToken creates a new builder of 'access_token' objects.
func NewAccessToken() *AccessTokenBuilder {
	return &AccessTokenBuilder{
		fieldSet_: make([]bool, 1),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccessTokenBuilder) Empty() bool {
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

// Auths sets the value of the 'auths' attribute to the given value.
func (b *AccessTokenBuilder) Auths(value map[string]*AccessTokenAuthBuilder) *AccessTokenBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 1)
	}
	b.auths = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccessTokenBuilder) Copy(object *AccessToken) *AccessTokenBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if len(object.auths) > 0 {
		b.auths = map[string]*AccessTokenAuthBuilder{}
		for k, v := range object.auths {
			b.auths[k] = NewAccessTokenAuth().Copy(v)
		}
	} else {
		b.auths = nil
	}
	return b
}

// Build creates a 'access_token' object using the configuration stored in the builder.
func (b *AccessTokenBuilder) Build() (object *AccessToken, err error) {
	object = new(AccessToken)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.auths != nil {
		object.auths = make(map[string]*AccessTokenAuth)
		for k, v := range b.auths {
			object.auths[k], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
