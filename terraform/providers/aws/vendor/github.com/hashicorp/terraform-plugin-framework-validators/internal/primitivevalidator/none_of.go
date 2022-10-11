package primitivevalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// NoneOf checks that the value held in the attribute
// is none of the given `unacceptableValues`.
//
// This validator can be used only against primitives like
// `types.String`, `types.Number`, `types.Int64`,
// `types.Float64` and `types.Bool`.
func NoneOf(unacceptableValues ...attr.Value) tfsdk.AttributeValidator {
	return &acceptablePrimitiveValuesAttributeValidator{
		acceptableValues: unacceptableValues,
		shouldMatch:      false,
	}
}
