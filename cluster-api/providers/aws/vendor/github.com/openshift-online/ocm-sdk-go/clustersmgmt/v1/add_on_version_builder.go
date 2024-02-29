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

// AddOnVersionBuilder contains the data and logic needed to build 'add_on_version' objects.
//
// Representation of an add-on version.
type AddOnVersionBuilder struct {
	bitmap_                  uint32
	id                       string
	href                     string
	additionalCatalogSources []*AdditionalCatalogSourceBuilder
	availableUpgrades        []string
	channel                  string
	config                   *AddOnConfigBuilder
	packageImage             string
	parameters               *AddOnParameterListBuilder
	pullSecretName           string
	requirements             []*AddOnRequirementBuilder
	sourceImage              string
	subOperators             []*AddOnSubOperatorBuilder
	enabled                  bool
}

// NewAddOnVersion creates a new builder of 'add_on_version' objects.
func NewAddOnVersion() *AddOnVersionBuilder {
	return &AddOnVersionBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AddOnVersionBuilder) Link(value bool) *AddOnVersionBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AddOnVersionBuilder) ID(value string) *AddOnVersionBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AddOnVersionBuilder) HREF(value string) *AddOnVersionBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnVersionBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// AdditionalCatalogSources sets the value of the 'additional_catalog_sources' attribute to the given values.
func (b *AddOnVersionBuilder) AdditionalCatalogSources(values ...*AdditionalCatalogSourceBuilder) *AddOnVersionBuilder {
	b.additionalCatalogSources = make([]*AdditionalCatalogSourceBuilder, len(values))
	copy(b.additionalCatalogSources, values)
	b.bitmap_ |= 8
	return b
}

// AvailableUpgrades sets the value of the 'available_upgrades' attribute to the given values.
func (b *AddOnVersionBuilder) AvailableUpgrades(values ...string) *AddOnVersionBuilder {
	b.availableUpgrades = make([]string, len(values))
	copy(b.availableUpgrades, values)
	b.bitmap_ |= 16
	return b
}

// Channel sets the value of the 'channel' attribute to the given value.
func (b *AddOnVersionBuilder) Channel(value string) *AddOnVersionBuilder {
	b.channel = value
	b.bitmap_ |= 32
	return b
}

// Config sets the value of the 'config' attribute to the given value.
//
// Representation of an add-on config.
// The attributes under it are to be used by the addon once its installed in the cluster.
func (b *AddOnVersionBuilder) Config(value *AddOnConfigBuilder) *AddOnVersionBuilder {
	b.config = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddOnVersionBuilder) Enabled(value bool) *AddOnVersionBuilder {
	b.enabled = value
	b.bitmap_ |= 128
	return b
}

// PackageImage sets the value of the 'package_image' attribute to the given value.
func (b *AddOnVersionBuilder) PackageImage(value string) *AddOnVersionBuilder {
	b.packageImage = value
	b.bitmap_ |= 256
	return b
}

// Parameters sets the value of the 'parameters' attribute to the given values.
func (b *AddOnVersionBuilder) Parameters(value *AddOnParameterListBuilder) *AddOnVersionBuilder {
	b.parameters = value
	b.bitmap_ |= 512
	return b
}

// PullSecretName sets the value of the 'pull_secret_name' attribute to the given value.
func (b *AddOnVersionBuilder) PullSecretName(value string) *AddOnVersionBuilder {
	b.pullSecretName = value
	b.bitmap_ |= 1024
	return b
}

// Requirements sets the value of the 'requirements' attribute to the given values.
func (b *AddOnVersionBuilder) Requirements(values ...*AddOnRequirementBuilder) *AddOnVersionBuilder {
	b.requirements = make([]*AddOnRequirementBuilder, len(values))
	copy(b.requirements, values)
	b.bitmap_ |= 2048
	return b
}

// SourceImage sets the value of the 'source_image' attribute to the given value.
func (b *AddOnVersionBuilder) SourceImage(value string) *AddOnVersionBuilder {
	b.sourceImage = value
	b.bitmap_ |= 4096
	return b
}

// SubOperators sets the value of the 'sub_operators' attribute to the given values.
func (b *AddOnVersionBuilder) SubOperators(values ...*AddOnSubOperatorBuilder) *AddOnVersionBuilder {
	b.subOperators = make([]*AddOnSubOperatorBuilder, len(values))
	copy(b.subOperators, values)
	b.bitmap_ |= 8192
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnVersionBuilder) Copy(object *AddOnVersion) *AddOnVersionBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.additionalCatalogSources != nil {
		b.additionalCatalogSources = make([]*AdditionalCatalogSourceBuilder, len(object.additionalCatalogSources))
		for i, v := range object.additionalCatalogSources {
			b.additionalCatalogSources[i] = NewAdditionalCatalogSource().Copy(v)
		}
	} else {
		b.additionalCatalogSources = nil
	}
	if object.availableUpgrades != nil {
		b.availableUpgrades = make([]string, len(object.availableUpgrades))
		copy(b.availableUpgrades, object.availableUpgrades)
	} else {
		b.availableUpgrades = nil
	}
	b.channel = object.channel
	if object.config != nil {
		b.config = NewAddOnConfig().Copy(object.config)
	} else {
		b.config = nil
	}
	b.enabled = object.enabled
	b.packageImage = object.packageImage
	if object.parameters != nil {
		b.parameters = NewAddOnParameterList().Copy(object.parameters)
	} else {
		b.parameters = nil
	}
	b.pullSecretName = object.pullSecretName
	if object.requirements != nil {
		b.requirements = make([]*AddOnRequirementBuilder, len(object.requirements))
		for i, v := range object.requirements {
			b.requirements[i] = NewAddOnRequirement().Copy(v)
		}
	} else {
		b.requirements = nil
	}
	b.sourceImage = object.sourceImage
	if object.subOperators != nil {
		b.subOperators = make([]*AddOnSubOperatorBuilder, len(object.subOperators))
		for i, v := range object.subOperators {
			b.subOperators[i] = NewAddOnSubOperator().Copy(v)
		}
	} else {
		b.subOperators = nil
	}
	return b
}

// Build creates a 'add_on_version' object using the configuration stored in the builder.
func (b *AddOnVersionBuilder) Build() (object *AddOnVersion, err error) {
	object = new(AddOnVersion)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.additionalCatalogSources != nil {
		object.additionalCatalogSources = make([]*AdditionalCatalogSource, len(b.additionalCatalogSources))
		for i, v := range b.additionalCatalogSources {
			object.additionalCatalogSources[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.availableUpgrades != nil {
		object.availableUpgrades = make([]string, len(b.availableUpgrades))
		copy(object.availableUpgrades, b.availableUpgrades)
	}
	object.channel = b.channel
	if b.config != nil {
		object.config, err = b.config.Build()
		if err != nil {
			return
		}
	}
	object.enabled = b.enabled
	object.packageImage = b.packageImage
	if b.parameters != nil {
		object.parameters, err = b.parameters.Build()
		if err != nil {
			return
		}
	}
	object.pullSecretName = b.pullSecretName
	if b.requirements != nil {
		object.requirements = make([]*AddOnRequirement, len(b.requirements))
		for i, v := range b.requirements {
			object.requirements[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.sourceImage = b.sourceImage
	if b.subOperators != nil {
		object.subOperators = make([]*AddOnSubOperator, len(b.subOperators))
		for i, v := range b.subOperators {
			object.subOperators[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
