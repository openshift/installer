package tf5dynamicvalue

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Must creates a *tfprotov5.DynamicValue or panics. This is intended only for
// simplifying testing code.
//
// The tftypes.Type parameter is separate to enable DynamicPsuedoType testing.
func Must(typ tftypes.Type, value tftypes.Value) *tfprotov5.DynamicValue {
	dynamicValue, err := tfprotov5.NewDynamicValue(typ, value)

	if err != nil {
		panic(fmt.Sprintf("unable to create DynamicValue: %s", err.Error()))
	}

	return &dynamicValue
}
