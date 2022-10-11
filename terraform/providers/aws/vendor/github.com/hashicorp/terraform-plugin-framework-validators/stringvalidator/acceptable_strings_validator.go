package stringvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// acceptableStringsCaseInsensitiveAttributeValidator is the underlying struct implementing OneOf and NoneOf.
type acceptableStringsCaseInsensitiveAttributeValidator struct {
	acceptableStrings []string
	shouldMatch       bool
}

var _ tfsdk.AttributeValidator = (*acceptableStringsCaseInsensitiveAttributeValidator)(nil)

func (av *acceptableStringsCaseInsensitiveAttributeValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av *acceptableStringsCaseInsensitiveAttributeValidator) MarkdownDescription(_ context.Context) string {
	if av.shouldMatch {
		return fmt.Sprintf("String must match one of: %q", av.acceptableStrings)
	} else {
		return fmt.Sprintf("String must match none of: %q", av.acceptableStrings)
	}
}

func (av *acceptableStringsCaseInsensitiveAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	value, ok := validateString(ctx, req, res)
	if !ok {
		return
	}

	if av.shouldMatch && !av.isAcceptableValue(value) || //< EITHER should match but it does not
		!av.shouldMatch && av.isAcceptableValue(value) { //< OR should not match but it does
		res.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			req.AttributePath,
			av.Description(ctx),
			value,
		))
	}
}

func (av *acceptableStringsCaseInsensitiveAttributeValidator) isAcceptableValue(v string) bool {
	for _, acceptableV := range av.acceptableStrings {
		if strings.EqualFold(v, acceptableV) {
			return true
		}
	}

	return false
}
