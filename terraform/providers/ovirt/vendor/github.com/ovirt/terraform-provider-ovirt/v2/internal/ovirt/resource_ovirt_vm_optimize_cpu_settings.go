package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v2"
)

var vmOptimizeCPUSettingsSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "oVirt ID of the VM to be started.",
	},
	"vm_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "oVirt ID of the VM to be optimized.",
		ForceNew:         true,
		ValidateDiagFunc: validateUUID,
	},
}

func (p *provider) vmOptimizeCPUSettingsResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.vmOptimizeCPUSettingsCreate,
		ReadContext:   p.vmOptimizeCPUSettingsRead,
		DeleteContext: p.vmOptimizeCPUSettingsDelete,
		Schema:        vmOptimizeCPUSettingsSchema,
		Description:   "The ovirt_vm_optimize_cpu_settings sets the CPU settings to automatically optimized for the specified VM.",
	}
}

func (p *provider) vmOptimizeCPUSettingsCreate(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	vmID := data.Get("vm_id").(string)
	err := p.client.WithContext(ctx).AutoOptimizeVMCPUPinningSettings(ovirtclient.VMID(vmID), true)
	if err == nil {
		data.SetId(vmID)
	}
	return errorToDiags("auto-optimizing CPU pinning settings", err)
}

func (p *provider) vmOptimizeCPUSettingsRead(
	_ context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	vmID := data.Get("vm_id").(string)
	data.SetId(vmID)
	return nil
}

func (p *provider) vmOptimizeCPUSettingsDelete(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	vmID := data.Get("vm_id").(string)
	err := p.client.WithContext(ctx).AutoOptimizeVMCPUPinningSettings(ovirtclient.VMID(vmID), false)
	if err == nil {
		data.SetId("")
	}
	return errorToDiags("auto-optimizing CPU pinning settings", err)
}
