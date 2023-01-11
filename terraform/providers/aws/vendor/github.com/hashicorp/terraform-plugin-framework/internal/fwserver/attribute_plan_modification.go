package fwserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type ModifyAttributePlanResponse struct {
	AttributePlan   attr.Value
	Diagnostics     diag.Diagnostics
	RequiresReplace path.Paths
	Private         *privatestate.ProviderData
}

// AttributeModifyPlan runs all AttributePlanModifiers
//
// TODO: Clean up this abstraction back into an internal Attribute type method.
// The extra Attribute parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
func AttributeModifyPlan(ctx context.Context, a fwschema.Attribute, req tfsdk.ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	ctx = logging.FrameworkWithAttributePath(ctx, req.AttributePath.String())

	var requiresReplace bool

	privateProviderData := privatestate.EmptyProviderData(ctx)

	if req.Private != nil {
		resp.Private = req.Private
		privateProviderData = req.Private
	}

	if attributeWithPlanModifiers, ok := a.(fwxschema.AttributeWithPlanModifiers); ok {
		for _, planModifier := range attributeWithPlanModifiers.GetPlanModifiers() {
			modifyResp := &tfsdk.ModifyAttributePlanResponse{
				AttributePlan:   req.AttributePlan,
				RequiresReplace: requiresReplace,
				Private:         privateProviderData,
			}

			logging.FrameworkDebug(
				ctx,
				"Calling provider defined AttributePlanModifier",
				map[string]interface{}{
					logging.KeyDescription: planModifier.Description(ctx),
				},
			)
			planModifier.Modify(ctx, req, modifyResp)
			logging.FrameworkDebug(
				ctx,
				"Called provider defined AttributePlanModifier",
				map[string]interface{}{
					logging.KeyDescription: planModifier.Description(ctx),
				},
			)

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

	if resp.Diagnostics.HasError() {
		return
	}

	if a.GetAttributes() == nil || len(a.GetAttributes().GetAttributes()) == 0 {
		return
	}

	nm := a.GetAttributes().GetNestingMode()
	switch nm {
	case fwschema.NestingModeList:
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

			for name, attr := range a.GetAttributes().GetAttributes() {
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

			planList.Elems[idx] = planObject
		}

		resp.AttributePlan = planList
	case fwschema.NestingModeSet:
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

			for name, attr := range a.GetAttributes().GetAttributes() {
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

			planSet.Elems[idx] = planObject
		}

		resp.AttributePlan = planSet
	case fwschema.NestingModeMap:
		configMap, diags := coerceMapValue(req.AttributePath, req.AttributeConfig)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		planMap, diags := coerceMapValue(req.AttributePath, req.AttributePlan)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		stateMap, diags := coerceMapValue(req.AttributePath, req.AttributeState)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		for key, planElem := range planMap.Elems {
			attrPath := req.AttributePath.AtMapKey(key)

			configObject, diags := mapElemObject(ctx, attrPath, configMap, key, fwschemadata.DataDescriptionConfiguration)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			planObject, diags := coerceObjectValue(attrPath, planElem)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			stateObject, diags := mapElemObject(ctx, attrPath, stateMap, key, fwschemadata.DataDescriptionState)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			for name, attr := range a.GetAttributes().GetAttributes() {
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

			planMap.Elems[key] = planObject
		}

		resp.AttributePlan = planMap
	case fwschema.NestingModeSingle:
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

		if len(planObject.Attrs) == 0 {
			return
		}

		for name, attr := range a.GetAttributes().GetAttributes() {
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

		resp.AttributePlan = planObject
	default:
		err := fmt.Errorf("unknown attribute nesting mode (%T: %v) at path: %s", nm, nm, req.AttributePath)
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Attribute Plan Modification Error",
			"Attribute plan modifier cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
		)

		return
	}
}

func attributePlanModificationValueError(ctx context.Context, value attr.Value, description fwschemadata.DataDescription, err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Attribute Plan Modification "+description.Title()+" Value Error",
		"An unexpected error occurred while fetching a "+value.Type(ctx).String()+" element value in the "+description.String()+". "+
			"This is an issue with the provider and should be reported to the provider developers.\n\n"+
			"Original Error: "+err.Error(),
	)
}

func attributePlanModificationWalkError(schemaPath path.Path, value attr.Value) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		schemaPath,
		"Attribute Plan Modification Walk Error",
		"An unexpected error occurred while walking the schema for attribute plan modification. "+
			"This is an issue with terraform-plugin-framework and should be reported to the provider developers.\n\n"+
			fmt.Sprintf("unknown attribute value type (%T) at path: %s", value, schemaPath),
	)
}
