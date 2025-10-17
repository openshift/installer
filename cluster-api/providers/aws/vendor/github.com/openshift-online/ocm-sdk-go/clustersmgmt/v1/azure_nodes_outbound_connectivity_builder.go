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

// AzureNodesOutboundConnectivityBuilder contains the data and logic needed to build 'azure_nodes_outbound_connectivity' objects.
//
// The configuration of the node outbound connectivity
type AzureNodesOutboundConnectivityBuilder struct {
	bitmap_      uint32
	outboundType string
}

// NewAzureNodesOutboundConnectivity creates a new builder of 'azure_nodes_outbound_connectivity' objects.
func NewAzureNodesOutboundConnectivity() *AzureNodesOutboundConnectivityBuilder {
	return &AzureNodesOutboundConnectivityBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureNodesOutboundConnectivityBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// OutboundType sets the value of the 'outbound_type' attribute to the given value.
func (b *AzureNodesOutboundConnectivityBuilder) OutboundType(value string) *AzureNodesOutboundConnectivityBuilder {
	b.outboundType = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureNodesOutboundConnectivityBuilder) Copy(object *AzureNodesOutboundConnectivity) *AzureNodesOutboundConnectivityBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.outboundType = object.outboundType
	return b
}

// Build creates a 'azure_nodes_outbound_connectivity' object using the configuration stored in the builder.
func (b *AzureNodesOutboundConnectivityBuilder) Build() (object *AzureNodesOutboundConnectivity, err error) {
	object = new(AzureNodesOutboundConnectivity)
	object.bitmap_ = b.bitmap_
	object.outboundType = b.outboundType
	return
}
