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

// OidcThumbprint represents the values of the 'oidc_thumbprint' type.
//
// Contains the necessary attributes to support oidc configuration thumbprint operations such as fetching/creation of a thumbprint
type OidcThumbprint struct {
	bitmap_      uint32
	href         string
	clusterId    string
	kind         string
	oidcConfigId string
	thumbprint   string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *OidcThumbprint) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// HREF returns the value of the 'HREF' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// HREF for the oidc config thumbprint, filled in response.
func (o *OidcThumbprint) HREF() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the value of the 'HREF' attribute and
// a flag indicating if the attribute has a value.
//
// HREF for the oidc config thumbprint, filled in response.
func (o *OidcThumbprint) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.href
	}
	return
}

// ClusterId returns the value of the 'cluster_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ClusterId is the for the cluster used, filled in response.
func (o *OidcThumbprint) ClusterId() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.clusterId
	}
	return ""
}

// GetClusterId returns the value of the 'cluster_id' attribute and
// a flag indicating if the attribute has a value.
//
// ClusterId is the for the cluster used, filled in response.
func (o *OidcThumbprint) GetClusterId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.clusterId
	}
	return
}

// Kind returns the value of the 'kind' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Kind is the resource type, filled in response.
func (o *OidcThumbprint) Kind() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.kind
	}
	return ""
}

// GetKind returns the value of the 'kind' attribute and
// a flag indicating if the attribute has a value.
//
// Kind is the resource type, filled in response.
func (o *OidcThumbprint) GetKind() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.kind
	}
	return
}

// OidcConfigId returns the value of the 'oidc_config_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// OidcConfigId is the ID for the oidc config used, filled in response.
func (o *OidcThumbprint) OidcConfigId() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.oidcConfigId
	}
	return ""
}

// GetOidcConfigId returns the value of the 'oidc_config_id' attribute and
// a flag indicating if the attribute has a value.
//
// OidcConfigId is the ID for the oidc config used, filled in response.
func (o *OidcThumbprint) GetOidcConfigId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.oidcConfigId
	}
	return
}

// Thumbprint returns the value of the 'thumbprint' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Thumbprint is the thumbprint itself, filled in response.
func (o *OidcThumbprint) Thumbprint() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.thumbprint
	}
	return ""
}

// GetThumbprint returns the value of the 'thumbprint' attribute and
// a flag indicating if the attribute has a value.
//
// Thumbprint is the thumbprint itself, filled in response.
func (o *OidcThumbprint) GetThumbprint() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.thumbprint
	}
	return
}

// OidcThumbprintListKind is the name of the type used to represent list of objects of
// type 'oidc_thumbprint'.
const OidcThumbprintListKind = "OidcThumbprintList"

// OidcThumbprintListLinkKind is the name of the type used to represent links to list
// of objects of type 'oidc_thumbprint'.
const OidcThumbprintListLinkKind = "OidcThumbprintListLink"

// OidcThumbprintNilKind is the name of the type used to nil lists of objects of
// type 'oidc_thumbprint'.
const OidcThumbprintListNilKind = "OidcThumbprintListNil"

// OidcThumbprintList is a list of values of the 'oidc_thumbprint' type.
type OidcThumbprintList struct {
	href  string
	link  bool
	items []*OidcThumbprint
}

// Len returns the length of the list.
func (l *OidcThumbprintList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *OidcThumbprintList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *OidcThumbprintList) Get(i int) *OidcThumbprint {
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
func (l *OidcThumbprintList) Slice() []*OidcThumbprint {
	var slice []*OidcThumbprint
	if l == nil {
		slice = make([]*OidcThumbprint, 0)
	} else {
		slice = make([]*OidcThumbprint, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *OidcThumbprintList) Each(f func(item *OidcThumbprint) bool) {
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
func (l *OidcThumbprintList) Range(f func(index int, item *OidcThumbprint) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
