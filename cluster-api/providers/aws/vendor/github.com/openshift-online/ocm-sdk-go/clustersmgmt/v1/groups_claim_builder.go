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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// GroupsClaimBuilder contains the data and logic needed to build 'groups_claim' objects.
type GroupsClaimBuilder struct {
	bitmap_ uint32
	claim   string
	prefix  string
}

// NewGroupsClaim creates a new builder of 'groups_claim' objects.
func NewGroupsClaim() *GroupsClaimBuilder {
	return &GroupsClaimBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GroupsClaimBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Claim sets the value of the 'claim' attribute to the given value.
func (b *GroupsClaimBuilder) Claim(value string) *GroupsClaimBuilder {
	b.claim = value
	b.bitmap_ |= 1
	return b
}

// Prefix sets the value of the 'prefix' attribute to the given value.
func (b *GroupsClaimBuilder) Prefix(value string) *GroupsClaimBuilder {
	b.prefix = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GroupsClaimBuilder) Copy(object *GroupsClaim) *GroupsClaimBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.claim = object.claim
	b.prefix = object.prefix
	return b
}

// Build creates a 'groups_claim' object using the configuration stored in the builder.
func (b *GroupsClaimBuilder) Build() (object *GroupsClaim, err error) {
	object = new(GroupsClaim)
	object.bitmap_ = b.bitmap_
	object.claim = b.claim
	object.prefix = b.prefix
	return
}
