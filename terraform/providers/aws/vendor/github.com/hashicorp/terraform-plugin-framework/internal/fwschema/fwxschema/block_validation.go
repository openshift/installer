package fwxschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// BlockWithValidators is an optional interface on Block which enables
// validation support.
type BlockWithValidators interface {
	// Implementations should include the fwschema.Block interface methods
	// for proper block handling.
	fwschema.Block

	// GetValidators should return a list of attribute-based validators. This
	// is named differently than Validators to prevent a conflict with the
	// tfsdk.Block field name.
	GetValidators() []tfsdk.AttributeValidator
}
