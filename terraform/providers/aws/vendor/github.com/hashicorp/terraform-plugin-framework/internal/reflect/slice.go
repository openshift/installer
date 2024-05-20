// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package reflect

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// build a slice of elements, matching the type of `target`, and fill it with
// the data in `val`.
func reflectSlice(ctx context.Context, typ attr.Type, val tftypes.Value, target reflect.Value, opts Options, path path.Path) (reflect.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	// this only works with slices, so check that out first
	if target.Kind() != reflect.Slice {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        val,
			TargetType: target.Type(),
			Err:        fmt.Errorf("expected a slice type, got %s", target.Type()),
		}))
		return target, diags
	}
	// TODO: check that the val is a list or set or tuple
	elemTyper, ok := typ.(attr.TypeWithElementType)
	if !ok {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        val,
			TargetType: target.Type(),
			Err:        fmt.Errorf("cannot reflect %s using type information provided by %T, %T must be an attr.TypeWithElementType", val.Type(), typ, typ),
		}))
		return target, diags
	}

	// we need our value to become a list of values so we can iterate over
	// them and handle them individually
	var values []tftypes.Value
	err := val.As(&values)
	if err != nil {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        val,
			TargetType: target.Type(),
			Err:        err,
		}))
		return target, diags
	}

	// we need to know the type the slice is wrapping
	elemType := target.Type().Elem()
	elemAttrType := elemTyper.ElementType()

	// we want an empty version of the slice
	slice := reflect.MakeSlice(target.Type(), 0, len(values))

	// go over each of the values passed in, create a Go value of the right
	// type for them, and add it to our new slice
	for pos, value := range values {
		// create a new Go value of the type that can go in the slice
		targetValue := reflect.Zero(elemType)

		// update our path so we can have nice errors
		valPath := path.AtListIndex(pos)

		if typ.TerraformType(ctx).Is(tftypes.Set{}) {
			attrVal, err := elemAttrType.ValueFromTerraform(ctx, value)

			if err != nil {
				diags.AddAttributeError(
					path,
					"Value Conversion Error",
					"An unexpected error was encountered trying to convert to slice value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
				)
				return target, diags
			}

			valPath = path.AtSetValue(attrVal)
		}

		// reflect the value into our new target
		val, valDiags := BuildValue(ctx, elemAttrType, value, targetValue, opts, valPath)
		diags.Append(valDiags...)

		if diags.HasError() {
			return target, diags
		}

		// add the new target to our slice
		slice = reflect.Append(slice, val)
	}

	return slice, diags
}

// FromSlice returns an attr.Value as produced by `typ` using the data in
// `val`. `val` must be a slice. `typ` must be an attr.TypeWithElementType or
// attr.TypeWithElementTypes. If the slice is nil, the representation of null
// for `typ` will be returned. Otherwise, FromSlice will recurse into FromValue
// for each element in the slice, using the element type or types defined on
// `typ` to construct values for them.
//
// It is meant to be called through FromValue, not directly.
func FromSlice(ctx context.Context, typ attr.Type, val reflect.Value, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	// TODO: support tuples, which are attr.TypeWithElementTypes
	tfType := typ.TerraformType(ctx)

	if val.IsNil() {
		tfVal := tftypes.NewValue(tfType, nil)

		if typeWithValidate, ok := typ.(xattr.TypeWithValidate); ok {
			diags.Append(typeWithValidate.Validate(ctx, tfVal, path)...)

			if diags.HasError() {
				return nil, diags
			}
		}

		attrVal, err := typ.ValueFromTerraform(ctx, tfVal)

		if err != nil {
			diags.AddAttributeError(
				path,
				"Value Conversion Error",
				"An unexpected error was encountered trying to convert from slice value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			)
			return nil, diags
		}

		return attrVal, diags
	}

	t, ok := typ.(attr.TypeWithElementType)
	if !ok {
		err := fmt.Errorf("cannot use type %T as schema type %T; %T must be an attr.TypeWithElementType to hold %T", val, typ, typ, val)
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert from slice value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
	}

	elemType := t.ElementType()
	tfElems := make([]tftypes.Value, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		// The underlying reflect.Slice is fetched by Index(). For set types,
		// the path is value-based instead of index-based. Since there is only
		// the index until the value is retrieved, this will pass the
		// technically incorrect index-based path at first for framework
		// debugging purposes, then correct the path afterwards.
		valPath := path.AtListIndex(i)

		val, valDiags := FromValue(ctx, elemType, val.Index(i).Interface(), valPath)
		diags.Append(valDiags...)

		if diags.HasError() {
			return nil, diags
		}

		tfVal, err := val.ToTerraformValue(ctx)
		if err != nil {
			return nil, append(diags, toTerraformValueErrorDiag(err, path))
		}

		if tfType.Is(tftypes.Set{}) {
			valPath = path.AtSetValue(val)
		}

		if typeWithValidate, ok := elemType.(xattr.TypeWithValidate); ok {
			diags.Append(typeWithValidate.Validate(ctx, tfVal, valPath)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		tfElems = append(tfElems, tfVal)
	}

	err := tftypes.ValidateValue(tfType, tfElems)
	if err != nil {
		return nil, append(diags, validateValueErrorDiag(err, path))
	}

	tfVal := tftypes.NewValue(tfType, tfElems)

	if typeWithValidate, ok := typ.(xattr.TypeWithValidate); ok {
		diags.Append(typeWithValidate.Validate(ctx, tfVal, path)...)

		if diags.HasError() {
			return nil, diags
		}
	}

	attrVal, err := typ.ValueFromTerraform(ctx, tfVal)

	if err != nil {
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert from slice value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
	}

	return attrVal, diags
}
