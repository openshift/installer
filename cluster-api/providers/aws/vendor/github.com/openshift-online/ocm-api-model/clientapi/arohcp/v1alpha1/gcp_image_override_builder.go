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

// GcpImageOverride specifies what a GCP VM Image should be used for a particular product and billing model
type GCPImageOverrideBuilder struct {
	fieldSet_    []bool
	id           string
	href         string
	billingModel *BillingModelItemBuilder
	imageID      string
	product      *ProductBuilder
	projectID    string
}

// NewGCPImageOverride creates a new builder of 'GCP_image_override' objects.
func NewGCPImageOverride() *GCPImageOverrideBuilder {
	return &GCPImageOverrideBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *GCPImageOverrideBuilder) Link(value bool) *GCPImageOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *GCPImageOverrideBuilder) ID(value string) *GCPImageOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *GCPImageOverrideBuilder) HREF(value string) *GCPImageOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GCPImageOverrideBuilder) Empty() bool {
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

// BillingModel sets the value of the 'billing_model' attribute to the given value.
//
// BillingModelItem represents a billing model
func (b *GCPImageOverrideBuilder) BillingModel(value *BillingModelItemBuilder) *GCPImageOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.billingModel = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// ImageID sets the value of the 'image_ID' attribute to the given value.
func (b *GCPImageOverrideBuilder) ImageID(value string) *GCPImageOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.imageID = value
	b.fieldSet_[4] = true
	return b
}

// Product sets the value of the 'product' attribute to the given value.
//
// Representation of an product that can be selected as a cluster type.
func (b *GCPImageOverrideBuilder) Product(value *ProductBuilder) *GCPImageOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.product = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// ProjectID sets the value of the 'project_ID' attribute to the given value.
func (b *GCPImageOverrideBuilder) ProjectID(value string) *GCPImageOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.projectID = value
	b.fieldSet_[6] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GCPImageOverrideBuilder) Copy(object *GCPImageOverride) *GCPImageOverrideBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.billingModel != nil {
		b.billingModel = NewBillingModelItem().Copy(object.billingModel)
	} else {
		b.billingModel = nil
	}
	b.imageID = object.imageID
	if object.product != nil {
		b.product = NewProduct().Copy(object.product)
	} else {
		b.product = nil
	}
	b.projectID = object.projectID
	return b
}

// Build creates a 'GCP_image_override' object using the configuration stored in the builder.
func (b *GCPImageOverrideBuilder) Build() (object *GCPImageOverride, err error) {
	object = new(GCPImageOverride)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.billingModel != nil {
		object.billingModel, err = b.billingModel.Build()
		if err != nil {
			return
		}
	}
	object.imageID = b.imageID
	if b.product != nil {
		object.product, err = b.product.Build()
		if err != nil {
			return
		}
	}
	object.projectID = b.projectID
	return
}
