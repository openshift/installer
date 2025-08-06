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

// Azure represents the values of the 'azure' type.
//
// Microsoft Azure settings of a cluster.
type Azure struct {
	fieldSet_                      []bool
	etcdEncryption                 *AzureEtcdEncryption
	managedResourceGroupName       string
	networkSecurityGroupResourceID string
	nodesOutboundConnectivity      *AzureNodesOutboundConnectivity
	operatorsAuthentication        *AzureOperatorsAuthentication
	resourceGroupName              string
	resourceName                   string
	subnetResourceID               string
	subscriptionID                 string
	tenantID                       string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Azure) Empty() bool {
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

// EtcdEncryption returns the value of the 'etcd_encryption' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Etcd encryption configuration.
// If not specified, etcd data is encrypted with platform managed keys.
// Currently etcd data encryption is only supported with customer managed keys.
// Creating a cluster with platform managed keys will result in a failure creating the cluster.
func (o *Azure) EtcdEncryption() *AzureEtcdEncryption {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.etcdEncryption
	}
	return nil
}

// GetEtcdEncryption returns the value of the 'etcd_encryption' attribute and
// a flag indicating if the attribute has a value.
//
// Etcd encryption configuration.
// If not specified, etcd data is encrypted with platform managed keys.
// Currently etcd data encryption is only supported with customer managed keys.
// Creating a cluster with platform managed keys will result in a failure creating the cluster.
func (o *Azure) GetEtcdEncryption() (value *AzureEtcdEncryption, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.etcdEncryption
	}
	return
}

// ManagedResourceGroupName returns the value of the 'managed_resource_group_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The desired name of the Azure Resource Group where the Azure Resources related
// to the cluster are created. It must not previously exist. The Azure Resource
// Group is created with the given value, within the Azure Subscription
// `subscription_id` of the cluster.
// `managed_resource_group_name` cannot be equal to the value of `managed_resource_group`.
// `managed_resource_group_name` is located in the same Azure location as the
// cluster's region.
// Not to be confused with `resource_group_name`, which is the Azure Resource Group Name
// where the own Azure Resource associated to the cluster resides.
// Required during creation.
// Immutable.
func (o *Azure) ManagedResourceGroupName() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.managedResourceGroupName
	}
	return ""
}

// GetManagedResourceGroupName returns the value of the 'managed_resource_group_name' attribute and
// a flag indicating if the attribute has a value.
//
// The desired name of the Azure Resource Group where the Azure Resources related
// to the cluster are created. It must not previously exist. The Azure Resource
// Group is created with the given value, within the Azure Subscription
// `subscription_id` of the cluster.
// `managed_resource_group_name` cannot be equal to the value of `managed_resource_group`.
// `managed_resource_group_name` is located in the same Azure location as the
// cluster's region.
// Not to be confused with `resource_group_name`, which is the Azure Resource Group Name
// where the own Azure Resource associated to the cluster resides.
// Required during creation.
// Immutable.
func (o *Azure) GetManagedResourceGroupName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.managedResourceGroupName
	}
	return
}

// NetworkSecurityGroupResourceID returns the value of the 'network_security_group_resource_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Resource ID of a pre-existing Azure Network Security Group.
// The Network Security Group specified in network_security_group_resource_id
// must already be associated to the Azure Subnet `subnet_resource_id`.
// It is the Azure Network Security Group associated to the cluster's subnet
// specified in `subnet_resource_id`.
// `network_security_group_resource_id` must be located in the same Azure
// location as the cluster's region.
// The Azure Subscription specified as part of
// `network_security_group_resource_id` must be located in the same Azure
// Subscription as `subscription_id`.
// The Azure Resource Group Name specified as part of `network_security_group_resource_id`
// must belong to the Azure Subscription `subscription_id`, and in the same
// Azure location as the cluster's region.
// The Azure Resource Group Name specified as part of `network_security_group_resource_id`
// must be a different Resource Group Name than the one specified in
// `managed_resource_group_name`.
// The Azure Resource Group Name specified as part of `network_security_group_resource_id`
// can be the same, or a different one than the one specified in
// `resource_group_name`.
// Required during creation.
// Immutable.
func (o *Azure) NetworkSecurityGroupResourceID() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.networkSecurityGroupResourceID
	}
	return ""
}

