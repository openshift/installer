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

// AddOnVersionKind is the name of the type used to represent objects
// of type 'add_on_version'.
const AddOnVersionKind = "AddOnVersion"

// AddOnVersionLinkKind is the name of the type used to represent links
// to objects of type 'add_on_version'.
const AddOnVersionLinkKind = "AddOnVersionLink"

// AddOnVersionNilKind is the name of the type used to nil references
// to objects of type 'add_on_version'.
const AddOnVersionNilKind = "AddOnVersionNil"

// AddOnVersion represents the values of the 'add_on_version' type.
//
// Representation of an add-on version.
type AddOnVersion struct {
	bitmap_                  uint32
	id                       string
	href                     string
	additionalCatalogSources []*AdditionalCatalogSource
	availableUpgrades        []string
	channel                  string
	config                   *AddOnConfig
	packageImage             string
	parameters               *AddOnParameterList
	pullSecretName           string
	requirements             []*AddOnRequirement
	sourceImage              string
	subOperators             []*AddOnSubOperator
	enabled                  bool
}

// Kind returns the name of the type of the object.
func (o *AddOnVersion) Kind() string {
	if o == nil {
		return AddOnVersionNilKind
	}
	if o.bitmap_&1 != 0 {
		return AddOnVersionLinkKind
	}
	return AddOnVersionKind
}

// Link returns true iif this is a link.
func (o *AddOnVersion) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *AddOnVersion) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AddOnVersion) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AddOnVersion) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AddOnVersion) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddOnVersion) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// AdditionalCatalogSources returns the value of the 'additional_catalog_sources' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Additional catalog sources associated with this addon version
func (o *AddOnVersion) AdditionalCatalogSources() []*AdditionalCatalogSource {
	if o != nil && o.bitmap_&8 != 0 {
		return o.additionalCatalogSources
	}
	return nil
}

// GetAdditionalCatalogSources returns the value of the 'additional_catalog_sources' attribute and
// a flag indicating if the attribute has a value.
//
// Additional catalog sources associated with this addon version
func (o *AddOnVersion) GetAdditionalCatalogSources() (value []*AdditionalCatalogSource, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.additionalCatalogSources
	}
	return
}

// AvailableUpgrades returns the value of the 'available_upgrades' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AvailableUpgrades is the list of versions this version can be upgraded to.
func (o *AddOnVersion) AvailableUpgrades() []string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.availableUpgrades
	}
	return nil
}

// GetAvailableUpgrades returns the value of the 'available_upgrades' attribute and
// a flag indicating if the attribute has a value.
//
// AvailableUpgrades is the list of versions this version can be upgraded to.
func (o *AddOnVersion) GetAvailableUpgrades() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.availableUpgrades
	}
	return
}

// Channel returns the value of the 'channel' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The specific addon catalog source channel of packages
func (o *AddOnVersion) Channel() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.channel
	}
	return ""
}

// GetChannel returns the value of the 'channel' attribute and
// a flag indicating if the attribute has a value.
//
// The specific addon catalog source channel of packages
func (o *AddOnVersion) GetChannel() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.channel
	}
	return
}

// Config returns the value of the 'config' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Additional configs to be used by the addon once its installed in the cluster.
func (o *AddOnVersion) Config() *AddOnConfig {
	if o != nil && o.bitmap_&64 != 0 {
		return o.config
	}
	return nil
}

// GetConfig returns the value of the 'config' attribute and
// a flag indicating if the attribute has a value.
//
// Additional configs to be used by the addon once its installed in the cluster.
func (o *AddOnVersion) GetConfig() (value *AddOnConfig, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.config
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this add-on version can be added to clusters.
func (o *AddOnVersion) Enabled() bool {
	if o != nil && o.bitmap_&128 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this add-on version can be added to clusters.
func (o *AddOnVersion) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.enabled
	}
	return
}

// PackageImage returns the value of the 'package_image' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The package image for this addon version
func (o *AddOnVersion) PackageImage() string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.packageImage
	}
	return ""
}

