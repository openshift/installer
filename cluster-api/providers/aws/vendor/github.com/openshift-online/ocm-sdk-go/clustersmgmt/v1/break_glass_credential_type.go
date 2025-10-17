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

import (
	time "time"
)

// BreakGlassCredentialKind is the name of the type used to represent objects
// of type 'break_glass_credential'.
const BreakGlassCredentialKind = "BreakGlassCredential"

// BreakGlassCredentialLinkKind is the name of the type used to represent links
// to objects of type 'break_glass_credential'.
const BreakGlassCredentialLinkKind = "BreakGlassCredentialLink"

// BreakGlassCredentialNilKind is the name of the type used to nil references
// to objects of type 'break_glass_credential'.
const BreakGlassCredentialNilKind = "BreakGlassCredentialNil"

// BreakGlassCredential represents the values of the 'break_glass_credential' type.
//
// Representation of a break glass credential.
type BreakGlassCredential struct {
	bitmap_             uint32
	id                  string
	href                string
	expirationTimestamp time.Time
	kubeconfig          string
	revocationTimestamp time.Time
	status              BreakGlassCredentialStatus
	username            string
}

// Kind returns the name of the type of the object.
func (o *BreakGlassCredential) Kind() string {
	if o == nil {
		return BreakGlassCredentialNilKind
	}
	if o.bitmap_&1 != 0 {
		return BreakGlassCredentialLinkKind
	}
	return BreakGlassCredentialKind
}

// Link returns true if this is a link.
func (o *BreakGlassCredential) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *BreakGlassCredential) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *BreakGlassCredential) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *BreakGlassCredential) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *BreakGlassCredential) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *BreakGlassCredential) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// ExpirationTimestamp returns the value of the 'expiration_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ExpirationTimestamp is the date and time when the credential will expire.
func (o *BreakGlassCredential) ExpirationTimestamp() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.expirationTimestamp
	}
	return time.Time{}
}

// GetExpirationTimestamp returns the value of the 'expiration_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// ExpirationTimestamp is the date and time when the credential will expire.
func (o *BreakGlassCredential) GetExpirationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.expirationTimestamp
	}
	return
}

// Kubeconfig returns the value of the 'kubeconfig' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Kubeconfig is the generated kubeconfig for this credential. It is only stored in memory
func (o *BreakGlassCredential) Kubeconfig() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.kubeconfig
	}
	return ""
}

// GetKubeconfig returns the value of the 'kubeconfig' attribute and
// a flag indicating if the attribute has a value.
//
// Kubeconfig is the generated kubeconfig for this credential. It is only stored in memory
func (o *BreakGlassCredential) GetKubeconfig() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.kubeconfig
	}
	return
}

// RevocationTimestamp returns the value of the 'revocation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// RevocationTimestamp is the date and time when the credential has been revoked.
func (o *BreakGlassCredential) RevocationTimestamp() time.Time {
	if o != nil && o.bitmap_&32 != 0 {
		return o.revocationTimestamp
	}
	return time.Time{}
}

// GetRevocationTimestamp returns the value of the 'revocation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// RevocationTimestamp is the date and time when the credential has been revoked.
func (o *BreakGlassCredential) GetRevocationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.revocationTimestamp
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Status is the status of this credential
func (o *BreakGlassCredential) Status() BreakGlassCredentialStatus {
	if o != nil && o.bitmap_&64 != 0 {
		return o.status
	}
	return BreakGlassCredentialStatus("")
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// Status is the status of this credential
func (o *BreakGlassCredential) GetStatus() (value BreakGlassCredentialStatus, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.status
	}
	return
}

// Username returns the value of the 'username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Username is the user which will be used for this credential.
func (o *BreakGlassCredential) Username() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.username
	}
	return ""
}

// GetUsername returns the value of the 'username' attribute and
// a flag indicating if the attribute has a value.
//
// Username is the user which will be used for this credential.
func (o *BreakGlassCredential) GetUsername() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.username
	}
	return
}

// BreakGlassCredentialListKind is the name of the type used to represent list of objects of
// type 'break_glass_credential'.
const BreakGlassCredentialListKind = "BreakGlassCredentialList"

// BreakGlassCredentialListLinkKind is the name of the type used to represent links to list
// of objects of type 'break_glass_credential'.
const BreakGlassCredentialListLinkKind = "BreakGlassCredentialListLink"

// BreakGlassCredentialNilKind is the name of the type used to nil lists of objects of
// type 'break_glass_credential'.
const BreakGlassCredentialListNilKind = "BreakGlassCredentialListNil"

// BreakGlassCredentialList is a list of values of the 'break_glass_credential' type.
type BreakGlassCredentialList struct {
	href  string
	link  bool
	items []*BreakGlassCredential
}

// Kind returns the name of the type of the object.
func (l *BreakGlassCredentialList) Kind() string {
	if l == nil {
		return BreakGlassCredentialListNilKind
	}
	if l.link {
		return BreakGlassCredentialListLinkKind
	}
	return BreakGlassCredentialListKind
}

// Link returns true iif this is a link.
func (l *BreakGlassCredentialList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *BreakGlassCredentialList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *BreakGlassCredentialList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *BreakGlassCredentialList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *BreakGlassCredentialList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *BreakGlassCredentialList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *BreakGlassCredentialList) SetItems(items []*BreakGlassCredential) {
	l.items = items
}

// Items returns the items of the list.
func (l *BreakGlassCredentialList) Items() []*BreakGlassCredential {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *BreakGlassCredentialList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *BreakGlassCredentialList) Get(i int) *BreakGlassCredential {
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
func (l *BreakGlassCredentialList) Slice() []*BreakGlassCredential {
	var slice []*BreakGlassCredential
	if l == nil {
		slice = make([]*BreakGlassCredential, 0)
	} else {
		slice = make([]*BreakGlassCredential, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *BreakGlassCredentialList) Each(f func(item *BreakGlassCredential) bool) {
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
func (l *BreakGlassCredentialList) Range(f func(index int, item *BreakGlassCredential) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
