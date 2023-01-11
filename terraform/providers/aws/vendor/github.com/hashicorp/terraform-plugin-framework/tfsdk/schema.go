package tfsdk

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/totftypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	// ErrPathInsideAtomicAttribute is used with AttributeAtPath is called
	// on a path that doesn't have a schema associated with it, because
	// it's an element, attribute, or block of a complex type, not a nested
	// attribute.
	ErrPathInsideAtomicAttribute = errors.New("path leads to element, attribute, or block of a schema.Attribute that has no schema associated with it")

	// ErrPathIsBlock is used with AttributeAtPath is called on a path is a
	// block, not an attribute. Use blockAtPath on the path instead.
	ErrPathIsBlock = errors.New("path leads to block, not an attribute")
)

// Schema must satify the fwschema.Schema interface.
var _ fwschema.Schema = Schema{}

// Schema is used to define the shape of practitioner-provider information,
// like resources, data sources, and providers. Think of it as a type
// definition, but for Terraform.
type Schema struct {
	// Attributes are value fields inside the resource, provider, or data
	// source that the schema is defining. The map key should be the name
	// of the attribute, and the body defines how it behaves. Names must
	// only contain lowercase letters, numbers, and underscores. Names must
	// not collide with any Blocks names.
	//
	// In practitioner configurations, an equals sign (=) is required to set
	// the value. See also:
	//   https://www.terraform.io/docs/language/syntax/configuration.html
	//
	// Attributes are strongly preferred over Blocks.
	Attributes map[string]Attribute

	// Blocks are structural fields inside the resource, provider, or data
	// source that the schema is defining. The map key should be the name
	// of the block, and the body defines how it behaves. Names must
	// only contain lowercase letters, numbers, and underscores. Names must
	// not collide with any Attributes names.
	//
	// Blocks are by definition, structural, meaning they are implicitly
	// required in values.
	//
	// In practitioner configurations, an equals sign (=) cannot be used to
	// set the value. Blocks are instead repeated as necessary, or require
	// the use of dynamic block expressions. See also:
	//   https://www.terraform.io/docs/language/syntax/configuration.html
	//   https://www.terraform.io/docs/language/expressions/dynamic-blocks.html
	//
	// Attributes are preferred over Blocks. Blocks should typically be used
	// for configuration compatibility with previously existing schemas from
	// an older Terraform Plugin SDK. Efforts should be made to convert from
	// Blocks to Attributes as a breaking change for practitioners.
	Blocks map[string]Block

	// Version indicates the current version of the schema. Schemas are
	// versioned to help with automatic upgrade process. This is not
	// typically required unless there is a change in the schema, such as
	// changing an attribute type, that needs manual upgrade handling.
	// Versions should only be incremented by one each release.
	Version int64

	DeprecationMessage  string
	Description         string
	MarkdownDescription string
}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// schema.
func (s Schema) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	a, ok := step.(tftypes.AttributeName)

	if !ok {
		return nil, fmt.Errorf("cannot apply AttributePathStep %T to schema", step)
	}

	attrName := string(a)

	if attr, ok := s.Attributes[attrName]; ok {
		return attr, nil
	}

	if block, ok := s.Blocks[attrName]; ok {
		return block, nil
	}

	return nil, fmt.Errorf("could not find attribute or block %q in schema", a)
}

// TypeAtPath returns the framework type at the given schema path.
func (s Schema) TypeAtPath(ctx context.Context, schemaPath path.Path) (attr.Type, diag.Diagnostics) {
	var diags diag.Diagnostics

	tftypesPath, tftypesDiags := totftypes.AttributePath(ctx, schemaPath)

	diags.Append(tftypesDiags...)

	if diags.HasError() {
		return nil, diags
	}

	attrType, err := s.TypeAtTerraformPath(ctx, tftypesPath)

	if err != nil {
		diags.AddAttributeError(
			schemaPath,
			"Invalid Schema Path",
			"When attempting to get the framework type associated with a schema path, an unexpected error was returned. "+
				"This is always an issue with the provider. Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Path: %s\n", schemaPath.String())+
				fmt.Sprintf("Original Error: %s", err),
		)
		return nil, diags
	}

	return attrType, diags
}

