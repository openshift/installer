package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Config returns the *tfprotov6.DynamicValue for a *tfsdk.Config.
func Config(ctx context.Context, fw *tfsdk.Config) (*tfprotov6.DynamicValue, diag.Diagnostics) {
	if fw == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	proto6, err := tfprotov6.NewDynamicValue(fw.Schema.Type().TerraformType(ctx), fw.Raw)

	if err != nil {
		diags.AddError(
			"Unable to Convert Configuration",
			"An unexpected error was encountered when converting the configuration to the protocol type. "+
				"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n"+
				"Please report this to the provider developer:\n\n"+err.Error(),
		)

		return nil, diags
	}

	return &proto6, nil
}
