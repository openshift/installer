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

// AddonBuilder contains the data and logic needed to build 'addon' objects.
//
// Representation of an addon that can be installed in a cluster.
type AddonBuilder struct {
	bitmap_              uint32
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
	parameters           *AddonParametersBuilder
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
	return &AddonBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AddonBuilder) Link(value bool) *AddonBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AddonBuilder) ID(value string) *AddonBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AddonBuilder) HREF(value string) *AddonBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CommonAnnotations sets the value of the 'common_annotations' attribute to the given value.
func (b *AddonBuilder) CommonAnnotations(value map[string]string) *AddonBuilder {
	b.commonAnnotations = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// CommonLabels sets the value of the 'common_labels' attribute to the given value.
func (b *AddonBuilder) CommonLabels(value map[string]string) *AddonBuilder {
	b.commonLabels = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// Config sets the value of the 'config' attribute to the given value.
//
// Representation of an addon config.
// The attributes under it are to be used by the addon once its installed in the cluster.
func (b *AddonBuilder) Config(value *AddonConfigBuilder) *AddonBuilder {
	b.config = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// CredentialsRequests sets the value of the 'credentials_requests' attribute to the given values.
func (b *AddonBuilder) CredentialsRequests(values ...*CredentialRequestBuilder) *AddonBuilder {
	b.credentialsRequests = make([]*CredentialRequestBuilder, len(values))
	copy(b.credentialsRequests, values)
	b.bitmap_ |= 64
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *AddonBuilder) Description(value string) *AddonBuilder {
	b.description = value
	b.bitmap_ |= 128
	return b
}

// DocsLink sets the value of the 'docs_link' attribute to the given value.
func (b *AddonBuilder) DocsLink(value string) *AddonBuilder {
	b.docsLink = value
	b.bitmap_ |= 256
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddonBuilder) Enabled(value bool) *AddonBuilder {
	b.enabled = value
	b.bitmap_ |= 512
	return b
}

// HasExternalResources sets the value of the 'has_external_resources' attribute to the given value.
func (b *AddonBuilder) HasExternalResources(value bool) *AddonBuilder {
	b.hasExternalResources = value
	b.bitmap_ |= 1024
	return b
}

// Hidden sets the value of the 'hidden' attribute to the given value.
func (b *AddonBuilder) Hidden(value bool) *AddonBuilder {
	b.hidden = value
	b.bitmap_ |= 2048
	return b
}

// Icon sets the value of the 'icon' attribute to the given value.
func (b *AddonBuilder) Icon(value string) *AddonBuilder {
	b.icon = value
	b.bitmap_ |= 4096
	return b
}

// InstallMode sets the value of the 'install_mode' attribute to the given value.
//
// Representation of an addon InstallMode field.
func (b *AddonBuilder) InstallMode(value AddonInstallMode) *AddonBuilder {
	b.installMode = value
	b.bitmap_ |= 8192
	return b
}

// Label sets the value of the 'label' attribute to the given value.
func (b *AddonBuilder) Label(value string) *AddonBuilder {
	b.label = value
	b.bitmap_ |= 16384
	return b
}

// ManagedService sets the value of the 'managed_service' attribute to the given value.
func (b *AddonBuilder) ManagedService(value bool) *AddonBuilder {
	b.managedService = value
	b.bitmap_ |= 32768
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddonBuilder) Name(value string) *AddonBuilder {
	b.name = value
	b.bitmap_ |= 65536
	return b
}

// Namespaces sets the value of the 'namespaces' attribute to the given values.
func (b *AddonBuilder) Namespaces(values ...*AddonNamespaceBuilder) *AddonBuilder {
	b.namespaces = make([]*AddonNamespaceBuilder, len(values))
	copy(b.namespaces, values)
	b.bitmap_ |= 131072
	return b
}

// OperatorName sets the value of the 'operator_name' attribute to the given value.
func (b *AddonBuilder) OperatorName(value string) *AddonBuilder {
	b.operatorName = value
	b.bitmap_ |= 262144
	return b
}

// Parameters sets the value of the 'parameters' attribute to the given value.
//
// Representation of AddonParameters
func (b *AddonBuilder) Parameters(value *AddonParametersBuilder) *AddonBuilder {
	b.parameters = value
	if value != nil {
		b.bitmap_ |= 524288
	} else {
		b.bitmap_ &^= 524288
	}
	return b
}

// Requirements sets the value of the 'requirements' attribute to the given values.
func (b *AddonBuilder) Requirements(values ...*AddonRequirementBuilder) *AddonBuilder {
	b.requirements = make([]*AddonRequirementBuilder, len(values))
	copy(b.requirements, values)
	b.bitmap_ |= 1048576
	return b
}

// ResourceCost sets the value of the 'resource_cost' attribute to the given value.
func (b *AddonBuilder) ResourceCost(value float64) *AddonBuilder {
	b.resourceCost = value
	b.bitmap_ |= 2097152
	return b
}

// ResourceName sets the value of the 'resource_name' attribute to the given value.
func (b *AddonBuilder) ResourceName(value string) *AddonBuilder {
	b.resourceName = value
	b.bitmap_ |= 4194304
	return b
}

// SubOperators sets the value of the 'sub_operators' attribute to the given values.
func (b *AddonBuilder) SubOperators(values ...*AddonSubOperatorBuilder) *AddonBuilder {
	b.subOperators = make([]*AddonSubOperatorBuilder, len(values))
	copy(b.subOperators, values)
	b.bitmap_ |= 8388608
	return b
}

// TargetNamespace sets the value of the 'target_namespace' attribute to the given value.
func (b *AddonBuilder) TargetNamespace(value string) *AddonBuilder {
	b.targetNamespace = value
	b.bitmap_ |= 16777216
	return b
}

// Version sets the value of the 'version' attribute to the given value.
//
// Representation of an addon version.
func (b *AddonBuilder) Version(value *AddonVersionBuilder) *AddonBuilder {
	b.version = value
	if value != nil {
		b.bitmap_ |= 33554432
	} else {
		b.bitmap_ &^= 33554432
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonBuilder) Copy(object *Addon) *AddonBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
		b.parameters = NewAddonParameters().Copy(object.parameters)
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
	object.bitmap_ = b.bitmap_
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
