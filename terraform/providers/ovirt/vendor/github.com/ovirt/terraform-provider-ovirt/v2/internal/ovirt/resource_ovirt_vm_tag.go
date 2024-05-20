package ovirt

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v3"
)

var vmTagSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"tag_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "ID for the tag to be attached",
		ForceNew:         true,
		ValidateDiagFunc: validateUUID,
	},
	"vm_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "ID for the VM to be attached",
		ForceNew:         true,
		ValidateDiagFunc: validateUUID,
	},
}

func (p *provider) vmTagResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.vmTagCreate,
		ReadContext:   p.vmTagRead,
		DeleteContext: p.vmTagDelete,
		Schema:        vmTagSchema,
		Description:   "The ovirt_vm_tag resource attaches a tag to a virtual machine.",
	}
}

func (p *provider) vmTagCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	tagID := data.Get("tag_id").(string)
	vmID := data.Get("vm_id").(string)
	err := client.AddTagToVM(ovirtclient.VMID(vmID), ovirtclient.TagID(tagID))
	if err != nil {
		return errorToDiags(fmt.Sprintf("add tag %s to VM %s", tagID, vmID), err)
	}
	data.SetId(fmt.Sprintf("%s_%s", tagID, vmID))
	return nil
}

func (p *provider) vmTagRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	id := strings.SplitN(data.Id(), "_", 2)
	if len(id) < 2 {
		return errorToDiags(
			fmt.Sprintf("failed to read tag attachmend %s", data.Id()),
			fmt.Errorf("incorrect ID format"),
		)
	}
	tagID := id[0]
	vmID := id[1]
	tags, err := client.ListVMTags(ovirtclient.VMID(vmID))
	if err != nil {
		if ovirtclient.HasErrorCode(err, ovirtclient.ENotFound) {
			data.SetId("")
			return nil
		}
		return errorToDiags(fmt.Sprintf("listing VM %s tags", vmID), err)
	}
	found := false
	for _, tag := range tags {
		if tag.ID() == ovirtclient.TagID(tagID) {
			found = true
		}
	}
	if !found {
		data.SetId("")
	}
	return nil
}

func (p *provider) vmTagDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	vmID := data.Get("vm_id").(string)
	tagID := data.Get("tag_id").(string)
	if err := client.RemoveTagFromVM(
		ovirtclient.VMID(vmID),
		ovirtclient.TagID(tagID),
	); err != nil {
		return errorToDiags(fmt.Sprintf("remove tag %s from vm %s", tagID, vmID), err)
	}
	return nil
}
