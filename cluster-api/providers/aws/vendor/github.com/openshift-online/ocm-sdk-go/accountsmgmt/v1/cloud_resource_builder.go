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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	time "time"
)

// CloudResourceBuilder contains the data and logic needed to build 'cloud_resource' objects.
type CloudResourceBuilder struct {
	bitmap_        uint32
	id             string
	href           string
	category       string
	categoryPretty string
	cloudProvider  string
	cpuCores       int
	createdAt      time.Time
	genericName    string
	memory         int
	memoryPretty   string
	namePretty     string
	resourceType   string
	sizePretty     string
	updatedAt      time.Time
	active         bool
}

// NewCloudResource creates a new builder of 'cloud_resource' objects.
func NewCloudResource() *CloudResourceBuilder {
	return &CloudResourceBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *CloudResourceBuilder) Link(value bool) *CloudResourceBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *CloudResourceBuilder) ID(value string) *CloudResourceBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *CloudResourceBuilder) HREF(value string) *CloudResourceBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CloudResourceBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Active sets the value of the 'active' attribute to the given value.
func (b *CloudResourceBuilder) Active(value bool) *CloudResourceBuilder {
	b.active = value
	b.bitmap_ |= 8
	return b
}

// Category sets the value of the 'category' attribute to the given value.
func (b *CloudResourceBuilder) Category(value string) *CloudResourceBuilder {
	b.category = value
	b.bitmap_ |= 16
	return b
}

// CategoryPretty sets the value of the 'category_pretty' attribute to the given value.
func (b *CloudResourceBuilder) CategoryPretty(value string) *CloudResourceBuilder {
	b.categoryPretty = value
	b.bitmap_ |= 32
	return b
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
func (b *CloudResourceBuilder) CloudProvider(value string) *CloudResourceBuilder {
	b.cloudProvider = value
	b.bitmap_ |= 64
	return b
}

// CpuCores sets the value of the 'cpu_cores' attribute to the given value.
func (b *CloudResourceBuilder) CpuCores(value int) *CloudResourceBuilder {
	b.cpuCores = value
	b.bitmap_ |= 128
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *CloudResourceBuilder) CreatedAt(value time.Time) *CloudResourceBuilder {
	b.createdAt = value
	b.bitmap_ |= 256
	return b
}

// GenericName sets the value of the 'generic_name' attribute to the given value.
func (b *CloudResourceBuilder) GenericName(value string) *CloudResourceBuilder {
	b.genericName = value
	b.bitmap_ |= 512
	return b
}

// Memory sets the value of the 'memory' attribute to the given value.
func (b *CloudResourceBuilder) Memory(value int) *CloudResourceBuilder {
	b.memory = value
	b.bitmap_ |= 1024
	return b
}

// MemoryPretty sets the value of the 'memory_pretty' attribute to the given value.
func (b *CloudResourceBuilder) MemoryPretty(value string) *CloudResourceBuilder {
	b.memoryPretty = value
	b.bitmap_ |= 2048
	return b
}

// NamePretty sets the value of the 'name_pretty' attribute to the given value.
func (b *CloudResourceBuilder) NamePretty(value string) *CloudResourceBuilder {
	b.namePretty = value
	b.bitmap_ |= 4096
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *CloudResourceBuilder) ResourceType(value string) *CloudResourceBuilder {
	b.resourceType = value
	b.bitmap_ |= 8192
	return b
}

// SizePretty sets the value of the 'size_pretty' attribute to the given value.
func (b *CloudResourceBuilder) SizePretty(value string) *CloudResourceBuilder {
	b.sizePretty = value
	b.bitmap_ |= 16384
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *CloudResourceBuilder) UpdatedAt(value time.Time) *CloudResourceBuilder {
	b.updatedAt = value
	b.bitmap_ |= 32768
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CloudResourceBuilder) Copy(object *CloudResource) *CloudResourceBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.active = object.active
	b.category = object.category
	b.categoryPretty = object.categoryPretty
	b.cloudProvider = object.cloudProvider
	b.cpuCores = object.cpuCores
	b.createdAt = object.createdAt
	b.genericName = object.genericName
	b.memory = object.memory
	b.memoryPretty = object.memoryPretty
	b.namePretty = object.namePretty
	b.resourceType = object.resourceType
	b.sizePretty = object.sizePretty
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'cloud_resource' object using the configuration stored in the builder.
func (b *CloudResourceBuilder) Build() (object *CloudResource, err error) {
	object = new(CloudResource)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.active = b.active
	object.category = b.category
	object.categoryPretty = b.categoryPretty
	object.cloudProvider = b.cloudProvider
	object.cpuCores = b.cpuCores
	object.createdAt = b.createdAt
	object.genericName = b.genericName
	object.memory = b.memory
	object.memoryPretty = b.memoryPretty
	object.namePretty = b.namePretty
	object.resourceType = b.resourceType
	object.sizePretty = b.sizePretty
	object.updatedAt = b.updatedAt
	return
}
