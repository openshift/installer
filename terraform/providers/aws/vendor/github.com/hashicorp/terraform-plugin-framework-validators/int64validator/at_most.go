package int64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.Int64 = atMostValidator{}

// atMostValidator validates that an integer Attribute's value is at most a certain value.
type atMostValidator struct {
	max int64
}

// Description describes the validation in plain text formatting.
func (validator atMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at most %d", validator.max)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator atMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// ValidateInt64 performs the validation.
func (v atMostValidator) ValidateInt64(ctx context.Context, request validator.Int64Request, response *validator.Int64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if request.ConfigValue.ValueInt64() > v.max {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", request.ConfigValue.ValueInt64()),
		))
	}
}

// AtMost returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a number, which can be represented by a 64-bit integer.
//   - Is less than or equal to the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtMost(max int64) validator.Int64 {
	return atMostValidator{
		max: max,
	}
}
