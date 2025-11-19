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

// AddOnKind is the name of the type used to represent objects
// of type 'add_on'.
const AddOnKind = "AddOn"

// AddOnLinkKind is the name of the type used to represent links
// to objects of type 'add_on'.
const AddOnLinkKind = "AddOnLink"

// AddOnNilKind is the name of the type used to nil references
// to objects of type 'add_on'.
const AddOnNilKind = "AddOnNil"

// AddOn represents the values of the 'add_on' type.
//
// Representation of an add-on that can be installed in a cluster.
type AddOn struct {
	fieldSet_            []bool
	id                   string
	href                 string
	commonAnnotations    map[string]string
	commonLabels         map[string]string
	config               *AddOnConfig
	credentialsRequests  []*CredentialRequest
	description          string
	docsLink             string
	icon                 string
	installMode          AddOnInstallMode
	label                string
	name                 string
	namespaces           []*AddOnNamespace
	operatorName         string
	parameters           *AddOnParameterList
	requirements         []*AddOnRequirement
	resourceCost         float64
	resourceName         string
	subOperators         []*AddOnSubOperator
	targetNamespace      string
	version              *AddOnVersion
	enabled              bool
	hasExternalResources bool
	hidden               bool
	managedService       bool
}

// Kind returns the name of the type of the object.
func (o *AddOn) Kind() string {
	if o == nil {
		return AddOnNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return AddOnLinkKind
	}
	return AddOnKind
}

// Link returns true if this is a link.
func (o *AddOn) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *AddOn) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AddOn) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AddOn) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AddOn) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddOn) Empty() bool {
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
// Common annotations to be applied to all resources created by this addon.
func (o *AddOn) CommonAnnotations() map[string]string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.commonAnnotations
	}
	return nil
}

// GetCommonAnnotations returns the value of the 'common_annotations' attribute and
// a flag indicating if the attribute has a value.
//
// Common annotations to be applied to all resources created by this addon.
func (o *AddOn) GetCommonAnnotations() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.commonAnnotations
	}
	return
}

// CommonLabels returns the value of the 'common_labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Common labels to be applied to all resources created by this addon.
func (o *AddOn) CommonLabels() map[string]string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.commonLabels
	}
	return nil
}

// GetCommonLabels returns the value of the 'common_labels' attribute and
// a flag indicating if the attribute has a value.
//
// Common labels to be applied to all resources created by this addon.
func (o *AddOn) GetCommonLabels() (value map[string]string, ok bool) {
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
func (o *AddOn) Config() *AddOnConfig {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.config
	}
	return nil
}

// GetConfig returns the value of the 'config' attribute and
// a flag indicating if the attribute has a value.
//
// Additional configs to be used by the addon once its installed in the cluster.
func (o *AddOn) GetConfig() (value *AddOnConfig, ok bool) {
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
func (o *AddOn) CredentialsRequests() []*CredentialRequest {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.credentialsRequests
	}
	return nil
}

// GetCredentialsRequests returns the value of the 'credentials_requests' attribute and
// a flag indicating if the attribute has a value.
//
// List of credentials requests to authenticate operators to access cloud resources.
func (o *AddOn) GetCredentialsRequests() (value []*CredentialRequest, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.credentialsRequests
	}
	return
}

// Description returns the value of the 'description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Description of the add-on.
func (o *AddOn) Description() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.description
	}
	return ""
}

// GetDescription returns the value of the 'description' attribute and
// a flag indicating if the attribute has a value.
//
// Description of the add-on.
func (o *AddOn) GetDescription() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.description
	}
	return
}

// DocsLink returns the value of the 'docs_link' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to documentation about the add-on.
func (o *AddOn) DocsLink() string {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.docsLink
	}
	return ""
}

// GetDocsLink returns the value of the 'docs_link' attribute and
// a flag indicating if the attribute has a value.
//
// Link to documentation about the add-on.
func (o *AddOn) GetDocsLink() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.docsLink
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this add-on can be added to clusters.
func (o *AddOn) Enabled() bool {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this add-on can be added to clusters.
func (o *AddOn) GetEnabled() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.enabled
	}
	return
}

// HasExternalResources returns the value of the 'has_external_resources' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this add-on has external resources associated with it
func (o *AddOn) HasExternalResources() bool {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.hasExternalResources
	}
	return false
}

// GetHasExternalResources returns the value of the 'has_external_resources' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this add-on has external resources associated with it
func (o *AddOn) GetHasExternalResources() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.hasExternalResources
	}
	return
}

