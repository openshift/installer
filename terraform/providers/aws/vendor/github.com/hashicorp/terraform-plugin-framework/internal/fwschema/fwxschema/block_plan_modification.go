package fwxschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// BlockWithPlanModifiers is an optional interface on Block which enables
// plan modification support.
type BlockWithPlanModifiers interface {
	// Implementations should include the fwschema.Block interface methods
	// for proper block handling.
	fwschema.Block

	// GetPlanModifiers should return a list of attribute-based plan modifiers.
	// This is named differently than PlanModifiers to prevent a conflict with
	// the tfsdk.Block field name.
	GetPlanModifiers() tfsdk.AttributePlanModifiers
}
