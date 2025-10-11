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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	time "time"
)

// VersionKind is the name of the type used to represent objects
// of type 'version'.
const VersionKind = "Version"

// VersionLinkKind is the name of the type used to represent links
// to objects of type 'version'.
const VersionLinkKind = "VersionLink"

// VersionNilKind is the name of the type used to nil references
// to objects of type 'version'.
const VersionNilKind = "VersionNil"

// Version represents the values of the 'version' type.
//
// Representation of an _OpenShift_ version.
type Version struct {
	bitmap_                   uint32
	id                        string
	href                      string
	availableUpgrades         []string
	channelGroup              string
	endOfLifeTimestamp        time.Time
	imageOverrides            *ImageOverrides
	rawID                     string
	releaseImage              string
	releaseImages             *ReleaseImages
	gcpMarketplaceEnabled     bool
	rosaEnabled               bool
	default_                  bool
	enabled                   bool
	hostedControlPlaneDefault bool
	hostedControlPlaneEnabled bool
	wifEnabled                bool
}

// Kind returns the name of the type of the object.
func (o *Version) Kind() string {
	if o == nil {
		return VersionNilKind
	}
	if o.bitmap_&1 != 0 {
		return VersionLinkKind
	}
	return VersionKind
}

// Link returns true if this is a link.
func (o *Version) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Version) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Version) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Version) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Version) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Version) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// GCPMarketplaceEnabled returns the value of the 'GCP_marketplace_enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCPMarketplaceEnabled indicates if this version can be used to create GCP Marketplace clusters.
func (o *Version) GCPMarketplaceEnabled() bool {
	if o != nil && o.bitmap_&8 != 0 {
		return o.gcpMarketplaceEnabled
	}
	return false
}

// GetGCPMarketplaceEnabled returns the value of the 'GCP_marketplace_enabled' attribute and
// a flag indicating if the attribute has a value.
//
// GCPMarketplaceEnabled indicates if this version can be used to create GCP Marketplace clusters.
func (o *Version) GetGCPMarketplaceEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.gcpMarketplaceEnabled
	}
	return
}

// ROSAEnabled returns the value of the 'ROSA_enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ROSAEnabled indicates whether this version can be used to create ROSA clusters.
func (o *Version) ROSAEnabled() bool {
	if o != nil && o.bitmap_&16 != 0 {
		return o.rosaEnabled
	}
	return false
}

// GetROSAEnabled returns the value of the 'ROSA_enabled' attribute and
// a flag indicating if the attribute has a value.
//
// ROSAEnabled indicates whether this version can be used to create ROSA clusters.
func (o *Version) GetROSAEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.rosaEnabled
	}
	return
}

// AvailableUpgrades returns the value of the 'available_upgrades' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AvailableUpgrades is the list of versions this version can be upgraded to.
func (o *Version) AvailableUpgrades() []string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.availableUpgrades
	}
	return nil
}

// GetAvailableUpgrades returns the value of the 'available_upgrades' attribute and
// a flag indicating if the attribute has a value.
//
// AvailableUpgrades is the list of versions this version can be upgraded to.
func (o *Version) GetAvailableUpgrades() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.availableUpgrades
	}
	return
}

// ChannelGroup returns the value of the 'channel_group' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ChannelGroup is the name of the group where this image belongs.
// ChannelGroup is a mechanism to partition the images to different groups,
// each image belongs to only a single group.
func (o *Version) ChannelGroup() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.channelGroup
	}
	return ""
}

// GetChannelGroup returns the value of the 'channel_group' attribute and
// a flag indicating if the attribute has a value.
//
// ChannelGroup is the name of the group where this image belongs.
// ChannelGroup is a mechanism to partition the images to different groups,
// each image belongs to only a single group.
func (o *Version) GetChannelGroup() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.channelGroup
	}
	return
}

// Default returns the value of the 'default' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this should be selected as the default version when a cluster is created
// without specifying explicitly the version.
func (o *Version) Default() bool {
	if o != nil && o.bitmap_&128 != 0 {
		return o.default_
	}
	return false
}

// GetDefault returns the value of the 'default' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this should be selected as the default version when a cluster is created
// without specifying explicitly the version.
func (o *Version) GetDefault() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.default_
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this version can be used to create clusters.
func (o *Version) Enabled() bool {
	if o != nil && o.bitmap_&256 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this version can be used to create clusters.
func (o *Version) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.enabled
	}
	return
}

// EndOfLifeTimestamp returns the value of the 'end_of_life_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// EndOfLifeTimestamp is the date and time when the version will get to End of Life, using the
// format defined in https://www.ietf.org/rfc/rfc3339.txt[RC3339].
func (o *Version) EndOfLifeTimestamp() time.Time {
	if o != nil && o.bitmap_&512 != 0 {
		return o.endOfLifeTimestamp
	}
	return time.Time{}
}

// GetEndOfLifeTimestamp returns the value of the 'end_of_life_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// EndOfLifeTimestamp is the date and time when the version will get to End of Life, using the
// format defined in https://www.ietf.org/rfc/rfc3339.txt[RC3339].
func (o *Version) GetEndOfLifeTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.endOfLifeTimestamp
	}
	return
}

