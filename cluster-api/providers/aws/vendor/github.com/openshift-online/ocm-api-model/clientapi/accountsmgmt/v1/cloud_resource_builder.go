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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

import (
	time "time"
)

type CloudResourceBuilder struct {
	fieldSet_      []bool
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
	return &CloudResourceBuilder{
		fieldSet_: make([]bool, 16),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *CloudResourceBuilder) Link(value bool) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *CloudResourceBuilder) ID(value string) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *CloudResourceBuilder) HREF(value string) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CloudResourceBuilder) Empty() bool {
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

// Active sets the value of the 'active' attribute to the given value.
func (b *CloudResourceBuilder) Active(value bool) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.active = value
	b.fieldSet_[3] = true
	return b
}

// Category sets the value of the 'category' attribute to the given value.
func (b *CloudResourceBuilder) Category(value string) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.category = value
	b.fieldSet_[4] = true
	return b
}

// CategoryPretty sets the value of the 'category_pretty' attribute to the given value.
func (b *CloudResourceBuilder) CategoryPretty(value string) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.categoryPretty = value
	b.fieldSet_[5] = true
	return b
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
func (b *CloudResourceBuilder) CloudProvider(value string) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.cloudProvider = value
	b.fieldSet_[6] = true
	return b
}

// CpuCores sets the value of the 'cpu_cores' attribute to the given value.
func (b *CloudResourceBuilder) CpuCores(value int) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.cpuCores = value
	b.fieldSet_[7] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *CloudResourceBuilder) CreatedAt(value time.Time) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.createdAt = value
	b.fieldSet_[8] = true
	return b
}

// GenericName sets the value of the 'generic_name' attribute to the given value.
func (b *CloudResourceBuilder) GenericName(value string) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.genericName = value
	b.fieldSet_[9] = true
	return b
}

// Memory sets the value of the 'memory' attribute to the given value.
func (b *CloudResourceBuilder) Memory(value int) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.memory = value
	b.fieldSet_[10] = true
	return b
}

// MemoryPretty sets the value of the 'memory_pretty' attribute to the given value.
func (b *CloudResourceBuilder) MemoryPretty(value string) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.memoryPretty = value
	b.fieldSet_[11] = true
	return b
}

// NamePretty sets the value of the 'name_pretty' attribute to the given value.
func (b *CloudResourceBuilder) NamePretty(value string) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.namePretty = value
	b.fieldSet_[12] = true
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *CloudResourceBuilder) ResourceType(value string) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.resourceType = value
	b.fieldSet_[13] = true
	return b
}

// SizePretty sets the value of the 'size_pretty' attribute to the given value.
func (b *CloudResourceBuilder) SizePretty(value string) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.sizePretty = value
	b.fieldSet_[14] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *CloudResourceBuilder) UpdatedAt(value time.Time) *CloudResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.updatedAt = value
	b.fieldSet_[15] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CloudResourceBuilder) Copy(object *CloudResource) *CloudResourceBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
