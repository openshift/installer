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

// SdnToOvnClusterMigration represents the values of the 'sdn_to_ovn_cluster_migration' type.
//
// Details for `SdnToOvn` cluster migrations.
type SdnToOvnClusterMigration struct {
	fieldSet_      []bool
	joinIpv4       string
	masqueradeIpv4 string
	transitIpv4    string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *SdnToOvnClusterMigration) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}
	for _, set := range o.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// JoinIpv4 returns the value of the 'join_ipv_4' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The IP address range to use for the internalJoinSubnet parameter of OVN-Kubernetes
// upon migration.
func (o *SdnToOvnClusterMigration) JoinIpv4() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.joinIpv4
	}
	return ""
}

// GetJoinIpv4 returns the value of the 'join_ipv_4' attribute and
// a flag indicating if the attribute has a value.
//
// The IP address range to use for the internalJoinSubnet parameter of OVN-Kubernetes
// upon migration.
func (o *SdnToOvnClusterMigration) GetJoinIpv4() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.joinIpv4
	}
	return
}

// MasqueradeIpv4 returns the value of the 'masquerade_ipv_4' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The IP address range to us for the internalMasqueradeSubnet parameter of OVN-Kubernetes
// upon migration.
func (o *SdnToOvnClusterMigration) MasqueradeIpv4() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.masqueradeIpv4
	}
	return ""
}

// GetMasqueradeIpv4 returns the value of the 'masquerade_ipv_4' attribute and
// a flag indicating if the attribute has a value.
//
// The IP address range to us for the internalMasqueradeSubnet parameter of OVN-Kubernetes
// upon migration.
func (o *SdnToOvnClusterMigration) GetMasqueradeIpv4() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.masqueradeIpv4
	}
	return
}

// TransitIpv4 returns the value of the 'transit_ipv_4' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The IP address range to use for the internalTransSwitchSubnet parameter of OVN-Kubernetes
// upon migration.
func (o *SdnToOvnClusterMigration) TransitIpv4() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.transitIpv4
	}
	return ""
}

// GetTransitIpv4 returns the value of the 'transit_ipv_4' attribute and
// a flag indicating if the attribute has a value.
//
// The IP address range to use for the internalTransSwitchSubnet parameter of OVN-Kubernetes
// upon migration.
func (o *SdnToOvnClusterMigration) GetTransitIpv4() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.transitIpv4
	}
	return
}

// SdnToOvnClusterMigrationListKind is the name of the type used to represent list of objects of
// type 'sdn_to_ovn_cluster_migration'.
const SdnToOvnClusterMigrationListKind = "SdnToOvnClusterMigrationList"

// SdnToOvnClusterMigrationListLinkKind is the name of the type used to represent links to list
// of objects of type 'sdn_to_ovn_cluster_migration'.
const SdnToOvnClusterMigrationListLinkKind = "SdnToOvnClusterMigrationListLink"

// SdnToOvnClusterMigrationNilKind is the name of the type used to nil lists of objects of
// type 'sdn_to_ovn_cluster_migration'.
const SdnToOvnClusterMigrationListNilKind = "SdnToOvnClusterMigrationListNil"

// SdnToOvnClusterMigrationList is a list of values of the 'sdn_to_ovn_cluster_migration' type.
type SdnToOvnClusterMigrationList struct {
	href  string
	link  bool
	items []*SdnToOvnClusterMigration
}

// Len returns the length of the list.
func (l *SdnToOvnClusterMigrationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *SdnToOvnClusterMigrationList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *SdnToOvnClusterMigrationList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *SdnToOvnClusterMigrationList) SetItems(items []*SdnToOvnClusterMigration) {
	l.items = items
}

// Items returns the items of the list.
func (l *SdnToOvnClusterMigrationList) Items() []*SdnToOvnClusterMigration {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *SdnToOvnClusterMigrationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SdnToOvnClusterMigrationList) Get(i int) *SdnToOvnClusterMigration {
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
func (l *SdnToOvnClusterMigrationList) Slice() []*SdnToOvnClusterMigration {
	var slice []*SdnToOvnClusterMigration
	if l == nil {
		slice = make([]*SdnToOvnClusterMigration, 0)
	} else {
		slice = make([]*SdnToOvnClusterMigration, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SdnToOvnClusterMigrationList) Each(f func(item *SdnToOvnClusterMigration) bool) {
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
func (l *SdnToOvnClusterMigrationList) Range(f func(index int, item *SdnToOvnClusterMigration) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
