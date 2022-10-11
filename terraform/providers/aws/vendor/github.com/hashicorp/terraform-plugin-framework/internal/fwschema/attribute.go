package fwschema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Attribute is the core interface required for implementing Terraform schema
// functionality that can accept a value. This is intended to be the first
// abstraction of tfsdk.Attribute functionality into data source, provider,
// and resource specific functionality.
//
// Refer to the internal/fwschema/fwxschema package for optional interfaces
// that define framework-specific functionality, such a plan modification and
// validation.
type Attribute interface {
	// Implementations should include the tftypes.AttributePathStepper
	// interface methods for proper path and data handling.
	tftypes.AttributePathStepper

	// Equal should return true if the other attribute is exactly equivalent.
	Equal(o Attribute) bool

	// FrameworkType should return the framework type, whether a direct type
	// or nested attributes type, for the attribute.
	//
	// When tfsdk.Attribute is removed, this should be deprecated and renamed
	// to Type() to match other interfaces.
	FrameworkType() attr.Type

	// GetAttributes should return the nested attributes of an attribute, if
	// applicable. This is named differently than Attribute to prevent a
	// conflict with the tfsdk.Attribute field name.
	GetAttributes() NestedAttributes

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

	// GetType should return the framework type of an attribute. This is named
	// differently than Type to prevent a conflict with the tfsdk.Attribute
	// field name.
	GetType() attr.Type

	// IsComputed should return true if the attribute configuration value is
	// computed. This is named differently than Computed to prevent a conflict
	// with the tfsdk.Attribute field name.
	IsComputed() bool

	// IsOptional should return true if the attribute configuration value is
	// optional. This is named differently than Optional to prevent a conflict
	// with the tfsdk.Attribute field name.
	IsOptional() bool

	// IsRequired should return true if the attribute configuration value is
	// required. This is named differently than Required to prevent a conflict
	// with the tfsdk.Attribute field name.
	IsRequired() bool

	// IsSensitive should return true if the attribute configuration value is
	// sensitive. This is named differently than Sensitive to prevent a
	// conflict with the tfsdk.Attribute field name.
	IsSensitive() bool
}
