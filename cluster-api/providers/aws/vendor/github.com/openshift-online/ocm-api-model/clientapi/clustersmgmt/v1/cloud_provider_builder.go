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

// Cloud provider.
type CloudProviderBuilder struct {
	fieldSet_   []bool
	id          string
	href        string
	displayName string
	name        string
	regions     []*CloudRegionBuilder
}

// NewCloudProvider creates a new builder of 'cloud_provider' objects.
func NewCloudProvider() *CloudProviderBuilder {
	return &CloudProviderBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *CloudProviderBuilder) Link(value bool) *CloudProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *CloudProviderBuilder) ID(value string) *CloudProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *CloudProviderBuilder) HREF(value string) *CloudProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CloudProviderBuilder) Empty() bool {
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

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *CloudProviderBuilder) DisplayName(value string) *CloudProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.displayName = value
	b.fieldSet_[3] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *CloudProviderBuilder) Name(value string) *CloudProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.name = value
	b.fieldSet_[4] = true
	return b
}

// Regions sets the value of the 'regions' attribute to the given values.
func (b *CloudProviderBuilder) Regions(values ...*CloudRegionBuilder) *CloudProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.regions = make([]*CloudRegionBuilder, len(values))
	copy(b.regions, values)
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CloudProviderBuilder) Copy(object *CloudProvider) *CloudProviderBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.displayName = object.displayName
	b.name = object.name
	if object.regions != nil {
		b.regions = make([]*CloudRegionBuilder, len(object.regions))
		for i, v := range object.regions {
			b.regions[i] = NewCloudRegion().Copy(v)
		}
	} else {
		b.regions = nil
	}
	return b
}

// Build creates a 'cloud_provider' object using the configuration stored in the builder.
func (b *CloudProviderBuilder) Build() (object *CloudProvider, err error) {
	object = new(CloudProvider)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.displayName = b.displayName
	object.name = b.name
	if b.regions != nil {
		object.regions = make([]*CloudRegion, len(b.regions))
		for i, v := range b.regions {
			object.regions[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
