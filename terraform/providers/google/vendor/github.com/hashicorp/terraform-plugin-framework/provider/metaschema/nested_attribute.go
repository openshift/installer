package metaschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
)

// Nested attributes are only compatible with protocol version 6.
type NestedAttribute interface {
	Attribute
	fwschema.NestedAttribute
}