// HostedControlPlaneDefault returns the value of the 'hosted_control_plane_default' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// HostedControlPlaneDefault is a flag that indicates if this should be selected as the default version when a
// HCP cluster is created without specifying explicitly the version.
func (o *Version) HostedControlPlaneDefault() bool {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.hostedControlPlaneDefault
	}
	return false
}

// GetHostedControlPlaneDefault returns the value of the 'hosted_control_plane_default' attribute and
// a flag indicating if the attribute has a value.
//
// HostedControlPlaneDefault is a flag that indicates if this should be selected as the default version when a
// HCP cluster is created without specifying explicitly the version.
func (o *Version) GetHostedControlPlaneDefault() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.hostedControlPlaneDefault
	}
	return
}

// HostedControlPlaneEnabled returns the value of the 'hosted_control_plane_enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// HostedControlPlaneEnabled indicates whether this version can be used to create HCP clusters.
func (o *Version) HostedControlPlaneEnabled() bool {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.hostedControlPlaneEnabled
	}
	return false
}

// GetHostedControlPlaneEnabled returns the value of the 'hosted_control_plane_enabled' attribute and
// a flag indicating if the attribute has a value.
//
// HostedControlPlaneEnabled indicates whether this version can be used to create HCP clusters.
func (o *Version) GetHostedControlPlaneEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.hostedControlPlaneEnabled
	}
	return
}

// ImageOverrides returns the value of the 'image_overrides' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ImageOverrides contains the lists of images per cloud provider.
func (o *Version) ImageOverrides() *ImageOverrides {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.imageOverrides
	}
	return nil
}

// GetImageOverrides returns the value of the 'image_overrides' attribute and
// a flag indicating if the attribute has a value.
//
// ImageOverrides contains the lists of images per cloud provider.
func (o *Version) GetImageOverrides() (value *ImageOverrides, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.imageOverrides
	}
	return
}

// RawID returns the value of the 'raw_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// RawID is the id of the version - without channel group and prefix.
func (o *Version) RawID() string {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.rawID
	}
	return ""
}

// GetRawID returns the value of the 'raw_ID' attribute and
// a flag indicating if the attribute has a value.
//
// RawID is the id of the version - without channel group and prefix.
func (o *Version) GetRawID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.rawID
	}
	return
}

// ReleaseImage returns the value of the 'release_image' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ReleaseImage contains the URI of Openshift release image for amd64 architecture.
func (o *Version) ReleaseImage() string {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.releaseImage
	}
	return ""
}

// GetReleaseImage returns the value of the 'release_image' attribute and
// a flag indicating if the attribute has a value.
//
// ReleaseImage contains the URI of Openshift release image for amd64 architecture.
func (o *Version) GetReleaseImage() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.releaseImage
	}
	return
}

// ReleaseImages returns the value of the 'release_images' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ReleaseImages contains the URI of OpenShift release images for arm64 and multi architectures.
func (o *Version) ReleaseImages() *ReleaseImages {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.releaseImages
	}
	return nil
}

// GetReleaseImages returns the value of the 'release_images' attribute and
// a flag indicating if the attribute has a value.
//
// ReleaseImages contains the URI of OpenShift release images for arm64 and multi architectures.
func (o *Version) GetReleaseImages() (value *ReleaseImages, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.releaseImages
	}
	return
}

// WifEnabled returns the value of the 'wif_enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// WifEnabled is a flag that indicates whether this version is enabled for Workload Identity Federation.
func (o *Version) WifEnabled() bool {
	if o != nil && o.bitmap_&65536 != 0 {
		return o.wifEnabled
	}
	return false
}

// GetWifEnabled returns the value of the 'wif_enabled' attribute and
// a flag indicating if the attribute has a value.
//
// WifEnabled is a flag that indicates whether this version is enabled for Workload Identity Federation.
func (o *Version) GetWifEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&65536 != 0
	if ok {
		value = o.wifEnabled
	}
	return
}

// VersionListKind is the name of the type used to represent list of objects of
// type 'version'.
const VersionListKind = "VersionList"

// VersionListLinkKind is the name of the type used to represent links to list
// of objects of type 'version'.
const VersionListLinkKind = "VersionListLink"

// VersionNilKind is the name of the type used to nil lists of objects of
// type 'version'.
const VersionListNilKind = "VersionListNil"

// VersionList is a list of values of the 'version' type.
type VersionList struct {
	href  string
	link  bool
	items []*Version
}

// Kind returns the name of the type of the object.
func (l *VersionList) Kind() string {
	if l == nil {
		return VersionListNilKind
	}
	if l.link {
		return VersionListLinkKind
	}
	return VersionListKind
}

// Link returns true iif this is a link.
func (l *VersionList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *VersionList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *VersionList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *VersionList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *VersionList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *VersionList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *VersionList) SetItems(items []*Version) {
	l.items = items
}

// Items returns the items of the list.
func (l *VersionList) Items() []*Version {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *VersionList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *VersionList) Get(i int) *Version {
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
func (l *VersionList) Slice() []*Version {
	var slice []*Version
	if l == nil {
		slice = make([]*Version, 0)
	} else {
		slice = make([]*Version, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *VersionList) Each(f func(item *Version) bool) {
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
func (l *VersionList) Range(f func(index int, item *Version) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
