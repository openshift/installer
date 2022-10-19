package fwserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// BlockModifyPlan performs all Block plan modification.
//
// TODO: Clean up this abstraction back into an internal Block type method.
// The extra Block parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
func BlockModifyPlan(ctx context.Context, b fwschema.Block, req tfsdk.ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	var requiresReplace bool

	privateProviderData := privatestate.EmptyProviderData(ctx)

	if req.Private != nil {
		resp.Private = req.Private
		privateProviderData = req.Private
	}

	if blockWithPlanModifiers, ok := b.(fwxschema.BlockWithPlanModifiers); ok {
		for _, planModifier := range blockWithPlanModifiers.GetPlanModifiers() {
			modifyResp := &tfsdk.ModifyAttributePlanResponse{
				AttributePlan:   req.AttributePlan,
				RequiresReplace: requiresReplace,
				Private:         privateProviderData,
			}

			planModifier.Modify(ctx, req, modifyResp)

			req.AttributePlan = modifyResp.AttributePlan
			resp.Diagnostics.Append(modifyResp.Diagnostics...)
			requiresReplace = modifyResp.RequiresReplace
			resp.AttributePlan = modifyResp.AttributePlan
			resp.Private = modifyResp.Private

			// Only on new errors.
			if modifyResp.Diagnostics.HasError() {
				return
			}
		}
	}

	if requiresReplace {
		resp.RequiresReplace = append(resp.RequiresReplace, req.AttributePath)
	}

	nm := b.GetNestingMode()
	switch nm {
	case fwschema.BlockNestingModeList:
		configList, diags := coerceListValue(req.AttributePath, req.AttributeConfig)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		planList, diags := coerceListValue(req.AttributePath, req.AttributePlan)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		stateList, diags := coerceListValue(req.AttributePath, req.AttributeState)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		for idx, planElem := range planList.Elems {
			attrPath := req.AttributePath.AtListIndex(idx)

			configObject, diags := listElemObject(ctx, attrPath, configList, idx, fwschemadata.DataDescriptionConfiguration)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			planObject, diags := coerceObjectValue(attrPath, planElem)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			stateObject, diags := listElemObject(ctx, attrPath, stateList, idx, fwschemadata.DataDescriptionState)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			for name, attr := range b.GetAttributes() {
				attrConfig, diags := objectAttributeValue(ctx, configObject, name, fwschemadata.DataDescriptionConfiguration)

				resp.Diagnostics.Append(diags...)

				if resp.Diagnostics.HasError() {
					return
				}

				attrPlan, diags := objectAttributeValue(ctx, planObject, name, fwschemadata.DataDescriptionPlan)

				resp.Diagnostics.Append(diags...)

				if resp.Diagnostics.HasError() {
					return
				}

				attrState, diags := objectAttributeValue(ctx, stateObject, name, fwschemadata.DataDescriptionState)

				resp.Diagnostics.Append(diags...)

				if resp.Diagnostics.HasError() {
					return
				}

				attrReq := tfsdk.ModifyAttributePlanRequest{
					AttributeConfig: attrConfig,
					AttributePath:   attrPath.AtName(name),
					AttributePlan:   attrPlan,
					AttributeState:  attrState,
					Config:          req.Config,
					Plan:            req.Plan,
					ProviderMeta:    req.ProviderMeta,
					State:           req.State,
					Private:         resp.Private,
				}
				attrResp := ModifyAttributePlanResponse{
					AttributePlan:   attrReq.AttributePlan,
					RequiresReplace: resp.RequiresReplace,
					Private:         attrReq.Private,
				}

				AttributeModifyPlan(ctx, attr, attrReq, &attrResp)

				planObject.Attrs[name] = attrResp.AttributePlan
				resp.Diagnostics.Append(attrResp.Diagnostics...)
				resp.RequiresReplace = attrResp.RequiresReplace
				resp.Private = attrResp.Private
			}

			for name, block := range b.GetBlocks() {
				attrConfig, diags := objectAttributeValue(ctx, configObject, name, fwschemadata.DataDescriptionConfiguration)

				resp.Diagnostics.Append(diags...)

				if resp.Diagnostics.HasError() {
					return
				}

				attrPlan, diags := objectAttributeValue(ctx, planObject, name, fwschemadata.DataDescriptionPlan)

				resp.Diagnostics.Append(diags...)

				if resp.Diagnostics.HasError() {
					return
				}

				attrState, diags := objectAttributeValue(ctx, stateObject, name, fwschemadata.DataDescriptionState)

				resp.Diagnostics.Append(diags...)

				if resp.Diagnostics.HasError() {
					return
				}

				blockReq := tfsdk.ModifyAttributePlanRequest{
					AttributeConfig: attrConfig,
					AttributePath:   req.AttributePath.AtListIndex(idx).AtName(name),
					AttributePlan:   attrPlan,
					AttributeState:  attrState,
					Config:          req.Config,
					Plan:            req.Plan,
					ProviderMeta:    req.ProviderMeta,
					State:           req.State,
					Private:         resp.Private,
				}
				blockResp := ModifyAttributePlanResponse{
					AttributePlan:   blockReq.AttributePlan,
					RequiresReplace: resp.RequiresReplace,
					Private:         blockReq.Private,
				}

				BlockModifyPlan(ctx, block, blockReq, &blockResp)

				planObject.Attrs[name] = blockResp.AttributePlan
				resp.Diagnostics.Append(blockResp.Diagnostics...)
				resp.RequiresReplace = blockResp.RequiresReplace
				resp.Private = blockResp.Private
			}

			planList.Elems[idx] = planObject
		}

		resp.AttributePlan = planList
	case fwschema.BlockNestingModeSet:
		configSet, diags := coerceSetValue(req.AttributePath, req.AttributeConfig)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		planSet, diags := coerceSetValue(req.AttributePath, req.AttributePlan)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		stateSet, diags := coerceSetValue(req.AttributePath, req.AttributeState)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		for idx, planElem := range planSet.Elems {
			attrPath := req.AttributePath.AtSetValue(planElem)

			configObject, diags := setElemObject(ctx, attrPath, configSet, idx, fwschemadata.DataDescriptionConfiguration)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			planObject, diags := coerceObjectValue(attrPath, planElem)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			stateObject, diags := setElemObject(ctx, attrPath, stateSet, idx, fwschemadata.DataDescriptionState)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			for name, attr := range b.GetAttributes() {
				attrConfig, diags := objectAttributeValue(ctx, configObject, name, fwschemadata.DataDescriptionConfiguration)

				resp.Diagnostics.Append(diags...)

				if resp.Diagnostics.HasError() {
					return
				}

				attrPlan, diags := objectAttributeValue(ctx, planObject, name, fwschemadata.DataDescriptionPlan)

				resp.Diagnostics.Append(diags...)

				if resp.Diagnostics.HasError() {
					return
				}

				attrState, diags := objectAttributeValue(ctx, stateObject, name, fwschemadata.DataDescriptionState)

				resp.Diagnostics.Append(diags...)

				if resp.Diagnostics.HasError() {
					return
				}

				attrReq := tfsdk.ModifyAttributePlanRequest{
					AttributeConfig: attrConfig,
					AttributePath:   attrPath.AtName(name),
					AttributePlan:   attrPlan,
					AttributeState:  attrState,
					Config:          req.Config,
					Plan:            req.Plan,
					ProviderMeta:    req.ProviderMeta,
					State:           req.State,
					Private:         resp.Private,
				}
				attrResp := ModifyAttributePlanResponse{
					AttributePlan:   attrReq.AttributePlan,
					RequiresReplace: resp.RequiresReplace,
					Private:         attrReq.Private,
				}

				AttributeModifyPlan(ctx, attr, attrReq, &attrResp)

				planObject.Attrs[name] = attrResp.AttributePlan
				resp.Diagnostics.Append(attrResp.Diagnostics...)
				resp.RequiresReplace = attrResp.RequiresReplace
				resp.Private = attrResp.Private
			}

			for name, block := range b.GetBlocks() {
				attrConfig, diags := objectAttributeValue(ctx, configObject, name, fwschemadata.DataDescriptionConfiguration)

				resp.Diagnostics.Append(diags...)

				if resp.Diagnostics.HasError() {
					return
				}

				attrPlan, diags := objectAttributeValue(ctx, planObject, name, fwschemadata.DataDescriptionPlan)

				resp.Diagnostics.Append(diags...)

				if resp.Diagnostics.HasError() {
					return
				}

				attrState, diags := objectAttributeValue(ctx, stateObject, name, fwschemadata.DataDescriptionState)

				resp.Diagnostics.Append(diags...)

				if resp.Diagnostics.HasError() {
					return
				}

				blockReq := tfsdk.ModifyAttributePlanRequest{
					AttributeConfig: attrConfig,
					AttributePath:   attrPath.AtName(name),
					AttributePlan:   attrPlan,
					AttributeState:  attrState,
					Config:          req.Config,
					Plan:            req.Plan,
					ProviderMeta:    req.ProviderMeta,
					State:           req.State,
					Private:         resp.Private,
				}
				blockResp := ModifyAttributePlanResponse{
					AttributePlan:   blockReq.AttributePlan,
					RequiresReplace: resp.RequiresReplace,
					Private:         blockReq.Private,
				}

				BlockModifyPlan(ctx, block, blockReq, &blockResp)

				planObject.Attrs[name] = blockResp.AttributePlan
				resp.Diagnostics.Append(blockResp.Diagnostics...)
				resp.RequiresReplace = blockResp.RequiresReplace
				resp.Private = blockResp.Private
			}

			planSet.Elems[idx] = planObject
		}

		resp.AttributePlan = planSet
	case fwschema.BlockNestingModeSingle:
		configObject, diags := coerceObjectValue(req.AttributePath, req.AttributeConfig)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		planObject, diags := coerceObjectValue(req.AttributePath, req.AttributePlan)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		stateObject, diags := coerceObjectValue(req.AttributePath, req.AttributeState)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		if planObject.Attrs == nil {
			planObject.Attrs = make(map[string]attr.Value)
		}

		for name, attr := range b.GetAttributes() {
			attrConfig, diags := objectAttributeValue(ctx, configObject, name, fwschemadata.DataDescriptionConfiguration)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			attrPlan, diags := objectAttributeValue(ctx, planObject, name, fwschemadata.DataDescriptionPlan)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			attrState, diags := objectAttributeValue(ctx, stateObject, name, fwschemadata.DataDescriptionState)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			attrReq := tfsdk.ModifyAttributePlanRequest{
				AttributeConfig: attrConfig,
				AttributePath:   req.AttributePath.AtName(name),
				AttributePlan:   attrPlan,
				AttributeState:  attrState,
				Config:          req.Config,
				Plan:            req.Plan,
				ProviderMeta:    req.ProviderMeta,
				State:           req.State,
				Private:         resp.Private,
			}
			attrResp := ModifyAttributePlanResponse{
				AttributePlan:   attrReq.AttributePlan,
				RequiresReplace: resp.RequiresReplace,
				Private:         attrReq.Private,
			}

			AttributeModifyPlan(ctx, attr, attrReq, &attrResp)

			planObject.Attrs[name] = attrResp.AttributePlan
			resp.Diagnostics.Append(attrResp.Diagnostics...)
			resp.RequiresReplace = attrResp.RequiresReplace
			resp.Private = attrResp.Private
		}

		for name, block := range b.GetBlocks() {
			attrConfig, diags := objectAttributeValue(ctx, configObject, name, fwschemadata.DataDescriptionConfiguration)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			attrPlan, diags := objectAttributeValue(ctx, planObject, name, fwschemadata.DataDescriptionPlan)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			attrState, diags := objectAttributeValue(ctx, stateObject, name, fwschemadata.DataDescriptionState)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			blockReq := tfsdk.ModifyAttributePlanRequest{
				AttributeConfig: attrConfig,
				AttributePath:   req.AttributePath.AtName(name),
				AttributePlan:   attrPlan,
				AttributeState:  attrState,
				Config:          req.Config,
				Plan:            req.Plan,
				ProviderMeta:    req.ProviderMeta,
				State:           req.State,
				Private:         resp.Private,
			}
			blockResp := ModifyAttributePlanResponse{
				AttributePlan:   blockReq.AttributePlan,
				RequiresReplace: resp.RequiresReplace,
				Private:         blockReq.Private,
			}

			BlockModifyPlan(ctx, block, blockReq, &blockResp)

			planObject.Attrs[name] = blockResp.AttributePlan
			resp.Diagnostics.Append(blockResp.Diagnostics...)
			resp.RequiresReplace = blockResp.RequiresReplace
			resp.Private = blockResp.Private
		}

		resp.AttributePlan = planObject
	default:
		err := fmt.Errorf("unknown block plan modification nesting mode (%T: %v) at path: %s", nm, nm, req.AttributePath)
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Block Plan Modification Error",
			"Block plan modification cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
		)

		return
	}
}
