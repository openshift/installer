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

// CloudRegionKind is the name of the type used to represent objects
// of type 'cloud_region'.
const CloudRegionKind = "CloudRegion"

// CloudRegionLinkKind is the name of the type used to represent links
// to objects of type 'cloud_region'.
const CloudRegionLinkKind = "CloudRegionLink"

// CloudRegionNilKind is the name of the type used to nil references
// to objects of type 'cloud_region'.
const CloudRegionNilKind = "CloudRegionNil"

// CloudRegion represents the values of the 'cloud_region' type.
//
// Description of a region of a cloud provider.
type CloudRegion struct {
	fieldSet_          []bool
	id                 string
	href               string
	kmsLocationID      string
	kmsLocationName    string
	cloudProvider      *CloudProvider
	displayName        string
	name               string
	ccsOnly            bool
	enabled            bool
	govCloud           bool
	supportsHypershift bool
	supportsMultiAZ    bool
}

// Kind returns the name of the type of the object.
func (o *CloudRegion) Kind() string {
	if o == nil {
		return CloudRegionNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return CloudRegionLinkKind
	}
	return CloudRegionKind
}

// Link returns true if this is a link.
func (o *CloudRegion) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *CloudRegion) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *CloudRegion) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *CloudRegion) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *CloudRegion) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *CloudRegion) Empty() bool {
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

// CCSOnly returns the value of the 'CCS_only' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// 'true' if the region is supported only for CCS clusters, 'false' otherwise.
func (o *CloudRegion) CCSOnly() bool {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.ccsOnly
	}
	return false
}

// GetCCSOnly returns the value of the 'CCS_only' attribute and
// a flag indicating if the attribute has a value.
//
// 'true' if the region is supported only for CCS clusters, 'false' otherwise.
func (o *CloudRegion) GetCCSOnly() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.ccsOnly
	}
	return
}

// KMSLocationID returns the value of the 'KMS_location_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// (GCP only) Comma-separated list of KMS location IDs that can be used with this region.
// E.g. "global,nam4,us". Order is not guaranteed.
func (o *CloudRegion) KMSLocationID() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.kmsLocationID
	}
	return ""
}

// GetKMSLocationID returns the value of the 'KMS_location_ID' attribute and
// a flag indicating if the attribute has a value.
//
// (GCP only) Comma-separated list of KMS location IDs that can be used with this region.
// E.g. "global,nam4,us". Order is not guaranteed.
func (o *CloudRegion) GetKMSLocationID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.kmsLocationID
	}
	return
}

// KMSLocationName returns the value of the 'KMS_location_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// (GCP only) Comma-separated list of display names corresponding to KMSLocationID.
// E.g. "Global,nam4 (Iowa, South Carolina, and Oklahoma),US". Order is not guaranteed but will match KMSLocationID.
// Unfortunately, this API doesn't allow robust splitting - Contact ocm-feedback@redhat.com if you want to rely on this.
func (o *CloudRegion) KMSLocationName() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.kmsLocationName
	}
	return ""
}

// GetKMSLocationName returns the value of the 'KMS_location_name' attribute and
// a flag indicating if the attribute has a value.
//
// (GCP only) Comma-separated list of display names corresponding to KMSLocationID.
// E.g. "Global,nam4 (Iowa, South Carolina, and Oklahoma),US". Order is not guaranteed but will match KMSLocationID.
// Unfortunately, this API doesn't allow robust splitting - Contact ocm-feedback@redhat.com if you want to rely on this.
func (o *CloudRegion) GetKMSLocationName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.kmsLocationName
	}
	return
}

// CloudProvider returns the value of the 'cloud_provider' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the cloud provider that the region belongs to.
func (o *CloudRegion) CloudProvider() *CloudProvider {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.cloudProvider
	}
	return nil
}

// GetCloudProvider returns the value of the 'cloud_provider' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the cloud provider that the region belongs to.
func (o *CloudRegion) GetCloudProvider() (value *CloudProvider, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.cloudProvider
	}
	return
}

// DisplayName returns the value of the 'display_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the region for display purposes, for example `N. Virginia`.
func (o *CloudRegion) DisplayName() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.displayName
	}
	return ""
}

