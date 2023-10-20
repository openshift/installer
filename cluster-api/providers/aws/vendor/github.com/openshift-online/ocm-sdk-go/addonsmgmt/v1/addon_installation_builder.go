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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

import (
	time "time"
)

// AddonInstallationBuilder contains the data and logic needed to build 'addon_installation' objects.
//
// Representation of addon installation
type AddonInstallationBuilder struct {
	bitmap_           uint32
	id                string
	href              string
	addon             *AddonBuilder
	addonVersion      *AddonVersionBuilder
	billing           *AddonInstallationBillingBuilder
	creationTimestamp time.Time
	csvName           string
	deletedTimestamp  time.Time
	operatorVersion   string
	parameters        *AddonInstallationParametersBuilder
	state             AddonInstallationState
	stateDescription  string
	subscription      *ObjectReferenceBuilder
	updatedTimestamp  time.Time
}

// NewAddonInstallation creates a new builder of 'addon_installation' objects.
func NewAddonInstallation() *AddonInstallationBuilder {
	return &AddonInstallationBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AddonInstallationBuilder) Link(value bool) *AddonInstallationBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AddonInstallationBuilder) ID(value string) *AddonInstallationBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AddonInstallationBuilder) HREF(value string) *AddonInstallationBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonInstallationBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Addon sets the value of the 'addon' attribute to the given value.
//
// Representation of an addon that can be installed in a cluster.
func (b *AddonInstallationBuilder) Addon(value *AddonBuilder) *AddonInstallationBuilder {
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
// Representation of an addon version.
func (b *AddonInstallationBuilder) AddonVersion(value *AddonVersionBuilder) *AddonInstallationBuilder {
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
func (b *AddonInstallationBuilder) Billing(value *AddonInstallationBillingBuilder) *AddonInstallationBuilder {
	b.billing = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *AddonInstallationBuilder) CreationTimestamp(value time.Time) *AddonInstallationBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 64
	return b
}

// CsvName sets the value of the 'csv_name' attribute to the given value.
func (b *AddonInstallationBuilder) CsvName(value string) *AddonInstallationBuilder {
	b.csvName = value
	b.bitmap_ |= 128
	return b
}

// DeletedTimestamp sets the value of the 'deleted_timestamp' attribute to the given value.
func (b *AddonInstallationBuilder) DeletedTimestamp(value time.Time) *AddonInstallationBuilder {
	b.deletedTimestamp = value
	b.bitmap_ |= 256
	return b
}

// OperatorVersion sets the value of the 'operator_version' attribute to the given value.
func (b *AddonInstallationBuilder) OperatorVersion(value string) *AddonInstallationBuilder {
	b.operatorVersion = value
	b.bitmap_ |= 512
	return b
}

// Parameters sets the value of the 'parameters' attribute to the given value.
//
// representation of addon installation parameter
func (b *AddonInstallationBuilder) Parameters(value *AddonInstallationParametersBuilder) *AddonInstallationBuilder {
	b.parameters = value
	if value != nil {
		b.bitmap_ |= 1024
	} else {
		b.bitmap_ &^= 1024
	}
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// representation of addon installation state
func (b *AddonInstallationBuilder) State(value AddonInstallationState) *AddonInstallationBuilder {
	b.state = value
	b.bitmap_ |= 2048
	return b
}

// StateDescription sets the value of the 'state_description' attribute to the given value.
func (b *AddonInstallationBuilder) StateDescription(value string) *AddonInstallationBuilder {
	b.stateDescription = value
	b.bitmap_ |= 4096
	return b
}

// Subscription sets the value of the 'subscription' attribute to the given value.
//
// representation of object reference/subscription
func (b *AddonInstallationBuilder) Subscription(value *ObjectReferenceBuilder) *AddonInstallationBuilder {
	b.subscription = value
	if value != nil {
		b.bitmap_ |= 8192
	} else {
		b.bitmap_ &^= 8192
	}
	return b
}

// UpdatedTimestamp sets the value of the 'updated_timestamp' attribute to the given value.
func (b *AddonInstallationBuilder) UpdatedTimestamp(value time.Time) *AddonInstallationBuilder {
	b.updatedTimestamp = value
	b.bitmap_ |= 16384
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonInstallationBuilder) Copy(object *AddonInstallation) *AddonInstallationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	b.operatorVersion = object.operatorVersion
	if object.parameters != nil {
		b.parameters = NewAddonInstallationParameters().Copy(object.parameters)
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
	object.csvName = b.csvName
	object.deletedTimestamp = b.deletedTimestamp
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
