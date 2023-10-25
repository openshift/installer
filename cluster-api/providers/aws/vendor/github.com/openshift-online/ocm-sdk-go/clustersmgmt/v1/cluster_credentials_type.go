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

// ClusterCredentialsKind is the name of the type used to represent objects
// of type 'cluster_credentials'.
const ClusterCredentialsKind = "ClusterCredentials"

// ClusterCredentialsLinkKind is the name of the type used to represent links
// to objects of type 'cluster_credentials'.
const ClusterCredentialsLinkKind = "ClusterCredentialsLink"

// ClusterCredentialsNilKind is the name of the type used to nil references
// to objects of type 'cluster_credentials'.
const ClusterCredentialsNilKind = "ClusterCredentialsNil"

// ClusterCredentials represents the values of the 'cluster_credentials' type.
//
// Credentials of the a cluster.
type ClusterCredentials struct {
	bitmap_    uint32
	id         string
	href       string
	kubeconfig string
}

// Kind returns the name of the type of the object.
func (o *ClusterCredentials) Kind() string {
	if o == nil {
		return ClusterCredentialsNilKind
	}
	if o.bitmap_&1 != 0 {
		return ClusterCredentialsLinkKind
	}
	return ClusterCredentialsKind
}

// Link returns true iif this is a link.
func (o *ClusterCredentials) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *ClusterCredentials) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ClusterCredentials) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ClusterCredentials) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ClusterCredentials) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterCredentials) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Kubeconfig returns the value of the 'kubeconfig' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Administrator _kubeconfig_ file for the cluster.
func (o *ClusterCredentials) Kubeconfig() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.kubeconfig
	}
	return ""
}

// GetKubeconfig returns the value of the 'kubeconfig' attribute and
// a flag indicating if the attribute has a value.
//
// Administrator _kubeconfig_ file for the cluster.
func (o *ClusterCredentials) GetKubeconfig() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.kubeconfig
	}
	return
}

// ClusterCredentialsListKind is the name of the type used to represent list of objects of
// type 'cluster_credentials'.
const ClusterCredentialsListKind = "ClusterCredentialsList"

// ClusterCredentialsListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_credentials'.
const ClusterCredentialsListLinkKind = "ClusterCredentialsListLink"

// ClusterCredentialsNilKind is the name of the type used to nil lists of objects of
// type 'cluster_credentials'.
const ClusterCredentialsListNilKind = "ClusterCredentialsListNil"

// ClusterCredentialsList is a list of values of the 'cluster_credentials' type.
type ClusterCredentialsList struct {
	href  string
	link  bool
	items []*ClusterCredentials
}

// Kind returns the name of the type of the object.
func (l *ClusterCredentialsList) Kind() string {
	if l == nil {
		return ClusterCredentialsListNilKind
	}
	if l.link {
		return ClusterCredentialsListLinkKind
	}
	return ClusterCredentialsListKind
}

// Link returns true iif this is a link.
func (l *ClusterCredentialsList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ClusterCredentialsList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ClusterCredentialsList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ClusterCredentialsList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ClusterCredentialsList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterCredentialsList) Get(i int) *ClusterCredentials {
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
func (l *ClusterCredentialsList) Slice() []*ClusterCredentials {
	var slice []*ClusterCredentials
	if l == nil {
		slice = make([]*ClusterCredentials, 0)
	} else {
		slice = make([]*ClusterCredentials, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterCredentialsList) Each(f func(item *ClusterCredentials) bool) {
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
func (l *ClusterCredentialsList) Range(f func(index int, item *ClusterCredentials) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
