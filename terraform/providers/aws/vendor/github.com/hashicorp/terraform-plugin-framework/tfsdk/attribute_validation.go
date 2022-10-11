package tfsdk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// AttributeValidator describes reusable Attribute validation functionality.
type AttributeValidator interface {
	// Description describes the validation in plain text formatting.
	//
	// This information may be automatically added to schema plain text
	// descriptions by external tooling.
	Description(context.Context) string

	// MarkdownDescription describes the validation in Markdown formatting.
	//
	// This information may be automatically added to schema Markdown
	// descriptions by external tooling.
	MarkdownDescription(context.Context) string

	// Validate performs the validation.
	Validate(context.Context, ValidateAttributeRequest, *ValidateAttributeResponse)
}

// ValidateAttributeRequest repesents a request for attribute validation.
type ValidateAttributeRequest struct {
	// AttributePath contains the path of the attribute. Use this path for any
	// response diagnostics.
	AttributePath path.Path

	// AttributePathExpression contains the expression matching the exact path
	// of the attribute.
	AttributePathExpression path.Expression

	// AttributeConfig contains the value of the attribute in the configuration.
	AttributeConfig attr.Value

	// Config contains the entire configuration of the data source, provider, or resource.
	Config Config
}

// ValidateAttributeResponse represents a response to a
// ValidateAttributeRequest. An instance of this response struct is
// automatically passed through to each AttributeValidator.
type ValidateAttributeResponse struct {
	// Diagnostics report errors or warnings related to validating the data
	// source configuration. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics
}
