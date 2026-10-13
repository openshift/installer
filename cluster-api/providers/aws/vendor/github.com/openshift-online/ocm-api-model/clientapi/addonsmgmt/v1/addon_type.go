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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1

// AddonKind is the name of the type used to represent objects
// of type 'addon'.
const AddonKind = "Addon"

// AddonLinkKind is the name of the type used to represent links
// to objects of type 'addon'.
const AddonLinkKind = "AddonLink"

// AddonNilKind is the name of the type used to nil references
// to objects of type 'addon'.
const AddonNilKind = "AddonNil"

// Addon represents the values of the 'addon' type.
//
// Representation of an addon that can be installed in a cluster.
type Addon struct {
	fieldSet_            []bool
	id                   string
	href                 string
	commonAnnotations    map[string]string
	commonLabels         map[string]string
	config               *AddonConfig
	credentialsRequests  []*CredentialRequest
	description          string
	docsLink             string
	icon                 string
	installMode          AddonInstallMode
	label                string
	name                 string
	namespaces           []*AddonNamespace
	operatorName         string
	parameters           *AddonParameterList
	requirements         []*AddonRequirement
	resourceCost         float64
	resourceName         string
	subOperators         []*AddonSubOperator
	targetNamespace      string
	version              *AddonVersion
	enabled              bool
	hasExternalResources bool
	hidden               bool
	managedService       bool
}

// Kind returns the name of the type of the object.
func (o *Addon) Kind() string {
	if o == nil {
		return AddonNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return AddonLinkKind
	}
	return AddonKind
}

// Link returns true if this is a link.
func (o *Addon) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *Addon) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Addon) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Addon) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Addon) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Addon) Empty() bool {
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

// CommonAnnotations returns the value of the 'common_annotations' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Common Annotations for this addon.
func (o *Addon) CommonAnnotations() map[string]string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.commonAnnotations
	}
	return nil
}

// GetCommonAnnotations returns the value of the 'common_annotations' attribute and
// a flag indicating if the attribute has a value.
//
// Common Annotations for this addon.
func (o *Addon) GetCommonAnnotations() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.commonAnnotations
	}
	return
}

// CommonLabels returns the value of the 'common_labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Common Labels for this addon.
func (o *Addon) CommonLabels() map[string]string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.commonLabels
	}
	return nil
}

// GetCommonLabels returns the value of the 'common_labels' attribute and
// a flag indicating if the attribute has a value.
//
// Common Labels for this addon.
func (o *Addon) GetCommonLabels() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.commonLabels
	}
	return
}

// Config returns the value of the 'config' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Additional configs to be used by the addon once its installed in the cluster.
func (o *Addon) Config() *AddonConfig {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.config
	}
	return nil
}

// GetConfig returns the value of the 'config' attribute and
// a flag indicating if the attribute has a value.
//
// Additional configs to be used by the addon once its installed in the cluster.
func (o *Addon) GetConfig() (value *AddonConfig, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.config
	}
	return
}

// CredentialsRequests returns the value of the 'credentials_requests' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of credentials requests to authenticate operators to access cloud resources.
func (o *Addon) CredentialsRequests() []*CredentialRequest {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.credentialsRequests
	}
	return nil
}

// GetCredentialsRequests returns the value of the 'credentials_requests' attribute and
// a flag indicating if the attribute has a value.
//
// List of credentials requests to authenticate operators to access cloud resources.
func (o *Addon) GetCredentialsRequests() (value []*CredentialRequest, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.credentialsRequests
	}
	return
}

// Description returns the value of the 'description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Description of the addon.
func (o *Addon) Description() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.description
	}
	return ""
}

// GetDescription returns the value of the 'description' attribute and
// a flag indicating if the attribute has a value.
//
// Description of the addon.
func (o *Addon) GetDescription() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.description
	}
	return
}

// DocsLink returns the value of the 'docs_link' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to documentation about the addon.
func (o *Addon) DocsLink() string {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.docsLink
	}
	return ""
}

// GetDocsLink returns the value of the 'docs_link' attribute and
// a flag indicating if the attribute has a value.
//
// Link to documentation about the addon.
func (o *Addon) GetDocsLink() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.docsLink
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this addon can be added to clusters.
func (o *Addon) Enabled() bool {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this addon can be added to clusters.
func (o *Addon) GetEnabled() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.enabled
	}
	return
}

// HasExternalResources returns the value of the 'has_external_resources' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this addon has external resources associated with it
func (o *Addon) HasExternalResources() bool {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.hasExternalResources
	}
	return false
}

