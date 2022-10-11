package fwxschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// AttributeWithValidators is an optional interface on Attribute which enables
// validation support.
type AttributeWithValidators interface {
	// Implementations should include the fwschema.Attribute interface methods
	// for proper attribute handling.
	fwschema.Attribute

	// GetValidators should return a list of attribute-based validators. This
	// is named differently than PlanModifiers to prevent a conflict with the
	// tfsdk.Attribute field name.
	GetValidators() []tfsdk.AttributeValidator
}
