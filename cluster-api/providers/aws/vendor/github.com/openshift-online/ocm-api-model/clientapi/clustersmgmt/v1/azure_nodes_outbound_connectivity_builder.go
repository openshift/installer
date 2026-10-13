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

// The configuration of the node outbound connectivity
type AzureNodesOutboundConnectivityBuilder struct {
	fieldSet_    []bool
	outboundType string
}

// NewAzureNodesOutboundConnectivity creates a new builder of 'azure_nodes_outbound_connectivity' objects.
func NewAzureNodesOutboundConnectivity() *AzureNodesOutboundConnectivityBuilder {
	return &AzureNodesOutboundConnectivityBuilder{
		fieldSet_: make([]bool, 1),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureNodesOutboundConnectivityBuilder) Empty() bool {
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

// OutboundType sets the value of the 'outbound_type' attribute to the given value.
func (b *AzureNodesOutboundConnectivityBuilder) OutboundType(value string) *AzureNodesOutboundConnectivityBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 1)
	}
	b.outboundType = value
	b.fieldSet_[0] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureNodesOutboundConnectivityBuilder) Copy(object *AzureNodesOutboundConnectivity) *AzureNodesOutboundConnectivityBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.outboundType = object.outboundType
	return b
}

// Build creates a 'azure_nodes_outbound_connectivity' object using the configuration stored in the builder.
func (b *AzureNodesOutboundConnectivityBuilder) Build() (object *AzureNodesOutboundConnectivity, err error) {
	object = new(AzureNodesOutboundConnectivity)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.outboundType = b.outboundType
	return
}