// Hidden returns the value of the 'hidden' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this add-on is hidden.
func (o *AddOn) Hidden() bool {
	if o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11] {
		return o.hidden
	}
	return false
}

// GetHidden returns the value of the 'hidden' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this add-on is hidden.
func (o *AddOn) GetHidden() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11]
	if ok {
		value = o.hidden
	}
	return
}

// Icon returns the value of the 'icon' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Base64-encoded icon representing an add-on. The icon should be in PNG format.
func (o *AddOn) Icon() string {
	if o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12] {
		return o.icon
	}
	return ""
}

// GetIcon returns the value of the 'icon' attribute and
// a flag indicating if the attribute has a value.
//
// Base64-encoded icon representing an add-on. The icon should be in PNG format.
func (o *AddOn) GetIcon() (value string, ok bool) {
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
func (o *AddOn) InstallMode() AddOnInstallMode {
	if o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13] {
		return o.installMode
	}
	return AddOnInstallMode("")
}

// GetInstallMode returns the value of the 'install_mode' attribute and
// a flag indicating if the attribute has a value.
//
// The mode in which the addon is deployed.
func (o *AddOn) GetInstallMode() (value AddOnInstallMode, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13]
	if ok {
		value = o.installMode
	}
	return
}

// Label returns the value of the 'label' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Label used to attach to a cluster deployment when add-on is installed.
func (o *AddOn) Label() string {
	if o != nil && len(o.fieldSet_) > 14 && o.fieldSet_[14] {
		return o.label
	}
	return ""
}

// GetLabel returns the value of the 'label' attribute and
// a flag indicating if the attribute has a value.
//
// Label used to attach to a cluster deployment when add-on is installed.
func (o *AddOn) GetLabel() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 14 && o.fieldSet_[14]
	if ok {
		value = o.label
	}
	return
}

// ManagedService returns the value of the 'managed_service' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if add-on is part of a managed service
func (o *AddOn) ManagedService() bool {
	if o != nil && len(o.fieldSet_) > 15 && o.fieldSet_[15] {
		return o.managedService
	}
	return false
}

// GetManagedService returns the value of the 'managed_service' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if add-on is part of a managed service
func (o *AddOn) GetManagedService() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 15 && o.fieldSet_[15]
	if ok {
		value = o.managedService
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the add-on.
func (o *AddOn) Name() string {
	if o != nil && len(o.fieldSet_) > 16 && o.fieldSet_[16] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the add-on.
func (o *AddOn) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 16 && o.fieldSet_[16]
	if ok {
		value = o.name
	}
	return
}

// Namespaces returns the value of the 'namespaces' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Namespaces which are required by this addon.
func (o *AddOn) Namespaces() []*AddOnNamespace {
	if o != nil && len(o.fieldSet_) > 17 && o.fieldSet_[17] {
		return o.namespaces
	}
	return nil
}

// GetNamespaces returns the value of the 'namespaces' attribute and
// a flag indicating if the attribute has a value.
//
// Namespaces which are required by this addon.
func (o *AddOn) GetNamespaces() (value []*AddOnNamespace, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 17 && o.fieldSet_[17]
	if ok {
		value = o.namespaces
	}
	return
}

// OperatorName returns the value of the 'operator_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the operator installed by this add-on.
func (o *AddOn) OperatorName() string {
	if o != nil && len(o.fieldSet_) > 18 && o.fieldSet_[18] {
		return o.operatorName
	}
	return ""
}

// GetOperatorName returns the value of the 'operator_name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the operator installed by this add-on.
func (o *AddOn) GetOperatorName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 18 && o.fieldSet_[18]
	if ok {
		value = o.operatorName
	}
	return
}

// Parameters returns the value of the 'parameters' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of parameters for this add-on.
func (o *AddOn) Parameters() *AddOnParameterList {
	if o != nil && len(o.fieldSet_) > 19 && o.fieldSet_[19] {
		return o.parameters
	}
	return nil
}

// GetParameters returns the value of the 'parameters' attribute and
// a flag indicating if the attribute has a value.
//
// List of parameters for this add-on.
func (o *AddOn) GetParameters() (value *AddOnParameterList, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 19 && o.fieldSet_[19]
	if ok {
		value = o.parameters
	}
	return
}

// Requirements returns the value of the 'requirements' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of requirements for this add-on.
func (o *AddOn) Requirements() []*AddOnRequirement {
	if o != nil && len(o.fieldSet_) > 20 && o.fieldSet_[20] {
		return o.requirements
	}
	return nil
}

