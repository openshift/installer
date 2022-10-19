package resource

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/totftypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// RequiresReplace returns an AttributePlanModifier specifying the attribute as
// requiring replacement. This behaviour is identical to the ForceNew behaviour
// in terraform-plugin-sdk and will result in the resource being destroyed and
// recreated when the following conditions are met:
//
// 1. The resource's state is not null; a null state indicates that we're
// creating a resource, and we never need to destroy and recreate a resource
// when we're creating it.
//
// 2. The resource's plan is not null; a null plan indicates that we're
// deleting a resource, and we never need to destroy and recreate a resource
// when we're deleting it.
//
// 3. The attribute's config is not null or the attribute is not computed; a
// computed attribute with a null config almost always means that the provider
// is changing the value, and practitioners are usually unpleasantly surprised
// when a resource is destroyed and recreated when their configuration hasn't
// changed. This has the unfortunate side effect that removing a computed field
// from the config will not trigger a destroy and recreate cycle, even when
// that is warranted. To get around this, provider developer can implement
// their own AttributePlanModifier that handles that behavior in the way that
// most makes sense for their use case.
//
// 4. The attribute's value in the plan does not match the attribute's value in
// the state.
func RequiresReplace() tfsdk.AttributePlanModifier {
	return requiresReplaceModifier{}
}

// requiresReplaceModifier is an AttributePlanModifier that sets RequiresReplace
// on the attribute.
type requiresReplaceModifier struct{}

// Modify fills the AttributePlanModifier interface. It sets RequiresReplace on
// the response to true if the following criteria are met:
//
// 1. The resource's state is not null; a null state indicates that we're
// creating a resource, and we never need to destroy and recreate a resource
// when we're creating it.
//
// 2. The resource's plan is not null; a null plan indicates that we're
// deleting a resource, and we never need to destroy and recreate a resource
// when we're deleting it.
//
// 3. The attribute's config is not null or the attribute is not computed; a
// computed attribute with a null config almost always means that the provider
// is changing the value, and practitioners are usually unpleasantly surprised
// when a resource is destroyed and recreated when their configuration hasn't
// changed. This has the unfortunate side effect that removing a computed field
// from the config will not trigger a destroy and recreate cycle, even when
// that is warranted. To get around this, provider developer can implement
// their own AttributePlanModifier that handles that behavior in the way that
// most makes sense for their use case.
//
// 4. The attribute's value in the plan does not match the attribute's value in
// the state.
func (r requiresReplaceModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if req.AttributeConfig == nil || req.AttributePlan == nil || req.AttributeState == nil {
		// shouldn't happen, but let's not panic if it does
		return
	}

	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to delete and
		// recreate it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and
		// recreate it
		return
	}

	// TODO: Remove after schema refactoring, Attribute is exposed in
	// ModifyAttributePlanRequest, or Computed is exposed in
	// ModifyAttributePlanRequest.
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/389
	tftypesPath, diags := totftypes.AttributePath(ctx, req.AttributePath)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	attrSchema, err := req.State.Schema.AttributeAtTerraformPath(ctx, tftypesPath)

	// Path may lead to block instead of attribute. Blocks cannot be Computed.
	// If ErrPathIsBlock, attrSchema.Computed will still be false later.
	if err != nil && !errors.Is(err, tfsdk.ErrPathIsBlock) {
		resp.Diagnostics.AddAttributeError(req.AttributePath,
			"Error finding attribute schema",
			fmt.Sprintf("An unexpected error was encountered retrieving the schema for this attribute. This is always a bug in the provider.\n\nError: %s", err),
		)
		return
	}

	if req.AttributeConfig.IsNull() && attrSchema.IsComputed() {
		// if the config is null and the attribute is computed, this
		// could be an out of band change, don't require replace
		return
	}

	if req.AttributePlan.Equal(req.AttributeState) {
		// if the plan and the state are in agreement, this attribute
		// isn't changing, don't require replace
		return
	}

	resp.RequiresReplace = true
}