// GetHasExternalResources returns the value of the 'has_external_resources' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this addon has external resources associated with it
func (o *Addon) GetHasExternalResources() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.hasExternalResources
	}
	return
}

// Hidden returns the value of the 'hidden' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this addon is hidden.
func (o *Addon) Hidden() bool {
	if o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11] {
		return o.hidden
	}
	return false
}

// GetHidden returns the value of the 'hidden' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this addon is hidden.
func (o *Addon) GetHidden() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11]
	if ok {
		value = o.hidden
	}
	return
}

// Icon returns the value of the 'icon' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Base64-encoded icon representing an addon. The icon should be in PNG format.
func (o *Addon) Icon() string {
	if o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12] {
		return o.icon
	}
	return ""
}

// GetIcon returns the value of the 'icon' attribute and
// a flag indicating if the attribute has a value.
//
// Base64-encoded icon representing an addon. The icon should be in PNG format.
func (o *Addon) GetIcon() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12]
	if ok {
		value = o.icon
	}
	return
}

// InstallMode returns the value of the 'install_mode' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The mode in which the addon is deployed.
func (o *Addon) InstallMode() AddonInstallMode {
	if o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13] {
		return o.installMode
	}
	return AddonInstallMode("")
}

// GetInstallMode returns the value of the 'install_mode' attribute and
// a flag indicating if the attribute has a value.
//
// The mode in which the addon is deployed.
func (o *Addon) GetInstallMode() (value AddonInstallMode, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13]
	if ok {
		value = o.installMode
	}
	return
}

// Label returns the value of the 'label' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Label used to attach to a cluster deployment when addon is installed.
func (o *Addon) Label() string {
	if o != nil && len(o.fieldSet_) > 14 && o.fieldSet_[14] {
		return o.label
	}
	return ""
}

// GetLabel returns the value of the 'label' attribute and
// a flag indicating if the attribute has a value.
//
// Label used to attach to a cluster deployment when addon is installed.
func (o *Addon) GetLabel() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 14 && o.fieldSet_[14]
	if ok {
		value = o.label
	}
	return
}

// ManagedService returns the value of the 'managed_service' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if addon is part of a managed service
func (o *Addon) ManagedService() bool {
	if o != nil && len(o.fieldSet_) > 15 && o.fieldSet_[15] {
		return o.managedService
	}
	return false
}

// GetManagedService returns the value of the 'managed_service' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if addon is part of a managed service
func (o *Addon) GetManagedService() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 15 && o.fieldSet_[15]
	if ok {
		value = o.managedService
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the addon.
func (o *Addon) Name() string {
	if o != nil && len(o.fieldSet_) > 16 && o.fieldSet_[16] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the addon.
func (o *Addon) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 16 && o.fieldSet_[16]
	if ok {
		value = o.name
	}
	return
}

// Namespaces returns the value of the 'namespaces' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of namespaces associated with this addon.
func (o *Addon) Namespaces() []*AddonNamespace {
	if o != nil && len(o.fieldSet_) > 17 && o.fieldSet_[17] {
		return o.namespaces
	}
	return nil
}

// GetNamespaces returns the value of the 'namespaces' attribute and
// a flag indicating if the attribute has a value.
//
// List of namespaces associated with this addon.
func (o *Addon) GetNamespaces() (value []*AddonNamespace, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 17 && o.fieldSet_[17]
	if ok {
		value = o.namespaces
	}
	return
}

// OperatorName returns the value of the 'operator_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the operator installed by this addon.
func (o *Addon) OperatorName() string {
	if o != nil && len(o.fieldSet_) > 18 && o.fieldSet_[18] {
		return o.operatorName
	}
	return ""
}

// GetOperatorName returns the value of the 'operator_name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the operator installed by this addon.
func (o *Addon) GetOperatorName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 18 && o.fieldSet_[18]
	if ok {
		value = o.operatorName
	}
	return
}

// Parameters returns the value of the 'parameters' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of parameters for this addon.
func (o *Addon) Parameters() *AddonParameterList {
	if o != nil && len(o.fieldSet_) > 19 && o.fieldSet_[19] {
		return o.parameters
	}
	return nil
}

// GetParameters returns the value of the 'parameters' attribute and
// a flag indicating if the attribute has a value.
//
// List of parameters for this addon.
func (o *Addon) GetParameters() (value *AddonParameterList, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 19 && o.fieldSet_[19]
	if ok {
		value = o.parameters
	}
	return
}

// Requirements returns the value of the 'requirements' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of requirements for this addon.
func (o *Addon) Requirements() []*AddonRequirement {
	if o != nil && len(o.fieldSet_) > 20 && o.fieldSet_[20] {
		return o.requirements
	}
	return nil
}

