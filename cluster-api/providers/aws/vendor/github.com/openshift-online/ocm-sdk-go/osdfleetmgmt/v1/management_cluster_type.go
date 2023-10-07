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

package v1 // github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1

import (
	time "time"
)

// ManagementClusterKind is the name of the type used to represent objects
// of type 'management_cluster'.
const ManagementClusterKind = "ManagementCluster"

// ManagementClusterLinkKind is the name of the type used to represent links
// to objects of type 'management_cluster'.
const ManagementClusterLinkKind = "ManagementClusterLink"

// ManagementClusterNilKind is the name of the type used to nil references
// to objects of type 'management_cluster'.
const ManagementClusterNilKind = "ManagementClusterNil"

// ManagementCluster represents the values of the 'management_cluster' type.
//
// Definition of an _OpenShift_ cluster.
//
// The `cloud_provider` attribute is a reference to the cloud provider. When a
// cluster is retrieved it will be a link to the cloud provider, containing only
// the kind, id and href attributes:
//
// ```json
//
//	{
//	  "cloud_provider": {
//	    "kind": "CloudProviderLink",
//	    "id": "123",
//	    "href": "/api/clusters_mgmt/v1/cloud_providers/123"
//	  }
//	}
//
// ```
//
// When a cluster is created this is optional, and if used it should contain the
// identifier of the cloud provider to use:
//
// ```json
//
//	{
//	  "cloud_provider": {
//	    "id": "123",
//	  }
//	}
//
// ```
//
// If not included, then the cluster will be created using the default cloud
// provider, which is currently Amazon Web Services.
//
// The region attribute is mandatory when a cluster is created.
//
// The `aws.access_key_id`, `aws.secret_access_key` and `dns.base_domain`
// attributes are mandatory when creation a cluster with your own Amazon Web
// Services account.
type ManagementCluster struct {
	bitmap_                    uint32
	id                         string
	href                       string
	dns                        *DNS
	cloudProvider              string
	clusterManagementReference *ClusterManagementReference
	creationTimestamp          time.Time
	labels                     []*Label
	name                       string
	parent                     *ManagementClusterParent
	region                     string
	sector                     string
	status                     string
	updateTimestamp            time.Time
}

// Kind returns the name of the type of the object.
func (o *ManagementCluster) Kind() string {
	if o == nil {
		return ManagementClusterNilKind
	}
	if o.bitmap_&1 != 0 {
		return ManagementClusterLinkKind
	}
	return ManagementClusterKind
}

// Link returns true iif this is a link.
func (o *ManagementCluster) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *ManagementCluster) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ManagementCluster) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ManagementCluster) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ManagementCluster) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ManagementCluster) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// DNS returns the value of the 'DNS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// DNS settings of the cluster.
func (o *ManagementCluster) DNS() *DNS {
	if o != nil && o.bitmap_&8 != 0 {
		return o.dns
	}
	return nil
}

// GetDNS returns the value of the 'DNS' attribute and
// a flag indicating if the attribute has a value.
//
// DNS settings of the cluster.
func (o *ManagementCluster) GetDNS() (value *DNS, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.dns
	}
	return
}

// CloudProvider returns the value of the 'cloud_provider' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cloud provider where the cluster is installed.
func (o *ManagementCluster) CloudProvider() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.cloudProvider
	}
	return ""
}

// GetCloudProvider returns the value of the 'cloud_provider' attribute and
// a flag indicating if the attribute has a value.
//
// Cloud provider where the cluster is installed.
func (o *ManagementCluster) GetCloudProvider() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.cloudProvider
	}
	return
}

// ClusterManagementReference returns the value of the 'cluster_management_reference' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster mgmt reference
func (o *ManagementCluster) ClusterManagementReference() *ClusterManagementReference {
	if o != nil && o.bitmap_&32 != 0 {
		return o.clusterManagementReference
	}
	return nil
}

