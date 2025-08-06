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

// Representation of an add-on that can be installed in a cluster.
type AddOnBuilder struct {
	fieldSet_            []bool
	id                   string
	href                 string
	commonAnnotations    map[string]string
	commonLabels         map[string]string
	config               *AddOnConfigBuilder
	credentialsRequests  []*CredentialRequestBuilder
	description          string
	docsLink             string
	icon                 string
	installMode          AddOnInstallMode
	label                string
	name                 string
	namespaces           []*AddOnNamespaceBuilder
	operatorName         string
	parameters           *AddOnParameterListBuilder
	requirements         []*AddOnRequirementBuilder
	resourceCost         float64
	resourceName         string
	subOperators         []*AddOnSubOperatorBuilder
	targetNamespace      string
	version              *AddOnVersionBuilder
	enabled              bool
	hasExternalResources bool
	hidden               bool
	managedService       bool
}

// NewAddOn creates a new builder of 'add_on' objects.
func NewAddOn() *AddOnBuilder {
	return &AddOnBuilder{
		fieldSet_: make([]bool, 26),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AddOnBuilder) Link(value bool) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AddOnBuilder) ID(value string) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AddOnBuilder) HREF(value string) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnBuilder) Empty() bool {
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
func (b *AddOnBuilder) CommonAnnotations(value map[string]string) *AddOnBuilder {
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
func (b *AddOnBuilder) CommonLabels(value map[string]string) *AddOnBuilder {
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
// Representation of an add-on config.
// The attributes under it are to be used by the addon once its installed in the cluster.
func (b *AddOnBuilder) Config(value *AddOnConfigBuilder) *AddOnBuilder {
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
func (b *AddOnBuilder) CredentialsRequests(values ...*CredentialRequestBuilder) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.credentialsRequests = make([]*CredentialRequestBuilder, len(values))
	copy(b.credentialsRequests, values)
	b.fieldSet_[6] = true
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *AddOnBuilder) Description(value string) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.description = value
	b.fieldSet_[7] = true
	return b
}

// DocsLink sets the value of the 'docs_link' attribute to the given value.
func (b *AddOnBuilder) DocsLink(value string) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.docsLink = value
	b.fieldSet_[8] = true
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddOnBuilder) Enabled(value bool) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.enabled = value
	b.fieldSet_[9] = true
	return b
}

// HasExternalResources sets the value of the 'has_external_resources' attribute to the given value.
func (b *AddOnBuilder) HasExternalResources(value bool) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.hasExternalResources = value
	b.fieldSet_[10] = true
	return b
}

// Hidden sets the value of the 'hidden' attribute to the given value.
func (b *AddOnBuilder) Hidden(value bool) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.hidden = value
	b.fieldSet_[11] = true
	return b
}

// Icon sets the value of the 'icon' attribute to the given value.
func (b *AddOnBuilder) Icon(value string) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.icon = value
	b.fieldSet_[12] = true
	return b
}

// InstallMode sets the value of the 'install_mode' attribute to the given value.
//
// Representation of an add-on InstallMode field.
func (b *AddOnBuilder) InstallMode(value AddOnInstallMode) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.installMode = value
	b.fieldSet_[13] = true
	return b
}

// Label sets the value of the 'label' attribute to the given value.
func (b *AddOnBuilder) Label(value string) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.label = value
	b.fieldSet_[14] = true
	return b
}

// ManagedService sets the value of the 'managed_service' attribute to the given value.
func (b *AddOnBuilder) ManagedService(value bool) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.managedService = value
	b.fieldSet_[15] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddOnBuilder) Name(value string) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.name = value
	b.fieldSet_[16] = true
	return b
}

// Namespaces sets the value of the 'namespaces' attribute to the given values.
func (b *AddOnBuilder) Namespaces(values ...*AddOnNamespaceBuilder) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.namespaces = make([]*AddOnNamespaceBuilder, len(values))
	copy(b.namespaces, values)
	b.fieldSet_[17] = true
	return b
}

