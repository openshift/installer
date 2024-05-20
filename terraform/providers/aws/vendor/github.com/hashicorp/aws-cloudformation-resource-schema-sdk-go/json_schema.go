package cfschema

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/xeipuuv/gojsonschema"
)

// jsonSchema is an internal implementation for shared JSON Schema functionality.
type jsonSchema struct {
	path   string
	loader gojsonschema.JSONLoader
	schema *gojsonschema.Schema
	source []byte
}

// validateDocument validates the provided document against the meta-schema.
func (s *jsonSchema) validateDocument(document string) error {
	documentLoader := gojsonschema.NewStringLoader(document)

	return s.validate(documentLoader)
}

// validateJsonSchema validates the provided jsonSchema against the meta-schema.
func (s *jsonSchema) validateJsonSchema(s2 jsonSchema) error {
	return s.validate(s2.loader)
}

// validatePath validates the document at the provided file path against the meta-schema.
func (s *jsonSchema) validatePath(path string) error {
	documentLoader := gojsonschema.NewReferenceLoader("file://" + path)

	return s.validate(documentLoader)
}

// validate performs common validation logic.
func (s *jsonSchema) validate(loader gojsonschema.JSONLoader) error {
	result, err := s.schema.Validate(loader)

	if err != nil {
		return fmt.Errorf("Unable to Validate JSON Schema: %w", err)
	}

	if !result.Valid() {
		var errs *multierror.Error

		for _, resultError := range result.Errors() {
			errs = multierror.Append(errs, fmt.Errorf("%s", resultError.String()))
		}

		return fmt.Errorf("Validation Errors: %w", errs)
	}

	return nil
}

// newJsonSchemaDocument returns a jsonSchema or any errors from a provided document.
func newJsonSchemaDocument(document string) (*jsonSchema, error) {
	schemaLoader := gojsonschema.NewStringLoader(document)

	schema, err := gojsonschema.NewSchema(schemaLoader)

	if err != nil {
		return nil, fmt.Errorf("loading JSON Schema (%s): %w", document, err)
	}

	return &jsonSchema{
		loader: schemaLoader,
		schema: schema,
		source: []byte(document),
	}, nil
}

// newJsonSchemaPath returns a jsonSchema or any errors from a provided document at the file path.
func newJsonSchemaPath(path string) (*jsonSchema, error) {
	// To prevent reading the file twice to populate source bytes,
	// manually read file path and use string handler.
	f, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("reading file (%s): %w", path, err)
	}

	cwd, err := os.Getwd()

	if err != nil {
		return nil, fmt.Errorf("getting current directory: %w", err)
	}

	defer func() {
		os.Chdir(cwd)
	}()

	// CD to the schema's directory so as to resolve any relative 'file://' URLs.
	os.Chdir(filepath.Dir(path))

	js, err := newJsonSchemaDocument(string(f))

	if err != nil {
		return nil, err
	}

	js.path = path

	return js, nil
}