// GetRequirements returns the value of the 'requirements' attribute and
// a flag indicating if the attribute has a value.
//
// List of requirements for this add-on.
func (o *AddOn) GetRequirements() (value []*AddOnRequirement, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 20 && o.fieldSet_[20]
	if ok {
		value = o.requirements
	}
	return
}

// ResourceCost returns the value of the 'resource_cost' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Used to determine how many units of quota an add-on consumes per resource name.
func (o *AddOn) ResourceCost() float64 {
	if o != nil && len(o.fieldSet_) > 21 && o.fieldSet_[21] {
		return o.resourceCost
	}
	return 0.0
}

// GetResourceCost returns the value of the 'resource_cost' attribute and
// a flag indicating if the attribute has a value.
//
// Used to determine how many units of quota an add-on consumes per resource name.
func (o *AddOn) GetResourceCost() (value float64, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 21 && o.fieldSet_[21]
	if ok {
		value = o.resourceCost
	}
	return
}

// ResourceName returns the value of the 'resource_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Used to determine from where to reserve quota for this add-on.
func (o *AddOn) ResourceName() string {
	if o != nil && len(o.fieldSet_) > 22 && o.fieldSet_[22] {
		return o.resourceName
	}
	return ""
}

// GetResourceName returns the value of the 'resource_name' attribute and
// a flag indicating if the attribute has a value.
//
// Used to determine from where to reserve quota for this add-on.
func (o *AddOn) GetResourceName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 22 && o.fieldSet_[22]
	if ok {
		value = o.resourceName
	}
	return
}

// SubOperators returns the value of the 'sub_operators' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of sub operators for this add-on.
func (o *AddOn) SubOperators() []*AddOnSubOperator {
	if o != nil && len(o.fieldSet_) > 23 && o.fieldSet_[23] {
		return o.subOperators
	}
	return nil
}

// GetSubOperators returns the value of the 'sub_operators' attribute and
// a flag indicating if the attribute has a value.
//
// List of sub operators for this add-on.
func (o *AddOn) GetSubOperators() (value []*AddOnSubOperator, ok bool) {
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
func (o *AddOn) TargetNamespace() string {
	if o != nil && len(o.fieldSet_) > 24 && o.fieldSet_[24] {
		return o.targetNamespace
	}
	return ""
}

// GetTargetNamespace returns the value of the 'target_namespace' attribute and
// a flag indicating if the attribute has a value.
//
// The namespace in which the addon CRD exists.
func (o *AddOn) GetTargetNamespace() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 24 && o.fieldSet_[24]
	if ok {
		value = o.targetNamespace
	}
	return
}

// Version returns the value of the 'version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the current default version of this add-on.
func (o *AddOn) Version() *AddOnVersion {
	if o != nil && len(o.fieldSet_) > 25 && o.fieldSet_[25] {
		return o.version
	}
	return nil
}

// GetVersion returns the value of the 'version' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the current default version of this add-on.
func (o *AddOn) GetVersion() (value *AddOnVersion, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 25 && o.fieldSet_[25]
	if ok {
		value = o.version
	}
	return
}

// AddOnListKind is the name of the type used to represent list of objects of
// type 'add_on'.
const AddOnListKind = "AddOnList"

// AddOnListLinkKind is the name of the type used to represent links to list
// of objects of type 'add_on'.
const AddOnListLinkKind = "AddOnListLink"

// AddOnNilKind is the name of the type used to nil lists of objects of
// type 'add_on'.
const AddOnListNilKind = "AddOnListNil"

// AddOnList is a list of values of the 'add_on' type.
type AddOnList struct {
	href  string
	link  bool
	items []*AddOn
}

// Kind returns the name of the type of the object.
func (l *AddOnList) Kind() string {
	if l == nil {
		return AddOnListNilKind
	}
	if l.link {
		return AddOnListLinkKind
	}
	return AddOnListKind
}

// Link returns true iif this is a link.
func (l *AddOnList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AddOnList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AddOnList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AddOnList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AddOnList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AddOnList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AddOnList) SetItems(items []*AddOn) {
	l.items = items
}

// Items returns the items of the list.
func (l *AddOnList) Items() []*AddOn {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AddOnList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddOnList) Get(i int) *AddOn {
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
func (l *AddOnList) Slice() []*AddOn {
	var slice []*AddOn
	if l == nil {
		slice = make([]*AddOn, 0)
	} else {
		slice = make([]*AddOn, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddOnList) Each(f func(item *AddOn) bool) {
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
func (l *AddOnList) Range(f func(index int, item *AddOn) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
