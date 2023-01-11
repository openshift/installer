package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// State returns the *tfprotov5.DynamicValue for a *tfsdk.State.
func State(ctx context.Context, fw *tfsdk.State) (*tfprotov5.DynamicValue, diag.Diagnostics) {
	if fw == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	proto5, err := tfprotov5.NewDynamicValue(fw.Schema.Type().TerraformType(ctx), fw.Raw)

	if err != nil {
		diags.AddError(
			"Unable to Convert State",
			"An unexpected error was encountered when converting the state to the protocol type. "+
				"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n"+
				"Please report this to the provider developer:\n\n"+err.Error(),
		)

		return nil, diags
	}

	return &proto5, nil
}
