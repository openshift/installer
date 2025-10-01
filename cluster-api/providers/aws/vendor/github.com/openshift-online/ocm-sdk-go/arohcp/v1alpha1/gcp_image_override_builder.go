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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// GCPImageOverrideBuilder contains the data and logic needed to build 'GCP_image_override' objects.
//
// GcpImageOverride specifies what a GCP VM Image should be used for a particular product and billing model
type GCPImageOverrideBuilder struct {
	bitmap_      uint32
	id           string
	href         string
	billingModel *v1.BillingModelItemBuilder
	imageID      string
	product      *v1.ProductBuilder
	projectID    string
}

// NewGCPImageOverride creates a new builder of 'GCP_image_override' objects.
func NewGCPImageOverride() *GCPImageOverrideBuilder {
	return &GCPImageOverrideBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *GCPImageOverrideBuilder) Link(value bool) *GCPImageOverrideBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *GCPImageOverrideBuilder) ID(value string) *GCPImageOverrideBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *GCPImageOverrideBuilder) HREF(value string) *GCPImageOverrideBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GCPImageOverrideBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// BillingModel sets the value of the 'billing_model' attribute to the given value.
//
// BillingModelItem represents a billing model
func (b *GCPImageOverrideBuilder) BillingModel(value *v1.BillingModelItemBuilder) *GCPImageOverrideBuilder {
	b.billingModel = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// ImageID sets the value of the 'image_ID' attribute to the given value.
func (b *GCPImageOverrideBuilder) ImageID(value string) *GCPImageOverrideBuilder {
	b.imageID = value
	b.bitmap_ |= 16
	return b
}

// Product sets the value of the 'product' attribute to the given value.
//
// Representation of an product that can be selected as a cluster type.
func (b *GCPImageOverrideBuilder) Product(value *v1.ProductBuilder) *GCPImageOverrideBuilder {
	b.product = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// ProjectID sets the value of the 'project_ID' attribute to the given value.
func (b *GCPImageOverrideBuilder) ProjectID(value string) *GCPImageOverrideBuilder {
	b.projectID = value
	b.bitmap_ |= 64
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GCPImageOverrideBuilder) Copy(object *GCPImageOverride) *GCPImageOverrideBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.billingModel != nil {
		b.billingModel = v1.NewBillingModelItem().Copy(object.billingModel)
	} else {
		b.billingModel = nil
	}
	b.imageID = object.imageID
	if object.product != nil {
		b.product = v1.NewProduct().Copy(object.product)
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
	object.bitmap_ = b.bitmap_
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
