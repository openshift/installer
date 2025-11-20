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

// Representation of an addon version.
type AddonVersionBuilder struct {
	fieldSet_                []bool
	id                       string
	href                     string
	additionalCatalogSources []*AdditionalCatalogSourceBuilder
	availableUpgrades        []string
	channel                  string
	config                   *AddonConfigBuilder
	metricsFederation        *MetricsFederationBuilder
	monitoringStack          *MonitoringStackBuilder
	packageImage             string
	parameters               *AddonParametersBuilder
	pullSecretName           string
	requirements             []*AddonRequirementBuilder
	sourceImage              string
	subOperators             []*AddonSubOperatorBuilder
	enabled                  bool
	upgradePlansCreated      bool
}

// NewAddonVersion creates a new builder of 'addon_version' objects.
func NewAddonVersion() *AddonVersionBuilder {
	return &AddonVersionBuilder{
		fieldSet_: make([]bool, 17),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AddonVersionBuilder) Link(value bool) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AddonVersionBuilder) ID(value string) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AddonVersionBuilder) HREF(value string) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonVersionBuilder) Empty() bool {
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

// AdditionalCatalogSources sets the value of the 'additional_catalog_sources' attribute to the given values.
func (b *AddonVersionBuilder) AdditionalCatalogSources(values ...*AdditionalCatalogSourceBuilder) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.additionalCatalogSources = make([]*AdditionalCatalogSourceBuilder, len(values))
	copy(b.additionalCatalogSources, values)
	b.fieldSet_[3] = true
	return b
}

// AvailableUpgrades sets the value of the 'available_upgrades' attribute to the given values.
func (b *AddonVersionBuilder) AvailableUpgrades(values ...string) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.availableUpgrades = make([]string, len(values))
	copy(b.availableUpgrades, values)
	b.fieldSet_[4] = true
	return b
}

// Channel sets the value of the 'channel' attribute to the given value.
func (b *AddonVersionBuilder) Channel(value string) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.channel = value
	b.fieldSet_[5] = true
	return b
}

// Config sets the value of the 'config' attribute to the given value.
//
// Representation of an addon config.
// The attributes under it are to be used by the addon once its installed in the cluster.
func (b *AddonVersionBuilder) Config(value *AddonConfigBuilder) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.config = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddonVersionBuilder) Enabled(value bool) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.enabled = value
	b.fieldSet_[7] = true
	return b
}

// MetricsFederation sets the value of the 'metrics_federation' attribute to the given value.
//
// Representation of Metrics Federation
func (b *AddonVersionBuilder) MetricsFederation(value *MetricsFederationBuilder) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.metricsFederation = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// MonitoringStack sets the value of the 'monitoring_stack' attribute to the given value.
//
// Representation of Monitoring Stack
func (b *AddonVersionBuilder) MonitoringStack(value *MonitoringStackBuilder) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.monitoringStack = value
	if value != nil {
		b.fieldSet_[9] = true
	} else {
		b.fieldSet_[9] = false
	}
	return b
}

// PackageImage sets the value of the 'package_image' attribute to the given value.
func (b *AddonVersionBuilder) PackageImage(value string) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.packageImage = value
	b.fieldSet_[10] = true
	return b
}

// Parameters sets the value of the 'parameters' attribute to the given value.
//
// Representation of AddonParameters
func (b *AddonVersionBuilder) Parameters(value *AddonParametersBuilder) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.parameters = value
	if value != nil {
		b.fieldSet_[11] = true
	} else {
		b.fieldSet_[11] = false
	}
	return b
}

// PullSecretName sets the value of the 'pull_secret_name' attribute to the given value.
func (b *AddonVersionBuilder) PullSecretName(value string) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.pullSecretName = value
	b.fieldSet_[12] = true
	return b
}

// Requirements sets the value of the 'requirements' attribute to the given values.
func (b *AddonVersionBuilder) Requirements(values ...*AddonRequirementBuilder) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.requirements = make([]*AddonRequirementBuilder, len(values))
	copy(b.requirements, values)
	b.fieldSet_[13] = true
	return b
}

// SourceImage sets the value of the 'source_image' attribute to the given value.
func (b *AddonVersionBuilder) SourceImage(value string) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.sourceImage = value
	b.fieldSet_[14] = true
	return b
}

// SubOperators sets the value of the 'sub_operators' attribute to the given values.
func (b *AddonVersionBuilder) SubOperators(values ...*AddonSubOperatorBuilder) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.subOperators = make([]*AddonSubOperatorBuilder, len(values))
	copy(b.subOperators, values)
	b.fieldSet_[15] = true
	return b
}

// UpgradePlansCreated sets the value of the 'upgrade_plans_created' attribute to the given value.
func (b *AddonVersionBuilder) UpgradePlansCreated(value bool) *AddonVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.upgradePlansCreated = value
	b.fieldSet_[16] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonVersionBuilder) Copy(object *AddonVersion) *AddonVersionBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
		b.config = NewAddonConfig().Copy(object.config)
	} else {
		b.config = nil
	}
	b.enabled = object.enabled
	if object.metricsFederation != nil {
		b.metricsFederation = NewMetricsFederation().Copy(object.metricsFederation)
	} else {
		b.metricsFederation = nil
	}
	if object.monitoringStack != nil {
		b.monitoringStack = NewMonitoringStack().Copy(object.monitoringStack)
	} else {
		b.monitoringStack = nil
	}
	b.packageImage = object.packageImage
	if object.parameters != nil {
		b.parameters = NewAddonParameters().Copy(object.parameters)
	} else {
		b.parameters = nil
	}
	b.pullSecretName = object.pullSecretName
	if object.requirements != nil {
		b.requirements = make([]*AddonRequirementBuilder, len(object.requirements))
		for i, v := range object.requirements {
			b.requirements[i] = NewAddonRequirement().Copy(v)
		}
	} else {
		b.requirements = nil
	}
	b.sourceImage = object.sourceImage
	if object.subOperators != nil {
		b.subOperators = make([]*AddonSubOperatorBuilder, len(object.subOperators))
		for i, v := range object.subOperators {
			b.subOperators[i] = NewAddonSubOperator().Copy(v)
		}
	} else {
		b.subOperators = nil
	}
	b.upgradePlansCreated = object.upgradePlansCreated
	return b
}

// Build creates a 'addon_version' object using the configuration stored in the builder.
func (b *AddonVersionBuilder) Build() (object *AddonVersion, err error) {
	object = new(AddonVersion)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
	if b.metricsFederation != nil {
		object.metricsFederation, err = b.metricsFederation.Build()
		if err != nil {
			return
		}
	}
	if b.monitoringStack != nil {
		object.monitoringStack, err = b.monitoringStack.Build()
		if err != nil {
			return
		}
	}
	object.packageImage = b.packageImage
	if b.parameters != nil {
		object.parameters, err = b.parameters.Build()
		if err != nil {
			return
		}
	}
	object.pullSecretName = b.pullSecretName
	if b.requirements != nil {
		object.requirements = make([]*AddonRequirement, len(b.requirements))
		for i, v := range b.requirements {
			object.requirements[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.sourceImage = b.sourceImage
	if b.subOperators != nil {
		object.subOperators = make([]*AddonSubOperator, len(b.subOperators))
		for i, v := range b.subOperators {
			object.subOperators[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.upgradePlansCreated = b.upgradePlansCreated
	return
}
