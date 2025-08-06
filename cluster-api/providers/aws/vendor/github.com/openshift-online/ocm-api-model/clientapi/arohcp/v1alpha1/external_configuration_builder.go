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
	v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
)

// Representation of cluster external configuration.
type ExternalConfigurationBuilder struct {
	fieldSet_ []bool
	labels    *v1.LabelListBuilder
	manifests *v1.ManifestListBuilder
	syncsets  *v1.SyncsetListBuilder
}

// NewExternalConfiguration creates a new builder of 'external_configuration' objects.
func NewExternalConfiguration() *ExternalConfigurationBuilder {
	return &ExternalConfigurationBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ExternalConfigurationBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// Labels sets the value of the 'labels' attribute to the given values.
func (b *ExternalConfigurationBuilder) Labels(value *v1.LabelListBuilder) *ExternalConfigurationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.labels = value
	b.fieldSet_[0] = true
	return b
}

// Manifests sets the value of the 'manifests' attribute to the given values.
func (b *ExternalConfigurationBuilder) Manifests(value *v1.ManifestListBuilder) *ExternalConfigurationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.manifests = value
	b.fieldSet_[1] = true
	return b
}

// Syncsets sets the value of the 'syncsets' attribute to the given values.
func (b *ExternalConfigurationBuilder) Syncsets(value *v1.SyncsetListBuilder) *ExternalConfigurationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.syncsets = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ExternalConfigurationBuilder) Copy(object *ExternalConfiguration) *ExternalConfigurationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.labels != nil {
		b.labels = v1.NewLabelList().Copy(object.labels)
	} else {
		b.labels = nil
	}
	if object.manifests != nil {
		b.manifests = v1.NewManifestList().Copy(object.manifests)
	} else {
		b.manifests = nil
	}
	if object.syncsets != nil {
		b.syncsets = v1.NewSyncsetList().Copy(object.syncsets)
	} else {
		b.syncsets = nil
	}
	return b
}

// Build creates a 'external_configuration' object using the configuration stored in the builder.
func (b *ExternalConfigurationBuilder) Build() (object *ExternalConfiguration, err error) {
	object = new(ExternalConfiguration)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.labels != nil {
		object.labels, err = b.labels.Build()
		if err != nil {
			return
		}
	}
	if b.manifests != nil {
		object.manifests, err = b.manifests.Build()
		if err != nil {
			return
		}
	}
	if b.syncsets != nil {
		object.syncsets, err = b.syncsets.Build()
		if err != nil {
			return
		}
	}
	return
}