// TypeAtTerraformPath returns the framework type at the given tftypes path.
func (s Schema) TypeAtTerraformPath(_ context.Context, path *tftypes.AttributePath) (attr.Type, error) {
	rawType, remaining, err := tftypes.WalkAttributePath(s, path)
	if err != nil {
		return nil, fmt.Errorf("%v still remains in the path: %w", remaining, err)
	}

	switch typ := rawType.(type) {
	case attr.Type:
		return typ, nil
	case fwschema.UnderlyingAttributes:
		return typ.Type(), nil
	case fwschema.NestedBlock:
		return typ.Block.Type(), nil
	case Attribute:
		return typ.FrameworkType(), nil
	case Block:
		return typ.Type(), nil
	case Schema:
		return typ.Type(), nil
	default:
		return nil, fmt.Errorf("got unexpected type %T", rawType)
	}
}

// GetAttributes satisfies the fwschema.Schema interface.
func (s Schema) GetAttributes() map[string]fwschema.Attribute {
	return schemaAttributes(s.Attributes)
}

// GetBlocks satisfies the fwschema.Schema interface.
func (s Schema) GetBlocks() map[string]fwschema.Block {
	return schemaBlocks(s.Blocks)
}

// GetDeprecationMessage satisfies the fwschema.Schema interface.
func (s Schema) GetDeprecationMessage() string {
	return s.DeprecationMessage
}

// GetDescription satisfies the fwschema.Schema interface.
func (s Schema) GetDescription() string {
	return s.Description
}

// GetMarkdownDescription satisfies the fwschema.Schema interface.
func (s Schema) GetMarkdownDescription() string {
	return s.MarkdownDescription
}

// GetVersion satisfies the fwschema.Schema interface.
func (s Schema) GetVersion() int64 {
	return s.Version
}

// Type returns the framework type of the schema.
func (s Schema) Type() attr.Type {
	attrTypes := map[string]attr.Type{}

	for name, attr := range s.Attributes {
		attrTypes[name] = attr.FrameworkType()
	}

	for name, block := range s.Blocks {
		attrTypes[name] = block.Type()
	}

	return types.ObjectType{AttrTypes: attrTypes}
}

// AttributeAtPath returns the Attribute at the passed path. If the path points
// to an element or attribute of a complex type, rather than to an Attribute,
// it will return an ErrPathInsideAtomicAttribute error.
func (s Schema) AttributeAtPath(ctx context.Context, schemaPath path.Path) (fwschema.Attribute, diag.Diagnostics) {
	var diags diag.Diagnostics

	tftypesPath, tftypesDiags := totftypes.AttributePath(ctx, schemaPath)

	diags.Append(tftypesDiags...)

	if diags.HasError() {
		return nil, diags
	}

	attribute, err := s.AttributeAtTerraformPath(ctx, tftypesPath)

	if err != nil {
		diags.AddAttributeError(
			schemaPath,
			"Invalid Schema Path",
			"When attempting to get the framework attribute associated with a schema path, an unexpected error was returned. "+
				"This is always an issue with the provider. Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Path: %s\n", schemaPath.String())+
				fmt.Sprintf("Original Error: %s", err),
		)
		return nil, diags
	}

	return attribute, diags
}

// AttributeAtPath returns the Attribute at the passed path. If the path points
// to an element or attribute of a complex type, rather than to an Attribute,
// it will return an ErrPathInsideAtomicAttribute error.
func (s Schema) AttributeAtTerraformPath(_ context.Context, path *tftypes.AttributePath) (fwschema.Attribute, error) {
	res, remaining, err := tftypes.WalkAttributePath(s, path)
	if err != nil {
		return Attribute{}, fmt.Errorf("%v still remains in the path: %w", remaining, err)
	}

	switch r := res.(type) {
	case attr.Type:
		return Attribute{}, ErrPathInsideAtomicAttribute
	case fwschema.UnderlyingAttributes:
		return Attribute{}, ErrPathInsideAtomicAttribute
	case fwschema.NestedBlock:
		return Attribute{}, ErrPathInsideAtomicAttribute
	case fwschema.Attribute:
		return r, nil
	case Block:
		return Attribute{}, ErrPathIsBlock
	default:
		return Attribute{}, fmt.Errorf("got unexpected type %T", res)
	}
}

// schemaAttributes is a tfsdk to fwschema type conversion function.
func schemaAttributes(attributes map[string]Attribute) map[string]fwschema.Attribute {
	result := make(map[string]fwschema.Attribute, len(attributes))

	for name, attribute := range attributes {
		result[name] = attribute
	}

	return result
}

// schemaBlocks is a tfsdk to fwschema type conversion function.
func schemaBlocks(blocks map[string]Block) map[string]fwschema.Block {
	result := make(map[string]fwschema.Block, len(blocks))

	for name, block := range blocks {
		result[name] = block
	}

	return result
}
