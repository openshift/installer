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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// HTPasswdUserListBuilder contains the data and logic needed to build
// 'HT_passwd_user' objects.
type HTPasswdUserListBuilder struct {
	items []*HTPasswdUserBuilder
}

// NewHTPasswdUserList creates a new builder of 'HT_passwd_user' objects.
func NewHTPasswdUserList() *HTPasswdUserListBuilder {
	return new(HTPasswdUserListBuilder)
}

// Items sets the items of the list.
func (b *HTPasswdUserListBuilder) Items(values ...*HTPasswdUserBuilder) *HTPasswdUserListBuilder {
	b.items = make([]*HTPasswdUserBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *HTPasswdUserListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *HTPasswdUserListBuilder) Copy(list *HTPasswdUserList) *HTPasswdUserListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*HTPasswdUserBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewHTPasswdUser().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'HT_passwd_user' objects using the
// configuration stored in the builder.
func (b *HTPasswdUserListBuilder) Build() (list *HTPasswdUserList, err error) {
	items := make([]*HTPasswdUser, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(HTPasswdUserList)
	list.items = items
	return
}
