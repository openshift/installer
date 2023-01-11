package fwtypes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type durationType uint8

const (
	DurationType durationType = iota
)

var (
	_ xattr.TypeWithValidate = DurationType
)

func (d durationType) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.String
}

func (d durationType) ValueFromTerraform(_ context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return Duration{Unknown: true}, nil
	}
	if in.IsNull() {
		return Duration{Null: true}, nil
	}
	var s string
	err := in.As(&s)
	if err != nil {
		return nil, err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return nil, err
	}
	return Duration{Value: dur}, nil
}

func (d durationType) ValueType(context.Context) attr.Value {
	return Duration{}
}

// Equal returns true if `o` is also a DurationType.
func (d durationType) Equal(o attr.Type) bool {
	_, ok := o.(durationType)
	return ok
}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// type.
func (d durationType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, d.String())
}

// String returns a human-friendly description of the DurationType.
func (d durationType) String() string {
	return "types.DurationType"
}

// Validate implements type validation.
func (d durationType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if !in.Type().Is(tftypes.String) {
		diags.AddAttributeError(
			path,
			"Duration Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Expected String value, received %T with value: %v", in, in),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var value string
	err := in.As(&value)
	if err != nil {
		diags.AddAttributeError(
			path,
			"Duration Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Cannot convert value to time.Duration: %s", err),
		)
		return diags
	}

	_, err = time.ParseDuration(value)
	if err != nil {
		diags.AddAttributeError(
			path,
			"Duration Type Validation Error",
			fmt.Sprintf("Value %q cannot be parsed as a Duration.", value),
		)
		return diags
	}

	return diags
}

func (d durationType) Description() string {
	return `A sequence of numbers with a unit suffix, "h" for hour, "m" for minute, and "s" for second.`
}

type Duration struct {
	// Unknown will be true if the value is not yet known.
	Unknown bool

	// Null will be true if the value was not set, or was explicitly set to
	// null.
	Null bool

	// Value contains the set value, as long as Unknown and Null are both
	// false.
	Value time.Duration
}

// Type returns a DurationType.
func (d Duration) Type(_ context.Context) attr.Type {
	return DurationType
}

// ToTerraformValue returns the data contained in the *String as a string. If
// Unknown is true, it returns a tftypes.UnknownValue. If Null is true, it
// returns nil.
func (d Duration) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	t := DurationType.TerraformType(ctx)
	if d.Null {
		return tftypes.NewValue(t, nil), nil
	}
	if d.Unknown {
		return tftypes.NewValue(t, tftypes.UnknownValue), nil
	}
	if err := tftypes.ValidateValue(tftypes.Number, d.Value); err != nil {
		return tftypes.NewValue(t, tftypes.UnknownValue), err
	}
	return tftypes.NewValue(t, d.Value), nil
}

// Equal returns true if `other` is a *Duration and has the same value as `d`.
func (d Duration) Equal(other attr.Value) bool {
	o, ok := other.(Duration)
	if !ok {
		return false
	}
	if d.Unknown != o.Unknown {
		return false
	}
	if d.Null != o.Null {
		return false
	}
	return d.Value == o.Value
}

// IsNull returns true if the Value is not set, or is explicitly set to null.
func (d Duration) IsNull() bool {
	return d.Null
}

// IsUnknown returns true if the Value is not yet known.
func (d Duration) IsUnknown() bool {
	return d.Unknown
}

// String returns a summary representation of either the underlying Value,
// or UnknownValueString (`<unknown>`) when IsUnknown() returns true,
// or NullValueString (`<null>`) when IsNull() return true.
//
// This is an intentionally lossy representation, that are best suited for
// logging and error reporting, as they are not protected by
// compatibility guarantees within the framework.
func (d Duration) String() string {
	if d.IsUnknown() {
		return attr.UnknownValueString
	}

	if d.IsNull() {
		return attr.NullValueString
	}

	return d.Value.String()
}
