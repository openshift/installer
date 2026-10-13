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

// Representation of a server config
type ServerConfigBuilder struct {
	fieldSet_  []bool
	id         string
	href       string
	awsShard   *AWSShardBuilder
	kubeconfig string
	server     string
	topology   ProvisionShardTopology
}

// NewServerConfig creates a new builder of 'server_config' objects.
func NewServerConfig() *ServerConfigBuilder {
	return &ServerConfigBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ServerConfigBuilder) Link(value bool) *ServerConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ServerConfigBuilder) ID(value string) *ServerConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ServerConfigBuilder) HREF(value string) *ServerConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ServerConfigBuilder) Empty() bool {
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

// AWSShard sets the value of the 'AWS_shard' attribute to the given value.
//
// Config for AWS provision shards
func (b *ServerConfigBuilder) AWSShard(value *AWSShardBuilder) *ServerConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.awsShard = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// Kubeconfig sets the value of the 'kubeconfig' attribute to the given value.
func (b *ServerConfigBuilder) Kubeconfig(value string) *ServerConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.kubeconfig = value
	b.fieldSet_[4] = true
	return b
}

// Server sets the value of the 'server' attribute to the given value.
func (b *ServerConfigBuilder) Server(value string) *ServerConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.server = value
	b.fieldSet_[5] = true
	return b
}

// Topology sets the value of the 'topology' attribute to the given value.
func (b *ServerConfigBuilder) Topology(value ProvisionShardTopology) *ServerConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.topology = value
	b.fieldSet_[6] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ServerConfigBuilder) Copy(object *ServerConfig) *ServerConfigBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.awsShard != nil {
		b.awsShard = NewAWSShard().Copy(object.awsShard)
	} else {
		b.awsShard = nil
	}
	b.kubeconfig = object.kubeconfig
	b.server = object.server
	b.topology = object.topology
	return b
}

// Build creates a 'server_config' object using the configuration stored in the builder.
func (b *ServerConfigBuilder) Build() (object *ServerConfig, err error) {
	object = new(ServerConfig)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.awsShard != nil {
		object.awsShard, err = b.awsShard.Build()
		if err != nil {
			return
		}
	}
	object.kubeconfig = b.kubeconfig
	object.server = b.server
	object.topology = b.topology
	return
}
