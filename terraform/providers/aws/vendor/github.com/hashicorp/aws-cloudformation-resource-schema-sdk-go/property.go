package cfschema

import (
	"encoding/json"
)

const (
	PropertyFormatDate                = "date"
	PropertyFormatDateTime            = "date-time"
	PropertyFormatEmail               = "email"
	PropertyFormatHostname            = "hostname"
	PropertyFormatIdnEmail            = "idn-email"
	PropertyFormatIdnHostname         = "idn-hostname"
	PropertyFormatIpv4                = "ipv4"
	PropertyFormatIpv6                = "ipv6"
	PropertyFormatIri                 = "iri"
	PropertyFormatIriReference        = "iri-reference"
	PropertyFormatJsonPointer         = "json-pointer"
	PropertyFormatRegex               = "regex"
	PropertyFormatRelativeJsonPointer = "relative-json-pointer"
	PropertyFormatTime                = "time"
	PropertyFormatUri                 = "uri"
	PropertyFormatUriReference        = "uri-reference"
	PropertyFormatUriTemplate         = "uri-template"
)

const (
	PropertyTypeArray   = "array"
	PropertyTypeBoolean = "boolean"
	PropertyTypeInteger = "integer"
	PropertyTypeNull    = "null"
	PropertyTypeNumber  = "number"
	PropertyTypeObject  = "object"
	PropertyTypeString  = "string"
)

// Property represents the CloudFormation Resource Schema customization for Definitions and Properties.
type Property struct {
	AdditionalProperties *bool                `json:"additionalProperties,omitempty"`
	AllOf                []*PropertySubschema `json:"allOf,omitempty"`
	AnyOf                []*PropertySubschema `json:"anyOf,omitempty"`
	Comment              *string              `json:"$comment,omitempty"`
	Default              interface{}          `json:"default,omitempty"`
	Description          *string              `json:"description,omitempty"`
	Enum                 []interface{}        `json:"enum,omitempty"`
	Examples             []interface{}        `json:"examples,omitempty"`
	Format               *string              `json:"format,omitempty"`
	InsertionOrder       *bool                `json:"insertionOrder,omitempty"`
	Items                *Property            `json:"items,omitempty"`
	Maximum              *int                 `json:"maximum,omitempty"`
	MaxItems             *int                 `json:"maxItems,omitempty"`
	MaxLength            *int                 `json:"maxLength,omitempty"`
	Minimum              *int                 `json:"minimum,omitempty"`
	MinItems             *int                 `json:"minItems,omitempty"`
	MinLength            *int                 `json:"minLength,omitempty"`
	OneOf                []*PropertySubschema `json:"oneOf,omitempty"`
	Pattern              *string              `json:"pattern,omitempty"`
	PatternProperties    map[string]*Property `json:"patternProperties,omitempty"`
	Properties           map[string]*Property `json:"properties,omitempty"`
	Ref                  *Reference           `json:"$ref,omitempty"`
	Required             []string             `json:"required,omitempty"`
	Type                 *Type                `json:"type,omitempty"`
	UniqueItems          *bool                `json:"uniqueItems,omitempty"`
}

// String returns a string representation of Property.
func (p *Property) String() string {
	if p == nil {
		return ""
	}

	b, _ := json.MarshalIndent(p, "", "  ")

	return string(b)
}

func (p *Property) IsRequired(name string) bool {
	if p == nil {
		return false
	}

	for _, req := range p.Required {
		if req == name {
			return true
		}
	}

	return false
}
