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

import (
	time "time"
)

// AddOnInstallationBuilder contains the data and logic needed to build 'add_on_installation' objects.
//
// Representation of an add-on installation in a cluster.
type AddOnInstallationBuilder struct {
	bitmap_           uint32
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
	return &AddOnInstallationBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AddOnInstallationBuilder) Link(value bool) *AddOnInstallationBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AddOnInstallationBuilder) ID(value string) *AddOnInstallationBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AddOnInstallationBuilder) HREF(value string) *AddOnInstallationBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnInstallationBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Addon sets the value of the 'addon' attribute to the given value.
//
// Representation of an add-on that can be installed in a cluster.
func (b *AddOnInstallationBuilder) Addon(value *AddOnBuilder) *AddOnInstallationBuilder {
	b.addon = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// AddonVersion sets the value of the 'addon_version' attribute to the given value.
//
// Representation of an add-on version.
func (b *AddOnInstallationBuilder) AddonVersion(value *AddOnVersionBuilder) *AddOnInstallationBuilder {
	b.addonVersion = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// Billing sets the value of the 'billing' attribute to the given value.
//
// Representation of an add-on installation billing.
func (b *AddOnInstallationBuilder) Billing(value *AddOnInstallationBillingBuilder) *AddOnInstallationBuilder {
	b.billing = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *AddOnInstallationBuilder) CreationTimestamp(value time.Time) *AddOnInstallationBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 64
	return b
}

// OperatorVersion sets the value of the 'operator_version' attribute to the given value.
func (b *AddOnInstallationBuilder) OperatorVersion(value string) *AddOnInstallationBuilder {
	b.operatorVersion = value
	b.bitmap_ |= 128
	return b
}

// Parameters sets the value of the 'parameters' attribute to the given values.
func (b *AddOnInstallationBuilder) Parameters(value *AddOnInstallationParameterListBuilder) *AddOnInstallationBuilder {
	b.parameters = value
	b.bitmap_ |= 256
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Representation of an add-on installation State field.
func (b *AddOnInstallationBuilder) State(value AddOnInstallationState) *AddOnInstallationBuilder {
	b.state = value
	b.bitmap_ |= 512
	return b
}

// StateDescription sets the value of the 'state_description' attribute to the given value.
func (b *AddOnInstallationBuilder) StateDescription(value string) *AddOnInstallationBuilder {
	b.stateDescription = value
	b.bitmap_ |= 1024
	return b
}

// UpdatedTimestamp sets the value of the 'updated_timestamp' attribute to the given value.
func (b *AddOnInstallationBuilder) UpdatedTimestamp(value time.Time) *AddOnInstallationBuilder {
	b.updatedTimestamp = value
	b.bitmap_ |= 2048
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnInstallationBuilder) Copy(object *AddOnInstallation) *AddOnInstallationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
