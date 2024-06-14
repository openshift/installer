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

// KubeletConfigBuilder contains the data and logic needed to build 'kubelet_config' objects.
//
// OCM representation of KubeletConfig, exposing the fields of Kubernetes
// KubeletConfig that can be managed by users
type KubeletConfigBuilder struct {
	bitmap_      uint32
	id           string
	href         string
	name         string
	podPidsLimit int
}

// NewKubeletConfig creates a new builder of 'kubelet_config' objects.
func NewKubeletConfig() *KubeletConfigBuilder {
	return &KubeletConfigBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *KubeletConfigBuilder) Link(value bool) *KubeletConfigBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *KubeletConfigBuilder) ID(value string) *KubeletConfigBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *KubeletConfigBuilder) HREF(value string) *KubeletConfigBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *KubeletConfigBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Name sets the value of the 'name' attribute to the given value.
func (b *KubeletConfigBuilder) Name(value string) *KubeletConfigBuilder {
	b.name = value
	b.bitmap_ |= 8
	return b
}

// PodPidsLimit sets the value of the 'pod_pids_limit' attribute to the given value.
func (b *KubeletConfigBuilder) PodPidsLimit(value int) *KubeletConfigBuilder {
	b.podPidsLimit = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *KubeletConfigBuilder) Copy(object *KubeletConfig) *KubeletConfigBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	object.name = b.name
	object.podPidsLimit = b.podPidsLimit
	return
}
