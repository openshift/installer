package explain

import (
	"github.com/pkg/errors"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func lookup(schema *apiextv1.JSONSchemaProps, path []string) (*apiextv1.JSONSchemaProps, error) {
	if len(path) == 0 {
		return schema, nil
	}

	properties := map[string]apiextv1.JSONSchemaProps{}
	if schema.Items != nil && schema.Items.Schema != nil {
		properties = schema.Items.Schema.Properties
	}
	if len(schema.Properties) > 0 {
		properties = schema.Properties
	}

	property, ok := properties[path[0]]
	if !ok {
		return nil, errors.Errorf("invalid field %s, no such property found", path[0])
	}
	return lookup(&property, path[1:])
}