// GetNetworkSecurityGroupResourceID returns the value of the 'network_security_group_resource_ID' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure Resource ID of a pre-existing Azure Network Security Group.
// The Network Security Group specified in network_security_group_resource_id
// must already be associated to the Azure Subnet `subnet_resource_id`.
// It is the Azure Network Security Group associated to the cluster's subnet
// specified in `subnet_resource_id`.
// `network_security_group_resource_id` must be located in the same Azure
// location as the cluster's region.
// The Azure Subscription specified as part of
// `network_security_group_resource_id` must be located in the same Azure
// Subscription as `subscription_id`.
// The Azure Resource Group Name specified as part of `network_security_group_resource_id`
// must belong to the Azure Subscription `subscription_id`, and in the same
// Azure location as the cluster's region.
// The Azure Resource Group Name specified as part of `network_security_group_resource_id`
// must be a different Resource Group Name than the one specified in
// `managed_resource_group_name`.
// The Azure Resource Group Name specified as part of `network_security_group_resource_id`
// can be the same, or a different one than the one specified in
// `resource_group_name`.
// Required during creation.
// Immutable.
func (o *Azure) GetNetworkSecurityGroupResourceID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.networkSecurityGroupResourceID
	}
	return
}

// NodesOutboundConnectivity returns the value of the 'nodes_outbound_connectivity' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// NodesOutboundConnectivity defines how the network outbound
// configuration of the Cluster's Node Pool's Nodes is performed.
// By default this is configured as Azure Load Balancer. This value is immutable.
func (o *Azure) NodesOutboundConnectivity() *AzureNodesOutboundConnectivity {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.nodesOutboundConnectivity
	}
	return nil
}

// GetNodesOutboundConnectivity returns the value of the 'nodes_outbound_connectivity' attribute and
// a flag indicating if the attribute has a value.
//
// NodesOutboundConnectivity defines how the network outbound
// configuration of the Cluster's Node Pool's Nodes is performed.
// By default this is configured as Azure Load Balancer. This value is immutable.
func (o *Azure) GetNodesOutboundConnectivity() (value *AzureNodesOutboundConnectivity, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.nodesOutboundConnectivity
	}
	return
}

// OperatorsAuthentication returns the value of the 'operators_authentication' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines how the operators of the cluster authenticate to Azure.
// Required during creation.
// Immutable.
func (o *Azure) OperatorsAuthentication() *AzureOperatorsAuthentication {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.operatorsAuthentication
	}
	return nil
}

// GetOperatorsAuthentication returns the value of the 'operators_authentication' attribute and
// a flag indicating if the attribute has a value.
//
// Defines how the operators of the cluster authenticate to Azure.
// Required during creation.
// Immutable.
func (o *Azure) GetOperatorsAuthentication() (value *AzureOperatorsAuthentication, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.operatorsAuthentication
	}
	return
}

// ResourceGroupName returns the value of the 'resource_group_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Resource Group Name of the cluster. It must be a pre-existing
// Azure Resource Group and it must exist within the Azure Subscription
// `subscription_id` of the cluster.
// `resource_group_name` is located in the same Azure location as the
// cluster's region.
// Required during creation.
// Immutable.
func (o *Azure) ResourceGroupName() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.resourceGroupName
	}
	return ""
}

// GetResourceGroupName returns the value of the 'resource_group_name' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure Resource Group Name of the cluster. It must be a pre-existing
// Azure Resource Group and it must exist within the Azure Subscription
// `subscription_id` of the cluster.
// `resource_group_name` is located in the same Azure location as the
// cluster's region.
// Required during creation.
// Immutable.
func (o *Azure) GetResourceGroupName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.resourceGroupName
	}
	return
}

// ResourceName returns the value of the 'resource_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Resource Name of the cluster. It must be within the
// Azure Resource Group Name `resource_group_name`.
// `resource_name` is located in the same Azure location as the cluster's region.
// Required during creation.
// Immutable.
func (o *Azure) ResourceName() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.resourceName
	}
	return ""
}

// GetResourceName returns the value of the 'resource_name' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure Resource Name of the cluster. It must be within the
// Azure Resource Group Name `resource_group_name`.
// `resource_name` is located in the same Azure location as the cluster's region.
// Required during creation.
// Immutable.
func (o *Azure) GetResourceName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.resourceName
	}
	return
}

