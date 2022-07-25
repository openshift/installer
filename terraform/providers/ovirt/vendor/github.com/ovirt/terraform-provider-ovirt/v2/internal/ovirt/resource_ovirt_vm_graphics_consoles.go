package ovirt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v2"
)

var vmGraphicsConsoleSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "UUID of the graphics console.",
	},
}

var vmGraphicsConsolesSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "oVirt ID of the VM to be started.",
	},
	"vm_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "oVirt ID of the VM to be started.",
		ForceNew:         true,
		ValidateDiagFunc: validateUUID,
	},
	"console": {
		Type:        schema.TypeSet,
		Optional:    true,
		ForceNew:    true,
		Description: "The list of consoles that should be on this VM. If a console is not in this list it will be removed from the VM.",
		MinItems:    0,
		Elem: &schema.Resource{
			Schema: vmGraphicsConsoleSchema,
		},
	},
}

func (p *provider) vmGraphicsConsolesResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.vmGraphicsConsolesCreate,
		ReadContext:   p.vmGraphicsConsolesRead,
		DeleteContext: p.vmGraphicsConsolesDelete,
		Schema:        vmGraphicsConsolesSchema,
		Description:   "The ovirt_vm_graphics_consoles controls all the graphic consoles of a VM.",
	}
}

func (p *provider) vmGraphicsConsolesCreate(
	ctx context.Context,
	data *schema.ResourceData,
	i interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	vmID := data.Get("vm_id").(string)
	consoles := data.Get("console").(*schema.Set)

	if consoles.Len() != 0 {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Creating consoles is not supported",
				Detail:   "Currently, only removing all graphics consoles is supported.",
			},
		}
	}

	cons, err := client.ListVMGraphicsConsoles(ovirtclient.VMID(vmID))
	if err != nil {
		return errorToDiags("listing graphics consoles", err)
	}
	var diags diag.Diagnostics
	for _, console := range cons {
		if err := console.Remove(); err != nil && !ovirtclient.HasErrorCode(err, ovirtclient.ENotFound) {
			diags = append(diags, errorToDiag(fmt.Sprintf("remove graphics console %s", console.ID()), err))
		}
	}
	data.SetId(vmID)
	return p.vmGraphicsConsolesToData(client, vmID, data, diags)
}

func (p *provider) vmGraphicsConsolesRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	vmID := data.Get("vm_id").(string)
	data.SetId(vmID)
	return p.vmGraphicsConsolesToData(client, vmID, data, nil)
}

func (p *provider) vmGraphicsConsolesToData(
	client ovirtclient.Client,
	vmID string,
	data *schema.ResourceData,
	diags diag.Diagnostics,
) diag.Diagnostics {
	cons, err := client.ListVMGraphicsConsoles(ovirtclient.VMID(vmID))
	if err != nil {
		return errorToDiags(fmt.Sprintf("list graphics consoles for VM %s", vmID), err)
	}
	result := make([]map[string]interface{}, len(cons))
	for i, con := range cons {
		result[i] = map[string]interface{}{
			"id": con.ID(),
		}
	}
	if err := data.Set("console", result); err != nil {
		diags = append(diags, errorToDiag("setting console on result", err))
	}
	return diags
}

func (p *provider) vmGraphicsConsolesDelete(
	_ context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	// Since we don't support adding graphics consoles there is nothing to do here.
	data.SetId("")
	return nil
}
