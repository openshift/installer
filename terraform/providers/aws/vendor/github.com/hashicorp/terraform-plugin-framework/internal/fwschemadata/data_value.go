package fwschemadata

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/totftypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ValueAtPath retrieves the attribute found at `path` and returns it as an
// attr.Value. Consumers should assert the type of the returned value with the
// desired attr.Type.
func (d Data) ValueAtPath(ctx context.Context, schemaPath path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	tftypesPath, tftypesPathDiags := totftypes.AttributePath(ctx, schemaPath)

	diags.Append(tftypesPathDiags...)

	if diags.HasError() {
		return nil, diags
	}

	attrType, err := d.Schema.TypeAtTerraformPath(ctx, tftypesPath)

	if err != nil {
		diags.AddAttributeError(
			schemaPath,
			d.Description.Title()+" Read Error",
			"An unexpected error was encountered trying to retrieve type information at a given path. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Error: "+err.Error(),
		)
		return nil, diags
	}

	// if the data is null, return a null value of the type
	if d.TerraformValue.IsNull() {
		attrValue, err := attrType.ValueFromTerraform(ctx, tftypes.NewValue(attrType.TerraformType(ctx), nil))

		if err != nil {
			diags.AddAttributeError(
				schemaPath,
				d.Description.Title()+" Read Error",
				"An unexpected error was encountered trying to create a null attribute value from the given path. "+
					"Please report the following to the provider developer:\n\n"+
					"Type: "+attrType.String()+"\n"+
					"Error:"+err.Error(),
			)
		}

		return attrValue, diags
	}

	tfValue, err := d.TerraformValueAtTerraformPath(ctx, tftypesPath)

	// Ignoring ErrInvalidStep will allow this method to return a null value of the type.
	if err != nil && !errors.Is(err, tftypes.ErrInvalidStep) {
		diags.AddAttributeError(
			schemaPath,
			d.Description.Title()+" Read Error",
			"An unexpected error was encountered trying to retrieve an attribute value from the given path. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
	}

	// TODO: If ErrInvalidStep, check parent paths for unknown value.
	//       If found, convert this value to an unknown value.
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/186

	if attrTypeWithValidate, ok := attrType.(xattr.TypeWithValidate); ok {
		logging.FrameworkTrace(ctx, "Type implements TypeWithValidate")
		logging.FrameworkDebug(ctx, "Calling provider defined Type Validate")
		diags.Append(attrTypeWithValidate.Validate(ctx, tfValue, schemaPath)...)
		logging.FrameworkDebug(ctx, "Called provider defined Type Validate")

		if diags.HasError() {
			return nil, diags
		}
	}

	attrValue, err := attrType.ValueFromTerraform(ctx, tfValue)

	if err != nil {
		diags.AddAttributeError(
			schemaPath,
			d.Description.Title()+" Read Error",
			"An unexpected error was encountered trying to convert an attribute value from the "+d.Description.String()+". This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Error: "+err.Error(),
		)
		return nil, diags
	}

	return attrValue, diags
}
