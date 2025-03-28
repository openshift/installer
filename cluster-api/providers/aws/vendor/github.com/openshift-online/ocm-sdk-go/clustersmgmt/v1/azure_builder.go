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

// AzureBuilder contains the data and logic needed to build 'azure' objects.
//
// Microsoft Azure settings of a cluster.
type AzureBuilder struct {
	bitmap_                        uint32
	managedResourceGroupName       string
	networkSecurityGroupResourceID string
	operatorsAuthentication        *AzureOperatorsAuthenticationBuilder
	resourceGroupName              string
	resourceName                   string
	subnetResourceID               string
	subscriptionID                 string
	tenantID                       string
}

// NewAzure creates a new builder of 'azure' objects.
func NewAzure() *AzureBuilder {
	return &AzureBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ManagedResourceGroupName sets the value of the 'managed_resource_group_name' attribute to the given value.
func (b *AzureBuilder) ManagedResourceGroupName(value string) *AzureBuilder {
	b.managedResourceGroupName = value
	b.bitmap_ |= 1
	return b
}

// NetworkSecurityGroupResourceID sets the value of the 'network_security_group_resource_ID' attribute to the given value.
func (b *AzureBuilder) NetworkSecurityGroupResourceID(value string) *AzureBuilder {
	b.networkSecurityGroupResourceID = value
	b.bitmap_ |= 2
	return b
}

// OperatorsAuthentication sets the value of the 'operators_authentication' attribute to the given value.
//
// The configuration that the operators of the
// cluster have to authenticate to Azure.
func (b *AzureBuilder) OperatorsAuthentication(value *AzureOperatorsAuthenticationBuilder) *AzureBuilder {
	b.operatorsAuthentication = value
	if value != nil {
		b.bitmap_ |= 4
	} else {
		b.bitmap_ &^= 4
	}
	return b
}

// ResourceGroupName sets the value of the 'resource_group_name' attribute to the given value.
func (b *AzureBuilder) ResourceGroupName(value string) *AzureBuilder {
	b.resourceGroupName = value
	b.bitmap_ |= 8
	return b
}

// ResourceName sets the value of the 'resource_name' attribute to the given value.
func (b *AzureBuilder) ResourceName(value string) *AzureBuilder {
	b.resourceName = value
	b.bitmap_ |= 16
	return b
}

// SubnetResourceID sets the value of the 'subnet_resource_ID' attribute to the given value.
func (b *AzureBuilder) SubnetResourceID(value string) *AzureBuilder {
	b.subnetResourceID = value
	b.bitmap_ |= 32
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *AzureBuilder) SubscriptionID(value string) *AzureBuilder {
	b.subscriptionID = value
	b.bitmap_ |= 64
	return b
}

// TenantID sets the value of the 'tenant_ID' attribute to the given value.
func (b *AzureBuilder) TenantID(value string) *AzureBuilder {
	b.tenantID = value
	b.bitmap_ |= 128
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureBuilder) Copy(object *Azure) *AzureBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.managedResourceGroupName = object.managedResourceGroupName
	b.networkSecurityGroupResourceID = object.networkSecurityGroupResourceID
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
	object.bitmap_ = b.bitmap_
	object.managedResourceGroupName = b.managedResourceGroupName
	object.networkSecurityGroupResourceID = b.networkSecurityGroupResourceID
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
