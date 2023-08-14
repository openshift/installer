package types

import "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

type Bool = basetypes.BoolValue

// BoolNull creates a Bool with a null value. Determine whether the value is
// null via the Bool type IsNull method.
func BoolNull() basetypes.BoolValue {
	return basetypes.NewBoolNull()
}

// BoolUnknown creates a Bool with an unknown value. Determine whether the
// value is unknown via the Bool type IsUnknown method.
func BoolUnknown() basetypes.BoolValue {
	return basetypes.NewBoolUnknown()
}

// BoolValue creates a Bool with a known value. Access the value via the Bool
// type ValueBool method.
func BoolValue(value bool) basetypes.BoolValue {
	return basetypes.NewBoolValue(value)
}
