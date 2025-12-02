package basetypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var (
	_ Float64Valuable = Float64Value{}
)

// Float64Valuable extends attr.Value for float64 value types.
// Implement this interface to create a custom Float64 value type.
type Float64Valuable interface {
	attr.Value

	// ToFloat64Value should convert the value type to a Float64.
	ToFloat64Value(ctx context.Context) (Float64Value, diag.Diagnostics)
}

// Float64Null creates a Float64 with a null value. Determine whether the value is
// null via the Float64 type IsNull method.
func NewFloat64Null() Float64Value {
	return Float64Value{
		state: attr.ValueStateNull,
	}
}

// Float64Unknown creates a Float64 with an unknown value. Determine whether the
// value is unknown via the Float64 type IsUnknown method.
//
// Setting the deprecated Float64 type Null, Unknown, or Value fields after
// creating a Float64 with this function has no effect.
func NewFloat64Unknown() Float64Value {
	return Float64Value{
		state: attr.ValueStateUnknown,
	}
}

// Float64Value creates a Float64 with a known value. Access the value via the Float64
// type ValueFloat64 method.
//
// Setting the deprecated Float64 type Null, Unknown, or Value fields after
// creating a Float64 with this function has no effect.
func NewFloat64Value(value float64) Float64Value {
	return Float64Value{
		state: attr.ValueStateKnown,
		value: value,
	}
}

// Float64Value represents a 64-bit floating point value, exposed as a float64.
type Float64Value struct {
	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState

	// value contains the known value, if not null or unknown.
	value float64
}

// Equal returns true if `other` is a Float64 and has the same value as `f`.
func (f Float64Value) Equal(other attr.Value) bool {
	o, ok := other.(Float64Value)

	if !ok {
		return false
	}

	if f.state != o.state {
		return false
	}

	if f.state != attr.ValueStateKnown {
		return true
	}

	return f.value == o.value
}

// ToTerraformValue returns the data contained in the Float64 as a tftypes.Value.
func (f Float64Value) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	switch f.state {
	case attr.ValueStateKnown:
		if err := tftypes.ValidateValue(tftypes.Number, f.value); err != nil {
			return tftypes.NewValue(tftypes.Number, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(tftypes.Number, f.value), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(tftypes.Number, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(tftypes.Number, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Float64 state in ToTerraformValue: %s", f.state))
	}
}

// Type returns a Float64Type.
func (f Float64Value) Type(ctx context.Context) attr.Type {
	return Float64Type{}
}

// IsNull returns true if the Float64 represents a null value.
func (f Float64Value) IsNull() bool {
	return f.state == attr.ValueStateNull
}

// IsUnknown returns true if the Float64 represents a currently unknown value.
func (f Float64Value) IsUnknown() bool {
	return f.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the Float64 value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (f Float64Value) String() string {
	if f.IsUnknown() {
		return attr.UnknownValueString
	}

	if f.IsNull() {
		return attr.NullValueString
	}

	return fmt.Sprintf("%f", f.value)
}

// ValueFloat64 returns the known float64 value. If Float64 is null or unknown, returns
// 0.0.
func (f Float64Value) ValueFloat64() float64 {
	return f.value
}

// ToFloat64Value returns Float64.
func (f Float64Value) ToFloat64Value(context.Context) (Float64Value, diag.Diagnostics) {
	return f, nil
}
