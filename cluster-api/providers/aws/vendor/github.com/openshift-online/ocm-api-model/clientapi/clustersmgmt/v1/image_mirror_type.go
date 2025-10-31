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

import (
	time "time"
)

// ImageMirrorKind is the name of the type used to represent objects
// of type 'image_mirror'.
const ImageMirrorKind = "ImageMirror"

// ImageMirrorLinkKind is the name of the type used to represent links
// to objects of type 'image_mirror'.
const ImageMirrorLinkKind = "ImageMirrorLink"

// ImageMirrorNilKind is the name of the type used to nil references
// to objects of type 'image_mirror'.
const ImageMirrorNilKind = "ImageMirrorNil"

// ImageMirror represents the values of the 'image_mirror' type.
//
// ImageMirror represents a container image mirror configuration for a cluster.
// This enables Day 2 image mirroring configuration for ROSA HCP clusters using
// HyperShift's native imageContentSources mechanism.
type ImageMirror struct {
	fieldSet_           []bool
	id                  string
	href                string
	creationTimestamp   time.Time
	lastUpdateTimestamp time.Time
	mirrors             []string
	source              string
	type_               string
}

// Kind returns the name of the type of the object.
func (o *ImageMirror) Kind() string {
	if o == nil {
		return ImageMirrorNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ImageMirrorLinkKind
	}
	return ImageMirrorKind
}

// Link returns true if this is a link.
func (o *ImageMirror) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *ImageMirror) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ImageMirror) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ImageMirror) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ImageMirror) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ImageMirror) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}

	// Check all fields except the link flag (index 0)
	for i := 1; i < len(o.fieldSet_); i++ {
		if o.fieldSet_[i] {
			return false
		}
	}
	return true
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// CreationTimestamp indicates when the image mirror was created.
func (o *ImageMirror) CreationTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// CreationTimestamp indicates when the image mirror was created.
func (o *ImageMirror) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.creationTimestamp
	}
	return
}

// LastUpdateTimestamp returns the value of the 'last_update_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// LastUpdateTimestamp indicates when the image mirror was last updated.
func (o *ImageMirror) LastUpdateTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.lastUpdateTimestamp
	}
	return time.Time{}
}

// GetLastUpdateTimestamp returns the value of the 'last_update_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// LastUpdateTimestamp indicates when the image mirror was last updated.
func (o *ImageMirror) GetLastUpdateTimestamp() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.lastUpdateTimestamp
	}
	return
}

// Mirrors returns the value of the 'mirrors' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Mirrors is the list of mirror registries that will serve content for the source.
// Mirrors array cannot be empty (must contain at least one mirror registry).
// Each mirror registry URL must conform to OpenShift's ImageDigestMirrorSet format specifications.
func (o *ImageMirror) Mirrors() []string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.mirrors
	}
	return nil
}

// GetMirrors returns the value of the 'mirrors' attribute and
// a flag indicating if the attribute has a value.
//
// Mirrors is the list of mirror registries that will serve content for the source.
// Mirrors array cannot be empty (must contain at least one mirror registry).
// Each mirror registry URL must conform to OpenShift's ImageDigestMirrorSet format specifications.
func (o *ImageMirror) GetMirrors() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.mirrors
	}
	return
}

// Source returns the value of the 'source' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Source is the source registry that will be mirrored.
// Source registry must be unique per cluster and is immutable after creation.
// Source is used to identify mirror entries in HostedCluster imageContentSources.
func (o *ImageMirror) Source() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.source
	}
	return ""
}

// GetSource returns the value of the 'source' attribute and
// a flag indicating if the attribute has a value.
//
// Source is the source registry that will be mirrored.
// Source registry must be unique per cluster and is immutable after creation.
// Source is used to identify mirror entries in HostedCluster imageContentSources.
func (o *ImageMirror) GetSource() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.source
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Type specifies the mirror type, currently only "digest" is supported.
func (o *ImageMirror) Type() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.type_
	}
	return ""
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
//
// Type specifies the mirror type, currently only "digest" is supported.
func (o *ImageMirror) GetType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.type_
	}
	return
}

// ImageMirrorListKind is the name of the type used to represent list of objects of
// type 'image_mirror'.
const ImageMirrorListKind = "ImageMirrorList"

// ImageMirrorListLinkKind is the name of the type used to represent links to list
// of objects of type 'image_mirror'.
const ImageMirrorListLinkKind = "ImageMirrorListLink"

// ImageMirrorNilKind is the name of the type used to nil lists of objects of
// type 'image_mirror'.
const ImageMirrorListNilKind = "ImageMirrorListNil"

// ImageMirrorList is a list of values of the 'image_mirror' type.
type ImageMirrorList struct {
	href  string
	link  bool
	items []*ImageMirror
}

// Kind returns the name of the type of the object.
func (l *ImageMirrorList) Kind() string {
	if l == nil {
		return ImageMirrorListNilKind
	}
	if l.link {
		return ImageMirrorListLinkKind
	}
	return ImageMirrorListKind
}

// Link returns true iif this is a link.
func (l *ImageMirrorList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ImageMirrorList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ImageMirrorList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ImageMirrorList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ImageMirrorList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ImageMirrorList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ImageMirrorList) SetItems(items []*ImageMirror) {
	l.items = items
}

// Items returns the items of the list.
func (l *ImageMirrorList) Items() []*ImageMirror {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ImageMirrorList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ImageMirrorList) Get(i int) *ImageMirror {
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
func (l *ImageMirrorList) Slice() []*ImageMirror {
	var slice []*ImageMirror
	if l == nil {
		slice = make([]*ImageMirror, 0)
	} else {
		slice = make([]*ImageMirror, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ImageMirrorList) Each(f func(item *ImageMirror) bool) {
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
func (l *ImageMirrorList) Range(f func(index int, item *ImageMirror) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
