package cfschema

// MetaJsonSchema represents the meta-schema for resource schemas
type MetaJsonSchema struct {
	jsonSchema
}

// ValidateResourceDocument validates the provided document against the meta-schema.
func (s *MetaJsonSchema) ValidateResourceDocument(document string) error {
	if s == nil {
		return nil
	}

	return s.validateDocument(document)
}

// ValidateResourceJsonSchema validates the provided ResourceJsonSchema against the meta-schema.
func (s *MetaJsonSchema) ValidateResourceJsonSchema(resourceJsonSchema *ResourceJsonSchema) error {
	if s == nil || resourceJsonSchema == nil {
		return nil
	}

	return s.validateJsonSchema(resourceJsonSchema.jsonSchema)
}

// ValidateResourcePath validates the provided document at the file path against the meta-schema.
func (s *MetaJsonSchema) ValidateResourcePath(path string) error {
	if s == nil {
		return nil
	}

	return s.validatePath(path)
}

// NewMetaJsonSchemaDocument returns a MetaJsonSchema or any errors from the provided document.
func NewMetaJsonSchemaDocument(document string) (*MetaJsonSchema, error) {
	js, err := newJsonSchemaDocument(document)

	if err != nil {
		return nil, err
	}

	return &MetaJsonSchema{
		jsonSchema: *js,
	}, nil
}

// NewMetaJsonSchemaPath returns a MetaJsonSchema or any errors from the provided document at the file path.
func NewMetaJsonSchemaPath(path string) (*MetaJsonSchema, error) {
	js, err := newJsonSchemaPath(path)

	if err != nil {
		return nil, err
	}

	return &MetaJsonSchema{
		jsonSchema: *js,
	}, nil
}
