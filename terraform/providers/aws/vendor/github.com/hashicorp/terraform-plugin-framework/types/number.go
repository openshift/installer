package types

import (
	"context"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Value = Number{}
)

func numberValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return Number{Unknown: true}, nil
	}
	if in.IsNull() {
		return Number{Null: true}, nil
	}
	n := big.NewFloat(0)
	err := in.As(&n)
	if err != nil {
		return nil, err
	}
	return Number{Value: n}, nil
}

// Number represents a number value, exposed as a *big.Float. Numbers can be
// floats or integers.
type Number struct {
	// Unknown will be true if the value is not yet known.
	Unknown bool

	// Null will be true if the value was not set, or was explicitly set to
	// null.
	Null bool

	// Value contains the set value, as long as Unknown and Null are both
	// false.
	Value *big.Float
}

// Type returns a NumberType.
func (n Number) Type(_ context.Context) attr.Type {
	return NumberType
}

// ToTerraformValue returns the data contained in the Number as a tftypes.Value.
func (n Number) ToTerraformValue(_ context.Context) (tftypes.Value, error) {
	if n.Null {
		return tftypes.NewValue(tftypes.Number, nil), nil
	}
	if n.Unknown {
		return tftypes.NewValue(tftypes.Number, tftypes.UnknownValue), nil
	}
	if n.Value == nil {
		return tftypes.NewValue(tftypes.Number, nil), nil
	}
	if err := tftypes.ValidateValue(tftypes.Number, n.Value); err != nil {
		return tftypes.NewValue(tftypes.Number, tftypes.UnknownValue), err
	}
	return tftypes.NewValue(tftypes.Number, n.Value), nil
}

// Equal returns true if `other` is a Number and has the same value as `n`.
func (n Number) Equal(other attr.Value) bool {
	o, ok := other.(Number)
	if !ok {
		return false
	}
	if n.Unknown != o.Unknown {
		return false
	}
	if n.Null != o.Null {
		return false
	}
	if n.Value == nil && o.Value == nil {
		return true
	}
	if n.Value == nil || o.Value == nil {
		return false
	}
	return n.Value.Cmp(o.Value) == 0
}

// IsNull returns true if the Number represents a null value.
func (n Number) IsNull() bool {
	return n.Null || (!n.Unknown && n.Value == nil)
}

// IsUnknown returns true if the Number represents a currently unknown value.
func (n Number) IsUnknown() bool {
	return n.Unknown
}

// String returns a human-readable representation of the Number value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (n Number) String() string {
	if n.Unknown {
		return attr.UnknownValueString
	}

	if n.IsNull() {
		return attr.NullValueString
	}

	return n.Value.String()
}
