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

// GCPEncryptionKey represents the values of the 'GCP_encryption_key' type.
//
// GCP Encryption Key for CCS clusters.
type GCPEncryptionKey struct {
	fieldSet_            []bool
	kmsKeyServiceAccount string
	keyLocation          string
	keyName              string
	keyRing              string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *GCPEncryptionKey) Empty() bool {
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

// KMSKeyServiceAccount returns the value of the 'KMS_key_service_account' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Service account used to access the KMS key
func (o *GCPEncryptionKey) KMSKeyServiceAccount() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.kmsKeyServiceAccount
	}
	return ""
}

// GetKMSKeyServiceAccount returns the value of the 'KMS_key_service_account' attribute and
// a flag indicating if the attribute has a value.
//
// Service account used to access the KMS key
func (o *GCPEncryptionKey) GetKMSKeyServiceAccount() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.kmsKeyServiceAccount
	}
	return
}

// KeyLocation returns the value of the 'key_location' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Location of the encryption key ring
func (o *GCPEncryptionKey) KeyLocation() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.keyLocation
	}
	return ""
}

// GetKeyLocation returns the value of the 'key_location' attribute and
// a flag indicating if the attribute has a value.
//
// Location of the encryption key ring
func (o *GCPEncryptionKey) GetKeyLocation() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.keyLocation
	}
	return
}

// KeyName returns the value of the 'key_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the encryption key
func (o *GCPEncryptionKey) KeyName() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.keyName
	}
	return ""
}

// GetKeyName returns the value of the 'key_name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the encryption key
func (o *GCPEncryptionKey) GetKeyName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.keyName
	}
	return
}

// KeyRing returns the value of the 'key_ring' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the key ring the encryption key is located on
func (o *GCPEncryptionKey) KeyRing() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.keyRing
	}
	return ""
}

// GetKeyRing returns the value of the 'key_ring' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the key ring the encryption key is located on
func (o *GCPEncryptionKey) GetKeyRing() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.keyRing
	}
	return
}

// GCPEncryptionKeyListKind is the name of the type used to represent list of objects of
// type 'GCP_encryption_key'.
const GCPEncryptionKeyListKind = "GCPEncryptionKeyList"

// GCPEncryptionKeyListLinkKind is the name of the type used to represent links to list
// of objects of type 'GCP_encryption_key'.
const GCPEncryptionKeyListLinkKind = "GCPEncryptionKeyListLink"

// GCPEncryptionKeyNilKind is the name of the type used to nil lists of objects of
// type 'GCP_encryption_key'.
const GCPEncryptionKeyListNilKind = "GCPEncryptionKeyListNil"

// GCPEncryptionKeyList is a list of values of the 'GCP_encryption_key' type.
type GCPEncryptionKeyList struct {
	href  string
	link  bool
	items []*GCPEncryptionKey
}

// Len returns the length of the list.
func (l *GCPEncryptionKeyList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *GCPEncryptionKeyList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *GCPEncryptionKeyList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *GCPEncryptionKeyList) SetItems(items []*GCPEncryptionKey) {
	l.items = items
}

// Items returns the items of the list.
func (l *GCPEncryptionKeyList) Items() []*GCPEncryptionKey {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *GCPEncryptionKeyList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *GCPEncryptionKeyList) Get(i int) *GCPEncryptionKey {
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
func (l *GCPEncryptionKeyList) Slice() []*GCPEncryptionKey {
	var slice []*GCPEncryptionKey
	if l == nil {
		slice = make([]*GCPEncryptionKey, 0)
	} else {
		slice = make([]*GCPEncryptionKey, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *GCPEncryptionKeyList) Each(f func(item *GCPEncryptionKey) bool) {
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
func (l *GCPEncryptionKeyList) Range(f func(index int, item *GCPEncryptionKey) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
