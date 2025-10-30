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

import (
	time "time"
)

// ProvisionShardKind is the name of the type used to represent objects
// of type 'provision_shard'.
const ProvisionShardKind = "ProvisionShard"

// ProvisionShardLinkKind is the name of the type used to represent links
// to objects of type 'provision_shard'.
const ProvisionShardLinkKind = "ProvisionShardLink"

// ProvisionShardNilKind is the name of the type used to nil references
// to objects of type 'provision_shard'.
const ProvisionShardNilKind = "ProvisionShardNil"

// ProvisionShard represents the values of the 'provision_shard' type.
//
// Contains the properties of the provision shard, including AWS and GCP related configurations
type ProvisionShard struct {
	bitmap_                  uint32
	id                       string
	href                     string
	awsAccountOperatorConfig *ServerConfig
	awsBaseDomain            string
	gcpBaseDomain            string
	gcpProjectOperator       *ServerConfig
	cloudProvider            *CloudProvider
	creationTimestamp        time.Time
	hiveConfig               *ServerConfig
	hypershiftConfig         *ServerConfig
	lastUpdateTimestamp      time.Time
	managementCluster        string
	region                   *CloudRegion
	status                   string
}

// Kind returns the name of the type of the object.
func (o *ProvisionShard) Kind() string {
	if o == nil {
		return ProvisionShardNilKind
	}
	if o.bitmap_&1 != 0 {
		return ProvisionShardLinkKind
	}
	return ProvisionShardKind
}

// Link returns true if this is a link.
func (o *ProvisionShard) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *ProvisionShard) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ProvisionShard) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ProvisionShard) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ProvisionShard) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ProvisionShard) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// AWSAccountOperatorConfig returns the value of the 'AWS_account_operator_config' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains the configuration for the AWS account operator.
func (o *ProvisionShard) AWSAccountOperatorConfig() *ServerConfig {
	if o != nil && o.bitmap_&8 != 0 {
		return o.awsAccountOperatorConfig
	}
	return nil
}

// GetAWSAccountOperatorConfig returns the value of the 'AWS_account_operator_config' attribute and
// a flag indicating if the attribute has a value.
//
// Contains the configuration for the AWS account operator.
func (o *ProvisionShard) GetAWSAccountOperatorConfig() (value *ServerConfig, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.awsAccountOperatorConfig
	}
	return
}

// AWSBaseDomain returns the value of the 'AWS_base_domain' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains the AWS base domain.
func (o *ProvisionShard) AWSBaseDomain() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.awsBaseDomain
	}
	return ""
}

// GetAWSBaseDomain returns the value of the 'AWS_base_domain' attribute and
// a flag indicating if the attribute has a value.
//
// Contains the AWS base domain.
func (o *ProvisionShard) GetAWSBaseDomain() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.awsBaseDomain
	}
	return
}

// GCPBaseDomain returns the value of the 'GCP_base_domain' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains the GCP base domain.
func (o *ProvisionShard) GCPBaseDomain() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.gcpBaseDomain
	}
	return ""
}

// GetGCPBaseDomain returns the value of the 'GCP_base_domain' attribute and
// a flag indicating if the attribute has a value.
//
// Contains the GCP base domain.
func (o *ProvisionShard) GetGCPBaseDomain() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.gcpBaseDomain
	}
	return
}

// GCPProjectOperator returns the value of the 'GCP_project_operator' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains the configuration for the GCP project operator.
func (o *ProvisionShard) GCPProjectOperator() *ServerConfig {
	if o != nil && o.bitmap_&64 != 0 {
		return o.gcpProjectOperator
	}
	return nil
}

// GetGCPProjectOperator returns the value of the 'GCP_project_operator' attribute and
// a flag indicating if the attribute has a value.
//
// Contains the configuration for the GCP project operator.
func (o *ProvisionShard) GetGCPProjectOperator() (value *ServerConfig, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.gcpProjectOperator
	}
	return
}

// CloudProvider returns the value of the 'cloud_provider' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains the cloud provider name.
func (o *ProvisionShard) CloudProvider() *CloudProvider {
	if o != nil && o.bitmap_&128 != 0 {
		return o.cloudProvider
	}
	return nil
}

// GetCloudProvider returns the value of the 'cloud_provider' attribute and
// a flag indicating if the attribute has a value.
//
// Contains the cloud provider name.
func (o *ProvisionShard) GetCloudProvider() (value *CloudProvider, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.cloudProvider
	}
	return
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the provision shard was initially created, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *ProvisionShard) CreationTimestamp() time.Time {
	if o != nil && o.bitmap_&256 != 0 {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the provision shard was initially created, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *ProvisionShard) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.creationTimestamp
	}
	return
}