// GetRequirements returns the value of the 'requirements' attribute and
// a flag indicating if the attribute has a value.
//
// List of requirements for this addon.
func (o *Addon) GetRequirements() (value []*AddonRequirement, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 20 && o.fieldSet_[20]
	if ok {
		value = o.requirements
	}
	return
}

// ResourceCost returns the value of the 'resource_cost' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Used to determine how many units of quota an addon consumes per resource name.
func (o *Addon) ResourceCost() float64 {
	if o != nil && len(o.fieldSet_) > 21 && o.fieldSet_[21] {
		return o.resourceCost
	}
	return 0.0
}

// GetResourceCost returns the value of the 'resource_cost' attribute and
// a flag indicating if the attribute has a value.
//
// Used to determine how many units of quota an addon consumes per resource name.
func (o *Addon) GetResourceCost() (value float64, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 21 && o.fieldSet_[21]
	if ok {
		value = o.resourceCost
	}
	return
}

// ResourceName returns the value of the 'resource_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Used to determine from where to reserve quota for this addon.
func (o *Addon) ResourceName() string {
	if o != nil && len(o.fieldSet_) > 22 && o.fieldSet_[22] {
		return o.resourceName
	}
	return ""
}

// GetResourceName returns the value of the 'resource_name' attribute and
// a flag indicating if the attribute has a value.
//
// Used to determine from where to reserve quota for this addon.
func (o *Addon) GetResourceName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 22 && o.fieldSet_[22]
	if ok {
		value = o.resourceName
	}
	return
}

// SubOperators returns the value of the 'sub_operators' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of sub operators for this addon.
func (o *Addon) SubOperators() []*AddonSubOperator {
	if o != nil && len(o.fieldSet_) > 23 && o.fieldSet_[23] {
		return o.subOperators
	}
	return nil
}

// GetSubOperators returns the value of the 'sub_operators' attribute and
// a flag indicating if the attribute has a value.
//
// List of sub operators for this addon.
func (o *Addon) GetSubOperators() (value []*AddonSubOperator, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 23 && o.fieldSet_[23]
	if ok {
		value = o.subOperators
	}
	return
}

// TargetNamespace returns the value of the 'target_namespace' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The namespace in which the addon CRD exists.
func (o *Addon) TargetNamespace() string {
	if o != nil && len(o.fieldSet_) > 24 && o.fieldSet_[24] {
		return o.targetNamespace
	}
	return ""
}

// GetTargetNamespace returns the value of the 'target_namespace' attribute and
// a flag indicating if the attribute has a value.
//
// The namespace in which the addon CRD exists.
func (o *Addon) GetTargetNamespace() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 24 && o.fieldSet_[24]
	if ok {
		value = o.targetNamespace
	}
	return
}

// Version returns the value of the 'version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the current default version of this addon.
func (o *Addon) Version() *AddonVersion {
	if o != nil && len(o.fieldSet_) > 25 && o.fieldSet_[25] {
		return o.version
	}
	return nil
}

// GetVersion returns the value of the 'version' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the current default version of this addon.
func (o *Addon) GetVersion() (value *AddonVersion, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 25 && o.fieldSet_[25]
	if ok {
		value = o.version
	}
	return
}

// AddonListKind is the name of the type used to represent list of objects of
// type 'addon'.
const AddonListKind = "AddonList"

// AddonListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon'.
const AddonListLinkKind = "AddonListLink"

// AddonNilKind is the name of the type used to nil lists of objects of
// type 'addon'.
const AddonListNilKind = "AddonListNil"

// AddonList is a list of values of the 'addon' type.
type AddonList struct {
	href  string
	link  bool
	items []*Addon
}

// Kind returns the name of the type of the object.
func (l *AddonList) Kind() string {
	if l == nil {
		return AddonListNilKind
	}
	if l.link {
		return AddonListLinkKind
	}
	return AddonListKind
}

// Link returns true iif this is a link.
func (l *AddonList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AddonList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AddonList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AddonList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AddonList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AddonList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AddonList) SetItems(items []*Addon) {
	l.items = items
}

// Items returns the items of the list.
func (l *AddonList) Items() []*Addon {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AddonList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonList) Get(i int) *Addon {
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
func (l *AddonList) Slice() []*Addon {
	var slice []*Addon
	if l == nil {
		slice = make([]*Addon, 0)
	} else {
		slice = make([]*Addon, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonList) Each(f func(item *Addon) bool) {
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
func (l *AddonList) Range(f func(index int, item *Addon) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
