package fwschema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Block is the core interface required for implementing Terraform schema
// functionality that structurally holds attributes and blocks. This is
// intended to be the first abstraction of tfsdk.Block functionality into
// data source, provider, and resource specific functionality.
//
// Refer to the internal/fwschema/fwxschema package for optional interfaces
// that define framework-specific functionality, such a plan modification and
// validation.
type Block interface {
	// Implementations should include the tftypes.AttributePathStepper
	// interface methods for proper path and data handling.
	tftypes.AttributePathStepper

	// Equal should return true if the other block is exactly equivalent.
	Equal(o Block) bool

	// GetAttributes should return the nested attributes of a block, if
	// applicable. This is named differently than Attributes to prevent a
	// conflict with the tfsdk.Block field name.
	GetAttributes() map[string]Attribute

	// GetBlocks should return the nested blocks of a block, if
	// applicable. This is named differently than Blocks to prevent a
	// conflict with the tfsdk.Block field name.
	GetBlocks() map[string]Block

	// GetDeprecationMessage should return a non-empty string if an attribute
	// is deprecated. This is named differently than DeprecationMessage to
	// prevent a conflict with the tfsdk.Attribute field name.
	GetDeprecationMessage() string

	// GetDescription should return a non-empty string if an attribute
	// has a plaintext description. This is named differently than Description
	// to prevent a conflict with the tfsdk.Attribute field name.
	GetDescription() string

	// GetMarkdownDescription should return a non-empty string if an attribute
	// has a Markdown description. This is named differently than
	// MarkdownDescription to prevent a conflict with the tfsdk.Attribute field
	// name.
	GetMarkdownDescription() string

	// GetMaxItems should return the max items of a block. This is named
	// differently than MaxItems to prevent a conflict with the tfsdk.Block
	// field name.
	GetMaxItems() int64

	// GetMinItems should return the min items of a block. This is named
	// differently than MinItems to prevent a conflict with the tfsdk.Block
	// field name.
	GetMinItems() int64

	// GetNestingMode should return the nesting mode of a block. This is named
	// differently than NestingMode to prevent a conflict with the tfsdk.Block
	// field name.
	GetNestingMode() BlockNestingMode

	// Type should return the framework type of a block.
	Type() attr.Type
}
