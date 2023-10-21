package setvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Set = sizeAtMostValidator{}

// sizeAtMostValidator validates that set contains at most max elements.
type sizeAtMostValidator struct {
	max int
}

// Description describes the validation in plain text formatting.
func (v sizeAtMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("set must contain at most %d elements", v.max)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v sizeAtMostValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v sizeAtMostValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elems := req.ConfigValue.Elements()

	if len(elems) > v.max {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		))
	}
}

// SizeAtMost returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a Set.
//   - Contains at most max elements.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func SizeAtMost(max int) validator.Set {
	return sizeAtMostValidator{
		max: max,
	}
}
