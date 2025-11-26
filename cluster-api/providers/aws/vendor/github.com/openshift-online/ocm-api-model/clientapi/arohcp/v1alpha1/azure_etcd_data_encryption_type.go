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

// AzureEtcdDataEncryption represents the values of the 'azure_etcd_data_encryption' type.
//
// Contains the necessary attributes to support data encryption for Azure based clusters.
type AzureEtcdDataEncryption struct {
	fieldSet_         []bool
	customerManaged   *AzureEtcdDataEncryptionCustomerManaged
	keyManagementMode string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AzureEtcdDataEncryption) Empty() bool {
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

// CustomerManaged returns the value of the 'customer_managed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Customer Managed encryption keys configuration.
// Required when key_management_mode is "customer_managed".
func (o *AzureEtcdDataEncryption) CustomerManaged() *AzureEtcdDataEncryptionCustomerManaged {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.customerManaged
	}
	return nil
}

// GetCustomerManaged returns the value of the 'customer_managed' attribute and
// a flag indicating if the attribute has a value.
//
// Customer Managed encryption keys configuration.
// Required when key_management_mode is "customer_managed".
func (o *AzureEtcdDataEncryption) GetCustomerManaged() (value *AzureEtcdDataEncryptionCustomerManaged, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.customerManaged
	}
	return
}

// KeyManagementMode returns the value of the 'key_management_mode' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The key management strategy used for the encryption key that encrypts the etcd data.
// Accepted values are: "customer_managed", "platform_managed".
// By default, "platform_managed" is used.
// Currently only "customer_managed" mode is supported.
func (o *AzureEtcdDataEncryption) KeyManagementMode() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.keyManagementMode
	}
	return ""
}

// GetKeyManagementMode returns the value of the 'key_management_mode' attribute and
// a flag indicating if the attribute has a value.
//
// The key management strategy used for the encryption key that encrypts the etcd data.
// Accepted values are: "customer_managed", "platform_managed".
// By default, "platform_managed" is used.
// Currently only "customer_managed" mode is supported.
func (o *AzureEtcdDataEncryption) GetKeyManagementMode() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.keyManagementMode
	}
	return
}

// AzureEtcdDataEncryptionListKind is the name of the type used to represent list of objects of
// type 'azure_etcd_data_encryption'.
const AzureEtcdDataEncryptionListKind = "AzureEtcdDataEncryptionList"

// AzureEtcdDataEncryptionListLinkKind is the name of the type used to represent links to list
// of objects of type 'azure_etcd_data_encryption'.
const AzureEtcdDataEncryptionListLinkKind = "AzureEtcdDataEncryptionListLink"

// AzureEtcdDataEncryptionNilKind is the name of the type used to nil lists of objects of
// type 'azure_etcd_data_encryption'.
const AzureEtcdDataEncryptionListNilKind = "AzureEtcdDataEncryptionListNil"

// AzureEtcdDataEncryptionList is a list of values of the 'azure_etcd_data_encryption' type.
type AzureEtcdDataEncryptionList struct {
	href  string
	link  bool
	items []*AzureEtcdDataEncryption
}

// Len returns the length of the list.
func (l *AzureEtcdDataEncryptionList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AzureEtcdDataEncryptionList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AzureEtcdDataEncryptionList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AzureEtcdDataEncryptionList) SetItems(items []*AzureEtcdDataEncryption) {
	l.items = items
}

// Items returns the items of the list.
func (l *AzureEtcdDataEncryptionList) Items() []*AzureEtcdDataEncryption {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AzureEtcdDataEncryptionList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AzureEtcdDataEncryptionList) Get(i int) *AzureEtcdDataEncryption {
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
func (l *AzureEtcdDataEncryptionList) Slice() []*AzureEtcdDataEncryption {
	var slice []*AzureEtcdDataEncryption
	if l == nil {
		slice = make([]*AzureEtcdDataEncryption, 0)
	} else {
		slice = make([]*AzureEtcdDataEncryption, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AzureEtcdDataEncryptionList) Each(f func(item *AzureEtcdDataEncryption) bool) {
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
func (l *AzureEtcdDataEncryptionList) Range(f func(index int, item *AzureEtcdDataEncryption) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
