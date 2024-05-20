package cfschema

import (
	"encoding/json"
	"fmt"
)

// ResourceJsonSchema represents the resource schema.
type ResourceJsonSchema struct {
	jsonSchema
}

// Resource parses the JSON Schema and returns Resource or an error.
func (s *ResourceJsonSchema) Resource() (*Resource, error) {
	if s == nil {
		return nil, nil
	}

	var result Resource

	err := json.Unmarshal(s.source, &result)

	if err != nil {
		return nil, fmt.Errorf("parsing JSON Schema into Resource: %w", err)
	}

	return &result, nil
}

// ValidateConfigurationDocument validates the provided document against the resource schema.
func (s *ResourceJsonSchema) ValidateConfigurationDocument(document string) error {
	if s == nil {
		return nil
	}

	return s.validateDocument(document)
}

// ValidateConfigurationPath validates the provided document at the file path against the resource schema.
func (s *ResourceJsonSchema) ValidateConfigurationPath(path string) error {
	if s == nil {
		return nil
	}

	return s.validatePath(path)
}

// NewResourceJsonSchemaDocument returns a ResourceJsonSchema or any errors from the provided document.
func NewResourceJsonSchemaDocument(document string) (*ResourceJsonSchema, error) {
	js, err := newJsonSchemaDocument(document)

	if err != nil {
		return nil, err
	}

	return &ResourceJsonSchema{
		jsonSchema: *js,
	}, nil
}

// NewResourceJsonSchemaPath returns a ResourceJsonSchema or any errors from the provided document at the file path.
func NewResourceJsonSchemaPath(path string) (*ResourceJsonSchema, error) {
	js, err := newJsonSchemaPath(path)

	if err != nil {
		return nil, err
	}

	return &ResourceJsonSchema{
		jsonSchema: *js,
	}, nil
}
