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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1

import (
	time "time"
)

// Representation of addon installation
type AddonInstallationBuilder struct {
	fieldSet_         []bool
	id                string
	href              string
	addon             *AddonBuilder
	addonVersion      *AddonVersionBuilder
	billing           *AddonInstallationBillingBuilder
	creationTimestamp time.Time
	csvName           string
	deletedTimestamp  time.Time
	desiredVersion    string
	operatorVersion   string
	parameters        *AddonInstallationParameterListBuilder
	state             AddonInstallationState
	stateDescription  string
	subscription      *ObjectReferenceBuilder
	updatedTimestamp  time.Time
}

// NewAddonInstallation creates a new builder of 'addon_installation' objects.
func NewAddonInstallation() *AddonInstallationBuilder {
	return &AddonInstallationBuilder{
		fieldSet_: make([]bool, 16),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AddonInstallationBuilder) Link(value bool) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AddonInstallationBuilder) ID(value string) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AddonInstallationBuilder) HREF(value string) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonInstallationBuilder) Empty() bool {
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
// Representation of an addon that can be installed in a cluster.
func (b *AddonInstallationBuilder) Addon(value *AddonBuilder) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
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
// Representation of an addon version.
func (b *AddonInstallationBuilder) AddonVersion(value *AddonVersionBuilder) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
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
func (b *AddonInstallationBuilder) Billing(value *AddonInstallationBillingBuilder) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
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
func (b *AddonInstallationBuilder) CreationTimestamp(value time.Time) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.creationTimestamp = value
	b.fieldSet_[6] = true
	return b
}

// CsvName sets the value of the 'csv_name' attribute to the given value.
func (b *AddonInstallationBuilder) CsvName(value string) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.csvName = value
	b.fieldSet_[7] = true
	return b
}

// DeletedTimestamp sets the value of the 'deleted_timestamp' attribute to the given value.
func (b *AddonInstallationBuilder) DeletedTimestamp(value time.Time) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.deletedTimestamp = value
	b.fieldSet_[8] = true
	return b
}

// DesiredVersion sets the value of the 'desired_version' attribute to the given value.
func (b *AddonInstallationBuilder) DesiredVersion(value string) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.desiredVersion = value
	b.fieldSet_[9] = true
	return b
}

// OperatorVersion sets the value of the 'operator_version' attribute to the given value.
func (b *AddonInstallationBuilder) OperatorVersion(value string) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.operatorVersion = value
	b.fieldSet_[10] = true
	return b
}

// Parameters sets the value of the 'parameters' attribute to the given values.
func (b *AddonInstallationBuilder) Parameters(value *AddonInstallationParameterListBuilder) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.parameters = value
	b.fieldSet_[11] = true
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// representation of addon installation state
func (b *AddonInstallationBuilder) State(value AddonInstallationState) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.state = value
	b.fieldSet_[12] = true
	return b
}

// StateDescription sets the value of the 'state_description' attribute to the given value.
func (b *AddonInstallationBuilder) StateDescription(value string) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.stateDescription = value
	b.fieldSet_[13] = true
	return b
}

// Subscription sets the value of the 'subscription' attribute to the given value.
//
// representation of object reference/subscription
func (b *AddonInstallationBuilder) Subscription(value *ObjectReferenceBuilder) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.subscription = value
	if value != nil {
		b.fieldSet_[14] = true
	} else {
		b.fieldSet_[14] = false
	}
	return b
}

// UpdatedTimestamp sets the value of the 'updated_timestamp' attribute to the given value.
func (b *AddonInstallationBuilder) UpdatedTimestamp(value time.Time) *AddonInstallationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.updatedTimestamp = value
	b.fieldSet_[15] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonInstallationBuilder) Copy(object *AddonInstallation) *AddonInstallationBuilder {
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
		b.addon = NewAddon().Copy(object.addon)
	} else {
		b.addon = nil
	}
	if object.addonVersion != nil {
		b.addonVersion = NewAddonVersion().Copy(object.addonVersion)
	} else {
		b.addonVersion = nil
	}
	if object.billing != nil {
		b.billing = NewAddonInstallationBilling().Copy(object.billing)
	} else {
		b.billing = nil
	}
	b.creationTimestamp = object.creationTimestamp
	b.csvName = object.csvName
	b.deletedTimestamp = object.deletedTimestamp
	b.desiredVersion = object.desiredVersion
	b.operatorVersion = object.operatorVersion
	if object.parameters != nil {
		b.parameters = NewAddonInstallationParameterList().Copy(object.parameters)
	} else {
		b.parameters = nil
	}
	b.state = object.state
	b.stateDescription = object.stateDescription
	if object.subscription != nil {
		b.subscription = NewObjectReference().Copy(object.subscription)
	} else {
		b.subscription = nil
	}
	b.updatedTimestamp = object.updatedTimestamp
	return b
}

// Build creates a 'addon_installation' object using the configuration stored in the builder.
func (b *AddonInstallationBuilder) Build() (object *AddonInstallation, err error) {
	object = new(AddonInstallation)
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
	object.csvName = b.csvName
	object.deletedTimestamp = b.deletedTimestamp
	object.desiredVersion = b.desiredVersion
	object.operatorVersion = b.operatorVersion
	if b.parameters != nil {
		object.parameters, err = b.parameters.Build()
		if err != nil {
			return
		}
	}
	object.state = b.state
	object.stateDescription = b.stateDescription
	if b.subscription != nil {
		object.subscription, err = b.subscription.Build()
		if err != nil {
			return
		}
	}
	object.updatedTimestamp = b.updatedTimestamp
	return
}
