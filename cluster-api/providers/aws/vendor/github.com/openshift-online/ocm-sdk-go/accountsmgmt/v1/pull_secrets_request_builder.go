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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// PullSecretsRequestBuilder contains the data and logic needed to build 'pull_secrets_request' objects.
type PullSecretsRequestBuilder struct {
	bitmap_            uint32
	externalResourceId string
}

// NewPullSecretsRequest creates a new builder of 'pull_secrets_request' objects.
func NewPullSecretsRequest() *PullSecretsRequestBuilder {
	return &PullSecretsRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *PullSecretsRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ExternalResourceId sets the value of the 'external_resource_id' attribute to the given value.
func (b *PullSecretsRequestBuilder) ExternalResourceId(value string) *PullSecretsRequestBuilder {
	b.externalResourceId = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *PullSecretsRequestBuilder) Copy(object *PullSecretsRequest) *PullSecretsRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.externalResourceId = object.externalResourceId
	return b
}

// Build creates a 'pull_secrets_request' object using the configuration stored in the builder.
func (b *PullSecretsRequestBuilder) Build() (object *PullSecretsRequest, err error) {
	object = new(PullSecretsRequest)
	object.bitmap_ = b.bitmap_
	object.externalResourceId = b.externalResourceId
	return
}
