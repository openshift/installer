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

// CloudProviderData represents the values of the 'cloud_provider_data' type.
//
// Description of a cloud provider data used for cloud provider inquiries.
type CloudProviderData struct {
	bitmap_           uint32
	aws               *AWS
	gcp               *GCP
	availabilityZones []string
	keyLocation       string
	keyRingName       string
	region            *CloudRegion
	subnets           []string
	version           *Version
	vpcIds            []string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *CloudProviderData) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AWS returns the value of the 'AWS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Amazon Web Services settings.
func (o *CloudProviderData) AWS() *AWS {
	if o != nil && o.bitmap_&1 != 0 {
		return o.aws
	}
	return nil
}

// GetAWS returns the value of the 'AWS' attribute and
// a flag indicating if the attribute has a value.
//
// Amazon Web Services settings.
func (o *CloudProviderData) GetAWS() (value *AWS, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.aws
	}
	return
}

// GCP returns the value of the 'GCP' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Google cloud platform settings.
func (o *CloudProviderData) GCP() *GCP {
	if o != nil && o.bitmap_&2 != 0 {
		return o.gcp
	}
	return nil
}

// GetGCP returns the value of the 'GCP' attribute and
// a flag indicating if the attribute has a value.
//
// Google cloud platform settings.
func (o *CloudProviderData) GetGCP() (value *GCP, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.gcp
	}
	return
}

// AvailabilityZones returns the value of the 'availability_zones' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Availability zone
func (o *CloudProviderData) AvailabilityZones() []string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.availabilityZones
	}
	return nil
}

// GetAvailabilityZones returns the value of the 'availability_zones' attribute and
// a flag indicating if the attribute has a value.
//
// Availability zone
func (o *CloudProviderData) GetAvailabilityZones() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.availabilityZones
	}
	return
}

// KeyLocation returns the value of the 'key_location' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Key location
func (o *CloudProviderData) KeyLocation() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.keyLocation
	}
	return ""
}

// GetKeyLocation returns the value of the 'key_location' attribute and
// a flag indicating if the attribute has a value.
//
// Key location
func (o *CloudProviderData) GetKeyLocation() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.keyLocation
	}
	return
}

// KeyRingName returns the value of the 'key_ring_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Key ring name
func (o *CloudProviderData) KeyRingName() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.keyRingName
	}
	return ""
}

// GetKeyRingName returns the value of the 'key_ring_name' attribute and
// a flag indicating if the attribute has a value.
//
// Key ring name
func (o *CloudProviderData) GetKeyRingName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.keyRingName
	}
	return
}

// Region returns the value of the 'region' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Region
func (o *CloudProviderData) Region() *CloudRegion {
	if o != nil && o.bitmap_&32 != 0 {
		return o.region
	}
	return nil
}

// GetRegion returns the value of the 'region' attribute and
// a flag indicating if the attribute has a value.
//
// Region
func (o *CloudProviderData) GetRegion() (value *CloudRegion, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.region
	}
	return
}

// Subnets returns the value of the 'subnets' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Subnets
func (o *CloudProviderData) Subnets() []string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.subnets
	}
	return nil
}

// GetSubnets returns the value of the 'subnets' attribute and
// a flag indicating if the attribute has a value.
//
// Subnets
func (o *CloudProviderData) GetSubnets() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.subnets
	}
	return
}

// Version returns the value of the 'version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Openshift version
func (o *CloudProviderData) Version() *Version {
	if o != nil && o.bitmap_&128 != 0 {
		return o.version
	}
	return nil
}

// GetVersion returns the value of the 'version' attribute and
// a flag indicating if the attribute has a value.
//
// Openshift version
func (o *CloudProviderData) GetVersion() (value *Version, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.version
	}
	return
}

// VpcIds returns the value of the 'vpc_ids' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// VPC ids
func (o *CloudProviderData) VpcIds() []string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.vpcIds
	}
	return nil
}

// GetVpcIds returns the value of the 'vpc_ids' attribute and
// a flag indicating if the attribute has a value.
//
// VPC ids
func (o *CloudProviderData) GetVpcIds() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.vpcIds
	}
	return
}

// CloudProviderDataListKind is the name of the type used to represent list of objects of
// type 'cloud_provider_data'.
const CloudProviderDataListKind = "CloudProviderDataList"

// CloudProviderDataListLinkKind is the name of the type used to represent links to list
// of objects of type 'cloud_provider_data'.
const CloudProviderDataListLinkKind = "CloudProviderDataListLink"

// CloudProviderDataNilKind is the name of the type used to nil lists of objects of
// type 'cloud_provider_data'.
const CloudProviderDataListNilKind = "CloudProviderDataListNil"

// CloudProviderDataList is a list of values of the 'cloud_provider_data' type.
type CloudProviderDataList struct {
	href  string
	link  bool
	items []*CloudProviderData
}

// Len returns the length of the list.
func (l *CloudProviderDataList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *CloudProviderDataList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *CloudProviderDataList) Get(i int) *CloudProviderData {
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
func (l *CloudProviderDataList) Slice() []*CloudProviderData {
	var slice []*CloudProviderData
	if l == nil {
		slice = make([]*CloudProviderData, 0)
	} else {
		slice = make([]*CloudProviderData, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *CloudProviderDataList) Each(f func(item *CloudProviderData) bool) {
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
func (l *CloudProviderDataList) Range(f func(index int, item *CloudProviderData) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
