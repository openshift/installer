package tfsdk

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
)

// ListNestedAttributes nests `attributes` under another attribute, allowing
// multiple instances of that group of attributes to appear in the
// configuration.
func ListNestedAttributes(attributes map[string]Attribute) fwschema.NestedAttributes {
	return fwschema.ListNestedAttributes{
		UnderlyingAttributes: schemaAttributes(attributes),
	}
}

// MapNestedAttributes nests `attributes` under another attribute, allowing
// multiple instances of that group of attributes to appear in the
// configuration. Each group will need to be associated with a unique string by
// the user.
func MapNestedAttributes(attributes map[string]Attribute) fwschema.NestedAttributes {
	return fwschema.MapNestedAttributes{
		UnderlyingAttributes: schemaAttributes(attributes),
	}
}

// SetNestedAttributes nests `attributes` under another attribute, allowing
// multiple instances of that group of attributes to appear in the
// configuration, while requiring each group of values be unique.
func SetNestedAttributes(attributes map[string]Attribute) fwschema.NestedAttributes {
	return fwschema.SetNestedAttributes{
		UnderlyingAttributes: schemaAttributes(attributes),
	}
}

// SingleNestedAttributes nests `attributes` under another attribute, only
// allowing one instance of that group of attributes to appear in the
// configuration.
func SingleNestedAttributes(attributes map[string]Attribute) fwschema.NestedAttributes {
	return fwschema.SingleNestedAttributes{
		UnderlyingAttributes: schemaAttributes(attributes),
	}
}
