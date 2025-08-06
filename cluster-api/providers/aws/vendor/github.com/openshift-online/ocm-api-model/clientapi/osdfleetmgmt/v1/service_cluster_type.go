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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/osdfleetmgmt/v1

// ServiceClusterKind is the name of the type used to represent objects
// of type 'service_cluster'.
const ServiceClusterKind = "ServiceCluster"

// ServiceClusterLinkKind is the name of the type used to represent links
// to objects of type 'service_cluster'.
const ServiceClusterLinkKind = "ServiceClusterLink"

// ServiceClusterNilKind is the name of the type used to nil references
// to objects of type 'service_cluster'.
const ServiceClusterNilKind = "ServiceClusterNil"

// ServiceCluster represents the values of the 'service_cluster' type.
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
type ServiceCluster struct {
	fieldSet_                  []bool
	id                         string
	href                       string
	dns                        *DNS
	cloudProvider              string
	clusterManagementReference *ClusterManagementReference
	labels                     []*Label
	name                       string
	provisionShardReference    *ProvisionShardReference
	region                     string
	sector                     string
	status                     string
}

// Kind returns the name of the type of the object.
func (o *ServiceCluster) Kind() string {
	if o == nil {
		return ServiceClusterNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ServiceClusterLinkKind
	}
	return ServiceClusterKind
}

// Link returns true if this is a link.
func (o *ServiceCluster) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *ServiceCluster) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ServiceCluster) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ServiceCluster) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ServiceCluster) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ServiceCluster) Empty() bool {
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

// DNS returns the value of the 'DNS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// DNS settings of the cluster.
func (o *ServiceCluster) DNS() *DNS {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.dns
	}
	return nil
}

// GetDNS returns the value of the 'DNS' attribute and
// a flag indicating if the attribute has a value.
//
// DNS settings of the cluster.
func (o *ServiceCluster) GetDNS() (value *DNS, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.dns
	}
	return
}

// CloudProvider returns the value of the 'cloud_provider' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cloud provider where the cluster is installed.
func (o *ServiceCluster) CloudProvider() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.cloudProvider
	}
	return ""
}

// GetCloudProvider returns the value of the 'cloud_provider' attribute and
// a flag indicating if the attribute has a value.
//
// Cloud provider where the cluster is installed.
func (o *ServiceCluster) GetCloudProvider() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.cloudProvider
	}
	return
}

// ClusterManagementReference returns the value of the 'cluster_management_reference' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster mgmt reference
func (o *ServiceCluster) ClusterManagementReference() *ClusterManagementReference {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.clusterManagementReference
	}
	return nil
}

// GetClusterManagementReference returns the value of the 'cluster_management_reference' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster mgmt reference
func (o *ServiceCluster) GetClusterManagementReference() (value *ClusterManagementReference, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.clusterManagementReference
	}
	return
}

// Labels returns the value of the 'labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Labels on service cluster
func (o *ServiceCluster) Labels() []*Label {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.labels
	}
	return nil
}

// GetLabels returns the value of the 'labels' attribute and
// a flag indicating if the attribute has a value.
//
// Labels on service cluster
func (o *ServiceCluster) GetLabels() (value []*Label, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.labels
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster name
func (o *ServiceCluster) Name() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster name
func (o *ServiceCluster) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.name
	}
	return
}

// ProvisionShardReference returns the value of the 'provision_shard_reference' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Provision shard reference for the service cluster
func (o *ServiceCluster) ProvisionShardReference() *ProvisionShardReference {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.provisionShardReference
	}
	return nil
}

// GetProvisionShardReference returns the value of the 'provision_shard_reference' attribute and
// a flag indicating if the attribute has a value.
//
// Provision shard reference for the service cluster
func (o *ServiceCluster) GetProvisionShardReference() (value *ProvisionShardReference, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.provisionShardReference
	}
	return
}

// Region returns the value of the 'region' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cloud provider region where the cluster is installed.
func (o *ServiceCluster) Region() string {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.region
	}
	return ""
}

// GetRegion returns the value of the 'region' attribute and
// a flag indicating if the attribute has a value.
//
// Cloud provider region where the cluster is installed.
func (o *ServiceCluster) GetRegion() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.region
	}
	return
}

// Sector returns the value of the 'sector' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Sector of cluster
func (o *ServiceCluster) Sector() string {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.sector
	}
	return ""
}

// GetSector returns the value of the 'sector' attribute and
// a flag indicating if the attribute has a value.
//
// Sector of cluster
func (o *ServiceCluster) GetSector() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.sector
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Status of cluster
func (o *ServiceCluster) Status() string {
	if o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11] {
		return o.status
	}
	return ""
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// Status of cluster
func (o *ServiceCluster) GetStatus() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11]
	if ok {
		value = o.status
	}
	return
}

// ServiceClusterListKind is the name of the type used to represent list of objects of
// type 'service_cluster'.
const ServiceClusterListKind = "ServiceClusterList"

// ServiceClusterListLinkKind is the name of the type used to represent links to list
// of objects of type 'service_cluster'.
const ServiceClusterListLinkKind = "ServiceClusterListLink"

// ServiceClusterNilKind is the name of the type used to nil lists of objects of
// type 'service_cluster'.
const ServiceClusterListNilKind = "ServiceClusterListNil"

// ServiceClusterList is a list of values of the 'service_cluster' type.
type ServiceClusterList struct {
	href  string
	link  bool
	items []*ServiceCluster
}

// Kind returns the name of the type of the object.
func (l *ServiceClusterList) Kind() string {
	if l == nil {
		return ServiceClusterListNilKind
	}
	if l.link {
		return ServiceClusterListLinkKind
	}
	return ServiceClusterListKind
}

// Link returns true iif this is a link.
func (l *ServiceClusterList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ServiceClusterList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ServiceClusterList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ServiceClusterList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ServiceClusterList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ServiceClusterList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ServiceClusterList) SetItems(items []*ServiceCluster) {
	l.items = items
}

// Items returns the items of the list.
func (l *ServiceClusterList) Items() []*ServiceCluster {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ServiceClusterList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ServiceClusterList) Get(i int) *ServiceCluster {
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
func (l *ServiceClusterList) Slice() []*ServiceCluster {
	var slice []*ServiceCluster
	if l == nil {
		slice = make([]*ServiceCluster, 0)
	} else {
		slice = make([]*ServiceCluster, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ServiceClusterList) Each(f func(item *ServiceCluster) bool) {
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
func (l *ServiceClusterList) Range(f func(index int, item *ServiceCluster) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
