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

// OCM representation of KubeletConfig, exposing the fields of Kubernetes
// KubeletConfig that can be managed by users
type KubeletConfigBuilder struct {
	fieldSet_    []bool
	id           string
	href         string
	name         string
	podPidsLimit int
}

// NewKubeletConfig creates a new builder of 'kubelet_config' objects.
func NewKubeletConfig() *KubeletConfigBuilder {
	return &KubeletConfigBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *KubeletConfigBuilder) Link(value bool) *KubeletConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *KubeletConfigBuilder) ID(value string) *KubeletConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *KubeletConfigBuilder) HREF(value string) *KubeletConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *KubeletConfigBuilder) Empty() bool {
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

// Name sets the value of the 'name' attribute to the given value.
func (b *KubeletConfigBuilder) Name(value string) *KubeletConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.name = value
	b.fieldSet_[3] = true
	return b
}

// PodPidsLimit sets the value of the 'pod_pids_limit' attribute to the given value.
func (b *KubeletConfigBuilder) PodPidsLimit(value int) *KubeletConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.podPidsLimit = value
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *KubeletConfigBuilder) Copy(object *KubeletConfig) *KubeletConfigBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.name = object.name
	b.podPidsLimit = object.podPidsLimit
	return b
}

// Build creates a 'kubelet_config' object using the configuration stored in the builder.
func (b *KubeletConfigBuilder) Build() (object *KubeletConfig, err error) {
	object = new(KubeletConfig)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.name = b.name
	object.podPidsLimit = b.podPidsLimit
	return
}
