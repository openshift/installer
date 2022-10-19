package tfsdk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// AttributePlanModifier represents a modifier for an attribute at plan time.
// An AttributePlanModifier can only modify the planned value for the attribute
// on which it is defined. For plan-time modifications that modify the values of
// several attributes at once, please instead use the ResourceWithModifyPlan
// interface by defining a ModifyPlan function on the resource.
type AttributePlanModifier interface {
	// Description is used in various tooling, like the language server, to
	// give practitioners more information about what this modifier is,
	// what it's for, and how it should be used. It should be written as
	// plain text, with no special formatting.
	Description(context.Context) string

	// MarkdownDescription is used in various tooling, like the
	// documentation generator, to give practitioners more information
	// about what this modifier is, what it's for, and how it should be
	// used. It should be formatted using Markdown.
	MarkdownDescription(context.Context) string

	// Modify is called when the provider has an opportunity to modify
	// the plan: once during the plan phase when Terraform is determining
	// the diff that should be shown to the user for approval, and once
	// during the apply phase with any unknown values from configuration
	// filled in with their final values.
	//
	// The Modify function has access to the config, state, and plan for
	// both the attribute in question and the entire resource, but it can
	// only modify the value of the one attribute.
	//
	// Any returned errors will stop further execution of plan modifications
	// for this Attribute and any nested Attribute. Other Attribute at the same
	// or higher levels of the Schema will still execute any plan modifications
	// to ensure all warnings and errors across all root Attribute are
	// captured.
	//
	// Please see the documentation for ResourceWithModifyPlan#ModifyPlan
	// for further details.
	Modify(context.Context, ModifyAttributePlanRequest, *ModifyAttributePlanResponse)
}

// AttributePlanModifiers represents a sequence of AttributePlanModifiers, in
// order.
type AttributePlanModifiers []AttributePlanModifier

// ModifyAttributePlanRequest represents a request for the provider to modify an
// attribute value, or mark it as requiring replacement, at plan time. An
// instance of this request struct is supplied as an argument to the Modify
// function of an attribute's plan modifier(s).
type ModifyAttributePlanRequest struct {
	// AttributePath is the path of the attribute. Use this path for any
	// response diagnostics.
	AttributePath path.Path

	// AttributePathExpression is the expression matching the exact path of the
	// attribute.
	AttributePathExpression path.Expression

	// Config is the configuration the user supplied for the resource.
	Config Config

	// State is the current state of the resource.
	State State

	// Plan is the planned new state for the resource.
	Plan Plan

	// AttributeConfig is the configuration the user supplied for the attribute.
	AttributeConfig attr.Value

	// AttributeState is the current state of the attribute.
	AttributeState attr.Value

	// AttributePlan is the planned new state for the attribute.
	AttributePlan attr.Value

	// ProviderMeta is metadata from the provider_meta block of the module.
	ProviderMeta Config

	// Private is provider-defined resource private state data which was previously
	// stored with the resource state. This data is opaque to Terraform and does
	// not affect plan output. Any existing data is copied to
	// ModifyAttributePlanResponse.Private to prevent accidental private state data loss.
	//
	// The private state data is always the original data when the schema-based plan
	// modification began or, is updated as the logic traverses deeper into underlying
	// attributes.
	//
	// Use the GetKey method to read data. Use the SetKey method on
	// ModifyAttributePlanResponse.Private to update or remove a value.
	Private *privatestate.ProviderData
}

// ModifyAttributePlanResponse represents a response to a
// ModifyAttributePlanRequest. An instance of this response struct is supplied
// as an argument to the Modify function of an attribute's plan modifier(s).
type ModifyAttributePlanResponse struct {
	// AttributePlan is the planned new state for the attribute.
	AttributePlan attr.Value

	// RequiresReplace indicates whether a change in the attribute
	// requires replacement of the whole resource.
	RequiresReplace bool

	// Private is the private state resource data following the ModifyAttributePlan operation.
	// This field is pre-populated from ModifyAttributePlanRequest.Private and
	// can be modified during the resource's ModifyAttributePlan operation.
	//
	// The private state data is always the original data when the schema-based plan
	// modification began or, is updated as the logic traverses deeper into underlying
	// attributes.
	Private *privatestate.ProviderData

	// Diagnostics report errors or warnings related to determining the
	// planned state of the requested resource. Returning an empty slice
	// indicates a successful validation with no warnings or errors
	// generated.
	Diagnostics diag.Diagnostics
}
