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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

// AddonVersionKind is the name of the type used to represent objects
// of type 'addon_version'.
const AddonVersionKind = "AddonVersion"

// AddonVersionLinkKind is the name of the type used to represent links
// to objects of type 'addon_version'.
const AddonVersionLinkKind = "AddonVersionLink"

// AddonVersionNilKind is the name of the type used to nil references
// to objects of type 'addon_version'.
const AddonVersionNilKind = "AddonVersionNil"

// AddonVersion represents the values of the 'addon_version' type.
//
// Representation of an addon version.
type AddonVersion struct {
	bitmap_                  uint32
	id                       string
	href                     string
	additionalCatalogSources []*AdditionalCatalogSource
	availableUpgrades        []string
	channel                  string
	config                   *AddonConfig
	metricsFederation        *MetricsFederation
	monitoringStack          *MonitoringStack
	packageImage             string
	parameters               *AddonParameters
	pullSecretName           string
	requirements             []*AddonRequirement
	sourceImage              string
	subOperators             []*AddonSubOperator
	enabled                  bool
	upgradePlansCreated      bool
}

// Kind returns the name of the type of the object.
func (o *AddonVersion) Kind() string {
	if o == nil {
		return AddonVersionNilKind
	}
	if o.bitmap_&1 != 0 {
		return AddonVersionLinkKind
	}
	return AddonVersionKind
}

// Link returns true if this is a link.
func (o *AddonVersion) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *AddonVersion) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AddonVersion) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AddonVersion) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AddonVersion) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonVersion) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// AdditionalCatalogSources returns the value of the 'additional_catalog_sources' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Additional catalog sources associated with this addon version
func (o *AddonVersion) AdditionalCatalogSources() []*AdditionalCatalogSource {
	if o != nil && o.bitmap_&8 != 0 {
		return o.additionalCatalogSources
	}
	return nil
}

// GetAdditionalCatalogSources returns the value of the 'additional_catalog_sources' attribute and
// a flag indicating if the attribute has a value.
//
// Additional catalog sources associated with this addon version
func (o *AddonVersion) GetAdditionalCatalogSources() (value []*AdditionalCatalogSource, ok bool) {
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
func (o *AddonVersion) AvailableUpgrades() []string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.availableUpgrades
	}
	return nil
}

// GetAvailableUpgrades returns the value of the 'available_upgrades' attribute and
// a flag indicating if the attribute has a value.
//
// AvailableUpgrades is the list of versions this version can be upgraded to.
func (o *AddonVersion) GetAvailableUpgrades() (value []string, ok bool) {
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
func (o *AddonVersion) Channel() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.channel
	}
	return ""
}

// GetChannel returns the value of the 'channel' attribute and
// a flag indicating if the attribute has a value.
//
// The specific addon catalog source channel of packages
func (o *AddonVersion) GetChannel() (value string, ok bool) {
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
func (o *AddonVersion) Config() *AddonConfig {
	if o != nil && o.bitmap_&64 != 0 {
		return o.config
	}
	return nil
}

// GetConfig returns the value of the 'config' attribute and
// a flag indicating if the attribute has a value.
//
// Additional configs to be used by the addon once its installed in the cluster.
func (o *AddonVersion) GetConfig() (value *AddonConfig, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.config
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this addon version can be added to clusters.
func (o *AddonVersion) Enabled() bool {
	if o != nil && o.bitmap_&128 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this addon version can be added to clusters.
func (o *AddonVersion) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.enabled
	}
	return
}

// MetricsFederation returns the value of the 'metrics_federation' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Configuration parameters to be injected in the ServiceMonitor used for federation.
func (o *AddonVersion) MetricsFederation() *MetricsFederation {
	if o != nil && o.bitmap_&256 != 0 {
		return o.metricsFederation
	}
	return nil
}

// GetMetricsFederation returns the value of the 'metrics_federation' attribute and
// a flag indicating if the attribute has a value.
//
// Configuration parameters to be injected in the ServiceMonitor used for federation.
func (o *AddonVersion) GetMetricsFederation() (value *MetricsFederation, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.metricsFederation
	}
	return
}

// MonitoringStack returns the value of the 'monitoring_stack' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Configuration parameters which will determine the underlying configuration of the MonitoringStack CR.
func (o *AddonVersion) MonitoringStack() *MonitoringStack {
	if o != nil && o.bitmap_&512 != 0 {
		return o.monitoringStack
	}
	return nil
}

// GetMonitoringStack returns the value of the 'monitoring_stack' attribute and
// a flag indicating if the attribute has a value.
//
// Configuration parameters which will determine the underlying configuration of the MonitoringStack CR.
func (o *AddonVersion) GetMonitoringStack() (value *MonitoringStack, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.monitoringStack
	}
	return
}

