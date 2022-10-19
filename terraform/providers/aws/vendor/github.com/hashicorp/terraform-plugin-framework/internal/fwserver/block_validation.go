package fwserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BlockValidate performs all Block validation.
//
// TODO: Clean up this abstraction back into an internal Block type method.
// The extra Block parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
func BlockValidate(ctx context.Context, b fwschema.Block, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	blockWithValidators, ok := b.(fwxschema.BlockWithValidators)

	if !ok {
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

	req.AttributeConfig = attributeConfig

	for _, validator := range blockWithValidators.GetValidators() {
		validator.Validate(ctx, req, resp)
	}

	nm := b.GetNestingMode()
	switch nm {
	case fwschema.BlockNestingModeList:
		l, ok := req.AttributeConfig.(types.List)

		if !ok {
			err := fmt.Errorf("unknown block value type (%s) for nesting mode (%T) at path: %s", req.AttributeConfig.Type(ctx), nm, req.AttributePath)
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Block Validation Error",
				"Block validation cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
			)

			return
		}

		for idx := range l.Elems {
			for name, attr := range b.GetAttributes() {
				nestedAttrReq := tfsdk.ValidateAttributeRequest{
					AttributePath:           req.AttributePath.AtListIndex(idx).AtName(name),
					AttributePathExpression: req.AttributePathExpression.AtListIndex(idx).AtName(name),
					Config:                  req.Config,
				}
				nestedAttrResp := &tfsdk.ValidateAttributeResponse{
					Diagnostics: resp.Diagnostics,
				}

				AttributeValidate(ctx, attr, nestedAttrReq, nestedAttrResp)

				resp.Diagnostics = nestedAttrResp.Diagnostics
			}

			for name, block := range b.GetBlocks() {
				nestedAttrReq := tfsdk.ValidateAttributeRequest{
					AttributePath:           req.AttributePath.AtListIndex(idx).AtName(name),
					AttributePathExpression: req.AttributePathExpression.AtListIndex(idx).AtName(name),
					Config:                  req.Config,
				}
				nestedAttrResp := &tfsdk.ValidateAttributeResponse{
					Diagnostics: resp.Diagnostics,
				}

				BlockValidate(ctx, block, nestedAttrReq, nestedAttrResp)

				resp.Diagnostics = nestedAttrResp.Diagnostics
			}
		}

		// Terraform 0.12 through 0.15.1 do not implement block MaxItems
		// validation.
		//
		// Terraform 0.15.2 and later implements MaxItems validation during
		// configuration decoding, so if this framework drops Terraform support
		// for earlier versions, this validation can be removed.
		if b.GetMaxItems() > 0 && int64(len(l.Elems)) > b.GetMaxItems() {
			resp.Diagnostics.Append(blockMaxItemsDiagnostic(req.AttributePath, b.GetMaxItems(), len(l.Elems)))
		}

		// Terraform 0.12 through 0.15.1 implement conservative block MinItems
		// validation, where the MinItems can be reset to 1 in certain
		// situations. This validation must ensure the list itself is not
		// unknown, which could erroneously trigger the error since the list
		// would have 0 elements.
		//
		// Terraform 0.15.2 and later implements proper MinItems validation
		// during configuration decoding, so if this framework drops Terraform
		// support for earlier versions, this validation can be removed.
		if b.GetMinItems() > 0 && int64(len(l.Elems)) < b.GetMinItems() && !l.IsUnknown() {
			resp.Diagnostics.Append(blockMinItemsDiagnostic(req.AttributePath, b.GetMinItems(), len(l.Elems)))
		}
	case fwschema.BlockNestingModeSet:
		s, ok := req.AttributeConfig.(types.Set)

		if !ok {
			err := fmt.Errorf("unknown block value type (%s) for nesting mode (%T) at path: %s", req.AttributeConfig.Type(ctx), nm, req.AttributePath)
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Block Validation Error",
				"Block validation cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
			)

			return
		}

		for _, value := range s.Elems {
			for name, attr := range b.GetAttributes() {
				nestedAttrReq := tfsdk.ValidateAttributeRequest{
					AttributePath:           req.AttributePath.AtSetValue(value).AtName(name),
					AttributePathExpression: req.AttributePathExpression.AtSetValue(value).AtName(name),
					Config:                  req.Config,
				}
				nestedAttrResp := &tfsdk.ValidateAttributeResponse{
					Diagnostics: resp.Diagnostics,
				}

				AttributeValidate(ctx, attr, nestedAttrReq, nestedAttrResp)

				resp.Diagnostics = nestedAttrResp.Diagnostics
			}

			for name, block := range b.GetBlocks() {
				nestedAttrReq := tfsdk.ValidateAttributeRequest{
					AttributePath:           req.AttributePath.AtSetValue(value).AtName(name),
					AttributePathExpression: req.AttributePathExpression.AtSetValue(value).AtName(name),
					Config:                  req.Config,
				}
				nestedAttrResp := &tfsdk.ValidateAttributeResponse{
					Diagnostics: resp.Diagnostics,
				}

				BlockValidate(ctx, block, nestedAttrReq, nestedAttrResp)

				resp.Diagnostics = nestedAttrResp.Diagnostics
			}
		}

		// Terraform 0.12 through 0.15.1 do not implement block MaxItems
		// validation.
		//
		// Terraform 0.15.2 and later implements MaxItems validation during
		// configuration decoding, so if this framework drops Terraform support
		// for earlier versions, this validation can be removed.
		if b.GetMaxItems() > 0 && int64(len(s.Elems)) > b.GetMaxItems() {
			resp.Diagnostics.Append(blockMaxItemsDiagnostic(req.AttributePath, b.GetMaxItems(), len(s.Elems)))
		}

		// Terraform 0.12 through 0.15.1 implement conservative block MinItems
		// validation, where the MinItems can be reset to 1 in certain
		// situations. This validation must ensure the set itself is not
		// unknown, which could erroneously trigger the error since the set
		// would have 0 elements.
		//
		// Terraform 0.15.2 and later implements proper MinItems validation
		// during configuration decoding, so if this framework drops Terraform
		// support for earlier versions, this validation can be removed.
		if b.GetMinItems() > 0 && int64(len(s.Elems)) < b.GetMinItems() && !s.IsUnknown() {
			resp.Diagnostics.Append(blockMinItemsDiagnostic(req.AttributePath, b.GetMinItems(), len(s.Elems)))
		}
	case fwschema.BlockNestingModeSingle:
		s, ok := req.AttributeConfig.(types.Object)

		if !ok {
			err := fmt.Errorf("unknown block value type (%s) for nesting mode (%T) at path: %s", req.AttributeConfig.Type(ctx), nm, req.AttributePath)
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Block Validation Error",
				"Block validation cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
			)

			return
		}

		for name, attr := range b.GetAttributes() {
			nestedAttrReq := tfsdk.ValidateAttributeRequest{
				AttributePath:           req.AttributePath.AtName(name),
				AttributePathExpression: req.AttributePathExpression.AtName(name),
				Config:                  req.Config,
			}
			nestedAttrResp := &tfsdk.ValidateAttributeResponse{
				Diagnostics: resp.Diagnostics,
			}

			AttributeValidate(ctx, attr, nestedAttrReq, nestedAttrResp)

			resp.Diagnostics = nestedAttrResp.Diagnostics
		}

		for name, block := range b.GetBlocks() {
			nestedAttrReq := tfsdk.ValidateAttributeRequest{
				AttributePath:           req.AttributePath.AtName(name),
				AttributePathExpression: req.AttributePathExpression.AtName(name),
				Config:                  req.Config,
			}
			nestedAttrResp := &tfsdk.ValidateAttributeResponse{
				Diagnostics: resp.Diagnostics,
			}

			BlockValidate(ctx, block, nestedAttrReq, nestedAttrResp)

			resp.Diagnostics = nestedAttrResp.Diagnostics
		}

		if b.GetMinItems() == 1 && s.IsNull() {
			resp.Diagnostics.Append(blockMinItemsDiagnostic(req.AttributePath, b.GetMinItems(), 0))
		}
	default:
		err := fmt.Errorf("unknown block validation nesting mode (%T: %v) at path: %s", nm, nm, req.AttributePath)
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Block Validation Error",
			"Block validation cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
		)

		return
	}

	// Show deprecation warning only on known values.
	if b.GetDeprecationMessage() != "" && !attributeConfig.IsNull() && !attributeConfig.IsUnknown() {
		resp.Diagnostics.AddAttributeWarning(
			req.AttributePath,
			"Block Deprecated",
			b.GetDeprecationMessage(),
		)
	}
}

