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

// Representation of an addon that can be installed in a cluster.
type AddonBuilder struct {
	fieldSet_            []bool
	id                   string
	href                 string
	commonAnnotations    map[string]string
	commonLabels         map[string]string
	config               *AddonConfigBuilder
	credentialsRequests  []*CredentialRequestBuilder
	description          string
	docsLink             string
	icon                 string
	installMode          AddonInstallMode
	label                string
	name                 string
	namespaces           []*AddonNamespaceBuilder
	operatorName         string
	parameters           *AddonParameterListBuilder
	requirements         []*AddonRequirementBuilder
	resourceCost         float64
	resourceName         string
	subOperators         []*AddonSubOperatorBuilder
	targetNamespace      string
	version              *AddonVersionBuilder
	enabled              bool
	hasExternalResources bool
	hidden               bool
	managedService       bool
}

// NewAddon creates a new builder of 'addon' objects.
func NewAddon() *AddonBuilder {
	return &AddonBuilder{
		fieldSet_: make([]bool, 26),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AddonBuilder) Link(value bool) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AddonBuilder) ID(value string) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AddonBuilder) HREF(value string) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonBuilder) Empty() bool {
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

// CommonAnnotations sets the value of the 'common_annotations' attribute to the given value.
func (b *AddonBuilder) CommonAnnotations(value map[string]string) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.commonAnnotations = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// CommonLabels sets the value of the 'common_labels' attribute to the given value.
func (b *AddonBuilder) CommonLabels(value map[string]string) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.commonLabels = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// Config sets the value of the 'config' attribute to the given value.
//
// Representation of an addon config.
// The attributes under it are to be used by the addon once its installed in the cluster.
func (b *AddonBuilder) Config(value *AddonConfigBuilder) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.config = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// CredentialsRequests sets the value of the 'credentials_requests' attribute to the given values.
func (b *AddonBuilder) CredentialsRequests(values ...*CredentialRequestBuilder) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.credentialsRequests = make([]*CredentialRequestBuilder, len(values))
	copy(b.credentialsRequests, values)
	b.fieldSet_[6] = true
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *AddonBuilder) Description(value string) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.description = value
	b.fieldSet_[7] = true
	return b
}

// DocsLink sets the value of the 'docs_link' attribute to the given value.
func (b *AddonBuilder) DocsLink(value string) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.docsLink = value
	b.fieldSet_[8] = true
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddonBuilder) Enabled(value bool) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.enabled = value
	b.fieldSet_[9] = true
	return b
}

// HasExternalResources sets the value of the 'has_external_resources' attribute to the given value.
func (b *AddonBuilder) HasExternalResources(value bool) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.hasExternalResources = value
	b.fieldSet_[10] = true
	return b
}

// Hidden sets the value of the 'hidden' attribute to the given value.
func (b *AddonBuilder) Hidden(value bool) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.hidden = value
	b.fieldSet_[11] = true
	return b
}

// Icon sets the value of the 'icon' attribute to the given value.
func (b *AddonBuilder) Icon(value string) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.icon = value
	b.fieldSet_[12] = true
	return b
}

// InstallMode sets the value of the 'install_mode' attribute to the given value.
//
// Representation of an addon InstallMode field.
func (b *AddonBuilder) InstallMode(value AddonInstallMode) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.installMode = value
	b.fieldSet_[13] = true
	return b
}

// Label sets the value of the 'label' attribute to the given value.
func (b *AddonBuilder) Label(value string) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.label = value
	b.fieldSet_[14] = true
	return b
}

// ManagedService sets the value of the 'managed_service' attribute to the given value.
func (b *AddonBuilder) ManagedService(value bool) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.managedService = value
	b.fieldSet_[15] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddonBuilder) Name(value string) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.name = value
	b.fieldSet_[16] = true
	return b
}

// Namespaces sets the value of the 'namespaces' attribute to the given values.
func (b *AddonBuilder) Namespaces(values ...*AddonNamespaceBuilder) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.namespaces = make([]*AddonNamespaceBuilder, len(values))
	copy(b.namespaces, values)
	b.fieldSet_[17] = true
	return b
}

// OperatorName sets the value of the 'operator_name' attribute to the given value.
func (b *AddonBuilder) OperatorName(value string) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.operatorName = value
	b.fieldSet_[18] = true
	return b
}

// Parameters sets the value of the 'parameters' attribute to the given values.
func (b *AddonBuilder) Parameters(value *AddonParameterListBuilder) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.parameters = value
	b.fieldSet_[19] = true
	return b
}

// Requirements sets the value of the 'requirements' attribute to the given values.
func (b *AddonBuilder) Requirements(values ...*AddonRequirementBuilder) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.requirements = make([]*AddonRequirementBuilder, len(values))
	copy(b.requirements, values)
	b.fieldSet_[20] = true
	return b
}

// ResourceCost sets the value of the 'resource_cost' attribute to the given value.
func (b *AddonBuilder) ResourceCost(value float64) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.resourceCost = value
	b.fieldSet_[21] = true
	return b
}

// ResourceName sets the value of the 'resource_name' attribute to the given value.
func (b *AddonBuilder) ResourceName(value string) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.resourceName = value
	b.fieldSet_[22] = true
	return b
}

