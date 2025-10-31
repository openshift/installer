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

// AzureEtcdDataEncryptionCustomerManaged represents the values of the 'azure_etcd_data_encryption_customer_managed' type.
//
// Contains the necessary attributes to support etcd data encryption with customer managed keys
// for Azure based clusters.
type AzureEtcdDataEncryptionCustomerManaged struct {
	fieldSet_      []bool
	encryptionType string
	kms            *AzureKmsEncryption
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AzureEtcdDataEncryptionCustomerManaged) Empty() bool {
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

// EncryptionType returns the value of the 'encryption_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The encryption type used.
// Accepted values are: "kms".
// By default, "kms" is used.
func (o *AzureEtcdDataEncryptionCustomerManaged) EncryptionType() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.encryptionType
	}
	return ""
}

// GetEncryptionType returns the value of the 'encryption_type' attribute and
// a flag indicating if the attribute has a value.
//
// The encryption type used.
// Accepted values are: "kms".
// By default, "kms" is used.
func (o *AzureEtcdDataEncryptionCustomerManaged) GetEncryptionType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.encryptionType
	}
	return
}

// Kms returns the value of the 'kms' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The KMS encryption configuration.
// Required when encryption_type is "kms".
func (o *AzureEtcdDataEncryptionCustomerManaged) Kms() *AzureKmsEncryption {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.kms
	}
	return nil
}

// GetKms returns the value of the 'kms' attribute and
// a flag indicating if the attribute has a value.
//
// The KMS encryption configuration.
// Required when encryption_type is "kms".
func (o *AzureEtcdDataEncryptionCustomerManaged) GetKms() (value *AzureKmsEncryption, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.kms
	}
	return
}

// AzureEtcdDataEncryptionCustomerManagedListKind is the name of the type used to represent list of objects of
// type 'azure_etcd_data_encryption_customer_managed'.
const AzureEtcdDataEncryptionCustomerManagedListKind = "AzureEtcdDataEncryptionCustomerManagedList"

// AzureEtcdDataEncryptionCustomerManagedListLinkKind is the name of the type used to represent links to list
// of objects of type 'azure_etcd_data_encryption_customer_managed'.
const AzureEtcdDataEncryptionCustomerManagedListLinkKind = "AzureEtcdDataEncryptionCustomerManagedListLink"

// AzureEtcdDataEncryptionCustomerManagedNilKind is the name of the type used to nil lists of objects of
// type 'azure_etcd_data_encryption_customer_managed'.
const AzureEtcdDataEncryptionCustomerManagedListNilKind = "AzureEtcdDataEncryptionCustomerManagedListNil"

// AzureEtcdDataEncryptionCustomerManagedList is a list of values of the 'azure_etcd_data_encryption_customer_managed' type.
type AzureEtcdDataEncryptionCustomerManagedList struct {
	href  string
	link  bool
	items []*AzureEtcdDataEncryptionCustomerManaged
}

// Len returns the length of the list.
func (l *AzureEtcdDataEncryptionCustomerManagedList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AzureEtcdDataEncryptionCustomerManagedList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AzureEtcdDataEncryptionCustomerManagedList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AzureEtcdDataEncryptionCustomerManagedList) SetItems(items []*AzureEtcdDataEncryptionCustomerManaged) {
	l.items = items
}

// Items returns the items of the list.
func (l *AzureEtcdDataEncryptionCustomerManagedList) Items() []*AzureEtcdDataEncryptionCustomerManaged {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AzureEtcdDataEncryptionCustomerManagedList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AzureEtcdDataEncryptionCustomerManagedList) Get(i int) *AzureEtcdDataEncryptionCustomerManaged {
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
func (l *AzureEtcdDataEncryptionCustomerManagedList) Slice() []*AzureEtcdDataEncryptionCustomerManaged {
	var slice []*AzureEtcdDataEncryptionCustomerManaged
	if l == nil {
		slice = make([]*AzureEtcdDataEncryptionCustomerManaged, 0)
	} else {
		slice = make([]*AzureEtcdDataEncryptionCustomerManaged, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AzureEtcdDataEncryptionCustomerManagedList) Each(f func(item *AzureEtcdDataEncryptionCustomerManaged) bool) {
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
func (l *AzureEtcdDataEncryptionCustomerManagedList) Range(f func(index int, item *AzureEtcdDataEncryptionCustomerManaged) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
