package types

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Value = Float64{}
)

func float64Validate(_ context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Equal(tftypes.Number) {
		diags.AddAttributeError(
			path,
			"Float64 Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Expected Number value, received %T with value: %v", in, in),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var value *big.Float
	err := in.As(&value)

	if err != nil {
		diags.AddAttributeError(
			path,
			"Float64 Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Cannot convert value to big.Float: %s", err),
		)
		return diags
	}

	_, accuracy := value.Float64()

	if accuracy != 0 {
		diags.AddAttributeError(
			path,
			"Float64 Type Validation Error",
			fmt.Sprintf("Value %s cannot be represented as a 64-bit floating point.", value),
		)
		return diags
	}

	return diags
}

func float64ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return Float64{Unknown: true}, nil
	}

	if in.IsNull() {
		return Float64{Null: true}, nil
	}

	var bigF *big.Float
	err := in.As(&bigF)

	if err != nil {
		return nil, err
	}

	f, accuracy := bigF.Float64()

	if accuracy != 0 {
		return nil, fmt.Errorf("Value %s cannot be represented as a 64-bit floating point.", bigF)
	}

	return Float64{Value: f}, nil
}

// Float64 represents a 64-bit floating point value, exposed as a float64.
type Float64 struct {
	// Unknown will be true if the value is not yet known.
	Unknown bool

	// Null will be true if the value was not set, or was explicitly set to
	// null.
	Null bool

	// Value contains the set value, as long as Unknown and Null are both
	// false.
	Value float64
}

// Equal returns true if `other` is a Float64 and has the same value as `f`.
func (f Float64) Equal(other attr.Value) bool {
	o, ok := other.(Float64)

	if !ok {
		return false
	}

	if f.Unknown != o.Unknown {
		return false
	}

	if f.Null != o.Null {
		return false
	}

	return f.Value == o.Value
}

// ToTerraformValue returns the data contained in the Float64 as a tftypes.Value.
func (f Float64) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if f.Null {
		return tftypes.NewValue(tftypes.Number, nil), nil
	}

	if f.Unknown {
		return tftypes.NewValue(tftypes.Number, tftypes.UnknownValue), nil
	}

	bf := big.NewFloat(f.Value)
	if err := tftypes.ValidateValue(tftypes.Number, bf); err != nil {
		return tftypes.NewValue(tftypes.Number, tftypes.UnknownValue), err
	}
	return tftypes.NewValue(tftypes.Number, bf), nil
}

// Type returns a Float64Type.
func (f Float64) Type(ctx context.Context) attr.Type {
	return Float64Type
}

// IsNull returns true if the Float64 represents a null value.
func (f Float64) IsNull() bool {
	return f.Null
}

// IsUnknown returns true if the Float64 represents a currently unknown value.
func (f Float64) IsUnknown() bool {
	return f.Unknown
}

// String returns a human-readable representation of the Float64 value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (f Float64) String() string {
	if f.Unknown {
		return attr.UnknownValueString
	}

	if f.Null {
		return attr.NullValueString
	}

	return fmt.Sprintf("%f", f.Value)
}
