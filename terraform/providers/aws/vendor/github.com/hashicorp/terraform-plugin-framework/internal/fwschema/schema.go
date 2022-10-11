package fwschema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Schema is the core interface required for data sources, providers, and
// resources.
type Schema interface {
	// Implementations should include the tftypes.AttributePathStepper
	// interface methods for proper path and data handling.
	tftypes.AttributePathStepper

	// AttributeAtPath should return the Attribute at the given path or return
	// an error.
	AttributeAtPath(context.Context, path.Path) (Attribute, diag.Diagnostics)

	// AttributeAtTerraformPath should return the Attribute at the given
	// Terraform path or return an error.
	AttributeAtTerraformPath(context.Context, *tftypes.AttributePath) (Attribute, error)

	// GetAttributes should return the attributes of a schema. This is named
	// differently than Attributes to prevent a conflict with the tfsdk.Schema
	// field name.
	GetAttributes() map[string]Attribute

	// GetBlocks should return the blocks of a schema. This is named
	// differently than Blocks to prevent a conflict with the tfsdk.Schema
	// field name.
	GetBlocks() map[string]Block

	// GetDeprecationMessage should return a non-empty string if a schema
	// is deprecated. This is named differently than DeprecationMessage to
	// prevent a conflict with the tfsdk.Schema field name.
	GetDeprecationMessage() string

	// GetDescription should return a non-empty string if a schema has a
	// plaintext description. This is named differently than Description
	// to prevent a conflict with the tfsdk.Schema field name.
	GetDescription() string

	// GetMarkdownDescription should return a non-empty string if a schema has
	// a Markdown description. This is named differently than
	// MarkdownDescription to prevent a conflict with the tfsdk.Schema field
	// name.
	GetMarkdownDescription() string

	// GetVersion should return the version of a schema. This is named
	// differently than Version to prevent a conflict with the tfsdk.Schema
	// field name.
	GetVersion() int64

	// Type should return the framework type of the schema.
	Type() attr.Type

	// TypeAtPath should return the framework type of the Attribute at the
	// the given path or return an error.
	TypeAtPath(context.Context, path.Path) (attr.Type, diag.Diagnostics)

	// AttributeTypeAtPath should return the framework type of the Attribute at
	// the given Terraform path or return an error.
	TypeAtTerraformPath(context.Context, *tftypes.AttributePath) (attr.Type, error)
}