// GetClusterManagementReference returns the value of the 'cluster_management_reference' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster mgmt reference
func (o *ManagementCluster) GetClusterManagementReference() (value *ClusterManagementReference, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.clusterManagementReference
	}
	return
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Creation timestamp of the cluster
func (o *ManagementCluster) CreationTimestamp() time.Time {
	if o != nil && o.bitmap_&64 != 0 {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Creation timestamp of the cluster
func (o *ManagementCluster) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.creationTimestamp
	}
	return
}

// Labels returns the value of the 'labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Labels on management cluster
func (o *ManagementCluster) Labels() []*Label {
	if o != nil && o.bitmap_&128 != 0 {
		return o.labels
	}
	return nil
}

// GetLabels returns the value of the 'labels' attribute and
// a flag indicating if the attribute has a value.
//
// Labels on management cluster
func (o *ManagementCluster) GetLabels() (value []*Label, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.labels
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster name
func (o *ManagementCluster) Name() string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster name
func (o *ManagementCluster) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.name
	}
	return
}

// Parent returns the value of the 'parent' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Service cluster handling the management cluster
func (o *ManagementCluster) Parent() *ManagementClusterParent {
	if o != nil && o.bitmap_&512 != 0 {
		return o.parent
	}
	return nil
}

// GetParent returns the value of the 'parent' attribute and
// a flag indicating if the attribute has a value.
//
// Service cluster handling the management cluster
func (o *ManagementCluster) GetParent() (value *ManagementClusterParent, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.parent
	}
	return
}

// Region returns the value of the 'region' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cloud provider region where the cluster is installed.
func (o *ManagementCluster) Region() string {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.region
	}
	return ""
}

// GetRegion returns the value of the 'region' attribute and
// a flag indicating if the attribute has a value.
//
// Cloud provider region where the cluster is installed.
func (o *ManagementCluster) GetRegion() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.region
	}
	return
}

// Sector returns the value of the 'sector' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Sector of cluster
func (o *ManagementCluster) Sector() string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.sector
	}
	return ""
}

// GetSector returns the value of the 'sector' attribute and
// a flag indicating if the attribute has a value.
//
// Sector of cluster
func (o *ManagementCluster) GetSector() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.sector
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Status of cluster
func (o *ManagementCluster) Status() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.status
	}
	return ""
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// Status of cluster
func (o *ManagementCluster) GetStatus() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.status
	}
	return
}

// UpdateTimestamp returns the value of the 'update_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Update timestamp of the cluster
func (o *ManagementCluster) UpdateTimestamp() time.Time {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.updateTimestamp
	}
	return time.Time{}
}

// GetUpdateTimestamp returns the value of the 'update_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Update timestamp of the cluster
func (o *ManagementCluster) GetUpdateTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.updateTimestamp
	}
	return
}

// ManagementClusterListKind is the name of the type used to represent list of objects of
// type 'management_cluster'.
const ManagementClusterListKind = "ManagementClusterList"

// ManagementClusterListLinkKind is the name of the type used to represent links to list
// of objects of type 'management_cluster'.
const ManagementClusterListLinkKind = "ManagementClusterListLink"

// ManagementClusterNilKind is the name of the type used to nil lists of objects of
// type 'management_cluster'.
const ManagementClusterListNilKind = "ManagementClusterListNil"

// ManagementClusterList is a list of values of the 'management_cluster' type.
type ManagementClusterList struct {
	href  string
	link  bool
	items []*ManagementCluster
}

// Kind returns the name of the type of the object.
func (l *ManagementClusterList) Kind() string {
	if l == nil {
		return ManagementClusterListNilKind
	}
	if l.link {
		return ManagementClusterListLinkKind
	}
	return ManagementClusterListKind
}

// Link returns true iif this is a link.
func (l *ManagementClusterList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ManagementClusterList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ManagementClusterList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ManagementClusterList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ManagementClusterList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ManagementClusterList) Get(i int) *ManagementCluster {
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
func (l *ManagementClusterList) Slice() []*ManagementCluster {
	var slice []*ManagementCluster
	if l == nil {
		slice = make([]*ManagementCluster, 0)
	} else {
		slice = make([]*ManagementCluster, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ManagementClusterList) Each(f func(item *ManagementCluster) bool) {
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
func (l *ManagementClusterList) Range(f func(index int, item *ManagementCluster) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