func blockMaxItemsDiagnostic(attrPath path.Path, maxItems int64, elements int) diag.Diagnostic {
	var details strings.Builder

	details.WriteString("The configuration should declare a maximum of ")

	if maxItems == 1 {
		details.WriteString("1 block")
	} else {
		details.WriteString(fmt.Sprintf("%d blocks", maxItems))
	}

	// Elements will always be greater than 1, so do not need to handle the
	// singular case.
	details.WriteString(fmt.Sprintf(", however %d blocks were configured.", elements))

	return diag.NewAttributeErrorDiagnostic(
		attrPath,
		"Extra Block Configuration",
		details.String(),
	)
}

func blockMinItemsDiagnostic(attrPath path.Path, minItems int64, elements int) diag.Diagnostic {
	var details strings.Builder

	details.WriteString("The configuration should declare a minimum of ")

	if minItems == 1 {
		details.WriteString("1 block")
	} else {
		details.WriteString(fmt.Sprintf("%d blocks", minItems))
	}

	details.WriteString(", however ")

	if elements == 1 {
		details.WriteString("1 block was")
	} else {
		details.WriteString(fmt.Sprintf("%d blocks were", elements))
	}

	details.WriteString(" configured.")

	return diag.NewAttributeErrorDiagnostic(
		attrPath,
		"Missing Block Configuration",
		details.String(),
	)
}
