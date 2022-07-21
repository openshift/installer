package ovirt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v2"
)

var vmDisksResizeSchema = map[string]*schema.Schema{
	"vm_id": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Resize all disks in this VM to the specified size.",
		ValidateDiagFunc: validateUUID,
	},
	"size": {
		Type:             schema.TypeInt,
		Required:         true,
		ForceNew:         true,
		Description:      "Disk size in bytes to set all disks to.",
		ValidateDiagFunc: validateDiskSize,
	},
}

func (p *provider) vmDisksResizeResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.vmDisksResizeCreate,
		ReadContext:   p.vmDisksResizeRead,
		DeleteContext: p.vmDisksResizeDelete,
		Schema:        vmDisksResizeSchema,
		Description: `The ovirt_vm_disks_resize resource resizes all disks in an oVirt VM to the specified size. 
		
~> Only use this resource with disks created from templates. Otherwise, two terraform resources will handle the same disk resource.  
		`,
	}
}

func (p *provider) vmDisksResizeCreate(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	return resizeAllDisks(client, data)
}

func (p *provider) vmDisksResizeRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	desiredSize := uint64(data.Get("size").(int))
	size := desiredSize

	vmID := ovirtclient.VMID(data.Get("vm_id").(string))

	diskAttachments, err := client.ListDiskAttachments(vmID)
	if err != nil {
		return errorToDiags(fmt.Sprintf("list disk attachments of VM %s", vmID), err)
	}
	for _, diskAttachment := range diskAttachments {
		disk, err := diskAttachment.Disk()
		if err != nil {
			return errorToDiags(fmt.Sprintf("get disk %s", diskAttachment.DiskID()), err)
		}
		if disk.ProvisionedSize() != desiredSize {
			// Set the reported size to the size differing so the resource can be refreshed and the disks resized.
			size = disk.ProvisionedSize()
		}
	}

	data.SetId(string(vmID))
	diags := diag.Diagnostics{}
	diags = setResourceField(data, "size", size, diags)

	return diags
}

func (p *provider) vmDisksResizeDelete(_ context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	data.SetId("")
	return nil
}

func resizeAllDisks(client ovirtclient.Client, data *schema.ResourceData) diag.Diagnostics {
	vmID := ovirtclient.VMID(data.Get("vm_id").(string))
	desiredSize := uint64(data.Get("size").(int))

	diskAttachments, err := client.ListDiskAttachments(vmID)
	if err != nil {
		return errorToDiags(fmt.Sprintf("list disk attachments of VM %s", vmID), err)
	}
	var diags diag.Diagnostics
	for _, diskAttachment := range diskAttachments {
		disk, err := diskAttachment.Disk()
		if err != nil {
			return errorToDiags(fmt.Sprintf("get disk %s", diskAttachment.DiskID()), err)
		}
		if disk.ProvisionedSize() == desiredSize {
			continue
		}
		params := ovirtclient.UpdateDiskParams()
		if _, err := params.WithProvisionedSize(desiredSize); err != nil {
			diags = append(
				diags,
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Failed to set parameters for updating disk %s size.", disk.ID()),
					Detail:   err.Error(),
				},
			)
			continue
		}
		updateFailedDiag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to update disk %s size.", disk.ID()),
		}
		diskUpdate, err := client.StartUpdateDisk(disk.ID(), params)
		if err != nil {
			updateFailedDiag.Detail = err.Error()
			diags = append(diags, updateFailedDiag)
			continue
		}
		_, err = diskUpdate.Wait()
		if err != nil {
			updateFailedDiag.Detail = err.Error()
			diags = append(diags, updateFailedDiag)
			continue
		}
	}

	data.SetId(string(vmID))
	if !diags.HasError() {
		diags = setResourceField(data, "size", desiredSize, diags)
	}
	return diags
}
