// Copyright 2023 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dcl

import (
	"fmt"
	"strings"
)

// FieldType is an enum of all the types a field can be.
type FieldType int

const (
	// UnknownType refers to a Field that does not have a proper type defined.
	UnknownType FieldType = iota
	// MapType refers to a Field that is a Map (typically from string to string).
	MapType
	// EnumType refers to a Field that is an Enum.
	EnumType
	// ArrayType refers to a Field that is an Array of any kind.
	ArrayType
	// ObjectType refers to a Field that is a dictionary with set subfields.
	ObjectType
	// ReferenceType refers to a Field that is referencing another GCP resource.
	ReferenceType
	// DoubleType refers to a Field that is a Double.
	DoubleType
	// StringType refers to a Field that is a String.
	StringType
	// TimeType refers to a Field that is a Timestamp.
	TimeType
	// IntegerType refers to a Field that is an Integer.
	IntegerType
	// BooleanType refers to a Field that is a Boolean.
	BooleanType
	// StatusType refers to a Field that is a Status.
	StatusType
	// ReusedType refers to a Field that does not require additional generation because it
	// is the same type as another field already being generated.
	ReusedType
	// UntypedType refers to a type that has no type (in Go-speak, that's an interface{}).
	// This can only be used for untyped maps (in proto-speak, google.protobuf.Struct) and cannot be used anywhere else.
	// This will not work properly if used outside of a map.
	UntypedType
)

// Schema is the Entire OpenAPI schema.
type Schema struct {
	Info       *Info       `yaml:"info"`
	Paths      *Paths      `yaml:"paths"`
	Components *Components `yaml:"components"`
}

// ResolveDefinition returns the schema component being referenced.
func (s *Schema) ResolveDefinition(ref string) (*Component, error) {
	if strings.HasPrefix(ref, "#/components/schemas/") {
		if props, ok := s.Components.Schemas[strings.TrimPrefix(ref, "#/components/schemas/")]; ok {
			return props, nil
		}
	}
	return nil, fmt.Errorf("could not resolve reference %q\v", ref)
}

// Link is a URL plus text that should be displayed to an end user, usually in docs.
type Link struct {
	Text string `yaml:"text"`
	URL  string `yaml:"url"`
}

// Info is the Info block for the OpenAPI schema.
type Info struct {
	Title       string  `yaml:"title"`
	Description string  `yaml:"description"`
	StructName  string  `yaml:"x-dcl-struct-name,omitempty"`
	HasIAM      bool    `yaml:"x-dcl-has-iam"`
	Mutex       string  `yaml:"x-dcl-mutex,omitempty"`
	Note        string  `yaml:"x-dcl-note,omitempty"`
	Warning     string  `yaml:"x-dcl-warning,omitempty"`
	Reference   *Link   `yaml:"x-dcl-ref,omitempty"`
	Guides      []*Link `yaml:"x-dcl-guides,omitempty"`
}

// ResourceTitle returns the title of this resource.
func (i *Info) ResourceTitle() string {
	return strings.Split(i.Title, "/")[1]
}

// Paths is the Paths block for the OpenAPI schema.
type Paths struct {
	Get       *Path `yaml:"get"`
	Apply     *Path `yaml:"apply"`
	Delete    *Path `yaml:"delete,omitempty"`
	DeleteAll *Path `yaml:"deleteAll,omitempty"`
	List      *Path `yaml:"list,omitempty"`
}

// Path is the Path used for a method.
type Path struct {
	Description string           `yaml:"description"`
	Parameters  []PathParameters `yaml:"parameters"`
}

// PathParameters is the Parameters for a given Path.
type PathParameters struct {
	Name        string                `yaml:"name"`
	Required    bool                  `yaml:"required"`
	Description string                `yaml:"description,omitempty"`
	Schema      *PathParametersSchema `yaml:"schema,omitempty"`
}

// PathParametersSchema is used to store the type. It is typically set to "string"
type PathParametersSchema struct {
	Type string `yaml:"type"`
}

// Components maps a Component name to the Component.
type Components struct {
	Schemas map[string]*Component
}