// OperatorName sets the value of the 'operator_name' attribute to the given value.
func (b *AddOnBuilder) OperatorName(value string) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.operatorName = value
	b.fieldSet_[18] = true
	return b
}

// Parameters sets the value of the 'parameters' attribute to the given values.
func (b *AddOnBuilder) Parameters(value *AddOnParameterListBuilder) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.parameters = value
	b.fieldSet_[19] = true
	return b
}

// Requirements sets the value of the 'requirements' attribute to the given values.
func (b *AddOnBuilder) Requirements(values ...*AddOnRequirementBuilder) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.requirements = make([]*AddOnRequirementBuilder, len(values))
	copy(b.requirements, values)
	b.fieldSet_[20] = true
	return b
}

// ResourceCost sets the value of the 'resource_cost' attribute to the given value.
func (b *AddOnBuilder) ResourceCost(value float64) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.resourceCost = value
	b.fieldSet_[21] = true
	return b
}

// ResourceName sets the value of the 'resource_name' attribute to the given value.
func (b *AddOnBuilder) ResourceName(value string) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.resourceName = value
	b.fieldSet_[22] = true
	return b
}

// SubOperators sets the value of the 'sub_operators' attribute to the given values.
func (b *AddOnBuilder) SubOperators(values ...*AddOnSubOperatorBuilder) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.subOperators = make([]*AddOnSubOperatorBuilder, len(values))
	copy(b.subOperators, values)
	b.fieldSet_[23] = true
	return b
}

// TargetNamespace sets the value of the 'target_namespace' attribute to the given value.
func (b *AddOnBuilder) TargetNamespace(value string) *AddOnBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 26)
	}
	b.targetNamespace = value
	b.fieldSet_[24] = true
	return b
}

// Version sets the value of the 'version' attribute to the given value.
//
// Representation of an add-on version.
func (b *AddOnBuilder) Version(value *AddOnVersionBuilder) *AddOnBuilder {
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
func (b *AddOnBuilder) Copy(object *AddOn) *AddOnBuilder {
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
		b.config = NewAddOnConfig().Copy(object.config)
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
		b.namespaces = make([]*AddOnNamespaceBuilder, len(object.namespaces))
		for i, v := range object.namespaces {
			b.namespaces[i] = NewAddOnNamespace().Copy(v)
		}
	} else {
		b.namespaces = nil
	}
	b.operatorName = object.operatorName
	if object.parameters != nil {
		b.parameters = NewAddOnParameterList().Copy(object.parameters)
	} else {
		b.parameters = nil
	}
	if object.requirements != nil {
		b.requirements = make([]*AddOnRequirementBuilder, len(object.requirements))
		for i, v := range object.requirements {
			b.requirements[i] = NewAddOnRequirement().Copy(v)
		}
	} else {
		b.requirements = nil
	}
	b.resourceCost = object.resourceCost
	b.resourceName = object.resourceName
	if object.subOperators != nil {
		b.subOperators = make([]*AddOnSubOperatorBuilder, len(object.subOperators))
		for i, v := range object.subOperators {
			b.subOperators[i] = NewAddOnSubOperator().Copy(v)
		}
	} else {
		b.subOperators = nil
	}
	b.targetNamespace = object.targetNamespace
	if object.version != nil {
		b.version = NewAddOnVersion().Copy(object.version)
	} else {
		b.version = nil
	}
	return b
}

// Build creates a 'add_on' object using the configuration stored in the builder.
func (b *AddOnBuilder) Build() (object *AddOn, err error) {
	object = new(AddOn)
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
		object.namespaces = make([]*AddOnNamespace, len(b.namespaces))
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
		object.requirements = make([]*AddOnRequirement, len(b.requirements))
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
		object.subOperators = make([]*AddOnSubOperator, len(b.subOperators))
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
