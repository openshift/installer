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

// Microsoft Azure settings of a cluster.
type AzureBuilder struct {
	fieldSet_                      []bool
	etcdEncryption                 *AzureEtcdEncryptionBuilder
	managedResourceGroupName       string
	networkSecurityGroupResourceID string
	nodesOutboundConnectivity      *AzureNodesOutboundConnectivityBuilder
	operatorsAuthentication        *AzureOperatorsAuthenticationBuilder
	resourceGroupName              string
	resourceName                   string
	subnetResourceID               string
	subscriptionID                 string
	tenantID                       string
}

// NewAzure creates a new builder of 'azure' objects.
func NewAzure() *AzureBuilder {
	return &AzureBuilder{
		fieldSet_: make([]bool, 10),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// EtcdEncryption sets the value of the 'etcd_encryption' attribute to the given value.
//
// Contains the necessary attributes to support etcd encryption for Azure based clusters.
func (b *AzureBuilder) EtcdEncryption(value *AzureEtcdEncryptionBuilder) *AzureBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.etcdEncryption = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// ManagedResourceGroupName sets the value of the 'managed_resource_group_name' attribute to the given value.
func (b *AzureBuilder) ManagedResourceGroupName(value string) *AzureBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.managedResourceGroupName = value
	b.fieldSet_[1] = true
	return b
}

// NetworkSecurityGroupResourceID sets the value of the 'network_security_group_resource_ID' attribute to the given value.
func (b *AzureBuilder) NetworkSecurityGroupResourceID(value string) *AzureBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.networkSecurityGroupResourceID = value
	b.fieldSet_[2] = true
	return b
}

// NodesOutboundConnectivity sets the value of the 'nodes_outbound_connectivity' attribute to the given value.
//
// The configuration of the node outbound connectivity
func (b *AzureBuilder) NodesOutboundConnectivity(value *AzureNodesOutboundConnectivityBuilder) *AzureBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.nodesOutboundConnectivity = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// OperatorsAuthentication sets the value of the 'operators_authentication' attribute to the given value.
//
// The configuration that the operators of the
// cluster have to authenticate to Azure.
func (b *AzureBuilder) OperatorsAuthentication(value *AzureOperatorsAuthenticationBuilder) *AzureBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.operatorsAuthentication = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// ResourceGroupName sets the value of the 'resource_group_name' attribute to the given value.
func (b *AzureBuilder) ResourceGroupName(value string) *AzureBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.resourceGroupName = value
	b.fieldSet_[5] = true
	return b
}

// ResourceName sets the value of the 'resource_name' attribute to the given value.
func (b *AzureBuilder) ResourceName(value string) *AzureBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.resourceName = value
	b.fieldSet_[6] = true
	return b
}

// SubnetResourceID sets the value of the 'subnet_resource_ID' attribute to the given value.
func (b *AzureBuilder) SubnetResourceID(value string) *AzureBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.subnetResourceID = value
	b.fieldSet_[7] = true
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *AzureBuilder) SubscriptionID(value string) *AzureBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.subscriptionID = value
	b.fieldSet_[8] = true
	return b
}

// TenantID sets the value of the 'tenant_ID' attribute to the given value.
func (b *AzureBuilder) TenantID(value string) *AzureBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.tenantID = value
	b.fieldSet_[9] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureBuilder) Copy(object *Azure) *AzureBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.etcdEncryption != nil {
		b.etcdEncryption = NewAzureEtcdEncryption().Copy(object.etcdEncryption)
	} else {
		b.etcdEncryption = nil
	}
	b.managedResourceGroupName = object.managedResourceGroupName
	b.networkSecurityGroupResourceID = object.networkSecurityGroupResourceID
	if object.nodesOutboundConnectivity != nil {
		b.nodesOutboundConnectivity = NewAzureNodesOutboundConnectivity().Copy(object.nodesOutboundConnectivity)
	} else {
		b.nodesOutboundConnectivity = nil
	}
	if object.operatorsAuthentication != nil {
		b.operatorsAuthentication = NewAzureOperatorsAuthentication().Copy(object.operatorsAuthentication)
	} else {
		b.operatorsAuthentication = nil
	}
	b.resourceGroupName = object.resourceGroupName
	b.resourceName = object.resourceName
	b.subnetResourceID = object.subnetResourceID
	b.subscriptionID = object.subscriptionID
	b.tenantID = object.tenantID
	return b
}

// Build creates a 'azure' object using the configuration stored in the builder.
func (b *AzureBuilder) Build() (object *Azure, err error) {
	object = new(Azure)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.etcdEncryption != nil {
		object.etcdEncryption, err = b.etcdEncryption.Build()
		if err != nil {
			return
		}
	}
	object.managedResourceGroupName = b.managedResourceGroupName
	object.networkSecurityGroupResourceID = b.networkSecurityGroupResourceID
	if b.nodesOutboundConnectivity != nil {
		object.nodesOutboundConnectivity, err = b.nodesOutboundConnectivity.Build()
		if err != nil {
			return
		}
	}
	if b.operatorsAuthentication != nil {
		object.operatorsAuthentication, err = b.operatorsAuthentication.Build()
		if err != nil {
			return
		}
	}
	object.resourceGroupName = b.resourceGroupName
	object.resourceName = b.resourceName
	object.subnetResourceID = b.subnetResourceID
	object.subscriptionID = b.subscriptionID
	object.tenantID = b.tenantID
	return
}
