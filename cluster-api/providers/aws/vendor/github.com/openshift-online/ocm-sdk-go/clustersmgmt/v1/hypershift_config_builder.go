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

// HypershiftConfigBuilder contains the data and logic needed to build 'hypershift_config' objects.
//
// Hypershift configuration.
type HypershiftConfigBuilder struct {
	bitmap_           uint32
	hcpNamespace      string
	managementCluster string
	enabled           bool
}

// NewHypershiftConfig creates a new builder of 'hypershift_config' objects.
func NewHypershiftConfig() *HypershiftConfigBuilder {
	return &HypershiftConfigBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *HypershiftConfigBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// HCPNamespace sets the value of the 'HCP_namespace' attribute to the given value.
func (b *HypershiftConfigBuilder) HCPNamespace(value string) *HypershiftConfigBuilder {
	b.hcpNamespace = value
	b.bitmap_ |= 1
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *HypershiftConfigBuilder) Enabled(value bool) *HypershiftConfigBuilder {
	b.enabled = value
	b.bitmap_ |= 2
	return b
}

// ManagementCluster sets the value of the 'management_cluster' attribute to the given value.
func (b *HypershiftConfigBuilder) ManagementCluster(value string) *HypershiftConfigBuilder {
	b.managementCluster = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *HypershiftConfigBuilder) Copy(object *HypershiftConfig) *HypershiftConfigBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.hcpNamespace = object.hcpNamespace
	b.enabled = object.enabled
	b.managementCluster = object.managementCluster
	return b
}

// Build creates a 'hypershift_config' object using the configuration stored in the builder.
func (b *HypershiftConfigBuilder) Build() (object *HypershiftConfig, err error) {
	object = new(HypershiftConfig)
	object.bitmap_ = b.bitmap_
	object.hcpNamespace = b.hcpNamespace
	object.enabled = b.enabled
	object.managementCluster = b.managementCluster
	return
}
