package stringvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateString ensures that the request contains a String value.
func validateString(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) (string, bool) {
	t := request.AttributeConfig.Type(ctx)
	if t != types.StringType {
		response.Diagnostics.Append(validatordiag.InvalidAttributeTypeDiagnostic(
			request.AttributePath,
			"expected value of type string",
			t.String(),
		))
		return "", false
	}

	s := request.AttributeConfig.(types.String)

	if s.Unknown || s.Null {
		return "", false
	}

	return s.Value, true
}