// Description returns a human-readable description of the plan modifier.
func (r requiresReplaceModifier) Description(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (r requiresReplaceModifier) MarkdownDescription(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

// RequiresReplaceIf returns an AttributePlanModifier that mimics
// RequiresReplace, but only when the passed function `f` returns true. The
// resource will be destroyed and recreated if `f` returns true and the
// following conditions are met:
//
// 1. The resource's state is not null; a null state indicates that we're
// creating a resource, and we never need to destroy and recreate a resource
// when we're creating it.
//
// 2. The resource's plan is not null; a null plan indicates that we're
// deleting a resource, and we never need to destroy and recreate a resource
// when we're deleting it.
//
// 3. The attribute's config is not null or the attribute is not computed; a
// computed attribute with a null config almost always means that the provider
// is changing the value, and practitioners are usually unpleasantly surprised
// when a resource is destroyed and recreated when their configuration hasn't
// changed. This has the unfortunate side effect that removing a computed field
// from the config will not trigger a destroy and recreate cycle, even when
// that is warranted. To get around this, provider developer can implement
// their own AttributePlanModifier that handles that behavior in the way that
// most makes sense for their use case.
//
// 4. The attribute's value in the plan does not match the attribute's value in
// the state.
//
// If `f` does not return true, RequiresReplaceIf will *not* override prior
// AttributePlanModifiers' determination of whether the resource needs to be
// recreated or not. This allows for multiple RequiresReplaceIf (or other
// modifiers that sometimes set RequiresReplace) to be used on a single
// attribute without the last one in the list always determining the outcome.
func RequiresReplaceIf(f RequiresReplaceIfFunc, description, markdownDescription string) tfsdk.AttributePlanModifier {
	return requiresReplaceIfModifier{
		f:                   f,
		description:         description,
		markdownDescription: markdownDescription,
	}
}

// RequiresReplaceIfFunc is a conditional function used in the RequiresReplaceIf
// plan modifier to determine whether the attribute requires replacement.
type RequiresReplaceIfFunc func(ctx context.Context, state, config attr.Value, path path.Path) (bool, diag.Diagnostics)

// requiresReplaceIfModifier is an AttributePlanModifier that sets RequiresReplace
// on the attribute if the conditional function returns true.
type requiresReplaceIfModifier struct {
	f                   RequiresReplaceIfFunc
	description         string
	markdownDescription string
}

// Modify fills the AttributePlanModifier interface. It sets RequiresReplace on
// the response to true if the following criteria are met:
//
// 1. `f` returns true. If `f` returns false, the response will not be modified
// at all.
//
// 2. The resource's state is not null; a null state indicates that we're
// creating a resource, and we never need to destroy and recreate a resource
// when we're creating it.
//
// 3. The resource's plan is not null; a null plan indicates that we're
// deleting a resource, and we never need to destroy and recreate a resource
// when we're deleting it.
//
// 4. The attribute's config is not null or the attribute is not computed; a
// computed attribute with a null config almost always means that the provider
// is changing the value, and practitioners are usually unpleasantly surprised
// when a resource is destroyed and recreated when their configuration hasn't
// changed. This has the unfortunate side effect that removing a computed field
// from the config will not trigger a destroy and recreate cycle, even when
// that is warranted. To get around this, provider developer can implement
// their own AttributePlanModifier that handles that behavior in the way that
// most makes sense for their use case.
//
// 5. The attribute's value in the plan does not match the attribute's value in
// the state.
func (r requiresReplaceIfModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if req.AttributeConfig == nil || req.AttributePlan == nil || req.AttributeState == nil {
		// shouldn't happen, but let's not panic if it does
		return
	}

	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to delete and
		// recreate it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and
		// recreate it
		return
	}

	// TODO: Remove after schema refactoring, Attribute is exposed in
	// ModifyAttributePlanRequest, or Computed is exposed in
	// ModifyAttributePlanRequest.
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/389
	tftypesPath, diags := totftypes.AttributePath(ctx, req.AttributePath)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	attrSchema, err := req.State.Schema.AttributeAtTerraformPath(ctx, tftypesPath)

	// Path may lead to block instead of attribute. Blocks cannot be Computed.
	// If ErrPathIsBlock, attrSchema.Computed will still be false later.
	if err != nil && !errors.Is(err, tfsdk.ErrPathIsBlock) {
		resp.Diagnostics.AddAttributeError(req.AttributePath,
			"Error finding attribute schema",
			fmt.Sprintf("An unexpected error was encountered retrieving the schema for this attribute. This is always a bug in the provider.\n\nError: %s", err),
		)
		return
	}

	if req.AttributeConfig.IsNull() && attrSchema.IsComputed() {
		// if the config is null and the attribute is computed, this
		// could be an out of band change, don't require replace
		return
	}

	if req.AttributePlan.Equal(req.AttributeState) {
		// if the plan and the state are in agreement, this attribute
		// isn't changing, don't require replace
		return
	}

	res, diags := r.f(ctx, req.AttributeState, req.AttributeConfig, req.AttributePath)
	resp.Diagnostics.Append(diags...)

	// If the function says to require replacing, we require replacing.
	// If the function says not to, we don't change the value that prior
	// plan modifiers may have set.
	if res {
		resp.RequiresReplace = true
	} else if resp.RequiresReplace {
		logging.FrameworkDebug(ctx, "Keeping previous attribute replacement requirement")
	}
}

// Description returns a human-readable description of the plan modifier.
func (r requiresReplaceIfModifier) Description(ctx context.Context) string {
	return r.description
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (r requiresReplaceIfModifier) MarkdownDescription(ctx context.Context) string {
	return r.markdownDescription
}

// UseStateForUnknown returns an AttributePlanModifier that copies the prior
// state value for an attribute into that attribute's plan, if that state is
// non-null.
//
// Computed attributes without the UseStateForUnknown attribute plan modifier
// will have their value set to Unknown in the plan by the framework to prevent
// Terraform errors, so their value always will be displayed as "(known after
// apply)" in the CLI plan output. Using this plan modifier will instead
// display the prior state value in the plan, unless a prior plan modifier
// adjusts the value.
func UseStateForUnknown() tfsdk.AttributePlanModifier {
	return useStateForUnknownModifier{}
}

// useStateForUnknownModifier implements the UseStateForUnknown
// AttributePlanModifier.
type useStateForUnknownModifier struct{}

// Modify copies the attribute's prior state to the attribute plan if the prior
// state value is not null.
func (r useStateForUnknownModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if req.AttributeState == nil || resp.AttributePlan == nil || req.AttributeConfig == nil {
		return
	}

	// if we have no state value, there's nothing to preserve
	if req.AttributeState.IsNull() {
		return
	}

	// if it's not planned to be the unknown value, stick with the concrete plan
	if !resp.AttributePlan.IsUnknown() {
		return
	}

	// if the config is the unknown value, use the unknown value otherwise, interpolation gets messed up
	if req.AttributeConfig.IsUnknown() {
		return
	}

	resp.AttributePlan = req.AttributeState
}

// Description returns a human-readable description of the plan modifier.
func (r useStateForUnknownModifier) Description(ctx context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (r useStateForUnknownModifier) MarkdownDescription(ctx context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}
