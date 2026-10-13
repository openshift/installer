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

// BackupListBuilder contains the data and logic needed to build
// 'backup' objects.
type BackupListBuilder struct {
	items []*BackupBuilder
}

// NewBackupList creates a new builder of 'backup' objects.
func NewBackupList() *BackupListBuilder {
	return new(BackupListBuilder)
}

// Items sets the items of the list.
func (b *BackupListBuilder) Items(values ...*BackupBuilder) *BackupListBuilder {
	b.items = make([]*BackupBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *BackupListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *BackupListBuilder) Copy(list *BackupList) *BackupListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*BackupBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewBackup().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'backup' objects using the
// configuration stored in the builder.
func (b *BackupListBuilder) Build() (list *BackupList, err error) {
	items := make([]*Backup, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(BackupList)
	list.items = items
	return
}