// HiveConfig returns the value of the 'hive_config' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains the configuration for Hive.
func (o *ProvisionShard) HiveConfig() *ServerConfig {
	if o != nil && o.bitmap_&512 != 0 {
		return o.hiveConfig
	}
	return nil
}

// GetHiveConfig returns the value of the 'hive_config' attribute and
// a flag indicating if the attribute has a value.
//
// Contains the configuration for Hive.
func (o *ProvisionShard) GetHiveConfig() (value *ServerConfig, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.hiveConfig
	}
	return
}

// HypershiftConfig returns the value of the 'hypershift_config' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains the configuration for Hypershift.
func (o *ProvisionShard) HypershiftConfig() *ServerConfig {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.hypershiftConfig
	}
	return nil
}

// GetHypershiftConfig returns the value of the 'hypershift_config' attribute and
// a flag indicating if the attribute has a value.
//
// Contains the configuration for Hypershift.
func (o *ProvisionShard) GetHypershiftConfig() (value *ServerConfig, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.hypershiftConfig
	}
	return
}

// LastUpdateTimestamp returns the value of the 'last_update_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the provision shard was last updated, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *ProvisionShard) LastUpdateTimestamp() time.Time {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.lastUpdateTimestamp
	}
	return time.Time{}
}

// GetLastUpdateTimestamp returns the value of the 'last_update_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the provision shard was last updated, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *ProvisionShard) GetLastUpdateTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.lastUpdateTimestamp
	}
	return
}

// ManagementCluster returns the value of the 'management_cluster' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains the name of the management cluster for Hypershift clusters that are assigned to this shard.
// This field is populated by OCM, and must not be overwritten via API.
func (o *ProvisionShard) ManagementCluster() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.managementCluster
	}
	return ""
}

// GetManagementCluster returns the value of the 'management_cluster' attribute and
// a flag indicating if the attribute has a value.
//
// Contains the name of the management cluster for Hypershift clusters that are assigned to this shard.
// This field is populated by OCM, and must not be overwritten via API.
func (o *ProvisionShard) GetManagementCluster() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.managementCluster
	}
	return
}

// Region returns the value of the 'region' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains the cloud-provider region in which the provisioner spins up the cluster.
func (o *ProvisionShard) Region() *CloudRegion {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.region
	}
	return nil
}

// GetRegion returns the value of the 'region' attribute and
// a flag indicating if the attribute has a value.
//
// Contains the cloud-provider region in which the provisioner spins up the cluster.
func (o *ProvisionShard) GetRegion() (value *CloudRegion, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.region
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Status of the provision shard. Possible values: active/maintenance/offline.
func (o *ProvisionShard) Status() string {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.status
	}
	return ""
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// Status of the provision shard. Possible values: active/maintenance/offline.
func (o *ProvisionShard) GetStatus() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.status
	}
	return
}

// ProvisionShardListKind is the name of the type used to represent list of objects of
// type 'provision_shard'.
const ProvisionShardListKind = "ProvisionShardList"

// ProvisionShardListLinkKind is the name of the type used to represent links to list
// of objects of type 'provision_shard'.
const ProvisionShardListLinkKind = "ProvisionShardListLink"

// ProvisionShardNilKind is the name of the type used to nil lists of objects of
// type 'provision_shard'.
const ProvisionShardListNilKind = "ProvisionShardListNil"

// ProvisionShardList is a list of values of the 'provision_shard' type.
type ProvisionShardList struct {
	href  string
	link  bool
	items []*ProvisionShard
}

// Kind returns the name of the type of the object.
func (l *ProvisionShardList) Kind() string {
	if l == nil {
		return ProvisionShardListNilKind
	}
	if l.link {
		return ProvisionShardListLinkKind
	}
	return ProvisionShardListKind
}

// Link returns true iif this is a link.
func (l *ProvisionShardList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ProvisionShardList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ProvisionShardList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ProvisionShardList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ProvisionShardList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ProvisionShardList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ProvisionShardList) SetItems(items []*ProvisionShard) {
	l.items = items
}

// Items returns the items of the list.
func (l *ProvisionShardList) Items() []*ProvisionShard {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ProvisionShardList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ProvisionShardList) Get(i int) *ProvisionShard {
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
func (l *ProvisionShardList) Slice() []*ProvisionShard {
	var slice []*ProvisionShard
	if l == nil {
		slice = make([]*ProvisionShard, 0)
	} else {
		slice = make([]*ProvisionShard, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ProvisionShardList) Each(f func(item *ProvisionShard) bool) {
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
func (l *ProvisionShardList) Range(f func(index int, item *ProvisionShard) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
