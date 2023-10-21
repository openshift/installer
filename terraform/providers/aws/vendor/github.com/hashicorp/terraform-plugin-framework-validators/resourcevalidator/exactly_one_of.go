package resourcevalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ExactlyOneOf checks that a set of path.Expression does not have more than
// one known value.
func ExactlyOneOf(expressions ...path.Expression) resource.ConfigValidator {
	return &configvalidator.ExactlyOneOfValidator{
		PathExpressions: expressions,
	}
}
