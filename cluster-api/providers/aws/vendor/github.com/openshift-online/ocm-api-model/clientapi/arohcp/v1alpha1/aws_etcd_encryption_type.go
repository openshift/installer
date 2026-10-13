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

// AwsEtcdEncryption represents the values of the 'aws_etcd_encryption' type.
//
// Contains the necessary attributes to support etcd encryption for AWS based clusters.
type AwsEtcdEncryption struct {
	fieldSet_ []bool
	kmsKeyARN string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AwsEtcdEncryption) Empty() bool {
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

// KMSKeyARN returns the value of the 'KMS_key_ARN' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ARN of the KMS to be used for the etcd encryption
func (o *AwsEtcdEncryption) KMSKeyARN() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.kmsKeyARN
	}
	return ""
}

// GetKMSKeyARN returns the value of the 'KMS_key_ARN' attribute and
// a flag indicating if the attribute has a value.
//
// ARN of the KMS to be used for the etcd encryption
func (o *AwsEtcdEncryption) GetKMSKeyARN() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.kmsKeyARN
	}
	return
}

// AwsEtcdEncryptionListKind is the name of the type used to represent list of objects of
// type 'aws_etcd_encryption'.
const AwsEtcdEncryptionListKind = "AwsEtcdEncryptionList"

// AwsEtcdEncryptionListLinkKind is the name of the type used to represent links to list
// of objects of type 'aws_etcd_encryption'.
const AwsEtcdEncryptionListLinkKind = "AwsEtcdEncryptionListLink"

// AwsEtcdEncryptionNilKind is the name of the type used to nil lists of objects of
// type 'aws_etcd_encryption'.
const AwsEtcdEncryptionListNilKind = "AwsEtcdEncryptionListNil"

// AwsEtcdEncryptionList is a list of values of the 'aws_etcd_encryption' type.
type AwsEtcdEncryptionList struct {
	href  string
	link  bool
	items []*AwsEtcdEncryption
}

// Len returns the length of the list.
func (l *AwsEtcdEncryptionList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AwsEtcdEncryptionList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AwsEtcdEncryptionList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AwsEtcdEncryptionList) SetItems(items []*AwsEtcdEncryption) {
	l.items = items
}

// Items returns the items of the list.
func (l *AwsEtcdEncryptionList) Items() []*AwsEtcdEncryption {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AwsEtcdEncryptionList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AwsEtcdEncryptionList) Get(i int) *AwsEtcdEncryption {
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
func (l *AwsEtcdEncryptionList) Slice() []*AwsEtcdEncryption {
	var slice []*AwsEtcdEncryption
	if l == nil {
		slice = make([]*AwsEtcdEncryption, 0)
	} else {
		slice = make([]*AwsEtcdEncryption, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AwsEtcdEncryptionList) Each(f func(item *AwsEtcdEncryption) bool) {
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
func (l *AwsEtcdEncryptionList) Range(f func(index int, item *AwsEtcdEncryption) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
