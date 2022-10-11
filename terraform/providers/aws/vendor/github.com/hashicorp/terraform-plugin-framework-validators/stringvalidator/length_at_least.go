package stringvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = lengthAtLeastValidator{}

// stringLenAtLeastValidator validates that a string Attribute's length is at least a certain value.
type lengthAtLeastValidator struct {
	minLength int
}

// Description describes the validation in plain text formatting.
func (validator lengthAtLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("string length must be at least %d", validator.minLength)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator lengthAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator lengthAtLeastValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	s, ok := validateString(ctx, request, response)

	if !ok {
		return
	}

	if l := len(s); l < validator.minLength {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueLengthDiagnostic(
			request.AttributePath,
			validator.Description(ctx),
			fmt.Sprintf("%d", l),
		))

		return
	}
}

// LengthAtLeast returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a string.
//   - Is of length exclusively greater than the given minimum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func LengthAtLeast(minLength int) tfsdk.AttributeValidator {
	if minLength < 0 {
		return nil
	}

	return lengthAtLeastValidator{
		minLength: minLength,
	}
}
