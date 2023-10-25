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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// TemplateParameterBuilder contains the data and logic needed to build 'template_parameter' objects.
//
// A template parameter is used in an email to replace placeholder content with
// values specific to the email recipient.
type TemplateParameterBuilder struct {
	bitmap_ uint32
	content string
	name    string
}

// NewTemplateParameter creates a new builder of 'template_parameter' objects.
func NewTemplateParameter() *TemplateParameterBuilder {
	return &TemplateParameterBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *TemplateParameterBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Content sets the value of the 'content' attribute to the given value.
func (b *TemplateParameterBuilder) Content(value string) *TemplateParameterBuilder {
	b.content = value
	b.bitmap_ |= 1
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *TemplateParameterBuilder) Name(value string) *TemplateParameterBuilder {
	b.name = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *TemplateParameterBuilder) Copy(object *TemplateParameter) *TemplateParameterBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.content = object.content
	b.name = object.name
	return b
}

// Build creates a 'template_parameter' object using the configuration stored in the builder.
func (b *TemplateParameterBuilder) Build() (object *TemplateParameter, err error) {
	object = new(TemplateParameter)
	object.bitmap_ = b.bitmap_
	object.content = b.content
	object.name = b.name
	return
}
