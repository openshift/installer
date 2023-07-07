// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fromtftypes"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
)

// TransformDefaults walks the schema and applies schema defined default values
// when configRaw contains a null value at the same path.
func (d *Data) TransformDefaults(ctx context.Context, configRaw tftypes.Value) diag.Diagnostics {
	var diags diag.Diagnostics

	configData := Data{
		Description:    DataDescriptionConfiguration,
		Schema:         d.Schema,
		TerraformValue: configRaw,
	}

	// Errors are handled as richer diag.Diagnostics instead.
	d.TerraformValue, _ = tftypes.Transform(d.TerraformValue, func(tfTypePath *tftypes.AttributePath, tfTypeValue tftypes.Value) (tftypes.Value, error) {
		fwPath, fwPathDiags := fromtftypes.AttributePath(ctx, tfTypePath, d.Schema)

		diags.Append(fwPathDiags...)

		// Do not transform if path cannot be converted.
		// Checking against fwPathDiags will capture all errors.
		if fwPathDiags.HasError() {
			return tfTypeValue, nil
		}

		configValue, configValueDiags := configData.ValueAtPath(ctx, fwPath)

		diags.Append(configValueDiags...)

		// Do not transform if rawConfig value cannot be retrieved.
		if configValueDiags.HasError() {
			return tfTypeValue, nil
		}

		// Do not transform if rawConfig value is not null.
		if !configValue.IsNull() {
			return tfTypeValue, nil
		}

		attrAtPath, err := d.Schema.AttributeAtTerraformPath(ctx, tfTypePath)

		if err != nil {
			if errors.Is(err, fwschema.ErrPathInsideAtomicAttribute) {
				// ignore attributes/elements inside schema.Attributes, they have no schema of their own
				logging.FrameworkTrace(ctx, "attribute is a non-schema attribute, not setting default")
				return tfTypeValue, nil
			}

			if errors.Is(err, fwschema.ErrPathIsBlock) {
				// ignore blocks, they do not have a computed field
				logging.FrameworkTrace(ctx, "attribute is a block, not setting default")
				return tfTypeValue, nil
			}

			return tftypes.Value{}, fmt.Errorf("couldn't find attribute in resource schema: %w", err)
		}

		switch a := attrAtPath.(type) {
		case fwschema.AttributeWithBoolDefaultValue:
			defaultValue := a.BoolDefaultValue()
			if defaultValue != nil {
				resp := defaults.BoolResponse{}
				defaultValue.DefaultBool(ctx, defaults.BoolRequest{}, &resp)
				logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath.String(), resp.PlanValue.String()))
				return resp.PlanValue.ToTerraformValue(ctx)
			}
		case fwschema.AttributeWithFloat64DefaultValue:
			defaultValue := a.Float64DefaultValue()
			if defaultValue != nil {
				resp := defaults.Float64Response{}
				defaultValue.DefaultFloat64(ctx, defaults.Float64Request{}, &resp)
				logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath.String(), resp.PlanValue.String()))
				return resp.PlanValue.ToTerraformValue(ctx)
			}
		case fwschema.AttributeWithInt64DefaultValue:
			defaultValue := a.Int64DefaultValue()
			if defaultValue != nil {
				resp := defaults.Int64Response{}
				defaultValue.DefaultInt64(ctx, defaults.Int64Request{}, &resp)
				logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath.String(), resp.PlanValue.String()))
				return resp.PlanValue.ToTerraformValue(ctx)
			}
		case fwschema.AttributeWithListDefaultValue:
			defaultValue := a.ListDefaultValue()
			if defaultValue != nil {
				resp := defaults.ListResponse{}
				defaultValue.DefaultList(ctx, defaults.ListRequest{}, &resp)
				logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath.String(), resp.PlanValue.String()))
				return resp.PlanValue.ToTerraformValue(ctx)
			}
		case fwschema.AttributeWithMapDefaultValue:
			defaultValue := a.MapDefaultValue()
			if defaultValue != nil {
				resp := defaults.MapResponse{}
				defaultValue.DefaultMap(ctx, defaults.MapRequest{}, &resp)
				logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath.String(), resp.PlanValue.String()))
				return resp.PlanValue.ToTerraformValue(ctx)
			}
		case fwschema.AttributeWithNumberDefaultValue:
			defaultValue := a.NumberDefaultValue()
			if defaultValue != nil {
				resp := defaults.NumberResponse{}
				defaultValue.DefaultNumber(ctx, defaults.NumberRequest{}, &resp)
				logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath.String(), resp.PlanValue.String()))
				return resp.PlanValue.ToTerraformValue(ctx)
			}
		case fwschema.AttributeWithObjectDefaultValue:
			defaultValue := a.ObjectDefaultValue()
			if defaultValue != nil {
				resp := defaults.ObjectResponse{}
				defaultValue.DefaultObject(ctx, defaults.ObjectRequest{}, &resp)
				logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath.String(), resp.PlanValue.String()))
				return resp.PlanValue.ToTerraformValue(ctx)
			}
		case fwschema.AttributeWithSetDefaultValue:
			defaultValue := a.SetDefaultValue()
			if defaultValue != nil {
				resp := defaults.SetResponse{}
				defaultValue.DefaultSet(ctx, defaults.SetRequest{}, &resp)
				logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath.String(), resp.PlanValue.String()))
				return resp.PlanValue.ToTerraformValue(ctx)
			}
		case fwschema.AttributeWithStringDefaultValue:
			defaultValue := a.StringDefaultValue()
			if defaultValue != nil {
				resp := defaults.StringResponse{}
				defaultValue.DefaultString(ctx, defaults.StringRequest{}, &resp)
				logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath.String(), resp.PlanValue.String()))
				return resp.PlanValue.ToTerraformValue(ctx)
			}
		}

		return tfTypeValue, nil
	})

	return diags
}
