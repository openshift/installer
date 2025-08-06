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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/servicemgmt/v1

// This represents the parameters needed by Managed Service to create a cluster.
type ClusterBuilder struct {
	fieldSet_   []bool
	api         *ClusterAPIBuilder
	aws         *AWSBuilder
	displayName string
	href        string
	id          string
	name        string
	network     *NetworkBuilder
	nodes       *ClusterNodesBuilder
	properties  map[string]string
	region      *CloudRegionBuilder
	state       string
	multiAZ     bool
}

// NewCluster creates a new builder of 'cluster' objects.
func NewCluster() *ClusterBuilder {
	return &ClusterBuilder{
		fieldSet_: make([]bool, 12),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterBuilder) Empty() bool {
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

// API sets the value of the 'API' attribute to the given value.
//
// Information about the API of a cluster.
func (b *ClusterBuilder) API(value *ClusterAPIBuilder) *ClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.api = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// AWS sets the value of the 'AWS' attribute to the given value.
//
// _Amazon Web Services_ specific settings of a cluster.
func (b *ClusterBuilder) AWS(value *AWSBuilder) *ClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.aws = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *ClusterBuilder) DisplayName(value string) *ClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.displayName = value
	b.fieldSet_[2] = true
	return b
}

// Href sets the value of the 'href' attribute to the given value.
func (b *ClusterBuilder) Href(value string) *ClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.href = value
	b.fieldSet_[3] = true
	return b
}

// Id sets the value of the 'id' attribute to the given value.
func (b *ClusterBuilder) Id(value string) *ClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.id = value
	b.fieldSet_[4] = true
	return b
}

// MultiAZ sets the value of the 'multi_AZ' attribute to the given value.
func (b *ClusterBuilder) MultiAZ(value bool) *ClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.multiAZ = value
	b.fieldSet_[5] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ClusterBuilder) Name(value string) *ClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.name = value
	b.fieldSet_[6] = true
	return b
}

// Network sets the value of the 'network' attribute to the given value.
//
// Network configuration of a cluster.
func (b *ClusterBuilder) Network(value *NetworkBuilder) *ClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.network = value
	if value != nil {
		b.fieldSet_[7] = true
	} else {
		b.fieldSet_[7] = false
	}
	return b
}

// Nodes sets the value of the 'nodes' attribute to the given value.
func (b *ClusterBuilder) Nodes(value *ClusterNodesBuilder) *ClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.nodes = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// Properties sets the value of the 'properties' attribute to the given value.
func (b *ClusterBuilder) Properties(value map[string]string) *ClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.properties = value
	if value != nil {
		b.fieldSet_[9] = true
	} else {
		b.fieldSet_[9] = false
	}
	return b
}

// Region sets the value of the 'region' attribute to the given value.
//
// Description of a region of a cloud provider.
func (b *ClusterBuilder) Region(value *CloudRegionBuilder) *ClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.region = value
	if value != nil {
		b.fieldSet_[10] = true
	} else {
		b.fieldSet_[10] = false
	}
	return b
}

// State sets the value of the 'state' attribute to the given value.
func (b *ClusterBuilder) State(value string) *ClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.state = value
	b.fieldSet_[11] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterBuilder) Copy(object *Cluster) *ClusterBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.api != nil {
		b.api = NewClusterAPI().Copy(object.api)
	} else {
		b.api = nil
	}
	if object.aws != nil {
		b.aws = NewAWS().Copy(object.aws)
	} else {
		b.aws = nil
	}
	b.displayName = object.displayName
	b.href = object.href
	b.id = object.id
	b.multiAZ = object.multiAZ
	b.name = object.name
	if object.network != nil {
		b.network = NewNetwork().Copy(object.network)
	} else {
		b.network = nil
	}
	if object.nodes != nil {
		b.nodes = NewClusterNodes().Copy(object.nodes)
	} else {
		b.nodes = nil
	}
	if len(object.properties) > 0 {
		b.properties = map[string]string{}
		for k, v := range object.properties {
			b.properties[k] = v
		}
	} else {
		b.properties = nil
	}
	if object.region != nil {
		b.region = NewCloudRegion().Copy(object.region)
	} else {
		b.region = nil
	}
	b.state = object.state
	return b
}

// Build creates a 'cluster' object using the configuration stored in the builder.
func (b *ClusterBuilder) Build() (object *Cluster, err error) {
	object = new(Cluster)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.api != nil {
		object.api, err = b.api.Build()
		if err != nil {
			return
		}
	}
	if b.aws != nil {
		object.aws, err = b.aws.Build()
		if err != nil {
			return
		}
	}
	object.displayName = b.displayName
	object.href = b.href
	object.id = b.id
	object.multiAZ = b.multiAZ
	object.name = b.name
	if b.network != nil {
		object.network, err = b.network.Build()
		if err != nil {
			return
		}
	}
	if b.nodes != nil {
		object.nodes, err = b.nodes.Build()
		if err != nil {
			return
		}
	}
	if b.properties != nil {
		object.properties = make(map[string]string)
		for k, v := range b.properties {
			object.properties[k] = v
		}
	}
	if b.region != nil {
		object.region, err = b.region.Build()
		if err != nil {
			return
		}
	}
	object.state = b.state
	return
}
