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

import (
	time "time"
)

// Representation of an add-on installation in a cluster.
type AddOnInstallationBuilder struct {
	fieldSet_         []bool
	id                string
	href              string
	addon             *AddOnBuilder
	addonVersion      *AddOnVersionBuilder
	billing           *AddOnInstallationBillingBuilder
	creationTimestamp time.Time
	operatorVersion   string
	parameters        *AddOnInstallationParameterListBuilder
	state             AddOnInstallationState
	stateDescription  string
	updatedTimestamp  time.Time
}

// NewAddOnInstallation creates a new builder of 'add_on_installation' objects.
func NewAddOnInstallation() *AddOnInstallationBuilder {
	return &AddOnInstallationBuilder{
		fieldSet_: make([]bool, 12),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AddOnInstallationBuilder) Link(value bool) *AddOnInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AddOnInstallationBuilder) ID(value string) *AddOnInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AddOnInstallationBuilder) HREF(value string) *AddOnInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnInstallationBuilder) Empty() bool {
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

// Addon sets the value of the 'addon' attribute to the given value.
//
// Representation of an add-on that can be installed in a cluster.
func (b *AddOnInstallationBuilder) Addon(value *AddOnBuilder) *AddOnInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.addon = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// AddonVersion sets the value of the 'addon_version' attribute to the given value.
//
// Representation of an add-on version.
func (b *AddOnInstallationBuilder) AddonVersion(value *AddOnVersionBuilder) *AddOnInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.addonVersion = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// Billing sets the value of the 'billing' attribute to the given value.
//
// Representation of an add-on installation billing.
func (b *AddOnInstallationBuilder) Billing(value *AddOnInstallationBillingBuilder) *AddOnInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.billing = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *AddOnInstallationBuilder) CreationTimestamp(value time.Time) *AddOnInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.creationTimestamp = value
	b.fieldSet_[6] = true
	return b
}

// OperatorVersion sets the value of the 'operator_version' attribute to the given value.
func (b *AddOnInstallationBuilder) OperatorVersion(value string) *AddOnInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.operatorVersion = value
	b.fieldSet_[7] = true
	return b
}

// Parameters sets the value of the 'parameters' attribute to the given values.
func (b *AddOnInstallationBuilder) Parameters(value *AddOnInstallationParameterListBuilder) *AddOnInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.parameters = value
	b.fieldSet_[8] = true
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Representation of an add-on installation State field.
func (b *AddOnInstallationBuilder) State(value AddOnInstallationState) *AddOnInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.state = value
	b.fieldSet_[9] = true
	return b
}

// StateDescription sets the value of the 'state_description' attribute to the given value.
func (b *AddOnInstallationBuilder) StateDescription(value string) *AddOnInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.stateDescription = value
	b.fieldSet_[10] = true
	return b
}

// UpdatedTimestamp sets the value of the 'updated_timestamp' attribute to the given value.
func (b *AddOnInstallationBuilder) UpdatedTimestamp(value time.Time) *AddOnInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.updatedTimestamp = value
	b.fieldSet_[11] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnInstallationBuilder) Copy(object *AddOnInstallation) *AddOnInstallationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.addon != nil {
		b.addon = NewAddOn().Copy(object.addon)
	} else {
		b.addon = nil
	}
	if object.addonVersion != nil {
		b.addonVersion = NewAddOnVersion().Copy(object.addonVersion)
	} else {
		b.addonVersion = nil
	}
	if object.billing != nil {
		b.billing = NewAddOnInstallationBilling().Copy(object.billing)
	} else {
		b.billing = nil
	}
	b.creationTimestamp = object.creationTimestamp
	b.operatorVersion = object.operatorVersion
	if object.parameters != nil {
		b.parameters = NewAddOnInstallationParameterList().Copy(object.parameters)
	} else {
		b.parameters = nil
	}
	b.state = object.state
	b.stateDescription = object.stateDescription
	b.updatedTimestamp = object.updatedTimestamp
	return b
}

// Build creates a 'add_on_installation' object using the configuration stored in the builder.
func (b *AddOnInstallationBuilder) Build() (object *AddOnInstallation, err error) {
	object = new(AddOnInstallation)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.addon != nil {
		object.addon, err = b.addon.Build()
		if err != nil {
			return
		}
	}
	if b.addonVersion != nil {
		object.addonVersion, err = b.addonVersion.Build()
		if err != nil {
			return
		}
	}
	if b.billing != nil {
		object.billing, err = b.billing.Build()
		if err != nil {
			return
		}
	}
	object.creationTimestamp = b.creationTimestamp
	object.operatorVersion = b.operatorVersion
	if b.parameters != nil {
		object.parameters, err = b.parameters.Build()
		if err != nil {
			return
		}
	}
	object.state = b.state
	object.stateDescription = b.stateDescription
	object.updatedTimestamp = b.updatedTimestamp
	return
}
