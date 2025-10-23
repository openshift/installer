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

// AzureKmsKey represents the values of the 'azure_kms_key' type.
//
// Contains the necessary attributes to support KMS encryption key for Azure based clusters
type AzureKmsKey struct {
	fieldSet_    []bool
	keyName      string
	keyVaultName string
	keyVersion   string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AzureKmsKey) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}
	for _, set := range o.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// KeyName returns the value of the 'key_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// key_name is the name of the Azure Key Vault Key
// Required during creation.
func (o *AzureKmsKey) KeyName() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.keyName
	}
	return ""
}

// GetKeyName returns the value of the 'key_name' attribute and
// a flag indicating if the attribute has a value.
//
// key_name is the name of the Azure Key Vault Key
// Required during creation.
func (o *AzureKmsKey) GetKeyName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.keyName
	}
	return
}

// KeyVaultName returns the value of the 'key_vault_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// key_vault_name is the name of the Azure Key Vault that contains the encryption key
// Required during creation.
func (o *AzureKmsKey) KeyVaultName() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.keyVaultName
	}
	return ""
}

// GetKeyVaultName returns the value of the 'key_vault_name' attribute and
// a flag indicating if the attribute has a value.
//
// key_vault_name is the name of the Azure Key Vault that contains the encryption key
// Required during creation.
func (o *AzureKmsKey) GetKeyVaultName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.keyVaultName
	}
	return
}

// KeyVersion returns the value of the 'key_version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// key_version is the version of the Azure Key Vault key
// Required during creation.
func (o *AzureKmsKey) KeyVersion() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.keyVersion
	}
	return ""
}

// GetKeyVersion returns the value of the 'key_version' attribute and
// a flag indicating if the attribute has a value.
//
// key_version is the version of the Azure Key Vault key
// Required during creation.
func (o *AzureKmsKey) GetKeyVersion() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.keyVersion
	}
	return
}

// AzureKmsKeyListKind is the name of the type used to represent list of objects of
// type 'azure_kms_key'.
const AzureKmsKeyListKind = "AzureKmsKeyList"

// AzureKmsKeyListLinkKind is the name of the type used to represent links to list
// of objects of type 'azure_kms_key'.
const AzureKmsKeyListLinkKind = "AzureKmsKeyListLink"

// AzureKmsKeyNilKind is the name of the type used to nil lists of objects of
// type 'azure_kms_key'.
const AzureKmsKeyListNilKind = "AzureKmsKeyListNil"

// AzureKmsKeyList is a list of values of the 'azure_kms_key' type.
type AzureKmsKeyList struct {
	href  string
	link  bool
	items []*AzureKmsKey
}

// Len returns the length of the list.
func (l *AzureKmsKeyList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AzureKmsKeyList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AzureKmsKeyList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AzureKmsKeyList) SetItems(items []*AzureKmsKey) {
	l.items = items
}

// Items returns the items of the list.
func (l *AzureKmsKeyList) Items() []*AzureKmsKey {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AzureKmsKeyList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AzureKmsKeyList) Get(i int) *AzureKmsKey {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *AzureKmsKeyList) Slice() []*AzureKmsKey {
	var slice []*AzureKmsKey
	if l == nil {
		slice = make([]*AzureKmsKey, 0)
	} else {
		slice = make([]*AzureKmsKey, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AzureKmsKeyList) Each(f func(item *AzureKmsKey) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *AzureKmsKeyList) Range(f func(index int, item *AzureKmsKey) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
