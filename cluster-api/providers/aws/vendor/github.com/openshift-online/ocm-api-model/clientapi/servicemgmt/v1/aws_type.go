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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/servicemgmt/v1

// AWS represents the values of the 'AWS' type.
//
// _Amazon Web Services_ specific settings of a cluster.
type AWS struct {
	fieldSet_       []bool
	sts             *STS
	accessKeyID     string
	accountID       string
	secretAccessKey string
	subnetIDs       []string
	tags            map[string]string
	privateLink     bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWS) Empty() bool {
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

// STS returns the value of the 'STS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains the necessary attributes to support role-based authentication on AWS.
func (o *AWS) STS() *STS {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.sts
	}
	return nil
}

// GetSTS returns the value of the 'STS' attribute and
// a flag indicating if the attribute has a value.
//
// Contains the necessary attributes to support role-based authentication on AWS.
func (o *AWS) GetSTS() (value *STS, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.sts
	}
	return
}

// AccessKeyID returns the value of the 'access_key_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AWS access key identifier.
func (o *AWS) AccessKeyID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.accessKeyID
	}
	return ""
}

// GetAccessKeyID returns the value of the 'access_key_ID' attribute and
// a flag indicating if the attribute has a value.
//
// AWS access key identifier.
func (o *AWS) GetAccessKeyID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.accessKeyID
	}
	return
}

// AccountID returns the value of the 'account_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AWS account identifier.
func (o *AWS) AccountID() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.accountID
	}
	return ""
}

// GetAccountID returns the value of the 'account_ID' attribute and
// a flag indicating if the attribute has a value.
//
// AWS account identifier.
func (o *AWS) GetAccountID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.accountID
	}
	return
}

// PrivateLink returns the value of the 'private_link' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// For PrivateLink-enabled clusters
func (o *AWS) PrivateLink() bool {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.privateLink
	}
	return false
}

// GetPrivateLink returns the value of the 'private_link' attribute and
// a flag indicating if the attribute has a value.
//
// For PrivateLink-enabled clusters
func (o *AWS) GetPrivateLink() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.privateLink
	}
	return
}

// SecretAccessKey returns the value of the 'secret_access_key' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AWS secret access key.
func (o *AWS) SecretAccessKey() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.secretAccessKey
	}
	return ""
}

// GetSecretAccessKey returns the value of the 'secret_access_key' attribute and
// a flag indicating if the attribute has a value.
//
// AWS secret access key.
func (o *AWS) GetSecretAccessKey() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.secretAccessKey
	}
	return
}

// SubnetIDs returns the value of the 'subnet_IDs' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The subnet ids to be used when installing the cluster.
func (o *AWS) SubnetIDs() []string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.subnetIDs
	}
	return nil
}

// GetSubnetIDs returns the value of the 'subnet_IDs' attribute and
// a flag indicating if the attribute has a value.
//
// The subnet ids to be used when installing the cluster.
func (o *AWS) GetSubnetIDs() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.subnetIDs
	}
	return
}

// Tags returns the value of the 'tags' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional keys and values that the installer will add as tags to all AWS resources it creates
func (o *AWS) Tags() map[string]string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.tags
	}
	return nil
}

// GetTags returns the value of the 'tags' attribute and
// a flag indicating if the attribute has a value.
//
// Optional keys and values that the installer will add as tags to all AWS resources it creates
func (o *AWS) GetTags() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.tags
	}
	return
}

// AWSListKind is the name of the type used to represent list of objects of
// type 'AWS'.
const AWSListKind = "AWSList"

// AWSListLinkKind is the name of the type used to represent links to list
// of objects of type 'AWS'.
const AWSListLinkKind = "AWSListLink"

// AWSNilKind is the name of the type used to nil lists of objects of
// type 'AWS'.
const AWSListNilKind = "AWSListNil"

// AWSList is a list of values of the 'AWS' type.
type AWSList struct {
	href  string
	link  bool
	items []*AWS
}

// Len returns the length of the list.
func (l *AWSList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AWSList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AWSList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AWSList) SetItems(items []*AWS) {
	l.items = items
}

// Items returns the items of the list.
func (l *AWSList) Items() []*AWS {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AWSList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AWSList) Get(i int) *AWS {
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
func (l *AWSList) Slice() []*AWS {
	var slice []*AWS
	if l == nil {
		slice = make([]*AWS, 0)
	} else {
		slice = make([]*AWS, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AWSList) Each(f func(item *AWS) bool) {
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
func (l *AWSList) Range(f func(index int, item *AWS) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
