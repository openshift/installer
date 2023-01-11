package stringvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = lengthBetweenValidator{}

// stringLenBetweenValidator validates that a string Attribute's length is in a range.
type lengthBetweenValidator struct {
	minLength, maxLength int
}

// Description describes the validation in plain text formatting.
func (validator lengthBetweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("string length must be between %d and %d", validator.minLength, validator.maxLength)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator lengthBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator lengthBetweenValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	s, ok := validateString(ctx, request, response)

	if !ok {
		return
	}

	if l := len(s); l < validator.minLength || l > validator.maxLength {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueLengthDiagnostic(
			request.AttributePath,
			validator.Description(ctx),
			fmt.Sprintf("%d", l),
		))

		return
	}
}

// LengthBetween returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a string.
//   - Is of length greater than the given minimum and less than the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func LengthBetween(minLength, maxLength int) tfsdk.AttributeValidator {
	if minLength < 0 || maxLength < 0 || minLength > maxLength {
		return nil
	}

	return lengthBetweenValidator{
		minLength: minLength,
		maxLength: maxLength,
	}
}
