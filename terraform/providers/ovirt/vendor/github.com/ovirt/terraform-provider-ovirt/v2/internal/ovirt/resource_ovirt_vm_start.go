package ovirt

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v2"
)

// VMStopBehavior describes the possible methods by which a VM can be stopped.
type VMStopBehavior string

const (
	// VMStopBehaviorStop causes the VM to be powered off.
	VMStopBehaviorStop VMStopBehavior = "stop"
	// VMStopBehaviorShutdown causes the VM to be shut down via an ACPI shutdown.
	VMStopBehaviorShutdown VMStopBehavior = "shutdown"
)

func vmStopBehaviorValues() []string {
	return []string{
		string(VMStopBehaviorStop),
		string(VMStopBehaviorShutdown),
	}
}

var vmStartSchema = map[string]*schema.Schema{
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
	"status": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "up",
		Description: "Desired status of the VM. The only valid value is \"up\".",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			if i == "up" {
				return nil
			}
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Invalid value.",
					Detail:        "The value must be \"up\".",
					AttributePath: path,
				},
			}
		},
		ForceNew: true,
	},
	"stop_behavior": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          VMStopBehaviorShutdown,
		Description:      "Use \"stop\" to power-off the machine, or \"shutdown\" (default) to send an ACPI shutdown.",
		ForceNew:         false,
		ValidateDiagFunc: validateEnum(vmStopBehaviorValues()),
	},
	"force_stop": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Force stop/shutdown even if a backup is in progress.",
		ForceNew:    false,
	},
}

func (p *provider) vmStartResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.vmStartCreate,
		ReadContext:   p.vmStartRead,
		UpdateContext: p.vmStartUpdate,
		DeleteContext: p.vmStartDelete,
		Importer: &schema.ResourceImporter{
			StateContext: p.vmStartImport,
		},
		Schema:      vmStartSchema,
		Description: "The ovirt_vm_start resource starts a VM in oVirt when created and stops the VM when destroyed. If additional resources should be created before the VM is started, please use the `depends_on` clause.",
	}
}

func (p *provider) vmStartCreate(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	id := data.Get("vm_id").(string)
	if err := client.StartVM(ovirtclient.VMID(id)); err != nil {
		return errorToDiags("start VM", err)
	}
	desiredStatus := data.Get("status").(string)
	vm, err := client.WaitForVMStatus(
		ovirtclient.VMID(id),
		ovirtclient.VMStatus(desiredStatus),
	)
	if err != nil {
		return errorToDiags("wait for VM start", err)
	}
	return vmStartResourceUpdate(vm, data)
}

func (p *provider) vmStartRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	id := data.Get("vm_id").(string)
	vm, err := client.GetVM(ovirtclient.VMID(id))
	if err != nil {
		if isNotFound(err) {
			data.SetId("")
			return nil
		}
		return errorToDiags("get VM status", err)
	}
	return vmStartResourceUpdate(vm, data)
}

func (p *provider) vmStartUpdate(
	_ context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	_ = data.Set("force_stop", data.Get("force_stop").(bool))
	_ = data.Set("stop_behavior", data.Get("stop_behavior").(string))
	return nil
}

func (p *provider) vmStartDelete(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	stopBehavior := data.Get("stop_behavior").(string)
	force := data.Get("force_stop").(bool)
	var err error
	if stopBehavior == string(VMStopBehaviorStop) {
		err = client.StopVM(ovirtclient.VMID(data.Id()), force)
	} else {
		err = client.ShutdownVM(ovirtclient.VMID(data.Id()), force)
	}
	if err != nil {
		return errorToDiags("shutdown VM status", err)
	}
	_, err = client.WaitForVMStatus(
		ovirtclient.VMID(data.Id()),
		ovirtclient.VMStatusDown,
	)
	if err != nil {
		return errorToDiags("wait for VM to stop", err)
	}
	data.SetId("")
	_ = data.Set("status", "")
	return nil
}

func (p *provider) vmStartImport(ctx context.Context, data *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData,
	error,
) {
	client := p.client.WithContext(ctx)
	vm, err := client.GetVM(ovirtclient.VMID(data.Id()))
	if err != nil {
		return nil, fmt.Errorf("failed to import VM %s (%w)", data.Id(), err)
	}
	d := vmStartResourceUpdate(vm, data)
	if err := diagsToError(d); err != nil {
		return nil, fmt.Errorf("failed to import VM %s (%w)", data.Id(), err)
	}
	return []*schema.ResourceData{
		data,
	}, nil
}

func vmStartResourceUpdate(vm ovirtclient.VM, data *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	data.SetId(string(vm.ID()))
	diags = setResourceField(data, "vm_id", vm.ID(), diags)
	diags = setResourceField(data, "status", vm.Status(), diags)
	return diags
}