// SubnetResourceID returns the value of the 'subnet_resource_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Resource ID of a pre-existing Azure Subnet. It is an Azure
// Subnet used for the Data Plane of the cluster. `subnet_resource_id`
// must be located in the same Azure location as the cluster's region.
// The Azure Subscription specified as part of the `subnet_resource_id`
// must be located in the same Azure Subscription as `subscription_id`.
// The Azure Resource Group Name specified as part of `subnet_resource_id`
// must belong to the Azure Subscription `subscription_id`, and in the same
// Azure location as the cluster's region.
// The Azure Resource Group Name specified as part of `subnet_resource_id`
// must be a different Resource Group Name than the one specified in
// `managed_resource_group_name`.
// The Azure Resource Group Name specified as part of the `subnet_resource_id`
// can be the same, or a different one than the one specified in
// `resource_group_name`.
// Required during creation.
// Immutable.
func (o *Azure) SubnetResourceID() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.subnetResourceID
	}
	return ""
}

// GetSubnetResourceID returns the value of the 'subnet_resource_ID' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure Resource ID of a pre-existing Azure Subnet. It is an Azure
// Subnet used for the Data Plane of the cluster. `subnet_resource_id`
// must be located in the same Azure location as the cluster's region.
// The Azure Subscription specified as part of the `subnet_resource_id`
// must be located in the same Azure Subscription as `subscription_id`.
// The Azure Resource Group Name specified as part of `subnet_resource_id`
// must belong to the Azure Subscription `subscription_id`, and in the same
// Azure location as the cluster's region.
// The Azure Resource Group Name specified as part of `subnet_resource_id`
// must be a different Resource Group Name than the one specified in
// `managed_resource_group_name`.
// The Azure Resource Group Name specified as part of the `subnet_resource_id`
// can be the same, or a different one than the one specified in
// `resource_group_name`.
// Required during creation.
// Immutable.
func (o *Azure) GetSubnetResourceID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.subnetResourceID
	}
	return
}

// SubscriptionID returns the value of the 'subscription_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Subscription ID associated with the cluster. It must belong to
// the Microsoft Entra Tenant ID `tenant_id`.
// Required during creation.
// Immutable.
func (o *Azure) SubscriptionID() string {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.subscriptionID
	}
	return ""
}

// GetSubscriptionID returns the value of the 'subscription_ID' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure Subscription ID associated with the cluster. It must belong to
// the Microsoft Entra Tenant ID `tenant_id`.
// Required during creation.
// Immutable.
func (o *Azure) GetSubscriptionID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.subscriptionID
	}
	return
}

// TenantID returns the value of the 'tenant_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Microsoft Entra Tenant ID where the cluster belongs.
// Required during creation.
// Immutable.
func (o *Azure) TenantID() string {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.tenantID
	}
	return ""
}

// GetTenantID returns the value of the 'tenant_ID' attribute and
// a flag indicating if the attribute has a value.
//
// The Microsoft Entra Tenant ID where the cluster belongs.
// Required during creation.
// Immutable.
func (o *Azure) GetTenantID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.tenantID
	}
	return
}

// AzureListKind is the name of the type used to represent list of objects of
// type 'azure'.
const AzureListKind = "AzureList"

// AzureListLinkKind is the name of the type used to represent links to list
// of objects of type 'azure'.
const AzureListLinkKind = "AzureListLink"

// AzureNilKind is the name of the type used to nil lists of objects of
// type 'azure'.
const AzureListNilKind = "AzureListNil"

// AzureList is a list of values of the 'azure' type.
type AzureList struct {
	href  string
	link  bool
	items []*Azure
}

// Len returns the length of the list.
func (l *AzureList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AzureList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AzureList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AzureList) SetItems(items []*Azure) {
	l.items = items
}

// Items returns the items of the list.
func (l *AzureList) Items() []*Azure {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AzureList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AzureList) Get(i int) *Azure {
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
func (l *AzureList) Slice() []*Azure {
	var slice []*Azure
	if l == nil {
		slice = make([]*Azure, 0)
	} else {
		slice = make([]*Azure, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AzureList) Each(f func(item *Azure) bool) {
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
func (l *AzureList) Range(f func(index int, item *Azure) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
