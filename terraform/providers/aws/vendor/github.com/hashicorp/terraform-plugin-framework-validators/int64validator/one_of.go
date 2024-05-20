package int64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.Int64 = oneOfValidator{}

// oneOfValidator validates that the value matches one of expected values.
type oneOfValidator struct {
	values []types.Int64
}

func (v oneOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v oneOfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be one of: %q", v.values)
}

func (v oneOfValidator) ValidateInt64(ctx context.Context, request validator.Int64Request, response *validator.Int64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	for _, otherValue := range v.values {
		if value.Equal(otherValue) {
			return
		}
	}

	response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
		request.Path,
		v.Description(ctx),
		value.String(),
	))
}

// OneOf checks that the Int64 held in the attribute
// is one of the given `values`.
func OneOf(values ...int64) validator.Int64 {
	frameworkValues := make([]types.Int64, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.Int64Value(value))
	}

	return oneOfValidator{
		values: frameworkValues,
	}
}