// GetPackageImage returns the value of the 'package_image' attribute and
// a flag indicating if the attribute has a value.
//
// The package image for this addon version
func (o *AddOnVersion) GetPackageImage() (value string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.packageImage
	}
	return
}

// Parameters returns the value of the 'parameters' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of parameters for this add-on version.
func (o *AddOnVersion) Parameters() *AddOnParameterList {
	if o != nil && o.bitmap_&512 != 0 {
		return o.parameters
	}
	return nil
}

// GetParameters returns the value of the 'parameters' attribute and
// a flag indicating if the attribute has a value.
//
// List of parameters for this add-on version.
func (o *AddOnVersion) GetParameters() (value *AddOnParameterList, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.parameters
	}
	return
}

// PullSecretName returns the value of the 'pull_secret_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The pull secret name used for this addon version.
func (o *AddOnVersion) PullSecretName() string {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.pullSecretName
	}
	return ""
}

// GetPullSecretName returns the value of the 'pull_secret_name' attribute and
// a flag indicating if the attribute has a value.
//
// The pull secret name used for this addon version.
func (o *AddOnVersion) GetPullSecretName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.pullSecretName
	}
	return
}

// Requirements returns the value of the 'requirements' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of requirements for this add-on version.
func (o *AddOnVersion) Requirements() []*AddOnRequirement {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.requirements
	}
	return nil
}

// GetRequirements returns the value of the 'requirements' attribute and
// a flag indicating if the attribute has a value.
//
// List of requirements for this add-on version.
func (o *AddOnVersion) GetRequirements() (value []*AddOnRequirement, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.requirements
	}
	return
}

// SourceImage returns the value of the 'source_image' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The catalog source image for this add-on version.
func (o *AddOnVersion) SourceImage() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.sourceImage
	}
	return ""
}

// GetSourceImage returns the value of the 'source_image' attribute and
// a flag indicating if the attribute has a value.
//
// The catalog source image for this add-on version.
func (o *AddOnVersion) GetSourceImage() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.sourceImage
	}
	return
}

// SubOperators returns the value of the 'sub_operators' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of sub operators for this add-on version.
func (o *AddOnVersion) SubOperators() []*AddOnSubOperator {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.subOperators
	}
	return nil
}

// GetSubOperators returns the value of the 'sub_operators' attribute and
// a flag indicating if the attribute has a value.
//
// List of sub operators for this add-on version.
func (o *AddOnVersion) GetSubOperators() (value []*AddOnSubOperator, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.subOperators
	}
	return
}

// AddOnVersionListKind is the name of the type used to represent list of objects of
// type 'add_on_version'.
const AddOnVersionListKind = "AddOnVersionList"

// AddOnVersionListLinkKind is the name of the type used to represent links to list
// of objects of type 'add_on_version'.
const AddOnVersionListLinkKind = "AddOnVersionListLink"

// AddOnVersionNilKind is the name of the type used to nil lists of objects of
// type 'add_on_version'.
const AddOnVersionListNilKind = "AddOnVersionListNil"

// AddOnVersionList is a list of values of the 'add_on_version' type.
type AddOnVersionList struct {
	href  string
	link  bool
	items []*AddOnVersion
}

// Kind returns the name of the type of the object.
func (l *AddOnVersionList) Kind() string {
	if l == nil {
		return AddOnVersionListNilKind
	}
	if l.link {
		return AddOnVersionListLinkKind
	}
	return AddOnVersionListKind
}

// Link returns true iif this is a link.
func (l *AddOnVersionList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AddOnVersionList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AddOnVersionList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AddOnVersionList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AddOnVersionList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddOnVersionList) Get(i int) *AddOnVersion {
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
func (l *AddOnVersionList) Slice() []*AddOnVersion {
	var slice []*AddOnVersion
	if l == nil {
		slice = make([]*AddOnVersion, 0)
	} else {
		slice = make([]*AddOnVersion, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddOnVersionList) Each(f func(item *AddOnVersion) bool) {
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
func (l *AddOnVersionList) Range(f func(index int, item *AddOnVersion) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