// GetDisplayName returns the value of the 'display_name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the region for display purposes, for example `N. Virginia`.
func (o *CloudRegion) GetDisplayName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.displayName
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Whether the region is enabled for deploying a managed cluster.
func (o *CloudRegion) Enabled() bool {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Whether the region is enabled for deploying a managed cluster.
func (o *CloudRegion) GetEnabled() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.enabled
	}
	return
}

// GovCloud returns the value of the 'gov_cloud' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Whether the region is an AWS GovCloud region.
func (o *CloudRegion) GovCloud() bool {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.govCloud
	}
	return false
}

// GetGovCloud returns the value of the 'gov_cloud' attribute and
// a flag indicating if the attribute has a value.
//
// Whether the region is an AWS GovCloud region.
func (o *CloudRegion) GetGovCloud() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.govCloud
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Human friendly identifier of the region, for example `us-east-1`.
//
// NOTE: Currently for all cloud providers and all regions `id` and `name` have exactly
// the same values.
func (o *CloudRegion) Name() string {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Human friendly identifier of the region, for example `us-east-1`.
//
// NOTE: Currently for all cloud providers and all regions `id` and `name` have exactly
// the same values.
func (o *CloudRegion) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.name
	}
	return
}

// SupportsHypershift returns the value of the 'supports_hypershift' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// 'true' if the region is supported for Hypershift deployments, 'false' otherwise.
func (o *CloudRegion) SupportsHypershift() bool {
	if o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11] {
		return o.supportsHypershift
	}
	return false
}

// GetSupportsHypershift returns the value of the 'supports_hypershift' attribute and
// a flag indicating if the attribute has a value.
//
// 'true' if the region is supported for Hypershift deployments, 'false' otherwise.
func (o *CloudRegion) GetSupportsHypershift() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11]
	if ok {
		value = o.supportsHypershift
	}
	return
}

// SupportsMultiAZ returns the value of the 'supports_multi_AZ' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Whether the region supports multiple availability zones.
func (o *CloudRegion) SupportsMultiAZ() bool {
	if o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12] {
		return o.supportsMultiAZ
	}
	return false
}

// GetSupportsMultiAZ returns the value of the 'supports_multi_AZ' attribute and
// a flag indicating if the attribute has a value.
//
// Whether the region supports multiple availability zones.
func (o *CloudRegion) GetSupportsMultiAZ() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12]
	if ok {
		value = o.supportsMultiAZ
	}
	return
}

// CloudRegionListKind is the name of the type used to represent list of objects of
// type 'cloud_region'.
const CloudRegionListKind = "CloudRegionList"

// CloudRegionListLinkKind is the name of the type used to represent links to list
// of objects of type 'cloud_region'.
const CloudRegionListLinkKind = "CloudRegionListLink"

// CloudRegionNilKind is the name of the type used to nil lists of objects of
// type 'cloud_region'.
const CloudRegionListNilKind = "CloudRegionListNil"

// CloudRegionList is a list of values of the 'cloud_region' type.
type CloudRegionList struct {
	href  string
	link  bool
	items []*CloudRegion
}

// Kind returns the name of the type of the object.
func (l *CloudRegionList) Kind() string {
	if l == nil {
		return CloudRegionListNilKind
	}
	if l.link {
		return CloudRegionListLinkKind
	}
	return CloudRegionListKind
}

// Link returns true iif this is a link.
func (l *CloudRegionList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *CloudRegionList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *CloudRegionList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *CloudRegionList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *CloudRegionList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *CloudRegionList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *CloudRegionList) SetItems(items []*CloudRegion) {
	l.items = items
}

// Items returns the items of the list.
func (l *CloudRegionList) Items() []*CloudRegion {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *CloudRegionList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *CloudRegionList) Get(i int) *CloudRegion {
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
func (l *CloudRegionList) Slice() []*CloudRegion {
	var slice []*CloudRegion
	if l == nil {
		slice = make([]*CloudRegion, 0)
	} else {
		slice = make([]*CloudRegion, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *CloudRegionList) Each(f func(item *CloudRegion) bool) {
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
func (l *CloudRegionList) Range(f func(index int, item *CloudRegion) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
