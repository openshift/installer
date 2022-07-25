package ovirt

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v2"
)

var nicSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"vnic_profile_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "ID of the VNIC profile to associate with this NIC.",
		ForceNew:         true,
		ValidateDiagFunc: validateUUID,
	},
	"vm_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "ID of the VM to attach this NIC to.",
		ForceNew:         true,
		ValidateDiagFunc: validateUUID,
	},
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Human-readable name for the NIC.",
		ForceNew:         true,
		ValidateDiagFunc: validateNonEmpty,
	},
}

func (p *provider) nicResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.nicCreate,
		ReadContext:   p.nicRead,
		DeleteContext: p.nicDelete,
		Importer: &schema.ResourceImporter{
			StateContext: p.nicImport,
		},
		Schema:      nicSchema,
		Description: "The ovirt_nic resource creates network interfaces in oVirt.",
	}
}

func (p *provider) nicCreate(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	vmID := data.Get("vm_id").(string)
	vnicProfileID := data.Get("vnic_profile_id").(string)
	name := data.Get("name").(string)

	nic, err := client.CreateNIC(
		ovirtclient.VMID(vmID),
		ovirtclient.VNICProfileID(vnicProfileID),
		name,
		nil,
	)
	if err != nil {
		return errorToDiags("create NIC", err)
	}

	return nicResourceUpdate(nic, data)
}

func (p *provider) nicRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	id := data.Id()
	vmID := data.Get("vm_id").(string)
	nic, err := client.GetNIC(
		ovirtclient.VMID(vmID),
		ovirtclient.NICID(id),
	)
	if err != nil {
		if isNotFound(err) {
			data.SetId("")
			return nil
		}
		return errorToDiags("get NIC", err)
	}
	return nicResourceUpdate(nic, data)
}

func (p *provider) nicDelete(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	id := data.Id()
	vmID := data.Get("vm_id").(string)
	if err := client.RemoveNIC(
		ovirtclient.VMID(vmID),
		ovirtclient.NICID(id),
	); err != nil {
		if !isNotFound(err) {
			return errorToDiags("get NIC", err)
		}
	}
	data.SetId("")
	return nil
}

func (p *provider) nicImport(ctx context.Context, data *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData,
	error,
) {
	client := p.client.WithContext(ctx)
	importID := data.Id()

	parts := strings.SplitN(importID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf(
			"invalid import specification, the ID should be specified as: VMID/NICID",
		)
	}
	nic, err := client.GetNIC(
		ovirtclient.VMID(parts[0]),
		ovirtclient.NICID(parts[1]),
	)
	if err != nil {
		return nil, err
	}
	if diags := nicResourceUpdate(nic, data); diags.HasError() {
		return nil, diagsToError(diags)
	}
	return []*schema.ResourceData{data}, nil
}

func nicResourceUpdate(nic ovirtclient.NIC, data *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	data.SetId(string(nic.ID()))
	diags = setResourceField(data, "vnic_profile_id", nic.VNICProfileID(), diags)
	diags = setResourceField(data, "name", nic.Name(), diags)
	diags = setResourceField(data, "vm_id", nic.VMID(), diags)
	return diags
}
