package fwserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AttributeValidate performs all Attribute validation.
//
// TODO: Clean up this abstraction back into an internal Attribute type method.
// The extra Attribute parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
func AttributeValidate(ctx context.Context, a fwschema.Attribute, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	ctx = logging.FrameworkWithAttributePath(ctx, req.AttributePath.String())

	if (a.GetAttributes() == nil || len(a.GetAttributes().GetAttributes()) == 0) && a.GetType() == nil {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Attribute Definition",
			"Attribute must define either Attributes or Type. This is always a problem with the provider and should be reported to the provider developer.",
		)

		return
	}

	if a.GetAttributes() != nil && len(a.GetAttributes().GetAttributes()) > 0 && a.GetType() != nil {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Attribute Definition",
			"Attribute cannot define both Attributes and Type. This is always a problem with the provider and should be reported to the provider developer.",
		)

		return
	}

	if !a.IsRequired() && !a.IsOptional() && !a.IsComputed() {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Attribute Definition",
			"Attribute missing Required, Optional, or Computed definition. This is always a problem with the provider and should be reported to the provider developer.",
		)

		return
	}

	configData := &fwschemadata.Data{
		Description:    fwschemadata.DataDescriptionConfiguration,
		Schema:         req.Config.Schema,
		TerraformValue: req.Config.Raw,
	}

	attributeConfig, diags := configData.ValueAtPath(ctx, req.AttributePath)
	resp.Diagnostics.Append(diags...)

	if diags.HasError() {
		return
	}

	// Terraform CLI does not automatically perform certain configuration
	// checks yet. If it eventually does, this logic should remain at least
	// until Terraform CLI versions 0.12 through the release containing the
	// checks are considered end-of-life.
	// Reference: https://github.com/hashicorp/terraform/issues/30669
	if a.IsComputed() && !a.IsOptional() && !attributeConfig.IsNull() {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Configuration for Read-Only Attribute",
			"Cannot set value for this attribute as the provider has marked it as read-only. Remove the configuration line setting the value.\n\n"+
				"Refer to the provider documentation or contact the provider developers for additional information about configurable and read-only attributes that are supported.",
		)
	}

	// Terraform CLI does not automatically perform certain configuration
	// checks yet. If it eventually does, this logic should remain at least
	// until Terraform CLI versions 0.12 through the release containing the
	// checks are considered end-of-life.
	// Reference: https://github.com/hashicorp/terraform/issues/30669
	if a.IsRequired() && attributeConfig.IsNull() {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Missing Configuration for Required Attribute",
			fmt.Sprintf("Must set a configuration value for the %s attribute as the provider has marked it as required.\n\n", req.AttributePath.String())+
				"Refer to the provider documentation or contact the provider developers for additional information about configurable attributes that are required.",
		)
	}

	req.AttributeConfig = attributeConfig

	if attributeWithValidators, ok := a.(fwxschema.AttributeWithValidators); ok {
		for _, validator := range attributeWithValidators.GetValidators() {
			logging.FrameworkDebug(
				ctx,
				"Calling provider defined AttributeValidator",
				map[string]interface{}{
					logging.KeyDescription: validator.Description(ctx),
				},
			)
			validator.Validate(ctx, req, resp)
			logging.FrameworkDebug(
				ctx,
				"Called provider defined AttributeValidator",
				map[string]interface{}{
					logging.KeyDescription: validator.Description(ctx),
				},
			)
		}
	}

	AttributeValidateNestedAttributes(ctx, a, req, resp)

	// Show deprecation warnings only for known values.
	if a.GetDeprecationMessage() != "" && !attributeConfig.IsNull() && !attributeConfig.IsUnknown() {
		resp.Diagnostics.AddAttributeWarning(
			req.AttributePath,
			"Attribute Deprecated",
			a.GetDeprecationMessage(),
		)
	}
}

