// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package reflect

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Number creates a *big.Float and populates it with the data in `val`. It then
// gets converted to the type of `target`, as long as `target` is a valid
// number type (any of the built-in int, uint, or float types, *big.Float, and
// *big.Int).
//
// Number will loudly fail when a number cannot be losslessly represented using
// the requested type, unless opts.AllowRoundingNumbers is set to true. This
// setting is mildly dangerous, because Terraform does not like when you round
// things, as a general rule of thumb.
//
// It is meant to be called through Into, not directly.
func Number(ctx context.Context, typ attr.Type, val tftypes.Value, target reflect.Value, opts Options, path path.Path) (reflect.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := big.NewFloat(0)
	err := val.As(&result)
	if err != nil {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Err:        err,
			TargetType: target.Type(),
			Val:        val,
		}))
		return target, diags
	}
	roundingError := fmt.Errorf("cannot store %s in %s", result.String(), target.Type())
	roundingErrorDiag := diag.NewAttributeErrorDiagnostic(
		path,
		"Value Conversion Error",
		"An unexpected error was encountered trying to convert to number. This is always an error in the provider. Please report the following to the provider developer:\n\n"+roundingError.Error(),
	)

	switch target.Type() {
	case reflect.TypeOf(big.NewFloat(0)):
		return reflect.ValueOf(result), diags
	case reflect.TypeOf(big.NewInt(0)):
		intResult, acc := result.Int(nil)
		if acc != big.Exact && !opts.AllowRoundingNumbers {
			return reflect.ValueOf(result), append(diags, roundingErrorDiag)
		}
		return reflect.ValueOf(intResult), diags
	}
	switch target.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		intResult, acc := result.Int64()
		if acc != big.Exact && !opts.AllowRoundingNumbers {
			return target, append(diags, roundingErrorDiag)
		}
		switch target.Kind() {
		case reflect.Int:
			if strconv.IntSize == 32 && intResult > math.MaxInt32 {
				if !opts.AllowRoundingNumbers {
					return target, append(diags, roundingErrorDiag)
				}
				intResult = math.MaxInt32
			}
			if strconv.IntSize == 32 && intResult < math.MinInt32 {
				if !opts.AllowRoundingNumbers {
					return target, append(diags, roundingErrorDiag)
				}
				intResult = math.MinInt32
			}
			return reflect.ValueOf(int(intResult)), diags
		case reflect.Int8:
			if intResult > math.MaxInt8 {
				if !opts.AllowRoundingNumbers {
					return target, append(diags, roundingErrorDiag)
				}
				intResult = math.MaxInt8
			}
			if intResult < math.MinInt8 {
				if !opts.AllowRoundingNumbers {
					return target, append(diags, roundingErrorDiag)
				}
				intResult = math.MinInt8
			}
			return reflect.ValueOf(int8(intResult)), diags
		case reflect.Int16:
			if intResult > math.MaxInt16 {
				if !opts.AllowRoundingNumbers {
					return target, append(diags, roundingErrorDiag)
				}
				intResult = math.MaxInt16
			}
			if intResult < math.MinInt16 {
				if !opts.AllowRoundingNumbers {
					return target, append(diags, roundingErrorDiag)
				}
				intResult = math.MinInt16
			}
			return reflect.ValueOf(int16(intResult)), diags
		case reflect.Int32:
			if intResult > math.MaxInt32 {
				if !opts.AllowRoundingNumbers {
					return target, append(diags, roundingErrorDiag)
				}
				intResult = math.MaxInt32
			}
			if intResult < math.MinInt32 {
				if !opts.AllowRoundingNumbers {
					return target, append(diags, roundingErrorDiag)
				}
				intResult = math.MinInt32
			}
			return reflect.ValueOf(int32(intResult)), diags
		case reflect.Int64:
			return reflect.ValueOf(intResult), diags
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		uintResult, acc := result.Uint64()
		if acc != big.Exact && !opts.AllowRoundingNumbers {
			return target, append(diags, roundingErrorDiag)
		}
		switch target.Kind() {
		case reflect.Uint:
			if strconv.IntSize == 32 && uintResult > math.MaxUint32 {
				if !opts.AllowRoundingNumbers {
					return target, append(diags, roundingErrorDiag)
				}
				uintResult = math.MaxUint32
			}
			return reflect.ValueOf(uint(uintResult)), diags
		case reflect.Uint8:
			if uintResult > math.MaxUint8 {
				if !opts.AllowRoundingNumbers {
					return target, append(diags, roundingErrorDiag)
				}
				uintResult = math.MaxUint8
			}
			return reflect.ValueOf(uint8(uintResult)), diags
		case reflect.Uint16:
			if uintResult > math.MaxUint16 {
				if !opts.AllowRoundingNumbers {
					return target, append(diags, roundingErrorDiag)
				}
				uintResult = math.MaxUint16
			}
			return reflect.ValueOf(uint16(uintResult)), diags
		case reflect.Uint32:
			if uintResult > math.MaxUint32 {
				if !opts.AllowRoundingNumbers {
					return target, append(diags, roundingErrorDiag)
				}
				uintResult = math.MaxUint32
			}
			return reflect.ValueOf(uint32(uintResult)), diags
		case reflect.Uint64:
			return reflect.ValueOf(uintResult), diags
		}
	case reflect.Float32:
		floatResult, acc := result.Float32()
		if acc != big.Exact && !opts.AllowRoundingNumbers {
			return target, append(diags, roundingErrorDiag)
		} else if acc == big.Above {
			floatResult = math.MaxFloat32
		} else if acc == big.Below {
			floatResult = math.SmallestNonzeroFloat32
		} else if acc != big.Exact {
			err := fmt.Errorf("unsure how to round %s and %f", acc, floatResult)
			diags.AddAttributeError(
				path,
				"Value Conversion Error",
				"An unexpected error was encountered trying to convert to number. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			)
			return target, diags
		}
		return reflect.ValueOf(floatResult), diags
	case reflect.Float64:
		floatResult, acc := result.Float64()
		if acc != big.Exact && !opts.AllowRoundingNumbers {
			return target, append(diags, roundingErrorDiag)
		}
		if acc == big.Above {
			if floatResult == math.Inf(1) || floatResult == math.MaxFloat64 {
				floatResult = math.MaxFloat64
			} else if floatResult == 0.0 || floatResult == math.SmallestNonzeroFloat64 {
				floatResult = -math.SmallestNonzeroFloat64
			} else {
				err := fmt.Errorf("not sure how to round %s and %f", acc, floatResult)
				diags.AddAttributeError(
					path,
					"Value Conversion Error",
					"An unexpected error was encountered trying to convert to number. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
				)
				return target, diags
			}
		} else if acc == big.Below {
			if floatResult == math.Inf(-1) || floatResult == -math.MaxFloat64 {
				floatResult = -math.MaxFloat64
			} else if floatResult == -0.0 || floatResult == -math.SmallestNonzeroFloat64 { //nolint:staticcheck
				floatResult = math.SmallestNonzeroFloat64
			} else {
				err := fmt.Errorf("not sure how to round %s and %f", acc, floatResult)
				diags.AddAttributeError(
					path,
					"Value Conversion Error",
					"An unexpected error was encountered trying to convert to number. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
				)
				return target, diags
			}
		} else if acc != big.Exact {
			err := fmt.Errorf("not sure how to round %s and %f", acc, floatResult)
			diags.AddAttributeError(
				path,
				"Value Conversion Error",
				"An unexpected error was encountered trying to convert to number. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			)
			return target, diags
		}
		return reflect.ValueOf(floatResult), diags
	}
	err = fmt.Errorf("cannot convert number to %s", target.Type())
	diags.AddAttributeError(
		path,
		"Value Conversion Error",
		"An unexpected error was encountered trying to convert to number. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
	)
	return target, diags
}

// FromInt creates an attr.Value using `typ` from an int64.
//
// It is meant to be called through FromValue, not directly.
func FromInt(ctx context.Context, typ attr.Type, val int64, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	err := tftypes.ValidateValue(tftypes.Number, val)
	if err != nil {
		return nil, append(diags, validateValueErrorDiag(err, path))
	}
	tfNum := tftypes.NewValue(tftypes.Number, val)

	if typeWithValidate, ok := typ.(xattr.TypeWithValidate); ok {
		diags.Append(typeWithValidate.Validate(ctx, tfNum, path)...)

		if diags.HasError() {
			return nil, diags
		}
	}

	num, err := typ.ValueFromTerraform(ctx, tfNum)
	if err != nil {
		return nil, append(diags, valueFromTerraformErrorDiag(err, path))
	}

	return num, diags
}

// FromUint creates an attr.Value using `typ` from a uint64.
//
// It is meant to be called through FromValue, not directly.
func FromUint(ctx context.Context, typ attr.Type, val uint64, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	err := tftypes.ValidateValue(tftypes.Number, val)
	if err != nil {
		return nil, append(diags, validateValueErrorDiag(err, path))
	}
	tfNum := tftypes.NewValue(tftypes.Number, val)

	if typeWithValidate, ok := typ.(xattr.TypeWithValidate); ok {
		diags.Append(typeWithValidate.Validate(ctx, tfNum, path)...)

		if diags.HasError() {
			return nil, diags
		}
	}

	num, err := typ.ValueFromTerraform(ctx, tfNum)
	if err != nil {
		return nil, append(diags, valueFromTerraformErrorDiag(err, path))
	}

	return num, diags
}

// FromFloat creates an attr.Value using `typ` from a float64.
//
// It is meant to be called through FromValue, not directly.
func FromFloat(ctx context.Context, typ attr.Type, val float64, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	err := tftypes.ValidateValue(tftypes.Number, val)
	if err != nil {
		return nil, append(diags, validateValueErrorDiag(err, path))
	}
	tfNum := tftypes.NewValue(tftypes.Number, val)

	if typeWithValidate, ok := typ.(xattr.TypeWithValidate); ok {
		diags.Append(typeWithValidate.Validate(ctx, tfNum, path)...)

		if diags.HasError() {
			return nil, diags
		}
	}

	num, err := typ.ValueFromTerraform(ctx, tfNum)
	if err != nil {
		return nil, append(diags, valueFromTerraformErrorDiag(err, path))
	}

	return num, diags
}

// FromBigFloat creates an attr.Value using `typ` from a *big.Float.
//
// It is meant to be called through FromValue, not directly.
func FromBigFloat(ctx context.Context, typ attr.Type, val *big.Float, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	err := tftypes.ValidateValue(tftypes.Number, val)
	if err != nil {
		return nil, append(diags, validateValueErrorDiag(err, path))
	}
	tfNum := tftypes.NewValue(tftypes.Number, val)

	if typeWithValidate, ok := typ.(xattr.TypeWithValidate); ok {
		diags.Append(typeWithValidate.Validate(ctx, tfNum, path)...)

		if diags.HasError() {
			return nil, diags
		}
	}

	num, err := typ.ValueFromTerraform(ctx, tfNum)
	if err != nil {
		return nil, append(diags, valueFromTerraformErrorDiag(err, path))
	}

	return num, diags
}

// FromBigInt creates an attr.Value using `typ` from a *big.Int.
//
// It is meant to be called through FromValue, not directly.
func FromBigInt(ctx context.Context, typ attr.Type, val *big.Int, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	fl := big.NewFloat(0).SetInt(val)
	err := tftypes.ValidateValue(tftypes.Number, fl)
	if err != nil {
		return nil, append(diags, validateValueErrorDiag(err, path))
	}
	tfNum := tftypes.NewValue(tftypes.Number, fl)

	if typeWithValidate, ok := typ.(xattr.TypeWithValidate); ok {
		diags.Append(typeWithValidate.Validate(ctx, tfNum, path)...)

		if diags.HasError() {
			return nil, diags
		}
	}

	num, err := typ.ValueFromTerraform(ctx, tfNum)
	if err != nil {
		return nil, append(diags, valueFromTerraformErrorDiag(err, path))
	}

	return num, diags
}
