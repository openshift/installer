package stringvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/primitivevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OneOf checks that the string held in the attribute
// is one of the given `acceptableStrings`.
func OneOf(acceptableStrings ...string) tfsdk.AttributeValidator {
	acceptableStringValues := make([]attr.Value, 0, len(acceptableStrings))
	for _, s := range acceptableStrings {
		acceptableStringValues = append(acceptableStringValues, types.String{Value: s})
	}

	return primitivevalidator.OneOf(acceptableStringValues...)
}

// OneOfCaseInsensitive checks that the string held in the attribute
// is one of the given `acceptableStrings`, irrespective of case sensitivity.
func OneOfCaseInsensitive(acceptableStrings ...string) tfsdk.AttributeValidator {
	return &acceptableStringsCaseInsensitiveAttributeValidator{
		acceptableStrings,
		true,
	}
}
