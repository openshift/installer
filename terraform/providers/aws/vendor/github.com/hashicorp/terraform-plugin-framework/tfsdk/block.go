package tfsdk

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ tftypes.AttributePathStepper = Block{}

// Block must satify the fwschema.Block interface. It must also satisfy
// fwxschema.BlockWithPlanModifiers and fwxschema.BlockWithValidators
// interfaces, however we cannot check that here or it would introduce an
// import cycle.
var _ fwschema.Block = Block{}

// Block defines the constraints and behaviors of a single structural field in a
// schema.
//
// The NestingMode field must be set or a runtime error will be raised by the
// framework when fetching the schema.
type Block struct {
	// Attributes are value fields inside the block. This map of attributes
	// behaves exactly like the map of attributes on the Schema type.
	Attributes map[string]Attribute

	// Blocks can have their own nested blocks. This nested map of blocks
	// behaves exactly like the map of blocks on the Schema type.
	Blocks map[string]Block

	// DeprecationMessage defines warning diagnostic details to display to
	// practitioners configuring this Block. The warning diagnostic summary
	// is automatically set to "Block Deprecated" along with configuration
	// source file and line information.
	//
	// This warning diagnostic is only displayed during Terraform's validation
	// phase when this field is a non-empty string and if the practitioner
	// configuration attempts to set the block value to a known or unknown
	// value (which may eventually be null).
	//
	// Set this field to a practitioner actionable message such as:
	//
	//     - "Configure other_attribute instead. This block will be removed
	//       in the next major version of the provider."
	//     - "Remove this block's configuration as it no longer is used and
	//       the block will be removed in the next major version of the
	//       provider."
	//
	DeprecationMessage string

	// Description is used in various tooling, like the language server, to
	// give practitioners more information about what this attribute is,
	// what it's for, and how it should be used. It should be written as
	// plain text, with no special formatting.
	Description string

	// MarkdownDescription is used in various tooling, like the
	// documentation generator, to give practitioners more information
	// about what this attribute is, what it's for, and how it should be
	// used. It should be formatted using Markdown.
	MarkdownDescription string

	// MaxItems is the maximum number of blocks that can be present in a
	// practitioner configuration.
	MaxItems int64

	// MinItems is the minimum number of blocks that must be present in a
	// practitioner configuration. Setting to 1 or above effectively marks
	// this configuration as required.
	MinItems int64

	// NestingMode indicates the block kind. This field must be set or a
	// runtime error will be raised by the framework when fetching the schema.
	NestingMode fwschema.BlockNestingMode

	// PlanModifiers defines a sequence of modifiers for this block at
	// plan time. Block-level plan modifications occur before any
	// resource-level plan modifications.
	//
	// Any errors will prevent further execution of this sequence
	// of modifiers and modifiers associated with any nested Attribute or
	// Block, but will not prevent execution of PlanModifiers on any
	// other Attribute or Block in the Schema.
	//
	// Plan modification only applies to resources, not data sources or
	// providers. Setting PlanModifiers on a data source or provider attribute
	// will have no effect.
	PlanModifiers AttributePlanModifiers

	// Validators defines validation functionality for the block.
	Validators []AttributeValidator
}

// ApplyTerraform5AttributePathStep allows Blocks to be walked using
// tftypes.Walk and tftypes.Transform.
func (b Block) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	switch b.NestingMode {
	case BlockNestingModeList:
		_, ok := step.(tftypes.ElementKeyInt)

		if !ok {
			return nil, fmt.Errorf("can't apply %T to block NestingModeList", step)
		}

		return fwschema.NestedBlock{Block: b}, nil
	case BlockNestingModeSet:
		_, ok := step.(tftypes.ElementKeyValue)

		if !ok {
			return nil, fmt.Errorf("can't apply %T to block NestingModeSet", step)
		}

		return fwschema.NestedBlock{Block: b}, nil
	case BlockNestingModeSingle:
		_, ok := step.(tftypes.AttributeName)

		if !ok {
			return nil, fmt.Errorf("can't apply %T to block NestingModeSingle", step)
		}

		return fwschema.NestedBlock{Block: b}.ApplyTerraform5AttributePathStep(step)
	default:
		return nil, fmt.Errorf("unsupported block nesting mode: %v", b.NestingMode)
	}
}

// Equal returns true if `b` and `o` should be considered Equal.
func (b Block) Equal(o fwschema.Block) bool {
	if !cmp.Equal(b.GetAttributes(), o.GetAttributes()) {
		return false
	}
	if !cmp.Equal(b.GetBlocks(), o.GetBlocks()) {
		return false
	}
	if b.GetDeprecationMessage() != o.GetDeprecationMessage() {
		return false
	}
	if b.GetDescription() != o.GetDescription() {
		return false
	}
	if b.GetMarkdownDescription() != o.GetMarkdownDescription() {
		return false
	}
	if b.GetMaxItems() != o.GetMaxItems() {
		return false
	}
	if b.GetMinItems() != o.GetMinItems() {
		return false
	}
	if b.GetNestingMode() != o.GetNestingMode() {
		return false
	}
	return true
}

// GetAttributes satisfies the fwschema.Block interface.
func (b Block) GetAttributes() map[string]fwschema.Attribute {
	return schemaAttributes(b.Attributes)
}

// GetBlocks satisfies the fwschema.Block interface.
func (b Block) GetBlocks() map[string]fwschema.Block {
	return schemaBlocks(b.Blocks)
}

// GetDeprecationMessage satisfies the fwschema.Block interface.
func (b Block) GetDeprecationMessage() string {
	return b.DeprecationMessage
}

// GetDescription satisfies the fwschema.Block interface.
func (b Block) GetDescription() string {
	return b.Description
}

// GetMarkdownDescription satisfies the fwschema.Block interface.
func (b Block) GetMarkdownDescription() string {
	return b.MarkdownDescription
}

// GetMaxItems satisfies the fwschema.Block interface.
func (b Block) GetMaxItems() int64 {
	return b.MaxItems
}

// GetMinItems satisfies the fwschema.Block interface.
func (b Block) GetMinItems() int64 {
	return b.MinItems
}

// GetNestingMode satisfies the fwschema.Block interface.
func (b Block) GetNestingMode() fwschema.BlockNestingMode {
	return b.NestingMode
}

// GetPlanModifiers satisfies the fwxschema.BlockWithPlanModifiers
// interface.
func (b Block) GetPlanModifiers() AttributePlanModifiers {
	return b.PlanModifiers
}

// GetValidators satisfies the fwxschema.BlockWithValidators interface.
func (b Block) GetValidators() []AttributeValidator {
	return b.Validators
}

// attributeType returns an attr.Type corresponding to the block.
func (b Block) Type() attr.Type {
	attrType := types.ObjectType{
		AttrTypes: map[string]attr.Type{},
	}

	for attrName, attr := range b.Attributes {
		attrType.AttrTypes[attrName] = attr.FrameworkType()
	}

	for blockName, block := range b.Blocks {
		attrType.AttrTypes[blockName] = block.Type()
	}

	switch b.NestingMode {
	case BlockNestingModeList:
		return types.ListType{
			ElemType: attrType,
		}
	case BlockNestingModeSet:
		return types.SetType{
			ElemType: attrType,
		}
	case BlockNestingModeSingle:
		return attrType
	default:
		panic(fmt.Sprintf("unsupported block nesting mode: %v", b.NestingMode))
	}
}