// AttributeValidateNestedAttributes performs all nested Attributes validation.
//
// TODO: Clean up this abstraction back into an internal Attribute type method.
// The extra Attribute parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
func AttributeValidateNestedAttributes(ctx context.Context, a fwschema.Attribute, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	if a.GetAttributes() == nil || len(a.GetAttributes().GetAttributes()) == 0 {
		return
	}

	nm := a.GetAttributes().GetNestingMode()
	switch nm {
	case fwschema.NestingModeList:
		l, ok := req.AttributeConfig.(types.List)

		if !ok {
			err := fmt.Errorf("unknown attribute value type (%T) for nesting mode (%T) at path: %s", req.AttributeConfig, nm, req.AttributePath)
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Attribute Validation Error",
				"Attribute validation cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
			)

			return
		}

		for idx := range l.Elems {
			for nestedName, nestedAttr := range a.GetAttributes().GetAttributes() {
				nestedAttrReq := tfsdk.ValidateAttributeRequest{
					AttributePath:           req.AttributePath.AtListIndex(idx).AtName(nestedName),
					AttributePathExpression: req.AttributePathExpression.AtListIndex(idx).AtName(nestedName),
					Config:                  req.Config,
				}
				nestedAttrResp := &tfsdk.ValidateAttributeResponse{
					Diagnostics: resp.Diagnostics,
				}

				AttributeValidate(ctx, nestedAttr, nestedAttrReq, nestedAttrResp)

				resp.Diagnostics = nestedAttrResp.Diagnostics
			}
		}
	case fwschema.NestingModeSet:
		s, ok := req.AttributeConfig.(types.Set)

		if !ok {
			err := fmt.Errorf("unknown attribute value type (%T) for nesting mode (%T) at path: %s", req.AttributeConfig, nm, req.AttributePath)
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Attribute Validation Error",
				"Attribute validation cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
			)

			return
		}

		for _, value := range s.Elems {
			for nestedName, nestedAttr := range a.GetAttributes().GetAttributes() {
				nestedAttrReq := tfsdk.ValidateAttributeRequest{
					AttributePath:           req.AttributePath.AtSetValue(value).AtName(nestedName),
					AttributePathExpression: req.AttributePathExpression.AtSetValue(value).AtName(nestedName),
					Config:                  req.Config,
				}
				nestedAttrResp := &tfsdk.ValidateAttributeResponse{
					Diagnostics: resp.Diagnostics,
				}

				AttributeValidate(ctx, nestedAttr, nestedAttrReq, nestedAttrResp)

				resp.Diagnostics = nestedAttrResp.Diagnostics
			}
		}
	case fwschema.NestingModeMap:
		m, ok := req.AttributeConfig.(types.Map)

		if !ok {
			err := fmt.Errorf("unknown attribute value type (%T) for nesting mode (%T) at path: %s", req.AttributeConfig, nm, req.AttributePath)
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Attribute Validation Error",
				"Attribute validation cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
			)

			return
		}

		for key := range m.Elems {
			for nestedName, nestedAttr := range a.GetAttributes().GetAttributes() {
				nestedAttrReq := tfsdk.ValidateAttributeRequest{
					AttributePath:           req.AttributePath.AtMapKey(key).AtName(nestedName),
					AttributePathExpression: req.AttributePathExpression.AtMapKey(key).AtName(nestedName),
					Config:                  req.Config,
				}
				nestedAttrResp := &tfsdk.ValidateAttributeResponse{
					Diagnostics: resp.Diagnostics,
				}

				AttributeValidate(ctx, nestedAttr, nestedAttrReq, nestedAttrResp)

				resp.Diagnostics = nestedAttrResp.Diagnostics
			}
		}
	case fwschema.NestingModeSingle:
		o, ok := req.AttributeConfig.(types.Object)

		if !ok {
			err := fmt.Errorf("unknown attribute value type (%T) for nesting mode (%T) at path: %s", req.AttributeConfig, nm, req.AttributePath)
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Attribute Validation Error",
				"Attribute validation cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
			)

			return
		}

		if !o.Null && !o.Unknown {
			for nestedName, nestedAttr := range a.GetAttributes().GetAttributes() {
				nestedAttrReq := tfsdk.ValidateAttributeRequest{
					AttributePath:           req.AttributePath.AtName(nestedName),
					AttributePathExpression: req.AttributePathExpression.AtName(nestedName),
					Config:                  req.Config,
				}
				nestedAttrResp := &tfsdk.ValidateAttributeResponse{
					Diagnostics: resp.Diagnostics,
				}

				AttributeValidate(ctx, nestedAttr, nestedAttrReq, nestedAttrResp)

				resp.Diagnostics = nestedAttrResp.Diagnostics
			}
		}
	default:
		err := fmt.Errorf("unknown attribute validation nesting mode (%T: %v) at path: %s", nm, nm, req.AttributePath)
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Attribute Validation Error",
			"Attribute validation cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
		)

		return
	}
}