// Component contains all the information for a component (resource or reused type)
type Component struct {
	Title           string   `yaml:"title,omitempty"`
	ID              string   `yaml:"x-dcl-id,omitempty"`
	Locations       []string `yaml:"x-dcl-locations,omitempty"`
	UsesStateHint   bool     `yaml:"x-dcl-uses-state-hint,omitempty"`
	ParentContainer string   `yaml:"x-dcl-parent-container,omitempty"`
	LabelsField     string   `yaml:"x-dcl-labels,omitempty"`
	HasCreate       bool     `yaml:"x-dcl-has-create"`
	HasIAM          bool     `yaml:"x-dcl-has-iam"`
	ReadTimeout     int      `yaml:"x-dcl-read-timeout"`
	ApplyTimeout    int      `yaml:"x-dcl-apply-timeout"`
	DeleteTimeout   int      `yaml:"x-dcl-delete-timeout"`

	// TODO: It appears that reused types are not fully conforming to the same spec as the rest of the components.
	// Reused Types seem to follow the property spec, but not the component spec.
	// This means that we need to have component "inline" all of the schema property fields to avoid having to override YAML parsing logic.
	SchemaProperty Property `yaml:",inline"`
}

// Property contains all information for a field (i.e. property)
type Property struct {
	Type                     string                       `yaml:"type,omitempty"`
	Format                   string                       `yaml:"format,omitempty"`
	AdditionalProperties     *Property                    `yaml:"additionalProperties,omitempty"`
	Ref                      string                       `yaml:"$ref,omitempty"`
	GoName                   string                       `yaml:"x-dcl-go-name,omitempty"`
	GoType                   string                       `yaml:"x-dcl-go-type,omitempty"`
	ReadOnly                 bool                         `yaml:"readOnly,omitempty"`
	Description              string                       `yaml:"description,omitempty"`
	Immutable                bool                         `yaml:"x-kubernetes-immutable,omitempty"`
	Conflicts                []string                     `yaml:"x-dcl-conflicts,omitempty"`
	Default                  interface{}                  `yaml:"default,omitempty"`
	ServerDefault            bool                         `yaml:"x-dcl-server-default,omitempty"`
	ServerGeneratedParameter bool                         `yaml:"x-dcl-server-generated-parameter,omitempty"`
	Sensitive                bool                         `yaml:"x-dcl-sensitive,omitempty"`
	ForwardSlashAllowed      bool                         `yaml:"x-dcl-forward-slash-allowed,omitempty"`
	SendEmpty                bool                         `yaml:"x-dcl-send-empty,omitempty"`
	ResourceReferences       []*PropertyResourceReference `yaml:"x-dcl-references,omitempty"`
	Enum                     []string                     `yaml:"enum,omitempty"`
	ListType                 string                       `yaml:"x-dcl-list-type,omitempty"`
	Items                    *Property                    `yaml:"items,omitempty"`
	Unreadable               bool                         `yaml:"x-dcl-mutable-unreadable,omitempty"`
	ExtractIfEmpty           bool                         `yaml:"x-dcl-extract-if-empty,omitempty"`
	Required                 []string                     `yaml:"required,omitempty"`
	Properties               map[string]*Property         `yaml:"properties,omitempty"`
	Deprecated               bool                         `yaml:"x-dcl-deprecated,omitempty"`
	OptionalType             bool                         `yaml:"x-dcl-optional-type,omitempty"`
}

// IsOptional returns if the type is an optional type.
func (p *Property) IsOptional() bool {
	return p.OptionalType
}

// TypeEnum returns an enum referring to the type.
func (p *Property) TypeEnum() FieldType {
	switch p.Type {
	case "string":
		if p.GoType != "" && p.GoType != "string" {
			return EnumType
		} else if len(p.ResourceReferences) > 0 {
			return ReferenceType
		}
		return StringType
	case "OptionalString":
		return StringType
	case "number", "OptionalFloat":
		return DoubleType
	case "integer", "OptionalInt":
		return IntegerType
	case "boolean", "OptionalBool":
		return BooleanType
	case "object":
		if p.AdditionalProperties != nil && p.AdditionalProperties.GoType != "" && len(p.AdditionalProperties.Properties) != 0 {
			return MapType
		}
		return ObjectType
	case "array":
		return ArrayType
	}
	return UnknownType
}

// PropertyResourceReference contains all resource reference information.
type PropertyResourceReference struct {
	Resource string `yaml:"resource"`
	Field    string `yaml:"field"`
	Format   string `yaml:"format,omitempty"`
	Parent   bool   `yaml:"parent,omitempty"`
}
