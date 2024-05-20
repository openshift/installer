package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v3"
)

var diskResizeSchema = map[string]*schema.Schema{
	"disk_id": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "ID of the disk to resize.",
		ValidateDiagFunc: validateUUID,
	},
	"size": {
		Type:             schema.TypeInt,
		Required:         true,
		ForceNew:         true,
		Description:      "Disk size in bytes.",
		ValidateDiagFunc: validateDiskSize,
	},
}

func (p *provider) diskResizeResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.diskResizeCreate,
		ReadContext:   p.diskResizeRead,
		DeleteContext: p.diskResizeDelete,
		Schema:        diskResizeSchema,
		Description: `The ovirt_disk_resize resource resizes disks in oVirt to the specified size. 
		
~> Only use this resource with disks created from templates. Otherwise, two terraform resources will handle the same disk resource.  
		`,
	}
}

func (p *provider) diskResizeCreate(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	return resizeDisk(client, data)
}

func (p *provider) diskResizeRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	diskID := data.Get("disk_id").(string)
	disk, err := client.GetDisk(ovirtclient.DiskID(diskID))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to retrieve disk.",
				Detail:   err.Error(),
			},
		}
	}

	data.SetId(diskID)
	diags := diag.Diagnostics{}
	diags = setResourceField(data, "size", disk.ProvisionedSize(), diags)

	return diags
}

func (p *provider) diskResizeDelete(_ context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	data.SetId("")
	return nil
}

func resizeDisk(client ovirtclient.Client, data *schema.ResourceData) diag.Diagnostics {
	diskID := data.Get("disk_id").(string)
	newSize := data.Get("size").(int)

	params := ovirtclient.UpdateDiskParams()
	_, err := params.WithProvisionedSize(uint64(newSize))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set parameters for updating disk size.",
				Detail:   err.Error(),
			},
		}
	}

	updateFailedDiag := diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Failed to update disk size.",
	}
	diskUpdate, err := client.StartUpdateDisk(ovirtclient.DiskID(diskID), params)
	if err != nil {
		if isNotFound(err) {
			data.SetId("")
		}
		updateFailedDiag.Detail = err.Error()
		return diag.Diagnostics{updateFailedDiag}
	}
	_, err = diskUpdate.Wait()
	if err != nil {
		if isNotFound(err) {
			data.SetId("")
		}
		updateFailedDiag.Detail = err.Error()
		return diag.Diagnostics{updateFailedDiag}
	}

	data.SetId(diskID)
	diags := diag.Diagnostics{}
	diags = setResourceField(data, "size", newSize, diags)

	return diags
}