// PackageImage returns the value of the 'package_image' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The url for the package image
func (o *AddonVersion) PackageImage() string {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.packageImage
	}
	return ""
}

// GetPackageImage returns the value of the 'package_image' attribute and
// a flag indicating if the attribute has a value.
//
// The url for the package image
func (o *AddonVersion) GetPackageImage() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.packageImage
	}
	return
}

// Parameters returns the value of the 'parameters' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of parameters for this addon version.
func (o *AddonVersion) Parameters() *AddonParameters {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.parameters
	}
	return nil
}

// GetParameters returns the value of the 'parameters' attribute and
// a flag indicating if the attribute has a value.
//
// List of parameters for this addon version.
func (o *AddonVersion) GetParameters() (value *AddonParameters, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.parameters
	}
	return
}

// PullSecretName returns the value of the 'pull_secret_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The pull secret name used for this addon version.
func (o *AddonVersion) PullSecretName() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.pullSecretName
	}
	return ""
}

// GetPullSecretName returns the value of the 'pull_secret_name' attribute and
// a flag indicating if the attribute has a value.
//
// The pull secret name used for this addon version.
func (o *AddonVersion) GetPullSecretName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.pullSecretName
	}
	return
}

// Requirements returns the value of the 'requirements' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of requirements for this addon version.
func (o *AddonVersion) Requirements() []*AddonRequirement {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.requirements
	}
	return nil
}

// GetRequirements returns the value of the 'requirements' attribute and
// a flag indicating if the attribute has a value.
//
// List of requirements for this addon version.
func (o *AddonVersion) GetRequirements() (value []*AddonRequirement, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.requirements
	}
	return
}

// SourceImage returns the value of the 'source_image' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The catalog source image for this addon version.
func (o *AddonVersion) SourceImage() string {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.sourceImage
	}
	return ""
}

// GetSourceImage returns the value of the 'source_image' attribute and
// a flag indicating if the attribute has a value.
//
// The catalog source image for this addon version.
func (o *AddonVersion) GetSourceImage() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.sourceImage
	}
	return
}

// SubOperators returns the value of the 'sub_operators' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of sub operators for this addon version.
func (o *AddonVersion) SubOperators() []*AddonSubOperator {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.subOperators
	}
	return nil
}

// GetSubOperators returns the value of the 'sub_operators' attribute and
// a flag indicating if the attribute has a value.
//
// List of sub operators for this addon version.
func (o *AddonVersion) GetSubOperators() (value []*AddonSubOperator, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.subOperators
	}
	return
}

// UpgradePlansCreated returns the value of the 'upgrade_plans_created' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if upgrade plans have been created for this addon version
func (o *AddonVersion) UpgradePlansCreated() bool {
	if o != nil && o.bitmap_&65536 != 0 {
		return o.upgradePlansCreated
	}
	return false
}

// GetUpgradePlansCreated returns the value of the 'upgrade_plans_created' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if upgrade plans have been created for this addon version
func (o *AddonVersion) GetUpgradePlansCreated() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&65536 != 0
	if ok {
		value = o.upgradePlansCreated
	}
	return
}

// AddonVersionListKind is the name of the type used to represent list of objects of
// type 'addon_version'.
const AddonVersionListKind = "AddonVersionList"

// AddonVersionListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_version'.
const AddonVersionListLinkKind = "AddonVersionListLink"

// AddonVersionNilKind is the name of the type used to nil lists of objects of
// type 'addon_version'.
const AddonVersionListNilKind = "AddonVersionListNil"

// AddonVersionList is a list of values of the 'addon_version' type.
type AddonVersionList struct {
	href  string
	link  bool
	items []*AddonVersion
}

// Kind returns the name of the type of the object.
func (l *AddonVersionList) Kind() string {
	if l == nil {
		return AddonVersionListNilKind
	}
	if l.link {
		return AddonVersionListLinkKind
	}
	return AddonVersionListKind
}

// Link returns true iif this is a link.
func (l *AddonVersionList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AddonVersionList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AddonVersionList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AddonVersionList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AddonVersionList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AddonVersionList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AddonVersionList) SetItems(items []*AddonVersion) {
	l.items = items
}

// Items returns the items of the list.
func (l *AddonVersionList) Items() []*AddonVersion {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AddonVersionList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonVersionList) Get(i int) *AddonVersion {
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
func (l *AddonVersionList) Slice() []*AddonVersion {
	var slice []*AddonVersion
	if l == nil {
		slice = make([]*AddonVersion, 0)
	} else {
		slice = make([]*AddonVersion, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonVersionList) Each(f func(item *AddonVersion) bool) {
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
func (l *AddonVersionList) Range(f func(index int, item *AddonVersion) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