// SubOperators sets the value of the 'sub_operators' attribute to the given values.
func (b *AddonBuilder) SubOperators(values ...*AddonSubOperatorBuilder) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.subOperators = make([]*AddonSubOperatorBuilder, len(values))
	copy(b.subOperators, values)
	b.fieldSet_[23] = true
	return b
}

// TargetNamespace sets the value of the 'target_namespace' attribute to the given value.
func (b *AddonBuilder) TargetNamespace(value string) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.targetNamespace = value
	b.fieldSet_[24] = true
	return b
}

// Version sets the value of the 'version' attribute to the given value.
//
// Representation of an addon version.
func (b *AddonBuilder) Version(value *AddonVersionBuilder) *AddonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.version = value
	if value != nil {
		b.fieldSet_[25] = true
	} else {
		b.fieldSet_[25] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonBuilder) Copy(object *Addon) *AddonBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if len(object.commonAnnotations) > 0 {
		b.commonAnnotations = map[string]string{}
		for k, v := range object.commonAnnotations {
			b.commonAnnotations[k] = v
		}
	} else {
		b.commonAnnotations = nil
	}
	if len(object.commonLabels) > 0 {
		b.commonLabels = map[string]string{}
		for k, v := range object.commonLabels {
			b.commonLabels[k] = v
		}
	} else {
		b.commonLabels = nil
	}
	if object.config != nil {
		b.config = NewAddonConfig().Copy(object.config)
	} else {
		b.config = nil
	}
	if object.credentialsRequests != nil {
		b.credentialsRequests = make([]*CredentialRequestBuilder, len(object.credentialsRequests))
		for i, v := range object.credentialsRequests {
			b.credentialsRequests[i] = NewCredentialRequest().Copy(v)
		}
	} else {
		b.credentialsRequests = nil
	}
	b.description = object.description
	b.docsLink = object.docsLink
	b.enabled = object.enabled
	b.hasExternalResources = object.hasExternalResources
	b.hidden = object.hidden
	b.icon = object.icon
	b.installMode = object.installMode
	b.label = object.label
	b.managedService = object.managedService
	b.name = object.name
	if object.namespaces != nil {
		b.namespaces = make([]*AddonNamespaceBuilder, len(object.namespaces))
		for i, v := range object.namespaces {
			b.namespaces[i] = NewAddonNamespace().Copy(v)
		}
	} else {
		b.namespaces = nil
	}
	b.operatorName = object.operatorName
	if object.parameters != nil {
		b.parameters = NewAddonParameterList().Copy(object.parameters)
	} else {
		b.parameters = nil
	}
	if object.requirements != nil {
		b.requirements = make([]*AddonRequirementBuilder, len(object.requirements))
		for i, v := range object.requirements {
			b.requirements[i] = NewAddonRequirement().Copy(v)
		}
	} else {
		b.requirements = nil
	}
	b.resourceCost = object.resourceCost
	b.resourceName = object.resourceName
	if object.subOperators != nil {
		b.subOperators = make([]*AddonSubOperatorBuilder, len(object.subOperators))
		for i, v := range object.subOperators {
			b.subOperators[i] = NewAddonSubOperator().Copy(v)
		}
	} else {
		b.subOperators = nil
	}
	b.targetNamespace = object.targetNamespace
	if object.version != nil {
		b.version = NewAddonVersion().Copy(object.version)
	} else {
		b.version = nil
	}
	return b
}

// Build creates a 'addon' object using the configuration stored in the builder.
func (b *AddonBuilder) Build() (object *Addon, err error) {
	object = new(Addon)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.commonAnnotations != nil {
		object.commonAnnotations = make(map[string]string)
		for k, v := range b.commonAnnotations {
			object.commonAnnotations[k] = v
		}
	}
	if b.commonLabels != nil {
		object.commonLabels = make(map[string]string)
		for k, v := range b.commonLabels {
			object.commonLabels[k] = v
		}
	}
	if b.config != nil {
		object.config, err = b.config.Build()
		if err != nil {
			return
		}
	}
	if b.credentialsRequests != nil {
		object.credentialsRequests = make([]*CredentialRequest, len(b.credentialsRequests))
		for i, v := range b.credentialsRequests {
			object.credentialsRequests[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.description = b.description
	object.docsLink = b.docsLink
	object.enabled = b.enabled
	object.hasExternalResources = b.hasExternalResources
	object.hidden = b.hidden
	object.icon = b.icon
	object.installMode = b.installMode
	object.label = b.label
	object.managedService = b.managedService
	object.name = b.name
	if b.namespaces != nil {
		object.namespaces = make([]*AddonNamespace, len(b.namespaces))
		for i, v := range b.namespaces {
			object.namespaces[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.operatorName = b.operatorName
	if b.parameters != nil {
		object.parameters, err = b.parameters.Build()
		if err != nil {
			return
		}
	}
	if b.requirements != nil {
		object.requirements = make([]*AddonRequirement, len(b.requirements))
		for i, v := range b.requirements {
			object.requirements[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.resourceCost = b.resourceCost
	object.resourceName = b.resourceName
	if b.subOperators != nil {
		object.subOperators = make([]*AddonSubOperator, len(b.subOperators))
		for i, v := range b.subOperators {
			object.subOperators[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.targetNamespace = b.targetNamespace
	if b.version != nil {
		object.version, err = b.version.Build()
		if err != nil {
			return
		}
	}
	return
}
