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

// Set of predefined properties of a cluster. For example, a _huge_ flavour can be a cluster
// with 10 infra nodes and 1000 compute nodes.
type FlavourBuilder struct {
	fieldSet_ []bool
	id        string
	href      string
	aws       *AWSFlavourBuilder
	gcp       *GCPFlavourBuilder
	name      string
	network   *NetworkBuilder
	nodes     *FlavourNodesBuilder
}

// NewFlavour creates a new builder of 'flavour' objects.
func NewFlavour() *FlavourBuilder {
	return &FlavourBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *FlavourBuilder) Link(value bool) *FlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *FlavourBuilder) ID(value string) *FlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *FlavourBuilder) HREF(value string) *FlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *FlavourBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// AWS sets the value of the 'AWS' attribute to the given value.
//
// Specification for different classes of nodes inside a flavour.
func (b *FlavourBuilder) AWS(value *AWSFlavourBuilder) *FlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.aws = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// GCP sets the value of the 'GCP' attribute to the given value.
//
// Specification for different classes of nodes inside a flavour.
func (b *FlavourBuilder) GCP(value *GCPFlavourBuilder) *FlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.gcp = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *FlavourBuilder) Name(value string) *FlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.name = value
	b.fieldSet_[5] = true
	return b
}

// Network sets the value of the 'network' attribute to the given value.
//
// Network configuration of a cluster.
func (b *FlavourBuilder) Network(value *NetworkBuilder) *FlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.network = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Nodes sets the value of the 'nodes' attribute to the given value.
//
// Counts of different classes of nodes inside a flavour.
func (b *FlavourBuilder) Nodes(value *FlavourNodesBuilder) *FlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.nodes = value
	if value != nil {
		b.fieldSet_[7] = true
	} else {
		b.fieldSet_[7] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *FlavourBuilder) Copy(object *Flavour) *FlavourBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.aws != nil {
		b.aws = NewAWSFlavour().Copy(object.aws)
	} else {
		b.aws = nil
	}
	if object.gcp != nil {
		b.gcp = NewGCPFlavour().Copy(object.gcp)
	} else {
		b.gcp = nil
	}
	b.name = object.name
	if object.network != nil {
		b.network = NewNetwork().Copy(object.network)
	} else {
		b.network = nil
	}
	if object.nodes != nil {
		b.nodes = NewFlavourNodes().Copy(object.nodes)
	} else {
		b.nodes = nil
	}
	return b
}

// Build creates a 'flavour' object using the configuration stored in the builder.
func (b *FlavourBuilder) Build() (object *Flavour, err error) {
	object = new(Flavour)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.aws != nil {
		object.aws, err = b.aws.Build()
		if err != nil {
			return
		}
	}
	if b.gcp != nil {
		object.gcp, err = b.gcp.Build()
		if err != nil {
			return
		}
	}
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
	return
}
