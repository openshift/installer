package cfschema

import (
	"fmt"
)

type Resource struct {
	AdditionalIdentifiers           []PropertyJsonPointers `json:"additionalIdentifiers,omitempty"`
	AdditionalProperties            *bool                  `json:"additionalProperties,omitempty"`
	AllOf                           []*PropertySubschema   `json:"allOf,omitempty"`
	AnyOf                           []*PropertySubschema   `json:"anyOf,omitempty"`
	ConditionalCreateOnlyProperties PropertyJsonPointers   `json:"conditionalCreateOnlyProperties,omitempty"`
	CreateOnlyProperties            PropertyJsonPointers   `json:"createOnlyProperties,omitempty"`
	Definitions                     map[string]*Property   `json:"definitions,omitempty"`
	DeprecatedProperties            PropertyJsonPointers   `json:"deprecatedProperties,omitempty"`
	Description                     *string                `json:"description,omitempty"`
	Handlers                        map[string]*Handler    `json:"handlers,omitempty"`
	OneOf                           []*PropertySubschema   `json:"oneOf,omitempty"`
	PrimaryIdentifier               PropertyJsonPointers   `json:"primaryIdentifier,omitempty"`
	Properties                      map[string]*Property   `json:"properties,omitempty"`
	PropertyTransform               map[string]string      `json:"propertyTransform,omitempty"`
	ReadOnlyProperties              PropertyJsonPointers   `json:"readOnlyProperties,omitempty"`
	ReplacementStrategy             *string                `json:"replacementStrategy,omitempty"`
	Required                        []string               `json:"required,omitempty"`
	ResourceLink                    *ResourceLink          `json:"resourceLink,omitempty"`
	SourceURL                       *string                `json:"sourceUrl,omitempty"`
	Taggable                        *bool                  `json:"taggable,omitempty"`
	Tagging                         *Tagging               `json:"tagging,omitempty"`
	TypeName                        *string                `json:"typeName,omitempty"`
	WriteOnlyProperties             PropertyJsonPointers   `json:"writeOnlyProperties,omitempty"`
}

func (r *Resource) IsCreateOnlyPropertyPath(path string) bool {
	if r == nil {
		return false
	}

	for _, createOnlyProperty := range r.CreateOnlyProperties {
		if createOnlyProperty.EqualsStringPath(path) {
			return true
		}
	}

	return false
}

func (r *Resource) IsRequired(name string) bool {
	if r == nil {
		return false
	}

	for _, req := range r.Required {
		if req == name {
			return true
		}
	}

	return false
}

// ResolveReference resolves a Reference (JSON Pointer) into a Property.
func (r *Resource) ResolveReference(ref Reference) (*Property, error) {
	if r == nil {
		return nil, nil
	}

	typ, err := ref.Type()

	if err != nil {
		return nil, err
	}

	var properties map[string]*Property

	switch typ {
	case ReferenceTypeDefinitions:
		properties = r.Definitions
	case ReferenceTypeProperties:
		properties = r.Properties
	default:
		return nil, fmt.Errorf("unexpected Reference type: %s", typ)
	}

	field, err := ref.Field()

	if err != nil {
		return nil, err
	}

	property, ok := properties[field]
	if !ok || property == nil {
		return nil, fmt.Errorf("%s/%s not found", typ, field)
	}

	return property, nil
}
