package tf5dynamicvalue

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// IsNull returns true if the given *tfprotov5.DynamicValue is nil or
// represents a null value.
func IsNull(schema *tfprotov5.Schema, dynamicValue *tfprotov5.DynamicValue) (bool, error) {
	if dynamicValue == nil {
		return true, nil
	}

	// Panic prevention
	if schema == nil {
		return false, fmt.Errorf("unable to unmarshal DynamicValue: missing Type")
	}

	tfValue, err := dynamicValue.Unmarshal(schema.ValueType())

	if err != nil {
		return false, fmt.Errorf("unable to unmarshal DynamicValue: %w", err)
	}

	return tfValue.IsNull(), nil
}
