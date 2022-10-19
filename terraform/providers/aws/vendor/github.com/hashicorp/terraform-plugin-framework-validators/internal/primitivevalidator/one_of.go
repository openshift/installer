package primitivevalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// OneOf checks that the value held in the attribute
// is one of the given `acceptableValues`.
//
// This validator can be used only against primitives like
// `types.String`, `types.Number`, `types.Int64`,
// `types.Float64` and `types.Bool`.
func OneOf(acceptableValues ...attr.Value) tfsdk.AttributeValidator {
	return &acceptablePrimitiveValuesAttributeValidator{
		acceptableValues: acceptableValues,
		shouldMatch:      true,
	}
}
