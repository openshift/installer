package types

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Value = Bool{}
)

func boolValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.IsNull() {
		return Bool{
			Null: true,
		}, nil
	}
	if !in.IsKnown() {
		return Bool{
			Unknown: true,
		}, nil
	}
	var b bool
	err := in.As(&b)
	if err != nil {
		return nil, err
	}
	return Bool{Value: b}, nil
}

// Bool represents a boolean value.
type Bool struct {
	// Unknown will be true if the value is not yet known.
	Unknown bool

	// Null will be true if the value was not set, or was explicitly set to
	// null.
	Null bool

	// Value contains the set value, as long as Unknown and Null are both
	// false.
	Value bool
}

// Type returns a BoolType.
func (b Bool) Type(_ context.Context) attr.Type {
	return BoolType
}

// ToTerraformValue returns the data contained in the Bool as a tftypes.Value.
func (b Bool) ToTerraformValue(_ context.Context) (tftypes.Value, error) {
	if b.Null {
		return tftypes.NewValue(tftypes.Bool, nil), nil
	}
	if b.Unknown {
		return tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue), nil
	}
	if err := tftypes.ValidateValue(tftypes.Bool, b.Value); err != nil {
		return tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue), err
	}
	return tftypes.NewValue(tftypes.Bool, b.Value), nil
}

// Equal returns true if `other` is a *Bool and has the same value as `b`.
func (b Bool) Equal(other attr.Value) bool {
	o, ok := other.(Bool)
	if !ok {
		return false
	}
	if b.Unknown != o.Unknown {
		return false
	}
	if b.Null != o.Null {
		return false
	}
	return b.Value == o.Value
}

// IsNull returns true if the Bool represents a null value.
func (b Bool) IsNull() bool {
	return b.Null
}

// IsUnknown returns true if the Bool represents a currently unknown value.
func (b Bool) IsUnknown() bool {
	return b.Unknown
}

// String returns a human-readable representation of the Bool value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (b Bool) String() string {
	if b.Unknown {
		return attr.UnknownValueString
	}

	if b.Null {
		return attr.NullValueString
	}

	return fmt.Sprintf("%t", b.Value)
}
