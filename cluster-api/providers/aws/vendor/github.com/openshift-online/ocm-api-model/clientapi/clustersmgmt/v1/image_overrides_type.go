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

// ImageOverridesKind is the name of the type used to represent objects
// of type 'image_overrides'.
const ImageOverridesKind = "ImageOverrides"

// ImageOverridesLinkKind is the name of the type used to represent links
// to objects of type 'image_overrides'.
const ImageOverridesLinkKind = "ImageOverridesLink"

// ImageOverridesNilKind is the name of the type used to nil references
// to objects of type 'image_overrides'.
const ImageOverridesNilKind = "ImageOverridesNil"

// ImageOverrides represents the values of the 'image_overrides' type.
//
// ImageOverrides holds the lists of available images per cloud provider.
type ImageOverrides struct {
	fieldSet_ []bool
	id        string
	href      string
	aws       []*AMIOverride
	gcp       []*GCPImageOverride
}

// Kind returns the name of the type of the object.
func (o *ImageOverrides) Kind() string {
	if o == nil {
		return ImageOverridesNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ImageOverridesLinkKind
	}
	return ImageOverridesKind
}

// Link returns true if this is a link.
func (o *ImageOverrides) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *ImageOverrides) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ImageOverrides) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ImageOverrides) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ImageOverrides) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ImageOverrides) Empty() bool {
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

// AWS returns the value of the 'AWS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ImageOverrides) AWS() []*AMIOverride {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.aws
	}
	return nil
}

// GetAWS returns the value of the 'AWS' attribute and
// a flag indicating if the attribute has a value.
func (o *ImageOverrides) GetAWS() (value []*AMIOverride, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.aws
	}
	return
}

// GCP returns the value of the 'GCP' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ImageOverrides) GCP() []*GCPImageOverride {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.gcp
	}
	return nil
}

// GetGCP returns the value of the 'GCP' attribute and
// a flag indicating if the attribute has a value.
func (o *ImageOverrides) GetGCP() (value []*GCPImageOverride, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.gcp
	}
	return
}

// ImageOverridesListKind is the name of the type used to represent list of objects of
// type 'image_overrides'.
const ImageOverridesListKind = "ImageOverridesList"

// ImageOverridesListLinkKind is the name of the type used to represent links to list
// of objects of type 'image_overrides'.
const ImageOverridesListLinkKind = "ImageOverridesListLink"

// ImageOverridesNilKind is the name of the type used to nil lists of objects of
// type 'image_overrides'.
const ImageOverridesListNilKind = "ImageOverridesListNil"

// ImageOverridesList is a list of values of the 'image_overrides' type.
type ImageOverridesList struct {
	href  string
	link  bool
	items []*ImageOverrides
}

// Kind returns the name of the type of the object.
func (l *ImageOverridesList) Kind() string {
	if l == nil {
		return ImageOverridesListNilKind
	}
	if l.link {
		return ImageOverridesListLinkKind
	}
	return ImageOverridesListKind
}

// Link returns true iif this is a link.
func (l *ImageOverridesList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ImageOverridesList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ImageOverridesList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ImageOverridesList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ImageOverridesList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ImageOverridesList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ImageOverridesList) SetItems(items []*ImageOverrides) {
	l.items = items
}

// Items returns the items of the list.
func (l *ImageOverridesList) Items() []*ImageOverrides {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ImageOverridesList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ImageOverridesList) Get(i int) *ImageOverrides {
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
func (l *ImageOverridesList) Slice() []*ImageOverrides {
	var slice []*ImageOverrides
	if l == nil {
		slice = make([]*ImageOverrides, 0)
	} else {
		slice = make([]*ImageOverrides, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ImageOverridesList) Each(f func(item *ImageOverrides) bool) {
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
func (l *ImageOverridesList) Range(f func(index int, item *ImageOverrides) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
