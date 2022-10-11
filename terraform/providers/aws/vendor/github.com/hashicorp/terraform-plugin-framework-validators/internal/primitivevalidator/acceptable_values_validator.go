package primitivevalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// acceptablePrimitiveValuesAttributeValidator is the underlying struct implementing OneOf and NoneOf.
type acceptablePrimitiveValuesAttributeValidator struct {
	acceptableValues []attr.Value
	shouldMatch      bool
}

var _ tfsdk.AttributeValidator = (*acceptablePrimitiveValuesAttributeValidator)(nil)

func (av *acceptablePrimitiveValuesAttributeValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av *acceptablePrimitiveValuesAttributeValidator) MarkdownDescription(_ context.Context) string {
	if av.shouldMatch {
		return fmt.Sprintf("Value must be one of: %q", av.acceptableValues)
	} else {
		return fmt.Sprintf("Value must be none of: %q", av.acceptableValues)
	}

}

func (av *acceptablePrimitiveValuesAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	if req.AttributeConfig.IsNull() || req.AttributeConfig.IsUnknown() {
		return
	}

	var value attr.Value
	switch typedAttributeConfig := req.AttributeConfig.(type) {
	case types.String, types.Bool, types.Int64, types.Float64, types.Number:
		value = typedAttributeConfig
	default:
		res.Diagnostics.AddAttributeError(
			req.AttributePath,
			"This validator should be used against primitive types (String, Bool, Number, Int64, Float64).",
			"This is always indicative of a bug within the provider.",
		)
		return
	}

	if av.shouldMatch && !av.isAcceptableValue(value) || //< EITHER should match but it does not
		!av.shouldMatch && av.isAcceptableValue(value) { //< OR should not match but it does
		res.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			req.AttributePath,
			av.Description(ctx),
			value.String(),
		))
	}
}

func (av *acceptablePrimitiveValuesAttributeValidator) isAcceptableValue(v attr.Value) bool {
	for _, acceptableV := range av.acceptableValues {
		if v.Equal(acceptableV) {
			return true
		}
	}

	return false
}
