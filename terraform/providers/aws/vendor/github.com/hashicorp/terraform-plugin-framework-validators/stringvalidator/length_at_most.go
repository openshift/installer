package stringvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = lengthAtMostValidator{}

// lengthAtMostValidator validates that a string Attribute's length is at most a certain value.
type lengthAtMostValidator struct {
	maxLength int
}

// Description describes the validation in plain text formatting.
func (validator lengthAtMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("string length must be at most %d", validator.maxLength)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator lengthAtMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator lengthAtMostValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	s, ok := validateString(ctx, request, response)

	if !ok {
		return
	}

	if l := len(s); l > validator.maxLength {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueLengthDiagnostic(
			request.AttributePath,
			validator.Description(ctx),
			fmt.Sprintf("%d", l),
		))

		return
	}
}

// LengthAtMost returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a string.
//   - Is of length exclusively less than the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func LengthAtMost(maxLength int) tfsdk.AttributeValidator {
	if maxLength < 0 {
		return nil
	}

	return lengthAtMostValidator{
		maxLength: maxLength,
	}
}
